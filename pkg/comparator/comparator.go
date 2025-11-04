package comparator

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"deepComparator/pkg/concurrent"
	"deepComparator/pkg/database"
	"deepComparator/pkg/models"
	"deepComparator/pkg/progress"
)

// convertBytesToString converts byte arrays to strings, keeping other types as-is
func convertBytesToString(value interface{}) interface{} {
	if byteArray, ok := value.([]byte); ok {
		return string(byteArray)
	}
	return value
}

// Comparator handles the comparison logic between two databases
type Comparator struct {
	DB1              *database.Connection
	DB2              *database.Connection
	ConcurrentWorker *concurrent.ConcurrentComparator
	MaxWorkers       int
	UUIDDecoder      *models.UUIDDecoder
}

// NewComparator creates a new comparator instance
func NewComparator(db1, db2 *database.Connection) *Comparator {
	maxWorkers := 4 // Default number of workers
	return &Comparator{
		DB1:              db1,
		DB2:              db2,
		ConcurrentWorker: concurrent.NewConcurrentComparator(db1, db2, maxWorkers),
		MaxWorkers:       maxWorkers,
		UUIDDecoder:      models.NewUUIDDecoder(true), // Enabled by default
	}
}

// NewConcurrentComparator creates a new comparator with specified worker count
func NewConcurrentComparator(db1, db2 *database.Connection, maxWorkers int) *Comparator {
	return &Comparator{
		DB1:              db1,
		DB2:              db2,
		ConcurrentWorker: concurrent.NewConcurrentComparator(db1, db2, maxWorkers),
		MaxWorkers:       maxWorkers,
		UUIDDecoder:      models.NewUUIDDecoder(true), // Enabled by default
	}
}

// NewComparatorWithUUIDecoding creates a new comparator with UUID decoding option
func NewComparatorWithUUIDDecoding(db1, db2 *database.Connection, maxWorkers int, decodeUUIDs bool) *Comparator {
	return &Comparator{
		DB1:              db1,
		DB2:              db2,
		ConcurrentWorker: concurrent.NewConcurrentComparator(db1, db2, maxWorkers),
		MaxWorkers:       maxWorkers,
		UUIDDecoder:      models.NewUUIDDecoder(decodeUUIDs),
	}
}

// CompareTable compares a table between two databases
func (c *Comparator) CompareTable(schema, tableName string, criteria *models.MatchCriteria) (*models.ComparisonResult, error) {
	// Check if table exists in both databases
	exists1, err := c.DB1.TableExists(schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence in DB1: %w", err)
	}
	if !exists1 {
		return nil, fmt.Errorf("table %s.%s does not exist in database 1", schema, tableName)
	}

	exists2, err := c.DB2.TableExists(schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence in DB2: %w", err)
	}
	if !exists2 {
		return nil, fmt.Errorf("table %s.%s does not exist in database 2", schema, tableName)
	}

	// Get table schemas
	schema1, err := c.DB1.GetTableSchema(schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema from DB1: %w", err)
	}

	_, err = c.DB2.GetTableSchema(schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema from DB2: %w", err)
	}

	// Get table data using concurrent operations
	data1, data2, _, err := c.ConcurrentWorker.ParallelDataFetch(schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get data using parallel fetch: %w", err)
	}

	// Perform comparison
	result := &models.ComparisonResult{
		TableName:         tableName,
		Schema:            schema,
		Timestamp:         time.Now(),
		TotalRowsDB1:      len(data1.Rows),
		TotalRowsDB2:      len(data2.Rows),
		OnlyInDB1:         []models.TableRow{},
		OnlyInDB2:         []models.TableRow{},
		Differences:       []models.RowDifference{},
		ForeignKeyResults: []models.ForeignKeyResult{},
	}

	// Create match criteria if not provided
	if criteria == nil {
		criteria = c.createDefaultMatchCriteria(schema1)
	}

	// Match rows between databases
	matchProgress := progress.NewSimpleProgress("Matching rows")
	matches, onlyInDB1, onlyInDB2 := c.matchRows(data1.Rows, data2.Rows, criteria)
	matchProgress.Finish(fmt.Sprintf("Found %d matches", len(matches)))

	result.OnlyInDB1 = onlyInDB1
	result.OnlyInDB2 = onlyInDB2
	result.MatchedRows = len(matches)
	result.UnmatchedRows = len(onlyInDB1) + len(onlyInDB2)

	// Compare matched rows for differences
	var comparisonProgress *progress.ProgressBar
	if len(matches) > 0 { // Show progress bar for any number of matches
		comparisonProgress = progress.NewProgressBar(int64(len(matches)), "Comparing matched rows")
	}

	for i, match := range matches {
		diff := c.compareRowsWithFK(match.row1, match.row2, criteria, schema1.ForeignKeys)
		if len(diff.ColumnDifferences) > 0 {
			diff.RowIdentifier = c.getRowIdentifier(match.row1, criteria)
			result.Differences = append(result.Differences, *diff)
		}

		if comparisonProgress != nil {
			// Update progress for every row, or batch for large datasets
			if len(matches) <= 100 || (i+1)%10 == 0 {
				comparisonProgress.SetProgress(int64(i + 1))
			}
		}
	}

	if comparisonProgress != nil {
		comparisonProgress.FinishWithMessage(fmt.Sprintf("Found %d differences", len(result.Differences)))
	}

	// Compare foreign key relationships
	for _, fk := range schema1.ForeignKeys {
		fkResult := c.compareForeignKey(fk, data1, data2, criteria)
		result.ForeignKeyResults = append(result.ForeignKeyResults, *fkResult)
	}

	// Process UUIDs if enabled
	result = c.UUIDDecoder.ProcessComparisonResult(result)

	return result, nil
}

// rowMatch represents a matched pair of rows
type rowMatch struct {
	row1 models.TableRow
	row2 models.TableRow
}

// matchRows matches rows between two datasets based on criteria
func (c *Comparator) matchRows(rows1, rows2 []models.TableRow, criteria *models.MatchCriteria) ([]rowMatch, []models.TableRow, []models.TableRow) {
	var matches []rowMatch
	var onlyInDB1 []models.TableRow
	var onlyInDB2 []models.TableRow

	// Create a map for faster lookup of rows2
	rows2Map := make(map[string]models.TableRow)
	for _, row2 := range rows2 {
		key := c.getRowKey(row2, criteria)
		rows2Map[key] = row2
	}

	// Track which rows from DB2 have been matched
	matchedInDB2 := make(map[string]bool)

	// Find matches and rows only in DB1
	for _, row1 := range rows1 {
		key := c.getRowKey(row1, criteria)
		if row2, exists := rows2Map[key]; exists {
			matches = append(matches, rowMatch{row1: row1, row2: row2})
			matchedInDB2[key] = true
		} else {
			onlyInDB1 = append(onlyInDB1, row1)
		}
	}

	// Find rows only in DB2
	for _, row2 := range rows2 {
		key := c.getRowKey(row2, criteria)
		if !matchedInDB2[key] {
			onlyInDB2 = append(onlyInDB2, row2)
		}
	}

	return matches, onlyInDB1, onlyInDB2
}

// getRowKey generates a key for matching rows based on criteria
func (c *Comparator) getRowKey(row models.TableRow, criteria *models.MatchCriteria) string {
	var keyParts []string

	// Build exclude map with explicit exclusions
	excludeMap := make(map[string]bool)
	for _, col := range criteria.ExcludeColumns {
		excludeMap[col] = true
	}

	// Add columns from file to exclude map if enabled
	fileColumns, err := c.getExcludeColumnsFromFile(criteria)
	if err != nil {
		// Log error but continue without file column exclusion
		fmt.Printf("Warning: Could not load exclude columns from file: %v\n", err)
	} else {
		for _, col := range fileColumns {
			excludeMap[col] = true
		}
	}

	// If specific columns are defined, use only those
	if len(criteria.Columns) > 0 {
		for _, col := range criteria.Columns {
			if !excludeMap[col] {
				if val, exists := row[col]; exists {
					keyParts = append(keyParts, fmt.Sprintf("%s:%v", col, val))
				}
			}
		}
	} else {
		// Use all columns except excluded ones and primary keys (unless specified)
		for col, val := range row {
			if !excludeMap[col] {
				// Skip primary key columns unless explicitly included
				if c.isPrimaryKeyColumn(col) && !criteria.IncludePrimaryKey {
					continue
				}
				keyParts = append(keyParts, fmt.Sprintf("%s:%v", col, val))
			}
		}
	}

	// Sort key parts to ensure consistent ordering
	sort.Strings(keyParts)
	return strings.Join(keyParts, "|")
}

// getRowIdentifier creates a human-readable identifier for a row
func (c *Comparator) getRowIdentifier(row models.TableRow, criteria *models.MatchCriteria) string {
	return c.getRowKey(row, criteria)
}

// compareRows compares two rows and returns differences
func (c *Comparator) compareRows(row1, row2 models.TableRow, criteria *models.MatchCriteria) *models.RowDifference {
	return c.compareRowsWithFK(row1, row2, criteria, nil)
}

func (c *Comparator) compareRowsWithFK(row1, row2 models.TableRow, criteria *models.MatchCriteria, foreignKeys []models.ForeignKey) *models.RowDifference {
	diff := &models.RowDifference{
		DB1Row:            row1,
		DB2Row:            row2,
		ColumnDifferences: []models.ColumnDifference{},
	}

	excludeMap := make(map[string]bool)
	for _, col := range criteria.ExcludeColumns {
		excludeMap[col] = true
	}

	// Add columns from file to exclude map if enabled
	fileColumns, err := c.getExcludeColumnsFromFile(criteria)
	if err != nil {
		// Log error but continue without file column exclusion
		fmt.Printf("Warning: Could not load exclude columns from file: %v\n", err)
	} else {
		for _, col := range fileColumns {
			excludeMap[col] = true
		}
	}

	// Create map for quick FK lookup
	fkMap := make(map[string]models.ForeignKey)
	for _, fk := range foreignKeys {
		fkMap[fk.ColumnName] = fk
	}

	// Compare all columns
	allColumns := make(map[string]bool)
	for col := range row1 {
		allColumns[col] = true
	}
	for col := range row2 {
		allColumns[col] = true
	}

	for col := range allColumns {
		if excludeMap[col] {
			continue
		}

		val1, exists1 := row1[col]
		val2, exists2 := row2[col]

		if !exists1 && !exists2 {
			continue
		}

		if !exists1 || !exists2 || !reflect.DeepEqual(val1, val2) {
			// Convert byte arrays to strings for consistent processing
			db1Value := convertBytesToString(val1)
			db2Value := convertBytesToString(val2)

			colDiff := models.ColumnDifference{
				ColumnName: col,
				DB1Value:   db1Value,
				DB2Value:   db2Value,
			}

			// Check if this column is a foreign key
			if fk, isForeignKey := fkMap[col]; isForeignKey {
				colDiff.IsForeignKey = true

				// Get referenced data if values exist
				if val1 != nil || val2 != nil {
					fkRef := c.getForeignKeyReference(fk, val1, val2)
					colDiff.ForeignKeyReference = fkRef
				}
			}

			diff.ColumnDifferences = append(diff.ColumnDifferences, colDiff)
		}
	}

	return diff
}

// getForeignKeyReference gets the referenced data for a foreign key
func (c *Comparator) getForeignKeyReference(fk models.ForeignKey, val1, val2 interface{}) *models.ForeignKeyReference {
	var values []interface{}
	if val1 != nil {
		values = append(values, val1)
	}
	if val2 != nil {
		values = append(values, val2)
	}

	if len(values) == 0 {
		return nil
	}

	// Get foreign key data from both databases
	fkData1, err1 := c.DB1.GetForeignKeyData(fk, values)
	fkData2, err2 := c.DB2.GetForeignKeyData(fk, values)

	if err1 != nil && err2 != nil {
		return nil
	}

	fkRef := &models.ForeignKeyReference{
		ForeignKey: fk,
	}

	// Find the referenced rows
	if val1 != nil && err1 == nil {
		for _, row := range fkData1 {
			if pkVal, exists := row[fk.ReferencedColumnName]; exists && reflect.DeepEqual(pkVal, val1) {
				fkRef.DB1Referenced = row
				break
			}
		}
	}

	if val2 != nil && err2 == nil {
		for _, row := range fkData2 {
			if pkVal, exists := row[fk.ReferencedColumnName]; exists && reflect.DeepEqual(pkVal, val2) {
				fkRef.DB2Referenced = row
				break
			}
		}
	}

	// Check if there are differences in the referenced data
	if len(fkRef.DB1Referenced) > 0 && len(fkRef.DB2Referenced) > 0 {
		fkRef.ReferencedDiff = !reflect.DeepEqual(fkRef.DB1Referenced, fkRef.DB2Referenced)
	} else {
		fkRef.ReferencedDiff = len(fkRef.DB1Referenced) != len(fkRef.DB2Referenced)
	}

	return fkRef
}

// compareForeignKey compares foreign key relationships
func (c *Comparator) compareForeignKey(fk models.ForeignKey, data1, data2 *models.TableData, criteria *models.MatchCriteria) *models.ForeignKeyResult {
	result := &models.ForeignKeyResult{
		ForeignKey:   fk,
		FKReferences: []models.ForeignKeyReference{},
	}

	// Get unique foreign key values from both datasets
	var fkValues1, fkValues2 []interface{}

	for _, row := range data1.Rows {
		if val, exists := row[fk.ColumnName]; exists && val != nil {
			fkValues1 = append(fkValues1, val)
		}
	}

	for _, row := range data2.Rows {
		if val, exists := row[fk.ColumnName]; exists && val != nil {
			fkValues2 = append(fkValues2, val)
		}
	}

	// Get all unique values for comparison
	allValues := c.getUniqueValues(append(fkValues1, fkValues2...))

	// Get foreign key data from both databases
	fkData1, err1 := c.DB1.GetForeignKeyData(fk, allValues)
	fkData2, err2 := c.DB2.GetForeignKeyData(fk, allValues)

	if err1 != nil || err2 != nil {
		result.Error = fmt.Sprintf("Error getting foreign key data: DB1=%v, DB2=%v", err1, err2)
		return result
	}

	// Create temporary table data objects
	tempData1 := &models.TableData{
		TableName: fk.ReferencedTable,
		Schema:    fk.ReferencedSchema,
		Rows:      fkData1,
	}

	tempData2 := &models.TableData{
		TableName: fk.ReferencedTable,
		Schema:    fk.ReferencedSchema,
		Rows:      fkData2,
	}

	// Get schema for the referenced table to create appropriate criteria
	referencedSchema, err := c.DB1.GetTableSchema(fk.ReferencedSchema, fk.ReferencedTable)
	var fkCriteria *models.MatchCriteria
	if err != nil {
		// If we can't get schema, use the provided criteria
		fkCriteria = criteria
	} else {
		// Create default criteria for the referenced table
		fkCriteria = c.createDefaultMatchCriteria(referencedSchema)
		// Copy file-based exclusion settings from original criteria
		fkCriteria.ExcludeColumnsFromFile = criteria.ExcludeColumnsFromFile
		fkCriteria.ExcludeColumnsFile = criteria.ExcludeColumnsFile
	}

	// Compare the foreign key data directly using row matching with appropriate criteria
	matches, onlyInDB1, onlyInDB2 := c.matchRows(tempData1.Rows, tempData2.Rows, fkCriteria)

	fkComparison := &models.ComparisonResult{
		TableName:     fk.ReferencedTable,
		Schema:        fk.ReferencedSchema,
		Timestamp:     time.Now(),
		TotalRowsDB1:  len(tempData1.Rows),
		TotalRowsDB2:  len(tempData2.Rows),
		MatchedRows:   len(matches),
		UnmatchedRows: len(onlyInDB1) + len(onlyInDB2),
		OnlyInDB1:     onlyInDB1,
		OnlyInDB2:     onlyInDB2,
		Differences:   []models.RowDifference{},
	}

	// Check for differences in matched foreign key rows and build FK references
	for _, match := range matches {
		diff := c.compareRows(match.row1, match.row2, criteria)
		hasDiff := len(diff.ColumnDifferences) > 0

		if hasDiff {
			diff.RowIdentifier = c.getRowIdentifier(match.row1, criteria)
			fkComparison.Differences = append(fkComparison.Differences, *diff)
		}

		// Create FK reference for this match
		fkRef := models.ForeignKeyReference{
			ForeignKey:     fk,
			DB1Referenced:  match.row1,
			DB2Referenced:  match.row2,
			ReferencedDiff: hasDiff,
		}
		result.FKReferences = append(result.FKReferences, fkRef)
	}

	// Add references for rows only in DB1
	for _, row := range onlyInDB1 {
		fkRef := models.ForeignKeyReference{
			ForeignKey:     fk,
			DB1Referenced:  row,
			ReferencedDiff: true, // Always different if only in one DB
		}
		result.FKReferences = append(result.FKReferences, fkRef)
	}

	// Add references for rows only in DB2
	for _, row := range onlyInDB2 {
		fkRef := models.ForeignKeyReference{
			ForeignKey:     fk,
			DB2Referenced:  row,
			ReferencedDiff: true, // Always different if only in one DB
		}
		result.FKReferences = append(result.FKReferences, fkRef)
	}

	result.ComparisonResult = *fkComparison
	return result
}

// createDefaultMatchCriteria creates default matching criteria based on table schema
func (c *Comparator) createDefaultMatchCriteria(schema *models.TableSchema) *models.MatchCriteria {
	var excludeColumns []string

	// Exclude primary key columns by default (IDs that might differ)
	for _, col := range schema.Columns {
		if col.IsPrimary {
			excludeColumns = append(excludeColumns, col.ColumnName)
		}
	}

	return &models.MatchCriteria{
		Columns:                []string{}, // Empty means use all columns
		ExcludeColumns:         excludeColumns,
		IncludePrimaryKey:      false,
		ExcludeColumnsFromFile: true,                  // Enable file column exclusion by default
		ExcludeColumnsFile:     "exclude_columns.txt", // Default exclude columns file
	}
}

// isPrimaryKeyColumn checks if a column is a primary key (simple heuristic)
func (c *Comparator) isPrimaryKeyColumn(columnName string) bool {
	// Simple heuristic - can be improved with schema information
	return columnName == "id" || columnName == "ID"
}

// getUniqueValues returns unique values from a slice
func (c *Comparator) getUniqueValues(values []interface{}) []interface{} {
	var unique []interface{}

	for _, val := range values {
		if val == nil {
			continue
		}

		// Check if this value already exists in unique slice
		found := false
		for _, existing := range unique {
			if reflect.DeepEqual(val, existing) {
				found = true
				break
			}
		}

		if !found {
			unique = append(unique, val)
		}
	}

	return unique
}

// FindReferences finds all references to a specific table/column across both databases
func (c *Comparator) FindReferences(schema, tableName, columnName string) (*models.MatchReferenceResult, error) {
	// Use parallel reference analysis for better performance
	concurrentResult, err := c.ConcurrentWorker.ParallelReferenceAnalysis(schema, tableName, columnName)
	if err != nil {
		return nil, fmt.Errorf("failed to perform parallel reference analysis: %w", err)
	}

	// Process UUIDs if enabled
	concurrentResult = c.UUIDDecoder.ProcessMatchReferenceResult(concurrentResult)

	// Use the concurrent result directly as it's already complete
	return concurrentResult, nil
}

// getExcludeColumnsFromFile loads exclude columns from file with error handling
func (c *Comparator) getExcludeColumnsFromFile(criteria *models.MatchCriteria) ([]string, error) {
	if !criteria.ExcludeColumnsFromFile {
		return []string{}, nil
	}

	return models.LoadExcludeColumnsFromFile(criteria.ExcludeColumnsFile)
}

// AnalyzeFKReferences finds all tables that reference a specific ID as foreign key
func (c *Comparator) AnalyzeFKReferences(schema, tableName, targetID string) (*models.FKAnalysisResult, error) {
	// Create connection progress
	connProgress := progress.NewSimpleProgress("Connecting to databases for FK analysis")

	// Ensure we have active connections
	if c.DB1 == nil || c.DB2 == nil {
		return nil, fmt.Errorf("database connections not initialized")
	}

	// Query to find all foreign key constraints that reference the target table
	// This query looks for FK constraints where:
	// - tc = table with the FK constraint (source table)
	// - kcu = column info for the FK constraint
	// - ccu = referenced table info (target table we're looking for)
	fkQuery := `
		SELECT 
			tc.constraint_name,
			tc.table_schema,
			tc.table_name,
			kcu.column_name,
			ccu.table_schema AS foreign_table_schema,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name 
		FROM information_schema.table_constraints AS tc 
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
			AND ccu.table_schema = tc.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY' 
			AND ccu.table_schema = $1 
			AND ccu.table_name = $2
		ORDER BY tc.table_schema, tc.table_name, kcu.column_name`

	connProgress.Finish("Connected successfully")

	// Create loading progress for FK discovery
	loadingProgress := progress.NewSimpleProgress("Discovering foreign key constraints")

	rows1, err := c.DB1.DB.Query(fkQuery, schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query foreign keys from DB1: %w", err)
	}
	defer rows1.Close()

	var fkConstraints []models.FKTableReference

	// Process formal foreign key constraints
	for rows1.Next() {
		var constraintName, tableSchema, table, columnName string
		var foreignTableSchema, foreignTableName, foreignColumnName string

		err := rows1.Scan(&constraintName, &tableSchema, &table, &columnName,
			&foreignTableSchema, &foreignTableName, &foreignColumnName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan foreign key constraint: %w", err)
		}

		fkConstraints = append(fkConstraints, models.FKTableReference{
			Schema:         tableSchema,
			TableName:      table,
			ColumnName:     columnName,
			ConstraintName: constraintName,
		})
	}

	// If no formal FK constraints found, look for potential FK relationships based on column naming patterns
	// This handles cases where FK relationships exist at the data level but formal constraints are not defined
	if len(fkConstraints) == 0 {
		// Search for columns that might reference this table based on naming conventions
		potentialFKQuery := `
			SELECT DISTINCT 
				'' as constraint_name,
				c.table_schema,
				c.table_name,
				c.column_name
			FROM information_schema.columns c
			WHERE c.table_schema = $1 
				AND c.column_name = $2
				AND c.table_name != $3`

		// Try common patterns: table_name + '_id'
		columnPatterns := []string{
			tableName + "_id",
			tableName + "Id",
			tableName + "_ID",
		}

		for _, pattern := range columnPatterns {
			rows, err := c.DB1.DB.Query(potentialFKQuery, schema, pattern, tableName)
			if err != nil {
				continue
			}

			for rows.Next() {
				var constraintName, tableSchema, table, columnName string
				err := rows.Scan(&constraintName, &tableSchema, &table, &columnName)
				if err != nil {
					continue
				}

				fkConstraints = append(fkConstraints, models.FKTableReference{
					Schema:         tableSchema,
					TableName:      table,
					ColumnName:     columnName,
					ConstraintName: fmt.Sprintf("potential_fk_%s_%s_%s", tableSchema, table, columnName),
				})
			}
			rows.Close()
		}
	}

	loadingProgress.Finish(fmt.Sprintf("Found %d FK constraints (formal + potential)", len(fkConstraints)))

	if len(fkConstraints) == 0 {
		return &models.FKAnalysisResult{
			TargetTable:       tableName,
			TargetSchema:      schema,
			TargetID:          targetID,
			Timestamp:         time.Now(),
			TotalConstraints:  0,
			ReferencingTables: []models.FKTableReference{},
		}, nil
	}

	// Create progress bar for analyzing references
	analysisProgress := progress.NewProgressBar(int64(len(fkConstraints)), "Analyzing FK references")

	// Now find references for each foreign key constraint
	for i := range fkConstraints {
		fk := &fkConstraints[i]

		// Find matching records in both databases
		query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s = $1 LIMIT 5",
			fk.ColumnName, fk.Schema, fk.TableName, fk.ColumnName)

		// Query DB1
		rows1, err := c.DB1.DB.Query(query, targetID)
		if err != nil {
			// If query fails, skip this FK
			analysisProgress.Update(1)
			continue
		}

		var db1Values []string
		for rows1.Next() {
			var value string
			if err := rows1.Scan(&value); err != nil {
				continue
			}
			db1Values = append(db1Values, value)
		}
		rows1.Close()

		// Query DB2
		rows2, err := c.DB2.DB.Query(query, targetID)
		if err != nil {
			// If query fails, skip this FK
			analysisProgress.Update(1)
			continue
		}

		var db2Values []string
		for rows2.Next() {
			var value string
			if err := rows2.Scan(&value); err != nil {
				continue
			}
			db2Values = append(db2Values, value)
		}
		rows2.Close()

		// Update FK constraint with found values
		fk.MatchesDB1 = len(db1Values)
		fk.MatchesDB2 = len(db2Values)
		fk.SampleRows = db1Values

		analysisProgress.Update(1)
	}

	analysisProgress.Finish()

	// Filter out FKs with no references
	var referencingTables []models.FKTableReference
	for _, fk := range fkConstraints {
		if fk.MatchesDB1 > 0 || fk.MatchesDB2 > 0 {
			referencingTables = append(referencingTables, fk)
		}
	}

	result := &models.FKAnalysisResult{
		TargetTable:       tableName,
		TargetSchema:      schema,
		TargetID:          targetID,
		Timestamp:         time.Now(),
		TotalConstraints:  len(fkConstraints),
		ReferencingTables: referencingTables,
	}

	return result, nil
}

// discoverFKConstraints finds all foreign key constraints that reference a specific table
func (c *Comparator) discoverFKConstraints(schema, tableName string) ([]models.FKTableReference, error) {
	query := `
		SELECT 
			tc.constraint_name,
			tc.table_schema,
			tc.table_name,
			kcu.column_name
		FROM information_schema.table_constraints AS tc 
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
			AND ccu.table_schema = tc.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY' 
			AND ccu.table_schema = $1 
			AND ccu.table_name = $2
		ORDER BY tc.table_schema, tc.table_name, kcu.column_name`

	// Use DB1 for schema information (both DBs should have same structure)
	rows, err := c.DB1.DB.Query(query, schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query foreign keys: %w", err)
	}
	defer rows.Close()

	var fkConstraints []models.FKTableReference

	// Process foreign key constraints
	for rows.Next() {
		var constraintName, tableSchema, table, columnName string

		err := rows.Scan(&constraintName, &tableSchema, &table, &columnName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan foreign key constraint: %w", err)
		}

		fkConstraints = append(fkConstraints, models.FKTableReference{
			Schema:         tableSchema,
			TableName:      table,
			ColumnName:     columnName,
			ConstraintName: constraintName,
		})
	}

	return fkConstraints, nil
}

// GenerateUpdateScript generates SQL script to update foreign key references from idTarget to idDestination
// and then delete the original record
func (c *Comparator) GenerateUpdateScript(schema, tableName, idTarget, idDestination string) (string, error) {
	// Discover FK constraints pointing to this table
	fkConstraints, err := c.discoverFKConstraints(schema, tableName)
	if err != nil {
		return "", fmt.Errorf("failed to discover FK constraints: %w", err)
	}

	var script strings.Builder

	// Header
	script.WriteString("-- Generated FK Update Script\n")
	script.WriteString(fmt.Sprintf("-- Target table: %s.%s\n", schema, tableName))
	script.WriteString(fmt.Sprintf("-- Update FK references from ID %s to ID %s\n", idTarget, idDestination))
	script.WriteString(fmt.Sprintf("-- Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	script.WriteString("-- WARNING: Review this script before execution!\n")
	script.WriteString("\n")

	// Begin transaction
	script.WriteString("BEGIN;\n\n")

	// Update foreign key references
	script.WriteString("-- Update foreign key references\n")
	for i, fk := range fkConstraints {
		if i > 0 {
			script.WriteString("\n")
		}
		script.WriteString(fmt.Sprintf("-- Table: %s.%s, Column: %s\n", fk.Schema, fk.TableName, fk.ColumnName))

		updateSQL := fmt.Sprintf("UPDATE %s.%s SET %s = %s WHERE %s = %s;",
			fk.Schema, fk.TableName, fk.ColumnName, idDestination, fk.ColumnName, idTarget)
		script.WriteString(updateSQL + "\n")
	}

	script.WriteString("\n")

	// Delete original record
	script.WriteString("-- Delete original record\n")
	deleteSQL := fmt.Sprintf("DELETE FROM %s.%s WHERE id = %s;", schema, tableName, idTarget)
	script.WriteString(deleteSQL + "\n")

	script.WriteString("\n")

	// Commit transaction
	script.WriteString("COMMIT;\n")

	script.WriteString("\n")
	script.WriteString("-- Script execution completed\n")
	script.WriteString("-- Verify results and check referential integrity\n")

	return script.String(), nil
}

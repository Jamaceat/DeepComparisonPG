package concurrent

import (
	"context"
	"deepComparator/pkg/database"
	"deepComparator/pkg/models"
	"deepComparator/pkg/progress"
	"fmt"
	"sync"
	"time"
)

// ConcurrentComparator handles concurrent database operations
type ConcurrentComparator struct {
	DB1        *database.Connection
	DB2        *database.Connection
	WorkerPool *WorkerPool
	maxWorkers int
}

// NewConcurrentComparator creates a new concurrent comparator
func NewConcurrentComparator(db1, db2 *database.Connection, maxWorkers int) *ConcurrentComparator {
	if maxWorkers <= 0 {
		maxWorkers = 4 // Default to 4 workers
	}

	return &ConcurrentComparator{
		DB1:        db1,
		DB2:        db2,
		maxWorkers: maxWorkers,
	}
}

// ParallelDataFetch fetches table data and schema concurrently
func (cc *ConcurrentComparator) ParallelDataFetch(schema, tableName string) (*models.TableData, *models.TableData, *models.TableSchema, error) {
	type fetchResult struct {
		data1  *models.TableData
		data2  *models.TableData
		schema *models.TableSchema
		err    error
	}

	resultChan := make(chan fetchResult, 1)

	// Show loading progress for data fetch
	loadProgress := progress.NewSimpleProgress("Loading table data")

	go func() {
		defer close(resultChan)

		// Use WaitGroup to synchronize goroutines
		var wg sync.WaitGroup
		var mu sync.Mutex

		var data1, data2 *models.TableData
		var tableSchema *models.TableSchema
		var errs []error // Fetch data from DB1
		wg.Add(1)
		go func() {
			defer wg.Done()
			d1, err := cc.DB1.GetTableData(schema, tableName)
			mu.Lock()
			data1 = d1
			if err != nil {
				errs = append(errs, fmt.Errorf("DB1 data fetch error: %w", err))
			}
			mu.Unlock()
		}()

		// Fetch data from DB2
		wg.Add(1)
		go func() {
			defer wg.Done()
			d2, err := cc.DB2.GetTableData(schema, tableName)
			mu.Lock()
			data2 = d2
			if err != nil {
				errs = append(errs, fmt.Errorf("DB2 data fetch error: %w", err))
			}
			mu.Unlock()
		}()

		// Fetch schema from DB1
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, err := cc.DB1.GetTableSchema(schema, tableName)
			mu.Lock()
			tableSchema = s
			if err != nil {
				errs = append(errs, fmt.Errorf("Schema fetch error: %w", err))
			}
			mu.Unlock()
		}()

		wg.Wait()

		// Complete the loading progress
		totalRows := 0
		if data1 != nil {
			totalRows += len(data1.Rows)
		}
		if data2 != nil {
			totalRows += len(data2.Rows)
		}
		loadProgress.Finish(fmt.Sprintf("Loaded %d total rows", totalRows))

		result := fetchResult{
			data1:  data1,
			data2:  data2,
			schema: tableSchema,
		}

		if len(errs) > 0 {
			result.err = fmt.Errorf("fetch errors: %v", errs)
		}

		resultChan <- result
	}()

	// Wait for result with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	select {
	case result := <-resultChan:
		return result.data1, result.data2, result.schema, result.err
	case <-ctx.Done():
		return nil, nil, nil, fmt.Errorf("data fetch timeout: %w", ctx.Err())
	}
}

// ParallelForeignKeyAnalysis analyzes foreign keys concurrently
func (cc *ConcurrentComparator) ParallelForeignKeyAnalysis(foreignKeys []models.ForeignKey, data1, data2 *models.TableData, criteria *models.MatchCriteria) ([]models.ForeignKeyResult, error) {
	if len(foreignKeys) == 0 {
		return []models.ForeignKeyResult{}, nil
	}

	// Create worker pool for FK analysis
	workerCount := cc.maxWorkers
	if len(foreignKeys) < workerCount {
		workerCount = len(foreignKeys)
	}

	wp := NewWorkerPool(workerCount, len(foreignKeys))
	wp.Start()
	defer wp.Stop()

	// Submit FK analysis jobs
	for i, fk := range foreignKeys {
		job := Job{
			ID:       fmt.Sprintf("fk_%d", i),
			TaskType: "foreign_key_analysis",
			Data: map[string]interface{}{
				"foreign_key": fk,
				"data1":       data1,
				"data2":       data2,
				"criteria":    criteria,
				"db1":         cc.DB1,
				"db2":         cc.DB2,
			},
			Timeout: 30 * time.Second,
		}
		wp.SubmitJob(job)
	}

	// Collect results
	results := make([]models.ForeignKeyResult, len(foreignKeys))
	resultMap := make(map[string]int)

	for i, fk := range foreignKeys {
		resultMap[fmt.Sprintf("fk_%d", i)] = i
		results[i] = models.ForeignKeyResult{
			ForeignKey: fk,
			Error:      "Processing...",
		}
	}

	// Process results as they come in
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	completed := 0
	for completed < len(foreignKeys) {
		select {
		case result := <-wp.GetResults():
			if idx, exists := resultMap[result.JobID]; exists {
				if result.Error != nil {
					results[idx].Error = result.Error.Error()
				} else if fkResult, ok := result.Data.(models.ForeignKeyResult); ok {
					results[idx] = fkResult
				}
				completed++
			}
		case <-ctx.Done():
			return results, fmt.Errorf("FK analysis timeout: %w", ctx.Err())
		}
	}

	return results, nil
}

// ParallelReferenceAnalysis analyzes references concurrently
func (cc *ConcurrentComparator) ParallelReferenceAnalysis(targetSchema, targetTable, targetColumn string) (*models.MatchReferenceResult, error) {
	result := &models.MatchReferenceResult{
		TargetTable:  targetTable,
		TargetSchema: targetSchema,
		TargetColumn: targetColumn,
		Timestamp:    time.Now(),
		References:   []models.ReferenceMatch{},
	}

	// Concurrently get referencing tables from both databases
	var referencingTables1, referencingTables2 []models.ForeignKey
	var err1, err2 error

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		referencingTables1, err1 = cc.DB1.GetReferencingTables(targetSchema, targetTable, targetColumn)
	}()

	go func() {
		defer wg.Done()
		referencingTables2, err2 = cc.DB2.GetReferencingTables(targetSchema, targetTable, targetColumn)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, fmt.Errorf("failed to get referencing tables from DB1: %w", err1)
	}
	if err2 != nil {
		return nil, fmt.Errorf("failed to get referencing tables from DB2: %w", err2)
	}

	// Combine all referencing tables
	allReferencingTables := make(map[string]models.ForeignKey)
	for _, fk := range referencingTables1 {
		key := fmt.Sprintf("%s.%s.%s", fk.ReferencedSchema, fk.ReferencedTable, fk.ColumnName)
		allReferencingTables[key] = fk
	}
	for _, fk := range referencingTables2 {
		key := fmt.Sprintf("%s.%s.%s", fk.ReferencedSchema, fk.ReferencedTable, fk.ColumnName)
		if _, exists := allReferencingTables[key]; !exists {
			allReferencingTables[key] = fk
		}
	}

	if len(allReferencingTables) == 0 {
		return result, nil
	}

	// Process each referencing table concurrently
	referenceChan := make(chan models.ReferenceMatch, len(allReferencingTables))
	errorChan := make(chan error, len(allReferencingTables))

	semaphore := make(chan struct{}, cc.maxWorkers)

	for _, fk := range allReferencingTables {
		go func(foreignKey models.ForeignKey) {
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			refMatch := models.ReferenceMatch{
				TableName:      foreignKey.ReferencedTable,
				Schema:         foreignKey.ReferencedSchema,
				ColumnName:     foreignKey.ColumnName,
				ConstraintName: foreignKey.ConstraintName,
			}

			// Concurrently get values from both databases
			var values1, values2 []interface{}
			var wg sync.WaitGroup

			wg.Add(2)
			go func() {
				defer wg.Done()
				v1, err := cc.DB1.GetColumnValues(foreignKey.ReferencedSchema, foreignKey.ReferencedTable, foreignKey.ColumnName)
				if err == nil {
					values1 = v1
				}
			}()

			go func() {
				defer wg.Done()
				v2, err := cc.DB2.GetColumnValues(foreignKey.ReferencedSchema, foreignKey.ReferencedTable, foreignKey.ColumnName)
				if err == nil {
					values2 = v2
				}
			}()

			wg.Wait()

			refMatch.DB1References = values1
			refMatch.DB2References = values2

			// Categorize values
			refMatch.CommonRefs, refMatch.OnlyInDB1, refMatch.OnlyInDB2 = cc.categorizeValues(values1, values2)

			referenceChan <- refMatch
		}(fk)
	}

	// Collect results
	references := make([]models.ReferenceMatch, 0, len(allReferencingTables))
	completed := 0

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for completed < len(allReferencingTables) {
		select {
		case ref := <-referenceChan:
			references = append(references, ref)
			completed++
		case err := <-errorChan:
			// Log error but continue
			fmt.Printf("Warning: error processing reference: %v\n", err)
			completed++
		case <-ctx.Done():
			return result, fmt.Errorf("reference analysis timeout: %w", ctx.Err())
		}
	}

	result.References = references
	result.ReferencingTables = len(references)

	totalRefs := 0
	for _, ref := range references {
		totalRefs += len(ref.DB1References) + len(ref.DB2References)
	}
	result.TotalReferences = totalRefs

	return result, nil
}

// categorizeValues separates values into common, only in first, only in second
func (cc *ConcurrentComparator) categorizeValues(values1, values2 []interface{}) (common, onlyInFirst, onlyInSecond []interface{}) {
	// Create maps for O(1) lookup
	map1 := make(map[string]interface{})
	map2 := make(map[string]interface{})

	// Populate maps concurrently for large datasets
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for _, val := range values1 {
			key := fmt.Sprintf("%v", val)
			map1[key] = val
		}
	}()

	go func() {
		defer wg.Done()
		for _, val := range values2 {
			key := fmt.Sprintf("%v", val)
			map2[key] = val
		}
	}()

	wg.Wait()

	// Find common values and values only in first
	for key, val := range map1 {
		if _, exists := map2[key]; exists {
			common = append(common, val)
		} else {
			onlyInFirst = append(onlyInFirst, val)
		}
	}

	// Find values only in second
	for key, val := range map2 {
		if _, exists := map1[key]; !exists {
			onlyInSecond = append(onlyInSecond, val)
		}
	}

	return common, onlyInFirst, onlyInSecond
}

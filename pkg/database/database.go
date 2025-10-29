package database

import (
	"database/sql"
	"fmt"
	"strings"

	"deepComparator/pkg/models"

	_ "github.com/lib/pq"
)

// Connection represents a database connection
type Connection struct {
	DB     *sql.DB
	Config models.DatabaseConfig
}

// NewConnection creates a new database connection
func NewConnection(config models.DatabaseConfig) (*Connection, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{
		DB:     db,
		Config: config,
	}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.DB.Close()
}

// GetTableSchema retrieves table schema information
func (c *Connection) GetTableSchema(schema, tableName string) (*models.TableSchema, error) {
	tableSchema := &models.TableSchema{
		TableName:   tableName,
		Schema:      schema,
		Columns:     []models.ColumnInfo{},
		ForeignKeys: []models.ForeignKey{},
	}

	// Get column information
	columnQuery := `
		SELECT 
			c.column_name, 
			c.data_type, 
			c.is_nullable = 'YES' as is_nullable,
			CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END as is_primary
		FROM information_schema.columns c
		LEFT JOIN (
			SELECT ku.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage ku ON tc.constraint_name = ku.constraint_name
			WHERE tc.constraint_type = 'PRIMARY KEY' 
			AND tc.table_schema = $1 
			AND tc.table_name = $2
		) pk ON c.column_name = pk.column_name
		WHERE c.table_schema = $1 AND c.table_name = $2
		ORDER BY c.ordinal_position`

	rows, err := c.DB.Query(columnQuery, schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get column information: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var col models.ColumnInfo
		if err := rows.Scan(&col.ColumnName, &col.DataType, &col.IsNullable, &col.IsPrimary); err != nil {
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}
		tableSchema.Columns = append(tableSchema.Columns, col)
	}

	// Get foreign key information
	fkQuery := `
		SELECT 
			kcu.column_name,
			ccu.table_name AS referenced_table_name,
			ccu.table_schema AS referenced_schema_name,
			ccu.column_name AS referenced_column_name,
			tc.constraint_name
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
			AND ccu.table_schema = tc.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY'
			AND tc.table_schema = $1 
			AND tc.table_name = $2`

	fkRows, err := c.DB.Query(fkQuery, schema, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get foreign key information: %w", err)
	}
	defer fkRows.Close()

	for fkRows.Next() {
		var fk models.ForeignKey
		if err := fkRows.Scan(&fk.ColumnName, &fk.ReferencedTable, &fk.ReferencedSchema, &fk.ReferencedColumnName, &fk.ConstraintName); err != nil {
			return nil, fmt.Errorf("failed to scan foreign key info: %w", err)
		}
		tableSchema.ForeignKeys = append(tableSchema.ForeignKeys, fk)
	}

	return tableSchema, nil
}

// GetTableData retrieves all data from a table
func (c *Connection) GetTableData(schema, tableName string) (*models.TableData, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s", schema, tableName)

	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table data: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	tableData := &models.TableData{
		TableName: tableName,
		Schema:    schema,
		Rows:      []models.TableRow{},
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(models.TableRow)
		for i, col := range columns {
			row[col] = values[i]
		}
		tableData.Rows = append(tableData.Rows, row)
	}

	return tableData, nil
}

// GetForeignKeyData retrieves data from a foreign key referenced table based on specific values
func (c *Connection) GetForeignKeyData(fk models.ForeignKey, values []interface{}) ([]models.TableRow, error) {
	if len(values) == 0 {
		return []models.TableRow{}, nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(
		"SELECT * FROM %s.%s WHERE %s IN (%s)",
		fk.ReferencedSchema,
		fk.ReferencedTable,
		fk.ReferencedColumnName,
		strings.Join(placeholders, ","),
	)

	rows, err := c.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to query foreign key data: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var result []models.TableRow
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(models.TableRow)
		for i, col := range columns {
			row[col] = values[i]
		}
		result = append(result, row)
	}

	return result, nil
}

// TableExists checks if a table exists in the database
func (c *Connection) TableExists(schema, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = $1 AND table_name = $2
		)`

	var exists bool
	err := c.DB.QueryRow(query, schema, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}

	return exists, nil
}

// GetReferencingTables finds all tables that have foreign keys pointing to the specified table/column
func (c *Connection) GetReferencingTables(targetSchema, targetTable, targetColumn string) ([]models.ForeignKey, error) {
	query := `
		SELECT 
			kcu.column_name,
			kcu.table_name,
			kcu.table_schema,
			tc.constraint_name
		FROM information_schema.key_column_usage kcu
		JOIN information_schema.table_constraints tc 
			ON kcu.constraint_name = tc.constraint_name 
			AND kcu.table_schema = tc.table_schema
		JOIN information_schema.referential_constraints rc 
			ON tc.constraint_name = rc.constraint_name
			AND tc.table_schema = rc.constraint_schema
		JOIN information_schema.key_column_usage rcu 
			ON rc.unique_constraint_name = rcu.constraint_name
			AND rc.unique_constraint_schema = rcu.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY'
			AND rcu.table_schema = $1 
			AND rcu.table_name = $2
			AND rcu.column_name = $3
		ORDER BY kcu.table_schema, kcu.table_name, kcu.column_name`

	rows, err := c.DB.Query(query, targetSchema, targetTable, targetColumn)
	if err != nil {
		return nil, fmt.Errorf("failed to query referencing tables: %w", err)
	}
	defer rows.Close()

	var foreignKeys []models.ForeignKey
	for rows.Next() {
		var fk models.ForeignKey
		err := rows.Scan(
			&fk.ColumnName,
			&fk.ReferencedTable,  // This is actually the referencing table
			&fk.ReferencedSchema, // This is actually the referencing schema
			&fk.ConstraintName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan foreign key: %w", err)
		}

		// Set the actual referenced table/schema/column (our target)
		actualReferencingTable := fk.ReferencedTable
		actualReferencingSchema := fk.ReferencedSchema

		fk.ReferencedTable = targetTable
		fk.ReferencedSchema = targetSchema
		fk.ReferencedColumnName = targetColumn

		// Create a temporary ForeignKey with correct structure
		referencingFK := models.ForeignKey{
			ColumnName:           fk.ColumnName,
			ReferencedTable:      actualReferencingTable,
			ReferencedSchema:     actualReferencingSchema,
			ReferencedColumnName: targetColumn,
			ConstraintName:       fk.ConstraintName,
		}

		foreignKeys = append(foreignKeys, referencingFK)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating foreign key rows: %w", err)
	}

	return foreignKeys, nil
}

// GetColumnValues gets all distinct values from a specific column in a table
func (c *Connection) GetColumnValues(schema, tableName, columnName string) ([]interface{}, error) {
	query := fmt.Sprintf(`SELECT DISTINCT "%s" FROM "%s"."%s" WHERE "%s" IS NOT NULL ORDER BY "%s"`,
		columnName, schema, tableName, columnName, columnName)

	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query column values: %w", err)
	}
	defer rows.Close()

	var values []interface{}
	for rows.Next() {
		var value interface{}
		err := rows.Scan(&value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan column value: %w", err)
		}
		values = append(values, value)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating column values: %w", err)
	}

	return values, nil
}

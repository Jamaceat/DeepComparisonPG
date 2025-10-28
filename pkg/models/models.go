package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// DatabaseConfig represents the database connection configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
}

// TableRow represents a single row from a table
type TableRow map[string]interface{}

// TableData represents all data from a table
type TableData struct {
	TableName string     `json:"table_name"`
	Schema    string     `json:"schema"`
	Rows      []TableRow `json:"rows"`
}

// ForeignKey represents a foreign key relationship
type ForeignKey struct {
	ColumnName           string `json:"column_name"`
	ReferencedTable      string `json:"referenced_table"`
	ReferencedSchema     string `json:"referenced_schema"`
	ReferencedColumnName string `json:"referenced_column_name"`
	ConstraintName       string `json:"constraint_name"`
}

// ColumnInfo represents column metadata
type ColumnInfo struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
	IsNullable bool   `json:"is_nullable"`
	IsPrimary  bool   `json:"is_primary"`
}

// TableSchema represents table structure and metadata
type TableSchema struct {
	TableName   string       `json:"table_name"`
	Schema      string       `json:"schema"`
	Columns     []ColumnInfo `json:"columns"`
	ForeignKeys []ForeignKey `json:"foreign_keys"`
}

// ComparisonResult represents the result of comparing two tables
type ComparisonResult struct {
	TableName         string             `json:"table_name"`
	Schema            string             `json:"schema"`
	Timestamp         time.Time          `json:"timestamp"`
	TotalRowsDB1      int                `json:"total_rows_db1"`
	TotalRowsDB2      int                `json:"total_rows_db2"`
	MatchedRows       int                `json:"matched_rows"`
	UnmatchedRows     int                `json:"unmatched_rows"`
	OnlyInDB1         []TableRow         `json:"only_in_db1"`
	OnlyInDB2         []TableRow         `json:"only_in_db2"`
	Differences       []RowDifference    `json:"differences"`
	ForeignKeyResults []ForeignKeyResult `json:"foreign_key_results"`
}

// RowDifference represents differences found between matching rows
type RowDifference struct {
	RowIdentifier     string             `json:"row_identifier"`
	DB1Row            TableRow           `json:"db1_row"`
	DB2Row            TableRow           `json:"db2_row"`
	ColumnDifferences []ColumnDifference `json:"column_differences"`
}

// ColumnDifference represents a difference in a specific column
type ColumnDifference struct {
	ColumnName string      `json:"column_name"`
	DB1Value   interface{} `json:"db1_value"`
	DB2Value   interface{} `json:"db2_value"`
}

// ForeignKeyResult represents the result of foreign key comparison
type ForeignKeyResult struct {
	ForeignKey       ForeignKey        `json:"foreign_key"`
	ComparisonResult *ComparisonResult `json:"comparison_result,omitempty"`
	Error            string            `json:"error,omitempty"`
}

// MatchCriteria represents the criteria used to match rows between tables
type MatchCriteria struct {
	Columns                []string `json:"columns"`
	ExcludeColumns         []string `json:"exclude_columns"`
	IncludePrimaryKey      bool     `json:"include_primary_key"`
	ExcludeColumnsFromFile bool     `json:"exclude_columns_from_file"`
	ExcludeColumnsFile     string   `json:"exclude_columns_file"`
}

// LoadExcludeColumnsFromFile loads column names to exclude from a file
func LoadExcludeColumnsFromFile(filename string) ([]string, error) {
	if filename == "" {
		return []string{}, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open exclude columns file %s: %w", filename, err)
	}
	defer file.Close()

	var columns []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		columns = append(columns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading exclude columns file: %w", err)
	}

	return columns, nil
}

// GetAllExcludeColumns returns exclude columns from file if configured
func (mc *MatchCriteria) GetAllExcludeColumns() []string {
	if !mc.ExcludeColumnsFromFile {
		return []string{}
	}

	columns, err := LoadExcludeColumnsFromFile(mc.ExcludeColumnsFile)
	if err != nil {
		// If there's an error loading the file, return empty slice
		// The error should be handled at a higher level
		return []string{}
	}

	return columns
}

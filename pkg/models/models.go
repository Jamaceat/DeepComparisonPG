package models

import (
	"bufio"
	"encoding/base64"
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
	ColumnName          string               `json:"column_name"`
	DB1Value            interface{}          `json:"db1_value"`
	DB2Value            interface{}          `json:"db2_value"`
	IsForeignKey        bool                 `json:"is_foreign_key,omitempty"`
	ForeignKeyReference *ForeignKeyReference `json:"foreign_key_reference,omitempty"`
}

// ForeignKeyReference represents the actual data referenced by a foreign key
type ForeignKeyReference struct {
	ForeignKey     ForeignKey `json:"foreign_key"`
	DB1Referenced  TableRow   `json:"db1_referenced,omitempty"`
	DB2Referenced  TableRow   `json:"db2_referenced,omitempty"`
	ReferencedDiff bool       `json:"referenced_diff"`
}

// ForeignKeyResult represents the result of a foreign key comparison
type ForeignKeyResult struct {
	ForeignKey       ForeignKey            `json:"foreign_key"`
	ComparisonResult ComparisonResult      `json:"comparison_result"`
	Error            string                `json:"error,omitempty"`
	FKReferences     []ForeignKeyReference `json:"fk_references,omitempty"`
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

// ReferenceMatch represents references found in different tables
type ReferenceMatch struct {
	TableName      string        `json:"table_name"`
	Schema         string        `json:"schema"`
	ColumnName     string        `json:"column_name"`
	ConstraintName string        `json:"constraint_name,omitempty"`
	DB1References  []interface{} `json:"db1_references"`
	DB2References  []interface{} `json:"db2_references"`
	CommonRefs     []interface{} `json:"common_references"`
	OnlyInDB1      []interface{} `json:"only_in_db1"`
	OnlyInDB2      []interface{} `json:"only_in_db2"`
}

// MatchReferenceResult represents the complete result of reference matching
type MatchReferenceResult struct {
	TargetTable       string           `json:"target_table"`
	TargetSchema      string           `json:"target_schema"`
	TargetColumn      string           `json:"target_column"`
	Timestamp         time.Time        `json:"timestamp"`
	TotalReferences   int              `json:"total_references"`
	ReferencingTables int              `json:"referencing_tables"`
	References        []ReferenceMatch `json:"references"`
}

// UUIDDecoder provides functionality to decode Base64 encoded UUIDs
type UUIDDecoder struct {
	DecodeEnabled bool
}

// NewUUIDDecoder creates a new UUID decoder
func NewUUIDDecoder(enabled bool) *UUIDDecoder {
	return &UUIDDecoder{
		DecodeEnabled: enabled,
	}
}

// IsBase64UUID checks if a string looks like a Base64 encoded UUID
func (u *UUIDDecoder) IsBase64UUID(value string) bool {
	if !u.DecodeEnabled || len(value) == 0 {
		return false
	}

	// Base64 encoded UUIDs are typically 48 characters long (36 chars UUID -> 48 chars base64)
	// Allow some flexibility but be more specific
	if len(value) < 44 || len(value) > 52 {
		return false
	}

	// Check if it contains only valid Base64 characters
	for _, char := range value {
		if !((char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '+' || char == '/' || char == '=') {
			return false
		}
	}

	// Try to decode and see if it looks like a UUID pattern
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return false
	}

	decodedStr := string(decoded)

	// Check if decoded string looks like a UUID (36 characters with hyphens at right positions)
	if len(decodedStr) == 36 &&
		decodedStr[8] == '-' &&
		decodedStr[13] == '-' &&
		decodedStr[18] == '-' &&
		decodedStr[23] == '-' {
		return true
	}

	return false
} // DecodeBase64UUID attempts to decode a Base64 string to UUID format
func (u *UUIDDecoder) DecodeBase64UUID(encoded string) string {
	if !u.DecodeEnabled || !u.IsBase64UUID(encoded) {
		return encoded
	} // Try to decode from Base64
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return encoded // Return original if decode fails
	}

	decodedStr := string(decoded)

	// Check if decoded string looks like a UUID (36 characters with hyphens)
	if len(decodedStr) == 36 &&
		decodedStr[8] == '-' &&
		decodedStr[13] == '-' &&
		decodedStr[18] == '-' &&
		decodedStr[23] == '-' {
		return decodedStr
	}

	return encoded // Return original if doesn't look like UUID
}

// ProcessTableRow decodes UUIDs in a table row if enabled
func (u *UUIDDecoder) ProcessTableRow(row TableRow) TableRow {
	if !u.DecodeEnabled {
		return row
	}

	processedRow := make(TableRow)
	for key, value := range row {
		if strValue, ok := value.(string); ok {
			processedRow[key] = u.DecodeBase64UUID(strValue)
		} else {
			processedRow[key] = value
		}
	}

	return processedRow
}

// ProcessComparisonResult processes a comparison result to decode UUIDs
func (u *UUIDDecoder) ProcessComparisonResult(result *ComparisonResult) *ComparisonResult {
	if !u.DecodeEnabled {
		return result
	}

	// Process differences
	for i, diff := range result.Differences {
		// Process DB1 row
		processedDB1Row := make(TableRow)
		for key, value := range diff.DB1Row {
			if strValue, ok := value.(string); ok {
				decodedValue := u.DecodeBase64UUID(strValue)

				processedDB1Row[key] = decodedValue
			} else {
				processedDB1Row[key] = value
			}
		}
		result.Differences[i].DB1Row = processedDB1Row

		// Process DB2 row
		processedDB2Row := make(TableRow)
		for key, value := range diff.DB2Row {
			if strValue, ok := value.(string); ok {
				processedDB2Row[key] = u.DecodeBase64UUID(strValue)
			} else {
				processedDB2Row[key] = value
			}
		}
		result.Differences[i].DB2Row = processedDB2Row

		// Process column differences
		for j, colDiff := range diff.ColumnDifferences {

			if db1StrValue, ok := colDiff.DB1Value.(string); ok {
				result.Differences[i].ColumnDifferences[j].DB1Value = u.DecodeBase64UUID(db1StrValue)
			}
			if db2StrValue, ok := colDiff.DB2Value.(string); ok {
				result.Differences[i].ColumnDifferences[j].DB2Value = u.DecodeBase64UUID(db2StrValue)
			}
		}
	}

	// Process only in DB1
	for i, row := range result.OnlyInDB1 {
		processedRow := make(TableRow)
		for key, value := range row {
			if strValue, ok := value.(string); ok {
				processedRow[key] = u.DecodeBase64UUID(strValue)
			} else {
				processedRow[key] = value
			}
		}
		result.OnlyInDB1[i] = processedRow
	}

	// Process only in DB2
	for i, row := range result.OnlyInDB2 {
		processedRow := make(TableRow)
		for key, value := range row {
			if strValue, ok := value.(string); ok {
				processedRow[key] = u.DecodeBase64UUID(strValue)
			} else {
				processedRow[key] = value
			}
		}
		result.OnlyInDB2[i] = processedRow
	}

	return result
}

// ProcessMatchReferenceResult processes a match reference result to decode UUIDs
func (u *UUIDDecoder) ProcessMatchReferenceResult(result *MatchReferenceResult) *MatchReferenceResult {
	if !u.DecodeEnabled {
		return result
	}

	// Process each reference
	for i := range result.References {
		ref := &result.References[i]

		// Process DB1 references
		for j, val := range ref.DB1References {
			// Convert bytes to string if needed
			var strValue string
			if byteArray, ok := val.([]byte); ok {
				strValue = string(byteArray)
			} else if str, ok := val.(string); ok {
				strValue = str
			} else {
				continue // Skip non-string/non-byte values
			}

			ref.DB1References[j] = u.DecodeBase64UUID(strValue)
		}

		// Process DB2 references
		for j, val := range ref.DB2References {
			// Convert bytes to string if needed
			var strValue string
			if byteArray, ok := val.([]byte); ok {
				strValue = string(byteArray)
			} else if str, ok := val.(string); ok {
				strValue = str
			} else {
				continue // Skip non-string/non-byte values
			}

			ref.DB2References[j] = u.DecodeBase64UUID(strValue)
		}

		// Process common references
		for j, val := range ref.CommonRefs {
			// Convert bytes to string if needed
			var strValue string
			if byteArray, ok := val.([]byte); ok {
				strValue = string(byteArray)
			} else if str, ok := val.(string); ok {
				strValue = str
			} else {
				continue // Skip non-string/non-byte values
			}

			ref.CommonRefs[j] = u.DecodeBase64UUID(strValue)
		}

		// Process only in DB1
		for j, val := range ref.OnlyInDB1 {
			// Convert bytes to string if needed
			var strValue string
			if byteArray, ok := val.([]byte); ok {
				strValue = string(byteArray)
			} else if str, ok := val.(string); ok {
				strValue = str
			} else {
				continue // Skip non-string/non-byte values
			}

			ref.OnlyInDB1[j] = u.DecodeBase64UUID(strValue)
		}

		// Process only in DB2
		for j, val := range ref.OnlyInDB2 {
			// Convert bytes to string if needed
			var strValue string
			if byteArray, ok := val.([]byte); ok {
				strValue = string(byteArray)
			} else if str, ok := val.(string); ok {
				strValue = str
			} else {
				continue // Skip non-string/non-byte values
			}

			ref.OnlyInDB2[j] = u.DecodeBase64UUID(strValue)
		}
	}

	return result
}

// FKTableReference represents a table that references a specific ID
type FKTableReference struct {
	Schema         string   `json:"schema"`
	TableName      string   `json:"table_name"`
	ColumnName     string   `json:"column_name"`
	ConstraintName string   `json:"constraint_name"`
	MatchesDB1     int      `json:"matches_db1"`
	MatchesDB2     int      `json:"matches_db2"`
	SampleRows     []string `json:"sample_rows,omitempty"`
}

// FKAnalysisResult represents the result of analyzing FK references for a specific ID
type FKAnalysisResult struct {
	TargetTable       string             `json:"target_table"`
	TargetSchema      string             `json:"target_schema"`
	TargetID          string             `json:"target_id"`
	Timestamp         time.Time          `json:"timestamp"`
	TotalConstraints  int                `json:"total_constraints"`
	ReferencingTables []FKTableReference `json:"referencing_tables"`
}

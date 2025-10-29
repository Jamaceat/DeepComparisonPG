package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"deepComparator/pkg/comparator"
	"deepComparator/pkg/config"
	"deepComparator/pkg/database"
	"deepComparator/pkg/models"
	"deepComparator/pkg/progress"
)

func main() {
	// Command line flags
	var (
		envFile         = flag.String("env", ".env", "Path to environment file")
		tableName       = flag.String("table", "", "Table name to compare (required)")
		schemaName      = flag.String("schema", "public", "Schema name (default: public)")
		outputFile      = flag.String("output", "", "Output file path (overrides env config)")
		excludeCols     = flag.String("exclude", "", "Comma-separated list of columns to exclude from comparison")
		includeCols     = flag.String("include", "", "Comma-separated list of columns to include in comparison (if empty, all columns are used)")
		includePK       = flag.Bool("include-pk", false, "Include primary key columns in comparison")
		excludeFromFile = flag.Bool("exclude-from-file", true, "Exclude columns from file")
		excludeFile     = flag.String("exclude-file", "exclude_columns.txt", "File containing columns to exclude (one per line)")
		showExcludeCols = flag.Bool("show-exclude-columns", false, "Show list of columns from exclude file and exit")
		verbose         = flag.Bool("verbose", false, "Enable verbose logging")
		findReferences  = flag.Bool("find-references", false, "Find all references to a table/column instead of comparing")
		targetColumn    = flag.String("target-column", "id", "Target column to find references for (used with -find-references)")
		maxWorkers      = flag.Int("max-workers", 4, "Maximum number of concurrent workers for parallel operations (default: 4)")
		decodeUUIDs     = flag.Bool("decode-uuids", true, "Decode Base64 encoded UUIDs in output for easier database searching (default: true)")
		progressDemo    = flag.Bool("progress-demo", false, "Run progress bar demonstration with simulated delays")
	)
	flag.Parse()

	// Run progress demo if requested
	if *progressDemo {
		fmt.Println("🎯 Running Progress Bar Demo...")
		fmt.Println("This will show you how the progress bars work with realistic timing.")
		fmt.Println()

		progress.DemoProgressBars()
		return
	}

	// Show exclude columns if requested
	if *showExcludeCols {
		columns, err := models.LoadExcludeColumnsFromFile(*excludeFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading exclude columns from %s: %v\n", *excludeFile, err)
			fmt.Println("\nMake sure the exclude columns file exists. You can create it with one column name per line.")
			fmt.Println("Example content:")
			fmt.Println("  created_at")
			fmt.Println("  updated_at")
			fmt.Println("  created_by")
			os.Exit(1)
		}

		fmt.Printf("Columns from %s that will be excluded:\n", *excludeFile)
		for _, col := range columns {
			fmt.Printf("  - %s\n", col)
		}
		fmt.Printf("\nTotal: %d columns\n", len(columns))
		fmt.Println("\nUse -exclude-from-file=false to include these columns in comparison")
		fmt.Printf("Use -exclude-file=<path> to specify a different exclude columns file\n")
		os.Exit(0)
	}

	if *tableName == "" {
		fmt.Fprintf(os.Stderr, "Error: table name is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Handle find-references mode
	if *findReferences {
		handleFindReferences(*envFile, *schemaName, *tableName, *targetColumn, *outputFile, *verbose, *maxWorkers, *decodeUUIDs)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig(*envFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Override output file if provided
	if *outputFile != "" {
		cfg.OutputFile = *outputFile
	}

	if *verbose {
		log.Printf("Loaded configuration from %s", *envFile)
		log.Printf("Comparing table: %s.%s", *schemaName, *tableName)
	}

	// Connect to databases
	db1, err := database.NewConnection(cfg.Database1)
	if err != nil {
		log.Fatalf("Failed to connect to database 1: %v", err)
	}
	defer db1.Close()

	db2, err := database.NewConnection(cfg.Database2)
	if err != nil {
		log.Fatalf("Failed to connect to database 2: %v", err)
	}
	defer db2.Close()

	if *verbose {
		log.Printf("Connected to both databases successfully")
	}

	// Create match criteria
	criteria := &models.MatchCriteria{
		IncludePrimaryKey:      *includePK,
		ExcludeColumnsFromFile: *excludeFromFile,
		ExcludeColumnsFile:     *excludeFile,
	}

	if *includeCols != "" {
		criteria.Columns = strings.Split(*includeCols, ",")
		// Trim whitespace from column names
		for i, col := range criteria.Columns {
			criteria.Columns[i] = strings.TrimSpace(col)
		}
	}

	if *excludeCols != "" {
		criteria.ExcludeColumns = strings.Split(*excludeCols, ",")
		// Trim whitespace from column names
		for i, col := range criteria.ExcludeColumns {
			criteria.ExcludeColumns[i] = strings.TrimSpace(col)
		}
	}

	if *verbose {
		log.Printf("Comparison settings:")
		log.Printf("  - Max concurrent workers: %d", *maxWorkers)
		log.Printf("  - Decode Base64 UUIDs: %v", *decodeUUIDs)
		log.Printf("  - Include primary keys: %v", *includePK)
		log.Printf("  - Exclude columns from file: %v", *excludeFromFile)
		if *excludeFromFile {
			excludeCols, err := models.LoadExcludeColumnsFromFile(*excludeFile)
			if err != nil {
				log.Printf("  - Warning: Could not load exclude columns from %s: %v", *excludeFile, err)
			} else if len(excludeCols) > 0 {
				log.Printf("  - Columns to exclude: %d columns from %s", len(excludeCols), *excludeFile)
			}
		}
		if len(criteria.ExcludeColumns) > 0 {
			log.Printf("  - Additional excluded columns: %v", criteria.ExcludeColumns)
		}
		if len(criteria.Columns) > 0 {
			log.Printf("  - Specific columns to include: %v", criteria.Columns)
		}
	}

	// Create comparator with concurrent support and UUID decoding
	comp := comparator.NewComparatorWithUUIDDecoding(db1, db2, *maxWorkers, *decodeUUIDs)
	result, err := comp.CompareTable(*schemaName, *tableName, criteria)
	if err != nil {
		log.Fatalf("Failed to compare table: %v", err)
	}

	if *verbose {
		log.Printf("Comparison completed. Total rows DB1: %d, DB2: %d", result.TotalRowsDB1, result.TotalRowsDB2)
		log.Printf("Matched rows: %d, Unmatched rows: %d", result.MatchedRows, result.UnmatchedRows)
		log.Printf("Differences found: %d", len(result.Differences))
		log.Printf("Foreign key comparisons: %d", len(result.ForeignKeyResults))
	}

	// Output results
	if err := outputResults(result, cfg.OutputFile, cfg.OutputFormat); err != nil {
		log.Fatalf("Failed to output results: %v", err)
	}

	// Print summary to console
	printSummary(result)
}

// outputResults writes the comparison results to a file
func outputResults(result *models.ComparisonResult, outputFile, format string) error {
	var data []byte
	var err error

	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(result, "", "  ")
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	if outputFile == "" {
		fmt.Println(string(data))
		return nil
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Results written to: %s\n", outputFile)
	return nil
}

// printSummary prints a summary of the comparison results to console
func printSummary(result *models.ComparisonResult) {
	fmt.Printf("\n=== COMPARISON SUMMARY ===\n")
	fmt.Printf("Table: %s.%s\n", result.Schema, result.TableName)
	fmt.Printf("Timestamp: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("\n--- Row Counts ---\n")
	fmt.Printf("Database 1: %d rows\n", result.TotalRowsDB1)
	fmt.Printf("Database 2: %d rows\n", result.TotalRowsDB2)
	fmt.Printf("Matched: %d rows\n", result.MatchedRows)
	fmt.Printf("Only in DB1: %d rows\n", len(result.OnlyInDB1))
	fmt.Printf("Only in DB2: %d rows\n", len(result.OnlyInDB2))
	fmt.Printf("Rows with differences: %d\n", len(result.Differences))

	if len(result.Differences) > 0 {
		fmt.Printf("\n--- Sample Differences ---\n")
		for i, diff := range result.Differences {
			if i >= 3 { // Show only first 3 differences in summary
				fmt.Printf("... and %d more differences\n", len(result.Differences)-3)
				break
			}
			fmt.Printf("Row %s: %d column differences\n", diff.RowIdentifier, len(diff.ColumnDifferences))
			for j, colDiff := range diff.ColumnDifferences {
				if j >= 2 { // Show only first 2 column differences per row
					fmt.Printf("  ... and %d more column differences\n", len(diff.ColumnDifferences)-2)
					break
				}
				fmt.Printf("  Column '%s': DB1='%v' vs DB2='%v'\n", colDiff.ColumnName, colDiff.DB1Value, colDiff.DB2Value)
			}
		}
	}

	if len(result.ForeignKeyResults) > 0 {
		fmt.Printf("\n--- Foreign Key Analysis ---\n")
		for _, fkResult := range result.ForeignKeyResults {
			if fkResult.Error != "" {
				fmt.Printf("FK %s -> %s.%s: ERROR - %s\n",
					fkResult.ForeignKey.ColumnName,
					fkResult.ForeignKey.ReferencedSchema,
					fkResult.ForeignKey.ReferencedTable,
					fkResult.Error)
			} else {
				fmt.Printf("FK %s -> %s.%s: %d matched, %d differences\n",
					fkResult.ForeignKey.ColumnName,
					fkResult.ForeignKey.ReferencedSchema,
					fkResult.ForeignKey.ReferencedTable,
					fkResult.ComparisonResult.MatchedRows,
					len(fkResult.ComparisonResult.Differences))
			}
		}
	}

	fmt.Printf("\n=========================\n")
}

// handleFindReferences handles the find-references mode
func handleFindReferences(envFile, schemaName, tableName, targetColumn, outputFile string, verbose bool, maxWorkers int, decodeUUIDs bool) {
	if verbose {
		log.Printf("Finding references to %s.%s.%s with %d concurrent workers (UUID decoding: %v)", schemaName, tableName, targetColumn, maxWorkers, decodeUUIDs)
	}

	// Load configuration
	cfg, err := config.LoadConfig(envFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	if verbose {
		log.Printf("Loaded configuration from %s", envFile)
	}

	// Connect to databases
	db1, err := database.NewConnection(cfg.Database1)
	if err != nil {
		log.Fatalf("Failed to connect to database 1: %v", err)
	}
	defer db1.Close()

	db2, err := database.NewConnection(cfg.Database2)
	if err != nil {
		log.Fatalf("Failed to connect to database 2: %v", err)
	}
	defer db2.Close()

	if verbose {
		log.Printf("Connected to both databases successfully")
	}

	// Create comparator with concurrent support and UUID decoding
	comp := comparator.NewComparatorWithUUIDDecoding(db1, db2, maxWorkers, decodeUUIDs)

	// Find references
	result, err := comp.FindReferences(schemaName, tableName, targetColumn)
	if err != nil {
		log.Fatalf("Failed to find references: %v", err)
	}

	if verbose {
		log.Printf("Found references in %d tables", len(result.References))
	}

	// Determine output file
	outputFileName := "match_reference_result.json"
	if outputFile != "" {
		outputFileName = outputFile
	} else if cfg.OutputFile != "" {
		// Change extension to indicate reference result
		if strings.HasSuffix(cfg.OutputFile, ".json") {
			outputFileName = strings.TrimSuffix(cfg.OutputFile, ".json") + "_references.json"
		} else {
			outputFileName = cfg.OutputFile + "_references.json"
		}
	}

	// Write results to JSON file
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result to JSON: %v", err)
	}

	err = os.WriteFile(outputFileName, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write result to file: %v", err)
	}

	fmt.Printf("Reference analysis results written to: %s\n", outputFileName)

	// Print summary
	printReferenceSummary(result)
}

// printReferenceSummary prints a summary of reference analysis results
func printReferenceSummary(result *models.MatchReferenceResult) {
	fmt.Printf("\n=== REFERENCE ANALYSIS SUMMARY ===\n")
	fmt.Printf("Target: %s.%s.%s\n", result.TargetSchema, result.TargetTable, result.TargetColumn)
	fmt.Printf("Timestamp: %s\n\n", result.Timestamp.Format("2006-01-02 15:04:05"))

	fmt.Printf("--- Reference Counts ---\n")
	fmt.Printf("Referencing tables: %d\n", result.ReferencingTables)
	fmt.Printf("Total references found: %d\n\n", result.TotalReferences)

	if len(result.References) == 0 {
		fmt.Printf("No referencing tables found.\n")
		fmt.Printf("This could mean:\n")
		fmt.Printf("- No foreign keys point to %s.%s.%s\n", result.TargetSchema, result.TargetTable, result.TargetColumn)
		fmt.Printf("- The table/column doesn't exist in one or both databases\n")
		fmt.Printf("- The foreign key constraints are not properly defined\n")
	} else {
		fmt.Printf("--- References by Table ---\n")
		for _, ref := range result.References {
			fmt.Printf("Table: %s.%s (column: %s)\n", ref.Schema, ref.TableName, ref.ColumnName)
			fmt.Printf("  DB1 values: %d\n", len(ref.DB1References))
			fmt.Printf("  DB2 values: %d\n", len(ref.DB2References))
			fmt.Printf("  Common: %d\n", len(ref.CommonRefs))
			fmt.Printf("  Only in DB1: %d\n", len(ref.OnlyInDB1))
			fmt.Printf("  Only in DB2: %d\n", len(ref.OnlyInDB2))
			if ref.ConstraintName != "" {
				fmt.Printf("  Constraint: %s\n", ref.ConstraintName)
			}
			fmt.Printf("\n")
		}
	}

	fmt.Printf("=================================\n")
}

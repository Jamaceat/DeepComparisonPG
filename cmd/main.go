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
		analyzeFKRefs   = flag.Bool("analyze-fk-references", false, "Find all tables that reference a specific ID as foreign key")
		targetID        = flag.String("id", "", "The ID value to search for in foreign key references (required with -analyze-fk-references)")
		generateScript  = flag.Bool("generate-update-script", false, "Generate SQL script to update foreign key references and delete original record")
		sourceDB        = flag.String("source-db", "db1", "Source database for script generation: 'db1' or 'db2' (default: db1)")
		idTarget        = flag.String("id-target", "", "Target ID to be replaced (required with -generate-update-script)")
		idDestination   = flag.String("id-destination", "", "Destination ID to replace with (required with -generate-update-script)")
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

	// Handle generate-update-script mode
	if *generateScript {
		if *sourceDB != "db1" && *sourceDB != "db2" {
			fmt.Fprintf(os.Stderr, "Error: source-db must be 'db1' or 'db2' when using -generate-update-script\n")
			fmt.Fprintf(os.Stderr, "Usage: -generate-update-script -table=<table> [-source-db=<db1|db2>] -id-target=<id> -id-destination=<id>\n")
			os.Exit(1)
		}
		if *idTarget == "" {
			fmt.Fprintf(os.Stderr, "Error: id-target is required when using -generate-update-script\n")
			fmt.Fprintf(os.Stderr, "Usage: -generate-update-script -table=<table> [-source-db=<db1|db2>] -id-target=<id> -id-destination=<id>\n")
			os.Exit(1)
		}
		if *idDestination == "" {
			fmt.Fprintf(os.Stderr, "Error: id-destination is required when using -generate-update-script\n")
			fmt.Fprintf(os.Stderr, "Usage: -generate-update-script -table=<table> [-source-db=<db1|db2>] -id-target=<id> -id-destination=<id>\n")
			os.Exit(1)
		}
		handleGenerateUpdateScript(*envFile, *schemaName, *tableName, *sourceDB, *idTarget, *idDestination, *outputFile, *verbose, *maxWorkers)
		return
	}

	// Handle analyze-fk-references mode
	if *analyzeFKRefs {
		if *targetID == "" {
			fmt.Fprintf(os.Stderr, "Error: ID value is required when using -analyze-fk-references\n")
			fmt.Fprintf(os.Stderr, "Usage: -analyze-fk-references -table=<table> -id=<id_value>\n")
			os.Exit(1)
		}
		handleAnalyzeFKReferences(*envFile, *schemaName, *tableName, *targetID, *outputFile, *verbose, *maxWorkers, *decodeUUIDs)
		return
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

// ensureGeneratedPath creates the generated directory if it doesn't exist and returns the full path
func ensureGeneratedPath(filename string) (string, error) {
	generatedDir := "generated"

	// Create generated directory if it doesn't exist
	if err := os.MkdirAll(generatedDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create generated directory: %w", err)
	}

	// If filename is empty, return just the directory
	if filename == "" {
		return generatedDir, nil
	}

	// If filename already starts with "generated/", return it as-is
	if strings.HasPrefix(filename, "generated/") {
		return filename, nil
	}

	return fmt.Sprintf("%s/%s", generatedDir, filename), nil
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

	// Ensure output file is in generated directory
	outputPath, err := ensureGeneratedPath(outputFile)
	if err != nil {
		return fmt.Errorf("failed to prepare output path: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Results written to: %s\n", outputPath)
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

	// Ensure output file is in generated directory
	outputPath, err := ensureGeneratedPath(outputFileName)
	if err != nil {
		log.Fatalf("Failed to prepare output path: %v", err)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write result to file: %v", err)
	}

	fmt.Printf("Reference analysis results written to: %s\n", outputPath)

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

// handleAnalyzeFKReferences handles the analyze-fk-references mode
func handleAnalyzeFKReferences(envFile, schemaName, tableName, targetID, outputFile string, verbose bool, maxWorkers int, decodeUUIDs bool) {
	if verbose {
		log.Printf("Analyzing foreign key references for ID '%s' in table %s.%s with %d concurrent workers (UUID decoding: %v)",
			targetID, schemaName, tableName, maxWorkers, decodeUUIDs)
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

	// Analyze FK references
	result, err := comp.AnalyzeFKReferences(schemaName, tableName, targetID)
	if err != nil {
		log.Fatalf("Failed to analyze FK references: %v", err)
	}

	if verbose {
		log.Printf("Found references in %d tables", len(result.ReferencingTables))
	}

	// Determine output file
	outputFileName := "id_matches_tables.json"
	if outputFile != "" {
		outputFileName = outputFile
	}

	// Write results to JSON file
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result to JSON: %v", err)
	}

	// Ensure output file is in generated directory
	outputPath, err := ensureGeneratedPath(outputFileName)
	if err != nil {
		log.Fatalf("Failed to prepare output path: %v", err)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write result to file: %v", err)
	}

	fmt.Printf("FK reference analysis results written to: %s\n", outputPath)

	// Print summary
	printFKAnalysisSummary(result)
}

// printFKAnalysisSummary prints a summary of FK analysis results
func printFKAnalysisSummary(result *models.FKAnalysisResult) {
	fmt.Printf("\n=== FK REFERENCE ANALYSIS SUMMARY ===\n")
	fmt.Printf("Target: %s.%s (ID: %s)\n", result.TargetSchema, result.TargetTable, result.TargetID)
	fmt.Printf("Timestamp: %s\n\n", result.Timestamp.Format("2006-01-02 15:04:05"))

	fmt.Printf("--- Reference Counts ---\n")
	fmt.Printf("Tables with references: %d\n", len(result.ReferencingTables))
	fmt.Printf("Total FK constraints: %d\n\n", result.TotalConstraints)

	if len(result.ReferencingTables) == 0 {
		fmt.Printf("No foreign key references found for ID '%s' in table %s.%s.\n",
			result.TargetID, result.TargetSchema, result.TargetTable)
		fmt.Printf("This could mean:\n")
		fmt.Printf("- No foreign keys reference this table\n")
		fmt.Printf("- The ID value doesn't exist in referencing tables\n")
		fmt.Printf("- Foreign key constraints are not properly defined\n")
	} else {
		fmt.Printf("--- Tables with References ---\n")
		for _, ref := range result.ReferencingTables {
			fmt.Printf("Table: %s.%s\n", ref.Schema, ref.TableName)
			fmt.Printf("  Column: %s\n", ref.ColumnName)
			fmt.Printf("  Constraint: %s\n", ref.ConstraintName)
			fmt.Printf("  DB1 matches: %d\n", ref.MatchesDB1)
			fmt.Printf("  DB2 matches: %d\n", ref.MatchesDB2)

			if len(ref.SampleRows) > 0 {
				fmt.Printf("  Sample references:\n")
				for i, sample := range ref.SampleRows {
					if i >= 3 { // Show only first 3 samples
						fmt.Printf("    ... and %d more\n", len(ref.SampleRows)-3)
						break
					}
					fmt.Printf("    %s\n", sample)
				}
			}
			fmt.Printf("\n")
		}
	}

	fmt.Printf("=====================================\n")
}

func handleGenerateUpdateScript(envFile, schemaName, tableName, sourceDB, idTarget, idDestination, outputFile string, verbose bool, maxWorkers int) {
	if verbose {
		fmt.Printf("🔧 Generating UPDATE script for FK references\n")
		fmt.Printf("Target table: %s.%s\n", schemaName, tableName)
		fmt.Printf("Source database: %s\n", sourceDB)
		fmt.Printf("Target ID: %s -> Destination ID: %s\n", idTarget, idDestination)
		fmt.Printf("\n")
	}

	// Load configuration
	cfg, err := config.LoadConfig(envFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Connect to the specified source database
	var db *database.Connection
	var dbName string

	if sourceDB == "db1" {
		db, err = database.NewConnection(cfg.Database1)
		dbName = "Database 1"
	} else {
		db, err = database.NewConnection(cfg.Database2)
		dbName = "Database 2"
	}

	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", dbName, err)
	}
	defer db.Close()

	if verbose {
		log.Printf("Connected to %s successfully", dbName)
	}

	// Create comparator instance (we only need one database for this operation)
	comp := comparator.NewComparator(db, nil) // Generate the update script
	script, err := comp.GenerateUpdateScript(schemaName, tableName, idTarget, idDestination)
	if err != nil {
		log.Fatalf("Failed to generate update script: %v", err)
	}

	// Determine output file name
	scriptFile := "update_fk_references.sql"
	if outputFile != "" {
		// Replace extension or add .sql
		if strings.HasSuffix(outputFile, ".json") {
			scriptFile = strings.TrimSuffix(outputFile, ".json") + ".sql"
		} else if strings.HasSuffix(outputFile, ".sql") {
			scriptFile = outputFile
		} else {
			scriptFile = outputFile + ".sql"
		}
	}

	// Ensure script file is in generated directory
	scriptPath, err := ensureGeneratedPath(scriptFile)
	if err != nil {
		log.Fatalf("Failed to prepare script path: %v", err)
	}

	// Write script to file
	err = os.WriteFile(scriptPath, []byte(script), 0644)
	if err != nil {
		log.Fatalf("Failed to write script to file %s: %v", scriptPath, err)
	}

	// Print summary
	fmt.Printf("\n🎯 UPDATE Script Generation Complete!\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("Script file: %s\n", scriptPath)
	fmt.Printf("Target table: %s.%s\n", schemaName, tableName)
	fmt.Printf("Source database: %s\n", dbName)
	fmt.Printf("Operation: Update FK %s -> %s\n", idTarget, idDestination)
	fmt.Printf("\n")
	fmt.Printf("⚠️  WARNING: Review the script before execution!\n")
	fmt.Printf("This script will:\n")
	fmt.Printf("  1. Update all foreign key references\n")
	fmt.Printf("  2. Delete the original record (%s)\n", idTarget)
	fmt.Printf("\n")
	fmt.Printf("To execute: psql -d <database> -f %s\n", scriptPath)
	fmt.Printf("=====================================\n")

	if verbose {
		fmt.Printf("\nScript content:\n")
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("%s\n", script)
		fmt.Printf("----------------------------------------\n")
	}
}

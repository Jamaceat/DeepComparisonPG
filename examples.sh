#!/bin/bash

# Example usage script for Deep Database Comparator

echo "=== Deep Database Comparator - Usage Examples ==="
echo

# Check if binary exists
if [ ! -f "./deepComparator" ]; then
    echo "Building the application..."
    go build -o deepComparator ./cmd
    echo "Build completed!"
    echo
fi

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "Creating .env file from example..."
    cp .env.example .env
    echo "Please edit .env with your actual database configurations before running comparisons."
    echo
    exit 1
fi

echo "Available examples (uncomment to run):"
echo

echo "# 1. Basic table comparison (excludes audit columns by default)"
echo "# ./deepComparator -table=billing_model -verbose"
echo

echo "# 2. Compare including all columns from file"
echo "# ./deepComparator -table=billing_model -exclude-from-file=false -verbose"
echo

echo "# 3. Compare with custom exclude columns file"
echo "# ./deepComparator -table=billing_model -exclude-file=\"my_exclude_columns.txt\" -verbose"
echo

echo "# 4. Compare excluding additional specific columns"
echo "# ./deepComparator -table=billing_model -exclude=\"notes,comments\" -verbose"
echo

echo "# 5. Compare specific columns only (ignores audit exclusions)" 
echo "# ./deepComparator -table=billing_model -include=\"description,order,status,concept_id\" -verbose"
echo

echo "# 6. Compare including primary keys"
echo "# ./deepComparator -table=billing_model -include-pk=true -verbose"
echo

echo "# 7. Compare with custom output file"
echo "# ./deepComparator -table=billing_model -output=billing_comparison.json -verbose"
echo

echo "# 8. Compare different schema"
echo "# ./deepComparator -table=users -schema=auth -exclude=\"password\" -verbose"
echo

echo "# 9. Show exclude columns from file"
echo "# ./deepComparator -show-exclude-columns"
echo

echo "# 10. Use different exclude columns file"
echo "# ./deepComparator -table=billing_model -exclude-file=\"custom_exclude.txt\" -verbose"
echo

echo "To run an example, uncomment the desired line above and execute it manually."
echo "Make sure your .env file contains valid database configurations."
#!/bin/bash
# Example: Generate mock contract data and analyze it

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Generating Mock Contract Data"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Generate order data
echo "📦 Generating order data..."
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output csv --rows 1000 2>/dev/null | tail -n +4 > mock_data/orders.csv
echo "✓ Created mock_data/orders.csv (1000 rows)"

# Generate user data
echo "👥 Generating user data..."
./contracts mock --domain demo --service user-service --type odcs --name user-data.yaml --output csv --rows 500 2>/dev/null | tail -n +4 > mock_data/users.csv
echo "✓ Created mock_data/users.csv (500 rows)"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Quick Analysis with Shell Tools"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo "Orders file:"
wc -l mock_data/orders.csv
head -n 3 mock_data/orders.csv

echo ""
echo "Order statuses (count):"
awk -F',' 'NR>1 {print $6}' mock_data/orders.csv | sort | uniq -c | sort -rn

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Mock data generated successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Files created in mock_data/ directory (git-ignored)"
echo ""
echo "Next steps:"
echo "  • Load in Python: import pandas as pd; df = pd.read_csv('mock_data/orders.csv')"
echo "  • Load in R: library(tidyverse); orders <- read_csv('mock_data/orders.csv')"
echo "  • Query with SQL: duckdb -c \"SELECT * FROM 'mock_data/orders.csv' LIMIT 10\""
echo "  • Open in Excel: File → Open → mock_data/orders.csv"
echo ""
echo "See examples/README.md for more examples"

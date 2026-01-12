# Working with Mock Contract Data

The contracts CLI generates CSV output that works with any data analysis tool. This keeps the tool simple and lets you use your preferred environment.

## Generate Mock Data

```bash
# Generate CSV data for ODCS contracts
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output csv --rows 1000 > mock_data/orders.csv

./contracts mock --domain demo --service user-service --type odcs --name user-data.yaml --output csv --rows 500 > mock_data/users.csv
```

The `mock_data/` directory is git-ignored for your convenience.

## Loading Data in Different Tools

### Python (Pandas)

```python
import pandas as pd

# Load data
orders = pd.read_csv('mock_data/orders.csv')

# Convert types
orders['created_at'] = pd.to_datetime(orders['created_at'])
orders['updated_at'] = pd.to_datetime(orders['updated_at'])
orders['total_amount'] = pd.to_numeric(orders['total_amount'])

# Query and analyze
high_value = orders[orders['total_amount'] > 500]
revenue_by_status = orders.groupby('status')['total_amount'].sum()

print(orders.describe())
```

### Python (Polars - faster for large datasets)

```python
import polars as pl

# Load data
orders = pl.read_csv('mock_data/orders.csv')

# Query with expression syntax
high_value = orders.filter(pl.col('total_amount') > 500)
revenue = orders.groupby('status').agg(pl.col('total_amount').sum())

print(orders.describe())
```

### R (tidyverse)

```r
library(tidyverse)

# Load data
orders <- read_csv('mock_data/orders.csv')

# Analyze
high_value <- orders %>% 
  filter(total_amount > 500)

revenue_by_status <- orders %>% 
  group_by(status) %>% 
  summarize(
    count = n(),
    total = sum(total_amount),
    average = mean(total_amount)
  )

summary(orders)
```

### DuckDB (SQL on CSV files)

```bash
# Install DuckDB: https://duckdb.org/docs/installation/
duckdb
```

```sql
-- Query CSV directly without loading
SELECT status, 
       COUNT(*) as orders,
       SUM(total_amount) as revenue,
       AVG(total_amount) as avg_order
FROM 'mock_data/orders.csv'
GROUP BY status
ORDER BY revenue DESC;

-- Join multiple CSVs
SELECT u.username, 
       COUNT(o.id) as order_count,
       SUM(o.total_amount) as total_spent
FROM 'mock_data/orders.csv' o
JOIN 'mock_data/users.csv' u ON o.user_id = u.id
GROUP BY u.username
ORDER BY total_spent DESC
LIMIT 10;
```

### Excel / Google Sheets

Simply open the CSV file directly:
- **Excel**: File → Open → Select `orders.csv`
- **Google Sheets**: File → Import → Upload → `orders.csv`

Then use built-in functions:
- Filter: Data → Filter
- Pivot Tables: Insert → Pivot Table
- Charts: Insert → Chart

### SQLite

```bash
# Import CSV to SQLite database
sqlite3 contracts.db
```

```sql
.mode csv
.import mock_data/orders.csv orders
.import mock_data/users.csv users

-- Now query with SQL
SELECT * FROM orders WHERE total_amount > 500;

SELECT status, COUNT(*), SUM(total_amount) 
FROM orders 
GROUP BY status;
```

### Jupyter Notebook

```python
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# Load data
orders = pd.read_csv('mock_data/orders.csv', parse_dates=['created_at', 'updated_at'])

# Visualize
plt.figure(figsize=(10, 6))
orders['total_amount'].hist(bins=50)
plt.title('Order Amount Distribution')
plt.xlabel('Amount')
plt.ylabel('Frequency')
plt.show()

# Box plot by status
plt.figure(figsize=(10, 6))
orders.boxplot(column='total_amount', by='status')
plt.title('Order Amount by Status')
plt.show()
```

## Example Analysis Workflows

### Quick Statistics (Bash)

```bash
# Generate data
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output csv --rows 10000 > mock_data/orders.csv

# Row count
wc -l mock_data/orders.csv

# Preview first rows
head -n 10 mock_data/orders.csv

# Get unique statuses
cut -d',' -f4 mock_data/orders.csv | sort | uniq -c
```

### Data Quality Check (Python)

```python
import pandas as pd

orders = pd.read_csv('mock_data/orders.csv')

# Check for nulls
print("Missing values:")
print(orders.isnull().sum())

# Check data types
print("\nData types:")
print(orders.dtypes)

# Check value ranges
print("\nValue ranges:")
print(orders.describe())

# Check unique values
print(f"\nUnique statuses: {orders['status'].nunique()}")
print(orders['status'].value_counts())
```

### Performance Testing (Python)

```python
import pandas as pd
import time

# Generate large dataset
# ./contracts mock ... --rows 1000000 > mock_data/large.csv

start = time.time()
df = pd.read_csv('mock_data/large.csv')
load_time = time.time() - start

print(f"Loaded {len(df)} rows in {load_time:.2f} seconds")
print(f"Memory usage: {df.memory_usage(deep=True).sum() / 1024**2:.2f} MB")

# Test query performance
start = time.time()
result = df[df['total_amount'] > 500].groupby('status')['total_amount'].sum()
query_time = time.time() - start

print(f"Query completed in {query_time:.2f} seconds")
```

## Tips

1. **Generate what you need**: Adjust `--rows` based on your testing needs
2. **Use appropriate tools**: 
   - Small datasets (< 100k rows): Any tool works
   - Medium (100k-10M rows): Pandas, R, DuckDB
   - Large (> 10M rows): Polars, DuckDB, Spark
3. **Version your test data**: Keep generated CSVs in `mock_data/` for reproducible tests
4. **Combine with real schemas**: Use mock data to test pipelines before real data arrives

## Why CSV?

- **Universal**: Works everywhere
- **Simple**: Easy to inspect and debug
- **Fast**: Efficient for most use cases
- **Flexible**: Choose your own analysis tools
- **Composable**: Follows Unix philosophy

The CLI does one thing well: generate realistic mock data from contracts. You control how to analyze it.

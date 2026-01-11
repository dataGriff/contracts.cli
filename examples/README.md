# Python Integration Quick Start

This directory contains Python examples for working with contract mock data interactively.

## Setup

```bash
# Install Python dependencies
pip install -r requirements.txt

# Make sure the Go CLI is built
go build -o contracts ./cmd/contracts
```

## Examples

### 1. Quick Demo (`interactive_demo.py`)
Shows 6 practical examples of data analysis:
```bash
python examples/interactive_demo.py
```

### 2. Comprehensive Analysis (`analyze_mock_data.py`)
Detailed analysis of orders and users with combined dataset queries:
```bash
python examples/analyze_mock_data.py
```

### 3. Interactive Python Session
Load the CLI wrapper and start exploring:
```bash
python -i contracts_cli.py
```

## Quick Code Snippets

### Generate Mock Data
```python
from contracts_cli import quick_mock

# Get a DataFrame instantly
df = quick_mock(rows=1000)
print(df.head())
```

### Filter and Query
```python
# SQL-like queries
high_value = df.query('total_amount > 500')
delivered = df.query('status == "delivered"')
recent = df[df['created_at'] > '2025-06-01']
```

### Aggregate Data
```python
# Group and summarize
by_status = df.groupby('status').agg({
    'total_amount': ['count', 'sum', 'mean', 'max']
})

# Multiple aggregations
summary = df.groupby(['status', df['created_at'].dt.month]).size()
```

### Export Data
```python
# Save to various formats
df.to_csv('mock_data.csv', index=False)
df.to_excel('mock_data.xlsx', index=False)
df.to_json('mock_data.json', orient='records')
df.to_parquet('mock_data.parquet')
```

### Statistical Analysis
```python
# Descriptive statistics
print(df.describe())
print(df['total_amount'].quantile([0.25, 0.5, 0.75, 0.9, 0.95]))

# Correlation
df[['total_amount']].corr()
```

### Visualization (requires matplotlib)
```python
import matplotlib.pyplot as plt

# Histogram
df['total_amount'].hist(bins=50)
plt.title('Order Amount Distribution')
plt.show()

# Box plot by status
df.boxplot(column='total_amount', by='status')
plt.show()

# Time series
df.set_index('created_at')['total_amount'].plot()
plt.show()
```

## Jupyter Notebook

For interactive exploration, use Jupyter:

```bash
pip install jupyter
jupyter notebook
```

Then create a notebook with:
```python
from contracts_cli import ContractsCLI
import pandas as pd
import matplotlib.pyplot as plt

cli = ContractsCLI()
orders = cli.mock_data('demo', 'order-service', 'order-data.yaml', rows=5000)

# Start exploring...
```

## API Reference

### `ContractsCLI` Class

```python
cli = ContractsCLI(cli_path='./contracts', domains_dir=None)
```

Methods:
- `list_domains()` → List all domains
- `list_services(domain)` → List services in a domain
- `list_contracts(domain, service)` → List contracts for a service
- `mock_data(domain, service, contract_name, rows=100)` → Generate mock DataFrame
- `print_contract(domain, service, contract_type, contract_name)` → Get contract content

### `quick_mock()` Function

```python
df = quick_mock(
    domain='demo',
    service='order-service',
    contract='order-data.yaml',
    rows=100
)
```

Convenience function to quickly generate mock data from the demo domain.

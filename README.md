# contracts.cli

A CLI tool that allows you to interact with different contract types including OpenAPI, Open Data Contract Specification (ODCS 3.1), and AsyncAPI.

## Features

- **List** domains, services, and contracts
- **Print** contract contents
- **Mock** contracts to generate sample data
- Support for multiple contract types:
  - OpenAPI 3.0
  - Open Data Contract Specification (ODCS 3.1)
  - AsyncAPI 2.6
- Organized folder structure: domains → services → contracts
- Built-in demo domain with example contracts
- **Python integration** for interactive data analysis with Pandas

## Installation

### From Source

```bash
git clone https://github.com/dataGriff/contracts.cli.git
cd contracts.cli
go build -o contracts ./cmd/contracts
```

### Using Go Install

```bash
go install github.com/dataGriff/contracts.cli/cmd/contracts@latest
```

## Usage

The CLI provides several commands to interact with contracts:

### List Domains

```bash
./contracts list domains
```

### List Services in a Domain

```bash
./contracts list services --domain demo
```

### List Contracts for a Service

```bash
./contracts list contracts --domain demo --service user-service
```

### List All Contracts

```bash
./contracts list contracts --all
```

### Print a Contract

```bash
./contracts print --domain demo --service user-service --type openapi --name user-api.yaml
```

### Generate Mock Data from a Contract

Generate mock data in different formats:

```bash
# Default JSON output for all contract types
./contracts mock --domain demo --service user-service --type openapi --name user-api.yaml

# For ODCS (data contracts), generate interactive dataframes
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output dataframe --rows 10

# Generate CSV output for easy export
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output csv --rows 20

# Customize the number of rows (default is 10)
./contracts mock --domain demo --service user-service --type odcs --name user-data.yaml --output dataframe --rows 50
```

#### Mock Output Formats

**For ODCS (Data Contracts):**
- `--output dataframe` (default for ODCS): Interactive tabular view with metadata
- `--output csv`: CSV format for easy export and analysis
- `--output json`: Traditional JSON format
- `--rows N`: Number of mock rows to generate (default: 10)

**For OpenAPI and AsyncAPI:**
- JSON format with mock endpoints/events

The dataframe output provides:
- Formatted table view of mock data
- Column names and data types
- Dimensions (rows × columns)
- Realistic mock data based on field types (UUIDs, timestamps, decimals, etc.)

#### Exporting Mock Data

You can redirect CSV output to a file for analysis:

```bash
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output csv --rows 100 > mock_data.csv
```

## Folder Structure

Contracts are organized in the following structure:

```
domains/
  └── {domain-name}/
      └── services/
          └── {service-name}/
              └── contracts/
                  ├── openapi/
                  │   └── *.yaml
                  ├── odcs/
                  │   └── *.yaml
                  └── asyncapi/
                      └── *.yaml
```

## Demo Domain

The repository includes a built-in demo domain with two services:
- **user-service**: User management API and data contracts
- **order-service**: Order management API and data contracts

Each service includes examples of all three contract types (OpenAPI, ODCS, AsyncAPI).

## Custom Domains

You can specify a custom domains directory using the `--domains-dir` flag:

```bash
./contracts --domains-dir /path/to/domains list domains
```

## Contract Types

### OpenAPI
OpenAPI specifications define REST API contracts including endpoints, request/response schemas, and operations.

### ODCS (Open Data Contract Specification)
ODCS 3.1 contracts define data models, schemas, and quality specifications for data products.

### AsyncAPI
AsyncAPI specifications define event-driven and message-based API contracts.

## Python Integration

For interactive data analysis and querying, use the Python wrapper with Pandas:

### Setup

```bash
# Install Python dependencies
pip install -r requirements.txt
```

### Quick Start

```python
from contracts_cli import ContractsCLI, quick_mock

# Quick way to get a DataFrame
df = quick_mock(domain='demo', service='order-service', contract='order-data.yaml', rows=1000)

# Analyze with Pandas
print(df.describe())
print(df.query('total_amount > 500'))
print(df.groupby('status')['total_amount'].sum())
```

### Advanced Usage

```python
from contracts_cli import ContractsCLI
import pandas as pd

# Initialize the wrapper
cli = ContractsCLI()

# List available contracts
domains = cli.list_domains()
services = cli.list_services('demo')
contracts = cli.list_contracts('demo', 'order-service')

# Generate and analyze mock data
orders = cli.mock_data('demo', 'order-service', 'order-data.yaml', rows=5000)
users = cli.mock_data('demo', 'user-service', 'user-data.yaml', rows=1000)

# Perform complex queries
high_value_orders = orders[orders['total_amount'] > 500]
revenue_by_status = orders.groupby('status').agg({
    'total_amount': ['sum', 'mean', 'count']
})

# Join datasets
combined = orders.merge(users, left_on='user_id', right_on='id')
user_spending = combined.groupby('username')['total_amount'].sum()
```

### Run Analysis Examples

```bash
# Run comprehensive data analysis example
python examples/analyze_mock_data.py

# Start interactive Python session with data loaded
python -i contracts_cli.py
```

### Why Python Integration?

- **Rich querying**: Use Pandas' powerful query syntax
- **Visualization**: Plot data with matplotlib, seaborn, plotly
- **Statistics**: Perform statistical analysis with scipy, statsmodels
- **Machine Learning**: Use scikit-learn, tensorflow for ML on mock data
- **Interactive exploration**: Use Jupyter notebooks for exploratory analysis

The Go CLI handles fast data generation, while Python provides the interactive analysis capabilities.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

MIT License

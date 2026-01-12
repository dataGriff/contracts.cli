# contracts.cli

A CLI tool that allows you to interact with different contract types including OpenAPI, Open Data Contract Specification (ODCS 3.1), and AsyncAPI.

## Features

- **List** domains, services, and contracts
- **Print** contract contents
- **Mock** contracts to generate sample data in CSV, JSON, or table format
- Support for multiple contract types:
  - OpenAPI 3.0
  - Open Data Contract Specification (ODCS 3.1)
  - AsyncAPI 2.6
- Organized folder structure: domains → services → contracts
- Built-in demo domain with example contracts
- **Universal CSV output** for use with any data analysis tool

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
# Generate CSV for data analysis (recommended)
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output csv --rows 1000 > mock_data/orders.csv

# View as formatted table
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output dataframe --rows 10

# Generate JSON for API testing
./contracts mock --domain demo --service user-service --type openapi --name user-api.yaml
```

#### Mock Output Formats

**For ODCS (Data Contracts):**
- `--output csv`: CSV format - universal, works with any tool
- `--output dataframe`: Formatted table view with metadata
- `--output json`: Traditional JSON format
- `--rows N`: Number of mock rows to generate (default: 10)

**For OpenAPI and AsyncAPI:**
- JSON format with mock endpoints/events

#### Working with Generated Data

Save CSV files to the `mock_data/` directory (git-ignored):

```bash
# Generate mock data
./contracts mock --domain demo --service order-service --type odcs --name order-data.yaml --output csv --rows 5000 > mock_data/orders.csv
./contracts mock --domain demo --service user-service --type odcs --name user-data.yaml --output csv --rows 1000 > mock_data/users.csv
```

Then load in your preferred tool:

**Python (Pandas):**
```python
import pandas as pd
orders = pd.read_csv('mock_data/orders.csv')
orders['created_at'] = pd.to_datetime(orders['created_at'])
print(orders.query('total_amount > 500'))
```

**R:**
```r
library(tidyverse)
orders <- read_csv('mock_data/orders.csv')
orders %>% filter(total_amount > 500)
```

**DuckDB (SQL):**
```sql
SELECT status, SUM(total_amount) 
FROM 'mock_data/orders.csv' 
GROUP BY status;
```

See [examples/README.md](examples/README.md) for complete examples in Python, R, DuckDB, Excel, and more.

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

## Why CSV Output?

The CLI generates universal CSV output that works with any analysis tool:

- **Universal**: Python, R, Excel, SQL databases, BI tools - everything reads CSV
- **Simple**: Easy to inspect, debug, and version control
- **Fast**: Efficient for datasets from 100 to millions of rows
- **Flexible**: Use your preferred tools and workflows
- **Composable**: Follows Unix philosophy - do one thing well

The CLI focuses on generating realistic mock data. You control how to analyze it.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

MIT License

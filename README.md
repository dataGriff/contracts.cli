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

```bash
./contracts mock --domain demo --service user-service --type openapi --name user-api.yaml
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

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

MIT License

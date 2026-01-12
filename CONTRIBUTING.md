# Contributing to contracts.cli

Thank you for your interest in contributing to contracts.cli! This document provides guidelines and instructions for contributing to the project.

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git

### Clone the Repository

```bash
git clone https://github.com/dataGriff/contracts.cli.git
cd contracts.cli
```

### Install Dependencies

```bash
go mod download
```

### Build the Project

```bash
go build -o contracts ./cmd/contracts
```

## Project Structure

```
contracts.cli/
├── cmd/
│   └── contracts/          # Main entry point
│       └── main.go
├── internal/
│   ├── cli/               # CLI commands and logic
│   │   └── root.go
│   ├── contracts/         # Contract repository management
│   │   └── repository.go
│   └── mock/              # Mock data generation
│       └── generator.go
├── domains/               # Contract storage (demo domain included)
│   └── demo/
│       └── services/
│           ├── user-service/
│           └── order-service/
├── go.mod
├── go.sum
├── README.md
└── CONTRIBUTING.md
```

## Making Changes

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

Follow these guidelines:

- **Isolation**: Keep CLI tooling (in `internal/cli`, `internal/mock`) separate from contract data (in `domains/`)
- **Code Style**: Follow standard Go conventions and use `gofmt` to format your code
- **Documentation**: Update README.md if you add new features
- **Testing**: Add tests for new functionality (when test infrastructure is added)

### 3. Test Your Changes

Build and test the CLI:

```bash
# Build
go build -o contracts ./cmd/contracts

# Test commands
./contracts list domains
./contracts list services --domain demo
./contracts list contracts --all
./contracts print --domain demo --service user-service --type openapi --name user-api.yaml
./contracts mock --domain demo --service user-service --type openapi --name user-api.yaml
```

### 4. Format Your Code

```bash
go fmt ./...
```

### 5. Run Go Mod Tidy

```bash
go mod tidy
```

### 6. Commit Your Changes

```bash
git add .
git commit -m "Brief description of your changes"
```

### 7. Push and Create a Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Adding New Contract Types

To add support for a new contract type:

1. Update the `ContractType` enum in `internal/contracts/repository.go`
2. Add a new mock generator method in `internal/mock/generator.go`
3. Update the `GenerateMock` switch statement to handle the new type
4. Add example contracts to the demo domain
5. Update README.md with documentation for the new contract type

## Adding New Commands

To add a new CLI command:

1. Create a new command function in `internal/cli/root.go` (or create a new file)
2. Register the command in the root command using `AddCommand()`
3. Follow the existing pattern using Cobra flags and command structure
4. Test the command thoroughly
5. Update README.md with usage examples

## Code Quality

- Run `go vet ./...` to check for common issues
- Run `go fmt ./...` to format code
- Keep functions small and focused
- Add comments for exported functions and types
- Use meaningful variable and function names

## Reporting Issues

If you find a bug or have a feature request:

1. Check if the issue already exists
2. If not, create a new issue with:
   - Clear description of the problem or feature
   - Steps to reproduce (for bugs)
   - Expected vs actual behavior
   - Your environment (Go version, OS, etc.)

## Questions?

Feel free to open an issue for questions or discussions about the project.

## License

By contributing to contracts.cli, you agree that your contributions will be licensed under the MIT License.
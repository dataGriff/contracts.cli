package contracts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ContractType represents the type of contract
type ContractType string

const (
	OpenAPI   ContractType = "openapi"
	ODCS      ContractType = "odcs"
	AsyncAPI  ContractType = "asyncapi"
)

// Contract represents a contract file
type Contract struct {
	Domain   string
	Service  string
	Type     ContractType
	Name     string
	Path     string
	Content  string
}

// Repository manages contract discovery and retrieval
type Repository struct {
	baseDir string
}

// NewRepository creates a new contract repository
func NewRepository(baseDir string) *Repository {
	return &Repository{baseDir: baseDir}
}

// ListDomains returns all domains in the repository
func (r *Repository) ListDomains() ([]string, error) {
	domains := []string{}
	entries, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read domains directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			domains = append(domains, entry.Name())
		}
	}

	return domains, nil
}

// ListServices returns all services in a domain
func (r *Repository) ListServices(domain string) ([]string, error) {
	services := []string{}
	servicesDir := filepath.Join(r.baseDir, domain, "services")
	
	if _, err := os.Stat(servicesDir); os.IsNotExist(err) {
		return services, nil
	}

	entries, err := os.ReadDir(servicesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read services directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			services = append(services, entry.Name())
		}
	}

	return services, nil
}

// ListContracts returns all contracts for a service
func (r *Repository) ListContracts(domain, service string) ([]Contract, error) {
	contracts := []Contract{}
	contractsDir := filepath.Join(r.baseDir, domain, "services", service, "contracts")

	if _, err := os.Stat(contractsDir); os.IsNotExist(err) {
		return contracts, nil
	}

	// Walk through each contract type directory
	contractTypes := []ContractType{OpenAPI, ODCS, AsyncAPI}
	for _, cType := range contractTypes {
		typeDir := filepath.Join(contractsDir, string(cType))
		if _, err := os.Stat(typeDir); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(typeDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml")) {
				contracts = append(contracts, Contract{
					Domain:  domain,
					Service: service,
					Type:    cType,
					Name:    entry.Name(),
					Path:    filepath.Join(typeDir, entry.Name()),
				})
			}
		}
	}

	return contracts, nil
}

// GetContract retrieves a specific contract with its content
func (r *Repository) GetContract(domain, service, contractType, name string) (*Contract, error) {
	contractPath := filepath.Join(r.baseDir, domain, "services", service, "contracts", contractType, name)
	
	if _, err := os.Stat(contractPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("contract not found: %s", contractPath)
	}

	content, err := os.ReadFile(contractPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read contract: %w", err)
	}

	return &Contract{
		Domain:  domain,
		Service: service,
		Type:    ContractType(contractType),
		Name:    name,
		Path:    contractPath,
		Content: string(content),
	}, nil
}

// ListAllContracts returns all contracts across all domains and services
func (r *Repository) ListAllContracts() ([]Contract, error) {
	allContracts := []Contract{}

	domains, err := r.ListDomains()
	if err != nil {
		return nil, err
	}

	for _, domain := range domains {
		services, err := r.ListServices(domain)
		if err != nil {
			continue
		}

		for _, service := range services {
			contracts, err := r.ListContracts(domain, service)
			if err != nil {
				continue
			}
			allContracts = append(allContracts, contracts...)
		}
	}

	return allContracts, nil
}

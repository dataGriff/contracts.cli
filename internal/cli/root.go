package cli

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/dataGriff/contracts.cli/internal/contracts"
	"github.com/dataGriff/contracts.cli/internal/mock"
	"github.com/spf13/cobra"
)

var (
	domainsDir string
	repo       *contracts.Repository
	mockGen    *mock.Generator
)

// NewRootCommand creates the root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "contracts",
		Short: "A CLI tool for managing and interacting with contracts",
		Long: `contracts.cli is a tool that allows you to list, print, and mock different contract types:
- OpenAPI specifications
- Open Data Contract Specification (ODCS 3.1)
- AsyncAPI specifications

Contracts are organized in a domains -> services -> contracts structure.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize repository and mock generator
			if domainsDir == "" {
				// Use built-in demo domain by default
				execPath, err := os.Executable()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error determining executable path: %v\n", err)
					os.Exit(1)
				}
				domainsDir = filepath.Join(filepath.Dir(execPath), "domains")

				// If domains dir doesn't exist, try current directory
				if _, err := os.Stat(domainsDir); os.IsNotExist(err) {
					cwd, err := os.Getwd()
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
						os.Exit(1)
					}
					domainsDir = filepath.Join(cwd, "domains")
				}
			}

			repo = contracts.NewRepository(domainsDir)
			mockGen = mock.NewGenerator()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&domainsDir, "domains-dir", "d", "", "Path to domains directory (default: ./domains or built-in demo)")

	// Add subcommands
	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newPrintCommand())
	rootCmd.AddCommand(newMockCommand())

	return rootCmd
}

// newListCommand creates the list command with subcommands
func newListCommand() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List domains, services, or contracts",
		Long:  "List domains, services, or contracts in the repository",
	}

	listCmd.AddCommand(newListDomainsCommand())
	listCmd.AddCommand(newListServicesCommand())
	listCmd.AddCommand(newListContractsCommand())

	return listCmd
}

func newListDomainsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "domains",
		Short: "List all domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			domains, err := repo.ListDomains()
			if err != nil {
				return err
			}

			if len(domains) == 0 {
				fmt.Println("No domains found")
				return nil
			}

			fmt.Println("Domains:")
			for _, domain := range domains {
				fmt.Printf("  - %s\n", domain)
			}

			return nil
		},
	}
}

func newListServicesCommand() *cobra.Command {
	var domain string

	cmd := &cobra.Command{
		Use:   "services",
		Short: "List all services in a domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" {
				return fmt.Errorf("domain is required (use --domain flag)")
			}

			services, err := repo.ListServices(domain)
			if err != nil {
				return err
			}

			if len(services) == 0 {
				fmt.Printf("No services found in domain '%s'\n", domain)
				return nil
			}

			fmt.Printf("Services in domain '%s':\n", domain)
			for _, service := range services {
				fmt.Printf("  - %s\n", service)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "Domain name (required)")
	cmd.MarkFlagRequired("domain")

	return cmd
}

func newListContractsCommand() *cobra.Command {
	var (
		domain  string
		service string
		all     bool
	)

	cmd := &cobra.Command{
		Use:   "contracts",
		Short: "List contracts for a service or all contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if all {
				contracts, err := repo.ListAllContracts()
				if err != nil {
					return err
				}

				if len(contracts) == 0 {
					fmt.Println("No contracts found")
					return nil
				}

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "DOMAIN\tSERVICE\tTYPE\tNAME")
				fmt.Fprintln(w, "------\t-------\t----\t----")
				for _, contract := range contracts {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", contract.Domain, contract.Service, contract.Type, contract.Name)
				}
				w.Flush()

				return nil
			}

			if domain == "" || service == "" {
				return fmt.Errorf("domain and service are required (use --domain and --service flags, or --all for all contracts)")
			}

			contracts, err := repo.ListContracts(domain, service)
			if err != nil {
				return err
			}

			if len(contracts) == 0 {
				fmt.Printf("No contracts found for service '%s' in domain '%s'\n", service, domain)
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "Contracts for service '%s' in domain '%s':\n", service, domain)
			fmt.Fprintln(w, "TYPE\tNAME")
			fmt.Fprintln(w, "----\t----")
			for _, contract := range contracts {
				fmt.Fprintf(w, "%s\t%s\n", contract.Type, contract.Name)
			}
			w.Flush()

			return nil
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "Domain name")
	cmd.Flags().StringVar(&service, "service", "", "Service name")
	cmd.Flags().BoolVar(&all, "all", false, "List all contracts across all domains and services")

	return cmd
}

func newPrintCommand() *cobra.Command {
	var (
		domain       string
		service      string
		contractType string
		name         string
	)

	cmd := &cobra.Command{
		Use:   "print",
		Short: "Print the content of a contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" || service == "" || contractType == "" || name == "" {
				return fmt.Errorf("all flags are required: --domain, --service, --type, --name")
			}

			contract, err := repo.GetContract(domain, service, contractType, name)
			if err != nil {
				return err
			}

			fmt.Printf("Contract: %s/%s/%s/%s\n", domain, service, contractType, name)
			fmt.Println("---")
			fmt.Println(contract.Content)

			return nil
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "Domain name (required)")
	cmd.Flags().StringVar(&service, "service", "", "Service name (required)")
	cmd.Flags().StringVar(&contractType, "type", "", "Contract type: openapi, odcs, or asyncapi (required)")
	cmd.Flags().StringVar(&name, "name", "", "Contract file name (required)")

	cmd.MarkFlagRequired("domain")
	cmd.MarkFlagRequired("service")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newMockCommand() *cobra.Command {
	var (
		domain       string
		service      string
		contractType string
		name         string
		rows         int
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "mock",
		Short: "Generate mock data from a contract",
		Long: `Generate mock data from a contract.

For ODCS (data contracts), you can specify:
  --output dataframe: Display mock data as an interactive table
  --output csv: Display mock data as CSV
  --output json: Display mock data as JSON (default)
  --rows N: Number of rows to generate (default: 10)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" || service == "" || contractType == "" || name == "" {
				return fmt.Errorf("all flags are required: --domain, --service, --type, --name")
			}

			contract, err := repo.GetContract(domain, service, contractType, name)
			if err != nil {
				return err
			}

			// For ODCS contracts with dataframe or csv output
			if contract.Type == "odcs" && (outputFormat == "dataframe" || outputFormat == "csv") {
				df, err := mockGen.GenerateMockDataFrame(contract, rows)
				if err != nil {
					return fmt.Errorf("failed to generate mock dataframe: %w", err)
				}

				fmt.Printf("Mock data for: %s/%s/%s/%s\n", domain, service, contractType, name)
				fmt.Printf("Generated %d rows\n", rows)
				fmt.Println("---")

				if outputFormat == "csv" {
					// Output as CSV
					writer := csv.NewWriter(os.Stdout)
					records := df.Records()
					for _, record := range records {
						writer.Write(record)
					}
					writer.Flush()
				} else {
					// Output as formatted table
					fmt.Println(df)
					fmt.Println()
					fmt.Println("DataFrame operations available:")
					rows, cols := df.Dims()
					fmt.Printf("  • Dimensions: %d rows × %d columns\n", rows, cols)
					fmt.Println("  • Column names:", df.Names())
					fmt.Println("  • Data types:", df.Types())
				}

				return nil
			}

			// Default JSON output
			mockData, err := mockGen.GenerateMock(contract)
			if err != nil {
				return fmt.Errorf("failed to generate mock: %w", err)
			}

			fmt.Printf("Mock data for: %s/%s/%s/%s\n", domain, service, contractType, name)
			fmt.Println("---")
			fmt.Println(mockData)

			return nil
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "Domain name (required)")
	cmd.Flags().StringVar(&service, "service", "", "Service name (required)")
	cmd.Flags().StringVar(&contractType, "type", "", "Contract type: openapi, odcs, or asyncapi (required)")
	cmd.Flags().StringVar(&name, "name", "", "Contract file name (required)")
	cmd.Flags().IntVar(&rows, "rows", 10, "Number of rows to generate for ODCS contracts (default: 10)")
	cmd.Flags().StringVar(&outputFormat, "output", "json", "Output format: json, csv, or dataframe (for ODCS only)")

	cmd.MarkFlagRequired("domain")
	cmd.MarkFlagRequired("service")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("name")

	return cmd
}

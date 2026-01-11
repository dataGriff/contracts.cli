#!/usr/bin/env python3
"""
Python wrapper for contracts.cli - provides seamless integration with Pandas.

This module allows you to generate mock data from ODCS contracts and immediately
work with it as a Pandas DataFrame for analysis, querying, and visualization.
"""

import subprocess
import pandas as pd
from io import StringIO
from typing import Optional, Literal
import os
import sys


class ContractsCLI:
    """Python wrapper for the contracts CLI tool."""
    
    def __init__(self, cli_path: str = "./contracts", domains_dir: Optional[str] = None):
        """
        Initialize the CLI wrapper.
        
        Args:
            cli_path: Path to the contracts CLI binary
            domains_dir: Optional custom domains directory
        """
        self.cli_path = cli_path
        self.domains_dir = domains_dir
        
        # Verify CLI exists
        if not os.path.exists(cli_path):
            raise FileNotFoundError(
                f"Contracts CLI not found at {cli_path}. "
                "Please build it first: go build -o contracts ./cmd/contracts"
            )
    
    def _run_command(self, args: list[str]) -> str:
        """Run a CLI command and return stdout."""
        cmd = [self.cli_path]
        if self.domains_dir:
            cmd.extend(['--domains-dir', self.domains_dir])
        cmd.extend(args)
        
        try:
            result = subprocess.run(
                cmd, 
                capture_output=True, 
                text=True, 
                check=True
            )
            return result.stdout
        except subprocess.CalledProcessError as e:
            raise RuntimeError(f"CLI command failed: {e.stderr}")
    
    def list_domains(self) -> list[str]:
        """List all available domains."""
        output = self._run_command(['list', 'domains'])
        domains = []
        for line in output.split('\n'):
            if line.strip().startswith('- '):
                domains.append(line.strip()[2:])
        return domains
    
    def list_services(self, domain: str) -> list[str]:
        """List all services in a domain."""
        output = self._run_command(['list', 'services', '--domain', domain])
        services = []
        for line in output.split('\n'):
            if line.strip().startswith('- '):
                services.append(line.strip()[2:])
        return services
    
    def list_contracts(self, domain: str, service: str) -> dict:
        """
        List all contracts for a service.
        
        Returns:
            Dictionary with contract types as keys and lists of contract names as values
        """
        output = self._run_command([
            'list', 'contracts', 
            '--domain', domain,
            '--service', service
        ])
        
        contracts = {'openapi': [], 'odcs': [], 'asyncapi': []}
        current_type = None
        
        for line in output.split('\n'):
            line = line.strip()
            if 'openapi:' in line.lower():
                current_type = 'openapi'
            elif 'odcs:' in line.lower():
                current_type = 'odcs'
            elif 'asyncapi:' in line.lower():
                current_type = 'asyncapi'
            elif line.startswith('- ') and current_type:
                contracts[current_type].append(line[2:])
        
        return contracts
    
    def mock_data(
        self, 
        domain: str,
        service: str,
        contract_name: str,
        rows: int = 100,
        contract_type: Literal['odcs'] = 'odcs'
    ) -> pd.DataFrame:
        """
        Generate mock data from an ODCS contract and return as a Pandas DataFrame.
        
        Args:
            domain: Domain name
            service: Service name
            contract_name: Contract file name (e.g., 'order-data.yaml')
            rows: Number of rows to generate
            contract_type: Contract type (currently only 'odcs' supported)
        
        Returns:
            Pandas DataFrame with mock data
        
        Example:
            >>> cli = ContractsCLI()
            >>> df = cli.mock_data('demo', 'order-service', 'order-data.yaml', rows=1000)
            >>> df.query('total_amount > 100')
        """
        if contract_type != 'odcs':
            raise ValueError("DataFrame generation only supported for ODCS contracts")
        
        output = self._run_command([
            'mock',
            '--domain', domain,
            '--service', service,
            '--type', contract_type,
            '--name', contract_name,
            '--output', 'csv',
            '--rows', str(rows)
        ])
        
        # Skip the header lines (Mock data for: ..., Generated X rows, ---)
        lines = output.split('\n')
        csv_start = 0
        for i, line in enumerate(lines):
            if line.strip() == '---':
                csv_start = i + 1
                break
        
        csv_content = '\n'.join(lines[csv_start:])
        
        # Parse CSV into DataFrame
        df = pd.read_csv(StringIO(csv_content))
        
        # Try to convert columns to appropriate types
        for col in df.columns:
            # Try datetime conversion for timestamp columns
            if 'created_at' in col or 'updated_at' in col or 'timestamp' in col:
                try:
                    df[col] = pd.to_datetime(df[col])
                except:
                    pass
            # Try numeric conversion for amount/price columns
            elif 'amount' in col or 'price' in col or 'total' in col:
                try:
                    df[col] = pd.to_numeric(df[col])
                except:
                    pass
        
        return df
    
    def print_contract(
        self,
        domain: str,
        service: str,
        contract_type: Literal['openapi', 'odcs', 'asyncapi'],
        contract_name: str
    ) -> str:
        """
        Print the content of a contract.
        
        Args:
            domain: Domain name
            service: Service name
            contract_type: Contract type
            contract_name: Contract file name
        
        Returns:
            Contract content as string
        """
        output = self._run_command([
            'print',
            '--domain', domain,
            '--service', service,
            '--type', contract_type,
            '--name', contract_name
        ])
        return output


def quick_mock(
    domain: str = 'demo',
    service: str = 'order-service', 
    contract: str = 'order-data.yaml',
    rows: int = 100
) -> pd.DataFrame:
    """
    Quick helper to generate mock data with default settings.
    
    Example:
        >>> df = quick_mock(rows=500)
        >>> df.describe()
    """
    cli = ContractsCLI()
    return cli.mock_data(domain, service, contract, rows)


if __name__ == '__main__':
    # Interactive mode
    print("Contracts CLI Python Wrapper")
    print("=" * 50)
    
    cli = ContractsCLI()
    
    # Show available domains
    print("\nAvailable domains:")
    domains = cli.list_domains()
    for domain in domains:
        print(f"  - {domain}")
    
    # Generate sample data
    print("\nGenerating sample data from demo domain...")
    df = cli.mock_data('demo', 'order-service', 'order-data.yaml', rows=10)
    
    print(f"\nGenerated DataFrame with {len(df)} rows:")
    print(df.head())
    
    print("\nDataFrame info:")
    print(df.info())
    
    print("\nSummary statistics:")
    print(df.describe())
    
    print("\n" + "=" * 50)
    print("Try: df = quick_mock(rows=1000)")
    print("Then: df.query('total_amount > 100')")
    print("Or:   df.groupby('status').agg({'total_amount': 'sum'})")

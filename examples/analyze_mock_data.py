#!/usr/bin/env python3
"""
Example script showing how to analyze mock contract data using Python and Pandas.

This demonstrates the power of combining the Go CLI's fast data generation
with Python's rich data analysis ecosystem.
"""

import sys
import os

# Add parent directory to path to import contracts_cli
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from contracts_cli import ContractsCLI, quick_mock
import pandas as pd


def analyze_orders():
    """Analyze mock order data."""
    print("=" * 60)
    print("Order Data Analysis")
    print("=" * 60)
    
    # Generate mock order data
    cli = ContractsCLI()
    orders = cli.mock_data('demo', 'order-service', 'order-data.yaml', rows=1000)
    
    print(f"\nGenerated {len(orders)} mock orders")
    print("\nFirst 5 orders:")
    print(orders.head())
    
    # Convert timestamps
    orders['created_at'] = pd.to_datetime(orders['created_at'])
    orders['updated_at'] = pd.to_datetime(orders['updated_at'])
    
    # Basic statistics
    print("\n" + "-" * 60)
    print("Summary Statistics")
    print("-" * 60)
    print(orders.describe())
    
    # Status breakdown
    print("\n" + "-" * 60)
    print("Orders by Status")
    print("-" * 60)
    status_counts = orders['status'].value_counts()
    print(status_counts)
    print(f"\nTotal: {status_counts.sum()}")
    
    # Amount analysis
    print("\n" + "-" * 60)
    print("Amount Analysis")
    print("-" * 60)
    print(f"Total revenue: ${orders['total_amount'].sum():,.2f}")
    print(f"Average order: ${orders['total_amount'].mean():,.2f}")
    print(f"Median order: ${orders['total_amount'].median():,.2f}")
    print(f"Largest order: ${orders['total_amount'].max():,.2f}")
    print(f"Smallest order: ${orders['total_amount'].min():,.2f}")
    
    # Revenue by status
    print("\n" + "-" * 60)
    print("Revenue by Status")
    print("-" * 60)
    revenue_by_status = orders.groupby('status')['total_amount'].agg([
        ('count', 'count'),
        ('total', 'sum'),
        ('average', 'mean')
    ])
    print(revenue_by_status.sort_values('total', ascending=False))
    
    # High-value orders
    print("\n" + "-" * 60)
    print("High-Value Orders (>$500)")
    print("-" * 60)
    high_value = orders[orders['total_amount'] > 500]
    print(f"Count: {len(high_value)} ({len(high_value)/len(orders)*100:.1f}%)")
    print(f"Total value: ${high_value['total_amount'].sum():,.2f}")
    
    # Time analysis
    print("\n" + "-" * 60)
    print("Time Analysis")
    print("-" * 60)
    orders['created_month'] = orders['created_at'].dt.to_period('M')
    monthly_orders = orders.groupby('created_month').size()
    print("Orders by month:")
    print(monthly_orders.head(10))
    
    return orders


def analyze_users():
    """Analyze mock user data."""
    print("\n\n" + "=" * 60)
    print("User Data Analysis")
    print("=" * 60)
    
    # Generate mock user data
    cli = ContractsCLI()
    users = cli.mock_data('demo', 'user-service', 'user-data.yaml', rows=500)
    
    print(f"\nGenerated {len(users)} mock users")
    print("\nFirst 5 users:")
    print(users.head())
    
    # Convert timestamps
    users['created_at'] = pd.to_datetime(users['created_at'])
    users['updated_at'] = pd.to_datetime(users['updated_at'])
    
    # Email domain analysis
    print("\n" + "-" * 60)
    print("Email Domain Analysis")
    print("-" * 60)
    users['email_domain'] = users['email'].str.split('@').str[1]
    domain_counts = users['email_domain'].value_counts()
    print(domain_counts.head(10))
    
    # User activity
    print("\n" + "-" * 60)
    print("User Activity")
    print("-" * 60)
    users['days_since_update'] = (pd.Timestamp.now() - users['updated_at']).dt.days
    users['active'] = users['days_since_update'] < 30
    
    active_count = users['active'].sum()
    print(f"Active users (updated in last 30 days): {active_count} ({active_count/len(users)*100:.1f}%)")
    print(f"Inactive users: {len(users) - active_count} ({(len(users) - active_count)/len(users)*100:.1f}%)")
    
    return users


def combined_analysis():
    """Perform combined analysis of orders and users."""
    print("\n\n" + "=" * 60)
    print("Combined Analysis: Orders & Users")
    print("=" * 60)
    
    cli = ContractsCLI()
    
    # Generate both datasets
    orders = cli.mock_data('demo', 'order-service', 'order-data.yaml', rows=1000)
    users = cli.mock_data('demo', 'user-service', 'user-data.yaml', rows=500)
    
    # Merge on user_id
    combined = orders.merge(
        users,
        left_on='user_id',
        right_on='id',
        how='left',
        suffixes=('_order', '_user')
    )
    
    print(f"\nCombined dataset: {len(combined)} rows")
    print("\nSample combined data:")
    print(combined[['username', 'email', 'total_amount', 'status']].head())
    
    # User spending analysis
    print("\n" + "-" * 60)
    print("User Spending Analysis")
    print("-" * 60)
    user_spending = combined.groupby('username').agg({
        'total_amount': ['count', 'sum', 'mean'],
        'id_order': 'first'
    })
    user_spending.columns = ['order_count', 'total_spent', 'avg_order', 'order_id']
    user_spending = user_spending.sort_values('total_spent', ascending=False)
    
    print("\nTop 10 spenders:")
    print(user_spending.head(10))
    
    print("\n" + "-" * 60)
    print("Spending Patterns")
    print("-" * 60)
    print(f"Average orders per user: {user_spending['order_count'].mean():.1f}")
    print(f"Average spending per user: ${user_spending['total_spent'].mean():,.2f}")
    print(f"Max orders by one user: {user_spending['order_count'].max()}")
    
    return combined


def interactive_session():
    """Start an interactive Python session with data loaded."""
    print("\n\n" + "=" * 60)
    print("Starting Interactive Session")
    print("=" * 60)
    print("\nLoading data...")
    
    # Load data
    cli = ContractsCLI()
    orders = cli.mock_data('demo', 'order-service', 'order-data.yaml', rows=1000)
    users = cli.mock_data('demo', 'user-service', 'user-data.yaml', rows=500)
    
    print("\nDataFrames available:")
    print("  - orders: 1000 order records")
    print("  - users: 500 user records")
    print("  - cli: ContractsCLI instance")
    
    print("\nTry these commands:")
    print("  orders.query('total_amount > 500')")
    print("  orders.groupby('status')['total_amount'].sum()")
    print("  users['email'].str.contains('example.com').sum()")
    print("  orders.plot.hist(column='total_amount', bins=50)")
    
    # Start interactive mode
    import code
    code.interact(local=locals())


if __name__ == '__main__':
    # Run all analyses
    orders_df = analyze_orders()
    users_df = analyze_users()
    combined_df = combined_analysis()
    
    # Offer interactive session
    print("\n\n" + "=" * 60)
    print("Analysis complete!")
    print("=" * 60)
    print("\nDataFrames created:")
    print(f"  - orders_df: {len(orders_df)} rows")
    print(f"  - users_df: {len(users_df)} rows")
    print(f"  - combined_df: {len(combined_df)} rows")
    
    response = input("\nStart interactive session? (y/n): ")
    if response.lower() == 'y':
        interactive_session()

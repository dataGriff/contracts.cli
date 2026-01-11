#!/usr/bin/env python3
"""
Quick interactive demo of the Python wrapper capabilities.
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from contracts_cli import quick_mock

print("=" * 70)
print("Contracts CLI - Interactive Data Demo")
print("=" * 70)

# Generate mock data
print("\n📊 Generating 500 mock orders...")
orders = quick_mock(rows=500)

print(f"✅ Generated {len(orders)} orders\n")

# Example 1: Basic filtering
print("-" * 70)
print("Example 1: High-value orders (>$700)")
print("-" * 70)
high_value = orders.query('total_amount > 700')
print(f"Found {len(high_value)} high-value orders")
print(high_value[['id', 'total_amount', 'status']].head())

# Example 2: Grouping and aggregation
print("\n" + "-" * 70)
print("Example 2: Revenue by status")
print("-" * 70)
revenue = orders.groupby('status')['total_amount'].agg(['count', 'sum', 'mean'])
revenue.columns = ['Orders', 'Total Revenue', 'Avg Order Value']
print(revenue.sort_values('Total Revenue', ascending=False))

# Example 3: Time-based analysis
print("\n" + "-" * 70)
print("Example 3: Most recent orders")
print("-" * 70)
recent = orders.nlargest(5, 'created_at')
print(recent[['id', 'status', 'total_amount', 'created_at']])

# Example 4: Statistical analysis
print("\n" + "-" * 70)
print("Example 4: Order amount distribution")
print("-" * 70)
print(orders['total_amount'].describe())

# Example 5: Conditional aggregation
print("\n" + "-" * 70)
print("Example 5: Completed vs Pending/Cancelled")
print("-" * 70)
orders['is_completed'] = orders['status'].isin(['delivered', 'shipped', 'confirmed'])
completion_stats = orders.groupby('is_completed')['total_amount'].agg(['count', 'sum', 'mean'])
completion_stats.index = ['Pending/Cancelled', 'Completed']
print(completion_stats)

# Example 6: Complex query
print("\n" + "-" * 70)
print("Example 6: High-value delivered orders from last 6 months")
print("-" * 70)
import pandas as pd
six_months_ago = pd.Timestamp.now(tz='UTC') - pd.Timedelta(days=180)
complex_query = orders.query(
    'total_amount > 600 and status == "delivered" and created_at > @six_months_ago'
)
print(f"Found {len(complex_query)} matching orders")
if len(complex_query) > 0:
    print(complex_query[['id', 'total_amount', 'created_at']].head())

print("\n" + "=" * 70)
print("🎉 Demo complete! You can now:")
print("  • Use df.query() for SQL-like filtering")
print("  • Use df.groupby() for aggregations")
print("  • Use df.plot() for visualizations")
print("  • Export with df.to_csv(), df.to_excel(), etc.")
print("=" * 70)

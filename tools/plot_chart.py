import pandas as pd
import matplotlib.pyplot as plt
import os
import argparse

def main():
    parser = argparse.ArgumentParser(description="Plot SSN Experiment Results")
    parser.add_argument("csv_file", nargs="?", default="game_analysis.csv", help="Name of the CSV file to plot")
    args = parser.parse_args()

    # Smart Path Finding
    possible_paths = [
        args.csv_file,
        f"results/{args.csv_file}",
        f"../results/{args.csv_file}"
    ]
    
    csv_path = None
    for p in possible_paths:
        if os.path.exists(p):
            csv_path = p
            break

    if csv_path is None:
        print(f"Error: File '{args.csv_file}' not found.")
        return

    print(f"Plotting: {csv_path}")

    try:
        df = pd.read_csv(csv_path)
    except Exception as e:
        print(f"Error reading CSV: {e}")
        return

    # Calculate Rolling Win Rate (window=50 or 500 dependent on size)
    window_size = 50
    if len(df) > 5000:
        window_size = 500
        
    df['RollingWR'] = df['WinRate'].rolling(window=window_size).mean()

    # Plot
    plt.figure(figsize=(12, 6))
    plt.plot(df['Step'], df['WinRate'], label='Cumulative Win Rate', color='blue', alpha=0.6)
    plt.plot(df['Step'], df['RollingWR'], label=f'Rolling Win Rate (Moving Avg {window_size})', color='red', linewidth=2)
    
    plt.axhline(y=0.33, color='gray', linestyle='--', label='Random Chance (33%)')
    plt.axhline(y=0.50, color='green', linestyle=':', label='Breakeven (50%)')
    
    plt.title(f'SSN Performance Analysis: {args.csv_file}')
    plt.xlabel('Game Step')
    plt.ylabel('Win Rate')
    plt.legend()
    plt.grid(True, alpha=0.3)
    
    output_png = csv_path.replace(".csv", ".png")
    plt.savefig(output_png)
    print(f"Plot saved to {output_png}")

if __name__ == "__main__":
    main()

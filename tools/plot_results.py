import pandas as pd
import matplotlib.pyplot as plt
import os

def main():
    possible_paths = [
        "results/game_analysis.csv",        # If run from root
        "../results/game_analysis.csv",     # If run from tools/
        "game_analysis.csv"                 # Fallback
    ]
    
    csv_path = None
    for p in possible_paths:
        if os.path.exists(p):
            csv_path = p
            break

    if csv_path is None:
        print(f"Error: game_analysis.csv not found. Run 'go run ./cmd/rps-context/main.go' first.")
        return

    try:
        df = pd.read_csv(csv_path)
    except FileNotFoundError:
        print(f"Error: game_analysis.csv not found at {csv_path}. Run 'go run ./cmd/rps-context/main.go' first.")
        return

    # Calculate Rolling Win Rate (window=50)
    df['Win'] = df['Reward'].apply(lambda x: 1 if x > 0 else 0)
    df['RollingWR'] = df['Win'].rolling(window=50).mean()

    plt.figure(figsize=(12, 6))
    
    # Plot Rolling Win Rate
    plt.plot(df['Step'], df['RollingWR'], label='Rolling Win Rate (50 steps)', color='blue')
    
    # Add Reference Lines
    plt.axhline(y=0.333, color='red', linestyle='--', label='Random Chance (33%)')
    plt.axhline(y=0.60, color='green', linestyle='--', label='Target (60%)')

    # Highlight Context Switches (Every 300 steps)
    for x in range(300, 1000, 300):
        plt.axvline(x=x, color='gray', linestyle=':', alpha=0.5)
        plt.text(x+10, 0.2, f'Switch', rotation=90, color='gray')

    plt.title('Contextual SSN Performance: Adaptation to Regime Switching')
    plt.xlabel('Game Step')
    plt.ylabel('Win Rate')
    plt.ylim(0, 1.0)
    plt.legend()
    plt.grid(True, alpha=0.3)
    
    output_file = "analysis_plot.png"
    plt.savefig(output_file)
    print(f"Plot saved to {output_file}")
    # plt.show() # Uncomment if running locally with UI

if __name__ == "__main__":
    main()

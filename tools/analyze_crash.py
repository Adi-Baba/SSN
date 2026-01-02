import pandas as pd
import numpy as np
import os

def analyze():
    print("--- Crash & Recovery Analysis ---")
    
    possible_paths = [
        "results/game_analysis.csv",
        "../results/game_analysis.csv",
        "game_analysis.csv"
    ]
    
    csv_path = None
    for p in possible_paths:
        if os.path.exists(p):
            csv_path = p
            break
            
    if csv_path is None:
        print("Error: game_analysis.csv not found.")
        return

    try:
        df = pd.read_csv(csv_path)
    except Exception as e:
        print(f"Error reading CSV: {e}")
        return

    df['Win'] = df['Reward'].apply(lambda x: 1 if x > 0 else 0)
    
    # Analyze around Switch Points (300, 600, 900)
    switch_points = [300, 600, 900]
    
    for switch in switch_points:
        print(f"\n[Regime Switch at Step {switch}]")
        
        # 1. State Before (Last 50 steps)
        pre_window = df[(df['Step'] > switch - 50) & (df['Step'] <= switch)]
        pre_wr = pre_window['Win'].mean()
        
        # 2. The Crash (Next 50 steps)
        post_window = df[(df['Step'] > switch) & (df['Step'] <= switch + 50)]
        post_wr = post_window['Win'].mean()
        
        drop = pre_wr - post_wr
        print(f"  WinRate Before: {pre_wr:.2%}")
        print(f"  WinRate After:  {post_wr:.2%}")
        print(f"  CRASH Magnitude: {drop:.2%} drop")
        
        # 3. Recovery Time
        # Find when we get back to 'Pre' levels (or at least > 50%)
        future = df[df['Step'] > switch]
        recovery_step = None
        
        # Simple rolling check for recovery
        for i in range(len(future) - 10):
            window = future.iloc[i:i+20] # Look at 20 step chunks
            if window['Win'].mean() >= 0.50:
                recovery_step = window.iloc[0]['Step']
                break
        
        if recovery_step:
            lag = recovery_step - switch
            print(f"  RECOVERY: Achieved >50% WR at Step {recovery_step} (Lag: {lag} steps)")
            print(f"  Interpretation: The network took {lag} generations of turnover to evolve the new strategy.")
        else:
            print("  RECOVERY: Failed to recover >50% before end of data/next switch.")

if __name__ == "__main__":
    analyze()

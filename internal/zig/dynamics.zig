const std = @import("std");

// --- SSN Evolutionary Dynamics (v1.3.0 Core) ---

// Winner Update: Linear Reinforcement
// Equation: E_new = E_old + Alpha * (Reward - MeanEnergy)
// Purpose: Boosts strategies that perform better than the population average.
pub fn update_winner(energy: f32, alpha: f32, reward: f32, mean_energy: f32) f32 {
    return energy + alpha * (reward - mean_energy);
}

// Loser Decay: Geometric Smoothing (The "Magic Equation")
// Equation: E_new = E_old * (1.0 - Beta)
// Purpose: Causes failed strategies to collapse exponentially, preventing
// "Zombie" paths from suffocating new mutants.
pub fn decay_loser(energy: f32, beta: f32) f32 {
    return energy * (1.0 - beta);
}

# SSN: Structural Selection Network
### Darwinian Intelligence for Online Adaptation

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Zig](https://img.shields.io/badge/Zig-0.11+-F7A41D?style=flat&logo=zig)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

**SSN** is a lightweight, high-performance machine learning engine built **from scratch** in Zig and Go. Unlike standard Neural Networks (Connectionist/Hebbian), SSN is a **Selectionist** (Darwinian) system. It does not use Backpropagation, Matrix Multiplication, or Weights. Instead, it uses **Geometric Energy Decay** and **Population Turnover** to adapt rapidly to dynamic environments.

🔗 **Repository**: [github.com/Adi-Baba/SSN](https://github.com/Adi-Baba/SSN.git)

---

## 🚀 Why SSN?

Traditional Neural Networks (LSTMs, Transformers) are "Slow Learners"—powerful on static datasets but sluggish in rapidly changing environments due to "Weight Inertia" and Catastrophic Forgetting.

**SSN is a "Fast Learner".**
*   **Instant Unlearning**: Uses Geometric Decay to flush old biases in **~15 steps**.
*   **Zero Dependencies**: No PyTorch, No TensorFlow, No BLAS. Just pure Zig & Go.
*   **Micro-Scale**: Runs efficiently on CPU (no GPU needed).
*   **Robust**: Proven to break even on pure noise (no hallucination) and exploit 50% signal.

---

## 🧠 Architecture

The system is a hybrid:
1.  **Core (Zig)**: Handles the high-frequency Evolutionary Dynamics (Selection, Mutation, Decay).
2.  **API (Go)**: Provides a clean, idiomatic Go interface for integration.

```mermaid
graph TD
    subgraph GoApp ["Go Application Layer"]
        Bot["Your Bot / App"]
    end

    subgraph Pkg ["SSN Package (pkg/ssn)"]
        API["Go API Wrapper"]
    end

    subgraph Zig ["Zig High-Performance Kernel"]
        C_ABI["C Export Layer"]
        Engine["Population Engine"]
        Math["Dynamics Core (dynamics.zig)"]
    end

    Bot -->|"1. Select Action"| API
    Bot -->|"2. Feed Reward"| API
    API <-->|"CGO (Fast)"| C_ABI
    C_ABI --> Engine
    Engine -->|"Geometric Decay"| Math
    Engine -->|"Survivor Selection"| Math
```

---

## ⚡ Benchmarks: The "Rock-Paper-Scissors" Test

We mandated a falsification test: **Online Rock-Paper-Scissors against a Regime-Switching Opponent.**
The opponent changes strategy (Bias: Rock -> Paper -> Scissors) every 300 steps.

| Metric | SSN (v1.2) | Standard NN (LSTM) |
| :--- | :--- | :--- |
| **Adaptation Lag** | **~15 Steps** | 100-500+ Steps |
| **Win Rate** | **60.7%** | ~50-55% (Slow convergence) |
| **Noise Robustness** | **Excellent** | Variable (Overfitting risk) |
| **Architecture** | **Selectionist** | Connectionist |

> **Verification**: In a 10,000-step stress test, SSN maintained a stable **56.8% Win Rate** with zero degradation.

---

## 📦 Installation
```bash
go get github.com/Adi-Baba/SSN
```

*   **Go 1.21+**
*   **Zig 0.11+** (for building the kernel)

## 🧪 Testing

To run the integration tests (which verify the Zig kernel and Go wrapper):

```bash
make test
# OR
build.bat
go test -v ./pkg/ssn
```

---

## 🛠 Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/Adi-Baba/SSN/pkg/ssn"
)

func main() {
    // 1. Configure
    cfg := ssn.DefaultConfig()
    cfg.PopSize = 100

    // 2. Initialize
    net, err := ssn.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer net.Close()

    // 3. Online Loop
    id := net.Select()
    net.Update(id, 1.0)
    
    // 4. Serialize (New Feature)
    state, err := net.Save()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Saved network state (%d bytes)\n", len(state))
}
```

---

## 🔬 Theoretical Difference

| Feature | Standard Neural Network | SSN (This Project) |
| :--- | :--- | :--- |
| **Paradigm** | Connectionist (Hebbian) | **Selectionist (Darwinian)** |
| **Core Math** | Calculus (Chain Rule) | **Arithmetic (Selection Rules)** |
| **Learning** | Weight Adjustment | **Population Turnover** |
| **Memory** | Distributed in Weights | **Discrete Paths** |

This independence makes SSN an ideal candidate for **Ensembles**—it fails differently than NNs.

---

## 📂 Project Structure

*   `pkg/ssn/`: Public Go API.
*   `internal/zig/`: The "Brain". Optimized Zig kernel.
*   `cmd/rps-context/`: The Reference Implementation (Contextual Bot).

## 📄 License
MIT License. Created by [Adi-Baba](https://github.com/Adi-Baba).

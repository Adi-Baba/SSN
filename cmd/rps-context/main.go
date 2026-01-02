package main

import (
    "encoding/csv"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "time"
    "github.com/Adi-Baba/SSN/pkg/ssn"
)

// Moves (Robust Encoding)
const (
    Rock     = 0
    Paper    = 1
    Scissors = -1
)

func getMoveName(move int) string {
    switch move {
    case Rock: return "Rock"
    case Paper: return "Paper"
    case Scissors: return "Scissors"
    default: return "Unknown"
    }
}

// Helper to map -1,0,1 to array index 0,1,2
func moveToIndex(m int) int {
    if m == -1 { return 2 }
    return m
}

func byteToMove(b byte) int {
    val := int(b) % 3
    if val == 2 { return Scissors }
    return val
}

const SwitchPeriod = 300

type Opponent struct {
    Step        int
    CurrentBias int 
}

func NewOpponent() *Opponent {
    return &Opponent{Step: 0, CurrentBias: Rock}
}

func (o *Opponent) GetMove() int {
    o.Step++
    if o.Step % SwitchPeriod == 0 {
        if o.CurrentBias == Rock { o.CurrentBias = Paper } else 
        if o.CurrentBias == Paper { o.CurrentBias = Scissors } else 
        { o.CurrentBias = Rock }
    }
    if rand.Float32() < 0.9 { return o.CurrentBias }
    r := rand.Intn(3)
    if r == 2 { return Scissors }
    return r
}

func solveWin(myMove, opMove int) float32 {
    if myMove == opMove { return 0.0 } // Tie = 0
    if (myMove == Rock && opMove == Scissors) ||
       (myMove == Paper && opMove == Rock) ||
       (myMove == Scissors && opMove == Paper) {
        return 1.0 // Win
    }
    return -0.5 // Loss
}

func main() {
    rand.Seed(time.Now().UnixNano())
    fmt.Printf("--- Contextual SSN (Ensemble) Experiment ---\n")
    
    // Log File
    // Log File Logic
    logDirs := []string{
        "../../results", // Standard project structure
        "results",       // If running from root
        ".",             // Fallback
    }
    
    logPath := "game_analysis.csv"
    for _, dir := range logDirs {
        if _, err := os.Stat(dir); err == nil {
            logPath = fmt.Sprintf("%s/game_analysis.csv", dir)
            break
        }
    }

    f, err := os.Create(logPath)
    if err != nil {
        fmt.Printf("Error creating log file: %v\n", err)
        return
    }
    defer f.Close()
    logger := csv.NewWriter(f)
    logger.Write([]string{"Step", "Context", "OpponentMove", "MyMove", "Reward", "WinRate"})
    defer logger.Flush()

    // Contextual Memory: 3 Independent Networks
    // Index 0: Context=Rock (Last move was Rock)
    // Index 1: Context=Paper
    // Index 2: Context=Scissors
    networks := make([]*ssn.Network, 3)
    
    cfg := ssn.DefaultConfig()
    cfg.PopSize = 100 // Smaller per context
    cfg.Alpha = 2.0
    cfg.Beta = 0.2
    cfg.Gamma = 0.1
    
    for i := 0; i < 3; i++ {
        cfg.Seed = uint64(time.Now().UnixNano() + int64(i*1000))
        net, err := ssn.New(cfg)
        if err != nil { panic(err) }
        networks[i] = net
        defer net.Close()
    }

    opp := NewOpponent()
    
    // Initial Context (Assume Rock to start)
    lastOppMove := Rock
    
    wins, losses, ties := 0, 0, 0
    
    for t := 1; t <= 1000; t++ {
        // 1. Determine Context
        ctxIdx := moveToIndex(lastOppMove)
        activeNet := networks[ctxIdx]
        
        // 2. Select (Conditional on Context)
        id := activeNet.Select()
        path := activeNet.GetPath(id)
        myMove := byteToMove(path[0])
        
        // 3. Opponent Moves
        opMove := opp.GetMove()
        
        // 4. Reward
        reward := solveWin(myMove, opMove)
        
        // 5. Update (The ACTIVE network learns)
        activeNet.Update(id, reward)
        
        // 6. Store Context
        lastOppMove = opMove

        // Stats
        if reward > 0 { wins++ } else if reward < 0 { losses++ } else { ties++ }
        
        wr := float64(wins) / float64(t)
        
        // Log
        logger.Write([]string{
            strconv.Itoa(t),
            getMoveName(lastOppMove), // Context for NEXT step, but technically 'Prev Move'
            getMoveName(opMove),
            getMoveName(myMove),
            fmt.Sprintf("%.1f", reward),
            fmt.Sprintf("%.3f", wr),
        })

        if t % 50 == 0 {
            fmt.Printf("Step %4d | Context: %-8s | WinRate: %.1f%%\n", 
                t, getMoveName(lastOppMove), wr*100)
        }
    }
    
    fmt.Printf("\nFinal Stats: Wins %d | Ties %d | Losses %d\n", wins, ties, losses)
    fmt.Printf("Analysis saved to %s\n", logPath)
}

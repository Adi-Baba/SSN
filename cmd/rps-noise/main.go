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

// Moves
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
    // 100% Random (Pure Chaos)
    r := rand.Intn(3)
    if r == 2 { return Scissors }
    return r
}

func solveWin(myMove, opMove int) float32 {
    if myMove == opMove { return 0.0 } 
    if (myMove == Rock && opMove == Scissors) ||
       (myMove == Paper && opMove == Rock) ||
       (myMove == Scissors && opMove == Paper) {
        return 1.0 
    }
    return -0.5 
}

func main() {
    rand.Seed(time.Now().UnixNano())
    fmt.Printf("--- Contextual SSN (Pure Random 100%%) ---\n")
    
    logDirs := []string{
        "../../results", 
        "results",       
        ".",             
    }
    
    logPath := "game_analysis_random.csv"
    for _, dir := range logDirs {
        if _, err := os.Stat(dir); err == nil {
            logPath = fmt.Sprintf("%s/game_analysis_random.csv", dir)
            break
        }
    }

    f, err := os.Create(logPath)
    if err != nil { panic(err) }
    defer f.Close()
    logger := csv.NewWriter(f)
    logger.Write([]string{"Step", "Context", "OpponentMove", "MyMove", "Reward", "WinRate"})
    defer logger.Flush()

    networks := make([]*ssn.Network, 3)
    
    cfg := ssn.DefaultConfig()
    cfg.PopSize = 100 
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
    lastOppMove := Rock
    wins, losses, ties := 0, 0, 0
    
    // 2000 steps is enough to see the chaos
    for t := 1; t <= 2000; t++ {
        ctxIdx := moveToIndex(lastOppMove)
        activeNet := networks[ctxIdx]
        
        id := activeNet.Select()
        path := activeNet.GetPath(id)
        myMove := byteToMove(path[0])
        opMove := opp.GetMove()
        reward := solveWin(myMove, opMove)
        activeNet.Update(id, reward)
        lastOppMove = opMove

        if reward > 0 { wins++ } else if reward < 0 { losses++ } else { ties++ }
        
        wr := float64(wins) / float64(t)
        
        logger.Write([]string{
            strconv.Itoa(t),
            getMoveName(lastOppMove),
            getMoveName(opMove),
            getMoveName(myMove),
            fmt.Sprintf("%.1f", reward),
            fmt.Sprintf("%.3f", wr),
        })

        if t % 100 == 0 {
            fmt.Printf("Step %4d | WinRate: %.1f%%\n", t, wr*100)
        }
    }
    
    fmt.Printf("\nFinal Stats: Wins %d | Ties %d | Losses %d\n", wins, ties, losses)
    fmt.Printf("Saved to %s\n", logPath)
}

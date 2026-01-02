package ssn_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Adi-Baba/SSN/pkg/ssn"
)

func init() {
	// Often helpful to seed Go's random for logic that uses it, 
	// though SSN has its own seed.
	rand.Seed(time.Now().UnixNano())
}

func TestLifecycle(t *testing.T) {
	cfg := ssn.DefaultConfig()
	net, err := ssn.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}
	defer net.Close()

	// 1. Select
	id := net.Select()
	if id >= uint32(cfg.PopSize) {
		t.Errorf("Selected ID %d out of bounds (PopSize %d)", id, cfg.PopSize)
	}

	// 2. GetPath
	path := net.GetPath(id)
	if len(path) != int(cfg.PathLen) {
		t.Errorf("Path length %d, expected %d", len(path), cfg.PathLen)
	}

	// 3. Update
	net.Update(id, 1.0)
}

func TestPersistence(t *testing.T) {
	cfg := ssn.DefaultConfig()
	cfg.Seed = 12345 // Deterministic seed
	net1, err := ssn.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create network 1: %v", err)
	}
	defer net1.Close()

	// Train net1 slightly
	for i := 0; i < 100; i++ {
		id := net1.Select()
		net1.Update(id, float32(i%2)) // Alternating reward
	}

	// Save
	state, err := net1.Save()
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}
	if len(state) == 0 {
		t.Fatalf("Saved state is empty")
	}

	// Create net2
	net2, err := ssn.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create network 2: %v", err)
	}
	defer net2.Close()

	// Verify Select differs initially (since net1 was trained)
	// Actually, with same seed, they start same. But net1 has evolved. 
	// Let's just load the state.
	err = net2.Load(state)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	// Now net2 should be in sync with net1 logic. 
	// We can't easily peek inside without more debug exports, 
	// but we can ensure it runs and potentially behaves similarly.
	id2 := net2.Select()
	if id2 >= uint32(cfg.PopSize) {
		t.Errorf("Net2 ID out of bounds")
	}
}

func BenchmarkSelect(b *testing.B) {
	cfg := ssn.DefaultConfig()
	net, err := ssn.New(cfg)
	if err != nil {
		b.Fatalf("Failed to create network: %v", err)
	}
	defer net.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		net.Select()
	}
}

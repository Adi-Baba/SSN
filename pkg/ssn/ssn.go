package ssn

import (
    "fmt"
    "math"
    "syscall"
    "unsafe"
)

// Version of the SSN package
const Version = "1.1.1"

var (
    kernel *syscall.LazyDLL
    procCreate *syscall.LazyProc
    procDestroy *syscall.LazyProc
    procSelect *syscall.LazyProc
    procUpdate *syscall.LazyProc
    procGet *syscall.LazyProc
)

// Config holds the hyperparameters for the Network.
// Layout must match Zig 'Config' extern struct.
type Config struct {
    PopSize int32
    PathLen int32
    Alpha   float32
    Beta    float32
    Gamma   float32
    Seed    uint64
}

// DefaultConfig returns safe defaults for the network.
func DefaultConfig() Config {
    return Config{
        PopSize: 100,
        PathLen: 10,
        Alpha:   0.5,
        Beta:    0.1,
        Gamma:   0.05,
        Seed:    0,
    }
}

// Network represents a Structural Selection Network instance.
type Network struct {
    handle uintptr // Pointer to Zig Population
    cfg    Config
}

// New creates and initializes a new SSN with the given configuration.
// It loads the 'ssn_kernel.dll' from the current working directory or system path.
func New(cfg Config) (*Network, error) {
    if kernel == nil {
        kernel = syscall.NewLazyDLL("ssn_kernel.dll")
        procCreate = kernel.NewProc("ssn_create")
        procDestroy = kernel.NewProc("ssn_destroy")
        procSelect = kernel.NewProc("ssn_select")
        procUpdate = kernel.NewProc("ssn_update")
        procGet = kernel.NewProc("ssn_get_path")
    }

    // Call Create(cfg*)
    // We pass a pointer to the Go struct. Since it's passed to C/Zig, we must ensure memory layout matches.
    // Go struct layout is not guaranteed to match C, but for simple types on x64 it usually works.
    // Ideally we'd use CGO types? No, pure Go syscall relies on standard ABI.
    r1, _, err := procCreate.Call(uintptr(unsafe.Pointer(&cfg)))
    if r1 == 0 {
        return nil, fmt.Errorf("failed to create network (returned null): %v", err)
    }

    return &Network{
        handle: r1,
        cfg:    cfg,
    }, nil
}

// Select asks the network to choose a path ID based on current energy states.
func (n *Network) Select() uint32 {
    r1, _, _ := procSelect.Call(n.handle)
    return uint32(r1)
}

// Update provides feedback to the network.
func (n *Network) Update(id uint32, reward float32) {
    rBits := math.Float32bits(reward)
    procUpdate.Call(n.handle, uintptr(id), uintptr(rBits))
}

// GetPath retrieves the decision sequence (bytes) for a given path ID.
func (n *Network) GetPath(id uint32) []byte {
    buf := make([]byte, n.cfg.PathLen)
    procGet.Call(n.handle, uintptr(id), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
    return buf
}

// Close cleans up the underlying kernel resources.
func (n *Network) Close() {
    if n.handle != 0 {
        procDestroy.Call(n.handle)
        n.handle = 0
    }
}

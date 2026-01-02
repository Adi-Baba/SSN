package ssn

import (
	"encoding/binary"
	"fmt"
	"math"
	"syscall"
	"unsafe"
)

// Version of the SSN package
const Version = "1.2.0"

var (
	kernel          *syscall.LazyDLL
	procCreate      *syscall.LazyProc
	procDestroy     *syscall.LazyProc
	procSelect      *syscall.LazyProc
	procUpdate      *syscall.LazyProc
	procGet         *syscall.LazyProc
	procGetStateSize *syscall.LazyProc
	procSaveState   *syscall.LazyProc
	procLoadState   *syscall.LazyProc
)

// Config holds the hyperparameters for the Network.
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
// It loads 'ssn_kernel.dll' from the current working directory or system path.
func New(cfg Config) (*Network, error) {
	if kernel == nil {
		kernel = syscall.NewLazyDLL("ssn_kernel.dll")
		procCreate = kernel.NewProc("ssn_create")
		procDestroy = kernel.NewProc("ssn_destroy")
		procSelect = kernel.NewProc("ssn_select")
		procUpdate = kernel.NewProc("ssn_update")
		procGet = kernel.NewProc("ssn_get_path")
		procGetStateSize = kernel.NewProc("ssn_get_state_size")
		procSaveState = kernel.NewProc("ssn_save_state")
		procLoadState = kernel.NewProc("ssn_load_state")
	}

	// Safety: Manually pack the Config struct to ensure C-ABI alignment
	// Zig Struct:
	// pop_size (4), path_len (4), alpha (4), beta (4), gamma (4), padding (4), seed (8)
	// Total: 32 bytes
	buf := make([]byte, 32)
	binary.LittleEndian.PutUint32(buf[0:], uint32(cfg.PopSize))
	binary.LittleEndian.PutUint32(buf[4:], uint32(cfg.PathLen))
	binary.LittleEndian.PutUint32(buf[8:], math.Float32bits(cfg.Alpha))
	binary.LittleEndian.PutUint32(buf[12:], math.Float32bits(cfg.Beta))
	binary.LittleEndian.PutUint32(buf[16:], math.Float32bits(cfg.Gamma))
	// Bytes 20-23 are padding for 8-byte alignment of seed
	binary.LittleEndian.PutUint64(buf[24:], cfg.Seed)

	r1, _, err := procCreate.Call(uintptr(unsafe.Pointer(&buf[0])))
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

// Save serializes the current network state to a byte slice.
func (n *Network) Save() ([]byte, error) {
	r1, _, _ := procGetStateSize.Call(n.handle)
	size := int(r1)
	if size == 0 {
		return nil, fmt.Errorf("failed to get state size")
	}

	buf := make([]byte, size)
	r2, _, _ := procSaveState.Call(n.handle, uintptr(unsafe.Pointer(&buf[0])), uintptr(size))
	if r2 == 0 { // false
		return nil, fmt.Errorf("failed to save state to buffer")
	}
	return buf, nil
}

// Load restores the network state from a byte slice.
// WARNING: The byte slice must match the structure of the version that saved it.
func (n *Network) Load(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data provided")
	}
	r1, _, _ := procLoadState.Call(n.handle, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	if r1 == 0 { // false
		return fmt.Errorf("failed to load state (corrupt data or size mismatch)")
	}
	return nil
}

// Close cleans up the underlying kernel resources.
func (n *Network) Close() {
	if n.handle != 0 {
		procDestroy.Call(n.handle)
		n.handle = 0
	}
}

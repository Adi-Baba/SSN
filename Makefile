.PHONY: all build clean test

# Default target
all: build

# Build the Zig kernel DLL
# Build the Zig kernel DLL
build:
	cd internal/zig && zig build-lib main.zig -dynamic -O ReleaseFast
	cmd /c move /Y internal\zig\main.dll ssn_kernel.dll
	cmd /c del internal\zig\main.lib
	cmd /c del internal\zig\main.pdb

# Run Go tests (requires DLL in path or same dir)
test: build
	cmd /c copy /Y ssn_kernel.dll pkg\ssn\ssn_kernel.dll
	go test -v ./pkg/ssn
	cmd /c del pkg\ssn\ssn_kernel.dll

# Clean up artifacts
clean:
	cmd /c del ssn_kernel.dll
	cmd /c del ssn_kernel.lib
	cmd /c del internal\zig\ssn_kernel.dll
	cmd /c del internal\zig\ssn_kernel.lib

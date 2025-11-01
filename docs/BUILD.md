# Build Guide - Hermes eUICC Manager

Comprehensive build instructions for all supported platforms and architectures.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Build Methods](#build-methods)
- [Platform-Specific Builds](#platform-specific-builds)
- [Optimization Levels](#optimization-levels)
- [Binary Size Optimization](#binary-size-optimization)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required

- **Go 1.21+**: Download from <https://go.dev/dl/>
- **Git**: For cloning repository
- **Linux**: Primary development platform

### Optional

- **UPX**: For binary compression (sudo apt install upx)
- **Make**: For using Makefile targets

### Verify Installation

```bash
go version    # Should show go1.21 or later
git --version
make --version  # Optional
upx --version   # Optional
```

## Quick Start

### Simple Build (Current Platform)

```bash
cd app
go build -o hermes-euicc .
./hermes-euicc --version
```

### Optimized Build (Current Platform)

```bash
cd app
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
```

Flags explained:

- `-ldflags="-s -w"`: Remove debug information and symbol table
- `-trimpath`: Remove file path information

## Build Methods

### Method 1: Using Makefile (Recommended)

From repository root:

```bash
# Build all platforms
make all

# Build specific platforms
make mipsle
make armv7
make arm64
make amd64

# Build OpenWRT popular platforms only
make openwrt

# Build + compress + checksum
make all compress checksum

# Clean build directory
make clean
```

### Method 2: Using build-all.sh

Generic platform builds from app/ directory:

```bash
cd app
./build-all.sh
```

Output directory: ../build/

Binaries created (8 total):

- hermes-euicc-mipsle (MIPS Little Endian)
- hermes-euicc-mips (MIPS Big Endian)
- hermes-euicc-armv5 (ARMv5)
- hermes-euicc-armv6 (ARMv6)
- hermes-euicc-armv7 (ARMv7)
- hermes-euicc-arm64 (ARM64)
- hermes-euicc-i386 (x86 32-bit)
- hermes-euicc-amd64 (x86 64-bit)

Features:

- Colored output
- Automatic UPX compression (if installed)
- SHA256 checksum generation
- Typical size: 6.5 MB → 2.8 MB (compressed)

### Method 3: Using build-openwrt.sh

Device-specific optimized builds from app/ directory:

```bash
cd app
./build-openwrt.sh
```

Output directory: ../build/openwrt/

Features:

- 41 device-specific binaries
- CPU-specific optimizations
- FPU detection for MIPS
- AVX2 support for modern x86
- Chipset information in build log

## Platform-Specific Builds

### MIPS Platforms

#### MIPS Little Endian (mipsle)

Most common for TP-Link, GL.iNet, Xiaomi routers:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-mipsle .
```

Devices: GL.iNet AR750S, TP-Link Archer C7/WR1043ND/WR841N, Xiaomi Mi Router 3G/4A, Ubiquiti EdgeRouter X

#### MIPS Big Endian (mips)

Older Broadcom-based routers:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-mips .
```

### ARM Platforms

#### ARMv7 (Recommended for most ARM routers)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-armv7 .
```

Features: VFPv3 floating point, NEON SIMD instructions, 30-40% faster than ARMv5/6

Devices: Raspberry Pi 2/3, GL.iNet B1300, Linksys WRT1900ACS, MikroTik hAP ac2

#### ARM64

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm64 .
```

Devices: Raspberry Pi 4/5, GL.iNet MT6000 (Flint 2), NanoPi R4S, Banana Pi BPI-R3/R4

### x86 Platforms

#### 64-bit (amd64)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-amd64 .
```

Devices: PC Engines APU, Protectli Vault, Qotom Mini PCs, Generic x86-64 servers

## Optimization Levels

### MIPS Optimizations

**GOMIPS=softfloat (Default)** - Safe for all MIPS devices

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build ...
```

**GOMIPS=hardfloat** - 20-30% faster on FPU-equipped devices

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=hardfloat go build ...
```

Use hardfloat for: QCA9563 (74Kc), MT7621 (1004Kc)
Use softfloat for: AR9341 (24Kc), QCA9531 (24Kc)

### x86 Optimizations

**GOAMD64 Levels:**

```bash
GOAMD64=v1  # Baseline x86-64 (all processors)
GOAMD64=v2  # +SSE4.2, POPCNT (2009+)
GOAMD64=v3  # +AVX2, BMI2 (2013+, 15-25% faster)
GOAMD64=v4  # +AVX512 (very new, rarely needed)
```

Examples:

```bash
# PC Engines APU (AMD GX-412TC, 2013)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v2 go build ...

# Qotom Q355G6 (Intel Celeron J3xxx)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build ...

# Generic (maximum compatibility)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v1 go build ...
```

## Binary Size Optimization

### Default Size

Without optimization: ~12 MB

### With Build Flags

```bash
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
```

Result: ~6.5 MB (45% reduction)

### With UPX Compression

```bash
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
upx --best --lzma hermes-euicc
```

Result: ~2.8 MB (57% reduction from 6.5 MB)

### Size Comparison

| Method | Size | Notes |
|--------|------|-------|
| Default | 12 MB | No optimization |
| -ldflags="-s -w" | 6.5 MB | Standard optimization |
| + UPX --best | 2.8 MB | Recommended |

## Troubleshooting

### "go: command not found"

Install Go 1.21+:

```bash
wget https://go.dev/dl/go1.23.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### "Illegal instruction" on Target Device

You used wrong optimization level. Rebuild with safe defaults (softfloat for MIPS, v1 for x86).

### Binary Too Large for Router

Use UPX compression:

```bash
upx --best --lzma hermes-euicc
```

### Build Script Permission Denied

Make scripts executable:

```bash
chmod +x build-all.sh build-openwrt.sh
```

## Verification

### Check Binary Architecture

```bash
file hermes-euicc-mipsle
# Output: ELF 32-bit LSB executable, MIPS, MIPS32 rel2 version 1...
```

### Verify SHA256

```bash
cd ../build
sha256sum -c SHA256SUMS
```

### Test Binary

```bash
./hermes-euicc --version

# On target device
scp hermes-euicc root@192.168.1.1:/tmp/
ssh root@192.168.1.1
chmod +x /tmp/hermes-euicc
/tmp/hermes-euicc --version
```

## Summary

**For quick development:**

```bash
go build -o hermes-euicc .
```

**For production (single platform):**

```bash
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
upx --best --lzma hermes-euicc
```

**For distribution (all platforms):**

```bash
cd app
./build-all.sh
```

**For device-specific (optimized):**

```bash
cd app
./build-openwrt.sh
```

**Using Makefile (recommended):**

```bash
make all compress checksum
```

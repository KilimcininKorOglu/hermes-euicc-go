# Build Guide - Hermes eUICC Manager

Comprehensive build instructions for 15 actively supported platforms and architectures.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Universal Build System](#universal-build-system)
- [Platform-Specific Builds](#platform-specific-builds)
- [Cross-Platform Builds](#cross-platform-builds)
- [Optimization Levels](#optimization-levels)
- [Binary Size Optimization](#binary-size-optimization)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required

- **Go 1.24+**: Download from <https://go.dev/dl/>
- **Git**: For cloning repository

### Verify Installation

```bash
go version    # Should show go1.24 or later
git --version
```

## Quick Start

### Simple Build (Current Platform)

```bash
cd app
go build -o hermes-euicc .
./hermes-euicc version
```

### Optimized Build (Current Platform)

```bash
cd app
go build -ldflags="-s -X main.Version=1.0.0 -X main.Release=1" -trimpath -o hermes-euicc .
```

Flags explained:

- `-ldflags="-s"`: Remove debug information and symbol table
- `-X main.Version=1.0.0`: Inject version number into binary
- `-X main.Release=1`: Inject release number into binary
- `-trimpath`: Remove file path information
- `CGO_ENABLED=0`: Static linking (no C dependencies)

## Universal Build System

### build-all.sh - One Script for All Platforms

The `build-all.sh` script builds for **15 active platforms** with a single command:

```bash
cd app
./build-all.sh
```

**Features:**
- Auto-installs Go 1.24.0 if not found
- **Auto-increments release numbers** using git commit count
- **Injects version info** into binaries via ldflags
- Builds 15 platforms (Linux, OpenWRT, Windows) - macOS/FreeBSD disabled due to upstream issues
- Platform-specific optimizations (GOMIPS=softfloat, GOAMD64=v2)
- **Version-specific output directory** (build/1.0.0/)
- **Versioned binary filenames** (hermes-euicc-1.0.0-48-linux-amd64)
- **IPK package generation** for all 10 OpenWRT architectures
- SHA256SUMS generation
- Colored output with progress indication
- Continues on errors (doesn't stop the entire build)

**Output structure:**
```
build/
└── 1.0.0/              # Version-specific directory
    ├── hermes-euicc-1.0.0-48-linux-amd64
    ├── hermes-euicc-1.0.0-48-linux-i386
    ├── hermes-euicc-1.0.0-48-linux-arm64
    ├── hermes-euicc-1.0.0-48-openwrt-mips
    ├── hermes-euicc-1.0.0-48-openwrt-mipsle
    ├── hermes-euicc-1.0.0-48-openwrt-mips64
    ├── hermes-euicc-1.0.0-48-openwrt-mips64le
    ├── hermes-euicc-1.0.0-48-openwrt-arm_v5
    ├── hermes-euicc-1.0.0-48-openwrt-arm_v6
    ├── hermes-euicc-1.0.0-48-openwrt-arm_v7
    ├── hermes-euicc-1.0.0-48-openwrt-arm64
    ├── hermes-euicc-1.0.0-48-openwrt-x86
    ├── hermes-euicc-1.0.0-48-openwrt-x86_64
    ├── hermes-euicc-1.0.0-48-windows-amd64.exe
    ├── hermes-euicc-1.0.0-48-windows-i386.exe
    ├── hermes-euicc-1.0.0-48-windows-arm64.exe
    ├── hermes-euicc_1.0.0-48_mips.ipk
    ├── hermes-euicc_1.0.0-48_mipsle.ipk
    ├── hermes-euicc_1.0.0-48_mips64.ipk
    ├── hermes-euicc_1.0.0-48_mips64le.ipk
    ├── hermes-euicc_1.0.0-48_arm_cortex-a5.ipk
    ├── hermes-euicc_1.0.0-48_arm_cortex-a7.ipk
    ├── hermes-euicc_1.0.0-48_arm_cortex-a9.ipk
    ├── hermes-euicc_1.0.0-48_aarch64.ipk
    ├── hermes-euicc_1.0.0-48_x86.ipk
    ├── hermes-euicc_1.0.0-48_x86_64.ipk
    └── SHA256SUMS
```

**Total: 16 binaries + 10 IPK packages**
- Linux: 3 (amd64, i386, arm64)
- OpenWRT: 10 (MIPS, ARM variants) + 10 IPK packages
- Windows: 3 (x64, x86, ARM64)
- macOS: Disabled (upstream library missing drivers)
- FreeBSD: Disabled (upstream library missing drivers)

**Version numbering:**
- Format: `{PKG_VERSION}-{PKG_RELEASE}` (e.g., `1.0.0-48`)
- PKG_VERSION: Manually set in build-all.sh (currently 1.0.0)
- PKG_RELEASE: Auto-incremented via `git rev-list --count HEAD`
- Each new git commit automatically increments the release number

**Build script optimizations:**
- Version injection via `-ldflags="-X main.Version=... -X main.Release=..."`
- OpenWRT builds use `-tags=openwrt` for UCI config support
- MIPS uses `GOMIPS=softfloat` for FPU-less routers
- MIPS64 uses `GOMIPS64=softfloat` for FPU-less 64-bit routers
- x86-64 uses `GOAMD64=v2` for SSE4.2 support (2009+ CPUs)
- CCID driver disabled on MIPS/32-bit Linux (purego v0.9.0 limitation)

## Platform-Specific Builds

### Linux Desktop/Server

#### x86-64 (amd64) - Modern PCs and Servers

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v2 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-amd64 .
```

**Devices:** Desktop PCs, servers, PC Engines APU, Protectli Vault

#### x86 32-bit (i386) - Legacy PCs

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-i386 .
```

#### ARM64 - Raspberry Pi 4+, Servers

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm64 .
```

**Devices:** Raspberry Pi 4/5, ARM servers

### OpenWRT/Embedded Linux (Routers)

**Note:** OpenWRT builds require `-tags=openwrt` flag to enable UCI config support.

#### MIPS Big Endian - Atheros AR/QCA

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-mips .
```

**Devices:** TP-Link Archer series, GL.iNet AR/XE series, Ubiquiti routers

#### MIPS Little Endian - MediaTek MT76xx

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-mipsle .
```

**Devices:** GL.iNet MT series, Xiaomi routers, Ralink-based devices

#### MIPS64 BE/LE - Cavium Octeon

```bash
# Big Endian
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 GOMIPS64=softfloat go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-mips64 .

# Little Endian
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le GOMIPS64=softfloat go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-mips64le .
```

**Devices:** Ubiquiti EdgeRouter, Cavium-based systems

#### ARM v5/v6/v7 - Various ARM Routers

```bash
# ARMv5 (Kirkwood, old NAS)
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm_v5 .

# ARMv6 (Raspberry Pi Zero/1)
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm_v6 .

# ARMv7 (Most ARM routers)
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm_v7 .
```

**ARMv7 Devices:** GL.iNet B1300, Linksys WRT series, Raspberry Pi 2/3, IPQ40xx routers

#### ARM64 - Modern ARM Routers

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm64 .
```

**Devices:** GL.iNet MT6000, BananaPi R3/R4, IPQ807x, MT7622/MT7986 routers

#### x86/x86-64 - PC-based Routers

```bash
# x86 32-bit
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-x86 .

# x86-64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v2 go build \
    -tags=openwrt -ldflags="-s -w" -trimpath \
    -o hermes-euicc-x86_64 .
```

**Devices:** PC Engines APU, Protectli Vault, x86 routers, VMs

### macOS

#### Intel (x86-64)

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GOAMD64=v2 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-amd64 .
```

**Devices:** Intel-based Macs (2006-2020)

#### Apple Silicon (ARM64)

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm64 .
```

**Devices:** M1/M2/M3/M4 Macs (2020+)

### Windows

#### x86-64 (64-bit)

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GOAMD64=v2 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-amd64.exe .
```

**Devices:** Modern Windows PCs (64-bit)

#### x86 (32-bit)

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-i386.exe .
```

**Devices:** Legacy Windows PCs (32-bit)

#### ARM64

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm64.exe .
```

**Devices:** Surface Pro X, ARM-based Windows laptops

### FreeBSD

#### x86-64 (amd64)

```bash
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 GOAMD64=v2 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-amd64 .
```

#### ARM64

```bash
CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm64 .
```

## Cross-Platform Builds

You can build for any platform from any platform using Go's cross-compilation:

```bash
# From Linux, build for macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o hermes-euicc-macos-arm64 .

# From macOS, build for Windows x64
GOOS=windows GOARCH=amd64 go build -o hermes-euicc-win64.exe .

# From Windows, build for Linux ARM64
set GOOS=linux
set GOARCH=arm64
go build -o hermes-euicc-linux-arm64 .
```

## Optimization Levels

### MIPS Optimizations

**GOMIPS=softfloat (Default)** - Safe for all MIPS devices

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build ...
```

Compatible with all MIPS routers (24Kc, 74Kc, 1004Kc cores).

**GOMIPS=hardfloat** - 20-30% faster on FPU-equipped devices

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=hardfloat go build ...
```

Use hardfloat for:
- QCA9563 (74Kc core with FPU)
- MT7621 (1004Kc core with FPU)

Use softfloat for:
- AR9341, QCA9531 (24Kc core without FPU)
- Most TP-Link, GL.iNet routers

### x86 Optimizations

**GOAMD64 Levels:**

```bash
GOAMD64=v1  # Baseline x86-64 (all processors, maximum compatibility)
GOAMD64=v2  # +SSE4.2, POPCNT (2009+, recommended)
GOAMD64=v3  # +AVX2, BMI2 (2013+, 15-25% faster)
GOAMD64=v4  # +AVX512 (very new, rarely needed)
```

**Recommendations:**

```bash
# For maximum compatibility (v1)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v1 go build ...

# For modern systems (v2, recommended)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v2 go build ...

# For latest CPUs (v3)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build ...
```

**CPU compatibility:**
- v1: All x86-64 CPUs (2003+)
- v2: Intel Nehalem (2009+), AMD Bulldozer (2011+)
- v3: Intel Haswell (2013+), AMD Excavator (2015+)

## Binary Size Optimization

### Default Size

Without optimization: ~12 MB

### With Build Flags

```bash
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
```

Result: ~6.4 MB (47% reduction)

### Size Comparison

| Method | Size | Reduction | Notes |
|--------|------|-----------|-------|
| Default build | 12 MB | - | No optimization |
| `-ldflags="-s -w"` | 6.4 MB | 47% | Standard optimization |
| `-trimpath` (included) | 6.4 MB | - | Removes build paths |

**Current binary sizes (from build-all.sh):**
- Linux amd64: 6.4M
- OpenWRT MIPS: 7.0M
- macOS ARM64: 6.1M
- Windows x64: 6.5M
- FreeBSD amd64: 6.1M

## Build Tags

### Build Tags and Configuration System

The application uses Go build tags to provide platform-specific configuration systems:

#### OpenWRT Builds (UCI Configuration)

OpenWRT builds require the `openwrt` build tag:

```bash
go build -tags=openwrt -o hermes-euicc .
```

**Files compiled:**
- `uci_openwrt.go` - Reads from `/etc/config/hermes-euicc`
- Excludes `config.go` and non-OpenWRT `uci_other.go`

**Features:**
- UCI configuration support
- LuCI web interface integration
- Automatic configuration from OpenWRT system

#### Non-OpenWRT Builds (Config File)

Standard builds (Linux, macOS, Windows, FreeBSD) don't need tags:

```bash
go build -o hermes-euicc .
```

**Files compiled:**
- `config.go` - Config file parser (key=value format)
- `uci_other.go` - Config file loader with auto-detection

**Features:**
- Config file support (./hermes-euicc.conf, ~/.config/hermes-euicc/config, etc.)
- `-config` flag for custom config file path
- Platform-specific config directory support

## Troubleshooting

### "go: command not found"

Install Go 1.24+:

```bash
# Linux/macOS
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Or let build-all.sh auto-install it
./build-all.sh
```

### "Illegal instruction" on Target Device

Wrong optimization level. Causes:

1. **MIPS hardfloat on softfloat device**
   - Solution: Use `GOMIPS=softfloat`

2. **GOAMD64=v3 on old CPU**
   - Solution: Use `GOAMD64=v2` or `GOAMD64=v1`

3. **Wrong GOARM level**
   - Solution: Use ARMv5 binary for maximum compatibility

### Build Fails on Windows

Use PowerShell or Git Bash, not CMD:

```powershell
# PowerShell
$env:CGO_ENABLED="0"
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -ldflags="-s -w" -trimpath -o hermes-euicc.exe .
```

### Cross-Compilation Issues

Ensure `CGO_ENABLED=0` for cross-compilation:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ...
```

## Verification

### Check Binary Architecture

```bash
# Linux
file hermes-euicc-mipsle
# Output: ELF 32-bit LSB executable, MIPS, MIPS32 rel2 version 1...

# macOS
file hermes-euicc-amd64
# Output: Mach-O 64-bit executable x86_64

# Windows (use Git Bash)
file hermes-euicc-amd64.exe
# Output: PE32+ executable (console) x86-64
```

### Verify SHA256

```bash
cd build/linux
sha256sum -c SHA256SUMS

# Or verify master file
cd build
sha256sum -c SHA256SUMS.txt
```

### Test Binary

```bash
# Local test
./hermes-euicc version

# Remote device test (OpenWRT example)
scp build/openwrt/hermes-euicc-mipsle root@192.168.1.1:/tmp/
ssh root@192.168.1.1
chmod +x /tmp/hermes-euicc-mipsle
/tmp/hermes-euicc-mipsle version
```

## Platform Selection Guide

**For OpenWRT/Embedded Routers:**

Check architecture on device:

```bash
ls -la /lib/ld-musl-*.so.1
```

Results:
- `ld-musl-mips-sf.so.1` → Use `hermes-euicc-mips`
- `ld-musl-mipsel-sf.so.1` → Use `hermes-euicc-mipsle`
- `ld-musl-armhf.so.1` → Use `hermes-euicc-arm_v7`
- `ld-musl-aarch64.so.1` → Use `hermes-euicc-arm64`
- `ld-musl-x86_64.so.1` → Use `hermes-euicc-x86_64`

**For Desktop/Server:**
- Linux x86-64 → `linux/hermes-euicc-amd64`
- macOS Intel → `darwin/hermes-euicc-amd64`
- macOS Apple Silicon → `darwin/hermes-euicc-arm64`
- Windows 64-bit → `windows/hermes-euicc-amd64.exe`

## Summary

**Quick development:**

```bash
go build -o hermes-euicc .
```

**Production (single platform):**

```bash
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
```

**All platforms (recommended):**

```bash
./build-all.sh
```

**OpenWRT with UCI support:**

```bash
go build -tags=openwrt -ldflags="-s -w" -trimpath -o hermes-euicc .
```

**Cross-platform example:**

```bash
# From any OS, build for Raspberry Pi
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o hermes-euicc-rpi4 .
```

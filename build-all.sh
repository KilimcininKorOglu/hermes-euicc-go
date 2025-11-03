#!/bin/bash

# Hermes eUICC Manager - Universal Multi-Platform Build Script
# Builds for: Linux (Desktop/Server/Embedded), macOS, Windows, FreeBSD

set -e  # Exit on error

# Colored output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Banner
echo -e "${BLUE}"
cat << 'EOF'
╔════════════════════════════════════════════════════════════╗
║     Hermes eUICC Manager - Universal Build System          ║
║   Desktop • Server • Embedded • macOS • Windows • BSD      ║
╚════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}\n"

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_DIR="${SCRIPT_DIR}"
BUILD_DIR="${SCRIPT_DIR}/build"
BINARY_NAME="hermes-euicc"
GO_VERSION="1.24.0"  # Required Go version

# Function to check and install Go if needed
check_go_installation() {
    if command -v go &> /dev/null; then
        CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        echo -e "${GREEN}✓ Go found:${NC} version ${CURRENT_GO_VERSION}"

        # Check if version is sufficient (basic check)
        REQUIRED_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
        CURRENT_MAJOR=$(echo $CURRENT_GO_VERSION | cut -d. -f1)
        REQUIRED_MINOR=$(echo $GO_VERSION | cut -d. -f2)
        CURRENT_MINOR=$(echo $CURRENT_GO_VERSION | cut -d. -f2)

        if [ "$CURRENT_MAJOR" -lt "$REQUIRED_MAJOR" ] || \
           ([ "$CURRENT_MAJOR" -eq "$REQUIRED_MAJOR" ] && [ "$CURRENT_MINOR" -lt "$REQUIRED_MINOR" ]); then
            echo -e "${YELLOW}⚠ Warning: Go ${GO_VERSION} or higher recommended (current: ${CURRENT_GO_VERSION})${NC}"
        fi
        return 0
    fi

    echo -e "${YELLOW}⚠ Go not found. Installing Go ${GO_VERSION}...${NC}\n"
    install_go
}

# Function to install Go
install_go() {
    local OS_TYPE=$(uname -s | tr '[:upper:]' '[:lower:]')
    local ARCH=$(uname -m)

    # Map architecture names
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        armv7l|armv6l) ARCH="armv6l" ;;
        i386|i686) ARCH="386" ;;
        *) echo -e "${RED}✗ Unsupported architecture: $ARCH${NC}"; exit 1 ;;
    esac

    local GO_PKG="go${GO_VERSION}.${OS_TYPE}-${ARCH}.tar.gz"
    local GO_URL="https://go.dev/dl/${GO_PKG}"
    local INSTALL_DIR="/usr/local"

    echo -e "${BLUE}Downloading Go ${GO_VERSION} for ${OS_TYPE}-${ARCH}...${NC}"

    # Download Go
    if command -v curl &> /dev/null; then
        curl -L -o "/tmp/${GO_PKG}" "$GO_URL" || {
            echo -e "${RED}✗ Failed to download Go${NC}"
            exit 1
        }
    elif command -v wget &> /dev/null; then
        wget -O "/tmp/${GO_PKG}" "$GO_URL" || {
            echo -e "${RED}✗ Failed to download Go${NC}"
            exit 1
        }
    else
        echo -e "${RED}✗ Neither curl nor wget found. Please install one of them.${NC}"
        exit 1
    fi

    echo -e "${BLUE}Installing Go to ${INSTALL_DIR}...${NC}"

    # Remove old Go installation
    if [ -d "${INSTALL_DIR}/go" ]; then
        echo -e "${YELLOW}Removing old Go installation...${NC}"
        sudo rm -rf "${INSTALL_DIR}/go"
    fi

    # Extract Go
    sudo tar -C "${INSTALL_DIR}" -xzf "/tmp/${GO_PKG}" || {
        echo -e "${RED}✗ Failed to extract Go${NC}"
        exit 1
    }

    # Clean up
    rm "/tmp/${GO_PKG}"

    # Add to PATH if not already there
    export PATH="${INSTALL_DIR}/go/bin:$PATH"

    # Check if PATH is in shell profile
    local PROFILE_FILE=""
    if [ -f "$HOME/.bashrc" ]; then
        PROFILE_FILE="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        PROFILE_FILE="$HOME/.bash_profile"
    elif [ -f "$HOME/.zshrc" ]; then
        PROFILE_FILE="$HOME/.zshrc"
    fi

    if [ -n "$PROFILE_FILE" ] && ! grep -q "${INSTALL_DIR}/go/bin" "$PROFILE_FILE"; then
        echo -e "${YELLOW}Adding Go to PATH in ${PROFILE_FILE}...${NC}"
        echo "" >> "$PROFILE_FILE"
        echo "# Go Programming Language" >> "$PROFILE_FILE"
        echo "export PATH=\"${INSTALL_DIR}/go/bin:\$PATH\"" >> "$PROFILE_FILE"
    fi

    echo -e "${GREEN}✓ Go ${GO_VERSION} installed successfully!${NC}"
    echo -e "${BLUE}Note: You may need to run 'source ${PROFILE_FILE}' or restart your terminal${NC}\n"
}

# Check Go installation
check_go_installation

# Create build directories
mkdir -p ${BUILD_DIR}/{linux,openwrt,darwin,windows,freebsd}

# Enhanced build function with optimization support
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local GOARM=$3
    local OUTPUT_DIR=$4
    local OUTPUT_NAME=$5
    local DESCRIPTION=$6
    local OPT_FLAGS=$7  # Optimization flags like GOMIPS=softfloat
    local BUILD_TAGS=$8  # Build tags like "openwrt"

    local OUTPUT_PATH="${BUILD_DIR}/${OUTPUT_DIR}/${OUTPUT_NAME}"

    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${MAGENTA}📦 Platform:${NC}  ${DESCRIPTION}"
    echo -e "${BLUE}🎯 Target:${NC}    ${GOOS}/${GOARCH}$([ -n "$GOARM" ] && echo "/v${GOARM}")"
    echo -e "${YELLOW}💾 Output:${NC}    ${OUTPUT_DIR}/${OUTPUT_NAME}"

    # Show optimization info
    if [ -n "$OPT_FLAGS" ]; then
        echo -e "${GREEN}⚡ Optimize:${NC}  ${OPT_FLAGS}"
        export $OPT_FLAGS
    fi

    # Build command
    local BUILD_CMD="CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH"
    [ -n "$GOARM" ] && BUILD_CMD="$BUILD_CMD GOARM=$GOARM"

    # Add build tags if specified
    local TAGS_FLAG=""
    [ -n "$BUILD_TAGS" ] && TAGS_FLAG="-tags=$BUILD_TAGS"

    if (cd ${SOURCE_DIR} && eval $BUILD_CMD go build \
        $TAGS_FLAG \
        -ldflags=-s \
        -trimpath \
        -o "${OUTPUT_PATH}" \
        . 2>&1 | grep -v "^#" || true); then

        # Unset optimization flags
        if [ -n "$OPT_FLAGS" ]; then
            unset $(echo "$OPT_FLAGS" | cut -d= -f1)
        fi

        if [ -f "${OUTPUT_PATH}" ]; then
            SIZE=$(ls -lh "${OUTPUT_PATH}" | awk '{print $5}')
            echo -e "${GREEN}✓ Success!${NC} Size: ${SIZE}"

            # Compress with UPX (if available and beneficial)
            if command -v upx &> /dev/null; then
                FILE_SIZE=$(stat -f%z "${OUTPUT_PATH}" 2>/dev/null || stat -c%s "${OUTPUT_PATH}" 2>/dev/null)
                if [ "$FILE_SIZE" -gt 500000 ]; then  # Only compress if > 500KB
                    echo -e "${YELLOW}  ⚡ Compressing with UPX...${NC}"
                    upx --best --lzma "${OUTPUT_PATH}" >/dev/null 2>&1 && {
                        SIZE_AFTER=$(ls -lh "${OUTPUT_PATH}" | awk '{print $5}')
                        echo -e "${GREEN}  ✓ Compressed!${NC} New size: ${SIZE_AFTER}"
                    } || {
                        echo -e "${YELLOW}  ⚠ UPX compression skipped${NC}"
                    }
                fi
            fi
            echo ""
            return 0
        fi
    fi

    # Unset optimization flags on failure too
    if [ -n "$OPT_FLAGS" ]; then
        unset $(echo "$OPT_FLAGS" | cut -d= -f1)
    fi
    echo -e "${RED}✗ Build failed!${NC}"
    echo -e "${YELLOW}   (Skipping - may not be supported on this platform)${NC}\n"
    return 0  # Continue with other builds
}

# =============================================================================
# Linux Desktop/Server Builds
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              Linux Desktop/Server Platforms                ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform "linux" "amd64" "" "linux" "${BINARY_NAME}-amd64" \
    "Linux x86-64 (Modern PCs, Servers)" \
    "GOAMD64=v2"

build_platform "linux" "386" "" "linux" "${BINARY_NAME}-i386" \
    "Linux x86 32-bit (Legacy PCs)" \
    ""

build_platform "linux" "arm64" "" "linux" "${BINARY_NAME}-arm64" \
    "Linux ARM 64-bit (Raspberry Pi 4+, Servers)" \
    ""

# =============================================================================
# OpenWRT/Embedded Linux Builds (Optimized for routers)
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║           OpenWRT/Embedded Linux (Routers/IoT)             ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

# MIPS Platforms (with openwrt build tag for UCI support)
build_platform "linux" "mips" "" "openwrt" "${BINARY_NAME}-mips" \
    "MIPS BE (Atheros AR/QCA, TP-Link, GL.iNet AR/XE, Ubiquiti)" \
    "GOMIPS=softfloat" "openwrt"

build_platform "linux" "mipsle" "" "openwrt" "${BINARY_NAME}-mipsle" \
    "MIPS LE (MediaTek MT76xx, GL.iNet MT, Ralink)" \
    "GOMIPS=softfloat" "openwrt"

build_platform "linux" "mips64" "" "openwrt" "${BINARY_NAME}-mips64" \
    "MIPS64 BE (Cavium Octeon, EdgeRouter Pro)" \
    "GOMIPS64=softfloat" "openwrt"

build_platform "linux" "mips64le" "" "openwrt" "${BINARY_NAME}-mips64le" \
    "MIPS64 LE (Cavium Octeon Little-Endian)" \
    "GOMIPS64=softfloat" "openwrt"

# ARM Embedded Platforms (with openwrt build tag for UCI support)
build_platform "linux" "arm" "5" "openwrt" "${BINARY_NAME}-arm_v5" \
    "ARM v5 (Kirkwood, Old NAS devices)" \
    "" "openwrt"

build_platform "linux" "arm" "6" "openwrt" "${BINARY_NAME}-arm_v6" \
    "ARM v6 (Raspberry Pi Zero/1, BCM2835)" \
    "" "openwrt"

build_platform "linux" "arm" "7" "openwrt" "${BINARY_NAME}-arm_v7" \
    "ARM v7 (IPQ40xx, GL.iNet B1300, Raspberry Pi 2/3)" \
    "" "openwrt"

build_platform "linux" "arm64" "" "openwrt" "${BINARY_NAME}-arm64" \
    "ARM64 (MT7622/MT7986, IPQ807x, BananaPi R3/R4, GL.iNet MT6000)" \
    "" "openwrt"

# x86 Embedded Platforms (with openwrt build tag for UCI support)
build_platform "linux" "386" "" "openwrt" "${BINARY_NAME}-x86" \
    "x86 32-bit (Legacy PC Engines, Old x86 routers)" \
    "" "openwrt"

build_platform "linux" "amd64" "" "openwrt" "${BINARY_NAME}-x86_64" \
    "x86-64 (PC Engines APU, Protectli, x86 routers, VMs)" \
    "GOAMD64=v2" "openwrt"

# =============================================================================
# macOS Builds
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                     macOS Platforms                        ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform "darwin" "amd64" "" "darwin" "${BINARY_NAME}-amd64" \
    "macOS Intel (x86-64)" \
    "GOAMD64=v2"

build_platform "darwin" "arm64" "" "darwin" "${BINARY_NAME}-arm64" \
    "macOS Apple Silicon (M1/M2/M3/M4)" \
    ""

# =============================================================================
# Windows Builds
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    Windows Platforms                       ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform "windows" "amd64" "" "windows" "${BINARY_NAME}-amd64.exe" \
    "Windows x86-64 (64-bit)" \
    "GOAMD64=v2"

build_platform "windows" "386" "" "windows" "${BINARY_NAME}-i386.exe" \
    "Windows x86 (32-bit)" \
    ""

build_platform "windows" "arm64" "" "windows" "${BINARY_NAME}-arm64.exe" \
    "Windows ARM64 (Surface Pro X, ARM laptops)" \
    ""

# =============================================================================
# FreeBSD Builds
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    FreeBSD Platforms                       ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform "freebsd" "amd64" "" "freebsd" "${BINARY_NAME}-amd64" \
    "FreeBSD x86-64" \
    "GOAMD64=v2"

build_platform "freebsd" "arm64" "" "freebsd" "${BINARY_NAME}-arm64" \
    "FreeBSD ARM 64-bit" \
    ""

# =============================================================================
# Generate Checksums
# =============================================================================
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Generating checksums...${NC}"

# Generate SHA256 checksums for each directory
for dir in linux openwrt darwin windows freebsd; do
    if [ -d "${BUILD_DIR}/${dir}" ] && [ "$(ls -A ${BUILD_DIR}/${dir} 2>/dev/null)" ]; then
        (cd ${BUILD_DIR}/${dir} && sha256sum ${BINARY_NAME}* > SHA256SUMS 2>/dev/null) && \
            echo -e "${GREEN}✓ ${dir}/SHA256SUMS created${NC}"
    fi
done

# Generate master checksum file
find ${BUILD_DIR} -name "${BINARY_NAME}*" -type f ! -name "SHA256SUMS" -exec sha256sum {} \; > ${BUILD_DIR}/SHA256SUMS.txt 2>/dev/null && \
    echo -e "${GREEN}✓ Master SHA256SUMS.txt created${NC}"

echo ""

# =============================================================================
# Build Summary
# =============================================================================
echo -e "${BLUE}"
cat << 'EOF'
╔════════════════════════════════════════════════════════════╗
║                   Build Complete!                          ║
╚════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

echo -e "${CYAN}Build Statistics:${NC}"
echo -e "  Output directory:    ${BUILD_DIR}"

for dir in linux openwrt darwin windows freebsd; do
    if [ -d "${BUILD_DIR}/${dir}" ]; then
        COUNT=$(ls -1 ${BUILD_DIR}/${dir}/${BINARY_NAME}* 2>/dev/null | grep -v SHA256SUMS | wc -l)
        [ $COUNT -gt 0 ] && echo -e "  ${dir}/ binaries:     ${COUNT}"
    fi
done

TOTAL=$(find ${BUILD_DIR} -name "${BINARY_NAME}*" -type f ! -name "SHA256SUMS*" | wc -l)
echo -e "  ${GREEN}Total binaries:      ${TOTAL}${NC}"
echo ""

# =============================================================================
# Usage Guide
# =============================================================================
echo -e "${YELLOW}Platform Selection Guide:${NC}"
echo ""
echo -e "${CYAN}For OpenWRT/Embedded Routers:${NC}"
echo -e "  Check architecture: ${MAGENTA}ls -la /lib/ld-musl-*.so.1${NC}"
echo -e "    ${GREEN}ld-musl-mips-sf.so.1${NC}     → Use: ${BLUE}openwrt/${BINARY_NAME}-mips${NC}"
echo -e "    ${GREEN}ld-musl-mipsel-sf.so.1${NC}   → Use: ${BLUE}openwrt/${BINARY_NAME}-mipsle${NC}"
echo -e "    ${GREEN}ld-musl-armhf.so.1${NC}       → Use: ${BLUE}openwrt/${BINARY_NAME}-arm_v7${NC}"
echo -e "    ${GREEN}ld-musl-aarch64.so.1${NC}     → Use: ${BLUE}openwrt/${BINARY_NAME}-arm64${NC}"
echo -e "    ${GREEN}ld-musl-x86_64.so.1${NC}      → Use: ${BLUE}openwrt/${BINARY_NAME}-x86_64${NC}"
echo ""
echo -e "${CYAN}For Desktop/Server:${NC}"
echo -e "  ${GREEN}Linux x86-64${NC}              → Use: ${BLUE}linux/${BINARY_NAME}-amd64${NC}"
echo -e "  ${GREEN}macOS Intel${NC}               → Use: ${BLUE}darwin/${BINARY_NAME}-amd64${NC}"
echo -e "  ${GREEN}macOS Apple Silicon${NC}       → Use: ${BLUE}darwin/${BINARY_NAME}-arm64${NC}"
echo -e "  ${GREEN}Windows 64-bit${NC}            → Use: ${BLUE}windows/${BINARY_NAME}-amd64.exe${NC}"
echo ""

echo -e "${GREEN}✓ All builds completed successfully!${NC}\n"

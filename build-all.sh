#!/bin/bash

# Hermes eUICC Manager - Multi-Platform Build Script
# Includes all drivers (QMI, MBIM, AT, CCID)

set -e  # Exit on error

# Colored output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Banner
echo -e "${BLUE}"
echo "╔════════════════════════════════════════════════════════════╗"
echo "║       Hermes eUICC Manager - Multi-Platform Build         ║"
echo "║              All Drivers Included (QMI/MBIM/AT/CCID)      ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo -e "${NC}\n"

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_DIR="${SCRIPT_DIR}"
BUILD_DIR="${SCRIPT_DIR}/build"
BINARY_NAME="hermes-euicc"

# Create build directory
mkdir -p ${BUILD_DIR}

# Build function
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local GOARM=$3
    local OUTPUT=$4
    local DESCRIPTION=$5

    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Platform:${NC} ${DESCRIPTION}"
    echo -e "${BLUE}Output:${NC}   ${BUILD_DIR}/${OUTPUT}"

    # Build
    if [ -n "$GOARM" ]; then
        (cd ${SOURCE_DIR} && GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM CGO_ENABLED=0 go build \
            -ldflags="-s -w" \
            -trimpath \
            -o "${BUILD_DIR}/${OUTPUT}" \
            . 2>&1 | grep -v "^#" || true)
    else
        (cd ${SOURCE_DIR} && GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build \
            -ldflags="-s -w" \
            -trimpath \
            -o "${BUILD_DIR}/${OUTPUT}" \
            . 2>&1 | grep -v "^#" || true)
    fi

    if [ $? -eq 0 ] && [ -f "${BUILD_DIR}/${OUTPUT}" ]; then
        SIZE=$(ls -lh "${BUILD_DIR}/${OUTPUT}" | awk '{print $5}')
        echo -e "${GREEN}✓ Success!${NC} Size: ${SIZE}"

        # Compress with UPX (if available)
        if command -v upx &> /dev/null; then
            echo -e "${YELLOW}  ⚡ Compressing with UPX...${NC}"
            upx --best --lzma "${BUILD_DIR}/${OUTPUT}" >/dev/null 2>&1 && {
                SIZE_AFTER=$(ls -lh "${BUILD_DIR}/${OUTPUT}" | awk '{print $5}')
                echo -e "${GREEN}  ✓ Compressed!${NC} New size: ${SIZE_AFTER}"
            } || {
                echo -e "${YELLOW}  ⚠ UPX compression failed, original binary preserved${NC}"
            }
        fi
        echo ""
    else
        echo -e "${RED}✗ Build failed!${NC}\n"
        return 1
    fi
}

# Main build process
echo -e "${BLUE}Source Directory:${NC} ${SOURCE_DIR}"
echo -e "${BLUE}Build Directory:${NC}  ${BUILD_DIR}\n"

# MIPS Platforms
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                     MIPS Platforms                         ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform linux mipsle "" "${BINARY_NAME}-mipsle" "MIPS Little Endian (TP-Link, GL.iNet, Xiaomi)"
build_platform linux mips "" "${BINARY_NAME}-mips" "MIPS Big Endian (Older Broadcom)"

# ARM Platforms
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                      ARM Platforms                         ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform linux arm 5 "${BINARY_NAME}-armv5" "ARMv5 (Older devices)"
build_platform linux arm 6 "${BINARY_NAME}-armv6" "ARMv6 (Raspberry Pi Zero)"
build_platform linux arm 7 "${BINARY_NAME}-armv7" "ARMv7 (Raspberry Pi 2/3, GL.iNet)"
build_platform linux arm64 "" "${BINARY_NAME}-arm64" "ARM64 (Raspberry Pi 3+/4, Modern routers)"

# x86 Platforms
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                      x86 Platforms                         ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform linux 386 "" "${BINARY_NAME}-i386" "x86 32-bit (Older PC)"
build_platform linux amd64 "" "${BINARY_NAME}-amd64" "x86-64 (PC Engines APU, Protectli, Modern PC)"

# Summary
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                   Build Completed!                         ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

# Binary list
echo -e "${BLUE}Generated Binaries:${NC}\n"
ls -lh ${BUILD_DIR}/${BINARY_NAME}-* 2>/dev/null | awk '{printf "  %s  %9s  %s\n", $6" "$7" "$8, $5, $9}' || echo "No binaries were created!"

# Generate SHA256 checksums
if ls ${BUILD_DIR}/${BINARY_NAME}-* >/dev/null 2>&1; then
    echo -e "\n${YELLOW}⚡ Generating SHA256 checksums...${NC}"
    cd ${BUILD_DIR} && sha256sum ${BINARY_NAME}-* > SHA256SUMS 2>/dev/null
    cd - >/dev/null
    echo -e "${GREEN}✓ Checksum file:${NC} ${BUILD_DIR}/SHA256SUMS"

    # Total size
    TOTAL_SIZE=$(du -sh ${BUILD_DIR} | awk '{print $1}')
    echo -e "\n${BLUE}Total Size:${NC} ${TOTAL_SIZE}"
fi

echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                         Success!                           ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

echo -e "${YELLOW}Tip:${NC} If UPX is not installed, you can install it with 'sudo apt install upx'."
echo -e "${YELLOW}Usage:${NC} You can find the binaries in the ${BUILD_DIR}/ directory.\n"

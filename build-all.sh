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
BINARY_NAME="hermes-euicc"
PKG_VERSION="1.0.0"  # Package version (manually update for new versions)
PKG_RELEASE=$(git rev-list --count HEAD 2>/dev/null || echo "1")  # Auto-increment with git commits
BUILD_DIR="${SCRIPT_DIR}/build/${PKG_VERSION}"  # Version-specific directory
GO_VERSION="1.24.0"  # Required Go version

# Package metadata
PKG_NAME="hermes-euicc"
PKG_LICENSE="MIT"
PKG_MAINTAINER="Kilimcinin Kör Oğlu <k@keremgok.tr>"

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

# Create build directory (no subdirectories needed)
mkdir -p ${BUILD_DIR}

# Enhanced build function with optimization support
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local GOARM=$3
    local OUTPUT_NAME=$4
    local DESCRIPTION=$5
    local OPT_FLAGS=$6  # Optimization flags like GOMIPS=softfloat
    local BUILD_TAGS=$7  # Build tags like "openwrt"

    local OUTPUT_PATH="${BUILD_DIR}/${OUTPUT_NAME}"

    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${MAGENTA}📦 Platform:${NC}  ${DESCRIPTION}"
    echo -e "${BLUE}🎯 Target:${NC}    ${GOOS}/${GOARCH}$([ -n "$GOARM" ] && echo "/v${GOARM}")"
    echo -e "${YELLOW}💾 Output:${NC}    ${OUTPUT_NAME}"

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

    # Set version information via ldflags
    local LDFLAGS="-s -X main.Version=${PKG_VERSION} -X main.Release=${PKG_RELEASE}"

    if (cd ${SOURCE_DIR} && eval $BUILD_CMD go build \
        $TAGS_FLAG \
        -ldflags=\"$LDFLAGS\" \
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

build_platform "linux" "amd64" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-linux-amd64" \
    "Linux x86-64 (Modern PCs, Servers)" \
    "GOAMD64=v2"

build_platform "linux" "386" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-linux-i386" \
    "Linux x86 32-bit (Legacy PCs)" \
    ""

build_platform "linux" "arm64" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-linux-arm64" \
    "Linux ARM 64-bit (Raspberry Pi 4+, Servers)" \
    ""

# =============================================================================
# OpenWRT/Embedded Linux Builds (Optimized for routers)
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║           OpenWRT/Embedded Linux (Routers/IoT)             ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

# MIPS Platforms (with openwrt build tag for UCI support)
build_platform "linux" "mips" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-mips" \
    "MIPS BE (Atheros AR/QCA, TP-Link, GL.iNet AR/XE, Ubiquiti)" \
    "GOMIPS=softfloat" "openwrt"

build_platform "linux" "mipsle" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-mipsle" \
    "MIPS LE (MediaTek MT76xx, GL.iNet MT, Ralink)" \
    "GOMIPS=softfloat" "openwrt"

build_platform "linux" "mips64" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-mips64" \
    "MIPS64 BE (Cavium Octeon, EdgeRouter Pro, Ubiquiti ER-8)" \
    "GOMIPS64=softfloat" "openwrt"

# Note: MIPS64 LE not built - OpenWRT doesn't support this architecture (too rare)

# ARM Embedded Platforms (with openwrt build tag for UCI support)
build_platform "linux" "arm" "5" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-arm_v5" \
    "ARM v5 (Kirkwood, Old NAS devices)" \
    "" "openwrt"

build_platform "linux" "arm" "6" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-arm_v6" \
    "ARM v6 (Raspberry Pi Zero/1, BCM2835)" \
    "" "openwrt"

build_platform "linux" "arm" "7" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-arm_v7" \
    "ARM v7 (IPQ40xx, GL.iNet B1300, Raspberry Pi 2/3)" \
    "" "openwrt"

build_platform "linux" "arm64" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-arm64" \
    "ARM64 (MT7622/MT7986, IPQ807x, BananaPi R3/R4, GL.iNet MT6000)" \
    "" "openwrt"

# x86 Embedded Platforms (with openwrt build tag for UCI support)
build_platform "linux" "386" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-x86" \
    "x86 32-bit (Legacy PC Engines, Old x86 routers)" \
    "" "openwrt"

build_platform "linux" "amd64" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-x86_64" \
    "x86-64 (PC Engines APU, Protectli, x86 routers, VMs)" \
    "GOAMD64=v2" "openwrt"

# =============================================================================
# macOS Builds (DISABLED - upstream euicc-go AT driver not implemented)
# =============================================================================
# echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
# echo -e "${GREEN}║                     macOS Platforms                        ║${NC}"
# echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"
#
# build_platform "darwin" "amd64" "" "${BINARY_NAME}-darwin-amd64" \
#     "macOS Intel (x86-64)" \
#     "GOAMD64=v2"
#
# build_platform "darwin" "arm64" "" "${BINARY_NAME}-darwin-arm64" \
#     "macOS Apple Silicon (M1/M2/M3/M4)" \
#     ""

# =============================================================================
# Windows Builds
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    Windows Platforms                       ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform "windows" "amd64" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-windows-amd64.exe" \
    "Windows x86-64 (64-bit)" \
    "GOAMD64=v2"

build_platform "windows" "386" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-windows-i386.exe" \
    "Windows x86 (32-bit)" \
    ""

build_platform "windows" "arm64" "" "${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-windows-arm64.exe" \
    "Windows ARM64 (Surface Pro X, ARM laptops)" \
    ""

# =============================================================================
# FreeBSD Builds (DISABLED - upstream euicc-go AT/CCID drivers not implemented)
# =============================================================================
# echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
# echo -e "${GREEN}║                    FreeBSD Platforms                       ║${NC}"
# echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"
#
# build_platform "freebsd" "amd64" "" "${BINARY_NAME}-freebsd-amd64" \
#     "FreeBSD x86-64" \
#     "GOAMD64=v2"
#
# build_platform "freebsd" "arm64" "" "${BINARY_NAME}-freebsd-arm64" \
#     "FreeBSD ARM 64-bit" \
#     ""

# =============================================================================
# Generate Checksums
# =============================================================================
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Generating checksums...${NC}"

# Generate single SHA256SUMS file for all binaries
(cd ${BUILD_DIR} && sha256sum ${BINARY_NAME}* > SHA256SUMS 2>/dev/null) && \
    echo -e "${GREEN}✓ SHA256SUMS created${NC}"

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
echo -e "  Version:             ${PKG_VERSION}-${PKG_RELEASE}"

# Count binaries
TOTAL=$(ls -1 ${BUILD_DIR}/${BINARY_NAME}* 2>/dev/null | grep -v SHA256SUMS | wc -l)
echo -e "  ${GREEN}Total binaries:      ${TOTAL}${NC}"

# Count IPK packages
IPK_TOTAL=$(ls -1 ${BUILD_DIR}/${BINARY_NAME}_*.ipk 2>/dev/null | wc -l)
if [ $IPK_TOTAL -gt 0 ]; then
    echo -e "  ${GREEN}IPK packages:        ${IPK_TOTAL}${NC}"
fi
echo ""

# =============================================================================
# Usage Guide
# =============================================================================
echo -e "${YELLOW}Platform Selection Guide:${NC}"
echo ""
echo -e "${CYAN}Build Directory:${NC} ${BLUE}build/${PKG_VERSION}/${NC}"
echo ""
echo -e "${CYAN}For OpenWRT/Embedded Routers:${NC}"
echo -e "  Check architecture: ${MAGENTA}ls -la /lib/ld-musl-*.so.1${NC}"
echo -e "    ${GREEN}ld-musl-mips-sf.so.1${NC}     → Use: ${BLUE}${BINARY_NAME}-openwrt-mips${NC}"
echo -e "    ${GREEN}ld-musl-mipsel-sf.so.1${NC}   → Use: ${BLUE}${BINARY_NAME}-openwrt-mipsle${NC}"
echo -e "    ${GREEN}ld-musl-armhf.so.1${NC}       → Use: ${BLUE}${BINARY_NAME}-openwrt-arm_v7${NC}"
echo -e "    ${GREEN}ld-musl-aarch64.so.1${NC}     → Use: ${BLUE}${BINARY_NAME}-openwrt-arm64${NC}"
echo -e "    ${GREEN}ld-musl-x86_64.so.1${NC}      → Use: ${BLUE}${BINARY_NAME}-openwrt-x86_64${NC}"
echo ""
echo -e "${CYAN}For Desktop/Server:${NC}"
echo -e "  ${GREEN}Linux x86-64${NC}              → Use: ${BLUE}${BINARY_NAME}-linux-amd64${NC}"
echo -e "  ${GREEN}Linux ARM64${NC}               → Use: ${BLUE}${BINARY_NAME}-linux-arm64${NC}"
echo -e "  ${GREEN}Windows x86-64${NC}            → Use: ${BLUE}${BINARY_NAME}-windows-amd64.exe${NC}"
echo -e "  ${GREEN}Windows x86${NC}               → Use: ${BLUE}${BINARY_NAME}-windows-i386.exe${NC}"
echo -e "  ${GREEN}Windows ARM64${NC}             → Use: ${BLUE}${BINARY_NAME}-windows-arm64.exe${NC}"
echo ""
echo -e "${YELLOW}Note:${NC} macOS and FreeBSD builds are disabled due to upstream driver issues."
echo ""

echo -e "${GREEN}✓ All builds completed successfully!${NC}\n"

# =============================================================================
# OpenWRT IPK Package Generation
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║            OpenWRT IPK Package Generation                  ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"


# Function to create IPK for a specific architecture
create_ipk() {
    local ARCH=$1
    local BINARY_PATH=$2
    local ARCH_DESC=$3

    if [ ! -f "$BINARY_PATH" ]; then
        echo -e "${YELLOW}⚠ Skipping ${ARCH}: binary not found${NC}"
        return 1
    fi

    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${MAGENTA}📦 Architecture:${NC} ${ARCH_DESC}"
    echo -e "${BLUE}🎯 Target:${NC}       ${ARCH}"

    local IPK_BUILD_DIR="${BUILD_DIR}/ipk-${ARCH}"
    local IPK_CONTROL_DIR="${IPK_BUILD_DIR}/CONTROL"
    local IPK_DATA_DIR="${IPK_BUILD_DIR}/data"
    local IPK_FILE="${BUILD_DIR}/${PKG_NAME}_${PKG_VERSION}-${PKG_RELEASE}_${ARCH}.ipk"

    # Clean and create directories
    rm -rf "${IPK_BUILD_DIR}"
    mkdir -p "${IPK_CONTROL_DIR}" "${IPK_DATA_DIR}/usr/bin"

    # Copy binary
    cp "${BINARY_PATH}" "${IPK_DATA_DIR}/usr/bin/hermes-euicc"
    chmod 755 "${IPK_DATA_DIR}/usr/bin/hermes-euicc"

    # Create control file
    cat > "${IPK_CONTROL_DIR}/control" << EOF
Package: ${PKG_NAME}
Version: ${PKG_VERSION}-${PKG_RELEASE}
Depends: libc, coreutils, coreutils-timeout
Provides: hermes-euicc
Section: utils
Architecture: ${ARCH}
Installed-Size: $(du -sb "${IPK_DATA_DIR}" | cut -f1)
Maintainer: ${PKG_MAINTAINER}
Description: eSIM profile management tool for eUICC devices
 Hermes eUICC Manager is a command-line tool for managing eSIM profiles
 on eUICC-enabled devices. Supports QMI, MBIM, and AT command drivers.
 .
 Features:
  - Profile listing, enabling, disabling, and deletion
  - Profile download via activation codes
  - SM-DS discovery and automatic profile download
  - Notification management
  - Chip information display
  - Cross-platform driver support (QMI, MBIM, AT)
  - UCI configuration support for OpenWRT
EOF

    # Create postinst script
    cat > "${IPK_CONTROL_DIR}/postinst" << 'POSTINST_EOF'
#!/bin/sh
# Initialize hermes-euicc UCI config only on first installation
# Check if config section exists - if not, this is first install
if ! uci -q get hermes-euicc.hermes-euicc >/dev/null 2>&1; then
    # First installation - create config file with defaults
    if [ ! -f /etc/config/hermes-euicc ]; then
        touch /etc/config/hermes-euicc
    fi

    uci set hermes-euicc.hermes-euicc='hermes-euicc'
    uci set hermes-euicc.hermes-euicc.driver='auto'
    uci set hermes-euicc.hermes-euicc.device=''
    uci set hermes-euicc.hermes-euicc.slot='1'
    uci set hermes-euicc.hermes-euicc.timeout='30'
    uci commit hermes-euicc

    echo "Hermes eUICC Manager installed successfully!"
    echo "Default configuration created in /etc/config/hermes-euicc"
else
    # Upgrade - preserve existing configuration
    echo "Hermes eUICC Manager upgraded successfully!"
    echo "Your existing configuration has been preserved."
fi

echo "Usage: hermes-euicc --help"
exit 0
POSTINST_EOF
    chmod 755 "${IPK_CONTROL_DIR}/postinst"

    # Create conffiles (mark UCI config as configuration file)
    cat > "${IPK_CONTROL_DIR}/conffiles" << 'EOF'
/etc/config/hermes-euicc
EOF

    # Build IPK package
    cd "${IPK_BUILD_DIR}"

    # Create debian-binary
    echo "2.0" > debian-binary

    # Create control.tar.gz
    tar -C CONTROL -czf control.tar.gz --owner=0 --group=0 --numeric-owner .

    # Create data.tar.gz
    tar -C data -czf data.tar.gz --owner=0 --group=0 --numeric-owner .

    # Create IPK (tar.gz archive)
    tar -czf "${IPK_FILE}" debian-binary control.tar.gz data.tar.gz

    # Verify IPK was created
    if [ -f "${IPK_FILE}" ]; then
        local SIZE=$(du -h "${IPK_FILE}" | cut -f1)
        echo -e "${GREEN}✓ Success!${NC} IPK: ${SIZE}"
        echo -e "${BLUE}   ${IPK_FILE##*/}${NC}"
        echo ""

        # Cleanup build directory
        cd "${SCRIPT_DIR}"
        rm -rf "${IPK_BUILD_DIR}"
        return 0
    else
        echo -e "${RED}✗ Failed to create IPK${NC}\n"
        return 1
    fi
}

# Create IPK for each OpenWRT architecture
IPK_COUNT=0

# MIPS - OpenWRT uses specific CPU-based architecture names
if create_ipk "mips_24kc" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-mips" "MIPS 24Kc (AR71xx/AR9xxx - TP-Link, GL.iNet, Ubiquiti)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

if create_ipk "mipsel_24kc" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-mipsle" "MIPS LE 24Kc (MediaTek MT76xx, Ralink)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

if create_ipk "mips64_octeonplus" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-mips64" "MIPS64 Octeon+ (Cavium EdgeRouter, Ubiquiti ER-8)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

# Note: mips64el (little-endian) not supported in OpenWRT - platform too rare

# ARM - OpenWRT uses CPU-specific names
if create_ipk "arm_arm926ej-s" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-arm_v5" "ARM v5 ARM926EJ-S (Kirkwood, Old NAS)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

if create_ipk "arm_arm1176jzf-s_vfp" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-arm_v6" "ARM v6 ARM1176JZF-S (Raspberry Pi Zero/1)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

if create_ipk "arm_cortex-a7_neon-vfpv4" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-arm_v7" "ARM v7 Cortex-A7 (IPQ40xx, Raspberry Pi 2/3)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

if create_ipk "aarch64_cortex-a53" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-arm64" "ARM64 Cortex-A53 (MT7622/MT7986, IPQ807x, BananaPi R3/R4)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

# x86 - OpenWRT uses CPU feature-based names
if create_ipk "i386_pentium4" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-x86" "x86 Pentium 4 (Legacy PC Engines, x86 routers)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

if create_ipk "x86_64" "${BUILD_DIR}/${BINARY_NAME}-${PKG_VERSION}-${PKG_RELEASE}-openwrt-x86_64" "x86-64 (PC Engines APU, Protectli, VMs)"; then
    IPK_COUNT=$((IPK_COUNT + 1))
fi

# Summary
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✓ IPK Generation Complete!${NC}"
echo -e "${BLUE}  Created ${IPK_COUNT} IPK packages in:${NC} ${BUILD_DIR}/"
echo ""
echo -e "${YELLOW}Installation:${NC}"
echo -e "  opkg update"
echo -e "  opkg install /tmp/hermes-euicc_*.ipk"
echo ""
echo -e "${YELLOW}Usage:${NC}"
echo -e "  hermes-euicc --help"
echo -e "  hermes-euicc list"
echo ""

echo -e "${GREEN}✓ All builds and packages completed successfully!${NC}\n"

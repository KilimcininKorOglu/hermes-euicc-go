#!/bin/bash

# Hermes eUICC Manager - Multi-Platform Build Script
# Includes all drivers (QMI, MBIM, AT, CCID)

set -e  # Hata durumunda dur

# Renkli çıktı
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Banner
echo -e "${BLUE}"
echo "╔════════════════════════════════════════════════════════════╗"
echo "║       Hermes eUICC Manager - Multi-Platform Build         ║"
echo "║              Tüm Sürücüler Dahil (All Drivers)             ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo -e "${NC}\n"

# Yapılandırma
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_DIR="${SCRIPT_DIR}"
BUILD_DIR="${SCRIPT_DIR}/build"
BINARY_NAME="hermes-euicc"

# Build dizini oluştur
mkdir -p ${BUILD_DIR}

# Derleme fonksiyonu
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local GOARM=$3
    local OUTPUT=$4
    local DESCRIPTION=$5

    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Platform:${NC} ${DESCRIPTION}"
    echo -e "${BLUE}Output:${NC}   ${BUILD_DIR}/${OUTPUT}"

    # Derleme
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
        echo -e "${GREEN}✓ Başarılı!${NC} Boyut: ${SIZE}"

        # UPX ile sıkıştır (varsa)
        if command -v upx &> /dev/null; then
            echo -e "${YELLOW}  ⚡ UPX ile sıkıştırılıyor...${NC}"
            upx --best --lzma "${BUILD_DIR}/${OUTPUT}" >/dev/null 2>&1 && {
                SIZE_AFTER=$(ls -lh "${BUILD_DIR}/${OUTPUT}" | awk '{print $5}')
                echo -e "${GREEN}  ✓ Sıkıştırıldı!${NC} Yeni boyut: ${SIZE_AFTER}"
            } || {
                echo -e "${YELLOW}  ⚠ UPX sıkıştırma başarısız, orijinal binary korundu${NC}"
            }
        fi
        echo ""
    else
        echo -e "${RED}✗ Derleme başarısız!${NC}\n"
        return 1
    fi
}

# Ana derleme süreci
echo -e "${BLUE}Kaynak Dizin:${NC} ${SOURCE_DIR}"
echo -e "${BLUE}Hedef Dizin:${NC}  ${BUILD_DIR}\n"

# MIPS Platformlar
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    MIPS Platformlar                        ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform linux mipsle "" "${BINARY_NAME}-mipsle" "MIPS Little Endian (TP-Link, GL.iNet, Xiaomi)"
build_platform linux mips "" "${BINARY_NAME}-mips" "MIPS Big Endian (Eski Broadcom)"

# ARM Platformlar
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                     ARM Platformlar                        ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform linux arm 5 "${BINARY_NAME}-armv5" "ARMv5 (Eski cihazlar)"
build_platform linux arm 6 "${BINARY_NAME}-armv6" "ARMv6 (Raspberry Pi Zero)"
build_platform linux arm 7 "${BINARY_NAME}-armv7" "ARMv7 (Raspberry Pi 2/3, GL.iNet)"
build_platform linux arm64 "" "${BINARY_NAME}-arm64" "ARM64 (Raspberry Pi 3+/4, Modern routerlar)"

# x86 Platformlar
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                     x86 Platformlar                        ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

build_platform linux 386 "" "${BINARY_NAME}-i386" "x86 32-bit (Eski PC)"
build_platform linux amd64 "" "${BINARY_NAME}-amd64" "x86-64 (PC Engines APU, Protectli, Modern PC)"

# Özet
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                  Derleme Tamamlandı!                       ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

# Binary listesi
echo -e "${BLUE}Oluşturulan Binary'ler:${NC}\n"
ls -lh ${BUILD_DIR}/${BINARY_NAME}-* 2>/dev/null | awk '{printf "  %s  %9s  %s\n", $6" "$7" "$8, $5, $9}' || echo "Hiç binary oluşturulamadı!"

# SHA256 checksum oluştur
if ls ${BUILD_DIR}/${BINARY_NAME}-* >/dev/null 2>&1; then
    echo -e "\n${YELLOW}⚡ SHA256 checksum'lar oluşturuluyor...${NC}"
    cd ${BUILD_DIR} && sha256sum ${BINARY_NAME}-* > SHA256SUMS 2>/dev/null
    cd - >/dev/null
    echo -e "${GREEN}✓ Checksum dosyası:${NC} ${BUILD_DIR}/SHA256SUMS"

    # Toplam boyut
    TOTAL_SIZE=$(du -sh ${BUILD_DIR} | awk '{print $1}')
    echo -e "\n${BLUE}Toplam Boyut:${NC} ${TOTAL_SIZE}"
fi

echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                        Başarılı!                           ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

echo -e "${YELLOW}İpucu:${NC} UPX kurulu değilse, 'sudo apt install upx' ile kurabilirsiniz."
echo -e "${YELLOW}Kullanım:${NC} Binary'leri ${BUILD_DIR}/ dizininden alabilirsiniz.\n"

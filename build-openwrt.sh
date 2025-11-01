#!/bin/bash

# OpenWRT Popüler Cihazlar İçin eUICC Go Derleme Script'i
# En yaygın kullanılan router mimarileri için optimize edilmiş

set -e  # Hata durumunda dur

# Renkli çıktı
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
║  Hermes eUICC Manager - OpenWRT Popular Devices Build     ║
║              En Yaygın Router Mimarileri                   ║
╚════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}\n"

# Yapılandırma
BINARY_NAME="hermes-euicc"

# Determine source directory based on where script is run from
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_DIR="${SCRIPT_DIR}"
BUILD_DIR="${SCRIPT_DIR}/build/openwrt"

LDFLAGS="-ldflags=\"-s -w\""
TRIMPATH="-trimpath"

# Build dizini oluştur
mkdir -p ${BUILD_DIR}

# Derleme fonksiyonu
build_for_device() {
    local DEVICE_NAME=$1
    local GOOS=$2
    local GOARCH=$3
    local GOARM=$4
    local OUTPUT=$5
    local DESCRIPTION=$6
    local OPT_FLAGS=$7  # Yeni: Optimizasyon flagleri (GOMIPS/GOAMD64/etc)

    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${MAGENTA}📱 Cihaz:${NC}     ${DEVICE_NAME}"
    echo -e "${BLUE}📦 Platform:${NC}  ${DESCRIPTION}"
    echo -e "${YELLOW}💾 Output:${NC}    ${BUILD_DIR}/${OUTPUT}"

    # Optimizasyon bilgisini göster
    if [ -n "$OPT_FLAGS" ]; then
        echo -e "${GREEN}⚡ Optimize:${NC}  ${OPT_FLAGS}"
    fi

    # Derleme - optimizasyon flagleriyle
    # Export optimization flags if present
    if [ -n "$OPT_FLAGS" ]; then
        export $OPT_FLAGS
    fi

    if [ -n "$GOARM" ]; then
        (cd ${SOURCE_DIR} && CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM go build \
            -ldflags="-s -w" \
            -trimpath \
            -o "${BUILD_DIR}/${OUTPUT}" \
            . 2>&1 | grep -v "^#" || true)
    else
        (cd ${SOURCE_DIR} && CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
            -ldflags="-s -w" \
            -trimpath \
            -o "${BUILD_DIR}/${OUTPUT}" \
            . 2>&1 | grep -v "^#" || true)
    fi

    # Unset optimization flags to avoid pollution
    if [ -n "$OPT_FLAGS" ]; then
        unset $(echo "$OPT_FLAGS" | cut -d= -f1)
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
                echo -e "${YELLOW}  ⚠ UPX sıkıştırma atlandı${NC}"
            }
        fi
        echo ""
    else
        echo -e "${RED}✗ Derleme başarısız!${NC}\n"
        return 1
    fi
}

# Ana derleme süreci
echo -e "${BLUE}📁 Kaynak Dizin:${NC} ${SOURCE_DIR}"
echo -e "${BLUE}📁 Hedef Dizin:${NC}  ${BUILD_DIR}\n"

# =============================================================================
# MIPS Cihazlar (En Yaygın OpenWRT Cihazları)
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              MIPS Tabanlı Cihazlar (En Yaygın)            ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

# GL.iNet Serisi (QCA9563 MIPS 74Kc - FPU var)
build_for_device "GL.iNet GL-AR750S (Slate)" "linux" "mipsle" "" \
    "${BINARY_NAME}-glinet-ar750s" \
    "MIPS LE (GL.iNet AR750S, AR300M, MT300N-V2)" \
    "GOMIPS=hardfloat"

build_for_device "GL.iNet GL-MT300N-V2 (Mango)" "linux" "mipsle" "" \
    "${BINARY_NAME}-glinet-mt300n-v2" \
    "MIPS LE (GL.iNet Mango serisi)" \
    "GOMIPS=hardfloat"

# TP-Link Serisi (QCA9558/QCA9563 MIPS 74Kc - FPU var)
build_for_device "TP-Link Archer C7 v2/v4/v5" "linux" "mipsle" "" \
    "${BINARY_NAME}-tplink-archerc7" \
    "MIPS LE (TP-Link Archer C7, C5, C2)" \
    "GOMIPS=hardfloat"

build_for_device "TP-Link TL-WR1043ND v1/v2/v3" "linux" "mipsle" "" \
    "${BINARY_NAME}-tplink-wr1043nd" \
    "MIPS LE (TP-Link WR1043ND, WR841N)" \
    "GOMIPS=hardfloat"

# Xiaomi Serisi (MT7621 MIPS 1004Kc - FPU var)
build_for_device "Xiaomi Mi Router 3G" "linux" "mipsle" "" \
    "${BINARY_NAME}-xiaomi-mir3g" \
    "MIPS LE (Xiaomi Mi Router 3G, 3, 4A)" \
    "GOMIPS=hardfloat"

build_for_device "Xiaomi Mi Router 4A Gigabit" "linux" "mipsle" "" \
    "${BINARY_NAME}-xiaomi-mir4ag" \
    "MIPS LE (Xiaomi Mi Router 4A Gigabit)" \
    "GOMIPS=softfloat"

# Ubiquiti Serisi (MT7621AT MIPS 1004Kc - FPU var)
build_for_device "Ubiquiti EdgeRouter X" "linux" "mipsle" "" \
    "${BINARY_NAME}-ubnt-erx" \
    "MIPS LE (Ubiquiti EdgeRouter X, ER-X-SFP)" \
    "GOMIPS=hardfloat"

# TP-Link Budget Serisi (QCA9563 MIPS 74Kc - FPU var)
build_for_device "TP-Link Archer A7" "linux" "mipsle" "" \
    "${BINARY_NAME}-tplink-archera7" \
    "MIPS LE (TP-Link Archer A7, Archer C7 successor)" \
    "GOMIPS=hardfloat"

# TP-Link Archer C6 (MT7621 bazı versiyonlarda FPU yok - güvenli seçim)
build_for_device "TP-Link Archer C6" "linux" "mipsle" "" \
    "${BINARY_NAME}-tplink-archerc6" \
    "MIPS LE (TP-Link Archer C6, budget-friendly)" \
    "GOMIPS=softfloat"

# TP-Link WR841N (AR9341 MIPS 24Kc - FPU YOK!)
build_for_device "TP-Link TL-WR841N" "linux" "mipsle" "" \
    "${BINARY_NAME}-tplink-wr841n" \
    "MIPS LE (TP-Link TL-WR841N, world's best-selling budget router)" \
    "GOMIPS=softfloat"

# GL.iNet 4G LTE Serisi (QCA9531 MIPS 24Kc - FPU YOK!)
build_for_device "GL.iNet GL-XE300 (Puli)" "linux" "mipsle" "" \
    "${BINARY_NAME}-glinet-xe300" \
    "MIPS LE (GL.iNet XE300 Puli, 4G LTE, 5000mAh battery)" \
    "GOMIPS=softfloat"

# =============================================================================
# ARM Cihazlar (Modern ve Popüler)
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║             ARM Tabanlı Cihazlar (Modern)                  ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

# GL.iNet ARM Serisi (Qualcomm IPQ4018 Cortex-A7 + NEON)
build_for_device "GL.iNet GL-B1300 (Convexa-B)" "linux" "arm" "7" \
    "${BINARY_NAME}-glinet-b1300" \
    "ARMv7 (GL.iNet B1300, S1300)" \
    ""

build_for_device "GL.iNet GL-AX1800 (Flint)" "linux" "arm64" "" \
    "${BINARY_NAME}-glinet-ax1800" \
    "ARM64 (GL.iNet AX1800, AXT1800)" \
    ""

# Raspberry Pi Serisi (Cortex-A53/A72/A76)
build_for_device "Raspberry Pi 2 Model B" "linux" "arm" "7" \
    "${BINARY_NAME}-rpi2" \
    "ARMv7 (Raspberry Pi 2, 3)" \
    ""

build_for_device "Raspberry Pi 3 Model B/B+" "linux" "arm" "7" \
    "${BINARY_NAME}-rpi3" \
    "ARMv7 (Raspberry Pi 3 B/B+)" \
    ""

build_for_device "Raspberry Pi 4 Model B" "linux" "arm64" "" \
    "${BINARY_NAME}-rpi4" \
    "ARM64 (Raspberry Pi 4, 400, CM4)" \
    ""

build_for_device "Raspberry Pi 5" "linux" "arm64" "" \
    "${BINARY_NAME}-rpi5" \
    "ARM64 (Raspberry Pi 5)" \
    ""

# Linksys ARM Serisi (Marvell Armada 385 Cortex-A9)
build_for_device "Linksys WRT1900ACS/WRT3200ACM" "linux" "arm" "7" \
    "${BINARY_NAME}-linksys-wrt1900" \
    "ARMv7 (Linksys WRT serisi)" \
    ""

# NanoPi Serisi (Rockchip RK3328/RK3399 Cortex-A53/A72)
build_for_device "NanoPi R2S" "linux" "arm64" "" \
    "${BINARY_NAME}-nanopi-r2s" \
    "ARM64 (NanoPi R2S, R4S, R5S)" \
    ""

build_for_device "NanoPi R4S" "linux" "arm64" "" \
    "${BINARY_NAME}-nanopi-r4s" \
    "ARM64 (NanoPi R4S - RK3399)" \
    ""

# GL.iNet WiFi 6/7 Flagship Serisi (MediaTek Filogic Cortex-A53)
build_for_device "GL.iNet GL-MT6000 (Flint 2)" "linux" "arm64" "" \
    "${BINARY_NAME}-glinet-mt6000" \
    "ARM64 (GL.iNet MT6000 Flint 2, WiFi 6, MediaTek MT7986, 2025 flagship)" \
    ""

build_for_device "GL.iNet GL-MT3000 (Beryl AX)" "linux" "arm64" "" \
    "${BINARY_NAME}-glinet-mt3000" \
    "ARM64 (GL.iNet MT3000 Beryl AX, WiFi 6 travel router)" \
    ""

build_for_device "GL.iNet Slate 7" "linux" "arm64" "" \
    "${BINARY_NAME}-glinet-slate7" \
    "ARM64 (GL.iNet Slate 7, WiFi 7 flagship, touchscreen)" \
    ""

# Linksys WiFi 6 Serisi (MediaTek MT7622 Cortex-A53)
build_for_device "Linksys E8450" "linux" "arm" "7" \
    "${BINARY_NAME}-linksys-e8450" \
    "ARMv7 (Linksys E8450, WiFi 6, MediaTek Filogic, 3.2 Gbps)" \
    ""

build_for_device "Linksys EA8300" "linux" "arm" "7" \
    "${BINARY_NAME}-linksys-ea8300" \
    "ARMv7 (Linksys EA8300, Tri-band, popular mid-range)" \
    ""

# Banana Pi Developer Boards (MediaTek Filogic Cortex-A53/A73)
build_for_device "Banana Pi BPI-R3" "linux" "arm64" "" \
    "${BINARY_NAME}-bananapi-r3" \
    "ARM64 (Banana Pi R3, WiFi 6, MediaTek MT7986, developer board)" \
    ""

build_for_device "Banana Pi BPI-R4" "linux" "arm64" "" \
    "${BINARY_NAME}-bananapi-r4" \
    "ARM64 (Banana Pi R4, WiFi 7, MediaTek MT7988A, cutting-edge)" \
    ""

build_for_device "Banana Pi BPI-R3 Mini" "linux" "arm64" "" \
    "${BINARY_NAME}-bananapi-r3mini" \
    "ARM64 (Banana Pi R3 Mini, WiFi 6, compact version)" \
    ""

# Dynalink WiFi 6 (Qualcomm IPQ8072A Cortex-A53)
build_for_device "Dynalink DL-WRX36" "linux" "arm" "7" \
    "${BINARY_NAME}-dynalink-wrx36" \
    "ARMv7 (Dynalink DL-WRX36, WiFi 6 AX3600, Qualcomm IPQ8072A)" \
    ""

# Cudy WiFi 6 Serisi (MediaTek MT7981B Cortex-A53)
build_for_device "Cudy WR3000" "linux" "arm64" "" \
    "${BINARY_NAME}-cudy-wr3000" \
    "ARM64 (Cudy WR3000, WiFi 6 AX3000, MediaTek MT7981B, budget)" \
    ""

build_for_device "Cudy P5 5G" "linux" "arm64" "" \
    "${BINARY_NAME}-cudy-p5" \
    "ARM64 (Cudy P5, WiFi 6 + 5G modem, unique feature)" \
    ""

# MikroTik Enterprise Serisi (Qualcomm IPQ4019 Cortex-A7 / Marvell Cortex-A72)
build_for_device "MikroTik hAP ac2" "linux" "arm" "7" \
    "${BINARY_NAME}-mikrotik-hap-ac2" \
    "ARMv7 (MikroTik hAP ac2, Qualcomm IPQ4019, enterprise)" \
    ""

build_for_device "MikroTik hAP ac3" "linux" "arm" "7" \
    "${BINARY_NAME}-mikrotik-hap-ac3" \
    "ARMv7 (MikroTik hAP ac3, enterprise, newer version)" \
    ""

build_for_device "MikroTik RB5009" "linux" "arm64" "" \
    "${BINARY_NAME}-mikrotik-rb5009" \
    "ARM64 (MikroTik RB5009, Cortex-A72, high-end enterprise)" \
    ""

# =============================================================================
# x86 Cihazlar (PC/Server Tabanlı)
# =============================================================================
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║            x86 Tabanlı Cihazlar (PC/Server)                ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

# PC Engines Serisi (AMD GX-412TC - 2013, Jaguar architecture)
build_for_device "PC Engines APU2/APU3/APU4" "linux" "amd64" "" \
    "${BINARY_NAME}-pcengines-apu" \
    "x86-64 (PC Engines APU serisi)" \
    "GOAMD64=v2"

# Protectli Serisi (Modern Intel - Haswell+)
build_for_device "Protectli Vault FW4B/FW6" "linux" "amd64" "" \
    "${BINARY_NAME}-protectli-vault" \
    "x86-64 (Protectli Vault serisi)" \
    "GOAMD64=v3"

# Genel x86 (Maksimum uyumluluk)
build_for_device "Generic x86-64 (64-bit)" "linux" "amd64" "" \
    "${BINARY_NAME}-x86-64" \
    "x86-64 (Generic 64-bit PC)" \
    "GOAMD64=v1"

build_for_device "Generic x86 (32-bit)" "linux" "386" "" \
    "${BINARY_NAME}-x86-32" \
    "x86 (Generic 32-bit PC)" \
    ""

# Qotom Mini PC (Intel Celeron J3xxx/J4xxx - Apollo Lake/Gemini Lake = Goldmont)
build_for_device "Qotom Q355G6" "linux" "amd64" "" \
    "${BINARY_NAME}-qotom-q355g6" \
    "x86-64 (Qotom Q355G6, 6x GbE, Intel Celeron, pfSense/OpenWRT box)" \
    "GOAMD64=v3"

# Minisforum Enterprise (Intel 12th/13th Gen - Alder Lake/Raptor Lake)
build_for_device "Minisforum MS-01" "linux" "amd64" "" \
    "${BINARY_NAME}-minisforum-ms01" \
    "x86-64 (Minisforum MS-01, 10GbE, enterprise homelab)" \
    "GOAMD64=v3"

# =============================================================================
# Özet ve Bilgilendirme
# =============================================================================
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                  Derleme Tamamlandı!                       ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

# Binary listesi
echo -e "${BLUE}📦 Oluşturulan Binary'ler:${NC}\n"
if ls ${BUILD_DIR}/${BINARY_NAME}-* >/dev/null 2>&1; then
    echo -e "${CYAN}Cihaz Modeli                              Boyut      Dosya${NC}"
    echo -e "${CYAN}────────────────────────────────────────────────────────────${NC}"
    ls -lh ${BUILD_DIR}/${BINARY_NAME}-* | awk '{
        # Dosya adından cihaz adını çıkar
        split($9, parts, "/");
        filename = parts[length(parts)];
        gsub("hermes-euicc-", "", filename);

        # Cihaz adlarını düzelt
        # MIPS devices
        if (filename ~ /glinet-ar750s/) name = "GL.iNet AR750S";
        else if (filename ~ /glinet-mt300n/) name = "GL.iNet MT300N-V2";
        else if (filename ~ /glinet-xe300/) name = "GL.iNet XE300 Puli";
        else if (filename ~ /tplink-archerc7/) name = "TP-Link Archer C7";
        else if (filename ~ /tplink-archera7/) name = "TP-Link Archer A7";
        else if (filename ~ /tplink-archerc6/) name = "TP-Link Archer C6";
        else if (filename ~ /tplink-wr1043nd/) name = "TP-Link WR1043ND";
        else if (filename ~ /tplink-wr841n/) name = "TP-Link WR841N";
        else if (filename ~ /xiaomi-mir3g/) name = "Xiaomi Mi Router 3G";
        else if (filename ~ /xiaomi-mir4ag/) name = "Xiaomi Mi Router 4A";
        else if (filename ~ /ubnt-erx/) name = "Ubiquiti EdgeRouter X";
        # ARM devices
        else if (filename ~ /glinet-b1300/) name = "GL.iNet B1300";
        else if (filename ~ /glinet-ax1800/) name = "GL.iNet AX1800";
        else if (filename ~ /glinet-mt6000/) name = "GL.iNet MT6000 Flint 2";
        else if (filename ~ /glinet-mt3000/) name = "GL.iNet MT3000 Beryl AX";
        else if (filename ~ /glinet-slate7/) name = "GL.iNet Slate 7";
        else if (filename ~ /rpi2/) name = "Raspberry Pi 2";
        else if (filename ~ /rpi3/) name = "Raspberry Pi 3";
        else if (filename ~ /rpi4/) name = "Raspberry Pi 4";
        else if (filename ~ /rpi5/) name = "Raspberry Pi 5";
        else if (filename ~ /linksys-wrt1900/) name = "Linksys WRT1900";
        else if (filename ~ /linksys-e8450/) name = "Linksys E8450";
        else if (filename ~ /linksys-ea8300/) name = "Linksys EA8300";
        else if (filename ~ /nanopi-r2s/) name = "NanoPi R2S";
        else if (filename ~ /nanopi-r4s/) name = "NanoPi R4S";
        else if (filename ~ /bananapi-r3mini/) name = "Banana Pi R3 Mini";
        else if (filename ~ /bananapi-r3/) name = "Banana Pi R3";
        else if (filename ~ /bananapi-r4/) name = "Banana Pi R4";
        else if (filename ~ /dynalink-wrx36/) name = "Dynalink DL-WRX36";
        else if (filename ~ /cudy-wr3000/) name = "Cudy WR3000";
        else if (filename ~ /cudy-p5/) name = "Cudy P5 5G";
        else if (filename ~ /mikrotik-hap-ac2/) name = "MikroTik hAP ac2";
        else if (filename ~ /mikrotik-hap-ac3/) name = "MikroTik hAP ac3";
        else if (filename ~ /mikrotik-rb5009/) name = "MikroTik RB5009";
        # x86 devices
        else if (filename ~ /pcengines-apu/) name = "PC Engines APU";
        else if (filename ~ /protectli-vault/) name = "Protectli Vault";
        else if (filename ~ /qotom-q355g6/) name = "Qotom Q355G6";
        else if (filename ~ /minisforum-ms01/) name = "Minisforum MS-01";
        else if (filename ~ /x86-64/) name = "Generic x86-64";
        else if (filename ~ /x86-32/) name = "Generic x86";
        else name = filename;

        printf "%-40s %9s  %s\n", name, $5, filename;
    }'

    # SHA256 checksum oluştur
    echo -e "\n${YELLOW}⚡ SHA256 checksum'lar oluşturuluyor...${NC}"
    cd ${BUILD_DIR} && sha256sum ${BINARY_NAME}-* > SHA256SUMS 2>/dev/null
    cd - >/dev/null
    echo -e "${GREEN}✓ Checksum dosyası:${NC} ${BUILD_DIR}/SHA256SUMS"

    # Toplam boyut
    TOTAL_SIZE=$(du -sh ${BUILD_DIR} | awk '{print $1}')
    echo -e "\n${BLUE}📊 Toplam Boyut:${NC} ${TOTAL_SIZE}"
else
    echo -e "${RED}Hiç binary oluşturulamadı!${NC}"
fi

# Kullanım talimatları
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                   Kullanım Talimatları                     ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"

cat << 'INSTRUCTIONS'
Router'a Yükleme:
  1. Cihazınıza uygun binary'yi seçin
  2. Router'a kopyalayın:
     scp build/openwrt/hermes-euicc-XXXX root@192.168.1.1:/usr/bin/hermes-euicc

  3. Çalıştırılabilir yapın:
     ssh root@192.168.1.1 "chmod +x /usr/bin/hermes-euicc"

  4. Kullanın!
     hermes-euicc list
     hermes-euicc eid

Daha Fazla Bilgi:
  - OPENWRT_INTEGRATION.md: UCI, init scripts, CLI wrapper
  - QUICK_START.md: Hızlı başlangıç rehberi
  - USAGE.md: Kapsamlı kullanım kılavuzu (Türkçe)

INSTRUCTIONS

echo -e "${YELLOW}💡 İpucu:${NC} UPX kurulu değilse: ${CYAN}sudo apt install upx${NC}"
echo -e "${YELLOW}💡 İpucu:${NC} Checksum doğrulama: ${CYAN}sha256sum -c ${BUILD_DIR}/SHA256SUMS${NC}\n"

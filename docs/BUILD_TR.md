# Build Kılavuzu - Hermes eUICC Manager

Tüm desteklenen platformlar ve mimariler için kapsamlı build talimatları.

## İçindekiler

- [Gereksinimler](#gereksinimler)
- [Hızlı Başlangıç](#hızlı-başlangıç)
- [Build Yöntemleri](#build-yöntemleri)
- [Platform-Spesifik Build](#platform-spesifik-build)
- [Optimizasyon Seviyeleri](#optimizasyon-seviyeleri)
- [Binary Boyut Optimizasyonu](#binary-boyut-optimizasyonu)
- [Sorun Giderme](#sorun-giderme)

## Gereksinimler

### Zorunlu

- **Go 1.21+**: <https://go.dev/dl/> adresinden indirin
- **Git**: Repository klonlamak için
- **Linux**: Birincil geliştirme platformu

### Opsiyonel

- **UPX**: Binary sıkıştırma için (sudo apt install upx)
- **Make**: Makefile hedeflerini kullanmak için

### Kurulum Doğrulama

```bash
go version    # go1.21 veya sonrası göstermeli
git --version
make --version  # Opsiyonel
upx --version   # Opsiyonel
```

## Hızlı Başlangıç

### Basit Build (Mevcut Platform)

```bash
cd app
go build -o hermes-euicc .
./hermes-euicc --version
```

### Optimize Build (Mevcut Platform)

```bash
cd app
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
```

Bayrakların açıklaması:

- `-ldflags="-s -w"`: Debug bilgilerini ve sembol tablosunu kaldır
- `-trimpath`: Dosya yolu bilgilerini kaldır

## Build Yöntemleri

### Yöntem 1: Makefile Kullanımı (Önerilen)

Repository kök dizininden:

```bash
# Tüm platformları build et
make all

# Belirli platformları build et
make mipsle
make armv7
make arm64
make amd64

# Sadece OpenWRT popüler platformları build et
make openwrt

# Build + sıkıştır + checksum
make all compress checksum

# Build dizinini temizle
make clean
```

### Yöntem 2: build-all.sh Kullanımı

app/ dizininden genel platform build'leri:

```bash
cd app
./build-all.sh
```

Çıktı dizini: build/

Oluşturulan binary'ler (toplam 8):

- hermes-euicc-mipsle (MIPS Little Endian)
- hermes-euicc-mips (MIPS Big Endian)
- hermes-euicc-armv5 (ARMv5)
- hermes-euicc-armv6 (ARMv6)
- hermes-euicc-armv7 (ARMv7)
- hermes-euicc-arm64 (ARM64)
- hermes-euicc-i386 (x86 32-bit)
- hermes-euicc-amd64 (x86 64-bit)

Özellikler:

- Renkli çıktı
- Otomatik UPX sıkıştırma (kuruluysa)
- SHA256 checksum oluşturma
- Tipik boyut: 6.5 MB → 2.8 MB (sıkıştırılmış)

### Yöntem 3: build-openwrt.sh Kullanımı

app/ dizininden cihaz-spesifik optimize build'ler:

```bash
cd app
./build-openwrt.sh
```

Çıktı dizini: build/openwrt/

Özellikler:

- 41 cihaz-spesifik binary
- CPU-spesifik optimizasyonlar
- MIPS için FPU tespiti
- Modern x86 için AVX2 desteği
- Build log'unda chipset bilgileri

## Platform-Spesifik Build

### MIPS Platformlar

#### MIPS Little Endian (mipsle)

TP-Link, GL.iNet, Xiaomi routerlar için en yaygın:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-mipsle .
```

Cihazlar: GL.iNet AR750S, TP-Link Archer C7/WR1043ND/WR841N, Xiaomi Mi Router 3G/4A, Ubiquiti EdgeRouter X

#### MIPS Big Endian (mips)

Eski Broadcom tabanlı routerlar:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-mips .
```

### ARM Platformlar

#### ARMv7 (Çoğu ARM router için önerilen)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-armv7 .
```

Özellikler: VFPv3 kayan nokta, NEON SIMD talimatları, ARMv5/6'dan %30-40 daha hızlı

Cihazlar: Raspberry Pi 2/3, GL.iNet B1300, Linksys WRT1900ACS, MikroTik hAP ac2

#### ARM64

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-arm64 .
```

Cihazlar: Raspberry Pi 4/5, GL.iNet MT6000 (Flint 2), NanoPi R4S, Banana Pi BPI-R3/R4

### x86 Platformlar

#### 64-bit (amd64)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" -trimpath \
    -o hermes-euicc-amd64 .
```

Cihazlar: PC Engines APU, Protectli Vault, Qotom Mini PC'ler, Genel x86-64 sunucular

## Optimizasyon Seviyeleri

### MIPS Optimizasyonları

**GOMIPS=softfloat (Varsayılan)** - Tüm MIPS cihazlar için güvenli

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build ...
```

**GOMIPS=hardfloat** - FPU donanımlı cihazlarda %20-30 daha hızlı

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=hardfloat go build ...
```

hardfloat kullanın: QCA9563 (74Kc), MT7621 (1004Kc)
softfloat kullanın: AR9341 (24Kc), QCA9531 (24Kc)

### x86 Optimizasyonları

**GOAMD64 Seviyeleri:**

```bash
GOAMD64=v1  # Baseline x86-64 (tüm işlemciler)
GOAMD64=v2  # +SSE4.2, POPCNT (2009+)
GOAMD64=v3  # +AVX2, BMI2 (2013+, %15-25 daha hızlı)
GOAMD64=v4  # +AVX512 (çok yeni, nadiren gerekli)
```

Örnekler:

```bash
# PC Engines APU (AMD GX-412TC, 2013)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v2 go build ...

# Qotom Q355G6 (Intel Celeron J3xxx)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build ...

# Genel (maksimum uyumluluk)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v1 go build ...
```

## Binary Boyut Optimizasyonu

### Varsayılan Boyut

Optimizasyon olmadan: ~12 MB

### Build Bayrakları İle

```bash
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
```

Sonuç: ~6.5 MB (%45 azalma)

### UPX Sıkıştırma İle

```bash
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
upx --best --lzma hermes-euicc
```

Sonuç: ~2.8 MB (6.5 MB'den %57 azalma)

### Boyut Karşılaştırması

| Yöntem | Boyut | Notlar |
|--------|------|--------|
| Varsayılan | 12 MB | Optimizasyon yok |
| -ldflags="-s -w" | 6.5 MB | Standart optimizasyon |
| + UPX --best | 2.8 MB | Önerilen |

## Sorun Giderme

### "go: command not found"

Go 1.21+ kurun:

```bash
wget https://go.dev/dl/go1.23.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Hedef Cihazda "Illegal instruction"

Yanlış optimizasyon seviyesi kullandınız. Güvenli varsayılanlarla yeniden build edin (MIPS için softfloat, x86 için v1).

### Binary Router İçin Çok Büyük

UPX sıkıştırma kullanın:

```bash
upx --best --lzma hermes-euicc
```

### Build Script İzin Hatası

Scriptleri çalıştırılabilir yapın:

```bash
chmod +x build-all.sh build-openwrt.sh
```

## Doğrulama

### Binary Mimarisini Kontrol Et

```bash
file hermes-euicc-mipsle
# Çıktı: ELF 32-bit LSB executable, MIPS, MIPS32 rel2 version 1...
```

### SHA256 Doğrula

```bash
cd build
sha256sum -c SHA256SUMS
```

### Binary'yi Test Et

```bash
./hermes-euicc --version

# Hedef cihazda
scp hermes-euicc root@192.168.1.1:/tmp/
ssh root@192.168.1.1
chmod +x /tmp/hermes-euicc
/tmp/hermes-euicc --version
```

## Özet

**Hızlı geliştirme için:**

```bash
go build -o hermes-euicc .
```

**Production için (tek platform):**

```bash
go build -ldflags="-s -w" -trimpath -o hermes-euicc .
upx --best --lzma hermes-euicc
```

**Dağıtım için (tüm platformlar):**

```bash
cd app
./build-all.sh
```

**Cihaz-spesifik için (optimize):**

```bash
cd app
./build-openwrt.sh
```

**Makefile kullanımı (önerilen):**

```bash
make all compress checksum
```

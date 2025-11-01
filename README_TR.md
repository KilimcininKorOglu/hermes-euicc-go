# Hermes eUICC Manager

Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>

**Hermes** - JSON çıktı veren tam özellikli eUICC (eSIM) yönetim CLI uygulaması.

## Özellikler

- ✅ **Tam JSON Çıktı** - Tüm komutlar JSON formatında yanıt verir
- ✅ **Otomatik Driver Tespiti** - QMI, MBIM, AT, CCID driver'larını otomatik algılar*
- ℹ️ **Platform Notu** - CCID desteği sadece amd64/arm64 platformlarında (MIPS/i386/ARMv5-7'de yok)
- ✅ **Tüm SGP.22 Fonksiyonları** - Kütüphanedeki tüm özellikleri destekler
- ✅ **Hata Yönetimi** - Yapılandırılmış hata mesajları
- ✅ **Komut Satırı Bayrakları** - Esnek konfigürasyon seçenekleri

## Kurulum

```bash
cd app
go build -o hermes-euicc
```

## Kullanım

### Global Seçenekler

```bash
-device string      # Cihaz yolu (örn: /dev/cdc-wdm0, /dev/ttyUSB2)
-driver string      # Driver türü: qmi, mbim, at, ccid (belirtilmezse otomatik)
-slot int           # SIM slot numarası (varsayılan: 1)
-timeout int        # HTTP timeout saniye cinsinden (varsayılan: 30)
-verbose            # Detaylı log çıktısı
```

### Komutlar

#### 0. Yardım ve Versiyon

```bash
# Yardım mesajı
./hermes-euicc help

# Versiyon bilgisi
./hermes-euicc version
```

**Version Çıktı:**

```json
{
  "success": true,
  "data": {
    "version": "1.0.0",
    "copyright": "Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>",
    "license": "MIT"
  }
}
```

#### 1. EID Öğrenme

```bash
./hermes-euicc eid
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "eid": "89049032003451234567890123456789"
  }
}
```

#### 2. eUICC Bilgileri

```bash
./hermes-euicc info
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "eid": "89049032003451234567890123456789",
    "euicc_info1": "...",
    "euicc_info2": "..."
  }
}
```

#### 3. Profil Listesi

```bash
./hermes-euicc list
```

**Çıktı:**

```json
{
  "success": true,
  "data": [
    {
      "iccid": "8944476500001224158",
      "isdp_aid": "A0000005591010FFFFFFFF8900000100",
      "profile_state": 1,
      "profile_name": "Personal",
      "profile_nickname": "My SIM",
      "service_provider_name": "Vodafone",
      "profile_class": "operational",
      "icon": "...",
      "icon_file_type": "image/png"
    }
  ]
}
```

#### 4. Profil Aktifleştirme

```bash
./hermes-euicc enable 8944476500001224158
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "message": "profile enabled successfully",
    "iccid": "8944476500001224158"
  }
}
```

#### 5. Profil Pasifleştirme

```bash
./hermes-euicc disable 8944476500001224158
```

#### 6. Profil Silme

```bash
./hermes-euicc delete 8944476500001224158
```

#### 7. Takma Ad Ayarlama

```bash
./hermes-euicc nickname 8944476500001224158 "Work SIM"
```

#### 8. Profil İndirme

```bash
./hermes-euicc download \
  --code "LPA:1$smdp.io$QR-G-5C-1LS" \
  --imei "356938035643809" \
  --confirm
```

**Seçenekler:**

- `--code`: Aktivasyon kodu (zorunlu)
- `--imei`: IMEI numarası (opsiyonel)
- `--confirmation-code`: Onay kodu (gerekirse)
- `--confirm`: Otomatik onaylama

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "isdp_aid": "A0000005591010FFFFFFFF8900000100",
    "notification": 1
  }
}
```

#### 9. Profil Keşfi

```bash
./hermes-euicc discovery --imei "356938035643809"
```

**Seçenekler:**

- `--server`: SM-DS sunucu URL'i (opsiyonel)
- `--imei`: IMEI numarası (opsiyonel)

**Çıktı:**

```json
{
  "success": true,
  "data": [
    {
      "event_id": "EVT-123456",
      "address": "LPA:1$smdp.io$QR-G-5C-1LS"
    }
  ]
}
```

#### 10. Bildirimleri Listeleme

```bash
./hermes-euicc notifications
```

**Çıktı:**

```json
{
  "success": true,
  "data": [
    {
      "sequence_number": 1,
      "profile_management_operation": 1,
      "address": "smdp.example.com",
      "iccid": "8944476500001224158"
    }
  ]
}
```

#### 11. Bildirim Kaldırma

```bash
./hermes-euicc notification-remove 1
```

#### 12. Bildirim İşleme

```bash
./hermes-euicc notification-handle 1
```

#### 13. Otomatik Bildirim İşleme

Bekleyen tüm bildirimleri otomatik olarak eşzamanlı işle.

```bash
./hermes-euicc auto-notification
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "total_notifications": 3,
    "processed": 2,
    "failed": 1,
    "processed_notifications": [
      {
        "sequence_number": 1,
        "profile_management_operation": "install"
      },
      {
        "sequence_number": 2,
        "profile_management_operation": "enable"
      }
    ],
    "failed_notifications": [
      {
        "sequence_number": 3,
        "error": "notification retrieve failed"
      }
    ]
  }
}
```

#### 14. Yapılandırılmış Adresler

```bash
./hermes-euicc configured-addresses
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "default_smdp_address": "smdp.example.com",
    "root_smds_address": "smds.example.com"
  }
}
```

#### 15. Varsayılan SM-DP+ Adresi Ayarlama

```bash
./hermes-euicc set-default-dp "smdp.example.com"
```

#### 16. eUICC Challenge

```bash
./hermes-euicc challenge
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "challenge": "1234567890abcdef"
  }
}
```

#### 17. Bellek Sıfırlama

```bash
./hermes-euicc memory-reset
```

## Hata Yönetimi

Tüm hatalar JSON formatında döner:

```json
{
  "success": false,
  "error": "failed to initialize driver: no compatible driver found"
}
```

## Driver Seçimi

### Otomatik Tespit

Uygulama aşağıdaki sırayla driver'ları test eder:

1. QMI (`/dev/cdc-wdm0`)
2. MBIM (`/dev/cdc-wdm0`)
3. AT (`/dev/ttyUSB2`, `/dev/ttyUSB3`, `/dev/ttyUSB1`)
4. CCID (USB smart card reader)

### Manuel Seçim

```bash
# QMI
./hermes-euicc --driver qmi --device /dev/cdc-wdm0 list

# MBIM
./hermes-euicc --driver mbim --device /dev/cdc-wdm0 list

# AT
./hermes-euicc --driver at --device /dev/ttyUSB2 list

# CCID
./hermes-euicc --driver ccid list
```

## Script ile Kullanım

### Bash Örneği

```bash
#!/bin/bash

# Profilleri JSON olarak al
PROFILES=$(./hermes-euicc list)

# jq ile parse et
ACTIVE_ICCID=$(echo "$PROFILES" | jq -r '.data[] | select(.profile_state == 1) | .iccid')

echo "Active profile ICCID: $ACTIVE_ICCID"
```

### Python Örneği

```python
import json
import subprocess

# Profil listesini al
result = subprocess.run(
    ["./hermes-euicc", "list"],
    capture_output=True,
    text=True
)

data = json.loads(result.stdout)

if data["success"]:
    for profile in data["data"]:
        print(f"Profile: {profile['profile_name']} - State: {profile['profile_state']}")
else:
    print(f"Error: {data['error']}")
```

## Örnekler

### Senaryo 1: Otomatik Profil Yönetimi

```bash
#!/bin/bash

# Tüm profilleri listele
PROFILES=$(./hermes-euicc list)

# Pasif profilleri bul
DISABLED=$(echo "$PROFILES" | jq -r '.data[] | select(.profile_state == 0) | .iccid')

# İlk pasif profili aktifleştir
FIRST_ICCID=$(echo "$DISABLED" | head -1)
./hermes-euicc enable "$FIRST_ICCID"
```

### Senaryo 2: Profil Keşfi ve İndirme

```bash
#!/bin/bash

# Profilleri keşfet
DISCOVERY=$(./hermes-euicc discovery --imei "356938035643809")

# İlk profili indir
CODE=$(echo "$DISCOVERY" | jq -r '.data[0].address')
./hermes-euicc download --code "$CODE" --confirm
```

### Senaryo 3: Bildirim İşleme

**Basit yaklaşım (önerilen):**

```bash
#!/bin/bash

# Tüm bekleyen bildirimleri otomatik olarak eşzamanlı işle
./hermes-euicc auto-notification
```

**Manuel yaklaşım:**

```bash
#!/bin/bash

# Bekleyen bildirimleri al
NOTIFICATIONS=$(./hermes-euicc notifications)

# Her bildirimi işle
echo "$NOTIFICATIONS" | jq -r '.data[].sequence_number' | while read seq; do
    ./hermes-euicc notification-handle "$seq"
done
```

## Profil State Değerleri

- `0`: Disabled (Pasif)
- `1`: Enabled (Aktif)

## Profile Management Operation Değerleri

- `0`: Install
- `1`: Enable
- `2`: Disable
- `3`: Delete

## Profile Class Değerleri

- `test`: Test profili
- `provisioning`: Provisioning profili
- `operational`: Operasyonel profil

## Lisans

MIT License - Detaylar için LICENSE dosyasına bakın.

## İletişim

Kilimcinin Kör Oğlu <k@keremgok.tr>

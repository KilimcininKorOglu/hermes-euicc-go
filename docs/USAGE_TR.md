# Hermes eUICC Manager - Kullanım Kılavuzu

Hermes eUICC Manager için kapsamlı kullanım dokümantasyonu. eUICC özellikli cihazlarda eSIM profil yönetimi için JSON tabanlı CLI uygulaması.

## İçindekiler

- [Genel Bakış](#genel-bakış)
- [Kurulum](#kurulum)
- [Global Seçenekler](#global-seçenekler)
- [Sürücü Seçimi](#sürücü-seçimi)
- [Komut Referansı](#komut-referansı)
- [JSON Çıktı Formatı](#json-çıktı-formatı)
- [Yaygın Kullanım Senaryoları](#yaygın-kullanım-senaryoları)
- [Sorun Giderme](#sorun-giderme)

## Genel Bakış

Hermes eUICC Manager, eUICC özellikli modem ve cihazlarda eSIM profilleri üzerinde tam kontrol sağlayan bir komut satırı aracıdır. Tüm işlemler betik ve otomasyon sistemleriyle kolay entegrasyon için yapılandırılmış JSON çıktısı döndürür.

**Ana Özellikler:**

- Otomasyon için yalnızca JSON çıktısı
- Donanım sürücülerinin otomatik tespiti (QMI, MBIM, AT, CCID)
- Tam SGP.22 protokol implementasyonu
- Tüm eSIM işlemlerini kapsayan 17 komut
- Çoklu platform desteği (MIPS, ARM, x86)

## Kurulum

### Kaynaktan Derleme

```bash
cd app
go build -o hermes-euicc .
```

### Kurulumu Doğrulama

```bash
hermes-euicc version
```

## Global Seçenekler

Bu seçenekler tüm komutlara uygulanır:

### -device string

Modem/okuyucu için cihaz yolu. Belirtilmezse yaygın yollar otomatik test edilir.

```bash
# QMI/MBIM modem
hermes-euicc -device /dev/cdc-wdm0 list

# AT modem
hermes-euicc -device /dev/ttyUSB2 list

# Otomatik tespit (varsayılan)
hermes-euicc list
```

### -driver string

Otomatik tespit yerine belirli bir sürücüyü zorla.

Desteklenen sürücüler:

- `qmi` - Qualcomm MSM Interface (en yaygın)
- `mbim` - Mobile Broadband Interface Model
- `at` - AT komutları (seri modemler)
- `ccid` - USB akıllı kart okuyucular (sadece Linux amd64/arm64)

```bash
hermes-euicc -driver qmi list
hermes-euicc -driver mbim list
hermes-euicc -driver at -device /dev/ttyUSB2 list
hermes-euicc -driver ccid list
```

### -slot int

Çoklu SIM cihazlarda SIM slot numarası (varsayılan: 1).

```bash
# SIM slot 2 kullan
hermes-euicc -slot 2 list
```

### -timeout int

SM-DP+ işlemleri için HTTP timeout saniye cinsinden (varsayılan: 30).

```bash
# Timeout'u 60 saniyeye çıkar
hermes-euicc -timeout 60 download --code "LPA:..."
```

### -verbose

Hata ayıklama için detaylı loglama aktif et.

```bash
hermes-euicc -verbose list
```

## Sürücü Seçimi

### Otomatik Tespit (Önerilen)

Sürücü belirtmeden komutları çalıştırın:

```bash
hermes-euicc list
```

Tespit sırası:

1. QMI: `/dev/cdc-wdm0`
2. MBIM: `/dev/cdc-wdm0`
3. AT: `/dev/ttyUSB2`, `/dev/ttyUSB3`, `/dev/ttyUSB1`
4. CCID (USB akıllı kart okuyucu)

### Manuel Seçim

Otomatik tespit başarısız olduğunda veya belirli sürücü gerektiğinde:

```bash
# QMI'yi zorla
hermes-euicc -driver qmi -device /dev/cdc-wdm0 list

# Özel cihaz yolu ile AT
hermes-euicc -driver at -device /dev/ttyUSB0 list
```

### Platform Bazında Sürücü Uygunluğu

| Platform | QMI | MBIM | AT | CCID |
|----------|-----|------|----|------|
| Linux amd64 | ✓ | ✓ | ✓ | ✓ |
| Linux arm64 | ✓ | ✓ | ✓ | ✓ |
| Linux MIPS | ✓ | ✓ | ✓ | ✗ |
| Linux ARMv7 | ✓ | ✓ | ✓ | ✗ |

## Komut Referansı

### help - Yardım Mesajı

Kullanım bilgisini görüntüle.

```bash
hermes-euicc help
```

### version - Versiyon Bilgisi

Uygulama versiyonu ve telif hakkını göster.

```bash
hermes-euicc version
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "name": "Hermes eUICC Manager",
    "version": "1.0.0",
    "copyright": "Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>",
    "license": "MIT"
  }
}
```

### eid - EID Al

eUICC Tanımlayıcısını (32 karakterlik hexadecimal) al.

```bash
hermes-euicc eid
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

### info - eUICC Bilgilerini Al

EID, EUICCInfo1 ve EUICCInfo2 dahil kapsamlı eUICC bilgilerini al.

```bash
hermes-euicc info
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "eid": "89049032003451234567890123456789",
    "euicc_info1": "a00b...",
    "euicc_info2": "bf2e..."
  }
}
```

### list - Profilleri Listele

eUICC'ye yüklenmiş tüm eSIM profillerini listele.

```bash
hermes-euicc list
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
      "icon": "base64-encoded-png",
      "icon_file_type": "image/png"
    }
  ]
}
```

**Profil Durumları:**

- `0` - Devre dışı
- `1` - Etkin (aktif)

**Profil Sınıfları:**

- `test` - Test profili
- `provisioning` - Provisioning profili
- `operational` - Normal operasyonel profil

### enable - Profil Aktifleştir

ICCID ile bir profili aktif et. Aynı anda sadece bir profil aktif olabilir.

```bash
hermes-euicc enable 8944476500001224158
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

**Not:** Bir profili aktifleştirmek mevcut aktif profili otomatik olarak devre dışı bırakır.

### disable - Profil Devre Dışı Bırak

ICCID ile bir profili deaktif et.

```bash
hermes-euicc disable 8944476500001224158
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "message": "profile disabled successfully",
    "iccid": "8944476500001224158"
  }
}
```

### delete - Profil Sil

eUICC'den bir profili kalıcı olarak kaldır.

```bash
hermes-euicc delete 8944476500001224158
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "message": "profile deleted successfully",
    "iccid": "8944476500001224158"
  }
}
```

**Uyarı:** Bu işlem geri alınamaz. Profil silinmeden önce devre dışı bırakılmalıdır.

### nickname - Profil Takma Adı Ayarla

Bir profile özel takma ad belirle.

```bash
hermes-euicc nickname 8944476500001224158 "İş SIM"
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "message": "nickname set successfully",
    "iccid": "8944476500001224158",
    "nickname": "İş SIM"
  }
}
```

### download - Profil İndir

Aktivasyon kodu kullanarak yeni eSIM profili indir ve yükle.

**Seçenekler:**

- `--code` (zorunlu) - LPA aktivasyon kodu
- `--imei` (opsiyonel) - Cihaz IMEI
- `--confirmation-code` (opsiyonel) - Gerekirse profil onay kodu
- `--confirm` (opsiyonel) - Sormadan otomatik onayla

```bash
# Otomatik onay ile basit indirme
hermes-euicc download \
  --code "LPA:1$smdp.io$MATCHING-ID" \
  --confirm

# IMEI ve onay kodu ile indirme
hermes-euicc download \
  --code "LPA:1$smdp.io$MATCHING-ID" \
  --imei "356938035643809" \
  --confirmation-code "1234" \
  --confirm
```

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

**Aktivasyon Kodu Formatı:**

```
LPA:1$<smdp-adresi>$<eşleşme-id>
```

Örnek: `LPA:1$smdp.io$QR-G-5C-1LS`

### discovery - Profil Keşfi

Mevcut profil indirmeleri için SM-DS sunucularını sorgula.

**Seçenekler:**

- `--server` (opsiyonel) - Özel SM-DS sunucusu (varsayılan: birden fazla sunucu dener)
- `--imei` (opsiyonel) - Cihaz IMEI

```bash
# IMEI ile keşif
hermes-euicc discovery --imei "356938035643809"

# Belirli sunucudan keşif
hermes-euicc discovery --server "lpa.ds.gsma.com" --imei "356938035643809"
```

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

### notifications - Bildirimleri Listele

eUICC'den bekleyen bildirimleri al.

```bash
hermes-euicc notifications
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

**Profil Yönetim İşlemleri:**

- `0` - Yükleme
- `1` - Aktifleştirme
- `2` - Devre dışı bırakma
- `3` - Silme

### notification-remove - Bildirim Kaldır

Sıra numarasına göre bildirimi sil.

```bash
hermes-euicc notification-remove 1
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "message": "notification removed successfully",
    "sequence_number": 1
  }
}
```

### notification-handle - Bildirim İşle

Bildirimi işle ve SM-DP+ sunucusuna gönder.

```bash
hermes-euicc notification-handle 1
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "message": "notification handled successfully",
    "sequence_number": 1
  }
}
```

### auto-notification - Tüm Bildirimleri Otomatik İşle

Bekleyen tüm bildirimleri otomatik olarak al ve eşzamanlı olarak işle. Bu komut tüm bekleyen bildirimleri listeler, her birini alır ve optimal performans için paralel olarak işler.

```bash
hermes-euicc auto-notification
```

**Çıktı (bildirimler var):**

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

**Çıktı (bildirim yok):**

```json
{
  "success": true,
  "data": {
    "total_notifications": 0,
    "processed": 0,
    "failed": 0,
    "processed_notifications": [],
    "failed_notifications": []
  }
}
```

**Notlar:**
- Daha iyi performans için bildirimleri eşzamanlı işler
- Başarılı ve başarısız işlemler için detaylı sonuçlar döndürür
- Bekleyen bildirimlerin otomatik toplu işlenmesi için kullanışlıdır

### configured-addresses - Yapılandırılmış Adresler

eUICC'de yapılandırılmış varsayılan SM-DP+ ve root SM-DS adreslerini al.

```bash
hermes-euicc configured-addresses
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

### set-default-dp - Varsayılan SM-DP+ Adresini Ayarla

Varsayılan SM-DP+ sunucu adresini yapılandır.

```bash
hermes-euicc set-default-dp "smdp.example.com"
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "message": "default DP address set successfully",
    "address": "smdp.example.com"
  }
}
```

### challenge - eUICC Challenge Al

eUICC'den kimlik doğrulama challenge'ı al.

```bash
hermes-euicc challenge
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

### memory-reset - eUICC Belleği Sıfırla

eUICC operasyonel belleği sıfırla (profilleri silmez).

```bash
hermes-euicc memory-reset
```

**Çıktı:**

```json
{
  "success": true,
  "data": {
    "message": "memory reset successfully"
  }
}
```

**Uyarı:** Bu işlem bildirim listelerini ve geçici verileri sıfırlar.

## JSON Çıktı Formatı

Tüm komutlar tutarlı formatta JSON döndürür:

### Başarılı Yanıt

```json
{
  "success": true,
  "data": {
    // Komuta özgü veri
  }
}
```

### Hata Yanıtı

```json
{
  "success": false,
  "error": "hata mesajı açıklaması"
}
```

### Parse Örnekleri

**Bash ile jq:**

```bash
# EID al
EID=$(hermes-euicc eid | jq -r '.data.eid')
echo "EID: $EID"

# Aktif profil ICCID'sini al
ACTIVE_ICCID=$(hermes-euicc list | jq -r '.data[] | select(.profile_state == 1) | .iccid')
echo "Aktif: $ACTIVE_ICCID"
```

**Python:**

```python
import json
import subprocess

result = subprocess.run(['hermes-euicc', 'list'], capture_output=True, text=True)
data = json.loads(result.stdout)

if data['success']:
    for profile in data['data']:
        print(f"{profile['profile_name']}: {profile['iccid']}")
```

## Yaygın Kullanım Senaryoları

### Profiller Arası Geçiş

```bash
#!/bin/bash

# Profilleri listele ve seç
PROFILES=$(hermes-euicc list)
echo "$PROFILES" | jq -r '.data[] | "\(.iccid) - \(.profile_name) (durum: \(.profile_state))"'

# Mevcut profili devre dışı bırak
CURRENT=$(echo "$PROFILES" | jq -r '.data[] | select(.profile_state == 1) | .iccid')
hermes-euicc disable "$CURRENT"

# Yeni profili aktif et
hermes-euicc enable "8944476500001224159"
```

### Otomatik Keşif ve İndirme

```bash
#!/bin/bash

# Mevcut profilleri keşfet
DISCOVERY=$(hermes-euicc discovery --imei "356938035643809")

# İlk aktivasyon kodunu al
CODE=$(echo "$DISCOVERY" | jq -r '.data[0].address')

# Profili indir
hermes-euicc download --code "$CODE" --confirm
```

### Toplu Profil Yönetimi

```bash
#!/bin/bash

# Tüm devre dışı profilleri listele
DISABLED=$(hermes-euicc list | jq -r '.data[] | select(.profile_state == 0) | .iccid')

# Tüm devre dışı profilleri sil
for ICCID in $DISABLED; do
    echo "$ICCID siliniyor..."
    hermes-euicc delete "$ICCID"
done
```

### Bildirim İşleme

```bash
#!/bin/bash

# Bekleyen bildirimleri al
NOTIFICATIONS=$(hermes-euicc notifications)

# Her bildirimi işle
echo "$NOTIFICATIONS" | jq -r '.data[].sequence_number' | while read SEQ; do
    echo "$SEQ numaralı bildirim işleniyor..."
    hermes-euicc notification-handle "$SEQ"
done
```

### Yeniden Deneme ile Profil İndirme

```bash
#!/bin/bash

CODE="LPA:1$smdp.io$MATCHING-ID"
MAX_RETRIES=3

for i in $(seq 1 $MAX_RETRIES); do
    echo "$MAX_RETRIES denemeden $i. deneme..."
    
    RESULT=$(hermes-euicc download --code "$CODE" --confirm 2>&1)
    
    if echo "$RESULT" | jq -e '.success' > /dev/null 2>&1; then
        echo "İndirme başarılı!"
        echo "$RESULT" | jq '.'
        exit 0
    fi
    
    echo "Başarısız: $(echo "$RESULT" | jq -r '.error')"
    sleep 5
done

echo "$MAX_RETRIES denemeden sonra indirme başarısız"
exit 1
```

## Sorun Giderme

### Cihaz Bulunamadı

**Hata:**

```json
{
  "success": false,
  "error": "failed to initialize driver: no compatible driver found"
}
```

**Çözümler:**

1. Cihaz yolunu kontrol edin:

```bash
ls -la /dev/cdc-wdm* /dev/ttyUSB*
```

2. İzinleri kontrol edin:

```bash
sudo chmod 666 /dev/cdc-wdm0
```

3. Manuel sürücü seçimi deneyin:

```bash
hermes-euicc -driver qmi -device /dev/cdc-wdm0 list
```

### CCID Sürücüsü Mevcut Değil

**Hata:**

```json
{
  "success": false,
  "error": "CCID driver not supported on this platform (requires amd64/arm64 + linux)"
}
```

**Çözüm:** CCID sadece Linux amd64/arm64'te mevcuttur. Bunun yerine QMI/MBIM/AT sürücülerini kullanın.

### Geçersiz ICCID

**Hata:**

```json
{
  "success": false,
  "error": "invalid ICCID: ..."
}
```

**Çözüm:** ICCID'nin tam olarak 19-20 rakam olduğundan emin olun. Doğru ICCID'yi `list` komutundan alın.

### Profil Zaten Etkin

**Hata:**

```json
{
  "success": false,
  "error": "profile already enabled"
}
```

**Çözüm:** Aynı anda sadece bir profil aktif olabilir. Önce mevcut profili devre dışı bırakın, veya enable komutu bunu otomatik yapacaktır.

### İndirme Onay Gerektiriyor

**Hata:**

```json
{
  "success": false,
  "error": "download requires user confirmation"
}
```

**Çözüm:** `--confirm` bayrağı ekleyin veya onay kodu sağlayın:

```bash
hermes-euicc download --code "LPA:..." --confirm
```

### İndirme Sırasında Timeout

**Hata:**

```json
{
  "success": false,
  "error": "context deadline exceeded"
}
```

**Çözüm:** Timeout'u artırın:

```bash
hermes-euicc -timeout 60 download --code "LPA:..." --confirm
```

### İzin Reddedildi

**Hata:**

```bash
permission denied: /dev/cdc-wdm0
```

**Çözümler:**

1. sudo ile çalıştırın:

```bash
sudo hermes-euicc list
```

2. Kullanıcıyı dialout grubuna ekleyin:

```bash
sudo usermod -aG dialout $USER
# Çıkış yapıp tekrar giriş yapın
```

3. Cihaz izinlerini ayarlayın:

```bash
sudo chmod 666 /dev/cdc-wdm0
```

## Entegrasyon Örnekleri

### Systemd Servisi

```ini
[Unit]
Description=eSIM Profil Değiştirici
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/hermes-euicc enable 8944476500001224158
User=root

[Install]
WantedBy=multi-user.target
```

### Cron Görevi

```bash
# Her gün sabah 6'da profil değiştir
0 6 * * * /usr/local/bin/hermes-euicc enable 8944476500001224158 >> /var/log/esim-switch.log 2>&1
```

### OpenWRT Init Script

```bash
#!/bin/sh /etc/rc.common

START=99

start() {
    /usr/bin/hermes-euicc enable 8944476500001224158
}
```

## Ek Kaynaklar

- **Derleme Kılavuzu:** [BUILD.md](BUILD.md)
- **Sürücü Dokümantasyonu:** [DRIVERS.md](DRIVERS.md)
- **Proje README:** [../README.md](../README.md)
- **GSMA SGP.22 Spesifikasyonu:** <https://aka.pw/sgp22/v2.5>

## Destek

Sorunlar, sorular veya katkılar için:

- GitHub Issues: <https://github.com/KilimcininKorOglu/euicc-go/issues>
- Dokümantasyon: Repository root'taki `/docs` dizinine bakın

## Lisans

MIT License - Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>

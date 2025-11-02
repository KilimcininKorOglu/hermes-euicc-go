# Hermes eUICC Manager Yeni Özellikler

Bu dokümant, Hermes eUICC Manager CLI uygulamasına eklenen ve euicc-go kütüphanesi v1.2.1'deki en son API'leri entegre eden yeni özellikleri açıklar.

**Son Güncelleme:** 2025-11-02

---

## Genel Bakış

CLI uygulamasına, euicc-go kütüphanesindeki yeni API'lerden yararlanan üç büyük özellik seti eklenmiştir:

1. **Geliştirilmiş Bildirim İşleme** - Otomatik toplu bildirim yönetimi
2. **Profil Keşfi ve İndirme** - SM-DS profil keşfi ve tek adımda indirme
3. **Chip Bilgileri** - Parse edilmiş detaylı eUICC chip bilgileri

Tüm özellikler kolay otomasyon ve script yazımı için JSON-only çıktı formatını korur.

---

## 1. Bildirim İşleme Özellikleri

### notification-process

Otomatik kaldırma ve hata yönetimi ile belirli bildirimleri sequence number'a göre işleyin.

**Kullanım:**
```bash
hermes-euicc notification-process <sequence_number> [<sequence_number> ...]
```

**Örnekler:**
```bash
# Tek bildirim işle
hermes-euicc notification-process 1

# Birden fazla bildirim işle
hermes-euicc notification-process 1 2 3 5
```

**Çıktı:**
```json
{
  "success": true,
  "data": {
    "message": "notification processing completed",
    "total": 3,
    "processed": 2,
    "failed": 1,
    "processed_list": [
      {
        "sequence_number": 1,
        "removed": true
      },
      {
        "sequence_number": 2,
        "removed": true
      }
    ],
    "failed_list": [
      {
        "sequence_number": 3,
        "error": "notification retrieve failed"
      }
    ]
  }
}
```

**Özellikler:**
- Tek komutta bir veya birden fazla bildirimi işleyin
- Başarılı işlemden sonra otomatik kaldırma
- Biri başarısız olsa bile işlemeye devam eder
- Bildirim başına detaylı sonuçlar
- Her bildirim için kaldırma durumunu gösterir

**Kullanım Alanları:**
- Seçici bildirim işleme
- Başarısız bildirimleri yeniden deneme
- Bildirimleri belirli sırada işleme
- Bildirim yönetimi üzerinde hassas kontrol

### auto-notification (Geliştirilmiş)

Mevcut `auto-notification` komutu, daha iyi performans ve güvenilirlik için kütüphanenin `ProcessAllNotifications()` fonksiyonunu kullanacak şekilde yükseltildi.

**Kullanım:**
```bash
hermes-euicc auto-notification
```

**İyileştirmeler:**
- Kütüphanenin optimize edilmiş işlemesini kullanır
- İşlenen bildirimlerin otomatik kaldırılması
- Bildirim başına daha iyi hata yönetimi
- Daha temiz implementasyon (~50 satır kod azaltıldı)

---

## 2. Profil Keşfi Özellikleri

### discovery (Geliştirilmiş)

Mevcut `discovery` komutu, kütüphanenin `DiscoverProfiles()` fonksiyonunu kullanacak şekilde yükseltildi.

**Kullanım:**
```bash
# Varsayılan GSMA SM-DS'den keşfet
hermes-euicc discovery

# Özel SM-DS sunucusundan keşfet
hermes-euicc discovery --server prod.smds.rsp.goog

# IMEI kimlik doğrulaması ile keşfet
hermes-euicc discovery --imei 123456789012345
```

**İyileştirmeler:**
- Kütüphanenin güçlü keşif implementasyonunu kullanır
- Basitleştirilmiş kod (~15 satır azaltıldı)
- Daha iyi hata yönetimi
- Manuel concurrent işlemeye gerek yok

### discover-download (Yeni)

SM-DS'den ilk uygun profilin tek adımda keşfi ve indirilmesi.

**Kullanım:**
```bash
# Varsayılan SM-DS'den keşfet ve indir
hermes-euicc discover-download

# Özel SM-DS'den keşfet ve indir
hermes-euicc discover-download --server prod.smds.rsp.goog

# IMEI kimlik doğrulaması ile
hermes-euicc discover-download --imei 123456789012345
```

**Çıktı (profil bulundu ve indirildi):**
```json
{
  "success": true,
  "data": {
    "message": "profile downloaded successfully"
  }
}
```

**Çıktı (profil yok):**
```json
{
  "success": true,
  "data": {
    "message": "no profiles available for download"
  }
}
```

**Özellikler:**
- Tek komutta otomatik keşif ve indirme
- Belirtilmezse varsayılan GSMA SM-DS sunucusu kullanılır
- Opsiyonel IMEI kimlik doğrulaması
- Özel SM-DS sunucu desteği
- Hata toleranslı (profil bulunamazsa bile success döner)

**Kullanım Alanları:**
- Otomatik provisioning scriptleri
- İlk cihaz kurulumu
- Aktivasyon kodu olmadan hızlı profil yükleme
- Sıfır dokunmatik provisioning iş akışları

**Manuel yaklaşımla karşılaştırma:**

**Eski yol:**
```bash
# Adım 1: Keşfet
PROFILES=$(hermes-euicc discovery)

# Adım 2: İlk profil adresini çıkar
CODE=$(echo "$PROFILES" | jq -r '.data[0].address')

# Adım 3: İndir
hermes-euicc download --code "$CODE" --confirm
```

**Yeni yol:**
```bash
# Tek komut her şeyi yapar
hermes-euicc discover-download
```

---

## 3. Chip Bilgileri Özelliği

### chip-info (Yeni)

Bellek, yetenekler, versiyonlar ve yapılandırma dahil detaylı, parse edilmiş chip bilgilerini alın.

**Kullanım:**
```bash
hermes-euicc chip-info
```

**Çıktı:**
```json
{
  "success": true,
  "data": {
    "eid": "89033023426200000000000123456789",
    "configured_addresses": {
      "default_smdp_address": "smdp.example.com",
      "root_smds_address": "smds.example.com"
    },
    "euicc_info2": {
      "profile_version": "2.3",
      "svn": "2.5.0",
      "euicc_firmware_ver": "1.2.3",
      "ts102241_version": "11.0.0",
      "global_platform_version": "2.3.1",
      "pp_version": "0201",
      "ext_card_resource": {
        "installed_application": 3,
        "free_non_volatile_memory": 524288,
        "free_volatile_memory": 16384
      },
      "uicc_capability": [
        "contactlessSupport",
        "usimSupport",
        "isimSupport",
        "javacard",
        "multipleUsimSupport"
      ],
      "rsp_capability": [
        "additionalProfile",
        "crlSupport",
        "rpmSupport"
      ],
      "euicc_ci_pkid_list_for_verification": [
        "ES_CERT_1",
        "ES_CERT_2"
      ],
      "euicc_ci_pkid_list_for_signing": [
        "ES_CERT_3"
      ],
      "forbidden_profile_policy_rules": [],
      "euicc_category": "basicEuicc",
      "sas_accreditation_number": "SAS-001",
      "certification_data_object": {
        "platform_label": "Platform v1.0",
        "discovery_base_url": "https://discovery.example.com"
      }
    },
    "rules_authorisation_table": [
      {
        "ppr_ids": ["ppr1", "ppr2"],
        "allowed_operators": [
          {
            "plmn": "310260",
            "gid1": "A1",
            "gid2": "B2"
          }
        ]
      }
    ]
  }
}
```

**Anahtar Bilgi Alanları:**

#### EID
- Benzersiz chip tanımlayıcısı (32 hex karakter)

#### Configured Addresses
- `default_smdp_address` - Varsayılan SM-DP+ sunucusu
- `root_smds_address` - Root SM-DS sunucusu

#### EUICCInfo2

**Versiyon Bilgileri:**
- `profile_version` - Profil spesifikasyon versiyonu
- `svn` - SGP.22 spesifikasyon versiyonu
- `euicc_firmware_ver` - Chip firmware versiyonu
- `ts102241_version` - JavaCard/ETSI versiyonu
- `global_platform_version` - GlobalPlatform versiyonu
- `pp_version` - Protection Profile versiyonu

**Bellek/Depolama (ÖNEMLİ):**
- `installed_application` - Yüklü uygulama sayısı
- `free_non_volatile_memory` - Mevcut kalıcı depolama (byte)
- `free_volatile_memory` - Mevcut RAM (byte)

**Yetenekler:**
- `uicc_capability` - Kart yetenekleri (USIM, ISIM, JavaCard, vb.)
- `rsp_capability` - Remote SIM Provisioning yetenekleri

**Güvenlik:**
- `euicc_ci_pkid_list_for_verification` - Doğrulama için public key ID'leri
- `euicc_ci_pkid_list_for_signing` - İmzalama için public key ID'leri
- `forbidden_profile_policy_rules` - Yasaklı policy kuralları

**Sınıflandırma:**
- `euicc_category` - Kategori (basicEuicc, mediumEuicc, contactlessEuicc, other)

**Sertifikasyon:**
- `sas_accreditation_number` - SAS akreditasyon numarası
- `certification_data_object` - Platform ve keşif bilgisi

#### Rules Authorisation Table
- `ppr_ids` - Profile Policy Rule ID'leri
- `allowed_operators` - İzin verilen operatör yapılandırmaları (PLMN, GID1, GID2)

**Kullanım Alanları:**

1. **Yüklemeden Önce Mevcut Depolamayı Kontrol Et:**
```bash
# Chip bilgisini al
INFO=$(hermes-euicc chip-info)

# Boş depolamayı çıkar
FREE_STORAGE=$(echo "$INFO" | jq -r '.data.euicc_info2.ext_card_resource.free_non_volatile_memory')

# Yeterli alan olup olmadığını kontrol et (500 KB gerekli)
if [ "$FREE_STORAGE" -lt 512000 ]; then
    echo "Uyarı: Düşük depolama! Sadece $FREE_STORAGE byte mevcut"
    exit 1
fi

# Profil yüklemesine devam et
hermes-euicc download --code "..."
```

2. **Chip Yeteneklerini Doğrula:**
```bash
# JavaCard desteklenip desteklenmediğini kontrol et
hermes-euicc chip-info | jq -r '.data.euicc_info2.uicc_capability[] | select(. == "javacard")'

# Firmware versiyonunu kontrol et
hermes-euicc chip-info | jq -r '.data.euicc_info2.euicc_firmware_ver'
```

3. **Depolama İstatistiklerini Göster:**
```bash
INFO=$(hermes-euicc chip-info)
FREE_NVM=$(echo "$INFO" | jq -r '.data.euicc_info2.ext_card_resource.free_non_volatile_memory')
FREE_RAM=$(echo "$INFO" | jq -r '.data.euicc_info2.ext_card_resource.free_volatile_memory')

echo "Depolama: $((FREE_NVM / 1024)) KB mevcut"
echo "RAM: $((FREE_RAM / 1024)) KB mevcut"
```

**Özellikler:**
- Tam parse edilmiş chip bilgileri (hex data yok)
- İnsan tarafından okunabilir alan adları
- Hata toleranslı (sadece EID gerekli)
- lpac'in `chip info` komutu ile uyumlu
- Kapsamlı versiyon ve yetenek bilgisi
- Kapasite planlaması için bellek/depolama detayları

**Mevcut `info` komutu ile karşılaştırma:**

| Özellik | `info` | `chip-info` |
|---------|--------|-------------|
| Çıktı Formatı | Ham hex data | Parse edilmiş JSON |
| EID | Hex string | Hex string |
| Info1 | Hex bytes | Dahil değil |
| Info2 | Hex bytes | Tam parse edilmiş yapı |
| Bellek Bilgisi | Hayır | Evet (byte) |
| Yetenekler | Hayır | Evet (diziler) |
| Versiyonlar | Hayır | Evet (tüm versiyonlar) |
| Configured Addresses | Hayır | Evet |
| RAT | Hayır | Evet |
| Okunabilir | Hayır | Evet |

---

## Tam Özellik Karşılaştırması

### Önce vs Sonra

| Kategori | Önce | Sonra | İyileştirme |
|----------|------|-------|-------------|
| **Bildirim İşleme** | Manuel döngü | Kütüphane API | ~50 satır azaltıldı |
| **Keşif** | Manuel concurrent | Kütüphane API | ~15 satır azaltıldı |
| **Chip Bilgisi** | Sadece ham hex | Parse edilmiş data | Yeni özellik |
| **Tek Adım İndirme** | Mevcut değil | discover-download | Yeni özellik |
| **Seçici İşleme** | Mevcut değil | notification-process | Yeni özellik |
| **Kod Bakımı** | Manuel concurrency | Kütüphane halleder | Çok daha iyi |
| **Hata Yönetimi** | Temel | Öğe başına detaylı | Çok daha iyi |

### Komut Sayısı

**Önce:** 17 komut
**Sonra:** 20 komut

**Yeni komutlar:**
1. `notification-process` - Belirli bildirimleri işle
2. `discover-download` - Tek adımda keşif ve indirme
3. `chip-info` - Detaylı chip bilgileri

---

## Göç Kılavuzu

### Mevcut Scriptler İçin

#### Bildirim İşleme

**Eski yaklaşım (hala çalışır):**
```bash
# Bildirimleri al
NOTIFS=$(hermes-euicc notifications)

# Her birini işle
echo "$NOTIFS" | jq -r '.data[].sequence_number' | while read seq; do
    hermes-euicc notification-handle "$seq"
done
```

**Yeni yaklaşım (önerilen):**
```bash
# Hepsini otomatik işle
hermes-euicc auto-notification

# Veya belirli olanları işle
hermes-euicc notification-process 1 2 3
```

#### Profil Keşfi

**Eski yaklaşım (manuel çok adımlı):**
```bash
DISCOVERY=$(hermes-euicc discovery)
CODE=$(echo "$DISCOVERY" | jq -r '.data[0].address')
hermes-euicc download --code "$CODE" --confirm
```

**Yeni yaklaşım (tek komut):**
```bash
hermes-euicc discover-download
```

#### Chip Bilgileri

**Eski yaklaşım (sadece hex data):**
```bash
hermes-euicc info  # Hex stringler döndürür
```

**Yeni yaklaşım (parse edilmiş data):**
```bash
hermes-euicc chip-info  # Parse edilmiş JSON döndürür

# Belirli alanları çıkarmak kolay
hermes-euicc chip-info | jq -r '.data.euicc_info2.ext_card_resource.free_non_volatile_memory'
```

---

## Performans Notları

### Bildirim İşleme
- Kütüphane dahili olarak concurrent işleme yapar
- Bildirim başına ~350-800ms (APDU + HTTPS + kaldırma)
- Toplu işleme optimize edilmiştir

### Profil Keşfi
- Keşif denemesi başına ~1-2 saniye
- Kütüphane kimlik doğrulama ve event alma işlemlerini yapar
- Concurrent SM-DS sorgulaması kaldırıldı (kütüphane halleder)

### Chip Bilgileri
- Birden fazla APDU çağrısı toplanır
- Tipik yürütme: donanıma bağlı olarak 500ms-2s
- Sonuçlar cache'lenebilir (chip bilgisi nadiren değişir)

---

## Sorun Giderme

### notification-process

**Hata: "invalid sequence number"**
- Sequence number'ların integer olduğundan emin olun
- Geçerli sequence number'ları listelemek için `notifications` komutunu kullanın

**Hata: "notification retrieve failed"**
- Bildirim artık mevcut olmayabilir
- Önce `notifications` komutu ile kontrol edin

### discover-download

**"no profiles available" döndürür**
- Cihaz EID'si SM-DS ile kayıtlı olmayabilir
- `--server` flag ile farklı SM-DS sunucuları deneyin
- IMEI kimlik doğrulaması gerekli olabilir (`--imei` kullanın)

**Keşiften sonra indirme başarısız olur**
- Ağ bağlantı sorunları
- SM-DP+ sunucusuna erişilemiyor
- Response'daki hata mesajını kontrol edin

### chip-info

**Sınırlı bilgi döndürür**
- Bazı alanlar opsiyoneldir ve nil olabilir
- Sadece EID garanti edilir
- Eski chip'ler tüm özellikleri desteklemeyebilir

**Bellek değerleri yanlış görünüyor**
- Değerler byte cinsinden, KB/MB değil
- Dönüştürme: `bytes / 1024 = KB`, `KB / 1024 = MB`

---

## Ayrıca Bakınız

- [NOTIFICATIONS.md](../NOTIFICATIONS.md) - Detaylı bildirim işleme API dokümantasyonu
- [PROFILE_DISCOVERY.md](../PROFILE_DISCOVERY.md) - Detaylı profil keşfi API dokümantasyonu
- [CHIP_INFO.md](../CHIP_INFO.md) - Detaylı chip bilgisi API dokümantasyonu
- [USAGE_TR.md](USAGE_TR.md) - Tam komut referansı
- [APP_JSON.md](APP_JSON.md) - JSON response format dokümantasyonu

---

## Kütüphane Versiyonu

Bu özellikler euicc-go kütüphanesi v1.2.1 veya üstünü gerektirir.

Kütüphane versiyonunuzu kontrol etmek için:
```bash
go list -m github.com/KilimcininKorOglu/euicc-go
```

Güncellemek için:
```bash
go get -u github.com/KilimcininKorOglu/euicc-go@latest
go mod tidy
```

---

## Değişiklik Günlüğü

### v1.0.0 (2025-11-02)
- ✅ `notification-process` komutu eklendi
- ✅ `auto-notification` komutu kütüphane API ile geliştirildi
- ✅ `discovery` komutu kütüphane API ile geliştirildi
- ✅ `discover-download` komutu eklendi
- ✅ `chip-info` komutu eklendi
- ✅ euicc-go v1.2.1'e yükseltildi
- ✅ Kod tabanı ~200 satır azaltıldı
- ✅ Tüm manuel concurrency yönetimi kaldırıldı
- ✅ Tüm komutlarda hata yönetimi iyileştirildi

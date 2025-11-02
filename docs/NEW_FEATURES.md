# New Features in Hermes eUICC Manager

This document describes the new features added to the Hermes eUICC Manager CLI application, integrating the latest APIs from euicc-go library v1.2.1.

**Last Updated:** 2025-11-02

---

## Overview

Three major feature sets have been added to the CLI application, leveraging new APIs from the euicc-go library:

1. **Enhanced Notification Processing** - Automated bulk notification handling
2. **Profile Discovery & Download** - SM-DS profile discovery with one-step download
3. **Chip Information** - Detailed eUICC chip information with parsed data

All features maintain the JSON-only output format for easy automation and scripting.

---

## 1. Notification Processing Features

### notification-process

Process specific notifications by sequence number with automatic removal and error handling.

**Usage:**
```bash
hermes-euicc notification-process <sequence_number> [<sequence_number> ...]
```

**Examples:**
```bash
# Process single notification
hermes-euicc notification-process 1

# Process multiple notifications
hermes-euicc notification-process 1 2 3 5
```

**Output:**
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

**Features:**
- Process one or multiple notifications in a single command
- Automatic removal after successful processing
- Continues processing even if one fails
- Detailed per-notification results
- Shows removal status for each notification

**Use Cases:**
- Selective notification processing
- Retry failed notifications
- Process notifications in specific order
- Fine-grained control over notification handling

### auto-notification (Enhanced)

The existing `auto-notification` command has been upgraded to use the library's `ProcessAllNotifications()` function for better performance and reliability.

**Usage:**
```bash
hermes-euicc auto-notification
```

**Improvements:**
- Uses library's optimized processing
- Automatic removal of processed notifications
- Better error handling per notification
- Cleaner implementation (~50 lines of code reduced)

---

## 2. Profile Discovery Features

### discovery (Enhanced)

The existing `discovery` command has been upgraded to use the library's `DiscoverProfiles()` function.

**Usage:**
```bash
# Discover from default GSMA SM-DS
hermes-euicc discovery

# Discover from custom SM-DS server
hermes-euicc discovery --server prod.smds.rsp.goog

# Discover with IMEI authentication
hermes-euicc discovery --imei 123456789012345
```

**Improvements:**
- Uses library's robust discovery implementation
- Simplified code (~15 lines reduced)
- Better error handling
- No manual concurrent processing needed

### discover-download (New)

One-step discovery and download of the first available profile from SM-DS.

**Usage:**
```bash
# Discover and download from default SM-DS
hermes-euicc discover-download

# Discover and download from custom SM-DS
hermes-euicc discover-download --server prod.smds.rsp.goog

# With IMEI authentication
hermes-euicc discover-download --imei 123456789012345
```

**Output (profile found and downloaded):**
```json
{
  "success": true,
  "data": {
    "message": "profile downloaded successfully"
  }
}
```

**Output (no profiles available):**
```json
{
  "success": true,
  "data": {
    "message": "no profiles available for download"
  }
}
```

**Features:**
- Automatic discovery and download in one command
- Uses default GSMA SM-DS server if not specified
- Optional IMEI authentication
- Custom SM-DS server support
- Fault-tolerant (returns success even if no profiles found)

**Use Cases:**
- Automated provisioning scripts
- Initial device setup
- Quick profile installation without activation codes
- Zero-touch provisioning workflows

**Comparison with manual approach:**

**Old way:**
```bash
# Step 1: Discover
PROFILES=$(hermes-euicc discovery)

# Step 2: Extract first profile address
CODE=$(echo "$PROFILES" | jq -r '.data[0].address')

# Step 3: Download
hermes-euicc download --code "$CODE" --confirm
```

**New way:**
```bash
# One command does everything
hermes-euicc discover-download
```

---

## 3. Chip Information Feature

### chip-info (New)

Get detailed, parsed chip information including memory, capabilities, versions, and configuration.

**Usage:**
```bash
hermes-euicc chip-info
```

**Output:**
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

**Key Information Fields:**

#### EID
- Unique chip identifier (32 hex characters)

#### Configured Addresses
- `default_smdp_address` - Default SM-DP+ server
- `root_smds_address` - Root SM-DS server

#### EUICCInfo2

**Version Information:**
- `profile_version` - Profile specification version
- `svn` - SGP.22 specification version
- `euicc_firmware_ver` - Chip firmware version
- `ts102241_version` - JavaCard/ETSI version
- `global_platform_version` - GlobalPlatform version
- `pp_version` - Protection Profile version

**Memory/Storage (IMPORTANT):**
- `installed_application` - Number of installed apps
- `free_non_volatile_memory` - Available persistent storage (bytes)
- `free_volatile_memory` - Available RAM (bytes)

**Capabilities:**
- `uicc_capability` - Card capabilities (USIM, ISIM, JavaCard, etc.)
- `rsp_capability` - Remote SIM Provisioning capabilities

**Security:**
- `euicc_ci_pkid_list_for_verification` - Public key IDs for verification
- `euicc_ci_pkid_list_for_signing` - Public key IDs for signing
- `forbidden_profile_policy_rules` - Forbidden policy rules

**Classification:**
- `euicc_category` - Category (basicEuicc, mediumEuicc, contactlessEuicc, other)

**Certification:**
- `sas_accreditation_number` - SAS accreditation number
- `certification_data_object` - Platform and discovery info

#### Rules Authorisation Table
- `ppr_ids` - Profile Policy Rule IDs
- `allowed_operators` - Allowed operator configurations (PLMN, GID1, GID2)

**Use Cases:**

1. **Check Available Storage Before Installation:**
```bash
# Get chip info
INFO=$(hermes-euicc chip-info)

# Extract free storage
FREE_STORAGE=$(echo "$INFO" | jq -r '.data.euicc_info2.ext_card_resource.free_non_volatile_memory')

# Check if enough space (500 KB required)
if [ "$FREE_STORAGE" -lt 512000 ]; then
    echo "Warning: Low storage! Only $FREE_STORAGE bytes available"
    exit 1
fi

# Proceed with profile installation
hermes-euicc download --code "..."
```

2. **Verify Chip Capabilities:**
```bash
# Check if JavaCard is supported
hermes-euicc chip-info | jq -r '.data.euicc_info2.uicc_capability[] | select(. == "javacard")'

# Check firmware version
hermes-euicc chip-info | jq -r '.data.euicc_info2.euicc_firmware_ver'
```

3. **Display Storage Statistics:**
```bash
INFO=$(hermes-euicc chip-info)
FREE_NVM=$(echo "$INFO" | jq -r '.data.euicc_info2.ext_card_resource.free_non_volatile_memory')
FREE_RAM=$(echo "$INFO" | jq -r '.data.euicc_info2.ext_card_resource.free_volatile_memory')

echo "Storage: $((FREE_NVM / 1024)) KB available"
echo "RAM: $((FREE_RAM / 1024)) KB available"
```

**Features:**
- Full parsed chip information (no hex data)
- Human-readable field names
- Fault-tolerant (only EID is required)
- Compatible with lpac's `chip info` command
- Comprehensive version and capability information
- Memory/storage details for capacity planning

**Comparison with existing `info` command:**

| Feature | `info` | `chip-info` |
|---------|--------|-------------|
| Output Format | Raw hex data | Parsed JSON |
| EID | Hex string | Hex string |
| Info1 | Hex bytes | Not included |
| Info2 | Hex bytes | Fully parsed structure |
| Memory Info | No | Yes (bytes) |
| Capabilities | No | Yes (arrays) |
| Versions | No | Yes (all versions) |
| Configured Addresses | No | Yes |
| RAT | No | Yes |
| Human Readable | No | Yes |

---

## Complete Feature Comparison

### Before vs After

| Category | Before | After | Improvement |
|----------|--------|-------|-------------|
| **Notification Processing** | Manual loop | Library API | ~50 lines reduced |
| **Discovery** | Manual concurrent | Library API | ~15 lines reduced |
| **Chip Info** | Raw hex only | Parsed data | New feature |
| **One-step Download** | Not available | discover-download | New feature |
| **Selective Processing** | Not available | notification-process | New feature |
| **Code Maintainability** | Manual concurrency | Library handles it | Much better |
| **Error Handling** | Basic | Per-item detailed | Much better |

### Command Count

**Before:** 17 commands
**After:** 20 commands

**New commands:**
1. `notification-process` - Process specific notifications
2. `discover-download` - One-step discovery and download
3. `chip-info` - Detailed chip information

---

## Migration Guide

### For Existing Scripts

#### Notification Processing

**Old approach (still works):**
```bash
# Get notifications
NOTIFS=$(hermes-euicc notifications)

# Process each
echo "$NOTIFS" | jq -r '.data[].sequence_number' | while read seq; do
    hermes-euicc notification-handle "$seq"
done
```

**New approach (recommended):**
```bash
# Process all automatically
hermes-euicc auto-notification

# Or process specific ones
hermes-euicc notification-process 1 2 3
```

#### Profile Discovery

**Old approach (manual multi-step):**
```bash
DISCOVERY=$(hermes-euicc discovery)
CODE=$(echo "$DISCOVERY" | jq -r '.data[0].address')
hermes-euicc download --code "$CODE" --confirm
```

**New approach (one command):**
```bash
hermes-euicc discover-download
```

#### Chip Information

**Old approach (hex data only):**
```bash
hermes-euicc info  # Returns hex strings
```

**New approach (parsed data):**
```bash
hermes-euicc chip-info  # Returns parsed JSON

# Easy to extract specific fields
hermes-euicc chip-info | jq -r '.data.euicc_info2.ext_card_resource.free_non_volatile_memory'
```

---

## Performance Notes

### Notification Processing
- Library handles concurrent processing internally
- ~350-800ms per notification (APDU + HTTPS + removal)
- Batch processing is optimized

### Profile Discovery
- ~1-2 seconds per discovery attempt
- Library handles authentication and event retrieval
- Concurrent SM-DS querying removed (library handles it)

### Chip Information
- Multiple APDU calls aggregated
- Typical execution: 500ms-2s depending on hardware
- Results can be cached (chip info rarely changes)

---

## Troubleshooting

### notification-process

**Error: "invalid sequence number"**
- Ensure sequence numbers are integers
- Use `notifications` command to list valid sequence numbers

**Error: "notification retrieve failed"**
- Notification may no longer exist
- Check with `notifications` command first

### discover-download

**Returns "no profiles available"**
- Device EID may not be registered with SM-DS
- Try different SM-DS servers with `--server` flag
- IMEI authentication may be required (use `--imei`)

**Download fails after discovery**
- Network connectivity issues
- SM-DP+ server unreachable
- Check error message in response

### chip-info

**Returns limited information**
- Some fields are optional and may be nil
- Only EID is guaranteed
- Older chips may not support all features

**Memory values seem incorrect**
- Values are in bytes, not KB/MB
- Convert: `bytes / 1024 = KB`, `KB / 1024 = MB`

---

## See Also

- [NOTIFICATIONS.md](../NOTIFICATIONS.md) - Detailed notification processing API documentation
- [PROFILE_DISCOVERY.md](../PROFILE_DISCOVERY.md) - Detailed profile discovery API documentation
- [CHIP_INFO.md](../CHIP_INFO.md) - Detailed chip information API documentation
- [USAGE.md](USAGE.md) - Complete command reference
- [APP_JSON.md](APP_JSON.md) - JSON response format documentation

---

## Library Version

These features require euicc-go library v1.2.1 or later.

To check your library version:
```bash
go list -m github.com/KilimcininKorOglu/euicc-go
```

To update:
```bash
go get -u github.com/KilimcininKorOglu/euicc-go@latest
go mod tidy
```

---

## Changelog

### v1.0.0 (2025-11-02)
- ✅ Added `notification-process` command
- ✅ Enhanced `auto-notification` command with library API
- ✅ Enhanced `discovery` command with library API
- ✅ Added `discover-download` command
- ✅ Added `chip-info` command
- ✅ Upgraded to euicc-go v1.2.1
- ✅ Reduced codebase by ~200 lines
- ✅ Removed all manual concurrency management
- ✅ Improved error handling across all commands

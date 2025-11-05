# Hermes eUICC Manager

**Hermes** - Full-featured eUICC (eSIM) management CLI application with JSON output.

## Features

- ✅ **Full JSON Output** - All commands respond in JSON format
- ✅ **Automatic Driver Detection** - Auto-detects QMI, MBIM, AT, CCID drivers*
- ℹ️ **Platform Note** - CCID support only on amd64/arm64 platforms (not available on MIPS/i386/ARMv5-7)
- ✅ **All SGP.22 Functions** - Supports all library features
- ✅ **Error Handling** - Structured error messages
- ✅ **Command Line Flags** - Flexible configuration options

## Installation

### Quick Build

```bash
cd app
go build -o hermes-euicc
```

### Build Script

**Universal build for all 19 platforms:**
```bash
./build-all.sh
```

**Build output:**
- All binaries in `build/1.0.0/` (version-specific directory)
- Includes 9 OpenWRT IPK packages with proper architecture names
- SHA256SUMS generated automatically

**Supported platforms:**
- Linux: 3 architectures (amd64, i386, arm64)
- OpenWRT: 9 architectures (MIPS, ARM, x86)
- Windows: 3 architectures (amd64, i386, arm64)

**Note:** Script automatically installs Go 1.24.0 if not found. macOS and FreeBSD builds disabled due to upstream driver issues.

### Pre-built Binaries

Check the [Releases](https://github.com/KilimcininKorOglu/euicc-go/releases) page for pre-built binaries.

## Usage

### Global Options

```bash
-device string      # Device path (e.g.: /dev/cdc-wdm0, /dev/ttyUSB2)
-driver string      # Driver type: qmi, mbim, at, ccid (auto if not specified)
-slot int           # SIM slot number (default: 1)
-timeout int        # HTTP timeout in seconds (default: 30)
-verbose            # Verbose log output
```

### Commands

#### 0. Help and Version

```bash
# Help message
./hermes-euicc help

# Version information
./hermes-euicc version
```

**Version Output:**

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

#### 1. Get EID

```bash
./hermes-euicc eid
```

**Output:**

```json
{
  "success": true,
  "data": {
    "eid": "89049032003451234567890123456789"
  }
}
```

#### 2. eUICC Information

```bash
./hermes-euicc info
```

**Output:**

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

#### 3. Profile List

```bash
./hermes-euicc list
```

**Output:**

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

#### 4. Enable Profile

```bash
./hermes-euicc enable 8944476500001224158
```

**Output:**

```json
{
  "success": true,
  "data": {
    "message": "profile enabled successfully",
    "iccid": "8944476500001224158"
  }
}
```

#### 5. Disable Profile

```bash
./hermes-euicc disable 8944476500001224158
```

#### 6. Delete Profile

```bash
./hermes-euicc delete 8944476500001224158
```

#### 7. Set Nickname

```bash
./hermes-euicc nickname 8944476500001224158 "Work SIM"
```

#### 8. Download Profile

```bash
./hermes-euicc download \
  --code "LPA:1$smdp.io$QR-G-5C-1LS" \
  --imei "356938035643809" \
  --confirm
```

**Options:**

- `--code`: Activation code (required)
- `--imei`: IMEI number (optional)
- `--confirmation-code`: Confirmation code (if needed)
- `--confirm`: Auto-confirmation

**Output:**

```json
{
  "success": true,
  "data": {
    "isdp_aid": "A0000005591010FFFFFFFF8900000100",
    "notification": 1
  }
}
```

#### 9. Profile Discovery

```bash
./hermes-euicc discovery --imei "356938035643809"
```

**Options:**

- `--server`: SM-DS server URL (optional)
- `--imei`: IMEI number (optional)

**Output:**

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

#### 10. List Notifications

```bash
./hermes-euicc notifications
```

**Output:**

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

#### 11. Remove Notification

```bash
./hermes-euicc notification-remove 1
```

#### 12. Handle Notification

```bash
./hermes-euicc notification-handle 1
```

#### 13. Auto-notification

Automatically process all pending notifications concurrently.

```bash
./hermes-euicc auto-notification
```

**Output:**

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

#### 14. Configured Addresses

```bash
./hermes-euicc configured-addresses
```

**Output:**

```json
{
  "success": true,
  "data": {
    "default_smdp_address": "smdp.example.com",
    "root_smds_address": "smds.example.com"
  }
}
```

#### 15. Set Default SM-DP+ Address

```bash
./hermes-euicc set-default-dp "smdp.example.com"
```

#### 16. eUICC Challenge

```bash
./hermes-euicc challenge
```

**Output:**

```json
{
  "success": true,
  "data": {
    "challenge": "1234567890abcdef"
  }
}
```

#### 17. Memory Reset

```bash
./hermes-euicc memory-reset
```

## Error Handling

All errors are returned in JSON format:

```json
{
  "success": false,
  "error": "failed to initialize driver: no compatible driver found"
}
```

## Driver Selection

### Auto-Detection

The application tests drivers in the following order:

1. QMI (`/dev/cdc-wdm0`)
2. MBIM (`/dev/cdc-wdm0`)
3. AT (`/dev/ttyUSB2`, `/dev/ttyUSB3`, `/dev/ttyUSB1`)
4. CCID (USB smart card reader)

### Manual Selection

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

## Scripting Usage

### Bash Example

```bash
#!/bin/bash

# Get profiles as JSON
PROFILES=$(./hermes-euicc list)

# Parse with jq
ACTIVE_ICCID=$(echo "$PROFILES" | jq -r '.data[] | select(.profile_state == 1) | .iccid')

echo "Active profile ICCID: $ACTIVE_ICCID"
```

### Python Example

```python
import json
import subprocess

# Get profile list
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

## Examples

### Scenario 1: Automatic Profile Management

```bash
#!/bin/bash

# List all profiles
PROFILES=$(./hermes-euicc list)

# Find disabled profiles
DISABLED=$(echo "$PROFILES" | jq -r '.data[] | select(.profile_state == 0) | .iccid')

# Enable first disabled profile
FIRST_ICCID=$(echo "$DISABLED" | head -1)
./hermes-euicc enable "$FIRST_ICCID"
```

### Scenario 2: Profile Discovery and Download

```bash
#!/bin/bash

# Discover profiles
DISCOVERY=$(./hermes-euicc discovery --imei "356938035643809")

# Download first profile
CODE=$(echo "$DISCOVERY" | jq -r '.data[0].address')
./hermes-euicc download --code "$CODE" --confirm
```

### Scenario 3: Notification Processing

**Simple approach (recommended):**

```bash
#!/bin/bash

# Automatically process all pending notifications concurrently
./hermes-euicc auto-notification
```

**Manual approach:**

```bash
#!/bin/bash

# Get pending notifications
NOTIFICATIONS=$(./hermes-euicc notifications)

# Process each notification
echo "$NOTIFICATIONS" | jq -r '.data[].sequence_number' | while read seq; do
    ./hermes-euicc notification-handle "$seq"
done
```

## Profile State Values

- `0`: Disabled
- `1`: Enabled (Active)

## Profile Management Operation Values

- `0`: Install
- `1`: Enable
- `2`: Disable
- `3`: Delete

## Profile Class Values

- `test`: Test profile
- `provisioning`: Provisioning profile
- `operational`: Operational profile

## OpenWRT Installation

### 1. Identify Your Device Architecture

On your OpenWRT device:

```bash
# Check architecture
uname -m

# Check musl libc architecture (most reliable)
ls -la /lib/ld-musl-*.so.1
```

### 2. Select the Correct Binary or IPK

| musl libc file | Architecture | Binary | IPK Package |
|----------------|--------------|--------|-------------|
| `ld-musl-mips-sf.so.1` | MIPS 24Kc BE | `hermes-euicc-*-openwrt-mips` | `hermes-euicc_*_mips_24kc.ipk` |
| `ld-musl-mipsel-sf.so.1` | MIPS 24Kc LE | `hermes-euicc-*-openwrt-mipsle` | `hermes-euicc_*_mipsel_24kc.ipk` |
| `ld-musl-armhf.so.1` | ARM Cortex-A7 | `hermes-euicc-*-openwrt-arm_v7` | `hermes-euicc_*_arm_cortex-a7_neon-vfpv4.ipk` |
| `ld-musl-aarch64.so.1` | ARM64 Cortex-A53 | `hermes-euicc-*-openwrt-arm64` | `hermes-euicc_*_aarch64_cortex-a53.ipk` |
| `ld-musl-x86_64.so.1` | x86-64 | `hermes-euicc-*-openwrt-x86_64` | `hermes-euicc_*_x86_64.ipk` |

All files are in `build/1.0.0/` directory.

**Common Devices:**
- **TP-Link Archer**, **GL.iNet AR/XE series**, **Ubiquiti EdgeRouter** → MIPS 24Kc BE
- **GL.iNet MT series**, **MediaTek MT76xx routers** → MIPS 24Kc LE
- **GL.iNet B1300**, **IPQ40xx devices** → ARM Cortex-A7
- **BananaPi R3/R4**, **GL.iNet MT6000**, **MT7622/MT7986** → ARM64 Cortex-A53
- **PC Engines APU**, **Protectli**, **x86 VMs** → x86-64

### 3. Installation Methods

**Method A: Install IPK Package (Recommended)**
```bash
# On your computer: Transfer IPK to router
scp build/1.0.0/hermes-euicc_1.0.0-*_mips_24kc.ipk root@192.168.1.1:/tmp/

# On OpenWRT router: Install via opkg
opkg install /tmp/hermes-euicc_*.ipk

# Test
hermes-euicc version
```

**Method B: Manual Binary Installation**
```bash
# On your computer: Transfer binary to router
scp build/1.0.0/hermes-euicc-1.0.0-*-openwrt-mips root@192.168.1.1:/tmp/

# On OpenWRT router: Install
chmod +x /tmp/hermes-euicc-*-openwrt-mips
mv /tmp/hermes-euicc-*-openwrt-mips /usr/bin/hermes-euicc

# Test
hermes-euicc version
```

### 4. Verify Binary Works

```bash
# Check binary format (requires file command - optional)
file /usr/bin/hermes-euicc

# Test functionality
hermes-euicc eid
hermes-euicc list
```

**Important:** Always transfer binaries in **binary mode** (not text/ASCII mode) to avoid corruption.

## License

MIT License - See LICENSE file for details.

## Contact

Kilimcinin Kör Oğlu <k@keremgok.tr>

# Hermes eUICC Manager - Usage Guide

Complete usage documentation for Hermes eUICC Manager, a JSON-based CLI application for managing eSIM profiles on eUICC-enabled devices.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Global Options](#global-options)
- [Driver Selection](#driver-selection)
- [Commands Reference](#commands-reference)
- [JSON Output Format](#json-output-format)
- [Common Use Cases](#common-use-cases)
- [Troubleshooting](#troubleshooting)

## Overview

Hermes eUICC Manager is a command-line tool that provides complete control over eSIM profiles on eUICC-enabled modems and devices. All operations return structured JSON output for easy integration with scripts and automation systems.

**Key Features:**

- JSON-only output for automation
- Auto-detection of hardware drivers (QMI, MBIM, AT, CCID)
- Complete SGP.22 protocol implementation
- 17 commands covering all eSIM operations
- Cross-platform support (MIPS, ARM, x86)

## Installation

### Build from Source

```bash
cd app
go build -o hermes-euicc .
```

### Verify Installation

```bash
hermes-euicc version
```

## Global Options

These options apply to all commands:

### -device string

Device path for modem/reader. If not specified, auto-detection attempts common paths.

```bash
# QMI/MBIM modem
hermes-euicc -device /dev/cdc-wdm0 list

# AT modem
hermes-euicc -device /dev/ttyUSB2 list

# Auto-detect (default)
hermes-euicc list
```

### -driver string

Force specific driver instead of auto-detection.

Supported drivers:

- `qmi` - Qualcomm MSM Interface (most common)
- `mbim` - Mobile Broadband Interface Model
- `at` - AT commands (serial modems)
- `ccid` - USB smart card readers (Linux amd64/arm64 only)

```bash
hermes-euicc -driver qmi list
hermes-euicc -driver mbim list
hermes-euicc -driver at -device /dev/ttyUSB2 list
hermes-euicc -driver ccid list
```

### -slot int

SIM slot number for multi-SIM devices (default: 1).

```bash
# Use SIM slot 2
hermes-euicc -slot 2 list
```

### -timeout int

HTTP timeout in seconds for SM-DP+ operations (default: 30).

```bash
# Increase timeout to 60 seconds
hermes-euicc -timeout 60 download --code "LPA:..."
```

### -verbose

Enable detailed logging for debugging.

```bash
hermes-euicc -verbose list
```

## Driver Selection

### Auto-Detection (Recommended)

Simply run commands without specifying driver:

```bash
hermes-euicc list
```

Detection order:

1. QMI at `/dev/cdc-wdm0`
2. MBIM at `/dev/cdc-wdm0`
3. AT at `/dev/ttyUSB2`, `/dev/ttyUSB3`, `/dev/ttyUSB1`
4. CCID (USB smart card reader)

### Manual Selection

When auto-detection fails or you need specific driver:

```bash
# Force QMI
hermes-euicc -driver qmi -device /dev/cdc-wdm0 list

# Force AT with custom device
hermes-euicc -driver at -device /dev/ttyUSB0 list
```

### Driver Availability by Platform

| Platform | QMI | MBIM | AT | CCID |
|----------|-----|------|----|------|
| Linux amd64 | ✓ | ✓ | ✓ | ✓ |
| Linux arm64 | ✓ | ✓ | ✓ | ✓ |
| Linux MIPS | ✓ | ✓ | ✓ | ✗ |
| Linux ARMv7 | ✓ | ✓ | ✓ | ✗ |

## Commands Reference

### help - Show Help Message

Display usage information.

```bash
hermes-euicc help
```

### version - Version Information

Show application version and copyright.

```bash
hermes-euicc version
```

**Output:**

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

### eid - Get EID

Retrieve eUICC Identifier (32-character hexadecimal).

```bash
hermes-euicc eid
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

### info - Get eUICC Information

Retrieve comprehensive eUICC information including EID, EUICCInfo1, and EUICCInfo2.

```bash
hermes-euicc info
```

**Output:**

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

### list - List Profiles

List all eSIM profiles installed on eUICC.

```bash
hermes-euicc list
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
      "icon": "base64-encoded-png",
      "icon_file_type": "image/png"
    }
  ]
}
```

**Profile States:**

- `0` - Disabled
- `1` - Enabled (active)

**Profile Classes:**

- `test` - Test profile
- `provisioning` - Provisioning profile
- `operational` - Normal operational profile

### enable - Enable Profile

Activate a profile by ICCID. Only one profile can be enabled at a time.

```bash
hermes-euicc enable 8944476500001224158
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

**Note:** Enabling a profile will automatically disable the currently active profile.

### disable - Disable Profile

Deactivate a profile by ICCID.

```bash
hermes-euicc disable 8944476500001224158
```

**Output:**

```json
{
  "success": true,
  "data": {
    "message": "profile disabled successfully",
    "iccid": "8944476500001224158"
  }
}
```

### delete - Delete Profile

Permanently remove a profile from eUICC.

```bash
hermes-euicc delete 8944476500001224158
```

**Output:**

```json
{
  "success": true,
  "data": {
    "message": "profile deleted successfully",
    "iccid": "8944476500001224158"
  }
}
```

**Warning:** This operation is irreversible. The profile must be disabled before deletion.

### nickname - Set Profile Nickname

Set a custom nickname for a profile.

```bash
hermes-euicc nickname 8944476500001224158 "Work SIM"
```

**Output:**

```json
{
  "success": true,
  "data": {
    "message": "nickname set successfully",
    "iccid": "8944476500001224158",
    "nickname": "Work SIM"
  }
}
```

### download - Download Profile

Download and install a new eSIM profile using activation code.

**Options:**

- `--code` (required) - LPA activation code
- `--imei` (optional) - Device IMEI
- `--confirmation-code` (optional) - Profile confirmation code if required
- `--confirm` (optional) - Auto-confirm download without prompting

```bash
# Basic download with auto-confirm
hermes-euicc download \
  --code "LPA:1$smdp.io$MATCHING-ID" \
  --confirm

# Download with IMEI and confirmation code
hermes-euicc download \
  --code "LPA:1$smdp.io$MATCHING-ID" \
  --imei "356938035643809" \
  --confirmation-code "1234" \
  --confirm
```

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

**Activation Code Format:**

```
LPA:1$<smdp-address>$<matching-id>
```

Example: `LPA:1$smdp.io$QR-G-5C-1LS`

### discovery - Discover Profiles

Query SM-DS servers for available profile downloads.

**Options:**

- `--server` (optional) - Custom SM-DS server (default: tries multiple servers)
- `--imei` (optional) - Device IMEI

```bash
# Discover with IMEI
hermes-euicc discovery --imei "356938035643809"

# Discover from specific server
hermes-euicc discovery --server "lpa.ds.gsma.com" --imei "356938035643809"
```

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

### notifications - List Notifications

Retrieve pending notifications from eUICC.

```bash
hermes-euicc notifications
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

**Profile Management Operations:**

- `0` - Install
- `1` - Enable
- `2` - Disable
- `3` - Delete

### notification-remove - Remove Notification

Delete a notification by sequence number.

```bash
hermes-euicc notification-remove 1
```

**Output:**

```json
{
  "success": true,
  "data": {
    "message": "notification removed successfully",
    "sequence_number": 1
  }
}
```

### notification-handle - Handle Notification

Process and send notification to SM-DP+ server.

```bash
hermes-euicc notification-handle 1
```

**Output:**

```json
{
  "success": true,
  "data": {
    "message": "notification handled successfully",
    "sequence_number": 1
  }
}
```

### auto-notification - Automatically Process All Notifications

Automatically retrieve and process all pending notifications concurrently. This command lists all pending notifications, retrieves each one, and handles them in parallel for optimal performance.

```bash
hermes-euicc auto-notification
```

**Output (with notifications):**

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

**Output (no notifications):**

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

**Notes:**
- Processes notifications concurrently for better performance
- Returns detailed results for both successful and failed operations
- Useful for automated batch processing of pending notifications

### configured-addresses - Get Configured Addresses

Retrieve default SM-DP+ and root SM-DS addresses configured in eUICC.

```bash
hermes-euicc configured-addresses
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

### set-default-dp - Set Default SM-DP+ Address

Configure default SM-DP+ server address.

```bash
hermes-euicc set-default-dp "smdp.example.com"
```

**Output:**

```json
{
  "success": true,
  "data": {
    "message": "default DP address set successfully",
    "address": "smdp.example.com"
  }
}
```

### challenge - Get eUICC Challenge

Retrieve authentication challenge from eUICC.

```bash
hermes-euicc challenge
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

### memory-reset - Reset eUICC Memory

Reset eUICC operational memory (does not delete profiles).

```bash
hermes-euicc memory-reset
```

**Output:**

```json
{
  "success": true,
  "data": {
    "message": "memory reset successfully"
  }
}
```

**Warning:** This operation resets notification lists and temporary data.

## JSON Output Format

All commands return JSON in consistent format:

### Success Response

```json
{
  "success": true,
  "data": {
    // Command-specific data
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": "error message description"
}
```

### Parsing Examples

**Bash with jq:**

```bash
# Get EID
EID=$(hermes-euicc eid | jq -r '.data.eid')
echo "EID: $EID"

# Get active profile ICCID
ACTIVE_ICCID=$(hermes-euicc list | jq -r '.data[] | select(.profile_state == 1) | .iccid')
echo "Active: $ACTIVE_ICCID"
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

## Common Use Cases

### Switch Between Profiles

```bash
#!/bin/bash

# List profiles and select
PROFILES=$(hermes-euicc list)
echo "$PROFILES" | jq -r '.data[] | "\(.iccid) - \(.profile_name) (state: \(.profile_state))"'

# Disable current
CURRENT=$(echo "$PROFILES" | jq -r '.data[] | select(.profile_state == 1) | .iccid')
hermes-euicc disable "$CURRENT"

# Enable new
hermes-euicc enable "8944476500001224159"
```

### Auto-Discovery and Download

```bash
#!/bin/bash

# Discover available profiles
DISCOVERY=$(hermes-euicc discovery --imei "356938035643809")

# Get first activation code
CODE=$(echo "$DISCOVERY" | jq -r '.data[0].address')

# Download profile
hermes-euicc download --code "$CODE" --confirm
```

### Batch Profile Management

```bash
#!/bin/bash

# List all disabled profiles
DISABLED=$(hermes-euicc list | jq -r '.data[] | select(.profile_state == 0) | .iccid')

# Delete all disabled profiles
for ICCID in $DISABLED; do
    echo "Deleting $ICCID..."
    hermes-euicc delete "$ICCID"
done
```

### Notification Processing

```bash
#!/bin/bash

# Get pending notifications
NOTIFICATIONS=$(hermes-euicc notifications)

# Process each notification
echo "$NOTIFICATIONS" | jq -r '.data[].sequence_number' | while read SEQ; do
    echo "Handling notification $SEQ..."
    hermes-euicc notification-handle "$SEQ"
done
```

### Profile Download with Retry

```bash
#!/bin/bash

CODE="LPA:1$smdp.io$MATCHING-ID"
MAX_RETRIES=3

for i in $(seq 1 $MAX_RETRIES); do
    echo "Attempt $i of $MAX_RETRIES..."
    
    RESULT=$(hermes-euicc download --code "$CODE" --confirm 2>&1)
    
    if echo "$RESULT" | jq -e '.success' > /dev/null 2>&1; then
        echo "Download successful!"
        echo "$RESULT" | jq '.'
        exit 0
    fi
    
    echo "Failed: $(echo "$RESULT" | jq -r '.error')"
    sleep 5
done

echo "Download failed after $MAX_RETRIES attempts"
exit 1
```

## Troubleshooting

### Device Not Found

**Error:**

```json
{
  "success": false,
  "error": "failed to initialize driver: no compatible driver found"
}
```

**Solutions:**

1. Check device path:

```bash
ls -la /dev/cdc-wdm* /dev/ttyUSB*
```

2. Check permissions:

```bash
sudo chmod 666 /dev/cdc-wdm0
```

3. Try manual driver selection:

```bash
hermes-euicc -driver qmi -device /dev/cdc-wdm0 list
```

### CCID Driver Not Available

**Error:**

```json
{
  "success": false,
  "error": "CCID driver not supported on this platform (requires amd64/arm64 + linux)"
}
```

**Solution:** CCID is only available on Linux amd64/arm64. Use QMI/MBIM/AT drivers instead.

### Invalid ICCID

**Error:**

```json
{
  "success": false,
  "error": "invalid ICCID: ..."
}
```

**Solution:** Ensure ICCID is exactly 19-20 digits. Get correct ICCID from `list` command.

### Profile Already Enabled

**Error:**

```json
{
  "success": false,
  "error": "profile already enabled"
}
```

**Solution:** Only one profile can be enabled at a time. Disable current profile first, or enable will do it automatically.

### Download Requires Confirmation

**Error:**

```json
{
  "success": false,
  "error": "download requires user confirmation"
}
```

**Solution:** Add `--confirm` flag or provide confirmation code:

```bash
hermes-euicc download --code "LPA:..." --confirm
```

### Timeout During Download

**Error:**

```json
{
  "success": false,
  "error": "context deadline exceeded"
}
```

**Solution:** Increase timeout:

```bash
hermes-euicc -timeout 60 download --code "LPA:..." --confirm
```

### Permission Denied

**Error:**

```bash
permission denied: /dev/cdc-wdm0
```

**Solutions:**

1. Run with sudo:

```bash
sudo hermes-euicc list
```

2. Add user to dialout group:

```bash
sudo usermod -aG dialout $USER
# Logout and login again
```

3. Set device permissions:

```bash
sudo chmod 666 /dev/cdc-wdm0
```

## Integration Examples

### Systemd Service

```ini
[Unit]
Description=eSIM Profile Switcher
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/hermes-euicc enable 8944476500001224158
User=root

[Install]
WantedBy=multi-user.target
```

### Cron Job

```bash
# Switch profile daily at 6 AM
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

## Additional Resources

- **Build Guide:** [BUILD.md](BUILD.md)
- **Driver Documentation:** [DRIVERS.md](DRIVERS.md)
- **Project README:** [../README.md](../README.md)
- **GSMA SGP.22 Specification:** <https://aka.pw/sgp22/v2.5>

## Support

For issues, questions, or contributions:

- GitHub Issues: <https://github.com/KilimcininKorOglu/euicc-go/issues>
- Documentation: See `/docs` directory in repository root

## License

MIT License - Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>

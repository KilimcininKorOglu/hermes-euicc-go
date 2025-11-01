# JSON Output Reference - Hermes eUICC Manager

Complete documentation of JSON outputs for all commands in Hermes eUICC Manager.

## Table of Contents

- [Response Format](#response-format)
- [Command Reference](#command-reference)
  - [version](#version)
  - [eid](#eid)
  - [info](#info)
  - [list](#list)
  - [enable](#enable)
  - [disable](#disable)
  - [delete](#delete)
  - [nickname](#nickname)
  - [download](#download)
  - [discovery](#discovery)
  - [notifications](#notifications)
  - [notification-remove](#notification-remove)
  - [notification-handle](#notification-handle)
  - [configured-addresses](#configured-addresses)
  - [set-default-dp](#set-default-dp)
  - [challenge](#challenge)
  - [memory-reset](#memory-reset)
- [Error Responses](#error-responses)

## Response Format

All commands return JSON with the following structure:

### Success Response

```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "error": "error message"
}
```

## Command Reference

### version

**Command:** `hermes-euicc version`

**Description:** Display version information

**Success Response:**

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

**Possible Errors:** None (this command cannot fail)

---

### eid

**Command:** `hermes-euicc eid`

**Description:** Get eUICC Identifier (EID)

**Success Response:**

```json
{
  "success": true,
  "data": {
    "eid": "89049032003451234567890123456789"
  }
}
```

**Fields:**

- `eid` (string): 32-character hexadecimal EID

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to initialize driver: no compatible driver found"
}
```

```json
{
  "success": false,
  "error": "failed to get EID: card error"
}
```

**Possible Errors:**

- Driver initialization failure
- Card communication error
- No eUICC present

---

### info

**Command:** `hermes-euicc info`

**Description:** Get complete eUICC information (EID + EUICCInfo1 + EUICCInfo2)

**Success Response:**

```json
{
  "success": true,
  "data": {
    "eid": "89049032003451234567890123456789",
    "euicc_info1": "bf20...",
    "euicc_info2": "bf22..."
  }
}
```

**Fields:**

- `eid` (string): 32-character hexadecimal EID
- `euicc_info1` (string): Hexadecimal encoded EUICCInfo1 (basic information)
- `euicc_info2` (string): Hexadecimal encoded EUICCInfo2 (detailed information)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to get EID: card error"
}
```

```json
{
  "success": false,
  "error": "failed to get EUICCInfo1: unsupported operation"
}
```

```json
{
  "success": false,
  "error": "failed to get EUICCInfo2: unsupported operation"
}
```

**Possible Errors:**

- Driver initialization failure
- Card communication error
- Unsupported eUICC version

---

### list

**Command:** `hermes-euicc list`

**Description:** List all profiles on eUICC

**Success Response (with profiles):**

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
      "icon": "base64encodedimage",
      "icon_file_type": "image/png"
    },
    {
      "iccid": "8944476500001224159",
      "isdp_aid": "A0000005591010FFFFFFFF8900000101",
      "profile_state": 0,
      "profile_name": "Work",
      "profile_nickname": "Office",
      "service_provider_name": "T-Mobile",
      "profile_class": "operational"
    }
  ]
}
```

**Success Response (no profiles):**

```json
{
  "success": true,
  "data": []
}
```

**Fields:**

- `iccid` (string): Integrated Circuit Card Identifier
- `isdp_aid` (string): ISD-P Application Identifier (optional)
- `profile_state` (int): Profile state (0=disabled, 1=enabled)
- `profile_name` (string): Profile name (optional)
- `profile_nickname` (string): User-set nickname (optional)
- `service_provider_name` (string): Operator name (optional)
- `profile_class` (string): Profile class ("test", "provisioning", "operational")
- `icon` (string): Base64 encoded icon (optional)
- `icon_file_type` (string): MIME type of icon (optional)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to list profiles: card error"
}
```

**Possible Errors:**

- Driver initialization failure
- Card communication error

---

### enable

**Command:** `hermes-euicc enable <iccid>`

**Description:** Enable a profile by ICCID

**Success Response:**

```json
{
  "success": true,
  "data": {
    "message": "profile enabled successfully",
    "iccid": "8944476500001224158"
  }
}
```

**Fields:**

- `message` (string): Success message
- `iccid` (string): The ICCID that was enabled

**Error Response Examples:**

```json
{
  "success": false,
  "error": "usage: enable <iccid>"
}
```

```json
{
  "success": false,
  "error": "invalid ICCID: invalid format"
}
```

```json
{
  "success": false,
  "error": "profile not found"
}
```

```json
{
  "success": false,
  "error": "failed to enable profile: another profile is already enabled"
}
```

**Possible Errors:**

- Missing ICCID argument
- Invalid ICCID format
- Profile not found
- Another profile already enabled
- Card communication error

---

### disable

**Command:** `hermes-euicc disable <iccid>`

**Description:** Disable a profile by ICCID

**Success Response:**

```json
{
  "success": true,
  "data": {
    "message": "profile disabled successfully",
    "iccid": "8944476500001224158"
  }
}
```

**Fields:**

- `message` (string): Success message
- `iccid` (string): The ICCID that was disabled

**Error Response Examples:**

```json
{
  "success": false,
  "error": "usage: disable <iccid>"
}
```

```json
{
  "success": false,
  "error": "invalid ICCID: invalid format"
}
```

```json
{
  "success": false,
  "error": "profile not found"
}
```

```json
{
  "success": false,
  "error": "profile is already disabled"
}
```

**Possible Errors:**

- Missing ICCID argument
- Invalid ICCID format
- Profile not found
- Profile already disabled
- Card communication error

---

### delete

**Command:** `hermes-euicc delete <iccid>`

**Description:** Delete a profile by ICCID

**Success Response:**

```json
{
  "success": true,
  "data": {
    "message": "profile deleted successfully",
    "iccid": "8944476500001224158"
  }
}
```

**Fields:**

- `message` (string): Success message
- `iccid` (string): The ICCID that was deleted

**Error Response Examples:**

```json
{
  "success": false,
  "error": "usage: delete <iccid>"
}
```

```json
{
  "success": false,
  "error": "invalid ICCID: invalid format"
}
```

```json
{
  "success": false,
  "error": "profile not found"
}
```

```json
{
  "success": false,
  "error": "cannot delete enabled profile"
}
```

**Possible Errors:**

- Missing ICCID argument
- Invalid ICCID format
- Profile not found
- Trying to delete enabled profile (must disable first)
- Card communication error

---

### nickname

**Command:** `hermes-euicc nickname <iccid> <nickname>`

**Description:** Set profile nickname

**Success Response:**

```json
{
  "success": true,
  "data": {
    "message": "nickname set successfully",
    "iccid": "8944476500001224158",
    "nickname": "My Work SIM"
  }
}
```

**Fields:**

- `message` (string): Success message
- `iccid` (string): The ICCID
- `nickname` (string): The new nickname

**Error Response Examples:**

```json
{
  "success": false,
  "error": "usage: nickname <iccid> <nickname>"
}
```

```json
{
  "success": false,
  "error": "invalid ICCID: invalid format"
}
```

```json
{
  "success": false,
  "error": "profile not found"
}
```

**Possible Errors:**

- Missing arguments
- Invalid ICCID format
- Profile not found
- Card communication error

---

### download

**Command:** `hermes-euicc download --code <activation-code> [--imei <imei>] [--confirmation-code <code>] [--confirm]`

**Description:** Download a new eSIM profile

**Success Response:**

```json
{
  "success": true,
  "data": {
    "isdp_aid": "A0000005591010FFFFFFFF8900000100",
    "notification": 1
  }
}
```

**Fields:**

- `isdp_aid` (string): ISD-P Application Identifier of downloaded profile
- `notification` (int): Profile management operation code (0=install, 1=enable, 2=disable, 3=delete)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "activation code required: use --code"
}
```

```json
{
  "success": false,
  "error": "invalid activation code: invalid format"
}
```

```json
{
  "success": false,
  "error": "invalid IMEI: invalid format"
}
```

```json
{
  "success": false,
  "error": "download failed: insufficient memory"
}
```

```json
{
  "success": false,
  "error": "download failed: invalid confirmation code"
}
```

```json
{
  "success": false,
  "error": "download cancelled by user"
}
```

**Possible Errors:**

- Missing activation code
- Invalid activation code format
- Invalid IMEI format
- Insufficient eUICC memory
- Invalid confirmation code
- User cancelled download
- Network/server errors
- Card communication error

**Activation Code Format:**

- `LPA:1$smdp.io$MATCHING-ID`
- `LPA:1$smdp.io$MATCHING-ID$CONFIRMATION-CODE`

---

### discovery

**Command:** `hermes-euicc discovery [--server <server>] [--imei <imei>]`

**Description:** Discover available profiles from SM-DS server

**Success Response (with profiles):**

```json
{
  "success": true,
  "data": [
    {
      "event_id": "EVT-123456",
      "address": "LPA:1$smdp.io$MATCHING-ID-1"
    },
    {
      "event_id": "EVT-789012",
      "address": "LPA:1$smdp.example.com$MATCHING-ID-2"
    }
  ]
}
```

**Success Response (no profiles):**

```json
{
  "success": true,
  "data": []
}
```

**Fields:**

- `event_id` (string): Event identifier from SM-DS
- `address` (string): Activation code address

**Error Response Examples:**

```json
{
  "success": false,
  "error": "invalid IMEI: invalid format"
}
```

```json
{
  "success": false,
  "error": "discovery failed: network error"
}
```

**Possible Errors:**

- Invalid IMEI format
- Network/server errors
- No SM-DS servers available

**Default SM-DS Servers:**

- `lpa.ds.gsma.com`
- `lpa.live.esimdiscovery.com`

---

### notifications

**Command:** `hermes-euicc notifications`

**Description:** List pending notifications

**Success Response (with notifications):**

```json
{
  "success": true,
  "data": [
    {
      "sequence_number": 1,
      "profile_management_operation": 1,
      "address": "smdp.example.com",
      "iccid": "8944476500001224158"
    },
    {
      "sequence_number": 2,
      "profile_management_operation": 0,
      "address": "smdp.io",
      "iccid": "8944476500001224159"
    }
  ]
}
```

**Success Response (no notifications):**

```json
{
  "success": true,
  "data": []
}
```

**Fields:**

- `sequence_number` (int): Notification sequence number
- `profile_management_operation` (int): Operation type (0=install, 1=enable, 2=disable, 3=delete)
- `address` (string): SM-DP+ server address (optional)
- `iccid` (string): Profile ICCID (optional)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to list notifications: card error"
}
```

**Possible Errors:**

- Card communication error

---

### notification-remove

**Command:** `hermes-euicc notification-remove <sequence-number>`

**Description:** Remove notification from list

**Success Response:**

```json
{
  "success": true,
  "data": {
    "message": "notification removed successfully",
    "sequence_number": 1
  }
}
```

**Fields:**

- `message` (string): Success message
- `sequence_number` (int): The sequence number that was removed

**Error Response Examples:**

```json
{
  "success": false,
  "error": "usage: notification-remove <sequence-number>"
}
```

```json
{
  "success": false,
  "error": "invalid sequence number: strconv.Atoi: parsing \"abc\": invalid syntax"
}
```

```json
{
  "success": false,
  "error": "notification not found"
}
```

**Possible Errors:**

- Missing sequence number argument
- Invalid sequence number format
- Notification not found
- Card communication error

---

### notification-handle

**Command:** `hermes-euicc notification-handle <sequence-number>`

**Description:** Process and handle notification

**Success Response:**

```json
{
  "success": true,
  "data": {
    "message": "notification handled successfully",
    "sequence_number": 1
  }
}
```

**Fields:**

- `message` (string): Success message
- `sequence_number` (int): The sequence number that was handled

**Error Response Examples:**

```json
{
  "success": false,
  "error": "usage: notification-handle <sequence-number>"
}
```

```json
{
  "success": false,
  "error": "invalid sequence number: strconv.Atoi: parsing \"abc\": invalid syntax"
}
```

```json
{
  "success": false,
  "error": "notification not found"
}
```

```json
{
  "success": false,
  "error": "failed to handle notification: network error"
}
```

**Possible Errors:**

- Missing sequence number argument
- Invalid sequence number format
- Notification not found
- Network/server errors
- Card communication error

---

### configured-addresses

**Command:** `hermes-euicc configured-addresses`

**Description:** Get configured SM-DP+ and SM-DS addresses

**Success Response:**

```json
{
  "success": true,
  "data": {
    "default_smdp_address": "smdp.example.com",
    "root_smds_address": "smds.example.com"
  }
}
```

**Success Response (no addresses configured):**

```json
{
  "success": true,
  "data": {
    "default_smdp_address": "",
    "root_smds_address": ""
  }
}
```

**Fields:**

- `default_smdp_address` (string): Default SM-DP+ server address (optional)
- `root_smds_address` (string): Root SM-DS server address (optional)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to get configured addresses: card error"
}
```

**Possible Errors:**

- Card communication error
- Unsupported operation

---

### set-default-dp

**Command:** `hermes-euicc set-default-dp <address>`

**Description:** Set default SM-DP+ server address

**Success Response:**

```json
{
  "success": true,
  "data": {
    "message": "default DP address set successfully",
    "address": "smdp.example.com"
  }
}
```

**Fields:**

- `message` (string): Success message
- `address` (string): The address that was set

**Error Response Examples:**

```json
{
  "success": false,
  "error": "usage: set-default-dp <address>"
}
```

```json
{
  "success": false,
  "error": "failed to set default DP address: card error"
}
```

**Possible Errors:**

- Missing address argument
- Card communication error
- Unsupported operation

---

### challenge

**Command:** `hermes-euicc challenge`

**Description:** Get eUICC challenge for authentication

**Success Response:**

```json
{
  "success": true,
  "data": {
    "challenge": "1234567890abcdef1234567890abcdef"
  }
}
```

**Fields:**

- `challenge` (string): Hexadecimal encoded challenge (typically 16 bytes = 32 hex chars)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to get challenge: card error"
}
```

**Possible Errors:**

- Card communication error

---

### memory-reset

**Command:** `hermes-euicc memory-reset`

**Description:** Reset eUICC memory (WARNING: Deletes all profiles!)

**Success Response:**

```json
{
  "success": true,
  "data": {
    "message": "memory reset successfully"
  }
}
```

**Fields:**

- `message` (string): Success message

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to reset memory: card error"
}
```

```json
{
  "success": false,
  "error": "operation not permitted"
}
```

**Possible Errors:**

- Card communication error
- Operation not permitted (some eUICCs don't allow reset)

⚠️ **WARNING:** This operation is irreversible and will delete ALL profiles on the eUICC!

---

## Error Responses

### Common Error Types

#### Driver Initialization Errors

```json
{
  "success": false,
  "error": "failed to initialize driver: no compatible driver found"
}
```

```json
{
  "success": false,
  "error": "failed to initialize driver: device path required for AT driver"
}
```

```json
{
  "success": false,
  "error": "failed to initialize driver: unknown driver type: invalid"
}
```

```json
{
  "success": false,
  "error": "failed to initialize driver: CCID driver not supported on this platform (requires amd64/arm64 + linux)"
}
```

#### Invalid Command Errors

```json
{
  "success": false,
  "error": "unknown command: invalid-command"
}
```

#### Card Communication Errors

```json
{
  "success": false,
  "error": "failed to communicate with card: device not found"
}
```

```json
{
  "success": false,
  "error": "card error: 6985 (conditions of use not satisfied)"
}
```

```json
{
  "success": false,
  "error": "card error: 6a82 (file not found)"
}
```

### Error Codes by Command

| Command | Common Error Scenarios |
|---------|------------------------|
| version | None (cannot fail) |
| eid | Driver init, card communication |
| info | Driver init, card communication, unsupported eUICC |
| list | Driver init, card communication |
| enable | Missing args, invalid ICCID, profile not found, another profile enabled |
| disable | Missing args, invalid ICCID, profile not found, already disabled |
| delete | Missing args, invalid ICCID, profile not found, enabled profile |
| nickname | Missing args, invalid ICCID, profile not found |
| download | Missing args, invalid code/IMEI, insufficient memory, network error |
| discovery | Invalid IMEI, network error |
| notifications | Card communication |
| notification-remove | Missing args, invalid seq number, not found |
| notification-handle | Missing args, invalid seq number, not found, network error |
| configured-addresses | Card communication |
| set-default-dp | Missing args, card communication |
| challenge | Card communication |
| memory-reset | Card communication, not permitted |

## JSON Parsing Examples

### Bash with jq

```bash
# Get EID
EID=$(./hermes-euicc eid | jq -r '.data.eid')
echo "EID: $EID"

# Check if command succeeded
if ./hermes-euicc list | jq -e '.success' > /dev/null; then
    echo "Success!"
else
    ERROR=$(./hermes-euicc list | jq -r '.error')
    echo "Error: $ERROR"
fi

# Get active profile ICCID
ACTIVE=$(./hermes-euicc list | jq -r '.data[] | select(.profile_state == 1) | .iccid')
```

### Python

```python
import json
import subprocess

# Get EID
result = subprocess.run(['./hermes-euicc', 'eid'], capture_output=True, text=True)
data = json.loads(result.stdout)

if data['success']:
    print(f"EID: {data['data']['eid']}")
else:
    print(f"Error: {data['error']}")

# List profiles
result = subprocess.run(['./hermes-euicc', 'list'], capture_output=True, text=True)
data = json.loads(result.stdout)

if data['success']:
    for profile in data['data']:
        state = "enabled" if profile['profile_state'] == 1 else "disabled"
        print(f"{profile['iccid']}: {profile['profile_name']} ({state})")
```

### JavaScript/Node.js

```javascript
const { execSync } = require('child_process');

// Get EID
const eidOutput = execSync('./hermes-euicc eid').toString();
const eidData = JSON.parse(eidOutput);

if (eidData.success) {
    console.log(`EID: ${eidData.data.eid}`);
} else {
    console.error(`Error: ${eidData.error}`);
}

// List profiles
const listOutput = execSync('./hermes-euicc list').toString();
const listData = JSON.parse(listOutput);

if (listData.success) {
    listData.data.forEach(profile => {
        const state = profile.profile_state === 1 ? 'enabled' : 'disabled';
        console.log(`${profile.iccid}: ${profile.profile_name} (${state})`);
    });
}
```

## Notes

- All JSON output is pretty-printed with 2-space indentation
- String fields may be empty (`""`) but are never null
- Optional fields may be omitted entirely from response
- Hexadecimal strings (EID, challenge, etc.) are lowercase
- ICCID format: 19-20 digits
- Profile states: 0 (disabled), 1 (enabled)
- Profile management operations: 0 (install), 1 (enable), 2 (disable), 3 (delete)
- Profile classes: "test", "provisioning", "operational"

# JSON Output Reference - Hermes eUICC Manager

Complete documentation of JSON outputs for all commands in Hermes eUICC Manager.

## Table of Contents

- [Response Format](#response-format)
- [Command Reference](#command-reference)
  - [version](#version)
  - [eid](#eid)
  - [info](#info)
  - [chip-info](#chip-info)
  - [list](#list)
  - [enable](#enable)
  - [disable](#disable)
  - [delete](#delete)
  - [nickname](#nickname)
  - [download](#download)
  - [discovery](#discovery)
  - [discover-download](#discover-download)
  - [notifications](#notifications)
  - [notification-remove](#notification-remove)
  - [notification-handle](#notification-handle)
  - [auto-notification](#auto-notification)
  - [notification-process](#notification-process)
  - [configured-addresses](#configured-addresses)
  - [set-default-dp](#set-default-dp)
  - [challenge](#challenge)
  - [memory-reset](#memory-reset)
- [Error Responses](#error-responses)
- [JSON Parsing Examples](#json-parsing-examples)

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

**Description:** Get complete eUICC information (EID + EUICCInfo1 + EUICCInfo2 as hex-encoded raw data)

**Success Response:**

```json
{
  "success": true,
  "data": {
    "eid": "89049032003451234567890123456789",
    "euicc_info1": "bf2040a02d0c0b6f6b6f64656e6576696c6c650c1c6f6b6f64656e6576696c6c6540746573742e636f6d",
    "euicc_info2": "bf2284a02d0c0b6f6b6f64656e6576696c6c650c1c6f6b6f64656e6576696c6c6540746573742e636f6d82020240"
  }
}
```

**Fields:**

- `eid` (string): 32-character hexadecimal EID
- `euicc_info1` (string): Hexadecimal encoded EUICCInfo1 (basic information, raw ASN.1 bytes)
- `euicc_info2` (string): Hexadecimal encoded EUICCInfo2 (detailed information, raw ASN.1 bytes)

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

**Note:** For parsed/human-readable chip information, use the `chip-info` command instead.

---

### chip-info

**Command:** `hermes-euicc chip-info`

**Description:** Get comprehensive parsed chip information including memory, capabilities, version info, and authorization rules

**Success Response:**

```json
{
  "success": true,
  "data": {
    "eid": "89049032003451234567890123456789",
    "configured_addresses": {
      "default_smdp_address": "smdp.example.com",
      "root_smds_address": "smds.example.com"
    },
    "euicc_info2": {
      "profile_version": "2.3",
      "svn": "3",
      "euicc_firmware_ver": "1.0.0",
      "ts102241_version": "13.0",
      "global_platform_version": "2.3.1",
      "pp_version": "1.0",
      "ext_card_resource": {
        "installed_application": 3,
        "free_non_volatile_memory": 524288,
        "free_volatile_memory": 8192
      },
      "uicc_capability": [
        "contactless",
        "usim",
        "isim"
      ],
      "rsp_capability": [
        "profileDownload",
        "profileManagement",
        "localProfileManagement"
      ],
      "euicc_ci_pkid_list_for_verification": [
        "B-1B7F-3B",
        "B-1B7F-1A"
      ],
      "euicc_ci_pkid_list_for_signing": [
        "B-1B7F-3B"
      ],
      "forbidden_profile_policy_rules": [
        "forbidContactlessInTestMode"
      ],
      "euicc_category": "basicEuicc",
      "sas_accreditation_number": "SAS-2023-001",
      "certification_data_object": {
        "platform_label": "Example eUICC v1",
        "discovery_base_url": "https://discovery.example.com"
      }
    },
    "rules_authorisation_table": [
      {
        "ppr_ids": [
          "forbidDisable",
          "forbidDelete"
        ],
        "allowed_operators": [
          {
            "plmn": "310260",
            "gid1": "A1",
            "gid2": ""
          },
          {
            "plmn": "310410",
            "gid1": "",
            "gid2": ""
          }
        ]
      }
    ]
  }
}
```

**Fields:**

- `eid` (string): eUICC Identifier
- `configured_addresses` (object, optional): Configured SM-DP+ and SM-DS addresses
  - `default_smdp_address` (string): Default SM-DP+ server address
  - `root_smds_address` (string): Root SM-DS server address
- `euicc_info2` (object, optional): Detailed eUICC information (parsed)
  - **Version Information:**
    - `profile_version` (string): SGP.22 profile version
    - `svn` (string): Security Version Number
    - `euicc_firmware_ver` (string): Firmware version
    - `ts102241_version` (string): ETSI TS 102 241 version
    - `global_platform_version` (string): GlobalPlatform version
    - `pp_version` (string): Protection Profile version
  - **Memory/Storage:**
    - `ext_card_resource` (object): Memory and resource information
      - `installed_application` (uint32): Number of installed applications
      - `free_non_volatile_memory` (uint32): Available non-volatile memory in bytes
      - `free_volatile_memory` (uint32): Available volatile memory in bytes
  - **Capabilities:**
    - `uicc_capability` (array of strings): UICC capabilities (e.g., "contactless", "usim", "isim")
    - `rsp_capability` (array of strings): RSP capabilities (e.g., "profileDownload", "profileManagement")
  - **Security:**
    - `euicc_ci_pkid_list_for_verification` (array of strings): CI PKI IDs for certificate verification
    - `euicc_ci_pkid_list_for_signing` (array of strings): CI PKI IDs for signing operations
    - `forbidden_profile_policy_rules` (array of strings): Policy rules that cannot be used
  - **Classification:**
    - `euicc_category` (string): eUICC category (e.g., "basicEuicc", "mediumEuicc")
  - **Certification:**
    - `sas_accreditation_number` (string): SAS accreditation number
    - `certification_data_object` (object): Certification information
      - `platform_label` (string): Platform identification label
      - `discovery_base_url` (string): Base URL for discovery services
- `rules_authorisation_table` (array of objects, optional): Profile authorization rules
  - `ppr_ids` (array of strings): Profile Policy Rule identifiers
  - `allowed_operators` (array of objects): Operators authorized to use these rules
    - `plmn` (string): Public Land Mobile Network identifier (MCC+MNC)
    - `gid1` (string): Group Identifier 1 (optional)
    - `gid2` (string): Group Identifier 2 (optional)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to get chip info: card error"
}
```

**Possible Errors:**

- Driver initialization failure
- Card communication error
- Unsupported eUICC version

**Use Cases:**

- Get comprehensive chip information for debugging
- Check available memory before profile download
- Verify chip capabilities and certifications
- Review operator authorization rules
- Display detailed chip information in management UIs

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
      "icon": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
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

- `iccid` (string): Integrated Circuit Card Identifier (19-20 digits)
- `isdp_aid` (string): ISD-P Application Identifier (hexadecimal, optional)
- `profile_state` (int): Profile state (0=disabled, 1=enabled)
- `profile_name` (string): Profile name from SM-DP+ (optional)
- `profile_nickname` (string): User-set nickname (optional)
- `service_provider_name` (string): Operator/provider name (optional)
- `profile_class` (string): Profile class ("test", "provisioning", "operational")
- `icon` (string): Base64 encoded profile icon (optional)
- `icon_file_type` (string): MIME type of icon, e.g., "image/png", "image/jpeg" (optional)

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
- Invalid ICCID format (must be 19-20 digits)
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
- Invalid ICCID format (must be 19-20 digits)
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
- Invalid ICCID format (must be 19-20 digits)
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
- Invalid ICCID format (must be 19-20 digits)
- Profile not found
- Card communication error

---

### download

**Command:** `hermes-euicc download --code <activation-code> [--imei <imei>] [--confirmation-code <code>] [--confirm]`

**Description:** Download a new eSIM profile

**Flags:**

- `--code` (required): Activation code in LPA format
- `--imei` (optional): Device IMEI for authentication
- `--confirmation-code` (optional): Confirmation code if required by profile
- `--confirm` (optional): Auto-confirm download without user prompt

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
- Invalid IMEI format (must be 15 digits)
- Insufficient eUICC memory
- Invalid confirmation code
- User cancelled download (when --confirm not provided)
- Network/server errors
- Card communication error

**Activation Code Format:**

- Basic: `LPA:1$smdp.io$MATCHING-ID`
- With confirmation: `LPA:1$smdp.io$MATCHING-ID$CONFIRMATION-CODE`

**Examples:**

```bash
# Download with auto-confirm
hermes-euicc download --code "LPA:1$smdp.io$ABC123" --confirm

# Download with IMEI
hermes-euicc download --code "LPA:1$smdp.io$ABC123" --imei "356938035643809" --confirm

# Download with confirmation code
hermes-euicc download --code "LPA:1$smdp.io$ABC123" --confirmation-code "1234" --confirm
```

---

### discovery

**Command:** `hermes-euicc discovery [--server <server>] [--imei <imei>]`

**Description:** Discover available profiles from SM-DS server

**Flags:**

- `--server` (optional): SM-DS server address (default: uses chip's configured SM-DS or GSMA default)
- `--imei` (optional): Device IMEI for authentication

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
- `address` (string): Activation code address (LPA format)

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

```json
{
  "success": false,
  "error": "discovery failed: SM-DS authentication failed"
}
```

**Possible Errors:**

- Invalid IMEI format (must be 15 digits)
- Network/server errors
- No SM-DS servers available
- SM-DS authentication failure
- Card communication error

**Default SM-DS Servers:**

- GSMA default: `lpa.ds.gsma.com`
- Chip's configured SM-DS (if available)

**Examples:**

```bash
# Discover using default SM-DS
hermes-euicc discovery

# Discover with custom SM-DS
hermes-euicc discovery --server "smds.example.com"

# Discover with IMEI authentication
hermes-euicc discovery --imei "356938035643809"
```

---

### discover-download

**Command:** `hermes-euicc discover-download [--server <server>] [--imei <imei>]`

**Description:** Discover available profiles and automatically download the first one

**Flags:**

- `--server` (optional): SM-DS server address
- `--imei` (optional): Device IMEI for authentication

**Success Response (profile downloaded):**

```json
{
  "success": true,
  "data": {
    "message": "profile downloaded successfully"
  }
}
```

**Success Response (no profiles available):**

```json
{
  "success": true,
  "data": {
    "message": "no profiles available for download"
  }
}
```

**Fields:**

- `message` (string): Status message

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

```json
{
  "success": false,
  "error": "download failed: insufficient memory"
}
```

**Possible Errors:**

- Invalid IMEI format (must be 15 digits)
- Network/server errors during discovery
- No profiles available from SM-DS
- Download failures (same as `download` command)
- Card communication error

**Use Cases:**

- Automated profile provisioning
- Quick setup flows where user doesn't need to choose
- Scripted deployments

**Examples:**

```bash
# Discover and download first available profile
hermes-euicc discover-download

# With custom SM-DS and IMEI
hermes-euicc discover-download --server "smds.example.com" --imei "356938035643809"
```

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

- `sequence_number` (int): Notification sequence number (used to reference specific notifications)
- `profile_management_operation` (int): Operation type
  - 0 = install
  - 1 = enable
  - 2 = disable
  - 3 = delete
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

**Use Cases:**

- Check for pending notifications after profile operations
- Verify notification delivery before handling
- Audit notification queue

---

### notification-remove

**Command:** `hermes-euicc notification-remove <sequence-number>`

**Description:** Remove notification from list without handling it

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
- Invalid sequence number format (must be integer)
- Notification not found
- Card communication error

---

### notification-handle

**Command:** `hermes-euicc notification-handle <sequence-number>`

**Description:** Process and handle notification (sends notification to SM-DP+)

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
- Invalid sequence number format (must be integer)
- Notification not found
- Network/server errors
- Card communication error

**Note:** This command retrieves the full notification details and sends it to the SM-DP+ server for processing.

---

### auto-notification

**Command:** `hermes-euicc auto-notification`

**Description:** Automatically process all pending notifications with auto-removal

**Success Response (with notifications):**

```json
{
  "success": true,
  "data": {
    "message": "auto notification processing completed",
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
        "error": "notification not found"
      }
    ]
  }
}
```

**Success Response (no notifications):**

```json
{
  "success": true,
  "data": {
    "message": "auto notification processing completed",
    "total": 0,
    "processed": 0,
    "failed": 0,
    "processed_list": [],
    "failed_list": []
  }
}
```

**Fields:**

- `message` (string): Status message
- `total` (int): Total number of notifications found
- `processed` (int): Number of successfully processed notifications
- `failed` (int): Number of failed notifications
- `processed_list` (array): Details of successfully processed notifications
  - `sequence_number` (int): Notification sequence number
  - `removed` (bool): Whether notification was removed from list after handling
- `failed_list` (array): Details of failed notifications
  - `sequence_number` (int): Notification sequence number
  - `error` (string): Error message

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to list notifications: card error"
}
```

**Possible Errors:**

- Card communication error during notification listing

**Notes:**

- Processes all pending notifications automatically
- Uses `AutoRemove: true` to remove notifications after successful handling
- Uses `ContinueOnError: true` to process remaining notifications even if one fails
- Individual notification failures do not stop processing
- Returns detailed results for both successful and failed notifications
- Useful for bulk notification processing after profile operations

**Use Cases:**

- Automated notification cleanup after profile downloads
- Scheduled notification processing
- Bulk operations in scripts

---

### notification-process

**Command:** `hermes-euicc notification-process <sequence-number> [<sequence-number> ...]`

**Description:** Process specific notifications by sequence number(s) with auto-removal

**Success Response:**

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
        "sequence_number": 3,
        "removed": true
      }
    ],
    "failed_list": [
      {
        "sequence_number": 2,
        "error": "failed to handle notification: network error"
      }
    ]
  }
}
```

**Fields:**

- `message` (string): Status message
- `total` (int): Total number of notifications to process
- `processed` (int): Number of successfully processed notifications
- `failed` (int): Number of failed notifications
- `processed_list` (array): Details of successfully processed notifications
  - `sequence_number` (int): Notification sequence number
  - `removed` (bool): Whether notification was removed from list after handling
- `failed_list` (array): Details of failed notifications
  - `sequence_number` (int): Notification sequence number
  - `error` (string): Error message

**Error Response Examples:**

```json
{
  "success": false,
  "error": "sequence number(s) required"
}
```

```json
{
  "success": false,
  "error": "invalid sequence number '2a': strconv.Atoi: parsing \"2a\": invalid syntax"
}
```

**Possible Errors:**

- Missing sequence number arguments
- Invalid sequence number format (must be integers)
- Notification not found (reported in `failed_list`, not as error)
- Network/server errors (reported in `failed_list`, not as error)
- Card communication error

**Notes:**

- Accepts multiple sequence numbers as arguments
- Uses `AutoRemove: true` to remove notifications after successful handling
- Uses `ContinueOnError: true` to process remaining notifications even if one fails
- Individual notification failures do not stop processing
- Returns detailed results for both successful and failed notifications

**Use Cases:**

- Process specific notifications selectively
- Retry failed notifications from `auto-notification`
- Handle high-priority notifications immediately

**Examples:**

```bash
# Process single notification
hermes-euicc notification-process 1

# Process multiple notifications
hermes-euicc notification-process 1 3 5

# Process all notifications from a list
hermes-euicc notifications | jq -r '.data[].sequence_number' | xargs hermes-euicc notification-process
```

---

### configured-addresses

**Command:** `hermes-euicc configured-addresses`

**Description:** Get configured SM-DP+ and SM-DS addresses from the eUICC

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

- `default_smdp_address` (string): Default SM-DP+ server address (empty if not configured)
- `root_smds_address` (string): Root SM-DS server address (empty if not configured)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to get configured addresses: card error"
}
```

**Possible Errors:**

- Card communication error
- Unsupported operation (older eUICC versions)

---

### set-default-dp

**Command:** `hermes-euicc set-default-dp <address>`

**Description:** Set default SM-DP+ server address on the eUICC

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
- Invalid address format
- Card communication error
- Unsupported operation (older eUICC versions)

---

### challenge

**Command:** `hermes-euicc challenge`

**Description:** Get eUICC challenge for authentication purposes

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

- `challenge` (string): Hexadecimal encoded challenge (typically 16 bytes = 32 hex characters)

**Error Response Examples:**

```json
{
  "success": false,
  "error": "failed to get challenge: card error"
}
```

**Possible Errors:**

- Card communication error

**Use Cases:**

- Cryptographic authentication flows
- Custom security implementations
- Testing eUICC challenge generation

---

### memory-reset

**Command:** `hermes-euicc memory-reset`

**Description:** Reset eUICC memory (⚠️ WARNING: Deletes all profiles!)

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

⚠️ **WARNING:** This operation is irreversible and will delete ALL profiles on the eUICC! Use with extreme caution.

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

```json
{
  "success": false,
  "error": "card error: 6a80 (incorrect parameters in data field)"
}
```

### Error Codes by Command

| Command | Common Error Scenarios |
|---------|------------------------|
| version | None (cannot fail) |
| eid | Driver init, card communication |
| info | Driver init, card communication, unsupported eUICC |
| chip-info | Driver init, card communication, unsupported eUICC |
| list | Driver init, card communication |
| enable | Missing args, invalid ICCID, profile not found, another profile enabled |
| disable | Missing args, invalid ICCID, profile not found, already disabled |
| delete | Missing args, invalid ICCID, profile not found, enabled profile |
| nickname | Missing args, invalid ICCID, profile not found |
| download | Missing args, invalid code/IMEI, insufficient memory, network error |
| discovery | Invalid IMEI, network error |
| discover-download | Invalid IMEI, network error, download errors |
| notifications | Card communication |
| notification-remove | Missing args, invalid seq number, not found |
| notification-handle | Missing args, invalid seq number, not found, network error |
| auto-notification | Card communication (listing phase) |
| notification-process | Missing args, invalid seq numbers |
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
echo "Active profile: $ACTIVE"

# Get free memory from chip-info
FREE_NV=$(./hermes-euicc chip-info | jq -r '.data.euicc_info2.ext_card_resource.free_non_volatile_memory')
echo "Free memory: $FREE_NV bytes"

# List all pending notification sequence numbers
./hermes-euicc notifications | jq -r '.data[].sequence_number'

# Process all notifications automatically
./hermes-euicc auto-notification | jq '.data | "Processed: \(.processed)/\(.total), Failed: \(.failed)"'
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

# Get chip info
result = subprocess.run(['./hermes-euicc', 'chip-info'], capture_output=True, text=True)
data = json.loads(result.stdout)

if data['success']:
    info2 = data['data']['euicc_info2']
    print(f"Firmware: {info2['euicc_firmware_ver']}")
    print(f"Free NV Memory: {info2['ext_card_resource']['free_non_volatile_memory']} bytes")
    print(f"Capabilities: {', '.join(info2['rsp_capability'])}")

# Download profile with error handling
result = subprocess.run([
    './hermes-euicc', 'download',
    '--code', 'LPA:1$smdp.io$MATCHING-ID',
    '--confirm'
], capture_output=True, text=True)
data = json.loads(result.stdout)

if data['success']:
    print(f"Downloaded: {data['data']['isdp_aid']}")
else:
    print(f"Download failed: {data['error']}")

# Auto-process all notifications
result = subprocess.run(['./hermes-euicc', 'auto-notification'], capture_output=True, text=True)
data = json.loads(result.stdout)

if data['success']:
    print(f"Processed: {data['data']['processed']}/{data['data']['total']}")
    if data['data']['failed'] > 0:
        print("Failed notifications:")
        for failed in data['data']['failed_list']:
            print(f"  - Seq {failed['sequence_number']}: {failed['error']}")
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

// Get chip capabilities
const chipOutput = execSync('./hermes-euicc chip-info').toString();
const chipData = JSON.parse(chipOutput);

if (chipData.success) {
    const info2 = chipData.data.euicc_info2;
    console.log('Capabilities:', info2.rsp_capability.join(', '));
    console.log('Free Memory:', info2.ext_card_resource.free_non_volatile_memory, 'bytes');
}

// Discover and download
try {
    const discoverOutput = execSync('./hermes-euicc discover-download').toString();
    const discoverData = JSON.parse(discoverOutput);

    if (discoverData.success) {
        console.log(discoverData.data.message);
    }
} catch (err) {
    console.error('Discovery failed');
}

// Process specific notifications
const notifOutput = execSync('./hermes-euicc notification-process 1 2 3').toString();
const notifData = JSON.parse(notifOutput);

if (notifData.success) {
    console.log(`Processed: ${notifData.data.processed}/${notifData.data.total}`);
    notifData.data.processed_list.forEach(p => {
        console.log(`  ✓ Notification ${p.sequence_number} ${p.removed ? '(removed)' : ''}`);
    });
    notifData.data.failed_list.forEach(f => {
        console.log(`  ✗ Notification ${f.sequence_number}: ${f.error}`);
    });
}
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "os/exec"
)

type Response struct {
    Success bool            `json:"success"`
    Data    json.RawMessage `json:"data,omitempty"`
    Error   string          `json:"error,omitempty"`
}

type EIDResponse struct {
    EID string `json:"eid"`
}

type ProfileResponse struct {
    ICCID        string `json:"iccid"`
    ProfileState int    `json:"profile_state"`
    ProfileName  string `json:"profile_name"`
}

func main() {
    // Get EID
    out, _ := exec.Command("./hermes-euicc", "eid").Output()

    var resp Response
    json.Unmarshal(out, &resp)

    if resp.Success {
        var eidData EIDResponse
        json.Unmarshal(resp.Data, &eidData)
        fmt.Printf("EID: %s\n", eidData.EID)
    }

    // List profiles
    out, _ = exec.Command("./hermes-euicc", "list").Output()
    json.Unmarshal(out, &resp)

    if resp.Success {
        var profiles []ProfileResponse
        json.Unmarshal(resp.Data, &profiles)

        for _, p := range profiles {
            state := "disabled"
            if p.ProfileState == 1 {
                state = "enabled"
            }
            fmt.Printf("%s: %s (%s)\n", p.ICCID, p.ProfileName, state)
        }
    }
}
```

## Notes

- All JSON output is pretty-printed with 2-space indentation
- String fields may be empty (`""`) but are never null
- Optional fields may be omitted entirely from response
- Hexadecimal strings (EID, challenge, AIDs, etc.) are lowercase
- ICCID format: 19-20 digits
- IMEI format: 15 digits
- Profile states: 0 (disabled), 1 (enabled)
- Profile management operations: 0 (install), 1 (enable), 2 (disable), 3 (delete)
- Profile classes: "test", "provisioning", "operational"
- Memory sizes are in bytes (uint32)
- Sequence numbers are integers starting from 1
- All commands respect the global `--verbose` flag for detailed logging to stderr
- JSON output always goes to stdout, logs and errors to stderr

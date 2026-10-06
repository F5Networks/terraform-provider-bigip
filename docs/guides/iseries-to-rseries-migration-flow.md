---
page_title: "i-Series to r-Series migration flow across terraform-provider-bigip and terraform-provider-f5os"
description: |-
  End-to-end phase flow for migrating a BIG-IP i-Series device to an r-Series/F5OS target using the discovery phases in terraform-provider-bigip and the target-configuration phases in terraform-provider-f5os.
---

# i-Series to r-Series migration flow across terraform-provider-bigip and terraform-provider-f5os

This document describes the end-to-end i-Series to r-Series migration flow
using the phases already defined in two Terraform provider repos:

- `terraform-provider-bigip`
- `terraform-provider-f5os`

The two repos play different roles in the workflow:

- `terraform-provider-bigip` handles source-side discovery and backup from the
  existing BIG-IP i-Series device.
- `terraform-provider-f5os` handles target-side configuration of the r-Series
  F5OS platform and tenant objects.

The phases below are based on the existing guides, scripts, and example
directories in those repos rather than a new migration model.

## Phase map

| Phase | Repo | Purpose |
|---|---|---|
| Phase 0 | `terraform-provider-bigip` | Inventory the source TMOS version, build, and hardware model; determine the recommended r-Series tenant image |
| Phase 1 | `terraform-provider-bigip` | Extract i-Series system and network settings into JSON artifacts |
| Phase 2 | `terraform-provider-bigip` | Generate and download a UCS backup from the source i-Series device |
| Licensing phase | `terraform-provider-f5os` | Apply a new F5OS/r-Series license on the target device |
| System-settings phase | `terraform-provider-f5os` | Configure target DNS, NTP, SNMP, auth order, and local platform users |
| Phase 3 | `terraform-provider-f5os` | Create VLANs on the target F5OS platform |
| Phase 4 | `terraform-provider-f5os` | Configure target F5OS interfaces with native/trunk VLAN assignments |
| Phase 5 | `terraform-provider-f5os` | Configure target F5OS LAGs from source trunks |
| Phase 6 | `terraform-provider-f5os` | Deploy BIG-IP tenants on the r-Series target |

## Repo responsibilities

### terraform-provider-bigip

Source-side migration responsibilities:

- identify source TMOS version/build and hardware
- determine whether the source TMOS train is suitable for r-Series tenant use
- extract reusable system and network settings from the i-Series device
- produce a UCS archive for restore or reference during tenant bring-up

Primary guides/scripts:

- `docs/guides/inventory-tmos-version.md`
- `docs/guides/extract-sys-settings.md`
- `docs/guides/generate-ucs-backup.md`
- `docs/guides/interface-trunk-mapping.md`
- `scripts/inventory-tmos-version.sh`
- `scripts/extract-sys-settings.sh`
- `scripts/generate-ucs-backup.sh`

### terraform-provider-f5os

Target-side migration responsibilities:

- apply F5OS platform licensing
- configure system settings on the r-Series platform layer
- create platform VLANs
- map and configure physical interfaces
- map and configure LAGs
- deploy BIG-IP tenants with the required VLAN attachments and sizing

Primary guides/examples/scripts:

- `examples/migration/license-from-iseries/`
- `examples/migration/system-settings-from-iseries/`
- `examples/migration/vlans-from-iseries/`
- `docs/guides/apply-license-from-iseries.md`
- `docs/guides/configure-system-settings-from-iseries.md`
- `docs/guides/create-vlans-from-iseries.md`
- `docs/guides/configure-interfaces-from-iseries.md`
- `docs/guides/configure-lags-from-iseries.md`
- `docs/guides/deploy-tenants-from-iseries.md`

## End-to-end flow

### Phase 0: inventory the source device

Run in `terraform-provider-bigip`.

Goal:

- capture the exact TMOS version/build running on the source i-Series device
- capture hardware/platform identity
- derive the recommended r-Series tenant image filename pattern
- determine whether the source version is even suitable for r-Series tenant use

Inputs:

- source i-Series management connectivity and credentials

Outputs:

- source inventory JSON
- recommended tenant image/version information
- version support warning or confirmation for r-Series

Use:

- `scripts/inventory-tmos-version.sh`

Why first:

- image selection is a prerequisite for tenant deployment later in Phase 6
- unsupported TMOS trains should be discovered before the rest of the
  migration work proceeds

### Phase 1: extract source system and network settings

Run in `terraform-provider-bigip`.

Goal:

- discover reusable source configuration and serialize it into downstream
  JSON artifacts

Extracted categories include:

- DNS
- NTP
- SNMP
- syslog
- auth settings
- local users
- VLANs
- self IPs
- routes
- physical interfaces
- trunks/LAGs

Outputs:

- a general extracted settings JSON file
- an `interface_vlans` JSON file derived from VLAN membership

Use:

- `scripts/extract-sys-settings.sh`

Why this phase matters:

- these JSON artifacts are the source material for several F5OS conversion
  scripts and example directories

### Phase 2: generate and download a UCS backup

Run in `terraform-provider-bigip`.

Goal:

- create a complete UCS archive from the source i-Series device
- download it locally for upload/reference during tenant migration

Use:

- `scripts/generate-ucs-backup.sh`

Outputs:

- a UCS file on local disk

Why separate from Phase 1:

- Phase 1 extracts structured settings for translation into F5OS resources
- Phase 2 captures a fuller source backup for restore/reference needs that do
  not map directly to F5OS platform resources

## Target-side phases in terraform-provider-f5os

The next phases happen in `terraform-provider-f5os`.

These phases are split into independent tracks where the target platform
allows it.

### Licensing phase

Run in `examples/migration/license-from-iseries/`.

Goal:

- apply a new F5OS/r-Series license to the target platform

Important note:

- the i-Series TMOS license key is not migrated or reused
- this phase requires a new registration key obtained for the target device

Dependency:

- independent of VLAN/system-settings migration
- can be performed before, after, or in parallel with those tracks

### System-settings phase

Run in `examples/migration/system-settings-from-iseries/`.

Goal:

- configure target platform DNS, NTP, SNMP, auth order, and local users

Input source:

- Phase 1 extracted settings JSON
- converted via `scripts/system-settings-from-iseries.sh` in the F5OS repo

Dependency:

- parallel to the VLAN/interface/LAG workflow
- not a prerequisite for VLAN creation itself

### Phase 3: create VLANs on the target F5OS platform

Run in `examples/migration/vlans-from-iseries/main.tf`.

Goal:

- create one F5OS VLAN per VLAN discovered on the source i-Series device
- preserve source VLAN names and tags

Input source:

- Phase 1 extracted settings JSON
- converted via `scripts/vlans-from-iseries.sh`

Dependency:

- first phase of the target network-configuration workflow
- prerequisite for Phases 4, 5, and 6 where VLAN IDs are referenced

### Phase 4: configure F5OS interfaces

Run in `examples/migration/vlans-from-iseries/interfaces.tf`.

Goal:

- configure each target interface with enabled state and native/trunk VLAN
  assignments matching the source layout

Input source:

- Phase 1 extracted settings JSON
- converted via `scripts/interfaces-from-iseries.sh`

Dependency:

- depends on Phase 3 VLAN creation

Special mapping note:

- TMOS interface names such as `1.1` do not map 1:1 to r-Series names
- for r-Series, the target naming convention is `<port>.0` such as `1.0`
- use `terraform-provider-bigip/docs/guides/interface-trunk-mapping.md` and
  the F5OS interface guide to validate mappings and detect collisions

### Phase 5: configure F5OS LAGs

Run in `examples/migration/vlans-from-iseries/lags.tf`.

Goal:

- convert source trunks into F5OS LAGs
- preserve trunk name, mode/type, mapped members, and VLAN assignment

Input source:

- Phase 1 extracted settings JSON
- converted via `scripts/lags-from-iseries.sh`

Dependency:

- depends on Phase 3 VLAN creation
- operationally follows the interface mapping work in Phase 4

### Phase 6: deploy BIG-IP tenants

Run in `examples/migration/vlans-from-iseries/tenant.tf`.

Goal:

- deploy BIG-IP tenant instances on the r-Series platform
- attach migrated VLANs
- configure tenant management connectivity
- size CPU, memory, disk, and nodes appropriately for the target workload

Inputs:

- tenant image identified from Phase 0
- VLAN IDs from Phase 3
- tenant sizing and management decisions supplied by the operator

Dependency:

- directly depends on VLAN creation
- does not declare a hard Terraform dependency on interfaces/LAGs, but those
  should be operationally ready before or by the time the tenant is deployed

Important note:

- tenant sizing is not derivable directly from i-Series source extraction
- it is an operator design decision that must be planned explicitly

## Recommended execution order

### Strict linear view

1. Phase 0: inventory TMOS version and hardware
2. Phase 1: extract system and network settings
3. Phase 2: generate/download UCS backup
4. Licensing phase on the target F5OS platform
5. System-settings phase on the target F5OS platform
6. Phase 3: create VLANs
7. Phase 4: configure interfaces
8. Phase 5: configure LAGs
9. Phase 6: deploy tenants

### Parallelized view

1. Phase 0 in `terraform-provider-bigip`
2. Phase 1 in `terraform-provider-bigip`
3. Phase 2 in `terraform-provider-bigip`
4. In `terraform-provider-f5os`, run these in parallel where desired:
   - licensing phase
   - system-settings phase
   - preparation of tenant sizing/image inputs
5. Run Phase 3 VLAN creation
6. Run Phase 4 interfaces and Phase 5 LAGs
7. Run Phase 6 tenant deployment

## Artifact handoff between repos

### From terraform-provider-bigip to terraform-provider-f5os

Primary handoff artifacts:

- Phase 0 inventory JSON
- Phase 1 extracted settings JSON
- Phase 1 interface/VLAN JSON
- Phase 2 UCS file

These are then consumed by F5OS-side conversion scripts such as:

- `scripts/system-settings-from-iseries.sh`
- `scripts/vlans-from-iseries.sh`
- `scripts/interfaces-from-iseries.sh`
- `scripts/lags-from-iseries.sh`

## Operational caveats

- TMOS interface and trunk naming does not map directly to r-Series naming.
- r-Series port groups are an F5OS-side prerequisite with no TMOS equivalent;
  plan them before interface/LAG rollout.
- Licensing is not migrated from the source device; obtain new registration
  keys for the target r-Series platform.
- Tenant sizing is not discoverable from the source device in a way that can
  be translated automatically; it requires operator planning.
- The tenant's VLAN attachments depend on Phase 3 resources, but correct data
  forwarding still depends on Phases 4 and 5 being operationally correct.

## Suggested operator checklist

1. Run Phase 0 and confirm the target tenant image/version strategy.
2. Run Phase 1 and validate the extracted JSON artifacts.
3. Run Phase 2 and archive the UCS backup securely.
4. Obtain target F5OS/r-Series registration key(s).
5. Prepare F5OS conversion outputs for system settings, VLANs, interfaces,
   and LAGs.
6. Configure target platform licensing and system settings.
7. Create VLANs, then configure interfaces and LAGs.
8. Upload/stage the tenant image selected from Phase 0.
9. Deploy the tenant(s) with the intended VLANs, management IP, and sizing.
10. Perform tenant-level restore/configuration and post-cutover validation.

## Source references

- `terraform-provider-bigip/docs/guides/inventory-tmos-version.md`
- `terraform-provider-bigip/docs/guides/extract-sys-settings.md`
- `terraform-provider-bigip/docs/guides/generate-ucs-backup.md`
- `terraform-provider-bigip/docs/guides/interface-trunk-mapping.md`
- `terraform-provider-f5os/docs/guides/apply-license-from-iseries.md`
- `terraform-provider-f5os/docs/guides/configure-system-settings-from-iseries.md`
- `terraform-provider-f5os/docs/guides/create-vlans-from-iseries.md`
- `terraform-provider-f5os/docs/guides/configure-interfaces-from-iseries.md`
- `terraform-provider-f5os/docs/guides/configure-lags-from-iseries.md`
- `terraform-provider-f5os/docs/guides/deploy-tenants-from-iseries.md`

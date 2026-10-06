---
page_title: "i-Series to r-Series migration operator checklist"
description: |-
  Quick operator checklist for executing an i-Series to r-Series migration using terraform-provider-bigip and terraform-provider-f5os.
---

# i-Series to r-Series migration operator checklist

Use this as the short operational companion to the full
[migration flow](iseries-to-rseries-migration-flow.html).

## Before you begin

1. Confirm source BIG-IP i-Series management access and credentials.
2. Confirm target r-Series/F5OS management access and credentials.
3. Confirm Terraform, `jq`, and required helper tools are installed.
4. Obtain a new target r-Series/F5OS registration key from F5.
5. Decide whether the target is r-Series only or VELOS/F5OS-C, since
   interface naming and platform behavior differ.

## Source-side work in terraform-provider-bigip

1. Run Phase 0 inventory.
Line items:
   - capture TMOS version/build
   - capture hardware/platform info
   - identify the recommended tenant image pattern
   - confirm r-Series supportability for the source TMOS train

2. Run Phase 1 extraction.
Line items:
   - export system settings JSON
   - export interface/VLAN mapping JSON
   - review extracted VLANs, interfaces, trunks, routes, self IPs, users,
     DNS, NTP, and SNMP

3. Run Phase 2 UCS backup.
Line items:
   - generate UCS on the source device
   - download UCS locally
   - store it in a controlled backup location

## Handoff to terraform-provider-f5os

1. Carry forward the following artifacts:
   - Phase 0 inventory output
   - Phase 1 extracted settings JSON
   - Phase 1 interface/VLAN mapping JSON
   - Phase 2 UCS archive

2. Generate the F5OS-side tfvars/config inputs using the F5OS migration
   scripts.

## Target-side work in terraform-provider-f5os

1. Apply the target F5OS/r-Series license.
2. Apply system settings if needed.
3. Create VLANs.
4. Configure interfaces.
5. Configure LAGs.
6. Stage the tenant image selected from Phase 0.
7. Deploy tenants with planned sizing, management addressing, and VLANs.

## Validation checkpoints

1. Verify the target license is active.
2. Verify platform DNS/NTP/SNMP and local users are correct.
3. Verify VLAN IDs and names match the migration plan.
4. Verify interface and LAG mappings match the intended r-Series names.
5. Verify the tenant image version matches the inventory decision.
6. Verify tenant management IP reachability after deployment.
7. Verify data-plane VLAN attachments align with the intended trunks/LAGs.

## Caveats to check explicitly

1. Do not reuse the i-Series TMOS license key on the r-Series target.
2. Do not assume TMOS interface names map directly to r-Series names.
3. Do not assume tenant sizing can be inferred automatically from the source.
4. Do not treat UCS backup as a substitute for the structured extraction
   artifacts; both are useful for different parts of the migration.

## Related guides

1. [Overall i-Series to r-Series migration flow](iseries-to-rseries-migration-flow.html)
2. [Inventorying TMOS version and hardware](inventory-tmos-version.html)
3. [Extracting i-Series system settings](extract-sys-settings.html)
4. [Generating and downloading a UCS backup](generate-ucs-backup.html)

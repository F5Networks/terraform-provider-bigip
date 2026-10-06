---
page_title: "Extracting i-Series system settings for F5OS/r-Series migration"
description: |-
  How to use scripts/extract-sys-settings.sh to extract DNS, NTP, SNMP, syslog, auth (LDAP/RADIUS/TACACS+), local user, VLAN (including VLAN-to-interface tagging), self IP, route, physical interface, and trunk (LAG) settings from a BIG-IP i-Series device for use in an r-Series (F5OS) migration.
---

# Extracting i-Series system settings for F5OS/r-Series migration

`scripts/extract-sys-settings.sh` is Phase 1 of an i-Series -> r-Series (F5OS)
migration workflow (see [Inventorying TMOS version and
hardware](inventory-tmos-version.html) for the prerequisite Phase 0 step of
identifying the correct F5OS tenant image before starting): it reads DNS,
NTP, SNMP, syslog, auth (LDAP/RADIUS/TACACS+), local user-account, VLAN,
self IP, route, physical interface, and trunk (LAG) settings from a live
BIG-IP i-Series device and writes them out as two JSON documents for a
downstream Phase 2 process to consume (e.g. applying equivalent settings
at the F5OS layer): a general settings dump, and a derived
`interface_vlans` map suited for `f5os_interface`/`f5os_lag`'s
`native_vlan`/`trunk_vlans` arguments (see [Interface VLAN
mapping](#interface-vlan-mapping) below). See also the [interface/trunk
naming-convention guide](interface-trunk-mapping.html) when mapping the
interface/trunk portion of the output to F5OS, and [Generating and
downloading a UCS backup](generate-ucs-backup.html) for the separate Phase 2
step that captures the full BIG-IP configuration as a UCS archive for
restoration on the r-Series tenant.

For the full cross-repo phase sequence, see the [overall i-Series to
r-Series migration flow](iseries-to-rseries-migration-flow.html).

This is a standalone shell script, not a Terraform resource or data source
you write into your own configuration -- run it directly:

```sh
BIGIP_HOST=10.1.1.1 BIGIP_USER=admin BIGIP_PASSWORD=secret \
  ./scripts/extract-sys-settings.sh [output.json] [interface_vlans.json] [workdir]
```

## Why a script instead of `terraform plan`/data sources

Most of these settings have no read-only data source at all in this
provider. Several of them also don't expose a fixed, discoverable ID
through the provider alone -- for example `bigip_sys_dns`'s ID is the
current value of its own `description` field, which is arbitrary. The
script queries the BIG-IP iControl REST API directly (via `curl`) purely
to *discover which objects exist and what their import IDs are*; the
actual extraction of every other field still goes through each resource's
`Read` function via `terraform import`.

The networking settings (`bigip_net_vlans`, `bigip_net_selfips`,
`bigip_net_routes`, `bigip_net_interfaces`, `bigip_net_trunks`) are the
exception: these already exist as proper data sources (there's no
meaningful "resource" to import for a bulk listing), so the script just
declares them directly in the generated config and populates them via a
`terraform apply -refresh-only` after all the `import` calls -- see
[Provider binary resolution](#provider-binary-resolution) for why they're
generated conditionally alongside the auth/syslog resources.

## What gets extracted

| Category | Resource/data source | Notes |
|---|---|---|
| DNS | `bigip_sys_dns` | Name servers and search domains |
| NTP | `bigip_sys_ntp` | NTP servers and timezone |
| SNMP | `bigip_sys_snmp`, `bigip_sys_snmp_traps` | Contact/location/allowed addresses, plus one entry per configured trap |
| Syslog | `bigip_sys_syslog` | Log levels and remote syslog servers |
| LDAP auth | `bigip_auth_ldap` | Only extracted if LDAP is configured (`auth/ldap` has an entry) |
| RADIUS auth | `bigip_auth_radius`, `bigip_auth_radius_server` | System config plus one entry per configured RADIUS server |
| TACACS+ auth | `bigip_auth_tacacs` | Only extracted if TACACS+ is configured |
| Users | `bigip_auth_user` | One entry per local user account, including partition/role assignments |
| VLANs | `bigip_net_vlans` (data source) | Tag, MTU, and per-VLAN tagged/untagged interface or trunk membership; also the source for the derived `interface_vlans` map below |
| Self IPs | `bigip_net_selfips` (data source) | Address, VLAN, traffic group, port lockdown (`allow_service`) |
| Routes | `bigip_net_routes` (data source) | Network, gateway, interface/tunnel, blackhole routes |
| Physical interfaces | `bigip_net_interfaces` (data source) | Status, negotiated speed/media, MTU, MAC, bundle support, LLDP/STP config |
| Trunks (LAGs) | `bigip_net_trunks` (data source) | Member interfaces, LACP mode/timeout, distribution hash |

## Prerequisites

- `terraform` CLI, `curl`, and `jq` in `PATH`.
- Network access from wherever you run the script to the i-Series
  management address.
- A BIG-IP account with read access to `/mgmt/tm/sys/*`, `/mgmt/tm/auth/*`,
  and `/mgmt/tm/net/*`.

## Environment variables

Same convention as the provider itself:

| Variable | Required | Default | Description |
|---|---|---|---|
| `BIGIP_HOST` | yes | -- | Management address of the i-Series device |
| `BIGIP_USER` | yes | -- | BIG-IP username |
| `BIGIP_PASSWORD` | yes | -- | BIG-IP password |
| `BIGIP_PORT` | no | `443` | Management port |
| `BIGIP_VERIFY_CERT_DISABLE` | no | `true` | Set to `false` to enforce TLS certificate verification |
| `BIGIP_PROVIDER_BINARY` | no | -- | Path to an already-built `terraform-provider-bigip` binary (see [Provider binary resolution](#provider-binary-resolution)) |
| `TEEM_DISABLE` | no | `true` | Passed through to the provider |

## Usage

```sh
./scripts/extract-sys-settings.sh [output.json] [interface_vlans.json] [workdir]
```

- `output.json` (optional, positional) -- where to write the extracted JSON.
  Defaults to `extracted-sys-settings.json` in the current directory.
- `interface_vlans.json` (optional, positional) -- where to write the
  derived `interface_vlans` map (see [Interface VLAN
  mapping](#interface-vlan-mapping) below). Defaults to
  `interface-vlans.json` in the current directory.
- `workdir` (optional, positional) -- the Terraform working directory the
  script creates and imports into. Defaults to a fresh `mktemp -d`. The
  script leaves this directory in place afterward (containing `main.tf`,
  `.terraform.lock.hcl`, and `terraform.tfstate`) so you can inspect it or
  re-run `terraform show` against it later; nothing is cleaned up
  automatically except the temporary provider plugin/CLI-config files (see
  below).

Example:

```sh
BIGIP_HOST=10.1.1.20 BIGIP_USER=admin BIGIP_PASSWORD='...' \
  ./scripts/extract-sys-settings.sh migration-output/dc1-switch01.json \
    migration-output/dc1-switch01-interface-vlans.json /tmp/dc1-switch01-tf
```

## Output format

`output.json` is the `values.root_module.resources` array from
`terraform show -json`: one entry per extracted object, each with `.type`
(the `bigip_*` resource type), `.address` (the Terraform resource address,
e.g. `bigip_auth_user.extracted_netops_admin`), and `.values` (every
attribute the resource's `Read` function populated, straight from the live
device). Sensitive fields (`bind_pw`, `secret`, `password`) are always
`null` -- BIG-IP never returns these on read, so they can't be extracted;
Phase 2 must supply them separately.

Example (abbreviated):

```json
[
  {
    "address": "bigip_sys_dns.extracted",
    "type": "bigip_sys_dns",
    "values": {
      "id": "configured-by-dhcp",
      "name_servers": ["8.8.8.8", "8.8.4.4"],
      "search": ["example.com", "corp.example.com"]
    }
  },
  {
    "address": "bigip_auth_user.extracted_netops_admin",
    "type": "bigip_auth_user",
    "values": {
      "name": "netops-admin",
      "description": "NetOps Administrator",
      "shell": "bash",
      "password": null,
      "partition_access": [{"partition": "all-partitions", "role": "admin"}]
    }
  },
  {
    "address": "data.bigip_net_trunks.extracted",
    "type": "bigip_net_trunks",
    "values": {
      "trunks": [
        {
          "name": "test-lag1",
          "interfaces": ["1.1", "1.2"],
          "lacp": "disabled",
          "lacp_mode": "active",
          "distribution_hash": "src-dst-ipport"
        }
      ]
    }
  },
  {
    "address": "data.bigip_net_vlans.extracted",
    "type": "bigip_net_vlans",
    "values": {
      "vlans": [
        {
          "name": "vlan10",
          "tag": 10,
          "interfaces": [
            {"name": "1.1", "tagged": false},
            {"name": "test-lag1", "tagged": true}
          ]
        }
      ]
    }
  }
]
```

Note the data source entries use the `data.<type>.<name>` address form
(e.g. `data.bigip_net_trunks.extracted`) and their `.values` is a single
attribute (e.g. `vlans`, `self_ips`, `routes`, `interfaces`, `trunks`)
containing the *list* of every object of that kind -- unlike the resource
entries above, which are one entry per object.

## Interface VLAN mapping

`bigip_net_vlans` reports membership VLAN-first: each VLAN lists its own
tagged/untagged interfaces and trunks. F5OS's `f5os_interface`/`f5os_lag`
resources (in the separate `F5Networks/f5os` provider) configure VLANs the
other way around -- interface-first, directly on the interface/LAG
resource, via a `native_vlan` argument (the single untagged/native VLAN
ID) and a `trunk_vlans` argument (a list of tagged VLAN IDs). After the
refresh-only apply, the script inverts `bigip_net_vlans`' output into that
interface-centric shape with a `jq` transform and writes it to
`interface_vlans.json`, so it can be consumed directly as those two
arguments once interface/trunk names are translated to F5OS naming (see
the [interface/trunk naming-convention guide](interface-trunk-mapping.html)).

Example:

```json
{
  "interface_vlans": {
    "1.1": {
      "type": "interface",
      "native_vlan": 10,
      "trunk_vlans": []
    },
    "test-lag1": {
      "type": "trunk",
      "native_vlan": null,
      "trunk_vlans": [10, 20]
    }
  },
  "ambiguous_native_vlans": []
}
```

- Keyed by TMOS interface/trunk name (translate to F5OS naming per the
  [interface/trunk naming-convention guide](interface-trunk-mapping.html)
  before applying to `f5os_interface`/`f5os_lag`).
- `type` is `"trunk"` for any name that also appears as a trunk in
  `bigip_net_trunks`, else `"interface"`.
- `native_vlan` is `null` when the interface/trunk has no untagged VLAN.
- An interface/trunk with **more than one** untagged VLAN is an invalid
  TMOS configuration the script can't silently resolve into a single
  `native_vlan` -- rather than guessing, that entry's `native_vlan` is
  left `null` and it's reported separately under
  `ambiguous_native_vlans` (a list of `{"interface": ..., "native_vlans":
  [...]}`) for manual review before applying to F5OS.
- If `bigip_net_vlans` wasn't extracted (registry fallback -- see
  [Registry fallback](#registry-fallback) below), `interface_vlans.json`
  is still written, just empty (`{"interface_vlans": {},
  "ambiguous_native_vlans": []}`), rather than not being written at all.

## Provider binary resolution

Several resources/data sources this script extracts (`bigip_sys_syslog`,
`bigip_auth_ldap`, `bigip_auth_radius`, `bigip_auth_radius_server`,
`bigip_auth_tacacs`, `bigip_auth_user`, `bigip_net_vlans`,
`bigip_net_selfips`, `bigip_net_routes`, `bigip_net_interfaces`,
`bigip_net_trunks`) may not exist yet in the version of `F5Networks/bigip`
published on the Terraform Registry. To make sure the script always talks
to a provider build that actually has them, it tries, in order:

1. **`BIGIP_PROVIDER_BINARY`**, if set -- copied in as-is.
2. **`go build`** from the repo this script lives in, resolved from the
   script's own file path (not your current working directory), so this
   works no matter where you invoke it from -- as long as the script is
   still sitting inside its original `scripts/` checkout.

Whichever binary is used, the script wires it up automatically via a
generated Terraform CLI [`dev_overrides`
config](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)
(`TF_CLI_CONFIG_FILE`), so you don't need to set that up yourself.

### Registry fallback

If neither option above is available -- `go` isn't installed, the script
has been copied out on its own (so its containing repo checkout is
missing or isn't a valid module), or the local build fails for any reason
-- the script falls back to letting `terraform init` install the **latest
published** `F5Networks/bigip` provider from the Terraform Registry
instead, with a warning printed to stderr.

In fallback mode, `bigip_sys_syslog`, every `bigip_auth_*` resource, and
all five `bigip_net_*` data sources (`bigip_net_vlans`, `bigip_net_selfips`,
`bigip_net_routes`, `bigip_net_interfaces`, `bigip_net_trunks`) are
**omitted from the extraction entirely** (not attempted and skipped one by
one): the script prints a single upfront note and only extracts
DNS/NTP/SNMP/SNMP-traps, which exist in every published version. This is
intentional, not a partial failure -- declaring a resource or data source
block of a type the loaded provider doesn't recognize can make *every*
`terraform import`/`apply` in the run fail (Terraform validates the entire
configuration graph, not just the resource/data source an operation
targets), so the script avoids declaring those blocks at all rather than
risk the whole run failing on one unsupported type. Since `bigip_net_vlans`
is omitted too, `interface_vlans.json` is written out empty in this mode
(see [Interface VLAN mapping](#interface-vlan-mapping) above) rather than
not being written at all. If you need the syslog/auth/networking
categories and don't have `go`/a repo checkout available, build (or
obtain) a provider binary elsewhere and pass it via `BIGIP_PROVIDER_BINARY`.

A separate, unrelated warning class -- "N object(s) unexpectedly failed to
import" -- covers genuine per-object import failures (e.g. an object was
deleted from the device mid-run), not the registry-fallback skip above.

## Cleanup

The script always removes its temporary provider plugin directory and
generated `dev_overrides` CLI config on exit (success or failure), via a
`trap`. It does **not** remove the Terraform working directory
(`workdir`) or its state -- that's left in place intentionally so you can
inspect exactly what was imported. Delete it yourself once you're done
with it; it references live resources only in local Terraform state, not
in any shared backend, so removing the directory has no effect on the
BIG-IP device itself.

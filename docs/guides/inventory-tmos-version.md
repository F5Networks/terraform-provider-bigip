---
page_title: "Inventorying TMOS version and hardware for F5OS/r-Series migration"
description: |-
  How to use scripts/inventory-tmos-version.sh to extract the exact TMOS version, build, and platform/hardware info from a BIG-IP i-Series device via the bigip_command resource, and derive the recommended F5OS tenant image filename for an r-Series migration.
---

# Inventorying TMOS version and hardware for F5OS/r-Series migration

`scripts/inventory-tmos-version.sh` is Phase 0 (a prerequisite step) of an
i-Series -> r-Series (F5OS) migration workflow: it runs `tmsh show sys
version` and `tmsh show sys hardware` on a live BIG-IP i-Series device via
this provider's [`bigip_command`](../resources/bigip_command.html)
resource, captures the platform model/hardware info and exact TMOS
version/build, and derives the recommended F5OS tenant image filename (and
whether the running TMOS version is even supported on rSeries at all) --
so the correct tenant image can be identified and staged *before* the rest
of the migration ([Phase 1](extract-sys-settings.html): extracting system
settings; [Phase 2](generate-ucs-backup.html): generating a UCS backup)
begins.

For the full cross-repo phase sequence, see the [overall i-Series to
r-Series migration flow](iseries-to-rseries-migration-flow.html).

This is a standalone shell script, not a Terraform resource or data
source you write into your own configuration -- run it directly:

```sh
BIGIP_HOST=10.1.1.1 BIGIP_USER=admin BIGIP_PASSWORD=secret \
  ./scripts/inventory-tmos-version.sh [output.json] [workdir]
```

## Why `bigip_command` instead of a dedicated data source

This provider has no `bigip_sys_version` or `bigip_sys_hardware` data
source (unlike `bigip_net_interfaces`/`bigip_net_trunks`/`bigip_net_vlans`,
used by [`extract-sys-settings.sh`](extract-sys-settings.html)), and the
vendored `go-bigip` client has no binding for either
`/mgmt/tm/sys/version` or `/mgmt/tm/sys/hardware`. Rather than wait on a
typed API, this script runs the two `tmsh show` commands the acceptance
criteria call for directly via `bigip_command` and parses their plain-text
output:

- `tmsh show sys version`'s tabular format (`Product`/`Version`/`Build`/
  `Edition`/`Date` under a `Main Package` heading) has been stable since
  early TMOS releases and is the same information AskF5 solution articles
  reference when asking users to identify their running version, so it's
  parsed with a portable regex.
- `tmsh show sys hardware`'s output varies more across platforms (VE vs.
  appliance vs. chassis blade all print different fields), so it's parsed
  generically: every `<key><2+ spaces><value>` line is captured into a
  flat JSON object, regardless of which exact fields a given platform
  happens to print. The raw text is also always preserved verbatim, so
  nothing the generic parse misses is lost.

Unlike [`generate-ucs-backup.sh`](generate-ucs-backup.html)'s `tmsh save
sys ucs`, both `show` commands here are pure reads that complete in a
fraction of a second, so none of that script's backgrounding-and-poll
machinery is needed -- a single `terraform apply` is sufficient.

## Recommended tenant image filename

F5's rSeries tenant images are named:

```
BIGIP-<version>-<build>.<TYPE>-F5OS.<format>.bundle
```

for example `BIGIP-17.1.1.2-0.0.10.ALL-F5OS.qcow2.zip.bundle` or
`BIGIP-21.1.0-0.0.38.ALL-F5OS.tar.bundle` (see F5's rSeries Planning
Guide, [Deploying an rSeries BIG-IP Tenant -- Tenant Image
Types](https://clouddocs.f5.com/training/community/rseries-training/html/rseries_deploying_a_tenant.html#tenant-image-types)).

The `<version>-<build>` portion comes directly from the parsed `tmsh show
sys version` output, so this script constructs it exactly. It does
**not** guess the remaining two segments, since neither is derivable from
the source i-Series device alone:

- **`<TYPE>`** is one of `T1`, `T2`, `T4`, or `ALL`, trading tenant disk
  footprint against in-place upgrade support -- this is an operator
  sizing decision (see [F5 solution article
  K45191957](https://my.f5.com/manage/s/article/K45191957) for guidance).
  `T1` supports no in-place upgrades at all; `T2`/`ALL`/`T4` do, with
  increasing default disk allocation.
- **`<format>`** depends on which F5 image-signing generation was used
  for the specific build you download from downloads.f5.com: F5
  transitioned from `.qcow2.zip.bundle` to `.tar.bundle` in October 2025
  (see [F5 solution article
  K000157005](https://my.f5.com/manage/s/article/K000157005)).

The script's output therefore includes a filename *pattern* with a
`<TYPE>` placeholder and a note about both possible extensions, not a
single fully-resolved filename guaranteed to exist verbatim on
downloads.f5.com.

## rSeries TMOS version support check

Per [F5 solution article
K86001294](https://my.f5.com/manage/s/article/K86001294) (F5OS
hardware/software support matrix), rSeries does not support TMOS 16.0.x,
16.1.x, or 17.0.x **at all** as an initial tenant image -- only 15.1.x
(minimum patch level varies by target rSeries model), or 17.1.x/21.x and
later, are valid. Upgrades to 16.x/17.0.x happen only *within* an
already-deployed tenant, never as the initial image loaded into F5OS.

This script performs a best-effort major.minor classification of the
extracted TMOS version against that matrix and includes the result (and
a warning if unsupported) in its output. This check cannot account for
platform-specific minimum patch levels (e.g. r5000/r10000 require
15.1.5+, r2000/r4000 require 15.1.6+) -- always confirm against
K86001294 for your specific target rSeries model before relying on this
alone.

## Prerequisites

- `terraform` CLI and `jq` in `PATH`.
- Network access from wherever you run the script to the i-Series
  management address.
- A BIG-IP account with permission to run `tmsh show sys version` and
  `tmsh show sys hardware` (available to any authenticated account by
  default).

## Environment variables

Same convention as the provider itself:

| Variable | Required | Default | Description |
|---|---|---|---|
| `BIGIP_HOST` | yes | -- | Management address of the i-Series device |
| `BIGIP_USER` | yes | -- | BIG-IP username |
| `BIGIP_PASSWORD` | yes | -- | BIG-IP password |
| `BIGIP_PORT` | no | `443` | Management port |
| `BIGIP_VERIFY_CERT_DISABLE` | no | `true` | Set to `false` to enforce TLS certificate verification |
| `BIGIP_TRUSTED_CERT_PATH` | no | -- | Path to a CA bundle PEM to verify the device's certificate against, when `BIGIP_VERIFY_CERT_DISABLE=false` |
| `BIGIP_PROVIDER_BINARY` | no | -- | Path to an already-built `terraform-provider-bigip` binary (see [Provider binary resolution](#provider-binary-resolution)) |
| `TEEM_DISABLE` | no | `true` | Passed through to the provider |

## Usage

```sh
./scripts/inventory-tmos-version.sh [output.json] [workdir]
```

- `output.json` (optional, positional) -- where to write the extracted
  JSON. Defaults to `tmos-inventory.json` in the current directory.
- `workdir` (optional, positional) -- the Terraform working directory
  used to run the `bigip_command` resources. Defaults to a fresh
  `mktemp -d`.

Example:

```sh
BIGIP_HOST=10.1.1.20 BIGIP_USER=admin BIGIP_PASSWORD='...' \
  ./scripts/inventory-tmos-version.sh migration-output/dc1-switch01-inventory.json /tmp/dc1-switch01-tf
```

## Output format

A single JSON document (default `tmos-inventory.json`) with four
top-level keys:

```json
{
  "version": {
    "product": "BIG-IP",
    "version": "17.1.1.2",
    "build": "0.0.10",
    "edition": "Point Release 2",
    "date": "Mon Jun 24 08:14:42 PDT 2024",
    "raw_output": "Sys::Version\nMain Package\n  Product     BIG-IP\n..."
  },
  "hardware": {
    "fields": {
      "Name": "BIG-IP r10600",
      "BIOS Revision": "1.0.0",
      "Base MAC": "00:94:a1:12:34:56",
      "Type": "r10600",
      "Chassis Serial": "f5-abcd-1234",
      "Appliance Serial": ["blade-serial-1", "blade-serial-2"]
    },
    "raw_output": "Sys::Hardware\nPlatform\n  Name  BIG-IP r10600\n..."
  },
  "recommended_tenant_image": {
    "version_build": "17.1.1.2-0.0.10",
    "filename_pattern": "BIGIP-17.1.1.2-0.0.10.<TYPE>-F5OS.<qcow2.zip|tar>.bundle",
    "type_options": ["T1", "T2", "T4", "ALL"],
    "note": "Replace <TYPE> with your chosen tenant image type ..."
  },
  "rseries_version_support": {
    "supported": true,
    "note": "TMOS 17.1.1.2 appears to be an rSeries-supported initial tenant image train per K86001294 ...",
    "reference": "K86001294: F5OS hardware/software support matrix (https://my.f5.com/manage/s/article/K86001294)"
  }
}
```

- `version.raw_output` / `hardware.raw_output` always contain the
  complete, unmodified `tmsh` output, regardless of how much of it the
  structured `version`/`hardware.fields` parsing was able to extract --
  useful for manual inspection if a platform's hardware output doesn't
  parse the way you expect.
- `hardware.fields` is best-effort and platform-dependent: a VE guest
  won't have a `Chassis Serial` field at all, for instance, since it has
  no physical chassis. Missing fields simply aren't present in the
  object rather than appearing as `null`. On multi-component platforms
  (chassis blades, power supplies, fans), the same key can legitimately
  appear once per component (e.g. `Appliance Serial` per blade, `Status`
  per power supply) -- when a key repeats, its value in `hardware.fields`
  is a JSON array of every occurrence (in the order `tmsh` printed them)
  rather than just the last one seen; a key that appears only once is
  still a plain scalar value, not a single-element array.
- `recommended_tenant_image` and `rseries_version_support` are described
  in detail above.

If `tmsh show sys version` returns no output at all, or its `Version`/
`Build` fields can't be parsed, the script exits with an error rather
than writing a partial/misleading `output.json`. `tmsh show sys
hardware` returning no output only produces a warning (empty
`hardware.fields`), since the version information alone is still useful.

## Provider binary resolution

Same mechanism as
[`extract-sys-settings.sh`](extract-sys-settings.html#provider-binary-resolution),
except `bigip_command` has been part of every published provider version
for a long time (same reasoning as
[`generate-ucs-backup.sh`](generate-ucs-backup.html#provider-binary-resolution)),
so falling back to the published registry provider is always safe here.
The script tries, in order:

1. **`BIGIP_PROVIDER_BINARY`**, if set -- copied in as-is.
2. **`go build`** from the repo this script lives in, resolved from the
   script's own file path (not your current working directory).
3. **The published `F5Networks/bigip` provider from the Terraform
   Registry**, if neither of the above is available (`go` not installed,
   the script copied out of its repo checkout, or the local build fails).

Whichever binary is used, the script wires it up automatically via a
generated Terraform CLI
[`dev_overrides` config](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)
(`TF_CLI_CONFIG_FILE`).

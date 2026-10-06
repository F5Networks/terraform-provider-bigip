---
page_title: "Generating and downloading a UCS backup for F5OS/r-Series migration"
description: |-
  How to use scripts/generate-ucs-backup.sh to generate a full UCS (User Configuration Set) archive on a BIG-IP i-Series device via the bigip_command resource, and download it locally via the iControl REST file-transfer API for upload to an r-Series tenant.
---

# Generating and downloading a UCS backup for F5OS/r-Series migration

`scripts/generate-ucs-backup.sh` is Phase 2 of an i-Series -> r-Series
(F5OS) migration workflow (see [Inventorying TMOS version and
hardware](inventory-tmos-version.html) for the prerequisite Phase 0 step,
and [Extracting i-Series system settings](extract-sys-settings.html) for
Phase 1): it generates a full UCS
(User Configuration Set) backup archive on a live BIG-IP i-Series device
via this provider's [`bigip_command`](../resources/bigip_command.html)
resource (`tmsh save sys ucs <filename>`), then downloads the resulting
archive to local disk via the iControl REST bulk file-transfer API
(`/mgmt/shared/file-transfer/ucs-downloads/<filename>`), ready for upload
to an r-Series tenant in Phase 3.

For the full cross-repo phase sequence, see the [overall i-Series to
r-Series migration flow](iseries-to-rseries-migration-flow.html).

This is a standalone shell script, not a Terraform resource or data
source you write into your own configuration -- run it directly:

```sh
BIGIP_HOST=10.1.1.1 BIGIP_USER=admin BIGIP_PASSWORD=secret \
  ./scripts/generate-ucs-backup.sh [ucs_filename] [output_dir] [workdir]
```

## Why not just `terraform apply` a plain `bigip_command` resource

A synchronous `tmsh save sys ucs` of any non-trivial configuration size
can run long enough (minutes, not seconds) that the device's own
`restjavad` request-processing timeout expires *before* the save
finishes. When that happens, `restjavad` kills the child `tmsh` process
outright (visible in `/var/log/restjavad.0.log` as `Killing child
process ...`) -- the UCS file is never written, and Terraform sees a
failed apply. This is a device-side timeout, independent of any
client-side HTTP timeout setting, so it can't be worked around from the
provider or `curl` side alone.

To avoid this, the `bigip_command` `commands` entry the script generates
backgrounds the actual `tmsh save` via a nested
`run util bash -c 'nohup bash -c "tmsh save sys ucs ... ; echo
UCS_SAVE_EXIT_CODE:$?" &'`, so the outer REST call that `terraform apply`
waits on returns almost immediately -- immune to `restjavad`'s internal
timeout. The script then polls a completion sentinel in the save's own
log file on the device (see [What the script
does](#what-the-script-does) below) rather than the outer `terraform
apply` call. This is still fundamentally "generated via `bigip_command`
running `tmsh save sys ucs`", just made reliable for archives that take
longer than a typical HTTP request timeout to produce.

## Prerequisites

- `terraform` CLI, `curl`, and `jq` in `PATH`.
- Network access from wherever you run the script to the i-Series
  management address.
- A BIG-IP account with permission to run `tmsh save sys ucs` and
  `tmsh delete sys ucs`, to read the save's log file under `/tmp` via
  `util/bash`, and to read/download from
  `/mgmt/shared/file-transfer/ucs-downloads/`.

## Environment variables

Same convention as the provider itself, plus a few script-specific knobs:

| Variable | Required | Default | Description |
|---|---|---|---|
| `BIGIP_HOST` | yes | -- | Management address of the i-Series device |
| `BIGIP_USER` | yes | -- | BIG-IP username |
| `BIGIP_PASSWORD` | yes | -- | BIG-IP password |
| `BIGIP_PORT` | no | `443` | Management port |
| `BIGIP_VERIFY_CERT_DISABLE` | no | `true` | Set to `false` to enforce TLS certificate verification -- applies to this script's own `curl` calls (pre-check, pre-deletion, polling, and chunked download) in addition to the Terraform provider's HTTP client |
| `BIGIP_TRUSTED_CERT_PATH` | no | -- | Path to a CA bundle PEM to verify the device's certificate against, when `BIGIP_VERIFY_CERT_DISABLE=false` and the certificate isn't already trusted by the system CA store. Passed to `curl` via `--cacert` and to the provider via its own support for this variable |
| `BIGIP_PROVIDER_BINARY` | no | -- | Path to an already-built `terraform-provider-bigip` binary (see [Provider binary resolution](#provider-binary-resolution)) |
| `TEEM_DISABLE` | no | `true` | Passed through to the provider |
| `UCS_POLL_INTERVAL_SECONDS` | no | `15` | How often to poll for save completion |
| `UCS_POLL_TIMEOUT_SECONDS` | no | `1800` | Give up waiting for the save after this long (30 minutes) |
| `UCS_CHUNK_SIZE_BYTES` | no | `1048576` | Download chunk size (1 MiB) |

## Usage

```sh
./scripts/generate-ucs-backup.sh [ucs_filename] [output_dir] [workdir]
```

- `ucs_filename` (optional, positional) -- name of the UCS archive on the
  device. Defaults to `migration-backup.ucs`; a `.ucs` extension is added
  automatically if not already present.
- `output_dir` (optional, positional) -- where to write the downloaded
  archive, as `<output_dir>/<ucs_filename>`. Defaults to the current
  directory.
- `workdir` (optional, positional) -- the Terraform working directory
  used to run the `bigip_command` resource. Defaults to a fresh
  `mktemp -d`.

Example:

```sh
BIGIP_HOST=10.1.1.20 BIGIP_USER=admin BIGIP_PASSWORD='...' \
  ./scripts/generate-ucs-backup.sh migration-backup.ucs migration-output/dc1-switch01 /tmp/dc1-switch01-ucs-tf
```

## What the script does

1. Checks whether a UCS archive with the same name already exists on the
   device via `GET /mgmt/tm/sys/ucs`; if so, deletes it first (`tmsh
   delete sys ucs <name>`) so the subsequent save starts clean, since
   some TMOS versions prompt/fail on re-saving over an existing name.
   This check (and the completion check in step 3) matches against
   whichever of `apiRawValues.filename`, a top-level `filename`, or
   `name` a given TMOS version actually populates for `sys/ucs` list
   items, since this has been observed to vary across versions.
2. Generates a minimal `main.tf` declaring a single `bigip_command`
   resource whose `commands` entry runs the backgrounded `tmsh save sys
   ucs` (see [above](#why-not-just-terraform-apply-a-plain-bigip_command-resource)),
   then runs `terraform init` and `terraform apply -auto-approve` against
   it. The apply itself only launches the background save -- it does not
   wait for the save to finish.
3. Polls the save's own log file on the device (via `util/bash`, every
   `UCS_POLL_INTERVAL_SECONDS`) for a `UCS_SAVE_EXIT_CODE:<n>` sentinel
   line that the backgrounded command appends only once `tmsh save`
   itself has actually exited. This is the authoritative completion
   signal -- checking only whether the archive appears in
   `GET /mgmt/tm/sys/ucs` isn't sufficient, since on large saves TMOS can
   create the destination file before archive assembly is actually
   finished, which would otherwise race the chunked download in step 4
   against an in-progress writer. A non-zero exit code, or the sentinel
   never appearing within `UCS_POLL_TIMEOUT_SECONDS`, fails the script
   with an error pointing at `/var/log/restjavad.0.log` and the save's
   log file on the device. Once the sentinel reports success, the script
   also double-checks the archive actually appears in
   `GET /mgmt/tm/sys/ucs` as a sanity check before downloading it.
4. Downloads the archive from
   `GET /mgmt/shared/file-transfer/ucs-downloads/<ucs_filename>` in
   `UCS_CHUNK_SIZE_BYTES`-sized chunks using the `Content-Range` request
   header (BIG-IP's documented bulk file-transfer pattern), discovering
   the archive's true total size from the response `Content-Range`
   header as it goes. Chunks are written to a `.tmp` file alongside the
   final destination, not to `<output_dir>/<ucs_filename>` directly.
5. Verifies the downloaded file's size on disk matches the size BIG-IP
   reported, then renames the `.tmp` file into place at
   `<output_dir>/<ucs_filename>`. If any step of the download or the
   size check fails, the `.tmp` file (and other temporary files) are
   removed and any file that already existed at
   `<output_dir>/<ucs_filename>` from a previous run is left untouched,
   rather than being truncated or overwritten with a partial download.

## Output

The downloaded UCS archive at `<output_dir>/<ucs_filename>`, verified
against the size BIG-IP itself reports for the archive -- ready for
upload to the r-Series tenant in Phase 3 of the migration workflow. The
Terraform working directory (`workdir`) is left in place afterward for
inspection/debugging; nothing in it is cleaned up automatically except
the temporary provider plugin/CLI-config files and download staging
files (see [Provider binary resolution](#provider-binary-resolution)).

## Provider binary resolution

Same mechanism as [`extract-sys-settings.sh`](extract-sys-settings.html#provider-binary-resolution),
except `bigip_command` has been part of every published provider version
for a long time, so falling back to the published registry provider is
always safe here -- there's no risk of hitting a resource type the
registry version doesn't support. The script tries, in order:

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

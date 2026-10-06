---
page_title: "Interface and trunk naming: TMOS (i-Series) vs F5OS (r-Series/VELOS)"
description: |-
  Naming-convention differences between TMOS physical interfaces/trunks and F5OS interfaces/LAGs, for mapping i-Series configuration to r-Series or VELOS.
---

# Interface and trunk naming: TMOS (i-Series) vs F5OS (r-Series/VELOS)

This guide documents the naming-convention differences between TMOS (BIG-IP i-Series) physical interfaces and trunks, and their F5OS equivalents, for use alongside the [`bigip_net_interfaces`](../data-sources/bigip_net_interfaces.html) and [`bigip_net_trunks`](../data-sources/bigip_net_trunks.html) (or, for a single already-known trunk name, [`bigip_net_trunk`](../data-sources/bigip_net_trunk.html)) data sources when planning an i-Series -> r-Series (or VELOS) migration.

~> **Important** The mapping is **not 1:1** in either direction: interface names differ, LAG configuration is structured differently (see [Trunk -> LAG mapping](#trunk-lag-mapping-is-not-11) below), and F5OS requires an additional up-front step (port groups) with no TMOS equivalent. Treat everything below as a starting point for manual/scripted mapping, not something to blindly substitute name-for-name.

## Interface naming

| Platform | Convention | Example | Notes |
|---|---|---|---|
| TMOS (i-Series, VE) | `<blade>.<port>` | `1.1`, `1.2` | Single-appliance platforms (i-Series, VE) are always blade `1`; the blade number only varies on chassis platforms (VIPRION). `mgmt` is the dedicated management interface, named literally, not numbered. |
| F5OS-A (rSeries: r2000/r4000/r5000/r10000/r12000) | `<port>.0` | `1.0`, `2.0`, ... `20.0` | **No blade prefix** -- rSeries appliances are single-node, so there's nothing to prefix. Port count varies by model: r2000/r4000 have 8 ports (`1.0`-`8.0`), r5000 has 10 (`1.0`-`10.0`), r10000/r12000 have 20 (`1.0`-`20.0`). |
| F5OS-C (VELOS, chassis-based) | `<blade>/<port>.<subport>` | `1/1.0` | VELOS chassis have multiple blades, so the blade number is a real, meaningful prefix here (unlike TMOS single-appliance platforms, where it's always `1`). `.<subport>` matters when a high-speed port is broken out into multiple lower-speed logical ports. |

Source: [F5 rSeries Planning Guide -- Initial Setup of the rSeries Network Layer](https://clouddocs.f5.com/training/community/rseries-training/html/initial_setup_of_rseries_network_layer.html#network-settings-interfaces).

If your migration target is **rSeries** (the common i-Series replacement), map `bigip_net_interfaces` entries to the `<port>.0` convention, not `<blade>/<port>.<subport>` -- despite what F5OS documentation for VELOS uses, a single rSeries appliance has no blade prefix. If your target is **VELOS**, use the `<blade>/<port>.<subport>` form instead, and additionally decide which blade a given interface should land on (there is no TMOS-side concept to derive this from -- it's a target-platform capacity-planning decision).

### Port groups: an rSeries-only prerequisite with no TMOS equivalent

Before an rSeries interface can be used at all, its **port group** must be configured for the desired mode (speed/bundling) -- there is nothing analogous on TMOS. Key constraints (see the same guide, [Port Groups](https://clouddocs.f5.com/training/community/rseries-training/html/initial_setup_of_rseries_network_layer.html#network-settings-port-groups) section):

- High-speed port pairs (`1.0`/`2.0`, and `11.0`/`12.0` on r10000/r12000) must be configured for the **same** mode/speed as each other (both 40Gb or both 100Gb; some F5OS-A versions also support a 4x10Gb breakout).
- Low-speed ports (`3.0`-`10.0`, and `13.0`-`20.0` on r10000/r12000) can each be configured independently as 10Gb or 25Gb.
- Changing a port group's mode requires an appliance **reboot** (it reloads a different FPGA bitstream).

Plan port group modes based on the negotiated speed (`media_active`) `bigip_net_interfaces` reports for each i-Series interface you're mapping over, before attempting to configure the equivalent rSeries interface.

## Trunk -> LAG mapping is not 1:1

TMOS and F5OS represent link aggregation with **inverted ownership**:

- **TMOS**: the trunk object itself owns the member list. `bigip_net_trunks` (or `bigip_net_trunk` for a single already-known trunk name) reports `interfaces = ["1.1", "1.2"]` on the trunk (`net/trunk`); the physical interfaces themselves have no reference back to the trunk.
- **F5OS**: a LAG is a special `interface` object (`type: ieee8023adLag`, arbitrarily named, e.g. `"Arista"`), and each *member* physical interface instead carries a reference to it (`ethernet config aggregate-id <lag-name>`). There is no single object that lists "all members of this LAG" the way a TMOS trunk does -- you derive membership by scanning every interface for a matching `aggregate-id`.

Other differences to account for when mapping `bigip_net_trunks`/`bigip_net_trunk` output to an F5OS LAG:

| TMOS (`bigip_net_trunks` field) | F5OS equivalent | Notes |
|---|---|---|
| `name` | LAG interface name (e.g. `"Arista"`) | F5OS LAG names are free-form identifiers, not numeric IDs -- there's no requirement (or ability) to preserve a TMOS trunk's `id` field. |
| `lacp` (enabled/disabled) | `aggregation config lag-type` (`LACP` vs static) | A TMOS trunk with `lacp = "disabled"` (a static/non-LACP trunk) maps to a non-LACP F5OS LAG type -- verify F5OS-A version support for static (non-LACP) LAGs before assuming this is available, since LACP is the primary documented path. |
| `lacp_mode` (active/passive) | Not directly configurable the same way | F5OS's documented LAG config primarily exposes `lag-type LACP` and distribution hash; active/passive negotiation role is not a directly analogous exposed field in the reference examples -- verify against your target F5OS-A version's schema. |
| `distribution_hash` | `aggregation config distribution-hash` | Same concept, e.g. `src-dst-ipport`; value strings should be compatible but verify against your F5OS-A version. |
| `interfaces` (member list) | Set `aggregate-id` on each corresponding member interface | Inverted, as above -- this is a per-interface change, not a per-LAG one. |
| *(none -- VLANs are separate TMOS objects tagged to the trunk)* | `native_vlan` / `trunk_vlans` (`f5os_lag`); `aggregation switched-vlan config trunk-vlans [ ... ]` at the API level | F5OS attaches VLANs directly to the LAG interface's config (one untagged `native_vlan`, plus tagged `trunk_vlans`); TMOS instead has independent VLAN objects that reference the trunk by name. There's no 1:1 field in `bigip_net_trunks` for this -- use `scripts/extract-sys-settings.sh`'s derived `interface_vlans` map (see the [extraction guide](extract-sys-settings.html#interface-vlan-mapping)), which inverts `bigip_net_vlans`' VLAN-centric output into a per-trunk `native_vlan`/`trunk_vlans` pair ready to feed `f5os_lag` directly (same shape applies to `f5os_interface` for non-trunked physical interfaces). |
| `bandwidth`, `working_member_count`, `stp`, `type` | Largely informational/derived on F5OS (`state` block, `oper-status`, counters) | These are runtime state on F5OS, not something you configure directly to match a TMOS value. |

Source: [F5 rSeries Planning Guide -- Network Settings -> LAGs](https://clouddocs.f5.com/training/community/rseries-training/html/initial_setup_of_rseries_network_layer.html#network-settings-lags) (see the `show running-config interfaces` and API examples showing `ieee8023adLag` / `aggregate-id`).

## Interface status/speed field mapping

| `bigip_net_interfaces` field | F5OS equivalent (`openconfig-interfaces:interfaces` state) | Notes |
|---|---|---|
| `status` (`up`/`down`/`uninit`) | `oper-status` (`UP`/`DOWN`) | Same concept; TMOS additionally has an `uninit` state (interface never came up / no stats yet) with no obvious single F5OS equivalent -- treat as `DOWN` for mapping purposes unless you have reason to distinguish. |
| `media_active` (e.g. `10000T-FD`, `100TX-FD`, `none`) | `openconfig-if-ethernet:ethernet.state.port-speed` (e.g. `SPEED_100GB`) | TMOS reports negotiated media+duplex as one string; F5OS reports speed alone via a distinct enum-style value, and duplex is not separately exposed the same way. Don't parse/split the TMOS string and expect a clean field-for-field match -- treat this as "what speed was this interface actually running at" only. |
| `enabled` | `config.enabled` (interface admin state) | Same concept. |
| `mac_address` | `state.hw-mac-address` | Same concept, informational only (F5OS assigns its own MACs; don't attempt to preserve TMOS MACs). |
| `mtu` | `state.mtu` | Same concept, but verify against the F5OS default (commonly `9600` in the reference examples above) before assuming TMOS's configured MTU should be copied over as-is. |

## Summary

1. Extract i-Series interfaces/trunks/VLANs with `bigip_net_interfaces`/`bigip_net_trunks`/`bigip_net_vlans` (or run `scripts/extract-sys-settings.sh`, which extracts all three and derives the `interface_vlans` map described below in one pass). Use the singular `bigip_net_trunk`/`bigip_net_vlan` data sources instead when you only need to look up one already-known trunk or VLAN by name.
2. Decide your migration target (rSeries vs VELOS) and translate names using the tables above -- **do not** assume a single naming scheme applies to both.
3. For rSeries, plan port group modes from each interface's `media_active` speed *before* interface/LAG configuration; changing a port group later requires a reboot.
4. For trunks, remember F5OS LAG membership is set on each *member interface* (`aggregate-id`), not on the LAG object itself -- there is no single-object copy from a TMOS trunk's `interfaces` list.
5. For VLAN tagging, use the [extraction guide's `interface_vlans` map](extract-sys-settings.html#interface-vlan-mapping) (derived from `bigip_net_vlans`) to get a per-interface/trunk `native_vlan`/`trunk_vlans` pair ready for `f5os_interface`/`f5os_lag`, after translating names per the tables above.
6. Treat everything in this guide as a starting point verified against F5's rSeries Planning Guide at the time of writing; confirm exact field names/values against your specific F5OS-A/VELOS software version before relying on them for an automated migration.

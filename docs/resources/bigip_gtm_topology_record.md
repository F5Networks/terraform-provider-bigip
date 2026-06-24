# bigip_gtm_topology_record

Manages F5 BIG-IP GTM (Global Traffic Manager) Topology Record resources.

A GTM topology record defines a routing rule for topology-based load balancing. Each record matches a source (LDNS) to a destination (server, datacenter, or pool) with an associated score. When the GTM pool load balancing mode is set to `topology`, these records determine how DNS queries are routed based on the geographic or network location of the requesting LDNS.

## Example Usage

### Route a Region to a Datacenter

```hcl
resource "bigip_gtm_topology_record" "east_to_dc1" {
  ldns {
    match_type  = "region"
    match_value = "/Common/east-coast"
  }

  server {
    match_type  = "datacenter"
    match_value = "/Common/dc1"
  }

  order = 1
  score = 100
}
```

### Route a Subnet to a Pool

```hcl
resource "bigip_gtm_topology_record" "internal_to_pool" {
  ldns {
    match_type  = "subnet"
    match_value = "10.0.0.0/8"
  }

  server {
    match_type  = "pool"
    match_value = "/Common/internal-pool"
  }

  order = 2
  score = 50
}
```

### Route a Country to a Datacenter

```hcl
resource "bigip_gtm_topology_record" "us_to_us_dc" {
  ldns {
    match_type  = "country"
    match_value = "US"
  }

  server {
    match_type  = "datacenter"
    match_value = "/Common/us-datacenter"
  }

  order = 3
  score = 75
}
```

### Using Negation

```hcl
resource "bigip_gtm_topology_record" "not_east_to_dc2" {
  ldns {
    match_type   = "region"
    match_value  = "/Common/east-coast"
    match_negate = true
  }

  server {
    match_type  = "datacenter"
    match_value = "/Common/dc2"
  }

  score = 50
}
```

### Using with Topology Regions

```hcl
resource "bigip_gtm_topology_region" "east_coast" {
  name      = "east-coast"
  partition = "Common"

  members {
    name = "state US/New-York"
  }

  members {
    name = "state US/Virginia"
  }
}

resource "bigip_gtm_topology_record" "east_to_dc1" {
  ldns {
    match_type  = "region"
    match_value = "/Common/east-coast"
  }

  server {
    match_type  = "datacenter"
    match_value = "/Common/dc1"
  }

  order = 1
  score = 100

  depends_on = [bigip_gtm_topology_region.east_coast]
}
```

## Argument Reference

* `ldns` - (Required) The LDNS (source) match criteria block. Cannot be changed after creation. Contains the following:
  - `match_type` - (Required) The type of match. Valid values:
    - `region` - e.g., match_value = `/Common/east-coast`
    - `datacenter` - e.g., match_value = `/Common/dc1`
    - `pool` - e.g., match_value = `/Common/my-pool`
    - `subnet` - e.g., match_value = `10.0.0.0/8`
    - `country` - e.g., match_value = `US`
    - `state` - e.g., match_value = `US/California`
    - `continent` - e.g., match_value = `NA`
    - `isp` - e.g., match_value = `Comcast`
  - `match_value` - (Required) The value to match against.
  - `match_negate` - (Optional) If `true`, the match is negated. Default is `false`.

* `server` - (Required) The server (destination) match criteria block. Cannot be changed after creation. Same nested arguments as `ldns`.

* `description` - (Optional) User defined description.

* `order` - (Optional) The order in which the topology record is evaluated. Lower values are evaluated first. Default is `0`.

* `score` - (Optional) The weight or preference given to this topology record. Higher scores indicate stronger preference. Default is `1`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The topology record identifier in the format `ldns: <type> <value> server: <type> <value>`.

## Import

GTM topology records can be imported using the topology definition string, e.g.

```
terraform import bigip_gtm_topology_record.example "ldns: region /Common/east-coast server: datacenter /Common/dc1"
```

## Notes

* Topology records are evaluated in order. Use the `order` field to control evaluation priority.
* The `ldns` and `server` blocks cannot be changed after creation. To modify the source/destination match, destroy and recreate the record.
* Regions referenced in topology records must exist before the record is created. Use `depends_on` to enforce ordering.
* Topology records are only effective when a WideIP's `pool_lb_mode` is set to `topology`.

# bigip_gtm_topology_record

Manages F5 BIG-IP GTM (Global Traffic Manager) Topology Record resources.

A GTM topology record defines a routing rule for topology-based load balancing. Each record matches a source (LDNS) to a destination (server, datacenter, or pool) with an associated score. When the GTM pool load balancing mode is set to `topology`, these records determine how DNS queries are routed based on the geographic or network location of the requesting LDNS.

## Example Usage

### Route a Region to a Datacenter

```hcl
resource "bigip_gtm_topology_record" "east_to_dc1" {
  description = "ldns: region /Common/east-coast server: datacenter /Common/dc1"
  order       = 1
  score       = 100
}
```

### Route a Subnet to a Pool

```hcl
resource "bigip_gtm_topology_record" "internal_to_pool" {
  description = "ldns: subnet 10.0.0.0/8 server: pool /Common/internal-pool"
  order       = 2
  score       = 50
}
```

### Route a Country to a Datacenter

```hcl
resource "bigip_gtm_topology_record" "us_to_us_dc" {
  description = "ldns: country US server: datacenter /Common/us-datacenter"
  order       = 3
  score       = 75
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
  description = "ldns: region /Common/east-coast server: datacenter /Common/dc1"
  order       = 1
  score       = 100

  depends_on = [bigip_gtm_topology_region.east_coast]
}
```

## Argument Reference

* `description` - (Required) The topology record description that defines the source and destination match. This follows the BIG-IP format: `ldns: <source> server: <destination>`. Cannot be changed after creation. Valid source/destination types include:
  - `region <path>` - e.g., `region /Common/east-coast`
  - `datacenter <path>` - e.g., `datacenter /Common/dc1`
  - `pool <path>` - e.g., `pool /Common/my-pool`
  - `subnet <cidr>` - e.g., `subnet 10.0.0.0/8`
  - `country <code>` - e.g., `country US`
  - `state <country>/<state>` - e.g., `state US/California`
  - `continent <code>` - e.g., `continent NA`
  - `isp <name>` - e.g., `isp Comcast`

* `order` - (Optional) The order in which the topology record is evaluated. Lower values are evaluated first. Default is `0`.

* `score` - (Optional) The weight or preference given to this topology record. Higher scores indicate stronger preference. Default is `1`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The description string of the topology record.

## Import

GTM topology records can be imported using the description string, e.g.

```
terraform import bigip_gtm_topology_record.example "ldns: region /Common/east-coast server: datacenter /Common/dc1"
```

## Notes

* Topology records are evaluated in order. Use the `order` field to control evaluation priority.
* The `description` field cannot be changed after creation. To modify the source/destination match, destroy and recreate the record.
* Regions referenced in topology records must exist before the record is created. Use `depends_on` to enforce ordering.
* Topology records are only effective when a WideIP's `pool_lb_mode` is set to `topology`.
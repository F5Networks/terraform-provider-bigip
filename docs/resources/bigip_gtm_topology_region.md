# bigip_gtm_topology_region

Manages F5 BIG-IP GTM (Global Traffic Manager) Topology Region resources.

A GTM topology region is a named group of network locations used in topology-based load balancing. Regions can contain subnets, countries, states, continents, datacenters, ISPs, or references to other regions. They are referenced by topology records to define source-to-destination routing preferences.

## Example Usage

### Region with Subnets

```hcl
resource "bigip_gtm_topology_region" "internal" {
  name      = "internal-networks"
  partition = "Common"

  members {
    name = "subnet 10.0.0.0/8"
  }

  members {
    name = "subnet 172.16.0.0/12"
  }

  members {
    name = "subnet 192.168.0.0/16"
  }
}
```

### Region with Geographic Locations

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

  members {
    name = "state US/Massachusetts"
  }
}
```

### Region with Mixed Member Types

```hcl
resource "bigip_gtm_topology_region" "primary_dc" {
  name      = "primary-region"
  partition = "Common"

  members {
    name = "datacenter /Common/dc1"
  }

  members {
    name = "subnet 10.10.0.0/16"
  }

  members {
    name = "country US"
  }
}
```

## Argument Reference

* `name` - (Required) Name of the GTM topology region. Cannot be changed after creation.

* `partition` - (Optional) Partition in which to create the region. Default is `Common`. Cannot be changed after creation.

* `members` - (Optional) The members that define this topology region. Each member is a block containing:
  * `name` - (Required) The member definition string. Valid formats include:
    - `subnet <cidr>` - e.g., `subnet 10.0.0.0/8`
    - `country <code>` - e.g., `country US`
    - `state <country>/<state>` - e.g., `state US/California`
    - `continent <code>` - e.g., `continent NA`
    - `datacenter <path>` - e.g., `datacenter /Common/my-dc`
    - `isp <name>` - e.g., `isp Comcast`
    - `region <path>` - e.g., `region /Common/other-region`
    - `pool <path>` - e.g., `pool /Common/my-pool`

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The full path of the region (e.g., `/Common/east-coast`)

## Import

GTM topology regions can be imported using the full path, e.g.

```
terraform import bigip_gtm_topology_region.example /Common/east-coast
```

## Notes

* Regions must be created before they can be referenced in topology records.
* A region can reference other regions as members, enabling hierarchical grouping.
* The `name` and `partition` cannot be changed after creation. You must destroy and recreate the resource to change these values.
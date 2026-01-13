# bigip_gtm_datacenter

Manages F5 BIG-IP GTM (Global Traffic Manager) Datacenter resources.

A GTM datacenter is a logical representation of a physical location where servers reside. Datacenters are fundamental building blocks in GTM topology and are used for geographic load balancing and disaster recovery.

## Example Usage

### Basic Datacenter

```hcl
resource "bigip_gtm_datacenter" "example" {
  name      = "example_datacenter"
  partition = "Common"

  location    = "Seattle, WA"
  contact     = "admin@example.com"
  description = "Primary datacenter for west coast operations"

  enabled = true

  prober_preference = "inside-datacenter"
  prober_fallback   = "any-available"
}
```

### Datacenter with Custom Prober Settings

```hcl
resource "bigip_gtm_datacenter" "custom_prober" {
  name      = "custom_datacenter"
  partition = "Common"

  location    = "New York, NY"
  contact     = "ops@example.com"
  description = "East coast datacenter with custom prober configuration"

  enabled = true

  prober_preference = "pool"
  prober_fallback   = "outside-datacenter"
}
```

### Disabled Datacenter

```hcl
resource "bigip_gtm_datacenter" "maintenance" {
  name      = "maintenance_dc"
  partition = "Common"

  location    = "Chicago, IL"
  contact     = "maintenance@example.com"
  description = "Datacenter currently under maintenance"

  enabled = false

  prober_preference = "inside-datacenter"
  prober_fallback   = "any-available"
}
```

## Argument Reference

* `name` - (Required) Name of the GTM datacenter. Cannot be changed after creation.

* `partition` - (Optional) Partition in which to create the datacenter. Default is `Common`. Cannot be changed after creation.

* `contact` - (Optional) Contact information for the datacenter administrator.

* `description` - (Optional) Description of the datacenter.

* `enabled` - (Optional) Enable or disable the datacenter. Default is `true`. When set to `false`, the datacenter is disabled and will not be used for load balancing decisions.

* `location` - (Optional) Physical location of the datacenter. This is a free-form text field.

* `prober_preference` - (Optional) Specifies the type of prober to prefer when monitoring resources in this datacenter. Default is `inside-datacenter`. Valid options are:
  - `inside-datacenter` - Prefer probers inside this datacenter
  - `outside-datacenter` - Prefer probers outside this datacenter
  - `pool` - Use a specific pool of probers
  - `inherit` - Inherit from parent configuration

* `prober_fallback` - (Optional) Specifies the type of prober to use as fallback when the preferred prober is unavailable. Default is `any-available`. Valid options are:
  - `any-available` - Use any available prober
  - `inside-datacenter` - Fallback to probers inside this datacenter
  - `outside-datacenter` - Fallback to probers outside this datacenter
  - `pool` - Fallback to a specific pool of probers
  - `inherit` - Inherit from parent configuration

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The full path of the datacenter (e.g., `/Common/example_datacenter`)

## Import

GTM datacenters can be imported using the full path, e.g.

```
terraform import bigip_gtm_datacenter.example /Common/example_datacenter
```

## Notes

* Datacenters are required before you can create GTM servers and pools.
* The `name` and `partition` cannot be changed after creation. You must destroy and recreate the resource to change these values.
* Disabling a datacenter (`enabled = false`) removes it from load balancing decisions but preserves its configuration.
* Prober settings control how BIG-IP monitors the health of resources in this datacenter.

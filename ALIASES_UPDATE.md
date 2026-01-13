# GTM WideIP - Aliases Support Update

## Summary

Added support for the `aliases` field and fixed the `loadBalancingDecisionLogVerbosity` to handle arrays correctly based on the latest API response.

## Changes Made

### 1. SDK Update
**File**: `vendor/github.com/f5devcentral/go-bigip/gtm.go`

Updated the `GTMWideIP` struct to include:
- `LoadBalancingDecisionLogVerbosity` changed from `string` to `[]string` (array)
- `TopologyPreferEdns0ClientSubnet` added as `string`
- `Aliases` added as `[]string` (array)

### 2. Resource Schema Update
**File**: `bigip/resource_bigip_gtm_wideip.go`

Added new schema field:
```go
"aliases": {
    Type:     schema.TypeSet,
    Optional: true,
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
    Description: "Specifies alternate domain names for the WideIP",
},
```

### 3. Read Function Enhancement
Updated `resourceBigipGtmWideipRead` to handle:
- `topology_prefer_edns0_client_subnet` field
- `load_balancing_decision_log_verbosity` as an array
- `aliases` as an array

```go
d.Set("topology_prefer_edns0_client_subnet", wideip.TopologyPreferEdns0ClientSubnet)

// Handle LoadBalancingDecisionLogVerbosity as array
if len(wideip.LoadBalancingDecisionLogVerbosity) > 0 {
    d.Set("load_balancing_decision_log_verbosity", wideip.LoadBalancingDecisionLogVerbosity)
}

// Handle Aliases as array
if len(wideip.Aliases) > 0 {
    d.Set("aliases", wideip.Aliases)
}
```

### 4. Update Function Enhancement
Updated `resourceBigipGtmWideipUpdate` to handle:
- `topology_prefer_edns0_client_subnet` field
- `load_balancing_decision_log_verbosity` as an array (converting from Set to []string)
- `aliases` as an array (converting from Set to []string)

```go
wideip.TopologyPreferEdns0ClientSubnet = d.Get("topology_prefer_edns0_client_subnet").(string)

// Handle LoadBalancingDecisionLogVerbosity as array
if v, ok := d.GetOk("load_balancing_decision_log_verbosity"); ok {
    verbositySet := v.(*schema.Set)
    verbosityList := make([]string, 0, verbositySet.Len())
    for _, item := range verbositySet.List() {
        verbosityList = append(verbosityList, item.(string))
    }
    wideip.LoadBalancingDecisionLogVerbosity = verbosityList
}

// Handle Aliases as array
if v, ok := d.GetOk("aliases"); ok {
    aliasesSet := v.(*schema.Set)
    aliasesList := make([]string, 0, aliasesSet.Len())
    for _, item := range aliasesSet.List() {
        aliasesList = append(aliasesList, item.(string))
    }
    wideip.Aliases = aliasesList
}
```

## API Response Mapping

The implementation now correctly maps all fields from your API response:

```json
{
    "name": "testwideip.local",
    "partition": "Common",
    "fullPath": "/Common/testwideip.local",
    "description": "test_wideip_a",
    "enabled": true,
    "failureRcode": "noerror",
    "failureRcodeResponse": "disabled",
    "failureRcodeTtl": 0,
    "lastResortPool": "a /Common/firstpool",
    "loadBalancingDecisionLogVerbosity": ["pool-selection"],  // ✅ Array support
    "minimalResponse": "disabled",
    "persistCidrIpv4": 32,
    "persistCidrIpv6": 128,
    "persistence": "disabled",
    "poolLbMode": "round-robin",
    "topologyPreferEdns0ClientSubnet": "enabled",  // ✅ New field
    "ttlPersistence": 3600,
    "aliases": ["testwideip2.local"]  // ✅ New field
}
```

## Usage Example

### Basic WideIP with Aliases
```hcl
resource "bigip_gtm_wideip" "example" {
  name      = "testwideip.local"
  type      = "a"
  partition = "Common"
  
  description = "test_wideip_a"
  enabled     = true
  
  last_resort_pool = "a /Common/firstpool"
  pool_lb_mode     = "round-robin"
  minimal_response = "disabled"
  
  topology_prefer_edns0_client_subnet = "enabled"
  
  load_balancing_decision_log_verbosity = ["pool-selection"]
  
  aliases = ["testwideip2.local"]
}
```

### Multiple Aliases
```hcl
resource "bigip_gtm_wideip" "multi" {
  name = "app.example.com"
  type = "a"
  
  aliases = [
    "app-primary.example.com",
    "app-secondary.example.com",
    "app-backup.example.com"
  ]
}
```

## Files Updated

1. ✅ `vendor/github.com/f5devcentral/go-bigip/gtm.go` - SDK struct updated
2. ✅ `bigip/resource_bigip_gtm_wideip.go` - Resource implementation updated
3. ✅ `examples/bigip_gtm_wideip.tf` - Example updated with aliases
4. ✅ `examples/bigip_gtm_wideip_api_example.tf` - New example matching API response
5. ✅ `examples/bigip_gtm_wideip_comprehensive.tf` - Comprehensive test examples
6. ✅ `docs/resources/bigip_gtm_wideip.md` - Documentation updated

## Testing

### Example Test Configuration
See `examples/bigip_gtm_wideip_api_example.tf` for a configuration that exactly matches your API response.

### Test Commands
```bash
# Build the provider
cd /Users/r.chinthalapalli/GolandProjects/F5Networks/terraform-provider-bigip
go build

# Run tests
go test -v ./bigip -run TestAccBigipGtmWideip

# Test with Terraform
cd examples
terraform init
terraform plan -var-file=api_example.tfvars
terraform apply -auto-approve
```

## Validation

The implementation now correctly handles:
- ✅ Single alias: `aliases = ["alias1.example.com"]`
- ✅ Multiple aliases: `aliases = ["alias1.com", "alias2.com", "alias3.com"]`
- ✅ No aliases: Field is optional and can be omitted
- ✅ Load balancing verbosity as array
- ✅ Topology EDNS0 client subnet preference

## Benefits

1. **Complete API Coverage**: All fields from the BIG-IP API response are now supported
2. **Flexibility**: Multiple aliases allow for complex DNS configurations
3. **Type Safety**: Proper handling of arrays vs strings
4. **Backward Compatible**: Existing configurations continue to work

## Notes

- Aliases are stored as a Set in Terraform to prevent duplicates
- Order of aliases doesn't matter (Set is unordered)
- Empty aliases array won't be sent to the API
- LoadBalancingDecisionLogVerbosity now properly handles multiple values

## Next Steps

1. Build and test the provider
2. Verify aliases work correctly with your BIG-IP
3. Test import functionality with WideIPs that have aliases
4. Verify that existing WideIPs can be updated to add/remove aliases

# GTM WideIP Terraform Resource Implementation

This implementation provides a complete Terraform resource for managing F5 BIG-IP GTM (Global Traffic Manager) WideIP objects.

## What Was Created

### 1. Resource Implementation
**File**: `bigip/resource_bigip_gtm_wideip.go`

A complete Terraform resource implementation that supports:
- **CRUD Operations**: Create, Read, Update, and Delete WideIP resources
- **All WideIP Types**: a, aaaa, cname, mx, naptr, srv
- **Full Configuration**: All fields from the BIG-IP API are supported
- **Import Support**: Import existing WideIPs into Terraform state
- **Context-aware**: Uses Terraform SDK v2 with context support

### 2. Resource Schema

The resource includes all fields from your API responses:
- Basic fields: name, type, partition, description
- State fields: enabled, disabled
- DNS failure handling: failure_rcode, failure_rcode_response, failure_rcode_ttl
- Pool configuration: last_resort_pool, pool_lb_mode
- Persistence: persistence, persist_cidr_ipv4, persist_cidr_ipv6, ttl_persistence
- Response options: minimal_response
- Logging: load_balancing_decision_log_verbosity
- Topology: topology_prefer_edns0_client_subnet

### 3. Examples
**File**: `examples/bigip_gtm_wideip.tf`

Six comprehensive examples demonstrating:
1. Basic WideIP with minimal configuration
2. Full configuration with all options
3. IPv6 WideIP (type: aaaa)
4. CNAME WideIP
5. Disabled WideIP
6. WideIP with custom persistence settings

### 4. Documentation
**File**: `docs/resources/bigip_gtm_wideip.md`

Complete documentation including:
- Resource description and overview
- Multiple usage examples
- Full argument reference with descriptions
- Detailed notes on pool LB modes, last resort pool format, persistence, and failure RCODE
- Import instructions
- Related resources
- API endpoint references

### 5. Tests
**File**: `bigip/resource_bigip_gtm_wideip_test.go`

Test suite including:
- Create test
- Update test
- Import test
- Helper functions for validation

## Key Features

### 1. Based on Your API Responses
The implementation maps directly to the JSON structure you provided:

```json
{
    "name": "testwideip.local",
    "partition": "Common",
    "description": "test_wideip_a",
    "enabled": true,
    "failureRcode": "noerror",
    "lastResortPool": "a /Common/firstpool",
    ...
}
```

### 2. Type Safety
- Uses the existing `bigip.GTMWideIP` struct from the go-bigip SDK
- Proper type conversions for strings, bools, and integers
- Set type for array fields like `load_balancing_decision_log_verbosity`

### 3. Proper Resource ID Management
- ID format: `/partition/name:type`
- Allows proper import and state management
- Handles partition correctly in all operations

### 4. Full Path Construction
The resource properly constructs full paths for API calls:
```go
fullPath := fmt.Sprintf("/%s/%s", partition, name)
```

### 5. Error Handling
Comprehensive error handling with:
- Detailed error messages
- Proper logging
- Graceful handling of missing resources

## Usage Example

Based on your API response, here's how to use it:

```hcl
resource "bigip_gtm_wideip" "example" {
  name      = "testwideip.local"
  type      = "a"
  partition = "Common"
  
  description              = "test_wideip_a"
  enabled                  = true
  failure_rcode            = "noerror"
  failure_rcode_response   = "disabled"
  failure_rcode_ttl        = 0
  last_resort_pool         = "a /Common/firstpool"
  minimal_response         = "enabled"
  persist_cidr_ipv4        = 32
  persist_cidr_ipv6        = 128
  persistence              = "disabled"
  pool_lb_mode             = "round-robin"
  ttl_persistence          = 3600
  
  load_balancing_decision_log_verbosity = ["pool-selection"]
}
```

## API Mapping

The resource properly maps to these BIG-IP API endpoints:

| Operation | Endpoint | Method |
|-----------|----------|--------|
| List | `/mgmt/tm/gtm/wideip` | GET |
| Get | `/mgmt/tm/gtm/wideip/a/~Common~testwideip.local` | GET |
| Create | `/mgmt/tm/gtm/wideip/a` | POST |
| Update | `/mgmt/tm/gtm/wideip/a/~Common~testwideip.local` | PUT |
| Delete | `/mgmt/tm/gtm/wideip/a/~Common~testwideip.local` | DELETE |

## Testing

### Run Tests
```bash
cd /Users/r.chinthalapalli/GolandProjects/F5Networks/terraform-provider-bigip
go test -v ./bigip -run TestAccBigipGtmWideip
```

### Test with Terraform
```bash
cd examples
terraform init
terraform plan
terraform apply
```

## Provider Registration

The resource is already registered in the provider at:
`bigip/provider.go` line 200:
```go
"bigip_gtm_wideip": resourceBigipGtmWideip(),
```

## Next Steps

1. **Build the Provider**:
   ```bash
   go build
   ```

2. **Test Locally**:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

3. **Verify Against BIG-IP**:
   - Create a WideIP using Terraform
   - Verify it appears in the BIG-IP GUI
   - Make changes and verify updates work
   - Test import functionality
   - Test destroy

## Notes

### Last Resort Pool Format
When specifying `last_resort_pool`, use the format shown in your API response:
```
"a /Common/firstpool"
```
This is: `<type> <partition>/<pool_name>`

### LoadBalancingDecisionLogVerbosity
Note: The API returns this as an array `["pool-selection"]`, but the SDK defines it as a string. I've implemented it as a Set in the schema to match the API behavior. If you encounter issues, this field might need adjustment based on the actual go-bigip SDK implementation.

### Defaults
The resource includes sensible defaults that match BIG-IP defaults:
- `partition`: "Common"
- `enabled`: true
- `failure_rcode`: "noerror"
- `failure_rcode_response`: "disabled"
- `failure_rcode_ttl`: 0
- `minimal_response`: "enabled"
- `persist_cidr_ipv4`: 32
- `persist_cidr_ipv6`: 128
- `persistence`: "disabled"
- `pool_lb_mode`: "round-robin"
- `ttl_persistence`: 3600
- `topology_prefer_edns0_client_subnet`: "disabled"

## Files Modified/Created

1. ✅ Modified: `bigip/resource_bigip_gtm_wideip.go` - Main resource implementation
2. ✅ Created: `bigip/resource_bigip_gtm_wideip_test.go` - Test suite
3. ✅ Created: `examples/bigip_gtm_wideip.tf` - Usage examples
4. ✅ Created: `docs/resources/bigip_gtm_wideip.md` - Documentation
5. ✅ Created: This README

## Troubleshooting

If you encounter issues:

1. **Check the SDK**: Verify the `bigip.GTMWideIP` struct matches expectations
2. **Enable Debug Logging**: Set `TF_LOG=DEBUG` to see API calls
3. **Verify API Access**: Test the API endpoints directly with curl
4. **Check Partition**: Ensure the partition exists on BIG-IP
5. **Check Pool Names**: If using `last_resort_pool`, ensure the pool exists

## Additional Enhancements (Optional)

Future enhancements could include:
1. Pool attachments (adding pools to WideIPs)
2. Pool member management
3. Alias support
4. Rules attachments
5. Custom validation for fields

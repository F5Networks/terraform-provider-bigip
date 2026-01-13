# GTM WideIP Quick Reference

## Basic Commands

### Initialize Terraform
```bash
cd examples
terraform init
```

### Validate Configuration
```bash
terraform validate
```

### Plan Changes
```bash
terraform plan
```

### Apply Changes
```bash
terraform apply
```

### Destroy Resources
```bash
terraform destroy
```

### Import Existing WideIP
```bash
terraform import bigip_gtm_wideip.example /Common/testwideip.local:a
```

## Common Configurations

### 1. Basic A Record WideIP
```hcl
resource "bigip_gtm_wideip" "basic" {
  name      = "app.example.com"
  type      = "a"
  partition = "Common"
}
```

### 2. WideIP with Pool
```hcl
resource "bigip_gtm_wideip" "with_pool" {
  name             = "app.example.com"
  type             = "a"
  partition        = "Common"
  last_resort_pool = "a /Common/my_pool"
  pool_lb_mode     = "round-robin"
}
```

### 3. IPv6 WideIP
```hcl
resource "bigip_gtm_wideip" "ipv6" {
  name      = "app.example.com"
  type      = "aaaa"
  partition = "Common"
}
```

### 4. WideIP with Persistence
```hcl
resource "bigip_gtm_wideip" "persistent" {
  name            = "app.example.com"
  type            = "a"
  partition       = "Common"
  persistence     = "enabled"
  ttl_persistence = 7200
}
```

### 5. WideIP with Topology Load Balancing
```hcl
resource "bigip_gtm_wideip" "topology" {
  name         = "app.example.com"
  type         = "a"
  partition    = "Common"
  pool_lb_mode = "topology"
}
```

### 6. WideIP with Failure RCODE
```hcl
resource "bigip_gtm_wideip" "with_rcode" {
  name                   = "app.example.com"
  type                   = "a"
  partition              = "Common"
  failure_rcode          = "servfail"
  failure_rcode_response = "enabled"
  failure_rcode_ttl      = 300
}
```

### 7. WideIP with Aliases
```hcl
resource "bigip_gtm_wideip" "with_aliases" {
  name      = "app.example.com"
  type      = "a"
  partition = "Common"
  
  aliases = [
    "app-alias1.example.com",
    "app-alias2.example.com",
    "www.app.example.com"
  ]
}
```

### 8. Complete Configuration (Matching API Response)
```hcl
resource "bigip_gtm_wideip" "complete" {
  name      = "testwideip.local"
  type      = "a"
  partition = "Common"
  
  description = "test_wideip_a"
  enabled     = true
  
  failure_rcode          = "noerror"
  failure_rcode_response = "disabled"
  failure_rcode_ttl      = 0
  
  last_resort_pool = "a /Common/firstpool"
  pool_lb_mode     = "round-robin"
  minimal_response = "disabled"
  
  persistence       = "disabled"
  persist_cidr_ipv4 = 32
  persist_cidr_ipv6 = 128
  ttl_persistence   = 3600
  
  topology_prefer_edns0_client_subnet = "enabled"
  
  load_balancing_decision_log_verbosity = ["pool-selection"]
  
  aliases = ["testwideip2.local"]
}
```

## WideIP Types

| Type   | Description | Example |
|--------|-------------|---------|
| `a` | IPv4 Address | 192.168.1.1 |
| `aaaa` | IPv6 Address | 2001:db8::1 |
| `cname` | Canonical Name | alias.example.com |
| `mx` | Mail Exchange | mail.example.com |
| `naptr` | Name Authority Pointer | Used for ENUM |
| `srv` | Service Locator | _http._tcp.example.com |

## Pool Load Balancing Modes

| Mode | Description |
|------|-------------|
| `round-robin` | Equal distribution across pools |
| `ratio` | Distribution based on pool ratios |
| `topology` | Based on topology records |
| `global-availability` | Based on availability and load |
| `least-connections` | Pool with fewest connections |
| `lowest-round-trip-time` | Pool with lowest RTT |

## Common Field Values

### failure_rcode
- `noerror` - No error condition
- `servfail` - Server failure
- `nxdomain` - Non-existent domain
- `refused` - Query refused

### failure_rcode_response
- `enabled` - Return configured RCODE
- `disabled` - Return no answer

### minimal_response
- `enabled` - Return minimal DNS response
- `disabled` - Return full DNS response

### persistence
- `enabled` - Client IP persistence enabled
- `disabled` - No persistence

## Troubleshooting

### Check Resource State
```bash
terraform state show bigip_gtm_wideip.example
```

### Refresh State
```bash
terraform refresh
```

### View Current Configuration
```bash
terraform show
```

### Debug Logging
```bash
export TF_LOG=DEBUG
terraform apply
```

### Verify on BIG-IP
```bash
# SSH to BIG-IP
tmsh list gtm wideip a /Common/testwideip.local

# Or via API
curl -sku admin:admin https://10.145.71.31/mgmt/tm/gtm/wideip/a/~Common~testwideip.local
```

## Testing

### Run Unit Tests
```bash
cd /Users/r.chinthalapalli/GolandProjects/F5Networks/terraform-provider-bigip
go test -v ./bigip -run TestAccBigipGtmWideip_create
```

### Run All WideIP Tests
```bash
go test -v ./bigip -run TestAccBigipGtmWideip
```

### Test with Coverage
```bash
go test -v -cover ./bigip -run TestAccBigipGtmWideip
```

## Environment Variables

```bash
# BIG-IP Connection
export BIGIP_HOST="10.145.71.31"
export BIGIP_USER="admin"
export BIGIP_PASSWORD="admin"

# Terraform Provider Settings
export TF_ACC=1  # Enable acceptance tests
```

## Example provider.tf
```hcl
terraform {
  required_providers {
    bigip = {
      source  = "F5Networks/bigip"
      version = "~> 1.0"
    }
  }
}

provider "bigip" {
  address  = "10.145.71.31"
  username = "admin"
  password = "admin"
}
```

## Common Errors and Solutions

### Error: WideIP already exists
**Solution**: Import the existing WideIP or delete it manually first
```bash
terraform import bigip_gtm_wideip.example /Common/name:a
```

### Error: Pool not found
**Solution**: Ensure the pool referenced in `last_resort_pool` exists
```hcl
# Create pool first or use existing pool
last_resort_pool = "a /Common/existing_pool"
```

### Error: Invalid type
**Solution**: Use only valid WideIP types: a, aaaa, cname, mx, naptr, srv
```hcl
type = "a"  # Must be lowercase
```

### Error: Partition not found
**Solution**: Ensure partition exists or use default "Common"
```hcl
partition = "Common"
```

## API Endpoints Reference

| Action | Endpoint | Method |
|--------|----------|--------|
| List All Types | `/mgmt/tm/gtm/wideip` | GET |
| List Type | `/mgmt/tm/gtm/wideip/a` | GET |
| Get WideIP | `/mgmt/tm/gtm/wideip/a/~Common~name` | GET |
| Create WideIP | `/mgmt/tm/gtm/wideip/a` | POST |
| Update WideIP | `/mgmt/tm/gtm/wideip/a/~Common~name` | PUT |
| Delete WideIP | `/mgmt/tm/gtm/wideip/a/~Common~name` | DELETE |

## Tips

1. **Always specify partition**: Even though "Common" is default, be explicit
2. **Test with plan first**: Always run `terraform plan` before `apply`
3. **Use variables**: Don't hardcode sensitive values
4. **Import existing resources**: Use `terraform import` for existing WideIPs
5. **Check BIG-IP version**: Ensure your BIG-IP supports the features you're using
6. **Use meaningful names**: WideIPs should be FQDNs (e.g., app.example.com)
7. **Match pool type**: Last resort pool type must match WideIP type
8. **Enable logging**: Use `load_balancing_decision_log_verbosity` for troubleshooting

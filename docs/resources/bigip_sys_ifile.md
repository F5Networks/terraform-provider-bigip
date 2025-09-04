---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_ifile"
subcategory: "System"
description: |-
  Provides details about bigip_sys_ifile resource
---

# bigip_sys_ifile

`bigip_sys_ifile` This resource uploads and manages system iFiles on F5 BIG-IP devices. 
System iFiles store file content on the BIG-IP that can be referenced by iRules, LTM policies, and other BIG-IP configurations for traffic processing and decision making.

## Example Usage

### Basic System iFile

```hcl
resource "bigip_sys_ifile" "config_file" {
  name      = "app-config"
  partition = "Common"
  content   = file("${path.module}/config.json")
}
```

### System iFile with Sub-path

```hcl
resource "bigip_sys_ifile" "template_file" {
  name      = "error-template"
  partition = "Common"
  sub_path  = "templates"
  content   = <<-EOT
    <html>
      <head><title>Service Unavailable</title></head>
      <body>
        <h1>503 - Service Temporarily Unavailable</h1>
        <p>Please try again later.</p>
      </body>
    </html>
  EOT
}
```

### Dynamic Content from Template

```hcl
locals {
  server_config = {
    database_host = var.database_host
    api_key       = var.api_key
    environment   = var.environment
  }
}

resource "bigip_sys_ifile" "app_config" {
  name      = "application-config"
  partition = "Production"
  content   = templatefile("${path.module}/templates/app-config.tpl", local.server_config)
}
```

### JSON Configuration File

```hcl
locals {
  server_list = jsonencode({
    servers = [
      { name = "web1", ip = "10.1.1.10", port = 80 },
      { name = "web2", ip = "10.1.1.11", port = 80 },
      { name = "web3", ip = "10.1.1.12", port = 80 }
    ]
  })
}

resource "bigip_sys_ifile" "server_config" {
  name      = "server-list"
  partition = "MyApp"
  content   = local.server_list
}
```

### Using System iFile with LTM iFile

```hcl
# Create system iFile with content
resource "bigip_sys_ifile" "lookup_table" {
  name      = "url-rewrite-map"
  partition = "Common"
  content   = <<-EOT
    /old-api/v1/ /api/v2/
    /legacy/ /new/
    /deprecated/ /current/
  EOT
}

# Create LTM iFile that references the system iFile
resource "bigip_ltm_ifile" "ltm_lookup" {
  name      = "ltm-url-rewrite-map"
  partition = "Common"
  file_name = bigip_sys_ifile.lookup_table.id
}

# Use in an iRule
resource "bigip_ltm_irule" "url_rewriter" {
  name = "url-rewrite-rule"
  irule = <<-EOT
    when HTTP_REQUEST {
      set uri [HTTP::uri]
      set mapping [ifile get ltm-url-rewrite-map]
      foreach line [split $mapping "\n"] {
        set parts [split $line " "]
        if {[string match [lindex $parts 0]* $uri]} {
          HTTP::uri [string map [list [lindex $parts 0] [lindex $parts 1]] $uri]
          break
        }
      }
    }
  EOT
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) Name of the system iFile to be created on BIG-IP. Changing this forces a new resource to be created.

* `content` - (Required, string) The content of the iFile. This can be inline text, file content loaded with `file()`, or dynamically generated content. This field is marked as sensitive.

* `partition` - (Optional, string) Partition where the iFile will be stored. Defaults to `Common`.

* `sub_path` - (Optional, string) Subdirectory within the partition for organizing iFiles hierarchically.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The full path identifier of the system iFile (e.g., `/Common/my-ifile` or `/Common/subpath/my-ifile`).

* `checksum` - MD5 checksum of the iFile content, automatically calculated by BIG-IP.

* `size` - Size of the iFile content in bytes.

## Import

System iFiles can be imported using their full path:

```bash
terraform import bigip_sys_ifile.example /Common/my-ifile
```

For iFiles with sub-paths:

```bash
terraform import bigip_sys_ifile.example /Common/templates/my-ifile
```

## Notes

* The `content` field is marked as sensitive and will not be displayed in Terraform logs or state output.
* Changes to `name` will force recreation of the resource since iFile names cannot be changed after creation.
* The `checksum` and `size` attributes are automatically computed by the BIG-IP system.
* iFile content is uploaded to the BIG-IP system and stored there permanently until the resource is destroyed.
* Use `file()` function to load content from local files or `templatefile()` for dynamic content generation.
* System iFiles can be referenced by `bigip_ltm_ifile` resources for use in LTM configurations.

## Path Structure

The full path of an iFile follows this pattern:
- Without sub-path: `/{partition}/{name}`
- With sub-path: `/{partition}/{sub_path}/{name}`

Examples:
- `/Common/config-file`
- `/Production/templates/error-page`
- `/MyApp/configs/database-settings`

## Related Resources

* [`bigip_ltm_ifile`](bigip_ltm_ifile.html) - Creates LTM iFiles that reference system iFiles
* [`bigip_ltm_irule`](bigip_ltm_irule.html) - Creates iRules that can access iFile content
* [`bigip_ltm_policy`](bigip_ltm_policy.html) - Creates LTM policies that can use iFile content

## Security Considerations

* iFile content is stored on the BIG-IP system and may contain sensitive information
* Use appropriate BIG-IP access controls to limit who can view or modify iFiles
* Consider using Terraform's sensitive variable handling for confidential content
* The `content` field is marked as sensitive in Terraform state to prevent accidental exposure
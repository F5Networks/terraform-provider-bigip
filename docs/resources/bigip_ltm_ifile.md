---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_ifile"
subcategory: "Local Traffic Manager (LTM)"
description: |-
  Provides details about bigip_ltm_ifile resource
---

# bigip_ltm_ifile

`bigip_ltm_ifile` This resource creates an LTM iFile on F5 BIG-IP that references an existing system iFile. 
LTM iFiles are used in iRules and LTM policies to access file content for traffic processing and decision making.

## Example Usage

### Basic LTM iFile

```hcl
# First create a system iFile
resource "bigip_sys_ifile" "config_file" {
  name      = "app-config"
  partition = "Common"
  content   = file("${path.module}/config.json")
}

# Create LTM iFile that references the system iFile
resource "bigip_ltm_ifile" "ltm_config" {
  name      = "ltm-app-config"
  partition = "Common"
  file_name = bigip_sys_ifile.config_file.id
}
```

### LTM iFile with Sub-path

```hcl
resource "bigip_sys_ifile" "template_file" {
  name      = "error-template"
  partition = "Common"
  sub_path  = "templates"
  content   = file("${path.module}/error.html")
}

resource "bigip_ltm_ifile" "ltm_template" {
  name      = "ltm-error-template"
  partition = "Common"
  sub_path  = "templates"
  file_name = bigip_sys_ifile.template_file.id
}
```

### Using LTM iFile in iRule

```hcl
resource "bigip_sys_ifile" "server_list" {
  name      = "server-mapping"
  partition = "Production"
  content   = <<-EOT
    web1:10.1.1.10
    web2:10.1.1.11
    web3:10.1.1.12
  EOT
}

resource "bigip_ltm_ifile" "ltm_servers" {
  name      = "ltm-server-mapping"
  partition = "Production"
  file_name = "/Production/server-mapping"
}

resource "bigip_ltm_irule" "server_selector" {
  name  = "select-server-rule"
  irule = <<-EOT
    when HTTP_REQUEST {
      set server_map [ifile get ltm-server-mapping]
      # Process server mapping logic
      foreach line [split $server_map "\n"] {
        set parts [split $line ":"]
        # Implement server selection logic
      }
    }
  EOT
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) Name of the LTM iFile to be created on BIG-IP.

* `file_name` - (Required, string) The system iFile name to reference (e.g., `/Common/my-sys-ifile`). This should reference an existing system iFile created with `bigip_sys_ifile`.

* `partition` - (Optional, string) Partition where the LTM iFile will be created. Defaults to `Common`.

* `sub_path` - (Optional, string) Subdirectory within the partition for organizing iFiles.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The full path identifier of the LTM iFile (e.g., `/Common/my-ltm-ifile`).

* `full_path` - The complete path of the LTM iFile on the BIG-IP system.

## Import

LTM iFiles can be imported using their full path:

```bash
terraform import bigip_ltm_ifile.example /Common/my-ltm-ifile
```

For iFiles with sub-paths:

```bash
terraform import bigip_ltm_ifile.example /Common/templates/my-ltm-ifile
```

## Notes

* The referenced system iFile (specified in `file_name`) must exist before creating the LTM iFile.
* LTM iFiles are primarily used in iRules and LTM policies for traffic processing.
* Changes to `name`, `partition`, or `sub_path` will force recreation of the resource.
* The LTM iFile acts as a reference to the system iFile and doesn't store content directly.
* Use `bigip_sys_ifile` to upload file content, then reference it with `bigip_ltm_ifile` for LTM usage.

## Related Resources

* [`bigip_sys_ifile`](bigip_sys_ifile.html) - Creates system iFiles with content
* [`bigip_ltm_irule`](bigip_ltm_irule.html) - Creates iRules that can reference LTM iFiles
* [`bigip_ltm_policy`](bigip_ltm_policy.html) - Creates LTM policies that can use LTM iFiles

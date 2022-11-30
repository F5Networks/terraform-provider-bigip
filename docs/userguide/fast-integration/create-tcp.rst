.. _fast-integration-tcp:

Scenario #2: Creating a TCP application
---------------------------------------

The goal of this template is to deploy a new TCP application on BIG-IP using Terraform as the orchestrator.

Pre-requisites
on the BIG-IP:

 version 16.1 minimal
 credentials with REST API access
on Terraform:

 use of F5 bigip provider version 1.16.0 minimal
 use of Hashicorp version following Link
Create TCP application
Create 5 files:

main.tf
variables.tf
inputs.tfvars
outputs.tf
providers.tf
variables.tf

variable bigip {}
variable username {}
variable password {}
inputs.tfvars

bigip = "10.1.1.9:443"
username = "admin"
password = "whatIsYourBigIPPassword?"
providers.tf

terraform {
  required_providers {
    bigip = {
      source = "F5Networks/bigip"
      version = ">= 1.16.0"
    }
  }
}
provider "bigip" {
  address  = var.bigip
  username = var.username
  password = var.password
}
main.tf

resource "bigip_fast_tcp_app" "this" {
  application               = "myApp2"
  tenant                    = "scenario2"
  virtual_server            {
    ip                        = "10.1.10.222"
    port                      = 80
  }
  pool_members  {
    addresses                 = ["10.1.10.120", "10.1.10.121", "10.1.10.122"]
    port                      = 80
  }
  load_balancing_mode       = "least-connections-member"
  existing_monitor          = "/Common/http"
}
outputs.tf

output "configJSON" {
	value	= bigip_fast_tcp_app.this
}
Now, run the following commands, so we can:

Initialize the terraform project
Plan the changes
Apply the changes
$ terraform init -upgrade

Initializing the backend...

Initializing provider plugins...
- Finding f5networks/bigip versions matching ">= 1.16.0"...
- Using previously-installed f5networks/bigip v1.16.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.


$ terraform plan -var-file=inputs.tfvars -out scenario2

Terraform used the selected providers to generate the following execution plan.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # bigip_fast_tcp_app.this will be created
  + resource "bigip_fast_tcp_app" "this" {
      + application         = "myApp2"
      + existing_monitor    = "/Common/http"
      + fast_tcp_json       = (known after apply)
      + id                  = (known after apply)
      + load_balancing_mode = "least-connections-member"
      + tenant              = "myTenant2"

      + pool_members {
          + addresses = [
              + "10.1.10.120",
              + "10.1.10.121",
              + "10.1.10.122",
            ]
          + port      = 80
        }

      + virtual_server {
          + ip   = "10.1.10.222"
          + port = 80
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + configJSON = {
      + application         = "myApp2"
      + existing_monitor    = "/Common/http"
      + existing_pool       = null
      + existing_snat_pool  = null
      + fast_tcp_json       = (known after apply)
      + id                  = (known after apply)
      + load_balancing_mode = "least-connections-member"
      + monitor             = []
      + pool_members        = [
          + {
              + addresses        = [
                  + "10.1.10.120",
                  + "10.1.10.121",
                  + "10.1.10.122",
                ]
              + connection_limit = null
              + port             = 80
              + priority_group   = null
              + share_nodes      = null
            },
        ]
      + slow_ramp_time      = null
      + snat_pool_address   = null
      + tenant              = "myTenant2"
      + virtual_server      = [
          + {
              + ip   = "10.1.10.222"
              + port = 80
            },
        ]
    }

───────────────────────────────────────────────────────────────────────────────

Saved the plan to: scenario2

To perform exactly these actions, run the following command to apply:
    terraform apply "scenario2"


$ terraform apply "scenario2"
bigip_fast_tcp_app.this: Creating...
bigip_fast_tcp_app.this: Still creating... [10s elapsed]
bigip_fast_tcp_app.this: Still creating... [20s elapsed]
bigip_fast_tcp_app.this: Creation complete after 27s [id=myApp2]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

configJSON = {
  "application" = "myApp2"
  "existing_monitor" = "/Common/http"
  "existing_pool" = tostring(null)
  "existing_snat_pool" = tostring(null)
  "fast_tcp_json" = tostring(null)
  "id" = "myApp2"
  "load_balancing_mode" = "least-connections-member"
  "monitor" = tolist([])
  "pool_members" = toset([
    {
      "addresses" = tolist([
        "10.1.10.120",
        "10.1.10.121",
        "10.1.10.122",
      ])
      "connection_limit" = tonumber(null)
      "port" = 80
      "priority_group" = tonumber(null)
      "share_nodes" = tobool(null)
    },
  ])
  "slow_ramp_time" = 0
  "snat_pool_address" = tolist(null) /* of string */
  "tenant" = "myTenant2"
  "virtual_server" = tolist([
    {
      "ip" = "10.1.10.222"
      "port" = 80
    },
  ])
}


$ terraform output -json
{
  "configJSON": {
    "sensitive": false,
    "type": [
      "object",
      {
        "application": "string",
        "existing_monitor": "string",
        "existing_pool": "string",
        "existing_snat_pool": "string",
        "fast_tcp_json": "string",
        "id": "string",
        "load_balancing_mode": "string",
        "monitor": [
          "list",
          [
            "object",
            {
              "interval": "number"
            }
          ]
        ],
        "pool_members": [
          "set",
          [
            "object",
            {
              "addresses": [
                "list",
                "string"
              ],
              "connection_limit": "number",
              "port": "number",
              "priority_group": "number",
              "share_nodes": "bool"
            }
          ]
        ],
        "slow_ramp_time": "number",
        "snat_pool_address": [
          "list",
          "string"
        ],
        "tenant": "string",
        "virtual_server": [
          "list",
          [
            "object",
            {
              "ip": "string",
              "port": "number"
            }
          ]
        ]
      }
    ],
    "value": {
      "application": "myApp2",
      "existing_monitor": "/Common/http",
      "existing_pool": null,
      "existing_snat_pool": null,
      "fast_tcp_json": null,
      "id": "myApp2",
      "load_balancing_mode": "least-connections-member",
      "monitor": [],
      "pool_members": [
        {
          "addresses": [
            "10.1.10.120",
            "10.1.10.121",
            "10.1.10.122"
          ],
          "connection_limit": 0,
          "port": 80,
          "priority_group": 0,
          "share_nodes": false
        }
      ],
      "slow_ramp_time": 0,
      "snat_pool_address": null,
      "tenant": "myTenant2",
      "virtual_server": [
        {
          "ip": "10.1.10.222",
          "port": 80
        }
      ]
    }
  }
}

$ terraform output -json > config_export1.json
The Terraform CLI output is designed to be parsed by humans. To get machine-readable format for automation, use the -json flag for JSON-formatted output.

Checking the virtual server and pool status you discover both down. Now update your terraform main.tf file with the following: main.tf

resource "bigip_fast_tcp_app" "this" {
  application               = "myApp2"
  tenant                    = "scenario2"
  virtual_server            {
    ip                        = "10.1.10.222"
    port                      = 80
  }
  pool_members  {
    addresses                 = ["10.1.10.120", "10.1.10.121", "10.1.10.122"]
    port                      = 80
  }
  load_balancing_mode       = "least-connections-member"
  existing_monitor          = "/Common/tcp"
}
Now, run the following commands, so we can:

Plan the changes
Apply the changes
$ terraform plan -var-file=inputs.tfvars -out scenario2
bigip_fast_tcp_app.this: Refreshing state... [id=myApp2]

Terraform used the selected providers to generate the following execution plan.
Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # bigip_fast_tcp_app.this will be updated in-place
  ~ resource "bigip_fast_tcp_app" "this" {
      ~ existing_monitor    = "/Common/http" -> "/Common/tcp"
        id                  = "myApp2"
        # (4 unchanged attributes hidden)


        # (2 unchanged blocks hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ configJSON = {
      ~ existing_monitor    = "/Common/http" -> "/Common/tcp"
        id                  = "myApp2"
        # (11 unchanged elements hidden)
    }

───────────────────────────────────────────────────────────────────────────────

Saved the plan to: scenario2

To perform exactly these actions, run the following command to apply:
    terraform apply "scenario2"

$ terraform apply "scenario2"
bigip_fast_tcp_app.this: Modifying... [id=myApp2]
bigip_fast_tcp_app.this: Still modifying... [id=myApp2, 10s elapsed]
bigip_fast_tcp_app.this: Modifications complete after 16s [id=myApp2]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

configJSON = {
  "application" = "myApp2"
  "existing_monitor" = "/Common/tcp"
  "existing_pool" = tostring(null)
  "existing_snat_pool" = tostring(null)
  "fast_tcp_json" = tostring(null)
  "id" = "myApp2"
  "load_balancing_mode" = "least-connections-member"
  "monitor" = tolist([])
  "pool_members" = toset([
    {
      "addresses" = tolist([
        "10.1.10.120",
        "10.1.10.121",
        "10.1.10.122",
      ])
      "connection_limit" = 0
      "port" = 80
      "priority_group" = 0
      "share_nodes" = false
    },
  ])
  "slow_ramp_time" = 0
  "snat_pool_address" = tolist(null) /* of string */
  "tenant" = "myTenant2"
  "virtual_server" = tolist([
    {
      "ip" = "10.1.10.222"
      "port" = 80
    },
  ])
}

$ terraform output -json > config_export2.json

$ diff config_export1.json config_export2.json 
59c59
<       "existing_monitor": "/Common/http",
---
>       "existing_monitor": "/Common/tcp",
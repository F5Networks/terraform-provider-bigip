.. _fast-integration-http:

Scenario #3: Creating a HTTP application
========================================

The goal of this template is to deploy a new HTTP application on BIG-IP using Terraform as the orchestrator.

Pre-requisites
--------------
On the BIG-IP:

- F5 BIG-IP version 16.1 or newer
- Credentials with REST API access

On Terraform:

- Using F5 BIG-IP provider version 1.16.0 or newer
- Using Hashicorp versions following :ref:`versions`

Create HTTP application
-----------------------
Create 5 files in folder ``cd ~/terraform/scenario3/app1``:

- main.tf
- variables.tf
- inputs.tfvars
- outputs.tf
- providers.tf

.. code-block:: json
   :caption: variables.tf
   :linenos:

   variable bigip {}
   variable username {}
   variable password {}

|

.. code-block:: json
   :caption: inputs.tfvars
   :linenos:

   bigip = "10.1.1.9:443"
   username = "admin"
   password = "whatIsYourBigIPPassword?"

|

.. code-block:: json
   :caption: providers.tf
   :linenos:

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

|

.. code-block:: json
   :caption: main.tf
   :linenos:

   resource "bigip_fast_http_app" "app1" {
     application               = "myApp3"
     tenant                    = "scenario3"
     virtual_server            {
       ip                        = "10.1.10.223"
       port                      = 80
     }
     pool_members  {
       addresses                 = ["10.1.10.120", "10.1.10.121", "10.1.10.122"]
       port                      = 80
     }
     load_balancing_mode       = "least-connections-member"
   }

|

.. code-block:: json
   :caption: outputs.tf
   :linenos:

   output "configJSON" {
           value		= bigip_fast_http_app.app1
           sensitive	= true
   }

|

Run the following commands, so you can:

1. Initialize the Terraform project
2. Plan the changes
3. Apply the changes

::

    $ cd ~/terraform/scenario3/app1

    $ terraform init -upgrade
    
    Initializing the backend...
    
    Initializing provider plugins...
    - Finding f5networks/bigip versions matching ">= 1.15.0"...
    - Installing f5networks/bigip v1.16.0...
    - Installed f5networks/bigip v1.16.0 (signed by a HashiCorp partner, key ID EBD2EE9544728437)

    Partner and community providers are signed by their developers.
    If you'd like to know more about provider signing, you can read about it here:
    https://www.terraform.io/docs/cli/plugins/signing.html
    
    Terraform has created a lock file .terraform.lock.hcl to record the provider
    selections it made above. Include this file in your version control repository
    so that Terraform can guarantee to make the same selections by default when
    you run "terraform init" in the future.
    
    Terraform has been successfully initialized!
    
    You may now begin working with Terraform. Try running "terraform plan" to see
    any changes that are required for your infrastructure. All Terraform commands
    should now work.
    
    If you ever set or change modules or backend configuration for Terraform,
    rerun this command to reinitialize your working directory. If you forget, other
    commands will detect it and remind you to do so if necessary.


    $ terraform plan -var-file=inputs.tfvars -out scenario3app1
    
    Terraform used the selected providers to generate the following execution plan.
    Resource actions are indicated with the following symbols:
      + create
    
    Terraform will perform the following actions:
    
      # bigip_fast_http_app.app1 will be created
      + resource "bigip_fast_http_app" "app1" {
          + application         = "myApp3"
          + existing_monitor    = "/Common/http"
          + fast_http_json      = (known after apply)
          + id                  = (known after apply)
          + load_balancing_mode = "least-connections-member"
          + tenant              = "scenario3"
    
          + pool_members {
              + addresses = [
                  + "10.1.10.120",
                  + "10.1.10.121",
                  + "10.1.10.122",
                ]
              + port      = 80
            }
    
          + virtual_server {
              + ip   = "10.1.10.223"
              + port = 80
            }
        }
    
    Plan: 1 to add, 0 to change, 0 to destroy.
    
    Changes to Outputs:
      + configJSON = (sensitive value)
    
    ───────────────────────────────────────────────────────────────────────────────
    
    Saved the plan to: scenario3app1
    
    To perform exactly these actions, run the following command to apply:
        terraform apply "scenario3app1"
    
    
    $ terraform apply "scenario3app1"
    bigip_fast_http_app.app1: Creating...
    bigip_fast_http_app.app1: Still creating... [10s elapsed]
    bigip_fast_http_app.app1: Creation complete after 19s [id=myApp3]

    Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

    Outputs:

    configJSON = <sensitive>


    $ terraform output -json > config_export1.json

|

Now you want to add a custom HTTP monitor and a snat pool. Update your Terraform main.tf file with the following:

.. code-block:: json
   :caption: main.tf
   :linenos:

   resource "bigip_fast_http_app" "app1" {
     application               = "myApp3"
     tenant                    = "scenario3"
     virtual_server            {
       ip                        = "10.1.10.223"
       port                      = 80
     }
     pool_members  {
       addresses                 = ["10.1.10.120", "10.1.10.121", "10.1.10.122"]
       port                      = 80
     }
     snat_pool_address = ["10.1.10.50", "10.1.10.51", "10.1.10.52"]
     load_balancing_mode       = "least-connections-member"
     monitor       {
       send_string               = "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConnection: Close\\r\\n\\r\\n"
       response                  = "200 OK"
     }
   }

|

Run the following commands so you can:

1. Plan the changes
2. Apply the changes

::

    $ terraform plan -var-file=inputs.tfvars -out scenario3app1
    bigip_fast_http_app.app1: Refreshing state... [id=myApp3]

    Note: Objects have changed outside of Terraform

    Terraform detected the following changes made outside of Terraform since the
    last "terraform apply" which may have affected this plan:

      # bigip_fast_http_app.app1 has changed
      ~ resource "bigip_fast_http_app" "app1" {
            id                    = "myApp3"
          + security_log_profiles = []
            # (5 unchanged attributes hidden)
    
          + pool_members {
              + addresses        = [
                  + "10.1.10.120",
                  + "10.1.10.121",
                  + "10.1.10.122",
                ]
              + connection_limit = 0
              + port             = 80
              + priority_group   = 0
              + share_nodes      = false
            }
          - pool_members {
              - addresses = [
                  - "10.1.10.120",
                  - "10.1.10.121",
                  - "10.1.10.122",
                ] -> null
              - port      = 80 -> null
            }
    
            # (1 unchanged block hidden)
        }
    
    
    Unless you have made equivalent changes to your configuration, or ignored the
    relevant attributes using ignore_changes, the following plan may include
    actions to undo or respond to these changes.
    
    ───────────────────────────────────────────────────────────────────────────────
    
    Terraform used the selected providers to generate the following execution plan.
    Resource actions are indicated with the following symbols:
      ~ update in-place
    
    Terraform will perform the following actions:
    
      # bigip_fast_http_app.app1 will be updated in-place
      ~ resource "bigip_fast_http_app" "app1" {
            id                    = "myApp3"
          + snat_pool_address     = [
              + "10.1.10.50",
              + "10.1.10.51",
              + "10.1.10.52",
            ]
            # (6 unchanged attributes hidden)
    
          + monitor {
              + monitor_auth = false
              + response     = "302"
              + send_string  = "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConnection: Close\\r\\n\\r\\n"
            }
    
            # (2 unchanged blocks hidden)
        }
    
    Plan: 0 to add, 1 to change, 0 to destroy.
    
    Changes to Outputs:
      ~ configJSON = (sensitive value)
    
    ───────────────────────────────────────────────────────────────────────────────
    
    Saved the plan to: scenario3app1

    To perform exactly these actions, run the following command to apply:
        terraform apply "scenario3app1"
    
    
    $ terraform apply "scenario3app1"
    bigip_fast_http_app.app1: Modifying... [id=myApp3]
    bigip_fast_http_app.app1: Still modifying... [id=myApp3, 10s elapsed]
    bigip_fast_http_app.app1: Still modifying... [id=myApp3, 20s elapsed]
    bigip_fast_http_app.app1: Modifications complete after 23s [id=myApp3]
    
    Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
    
    Outputs:
    
    configJSON = <sensitive>
    
    $ terraform output -json > config_export2.json
    
    $ diff config_export1.json config_export2.json
    68c68,77
    <       "monitor": [],
    ---
    >       "monitor": [
    >         {
    >           "interval": null,
    >           "monitor_auth": false,
    >           "password": null,
    >           "response": "302",
    >           "send_string": "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConnection: Close\\r\\n\\r\\n",
    >           "username": null
    >         }
    >       ],
    76c85
    <           "connection_limit": null,
    ---
    >           "connection_limit": 0,
    78,79c87,88
    <           "priority_group": null,
    <           "share_nodes": null
    ---
    >           "priority_group": 0,
    >           "share_nodes": false
    83c92,96
    <       "snat_pool_address": null,
    ---
    >       "snat_pool_address": [
    >         "10.1.10.50",
    >         "10.1.10.51",
    >         "10.1.10.52"
    >       ],

|

Now you want to add a second virtual server or application in the same tenant. Create a second main.tf file in the app2 folder with the following:

Create 5 files in folder ``cd ~/terraform/scenario3/app2``:

- main.tf
- variables.tf
- inputs.tfvars
- outputs.tf
- providers.tf

.. code-block:: json
   :caption: variables.tf
   :linenos:

   variable bigip {}
   variable username {}
   variable password {}

|

.. code-block:: json
   :caption: inputs.tfvars
   :linenos:

   bigip = "10.1.1.9:443"
   username = "admin"
   password = "whatIsYourBigIPPassword?"

|

.. code-block:: json
   :caption: providers.tf
   :linenos:

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

|

.. code-block:: json
   :caption: main.tf
   :linenos:

   resource "bigip_fast_http_app" "app2" {
     application               = "myApp3-1"
     tenant                    = "scenario3"
     virtual_server            {
       ip                        = "10.1.10.233"
       port                      = 80
     }
     pool_members  {
       addresses                 = ["10.1.10.130", "10.1.10.131", "10.1.10.132"]
       port                      = 80
     }
     snat_pool_address           = ["10.1.10.53", "10.1.10.54", "10.1.10.55"]
     load_balancing_mode         = "round-robin"
     monitor       {
       send_string               = "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConnection: Close\\r\\n\\r\\n"
       response                  = "302"
     }
   }

|

.. code-block:: json
   :caption: outputs.tf
   :linenos:

   output "configJSON2" {
   	value		= bigip_fast_http_app.app2
   	sensitive	= true
   }
      
|

Run the following commands so you can:

1. Plan the changes
2. Apply the changes

::

    $ cd ~/terraform/scenario3/app2

    $ terraform init -upgrade
    
    Initializing the backend...
    
    Initializing provider plugins...
    - Finding f5networks/bigip versions matching ">= 1.16.0"...
    - Installing f5networks/bigip v1.16.0...
    - Installed f5networks/bigip v1.16.0 (signed by a HashiCorp partner, key ID EBD2EE9544728437)
    
    Partner and community providers are signed by their developers.
    If you'd like to know more about provider signing, you can read about it here:
    https://www.terraform.io/docs/cli/plugins/signing.html
    
    Terraform has created a lock file .terraform.lock.hcl to record the provider
    selections it made above. Include this file in your version control repository
    so that Terraform can guarantee to make the same selections by default when
    you run "terraform init" in the future.
    
    Terraform has been successfully initialized!
    
    You may now begin working with Terraform. Try running "terraform plan" to see
    any changes that are required for your infrastructure. All Terraform commands
    should now work.
    
    If you ever set or change modules or backend configuration for Terraform,
    rerun this command to reinitialize your working directory. If you forget, other
    commands will detect it and remind you to do so if necessary.


    $ terraform plan -var-file=inputs.tfvars -out scenario3app2

    Terraform used the selected providers to generate the following execution plan.
    Resource actions are indicated with the following symbols:
      + create
    
    Terraform will perform the following actions:
    
      # bigip_fast_http_app.app2 will be created
      + resource "bigip_fast_http_app" "app2" {
          + application         = "myApp3-1"
          + existing_monitor    = "/Common/http"
          + fast_http_json      = (known after apply)
          + id                  = (known after apply)
          + load_balancing_mode = "round-robin"
          + snat_pool_address   = [
              + "10.1.10.53",
              + "10.1.10.54",
              + "10.1.10.55",
            ]
          + tenant              = "scenario3"
    
          + monitor {
              + monitor_auth = false
              + response     = "302"
              + send_string  = "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConnection: Close\\r\\n\\r\\n"
            }
    
          + pool_members {
              + addresses = [
                  + "10.1.10.130",
                  + "10.1.10.131",
                  + "10.1.10.132",
                ]
              + port      = 80
            }
    
          + virtual_server {
              + ip   = "10.1.10.233"
              + port = 80
            }
        }
    
    Plan: 1 to add, 0 to change, 0 to destroy.
    
    Changes to Outputs:
      + configJSON2 = (sensitive value)
    
    ───────────────────────────────────────────────────────────────────────────────

    Saved the plan to: scenario3app2
    
    To perform exactly these actions, run the following command to apply:
        terraform apply "scenario3app2"
    
    $ terraform apply "scenario3app2"
    bigip_fast_http_app.app2: Creating...
    bigip_fast_http_app.app2: Still creating... [10s elapsed]
    bigip_fast_http_app.app2: Still creating... [20s elapsed]
    bigip_fast_http_app.app2: Creation complete after 23s [id=myApp3-1]
    
    Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
    
    Outputs:
    
    configJSON2 = <sensitive>


    $ terraform output -json
    {
      "configJSON2": {
        "sensitive": true,
        "type": [
          "object",
          {
            "application": "string",
            "endpoint_ltm_policy": [
              "list",
              "string"
            ],
            "existing_monitor": "string",
            "existing_pool": "string",
            "existing_snat_pool": "string",
            "existing_waf_security_policy": "string",
            "fast_http_json": "string",
            "id": "string",
            "load_balancing_mode": "string",
            "monitor": [
              "list",
              [
                "object",
                {
                  "interval": "number",
                  "monitor_auth": "bool",
                  "password": "string",
                  "response": "string",
                  "send_string": "string",
                  "username": "string"
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
            "security_log_profiles": [
              "list",
              "string"
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
            ],
            "waf_security_policy": [
              "list",
              [
                "object",
                {
                  "enable": "bool"
                }
              ]
            ]
          }
        ],
        "value": {
          "application": "myApp3-1",
          "endpoint_ltm_policy": null,
          "existing_monitor": "/Common/http",
          "existing_pool": "",
          "existing_snat_pool": "",
          "existing_waf_security_policy": null,
          "fast_http_json": "{\"app_name\":\"myApp3-1\",\"enable_asm_logging\":false,\"enable_monitor\":true,\"enable_pool\":true,\"enable_snat\":true,\"enable_tls_client\":false,\"enable_tls_server\":false,\"enable_waf_policy\":false,\"load_balancing_mode\":\"round-robin\",\"make_monitor\":true,\"make_pool\":true,\"make_snatpool\":true,\"make_tls_client_profile\":false,\"make_tls_server_profile\":false,\"make_waf_policy\":false,\"monitor_credentials\":false,\"monitor_expected_response\":\"302\",\"monitor_name_http\":\"/Common/http\",\"monitor_send_string\":\"GET / HTTP/1.1\\\\r\\\\nHost: example.com\\\\r\\\\nConnection: Close\\\\r\\\\n\\\\r\\\\n\",\"pool_members\":[{\"connectionLimit\":0,\"priorityGroup\":0,\"serverAddresses\":[\"10.1.10.130\",\"10.1.10.131\",\"10.1.10.132\"],\"servicePort\":80,\"shareNodes\":true}],\"snat_addresses\":[\"10.1.10.53\",\"10.1.10.54\",\"10.1.10.55\"],\"snat_automap\":false,\"tenant_name\":\"scenario3\",\"virtual_address\":\"10.1.10.233\",\"virtual_port\":80}",
          "id": "myApp3-1",
          "load_balancing_mode": "round-robin",
          "monitor": [
            {
              "interval": 0,
              "monitor_auth": false,
              "password": "",
              "response": "302",
              "send_string": "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConnection: Close\\r\\n\\r\\n",
              "username": ""
            }
          ],
          "pool_members": [
            {
              "addresses": [
                "10.1.10.130",
                "10.1.10.131",
                "10.1.10.132"
              ],
              "connection_limit": null,
              "port": 80,
              "priority_group": null,
              "share_nodes": null
            }
          ],
          "security_log_profiles": null,
          "slow_ramp_time": null,
          "snat_pool_address": [
            "10.1.10.53",
            "10.1.10.54",
            "10.1.10.55"
          ],
          "tenant": "scenario3",
          "virtual_server": [
            {
              "ip": "10.1.10.233",
              "port": 80
            }
          ],
          "waf_security_policy": []
        }
      }
    }


.. Note::

   You created two different application definitions sharing the same tenant in two different Terraform
   projects. The FAST plugin makes the AS3 declarations reconcile on the BIG-IP so you do not have to manage the stacking of
   them for a single tenant.

.. The |fast-serviceDiscovery| enables you to discover pool members using account AWS, GCP, and Azure tags.

.. |fast-serviceDiscovery| raw:: html

   <a href="https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_fast_http_app#service-discovery" target="_blank">Service Discovery parameters</a>

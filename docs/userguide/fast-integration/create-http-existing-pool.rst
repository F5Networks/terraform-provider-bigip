.. _fast-integration-http-existing-pool:

Scenario #5: Creating a HTTP application using existing pool and snat pool
==========================================================================

The goal of this template is to deploy a new HTTP application on BIG-IP using a pool and a snat pool that have already been created.

Pre-requisites
--------------
 On the BIG-IP:

- F5 BIG-IP version 16.1 or newer
- Credentials with REST API access

On Terraform:

- Using F5 BIG-IP provider version 1.16.0 or newer
- Using Hashicorp versions following :ref:`versions`


Create HTTPS application with existing pool and SNAT pool
---------------------------------------------------------
Create 5 files:

- main.tf
- variables.tf
- inputs.auto.tfvars
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

   resource "bigip_fast_https_app" "this" {
     application               = "myApp5"
     tenant                    = "scenario5"
     virtual_server            {
       ip                        = "10.1.10.225"
       port                      = 443
     }
     tls_server_profile {
       tls_cert_name             = "/Common/default.crt"
       tls_key_name              = "/Common/default.key"
     }
     existing_snat_pool        = "/Common/snat-pool-90"
     existing_pool             = "/Common/dvwa"
     load_balancing_mode       = "round-robin"
   }

|

.. code-block:: json
   :caption: outputs.tf
   :linenos:

   output "configJSON" {
   	value		= bigip_fast_https_app.this
   	sensitive	= true
   }

|

Run the following commands so you can:

1. Initialize the Terraform project
2. Plan the changes
3. Apply the changes

::

    $ terraform init -upgrade

    Initializing the backend...

    Initializing provider plugins...
    - Finding f5networks/bigip versions matching ">= 1.16.0"...
    - Installing f5networks/bigip v1.16.0...
    - Installed f5networks/bigip v1.16.0 (signed by a HashiCorp partner, key ID EBD2EE9544728437)
    
    Partner and community providers are signed by their developers.
    If you'd like to know more about provider signing, you can read about it here:
    https://www.terraform.io/docs/cli/plugins/signing.html

    Terraform has made some changes to the provider dependency selections recorded
    in the .terraform.lock.hcl file. Review those changes and commit them to your
    version control system if they represent changes you intended to make.

    Terraform has been successfully initialized!

    You may now begin working with Terraform. Try running "terraform plan" to see
    any changes that are required for your infrastructure. All Terraform commands
    should now work.

    If you ever set or change modules or backend configuration for Terraform,
    rerun this command to reinitialize your working directory. If you forget, other
    commands will detect it and remind you to do so if necessary.


    $ terraform plan -out scenario5
    
    Terraform used the selected providers to generate the following execution plan.
    Resource actions are indicated with the following symbols:
      + create
    
    Terraform will perform the following actions:
    
      # bigip_fast_https_app.this will be created
      + resource "bigip_fast_https_app" "this" {
          + application         = "myApp5"
          + existing_pool       = "/Common/dvwa"
          + existing_snat_pool  = "/Common/snat-pool-90"
          + id                  = (known after apply)
          + load_balancing_mode = "round-robin"
          + tenant              = "scenario5"
    
          + tls_server_profile {
              + tls_cert_name = "/Common/default.crt"
              + tls_key_name  = "/Common/default.key"
            }
    
          + virtual_server {
              + ip   = "10.1.10.225"
              + port = 443
            }
        }
    
    Plan: 1 to add, 0 to change, 0 to destroy.

    Changes to Outputs:
      + configJSON = (sensitive value)
    
    ───────────────────────────────────────────────────────────────────────────────
    
    Saved the plan to: scenario5
    
    To perform exactly these actions, run the following command to apply:
        terraform apply "scenario5"
    
    
    $ terraform apply "scenario5"
    bigip_fast_https_app.this: Creating...
    bigip_fast_https_app.this: Still creating... [10s elapsed]
    bigip_fast_https_app.this: Creation complete after 18s [id=myApp5]
    
    Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
    
    Outputs:
    
    configJSON = <sensitive>
    
    $ terraform output -json
    {
      "configJSON": {
        "sensitive": true,
        "type": [
          "object",
          {
            "application": "string",
            "existing_monitor": "string",
            "existing_pool": "string",
            "existing_snat_pool": "string",
            "existing_tls_client_profile": "string",
            "existing_tls_server_profile": "string",
            "existing_waf_security_policy": "string",
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
            "slow_ramp_time": "number",
            "snat_pool_address": [
              "list",
              "string"
            ],
            "tenant": "string",
            "tls_client_profile": [
              "list",
              [
                "object",
                {
                  "tls_cert_name": "string",
                  "tls_key_name": "string"
                }
              ]
            ],
            "tls_server_profile": [
              "list",
              [
                "object",
                {
                  "tls_cert_name": "string",
                  "tls_key_name": "string"
                }
              ]
            ],
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
          "application": "myApp5",
          "existing_monitor": "",
          "existing_pool": "/Common/dvwa",
          "existing_snat_pool": "/Common/snat-pool-90",
          "existing_tls_client_profile": null,
          "existing_tls_server_profile": null,
          "existing_waf_security_policy": null,
          "id": "myApp5",
          "load_balancing_mode": "round-robin",
          "monitor": [],
          "pool_members": [],
          "slow_ramp_time": 0,
          "snat_pool_address": null,
          "tenant": "scenario5",
          "tls_client_profile": [],
          "tls_server_profile": [
            {
              "tls_cert_name": "/Common/default.crt",
              "tls_key_name": "/Common/default.key"
            }
          ],
          "virtual_server": [
            {
              "ip": "10.1.10.225",
              "port": 443
            }
          ],
          "waf_security_policy": []
        }
      }
    }
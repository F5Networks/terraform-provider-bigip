.. _fast-integration-https:

Scenario #4: Creating a HTTPS application
=========================================
The goal of this template is to deploy a new HTTPS application on BIG-IP using Terraform as the orchestrator.

Pre-requisites
--------------

- F5 BIG-IP version 16.1 or newer
- Credentials with REST API access

On Terraform:

- Using F5 BIG-IP provider version 1.16.0 or newer
- Using Hashicorp versions following :ref:`versions`
- Certificate and key files should be available under ~/terraform/scenario4:

  - app4.crt
  - app4.key


Create HTTPS application
------------------------
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

   resource "bigip_ssl_certificate" "app4crt" {
     name      = "app4.crt"
     content   = file("app4.crt")
     partition = "Common"
   }
   
   resource "bigip_ssl_key" "app4key" {
     name      = "app4.key"
     content   = file("app4.key")
     partition = "Common"
   }
   
   resource "bigip_fast_https_app" "this" {
     application               = "myApp4"
     tenant                    = "scenario4"
     virtual_server            {
       ip                        = "10.1.10.224"
       port                      = 443
     }
     tls_server_profile {
       tls_cert_name             = "/Common/app4.crt"
       tls_key_name              = "/Common/app4.key"
     }
     pool_members  {
       addresses                 = ["10.1.10.120", "10.1.10.121", "10.1.10.122"]
       port                      = 80
     }
     snat_pool_address = ["10.1.10.50", "10.1.10.51", "10.1.10.52"]
     load_balancing_mode       = "least-connections-member"
     monitor       {
       send_string               = "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConn
   ection: Close\\r\\n\\r\\n"
       response                  = "200 OK"
     }
     depends_on		      = [bigip_ssl_certificate.app4crt, bigip_ssl_key.ap
   p4key]
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

Now, run the following commands, so you can:

1. Initialize the terraform project
2. Plan the changes
3. Apply the changes

::

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
    
    
    $ terraform plan -out scenario4
    
    Terraform used the selected providers to generate the following execution plan.
    Resource actions are indicated with the following symbols:
      + create
    
    Terraform will perform the following actions:
    
      # bigip_fast_https_app.this will be created
      + resource "bigip_fast_https_app" "this" {
          + application         = "myApp4"
          + fast_https_json     = (known after apply)
          + id                  = (known after apply)
          + load_balancing_mode = "least-connections-member"
          + snat_pool_address   = [
              + "10.1.10.50",
              + "10.1.10.51",
              + "10.1.10.52",
            ]
          + tenant              = "scenario4"
    
          + monitor {
              + monitor_auth = false
              + response     = "200 OK"
              + send_string  = "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConnection: Close\\r\\n\\r\\n"
            }
    
          + pool_members {
              + addresses = [
                  + "10.1.10.120",
                  + "10.1.10.121",
                  + "10.1.10.122",
                ]
              + port      = 80
            }
    
          + tls_server_profile {
              + tls_cert_name = "/Common/app4.crt"
              + tls_key_name  = "/Common/app4.key"
            }
    
          + virtual_server {
              + ip   = "10.1.10.224"
              + port = 443
            }
        }
    
      # bigip_ssl_certificate.app4crt will be created
      + resource "bigip_ssl_certificate" "app4crt" {
          + content   = (sensitive value)
          + full_path = (known after apply)
          + id        = (known after apply)
          + name      = "app4.crt"
          + partition = "Common"
        }
    
      # bigip_ssl_key.app4key will be created
      + resource "bigip_ssl_key" "app4key" {
          + content   = (sensitive value)
          + full_path = (known after apply)
          + id        = (known after apply)
          + name      = "app4.key"
          + partition = "Common"
        }
    
    Plan: 3 to add, 0 to change, 0 to destroy.
    
    Changes to Outputs:
      + configJSON = (sensitive value)
    
    ───────────────────────────────────────────────────────────────────────────────
    
    Saved the plan to: scenario4
    
    To perform exactly these actions, run the following command to apply:
        terraform apply "scenario4"
    
    
    $ terraform apply "scenario4"
    bigip_ssl_certificate.app4crt: Creating...
    bigip_ssl_key.app4key: Creating...
    bigip_ssl_key.app4key: Creation complete after 1s [id=app4.key]
    bigip_ssl_certificate.app4crt: Creation complete after 1s [id=app4.crt]
    bigip_fast_https_app.this: Creating...
    bigip_fast_https_app.this: Still creating... [10s elapsed]
    bigip_fast_https_app.this: Creation complete after 17s [id=myApp4]
    
    Apply complete! Resources: 3 added, 0 changed, 0 destroyed.
    
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
            "endpoint_ltm_policy": [
              "list",
              "string"
            ],
            "existing_monitor": "string",
            "existing_pool": "string",
            "existing_snat_pool": "string",
            "existing_tls_client_profile": "string",
            "existing_tls_server_profile": "string",
            "existing_waf_security_policy": "string",
            "fast_https_json": "string",
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
          "application": "myApp4",
          "endpoint_ltm_policy": null,
          "existing_monitor": "",
          "existing_pool": "",
          "existing_snat_pool": "",
          "existing_tls_client_profile": null,
          "existing_tls_server_profile": null,
          "existing_waf_security_policy": null,
          "fast_https_json": null,
          "id": "myApp4",
          "load_balancing_mode": "least-connections-member",
          "monitor": [
            {
              "interval": 0,
              "monitor_auth": false,
              "password": "",
              "response": "200 OK",
              "send_string": "GET / HTTP/1.1\\r\\nHost: example.com\\r\\nConnection: Close\\r\\n\\r\\n",
              "username": ""
            }
          ],
          "pool_members": [
            {
              "addresses": [
                "10.1.10.120",
                "10.1.10.121",
                "10.1.10.122"
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
            "10.1.10.50",
            "10.1.10.51",
            "10.1.10.52"
          ],
          "tenant": "scenario4",
          "tls_client_profile": [],
          "tls_server_profile": [
            {
              "tls_cert_name": "/Common/app4.crt",
              "tls_key_name": "/Common/app4.key"
            }
          ],
          "virtual_server": [
            {
              "ip": "10.1.10.224",
              "port": 443
            }
          ],
          "waf_security_policy": []
        }
      }
    }
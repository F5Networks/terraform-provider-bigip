.. _fast-integration-udp:

Scenario #1: Creating a UDP application
=======================================

The goal of this template is to deploy a new UDP application on BIG-IP using Terraform as the orchestrator.

Pre-requisites
--------------
On the BIG-IP:

- F5 BIG-IP version 16.1 or newer
- Credentials with REST API access

On Terraform:

- Using F5 BIG-IP provider version 1.16.0 or newer
- Using Hashicorp version following :ref:`versions`


Create a UDP application
------------------------

Create 4 files:

- main.tf
- variables.tf
- inputs.tfvars
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

.. code-block:: json
   :caption: main.tf
   :linenos:

   resource "bigip_fast_udp_app" "this" {
     application               = "myApp"
     tenant                    = "scenario1"
     virtual_server            {
       ip                        = "10.1.10.101"
       port                      = "80"
     }
     pool_members  	    {
       addresses                 = ["10.1.10.120", "10.1.10.121", "10.1.10.122"]
       port                      = "80"
     }
     load_balancing_mode       = "least-connections-member"
     existing_monitor          = "/Common/gateway_icmp"
     enable_fastl4		    = true
   }


Now, run the following commands, so we can:

1. Initialize the terraform project
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


   $ terraform plan -var-file=inputs.tfvars -out scenario1

   Terraform used the selected providers to generate the following execution plan.
   Resource actions are indicated with the following symbols:
     + create
   
   Terraform will perform the following actions:
   
     # bigip_fast_udp_app.this will be created
     + resource "bigip_fast_udp_app" "this" {
         + application         = "myApp"
         + enable_fastl4       = true
         + existing_monitor    = "/Common/gateway_icmp"
         + fast_udp_json       = (known after apply)
         + id                  = (known after apply)
         + load_balancing_mode = "least-connections-member"
         + tenant              = "scenario1"
   
         + pool_members {
             + addresses = [
                 + "10.1.10.120",
                 + "10.1.10.121",
                 + "10.1.10.122",
               ]
             + port      = 80
           }
   
         + virtual_server {
             + ip   = "10.1.10.101"
             + port = 80
           }
       }
   
   Plan: 1 to add, 0 to change, 0 to destroy.
   
   ───────────────────────────────────────────────────────────────────────────────
   
   Saved the plan to: scenario1
   
   To perform exactly these actions, run the following command to apply:
       terraform apply "scenario1"
   
   
   $ terraform apply "scenario1"
   bigip_fast_udp_app.this: Creating...
   bigip_fast_udp_app.this: Still creating... [10s elapsed]
   bigip_fast_udp_app.this: Creation complete after 15s [id=myApp]
   
   Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
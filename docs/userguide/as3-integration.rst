.. _as3-integration:

AS3 Integration with Terraform
==============================

You can use Terraform with AS3 for managing application-specific configurations on a BIG-IP system. AS3 uses a declarative model, meaning you provide a JSON declaration rather than a set of imperative commands. The declaration represents the configuration which AS3 is responsible for creating on a BIG-IP system. Terraform provides resources to configure AS3 declarative JSON on BIG-IP.

Prerequisites
-------------

To use AS3 Extensions with Terraform, ensure you meet the following requirements:
- The BIG-IP system is running software version 12.1.x or higher
- The BIG-IP system has AS3 Extension version 3.10 or higher installed
- You have a BIG-IP system user account with the Administrator role


Example
-------

.. code-block:: json

    resource "bigip_as3"  "as3-example" {
        as3_json = "${file("example.json")}"
        tenant_name = "as3"
    }


Argument Reference
------------------

- |as3_json| - (Required) Name of the Declarative AS3 JSON file

- tenant_name - (Required) The partition name where the application services will be configured

- example.json - Example of AS3 Declarative JSON


.. code-block:: json

    {
        "class": "AS3",
        "action": "deploy",
        "persist": true,
        "declaration": {
            "class": "ADC",
            "schemaVersion": "3.0.0",
            "id": "urn:uuid:33045210-3ab8-4636-9b2a-c98d22ab915d",
            "label": "Sample 1",
            "remark": "Simple HTTP application with RR pool",
            "as3": {
                "class": "Tenant",
                "A1": {
                    "class": "Application",
                    "template": "http",
                    "serviceMain": {
                        "class": "Service_HTTP",
                        "virtualAddresses": [
                            "10.0.1.10"
                        ],
                        "pool": "web_pool"
                    },
                    "web_pool": {
                        "class": "Pool",
                        "monitors": [
                            "http"
                        ],
                    "members": [{
                        "servicePort": 80,
                        "serverAddresses": [
                            "192.0.1.10",
                            "192.0.1.11"
                        ]
                    }]
                    }
                }
            }
        }
    }


AS3 Installation
----------------

Use the following terraform provisioner to download AS3 RPM from GitHub and install the RPM on BIG-IP.

::

   resource "null_resource" "install_as3" {

     provisioner "local-exec" {

        command = "./install-as3-rpm.sh x.x.x.x xxxx:xxxx"

        }

   }


You will need to pass BIG-IP and its credentials as an argument to the install script. This script is available in the `examples section <https://github.com/F5Networks/terraform-provider-bigip/tree/master/examples>`_ of AS3 in the Terraform repo.


.. NOTE::vAS3 tenants are BIG-IP administrative partitions used to group configurations that support specific AS3 applications. An AS3 application may support a network-based business application or system. AS3 tenants may also include resources shared by applications in other tenants.



.. |as3_json| raw:: html

   <a href="https://www.terraform.io/docs/providers/bigip/r/bigip_as3.html#as3_json" target="_blank">as3_json</a>


.. |tenant_name| raw:: html

   <a href="https://www.terraform.io/docs/providers/bigip/r/bigip_as3.html#tenant_name" target="_blank">tenant_name</a>


.. |example.json| raw:: html

   <a href="https://www.terraform.io/docs/providers/bigip/r/bigip_as3.html#example-json" target="_blank">example.json</a>
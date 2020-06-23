.. _as3-integration:

AS3 Integration with Terraform
==============================

You can use Terraform with AS3 for managing application-specific configurations on a BIG-IP system. AS3 uses a declarative model, meaning you provide a JSON declaration rather than a set of imperative commands. The declaration represents the configuration which AS3 is responsible for creating on a BIG-IP system. Terraform provides resources to configure AS3 declarative JSON on BIG-IP.

Prerequisites
-------------

To use AS3 Extensions with Terraform, ensure you meet the following requirements:

- The BIG-IP system is running software version 12.1.x or newer
- The BIG-IP system has AS3 Extension version 3.10 or newer installed
- You have a BIG-IP system user account with the Administrator role


Example Usage
-------------

.. code-block:: json
   :caption: Example usage for json file
   :linenos:

    resource "bigip_as3" "as3-example1" {
        Unknown macro: { as3_json = "${file("example1.json")}" 
    }



.. code-block:: json
   :caption: Example usage for json file with tenant filter
   :linenos:

    resource "bigip_as3" "as3-example1" {
        Unknown macro: { as3_json = "${file("example2.json")}" tenant_filter = "Sample_03" 
    }


Argument Reference
------------------


- `as3_json <https://www.terraform.io/docs/providers/bigip/r/bigip_as3.html#as3_json>`_ - (Required) Path/Filename of Declarative AS3 JSON which is a json file used with builtin ``file`` function

- `tenant_filter <https://www.terraform.io/docs/providers/bigip/r/bigip_as3.html#tenant_filter>`_ - (Optional) If there are muntiple tenants in a json this attribute helps the user to set a particular tenant to which he want to reflect the changes. Other tenants will neither be created nor be modified

- `as3_example1.json <https://www.terraform.io/docs/providers/bigip/r/bigip_as3.html#as3_example1-json>`_ - Example AS3 Declarative JSON file with single tenant


.. code-block:: json
   :linenos:

    {
        "class": "AS3",
        "action": "deploy",
        "persist": true,
        "declaration": {
            "class": "ADC",
            "schemaVersion": "3.0.0",
            "id": "example-declaration-01",
            "label": "Sample 1",
            "remark": "Simple HTTP application with round robin pool",
            "Sample_01": {
                "class": "Tenant",
                "defaultRouteDomain": 0,
                "Application_1": {
                    "class": "Application",
                    "template": "http",
                "serviceMain": {
                    "class": "Service_HTTP",
                    "virtualAddresses": [
                        "10.0.2.10"
                    ],
                    "pool": "web_pool"
                    },
                    "web_pool": {
                        "class": "Pool",
                        "monitors": [
                            "http"
                        ],
                        "members": [
                            {
                                "servicePort": 80,
                                "serverAddresses": [
                                    "192.0.1.100",
                                    "192.0.1.110"
                                ]
                            }
                        ]
                    }
                }
            }
        }
    }

- `as3_example2.json <https://www.terraform.io/docs/providers/bigip/r/bigip_as3.html#as3_example2-json>`_ - Example AS3 Declarative JSON file with multiple tenants

.. code-block:: json
   :linenos:

    
    {
        "class": "AS3",
        "action": "deploy",
        "persist": true,
        "declaration": {
            "class": "ADC",
            "schemaVersion": "3.0.0",
            "id": "example-declaration-01",
            "label": "Sample 1",
            "remark": "Simple HTTP application with round robin pool",
            "Sample_02": {
                "class": "Tenant",
                "defaultRouteDomain": 0,
                "Application_2": {
                    "class": "Application",
                    "template": "http",
                "serviceMain": {
                    "class": "Service_HTTP",
                    "virtualAddresses": [
                        "10.2.2.10"
                    ],
                    "pool": "web_pool2"
                    },
                    "web_pool2": {
                        "class": "Pool",
                        "monitors": [
                            "http"
                        ],
                        "members": [
                            {
                                "servicePort": 80,
                                "serverAddresses": [
                                    "192.2.1.100",
                                    "192.2.1.110"
                                ]
                            }
                        ]
                    }
                }
            },
            "Sample_03": {
                "class": "Tenant",
                "defaultRouteDomain": 0,
                "Application_3": {
                    "class": "Application",
                    "template": "http",
                "serviceMain": {
                    "class": "Service_HTTP",
                    "virtualAddresses": [
                        "10.1.2.10"
                    ],
                    "pool": "web_pool3"
                    },
                    "web_pool3": {
                        "class": "Pool",
                        "monitors": [
                            "http"
                        ],
                        "members": [
                            {
                                "servicePort": 80,
                                "serverAddresses": [
                                    "192.3.1.100",
                                    "192.3.1.110"
                                ]
                            }
                        ]
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


.. NOTE:: AS3 tenants are BIG-IP administrative partitions used to group configurations that support specific AS3 applications. An AS3 application may support a network-based business application or system. AS3 tenants may also include resources shared by applications in other tenants.
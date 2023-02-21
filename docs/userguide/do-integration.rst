.. _do-integration:

Declarative Onboarding Integration with Terraform
=================================================

Terraform is used with Declarative Onboarding (DO) to automate network (L1-L3) onboarding for the F5 BIG-IP system, making
it available and ready to accept application services configurations.

F5 BIG-IP Declarative Onboarding uses a declarative model to initially configure a BIG-IP device with all of the required settings to get up and running. This includes system settings such as licensing and provisioning, network settings such as VLANs and Self IPs, and clustering settings if you are using more than one BIG-IP system.

Prerequisites
-------------

To use DO Extensions with Terraform, ensure you meet the following requirements:

- The F5 BIG-IP must be running version 13.1.0 or higher. Due to changes in TMOS v13.1.1.5, the Declarative Onboarding (DO) Extension is not compatible with this specific TMOS version. Versions before and after 13.1.1.5 are compatible.
- Domain name resolution is used anywhere the declaration accepts a hostname. DO makes sure that any hostnames are resolvable and fails if they are not. The exception is deviceGroup.members, which do not require hostname resolution as they have been added to the trust.
- You must have an existing BIG-IP device with a management IP address
- You must have an existing user account with the Administrator role. If you are using 13.1.x, the BIG-IP contains an admin user by default. If you are using 14.x, you must reset the admin password before installing Declarative Onboarding. See `the documentation <https://clouddocs.f5.com/products/extensions/f5-declarative-onboarding/latest/installation.html#if-using-big-ip-14-0-or-later>`_ for instructions.
- While Declarative Onboarding is supported on F5 vCMP systems, network stitching to vCMP Guests or Hosts is not supported.
- If you are using an F5 BYOL license, you must have a valid F5 Networks License Registration Key to include in your declaration. If you do not have one, contact your F5 sales representative. If you do not use a valid F5 license key, your declaration will fail. This is not a requirement if you are using a BIG-IP with pay-as-you-go licensing.
- If you are using a single NIC BIG-IP system, you must include port 8443 after the IP address of the BIG-IP in your POST and GET requests, for example, ``https://<BIG-IP>:8443/mgmt/shared/declarative-onboarding``.

Limitations:

- DO does not support ``DELETE`` operation: ``terraform destroy`` will raise an error.
- If you POST a declaration that modifies the password for the admin account, even if the declaration returns an error, the password can be changed. Therefore you may need to update the admin password in the client you are using to send the declaration.
- The first time you POST a Declarative Onboarding declaration, the system records the configuration that exists prior to processing the declaration. Declarative Onboarding is meant to initially configure a BIG-IP device. However, if you POST subsequent declarations to the same BIG-IP system, and leave out some of the properties you initially used, the system restores the original properties for those items. 

.. IMPORTANT:: No matter what you send in a subsequent declaration, Declarative Onboarding will never unlicense a BIG-IP device, it will never delete a user, and it WILL never break the device trust once it has been established.

Example Usage
-------------

.. code-block:: json

    resource "bigip_do" "do-example" {
        do_json = "${file("example.json")}"
        timeout = 15
    }


Argument Reference
------------------

- `do_json <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_do#do_json>`_ - (Required) Name of the of the Declarative DO JSON file

- `timeout(minutes) <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_do#timeout(minutes)>`_ - (Optional) timeout to keep polling DO endpoint until the BIG-IP is provisioned by DO. The default timeout is 20 minutes

  .. NOTE:: Delete method is not supported by DO, so terraform destroy won't delete configuration in BIG-IP but we will set the terrform state to empty and won't throw an error.

- `example.json <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_do#example.json>`_ - Example of DO Declarative JSON


.. code-block:: json

    {
        "schemaVersion": "1.0.0",
        "class": "Device",
        "async": true,  
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip.example.com",
            "myLicense": {
                "class": "License",
                "licenseType": "regKey",
                "regKey": "xxxx"
            }, 
            "admin": {
                "class": "User",
                "userType": "regular",
                "password": "xxxx",
                "shell": "bash"
            },
            "myProvisioning": {
                "class": "Provision",
                "ltm": "nominal",
                "gtm": "minimum"
            },
            "external": {
                "class": "VLAN",
                "tag": 4093,
                "mtu": 1500,
                "interfaces": [
                    {
                        "name": "1.1",
                        "tagged": true
                    }
                ],
                "cmpHash": "dst-ip"
            },
            "external-self": {
                "class": "SelfIp",
                "address": "x.x.x.x",
                "vlan": "external",
                "allowService": "default",
                "trafficGroup": "traffic-group-local-only"
            }

        }
    }


DO Installation
----------------

Use the following terraform provisioner to download DO RPM from GitHub and install the RPM on BIG-IP.

::

    resource "null_resource" "install_do" {

    provisioner "local-exec" {

        command = "./install-do-rpm.sh x.x.x.x xxxx:xxxx"

        }

    }


You will need to pass BIG-IP and its credentials as an argument to the install script. This script is available in the `examples section <https://github.com/F5Networks/terraform-provider-bigip/tree/master/examples>`_ of DO in the Terraform repo.


.. NOTE:: DO tenants are BIG-IP administrative partitions used to group configurations and also resources shared by applications in other tenants.
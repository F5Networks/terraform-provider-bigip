.. _bigip-modules-gcp:

GCP
===
This Terraform module deploys N-NIC F5 BIG-IP in Google Cloud Provider (GCP). You can deploy multiple instances of BIG-IP with the module count feature.

Prerequisites
-------------
.. sidebar:: :fonticon:`fa fa-info-circle fa-lg` Version Notice:

   This module is supported on Terraform version 0.13 and newer.

.. seealso::

   `Getting Started with the Google Provider <https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/getting_started>`_

The templates below are tested on Terraform v0.14.0:

- provider registry.terraform.io/hashicorp/google v3.51.0
- provider registry.terraform.io/hashicorp/null v2.1.2
- provider registry.terraform.io/hashicorp/random v3.0.1
- provider registry.terraform.io/hashicorp/template v2.2.0

+-------------------------+----------------------+
| BIG-IP version          | Terraform v0.14      |
+=========================+======================+
| BIG-IP 15.x             | X                    |
+-------------------------+----------------------+
| BIG-IP 14.x             | X                    |
+-------------------------+----------------------+
| BIG-IP 13.x             | X                    |
+-------------------------+----------------------+


Password Management
-------------------
By default, the BIG-IP module dynamically generates passwords. Users can provide a password as input to the module using the optional variable ``f5_password``. To use GCP secret manager, you must enable the variable ``gcp_secret_manager_authentication`` to ``true`` and supply the variables with secret name and version.


Example Usage
-------------
.. seealso::

   `Common deployment examples <https://github.com/f5devcentral/terraform-gcp-bigip-module/tree/main/examples>`_.

You can use dynamic or static private IP allocation. If the primary or secondary private IP value is null, it will default to dynamic IP allocation. With static private IP allocation, you can assign primary and secondary private IPs for external interfaces. If you are using static private IP allocation, you cannot use module count because the same private IPs will be allocated for multiple BIG-IP instances based on module count. If you are using dynamic private IP allocation, you must pass a null value to primary/secondary private IP declaration and module count will be supported.

.. code-block:: javascript
   :caption: Example Deployment with Dynamic Private IP Allocation

    Example 1-NIC Deployment Module usage

    module bigip {
      count           = var.instance_count
      source          = "../.."
      prefix          = "bigip-gcp-1nic"
      project_id      = var.project_id
      zone            = var.zone
      image           = var.image
      service_account = var.service_account
      mgmt_subnet_ids = [{ "subnet_id" = google_compute_subnetwork.mgmt_subnetwork.id, "public_ip" = true, "private_ip_primary" = "" }]
    }
    
    Example 2-NIC Deployment Module usage
    
    module "bigip" {
      count               = var.instance_count
      source              = "../.."
      prefix              = "bigip-gcp-2nic"
      project_id          = var.project_id
      zone                = var.zone
      image               = var.image
      service_account     = var.service_account
      mgmt_subnet_ids     = [{ "subnet_id" = google_compute_subnetwork.mgmt_subnetwork.id, "public_ip" = true, "private_ip_primary" = "" }]
      external_subnet_ids = [{ "subnet_id" = google_compute_subnetwork.external_subnetwork.id, "public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = "" }]
    }
    
    
    Example 3-NIC Deployment  Module usage 
    
    module bigip {
      count               = var.instance_count
      source              = "../.."
      prefix              = "bigip-gcp-3nic"
      project_id          = var.project_id
      zone                = var.zone
      image               = var.image
      service_account     = var.service_account
      mgmt_subnet_ids     = [{ "subnet_id" = google_compute_subnetwork.mgmt_subnetwork.id, "public_ip" = true, "private_ip_primary" = "" }]
      external_subnet_ids = [{ "subnet_id" = google_compute_subnetwork.external_subnetwork.id, "public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = "" }]
      internal_subnet_ids = [{ "subnet_id" = google_compute_subnetwork.internal_subnetwork.id, "public_ip" = false, "private_ip_primary" = "", "private_ip_secondary" = "" }]
    }
    
        
    Example 4-NIC Deployment  Module usage(with 2 external public interfaces,one management and internal interfaces)
    
    module bigip s
      count               = vas.instance_count
      source              = "../.."
      prefix              = "bigip-gcp-4nic"
      project_id          = var.project_id
      zone                = var.zone
      image               = var.image
      service_account     = var.service_account
      mgmt_subnet_ids     = [{ "subnet_id" = google_compute_subnetwork.mgmt_subnetwork.id, "public_ip" = true, "private_ip_primary" = "" }]
      external_subnet_ids = ([{ "subnet_id" = google_compute_subnetwork.external_subnetwork.id, "public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = ""  },                                         { "subnet_id" = google_compute_subnetwork.external_subnetwork2.id, "public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = ""  }])
      internal_subnet_ids = [{ "subnet_id" = google_compute_subnetwork.internal_subnetwork.id, "public_ip" = false, "private_ip_primary" = "" }]
    }

    .............
    
    Similarly we can have N-nic deployments based on user-provided subnet_ids.
    With module count, you can deploy multiple BIG-IP instances in the GCP cloud (with the default value of count being 1).
    
|

.. code-block:: javascript
   :caption: Example Deployment for Private IP Allocation

    Example 3-NIC Deployment with Static Private IP Allocation

    module bigip {
      count               = var.instance_count
      source              = "../.."
      prefix              = "bigip-gcp-3nic"
      project_id          = var.project_id
      zone                = var.zone
      image               = var.image
      service_account     = var.service_account
      mgmt_subnet_ids     = [{ "subnet_id" = google_compute_subnetwork.mgmt_subnetwork.id, "public_ip" = true, "private_ip_primary" = "" }]
      external_subnet_ids = [{ "subnet_id" = google_compute_subnetwork.external_subnetwork.id, "public_ip" = true, "private_ip_primary" = "10.2.1.2", "private_ip_secondary" = "10.2.1.3" }]
      internal_subnet_ids = [{ "subnet_id" = google_compute_subnetwork.internal_subnetwork.id, "public_ip" = false, "private_ip_primary" = "", "private_ip_secondary" = "" }]
    }


|

Required and Optional Input Variables
`````````````````````````````````````
Required variables must be set in the module block when using this module. Optional variables have default values and do not need to be set to use this module. You may set these variables to override their default values.

+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| Parameter                   | Type    | Required | Default                     | Description                             |
+=============================+=========+==========+=============================+=========================================+
| prefix                      | String  | Required | N/A                         | This value is inserted in the beginning |
|                             |         |          |                             | of each GCP object.                     |
|                             |         |          |                             | Note: Requires alpha-numeric without    |
|                             |         |          |                             | special characters.                     |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| project_id                  | String  | Required | N/A                         | The GCP project identifier where the    |
|                             |         |          |                             | cluster will be created.                |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| zone                        | String  | Required | N/A                         | The compute zones which will host the   |
|                             |         |          |                             | BIG-IP Virtual Machines.                |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| mgmt_subnet_ids             | List of | Required | N/A                         | Map with Subnet-id and public_ip as     |
|                             | Maps    |          |                             | keys for the management subnet.         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| service_account             | String  | Required | N/A                         | Service account email to use with       |
|                             |         |          |                             | the BIG-IP system.                      |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| f5_username                 | String  | Optional | ``bigipuser``               | The admin username of the F5 BIG-IP     |
|                             |         |          |                             | that will be deployed.                  |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| f5_password                 | String  | Optional | m5.large                    | Password of the F5 BIG-IP that will be  |
|                             |         |          |                             | deployed. If this is not specified,     |
|                             |         |          |                             | a random password will be generated.    |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| image                       | String  | Optional | "projects/f5-7626-networks- | The self-link URI for a BIG-IP image    |
|                             |         |          | public/global/images/f5-    | to use as a base for the VM cluster.    |
|                             |         |          | bigip-16-0-1-1-0-0-6-payg-  |                                         |
|                             |         |          | good-25mbps-210129040032"   |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| min_cpu_platform            | String  | Optional | Intel Skylake               | Minimum CPU platform for the VM         |
|                             |         |          |                             | instance such as Intel Haswell or       |
|                             |         |          |                             | Intel Skylake.                          |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| machine_type                | String  | Optional | n1-standard-4               | The machine type to create. If you want |
|                             |         |          |                             | to update this value (resize the VM)    |
|                             |         |          |                             | after initial creation, you must set    |
|                             |         |          |                             | ``allow_stopping_for_update`` to        |
|                             |         |          |                             | ``true``.                               |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| automatic_restart           | Boolean | Optional | true                        | Specifies if the instance should be     |
|                             |         |          |                             | restarted if it was terminated by       |
|                             |         |          |                             | Compute Engine (not a user).            |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| preemptible                 | Boolean | Optional | false                       | Specifies if the instance is            |
|                             |         |          |                             | preemptible. If this field is set to    |
|                             |         |          |                             | true, then automatic_restart must be    |
|                             |         |          |                             | set to false.                           |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| disk_type                   | String  | Optional | pd-ssd                      | The GCE disk type. May be set to        |
|                             |         |          |                             | pd-standard, pd-balanced or pd-ssd.     |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| disk_size_gb                | Number  | Optional | null                        | The size of the image in gigabytes. If  |
|                             |         |          |                             | not specified, it will inherit the size |
|                             |         |          |                             | of its base image.                      |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| gcp_secret_manager_         | Boolean | Optional | false                       | Whether to use secret manager to pass   |
| authentication              |         |          |                             | authentication.                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| gcp_secret_name             | String  | Optional | null                        | The secret to get the secret version    |
|                             |         |          |                             | for.                                    |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| gcp_secret_version          | String  | Optional | latest                      | The version of the secret to get. If it |
|                             |         |          |                             | is not provided, the latest version is  |
|                             |         |          |                             | retrieved.                              |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| libs_dir                    | String  | Optional | /config/cloud/gcp/node      | Directory on the BIG-IP to download the |
|                             |         |          | _modules                    | A&O Toolchain into.                     |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| onboard_log	              | String  | Optional | /var/log/startup-script.log | Directory on the BIG-IP to store the    |
|                             |         |          |                             | cloud-init logs.                        |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| mgmt_subnet_ids             | List of | Optional | [{ "subnet_id" = null,      | The list of maps of subnet IDs of the   |
|                             | Maps    |          | "public_ip" = null,"private | virtual network where the virtual       |
|                             |         |          | _ip_primary" = "" }]        | machines will reside.                   |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| external_subnet_ids         | List of | Optional | [{ "subnet_id" = null,      | The list of maps of subnet IDs of the   |
|                             | Maps    |          | "public_ip" = null,"private | virtual network where the virtual       |
|                             |         |          | _ip_primary" = "", "private | machines will reside.                   |
|                             |         |          | _ip_secondary" = "" }]      |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| internal_subnet_ids         | List of | Optional | [{ "subnet_id" = null,      | The list of maps of subnet IDs of the   |
|                             | Maps    |          | "public_ip" = null,"private | virtual network where the virtual       |
|                             |         |          | _ip_primary" = "" }]        | machines will reside.                   |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| DO_URL                      | String  | Optional | latest                      | URL to download the BIG-IP Declarative  |
|                             |         |          |                             | Onboarding module.                      |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| AS3_URL                     | String  | Optional | latest                      | URL to download the BIG-IP Application  |
|                             |         |          |                             | Service Extension 3 (AS3) module.       |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| TS_URL                      | String  | Optional | latest                      | URL to download the BIG-IP Telemetry    |
|                             |         |          |                             | Streaming module.                       |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| FAST_URL                    | String  | Optional | latest                      | URL to download the BIG-IP FAST module. |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| CFE_URL                     | String  | Optional | latest                      | URL to download the BIG-IP Cloud        |
|                             |         |          |                             | Failover Extension module.              |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+
| INIT_URL                    | String  | Optional | latest                      | URL to download the BIG-IP runtime init |
|                             |         |          |                             | module.                                 |
|                             |         |          |                             |                                         |
|                             |         |          |                             |                                         |
+-----------------------------+---------+----------+-----------------------------+-----------------------------------------+



Output Variables
````````````````
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Parameter          | Description/Notes                                                                                                                                                                    |
+====================+======================================================================================================================================================================================+
| mgmtPublicIP       | The actual IP address allocated for the resource.                                                                                                                                    |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| mgmtPort           | The Mgmt Port.                                                                                                                                                                       |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| f5_username        | BIG-IP username.                                                                                                                                                                     |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| bigip_password     | The BIG-IP Password.                                                                                                                                                                 |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| public_addresses   | List of BIG-IP public addresses.                                                                                                                                                     |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| private_addresses  | List of BIG-IP private addresses.                                                                                                                                                    |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| service_account    | The service account that will be used for the BIG-IP VMs.                                                                                                                            |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
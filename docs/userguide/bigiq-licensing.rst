.. _bigiq-licensing:

BIG-IQ Licensing with Terraform
===============================

.. seealso::
   :class: sidebar

   `Terraform documentation <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html>`_.

In this section you can see examples of the ``bigip_common_license_manage_bigiq`` resource module. This resource is used for BIG-IP provider license managment from BIG-IQ.

There are 3 different managed modes:

- MANAGED
- UNMANAGED
- UNREACHABLE

Prerequisites
-------------

To license BIG-IQ with Terraform, ensure you meet the following requirements:

- The BIG-IP system is running software version 12.1.x or higher
- You have a BIG-IP system user account with the Administrator role


Example Usage
-------------

.. code-block:: json


    # MANAGED Regkey Pool
    resource "bigip_common_license_manage_bigiq" "test_example" {
      bigiq_address = var.bigiq
      bigiq_user = var.bigiq_un
      bigiq_password = var.bigiq_pw
      license_poolname = "regkeypool_name"
      assignment_type = "MANAGED"
    }

    # UNMANAGED Regkey Pool
    resource "bigip_common_license_manage_bigiq" "test_example" {
      bigiq_address = var.bigiq
      bigiq_user = var.bigiq_un
      bigiq_password = var.bigiq_pw
      license_poolname = "regkeypool_name"
      assignment_type = "UNMANAGED"
    } 

    # UNMANAGED Utility Pool
    resource "bigip_common_license_manage_bigiq" "test_example" {
      bigiq_address = var.bigiq
      bigiq_user = var.bigiq_un
      bigiq_password = var.bigiq_pw
      license_poolname = "utilitypool_name"
      assignment_type = "UNMANAGED"
      unit_of_measure = "yearly"
      skukeyword1 = "BTHSM200M"
    }

    # UNREACHABLE Regkey Pool
    resource "bigip_common_license_manage_bigiq" "test_example" {
      bigiq_address="xxx.xxx.xxx.xxx"
      bigiq_user="xxxx"
      bigiq_password="xxxxx"
      license_poolname = "regkey_pool_name"
      assignment_type = "UNREACHABLE"
      mac_address = "FA:16:3E:1B:6D:32"
      hypervisor = "azure"
    }


Argument Reference
------------------

- `big-iq_address <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#bigiq_address>`_ - (Required) BIGIQ License Manager IP Address, variable type ``string``

- `bigiq_user <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#bigiq_user>`_ - (Required) BIGIQ License Manager username, variable type ``string``

- `bigiq_password <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#bigiq_password>`_ - (Required) BIGIQ License Manager password. variable type ``string``

- `bigiq_port <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#bigiq_port>`_ - (Optional) type ``int``, BIGIQ License Manager Port number, specify if port is other than ``443``

- `bigiq_token_auth <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#bigiq_token_auth>`_ - (Optional) type ``bool``, if set to ``true`` enables Token based Authentication,default is ``false``

- `bigiq_login_ref <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#bigiq_login_ref>`_ - (Optional) BIGIQ Login reference for token authentication

- `assignment_type <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#assignment_type>`_ - (Required) The type of assignment, which is determined by whether the BIG-IP is unreachable, unmanaged, or managed by BIG-IQ. Possible values: “UNREACHABLE”, “UNMANAGED”, or “MANAGED”.

- `license_poolname <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#license_poolname>`_ - (Required) A name given to the license pool. type ``string``

- `unit_of_measure <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#unit_of_measure>`_ - (Optional, Required for ``Utility`` licenseType) The units used to measure billing. For example, “hourly” or “daily”. Type ``string``

- `skukeyword1 <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#skukeyword1>`_ - (Optional) An optional offering name. type ``string``

- `skukeyword2 <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#skukeyword2>`_ - (Optional) An optional offering name. type ``string``

- `mac_address <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#mac_address>`_ - (Optional, Required Only for ``unreachable BIG-IP``) MAC address of the BIG-IP. type ``string``

- `hypervisor <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#hypervisor>`_ - (Optional,Required Only for ``unreachable BIG-IP``) Identifies the platform running the BIG-IP VE. Possible values: “aws”, “azure”, “gce”, “vmware”, “hyperv”, “kvm”, or “xen”. type ``string``

- `tenant <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#tenant>`_ - (Optional) For an unreachable BIG-IP, you can provide an optional description for the assignment in this field.

- `key <https://www.terraform.io/docs/providers/bigip/r/bigip_common_license_manage_bigiq.html#key>`_ - Optional) License Assignment is done with specified ``key``, supported only with RegKeypool type License assignement. type ``string``
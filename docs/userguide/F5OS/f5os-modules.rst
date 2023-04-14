F5OS Tenant Reference
=============================

This topic describes how to manage the F5OS tenant.

.. _f5os_prereqs:

Prerequisites
--------------

- |f5OS_terraform-Dwnlds|  v0.13.x or later
- |f5OS_go|  v1.18 or later

.. _f5os_tenantEx:

F5OS tenant example
--------------------

The following example demonstrates a resource used to manage F5OS tenant.

.. code-block:: console

   resource "f5os_tenant" "test2" {
     name              = "testtenant-ecosys2"
     image_name        = "BIGIP-17.1.0-0.0.16.ALL-F5OS.qcow2.zip.bundle"
     mgmt_ip           = "10.10.10.26"
     mgmt_gateway      = "10.10.10.1"
     mgmt_prefix       = 24
     type              = "BIG-IP"
     cpu_cores         = 8
     running_state     = "configured"
     virtual_disk_size = 82
   }

**Required parameters**

- ``image_name`` (string) - Name of the tenant image used. Required for create operations.
- ``mgmt_gateway`` (string) - Tenant management gateway.
- ``mgmt_ip`` (string) - IP address used to connect to the deployed tenant. Required for create operations.
- ``mgmt_prefix`` (integer) - Tenant management CIDR prefix.
- ``name`` (string) - Name of the tenant. The first character must be a letter. Only lowercase alphanumeric characters
  are allowed. No special or extended characters are allowed except for hyphens. The name cannot exceed 50 characters.

**Optional parameters**

- ``cpu_cores`` (integer) - The number of vCPUs you want added to the tenant. Required for create operations.
- ``cryptos`` (string) - Whether crypto and compression hardware offload should be enabled on the tenant. We recommend it is enabled, otherwise crypto and compression may be processed in CPU.
- ``deployment_file`` (string) - Deployment file used for BIG-IP-Next . Required for if type is BIG-IP-Next.
- ``nodes`` (integer list) - List of integers. Specifies on which blades nodes the tenants are deployed. Required for create operations. For single blade platforms like rSeries only the value of 1 should be provided.
- ``running_state`` (string) - Desired running_state of the tenant.
- ``timeout`` (integer) - The number of seconds to wait for image import to finish.
- ``type`` (string) - Name of the tenant image to be used. Required for create operations
- ``virtual_disk_size`` (integer) - Minimum virtual disk size required for Tenant deployment
- ``vlans`` (integer list) - The existing VLAN IDs in the chassis partition that should be added to the tenant. The order
  of these VLANs is ignored. This module orders the VLANs automatically, if you deliberately re-order them in subsequent tasks, this module will not register a change. Required for create operations

**Read-only parameters**

- ``id`` (string) - Tenant identifier
- ``status`` (string) - Tenant status

.. _f5os_tenantExImg:

F5OS tenant image example
---------------------------

The following example demonstrates a resource used to manage F5OS tenant image.

.. code-block:: console

    resource "f5os_tenant_image" "test" {
      image_name  = "BIGIP-17.1.0-0.0.16.ALL-F5OS.qcow2.zip.bundle"
      remote_host = "remote-host"
      remote_path = "remote-path"
      local_path  = "images"
      timeout     = 360
    }

**Required parameters**

- ``image_name`` (string) - Name of the tenant image.

**Optional parameters**

- ``local_path`` (string) - The path on the F5OS where the the tenant image is to be uploaded.
- ``protocol`` (string) - Protocol for image transfer.
- ``remote_host`` (string) - The hostname or IP address of the remote server on which the tenant image is stored. The server must make the image accessible via the specified protocol.
- ``remote_password`` (string, sensitive) Password for the user on the remote server on which the tenant image is stored.
- ``remote_path`` (string) - The path to the tenant image on the remote server.
- ``remote_port`` (integer) - The port on the remote host to which you want to connect. If the port is not provided, a default port for the selected protocol is used.
- ``remote_user`` (string) - User name for the remote server on which the tenant image is stored.
- ``timeout`` (integer) - The number of seconds to wait for image import to finish.

**Read-only parameters**

- ``id`` (string) - Example identifier.
- ``status`` (string) - Status of imported image




What’s Next?

- :doc:`Quick start guide <f5os-qsg>`
- :doc:`Support <../support>`







.. |f5OS_go| raw:: html

   <a href="https://golang.org/doc/install" target="_blank">Go</a>



.. |f5OS_terraform-Dwnlds| raw:: html

   <a href="https://www.terraform.io/downloads.html" target="_blank">Terraform</a>
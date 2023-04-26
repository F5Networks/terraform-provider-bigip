F5OS Resources for Terraform
=============================

Welcome to the F5OS Resources for Terraform User Guide. The F5 F5OS is the operating system for F5 VELOS and F5 rSeries.

The VELOS platform is the next generation of F5’s industry-leading chassis-based systems, which delivers unprecedented
performance and scalability in a single Application Delivery Controller (ADC). The next-generation Application Delivery
Controller (ADC) solution, F5 rSeries, bridges the gap between traditional and modern infrastructures with a re-architected,
API-first platform designed to meet the needs of your traditional and emerging applications.

F5 VELOS and F5 rSeries rely on a Kubernetes-based platform layer (F5OS) that is tightly integrated with F5’s Traffic
Management Operating System (TMOS) software, aligning with your modern architecture plans. This new microservices platform
layer powers the next-generation of BIG-IP software, BIG-IP Next, which is built to offer greater automatability, scalability,
and ease-of-use for organizations running applications on-premises, in the cloud, or at the edge.

F5OS Terraform provider for F5 VELOS and F5 rSeries helps you automate configurations and interactions with various services
provided by F5 VELOS platform and F5 rSeries appliances. F5OS Terraform provider is open source and |f5_terraform_F5OSgithub|.

For more information:

- |f5_velos| and |f5_velosOvr|
- |f5_rSeries| and |f5_rSeriesOvr|


.. _versions-F5os:

Releases and Versioning
-----------------------
These F5OS versions are supported in these Terraform versions.

**F5 VELOS**

+-------------------------+----------------------+----------------------+
| VELOS F5OS version      | Terraform 0.14       | Terraform 0.13       |
+=========================+======================+======================+
| 1.4.1                   | X                    | X                    |
+-------------------------+----------------------+----------------------+
| 1.3.1                   | X                    | X                    |
+-------------------------+----------------------+----------------------+

**F5 rSeries**

+-------------------------+----------------------+----------------------+
| rSeries F5OS version    | Terraform 0.14       | Terraform 0.13       |
+=========================+======================+======================+
| 1.2.0                   | X                    | X                    |
+-------------------------+----------------------+----------------------+
| 1.1.0                   | X                    | X                    |
+-------------------------+----------------------+----------------------+

**Supported Resources**

F5OS Terraform provider currently supports the following modules for F5OS-based platforms.

+---------------------------------+--------------------------------------------------------------------------------------------+
| Resource Name                   | Description                                                                                |
+=================================+============================================================================================+
| velos_partition_image           | Upload/copy partition image on VELOS chassis.                                              |
+---------------------------------+--------------------------------------------------------------------------------------------+
| velos_partition                 | Manage partitions on the VELOS chassis.                                                    |
+---------------------------------+--------------------------------------------------------------------------------------------+
| velos_partition_wait            | Waits for the specified timeout value for the partition to match the desired state.        |
+---------------------------------+--------------------------------------------------------------------------------------------+
| velos_partition_change_password | Manage password of a particular user on a partition.                                       |
+---------------------------------+--------------------------------------------------------------------------------------------+
| f5os_interface                  | Configure properties related to physical interfaces on a partition \* or appliance.        |
+---------------------------------+--------------------------------------------------------------------------------------------+
| f5os_lag                        | Manage LAG interfaces - trunk/native vlans and physical interfaces to the LAG interfaces.  |
+---------------------------------+--------------------------------------------------------------------------------------------+
| f5os_vlan                       | Manage VLANs on a partition or appliance.                                                  |
+---------------------------------+--------------------------------------------------------------------------------------------+
| f5os_tenant_image               | Manage tenant images on a partition or appliance.                                          |
+---------------------------------+--------------------------------------------------------------------------------------------+
| f5os_tenant                     | Manage F5OS tenant configuration.                                                          |
+---------------------------------+--------------------------------------------------------------------------------------------+
| velos_tenant_wait               | Waits for a tenant to be in a desired running state.                                       |
+---------------------------------+--------------------------------------------------------------------------------------------+

\* Use this module with partitions that have existing physical interfaces.

What’s Next?

- :doc:`F5OS Terraform Tenant Reference <f5os-qsg>`
- :doc:`Quick start guide <f5os-qsg>`






.. |f5_velos| raw:: html

   <a href="https://www.f5.com/products/big-ip-services/velos-hardware-chassis-and-blades" target="_blank">F5 VELOS hardware</a>

.. |f5_velosOvr| raw:: html

   <a href="https://techdocs.f5.com/en-us/velos-1-1-0/velos-systems-administration-configuration/title-velos-system-overview.html" target="_blank">system overview</a>

.. |f5_rSeries| raw:: html

   <a href="https://www.f5.com/products/big-ip-services/rseries-adc-hardware-appliance" target="_blank">F5 rSeries hardware</a>

.. |f5_rSeriesOvr| raw:: html

   <a href="https://techdocs.f5.com/en-us/hardware/f5-rseries-systems-getting-started.html" target="_blank">system overview</a>

.. |f5_terraform_github| raw:: html

   <a href="https://github.com/F5Networks/terraform-provider-bigip" target="_blank">available on GitHub</a>

.. |f5_terraform_F5OSgithub| raw:: html

   <a href="https://github.com/F5Networks/terraform-provider-F5OS" target="_blank">available on GitHub</a>


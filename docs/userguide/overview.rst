F5 BIG-IP Resources for Terraform Overview
==========================================

Welcome to the F5 BIG-IP Resources for Terraform User Guide. Terraform is a tool for building, changing, and versioning infrastructure safely and efficiently. Terraform can manage existing and popular service providers as well as custom in-house solutions.

Use these resources to create, edit, update, and delete configuration objects on BIG-IP 12.1.1 and later.

The code is open source and `available on GitHub <https://github.com/F5Networks/terraform-provider-bigip>`_.

Configuration files describe to Terraform the components needed to run a single application or your entire datacenter. Terraform generates an execution plan describing what it will do to reach the desired state, and then executes it to build the described infrastructure. As the configuration changes, Terraform is able to determine what changed and create incremental execution plans which can be applied.

Terraform can manage infrastructure including low-level components such as compute instances, storage, and networking, as well as high-level components such as DNS entries, SaaS features, and more.

.. _versions:

Releases and Versioning
-----------------------
These BIG-IP versions are supported in these Terraform versions.

+-------------------------+----------------------+----------------------+----------------------+----------------------+
| F5 BIG-IP version       | Terraform 0.14       | Terraform 0.13       | Terraform 0.12       | Terraform 0.11       |
+=========================+======================+======================+======================+======================+
| BIG-IP 16.x             | X                    | X                    | X                    | X                    | 
+-------------------------+----------------------+----------------------+----------------------+----------------------+
| BIG-IP 15.x             | X                    | X                    | X                    | X                    | 
+-------------------------+----------------------+----------------------+----------------------+----------------------+
| BIG-IP 14.x             | X                    | X                    | X                    | X                    |
+-------------------------+----------------------+----------------------+----------------------+----------------------+
| BIG-IP 13.x             | X                    | X                    | X                    | X                    | 
+-------------------------+----------------------+----------------------+----------------------+----------------------+
| BIG-IP 12.x             | X                    | X                    | X                    | X                    | 
+-------------------------+----------------------+----------------------+----------------------+----------------------+


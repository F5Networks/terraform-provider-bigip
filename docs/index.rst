F5 Modules for Terraform
========================

Welcome to the F5 Modules for Terraform documentation. Terraform is a tool for building, changing, and versioning infrastructure safely and efficiently. Terraform can manage existing and popular service providers as well as custom in-house solutions.

Use these modules to create, edit, update, and delete configuration objects on BIG-IP 12.1.1 and later.

The code is open source and |f5_terraform_github|.

Configuration files describe to Terraform the components needed to run a single application or your entire datacenter. Terraform generates an execution plan describing what it will do to reach the desired state, and then executes it to build the described infrastructure. As the configuration changes, Terraform is able to determine what changed and create incremental execution plans which can be applied.

The infrastructure Terraform can manage includes low-level components such as compute instances, storage, and networking, as well as high-level components such as DNS entries, SaaS features, and more.

User Guide Index
----------------

.. toctree::
   :maxdepth: 2
   :includehidden:
   :glob:

   /userguide/installing
   /userguide/configuring
   /userguide/modifying
   /userguide/support
   /userguide/revision_history



.. |f5_terraform_github| raw:: html

   <a href="https://github.com/F5Networks/terraform-provider-bigip" target="_blank">available on GitHub</a>


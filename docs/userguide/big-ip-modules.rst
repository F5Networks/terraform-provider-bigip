BIG-IP Modules
==============
This page covers modules for :ref:`big-ip-modules-aws` and :ref:`big-ip-modules-azure`.

.. _big-ip-modules-aws:

AWS
---
This Terraform module deploys N-NIC F5 BIG-IP in AWS cloud. You can deploy multiple instances of BIG-IP with the module count feature.

Prerequisites
`````````````
This module is supported on Terraform version 0.13 and newer.

The templates below are tested on Terraform v0.13.0

- provider registry.terraform.io/hashicorp/aws v3.8.0
- provider registry.terraform.io/hashicorp/random v2.3.0
- provider registry.terraform.io/hashicorp/template v2.1.2
- provider registry.terraform.io/hashicorp/null v2.1.2

|

+-------------------------+----------------------+
| BIG-IP version          | Terraform v0.13      |
+=========================+======================+
| BIG-IP 15.x             | X                    | 
+-------------------------+----------------------+
| BIG-IP 14.x             | X                    |
+-------------------------+----------------------+
| BIG-IP 13.x             | X                    |
+-------------------------+----------------------+

Password Management
```````````````````
By default, the BIG-IP module dynamically generates passwords.

.. code-block:: javascript

    cat terraform-aws-bigip-module/variables.tf

    variable aws_secretmanager_auth {
      description = "Whether to use key vault to pass authentication"
      type        = bool
      default     = false
    }

    Outputs:

    bigip_password = [
      "xxxxxxxxxxxxxxxxxx",
    ]

To use AWS secret manager password, you must enable the variable ``aws_secretmanager_auth`` to ``true`` and supply the secret name to variable ``aws_secretmanager_secret_id``.

.. code-block:: javascript

    cat terraform-aws-bigip-module/variables.tf

    variable aws_secretmanager_auth {
      description = "Whether to use key vault to pass authentication"
      type        = bool
      default     = true
    }

    variable aws_secretmanager_secret_id {
      description = "AWS Secret Manager Secret ID that stores the BIG-IP password"
      type        = string
      default     = "tf-aws-bigip-bigip-secret-9759"
    } 

    Outputs:

    bigip_password = [
      "xxxxxxxxxxxxxxx",
    ]


Example Usage
`````````````
`See more common deployment examples here<https://github.com/f5devcentral/terraform-aws-bigip-module/tree/master/examples>`_. 

There should be one-to-one mapping between subnet_ids and securitygroup_ids. For example if you have two or more external subnet_ids, you must give the same number of external securitygroup_ids to the module.

Users can have dynamic or static private IP allocation. If the primary/secondary private IP value is null, it will be dynamic or else static private IP allocation. With Static private IP allocation, you can assign primary and secondary private IPs for external interfaces, whereas the primary private IP is for management
and internal interfaces.

If it is static private ip allocation we can't use module count as same private ips will be tried to allocate for multiple bigip instances based on module count.

With Dynamic private ip allocation,we have to pass null value to primary/secondary private ip declaration and module count will be supported.

.. Note:: Sometimes the given static primary and secondary private IPs may get exchanged. This limitation is present in AWS.


.. code-block:: javascript
   :caption: Dynamic Private IP Allocation

    #
    #Example 1-NIC Deployment Module usage
    #
    module bigip {
      count                  = var.instance_count
      source                 = "../../"
      prefix                 = "bigip-aws-1nic"
      ec2_key_name           = aws_key_pair.generated_key.key_name
      mgmt_subnet_ids        = [{ "subnet_id" = "subnet_id_mgmt", "public_ip" = true, "private_ip_primary" =  ""}]
      mgmt_securitygroup_ids = ["securitygroup_id_mgmt"]
    }

    #
    #Example 2-NIC Deployment Module usage
    #
    module bigip {
      count                  = var.instance_count
      source                      = "../../"
      prefix                      = "bigip-aws-2nic"
      ec2_key_name                = aws_key_pair.generated_key.key_name
      mgmt_subnet_ids             = [{ "subnet_id" = "subnet_id_mgmt", "public_ip" = true, "private_ip_primary" =  ""}]
      mgmt_securitygroup_ids      = ["securitygroup_id_mgmt"]
      external_subnet_ids         = [{ "subnet_id" = "subnet_id_external", "public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = ""}]
      external_securitygroup_ids  = ["securitygroup_id_external"]
    }

    #
    #Example 3-NIC Deployment  Module usage
    #
    module bigip {
      count                  = var.instance_count
      source                      = "../../"
      prefix                      = "bigip-aws-3nic"
      ec2_key_name                = aws_key_pair.generated_key.key_name
      mgmt_subnet_ids             = [{ "subnet_id" = "subnet_id_mgmt", "public_ip" = true, "private_ip_primary" =  ""}]
      mgmt_securitygroup_ids      = ["securitygroup_id_mgmt"]
      external_subnet_ids         = [{ "subnet_id" = "subnet_id_external", "public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = ""}]
      external_securitygroup_ids  = ["securitygroup_id_external"]
      internal_subnet_ids         = [{"subnet_id" =  "subnet_id_internal", "public_ip"=false, "private_ip_primary" = ""}]
      internal_securitygroup_ids  = ["securitygropu_id_internal"]
    }

    #
    #Example 4-NIC Deployment  Module usage(with 2 external public interfaces,one management and internal interface.There should be one to one mapping between subnet_ids and securitygroupids)
    #

    module bigip {
      count                  = var.instance_count
      source                      = "../../"
      prefix                      = "bigip-aws-4nic"
      ec2_key_name                = aws_key_pair.generated_key.key_name
      mgmt_subnet_ids             = [{ "subnet_id" = "subnet_id_mgmt", "public_ip" = true }]
      mgmt_securitygroup_ids      = ["securitygroup_id_mgmt"]
      external_subnet_ids         = [{ "subnet_id" = "subnet_id_external", "public_ip" = true },{"subnet_id" =  "subnet_id_external2", "public_ip" = true }]
      external_securitygroup_ids  = ["securitygroup_id_external","securitygroup_id_external"]
      internal_subnet_ids         = [{"subnet_id" =  "subnet_id_internal", "public_ip"=false }]
      internal_securitygroup_ids  = ["securitygropu_id_internal"]
    }

Similarly, you can have N-nic deployments based on user provided subnet_ids and securitygroup_ids. With module count, you can deploy multiple BIG-IP instances in the AWS cloud (with the default value of count being one).



.. code-block:: javascript
   :caption: Private IP Allocation

    Example 3-NIC Deployment with static private ip allocation

    module bigip {
      source                      = "../../"
      count                       = var.instance_count
      prefix                      = format("%s-3nic", var.prefix)
      ec2_key_name                = aws_key_pair.generated_key.key_name
      aws_secretmanager_secret_id = aws_secretsmanager_secret.bigip.id
      mgmt_subnet_ids             = [{ "subnet_id" = aws_subnet.mgmt.id, "public_ip" = true, "private_ip_primary" = "10.0.1.4"}]
      mgmt_securitygroup_ids      = [module.mgmt-network-security-group.this_security_group_id]
      external_securitygroup_ids  = [module.external-network-security-group-public.this_security_group_id]
      internal_securitygroup_ids  = [module.internal-network-security-group-public.this_security_group_id]
      external_subnet_ids         = [{ "subnet_id" = aws_subnet.external-public.id, "public_ip" = true, "private_ip_primary" = "10.0.2.4", "private_ip_secondary" = "10.0.2.5"}]
      internal_subnet_ids         = [{ "subnet_id" = aws_subnet.internal.id, "public_ip" = false, "private_ip_primary" = "10.0.3.4"}]
    }


InSpec Tool
```````````
The BIG-IP Automation Toolchain InSpec Profile is used for testing the readiness of Automation Tool Chain (ATC) components. After the module deployment, you can use the InSpec tool to verify BIG-IP connectivity with ATC components.

This InSpec profile evaluates the following:

- Basic connectivity to a BIG-IP management endpoint: ``bigip-connectivity``
- Availability of the Declarative Onboarding (DO) service: ``bigip-declarative-onboarding``
- Version reported by the Declarative Onboarding (DO) service: ``bigip-declarative-onboarding-version``
- Availability of the Application Services (AS3) service: ``bigip-application-services``
- Version reported by the Application Services (AS3) service: ``bigip-application-services-version``
- Availability of the Telemetry Streaming (TS) service: ``bigip-telemetry-streaming``
- Version reported by the Telemetry Streaming (TS) service: ``bigip-telemetry-streaming-version``
- Availability of the Cloud Failover Extension (CFE) service: ``bigip-cloud-failover-extension``
- Version reported by the Cloud Failover Extension (CFE) service: ``bigip-cloud-failover-extension-version``


To run InSpec tests, you can either run the inspec exec command or execute runtests.sh in any one of example NIC folders which will run below the inspec command. For example:

``inspec exec inspec/bigip-ready --input bigip_address=$BIGIP_MGMT_IP bigip_port=$BIGIP_MGMT_PORT user=$BIGIP_USER password=$BIGIP_PASSWORD do_version=$DO_VERSION as3_version=$AS3_VERSION ts_version=$TS_VERSION fast_version=$FAST_VERSION cfe_version=$CFE_VERSION``


Required and Optional Input Variables
`````````````````````````````````````
Required variables must be set in the module block when using this module. Optional variables have default values and do not need to be set to use this module. You may set these variables to override their default values.

+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| Parameter                   | Type    | Required | Default               | Description                             |
+=============================+=========+==========+=======================+=========================================+
| prefix                      | String  | Required | N/A                   | This value is inserted in the beginning |
|                             |         |          |                       | of each AWS object.                     |
|                             |         |          |                       | Note: Requires alpha-numeric without    |
|                             |         |          |                       | special characters.                     |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| ec2_key_name	              | String  | Required | N/A                   | AWS EC2 Key name for SSH access.        |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| mgmt_subnet_ids             | List of | Required | N/A                   | Map with Subnet-id and public_ip as     |
|                             | maps    |          |                       | keys for the management subnet.         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| mgmt_securitygroup_ids      | List    | Required | N/A                   | securitygroup_ids for the management    |
|                             |         |          |                       | interface.                              |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| instance_count              | Number  | Required | false                 | Number of BIG-IP instances to spin up.  |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| f5_username                 | String  | Optional | bigipuser             | The admin username of the F5 BIG-IP     |
|                             |         |          |                       | that will be deployed                   |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| ec2_instance_type           | String  | Optional | m5.large              | AWS EC2 instance type.                  |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| f5_ami_search_name	      | String  | Optional | ``F5 Networks``       | BIG-IP AMI name to search for.          |
|                             |         |          | ``BIGIP-14.* PAYG``   |                                         |
|                             |         |          | ``- Best 200Mbps*``   |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| mgmt_eip                    | Boolean | Optional | True                  | Enable an Elastic IP address on the     |
|                             |         |          |                       | management interface.                   |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| aws_secretmanager_auth      | Boolean | Optional | False                 | Whether to use key vault to pass        |
|                             |         |          |                       | authentication.                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| aws_secretmanager_secret_id | String  | Optional | N/A                   | AWS Secret Manager Secret ID that       |
|                             |         |          |                       | stores the BIG-IP password.             |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| aws_iam_instance_profile    | String  | Optional | N/A                   | AWS IAM instance profile that can be    |
|                             |         |          |                       | associated for BIG-IP with required     |
|                             |         |          |                       | permissions.                            |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| DO_URL                      | String  | Optional | latest                | URL to download the BIG-IP Declarative  |
|                             |         |          |                       | Onboarding module.                      |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| AS3_URL                     | String  | Optional | latest                | URL to download the BIG-IP Application  |
|                             |         |          |                       | Service Extension 3 (AS3) module.       |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| TS_URL                      | String  | Optional | latest                | URL to download the BIG-IP Telemetry    |
|                             |         |          |                       | Streaming module.                       |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| fastPackageUrl              | String  | Optional | latest                | URL to download the BIG-IP FAST module. |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| CFE_URL                     | String  | Optional | latest                | URL to download the BIG-IP Cloud        |
|                             |         |          |                       | Failover Extension module.              |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| libs_dir                    | String  | Optional | /config/cloud/aws     | Directory on the BIG-IP to download the |
|                             |         |          | /node_modules         | A&O Toolchain into.                     |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| onboard_log	              | String  | Optional | /var/log/startup      | Directory on the BIG-IP to store the    |
|                             |         |          | -script.log           | cloud-init logs.                        |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| external_subnet_ids         | List of | Optional | ``[{ "subnet_id" =``  | The subnet ID of the virtual network    |
|                             | Maps    |          | ``null, "public_ip"`` | where the virtual machines will reside. |
|                             |         |          | ``= null }]``         |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| internal_subnet_ids         | List of | Optional | ``[{ "subnet_id" =``  | The subnet ID of the virtual network    |
|                             | Maps    |          | ``null, "public_ip"`` | where the virtual machines will reside. |
|                             |         |          | ``= null }]``         |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| external_securitygroup_ids  | List    | Optional | ``[]``                | The Network Security Group IDs for      |
|                             |         |          |                       | external network.                       |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+
| internal_securitygroup_ids  | List    | Optional | ``[]``                | The Network Security Group IDs for      |
|                             |         |          |                       | internal network.                       |
|                             |         |          |                       |                                         |
|                             |         |          |                       |                                         |
+-----------------------------+---------+----------+-----------------------+-----------------------------------------+

.. Note:: For each external interface there will be one primary, secondary private IP will be assigned.

Output Variables
````````````````
+--------------------+---------------------------------------------------------------------------------------------------------------------------+
| Parameter          | Description/Notes                                                                                                         |
+====================+===========================================================================================================================+
| mgmtPublicIP       | Describes the name of the policy.                                                                                         |
+--------------------+---------------------------------------------------------------------------------------------------------------------------+
| mgmtPublicDNS      | This value specifies the match strategy.                                                                                  |
+--------------------+---------------------------------------------------------------------------------------------------------------------------+
| mgmtPort           | This value specifies the protocol.                                                                                        |
+--------------------+---------------------------------------------------------------------------------------------------------------------------+
| f5_username        | This value determines if you want to publish the policy else it will be deployed in Drafts mode.                          |
+--------------------+---------------------------------------------------------------------------------------------------------------------------+
| bigip_password     | This value specifies the controls.                                                                                        |
+--------------------+---------------------------------------------------------------------------------------------------------------------------+
| private_addresses  | Use this policy to apply rules.                                                                                           |
+--------------------+---------------------------------------------------------------------------------------------------------------------------+
| public_addresses   | If Rule is used, then you need to provide the tm_name. It can be any value.                                               |
+--------------------+---------------------------------------------------------------------------------------------------------------------------+

.. Note:: A local json file that contains the DO declaration will be generated.

|

.. _big-ip-modules-azure:

Azure
-----
This Terraform module deploys N-nic F5 BIG-IP in Azure cloud. You can deploy multiple instances of BIG-IP with the module count feature.

Prerequisites
`````````````
This module is supported on Terraform version 0.13 and newer.

The templates below are tested on Terraform v0.13.0:

- provider registry.terraform.io/hashicorp/azurerm v2.28.0
- provider registry.terraform.io/hashicorp/null v2.1.2
- provider registry.terraform.io/hashicorp/random v2.3.0
- provider registry.terraform.io/hashicorp/template v2.1.2


|

+-------------------------+----------------------+
| BIG-IP version          | Terraform v0.13      |
+=========================+======================+
| BIG-IP 15.x             | X                    | 
+-------------------------+----------------------+
| BIG-IP 14.x             | X                    |
+-------------------------+----------------------+
| BIG-IP 13.x             | X                    |
+-------------------------+----------------------+

Password Management
```````````````````
By default, the BIG-IP module dynamically generates passwords.

.. code-block:: javascript

    variable az_key_vault_authentication {
      description = "Whether to use key vault to pass authentication"
      type        = bool
      default     = false
    }

    Outputs:
    bigip_password = [
      "xxxxxxxxxxxxxxxxxx",
    ]
    
To use Azure secret key vault, you must enable the variable ``az_key_vault_authentication`` to ``true`` and supply the variables (shown below) with key_vault secret name along with resource group name where the Azure key vault is defined.

.. code-block:: javascript

    variable az_key_vault_authentication {
      description = "Whether to use key vault to pass authentication"
      type        = bool
      default     = false
    }

    variable azure_secret_rg {
      description = "The name of the resource group in which the Azure Key Vault exists"
      type        = string
      default     = ""
    }

    variable azure_keyvault_name {
      description = "The name of the Azure Key Vault to use"
      type        = string
      default     = ""
    }

    variable azure_keyvault_secret_name {
      description = "The name of the Azure Key Vault secret containing the password"
      type        = string
      default     = ""
    }

    Outputs:
    bigip_password = [
      "xxxxxxxxxxxxxxxxxx",
    ]


Example Usage
`````````````
`See more common deployment examples here<https://github.com/f5devcentral/terraform-azure-bigip-module/tree/master/examples>`_. 

There should be one-to-one mapping between subnet_ids and securitygroup_ids. For example, if you have two or more external subnet_ids, you must give the same number of external securitygroup_ids to the module.

Users can have dynamic or static private IP allocation. If the primary/secondary private IP value is null, it will be dynamic or else static private IP allocation. With Static private IP allocation, you can assign primary and secondary private IPs for external interfaces, whereas the primary private IP is for management
and internal interfaces.

If it is static private IP allocation, you cannot use module count as same private IPs will be tried to allocate for multiple BIG-IP instances based on module count.

With Dynamic Private IP Allocation, you must pass null value to primary/secondary private IP declaration and module count will be supported.

.. code-block:: javascript
   :caption: Example of 1-NIC Deployment with Dynamic Private IP Allocation

    Example 1-NIC Deployment Module usage

    module bigip {
      count 		      = var.instance_count
      source                      = "../../"
      prefix                      = "bigip-azure-1nic"
      resource_group_name         = "testbigip"
      mgmt_subnet_ids             = [{"subnet_id" = "subnet_id_mgmt" , "public_ip" = true,"private_ip_primary" =  ""}]
      mgmt_securitygroup_ids      = ["securitygroup_id_mgmt"]
      availabilityZones           =  var.availabilityZones


    }


    Example 2-NIC Deployment Module usage

    module bigip {
      count                       = var.instance_count
      source                      = "../../"
      prefix                      = "bigip-azure-2nic"
      resource_group_name         = "testbigip"
      mgmt_subnet_ids             = [{"subnet_id" = "subnet_id_mgmt" , "public_ip" = true, "private_ip_primary" =  ""}]
      mgmt_securitygroup_ids      = ["securitygroup_id_mgmt"]
      external_subnet_ids         = [{"subnet_id" =  "subnet_id_external", "public_ip" = true,"private_ip_primary" = "", "private_ip_secondary" = "" }]
      external_securitygroup_ids  = ["securitygroup_id_external"]
      availabilityZones           =  var.availabilityZones
    }


    Example 3-NIC Deployment  Module usage 

    module bigip {
      count                       = var.instance_count 
      source                      = "../../"
      prefix                      = "bigip-azure-3nic"
      resource_group_name         = "testbigip"
      mgmt_subnet_ids             = [{"subnet_id" = "subnet_id_mgmt" , "public_ip" = true, "private_ip_primary" =  ""}]
      mgmt_securitygroup_ids      = ["securitygroup_id_mgmt"]
      external_subnet_ids         = [{"subnet_id" =  "subnet_id_external", "public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = "" }]
      external_securitygroup_ids  = ["securitygroup_id_external"]
      internal_subnet_ids         = [{"subnet_id" =  "subnet_id_internal", "public_ip"=false, "private_ip_primary" = "" }]
      internal_securitygroup_ids  = ["securitygropu_id_internal"]
      availabilityZones           =  var.availabilityZones
    }

    Example 4-NIC Deployment  Module usage(with 2 external public interfaces,one management and internal interface.There should be one to one mapping between subnet_ids and securitygroupids)

    module bigip {
      count                       = var.instance_count
      source                      = "../../"
      prefix                      = "bigip-azure-4nic"
      resource_group_name         = "testbigip"
      mgmt_subnet_ids             = [{"subnet_id" = "subnet_id_mgmt" , "public_ip" = true, "private_ip_primary" =  ""}]
      mgmt_securitygroup_ids      = ["securitygroup_id_mgmt"]
      external_subnet_ids         = [{"subnet_id" = "subnet_id_external", public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = "" },{"subnet_id" = subnet_id_external2", public_ip" = true, "private_ip_primary" = "", "private_ip_secondary" = "" }]
      external_securitygroup_ids  = ["securitygroup_id_external","securitygroup_id_external"]
      internal_subnet_ids         = [{"subnet_id" =  "subnet_id_internal", "public_ip"=false, "private_ip_primary" = "" }]
      internal_securitygroup_ids  = ["securitygropu_id_internal"]
      availabilityZones           =  var.availabilityZones
    }


Similarly, you can have N-NIC deployments based on user-provided subnet_ids and securitygroup_ids.
With module count, user can deploy multiple bigip instances in the azure cloud (with the default value of count being one )


.. code-block:: javascript
   :caption: Example 3-NIC Deployment with Static Private IP Allocation
   
    module bigip {
      count                      = var.instance_count
      source                     = "../../"
      prefix                     = format("%s-3nic", var.prefix)
      resource_group_name        = azurerm_resource_group.rg.name
      mgmt_subnet_ids            = [{ "subnet_id" = data.azurerm_subnet.mgmt.id, "public_ip" = true, "private_ip_primary" =  "10.2.1.5"}]
      mgmt_securitygroup_ids     = [module.mgmt-network-security-group.network_security_group_id]
      external_subnet_ids        = [{ "subnet_id" = data.azurerm_subnet.external-public.id, "public_ip" = true, 
                                    "private_ip_primary" = "10.2.2.40","private_ip_secondary" = "10.2.2.50" }]
      external_securitygroup_ids = [module.external-network-security-group-public.network_security_group_id]
      internal_subnet_ids        = [{ "subnet_id" = data.azurerm_subnet.internal.id, "public_ip" = false, "private_ip_primary" = "10.2.3.40"}]
      internal_securitygroup_ids = [module.internal-network-security-group.network_security_group_id]
      availabilityZones          = var.availabilityZones
    }
    
|

InSpec Tool
```````````
The BIG-IP Automation Toolchain InSpec Profile is used for testing the readiness of Automation Tool Chain (ATC) components. After the module deployment, you can use the InSpec tool to verify BIG-IP connectivity with ATC components.

This InSpec profile evaluates the following:

- Basic connectivity to a BIG-IP management endpoint: ``bigip-connectivity``
- Availability of the Declarative Onboarding (DO) service: ``bigip-declarative-onboarding``
- Version reported by the Declarative Onboarding (DO) service: ``bigip-declarative-onboarding-version``
- Availability of the Application Services (AS3) service: ``bigip-application-services``
- Version reported by the Application Services (AS3) service: ``bigip-application-services-version``
- Availability of the Telemetry Streaming (TS) service: ``bigip-telemetry-streaming``
- Version reported by the Telemetry Streaming (TS) service: ``bigip-telemetry-streaming-version``
- Availability of the Cloud Failover Extension (CFE) service: ``bigip-cloud-failover-extension``
- Version reported by the Cloud Failover Extension (CFE) service: ``bigip-cloud-failover-extension-version``

To run InSpec tests, you can either run the inspec exec command or execute runtests.sh in any one of example NIC folders which will run below the inspec command. For example:

::

    inspec exec inspec/bigip-ready --input bigip_address=$BIGIP_MGMT_IP bigip_port=$BIGIP_MGMT_PORT user=$BIGIP_USER password=$BIGIP_PASSWORD do_version=$DO_VERSION as3_version=$AS3_VERSION ts_version=$TS_VERSION fast_version=$FAST_VERSION cfe_version=$CFE_VERSION



|

Required and Optional Input Variables
`````````````````````````````````````
Required variables must be set in the module block when using this module. Optional variables have default values and do not have to be set to use this module. You may set these variables to override their default values.

+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| Parameter                     | Type    | Required | Default               | Description                             |
+===============================+=========+==========+=======================+=========================================+
| prefix                        | String  | Required | N/A                   | This value is inserted in the beginning |
|                               |         |          |                       | of each Azure object.                   |
|                               |         |          |                       | Note: Requires alpha-numeric without    |
|                               |         |          |                       | special characters.                     |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| resource_group_name           | String  | Required | N/A                   | The name of the resource group in which |
|                               |         |          |                       | the resources will be created.          |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| mgmt_subnet_ids               | List of | Required | N/A                   | Map with Subnet-id and public_ip as     |
|                               | maps    |          |                       | keys for the management subnet.         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| mgmt_securitygroup_ids        | List    | Required | N/A                   | securitygroup_ids for the management    |
|                               |         |          |                       | interface.                              |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| availabilityZones             | List    | Required | N/A                   | availabilityZones                       |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| instance_count                | Number  | Required | N/A                   | Number of BIG-IP instances to spin up.  |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| f5_username                   | String  | Optional | bigipuser             | The admin username of the F5 BIG-IP     |
|                               |         |          |                       | that will be deployed                   |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| f5_instance_type              | String  | Optional | Standard_DS3_v2       | Specifies the size of the virtual       |
|                               |         |          |                       | machine.                                |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| f5_image_name                 | String  | Optional | f5-bigip-virtual-edit | 5 SKU (image) you want to deploy.       |
|                               |         |          | ion-200m-best-hourly  | Note: The disk size of the VM will be   |
|                               |         |          |                       | determined based on the option you      |
|                               |         |          |                       | select.                                 |
|                               |         |          |                       | Important: If intending to provision    |
|                               |         |          |                       | multiple modules, ensure the            |
|                               |         |          |                       | appropriate value is selected, such as  |
|                               |         |          |                       | AllTwoBootLocations or                  |
|                               |         |          |                       | AllOneBootLocation.                     |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| f5_version                    | String  | Optional | latest                | It is set to default to use the latest  |
|                               |         |          |                       | software.                               |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| f5_product_name               | String  | Optional | f5-big-ip-best        | Azure BIG-IP VE Offer.                  |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| storage_account_type          | String  | Optional | Standard_LRS          | Defines the type of storage account to  |
|                               |         |          |                       | be created. Valid options are           |
|                               |         |          |                       | Standard_LRS, Standard_ZRS,             |
|                               |         |          |                       | Standard_GRS, Standard_RAGRS, and       |
|                               |         |          |                       | Premium_LRS.                            |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| enable_accelerated_networking | Boolean | Optional | FALSE                 | Enable accelerated networking on        |
|                               |         |          |                       | Network interface.                      |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| enable_ssh_key                | Boolean | Optional | TRUE                  | Enable ssh key authentication in Linux  |
|                               |         |          |                       | Virtual Machine.                        |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| f5_ssh_publickey              | String  | Optional | ~/.ssh/id_rsa.pub     | Path to the public key to be used for   |
|                               |         |          |                       | SSH access to the VM. Only used with    |
|                               |         |          |                       | non-Windows VMs and can be left as-is   |
|                               |         |          |                       | even if using Windows VMs. If you are   |
|                               |         |          |                       | specifying a path to a certification on |
|                               |         |          |                       | a Windows machine to provision a linux  |
|                               |         |          |                       | VM, use the ``/`` in the path instead   |
|                               |         |          |                       | of a backslash. For example:            |
|                               |         |          |                       | ``c:/home/id_rsa.pub``                  |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| doPackageUrl                  | String  | Optional | latest                | URL to download the BIG-IP Declarative  |
|                               |         |          |                       | Onboarding module.                      |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| as3PackageUrl                 | String  | Optional | latest                | URL to download the BIG-IP Application  |
|                               |         |          |                       | Service Extension 3 (AS3) module.       |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| tsPackageUrl                  | String  | Optional | latest                | URL to download the BIG-IP Telemetry    |
|                               |         |          |                       | Streaming module.                       |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| fastPackageUrl                | String  | Optional | latest                | URL to download the BIG-IP FAST module. |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| cfePackageUrlL                | String  | Optional | latest                | URL to download the BIG-IP Cloud        |
|                               |         |          |                       | Failover Extension module.              |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| libs_dir                      | String  | Optional | /config/cloud/azure   | Directory on the BIG-IP to download the |
|                               |         |          | /node_modules         | A&O Toolchain.                          |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| onboard_log                   | String  | Optional | /var/log/startup      | Directory on the BIG-IP to store the    |
|                               |         |          | -script.log           | cloud-init logs.                        |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| azure_secret_rg               | String  | Optional | ``""``                | The name of the resource group in which |
|                               |         |          |                       | the Azure Key Vault exists.             |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| az_key_vault_authentication   | String  | Optional | false                 | Whether to use key vault to pass        |
|                               |         |          |                       | authentications.                        |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| azure_keyvault_name           | String  | Optional | ``""``                | Directory on the BIG-IP to store the    |
|                               |         |          |                       | cloud-init logs.                        |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| azure_keyvault_secret_name    | String  | Optional | ``""``                | The name of the Azure Key Vault secret  |
|                               |         |          |                       | containing the password.                |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| external_subnet_ids           | List of | Optional | [{ "subnet_id" = null | The subnet ID of the virtual network    |
|                               | Maps    |          | , "public_ip" = null, | where the virtual machines will reside. |
|                               |         |          | "private_ip_primary"  |                                         |
|                               |         |          | = "", "private_ip_sec |                                         |
|                               |         |          | ondary" = "" }]       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| internal_subnet_ids           | List of | Optional | [{ "subnet_id" =      | List of maps of subnet IDs of the       |
|                               | Maps    |          | null, "public_ip" =   | virtual network where the virtual       |
|                               |         |          | null,"private_ip_prim | machines will reside.                   |
|                               |         |          | ary" = "" }]          |                                         |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| external_securitygroup_ids    | List    | Optional | ``[]``                | List of Network Security Group IDs for  |
|                               |         |          |                       | external network.                       |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+
| internal_securitygroup_ids    | List    | Optional | ``[]``                | List of Network Security Group IDs for  |
|                               |         |          |                       | internal network.                       |
|                               |         |          |                       |                                         |
|                               |         |          |                       |                                         |
+-------------------------------+---------+----------+-----------------------+-----------------------------------------+

|

Output Variables
````````````````
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Parameter          | Description/Notes                                                                                                                                                                    |
+====================+======================================================================================================================================================================================+
| mgmtPublicIP       | The actual IP address allocated for the resource.                                                                                                                                    |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| mgmtPublicDNS      | FQDN to connect to the first VM provisioned.                                                                                                                                         |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| mgmtPort           | The Mgmt Port.                                                                                                                                                                       |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| f5_username        | BIG-IP username.                                                                                                                                                                     |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| bigip_password     | The BIG-IP Password. If ``dynamic_password`` is selected, then it will be a randomly generated password. If ``azure_keyvault`` is selected, then it will be a key vault secret name. |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| private_addresses  | List of BIG-IP private addresses.                                                                                                                                                    |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| public_addresses   | List of BIG-IP public addresses.                                                                                                                                                     |
+--------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

.. Note:: A local json file will be generated which contains the DO declaration (for 1,2,3 NICs as provided in the examples).

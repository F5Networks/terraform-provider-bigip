F5OS Quick Start Guide
=============================

In this quick start guide, you will deploy a BIG-IP on VELOS hardware using Terraform F5OS resources.

You require the following five files:

- :ref:`Main.tf <f5os_main>` - Contains resource definitions. Use the Single Responsibility principle to minimize the file size.
- :ref:`Variables.tf <f5os_variables>` - Stores all input variables, including all applicable default values.
- :ref:`Inputs.auto.tfvars <f5os_inputs>` - Defines and manages variables values.
- :ref:`Providers.tf <f5os_providers>` - Declares which providers and version to use in the project. Use only on the top-level directory.
- :ref:`Outputs.tf <f5os_outputs>`- Stores exported data items.

.. _f5os_examples:

Example files
----------------------

Use the following file examples to help you automate configurations and interactions for various services
provided by F5 VELOS platform and F5 rSeries appliances.

.. _f5os_variables:

**Example variables file**

.. code-block:: javascript

    variable "tenant_name" {
      description = "the name of the tenant"
    }

    variable "tenant_image_name" {
      description = "the name of the image that is to be used for the tenant"
    }

    variable "tenant_deployment_file" {
      description = "the name of deployment file used for deploying the tenant"
    }

    variable "tenant_type" {
      description = "type of the tenant"
      default = "BIG-IP-Next"
    }

    variable "cpu_cores" {
      description = "the number of vCPUs to be added to the tenant"
      default = 8
    }

    variable "running_state" {
      description = "the desired state for the tenant, should be either 'configured' or 'deployed'"
      default = "deployed"
    }

    variable "mgmt_ip" {
      description = "ip address of the tenant"
    }

    variable "mgmt_gateway" {
      description = "management gateway for the tenant"
    }

    variable "mgmt_prefix" {
      description = "tenant management CIDR prefix"
    }

    variable "timeout" {
      description = "the number of seconds to wait for the tenant to have the desired running state"
      default = 360
    }

    variable "disk_size" {
      description = "minimum virtual disk size required for tenant deployment"
      default = 15
    }

    variable "new_passwd" {
      description = "new password for the bigip next tenant"
      default = ""
    }

.. _f5os_inputs:

**Example inputs file**

.. code-block:: javascript

    cpu_cores 				= 8
    cryptos		 			= "disabled"
    tenant_deployment_file 	= "BIG-IP-Next-0.10.0-4.38.5+0.0.14.yaml"
    id 						= (known after apply)
    tenant_image_name 		= "BIG-IP-Next-0.10.0-4.38.5+0.0.14"
    mgmt_gateway 			= "10.1.10.253"
    mgmt_ip 				= "10.1.10.1"
    mgmt_prefix 			= 24
    tenant_name 			= "testnext"
    running_state 			= "deployed"
    timeout 				= 360
    tenant_type 			= "BIG-IP-Next"
    disk_size 				= 15

.. _f5os_outputs:

**Example outputs file**

.. code-block:: javascript

    output "tenant_status" {
      value = tenant_status
    }

.. _f5os_providers:

**Example providers file**

.. code-block:: javascript

    terraform {
      required_providers {
        f5os = {
          source  = "F5Networks/f5os"
          version = "1.0.0"
        }
      }
    }

    provider "f5os" {
      host 		= "10.10.100.100"
      username 	= "username"
      password 	= "passwd"
    }

.. _f5os_main:

**Example main file**

.. code-block:: javascript

    resource "random_string" "dynamic_password" {
      length      = 16
      min_upper   = 1
      min_lower   = 1
      min_numeric = 1
      special     = false
    }

    resource "f5os_tenant" "bigip_next_tenant" {
      name 				= var.tenant_name
      image_name 		= var.tenant_image_name
      deployment_file 	= var.tenant_deployment_file
      mgmt_ip 			= var.mgmt_ip
      mgmt_prefix 		= var.mgmt_prefix
      mgmt_gateway 		= var.mgmt_gateway
      cpu_cores 		= var.cpu_cores
      running_state 	= var.running_state
      type 				= var.tenant_type
      virtual_disk_size = var.disk_size

      provisioner "local-exec" {
        command = <<EOF
          if [ ${var.running_state} = "deployed" ]
          then
            num_seconds=100
            expected_http=200
            endpoint="https://${var.mgmt_ip}:5443/gui"
            for((i=0; i<$num_seconds; i++)); do
              http_resp=$(curl -k -s -o /dev/null -w "%%{http_code}" $endpoint)
              if [ $http_resp -eq $expected_http ]; then
                curl -k -u admin:admin \
                --header 'Content-Type: application/json' \
                -X PUT https://${var.mgmt_ip}:5443/api/v1/me \
                --data '{"newPassword": "${var.new_passwd != "" ? var.new_passwd : random_string.dynamic_password.result}", "currentPassword": "admin"}'
                exit 0
              fi
              sleep 2
            done
            echo "Could not change the password, maybe the tenant is not yet in the running state"
            exit 1
          fi
        EOF
      }
    }

.. _f5os_deploy:

Deploying BIG-IP on VELOS
---------------------------

1. Use ``Terraform Initialize`` to prepare the working directory so Terraform can run the configuration.

   .. code-block:: console

      $ terraform init
      Initializing the backend...

      Initializing provider plugins...
      - Finding f5networks/f5os versions matching "1.0.0"...

      Terraform has been successfully initialized!

2. Use ``Terraform Plan`` to preview any changes that are required for your infrastructure before applying.

   .. code-block:: console

      $ terraform plan -out bigip-velos

   .. tip::

      If you change modules or change backend configuration for Terraform,
      rerun this command to reinitialize your working directory. If you forget, other
      commands will detect, and then prompt you to rerun ``plan`` (if necessary).

   Terraform uses the selected providers to generate the following example execution plan. Resource actions are indicated with
   the ``+ create`` symbols.

   .. code-block:: console

      # f5os_tenant.bigip_next_tenant will be created
      + resource "f5os_tenant" "bigip_next_tenant" {
              + cpu_cores = 8
              + cryptos = "disabled"
              + deployment_file = "BIG-IP-Next-0.10.0-4.38.5+0.0.14.yaml"
              + id = (known after apply)
              + image_name = "BIG-IP-Next-0.10.0-4.38.5+0.0.14"
              + mgmt_gateway = "10.1.10.253"
              + mgmt_ip = "10.1.10.1"
              + mgmt_prefix = 24
              + name = "testnext"
              + running_state = "deployed"
              + status = (known after apply)
              + timeout = 360
              + type = "BIG-IP-Next"
              + virtual_disk_size = 15
 	          }
      # random_string.dynamic_password will be created
      + resource "random_string" "dynamic_password" {
              + id = (known after apply)
              + length = 16
              + lower = true
              + min_lower = 1
              + min_numeric = 1
              + min_special = 0
              + min_upper = 1
              + number = true
              + numeric = true
              + result = (known after apply)
              + special = false
              + upper = true
            }

      Plan: 2 to add, 0 to change, 0 to destroy.
      Changes to Outputs: ``+ tenant_status = (known after apply)``

   a. Use ``bigip-velos`` to save your plan.

3. Use ``Terraform Apply`` to execute the changes defined by your Terraform configuration and create, update, or destroy resources.
   To perform exactly the previous example actions, run the following command to apply the plan.

   ``terraform apply "bigip-velos"``

   For example:

   .. code-block:: console

        $ terraform apply "bigip-velos"
        random_string.dynamic_password: Creating...
        random_string.dynamic_password: Creation complete after 0s [id=TlROhi9CjZVUPq6E]
        f5os_tenant.bigip_next_tenant: Creating...
        f5os_tenant.bigip_next_tenant: Still creating... [10s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [20s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [30s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [40s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [50s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [1m0s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [1m10s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [1m20s elapsed]
        f5os_tenant.bigip_next_tenant: Provisioning with 'local-exec'...
        f5os_tenant.bigip_next_tenant (local-exec): Executing: ["/bin/sh" "-c" " if [ deployed = \"deployed\" ]\n then\n num_seconds=100\n expected_http=200\n endpoint=\"https://10.1.10.1:5443/gui\"\n for((i=0; i<$num_seconds; i++)); do\n http_resp=$(curl -k -s -o /dev/null -w \"%{http_code}\" $endpoint)\n if [ $http_resp -eq $expected_http ]; then\n curl -k -u admin:admin \\\n --header 'Content-Type: application/json' \\\n -X PUT https://10.1.10.1:5443/api/v1/me \\\n --data '{\"newPassword\": \"F5site02\", \"currentPassword\": \"admin\"}'\n exit 0\n fi\n sleep 2\n done\n echo \"Could not change the password, maybe the tenant is not yet in the running state\"\n exit 1\n fi\n"]
        f5os_tenant.bigip_next_tenant: Still creating... [1m30s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [1m40s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [1m50s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [2m0s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [2m10s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [2m20s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [2m30s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [2m40s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [2m50s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [3m0s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [3m10s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [3m20s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [3m30s elapsed]
        f5os_tenant.bigip_next_tenant: Still creating... [3m40s elapsed]
        f5os_tenant.bigip_next_tenant (local-exec):	  % Total	 % Received	% Xferd	Average Speed	Time	Time	Time	Current
        f5os_tenant.bigip_next_tenant (local-exec):    								Dload  Upload 	Total 	Spent 	Left 	Speed
        f5os_tenant.bigip_next_tenant (local-exec):   0     0	 0     0 	0 	  0 	0 	   0 --:--:-- --:--:-- --:--:--      0
        f5os_tenant.bigip_next_tenant (local-exec):   0     0	 0     0 	0 	  0 	0 	   0 --:--:-- --:--:-- --:--:--      0
        f5os_tenant.bigip_next_tenant (local-exec): 100    55    0     0  100    55 	0 	  36  0:00:01  0:00:01 --:--:--     36
        f5os_tenant.bigip_next_tenant: Creation complete after 3m42s [id=testnext]

        Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

        Outputs:

        tenant_status = "Configured"









What’s Next?

- |f5_terraform_F5OSgithub|
- :doc:`Support <../support>`











.. |f5_terraform_F5OSgithub| raw:: html

   <a href="https://github.com/F5Networks/terraform-provider-F5OS" target="_blank">GitHub</a>
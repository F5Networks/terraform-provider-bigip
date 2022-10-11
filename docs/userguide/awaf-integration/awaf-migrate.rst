.. _awaf-migrate:

Scenario #3: Migrating a WAF Policy from one BIG-IP to another BIG-IP
=====================================================================

.. seealso:: https://github.com/fchmainy/awaf_tf_docs/tree/main/3.migrate

This lab is a variant of the previous one. It takes a manually managed Advanced WAF Policy from an existing BIG-IP and migrates it to a different BIG-IP through Terraform resources.

You can meet this scenario in multiple use-cases:

- Migrating from a BIG-IP to another (platform refresh).
- Re-Hosting (aka Lift&Shift) in a Cloud migration project.
- Back-and-Forth importing/exporting WAF Policies between environments (dev/test/QA/Production)

The goal is to leverage the previous import scenario in order to carry and ingest the WAF Policy from one BIG-IP to another while keeping its state through Terraform.

The WAF Policy and its children objects (parameters, URLs, attack signatures, exceptions, etc.) can be tightly coupled to a BIG-IP and/or can be shared across multiple policies, depending on the use case.

Pre-requisites
--------------
On the BIG-IP:

- BIG-IP version 16.1 or newer
- Advanced WAF Provisioned
- Credentials with REST API access

On Terraform:

- Using F5 BIG-IP provider version 1.15.0 or newer
- Using Hashicorp versions following :ref:`versions`

Migrating a Policy
------------------
Create 4 files:

- variables.tf
- inputs.auto.tfvars
- main.tf
- outputs.tf

.. code-block:: json
   :caption: variables.tf
   :linenos:

   variable previous_bigip {}
   variable new_bigip {}
   variable username {}
   variable password {}

|

.. code-block:: json
   :caption: inputs.auto.tfvars
   :linenos:

   previous_bigip = "10.1.1.8:443"
   new_bigip = "10.1.1.9:443"
   username = "admin"
   password = "whatIsYourBigIPPassword?"

|

.. code-block:: json
   :caption: main.tf
   :linenos:

   terraform {
     required_providers {
       bigip = {
         source = "F5Networks/bigip"
         version = "1.15"
       }
     }
   }
   provider "bigip" {
     alias    = "old"
     address  = var.previous_bigip
     username = var.username
     password = var.password
   }
   provider "bigip" {
     alias    = "new"
     address  = var.new_bigip
     username = var.username
     password = var.password
   }


   resource "bigip_waf_policy" "current" {
     provider	       = bigip.old
     partition            = "Common"
     name                 = "scenario3"
     template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
   }

.. Note:: The template name can be set to anything. When it is imported, the value is overwritten.

|

.. code-block:: json
   :caption: outputs.tf
   :linenos:

   output "policyId" {
   	value	= bigip_waf_policy.current.policy_id
   }

   output "policyJSON" {
           value   = bigip_waf_policy.current.policy_export_json
   }

|

Here we defined two BIG-IPs: "old" and "new". The "old" BIG-IP has the existing Advanced WAF Policies, the "new" is our target.

Similar to :ref:`awaf-import`, you need the Advanced WAF Policy ID to make the initial import:

- Check on the iControl REST API Endpoint: ``/mgmt/tm/asm/policies?$filter=name+eq+scenario3&$select=id``
- Get a script example in the ``lab/scripts/`` folder
- Run the following piece of code in the Go PlayGround


::

   package main

   import (
       "crypto/md5"
       b64 "encoding/base64"
       "fmt"
       "strings"
   )

   func Hasher(policyName string) string {
       hasher := md5.New()
       hasher.Write([]byte(policyName))
       encodedString := b64.StdEncoding.EncodeToString(hasher.Sum(nil))

       return strings.TrimRight(encodedString, "=")
   }

   func main() {
       var partition string = "Common"
       var policyName string = "scenario3"

       fullName := "/" + partition + "/" + policyName
       policyId := Hasher(fullName)

       r := strings.NewReplacer("/", "_", "-", "_", "+", "-")
       fmt.Println("Policy Id: ", r.Replace(policyId))
   }


Run the following commands to:

1. Initialize the Terraform Project.
2. Import the current WAF policy from the "old" BIG-IP into your state.
3. Create the Advanced WAF Policy resource for the "BIG-IP" pointing to the imported state.
4. Configure the lifecycle of our WAF Policy.

:: 

   foo@bar:~$ terraform init
   Initializing the backend...

   Initializing provider plugins...
   [...]
   Terraform has been successfully initialized!

   foo@bar:~$ terraform import bigip_waf_policy.current YiEQ4l1Fw1U9UnB2-mTKWA
   bigip_waf_policy.this: Importing from ID "YiEQ4l1Fw1U9UnB2-mTKWA"...
   bigip_waf_policy.this: Import prepared!
     Prepared bigip_waf_policy for import
   bigip_waf_policy.this: Refreshing state... [id=YiEQ4l1Fw1U9UnB2-mTKWA]

   Import successful!

   The resources that were imported are shown above. These resources are now in
   your Terraform state and will henceforth be managed by Terraform.


Update your **terraform main.tf** file with the ouputs of the following two commands:

::

   foo@bar:~$ terraform show -json | jq '.values.root_module.resources[].values.policy_export_json | fromjson' > currentWAFPolicy.json

   foo@bar:~$ terraform show -no-color
   # bigip_waf_policy.this:
   resource "bigip_waf_policy" "this" {
       application_language = "utf-8"
       id                   = "YiEQ4l1Fw1U9UnB2-mTKWA"
       name                 = "/Common/scenario3"
       policy_export_json   = jsonencode(
           {
               [...]
           }
       )
       policy_id            = "YiEQ4l1Fw1U9UnB2-mTKWA"
       template_name        = "POLICY_TEMPLATE_COMPREHENSIVE"
       type                 = "security"
   }


This a migration use case so you do not need the current WAF Policy from the existing BIG-IP. Using the collected data from the Terraform import, you can now update your **main.tf** file:

::

   resource "bigip_waf_policy" "migrated" {
       provider	           = bigip.new
       application_language = "utf-8"
       partition            = "Common"
       name                 = "scenario3"
       policy_id            = "YiEQ4l1Fw1U9UnB2-mTKWA"
       template_name        = "POLICY_TEMPLATE_COMPREHENSIVE"
       type                 = "security"
       policy_import_json   = file("${path.module}/currentWAFPolicy.json")
   }



Note that F5 replaced the "policy_export_json" argument with "policy_import_json" pointing to the imported WAF Policy JSON file.

Finally, you can plan and apply your new project.

:: 

   foo@bar:~$ terraform plan -out scenario3
   bigip_waf_policy.migrated: Refreshing state... [id=YiEQ4l1Fw1U9UnB2-mTKWA]
   
   Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
     ~ update in-place
   [...]
   ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
   
   Saved the plan to: scenario3
   
   To perform exactly these actions, run the following command to apply:
       terraform apply "scenario3"
   
   foo@bar:~$ terraform apply "scenario3"
   bigip_waf_policy.this: Modifying... [id=YiEQ4l1Fw1U9UnB2-mTKWA]
   bigip_waf_policy.this: Still modifying... [id=EdchwjSqo9cFtYP-iWUJmw, 10s elapsed]
   bigip_waf_policy.this: Modifications complete after 16s [id=EdchwjSqo9cFtYP-iWUJmw]
   
   Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
   
   Outputs:
   
   policyId = "EdchwjSqo9cFtYP-iWUJmw"
   policyJSON = "{[...]}"

|

Policy lifecycle management
---------------------------
You can manage your WAF Policy as shown in the previous lab. You can check your WAF Policy on your BIG-IP after each terraform apply.

Defining parameters
```````````````````
Create a **parameters.tf** file:

:: 

   data "bigip_waf_entity_parameter" "P1" {
     name            = "Parameter1"
     type            = "explicit"
     data_type       = "alpha-numeric"
     perform_staging = true
     signature_overrides_disable = [200001494, 200001472]
   }


Add references to these parameters in the ``bigip_waf_policy`` TF resource in the **main.tf** file:

:: 

   resource "bigip_waf_policy" "migrated" {
     [...]
     parameters           = [data.bigip_waf_entity_parameter.P1.json]
   }

:: 

   foo@bar:~$ terraform plan -out scenario3
   foo@bar:~$ terraform apply "scenario3"

|

Defining URLs
`````````````
Create a **urls.tf** file:

::

   data "bigip_waf_entity_url" "U1" {
     name		              = "/URL1"
     description                 = "this is a test for URL1"
     type                        = "explicit"
     protocol                    = "http"
     perform_staging             = true
     signature_overrides_disable = [12345678, 87654321]
     method_overrides {
       allow  = false
       method = "BCOPY"
     }
     method_overrides {
       allow  = true
       method = "BDELETE"
     }
   }
   
   data "bigip_waf_entity_url" "U2" {
     name                        = "/URL2"
   }



Add references to this URL in the ``bigip_waf_policy`` TF resource in the **main.tf** file:

:: 

   resource "bigip_waf_policy" "migrated" {
     [...]
     urls                 = [data.bigip_waf_entity_url.U1.json, data.bigip_waf_entity_url.U2.json]
   }


Run it:

:: 

   foo@bar:~$ terraform plan -out scenario3
   foo@bar:~$ terraform apply "scenario3"

|

Defining Attack Signatures
``````````````````````````
Create a **signatures.tf** file:

:: 

   data "bigip_waf_signatures" "S1" {
     provider         = bigip.new
     signature_id     = 200104004
     description      = "Java Code Execution"
     enabled          = true
     perform_staging  = true
   }
   
   data "bigip_waf_signatures" "S2" {
     provider         = bigip.new
     signature_id     = 200104005
     enabled          = false
   }

Add references to this URL in the ``bigip_waf_policy`` TF resource in the **main.tf** file:

:: 

   resource "bigip_waf_policy" "migrated" {
     [...]
     signatures       = [data.bigip_waf_signatures.S1.json, data.bigip_waf_signatures.S2.json]
   }

Run it:

:: 

   foo@bar:~$ terraform plan -out scenario3
   foo@bar:~$ terraform apply "scenario3"
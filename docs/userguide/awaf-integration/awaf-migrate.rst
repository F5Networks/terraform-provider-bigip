.. _awaf-integration:

**`Scenario #3: Migrating a WAF Policy from a BIG-IP to another BIG-IP <https://github.com/fchmainy/awaf_tf_docs/tree/main/3.migrate>`_**
 
This lab is a variant of the previous one. It takes a manually managed A.WAF Policy from an existing BIG-IP and migrate it to a different BIG-IP through Terraform resources.

Goals
You can meet this scenario in multiple use-cases:

migrating from a BIG-IP to another (platform refresh)
Re-Hosting (aka Lift&Shift) in a Cloud migration project
Back-and-Forth importing / exporting WAF Policies between environments (dev / test / QA / Production)
The goal is to leverage the previous import scenario in order to carry and ingest the WAF Policy from one BIG-IP to another while keeping its state through Terraform.

The WAF Policy and its children objects (parameters, urls, attack signatures, exceptions...) can be tightly coupled to a BIG-IP AND / OR can be shared across multiple policies depending on the use case.

Pre-requisites
--------------
on the BIG-IP:

 version 16.1 minimal
 A.WAF Provisioned
 credentials with REST API access

on Terraform:

 use of F5 BIG-IP provider version 1.15.0 minimal
 use of Hashicorp version followinf Link


Policy Migration
----------------
Create 4 files:

main.tf
variables.tf
inputs.auto.tfvars
outputs.tf

.. code-block:: json
   :caption: variables.tf
   :linenos:

   variable previous_bigip {}
   variable new_bigip {}
   variable username {}
   variable password {}


.. code-block:: json
   :caption: inputs.auto.tfvars
   :linenos:

   previous_bigip = "10.1.1.8:443"
   new_bigip = "10.1.1.9:443"
   username = "admin"
   password = "whatIsYourBigIPPassword?"

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

Note: the template name can be set to anything. When it is imported, we will overwrite the value

.. code-block:: json
   :caption: outputs.tf
   :linenos:

   output "policyId" {
   	value	= bigip_waf_policy.current.policy_id
   }

   output "policyJSON" {
           value   = bigip_waf_policy.current.policy_export_json
   }


Here we defined two Big-IPs: "old" and "new". The "old" BIG-IP has the existing A.WAF Policies, the "new" is our target.

Same as scenario #2 we need the A.WAF Policy ID to make the initial import:

- check on the iControl REST API Endpoint: /mgmt/tm/asm/policies?$filter=name+eq+scenario3&$select=id
- get a script example in the lab/scripts/ folder
- run the following piece of code in the Go PlayGround


.. code-block:: json
   :caption: 
   :linenos:

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


Now, run the following commands, so we can:

1. Initialize the terraform project
2. Import the current WAF policy from the "old" BIG-IP into our state
3. Create the A.WAF Policy resource for the "BIG-IP" pointing to the imported state
4. Configure the lifecycle of our WAF Policy

.. code-block:: json
   :caption: 
   :linenos:

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


Now update your terraform main.tf file with the ouputs of the following two commands:

.. code-block:: json
   :caption: 
   :linenos:

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


This a migration use case so we don't need anymore the current WAF Policy from the existing BIG-IP. So, using the collected data from the terraform import, we are now updating our main.tf file: If you want to keep the policy on both BIG-IPs, [please get there](UPDATE LINK!!!!)

.. code-block:: json
   :caption: 
   :linenos:

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



Note: You can note that we replaced the "policy_export_json" argument with "policy_import_json" pointing to the imported WAF Policy JSON file.

Finally, we can plan & apply our new project.

.. code-block:: json
   :caption: 
   :linenos:

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



Policy lifecycle management
---------------------------
Now you can manage your WAF Policy as we did in the previous lab

You can check your WAF Policy on your BIG-IP after each terraform apply.

Defining parameters
```````````````````
Create a parameters.tf file:

.. code-block:: json
   :caption: 
   :linenos:

   data "bigip_waf_entity_parameter" "P1" {
     name            = "Parameter1"
     type            = "explicit"
     data_type       = "alpha-numeric"
     perform_staging = true
     signature_overrides_disable = [200001494, 200001472]
   }


And add references to these parameters in the "bigip_waf_policy" TF resource in the main.tf file:

.. code-block:: json
   :caption: 
   :linenos:

   resource "bigip_waf_policy" "migrated" {
     [...]
     parameters           = [data.bigip_waf_entity_parameter.P1.json]
   }
   foo@bar:~$ terraform plan -out scenario3
   foo@bar:~$ terraform apply "scenario3"


Defining URLs
`````````````
Create a urls.tf file:

.. code-block:: json
   :caption: 
   :linenos:

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



And add references to this URL in the "bigip_waf_policy" TF resource in the main.tf file:

.. code-block:: json
   :caption: 
   :linenos:

   resource "bigip_waf_policy" "migrated" {
     [...]
     urls                 = [data.bigip_waf_entity_url.U1.json, data.bigip_waf_entity_url.U2.json]
   }


and run it:

.. code-block:: json
   :caption: 
   :linenos:

   foo@bar:~$ terraform plan -out scenario3
   foo@bar:~$ terraform apply "scenario3"


Defining Attack Signatures
``````````````````````````
Create a signatures.tf file:

.. code-block:: json
   :caption: 
   :linenos:

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

And add references to this URL in the "bigip_waf_policy" TF resource in the main.tf file:

.. code-block:: json
   :caption: 
   :linenos:

   resource "bigip_waf_policy" "migrated" {
     [...]
     signatures       = [data.bigip_waf_signatures.S1.json, data.bigip_waf_signatures.S2.json]
   }

and run it:

.. code-block:: json
   :caption: 
   :linenos:

   foo@bar:~$ terraform plan -out scenario3
   foo@bar:~$ terraform apply "scenario3"
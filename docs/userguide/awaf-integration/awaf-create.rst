.. _awaf-create:

Scenario #1: Creating a WAF Policy
==================================

This page shows you how to create a new Advanced WAF Policy from scratch and manage some entities additions.

Pre-requisites
--------------

- BIG-IP version 16.1 or newer
- Credentials with REST API access

On Terraform:

- Using F5 BIG-IP provider version 1.15.0 or newer
- Using Hashicorp versions following :ref:`versions`


Creating a Policy
-----------------

Create 3 files:

- main.tf
- variables.tf
- inputs.auto.tfvars

.. code-block:: json
   :caption: variables.tf
   :linenos:

   variable bigip {}
   variable username {}
   variable password {}

|

.. code-block:: json
   :caption: inputs.auto.tfvars
   :linenos:

   bigip = "10.1.1.9:443"
   username = "admin"
   password = "yYyYyYy"

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
     address  = var.bigip
     username = var.username
     password = var.password
   }
   
   resource "bigip_waf_policy" "this" {
     name                 = "scenario1"
     partition             = "Common"
     template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
     application_language = "utf-8"
     enforcement_mode     = "blocking"
     server_technologies  = ["Apache Tomcat", "MySQL", "Unix/Linux"]
   }

|   

Now that you have your terraform project all set up, you can run it:


:: 

   foo@bar:~$ terraform init

   Initializing the backend...
   
   Initializing provider plugins...
   - Finding f5networks/bigip versions matching "1.14.0"...
   - Installing f5networks/bigip v1.14.0...
   - Installed f5networks/bigip v1.14.0 (signed by a HashiCorp partner, key ID 0F284A6527D73A63)
   
   Partner and community providers are signed by their developers.
   If you'd like to know more about provider signing, you can read about it here:
   https://www.terraform.io/docs/cli/plugins/signing.html
   
   Terraform has created a lock file .terraform.lock.hcl to record the provider
   selections it made above. Include this file in your version control repository
   so that Terraform can guarantee to make the same selections by default when
   you run "terraform init" in the future.
   
   Terraform has been successfully initialized!
   
   You may now begin working with Terraform. Try running "terraform plan" to see
   any changes that are required for your infrastructure. All Terraform commands
   should now work.
   
   If you ever set or change modules or backend configuration for Terraform,
   rerun this command to reinitialize your working directory. If you forget, other
   commands will detect it and remind you to do so if necessary.
   
   foo@bar:~$ terraform plan -out scenario1
   
   Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
     + create
   
   Terraform will perform the following actions:
   
     # bigip_waf_policy.this will be created
     + resource "bigip_waf_policy" "this" {
         + application_language = "utf-8"
         + case_insensitive     = false
         + enable_passivemode   = false
         + enforcement_mode     = "blocking"
         + id                   = (known after apply)
         + name                 = "/Common/scenario1"
         + policy_export_json   = (known after apply)
         + policy_id            = (known after apply)
         + server_technologies  = [
             + "MySQL",
             + "Unix/Linux",
             + "MongoDB",
           ]
         + template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
         + type                 = "security"
       }

   Plan: 1 to add, 0 to change, 0 to destroy.

   ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

   Saved the plan to: scenario1

   To perform exactly these actions, run the following command to apply:
       terraform apply "scenario1"

   foo@bar:~$ terraform apply "scenario1"
   bigip_waf_policy.this: Creating...
   bigip_waf_policy.this: Still creating... [10s elapsed]
   bigip_waf_policy.this: Creation complete after 17s [id=41UMLL7yDtzoa0000Wimzw]

   Apply complete! Resources: 1 added, 0 changed, 0 destroyed.


Your WAF Policy might evolve over time, so you may want to add entities, manage attack signature exceptions, etc.

Policy Lifecycle Management
---------------------------

Server Technologies
```````````````````
You now want to add a **MongoDB** server technology into your WAF Policy. The allowed values for server technologies are listed in the `Declarative WAF API documentation <https://clouddocs.f5.com/products/waf-declarative-policy/declarative_policy_v16_1.html#server-technologies>`_.

Edit the main.tf file:

.. code-block:: json
   :caption: main.tf
   :linenos:

   resource "bigip_waf_policy" "this" {
     name                 = "localS1"
     partition	           = "Common"
     template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
     application_language = "utf-8"
     enforcement_mode     = "blocking"
     server_technologies  = ["Apache Tomcat", "MySQL", "Unix/Linux", "MongoDB"]
   }

|

Parameters management
`````````````````````
Create a parameters.tf file:

.. code-block:: json
   :caption: parameters.tf
   :linenos:

   data "bigip_waf_entity_parameter" "P1" {
     name            = "Parameter1"
     type            = "explicit"
     data_type       = "alpha-numeric"
     perform_staging = true
   }
   
   data "bigip_waf_entity_parameter" "P2" {
     name            = "Parameter2"
     type            = "wildcard"
     data_type       = "alpha-numeric"
     perform_staging = false
     signature_overrides_disable = [200001494, 200001472]
   }
   
   data "bigip_waf_entity_parameter" "P3" {
     name            = "Parameter3"
     type            = "explicit"
     data_type       = "alpha-numeric"
     is_header	  = true
     sensitive_parameter = true
     perform_staging = true
   }

|

And add references to these parameters in the "bigip_waf_policy" TF resource in the main.tf file:

.. code-block:: json
   :caption: main.tf
   :linenos:
   
   resource "bigip_waf_policy" "this" {
     name                 = "scenario1"
     partition            = "Common"
     template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
     application_language = "utf-8"
     enforcement_mode     = "blocking"
     server_technologies  = ["Apache Tomcat", "MySQL", "Unix/Linux", "MongoDB"]
     parameters           = [data.bigip_waf_entity_parameter.P1.json, data.bigip_waf_entity_parameter.P2.json, data.bigip_waf_entity_parameter.P3.json]
   }

|

Run it:

::

   foo@bar:~$ terraform plan -out scenario1
   [...]
   
   Plan: 0 to add, 1 to change, 0 to destroy.
   
   ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
   
   Saved the plan to: scenario1
   
   To perform exactly these actions, run the following command to apply:
       terraform apply "scenario1"
   
   foo@bar:~$ terraform apply "scenario1"
   bigip_waf_policy.this: Modifying... [id=41UMLL7yDtzoa0000Wimzw]
   bigip_waf_policy.this: Still modifying... [id=41UMLL7yDtzoa0000Wimzw, 10s elapsed]
   bigip_waf_policy.this: Modifications complete after 17s [id=41UMLL7yDtzoa0000Wimzw]
   
   Apply complete! Resources: 0 added, 1 changed, 0 destroyed.


Signatures Management
`````````````````````
Create a separate signature definition file with 3 signatures:

- S1 enables and perform staging on the **200010293** attack signature.
- S2 disables the **200009024** attack signature.
- S3 enables and enforce the **200014009** attack signature.

Create a signatures.tf file:

.. code-block:: json
   :caption: signatures.tf
   :linenos:

   data "bigip_waf_signatures" "S1" {
     signature_id     = 200010293
     description      = "Java Code Execution"
     enabled          = true
     perform_staging  = true
   }
   
   data "bigip_waf_signatures" "S2" {
     signature_id      = 200009024
     enabled          = false
   }
   
   data "bigip_waf_signatures" "S3" {
     signature_id      = 200014009
     description      = "src http: (Header)"
     enabled          = true
     perform_staging  = false
   }

|

Add references to these attack signatures in the "bigip_waf_policy" TF resource in the main.tf file:

::

   resource "bigip_waf_policy" "this" {
     partition            = "Common"
     name                 = "scenario1"
     template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
     application_language = "utf-8"
     enforcement_mode     = "blocking"
     server_technologies  = ["Apache Tomcat", "MySQL", "Unix/Linux", "MongoDB"]
     parameters           = [data.bigip_waf_entity_parameter.P1.json, data.bigip_waf_entity_parameter.P2.json, data.bigip_waf_entity_parameter.P3.json]
     signatures           = [data.bigip_waf_signatures.S1.json, data.bigip_waf_signatures.S2.json, data.bigip_waf_signatures.S3.json]
   }


Run it:

::

   foo@bar:~$ terraform plan -out scenario1
   [...]
     # bigip_waf_policy.this will be updated in-place
     ~ resource "bigip_waf_policy" "this" {
           id                   = "tCwXEedPDS-S35Bl4TSU5w"
           name                 = "localS1"
         + signatures           = [
             + jsonencode(
                   {
                     + enabled        = true
                     + performStaging = true
                     + signatureId    = 200010293
                   }
               ),
             + jsonencode(
                   {
                     + performStaging = false
                     + signatureId    = 200009024
                   }
               ),
             + jsonencode(
                   {
                     + enabled        = true
                     + performStaging = false
                     + signatureId    = 200014009
                   }
               ),
           ]
           # (11 unchanged attributes hidden)
       }

   Plan: 0 to add, 1 to change, 0 to destroy.
   
   ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
   
   Saved the plan to: scenario1
   
   To perform exactly these actions, run the following command to apply:
       terraform apply "scenario1"
   
   foo@bar:~$ terraform apply "scenario1"
   bigip_waf_policy.this: Modifying... [id=41UMLL7yDtzoa0000Wimzw]
   bigip_waf_policy.this: Still modifying... [id=41UMLL7yDtzoa0000Wimzw, 10s elapsed]
   bigip_waf_policy.this: Modifications complete after 17s [id=41UMLL7yDtzoa0000Wimzw]
   
   Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
   
|

Use the following command to check the details on a specific Attack signature:

:: 

   $ terraform show -json | jq '.values.root_module.resources[] | select(.name == "S3")'

::

   {
     "address": "data.bigip_waf_signatures.S3",
     "mode": "data",
     "type": "bigip_waf_signatures",
     "name": "S3",
     "provider_name": "terraform.local/local/bigip",
     "schema_version": 0,
     "values": {
       "accuracy": "medium",
       "description": "Summary:\nThis event is generated when an attempt is made to abuse a web server functionality. This is a general detection signature (i.e. it is not specific to any web application).\n\nImpact:\nVary from information gathering to web server compromise.\n\nDetailed Information:\nAbuse of Functionality is an attack technique that uses a web site's own features and functionality to consume, defraud, or circumvents access controls mechanisms\n\nAffected Systems:\nAll systems.\n\nAttack Scenarios:\nThere are many possible.\n\nEase Of Attack:\nSimple to medium.\n\nFalse Positives:\nSome applications may accept valid input which matches these signatures.\n\nFalse Negatives:\nNone known.\n\nCorrective Action:\nEnsure the system is using an up to date version of the software and has had all vendor supplied patches applied. Utilize \"Positive Security Model\" by accepting only known types of input in web application.\n\nAdditional References:\nhttp://www.webappsec.org/projects/threat/classes/abuse_of_functionality.shtml\n\n",
       "enabled": true,
       "id": "200014009",
       "json": "{\"signatureId\":200014009,\"performStaging\":false,\"enabled\":true}",
       "name": "Unix \"cmd\" parameter execution attempt",
       "perform_staging": false,
       "risk": "high",
       "signature_id": 200014009,
       "system_signature_id": "GTK2ItJX6pnKHXBqiwtlxQ",
       "tag": null,
       "type": "request"
     },
     "sensitive_values": {}
   }


.. NOTE:: If you have multiple entities to manage, the entity lists in the bigip_waf_policy can be difficult to use. In that case, F5 recommends using Terraform HCL maps as presented in `lab 4 <https://github.com/fchmainy/awaf_tf_docs/blob/main/4.multiple/README.md#enforcing-attack-signatures-on-the-qa-environment>`_.

.. seealso:: `F5 BIG-IP Terraform Provider official documentation <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs>`_

Creating a Policy via OpenAPI file
----------------------------------

Create 3 files:

- main.tf
- variables.tf
- inputs.tfvars

.. code-block:: json
   :caption: variables.tf
   :linenos:

   variable bigip {}
   variable username {}
   variable password {}

|

.. code-block:: json
   :caption: inputs.auto.tfvars
   :linenos:

   bigip = "10.1.1.9:443"
   username = "admin"
   password = "yYyYyYy"

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
     address  = var.bigip
     username = var.username
     password = var.password
   }
 
   resource "bigip_waf_policy" "this" {
     partition                 = "Common"
     name                      = "scenario1.swagger"
     template_name             = "POLICY_TEMPLATE_API_SECURITY"
     application_language      = "utf-8"
     enforcement_mode          = "blocking"
     server_technologies       = ["MySQL", "Unix/Linux", "MongoDB"]
     open_api_files            = ["https://api.swaggerhub.com/apis/F5EMEASSA/API-Sentence/3.0.1"]
     parameters                = [data.bigip_waf_entity_parameter.P1.json, data.bigip_waf_entity_parameter.P2.json, data.bigip_waf_entity_parameter.P3.json]
     signatures                = [data.bigip_waf_signatures.S1.json, data.bigip_waf_signatures.S2.json]


.. seealso:: `How to create an OpenAPI security policy using a Swagger file <https://support.f5.com/csp/article/K07241201>`_

Run it:

::

   foo@bar:~$ terraform plan -out scenario1.swagger
   [...]
     # bigip_waf_policy.this must be replaced
   -/+ resource "bigip_waf_policy" "this" {
         ~ id                   = "41UMLL7yDtzoa0000Wimzw" -> (known after apply)
         ~ name                 = "scenario1" -> "scenario1-2.swagger" # forces replacement
         + open_api_files       = [
             + "https://api.swaggerhub.com/apis/F5EMEASSA/API-Sentence/3.0.1",
           ]
         ~ policy_export_json   = jsonencode(
               {
                 - applicationLanguage = "utf-8"
                 - description         = "FCH Testing WAF Policy from RDP Template"
                 - enforcementMode     = "blocking"
                 - fullPath            = "/Common/scenario1"
                 - name                = "scenario1"
                 - parameters          = [
                     - {
                         - allowEmptyValue            = true
                         - allowRepeatedParameterName = true
                         - attackSignaturesCheck      = true
                         - dataType                   = "alpha-numeric"
                         - isHeader                   = true
                         - level                      = "global"
                         - name                       = "parameter3"
                         - performStaging             = true
                         - sensitiveParameter         = true
                         - type                       = "explicit"
                         - valueType                  = "user-input"
                       },
                     - {
                         - allowEmptyValue            = true
                         - allowRepeatedParameterName = true
                         - attackSignaturesCheck      = true
                         - dataType                   = "alpha-numeric"
                         - level                      = "global"
                         - name                       = "Parameter2"
                         - parameterLocation          = "any"
                         - signatureOverrides         = [
                             - {
                                 - enabled     = false
                                 - name        = "\"style :expression (\" (Parameter)(1)"
                                 - signatureId = 200001494
                               },
                           ]
                         - type                       = "wildcard"
                         - valueType                  = "user-input"
                       },
                     - {
                         - allowEmptyValue            = true
                         - allowRepeatedParameterName = true
                         - attackSignaturesCheck      = true
                         - dataType                   = "alpha-numeric"
                         - level                      = "global"
                         - name                       = "Parameter1"
                         - parameterLocation          = "any"
                         - performStaging             = true
                         - type                       = "explicit"
                         - valueType                  = "user-input"
                       },
                     - {
                         - allowEmptyValue            = true
                         - allowRepeatedParameterName = true
                         - attackSignaturesCheck      = true
                         - level                      = "global"
                         - name                       = "*"
                         - parameterLocation          = "any"
                         - type                       = "wildcard"
                         - valueType                  = "auto-detect"
                       },
                     - {
                         - allowEmptyValue   = true
                         - level             = "global"
                         - name              = "__VIEWSTATE"
                         - parameterLocation = "any"
                         - type              = "explicit"
                         - valueType         = "ignore"
                       },
                   ]
                 - server-technologies = [
                     - {
                         - serverTechnologyName = "MongoDB"
                       },
                     - {
                         - serverTechnologyName = "MySQL"
                       },
                     - {
                         - serverTechnologyName = "Apache Tomcat"
                       },
                     - {
                         - serverTechnologyName = "Unix/Linux"
                       },
                   ]
                 - signature-sets      = [
                     - {
                         - alarm        = true
                         - block        = true
                         - learn        = true
                         - name         = "Apache Tomcat Signatures (High/Medium Accuracy)"
                         - signatureSet = {
                             - filter  = {
                                 - accuracyFilter    = "ge"
                                 - accuracyValue     = "medium"
                                 - hasCve            = "all"
                                 - lastUpdatedFilter = "all"
                                 - riskFilter        = "all"
                                 - riskValue         = "all"
                                 - signatureType     = "all"
                                 - tagFilter         = "all"
                                 - userDefinedFilter = "all"
                               }
                             - systems = [
                                 - {
                                     - name = "Apache Tomcat"
                                   },
                               ]
                             - type    = "filter-based"
                           }
                       },
                     - {
                         - alarm        = true
                         - block        = true
                         - learn        = true
                         - name         = "MongoDB Signatures (High/Medium Accuracy)"
                         - signatureSet = {
                             - filter  = {
                                 - accuracyFilter    = "ge"
                                 - accuracyValue     = "medium"
                                 - hasCve            = "all"
                                 - lastUpdatedFilter = "all"
                                 - riskFilter        = "all"
                                 - riskValue         = "all"
                                 - signatureType     = "all"
                                 - tagFilter         = "all"
                                 - userDefinedFilter = "all"
                               }
                             - systems = [
                                 - {
                                     - name = "MongoDB"
                                   },
                               ]
                             - type    = "filter-based"
                           }
                       },
                     - {
                         - alarm        = true
                         - block        = true
                         - learn        = true
                         - name         = "Unix/Linux Signatures (High/Medium Accuracy)"
                         - signatureSet = {
                             - filter  = {
                                 - accuracyFilter    = "ge"
                                 - accuracyValue     = "medium"
                                 - hasCve            = "all"
                                 - lastUpdatedFilter = "all"
                                 - riskFilter        = "all"
                                 - riskValue         = "all"
                                 - signatureType     = "all"
                                 - tagFilter         = "all"
                                 - userDefinedFilter = "all"
                               }
                             - systems = [
                                 - {
                                     - name = "Unix/Linux"
                                   },
                               ]
                             - type    = "filter-based"
                           }
                       },
                     - {
                         - alarm        = true
                         - block        = true
                         - learn        = true
                         - name         = "MySQL Signatures (High/Medium Accuracy)"
                         - signatureSet = {
                             - filter  = {
                                 - accuracyFilter    = "ge"
                                 - accuracyValue     = "medium"
                                 - hasCve            = "all"
                                 - lastUpdatedFilter = "all"
                                 - riskFilter        = "all"
                                 - riskValue         = "all"
                                 - signatureType     = "all"
                                 - tagFilter         = "all"
                                 - userDefinedFilter = "all"
                               }
                             - systems = [
                                 - {
                                     - name = "MySQL"
                                   },
                               ]
                             - type    = "filter-based"
                           }
                       },
                     - {
                         - alarm        = true
                         - block        = true
                         - learn        = true
                         - name         = "Generic Detection Signatures (High/Medium Accuracy)"
                         - signatureSet = {
                             - filter = {}
                           }
                       },
                   ]
                 - signature-settings  = {
                     - signatureStaging = true
                   }
                 - signatures          = [
                     - {
                         - enabled        = true
                         - performStaging = false
                         - signatureId    = "200014009"
                       },
                     - {
                         - enabled        = true
                         - performStaging = false
                         - signatureId    = "200009024"
                       },
                   ]
                 - template            = {
                     - name = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
                   }
                 - type                = "security"
                 - urls                = [
                     - {
                         - attackSignaturesCheck = true
                         - isAllowed             = true
                         - method                = "*"
                         - name                  = "*"
                         - protocol              = "http"
                         - type                  = "wildcard"
                       },
                     - {
                         - attackSignaturesCheck = true
                         - isAllowed             = true
                         - method                = "*"
                         - name                  = "*"
                         - protocol              = "https"
                         - type                  = "wildcard"
                       },
                   ]
               }
           ) -> (known after apply)
         ~ policy_id            = "41UMLL7yDtzoa0000Wimzw" -> (known after apply)
         ~ server_technologies  = [
             - "Apache Tomcat",
               "MySQL",
               # (2 unchanged elements hidden)
           ]
         ~ signatures           = [
               # (1 unchanged element hidden)
               jsonencode(
                   {
                       performStaging = false
                       signatureId    = 200009024
                   }
               ),
             - jsonencode(
                   {
                     - enabled        = true
                     - performStaging = false
                     - signatureId    = 200014009
                   }
               ),
           ]
         ~ template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT" -> "POLICY_TEMPLATE_API_SECURITY" # forces replacement
           # (7 unchanged attributes hidden)
       }

   Plan: 1 to add, 0 to change, 1 to destroy.

   ───────────────────────────────────────────────────────────────────────────────

   Saved the plan to: scenario1.swagger

   To perform exactly these actions, run the following command to apply:
       terraform apply "scenario1.swagger"
   
   foo@bar:~$ terraform apply "scenario1"
   bigip_waf_policy.this: Modifying... [id=41UMLL7yDtzoa0000Wimzw]
   bigip_waf_policy.this: Still modifying... [id=41UMLL7yDtzoa0000Wimzw, 10s elapsed]
   bigip_waf_policy.this: Modifications complete after 17s [id=41UMLL7yDtzoa0000Wimzw]
   
   Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
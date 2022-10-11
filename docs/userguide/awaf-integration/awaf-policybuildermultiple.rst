.. _awaf-policybuildermultiple:

Scenario #6: Managing an Advanced WAF Policy with Policy Builder on multiple devices
====================================================================================
 
The goal of this lab is to manage Policy Builder Suggestions an Advanced WAF Policy from on multiple devices or clusters.

Goals
The goal of this lab is to manage Policy Builder Suggestions an Advanced WAF Policy from on multiple devices or clusters. Several use cases are covered here:

Multiple devices serving and protecting the same application (multiple datacenters, application spanned across multiple clouds... By nature, each standalone device or clusters can see different traffic patterns so the suggestions can be somehow differents. The goal here is to consolidate the suggestions before enforcing them.
Production BIG-IPs protecting the application therefore seeing the real life traffic flow for seeding the Policy Builder but all changes need to be first validated in the qualification environment before enforcing into production.
.. Note:: The two uses cases aforementioned are not mutually exclusive and can be managed within a single workflow


Pre-requisites
--------------

on the BIG-IPs:

 version 16.1 minimal
 Advanced WAF Provisioned
 credentials with REST API access
 an Advanced WAF Policy with Policy Builder enabled and Manual traffic Learning
on Terraform:

 use of F5 bigip provider version 1.16.0 minimal
 use of Hashicorp version following Link



Policy Creation
You need to create the 4 following files:

variables.tf

variable prod_dc1_bigip {}
variable prod_cloud_bigip {}
variable qa_bigip {}
variable username {}
variable password {}
inputs.auto.tfvars

prod_dc1_bigip = "10.1.1.8"
prod_cloud_bigip = "10.1.1.7"
qa_bigip = "10.1.1.9"
username = "admin"
password = "WhatisYourPassword?"
main.tf

terraform {
  required_providers {
    bigip = {
      source 			= "F5Networks/bigip"
      version 			= "1.15"
    }
  }
}

provider "bigip" {
  alias    			= "prod1"
  address  			= var.prod_dc1_bigip
  username 			= var.username
  password 			= var.password
}

provider "bigip" {
  alias    			= "prod2"
  address  			= var.prod_cloud_bigip
  username 			= var.username
  password 			= var.password
}

provider "bigip" {
  alias    			= "qa"
  address  			= var.qa_bigip
  username 			= var.username
  password 			= var.password
}

resource "bigip_waf_policy" "P1S6" {
    provider	           	= bigip.prod1
    application_language 	= "utf-8"
    partition			= "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
}

resource "bigip_waf_policy" "P2S6" {
    provider                    = bigip.prod2
    application_language 	= "utf-8"
    partition                   = "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
}

resource "bigip_waf_policy" "QAS6" {
    provider                    = bigip.qa
    application_language 	= "utf-8"
    partition                   = "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
}
outputs.tf

output "P1S6Id" {
	value	= bigip_waf_policy.P1S6.policy_id
}
output "P1S6JSON" {
	value   = bigip_waf_policy.P1S6.policy_export_json
}
output "P2S6Id" {
	value	= bigip_waf_policy.P2S6.policy_id
}
output "P2S6JSON" {
	value   = bigip_waf_policy.P2S6.policy_export_json
}
output "QAS6Id" {
	value	= bigip_waf_policy.QAS6.policy_id
}
output "QAS6JSON" {
	value   = bigip_waf_policy.QAS6.policy_export_json
}



Simulate a WAF Policy workflow
Here is a typical workflow:

On each BIG-IP, there is a scenario6.vs Virtual Server.

We will create and associate the same WAF Policy to these Virtual Servers.
Runing traffic on Production devices. We will make sure we are not running the same requests on both Production devices so we get distinct suggestions.
Test the suggestions from Prod1 and Prod2 devices on the QA device and check that the application is not broken.
Enforce suggestions on the Production devices.
.. Note:: There are some changes that may be specific to the QA env, such as setting Trusted IP addresses. So we will make the specific tuning first. *

1. Policy creation and association
Plan and apply your new Terraform project.

foo@bar:~$ terraform init

foo@bar:~$ terraform plan -out scenario6

foo@bar:~$ terraform apply "scenario6"
Now go on your WebUI and associate the WAF Policies to the scenario6.vs Virtual Servers. (here we do it manually but it can definitely be done using the "bigip_as3" terraform resource from the same Terraform "F5Networks/bigip" provider).




2. Running Real life traffic
Now, run both legitimate AND illegitimate traffic against your two production BIG-IP devices (scenario6 virtual servers on PROD1 and PROD2 BIG-IPs). Try to throw different attacks on each devices so we make sure we collect different Policy Builder suggestions (checkout the recommended steps described on Module5).

You may have to run multiple time the same request to make sure we get a satisfying learning score.




3. Collect and test the Policy Builder suggestions.
Create a pb_suggestions.tf file:

data "bigip_waf_pb_suggestions" "S6_22AUG20221800_P1" {
  provider	         = bigip.prod1
  policy_name            = "scenario6"
  partition              = "Common"
  minimum_learning_score = 100
}

data "bigip_waf_pb_suggestions" "S6_22AUG20221800_P2" {
  provider	         = bigip.prod2
  policy_name            = "scenario6"
  partition              = "Common"
  minimum_learning_score = 100
}

output "PB_S6_22AUG20221800_P1" {
	value	= data.bigip_waf_pb_suggestions.S6_22AUG20221800_P1.json
}

output "PB_S6_22AUG20221800_P2" {
	value	= data.bigip_waf_pb_suggestions.S6_22AUG20221800_P2.json
}
and update the main.tf file on the scenario6 QA WAF Policy resource:

resource "bigip_waf_policy" "QAS6" {
    provider	         = bigip.qa
    application_language = "utf-8"
    name                 = "scenario6"
    partition            = "Common"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario6.body
    modifications        = [data.bigip_waf_pb_suggestions.S6_22AUG20221800_P1.json, data.bigip_waf_pb_suggestions.S6_22AUG20221800_P2.json]
}
.. Note:: There are obviously some redundant learning suggestions on both data sources but the Declarative WAF API automatically removes them.

Now you can test your application through the QA device.

For UDF users: check https://qa.f5demo.fch and see that the application is not broken and attacks are blocked

4. Enforce suggestions on the Production devices
In a real life scenario, there are two ways we can consider this step:

a) the QA device WAF Policy should be 100% consistent with production devices

b) the QA device WAF Policy may have settings differences with production devices (Trusted IP exceptions for example)

a) QA.WAF == PROD.WAF
That is the easiest way. After validating the suggestions and removing the potential False Positives, just output the JSON policy from QA and refer to it as a policy_import_json argument in the production BIG-IPs

In this case, update the main.tf file

resource "bigip_waf_policy" "P1S6" {
    provider	           	= bigip.prod1
    application_language 	= "utf-8"
    partition			= "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
    policy_import_json          = bigip_waf_policy.QAS6.policy_export_json
}

resource "bigip_waf_policy" "P2S6" {
    provider	           	= bigip.prod2
    application_language 	= "utf-8"
    partition		    	= "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
    policy_import_json          = bigip_waf_policy.QAS6.policy_export_json
}

resource "bigip_waf_policy" "QAS6" {
    provider	           	= bigip.qa
    application_language 	= "utf-8"
    partition                   = "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
}
now, plan & apply!:

foo@bar:~$ terraform plan -out scenario6
foo@bar:~$ terraform apply "scenario6"
b) QA.WAF != PROD.WAF
In this case, we need to manage the learning suggestions as a separate modifications entity that has to move between WAF Policies.

The learning suggestions, when imported into the QA WAF Policy, are deduplicated and ingested into the WAF Policy. However, they remain in a dedicated space of the Declarative REST JSON: the modifications array. So, the goal is to import only this section back to the production devices, so any differences in the core entities are not affected.

resource "bigip_waf_policy" "P1S6" {
    provider	           	= bigip.prod1
    application_language 	= "utf-8"
    partition                   = "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
    policy_import_json          = bigip_waf_policy.QAS6.policy_export_json.modifications
}

resource "bigip_waf_policy" "P2S6" {
    provider	           	= bigip.prod2
    application_language 	= "utf-8"
    partition		    	= "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
    policy_import_json          = bigip_waf_policy.QAS6.policy_export_json.modifications
}

resource "bigip_waf_policy" "QAS6" {
    provider	           	= bigip.qa
    application_language 	= "utf-8"
    partition                   = "Common"
    name                 	= "scenario6"
    enforcement_mode     	= "blocking"
    template_name        	= "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
}
now, plan & apply!:

foo@bar:~$ terraform plan -out scenario6
foo@bar:~$ terraform apply "scenario6"
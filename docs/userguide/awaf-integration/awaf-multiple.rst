.. _awaf-multiple:

Scenario #4: Managing an Advanced WAF Policy on different devices
=================================================================

.. seealso:: https://github.com/fchmainy/awaf_tf_docs/tree/main/4.multiple

The goal of this lab is to manage an Advanced WAF Policy on multiple devices.

It can be:

different standalone devices serving the same applications
different devices serving different purposes, for example changes tested first on a QA/Test BIG-IP before applying into production.
Pre-requisites
on the BIG-IP:

 version 16.1 minimal
 Advanced WAF Provisioned
 credentials with REST API access
on Terraform:

 use of F5 bigip provider version 1.15.0 minimal
 use of Hashicorp version following Link
Policy Creation
Create 4 files:

variables.tf

variable qa_bigip {}
variable prod_bigip {}
variable username {}
variable password {}
inputs.auto.tfvars

qa_bigip = "10.1.1.9:443"
prod_bigip = "10.1.1.8:443"
username = "admin"
password = "whatIsYourBigIPPassword?"
main.tf

terraform {
  required_providers {
    bigip = {
      source = "F5Networks/bigip"
      version = "1.15"
    }
  }
}
provider "bigip" {
  alias    = "qa"
  address  = var.qa_bigip
  username = var.username
  password = var.password
}
provider "bigip" {
  alias    = "prod"
  address  = var.prod_bigip
  username = var.username
  password = var.password
}

data "http" "scenario4" {
  url = "https://raw.githubusercontent.com/fchmainy/awaf_tf_docs/main/0.Appendix/scenario4.json"
  request_headers = {
  	Accept = "application/json"
  }
}

resource "bigip_waf_policy" "s4_qa" {
    provider	    	 = bigip.qa
    application_language = "utf-8"
    partition            = "Common"
    name                 = "scenario4"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario4.body
}

resource "bigip_waf_policy" "s4_prod" {
    provider	         = bigip.prod
    application_language = "utf-8"
    partition            = "Common"
    name                 = "scenario4"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario4.body
}
Note: the template name can be set to anything. When it is imported, we will overwrite the value

Here, we are referencing an existing policy from a GitHub repository but it can also be created from zero on both BIG-IPs.

Now initialize, plan and apply your new Terraform project.

foo@bar:~$ terraform init
Initializing the backend...

Initializing provider plugins...
[...]
Terraform has been successfully initialized!

foo@bar:~$ terraform plan -out scenario4 > output_scenario4.1
foo@bar:~$ more output_scenario4.1
foo@bar:~$ terraform apply "scenario4"
You can check on both BIG-IPs, the two policies are here and very consistent.

Simulate a WAF Policy workflow
Here is a common workflow:

enforcing attack signatures on the QA environment
checking if these changes does not break the application and identify potential False Positives
applying the changes on QA before applying them on Production
Enforcing attack signatures on the QA environment
In order to facilitate the tracking of attack signature changes, we are using here a terraform hcl map. Add this signature list definition in the inputs.auto.tfvars file:

signatures = {
    200101559 = {
        signature_id    = 200101559
        description     = "src http: (Header)"
        enabled         = true
        perform_staging = false
    }
    200101558 = {
        signature_id    = 200101558
        description     = "src http: (Parameter)"
        enabled         = true
        perform_staging = false
    }
    200003067 = {
        signature_id    = 200003067
        description     = "\"/..namedfork/data\" execution attempt (Headers)"
        enabled         = true
        perform_staging = false
    }
    200003066 = {
        signature_id    = 200003066
        description     = "\"/..namedfork/data\" execution attempt (Parameters)"
        enabled         = true
        perform_staging = false
    }
    200003068 = {
        signature_id    = 200003068
        description     = "\"/..namedfork/data\" execution attempt (URI)"
        enabled         = true
        perform_staging = false
    }
}
Now, we create a signatures.tf file with a map to all the attack signatures defied previously:

variable "signatures" {
  type = map(object({
        signature_id    = number
	enabled		= bool
	perform_staging	= bool
        description     = string
  }))
}


data "bigip_waf_signatures" "map_qa" {
  provider	        = bigip.qa
  for_each		= var.signatures
  signature_id		= each.value["signature_id"]
  description		= each.value["description"]
  enabled		= each.value["enabled"]
  perform_staging	= each.value["perform_staging"]
}

data "bigip_waf_signatures" "map_prod" {
  provider	        = bigip.prod
  for_each		= var.signatures
  signature_id		= each.value["signature_id"]
  description		= each.value["description"]
  enabled		= each.value["enabled"]
  perform_staging	= each.value["perform_staging"]
}
As you can see, we defined 2 different maps: one for the QA BIG-IP and one for the PRODUCTION BIG-IP because the "bigip_waf_signatures" data source are linked to their BIG-IP in order to have consistencies. Unlike the parameters and urls data sources which are just "json payload generators", the attack signature data sources has to read first the existence of the signatures id and their status on the BIG-IP before applying a configuration change.

Now finally, update the main.tf file:

resource "bigip_waf_policy" "s4_qa" {
    provider	    	 = bigip.qa
    application_language = "utf-8"
    partition            = "Common"
    name                 = "scenario4"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario4.body
    signatures           = [ for k,v in data.bigip_waf_signatures.map_qa: v.json ]
}

resource "bigip_waf_policy" "s4_prod" {
    provider	    	 = bigip.prod
    application_language = "utf-8"
    partition            = "Common"
    name                 = "scenario4"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario4.body
}
now, plan & apply!:

foo@bar:~$ terraform plan -out scenario4 > output_scenario4.2
foo@bar:~$ more output_scenario4.2
foo@bar:~$ terraform apply "scenario4"
We can verify that the 5 attack signatures have been enabled and enforced on the scenario4 WAF Policy on the QA BIG-IP (first 5 lines in the attack signatures list of the Advanced WAF Policy).

Now, the applicatiopn owner identified that these last changes on the QA device have introduced some FP. Using the log events on the Advanced WAF GUI, we identified that :

the attack signature "200101558" should be disabled globally
the attack signature "200003068" should be disabled for the "/U1" URL
the attack signaure "200003067" should be enabled globally but disabled specifically for the parameter "P1".
so we can proceed to the final changes before enforcing into production:

inputs.auto.tfvars file:

signatures = {
    200101559 = {
        signature_id    = 200101559
        description     = "src http: (Header)"
        enabled         = true
        perform_staging = false
    }
    200101558 = {
        signature_id    = 200101558
        description     = "src http: (Parameter)"
        enabled         = false
        perform_staging = false
    }
    200003067 = {
        signature_id    = 200003067
        description     = "\"/..namedfork/data\" execution attempt (Headers)"
        enabled         = true
        perform_staging = false
    }
    200003066 = {
        signature_id    = 200003066
        description     = "\"/..namedfork/data\" execution attempt (Parameters)"
        enabled         = true
        perform_staging = false
    }
    200003068 = {
        signature_id    = 200003068
        description     = "\"/..namedfork/data\" execution attempt (URI)"
        enabled         = true
        perform_staging = false
    }
}
parameters.tf file:

data "bigip_waf_entity_parameter" "P1" {
  name            		= "P1"
  type            		= "explicit"
  data_type       		= "alpha-numeric"
  perform_staging 		= true
  signature_overrides_disable 	= [200003067]
  //url		  		= data.bigip_waf_entity_url.U1
}
urls.tf file:

data "bigip_waf_entity_url" "U1" {
  name		              	= "/U1"
  type                        	= "explicit"
  perform_staging             	= false
  signature_overrides_disable 	= [200003068]
}
update the main.tf file:

resource "bigip_waf_policy" "s4_qa" {
    provider	    	 = bigip.qa
    application_language = "utf-8"
    partition            = "Common"
    name                 = "scenario4"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario4.body
    signatures		 = [ for k,v in data.bigip_waf_signatures.map_qa: v.json ]
    parameters		 = [data.bigip_waf_entity_parameter.P1.json]
    urls		 = [data.bigip_waf_entity_url.U1.json]
}

resource "bigip_waf_policy" "s4_prod" {
    provider	    	 = bigip.prod
    application_language = "utf-8"
    partition            = "Common"
    name                 = "scenario4"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario4.body
    signatures		 = [ for k,v in data.bigip_waf_signatures.map_prod: v.json ]
    parameters		 = [data.bigip_waf_entity_parameter.P1.json]
    urls		 = [data.bigip_waf_entity_url.U1.json]
}
now, plan & apply!:

foo@bar:~$ terraform plan -out scenario4 > output_scenario4.3
foo@bar:~$ more output_scenario4.3
foo@bar:~$ terraform apply "scenario4"
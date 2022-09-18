.. _awaf-integration:

**`Scenario #5: Managing an A.WAF Policy with Policy Builder on a single device <https://github.com/fchmainy/awaf_tf_docs/tree/main/5.policyBuilderSingle>`_**
 
The goal of this lab is to manage Policy Builder Suggestions an A.WAF Policy on a single device or cluster.

Goals
The goal of this lab is to manage Policy Builder Suggestions an A.WAF Policy on a single device or cluster. As the traffic flows through the BIG-IP, it is easy to manage suggestions from the Policy Builder and enforce them on the WAF Policy. It also shows what can be the management workflow:

the security engineer regularly checks the sugestions directly on the BIG-IP WebUI and clean the irrelevant suggestions.
once the cleaning is done, the terraform engineer (who can also be the security engineer btw) issue a terraform apply for the current suggestions. You can filter the suggestions on their scoring level (from 5 to 100% - 100% having the highest confidence level).
Every suggestions application can be tracked on Terraform and can easily be roll-backed if needed.



Pre-requisites
on the BIG-IP:

 version 16.1 minimal
 A.WAF Provisioned
 credentials with REST API access
 an A.WAF Policy with Policy Builder enabled and Manual traffic Learning
on Terraform:

 use of F5 bigip provider version 1.15.0 minimal
 use of Hashicorp version followinf Link



Policy Creation
We already have exported a WAF Policy called scenario5.json available here including several Policy Builder Suggestions so you won't have to generate traffic.

So you have to create 4 files:

variables.tf

variable prod_bigip {}
variable username {}
variable password {}
inputs.auto.tfvars

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
  alias    = "prod"
  address  = var.prod_bigip
  username = var.username
  password = var.password
}

data "http" "scenario5" {
  url = "https://raw.githubusercontent.com/fchmainy/awaf_tf_docs/main/0.Appendix/Common_scenario5__2022-8-12_15-49-28__prod1.f5demo.com.json"
  request_headers = {
  	Accept = "application/json"
  }
}

resource "bigip_waf_policy" "this" {
    provider	           = bigip.prod
    application_language = "utf-8"
    partition            = "Common"
    name                 = "scenario5"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario5.body
}
Note: the template name can be set to anything. When it is imported, we will overwrite the value

outputs.tf

output "policyId" {
	value	= bigip_waf_policy.this.policy_id
}

output "policyJSON" {
        value   = bigip_waf_policy.this.policy_export_json
}
Now initialize, plan and apply your new Terraform project.

foo@bar:~$ terraform init

foo@bar:~$ terraform plan -out scenario5

foo@bar:~$ terraform apply "scenario5"
Now you can go on your BIG-IP UI and associate the A.WAF Policy scenario5 to the Virtual Server scenario5.vs.

Note: remember, the Virtual Server and the whole application service can be automated using the BIG-IP provider with the AS3 or FAST resources.




Simulate a WAF Policy workflow
Change the Policy Builder process (For testing and demoing purpose only):
First, go to the DVWA WAF Policy on your BIG-IP TMUI (if you are using UDF, the WAF policy is called scenario5 and is located under the Common partition.
In the Learning and blocking Settings (Security ›› Application Security : Policy Building : Learning and Blocking Settings), at the very bottom of the page, go on the Loosen Policy settings in the Advanced view of the Policy Building Process.
Change the different sources, spread out over a time period of at least value from 10 to 1 so the policy builder generates learning suggestions more rapidely.
Browse the Vulnerable Application
Now browse the DVWA web application through the AWAF Virtual Server. The credentials to log in to DVWA is admin/password.

Go on the *DVWA Security menu and change the level to Low then Submit
Browse the DVWA website by clicking into any menus.
Then generate some attacks:
SQL Injection: %' or 1='1 ' and 1=0 union select null, concat(first_name,0x0a,last_name,0x0a,user,0x0a,password) from users #
XSS Reflected: <script>alert('hello')</script>
Check Learning Suggestions
Now, if you go to the WAF Policy learning suggestions, you will find multiple suggestions with a high score of 100% (because we have not been picky in the learning process settings).

Here is a typical workflow in a real life:

the security engineer (yourself) regularly checks the sugestions directly on the BIG-IP WebUI and clean the irrelevant suggestions.
once the cleaning is done, the terraform engineer (can either be the same person or different) creates a unique bigip_waf_pb_suggestions data source before issuing a terraform apply for the current suggestions. You can filter the suggestions on their scoring level (from 5 to 100% - 100% having the highest confidence level).
Note: Every suggestions application can be tracked on Terraform and can easily be roll-backed if needed.


1. Go to your BIG-IP WebUI and clean the irrelevant suggestions
⚠️ IMPORTANT you can ignore suggestions but you should never accept them on the WebUI, otherwise you will then have to reconciliate the changes between the WAF Policy on the BIG-IP and the latest known WAF Policy in your terraform state.

For example, remove all the suggestions with a scoring = 1%


2. Use Terraform to enforce the policy builder suggestions
Create a suggestions.tf file:

the name of the bigip_waf_pb_suggestions data source should be unique so we can track what modifications have been enforced and when it was.

data "bigip_waf_pb_suggestions" "AUG3rd20221715" {
  provider	           = bigip.prod 
  policy_name            = "scenario5"
  partition              = "Common"
  minimum_learning_score = 100
}

output "AUG3rd20221715" {
	value	= data.bigip_waf_pb_suggestions.AUG3rd20221715.json
}
You can check here the suggestions before they are applied to the BIG-IP:

foo@bar:~$ terraform plan -out scenario5

foo@bar:~$ terraform apply "scenario5"

foo@bar:~$ terraform output AUG3rd20221715 | jq '. | fromjson'
You will get the JSON list of suggestions that have a learning score of 100%.

{
    "suggestions": [
      {
        "action": "update-append",
        "description": "Add/Update Parameter. Disable the matched signature on the matched Parameter",
        "entity": {
          "level": "global",
          "name": "id"
        },
        "entityChanges": {
          "signatureOverrides": [
            {
              "enabled": false,
              "name": "SQL-INJ ' UNION SELECT (Parameter)",
              "signatureId": 200002736
            }
          ],
          "type": "explicit"
        },
        "entityType": "parameter"
      },
[...],      
      {
        "action": "add-or-update",
        "description": "Add Policy Server Technology",
        "entity": {
          "serverTechnologyName": "Unix/Linux"
        },
        "entityType": "server-technology"
      }
    ]
  }
update the main.tf file:

resource "bigip_waf_policy" "this" {
    provider             = bigip.prod
    application_language = "utf-8"
    partition            = "Common"
    name                 = "scenario5"
    template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
    type                 = "security"
    policy_import_json   = data.http.scenario5.body
    modifications        = [data.bigip_waf_pb_suggestions.AUG3rd20221715.json]
}
now, plan & apply!:

foo@bar:~$ terraform plan -out scenario5

foo@bar:~$ terraform apply "scenario5"
You can check on your BIGIP UI that the server technologies and other suggestions have been succesfully enforced to your WAF Policy.
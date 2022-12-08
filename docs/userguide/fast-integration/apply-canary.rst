.. _fast-integration-apply-canary:

Scenario #7: Applying Canary deployment strategy for HTTPS application with Web Application Firewall policy
===========================================================================================================

The goal of this template is to deploy a new HTTPS application using canary deployment strategy with Web Application Firewall policy on BIG-IP using Terraform as the orchestrator. Canary strategy will be based on HTTP header.


Pre-requisites
--------------

on the BIG-IP:

 version 16.1 minimal
 credentials with REST API access
on Terraform:

 use of F5 bigip provider version 1.16.0 minimal
 use of Hashicorp version following Link


Create HTTPS application
------------------------
Create 4 files:

- main.tf
- variables.tf
- inputs.auto.tfvars
- providers.tf

.. code-block:: json
   :caption: 
   :linenos:

variables.tf

variable "bigip" {}
variable "username" {}
variable "password" {}
variable "policyname" {
  type    = string
  default = ""

}
variable "partition" {
  type    = string
  default = "Common"
}

|

.. code-block:: json
   :caption: 
   :linenos:

inputs.tfvars

bigip      = "10.1.1.9:443"
username   = "admin"
password   = "A7U+=$vJ"
partition  = "Common"
policyname = "myApp7_ltm_policy"

|

.. code-block:: json
   :caption: 
   :linenos:

providers.tf

terraform {
  required_providers {
    bigip = {
      source = "F5Networks/bigip"
      version = ">= 1.16.0"
    }
  }
}
provider "bigip" {
  address  = var.bigip
  username = var.username
  password = var.password
}

|

.. code-block:: json
   :caption: 
   :linenos:

main.tf

resource "bigip_waf_policy" "app1_waf_v1" {
  provider             = bigip
  description          = "Current version of the WAF Policy"
  name                 = "v1"
  partition            = "Common"
  template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  server_technologies  = ["Apache Tomcat", "MySQL", "Unix/Linux"]
}

resource "bigip_waf_policy" "app1_waf_v2" {
  provider             = bigip
  description          = "new version of the WAF Policy"
  name                 = "v2"
  partition            = "Common"
  template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  server_technologies  = ["Apache Tomcat", "MySQL", "Unix/Linux", "MongoDB"]
}

module "canary_app1" {
  source = "github.com/f5devcentral/fast-terraform//canary_policy_header?ref=v1.0.0"
  providers = {
    bigip = bigip
  }
  name               = var.policyname
  partition          = var.partition
  header_name        = "user_profile"
  header_value       = "earlyAdopter"
  new_waf_policy     = bigip_waf_policy.app1_waf_v2.name
  current_waf_policy = bigip_waf_policy.app1_waf_v1.name
  depends_on         = [bigip_waf_policy.app1_waf_v1, bigip_waf_policy.app1_waf_v2]
}

resource "bigip_fast_https_app" "this" {
  application = "myApp7"
  tenant      = "scenario7"
  virtual_server {
    ip   = "10.1.10.227"
    port = 443
  }
  tls_server_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
  pool_members {
    addresses = ["10.1.10.120", "10.1.10.121", "10.1.10.122"]
    port      = 80
  }
  snat_pool_address     = ["10.1.10.50", "10.1.10.51", "10.1.10.52"]
  endpoint_ltm_policy   = ["${module.canary_app1.ltmPolicyName}"]
  security_log_profiles = ["/Common/Log all requests"]
  depends_on            = [bigip_waf_policy.app1_waf_v1, bigip_waf_policy.app1_waf_v2, module.canary_app1.lt
mPolicyName]
}



Run the following commands so you can:

1. Initialize the terraform project
2. Plan the changes
3. Apply the changes

::

  
$ terraform init -upgrade
Upgrading modules...
Downloading git::https://github.com/fchmainy/waf_modules.git?ref=v1.0.8 for canary_app1...
- canary_app1 in .terraform/modules/canary_app1/canary_policy_header

Initializing the backend...

Initializing provider plugins...
- Finding f5networks/bigip versions matching ">= 1.16.0"...
- Using previously-installed f5networks/bigip v1.16.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.


$ terraform plan -out scenario7

Terraform used the selected providers to generate the following execution plan. Resource actions are
indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # bigip_fast_https_app.this will be created
  + resource "bigip_fast_https_app" "this" {
      + application           = "myApp7"
      + endpoint_ltm_policy   = [
          + "/Common/myApp7_ltm_policy",
        ]
      + fast_https_json       = (known after apply)
      + id                    = (known after apply)
      + load_balancing_mode   = "least-connections-member"
      + security_log_profiles = [
          + "/Common/Log all requests",
        ]
      + snat_pool_address     = [
          + "10.1.10.50",
          + "10.1.10.51",
          + "10.1.10.52",
        ]
      + tenant                = "scenario7"

      + pool_members {
          + addresses = [
              + "10.1.10.120",
              + "10.1.10.121",
              + "10.1.10.122",
            ]
          + port      = 80
        }

      + tls_server_profile {
          + tls_cert_name = "/Common/default.crt"
          + tls_key_name  = "/Common/default.key"
        }

      + virtual_server {
          + ip   = "10.1.10.227"
          + port = 443
        }
    }

  # bigip_waf_policy.app1_waf_v1 will be created
  + resource "bigip_waf_policy" "app1_waf_v1" {
      + application_language = "utf-8"
      + case_insensitive     = false
      + description          = "Current version of the WAF Policy"
      + enable_passivemode   = false
      + enforcement_mode     = "blocking"
      + id                   = (known after apply)
      + name                 = "v1"
      + partition            = "Common"
      + policy_export_json   = (known after apply)
      + policy_id            = (known after apply)
      + server_technologies  = [
          + "Apache Tomcat",
          + "MySQL",
          + "Unix/Linux",
        ]
      + template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
      + type                 = "security"
    }

  # bigip_waf_policy.app1_waf_v2 will be created
  + resource "bigip_waf_policy" "app1_waf_v2" {
      + application_language = "utf-8"
      + case_insensitive     = false
      + description          = "new version of the WAF Policy"
      + enable_passivemode   = false
      + enforcement_mode     = "blocking"
      + id                   = (known after apply)
      + name                 = "v2"
      + partition            = "Common"
      + policy_export_json   = (known after apply)
      + policy_id            = (known after apply)
      + server_technologies  = [
          + "Apache Tomcat",
          + "MySQL",
          + "Unix/Linux",
          + "MongoDB",
        ]
      + template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
      + type                 = "security"
    }

  # module.canary_app1.bigip_ltm_policy.canary will be created
  + resource "bigip_ltm_policy" "canary" {
      + controls = [
          + "asm",
        ]
      + id       = (known after apply)
      + name     = "/Common/myApp7_ltm_policy"
      + requires = [
          + "http",
        ]
      + strategy = "first-match"

      + rule {
          + name = "ea"

          + action {
              + app_service          = (known after apply)
              + application          = (known after apply)
              + asm                  = true
              + avr                  = (known after apply)
              + cache                = (known after apply)
              + carp                 = (known after apply)
              + category             = (known after apply)
              + classify             = (known after apply)
              + clone_pool           = (known after apply)
              + code                 = (known after apply)
              + compress             = (known after apply)
              + connection           = false
              + content              = (known after apply)
              + cookie_hash          = (known after apply)
              + cookie_insert        = (known after apply)
              + cookie_passive       = (known after apply)
              + cookie_rewrite       = (known after apply)
              + decompress           = (known after apply)
              + defer                = (known after apply)
              + destination_address  = (known after apply)
              + disable              = (known after apply)
              + domain               = (known after apply)
              + enable               = (known after apply)
              + expiry               = (known after apply)
              + expiry_secs          = (known after apply)
              + expression           = (known after apply)
              + extension            = (known after apply)
              + facility             = (known after apply)
              + forward              = false
              + from_profile         = (known after apply)
              + hash                 = (known after apply)
              + host                 = (known after apply)
              + http                 = (known after apply)
              + http_basic_auth      = (known after apply)
              + http_cookie          = (known after apply)
              + http_header          = (known after apply)
              + http_referer         = (known after apply)
              + http_reply           = (known after apply)
              + http_set_cookie      = (known after apply)
              + http_uri             = (known after apply)
              + ifile                = (known after apply)
              + insert               = (known after apply)
              + internal_virtual     = (known after apply)
              + ip_address           = (known after apply)
              + key                  = (known after apply)
              + l7dos                = (known after apply)
              + length               = (known after apply)
              + location             = (known after apply)
              + log                  = (known after apply)
              + ltm_policy           = (known after apply)
              + member               = (known after apply)
              + message              = (known after apply)
              + netmask              = (known after apply)
              + nexthop              = (known after apply)
              + node                 = (known after apply)
              + offset               = (known after apply)
              + path                 = (known after apply)
              + pem                  = (known after apply)
              + persist              = (known after apply)
              + pin                  = (known after apply)
              + policy               = "/Common/v2"
              + pool                 = (known after apply)
              + port                 = (known after apply)
              + priority             = (known after apply)
              + profile              = (known after apply)
              + protocol             = (known after apply)
              + query_string         = (known after apply)
              + rateclass            = (known after apply)
              + redirect             = (known after apply)
              + remove               = (known after apply)
              + replace              = (known after apply)
              + request              = true
              + request_adapt        = (known after apply)
              + reset                = (known after apply)
              + response             = (known after apply)
              + response_adapt       = (known after apply)
              + scheme               = (known after apply)
              + script               = (known after apply)
              + select               = (known after apply)
              + server_ssl           = (known after apply)
              + set_variable         = (known after apply)
              + snat                 = (known after apply)
              + snatpool             = (known after apply)
              + source_address       = (known after apply)
              + ssl_client_hello     = (known after apply)
              + ssl_server_handshake = (known after apply)
              + ssl_server_hello     = (known after apply)
              + ssl_session_id       = (known after apply)
              + status               = (known after apply)
              + tcl                  = (known after apply)
              + tcp_nagle            = (known after apply)
              + text                 = (known after apply)
              + timeout              = (known after apply)
              + tm_name              = (known after apply)
              + uie                  = (known after apply)
              + universal            = (known after apply)
              + value                = (known after apply)
              + virtual              = (known after apply)
              + vlan                 = (known after apply)
              + vlan_id              = (known after apply)
              + wam                  = (known after apply)
              + write                = (known after apply)
            }

          + condition {
              + address                 = (known after apply)
              + all                     = true
              + app_service             = (known after apply)
              + browser_type            = (known after apply)
              + browser_version         = (known after apply)
              + case_insensitive        = true
              + case_sensitive          = (known after apply)
              + cipher                  = (known after apply)
              + cipher_bits             = (known after apply)
              + client_accepted         = (known after apply)
              + client_ssl              = (known after apply)
              + code                    = (known after apply)
              + common_name             = (known after apply)
              + contains                = (known after apply)
              + continent               = (known after apply)
              + country_code            = (known after apply)
              + country_name            = (known after apply)
              + cpu_usage               = (known after apply)
              + device_make             = (known after apply)
              + device_model            = (known after apply)
              + domain                  = (known after apply)
              + ends_with               = (known after apply)
              + equals                  = true
              + exists                  = (known after apply)
              + expiry                  = (known after apply)
              + extension               = (known after apply)
              + external                = true
              + geoip                   = (known after apply)
              + greater                 = (known after apply)
              + greater_or_equal        = (known after apply)
              + host                    = (known after apply)
              + http_basic_auth         = (known after apply)
              + http_cookie             = (known after apply)
              + http_header             = true
              + http_host               = (known after apply)
              + http_method             = (known after apply)
              + http_referer            = (known after apply)
              + http_set_cookie         = (known after apply)
              + http_status             = (known after apply)
              + http_uri                = (known after apply)
              + http_user_agent         = (known after apply)
              + http_version            = (known after apply)
              + index                   = (known after apply)
              + internal                = (known after apply)
              + isp                     = (known after apply)
              + last_15secs             = (known after apply)
              + last_1min               = (known after apply)
              + last_5mins              = (known after apply)
              + less                    = (known after apply)
              + less_or_equal           = (known after apply)
              + local                   = (known after apply)
              + major                   = (known after apply)
              + matches                 = (known after apply)
              + minor                   = (known after apply)
              + missing                 = (known after apply)
              + mss                     = (known after apply)
              + not                     = (known after apply)
              + org                     = (known after apply)
              + password                = (known after apply)
              + path                    = (known after apply)
              + path_segment            = (known after apply)
              + port                    = (known after apply)
              + present                 = true
              + protocol                = (known after apply)
              + query_parameter         = (known after apply)
              + query_string            = (known after apply)
              + region_code             = (known after apply)
              + region_name             = (known after apply)
              + remote                  = true
              + request                 = true
              + response                = (known after apply)
              + route_domain            = (known after apply)
              + rtt                     = (known after apply)
              + scheme                  = (known after apply)
              + server_name             = (known after apply)
              + ssl_cert                = (known after apply)
              + ssl_client_hello        = (known after apply)
              + ssl_extension           = (known after apply)
              + ssl_server_handshake    = (known after apply)
              + ssl_server_hello        = (known after apply)
              + starts_with             = (known after apply)
              + tcp                     = (known after apply)
              + text                    = (known after apply)
              + tm_name                 = "user_profile"
              + unnamed_query_parameter = (known after apply)
              + user_agent_token        = (known after apply)
              + username                = (known after apply)
              + value                   = (known after apply)
              + values                  = [
                  + "earlyAdopter",
                ]
              + version                 = (known after apply)
              + vlan                    = (known after apply)
              + vlan_id                 = (known after apply)
            }
        }
      + rule {
          + name = "default"

          + action {
              + app_service          = (known after apply)
              + application          = (known after apply)
              + asm                  = true
              + avr                  = (known after apply)
              + cache                = (known after apply)
              + carp                 = (known after apply)
              + category             = (known after apply)
              + classify             = (known after apply)
              + clone_pool           = (known after apply)
              + code                 = (known after apply)
              + compress             = (known after apply)
              + connection           = false
              + content              = (known after apply)
              + cookie_hash          = (known after apply)
              + cookie_insert        = (known after apply)
              + cookie_passive       = (known after apply)
              + cookie_rewrite       = (known after apply)
              + decompress           = (known after apply)
              + defer                = (known after apply)
              + destination_address  = (known after apply)
              + disable              = (known after apply)
              + domain               = (known after apply)
              + enable               = true
              + expiry               = (known after apply)
              + expiry_secs          = (known after apply)
              + expression           = (known after apply)
              + extension            = (known after apply)
              + facility             = (known after apply)
              + forward              = false
              + from_profile         = (known after apply)
              + hash                 = (known after apply)
              + host                 = (known after apply)
              + http                 = (known after apply)
              + http_basic_auth      = (known after apply)
              + http_cookie          = (known after apply)
              + http_header          = (known after apply)
              + http_referer         = (known after apply)
              + http_reply           = (known after apply)
              + http_set_cookie      = (known after apply)
              + http_uri             = (known after apply)
              + ifile                = (known after apply)
              + insert               = (known after apply)
              + internal_virtual     = (known after apply)
              + ip_address           = (known after apply)
              + key                  = (known after apply)
              + l7dos                = (known after apply)
              + length               = (known after apply)
              + location             = (known after apply)
              + log                  = (known after apply)
              + ltm_policy           = (known after apply)
              + member               = (known after apply)
              + message              = (known after apply)
              + netmask              = (known after apply)
              + nexthop              = (known after apply)
              + node                 = (known after apply)
              + offset               = (known after apply)
              + path                 = (known after apply)
              + pem                  = (known after apply)
              + persist              = (known after apply)
              + pin                  = (known after apply)
              + policy               = "/Common/v1"
              + pool                 = (known after apply)
              + port                 = (known after apply)
              + priority             = (known after apply)
              + profile              = (known after apply)
              + protocol             = (known after apply)
              + query_string         = (known after apply)
              + rateclass            = (known after apply)
              + redirect             = (known after apply)
              + remove               = (known after apply)
              + replace              = (known after apply)
              + request              = true
              + request_adapt        = (known after apply)
              + reset                = (known after apply)
              + response             = (known after apply)
              + response_adapt       = (known after apply)
              + scheme               = (known after apply)
              + script               = (known after apply)
              + select               = (known after apply)
              + server_ssl           = (known after apply)
              + set_variable         = (known after apply)
              + snat                 = (known after apply)
              + snatpool             = (known after apply)
              + source_address       = (known after apply)
              + ssl_client_hello     = (known after apply)
              + ssl_server_handshake = (known after apply)
              + ssl_server_hello     = (known after apply)
              + ssl_session_id       = (known after apply)
              + status               = (known after apply)
              + tcl                  = (known after apply)
              + tcp_nagle            = (known after apply)
              + text                 = (known after apply)
              + timeout              = (known after apply)
              + tm_name              = (known after apply)
              + uie                  = (known after apply)
              + universal            = (known after apply)
              + value                = (known after apply)
              + virtual              = (known after apply)
              + vlan                 = (known after apply)
              + vlan_id              = (known after apply)
              + wam                  = (known after apply)
              + write                = (known after apply)
            }
        }
    }

Plan: 4 to add, 0 to change, 0 to destroy.

───────────────────────────────────────────────────────────────────────────────────────────────────────────

Saved the plan to: scenario7

To perform exactly these actions, run the following command to apply:
    terraform apply "scenario7"


$ terraform apply "scenario7"
bigip_waf_policy.app1_waf_v1: Creating...
bigip_waf_policy.app1_waf_v2: Creating...
bigip_waf_policy.app1_waf_v1: Still creating... [10s elapsed]
bigip_waf_policy.app1_waf_v2: Still creating... [10s elapsed]
bigip_waf_policy.app1_waf_v1: Creation complete after 17s [id=dmxiH2VYPedQA-63JPJmNA]
bigip_waf_policy.app1_waf_v2: Still creating... [20s elapsed]
bigip_waf_policy.app1_waf_v2: Creation complete after 22s [id=3FMicDmDaJZ9OxCV35PDjw]
module.canary_app1.bigip_ltm_policy.canary: Creating...
module.canary_app1.bigip_ltm_policy.canary: Creation complete after 2s [id=/Common/myApp7_ltm_policy]
bigip_fast_https_app.this: Creating...
bigip_fast_https_app.this: Still creating... [10s elapsed]
bigip_fast_https_app.this: Creation complete after 17s [id=myApp7]

Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
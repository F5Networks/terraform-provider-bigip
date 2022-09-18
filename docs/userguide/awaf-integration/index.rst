.. _awaf-integration:

AWAF Integration with Terraform
===============================

Introduction
------------
.. seealso::
   :class: sidebar

   Read `Getting Started with Declarative Policies <https://techdocs.f5.com/en-us/bigip-15-1-0/big-ip-declarative-security-policy/declarative-policy-getting-started.html#concept-1626>`_.

The Advanced Web Application Firewall (Advanced WAF) or Application Security Manager (ASM) security policies can be deployed using the declarative JSON format, facilitating easy integration into a CI/CD pipeline. The declarative policies are extracted from a source control system, for example Git, and imported into the BIG-IP.
Using the provided declarative policy templates, you can modify the necessary parameters, save the JSON file, and import the updated security policy into your BIG-IP devices. The declarative policy copies the content of the template and adds the adjustments and modifications on to it. The templates therefore allow you to concentrate only on the specific settings that need to be adapted for the specific application that the policy protects.

Terraform can be used to manage AWAF policy resource with its adjustments and modifications on a BIG-IP. It outputs an up-to-date WAF Policy in a JSON format so you can store it in a registry and/or push it to your BIG-IP.
 
AWAF Policy structure
`````````````````````
 
The supported declarative policy structure includes three logical sections: 
 
- **The "core" section** includes all the building parameters of the policy (name, description, enforcement mode, server technologies…).

- **The “adjustment" section** includes attributes of the policy that override or add to those defined in the template. Attributes included in this section can include both properties that are particular to the protected application, such as server technologies, URLs, or parameters; and modifications to settings defined by the template, such as enabling the Data Guard if it is disabled in the template and specifying Data Guard attributes. In general, the adjustments section is used for defining major features of the policy which are different from template.
 
- **The modifications section** includes actions that modify the declarative policy as it is defined in the adjustments section. In general, while the modifications section is used for frequent and granular changes that are required to tune the policy, such as reducing false positives, patching vulnerabilities, etc.
 
When an attribute is defined in both the adjustments and modifications sections, the policy is deployed with the value as it appears in the modifications section. When an attribute appears in the modifications section multiple times, the policy is deployed with the value as it appears in the latest definition in the modifications section.
 
Example Usage:



Prerequisites
-------------

Before working with declarative policies, make sure you are familiar with the F5 BIG-IP Application Security Manager and general BIG-IP terminology. Information is found in the F5 Knowledge Centers.
To read, modify, or import declarative policies, you must have:

- BIG-IP devices running version 15.1.x or later
- BIG-IP Administrator role permissions
- An active ASM or Advanced WAF license on the BIG-IP devices
- Terraform provider BIG-IP v1.15.0 and above

Example Usage
-------------

.. code-block:: json
   :caption: Example usage for json file
   :linenos:


   data "bigip_waf_entity_parameter" "Param1" {
     name            = "Param1"
     type            = "explicit"
     data_type       = "alpha-numeric"
     perform_staging = true
   }
    
   data "bigip_waf_entity_parameter" "Param2" {
     name            = "Param2"
     type            = "explicit"
     data_type       = "alpha-numeric"
     perform_staging = true
   }
    
   data "bigip_waf_entity_url" "URL" {
     name     = "URL1"
     protocol = "http"
   }
    
   data "bigip_waf_entity_url" "URL2" {
     name = "URL2"
   }
    
   resource "bigip_waf_policy" "test-awaf" {
     name                 = "testpolicyravi"
     partition            = "Common"
     template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
     application_language = "utf-8"
     enforcement_mode     = "blocking"
     server_technologies  = ["MySQL", "Unix/Linux", "MongoDB"]
     parameters           = [data.bigip_waf_entity_parameter.Param1.json, data.bigip_waf_entity_parameter.Param2.json]
     urls                 = [data.bigip_waf_entity_url.URL.json, data.bigip_waf_entity_url.URL2.json]
   }







Terraform Integration Resources and Data Sources
------------------------------------------------

AWAF Resources:

- `bigip_waf_policy <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_waf_policy>`_

AWAF Data Sources:

- `bigip_waf_entity_parameters <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/data-sources/bigip_waf_entity_parameters>`_
- `bigip_waf_entity_url <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/data-sources/bigip_waf_entity_url>`_ 
- `bigip_waf_pb_suggestions <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/data-sources/bigip_waf_pb_suggestions>`_ 
- `bigip_waf_policy <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/data-sources/bigip_waf_policy>`_ 
- `bigip_waf_signatures <https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/data-sources/bigip_waf_signatures>`_ 
 

Quickstart Guide
----------------


Additional lab guides
---------------------


.. toctree::
   :maxdepth: 2
   :includehidden:
   :glob:

   /userguide/awaf-integration/awaf-create.rst
   /userguide/awaf-integration/awaf-import.rst
   /userguide/awaf-integration/awaf-migrate.rst
   /userguide/awaf-integration/awaf-multiple.rst
   /userguide/awaf-integration/awaf-policybuildersingle.rst
   /userguide/awaf-integration/awaf-policybuildermultiple.rst
 
**`Scenario #1: Creating a WAF Policy <https://github.com/fchmainy/awaf_tf_docs/tree/main/1.create>`_**
 
The goal of this lab is to create a new A.WAF Policy from scratch and manage some entities additions.
 
**`Scenario #2: Managing with terraform an existing WAF Policy <https://github.com/fchmainy/awaf_tf_docs/tree/main/2.import>`_**
 
The goal of this lab is to take an existing A.WAF Policy -- that have been created and managed on a BIG-IP outside of Terraform -- and to import and manage its lifecycle using the F5’s BIG-IP terraform provider.
 
**`Scenario #3: Migrating a WAF Policy from a BIG-IP to another BIG-IP <https://github.com/fchmainy/awaf_tf_docs/tree/main/3.migrate>`_**
 
This lab is a variant of the previous one. It takes a manually managed A.WAF Policy from an existing BIG-IP and migrate it to a different BIG-IP through Terraform resources.
 
**`Scenario #4: Managing an A.WAF Policy on different devices <https://github.com/fchmainy/awaf_tf_docs/tree/main/4.multiple>`_**
 
The goal of this lab is to manage an A.WAF Policy on multiple devices.
 
**`Scenario #5: Managing an A.WAF Policy with Policy Builder on a single device <https://github.com/fchmainy/awaf_tf_docs/tree/main/5.policyBuilderSingle>`_**
 
The goal of this lab is to manage Policy Builder Suggestions an A.WAF Policy on a single device or cluster.
 
**`Scenario #6: Managing an A.WAF Policy with Policy Builder on multiple device <https://github.com/fchmainy/awaf_tf_docs/tree/main/6.policyBuilderMultiple>`_**
 
The goal of this lab is to manage Policy Builder Suggestions an A.WAF Policy from on multiple devices or clusters.
.. _fast-integration:

FAST Integration with Terraform
===============================

Introduction
------------
The objective of this templates is to demonstrate how FAST can be used to manage, deploy, and log changes in applications using Terraform as a resource manager through their API.

F5 BIG-IP Application Services Templates (FAST) are an easy and effective way to deploy applications on the BIG-IP system using AS3.
The FAST Extension provides a toolset for templating and managing AS3 Applications on BIG-IP.

For more information about FAST, including installation and usage information, see the FAST Documentation
https://clouddocs.f5.com/products/extensions/f5-appsvcs-templates/latest/
https://github.com/F5Networks/f5-appsvcs-templates

Example Usage

.. code-block:: json
   :caption: 
   :linenos:

   resource "bigip_fast_http_app" "myapp" {
     tenant = "Mytenant"
     application = "myapp"
     virtual_server {
       ip = "10.1.1.1"
       port = "80"
     }
     fast_create_snat_pool_address = [ "10.2.2.2" ]
     fast_create_pool_members {
       addresses = [ "10.2.2.100" ]
       port = "80"
     }
     load_balancing_mode = "least-connections-member"
     fast_create_monitor {
       send_string = "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: Close\r\n\r\n"
       response = "200 OK"
     }
   }

Terraform integration resouces/data source
             Here will have list with hyperlink to registry link FAST resource:
                               https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_fast_http_app
                               https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_fast_https_app
                               https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_fast_tcp_app


Quick start Guide
Additional lab guides

.. toctree::
   :maxdepth: 2
   :includehidden:
   :glob:

   /userguide/fast-integration/create-udp.rst
   /userguide/fast-integration/create-tcp.rst
   /userguide/fast-integration/create-http.rst
   /userguide/fast-integration/create-https.rst
   /userguide/fast-integration/create-http-existing-pool.rst












.. _fast-integration:

FAST Integration with Terraform
===============================

The objective of this template is to demonstrate how FAST can be used to manage, deploy, and log changes in applications using Terraform as a resource manager through their API.

F5 BIG-IP Application Services Templates (FAST) are an easy and effective way to deploy applications on the BIG-IP system using AS3.
The FAST Extension provides a toolset for templating and managing AS3 Applications on BIG-IP.

For more information about FAST, including installation and usage information, see the FAST Documentation
https://clouddocs.f5.com/products/extensions/f5-appsvcs-templates/latest/
https://github.com/F5Networks/f5-appsvcs-templates

Example Usage

::

   resource “bigip_fast_http_app” “app1” {
       application = “myApp3”
       tenant = “scenario3”
       virtual_server {
         ip = “10.1.10.223”
         port = 80
       }
     pool_members {
       addresses = [“10.1.10.120”, “10.1.10.121”, “10.1.10.122”]
       port = 80
       }
     snat_pool_address = [“10.1.10.50”, “10.1.10.51”, “10.1.10.52”]
     load_balancing_mode = “least-connections-member”
     monitor {
       send_string = “GET / HTTP/1.1\r\nHost: example.com\r\nConnection: Close\r\n\r\n” response = “200 OK”
     }
   }


Terraform integration resouces/data source
``````````````````````````````````````````
- https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_fast_http_app
- https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_fast_https_app
- https://registry.terraform.io/providers/F5Networks/bigip/latest/docs/resources/bigip_fast_tcp_app


**Additional lab guides**

.. toctree::
   :maxdepth: 2
   :includehidden:
   :glob:

   /userguide/fast-integration/create-udp.rst
   /userguide/fast-integration/create-tcp.rst
   /userguide/fast-integration/create-http.rst
   /userguide/fast-integration/create-https.rst
   /userguide/fast-integration/create-http-existing-pool.rst
   /userguide/fast-integration/create-awaf.rst
   /userguide/fast-integration/apply-canary.rst
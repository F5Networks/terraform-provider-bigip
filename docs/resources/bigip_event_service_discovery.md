---
layout: "bigip"
page_title: "BIG-IP: bigip_event_service_discovery"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_event_service_discovery resource
---

# bigip_event_service_discovery

`bigip_event_service_discovery` Terraform resource to update pool members based on event driven Service Discovery.

The API endpoint for Service discovery tasks should be available before using the resource and with this resource,we will be able to connect to a specific endpoint related to event based service discovery that will allow us to update the list of pool members


## Example Usage


```hcl
resource "bigip_event_service_discovery" "test" {
  taskid = "~Sample_event_sd~My_app~My_pool"
  node {
    id   = "newNode1"
    ip   = "192.168.2.3"
    port = 8080
  }
  node {
    id   = "newNode2"
    ip   = "192.168.2.4"
    port = 8080
  }
}
```      

## Argument Reference

* `taskid` - (Required) servicediscovery endpoint ( Below example shows how to create endpoing using AS3 )

* `node` - (Required) Map of node which will be added to pool which will be having node name(id),node address(ip) and node port(port)

For more information, please refer below document
https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/declarations/discovery.html?highlight=service%20discovery#event-driven-service-discovery

Below example shows how to use event-driven service discovery, introduced in AS3 3.9.0.

With event-driven service discovery, you POST a declaration with the addressDiscovery property set to event. This creates a new endpoint which you can use to add nodes that does not require an AS3 declaration, so it can be more efficient than using PATCH or POST to add nodes. 

When you use the event value for addressDiscovery, the system creates the new endpoint with the following syntax: https://<host>/mgmt/shared/service-discovery/task/~<tenant name>~<application name>~<pool name>/nodes.

For example, in the following declaration, assuming 192.0.2.14 is our BIG-IP, the endpoint that is created is: https://192.0.2.14/mgmt/shared/service-discovery/task/~Sample_event_sd~My_app~My_pool/nodes

Once the endpoint is created( taskid ), you can use it to add nodes to the BIG-IP pool
First we show the initial declaration to POST to the BIG-IP system.

{
    "class": "ADC",
    "schemaVersion": "3.9.0",
    "id": "Pool",
    "Sample_event_sd": {
        "class": "Tenant",
        "My_app": {
            "class": "Application",
            "My_pool": {
                "class": "Pool",
                "members": [
                    {
                        "servicePort": 8080,
                        "addressDiscovery": "static",
                        "serverAddresses": [
                            "192.0.2.2"
                        ]
                    },
                    {
                        "servicePort": 8080,
                        "addressDiscovery": "event"
                    }
                ]
            }
        }
    }
}


Once the declaration has been sent to the BIG-IP, we can use taskid/id ( ~Sample_event_sd~My_app~My_pool" ) and node list for the resource to dynamically update the node list.



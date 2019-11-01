---
layout: "bigip"
page_title: "BIG-IP: bigip_do"
sidebar_current: "docs-bigip-resource-x"
description: |-
    Provides details about bigip do resource
---

# bigip_do

`bigip_do` provides details about bigip do resource

This resource is helpful to configure deploy do declarative JSON on BIG-IP.
## Example Usage


```hcl

resource "bigip_do"  "do-example" {
     do_json = "${file("example.json")}"
     tenant_name = "Common"
 }

```

## Argument Reference


* `do_json` - (Required) Name of the of the Declarative DO JSON file

* `tenant_name` - (Required) This is the partition name where the application services will be configured.

* `example.json` - Example of DO Declarative JSON

```hcl
{
    "schemaVersion": "1.0.0",
    "class": "Device",
    "async": true,
    "webhook": "https://example.com/myHook",
    "label": "my BIG-IP declaration for declarative onboarding",
    "Common": {
        "class": "Tenant",
        "hostname": "bigip.example.com",
        "myLicense": {
            "class": "License",
            "licenseType": "regKey",
            "regKey": "AAAAA-BBBBB-CCCCC-DDDDD-EEEEEEE"
        },
        "myDns": {
            "class": "DNS",
            "nameServers": [
                "8.8.8.8",
                "2001:4860:4860::8844"
            ],
            "search": [
                "f5.com"
            ]
        },
        "myNtp": {
            "class": "NTP",
            "servers": [
                "0.pool.ntp.org",
                "1.pool.ntp.org",
                "2.pool.ntp.org"
            ],
            "timezone": "UTC"
        },
        "root": {
            "class": "User",
            "userType": "root",
            "oldPassword": "default",
            "newPassword": "myNewPass1word"
        },
        "admin": {
            "class": "User",
            "userType": "regular",
            "password": "asdfjkl",
            "shell": "bash"
        },
        "guestUser": {
            "class": "User",
            "userType": "regular",
            "password": "guestNewPass1",
            "partitionAccess": {
                "Common": {
                    "role": "guest"
                }
            }
        },
        "anotherUser": {
            "class": "User",
            "userType": "regular",
            "password": "myPass1word",
            "shell": "none",
            "partitionAccess": {
                "all-partitions": {
                    "role": "guest"
                }
            }
        },
        "myProvisioning": {
            "class": "Provision",
            "ltm": "nominal",
            "gtm": "minimum"
        },
        "internal": {
            "class": "VLAN",
            "tag": 4093,
            "mtu": 1500,
            "interfaces": [
                {
                    "name": "1.2",
                    "tagged": true
                }
            ],
            "cmpHash": "dst-ip"
        },
        "internal-self": {
            "class": "SelfIp",
            "address": "10.10.0.100/24",
            "vlan": "internal",
            "allowService": "default",
            "trafficGroup": "traffic-group-local-only"
        },
        "external": {
            "class": "VLAN",
            "tag": 4094,
            "mtu": 1500,
            "interfaces": [
                {
                    "name": "1.1",
                    "tagged": true
                }
            ],
            "cmpHash": "src-ip"
        },
        "external-self": {
            "class": "SelfIp",
            "address": "10.20.0.100/24",
            "vlan": "external",
            "allowService": "none",
            "trafficGroup": "traffic-group-local-only"
        },
        "default": {
            "class": "Route",
            "gw": "10.10.0.1",
            "network": "default",
            "mtu": 1500
        },
        "managementRoute": {
            "class": "ManagementRoute",
            "gw": "1.2.3.4",
            "network": "default",
            "mtu": 1500
        },
        "myRouteDomain": {
            "class": "RouteDomain",
            "id": 100,
            "bandWidthControllerPolicy": "bwcPol",
            "connectionLimit": 5432991,
            "flowEvictionPolicy": "default-eviction-policy",
            "ipIntelligencePolicy": "ip-intelligence",
            "enforcedFirewallPolicy": "enforcedPolicy",
            "stagedFirewallPolicy": "stagedPolicy",
            "securityNatPolicy": "securityPolicy",
            "servicePolicy": "servicePolicy",
            "strict": false,
            "routingProtocols": [
                "RIP"
            ],
            "vlans": [
                "external"
            ]
        },
        "dbvars": {
        	"class": "DbVariables",
        	"ui.advisory.enabled": true,
        	"ui.advisory.color": "green",
        	"ui.advisory.text": "/Common/hostname"
        }
    }
}
```
* `DO documentation` - https://clouddocs.f5.com/products/extensions/f5-declarative-onboarding/latest/bigip-examples.html#standalone-declaration

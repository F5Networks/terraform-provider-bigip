---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_iapp"
sidebar_current: "docs-bigip-resource-iapp-x"
description: |-
    Provides details about bigip  iapp resource for BIG-IP
---

# bigip\_iapp

`bigip_sys_iapp` resource helps you to deploy Application Services template that can be used to automate and orchestrate Layer 4-7 applications service deployments using F5 Network. More information on iApp 2.0 is at https://devcentral.f5.com/wiki/iApp.AppSvcsiApp_userguide_userguide.ashx

## Example Usage


```hcl
 resource "bigip_sys_iapp" "waf_asm" {
  name = "policywaf"
  jsonfile = "${file("policywaf.json")}"
}
```

## Argument Reference

* `name` -  Name of the iApp.

* `jsonfile` - Refer to the Json file which will be deployed on F5 BIG-IP.



## Example Usage of Json file

{
 "name":"policywaf",
  "partition": "Common",
  "inheritedDevicegroup": "true",
  "inheritedTrafficGroup": "true",
  "strictUpdates": "enabled",
  "template": "/Common/appsvcs_integration_v2.0.003",
  "execute-action": "definition",
        "tables": [{
                        "name": "feature__easyL4FirewallBlacklist",
                        "columnNames": [
                                "CIDRRange"
                        ],
                        "rows": [

                        ]
                },
                {
                        "name": "feature__easyL4FirewallSourceList",
                        "columnNames": [
                                "CIDRRange"
                        ],
                        "rows": [{
                                "row": [
                                        "0.0.0.0/0"
                                ]
                        }]
                },
                {
                        "name": "l7policy__rulesAction",
                        "columnNames": [
                                "Group",
                                "Target",
                                "Parameter"
                        ],
                        "rows": [
                                {"row": ["0", "asm/request/enable/policy", "/Common/Demo"]},
                                {"row": ["0", "forward/request/select/pool", "pool:0"]},
                                {"row": ["default", "forward/request/select/pool", "pool:0"]}
                        ]
                },
                {
                        "name": "l7policy__rulesMatch",
                        "columnNames": [
                                "Group",
                                "Operand",
                                "Negate",
                                "Condition",
                                "Value",
                                "CaseSensitive",
                                "Missing"
                        ],
                        "rows": [
                                {"row": ["0","http-uri/request/path","no","equals","/","no","no"]},
                                {"row": ["default","","no","equals","","no","no"]}
                        ]
                },
                {
                        "name": "monitor__Monitors",
                        "columnNames": [
                                "Index",
                                "Name",
                                "Type",
                                "Options"
                        ],
                        "rows": [{
                                "row": [
                                        "0",
                                        "/Common/http",
                                        "none",
                                        "none"
                                ]
                        }]
                },
                {
                        "name": "pool__Members",
                        "columnNames": [
                                "Index",
                                "IPAddress",
                                "Port",
                                "ConnectionLimit",
                                "Ratio",
                                "PriorityGroup",
                                "State",
                                "AdvOptions"
                        ],
                        "rows": [
                                {"row": ["0","192.168.69.140","80","0","1","0","enabled","none"]},
                                {"row": ["0","192.168.69.141","80","0","1","0","enabled","none"]},
                                {"row": ["0","192.168.68.142","80","0","1","0","enabled","none"]},
                                {"row": ["0","192.168.68.143","80","0","1","0","enabled","none"]},
                                {"row": ["0","192.168.68.144","80","0","1","0","enabled","none"]}
                        ]
                },
                {
                        "name": "pool__Pools",
                        "columnNames": [
                                "Index",
                                "Name",
                                "Description",
                                "LbMethod",
                                "Monitor",
                                "AdvOptions"
                        ],
                        "rows": [{
                                "row": [
                                        "0",
                                        "",
                                        "",
                                        "round-robin",
                                        "0",
                                        "none"
                                ]
                        }]
                },
                {
                        "name": "vs__BundledItems",
                        "columnNames": [
                                "Resource"
                        ],
                        "rows": [

                        ]
                },
                {
                        "name": "vs__Listeners",
                        "columnNames": [
                                "Listener",
                                "Destination"
                        ],
                        "rows": [

                        ]
                }
        ],
        "variables": [{
                        "name": "extensions__Field1",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "extensions__Field2",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "extensions__Field3",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "feature__easyL4Firewall",
                        "encrypted": "no",
                        "value": "auto"
                },
                {
                        "name": "feature__insertXForwardedFor",
                        "encrypted": "no",
                        "value": "auto"
                },
                {
                        "name": "feature__redirectToHTTPS",
                        "encrypted": "no",
                        "value": "auto"
                },
                {
                        "name": "feature__securityEnableHSTS",
                        "encrypted": "no",
                        "value": "disabled"
                },
                {
                        "name": "feature__sslEasyCipher",
                        "encrypted": "no",
                        "value": "disabled"
                },
                {
                        "name": "feature__statsHTTP",
                        "encrypted": "no",
                        "value": "auto"
                },
                {
                        "name": "feature__statsTLS",
                        "encrypted": "no",
                        "value": "auto"
                },
                {
                        "name": "iapp__apmDeployMode",
                        "encrypted": "no",
                        "value": "preserve-bypass"
                },
                {
                        "name": "iapp__appStats",
                        "encrypted": "no",
                        "value": "enabled"
                },
                {
                        "name": "iapp__asmDeployMode",
                        "encrypted": "no",
                        "value": "preserve-bypass"
                },
                {
                        "name": "iapp__logLevel",
                        "encrypted": "no",
                        "value": "7"
                },
                {
                        "name": "iapp__mode",
                        "encrypted": "no",
                        "value": "auto"
                },
                {
                        "name": "iapp__routeDomain",
                        "encrypted": "no",
                        "value": "auto"
                },
                {
                        "name": "iapp__strictUpdates",
                        "encrypted": "no",
                        "value": "enabled"
                },
                {
                        "name": "l7policy__defaultASM",
                        "encrypted": "no",
                        "value": "bypass"
                },
                {
                        "name": "l7policy__defaultL7DOS",
                        "encrypted": "no",
                        "value": "bypass"
                },
                {
                        "name": "l7policy__strategy",
                        "encrypted": "no",
                        "value": "/Common/first-match"
                },
                {
                        "name": "pool__DefaultPoolIndex",
                        "encrypted": "no",
                        "value": "0"
                },
                {
                        "name": "pool__MemberDefaultPort",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "pool__addr",
                        "encrypted": "no",
                        "value": "10.168.68.100"
                },
                {
                        "name": "pool__mask",
                        "encrypted": "no",
                        "value": "255.255.255.255"
                },
                {
                        "name": "pool__port",
                        "encrypted": "no",
                        "value": "80"
                },
                {
                        "name": "vs__AdvOptions",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__AdvPolicies",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__AdvProfiles",
                        "value": "/Common/websecurity",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ConnectionLimit",
                        "encrypted": "no",
                        "value": "0"
                },
                {
                        "name": "vs__Description",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__IpProtocol",
                        "encrypted": "no",
                        "value": "tcp"
                },
                {
                        "name": "vs__Irules",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__Name",
                        "encrypted": "no",
                        "value": "VS_80"
                },
                {
                        "name": "vs__OptionConnectionMirroring",
                        "encrypted": "no",
                        "value": "disabled"
                },
                {
                        "name": "vs__OptionSourcePort",
                        "encrypted": "no",
                        "value": "preserve"
                },
                {
                        "name": "vs__ProfileAccess",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileAnalytics",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileClientProtocol",
                        "encrypted": "no",
                        "value": "/Common/tcp"
                },
                {
                        "name": "vs__ProfileClientSSL",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileClientSSLAdvOptions",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileClientSSLCert",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileClientSSLChain",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileClientSSLCipherString",
                        "encrypted": "no",
                        "value": "DEFAULT"
                },
                {
                        "name": "vs__ProfileClientSSLKey",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileCompression",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileConnectivity",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileDefaultPersist",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileFallbackPersist",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileHTTP",
                        "value": "/Common/http",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileOneConnect",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfilePerRequest",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileRequestLogging",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileSecurityDoS",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileSecurityIPBlacklist",
                        "encrypted": "no",
                        "value": "none"
                },
                {
                        "name": "vs__ProfileSecurityLogProfiles",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileServerProtocol",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__ProfileServerSSL",
                        "value": "",
                        "encrypted": "no"
                },
                {
                        "name": "vs__RouteAdv",
                        "encrypted": "no",
                        "value": "disabled"
                },
                {
                        "name": "vs__SNATConfig",
                        "encrypted": "no",
                        "value": "automap"
                },
                {
                        "name": "vs__SourceAddress",
                        "encrypted": "no",
                        "value": "0.0.0.0/0"
                },
                {
                        "name": "vs__VirtualAddrAdvOptions",
                        "value": "",
                        "encrypted": "no"
                }
        ]
}


 * `description` - User defined description.
 * `deviceGroup` - The name of the device group that the application service is assigned to.
 * `executeAction` - Run the specified template action associated with the application.
 * `inheritedDevicegroup`- Read-only. Shows whether the application folder will automatically remain with the same device-group as its parent folder. Use 'device-group default' or 'device-group non-default' to set this.
 * `inheritedTrafficGroup` - Read-only. Shows whether the application folder will automatically remain with the same traffic-group as its parent folder. Use 'traffic-group default' or 'traffic-group non-default' to set this.
 * `partition` - Displays the administrative partition within which the application resides.
 * `strictUpdates` - Specifies whether configuration objects contained in the application may be directly modified, outside the context of the system's application management interfaces.
 * `template` - The template defines the configuration for the application. This may be changed after the application has been created to move the application to a new template.
 * `templateModified` - Indicates that the application template used to deploy the application has been modified. The application should be updated to make use of the latest changes.
 * `templatePrerequisiteErrors` - Indicates any missing prerequisites associated with the template that defines this application.
 * `trafficGroup` - The name of the traffic group that the application service is assigned to.
 * `lists` - string values
 * `metadata` - User defined generic data for the application service. It is a name and value pair.
 * `tables` - Values provided like pool name, nodes etc.
 * `variables` - Name, values, encrypted or not 

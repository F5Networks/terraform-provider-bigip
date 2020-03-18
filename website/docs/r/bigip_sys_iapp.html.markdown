---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_iapp"
sidebar_current: "docs-bigip-resource-iapp-x"
description: |-
    Provides details about bigip  iapp resource for BIG-IP
---

# bigip\_sys\_iapp

`bigip_sys_iapp` resource helps you to deploy Application Services template that can be used to automate and orchestrate Layer 4-7 applications service deployments using F5 Network.  

## Example Usage


```hcl
 resource "bigip_sys_iapp" "simplehttp" {
   name     = "simplehttp"
   jsonfile = file("simplehttp.json")
 }
```

## Argument Reference

* `name` -  Name of the iApp.

* `jsonfile` - Refer to the Json file which will be deployed on F5 BIG-IP.



## Example Usage of Json file
```
{  
"fullPath":"/Common/simplehttp.app/simplehttp",
"generation":222,
"inheritedDevicegroup":"true",
"inheritedTrafficGroup":"true",
"kind":"tm:sys:application:service:servicestate",
"name":"simplehttp",
"partition":"Common",
"selfLink":"https://localhost/mgmt/tm/sys/application/service/~Common~simplehttp.app~simplehttp?ver=13.0.0",
"strictUpdates":"enabled",
"subPath":"simplehttp.app",
"tables":[  
   {  
      "name":"basic__snatpool_members"
   },
   {  
      "name":"net__snatpool_members"
   },
   {  
      "name":"optimizations__hosts"
   },
   {  
      "columnNames":[  
         "name"
      ],
      "name":"pool__hosts",
      "rows":[  
         {  
            "row":[  
               "f5.cisco.com"
            ]
         }
      ]
   },
   {  
      "columnNames":[  
         "addr",
         "port",
         "connection_limit"
      ],
      "name":"pool__members",
      "rows":[  
         {  
            "row":[  
               "10.0.2.167",
               "80",
               "0"
            ]
         },
         {  
            "row":[  
               "10.0.2.168",
               "80",
               "0"
            ]
         }
      ]
   },
   {  
      "name":"server_pools__servers"
   }
],
"template":"/Common/f5.http",
"templateModified":"no",
"templateReference":{  
   "link":"https://localhost/mgmt/tm/sys/application/template/~Common~f5.http?ver=13.0.0"
},
"trafficGroup":"/Common/traffic-group-1",
"trafficGroupReference":{  
   "link":"https://localhost/mgmt/tm/cm/traffic-group/~Common~traffic-group-1?ver=13.0.0"
},
"variables":[  
   {  
      "encrypted":"no",
      "name":"client__http_compression",
      "value":"/#create_new#"
   },
   {  
      "encrypted":"no",
      "name":"monitor__monitor",
      "value":"/Common/http"
   },
   {  
      "encrypted":"no",
      "name":"net__client_mode",
      "value":"wan"
   },
   {  
      "encrypted":"no",
      "name":"net__server_mode",
      "value":"lan"
   },
   {  
      "encrypted":"no",
      "name":"net__v13_tcp",
      "value":"warn"
   },
   {  
      "encrypted":"no",
      "name":"pool__addr",
      "value":"10.0.1.100"
   },
   {  
      "encrypted":"no",
      "name":"pool__pool_to_use",
      "value":"/#create_new#"
   },
   {  
      "encrypted":"no",
      "name":"pool__port",
      "value":"80"
   },
   {  
      "encrypted":"no",
      "name":"ssl__mode",
      "value":"no_ssl"
   },
   {  
      "encrypted":"no",
      "name":"ssl_encryption_questions__advanced",
      "value":"no"
   },
   {  
      "encrypted":"no",
      "name":"ssl_encryption_questions__help",
      "value":"hide"
   }
]
}
```

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

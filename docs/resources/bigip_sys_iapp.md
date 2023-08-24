---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_iapp"
subcategory: "System"
description: |-
  Provides details about bigip_sys_iapp resource
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

* `name` - (Required, type `string`) Name of the iApp.
* `jsonfile` - (Required, type `string`) Refer to the Json file which will be deployed on F5 BIG-IP.
* `description` - (Optional type `string`) - User defined description.
* `partition` - (Optional type `string`) - Displays the administrative partition within which the application resides.
* `execute_action` - (Optional type `string`) - Run the specified template action associated with the application, this option can be specified in `json` with `executeAction`, value specified with `execute_action` attribute take precedence over `json` value

## Example Usage of Json file
```json
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

## Attribute Reference

* `inherited_devicegroup` - Read-only. Shows whether the application folder will automatically remain with the same device-group as its parent folder. Use 'device-group default' or 'device-group non-default' to set this.
* `inherited_traffic_group` - Read-only. Shows whether the application folder will automatically remain with the same traffic-group as its parent folder. Use 'traffic-group default' or 'traffic-group non-default' to set this.
* `strict_updates` - Specifies whether configuration objects contained in the application may be directly modified, outside the context of the system's application management interfaces.
* `template` - The template defines the configuration for the application. This may be changed after the application has been created to move the application to a new template.
* `template_modified` - Indicates that the application template used to deploy the application has been modified. The application should be updated to make use of the latest changes.
* `template_prerequisite_errors` - Indicates any missing prerequisites associated with the template that defines this application.
* `traffic_group` - The name of the traffic group that the application service is assigned to.
* `lists` - string values
* `metadata` - User defined generic data for the application service. It is a name and value pair.

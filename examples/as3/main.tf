#
#resource "bigip_as3" "as3-example1" {
#  as3_json = file("../examples/as3/as3_example1.json")
#}
#
#resource "bigip_as3" "as3-example2" {
#  as3_json = file("../examples/as3/as3_example2.json")
#}

resource "bigip_as3" "app-as3-irule" {
  as3_json = <<EOT
{
  "class": "AS3",
  "action": "deploy",
  "persist": true,
  "declaration": {
    "class": "ADC",
    "schemaVersion": "3.0.0",
    "id": "urn:uuid:a858e55e-bbe6-42ce-a9b9-0f4ab33e3bf7",
    "label": "Sample 2",
    "remark": "HTTP with no compression, BIG-IP tcp profile, iRule for pool",
    "Sample_http_02": {
      "class": "Tenant",
      "A1": {
        "class": "Application",
        "service": {
          "class": "Service_HTTP",
          "virtualAddresses": [
            "10.0.3.10"
          ],
          "pool": "dfl_pool",
          "profileHTTPCompression": "basic",
          "iRules": [
            "choose_pool"
          ],
          "profileTCP": {
            "bigip": "/Common/mptcp-mobile-optimized"
          }
        },
        "dfl_pool": {
          "class": "Pool",
          "monitors": [
            "http"
          ],
          "members": [{
            "servicePort": 80,
            "serverAddresses": [
              "192.0.3.10",
              "192.0.3.11"
            ]
          }]
        },
        "pvt_pool": {
          "class": "Pool",
          "monitors": [
            "http"
          ],
          "members": [{
            "servicePort": 80,
            "serverAddresses": [
              "192.0.3.20",
              "192.0.3.21"
            ]
          }]
        },
        "choose_pool": {
          "class": "iRule",
          "remark": "choose private pool based on IP",
          "iRule": "when CLIENT_ACCEPTED {\nif {[class match [IP::remote_addr] equals IPs ]} {\n pool `*pvt_pool`\n }\n}"
        }
      }
    }
  }
}

EOT
}


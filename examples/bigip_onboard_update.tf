resource "bigip_do" "do-example1" {
  do_json = <<EOT
  {
    "schemaVersion": "1.38.0",
    "class": "Device",
    "async": true,
    "label": "my BIG-IP declaration for declarative onboarding",
    "Common": {
        "class": "Tenant",
        "hostname": "ecosyshyd-bigip03.com",
        "guestUser": {
            "class": "User",
            "userType": "regular",
            "partitionAccess": {
                "Common": {
                    "role": "guest"
                }
            },
            "shell": "tmsh"
        },
        "anotherUser": {
            "class": "User",
            "userType": "regular",
            "shell": "none",
            "partitionAccess": {
                "all-partitions": {
                    "role": "guest"
                }
            }
        },
        "dbvars": {
            "class": "DbVariables",
            "ui.advisory.enabled": true,
            "ui.advisory.color": "green",
            "ui.advisory.text": "/Common/hostname"
        }
    }
}

EOT
}

output "do_json" {
  value = bigip_do.do-example1.do_json
}
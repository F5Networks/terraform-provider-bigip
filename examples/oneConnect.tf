provider "bigip" {
  address = "10.192.74.61"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_oneconnect" "oneconnect-sanjose"

        {  
            name = "sanjose"
            partition = "Common"
            defaults_from = "/Common/oneconnect"
            idler_timeout_override = "disabled"
            max_age = 3600
            max_reuse = 1000
            max_size = 1000
            sharer_pools = "disabled"
            source_mask = "255.255.255.255"
        }



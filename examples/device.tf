provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_cm_device" "my_new_device"

        {
            name = "bigip300.f5.com"
            configsync_ip = "2.2.2.2"
            mirror_ip = "10.10.10.10"
            mirror_secondary_ip = "11.11.11.11"
        }
 

provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_device" "my_new_device"

        {
            name = "bigip200.f5.com"
            configsyncIp = "2.2.2.2"
            mirrorIp = "10.10.10.10"
            mirrorSecondaryIp = "11.11.11.11"
        }

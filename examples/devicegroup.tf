provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_devicegroup" "my_new_devicegroup"

        {
            name = "bigip20r10.f5.com"
            autoSync = "enabled"
            fullLoadOnSync = "true"
            type = "sync-only"
        }

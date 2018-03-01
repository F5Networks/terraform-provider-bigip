provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_cm_devicegroup" "my_new_devicegroup"

        {
            name = "deadlygroup"
            auto_sync = "enabled"
            full_load_on_sync = "true"
            type = "sync-only"
            device  { name = "bigip1.cisco.com"}
            device  { name = "bigip200.f5.com"}
        }

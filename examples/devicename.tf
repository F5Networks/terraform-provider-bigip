provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_device_name" "my_new_device_in_LA"

        {
            name = "bigip1i00.f5.com"
            command = "mv"
            target = "bigipNew.f5.com"
        }

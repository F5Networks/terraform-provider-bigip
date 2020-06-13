provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxx"
}

resource "bigip_ssl_key" "test-cert" {
  name      = "serverkey.key"
  content   = file("serverkey.key")
  partition = "Common"
}


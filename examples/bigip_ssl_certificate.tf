provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxxx"
}

resource "bigip_ssl_certificate" "test-cert" {
        name = "servercert"
        cert_path = "servercert.crt"
        partition = "Common"
}

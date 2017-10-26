provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_http2_profile" "nyhttp2"

        {
            name = "/Common/NewYork_http2"
            defaults_from = "/Common/http2"
            concurrent_streams_per_connection = 10
            connection_idle_timeout= 30
            activation_modes = ["alpn","npn"]
        }

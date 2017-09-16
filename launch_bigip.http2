provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_http2_profile" "nyhttp2"

        {
            name = "/Common/NewYork_http2"
            defaultsFrom = "/Common/http2"
            concurrentStreamsPerConnection = 10
            connectionIdleTimeout= 30
            activationModes = ["alpn","npn"]
        }

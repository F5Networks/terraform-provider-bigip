provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_profile_httpcompress" "sjhttpcompression"

        {
            name = "/Common/sjhttpcompression2"
            defaults_from = "/Common/httpcompression"
            uri_exclude   = ["www.abc.f5.com", "www.abc2.f5.com"]
            uri_include   = ["www.xyzbc.cisco.com"]
        }

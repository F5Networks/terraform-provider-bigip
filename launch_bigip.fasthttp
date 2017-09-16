provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_fasthttp_profile" "sjfasthttpprofile"

        {
            name = "sjfasthttpprofile"
            defaults_from = "/Common/fasthttp"
            idle_timeout = 300
            connpoolidle_timeoutoverride	= 0
            connpool_maxreuse = 2
            connpool_maxsize  = 2048
            connpool_minsize = 0
            connpool_replenish = "enabled"
            connpool_step = 4
            forcehttp_10response = "disabled"
            maxheader_size = 32768
      }

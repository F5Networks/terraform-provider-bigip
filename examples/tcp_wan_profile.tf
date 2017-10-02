provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_tcp_profile" "sanjose-tcp-wan-profile"

        {  
            name = "sanjose-tcp-wan-profile"
            defaultsFrom = "/Common/tcp-wan-optimized"
            idleTimeout = 300
            closeWaitTimeout = 5
            finWait_2Timeout = 5
            finWaitTimeout = 300
            keepAliveInterval = 1700
            deferredAccept = "enabled"
            fastOpen = "enabled"
        }



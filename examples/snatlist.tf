provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_ltm_snat" "snat_list" {
 name = "NewSnatList"
 translation = "136.1.1.1"
 origins =  {name = "2.2.2.2"}
 origins =  {name = "3.3.3.3"} 
}


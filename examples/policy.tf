provider "bigip" {
  address = "10.192.74.61"
  username = "admin"
  password = "admin"
}

resource "bigip_ltm_policy" "test-policy" {
	name = "/Common/Drafts/newp"
        strategy = "/Common/first-match"
	controls = ["forwarding"]
	requires = ["http"]
	rule {
		name = "/Common/rulepolicy"
			}
}




resource "bigip_ltm_profile_rewrite" "tftest" {
  name          = "/Common/tf_profile"
  defaults_from = "/Common/rewrite"
  rewrite_mode  = "uri-translation"
}

resource "bigip_ltm_profile_rewrite_uri_rules" "tftestrule1" {
  profile_name = bigip_ltm_profile_rewrite.tftest.name
  rule_name    = "tf_rule"
  rule_type    = "request"
  client {
    host   = "www.foo.com"
    scheme = "https"
  }
  server {
    host   = "www.bar.com"
    path   = "/this/"
    scheme = "https"
    port   = "8888"
  }
}

resource "bigip_ltm_profile_rewrite_uri_rules" "tftestrule2" {
  profile_name = bigip_ltm_profile_rewrite.tftest.name
  rule_name    = "tf_rule2"
  client {
    host   = "www.baz.com"
    path   = "/that/"
    scheme = "ftp"
    port   = "8888"
  }
  server {
    host   = "www.buz.com"
    path   = "/those/"
    scheme = "ftps"
  }
}

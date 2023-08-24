resource "bigip_waf_policy" "test-APPLANG01" {
  name               = "testapplang01"
  template_name      = "POLICY_TEMPLATE_FUNDAMENTAL"
  policy_import_json = file("/Users/r.chinthalapalli/Downloads/POLICY_IMPORT_APPLANG.json")
}

resource "bigip_waf_policy" "test-APPLANG02" {
  name                 = "testapplang02"
  template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
  application_language = "utf-8"
  policy_import_json   = file("/Users/r.chinthalapalli/Downloads/POLICY_IMPORT_APPLANG.json")
}

resource "bigip_waf_policy" "test-APPLANG03" {
  name          = "testapplang03"
  template_name = "POLICY_TEMPLATE_FUNDAMENTAL"
}

resource "bigip_waf_policy" "test-APPLANG04" {
  name                 = "testapplang04"
  template_name        = "POLICY_TEMPLATE_FUNDAMENTAL"
  application_language = "utf-8"
}
# This is the As3 object model for terraform ( Another way to post As3 json to bigip ).
# This model will construct AS3 json from the user inputs and post it to bigip
# In this model top As3 class is defined as resource "bigip_as3_class" and rest as data sources.
# resource "bigip_as3_class" uses data source "bigip_as3_adc" which in turn uses other data sources.
# Below is the flow of code between data sources and resource
# 
# 
# For "bigip_as3_app" https declartion, we need pool,service,cert,tls_server to attach it,hence it takes input from those data sources.
# All the App declaration will be logically moved under tenant, hence "bigip_as3_tenant" will consume app data source
# All the tenant decalaration will go under ADC class, hence "bigip_as3_adc" will consume tenant data source
# Finally As3 class ( "bigip_as3_class" ) will be top class in AS3 declaration which will consume ADC data source.
#
# This is as per AS3 schema ( https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/composing-a-declaration.html )
# For more info Please refer to F5 cloud docs
provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}
data "bigip_as3_pool" "mydataas3pool" {
  name = "web_pool3"
  loadbalancing_mode = "round-robin"
  servicedown_action = "none"
  pool_members {
    connection_limit = 10
    rate_limit=10
    dynamic_ratio=100
    service_port=8080
    ratio=90
    priority_group=30
    sharenodes=true
    adminstate="enable"
    #address_discovery="enabled"
    server_addresses=["192.168.30.1","192.168.25.1"]
  }
  minimummembers_active=1
  reselect_tries=0
  slowramp_time=10
  minimum_monitors=1
  monitors=["http"]
}
data "bigip_as3_service" "myservice" {
  name = "serviceMain"
  virtual_addresses=["10.0.10.10"]
  pool_name = "${data.bigip_as3_pool.mydataas3pool.name}"
  server_tls = "${data.bigip_as3_tls_server.exmpserver.name}"
  service_type = "https"
  virtual_port = 443
  persistence_methods = ["cookie"]
}
data "bigip_as3_cert" "exmpcert" {
  name = "exmpcert"
  remark = "in practice we recommend using a passphrase"
  certificate = "${file("servercert.crt")}"
  private_key = "${file("serverkey.key")}"
  passphrase {
    ciphertext = "ZjVmNQ=="
    protected = "eyJhbGciOiJkaXIiLCJlbmMiOiJub25lIn0"
  }
}
data "bigip_as3_tls_server" "exmpserver" {
  name = "exmpserver"
  certificates {
    certificate = "exmpcert"
  }
}

data "bigip_as3_app" "App1" {
  name = "App1"
  template = "https"
  pool_class = "${data.bigip_as3_pool.mydataas3pool.id}"
  service_class = "${data.bigip_as3_service.myservice.id}"
  cert_class = "${data.bigip_as3_cert.exmpcert.id}"
  tls_server_class = "${data.bigip_as3_tls_server.exmpserver.id}"
  enable = true
}
data "bigip_as3_tenant" "sample"{
  name = "sample_01"
  app_class_list = ["${data.bigip_as3_app.App1.id}"]
  defaultroutedomain = 0
  enable = true
  label = "this is label for tenant"
  optimisticlockkey = "dfghj"
  remark = "dfghjk"
}

data "bigip_as3_adc" "exmpadc"{
  name = "exmpadc"
  //label = "asdfghj"
  tenant_class_list = ["${data.bigip_as3_tenant.sample.id}"]
}
resource "bigip_as3_class" "as3-example" {
  name = "as3-example"
  declaration="${data.bigip_as3_adc.exmpadc.id}"
  tenants = ["sample_01"]
}

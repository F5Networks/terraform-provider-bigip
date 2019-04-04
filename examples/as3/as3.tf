//We can use null_resource to  deploy As3 templates, below is simple example to install the as3 rpm and another resource which deploys the example1.json ( which has the http VS configuration) More details on As3 please refer to https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/


provider "bigip" {
  address = "X.X.X.X"
  username = "admin"
  password = "pass"
}


// As3 using terraform provisoner
resource "null_resource" "install_as3" {
  provisioner "local-exec" {
    command = "sh install_as3.sh X.X.X.X  admin:pass f5-appsvcs-3.9.0-3.noarch.rpm"
  }
}
 resource "null_resource" "deploy_as3_http" {
 provisioner "local-exec" {
    command = "sh as3_http.sh"
  }
depends_on = ["null_resource.install_as3"]

}


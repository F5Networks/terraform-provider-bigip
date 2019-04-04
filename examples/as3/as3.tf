provider "bigip" {
  address = "54.177.207.51"
  username = "admin"
  password = "1LStm3545"
}


// As3 using terraform provisoner
resource "null_resource" "install_as3" {
  provisioner "local-exec" {
    command = "sh install_as3.sh 54.177.207.51 admin:1LStm3545 f5-appsvcs-3.9.0-3.noarch.rpm"
  }
}
 resource "null_resource" "deploy_as3_http" {
 provisioner "local-exec" {
    command = "sh as3_http.sh"
  }
depends_on = ["null_resource.install_as3"]

}


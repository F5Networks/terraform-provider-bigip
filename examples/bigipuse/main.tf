variable "address" {}

provider bigip {
 address = "$address"
 username = "admin"
 password = "admin"
}

terraform {
  required_providers {
    bigip = {
      source  = "F5Networks/bigip"
      version = "1.7.0"
    }
  }
}

provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxx"
}


resource "bigip_fast_template" "foo-template" {
  name   = "foo_template"
  source = "foo_template.zip"
}



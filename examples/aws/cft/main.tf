/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "aws" {
  region = "us-east-1"
}

resource "aws_cloudformation_stack"  "network" {
 name = "networking-stack"
  parameters {
   sshKey = "xxx"
   availabilityZone1 = "us-east-1a"
   adminPassword = "cisco123"
   imageName = "Best1000Mbps"
   instanceType = "m3.2xlarge"
   managementGuiPort =  "8443"
}
 template_body = "${file("cft.json")}"

}

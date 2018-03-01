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

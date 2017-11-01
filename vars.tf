variable "AWS_REGION" {
  default =  "us-east-1"
}
variable "PATH_TO_PRIVATE_KEY" {
  default = "mykey"
}
variable "PATH_TO_PUBLIC_KEY" {
  default = "mykey.pub"
}
variable "AMIS" {
  type = "map"
  default = {
    us-east-1 = "ami-8f007b98"
  }
}

variable "availabilty_zone" {
  default = "us-east-1a"
} 


variable "instance_type" {
  description = "AWS instance type"
  default = "m4.xlarge"
}







variable "username" {
  default =  "admin"
}

variable "password" {
  default = "admin"
}


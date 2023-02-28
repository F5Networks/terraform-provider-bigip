/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
variable "public_key_path" {
  default = "f5key.pub"
}

variable "public_key_path1" {
  default = "server1.pub"
}

variable "public_key_path2" {
  default = "server2.pub"
}

variable "key_name2" {
  default = "server2"
}

variable "key_name1" {
  default = "server1"
}

variable "key_name" {
  default = "f5key"
}

variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "us-east-1"
}

variable "availabilty_zone" {
  default = "us-east-1a"
}

# F5 Networks Hourly BIGIP-12.1.1.1.0.196 - Better 25Mbps - built on Sep 07 20-6f7c56e1-c69f-4c47-9659-e26e27406220-ami-1d31460a.3 (ami-8f007b98)
variable "aws_amis" {
  default = {
    us-east-1 = "ami-e56d4b85"
    us-east-1 = "ami-8f007b988f007b988f007b988f007b988f007b988f007b988f007b988f007b98"
    us-east-1 = "ami-9be6f38c"
  }
}

variable "instance_type" {
  description = "AWS instance type"
  default     = "m4.xlarge"
}


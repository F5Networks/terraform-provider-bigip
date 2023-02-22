/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
 */
# Specify the provider and access details
#https://www.terraform.io/docs/providers/aws/

# Create a VPC to launch our instances into
resource "aws_vpc" "abc" {
  cidr_block = "10.0.0.0/16"
  tags {
    Name = "abc"
  }
}

# Create an internet gateway to give our subnet access to the outside world
resource "aws_internet_gateway" "default" {
  vpc_id = aws_vpc.abc.id
  tags {
    Name = "default"
  }
}

# Grant the VPC internet access on its main route table
resource "aws_route" "internet_access" {
  route_table_id         = aws_vpc.abc.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.default.id
}

resource "aws_route_table_association" "route_table_external" {
  subnet_id      = aws_subnet.external.id
  route_table_id = aws_vpc.abc.main_route_table_id
}

resource "aws_route_table_association" "route_table_internal" {
  subnet_id      = aws_subnet.internal.id
  route_table_id = aws_vpc.abc.main_route_table_id
}

# Create a management subnet to launch our instances into
resource "aws_subnet" "management" {
  vpc_id                  = aws_vpc.abc.id
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true
  availability_zone       = var.availabilty_zone
  tags {
    Name = "management"
  }
}

# Create an external subnet to launch our instances into
resource "aws_subnet" "external" {
  vpc_id                  = aws_vpc.abc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = var.availabilty_zone
  tags {
    Name = "external"
  }
}

# Create an internal subnet to launch our instances into
resource "aws_subnet" "internal" {
  vpc_id                  = aws_vpc.abc.id
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = true
  availability_zone       = var.availabilty_zone
  tags {
    Name = "internal"
  }
}

resource "aws_network_interface" "external" {
  subnet_id       = aws_subnet.external.id
  private_ips     = ["10.0.1.10", "10.0.1.100"]
  security_groups = ["${aws_security_group.allow_all.id}"]
  attachment {
    instance     = aws_instance.SCS_F5.id
    device_index = 1
  }
}

resource "aws_network_interface" "internal" {
  subnet_id       = aws_subnet.internal.id
  private_ips     = ["10.0.2.10", "10.0.2.183"]
  security_groups = ["${aws_security_group.allow_all.id}"]
  attachment {
    instance     = aws_instance.SCS_F5.id
    device_index = 2
  }
}

resource "aws_eip" "eip_vip" {
  vpc                       = true
  network_interface         = aws_network_interface.external.id
  associate_with_private_ip = "10.0.1.100"
}

resource "aws_security_group" "allow_all" {
  name        = "allow_all"
  description = "Used in the terraform"
  vpc_id      = aws_vpc.abc.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

#A key pair is used to control login access to EC2 instances
resource "aws_key_pair" "auth" {
  key_name   = var.key_name
  public_key = file(var.public_key_path)
}

resource "aws_key_pair" "auth1" {
  key_name   = var.key_name1
  public_key = file(var.public_key_path1)
}
resource "aws_key_pair" "auth2" {
  key_name   = var.key_name2
  public_key = file(var.public_key_path2)
}


resource "aws_instance" "SCS_F5" {
  ami                         = "ami-03b18799f246d2695"
  instance_type               = var.instance_type
  associate_public_ip_address = true
  private_ip                  = "10.0.0.10"
  availability_zone           = aws_subnet.management.availability_zone
  subnet_id                   = aws_subnet.management.id
  security_groups             = ["${aws_security_group.allow_all.id}"]
  vpc_security_group_ids      = ["${aws_security_group.allow_all.id}"]
  user_data                   = file("userdata.sh")
  key_name                    = var.key_name
  root_block_device { delete_on_termination = true }
  tags {
    Name = "SCS_F5"
  }
}

resource "aws_instance" "SCS_appserver1" {
  ami                         = "ami-9be6f38c"
  instance_type               = "t2.micro"
  associate_public_ip_address = true
  private_ip                  = "10.0.2.167"
  availability_zone           = aws_subnet.internal.availability_zone
  subnet_id                   = aws_subnet.internal.id
  security_groups             = ["${aws_security_group.allow_all.id}"]
  vpc_security_group_ids      = ["${aws_security_group.allow_all.id}"]
  user_data                   = file("userdata_ami.sh")
  key_name                    = var.key_name1
  root_block_device { delete_on_termination = true }
  tags {
    Name = "SCS_appserver1"
  }
}

resource "aws_instance" "SCS_appserver2" {
  ami                         = "ami-9be6f38c"
  instance_type               = "t2.micro"
  associate_public_ip_address = true
  private_ip                  = "10.0.2.168"
  availability_zone           = aws_subnet.internal.availability_zone
  subnet_id                   = aws_subnet.internal.id
  security_groups             = ["${aws_security_group.allow_all.id}"]
  vpc_security_group_ids      = ["${aws_security_group.allow_all.id}"]
  user_data                   = file("userdata_ami.sh")
  key_name                    = var.key_name2
  root_block_device { delete_on_termination = true }
  tags {
    Name = "SCS_appserver2"
  }
}



output "SCS_F5public_ip" {
  value = aws_instance.SCS_F5.public_ip
}

output "SCS_F5_Virtual_Server_IP" {
  value = aws_eip.eip_vip.public_ip
}

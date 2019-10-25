
[//]: # (Copyright 2019 F5 Networks Inc.)
[//]: # (This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.)
[//]: # (If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.)
# AWS Terraform example TF files
The above TF files shows example of creating one F5 instance and 2 backend servers. It also creates
VPC, subnets external, internal and management. Once terraform apply is done the infrastructure is
created and the management IP addresses for BIG-IP and servers are displayed. BIG-IP is configured
with VIP and Pool members.

### master.tf
Does the VPC, subnet and instance configuration.

### providers.tf
has all the aws details required like availability zone etc.

### userdata.sh
Is the script which configures BIG-IP for selfIPs, credentials etc. Which is executed after the instance
is created.

### userdata_ami.sh
Is the script which configures the backend servers for apache etc.

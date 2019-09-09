# End to End workflow
deploy.sh shell script will execute two .TF files, the first TF file will manifest the instances on AWS for
f5 big-ip and App servers, doing so it will also create VPC, subnets, routing & security groups required for the VPC.
The second TF file will take management IP for big-ip which got created after the execution of first .TF file and configure other services using f5 terraform resources like bigip_sys_iapp, bigip_ltm_virtual_server etc.

## First .TF file
To manifest VPC, instances like f5 bigip, appservers, security groups in AWS
https://github.com/f5devcentral/terraform-provider-bigip/blob/master/examples/aws/master.tf

## Second .TF file
This TF file uses the bigip_management_ip which is extracted by deploy.sh script and configures additional services
using resources like bigip_sys_iapp, bigip_ltm_node, bigip_ltm_virtual_server resources.
https://github.com/f5devcentral/terraform-provider-bigip/master.tf
Both the .TF files need to be in different directories.

## Deploy shell Script
Deploy shell script located at https://github.com/f5devcentral/terraform-provider-bigip/blob/master/examples/aws/deploy.sh execute the First .TF file which is located at https://github.com/f5devcentral/terraform-provider-bigip/blob/master/examples/aws/master.tf
while doing so also create a dump file so as to extract the Public IP for bigip management and use this public IP in the Second TF file. It creates variables.tf file as shown below

variable "bigip_management_ip" {
  default = "54.166.104.85"
}
It waits around 16 mins so that bigip and app server instances are created in AWS.
Then it executes the second .TF file for services, in this example it uses bigip_sys_iapp resource to create VIP, Pools, node and monitor.

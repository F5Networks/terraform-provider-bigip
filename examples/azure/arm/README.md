# Azure ARM  Terraform example TF files
Below example shows how you can deploy f5 BIG-IP one NIC Azure Template using Terraform.

### main.tf
provider "azurerm" refers to Azure Cloud

resource "azurerm_resource_group" is required to define the resource group which wil have all the elements, like name of resource group and location of the instance. Eg its "East US" here 

resource "azurerm_template_deployment"  is required to include the name of the ARM template, refers to the resource group and the body of template followed by <<DEPLOY to indicate start of the template and at the end you need to specify or end with DEPLOY and specify demployment mode as "complete" or "incremental"


### variables.tf
Need variable file as well to define the BIG-IP username and password
variable "admin_username" {
  default = "xxx"
}
variable "admin_password" {
  default = "xxxxx"
}


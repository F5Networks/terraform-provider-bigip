provider "azurerm" {
}
resource "azurerm_resource_group" "scs" {
        name = "SCSResourceGroup"
        location = "eastus"
}

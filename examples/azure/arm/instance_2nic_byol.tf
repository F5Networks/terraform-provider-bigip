/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "azurerm" {
}
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-01"
  location = "East US"
}

resource "azurerm_template_deployment" "test" {
  name                = "acctesttemplate-01"
  resource_group_name = "${azurerm_resource_group.test.name}"

  template_body = <<DEPLOY

{
    "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json", 
    "contentVersion": "4.4.0.1", 
    "parameters": {
        "adminUsername": {
            "defaultValue": "${var.admin_username}", 
            "metadata": {
                "description": "User name for the Virtual Machine."
            }, 
            "type": "string"
        }, 
        "adminPassword": {
            "defaultValue": "${var.admin_password}",
            "metadata": {
                "description": "Password to login to the Virtual Machine."
            }, 
            "type": "securestring"
        }, 
        "dnsLabel": {
            "defaultValue": "REQUIRED", 
            "metadata": {
                "description": "Unique DNS Name for the Public IP address used to access the Virtual Machine."
            }, 
            "type": "string"
        }, 
        "instanceName": {
            "defaultValue": "scsf5vm_2nic_byol", 
            "metadata": {
                "description": "Name of the Virtual Machine."
            }, 
            "type": "string"
        }, 
        "instanceType": {
            "allowedValues": [
                "Standard_A2", 
                "Standard_A3", 
                "Standard_A4", 
                "Standard_A5", 
                "Standard_A6", 
                "Standard_A7", 
                "Standard_D2", 
                "Standard_D3", 
                "Standard_D4", 
                "Standard_D11", 
                "Standard_D12", 
                "Standard_D13", 
                "Standard_D14", 
                "Standard_DS2", 
                "Standard_DS3", 
                "Standard_DS4", 
                "Standard_DS11", 
                "Standard_DS12", 
                "Standard_DS13", 
                "Standard_DS14", 
                "Standard_D2_v2", 
                "Standard_D3_v2", 
                "Standard_D4_v2", 
                "Standard_D5_v2", 
                "Standard_D11_v2", 
                "Standard_D12_v2", 
                "Standard_D13_v2", 
                "Standard_D14_v2", 
                "Standard_D15_v2", 
                "Standard_DS2_v2", 
                "Standard_DS3_v2", 
                "Standard_DS4_v2", 
                "Standard_DS5_v2", 
                "Standard_DS11_v2", 
                "Standard_DS12_v2", 
                "Standard_DS13_v2", 
                "Standard_DS14_v2", 
                "Standard_DS15_v2", 
                "Standard_F2", 
                "Standard_F4", 
                "Standard_F8", 
                "Standard_F2S", 
                "Standard_F4S", 
                "Standard_F8S", 
                "Standard_F16S", 
                "Standard_G2", 
                "Standard_G3", 
                "Standard_G4", 
                "Standard_G5", 
                "Standard_GS2", 
                "Standard_GS3", 
                "Standard_GS4", 
                "Standard_GS5"
            ], 
            "defaultValue": "Standard_DS2_v2", 
            "metadata": {
                "description": "Azure instance size of the Virtual Machine."
            }, 
            "type": "string"
        }, 
        "imageName": {
            "allowedValues": [
                "Good", 
                "Better", 
                "Best"
            ], 
            "defaultValue": "Good", 
            "metadata": {
                "description": "F5 SKU (IMAGE) to you want to deploy. Note: The disk size of the VM will be determined based on the option you select."
            }, 
            "type": "string"
        }, 
        "bigIpVersion": {
            "allowedValues": [
                "13.1.0200", 
                "latest"
            ], 
            "defaultValue": "13.1.0200", 
            "metadata": {
                "description": "F5 BIG-IP version you want to use."
            }, 
            "type": "string"
        }, 
        "licenseKey1": {
            "defaultValue": "REQUIRED", 
            "metadata": {
                "description": "The license token for the F5 BIG-IP VE (BYOL)."
            }, 
            "type": "string"
        }, 
        "numberOfExternalIps": {
            "allowedValues": [
                0, 
                1, 
                2, 
                3, 
                4, 
                5, 
                6, 
                7, 
                8, 
                9, 
                10, 
                11, 
                12, 
                13, 
                14, 
                15, 
                16, 
                17, 
                18, 
                19, 
                20
            ], 
            "defaultValue": 1, 
            "metadata": {
                "description": "The number of public/private IP addresses you want to deploy for the application traffic (external) NIC on the BIG-IP VE to be used for virtual servers."
            }, 
            "type": "int"
        }, 
        "vnetAddressPrefix": {
            "defaultValue": "10.0", 
            "metadata": {
                "description": "The start of the CIDR block the BIG-IP VEs use when creating the Vnet and subnets.  You MUST type just the first two octets of the /16 virtual network that will be created, for example '10.0', '10.100', 192.168'."
            }, 
            "type": "string"
        }, 
        "ntpServer": {
            "defaultValue": "0.pool.ntp.org", 
            "metadata": {
                "description": "Leave the default NTP server the BIG-IP uses, or replace the default NTP server with the one you want to use."
            }, 
            "type": "string"
        }, 
        "timeZone": {
            "defaultValue": "UTC", 
            "metadata": {
                "description": "If you would like to change the time zone the BIG-IP uses, enter the time zone you want to use. This is based on the tz database found in /usr/share/zoneinfo. Example values: UTC, US/Pacific, US/Eastern, Europe/London or Asia/Singapore."
            }, 
            "type": "string"
        }, 
        "restrictedSrcAddress": {
            "defaultValue": "*", 
            "metadata": {
                "description": "This field restricts management access to a specific network or address. Enter an IP address or address range in CIDR notation, or asterisk for all sources"
            }, 
            "type": "string"
        }, 
        "tagValues": {
            "defaultValue": {
                "application": "APP", 
                "cost": "COST", 
                "environment": "ENV", 
                "group": "GROUP", 
                "owner": "OWNER"
            }, 
            "metadata": {
                "description": "Default key/value resource tags will be added to the resources in this deployment, if you would like the values to be unique adjust them as needed for each key."
            }, 
            "type": "object"
        }, 
        "allowUsageAnalytics": {
            "allowedValues": [
                "Yes", 
                "No"
            ], 
            "defaultValue": "Yes", 
            "metadata": {
                "description": "This deployment can send anonymous statistics to F5 to help us determine how to improve our solutions. If you select **No** statistics are not sent."
            }, 
            "type": "string"
        }
    }, 
    "variables": {
        "bigIpNicPortMap": {
            "1": {
                "Port": "[parameters('bigIpVersion')]"
            }, 
            "2": {
                "Port": "443"
            }, 
            "3": {
                "Port": "443"
            }, 
            "4": {
                "Port": "443"
            }, 
            "5": {
                "Port": "443"
            }, 
            "6": {
                "Port": "443"
            }
        }, 
        "bigIpVersionPortMap": {
            "12.1.2200": {
                "Port": 443
            }, 
            "13.1.0200": {
                "Port": 8443
            }, 
            "443": {
                "Port": 443
            }, 
            "latest": {
                "Port": 8443
            }
        }, 
        "apiVersion": "2015-06-15", 
        "computeApiVersion": "2017-12-01", 
        "networkApiVersion": "2017-11-01", 
        "storageApiVersion": "2017-10-01", 
        "location": "[resourceGroup().location]", 
        "subscriptionID": "[subscription().subscriptionId]", 
        "resourceGroupName": "[resourceGroup().name]", 
        "singleQuote": "'", 
        "f5CloudLibsTag": "v3.6.2", 
        "f5CloudLibsAzureTag": "v1.5.0", 
        "f5NetworksTag": "v4.4.0.1", 
        "f5CloudIappsTag": "v1.2.1", 
        "verifyHash": "[concat(variables('singleQuote'), 'cli script /Common/verifyHash {\nproc script::run {} {\n        if {[catch {\n            set hashes(f5-cloud-libs.tar.gz) 4cf5edb76d2e8dd0493f4892ff3679a58c8c79b1c02e550b55150d9002228c24c6d841095f1edd33fb49c5aaea518771252b4fb6d423a8a4ba8d94a0baf0f77a\n            set hashes(f5-cloud-libs-aws.tar.gz) 1a4ba191e997b2cfaaee0104deccc0414a6c4cc221aedc65fbdec8e47a72f1d5258b047d6487a205fa043fdbd6c8fcb1b978cac36788e493e94a4542f90bd92b\n            set hashes(f5-cloud-libs-azure.tar.gz) 5c256d017d0a57f5c96c2cb43f4d8b76297ae0b91e7a11c6d74e5c14268232f6a458bf0c16033b992040be076e934392c69f32fc8beffe070b5d84924ec7b947\n            set hashes(f5-cloud-libs-gce.tar.gz) 6ef33cc94c806b1e4e9e25ebb96a20eb1fe5975a83b2cd82b0d6ccbc8374be113ac74121d697f3bfc26bf49a55e948200f731607ce9aa9d23cd2e81299a653c1\n            set hashes(f5-cloud-libs-openstack.tar.gz) fb6d63771bf0c8d9cae9271553372f7fb50ce2e7a653bb3fb8b7d57330a18d72fa620e844b579fe79c8908a3873b2d33ee41803f23ea6c5dc9f7d7e943e68c3a\n            set hashes(asm-policy-linux.tar.gz) 63b5c2a51ca09c43bd89af3773bbab87c71a6e7f6ad9410b229b4e0a1c483d46f1a9fff39d9944041b02ee9260724027414de592e99f4c2475415323e18a72e0\n            set hashes(f5.http.v1.2.0rc4.tmpl) 47c19a83ebfc7bd1e9e9c35f3424945ef8694aa437eedd17b6a387788d4db1396fefe445199b497064d76967b0d50238154190ca0bd73941298fc257df4dc034\n            set hashes(f5.http.v1.2.0rc6.tmpl) 811b14bffaab5ed0365f0106bb5ce5e4ec22385655ea3ac04de2a39bd9944f51e3714619dae7ca43662c956b5212228858f0592672a2579d4a87769186e2cbfe\n            set hashes(f5.http.v1.2.0rc7.tmpl) 21f413342e9a7a281a0f0e1301e745aa86af21a697d2e6fdc21dd279734936631e92f34bf1c2d2504c201f56ccd75c5c13baa2fe7653213689ec3c9e27dff77d\n            set hashes(f5.aws_advanced_ha.v1.3.0rc1.tmpl) 9e55149c010c1d395abdae3c3d2cb83ec13d31ed39424695e88680cf3ed5a013d626b326711d3d40ef2df46b72d414b4cb8e4f445ea0738dcbd25c4c843ac39d\n            set hashes(f5.aws_advanced_ha.v1.4.0rc1.tmpl) de068455257412a949f1eadccaee8506347e04fd69bfb645001b76f200127668e4a06be2bbb94e10fefc215cfc3665b07945e6d733cbe1a4fa1b88e881590396\n            set hashes(f5.aws_advanced_ha.v1.4.0rc2.tmpl) 6ab0bffc426df7d31913f9a474b1a07860435e366b07d77b32064acfb2952c1f207beaed77013a15e44d80d74f3253e7cf9fbbe12a90ec7128de6facd097d68f\n            set hashes(asm-policy.tar.gz) 2d39ec60d006d05d8a1567a1d8aae722419e8b062ad77d6d9a31652971e5e67bc4043d81671ba2a8b12dd229ea46d205144f75374ed4cae58cefa8f9ab6533e6\n            set hashes(deploy_waf.sh) eebaf8593a29fa6e28bb65942d2b795edca0da08b357aa06277b0f4d2f25fe416da6438373f9955bdb231fa1de1a7c8d0ba7c224fa1f09bd852006070d887812\n            set hashes(f5.policy_creator.tmpl) 06539e08d115efafe55aa507ecb4e443e83bdb1f5825a9514954ef6ca56d240ed00c7b5d67bd8f67b815ee9dd46451984701d058c89dae2434c89715d375a620\n            set hashes(f5.service_discovery.tmpl) acc7c482a1eb8787a371091f969801b422cb92830b46460a3313b6a8e1cda0759f8013380e0c46d5214a351a248c029ec3ff04220aaef3e42a66badf9804041f\n\n            set file_path [lindex $tmsh::argv 1]\n            set file_name [file tail $file_path]\n\n            if {![info exists hashes($file_name)]} {\n                tmsh::log err \"No hash found for $file_name\"\n                exit 1\n            }\n\n            set expected_hash $hashes($file_name)\n            set computed_hash [lindex [exec /usr/bin/openssl dgst -r -sha512 $file_path] 0]\n            if { $expected_hash eq $computed_hash } {\n                exit 0\n            }\n            tmsh::log err \"Hash does not match for $file_path\"\n            exit 1\n        }]} {\n            tmsh::log err {Unexpected error in verifyHash}\n            exit 1\n        }\n    }\n    script-signature Kir5DhV/uRo0SwVRgPGrnNnAJBgHZ3XYraih5T90VbRZii5vPt0q3codJUdgoWiByQGpFREsa5Gy+v0+yYDAdYBzyZlThwRe+6RjWYfxP2+cKAC28wByJ0x6En1UD9kscj7ILUON5yv771izvIrxJ7x4Fd4RHcqB5++hWLvOLxXMiyJAYh2aUSOgdc+kx4lCHS6IU0aXtUxAQYpq510k4eS4UZJrfE7GPmpYkpRDJivR8UUyUWtuj0CAt3pWQEijKnC5zHhH6q5ikvQFn05PugcZO7RzOaA/a2gZw609wYAkXODMA6L49l+IKB31Y+/5ROB1w9/wf/H5RiP/kXC5/A==\n    signing-key /Common/f5-irule\n}', variables('singleQuote'))]", 
        "installCloudLibs": "[concat(variables('singleQuote'), '#!/bin/bash\necho about to execute\nchecks=0\nwhile [ $checks -lt 120 ]; do echo checking mcpd\n/usr/bin/tmsh -a show sys mcp-state field-fmt | grep -q running\nif [ $? == 0 ]; then\necho mcpd ready\nbreak\nfi\necho mcpd not ready yet\nlet checks=checks+1\nsleep 1\ndone\necho loading verifyHash script\n/usr/bin/tmsh load sys config merge file /config/verifyHash\nif [ $? != 0 ]; then\necho cannot validate signature of /config/verifyHash\nexit 1\nfi\necho loaded verifyHash\n\nconfig_loc=\"/config/cloud/\"\nhashed_file_list=\"$${config_loc}f5-cloud-libs.tar.gz f5.service_discovery.tmpl\"\nfor file in $hashed_file_list; do\necho \"verifying $file\"\n/usr/bin/tmsh run cli script verifyHash $file\nif [ $? != 0 ]; then\necho \"$file is not valid\"\nexit 1\nfi\necho \"verified $file\"\ndone\necho \"expanding $hashed_file_list\"\ntar xfz /config/cloud/f5-cloud-libs.tar.gz -C /config/cloud/azure/node_modules\ntouch /config/cloud/cloudLibsReady', variables('singleQuote'))]", 
        "dnsLabel": "[toLower(parameters('dnsLabel'))]", 
        "imageNameToLower": "[toLower(parameters('imageName'))]", 
        "skuToUse": "[concat('f5-bigip-virtual-edition-', variables('imageNameToLower'),'-byol')]", 
        "offerToUse": "[if(equals(parameters('bigIpVersion'), '12.1.2200'), 'f5-big-ip', concat('f5-big-ip-', variables('imageNameToLower')))]", 
        "bigIpNicPortValue": "[variables('bigIpNicPortMap')['2'].Port]", 
        "bigIpMgmtPort": "[variables('bigIpVersionPortMap')[variables('bigIpNicPortValue')].Port]", 
        "instanceName": "[toLower(parameters('instanceName'))]", 
        "availabilitySetName": "[concat(variables('dnsLabel'), '-avset')]", 
        "virtualNetworkName": "[concat(variables('dnsLabel'), '-vnet')]", 
        "vnetId": "[resourceId('Microsoft.Network/virtualNetworks', variables('virtualNetworkName'))]", 
        "vnetAddressPrefix": "[concat(parameters('vnetAddressPrefix'),'.0.0/16')]", 
        "publicIPAddressType": "Static", 
        "mgmtPublicIPAddressName": "[concat(variables('dnsLabel'), '-mgmt-pip')]", 
        "mgmtPublicIPAddressId": "[resourceId('Microsoft.Network/publicIPAddresses', variables('mgmtPublicIPAddressName'))]", 
        "mgmtNsgID": "[resourceId('Microsoft.Network/networkSecurityGroups/',concat(variables('dnsLabel'),'-mgmt-nsg'))]", 
        "mgmtNicName": "[concat(variables('dnsLabel'), '-mgmt')]", 
        "mgmtNicID": "[resourceId('Microsoft.Network/NetworkInterfaces', variables('mgmtNicName'))]", 
        "mgmtSubnetName": "mgmt", 
        "mgmtSubnetId": "[concat(variables('vnetId'), '/subnets/', variables('mgmtSubnetName'))]", 
        "mgmtSubnetPrefix": "[concat(parameters('vnetAddressPrefix'), '.1.0/24')]", 
        "mgmtSubnetPrivateAddress": "[concat(parameters('vnetAddressPrefix'), '.1.4')]", 
        "extSelfPublicIpAddressNamePrefix": "[concat(variables('dnsLabel'), '-self-pip')]", 
        "extSelfPublicIpAddressIdPrefix": "[resourceId('Microsoft.Network/publicIPAddresses', variables('extSelfPublicIpAddressNamePrefix'))]", 
        "extpublicIPAddressNamePrefix": "[concat(variables('dnsLabel'), '-ext-pip')]", 
        "extPublicIPAddressIdPrefix": "[resourceId('Microsoft.Network/publicIPAddresses', variables('extPublicIPAddressNamePrefix'))]", 
        "extNsgID": "[resourceId('Microsoft.Network/networkSecurityGroups/',concat(variables('dnsLabel'),'-ext-nsg'))]", 
        "extNicName": "[concat(variables('dnsLabel'), '-ext')]", 
        "extSubnetName": "external", 
        "extSubnetPrefix": "[concat(parameters('vnetAddressPrefix'), '.2.0/24')]", 
        "extSubnetId": "[concat(variables('vnetId'), '/subnets/', variables('extsubnetName'))]", 
        "extSubnetPrivateAddress": "[concat(parameters('vnetAddressPrefix'), '.2.4')]", 
        "extSubnetPrivateAddressPrefix": "[concat(parameters('vnetAddressPrefix'), '.2.')]", 
        "numberOfExternalIps": "[parameters('numberOfExternalIps')]", 
        "mgmtRouteGw": "[concat(parameters('vnetAddressPrefix'), '.1.1')]", 
        "tmmRouteGw": "[concat(parameters('vnetAddressPrefix'), '.2.1')]", 
        "routeCmdArray": {
            "12.1.2200": "[concat('tmsh create sys management-route waagent_route network 168.63.129.16/32 gateway ', variables('mgmtRouteGw'), '; tmsh save sys config')]", 
            "13.1.0200": "route", 
            "latest": "route"
        }, 
        "instanceTypeMap": {
            "Standard_A3": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_A4": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_A5": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_A6": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_A7": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D11": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D11_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D12": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D12_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D13": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D13_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D14": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D14_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D15_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D2_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D3": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D3_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D4": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D4_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_D5_v2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_DS1": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS11": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS11_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS12": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS12_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS13": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS13_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS14": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS14_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS15_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS1_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS2_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS3": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS3_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS4": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS4_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_DS5_v2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_F2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_F4": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_G1": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_G2": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_G3": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_G4": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_G5": {
                "storageAccountTier": "Standard", 
                "storageAccountType": "Standard_LRS"
            }, 
            "Standard_GS1": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_GS2": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_GS3": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_GS4": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }, 
            "Standard_GS5": {
                "storageAccountTier": "Premium", 
                "storageAccountType": "Premium_LRS"
            }
        }, 
        "tagValues": "[parameters('tagValues')]", 
        "newStorageAccountName0": "[concat(uniqueString(variables('dnsLabel'), resourceGroup().id, deployment().name), 'stor0')]", 
        "storageAccountType": "[variables('instanceTypeMap')[parameters('instanceType')].storageAccountType]", 
        "storageAccountTier": "[variables('instanceTypeMap')[parameters('instanceType')].storageAccountTier]", 
        "newDataStorageAccountName": "[concat(uniqueString(variables('dnsLabel'), resourceGroup().id, deployment().name), 'data000')]", 
        "dataStorageAccountType": "Standard_LRS", 
        "deploymentId": "[concat(variables('subscriptionId'), resourceGroup().id, deployment().name, variables('dnsLabel'))]", 
        "allowUsageAnalytics": {
            "No": {
                "hashCmd": "echo AllowUsageAnalytics:No", 
                "metricsCmd": ""
            }, 
            "Yes": {
                "hashCmd": "[concat('custId=`echo \"', variables('subscriptionId'), '\"|sha512sum|cut -d \" \" -f 1`; deployId=`echo \"', variables('deploymentId'), '\"|sha512sum|cut -d \" \" -f 1`')]", 
                "metricsCmd": "[concat(' --metrics customerId:$${custId},deploymentId:$${deployId},templateName:standalone_2nic-new_stack-experimental,templateVersion:4.4.0.1,region:', variables('location'), ',bigIpVersion:', parameters('bigIpVersion') ,',licenseType:BYOL,cloudLibsVersion:', variables('f5CloudLibsTag'), ',cloudName:azure')]"
            }
        }, 
        "customConfig": "### START (INPUT) CUSTOM CONFIGURATION HERE\n", 
        "installCustomConfig": "[concat(variables('singleQuote'), '#!/bin/bash\n', variables('customConfig'), variables('singleQuote'))]"
    }, 
    "resources": [
        {
            "apiVersion": "[variables('networkApiVersion')]", 
            "location": "[variables('location')]", 
            "name": "[variables('mgmtPublicIPAddressName')]", 
            "properties": {
                "dnsSettings": {
                    "domainNameLabel": "[variables('dnsLabel')]"
                }, 
                "idleTimeoutInMinutes": 30, 
                "publicIPAllocationMethod": "[variables('publicIPAddressType')]"
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Network/publicIPAddresses"
        }, 
        {
            "apiVersion": "[variables('networkApiVersion')]", 
            "location": "[variables('location')]", 
            "name": "[concat(variables('extSelfPublicIpAddressNamePrefix'), '0')]", 
            "properties": {
                "idleTimeoutInMinutes": 30, 
                "publicIPAllocationMethod": "[variables('publicIPAddressType')]"
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Network/publicIPAddresses"
        }, 
        {
            "apiVersion": "[variables('networkApiVersion')]", 
            "condition": "[not(equals(variables('numberOfExternalIps'),0))]", 
            "copy": {
                "count": "[if(not(equals(variables('numberOfExternalIps'), 0)), variables('numberOfExternalIps'), 1)]", 
                "name": "extpipcopy"
            }, 
            "location": "[variables('location')]", 
            "name": "[concat(variables('extPublicIPAddressNamePrefix'), copyIndex())]", 
            "properties": {
                "dnsSettings": {
                    "domainNameLabel": "[concat(variables('dnsLabel'), copyIndex(0))]"
                }, 
                "idleTimeoutInMinutes": 30, 
                "publicIPAllocationMethod": "[variables('publicIPAddressType')]"
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Network/publicIPAddresses"
        }, 
        {
            "apiVersion": "[variables('apiVersion')]", 
            "location": "[variables('location')]", 
            "name": "[variables('virtualNetworkName')]", 
            "properties": {
                "addressSpace": {
                    "addressPrefixes": [
                        "[variables('vnetAddressPrefix')]"
                    ]
                }, 
                "subnets": [
                    {
                        "name": "[variables('mgmtSubnetName')]", 
                        "properties": {
                            "addressPrefix": "[variables('mgmtSubnetPrefix')]"
                        }
                    }, 
                    {
                        "name": "[variables('extSubnetName')]", 
                        "properties": {
                            "addressPrefix": "[variables('extSubnetPrefix')]"
                        }
                    }
                ]
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Network/virtualNetworks"
        }, 
        {
            "apiVersion": "[variables('apiVersion')]", 
            "dependsOn": [
                "[variables('vnetId')]", 
                "[variables('mgmtPublicIPAddressId')]", 
                "[variables('mgmtNsgID')]"
            ], 
            "location": "[variables('location')]", 
            "name": "[variables('mgmtNicName')]", 
            "properties": {
                "ipConfigurations": [
                    {
                        "name": "[concat(variables('instanceName'), '-ipconfig1')]", 
                        "properties": {
                            "PublicIpAddress": {
                                "Id": "[variables('mgmtPublicIPAddressId')]"
                            }, 
                            "privateIPAddress": "[variables('mgmtSubnetPrivateAddress')]", 
                            "privateIPAllocationMethod": "Static", 
                            "subnet": {
                                "id": "[variables('mgmtSubnetId')]"
                            }
                        }
                    }
                ], 
                "networkSecurityGroup": {
                    "id": "[variables('mgmtNsgID')]"
                }
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Network/networkInterfaces"
        }, 
        {
            "apiVersion": "[variables('apiVersion')]", 
            "dependsOn": [
                "[variables('vnetId')]", 
                "[variables('extNsgID')]", 
                "extpipcopy", 
                "[concat('Microsoft.Network/publicIPAddresses/', variables('extSelfPublicIpAddressNamePrefix'), '0')]"
            ], 
            "location": "[variables('location')]", 
            "name": "[variables('extNicName')]", 
            "properties": {
                "copy": [
                    {
                        "count": "[add(variables('numberOfExternalIps'), 1)]", 
                        "input": {
                            "name": "[if(equals(copyIndex('ipConfigurations', 1), 1), concat(variables('instanceName'), '-self-ipconfig'), concat(variables('resourceGroupName'), '-ext-ipconfig', sub(copyIndex('ipConfigurations', 1), 2)))]", 
                            "properties": {
                                "PublicIpAddress": {
                                    "Id": "[if(equals(copyIndex('ipConfigurations', 1), 1), concat(variables('extSelfPublicIpAddressIdPrefix'), '0'), concat(variables('extPublicIPAddressIdPrefix'), sub(copyIndex('ipConfigurations', 1), 2)))]"
                                }, 
                                "primary": "[if(equals(copyIndex('ipConfigurations', 1), 1), 'True', 'False')]", 
                                "privateIPAddress": "[if(equals(copyIndex('ipConfigurations', 1), 1), variables('extSubnetPrivateAddress'), concat(variables('extSubnetPrivateAddressPrefix'), 1, sub(copyIndex('ipConfigurations', 1), 2)))]", 
                                "privateIPAllocationMethod": "Static", 
                                "subnet": {
                                    "id": "[variables('extSubnetId')]"
                                }
                            }
                        }, 
                        "name": "ipConfigurations"
                    }
                ], 
                "networkSecurityGroup": {
                    "id": "[concat(variables('extNsgID'))]"
                }
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Network/networkInterfaces"
        }, 
        {
            "apiVersion": "[variables('apiVersion')]", 
            "location": "[variables('location')]", 
            "name": "[concat(variables('dnsLabel'), '-mgmt-nsg')]", 
            "properties": {
                "securityRules": [
                    {
                        "name": "mgmt_allow_https", 
                        "properties": {
                            "access": "Allow", 
                            "description": "", 
                            "destinationAddressPrefix": "*", 
                            "destinationPortRange": "[variables('bigIpMgmtPort')]", 
                            "direction": "Inbound", 
                            "priority": 101, 
                            "protocol": "Tcp", 
                            "sourceAddressPrefix": "[parameters('restrictedSrcAddress')]", 
                            "sourcePortRange": "*"
                        }
                    }, 
                    {
                        "name": "ssh_allow_22", 
                        "properties": {
                            "access": "Allow", 
                            "description": "", 
                            "destinationAddressPrefix": "*", 
                            "destinationPortRange": "22", 
                            "direction": "Inbound", 
                            "priority": 102, 
                            "protocol": "Tcp", 
                            "sourceAddressPrefix": "[parameters('restrictedSrcAddress')]", 
                            "sourcePortRange": "*"
                        }
                    }
                ]
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Network/networkSecurityGroups"
        }, 
        {
            "apiVersion": "[variables('apiVersion')]", 
            "location": "[variables('location')]", 
            "name": "[concat(variables('dnsLabel'), '-ext-nsg')]", 
            "properties": {
                "securityRules": []
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Network/networkSecurityGroups"
        }, 
        {
            "apiVersion": "[variables('apiVersion')]", 
            "location": "[variables('location')]", 
            "name": "[variables('availabilitySetName')]", 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Compute/availabilitySets"
        }, 
        {
            "apiVersion": "[variables('storageApiVersion')]", 
            "kind": "Storage", 
            "location": "[variables('location')]", 
            "name": "[variables('newStorageAccountName0')]", 
            "sku": {
                "name": "[variables('storageAccountType')]", 
                "tier": "[variables('storageAccountTier')]"
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Storage/storageAccounts"
        }, 
        {
            "apiVersion": "[variables('storageApiVersion')]", 
            "kind": "Storage", 
            "location": "[variables('location')]", 
            "name": "[variables('newDataStorageAccountName')]", 
            "sku": {
                "name": "[variables('dataStorageAccountType')]", 
                "tier": "Standard"
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Storage/storageAccounts"
        }, 
        {
            "apiVersion": "[variables('apiVersion')]", 
            "dependsOn": [
                "[concat('Microsoft.Storage/storageAccounts/', variables('newStorageAccountName0'))]", 
                "[concat('Microsoft.Storage/storageAccounts/', variables('newDataStorageAccountName'))]", 
                "[concat('Microsoft.Compute/availabilitySets/', variables('availabilitySetName'))]", 
                "[concat('Microsoft.Network/networkInterfaces/', variables('mgmtNicName'))]", 
                "[concat('Microsoft.Network/networkInterfaces/', variables('extNicName'))]"
            ], 
            "location": "[variables('location')]", 
            "name": "[variables('instanceName')]", 
            "plan": {
                "name": "[variables('skuToUse')]", 
                "product": "[variables('offerToUse')]", 
                "publisher": "f5-networks"
            }, 
            "properties": {
                "availabilitySet": {
                    "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('availabilitySetName'))]"
                }, 
                "diagnosticsProfile": {
                    "bootDiagnostics": {
                        "enabled": true, 
                        "storageUri": "[reference(concat('Microsoft.Storage/storageAccounts/', variables('newDataStorageAccountName')), providers('Microsoft.Storage', 'storageAccounts').apiVersions[0]).primaryEndpoints.blob]"
                    }
                }, 
                "hardwareProfile": {
                    "vmSize": "[parameters('instanceType')]"
                }, 
                "networkProfile": {
                    "networkInterfaces": [
                        {
                            "id": "[resourceId('Microsoft.Network/networkInterfaces', variables('mgmtNicName'))]", 
                            "properties": {
                                "primary": true
                            }
                        }, 
                        {
                            "id": "[resourceId('Microsoft.Network/networkInterfaces', variables('extNicName'))]", 
                            "properties": {
                                "primary": false
                            }
                        }
                    ]
                }, 
                "osProfile": {
                    "adminPassword": "[parameters('adminPassword')]", 
                    "adminUsername": "[parameters('adminUsername')]", 
                    "computerName": "[variables('instanceName')]"
                }, 
                "storageProfile": {
                    "imageReference": {
                        "offer": "[variables('offerToUse')]", 
                        "publisher": "f5-networks", 
                        "sku": "[variables('skuToUse')]", 
                        "version": "[parameters('bigIpVersion')]"
                    }, 
                    "osDisk": {
                        "caching": "ReadWrite", 
                        "createOption": "FromImage", 
                        "name": "osdisk", 
                        "vhd": {
                            "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/', variables('newStorageAccountName0')), providers('Microsoft.Storage', 'storageAccounts').apiVersions[0]).primaryEndpoints.blob, 'vhds/', variables('instanceName'),'.vhd')]"
                        }
                    }
                }
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Compute/virtualMachines"
        }, 
        {
            "apiVersion": "[variables('computeApiVersion')]", 
            "dependsOn": [
                "[concat('Microsoft.Compute/virtualMachines/', variables('instanceName'))]"
            ], 
            "location": "[variables('location')]", 
            "name": "[concat(variables('instanceName'),'/start')]", 
            "properties": {
                "protectedSettings": {
                    "commandToExecute": "[concat('mkdir -p /config/cloud/azure/node_modules && cp f5-cloud-libs.tar.gz* /config/cloud; mkdir -p /var/log/cloud/azure; function cp_logs() { cd /var/lib/waagent/custom-script/download && cp `ls -r | head -1`/std* /var/log/cloud/azure; }; TMP_DIR=/mnt/creds; TMP_CREDENTIALS_FILE=$TMP_DIR/.passwd; BIG_IP_CREDENTIALS_FILE=/config/cloud/.passwd; /usr/bin/install -b -m 755 /dev/null /config/verifyHash; /usr/bin/install -b -m 755 /dev/null /config/installCloudLibs.sh; /usr/bin/install -b -m 400 /dev/null $BIG_IP_CREDENTIALS_FILE; IFS=', variables('singleQuote'), '%', variables('singleQuote'), '; echo -e ', variables('verifyHash'), ' >> /config/verifyHash; echo -e ', variables('installCloudLibs'), ' >> /config/installCloudLibs.sh; echo -e ', variables('installCustomConfig'), ' >> /config/customConfig.sh; unset IFS; bash /config/installCloudLibs.sh; . /config/cloud/azure/node_modules/f5-cloud-libs/scripts/util.sh; create_temp_dir $TMP_DIR; echo ', variables('singleQuote'), parameters('adminPassword'), variables('singleQuote'), '|sha512sum|cut -d \" \" -f 1|tr -d \"\n\" > $TMP_CREDENTIALS_FILE; bash /config/cloud/azure/node_modules/f5-cloud-libs/scripts/createUser.sh --user svc_user --password-file $TMP_CREDENTIALS_FILE; f5-rest-node /config/cloud/azure/node_modules/f5-cloud-libs/scripts/encryptDataToFile.js --data-file $TMP_CREDENTIALS_FILE --out-file $BIG_IP_CREDENTIALS_FILE; wipe_temp_dir $TMP_DIR;', variables('allowUsageAnalytics')[parameters('allowUsageAnalytics')].hashCmd, '; /usr/bin/f5-rest-node /config/cloud/azure/node_modules/f5-cloud-libs/scripts/onboard.js --output /var/log/cloud/azure/onboard.log --log-level debug --host ', variables('mgmtSubnetPrivateAddress'), ' --ssl-port ', variables('bigIpMgmtPort'), ' -u svc_user --password-url file:///config/cloud/.passwd --password-encrypted --hostname ', concat(variables('instanceName'), '.', resourceGroup().location, '.cloudapp.azure.com'), ' --license ', parameters('licenseKey1'), ' --ntp ', parameters('ntpServer'), ' --tz ', parameters('timeZone'), ' --db tmm.maxremoteloglength:2048', variables('allowUsageAnalytics')[parameters('allowUsageAnalytics')].metricsCmd, ' --module ltm:nominal --module afm:none; /usr/bin/f5-rest-node /config/cloud/azure/node_modules/f5-cloud-libs/scripts/network.js --output /var/log/cloud/azure/network.log --host ', variables('mgmtSubnetPrivateAddress'), ' --port ', variables('bigIpMgmtPort'), ' -u svc_user --password-url file:///config/cloud/.passwd --password-encrypted --default-gw ', variables('tmmRouteGw'), ' --vlan name:external,nic:1.1 --self-ip name:self_2nic,address:', variables('extSubnetPrivateAddress'),  ',vlan:external --log-level debug', '; if [[ $? == 0 ]]; then tmsh load sys application template f5.service_discovery.tmpl; ', variables('routeCmdArray')[parameters('bigIpVersion')], '; bash /config/customConfig.sh; $(cp_logs); else $(cp_logs); exit 1; fi', '; if grep -i \"PUT failed\" /var/log/waagent.log -q; then echo \"Killing waagent exthandler, daemon should restart it\"; pkill -f \"python -u /usr/sbin/waagent -run-exthandlers\"; fi')]"
                }, 
                "publisher": "Microsoft.Azure.Extensions", 
                "settings": {
                    "fileUris": [
                        "[concat('https://raw.githubusercontent.com/F5Networks/f5-cloud-libs/', variables('f5CloudLibsTag'), '/dist/f5-cloud-libs.tar.gz')]", 
                        "[concat('https://raw.githubusercontent.com/F5Networks/f5-cloud-iapps/', variables('f5CloudIappsTag'), '/f5-service-discovery/f5.service_discovery.tmpl')]"
                    ]
                }, 
                "type": "CustomScript", 
                "typeHandlerVersion": "2.0"
            }, 
            "tags": "[if(empty(variables('tagValues')), json('null'), variables('tagValues'))]", 
            "type": "Microsoft.Compute/virtualMachines/extensions"
        }
    ], 
    "outputs": {
        "GUI-URL": {
            "type": "string", 
            "value": "[concat('https://', reference(variables('mgmtPublicIPAddressId')).dnsSettings.fqdn, ':', variables('bigIpMgmtPort'))]"
        }, 
        "SSH-URL": {
            "type": "string", 
            "value": "[concat(reference(variables('mgmtPublicIPAddressId')).dnsSettings.fqdn, ' ',22)]"
        }
    }
}
DEPLOY
deployment_mode = "Complete"
}


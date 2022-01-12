//Original work from https://github.com/DealerDotCom/terraform-provider-bigip
//Modifications Copyright 2019 F5 Networks Inc.
//This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
//If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.

module github.com/F5Networks/terraform-provider-bigip

require (
	github.com/Azure/azure-sdk-for-go v53.4.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.13.0
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/adal v0.9.13
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/f5devcentral/go-bigip v0.0.0-20211208144806-c9b2472e4619
	github.com/f5devcentral/go-bigip/f5teem v0.0.0-20211208144806-c9b2472e4619
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.3
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/zclconf/go-cty v1.5.1 // indirect
)

go 1.16

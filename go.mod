//Original work from https://github.com/DealerDotCom/terraform-provider-bigip
//Modifications Copyright 2019 F5 Networks Inc.
//This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
//If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.

module github.com/F5Networks/terraform-provider-bigip

require (
	github.com/f5devcentral/go-bigip v0.0.0-20210218015208-fde2d84238d6
	github.com/f5devcentral/go-bigip/f5teem v0.0.0-20210218015208-fde2d84238d6
	github.com/google/uuid v1.1.1
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.1.0
	github.com/stretchr/testify v1.3.0
)

go 1.13

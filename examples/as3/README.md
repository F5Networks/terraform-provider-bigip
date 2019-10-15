
[//]: # (Copyright 2019 F5 Networks Inc.)
[//]: # (This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.)
[//]: # (If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.)
# AS3 Deployment using terraform  null resource
- This shows details about how you can deploy AS3 RPM using null resource calling shell scripts to deploy AS3 RPM and AS3 Json payload. 
# How you can use AS3 with null resource ?
- Look at the ``as3.tf`` file it uses null resources twice 
- First null resource ``install_as3`` uses script ``install_as3.sh`` script to load the as3 rpm on the BIG-IP
- Second null resource uses shell script ``as3_http.sh`` to deploy example1.json.

For more information on AS3 please look at https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/ 
More information on null resource and provisioner refer to https://www.terraform.io/docs/provisioners/null_resource.html

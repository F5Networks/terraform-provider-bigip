## 1.1.1 (December 19, 2019)

## Bug Fixes
- bigip_ssl resources not over writing existing cert/key #218
- Content argument of `bigip_ssl_key` should be marked sensitive #208
- Pool attachment docs is not updated upto date #207
- Bigip provider - add a parameter to specify the mgmt port #23
- AS3 module - tenant_name usage #24
- [doc] DO module - mistake in documentation #25
- creating Client SSL Profile with non-default partition Failed using terraform #27
- creating Server SSL Profile with non-default partition Failed using terraform #28

## 1.1.0 (November 22, 2019)

## Added Functionalities
- Terraform resources for DO( Declarative Onboarding )
- Docs for DO resources
- Terraform Provisioner for DO/AS3 installation mentioned in examples section of repo
- Docs for terraform Client/Server SSL resource profiles
- Terraform resource for importing SSL Certificates on bigip with docs
- Terraform resource for importing SSL Keys on bigip with docs

## Bug Fixes

- build ssl profile #17
- make build failed #14
- AWS example needs to be updated #15
- Having trouble logging into f5 #18
- Bigip_ltm_virtual_server attribute name not updating on apply terraform-providers/terraform-provider-bigip#178
- Docs for ltm_virtual_server incorrect terraform-providers/terraform-provider-bigip#171
- Missing Documentation for SSL Client/Server Profiles terraform-providers/terraform-provider-bigip#188
- Can’t change virtual server name? # terraform-providers/terraform-provider-bigip#186
- Terraform Official docs way behind # terraform-providers/terraform-provider-bigip#182
- Error: Unsupported argument on bigip_ltm_policy # terraform-providers/terraform-provider-bigip#176
- Not possible to remove persistence profile for a Virtual Server # terraform-providers/terraform-provider-bigip#169
- Cannot Modify Datagroup # terraform-providers/terraform-provider-bigip#180

## 1.0.0 (October 25, 2019)
- Added membership based monitor map
- Fix a URL issue in readme
- Added SSL code
- Added test conditions for udp
- Added License banner to shell scripts, travis.yml, goreleaser.yml
- Updated AS3 docs
- Added banner to resource files


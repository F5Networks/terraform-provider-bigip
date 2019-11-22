## 1.1.0 (Unreleased)

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

## 0.12.4 (August 14, 2019)
- Fix #139 changing required parameters to optional in _bigip_ltm_policy
- Added #134 SSL Ssl Client Profile and Server Profile feature
- Added #137 Bigip AS3 integration
- Fix Changed Required to optional for tenant name
- Fix #128 Addition of description field for virtual server/pool/node
- Fix #126 Fix for Changing name in 'bigip_ltm_profile_http2' causes Terraform Crash
- Added #116 Add Node/Virtualserver with Routedomain

## 0.12.3 (June 06, 2019)
- Fix for terraform 0.12 
- Fix the test TF files for terraform 12
## 0.12.2 (May 02, 2019)
- go-bigip vendor update for vxlan, tunnel interfaces
- Changed defaults to Computed for couple of resources
## 0.12.1 (April 23, 2019) (April 2019)
- Fixed #80 #81
- Added http profile resource with documentation
- Fixed #67 issue Unable to pass username and password to monitor
- Fixed #63 added documentation for data datagroup
- Fixed #59 Created Ftp monitor resource
- Fixed #58 Ability to provision FTP virtual servers and monitors
- Fixed #54  Switch to Go Modules
- Fixed #49 Docs updated for ltm node resources
- Fixed #46 Unable to set Alias Service Port on HTTPS monitor
- Fixed #35 bigip_ltm_snat missing functionality
- Fixed #25 add "content list" For bigip_ltm_profile_httpcompress resource
- Added include - exclude to the resource httpcompress profile
- Added Valid function to node resource
- Added  pool_attachement resource doc
- Improvement to node resource , interval
- Simplified some parts with new utility methods (SelfIP & Vlan)
- Fix SelfIP and VLAN Read methods + other minor fixes
- Implement import for pool attachement resource #84

## 0.12.0 (September 26, 2018)
- Added couple of resources like snat, snmp, profiles, test modules etc.

## 0.3.0
- iRule creation support
- **Breaking Change** - rules property on bigip_ltm_virtual_server renamed to irules

## 0.2.0

- Added profiles, irules, source_address_translation to virtual servers
- Cleaned up handling of lists

## 0.1.0

- Initial release

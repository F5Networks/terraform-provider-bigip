## 1.4.0 (Unreleased)
## 1.3.0 (July 23, 2020)

# Features additions:

 - AS3 resource for BIGIQ

# Bug Fixes:

1. F5 LTM default custom profile values should calculated rather than hardcoded in code [https://github.com/terraform-providers/terraform-provider-bigip/issues/298]
2. Handling Common/Shared tenant created via AS3
3. DO declaration fails with CRASH error
4. Terraform crashes when the structure of response from bigip changes
5. Error: produced an unexpected new value for was present, but new absent.#305 [https://github.com/terraform-providers/terraform-provider-bigip/issues/305] 

## 1.2.1 (June 11, 2020)

# Bug Fixes

1. Provider shows passwords in clear text when issuing terraform plan.#279
2. Terraforn apply crash with bigip_as3 and F5 VE 15 #291 
3. AS3 apply fails on Terraform #294 
4. Documentation Link on readme.md [F5Networks/terraform-provider-bigip/issues/85]

# Additional Changes

1. Added Acceptance test for terraform resource ""bigip_ltm_pool_attachment"
2. Documentation Update for terraform resource ""bigip_ltm_pool_attachment"
3. Update Example for terraform resource "bigip_ltm_pool_attachment"
4. New test scenerios for terraform resource "bigip_as3"

## 1.2.0 (May 11, 2020)

# Feature additions:

- Terraform resource module for BIGIP Licence management through BIGIQ
- As3 Schema validation.
- AS3 TEEM control Agent additions.
- Terraform resource module for bigip_command /Run TMSH and BASH commands on F5 devices

# Bug Fixes:
1. big_ltm_virtual_server does not work with IPv6 address [https://github.com/F5Networks/terraform-provider-bigip/issues/62, #278]
2. declaring virtual addresses in /Common/Shared via AS3 fails [ https://github.com/F5Networks/terraform-provider-bigip/issues/48]
3. BIG-IQ Licensing  [https://github.com/F5Networks/terraform-provider-bigip/issues/44]
4. DO Error on Terraform destroy [https://github.com/F5Networks/terraform-provider-bigip/issues/43]
5. The provider provider.bigip does not support resource type "bigip_command".[https://github.com/F5Networks/terraform-provider-bigip/issues/63]
6. `bigip_as3` Read/Exists/Update actions should be restricted to target tenant #253
7. `bigip_as3` resource `resourceBigipAs3Read` action does not store actual value in state #254 
8. v1.1.2 changed the contract of the `bigip_as3` resource #267 
9. autopopulate not passed to pool attachment #242
10. Error while Sending/Posting http request with DO json :{"code":404 #243
11. Rework ltm policy (#241)

## 1.1.2 (March 19, 2020)

# Bug Fixes

- Missing "database" entry for PostgreSQL monitor #224 ( https://github.com/terraform-providers/terraform-provider-bigip/issues/224 )
- `bigip_as3` resource should validate JSON #227 ( https://github.com/terraform-providers/terraform-provider-bigip/issues/227 )
- bigip_as3 - doesn't delete resource #38 ( https://github.com/F5Networks/terraform-provider-bigip/issues/38 )
- examples for bigip resources in repo are not compatible with terraform 0.12 #40.(https://github.com                    /F5Networks/terraform-provider-bigip/issues/40 )
- Looks like provisioner resource in sys.go is not complete #244 ( https://github.com/terraform-providers/terraform-provider-bigip/issues/244 )
- bigip_as3 - doesn't delete resource #38 (https://github.com/F5Networks/terraform-provider-bigip/issues/38)
- `bigip_as3` resource `resourceBigipAs3Read` action does not store actual value in state #254  ( https://github.com/terraform-providers/terraform-provider-bigip/issues/254 )
- Unable to modify/update data group #248 ( https://github.com/terraform-providers/terraform-provider-bigip/issues/248 )
- Terraform crash when creating SSL certificate resources on F5 BIG-IP #255 (https://github.com/terraform-providers/terraform-provider-bigip/issues/255 )

# Other Notes:

- Any Documentation changes for terraform resources w.r.t above bug fixes are updated
- Bigip_as3 resource now read as3 json from bigip and set the terraform state file, but as3 json from bigip will not have all the standard as3 classes as given from user json as input to tf file. So sometimes though there may not be actual changes between user as3 json and bigip as3 json ,( Top level AS3 class will not be there in bigip as3 json ) terraform will detect as change and when we do terraform apply it will says 1 changed. But it will be same json and there will be no change in bigip ( as3 is idempotent ).


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

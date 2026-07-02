## 1.28.0 (July 1st, 2026)

# Features additions:

 - Added `bigip_gtm_pool` resource for creating and managing GTM pools
 - Added `bigip_gtm_pool_attachment` (or `pools` attribute improvements) to attach pools to WideIPs more reliably
 - Added `bigip_gtm_virtual_server` data source for reading existing GTM virtual servers
 - Improvements to GTM topology handling: better import/export and topology-region association

## 1.27.0 (June 3rd, 2026)

# Features additions:

 - Added `pools` attribute to `bigip_gtm_wideip` resource for associating GTM pools with a WideIP
 - Added `bigip_gtm_datacenter` data source for reading existing GTM datacenters
 - Added `bigip_gtm_server` data source for reading existing GTM servers
 - Added `bigip_gtm_topology_region` resource for managing GTM topology regions (subnets, countries, states, datacenters, etc.)
 - Added `bigip_gtm_topology_record` resource for managing GTM topology-based routing rules

## 1.26.0 (March 25, 2026)

  - Add GTM (Global Traffic Manager) monitor resources support
  - bigip_gtm_monitor_http
  - bigip_gtm_monitor_https
  - bigip_gtm_monitor_tcp
  - bigip_gtm_monitor_postgresql
  - bigip_gtm_monitor_bigip

## 1.25.1 (March 18, 2026)

# Bug Fixes:

- fix ltm policy rule removal corruption #1128

## 1.25.0 (March 3, 2026)

- Add GTM (Global Traffic Manager) resources support
- Bump crazy-max/ghaction-import-gpg from 6.3.0 to 7.0.0
- Bump goreleaser/goreleaser-action from 6.4.0 to 7.0.0
- Fix formatting of subcategory in bigip_ltm_ifile.md
- Add session param to docs
- Made a fqdn configuration block guide tweak

## 1.24.2 (January 9, 2026)

- Bump actions/checkout from 5 to 6
- Bump golang.org/x/crypto from 0.36.0 to 0.45.0
- Bump golangci/golangci-lint-action from 8 to 9
- Adding fix for bug
- Adding monitor doc fixes

## 1.24.1 (October 24, 2025)

- Bump actions/checkout from 4 to 5
- Bump actions/setup-go from 5 to 6
- Bump goreleaser/goreleaser-action from 6.3.0 to 6.4.0
- Adding ifile for ltm
- Adding new ifile resource
- Fix for issue967
- Vendor sync and version changes

## 1.24.0 (August 12, 2025)

# Bug Fixes:

- Fix crash and add tests
- Fix crash issues
- Fix data source issues
- Fixed vendor sync issues and lint issues

## 1.23.1 (July 7, 2025)

# Bug Fixes:

- Adding patch fix for crash

## 1.23.0 (July 3, 2025)

# Features:

- Added Support for Domain in SMTP monitors
- Added the support for deletion of app from a tenant
- Create SECURITY.md

# Bug Fixes:

- Fix server ssl profile issue
- Adding vendor changes

## 1.22.10 (May 19, 2025)

- Bump golang.org/x/net from 0.36.0 to 0.38.0
- Bump golangci/golangci-lint-action from 7 to 8
- Fixed the persistent profiles config drift issue
- Fix as3 idempotency issue for schema in json

## 1.22.9 (April 7, 2025)

# Bug Fixes:

- Fix awaf issues

# Other:

- Bump crazy-max/ghaction-import-gpg from 6.2.0 to 6.3.0
- Bump golang.org/x/net from 0.33.0 to 0.36.0
- Bump golangci/golangci-lint-action from 6 to 7
- Bump goreleaser/goreleaser-action from 6.2.1 to 6.3.0
- Adding ci/cd changes
- Adding vendor sync changes

## 1.22.8 (February 26, 2025)

# Bug Fixes:

- Cannot assign non-RSA (ECC) certificate and key to client ssl profile #839

# Other:

- Bump goreleaser/goreleaser-action from 6.1.0 to 6.2.1
- Adding cerkey chain docs
- Updated release yaml

## 1.22.7 (January 16, 2025)

# Bug Fixes:

- Need support for Distributed Cloud Bot Defense profiles #1021
- Failing to create external bigip_ltm_datagroup resource #1024
- Add dry run to AS3 resource
- bigip_ltm_profile_http removing fallback_host does not result in change #1009

## 1.22.6 (December 3, 2024)

# Bug Fixes:

- bigip_ltm_policy: can't create new plan when initial creation failed #1007
- bigip_ltm_pool_attachment import block does not populate most properties #1031
- Issue #1014 - vCMP guest virtual-disk removal works when there's more than 1 entry #1015
- Add Bash example for bigip_command #1017
- Fix handling of single quotes in shell commands #1019

## 1.22.5 (October 21, 2024)

# Bug Fixes:

- Fix virtual address mask issue
- Use of single quotes causes unexpected behaviour in bigip_command #1018
- VCMP Guest deleteVirtualDisk only checks first disk name #1014
- resource bigip_ltm_profile_client_ssl, CRL file and new option Allow Expired CRL File #997

## 1.22.4 (September 10, 2024)

# Bug Fixes:

- persistence_profiles changes are not handled correctly on import #993
- bigip_as3 returns Successful deployment when AS3 actually fails and does NOT deploy #1000

## 1.22.3 (August 1, 2024)

# Bug Fixes:

- Application_list is forcing Per-App deployments to constantly change #987

## 1.22.2 (June 19, 2024)

# Features:

- Addition of "Enforcement Mode" for Cross Domain Request Enforcement in URL properties #958

# Bug Fixes:

- Concurrent execution of bigip_ssl_key_cert fails #968
- Request-log templates escaped double-quotes not being configured on the BIG-IP #960
- Wrong value saved in state for request log profile #958
- Fake changes on security log profile in virtual server resource #947
- Provider detects changes when none are present for waf policy resource #994
- AS3 per-app bug-fixes #987

## 1.22.1 (May 9, 2024)

# Bug Fixes:

- Failure to refresh state #954
- AS3 deployment with v1.22 always detects changes #972
- Addition of "Enforcement Mode" for Cross Domain Request Enforcement in URL properties #958
- Big-IQ bigip_bigiq_as3 tenant resource removal/destruction fails with 500 error #959
- bigip_ltm_profile_http encrypt_cookie_secret plan diffs #941

## 1.22.0 (March 27, 2024)

# Features:

- RFE: Added Perapp AS3 Declaration support in AS3 Deploy resource #938

# Bug Fixes:

- Issue with virtual address deletion or replacement when the name field contains the route domain #936

## 1.21.0 (February 14, 2024)

# Features:

- RFE: Add new resource to create/manage Rewrite profiles #924
- Missing Terraform resource to "Create profile for request logging" #892
- Add Bot Profile resource with template argument #775

# Bug Fixes:

- Support for check_max_value_length in parameters for declarative WAF policy #925
- Common tenant is trying to be updated without modification in the json file #869
- Deleting iRules from Virtual Server Configs #923
- Fix Virtual Server 0.0.0.0 Mask Conversion #922

## 1.20.2 (January 3, 2024)

# Bug Fixes:

- Issue with clientssl profile when using parent/child (defaults-from) #902
- Possibility to configure the MTU size on a vlan #908

## 1.20.1 (November 23, 2023)

# Bug Fixes:

- EOF Error with BigIP v17.10.1 running FIPS #817
- .mgmt/shared/authn/login: EOF when connecting to bigip with 1.19.0 and 1.20.0 #894
- Fix add Server Agent Name in resource bigip_ltm_profile_http #836
- RFE: Add support for Trusted/Advertised CAs in resource bigip_ltm_profile_client_ssl #872
- HSTS support in resource bigip_ltm_profile_http #834
- Add persistence_profiles argument for FAST TCP/HTTP(S) applications #883

## 1.20.0 (October 13, 2023)

# Enhancements:

- Add HTTP Methods support in resource bigip_ltm_profile_http #833
- Add HSTS support in resource bigip_ltm_profile_http #834
- Add Maximum Header size/count support in resource bigip_ltm_profile_http #835
- Add Server Agent Name support in resource bigip_ltm_profile_http #836
- Add support for OCSP in certificate resource #850
- F5 Bigip Cipher Rules and Cipher Groups #654
- bigip_ltm_profile_http terraform provider resource type missing Argument request #681

# Bug Fixes:

- Issue importing full json WAF policy with v1.19.0 - resource bigip_waf_policy #858

## 1.19.0 (August 28, 2023)

# Enhancements:

- RFE: add support for description attribute in resource bigip_ltm_policy #838
- RFE: new resource bigip_ssl_key_certificate #832

# Bug Fixes:

- bigip_sys_iapp update, change and destroy issue when partition is not Common #825
- SSL cert struct using 32 bit signed int vulnerable to Y2038 problem #818

## 1.18.1 (July 14, 2023)

# Bug Fixes:

- Unable to deploy a parameter when the level is set to "global" using data.bigip_waf_entity_parameter json #822
- Is v17.x supported? #820
- Refresh of profiles in resource bigip_ltm_virtual_server #679
- invalid memory address or nil pointer dereference #674

## 1.18.0 (June 5, 2023)

# Enhancements:

- Add ability to allow/disallow file_types for WAF policy #799
- Add support for WAF IP address exception #762
- Add support WAF GraphQL security settings #763

# Bug Fixes:

- bigip_sys_iapp An argument named "executeAction" is not expected here #801
- AWAF host_names are ignored when setting policy_import_json #792
- Issue with declaring a replace http_uri action inside policy rule #794
- bigip_ltm_monitor - custom parent monitor is not allowed #721

## 1.17.1 (April 26, 2023)

# Enhancements:

- Adding consul service discovery for FAST HTTP app
- Adding consul service discovery for FAST HTTPs app

# Bug Fixes:

- SR- Plugin crash on 1.16.2 version
- Import of bigip_ltm_virtual_server crashes terraform - #2 #764
- AWAF host_names are ignored when setting policy_import_json #792
- panic: runtime error: failed to respond to the plugin.(*GRPCProvider).ConfigureProvider call #785
- Issue with declaring a replace http_uri action inside policy rule #794

## 1.17.0 (March 16, 2023)

# Enhancements:

- Adding service discovery for FAST HTTP app
- Adding service discovery for FAST HTTPs app
- Plugin migration to sdkv2

# Bug Fixes:

- SR-Terraform | policy_import_json failing to import all the JSON elements
- SR- invalid memory address or nil pointer dereference #674
- Error during WAF policy deployment when setting file_types #765
- Providing policy_import_json with other arguments causes crash #766
- bigip_ltm_policy fails applying new rules when asm enable is changed #737
- SR[Terraform]grpc error "received message larger than max"

## 1.16.2 (January 31, 2023)

# Bug Fixes:

- SR-Terraform | policy_import_json failing to import all the JSON elements
- SR- invalid memory address or nil pointer dereference #674
- missing source_address_translation property when doing "terraform import" of bigip_ltm_virtual_server #755
- bigip_net_vlan add the option for cmp-hash #751
- Deleting AS3 declaration with bigip-as3 resource still throws an error #735
- Adding smtp monitor to parentMonitors #756
- Add support to restrict allowed WAF policy host names #748
- Import of bigip_ltm_virtual_server crashes terraform #729

## 1.16.1 (December 22, 2022)

# Bug Fixes:

- bigip_ltm_virtual_server - fw-enforced-policy, source-port are not available #722
- invalid memory address or nil pointer dereference #674
- ltm policy rule in order to permit the configuration of a datagroup for TCP address filtering #706
- bigip_sys_dns and bigip_sys_ntp resources remain after terraform destroy #708
- Import methodsOverrideOnUrl not set when method_overrides #715
- Add support for verified accept on LTM TCP profile #724
- 1.16.0 starts throwing certificate error #720
- VCMP Support? #157

## 1.16.0 (November 8, 2022)

# Bug Fixes:

- Unable to create VIP due to illegally sharing resources #712
- bigip_ltm_profile_client_ssl does not update property "chain" #710

# Feature Requests:

- bigip_ltm_profile_httpcompress missing Argument request #684
- bigip_ltm_profile_tcp missing Argument request #683
- bigip_ltm_profile_fastl4 missing Argument request #682

## 1.15.2 (September 29, 2022)

# Bug Fixes:

- bigip_ltm_policy cannot be created with an additional path #634
- bigip_ltm_policy; operand 'ssl-extension' is not available during event 'request' #648
- bigip_ltm_policy - action replace #591
- Crash when running terraform apply #26

# Other Fixes:

- FAST Resources naming fixes
- Added WAF policy support FAST Https Resource

## 1.15.1 (August 18, 2022)

# Bug Fixes:

- bigip_ltm_node does not show difference after monitor modification #526
- Provider produced inconsistent final plan on ltm_pool_attachments with route domains #660
- Cannot create ltm_pool_attachments with route domains #661
- bigip_ltm_monitor based on ldap parent monitor #632
- bigip_ltm_policy; While injecting a new Rule on apply error is received #649
- v1.15.0 Missing darwin_arm64 build - Broken on M1 Mac #664
- Terraform AWAF Integration enhancements like filetype support and graphql profile addition

## 1.15.0 (July 11, 2022)

# New Features:

- Terraform FAST TCP, HTTP and HTTPS application integration
- Terraform AWAF Integration enhancements

# Bug Fixes:

- bigip_do timeout issue #607
- Handling of Self IPs within a route domain #608
- resource bigip_as3 plugin does not respond and throws a runtime error and crashes #611

## 1.14.0 (May 26, 2022)

# Bug Fixes:

- Terraform AWAF Integration new resources
- Provider crash when renaming AS3 partition #601
- Terraform state inconsistency for LTM and DNS AS3 declaration #604
- bigip_ltm_virtual_address should ForceNew when changing name #612
- Update the terraform documentation for BIG-IQ #606

## 1.13.1 (April 13, 2022)

# Bug Fixes:

- Basic Auth sent with token request to mgmt/shared/authn/login #602
- Network resource bigip_net_selfip documentation missing correct port_lockdown explanation #594
- Error: no change not handled by provider #600

## 1.13.0 (March 3, 2022)

# Feature Request:

- Adding Support for external data group creation fixes #583

# Bug Fixes:

- bigip_ltm_monitor of parent type mysql/mssql #580
- Additional documentation fixes

## 1.12.2 (January 19, 2022)

# Bug Fixes:

- Fixed vlan_disabled issue for bigip_ltm_snat and bigip_ltm_virtual_server
- Fixed issue with sending AS3 json payload using bigip_as3

## 1.12.1 (December 9, 2021)

# Bug Fixes:

- Incorrect 'bigip_ltm_monitor' parent profile name validation #553
- bigip_do showing json difference even when there is not any #557

# Enhancements:

- Support Port Lockdown for self-ip
- Support BIGIP access data in bigip_do

## 1.12.0 (October 27, 2021)

# Bug Fixes:

- Support new Mac Laptops - aka build arm64 modules #545
- terraform provider plugin crashed when deploying AS3 via BIG-IQ with defaultRouteDomain property #558
- ADC id being overwritten by subsequent deployments #560

# Enhancements:

- bigip_ltm_policy datasource #547
- Documentation updates

## 1.11.1 (September 16, 2021)

# Bug Fixes:

- Added TEEM userAgent for FAST Application resource
- Fixed testing issues for http2 profile

# Misc Fixes:

- Added golang-lint into github actions
- Added plugin binary generation into github action

## 1.11.0 (August 4, 2021)

# Bug Fixes:

- Fixed BIG-IQ resource use requires BIG-IP login #442
- Fixed Allow managing default profiles #461
- Created New resource for bigip_ltm_profile_ftp to allow Custom FTP Profiles #480
- Fixed Export full_path when a resource name isn't its full path #490
- Fixed feature request: license revocation #417
- Fixed bigip_ltm_node fails when monitor is set and session is not defined #521
- terraform destroy fails while trying to destroy big_ltm_virtual_address resource #519
- Re-open: Inconsistent documentation locations, docs vs website/docs #518
- bigip_fast_application fails with null exception error #527

# Other Fixes:

- Validate Terraform v1.X
- Testcase addition for profile http

## 1.10.0 (June 23, 2021)

# Bug Fixes:

- Fixed Inconsistent documentation locations, docs vs website/docs #515
- Fixed posting as3 config failed for tenants:(Sample_01) with error #511
- Fixed bigip_net_route arguments typo #508
- Fixed Change the name field of bigip_ltm_profile_client_ssl causes provider fatal error #505
- Fixed Update Command Reference Guide #498
- Fixed updating FAST template does not trigger an update #492
- Possibility to configure a ltm policy in order to reset traffic connection #463
- Update a client_ssl profile breaks inheritance from parent profile #450
- Error: Error create profile Ssl (/partition/profile_5601): HTTP 400 #449
- terraform import on bigip_ltm_monitor does not properly import parent object #443

# Enhancements:

- Support for Security Log Profiles in Virtual Server resource #447
- Added checksum field for resource for bigip_fast_template to identify if content modifies

## 1.9.0 (May 13, 2021)

# Bug Fixes:

- Fixed bigip_bigiq_as3 doesn't work with auth token #436
- Fixed bigip_bigiq_as3 fails with multiple applications #437
- Fixed bigip_sys_snmp_traps fails when changing port #448
- Fixed bigip_ltm_virtual_server - does not remove unrequired profile #467
- Fixed Panic crash #470
- Fix Node Datasource acceptance test

# Enhancements:

- Terraform resource for bigip_fast_template
- Terraform resource for bigip_fast_application
- Terraform resource for bigip_net_ike_peer
- Add support to define resources inside directories for other resources

## 1.8.0 (April 1, 2021)

# Bug Fixes:

- clientssl profile options list doesnt seem to be working #424
- serverssl profile options list doesnt seem to be working
- bigip_ltm_persistence_profile_dstaddr - Changes to a custom PARENT profile values and not pushed on to CHILD profile #299
- bigip_ssl_certificate and bigip_ssl_key import not supported properly #428
- bigip_bigiq_as3 resource fails with multiple application #437
- bigip_ltm_virtual_address erroring out while terraform apply #432

# Enhancements:

- Added support for ignore_metadata for bigip_as3 resource
- Terraform resource for bigip_ipsec_policy
- Terraform resource for bigip_tunnel
- Terraform resource for bigip_traffic_selector

## 1.7.0 (February 18, 2021)

# Bug Fixes:

- Unable to import Policy with "." in the name #407
- empty tenant_list deletes all tenants #423
- bigip_ltm_profile_client_ssl - Changes to a custom PARENT profile values and not pushed on to CHILD profile
- bigip_ltm_profile_fasthttp - Changes to a custom PARENT profile values and not pushed on to CHILD profile
- bigip_ltm_profile_server_ssl - Changes to a custom PARENT profile values and not pushed on to CHILD profile

# Enhancements:

- [FEATURE REQUEST] Node data resource #421
- Add support to define resources inside directories #411
- Documentation updates

## 1.6.0 (January 7, 2021)

# Bug Fixes:

- bigip_event_service_discovery does not reconcile state #400
- Import of bigip_as3 resource results in complete as3 definition in state #385
- Unable to edit a policy after creation #384
- Unable to remove policy from virtual server #383
- Unable to set priority group on pool member #381
- bigip_ltm_persistence_profile_srcaddr - Changes to a custom PARENT profile values and not pushed on to CHILD profile
- bigip_ltm_persistence_profile_cookie - Changes to a custom PARENT profile values and not pushed on to CHILD profile
- Fixing import tests for ltm_virtual_server resource

# Enhancements:

- Data source for bigip pool
- Documentation updates for as3 resource

## 1.5.0 (November 26, 2020)

# Bug Fixes:

- bigip_ltm_profile_oneconnect - Changes to a custom PARENT profile values and not pushed on to CHILD profile
- bigip_ltm_profile_httpcompress - Changes to a custom PARENT profile values and not pushed on to CHILD profile

# Enhancements:

- [Feature] Provide data resources #349 (data resources for iRule, datagroup, monitor, ssl certs)
- Add support for event driven service discovery #343
- ssl forward proxy configuration is not supported #341
- bigip_common_license_manage_bigiq resource fails to locate Purchased Pool #376
- Add feature - bigip_ltm_persistence_profile_ssl - cookie_method #351
- Documentation fixes

## 1.4.0 (October 26, 2020)

# Bug Fixes:

- bigip_ltm_pool_attachment Can't upgrade to 1.3.3 from 1.3.2 #361
- bigip_ltm_policy fails with "does not have update access to partition (Common)" #333
- Feature Request: Add Pool Priority Group Activation #16
- bigip_ltm_profile_http2 - Changes to a custom PARENT profile values and not pushed on to CHILD profile #299
- bigip_ltm_profile_tcp - Changes to a custom PARENT profile values and not pushed on to CHILD profile #299

# Enhancements:

- ltm_pool_attachment resource will support the legacy behaviour of taking node reference from ltm_node along with new behaviour of node taking input directly in the format ip:port/fqdn:port

## 1.3.3 (October 15, 2020)

# Bug Fixes:

- Cannot create a virtual server that differs in only the source address #317
- Monitor: destination port is ignored #319
- resourceBigipLtmVirtualAddressRead/VirtualAddresses() failed to parse json response if RouteAdvertisement is "selective" #323
- Rerunning terraform - cant update node #336

# Other Fixes:

- [BUG][DOC]bigip_ltm_policy example in doc fail

## 1.3.2 (September 3, 2020)

# Bug Fixes:

- Profile_fastL4 - Changes to a custom PARENT profile values and not pushed on to CHILD profile #299
- Documentation updates
- DO Resource enhancements

## 1.3.1 (August 18, 2020)

- Terraform v0.13 Support
- Repo move from terraform-providers to F5Networks Organization

## 1.3.0 (July 24, 2020)

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

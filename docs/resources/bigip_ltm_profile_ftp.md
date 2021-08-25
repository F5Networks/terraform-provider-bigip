---
   layout: "bigip"
   page_title: "BIG-IP: bigip_ltm_profile_ftp"
   subcategory: "Local Traffic Manager(LTM)"
   description: |-
     Provides details about bigip_ltm_profile_ftp resource
---

# bigip\_ltm\_profile_ftp

`bigip_ltm_profile_ftp` Configures a custom profile_ftp.

Resources should be named with their "full path". The full path is the combination of the partition + name (example: /Common/my-pool ) or  partition + directory + name of the resource  (example: /Common/test/my-pool )

## Example Usage


### For Bigip versions (14.x - 16.x)

```hcl
resource "bigip_ltm_profile_ftp" "sanjose-ftp-profile" {
  name                     = "/Common/sanjose-ftp-profile"
  defaults_from            = "/Common/ftp"
  port                     = 2020
  description              = "test-tftp-profile"
  ftps_mode                = "allow"
  enforce_tlssession_reuse = "enabled"
  allow_active_mode        = "enabled"
}

```      

### For Bigip versions (12.x - 13.x)

```hcl
resource "bigip_ltm_profile_ftp" "sanjose-ftp-profile" {
  name               = "/Common/sanjose-ftp-profile"
  defaults_from      = "/Common/ftp"
  port               = 2020
  description        = "test-tftp-profile"
  allow_ftps         = "enabled"
  translate_extended = "enabled"
}

```


## Argument Reference

* `name` (Required) Name of the profile_ftp

* `partition` - (Optional) Displays the administrative partition within which this profile resides

* `defaults_from` - (Optional) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.



### Arguments which are updated/available in Bigip versions (14.x - 16.x) for FTP profile.For more information refer below KB article
https://support.f5.com/csp/article/K08859735

* `ftps_mode` - (Optional) Specifies if you want to Disallow, Allow, or Require FTPS mode. The default is Disallow

* `enforce_tlssession_reuse` - (Optional) Specifies, when selected (enabled), that the system enforces the data connection to reuse a TLS session. The default value is unchecked (disabled)

* `allow_active_mode` - (Optional)Specifies, when selected (enabled), that the system allows FTP Active Transfer mode. The default value is enabled



### Arguments which are updated/available in Bigip versions (12.x - 13.x) for FTP profile.For more information refer below KB article
https://support.f5.com/csp/article/K13044205

* `allow_ftps` - (Optional)Allow explicit FTPS negotiation. The default is disabled.When enabled (selected), that the system allows explicit FTPS negotiation for SSL or TLS. 

* `translate_extended` - (Optional)Specifies, when selected (enabled), that the system uses ensures compatibility between IP version 4 and IP version 6 clients and servers when using the FTP protocol. The default is selected (enabled).



## Common Arguments for all versions

* `security` - (Optional)Specifies, when checked (enabled), that the system inspects FTP traffic for security vulnerabilities using an FTP security profile. This option is available only on systems licensed for BIG-IP ASM.

* `port` - (Optional)Allows you to configure the FTP service to run on an alternate port. The default is 20.

* `log_profile` - (Optional)Configures the ALG log profile that controls logging

* `log_publisher` - (Optional)Configures the log publisher that handles events logging for this profile

*  `inherit_parent_profile` - (Optional)Enables the FTP data channel to inherit the TCP profile used by the control channel.If disabled,the data channel uses FastL4 only.

* `description` - (Optional)User defined description for FTP profile



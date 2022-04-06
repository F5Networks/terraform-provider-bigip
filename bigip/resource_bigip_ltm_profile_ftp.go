/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"regexp"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmProfileFtp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileFtpCreate,
		Update: resourceBigipLtmProfileFtpUpdate,
		Read:   resourceBigipLtmProfileFtpRead,
		Delete: resourceBigipLtmProfileFtpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the FTP Profile",
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Use the parent ftp profile",
			},
			"allow_ftps": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Allows explicit FTPS negotiation",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description",
			},
			"app_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application service to which the object belongs.",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "name of partition",
			},
			"inherit_parent_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Enables the FTP data channel to inherit the TCP profile used by the control channel.If disabled,the data channel uses FastL4 only.",
			},
			"inherit_vlan_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "inherent vlan list",
			},
			"log_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Configures the ALG log profile that controls logging",
			},
			"log_publisher": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Configures the log publisher that handles events logging for this profile",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				//	Computed:    true,
				Default:     "20",
				Description: "Specifies a service for the data channel port used for this FTP profile. The default port is ftp-data.",
			},
			"security": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables secure FTP traffic for the BIG-IP Application Security Manager. You can set the security option only if the system is licensed for the BIG-IP Application Security Manager. The default value is disabled.",
			},
			"ftps_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disallow",
				Description: "Allows explicit FTPS negotiation",
			},
			"enforce_tlssession_reuse": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Specifies, when selected (enabled), that the system enforces the data connection to reuse a TLS session. The default value is unchecked (disabled).",
			},
			"allow_active_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "Specifies, when selected (enabled), that the system allows FTP Active Transfer mode. The default value is enabled.",
			},
			"translate_extended": {
				Type:     schema.TypeString,
				Optional: true,
				//	Computed:     true,
				Default:     "enabled",
				Description: "This setting is enabled by default, and thus, automatically translates RFC 2428 extended requests EPSV and EPRT to PASV and PORT when communicating with IPv4 servers.",
			},
		},
	}
}

func resourceBigipLtmProfileFtpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	ver, err := client.BigipVersion()
	if err != nil {
		log.Printf("[ERROR] Unable to get bigip version  (%v)", err)
		return err
	}
	bigipversion := ver.Entries.HTTPSLocalhostMgmtTmCliVersion0.NestedStats.Entries.Active.Description
	re := regexp.MustCompile(`^(12)|(13).*`)
	matchresult := re.MatchString(bigipversion)
	regversion := re.FindAllString(bigipversion, -1)

	if !matchresult {
		log.Printf("[DEBUG] Bigip version is : %s", regversion)
		ftpProfileConfig := &bigip.Ftp{
			Name:                  name,
			AllowFtps:             d.Get("allow_ftps").(string),
			AppService:            d.Get("app_service").(string),
			DefaultsFrom:          d.Get("defaults_from").(string),
			Description:           d.Get("description").(string),
			InheritParentProfile:  d.Get("inherit_parent_profile").(string),
			InheritVlanList:       d.Get("inherit_vlan_list").(string),
			LogProfile:            d.Get("log_profile").(string),
			LogPublisher:          d.Get("log_publisher").(string),
			Port:                  d.Get("port").(int),
			Security:              d.Get("security").(string),
			FtpsMode:              d.Get("ftps_mode").(string),
			EnforceTlsSesionReuse: d.Get("enforce_tlssession_reuse").(string),
			AllowActiveMode:       d.Get("allow_active_mode").(string),
			TranslateExtended:     d.Get("translate_extended").(string),
		}

		log.Println("[INFO] Creating FTP profile")
		err := client.CreateFtp(ftpProfileConfig)
		if err != nil {
			log.Printf("[ERROR] Unable to Create ftp Profile  (%s) (%v)", name, err)
			return err
		}
	} else {
		log.Printf("[DEBUG] Bigip version is : %s", regversion)
		ftpProfileConfig := &bigip.Ftp{
			Name:                 name,
			AllowFtps:            d.Get("allow_ftps").(string),
			AppService:           d.Get("app_service").(string),
			DefaultsFrom:         d.Get("defaults_from").(string),
			Description:          d.Get("description").(string),
			InheritParentProfile: d.Get("inherit_parent_profile").(string),
			InheritVlanList:      d.Get("inherit_vlan_list").(string),
			LogProfile:           d.Get("log_profile").(string),
			LogPublisher:         d.Get("log_publisher").(string),
			Port:                 d.Get("port").(int),
			Security:             d.Get("security").(string),
			TranslateExtended:    d.Get("translate_extended").(string),
		}
		log.Println("[INFO] Creating FTP profile")
		err := client.CreateFtp(ftpProfileConfig)
		if err != nil {
			log.Printf("[ERROR] Unable to Create ftp Profile  (%s) (%v)", name, err)
			return err
		}
	}
	d.SetId(name)
	return resourceBigipLtmProfileFtpRead(d, meta)
}

func resourceBigipLtmProfileFtpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	ver, err := client.BigipVersion()
	if err != nil {
		log.Printf("[ERROR] Unable to get bigip version  (%v)", err)
		return err
	}
	bigipversion := ver.Entries.HTTPSLocalhostMgmtTmCliVersion0.NestedStats.Entries.Active.Description
	re := regexp.MustCompile(`^(12)|(13).*`)
	matchresult := re.MatchString(bigipversion)
	regversion := re.FindAllString(bigipversion, -1)

	if !matchresult {
		log.Printf("[DEBUG] Bigip version is : %s", regversion)
		log.Println("[INFO] Updating TCP Profile Route " + name)
		ftpProfileConfig := &bigip.Ftp{
			Name:                  name,
			AllowFtps:             d.Get("allow_ftps").(string),
			AppService:            d.Get("app_service").(string),
			DefaultsFrom:          d.Get("defaults_from").(string),
			Description:           d.Get("description").(string),
			InheritParentProfile:  d.Get("inherit_parent_profile").(string),
			InheritVlanList:       d.Get("inherit_vlan_list").(string),
			LogProfile:            d.Get("log_profile").(string),
			LogPublisher:          d.Get("log_publisher").(string),
			Port:                  d.Get("port").(int),
			Security:              d.Get("security").(string),
			FtpsMode:              d.Get("ftps_mode").(string),
			EnforceTlsSesionReuse: d.Get("enforce_tlssession_reuse").(string),
			AllowActiveMode:       d.Get("allow_active_mode").(string),
			TranslateExtended:     d.Get("translate_extended").(string),
		}
		err := client.ModifyFtp(name, ftpProfileConfig)
		if err != nil {
			return fmt.Errorf("Error create profile ftp (%s): %s ", name, err)
		}
	} else {
		log.Printf("[DEBUG] Bigip version is : %s", regversion)
		log.Println("[INFO] Updating TCP Profile Route " + name)
		ftpProfileConfig := &bigip.Ftp{
			Name:                 name,
			AllowFtps:            d.Get("allow_ftps").(string),
			AppService:           d.Get("app_service").(string),
			DefaultsFrom:         d.Get("defaults_from").(string),
			Description:          d.Get("description").(string),
			InheritParentProfile: d.Get("inherit_parent_profile").(string),
			InheritVlanList:      d.Get("inherit_vlan_list").(string),
			LogProfile:           d.Get("log_profile").(string),
			LogPublisher:         d.Get("log_publisher").(string),
			Port:                 d.Get("port").(int),
			Security:             d.Get("security").(string),
			TranslateExtended:    d.Get("translate_extended").(string),
		}
		err := client.ModifyFtp(name, ftpProfileConfig)
		if err != nil {
			return fmt.Errorf("Error create profile ftp (%s): %s ", name, err)
		}
	}
	return resourceBigipLtmProfileFtpRead(d, meta)
}

func resourceBigipLtmProfileFtpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetFtp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve ftp Profile  (%s) (%v)", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] ftp  Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", obj.DefaultsFrom)

	ver, err := client.BigipVersion()
	if err != nil {
		log.Printf("[ERROR] Unable to get bigip version  (%v)", err)
		return err
	}
	bigipversion := ver.Entries.HTTPSLocalhostMgmtTmCliVersion0.NestedStats.Entries.Active.Description
	re := regexp.MustCompile(`^(12)|(13).*`)
	matchresult := re.MatchString(bigipversion)
	regversion := re.FindAllString(bigipversion, -1)

	if !matchresult {
		log.Printf("[DEBUG] Bigip version is : %s", regversion)

		if _, ok := d.GetOk("ftps_mode"); ok {
			_ = d.Set("ftps_mode", obj.FtpsMode)
		}

		if _, ok := d.GetOk("enforce_tlssession_reuse"); ok {
			_ = d.Set("enforce_tlssession_reuse", obj.EnforceTlsSesionReuse)
		}

		if _, ok := d.GetOk("allow_active_mode"); ok {
			_ = d.Set("allow_active_mode", obj.AllowActiveMode)
		}
	}

	if _, ok := d.GetOk("allow_ftps"); ok {
		_ = d.Set("allow_ftps", obj.AllowFtps)
	}

	if _, ok := d.GetOk("app_service"); ok {
		_ = d.Set("app_service", obj.AppService)
	}

	if _, ok := d.GetOk("description"); ok {
		_ = d.Set("description", obj.Description)
	}

	if _, ok := d.GetOk("inherit_parent_profile"); ok {
		_ = d.Set("inherit_parent_profile", obj.InheritParentProfile)
	}

	if _, ok := d.GetOk("log_profile"); ok {
		_ = d.Set("log_profile", obj.LogProfile)
	}

	if _, ok := d.GetOk("inherit_vlan_list"); ok {
		_ = d.Set("inherit_vlan_list", obj.InheritVlanList)
	}

	if _, ok := d.GetOk("log_publisher"); ok {
		_ = d.Set("log_publisher", obj.LogPublisher)
	}

	if _, ok := d.GetOk("port"); ok {
		_ = d.Set("port", obj.Port)
	}

	if _, ok := d.GetOk("security"); ok {
		_ = d.Set("security", obj.Security)
	}

	if _, ok := d.GetOk("translate_extended"); ok {
		_ = d.Set("translate_extended", obj.TranslateExtended)
	}

	return nil
}

func resourceBigipLtmProfileFtpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Ftp Profile " + name)

	err := client.DeleteFtp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete ftp Profile (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}

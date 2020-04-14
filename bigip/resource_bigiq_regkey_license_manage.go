/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"reflect"
)

func resourceBigiqLicenseManage() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigiqLicenseManageCreate,
		Read:   resourceBigiqLicenseManageRead,
		Update: resourceBigiqLicenseManageUpdate,
		Delete: resourceBigiqLicenseManageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bigiq_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"bigiq_user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"bigiq_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The registration key pool to use",
			},
			"bigiq_password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"bigiq_token_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_TOKEN_AUTH", nil),
			},
			"bigiq_login_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tmos",
				Description: "Login reference for token authentication (see BIG-IQ REST docs for details)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_LOGIN_REF", nil),
			},
			"pool": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key that you want to assign from the pool",
			},
			"managed": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the specified device is a managed or un-managed device",
			},
			"device": {
				Type:         schema.TypeString,
				Required:     true,
				ConfigMode:   schema.SchemaConfigModeAttr,
				ValidateFunc: getDevicedetails,
				Description:  "When managed is yes, specifies the managed device, or device UUID, that you want to register.",
			},
			"device_username": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: getDevicedetails,
				Description:  "The username used to connect to the remote device.",
			},
			"device_password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: getDevicedetails,
				Description:  "The password of the device_username.When managed is no, this parameter is required.",
			},
		},
	}
}

func resourceBigiqLicenseManageCreate(d *schema.ResourceData, meta interface{}) error {
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	meta = bigiqRef
	//log.Printf("bigiqRef = %+v", meta)
	client := meta.(*bigip.BigIP)
	poolName := d.Get("pool").(string)
	regKey := d.Get("key").(string)
	poolId, err := client.GetRegkeyPoolId(poolName)
	if err != nil && poolId == "" {
		log.Printf("Getting PoolID failed with :%v", err)
		return err
	}
	//log.Printf("Pool ID = %+v", poolId)
	deRef := bigip.DeviceRef{
		Link: "https://localhost/mgmt/shared/resolver/device-groups/cm-bigip-allBigIpDevices/devices/5c1e6fa1-ae98-4d65-b7c4-2872c21d5fa3",
	}
	config := &bigip.ManagedDevice{
		DeviceReference: deRef,
	}
	resp, err := client.RegkeylicenseAssign(config, poolId, regKey)
	if err != nil {
		log.Printf("Assigning License failed from regKey Pool:%v", err)
		return err
	}
	log.Printf("Resp from Post = %+v", resp)
	d.SetId(resp.ID)
	return resourceBigiqLicenseManageRead(d, meta)
}

func resourceBigiqLicenseManageRead(d *schema.ResourceData, meta interface{}) error {
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	meta = bigiqRef
	client := meta.(*bigip.BigIP)
	log.Printf("meta in Read = %+v", meta)
	memID := d.Id()
	//log.Printf("bigiqRef = %+v", bigiqRef)
	poolName := d.Get("pool").(string)
	regKey := d.Get("key").(string)
	poolId, err := client.GetRegkeyPoolId(poolName)
	if err != nil && poolId == "" {
		log.Printf("Getting PoolID failed with :%v", err)
		return err
	}
	client.GetMemberStatus(poolId, regKey, memID)
	return nil
}

func resourceBigiqLicenseManageUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	return nil

}

func resourceBigiqLicenseManageDelete(d *schema.ResourceData, meta interface{}) error {
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	meta = bigiqRef
	client := meta.(*bigip.BigIP)
	log.Printf("meta in delete = %+v", meta)
	memID := d.Id()
	poolName := d.Get("pool").(string)
	regKey := d.Get("key").(string)
	poolId, err := client.GetRegkeyPoolId(poolName)
	if err != nil && poolId == "" {
		log.Printf("Getting PoolID failed with :%v", err)
		return err
	}
	//log.Printf("Pool ID = %+v", poolId)
	client.RegkeylicenseRevoke(poolId, regKey, memID)
	return nil
}

func connectBigIq(d *schema.ResourceData) (*bigip.BigIP, error) {
	bigiqConfig := Config{
		Address:  d.Get("bigiq_address").(string),
		Port:     d.Get("bigiq_port").(string),
		Username: d.Get("bigiq_user").(string),
		Password: d.Get("bigiq_password").(string),
	}
	if d.Get("bigiq_token_auth").(bool) {
		bigiqConfig.LoginReference = d.Get("bigiq_login_ref").(string)
	}
	//log.Printf("bigiqConfig = %+v", bigiqConfig)
	return bigiqConfig.Client()
}
func getDevicedetails(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	log.Printf("Value type:%T and Value:%+v", value, value)
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
		break
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateF5Name", reflect.TypeOf(value)))
	}

	for _, v := range values {
		log.Printf("---Value type:%T and Value:%+v", v, v)
		//match, _ := regexp.MatchString("^/[\\w_\\-.]+/[\\w_\\-.]+$", v)
		//if !match {
		//	errors = append(errors, fmt.Errorf("%q must match /Partition/Name and contain letters, numbers or [._-]. e.g. /Common/my-pool", field))
		//}
	}
	return
}

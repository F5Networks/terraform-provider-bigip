/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	//"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
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
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_token_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Sensitive:   true,
				Default:     false,
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_TOKEN_AUTH", nil),
			},
			"bigiq_login_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "tmos",
				Description: "Login reference for token authentication (see BIG-IQ REST docs for details)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_LOGIN_REF", nil),
			},
			"pool_license_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePoolLicenseType,
				Description:  "This will specify Utility/regKey Licence pool type",
			},
			"assignment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAssignmentType,
				Description:  "Whether the specified device is a managed or un-managed device",
			},
			"pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The registration key pool to use",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The registration key that you want to assign from the pool",
			},
			"unit_of_measure": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the rate at which this license usage is billed",
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
	bigipRef := meta.(*bigip.BigIP)
	log.Printf("bigipRef = %+v", bigipRef)
	var poolId, regKey string
	if v, ok := d.GetOk("pool"); ok {
		poolId, err = bigiqRef.GetRegkeyPoolId(v.(string))
		if (err != nil) && (poolId == "") {
			log.Printf("Getting PoolID failed with :%v", err)
			return err
		}
	}
	if v, ok := d.GetOk("key"); ok {
		regKey = v.(string)
	}
	log.Printf("Pool ID = %+v", poolId)
	deviceID, err := bigiqRef.GetDeviceId("10.145.65.170")
	if (err != nil) && (deviceID == "") {
		log.Printf("Getting deviceID failed with :%v", err)
		return err
	}
	log.Printf("deviceID = %+v", deviceID)
	//Link: "https://localhost/mgmt/shared/resolver/device-groups/cm-bigip-allBigIpDevices/devices/5c1e6fa1-ae98-4d65-b7c4-2872c21d5fa3",
	deRef := bigip.DeviceRef{
		Link: deviceID,
	}
	config := &bigip.ManagedDevice{
		DeviceReference: deRef,
	}
	resp, err := bigiqRef.RegkeylicenseAssign(config, poolId, regKey)
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
	bigipRef := meta.(*bigip.BigIP)
	log.Printf("bigipRef = %+v", bigipRef)
	memID := d.Id()
	//log.Printf("bigiqRef = %+v", bigiqRef)
	poolName := d.Get("pool").(string)
	regKey := d.Get("key").(string)
	poolId, err := bigiqRef.GetRegkeyPoolId(poolName)
	if err != nil && poolId == "" {
		log.Printf("Getting PoolID failed with :%v", err)
		return err
	}
	bigiqRef.GetMemberStatus(poolId, regKey, memID)
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
	//meta = bigiqRef
	bigipRef := meta.(*bigip.BigIP)
	log.Printf("bigipRef = %+v", bigipRef)
	memID := d.Id()
	poolName := d.Get("pool").(string)
	regKey := d.Get("key").(string)
	poolId, err := bigiqRef.GetRegkeyPoolId(poolName)
	if err != nil && poolId == "" {
		log.Printf("Getting PoolID failed with :%v", err)
		return err
	}
	//log.Printf("Pool ID = %+v", poolId)
	bigiqRef.RegkeylicenseRevoke(poolId, regKey, memID)
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

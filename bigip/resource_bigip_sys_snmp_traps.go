/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// this module does not have DELETE function as there is no API for Delete
func resourceBigipSysSnmpTraps() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysSnmpTrapsCreate,
		UpdateContext: resourceBigipSysSnmpTrapsUpdate,
		ReadContext:   resourceBigipSysSnmpTrapsRead,
		DeleteContext: resourceBigipSysSnmpTrapsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name",
			},
			"auth_passwordencrypted": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Encrypted password ",
			},

			"auth_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the protocol used to authenticate the user.",
			},

			"community": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the community string used for this trap. ",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description.",
			},

			"engine_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the authoritative security engine for SNMPv3.",
			},

			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The host the trap will be sent to.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port that the trap will be sent to.",
			},
			"privacy_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the clear text password used to encrypt traffic. This field will not be displayed. ",
			},
			"privacy_password_encrypted": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the encrypted password used to encrypt traffic. ",
			},
			"privacy_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the protocol used to encrypt traffic. ",
			},
			"security_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether or not traffic is encrypted and whether or not authentication is required.",
			},

			"security_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security name used in conjunction with SNMPv3.",
			},

			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SNMP version used for sending the trap. ",
			},
		},
	}

}

func resourceBigipSysSnmpTrapsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	authPasswordEncrypted := d.Get("auth_passwordencrypted").(string)
	authProtocol := d.Get("auth_protocol").(string)
	community := d.Get("community").(string)
	description := d.Get("description").(string)
	engineId := d.Get("engine_id").(string)
	host := d.Get("host").(string)
	port := d.Get("port").(int)
	privacyPassword := d.Get("privacy_password").(string)
	privacyPasswordEncrypted := d.Get("privacy_password_encrypted").(string)
	privacyProtocol := d.Get("privacy_protocol").(string)
	securityLevel := d.Get("security_level").(string)
	securityName := d.Get("security_name").(string)
	version := d.Get("version").(string)

	log.Println("[INFO] Creating Snmp traps ")

	err := client.CreateTRAP(
		name,
		authPasswordEncrypted,
		authProtocol,
		community,
		description,
		engineId,
		host,
		port,
		privacyPassword,
		privacyPasswordEncrypted,
		privacyProtocol,
		securityLevel,
		securityName,
		version,
	)

	if err != nil {
		log.Printf("[ERROR] Unable to Create SNMP trap (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipSysSnmpTrapsRead(ctx, d, meta)
}

func resourceBigipSysSnmpTrapsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating SNMP Traps " + name)

	r := &bigip.TRAP{
		Name:                     name,
		Host:                     d.Get("host").(string),
		AuthPasswordEncrypted:    d.Get("auth_passwordencrypted").(string),
		AuthProtocol:             d.Get("auth_protocol").(string),
		Community:                d.Get("community").(string),
		Description:              d.Get("description").(string),
		EngineId:                 d.Get("engine_id").(string),
		PrivacyPassword:          d.Get("privacy_password").(string),
		PrivacyPasswordEncrypted: d.Get("privacy_password_encrypted").(string),
		PrivacyProtocol:          d.Get("privacy_protocol").(string),
		SecurityLevel:            d.Get("security_level").(string),
		SecurityName:             d.Get("security_name").(string),
		Version:                  d.Get("version").(string),
		Port:                     d.Get("port").(int),
	}

	err := client.ModifyTRAP(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify SNMP trap (%v) ", err)
		return diag.FromErr(err)
	}
	return resourceBigipSysSnmpTrapsRead(ctx, d, meta)
}

func resourceBigipSysSnmpTrapsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	host := d.Id()

	log.Println("[INFO] Reading SNMP traps " + host)

	traps, err := client.TRAPs(host)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve SNMP trap (%v) ", err)
		return diag.FromErr(err)
	}
	if traps == nil {
		log.Printf("[WARN] SNMP traps (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("name", traps.Name)
	_ = d.Set("auth_passwordencrypted", traps.AuthPasswordEncrypted)
	_ = d.Set("auth_protocol", traps.AuthProtocol)
	_ = d.Set("community", traps.Community)
	_ = d.Set("description", traps.Description)
	_ = d.Set("engine_id", traps.EngineId)
	_ = d.Set("host", traps.Host)
	_ = d.Set("port", traps.Port)
	_ = d.Set("privacy_password", traps.PrivacyPassword)
	_ = d.Set("privacy_password_encrypted", traps.PrivacyPasswordEncrypted)
	_ = d.Set("privacy_protocol", traps.PrivacyProtocol)
	_ = d.Set("security_level", traps.SecurityLevel)
	_ = d.Set("security_name", traps.SecurityName)
	_ = d.Set("version", traps.Version)

	return nil
}

func resourceBigipSysSnmpTrapsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting snmp host " + name)

	err := client.DeleteTRAP(name)
	if err != nil {
		log.Printf("[ERROR] Unable to delete SNMP trap (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

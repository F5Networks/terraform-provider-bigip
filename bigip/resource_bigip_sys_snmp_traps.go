package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

// this module does not have DELETE function as there is no API for Delete
func resourceBigipSysSnmpTraps() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysSnmpTrapsCreate,
		Update: resourceBigipSysSnmpTrapsUpdate,
		Read:   resourceBigipSysSnmpTrapsRead,
		Delete: resourceBigipSysSnmpTrapsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name",
				//ValidateFunc: validateF5Name,
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
				Description: "SNMP version used for sending the trap. ",
			},
		},
	}

}

func resourceBigipSysSnmpTrapsCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId(name)
	return resourceBigipSysSnmpTrapsRead(d, meta)
}

func resourceBigipSysSnmpTrapsUpdate(d *schema.ResourceData, meta interface{}) error {
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
	}

	err := client.ModifyTRAP(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify SNMP trap (%v) ", err)
		return err
	}
	return resourceBigipSysSnmpTrapsRead(d, meta)
}

func resourceBigipSysSnmpTrapsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	host := d.Id()

	log.Println("[INFO] Reading SNMP traps " + host)

	traps, err := client.TRAPs()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve SNMP trap (%v) ", err)
		return err
	}
	if traps == nil {
		log.Printf("[WARN] SNMP traps (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", traps.Name)
	if err := d.Set("auth_passwordencrypted", traps.AuthPasswordEncrypted); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AuthPasswordEncrypted to state for Snmp Traps  (%s): %s", d.Id(), err)
	}
	if err := d.Set("auth_protocol", traps.AuthProtocol); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AuthProtocol to state for Snmp Traps (%s): %s", d.Id(), err)
	}
	if err := d.Set("community", traps.Community); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Community to state for Snmp Traps  (%s): %s", d.Id(), err)
	}
	d.Set("description", traps.Description)
	d.Set("engine_id", traps.EngineId)
	if err := d.Set("host", traps.Host); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Host to state for Snmp Traps  (%s): %s", d.Id(), err)
	}
	d.Set("port", traps.Port)
	d.Set("privacy_password", traps.PrivacyPassword)
	if err := d.Set("privacy_password_encrypted", traps.PrivacyPasswordEncrypted); err != nil {
		return fmt.Errorf("[DEBUG] Error saving PrivacyPasswordEncrypted to state for Snmp Traps (%s): %s", d.Id(), err)
	}
	d.Set("privacy_protocol", traps.PrivacyProtocol)
	d.Set("security_level", traps.SecurityLevel)
	d.Set("security_name", traps.SecurityName)
	d.Set("version", traps.Version)

	return nil
}

func resourceBigipSysSnmpTrapsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting snmp host " + name)

	err := client.DeleteTRAP(name)
	if err != nil {
		log.Printf("[ERROR] Unable to delete SNMP trap (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

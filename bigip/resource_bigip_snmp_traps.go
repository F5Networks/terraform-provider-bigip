package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

// this module does not have DELETE function as there is no API for Delete
func resourceBigipLtmSnmpTraps() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmSnmpTrapsCreate,
		Update: resourceBigipLtmSnmpTrapsUpdate,
		Read:   resourceBigipLtmSnmpTrapsRead,
		Delete: resourceBigipLtmSnmpTrapsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmSnmpTrapsImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name",
				//ValidateFunc: validateF5Name,
			},
			"auth_passwordencrypted": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Encrypted password ",
			},

			"auth_protocol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the protocol used to authenticate the user.",
			},

			"community": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the community string used for this trap. ",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description.",
			},

			"engine_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the authoritative security engine for SNMPv3.",
			},

			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The host the trap will be sent to.",
			},
			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port that the trap will be sent to.",
			},
			"privacy_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the clear text password used to encrypt traffic. This field will not be displayed. ",
			},
			"privacy_password_encrypted": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the encrypted password used to encrypt traffic. ",
			},
			"privacy_protocol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the protocol used to encrypt traffic. ",
			},
			"security_level": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether or not traffic is encrypted and whether or not authentication is required.",
			},

			"security_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security name used in conjunction with SNMPv3.",
			},

			"version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SNMP version used for sending the trap. ",
			},
		},
	}

}

func resourceBigipLtmSnmpTrapsCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipLtmSnmpTrapsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating SNMP Traps " + name)

	r := &bigip.TRAP{
		Name: name,
		Host: d.Get("host").(string),
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

	return client.ModifyTRAP(r)
	return nil
}

func resourceBigipLtmSnmpTrapsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	host := d.Id()

	log.Println("[INFO] Reading SNMP traps " + host)

	traps, err := client.TRAPs()
	if err != nil {
		return err
	}

	d.Set("name", traps.Name)
	d.Set("auth_passwordencrypted", traps.AuthPasswordEncrypted)
	d.Set("auth_protocol", traps.AuthProtocol)
	d.Set("community", traps.Community)
	d.Set("description", traps.Description)
	d.Set("engine_id", traps.EngineId)
	d.Set("host", traps.Host)
	d.Set("port", traps.Port)
	d.Set("privacy_password", traps.PrivacyPassword)
	d.Set("privacy_password_encrypted", traps.PrivacyPasswordEncrypted)
	d.Set("privacy_protocol", traps.PrivacyProtocol)
	d.Set("security_level", traps.SecurityLevel)
	d.Set("security_name", traps.SecurityName)
	d.Set("version", traps.Version)

	return nil
}

func resourceBigipLtmSnmpTrapsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting snmp host " + name)

	return client.DeleteTRAP(name)
}

func resourceBigipLtmSnmpTrapsImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

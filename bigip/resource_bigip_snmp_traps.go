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
			"authPasswordEncrypted": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Encrypted password ",
			},

			"authProtocol": &schema.Schema{
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

			"engineId": &schema.Schema{
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
			"privacyPassword": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the clear text password used to encrypt traffic. This field will not be displayed. ",
			},
			"privacyPasswordEncrypted": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the encrypted password used to encrypt traffic. ",
			},
			"privacyProtocol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the protocol used to encrypt traffic. ",
			},
			"securityLevel": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether or not traffic is encrypted and whether or not authentication is required.",
			},

			"securityName": &schema.Schema{
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
	authPasswordEncrypted := d.Get("authPasswordEncrypted").(string)
	authProtocol := d.Get("authProtocol").(string)
	community := d.Get("community").(string)
	description := d.Get("description").(string)
	engineId := d.Get("engineId").(string)
	host := d.Get("host").(string)
	port := d.Get("port").(int)
	privacyPassword := d.Get("privacyPassword").(string)
	privacyPasswordEncrypted := d.Get("privacyPasswordEncrypted").(string)
	privacyProtocol := d.Get("privacyProtocol").(string)
	securityLevel := d.Get("securityLevel").(string)
	securityName := d.Get("securityName").(string)
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
		AuthPasswordEncrypted:    d.Get("authPasswordEncrypted").(string),
		AuthProtocol:             d.Get("authProtocol").(string),
		Community:                d.Get("community").(string),
		Description:              d.Get("description").(string),
		EngineId:                 d.Get("engineId").(string),
		PrivacyPassword:          d.Get("privacyPassword").(string),
		PrivacyPasswordEncrypted: d.Get("privacyPasswordEncrypted").(string),
		PrivacyProtocol:          d.Get("privacyProtocol").(string),
		SecurityLevel:            d.Get("securityLevel").(string),
		SecurityName:             d.Get("securityName").(string),
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
	d.Set("authPasswordEncrypted", traps.AuthPasswordEncrypted)
	d.Set("authProtocol", traps.AuthProtocol)
	d.Set("community", traps.Community)
	d.Set("description", traps.Description)
	d.Set("engineId", traps.EngineId)
	d.Set("host", traps.Host)
	d.Set("port", traps.Port)
	d.Set("privacyPassword", traps.PrivacyPassword)
	d.Set("privacyPasswordEncrypted", traps.PrivacyPasswordEncrypted)
	d.Set("privacyProtocol", traps.PrivacyProtocol)
	d.Set("securityLevel", traps.SecurityLevel)
	d.Set("securityName", traps.SecurityName)
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

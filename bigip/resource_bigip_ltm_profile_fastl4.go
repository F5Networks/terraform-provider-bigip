package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmFastl4() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmFastl4Create,
		Update: resourceBigipLtmFastl4Update,
		Read:   resourceBigipLtmFastl4Read,
		Delete: resourceBigipLtmFastl4Delete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmFastl4Importer,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Fastl4 Profile",
				//ValidateFunc: validateF5Name,
			},
			"partition": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaultsFrom": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"clientTimeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"explicitFlowMigration": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"hardwareSynCookie": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"idleTimeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"ipTosToClient": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"ipTosToServer": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"keepAliveInterval": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
		},
	}

}

func resourceBigipLtmFastl4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	defaultsFrom := d.Get("defaultsFrom").(string)
	clientTimeout := d.Get("clientTimeout").(int)
	explicitFlowMigration := d.Get("explicitFlowMigration").(string)
	hardwareSynCookie := d.Get("hardwareSynCookie").(string)
	idleTimeout := d.Get("idleTimeout").(int)
	ipTosToClient := d.Get("ipTosToClient").(string)
	ipTosToServer := d.Get("ipTosToServer").(string)
	keepAliveInterval := d.Get("keepAliveInterval").(string)

	log.Println("[INFO] Creating Fastl4 profile")

	err := client.CreateFastl4(
		name,
		partition,
		defaultsFrom,
		clientTimeout,
		explicitFlowMigration,
		hardwareSynCookie,
		idleTimeout,
		ipTosToClient,
		ipTosToServer,
		keepAliveInterval,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipLtmFastl4Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Fastl4{
		Name:                  name,
		Partition:             d.Get("partition").(string),
		DefaultsFrom:          d.Get("defaultsFrom").(string),
		ClientTimeout:         d.Get("clientTimeout").(int),
		ExplicitFlowMigration: d.Get("explicitFlowMigration").(string),
		HardwareSynCookie:     d.Get("hardwareSynCookie").(string),
		IdleTimeout:           d.Get("idleTimeout").(int),
		IpTosToClient:         d.Get("ipTosToClient").(string),
		IpTosToServer:         d.Get("ipTosToServer").(string),
		KeepAliveInterval:     d.Get("keepAliveInterval").(string),
	}

	return client.ModifyFastl4(name, r)
}

func resourceBigipLtmFastl4Read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBigipLtmFastl4Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Fastl4 Profile " + name)

	return client.DeleteFastl4(name)
}

func resourceBigipLtmFastl4Importer(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmFastl4() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmFastl4Create,
		Update: resourceBigipLtmFastl4Update,
		Read:   resourceBigipLtmFastl4Read,
		Delete: resourceBigipLtmFastl4Delete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmFastl4Importer,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Fastl4 Profile",
				//ValidateFunc: validateF5Name,
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"client_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"explicitflow_migration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"hardware_syncookie": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"idle_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"iptos_toclient": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"iptos_toserver": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"keepalive_interval": {
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
	defaultsFrom := d.Get("defaults_from").(string)
	clientTimeout := d.Get("client_timeout").(int)
	explicitFlowMigration := d.Get("explicitflow_migration").(string)
	hardwareSynCookie := d.Get("hardware_syncookie").(string)
	idleTimeout := d.Get("idle_timeout").(string)
	ipTosToClient := d.Get("iptos_toclient").(string)
	ipTosToServer := d.Get("iptos_toserver").(string)
	keepAliveInterval := d.Get("keepalive_interval").(string)

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
	return resourceBigipLtmFastl4Read(d, meta)
}

func resourceBigipLtmFastl4Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Fastl4{
		Name:                  name,
		Partition:             d.Get("partition").(string),
		DefaultsFrom:          d.Get("defaults_from").(string),
		ClientTimeout:         d.Get("client_timeout").(int),
		ExplicitFlowMigration: d.Get("explicitflow_migration").(string),
		HardwareSynCookie:     d.Get("hardware_syncookie").(string),
		IdleTimeout:           d.Get("idle_timeout").(string),
		IpTosToClient:         d.Get("iptos_toclient").(string),
		IpTosToServer:         d.Get("iptos_toserver").(string),
		KeepAliveInterval:     d.Get("keepalive_interval").(string),
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

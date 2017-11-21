package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipGtmPool_a() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmPool_aCreate,
		Update: resourceBigipGtmPool_aUpdate,
		Read:   resourceBigipGtmPool_aRead,
		Delete: resourceBigipGtmPool_aDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipGtmPool_aImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the pool_a Servers",
				ValidateFunc: validateF5Name,
			},

			"monitor": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the pool_a Servers",
			},

			"load_balancing_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the pool_a Servers",
			},
			"max_answers_returned": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Name of the pool_a Servers",
			},

			"alternate_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the pool_a Servers",
			},

			"fallback_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the pool_a Servers",
			},

			"fallback_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the pool_a Servers",
			},

			"members": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},

		},
	}

}

func resourceBigipGtmPool_aCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	monitor := d.Get("monitor").(string)
	load_balancing_mode := d.Get("load_balancing_mode").(string)
	max_answers_returned := d.Get("max_answers_returned").(int)
	alternate_mode:= d.Get("alternate_mode").(string)
	fallback_ip := d.Get("fallback_ip").(string)
	fallback_mode := d.Get("fallback_mode").(string)
	members := setToStringSlice(d.Get("members").(*schema.Set))

	log.Println("[INFO] Creating gtm pool ")

	err := client.CreatePool_a(
		name,
		monitor,
		load_balancing_mode,
		max_answers_returned,
		alternate_mode,
		fallback_ip,
		fallback_mode,
		members,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipGtmPool_aUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Pool_a " + name)

	r := &bigip.Pool_a{
		Name: name,
		Monitor:  d.Get("monitor").(string),
		Load_balancing_mode:  d.Get("load_balancing_mode").(string),
		Max_answers_returned:  d.Get("max_answers_returned").(int),
		Alternate_mode:  d.Get("alternate_mode").(string),
		Fallback_ip:  d.Get("Fallback_ip").(string),
		Fallback_mode:  d.Get("fallback_mode").(string),
		Members: setToStringSlice(d.Get("members").(*schema.Set)),
	}

	return client.ModifyPool_a(r)
}

func resourceBigipGtmPool_aRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Pool_a " + name)

	pool_a, err := client.Pool_as()
	if err != nil {
		return err
	}

	d.Set("name", pool_a.Name)
	d.Set("monitor", pool_a.Monitor)
	d.Set("load_balancing_mode", pool_a.Load_balancing_mode)

	d.Set("max_answers_returned", pool_a.Max_answers_returned)
	d.Set("alternate_mode", pool_a.Alternate_mode)
	d.Set("fallback_ip", pool_a.Fallback_ip)
	d.Set("fallback_mode", pool_a.Fallback_mode)
	d.Set("members", pool_a.Members)

	return nil
}

func resourceBigipGtmPool_aDelete(d *schema.ResourceData, meta interface{}) error {
	/* This function is not supported on BIG-IP, you cannot DELETE pool_a API is not supported */
	return nil
}

func resourceBigipGtmPool_aImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

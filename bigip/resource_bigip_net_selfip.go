package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipNetSelfIP() *schema.Resource {

	return &schema.Resource{
		Create: resourceBigipNetSelfIPCreate,
		Read:   resourceBigipNetSelfIPRead,
		Update: resourceBigipNetSelfIPUpdate,
		Delete: resourceBigipNetSelfIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the SelfIP",
			},

			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SelfIP IP address",
			},

			"vlan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the vlan",
			},

			"traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group, defaults to traffic-group-local-only if not specified",
				Default:     "traffic-group-local-only",
			},
		},
	}
}

func resourceBigipNetSelfIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	ip := d.Get("ip").(string)
	vlan := d.Get("vlan").(string)

	log.Println("[INFO] Creating SelfIP ")

	err := client.CreateSelfIP(name, ip, vlan)

	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceBigipNetSelfIPUpdate(d, meta)
}

func resourceBigipNetSelfIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching SelfIP " + name)

	selfIPs, err := client.SelfIPs()
	if err != nil {
		return err
	}
	d.Set("name", name)
	if selfIPs == nil {
		log.Printf("[WARN] SelfIP (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	for _, selfip := range selfIPs.SelfIPs {
		log.Println(selfip.Name)
		if selfip.Name == name {
			return nil
		}
	}

	return nil
}

func resourceBigipNetSelfIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating SelfIP " + name)

	r := &bigip.SelfIP{
		Name:         name,
		Address:      d.Get("ip").(string),
		Vlan:         d.Get("vlan").(string),
		TrafficGroup: d.Get("traffic_group").(string),
	}

	err := client.ModifySelfIP(name, r)
	if err != nil {
		return err
	}

	return resourceBigipNetSelfIPRead(d, meta)

}

func resourceBigipNetSelfIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Deleting selfIP " + name)

	err := client.DeleteSelfIP(name)
	if err != nil {
		return err
	}
	if err == nil {
		log.Printf("[WARN] Selfip (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

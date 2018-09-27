package bigip

import (
	"fmt"
	"log"
	"regexp"

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

	log.Printf("[DEBUG] Creating SelfIP %s", name)

	err := client.CreateSelfIP(name, ip, vlan)

	if err != nil {
		return fmt.Errorf("Error creating SelfIP %s: %v", name, err)
	}

	d.SetId(name)

	return resourceBigipNetSelfIPRead(d, meta)
}

func resourceBigipNetSelfIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[DEBUG] Reading SelfIP %s", name)

	selfIPs, err := client.SelfIPs()
	if err != nil {
		return fmt.Errorf("Error retrieving SelfIPs: %v", err)
	}
	if selfIPs == nil {
		log.Printf("[DEBUG] SelfIPs not found, removing SelfIP %s from state", name)
		d.SetId("")
		return nil
	}

	found := false
	for _, selfip := range selfIPs.SelfIPs {
		if selfip.FullPath == name {
			d.Set("name", selfip.FullPath)
			d.Set("vlan", selfip.Vlan)

			// Extract Self IP address from "(selfip_address)[%route_domain](/mask)" groups 1 + 2
			regex := regexp.MustCompile(`((?:[0-9]{1,3}\.){3}[0-9]{1,3})(?:\%\d+)?(\/\d+)`)
			selfipAddress := regex.FindStringSubmatch(selfip.Address)
			parsedSelfipAddress := selfipAddress[1] + selfipAddress[2]
			d.Set("ip", parsedSelfipAddress)

			// Extract Traffic Group name from the full path (ignoring /Common/ prefix)
			regex = regexp.MustCompile(`\/Common\/(.+)`)
			trafficGroup := regex.FindStringSubmatch(selfip.TrafficGroup)
			parsedTrafficGroup := trafficGroup[1]
			d.Set("traffic_group", parsedTrafficGroup)

			found = true
			break
		}
	}

	if !found {
		log.Printf("[DEBUG] SelfIP %s not found, removing from state", name)
		d.SetId("")
	}

	return nil
}

func resourceBigipNetSelfIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Updating SelfIP %s", name)

	r := &bigip.SelfIP{
		Name:         name,
		Address:      d.Get("ip").(string),
		Vlan:         d.Get("vlan").(string),
		TrafficGroup: d.Get("traffic_group").(string),
	}

	err := client.ModifySelfIP(name, r)
	if err != nil {
		return fmt.Errorf("Error modifying SelfIP %s: %v", name, err)
	}

	return resourceBigipNetSelfIPRead(d, meta)

}

func resourceBigipNetSelfIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[DEBUG] Deleting SelfIP %s", name)

	err := client.DeleteSelfIP(name)
	if err != nil {
		return fmt.Errorf("Error deleting SelfIP %s: %v", name, err)
	}

	d.SetId("")
	return nil
}

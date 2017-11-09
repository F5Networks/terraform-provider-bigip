package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipNetRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipNetRouteCreate,
		Update: resourceBigipNetRouteUpdate,
		Read:   resourceBigipNetRouteRead,
		Delete: resourceBigipNetRouteDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipNetRouteImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the route",
				//ValidateFunc: validateF5Name,
			},

			"network": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination network",
			},

			"gw": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gw address",
			},
		},
	}

}

func resourceBigipNetRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	network := d.Get("network").(string)
	gw := d.Get("gw").(string)

	log.Println("[INFO] Creating Route")

	err := client.CreateRoute(
		name,
		network,
		gw,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipNetRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Route " + name)

	r := &bigip.Route{
		Name:    name,
		Network: d.Get("network").(string),
	}

	return client.ModifyRoute(name, r)
}

func resourceBigipNetRouteRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipNetRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Route " + name)

	return client.DeleteRoute(name)
}

func resourceBigipNetRouteImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

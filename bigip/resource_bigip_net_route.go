package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceBigipNetRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipNetRouteCreate,
		Update: resourceBigipNetRouteUpdate,
		Read:   resourceBigipNetRouteRead,
		Delete: resourceBigipNetRouteDelete,
		Importer: &schema.ResourceImporter{
		 State: schema.ImportStatePassthrough,
	 },


		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the route",
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
	return resourceBigipNetRouteRead(d, meta)
}

func resourceBigipNetRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Route " + name)

	r := &bigip.Route{
		Name:    name,
		Network: d.Get("network").(string),
	}

	err := client.ModifyRoute(name, r)
	if err != nil {
		return err
	}
	return resourceBigipNetRouteRead(d, meta)
}

func resourceBigipNetRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetRoute(name)
	if err != nil {
		d.SetId("")
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Route (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	if err := d.Set("network", obj.Network); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Network to state for Route  (%s): %s", d.Id(), err)
	}
	d.Set("gw", obj.Gateway)
	if err := d.Set("gw", obj.Gateway); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Gateway to state for Route  (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceBigipNetRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Route " + name)

	err := client.DeleteRoute(name)
	if err == nil {
		log.Printf("[WARN] Node (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

 

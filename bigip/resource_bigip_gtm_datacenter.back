package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipGtmDatacenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmDatacenterCreate,
		Update: resourceBigipGtmDatacenterUpdate,
		Read:   resourceBigipGtmDatacenterRead,
		Delete: resourceBigipGtmDatacenterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipGtmDatacenterImporter,
		},

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:    true,
				Description:  "Name of the datacenter",
				ValidateFunc: validateF5Name,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description.",
			},

			"contact": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact person",
			},


			"prober_pool": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional prober pool used to monitor the data center's virtual servers.",
			},
		"app_service": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application service that the object belongs",
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies that the data center and its resources are not available for load balancing.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies that the data center and its resources are available for load balancing.",
			},
		},
	}
}

func resourceBigipGtmDatacenterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	contact := d.Get("contact").(string)
	app_service := d.Get("app_service").(string)
	enabled := d.Get("enabled").(bool)
	disabled := d.Get("disabled").(bool)
	prober_pool := d.Get("prober_pool").(string)
	log.Println("[INFO] Creating Datacenter ", name)

	err := client.CreateDatacenter(
		name,
		description,
		contact,
		app_service,
		enabled,
		disabled,
		prober_pool,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return resourceBigipGtmDatacenterRead(d, meta)
}

func resourceBigipGtmDatacenterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Datacenter " + name)

	datacenter, err := client.Datacenters()
	if err != nil {
		return err
	}

	d.Set("name", datacenter.Name)
	d.Set("description", datacenter.Description)
	d.Set("contact", datacenter.Contact)
	d.Set("prober_pool", datacenter.Prober_pool)

	return nil
}


func resourceBigipGtmDatacenterExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Checking Datacenter " + name + " exists.")

	datacenter, err := client.Datacenters()
	if err != nil {
		return false, err
	}

	d.Set("name", datacenter.Name)
	d.Set("description", datacenter.Description)
	d.Set("contact", datacenter.Contact)
	d.Set("app_service", datacenter.App_service)
	d.Set("enabled", datacenter.Enabled)
	d.Set("disabled", datacenter.Disabled)
	d.Set("prober_pool", datacenter.Prober_pool)

	return datacenter != nil, nil
}

func resourceBigipGtmDatacenterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Datacenter " + name)

	r := &bigip.Datacenter{
		Name:    name,
		Description: d.Get("description").(string),
		Contact: d.Get("contact").(string),
		App_service:  d.Get("app_service").(string),
		Enabled:  d.Get("enabled").(bool),
		Disabled:  d.Get("disabled").(bool),
		Prober_pool:  d.Get("prober_pool").(string),
	}

	return client.ModifyDatacenter(r)
}

func resourceBigipGtmDatacenterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting datacenter " + name)

	return client.DeleteDatacenter(name)
}

func resourceBigipGtmDatacenterImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

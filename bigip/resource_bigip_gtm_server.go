package bigip

import (
	"log"
"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipGtmServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmServerCreate,
		Update: resourceBigipGtmServerUpdate,
		Read:   resourceBigipGtmServerRead,
		Delete: resourceBigipGtmServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipGtmServerImporter,
		},

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:    true,
				Description:  "Name of the datacenter",
				ValidateFunc: validateF5Name,
			},
			"datacenter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description.",
			},

			"monitor": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact person",
			},


			"virtual_server_discovery": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional prober pool used to monitor the data center's virtual servers.",
			},
		"product": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application service that the object belongs",
			},

			"addresses": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule name",
							//ValidateFunc: validateF5Name,
						},
						"device_name": {
	 					 Type:        schema.TypeString,
	 					 Required:    true,
	 					 Description: "Rule name",
	 					 //ValidateFunc: validateF5Name,
	 				 },
					 "translation": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Rule name",
						//ValidateFunc: validateF5Name,
					},
				},
			},
		},

		"virtual_server": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Rule name",
						ValidateFunc: validateF5Name,
					},
					"destination": {
					 Type:        schema.TypeString,
					 Required:    true,
					 Description: "Rule name",
					 //ValidateFunc: validateF5Name,
				 },

			},
		},
		},


		},
	}
}

func resourceBigipGtmServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating Server" + name)

	p := dataToServer(name, d)
	d.SetId(name)
	err := client.CreateGtmserver(&p)
	if err != nil {
		return err
	}
	return resourceBigipGtmServerRead(d, meta)
}

func resourceBigipGtmServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Fetching GTM server " + name)
	p, err := client.GetGtmserver(name)

	if err != nil {
		return err
	}

	return ServerToData(p, d)
}


func resourceBigipGtmServerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Gtmserver " + name)

	p, err := client.GetGtmserver(name)
	if err != nil {
		return false, err
	}

	return p != nil, nil
}

func resourceBigipGtmServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating Gtmserver " + name)
	p := dataToServer(name, d)
	return client.UpdateGtmserver(name, &p)
}

func resourceBigipGtmServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	return client.DeleteGtmserver(name)
}

func resourceBigipGtmServerImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}



func dataToServer(name string, d *schema.ResourceData) bigip.Server {
	var p bigip.Server


	p.Name = name
	p.Datacenter = d.Get("datacenter").(string)
	p.Monitor = d.Get("monitor").(string)
	p.Virtual_server_discovery = d.Get("virtual_server_discovery").(bool)
	p.Product = d.Get("product").(string)
	addressCount := d.Get("addresses.#").(int)
	p.Addresses = make([]bigip.ServerAddresses, 0, addressCount)
	for i := 0; i < addressCount; i++ {
		var  r  bigip.ServerAddresses
		log.Println("I am in dattoserver policy ", p, addressCount, i)
		prefix := fmt.Sprintf("addresses.%d", i)
		r.Name = d.Get(prefix + ".name").(string)

		r.Device_name = d.Get(prefix + ".device_name").(string)
		log.Println("I am in dattoserver policy                                                          ", p, r.Device_name,  prefix)

		r.Translation = d.Get(prefix + ".translation").(string)
		p.Addresses = append(p.Addresses, r)
}

		vsCount := d.Get("virtual_server.#").(int)

		p.GTMVirtual_Server = make([]bigip.VSrecord, 0, vsCount)
		for i := 0; i < vsCount; i++ {
			var k bigip.VSrecord
			prefix := fmt.Sprintf("virtual_server.%d", i)

			k.Name = d.Get(prefix + ".name").(string)

			k.Destination = d.Get(prefix + ".destination").(string)


			p.GTMVirtual_Server = append(p.GTMVirtual_Server, k)
}
log.Println("I am in VS   value of p                                                   ", p)

	return p
}

func ServerToData(p *bigip.Server, d *schema.ResourceData) error {
	d.Set("datacenter", p.Datacenter)
	d.Set("monitor", p.Monitor)
	d.Set("virtual_server_discovery", p.Virtual_server_discovery)
	d.Set("product", p.Product)

	for i, r := range p.Addresses {
		addresses := fmt.Sprintf("addresses.%d", i)
		d.Set(fmt.Sprintf("%s.name", addresses), r.Device_name)
		d.Set(fmt.Sprintf("%s.name", addresses), r.Translation)
	}

	for i, k := range p.GTMVirtual_Server {
		virtual_server := fmt.Sprintf("virtual_server.%d", i)

		d.Set(fmt.Sprintf("%s.name", virtual_server), k.Name)
		d.Set(fmt.Sprintf("%s.name", virtual_server), k.Destination)
		log.Println("I am in GTMVirtual_Server    ", virtual_server)

	}

	return nil
}

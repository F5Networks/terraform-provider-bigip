package bigip

import (
	"fmt"
	"log"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmNodeCreate,
		Read:   resourceBigipLtmNodeRead,
		Update: resourceBigipLtmNodeUpdate,
		Delete: resourceBigipLtmNodeDelete,
		Exists: resourceBigipLtmNodeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the node",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},

			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the node",
				ForceNew:    true,
			},
			"rate_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the maximum number of connections per second allowed for a node or node address. The default value is 'disabled'.",
			},

			"connection_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the maximum number of connections allowed for the node or node address.",
				Default:     0,
			},
			"dynamic_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the dynamic ratio number for the node. Used for dynamic ratio load balancing. ",
				Default:     0,
			},
			"monitor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the monitor or monitor rule that you want to associate with the node.",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "user-up",
				Description: "Marks the node up or down. The default value is user-up.",
			},
			"fqdn": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_family": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the node's address family. The default is 'unspecified', or IP-agnostic",
						},
						"fqdn_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the fully qualified domain name of the node.",
						},
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "3600",
							Description: "Specifies the amount of time before sending the next DNS query.",
						},
						"downinterval": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "5",
							Description: "Specifies the number of attempts to resolve a domain name. The default is 5.",
						},
						"autopopulate": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "disabled",
							Description: "Specifies whether the node should scale to the IP address set returned by DNS.",
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	p := dataToNode(name, d)
	d.Partial(true)
err := client.CreateNode(&p)
if err != nil {
	log.Printf("[ERROR] Unable to Create Node  (%s) (%v) ", name, err)
	return err
}
return nil
	d.SetId(name)
d.SetPartial("fqdn")
d.Partial(false)
	return resourceBigipLtmNodeRead(d, meta)
}

func resourceBigipLtmNodeRead(d *schema.ResourceData, meta interface{}) error {



	return nil
}

func resourceBigipLtmNodeExists(d *schema.ResourceData, meta interface{}) (bool, error) {


	return true, nil
}

func resourceBigipLtmNodeUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting node " + name)
	err := client.DeleteNode(name)

	if err != nil {
		log.Printf("[ERROR] Unable to Delete Node %s  %v : ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func dataToNode(name string, d *schema.ResourceData) bigip.Node {
	var p bigip.Node
	p.Name = name
	p.Address = d.Get("address").(string)
	p.RateLimit = d.Get("rate_limit").(string)
	p.ConnectionLimit = d.Get("connection_limit").(int)
	p.DynamicRatio = d.Get("dynamic_ratio").(int)
	p.Monitor = d.Get("monitor").(string)
	p.State = d.Get("state").(string)
	fqdnCount := d.Get("fqdn.#").(int)
	p.Fqdn = make([]bigip.Fqdnrecord, 0, fqdnCount)
	//p.Fqdn = setToStringSlice(d.Get("fqdn").(*schema.Set))

	for i := 0; i < fqdnCount; i++ {
		var r bigip.Fqdnrecord
		prefix := fmt.Sprintf("fqdn.%d", i)
		r.FqdnName = d.Get(prefix + ".fqdn_name").(string)
		r.Interval = d.Get(prefix + ".interval").(string)
		p.Fqdn = append(p.Fqdn, r)
	}
	log.Println("I am in data to Node value of p                                                   ", p)

	return p
}

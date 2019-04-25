package bigip

import (
	"fmt"
	"log"
	"regexp"

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
				Computed:    true,
			},

			"connection_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the maximum number of connections allowed for the node or node address.",
				Computed:    true,
			},
			"dynamic_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the dynamic ratio number for the node. Used for dynamic ratio load balancing. ",
				Computed:    true,
			},
			"monitor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the monitor or monitor rule that you want to associate with the node.",
				Computed:    true,
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
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the fully qualified domain name of the node.",
						},
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the amount of time before sending the next DNS query.",
						},
						"downinterval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the number of attempts to resolve a domain name. The default is 5.",
						},
						"autopopulate": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
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
	address := d.Get("address").(string)
	rate_limit := d.Get("rate_limit").(string)
	connection_limit := d.Get("connection_limit").(int)
	dynamic_ratio := d.Get("dynamic_ratio").(int)
	monitor := d.Get("monitor").(string)
	state := d.Get("state").(string)

	r, _ := regexp.Compile("^((?:[0-9]{1,3}.){3}[0-9]{1,3})|(.*:.*)$")

	log.Println("[INFO] Creating node " + name + "::" + address)
	var err error
	if r.MatchString(address) {
		err = client.CreateNode(
			name,
			address,
			rate_limit,
			connection_limit,
			dynamic_ratio,
			monitor,
			state,
		)
	} else {
		interval := d.Get("fqdn.0.interval").(string)
		address_family := d.Get("fqdn.0.address_family").(string)
		autopopulate := d.Get("fqdn.0.autopopulate").(string)
		downinterval := d.Get("fqdn.0.downinterval").(int)

		err = client.CreateFQDNNode(
			name,
			address,
			rate_limit,
			connection_limit,
			dynamic_ratio,
			monitor,
			state,
			interval,
			address_family,
			autopopulate,
			downinterval,
		)
	}

	if err != nil {
		return fmt.Errorf("Error modifying node %s: %v", name, err)
	}

	d.SetId(name)

	return resourceBigipLtmNodeRead(d, meta)
}

func resourceBigipLtmNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching node " + name)

	node, err := client.GetNode(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve node %s  %v :", name, err)
		return err
	}
	if node == nil {
		log.Printf("[WARN] Node (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if node.FQDN.Name != "" {
		if err := d.Set("address", node.FQDN.Name); err != nil {
			return fmt.Errorf("[DEBUG] Error saving address to state for Node (%s): %s", d.Id(), err)
		}
	} else {
		// xxx.xxx.xxx.xxx(%x)
		regex := regexp.MustCompile(`((?:[0-9]{1,3}\.){3}[0-9]{1,3})(?:\%\d+)?`)
		address := regex.FindStringSubmatch(node.Address)
		if err := d.Set("address", address[1]); err != nil {
			return fmt.Errorf("[DEBUG] Error saving address to state for Node (%s): %s", d.Id(), err)
		}
	}
	d.Set("name", name)
	if err := d.Set("monitor", node.Monitor); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Monitor to state for Node (%s): %s", d.Id(), err)
	}
	if err := d.Set("rate_limit", node.RateLimit); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Monitor to state for Node (%s): %s", d.Id(), err)
	}

	d.Set("connection_limit", node.ConnectionLimit)
	d.Set("dynamic_ratio", node.DynamicRatio)
	d.Set("fqdn.0.interval", node.FQDN.Interval)
	d.Set("fqdn.0.downinterval", node.FQDN.DownInterval)
	d.Set("fqdn.0.autopopulate", node.FQDN.AutoPopulate)
	d.Set("fqdn.0.address_family", node.FQDN.AddressFamily)

	return nil
}

func resourceBigipLtmNodeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching node " + name)

	node, err := client.GetNode(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve node %s  %v :", name, err)
		return false, err
	}

	if node == nil {
		log.Printf("[WARN] node (%s) not found, removing from state", d.Id())
		return false, nil
	}
	return node != nil, nil
}

func resourceBigipLtmNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	address := d.Get("address").(string)
	r, _ := regexp.Compile("^((?:[0-9]{1,3}.){3}[0-9]{1,3})|(.*:.*)$")

	var node *bigip.Node
	if r.MatchString(address) {
		node = &bigip.Node{
			Address:         address,
			ConnectionLimit: d.Get("connection_limit").(int),
			DynamicRatio:    d.Get("dynamic_ratio").(int),
			Monitor:         d.Get("monitor").(string),
			RateLimit:       d.Get("rate_limit").(string),
			State:           d.Get("state").(string),
		}
	} else {
		node = &bigip.Node{
			ConnectionLimit: d.Get("connection_limit").(int),
			DynamicRatio:    d.Get("dynamic_ratio").(int),
			Monitor:         d.Get("monitor").(string),
			RateLimit:       d.Get("rate_limit").(string),
			State:           d.Get("state").(string),
		}

		err := client.ModifyNode(name, node)
		if err != nil {
			return fmt.Errorf("Error modifying node %s: %v", name, err)
		}
	}
	return resourceBigipLtmNodeRead(d, meta)
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

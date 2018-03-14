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
			State: resourceBigipLtmNodeImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the node",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},

			"address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the node",
				ForceNew:    true,
			},
			"rate_limit": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the maximum number of connections per second allowed for a node or node address. The default value is 'disabled'.",
			},

			"connection_limit": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the maximum number of connections allowed for the node or node address.",
				Default:     0,
			},
			"dynamic_ratio": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the dynamic ratio number for the node. Used for dynamic ratio load balancing. ",
				Default:     0,
			},
			"monitor": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the monitor or monitor rule that you want to associate with the node.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "user-up",
				Description: "Marks the node up or down. The default value is user-up.",
			},
			"fqdn": {
				Type:     schema.TypeList,
				Optional: true,
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
		err = client.CreateFQDNNode(
			name,
			address,
			rate_limit,
			connection_limit,
			dynamic_ratio,
			monitor,
			state,
		)
	}
	if err != nil {
		return err
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
		if err := d.Set("address", node.Address); err != nil {
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
	return nil
}

func resourceBigipLtmNodeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching node " + name)

	node, err := client.GetNode(name)
	if err != nil {
		return false, err
	}

	if node == nil {
		d.SetId("")
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
			State:       d.Get("state").(string),
		}
	} else {
		node = &bigip.Node{
			ConnectionLimit: d.Get("connection_limit").(int),
			DynamicRatio:    d.Get("dynamic_ratio").(int),
			Monitor:         d.Get("monitor").(string),
			RateLimit:       d.Get("rate_limit").(string),
			State:       d.Get("state").(string),
		}
		node.FQDN.Name = address
	}

	err := client.ModifyNode(name, node)
	if err != nil {
		return nil
	}
	return resourceBigipLtmNodeRead(d, meta)
}

func resourceBigipLtmNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting node " + name)

	err := client.DeleteNode(name)
	if err == nil {
		log.Printf("[WARN] Node (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	regex := regexp.MustCompile("referenced by a member of pool '\\/\\w+/([\\w-_.]+)")
	for err != nil {
		log.Printf("[INFO] Deleting %s from pools...\n", name)
		parts := regex.FindStringSubmatch(err.Error())
		if len(parts) > 1 {
			poolName := parts[1]
			members, e := client.PoolMembers(poolName)
			if e != nil {
				return e
			}
			for _, member := range members.PoolMembers {
				e = client.DeletePoolMember(poolName, member.Name)
				if e != nil {
					return e
				}
			}
			err = client.DeleteNode(name)
		} else {
			break
		}
	}
	return err
}

func resourceBigipLtmNodeImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

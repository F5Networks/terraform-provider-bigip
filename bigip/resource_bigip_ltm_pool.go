package bigip

import (
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

var NODE_VALIDATION = regexp.MustCompile(":\\d{2,5}$")

func resourceBigipLtmPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPoolCreate,
		Read:   resourceBigipLtmPoolRead,
		Update: resourceBigipLtmPoolUpdate,
		Delete: resourceBigipLtmPoolDelete,
		Exists: resourceBigipLtmPoolExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigIpLtmPoolImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the pool",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"nodes": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Optional:    true,
				Description: "Nodes to add to the pool. Format node_name:port. e.g. node01:443",
			},

			"monitors": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Optional:    true,
				Description: "Assign monitors to a pool.",
			},

			"allow_nat": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow NAT",
			},

			"allow_snat": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow SNAT",
			},

			"load_balancing_mode": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "round-robin",
				Description: "Possible values: round-robin, ...",
			},
		},
	}
}

func resourceBigipLtmPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating pool " + name)
	err := client.CreatePool(name)
	if err != nil {
		return err
	}
	d.SetId(name)

	err = resourceBigipLtmPoolUpdate(d, meta)
	if err != nil {
		client.DeletePool(name)
		return err
	}

	return resourceBigipLtmPoolRead(d, meta)
}

func resourceBigipLtmPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading pool " + name)

	pool, err := client.GetPool(name)
	if err != nil {
		return err
	}
	nodes, err := client.PoolMembers(name)
	if err != nil {
		return err
	}

	d.Set("allow_nat", pool.AllowNAT)
	d.Set("allow_snat", pool.AllowSNAT)
	d.Set("load_balancing_mode", pool.LoadBalancingMode)
	d.Set("nodes", makeStringSet(&nodes))
	d.Set("name", name)

	monitors := strings.Split(strings.TrimSpace(pool.Monitor), " and ")
	d.Set("monitors", makeStringSet(&monitors))

	return nil
}

func resourceBigipLtmPoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Checking pool " + name + " exists.")

	pool, err := client.GetPool(name)
	if err != nil {
		return false, err
	}

	if pool == nil {
		d.SetId("")
	}

	return pool != nil, nil
}

func resourceBigipLtmPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//monitors
	var monitors []string
	if m, ok := d.GetOk("monitors"); ok {
		for _, monitor := range m.(*schema.Set).List() {
			monitors = append(monitors, monitor.(string))
		}
	}

	pool := &bigip.Pool{
		AllowNAT:          d.Get("allow_nat").(bool),
		AllowSNAT:         d.Get("allow_snat").(bool),
		LoadBalancingMode: d.Get("load_balancing_mode").(string),
		Monitor:           strings.Join(monitors, " and "),
	}

	err := client.ModifyPool(name, pool)
	if err != nil {
		return err
	}

	//members
	nodes, err := client.PoolMembers(name)
	if err != nil {
		return err
	}
	existing := makeStringSet(&nodes)
	incoming := d.Get("nodes").(*schema.Set)
	delete := existing.Difference(incoming)
	add := incoming.Difference(existing)
	if delete.Len() > 0 {
		for _, d := range delete.List() {
			client.DeletePoolMember(name, d.(string))
		}
	}
	if add.Len() > 0 {
		for _, d := range add.List() {
			client.AddPoolMember(name, d.(string))
		}
	}

	return nil
}

func resourceBigipLtmPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting pool " + name)

	return client.DeletePool(name)
}

func resourceBigIpLtmPoolImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

package bigip

import (
	"log"
	"regexp"
	"fmt"
	"strings"

	"github.com/DealerDotCom/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

var NODE_VALIDATION = regexp.MustCompile(":\\d{2,5}$")

func resourceBigipLtmPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPoolCreate,
		Read:   resourceBigipLtmPoolRead,
		Update: resourceBigipLtmPoolUpdate,
		Delete: resourceBigipLtmPoolDelete,
		Exists: resourceBigipLtmPoolExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Name of the pool",
				ForceNew: true,
			},

			"nodes": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "Nodes to add to the pool. Format node_name:port. e.g. node01:443",
			},

			"monitor": &schema.Schema{
				Type:     schema.TypeString,
				//Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "Assign a monitor to a pool.",
			},

			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: DEFAULT_PARTITION,
				Description: "LTM Partition",
				ForceNew: true,
			},

			"allow_nat": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
				Description: "Allow NAT",
			},

			"allow_snat": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
				Description: "Allow SNAT",
			},

			"load_balancing_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "round-robin",
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
	partition := pool.Partition
	if partition == "" {
		partition = DEFAULT_PARTITION
	}

	d.Set("name", pool.Name)
	d.Set("partition", partition)
	d.Set("allow_nat", pool.AllowNAT)
	d.Set("allow_snat", pool.AllowSNAT)
	d.Set("load_balancing_mode", pool.LoadBalancingMode)
	d.Set("nodes", nodes)
	parts := strings.Split(strings.TrimSpace(pool.Monitor),"/")
	d.Set("monitor", parts[len(parts)-1])
//	monitors := strings.Split(pool.Monitor, " and ")
//	for i := range monitors {
//		// strip off the partition name
//		parts := strings.Split(monitors[i],"/")
//		monitors[i] = strings.TrimSpace(parts[len(parts)-1])
//	}
//	sort.Sort(sort.StringSlice(monitors))
//	d.Set("monitors", monitors)

	return nil;
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

	pool := &bigip.Pool{
		Name: name,
		AllowNAT: d.Get("allow_nat").(bool),
		AllowSNAT: d.Get("allow_snat").(bool),
		LoadBalancingMode: d.Get("load_balancing_mode").(string),
		Monitor: d.Get("monitor").(string),
		//Partition: d.Get("partition").(string),
	}

	err := client.ModifyPool(name, pool)
	if err != nil {
		return err
	}

	//monitors
//	m := d.Get("monitors").([]interface{})
//	if len(m) > 0 {
//		for _, monitor := range m {
//			fmt.Println("Adding monitor " + monitor.(string))
//			client.AddMonitorToPool(monitor.(string), name)
//		}
//	}

	//members
	existing_nodes, err := client.PoolMembers(name)
	existing_set := mapify(existing_nodes)

	incoming_list := d.Get("nodes").([]interface{})
	if len(incoming_list) > 0 {
		for i := range incoming_list {
			incoming := incoming_list[i].(string)
			if _, ok := existing_set[incoming]; ok {
				delete(existing_set, incoming)
			} else {
				if !NODE_VALIDATION.MatchString(incoming) {
					return fmt.Errorf("%s must match spec <node_name>:<port>", incoming);
				}
				err := client.AddPoolMember(name, incoming)
				if err != nil {
					return err
				}
			}
		}
	}
	for key, _ := range existing_set {
		err := client.DeletePoolMember(name, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func mapify(s []string) (map[string]struct{}) {
	set := make(map[string] struct{}, len(s))
	for _, s := range s {
		set[s] = struct{}{}
	}
	return set
}

func resourceBigipLtmPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting pool " + name)

	return client.DeletePool(name)
}

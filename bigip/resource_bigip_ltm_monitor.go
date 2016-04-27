package bigip

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmMonitorCreate,
		Read:   resourceBigipLtmMonitorRead,
		Update: resourceBigipLtmMonitorUpdate,
		Delete: resourceBigipLtmMonitorDelete,
		Exists: resourceBigipLtmMonitorExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the monitor",
				ForceNew:    true,
			},

			"parent": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateParent,
				ForceNew:     true,
				Description:  "Existing monitor to inherit from. Must be one of http, https, icmp or gateway-icmp.",
			},

			"interval": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Check interval in seconds",
				Default:     3,
			},

			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout in seconds",
				Default:     16,
			},

			"send": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "GET /\\r\\n",
				Description: "Request string to send.",
			},

			"receive": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expected response string.",
			},

			"receive_disable": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expected response string.",
			},

			"partition": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     DEFAULT_PARTITION,
				Description: "LTM Partition",
				ForceNew:    true,
			},

			"reverse": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"transparent": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"manual_resume": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ip_dscp": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"time_until_up": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Time in seconds",
			},
		},
	}
}

func resourceBigipLtmMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Creating monitor " + name)

	client.CreateMonitor(
		name,
		d.Get("parent").(string),
		d.Get("interval").(int),
		d.Get("timeout").(int),
		d.Get("send").(string),
		d.Get("receive").(string),
	)

	d.SetId(name)

	resourceBigipLtmMonitorUpdate(d, meta)
	return resourceBigipLtmMonitorRead(d, meta)
}

func resourceBigipLtmMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	monitors, err := client.Monitors()
	if err != nil {
		return err
	}
	for _, m := range monitors {
		if m.Name == name {
			d.Set("interval", m.Interval)
			d.Set("timeout", m.Timeout)
			d.Set("send", m.SendString)
			d.Set("receive", m.ReceiveString)
			d.Set("receive_disable", m.ReceiveDisable)
			d.Set("partition", m.Partition)
			d.Set("reverse", m.Reverse)
			d.Set("transparent", m.Transparent)
			d.Set("ip_dscp", m.IPDSCP)
			d.Set("time_until_up", m.TimeUntilUp)
			d.Set("manual_resume", m.ManualResume)
			return nil
		}
	}
	return nil
}

func resourceBigipLtmMonitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching monitor " + name)

	monitors, err := client.Monitors()
	if err != nil {
		return false, err
	}

	for _, m := range monitors {
		if m.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func resourceBigipLtmMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	m := &bigip.Monitor{
		Interval:       d.Get("interval").(int),
		Timeout:        d.Get("timeout").(int),
		SendString:     d.Get("send").(string),
		ReceiveString:  d.Get("receive").(string),
		ReceiveDisable: d.Get("receive_disable").(string),
		Reverse:        d.Get("reverse").(bool),
		Transparent:    d.Get("transparent").(bool),
		IPDSCP:         d.Get("ip_dscp").(int),
		TimeUntilUp:    d.Get("time_until_up").(int),
		ManualResume:   d.Get("manual_resume").(bool),
	}

	return client.ModifyMonitor(name, d.Get("parent").(string), m)
}

func resourceBigipLtmMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	parent := d.Get("parent").(string)
	log.Println("[Info] Deleting monitor " + name + "::" + parent)
	return client.DeleteMonitor(name, parent)
}

func validateParent(v interface{}, k string) ([]string, []error) {
	p := v.(string)

	if p == "http" || p == "https" || p == "icmp" || p == "gateway-icmp" || p == "tcp" {
		return nil, nil
	}

	return nil, []error{fmt.Errorf("parent must be one of http, https, icmp or gateway-icmp")}
}

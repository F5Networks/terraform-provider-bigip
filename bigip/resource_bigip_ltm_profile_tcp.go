package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmTcp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmTcpCreate,
		Update: resourceBigipLtmTcpUpdate,
		Read:   resourceBigipLtmTcpRead,
		Delete: resourceBigipLtmTcpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmTcpImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the TCP Profile",
				//ValidateFunc: validateF5Name,
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent tcp profile",
			},

			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "idle_timeout can be given value",
			},

			"close_wait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "close wait timer integer",
			},

			"finwait_2timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "timer integer",
			},

			"finwait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "fin wait timer integer",
			},

			"keepalive_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "keepalive_interval timer integer",
			},

			"deferred_accept": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defferred accept",
			},
			"fast_open": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "fast_open value ",
			},
		},
	}

}

func resourceBigipLtmTcpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	idleTimeout := d.Get("idle_timeout").(int)
	closeWaitTimeout := d.Get("close_wait_timeout").(int)
	finWait_2Timeout := d.Get("finwait_2timeout").(int)
	finWaitTimeout := d.Get("finwait_timeout").(int)
	keepAliveInterval := d.Get("keepalive_interval").(int)
	deferredAccept := d.Get("deferred_accept").(string)
	fastOpen := d.Get("fast_open").(string)
	log.Println("[INFO] Creating TCP profile")

	err := client.CreateTcp(
		name,
		partition,
		defaultsFrom,
		idleTimeout,
		closeWaitTimeout,
		finWait_2Timeout,
		finWaitTimeout,
		keepAliveInterval,
		deferredAccept,
		fastOpen,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipLtmTcpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Tcp{
		Name:              name,
		Partition:         d.Get("partition").(string),
		DefaultsFrom:      d.Get("defaults_from").(string),
		IdleTimeout:       d.Get("idle_timeout").(int),
		CloseWaitTimeout:  d.Get("close_wait_timeout").(int),
		FinWait_2Timeout:  d.Get("finwait_2timeout").(int),
		FinWaitTimeout:    d.Get("finwait_timeout").(int),
		KeepAliveInterval: d.Get("keepalive_interval").(int),
		DeferredAccept:    d.Get("deferred_accept").(string),
		FastOpen:          d.Get("fast_open").(string),
	}

	return client.ModifyTcp(name, r)
}

func resourceBigipLtmTcpRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmTcpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Tcp Profile " + name)

	return client.DeleteTcp(name)
}

func resourceBigipLtmTcpImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmTcp() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmTcpCreate,
		Update: resourceBigipLtmTcpUpdate,
		Read:   resourceBigipLtmTcpRead,
		Delete: resourceBigipLtmTcpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmTcpImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the TCP Profile",
				//ValidateFunc: validateF5Name,
			},
			"partition": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaultsFrom": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent oneconnect profile",
			},

			"idleTimeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "idleTimeout can be given value",
			},

			"closeWaitTimeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "close wait timer integer",
			},

			"finWait_2Timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "timer integer",
			},

			"finWaitTimeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "fin wait timer integer",
			},

			"keepAliveInterval": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "keepAliveInterval timer integer",
			},

			"deferredAccept": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defferred accept",
			},
			"fastOpen": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "fastopen value ",
			},
		},
	}

}

func resourceBigipLtmTcpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	defaultsFrom := d.Get("defaultsFrom").(string)
	idleTimeout := d.Get("idleTimeout").(int)
	closeWaitTimeout := d.Get("closeWaitTimeout").(int)
	finWait_2Timeout := d.Get("finWait_2Timeout").(int)
	finWaitTimeout := d.Get("finWaitTimeout").(int)
	keepAliveInterval := d.Get("keepAliveInterval").(int)
	deferredAccept := d.Get("deferredAccept").(string)
	fastOpen := d.Get("fastOpen").(string)
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
		DefaultsFrom:      d.Get("defaultsFrom").(string),
		IdleTimeout:       d.Get("idleTimeout").(int),
		CloseWaitTimeout:  d.Get("closeWaitTimeout").(int),
		FinWait_2Timeout:  d.Get("finWait_2Timeout").(int),
		FinWaitTimeout:    d.Get("finWaitTimeout").(int),
		KeepAliveInterval: d.Get("keepAliveInterval").(int),
		DeferredAccept:    d.Get("deferredAccept").(string),
		FastOpen:          d.Get("fastOpen").(string),
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

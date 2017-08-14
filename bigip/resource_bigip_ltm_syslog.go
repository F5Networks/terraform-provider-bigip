package bigip

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmSyslog() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmSyslogCreate,
		Update: resourceBigipLtmSyslogUpdate,
		Read:   resourceBigipLtmSyslogRead,
		Delete: resourceBigipLtmSyslogDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmSyslogImporter,
		},

		Schema: map[string]*schema.Schema{
			"authPrivFrom": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "authPrivFrom port",
			},

			"remoteServers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of syslog Server",
						},

						"host": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Destination syslog host",
						},

						"remotePort": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "RemotePort port",
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmSyslogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	authPrivFrom := d.Get("authPrivFrom").(string)
	remoteServers := listToSlice(d)

	log.Println("[INFO] Creating Syslog servers ")

	err := client.CreateSyslog(
		authPrivFrom,
		remoteServers,
	)

	if err != nil {
		return err
	}

	return nil
}

func resourceBigipLtmSyslogUpdate(d *schema.ResourceData, meta interface{}) error {
	/*client := meta.(*bigip.BigIP)

	servers := d.Id()

	log.Println("[INFO] Updating Syslog " + description)

	r := &bigip.Syslog{
		remoteServers: setToStringSlice(d.Get("remoteServers").(*schema.Set)),
	}

	return client.ModifySyslog(r)  */
	return nil
}

func resourceBigipLtmSyslogRead(d *schema.ResourceData, meta interface{}) error {
	/*	client := meta.(*bigip.BigIP)

		servers := d.Id()

		log.Println("[INFO] Reading Syslog " + description)

		syslog, err := client.Syslogs()
		if err != nil {
			return err
		}

		d.Set("remoteServers", syslog.Servers)
	*/
	return nil
}

func resourceBigipLtmSyslogDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmSyslogImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func listToSlice(d *schema.ResourceData) []bigip.RemoteServer {
	remoteServerCount := d.Get("remoteServers.#").(int)
	var r = make([]bigip.RemoteServer, remoteServerCount, remoteServerCount)

	for i := 0; i < remoteServerCount; i++ {
		prefix := fmt.Sprintf("remoteServers.%d", i)
		r[i].Name = d.Get(prefix + ".name").(string)
		r[i].Host = d.Get(prefix + ".host").(string)
		r[i].RemotePort = d.Get(prefix + ".remotePort").(int)
	}

	return r
}

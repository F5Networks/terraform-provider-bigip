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
			"auth_privfrom": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "authPrivFrom port",
			},

			"remote_servers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Destination syslog host",
						},
						"host": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Destination syslog host",
						},

						"remote_port": &schema.Schema{
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
	log.Println("[INFO] Creating Syslog servers.")

	r := dataToSysLog(d)
	log.Println("[INFO] Data to Syslog.", r)

	err := client.CreateSyslog(&r)
	if err != nil {
		return err
	}

	return nil
}

func resourceBigipLtmSyslogUpdate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Updating Syslog " + name)

	// r := dataToSysLog(d)
	//
	// return client.ModifySyslog(&r)
	return nil
}

func resourceBigipLtmSyslogRead(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*bigip.BigIP)
	// s, err := client.Syslogs()
	// log.Println("[INFO] Reading Syslog.")
	//
	// if err != nil {
	// 	return err
	// }
	//
	// return syslogToData(s, d)
	return nil
}

func resourceBigipLtmSyslogDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmSyslogImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func dataToSysLog(d *schema.ResourceData) bigip.Syslog {
	var r bigip.Syslog

	r.AuthPrivFrom = d.Get("auth_privfrom").(string)
	remoteServerCount := d.Get("remote_servers.#").(int)
	r.RemoteServers = make([]bigip.RemoteServer, remoteServerCount, remoteServerCount)

	for i := 0; i < remoteServerCount; i++ {
		prefix := fmt.Sprintf("remote_servers.%d", i)
		r.RemoteServers[i].Host = d.Get(prefix + ".host").(string)
		r.RemoteServers[i].Name = d.Get(prefix + ".name").(string)
		r.RemoteServers[i].RemotePort = d.Get(prefix + ".remote_port").(int)
	}

	return r
}

func syslogToData(p *bigip.Syslog, d *schema.ResourceData) error {
	d.Set("auth_privfrom", p.AuthPrivFrom)

	for i, r := range p.RemoteServers {
		prefix := fmt.Sprintf("remote_servers.%d", i)
		d.Set(fmt.Sprintf("%s.host", prefix), r.Host)
		d.Set(fmt.Sprintf("%s.name", prefix), r.Name)
		d.Set(fmt.Sprintf("%s.remote_port", prefix), r.RemotePort)
	}
	return nil
}

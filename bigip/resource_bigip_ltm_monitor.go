package bigip

import (
	"fmt"
	"log"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmMonitorCreate,
		Read:   resourceBigipLtmMonitorRead,
		Update: resourceBigipLtmMonitorUpdate,
		Delete: resourceBigipLtmMonitorDelete,
		Exists: resourceBigipLtmMonitorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the monitor",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},

			"parent": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateParent,
				ForceNew:     true,
				Description:  "Existing monitor to inherit from. Must be one of /Common/http, /Common/https, /Common/icmp or /Common/gateway-icmp.",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Existing monitor to inherit from. Must be one of /Common/http, /Common/https, /Common/icmp or /Common/gateway-icmp.",
			},

			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Check interval in seconds",
				Default:     3,
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout in seconds",
				Default:     16,
			},

			"send": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Request string to send.",
				StateFunc: func(s interface{}) string {
					return strings.Replace(s.(string), "\r\n", "\\r\\n", -1)
				},
			},

			"receive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expected response string.",
			},

			"receive_disable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expected response string.",
			},

			"reverse": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"transparent": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "disabled",
			},

			"manual_resume": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "disabled",
			},

			"ip_dscp": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"time_until_up": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Time in seconds",
			},

			"destination": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*:*",
				Description: "Alias for the destination",
			},
			"compatibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies, when enabled, that the SSL options setting (in OpenSSL) is set to ALL. The default value is enabled.",
				ValidateFunc: validateEnabledDisabled,
			},
			"filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the full path and file name of the file that the system attempts to download. The health check is successful if the system can download the file.",
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the data transfer process (DTP) mode. The default value is passive.",
			},
			"adaptive": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "ftp adaptive",
			},
			"adaptive_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     200,
				Description: "Integer value",
			},
		},
	}
}

func resourceBigipLtmMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	log.Println("[INFO] Creating monitor " + name + " :: " + monitorParent(d.Get("parent").(string)))

	err := client.CreateMonitor(
		name,
		monitorParent(d.Get("parent").(string)),
		d.Get("defaults_from").(string),
		d.Get("interval").(int),
		d.Get("timeout").(int),
		d.Get("send").(string),
		d.Get("receive").(string),
		d.Get("receive_disable").(string),
		d.Get("compatibility").(string),
	)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Monitor (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)

	resourceBigipLtmMonitorUpdate(d, meta)
	return resourceBigipLtmMonitorRead(d, meta)
}

func resourceBigipLtmMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	monitors, err := client.Monitors()
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Monitor (%s) (%v) ", name, err)
		return err
	}
	if monitors == nil {
		log.Printf("[WARN] Monitor (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	for _, m := range monitors {
		if m.FullPath == name {
			d.Set("defaults_from", m.DefaultsFrom)
			d.Set("interval", m.Interval)
			d.Set("timeout", m.Timeout)
			if err := d.Set("send", m.SendString); err != nil {
				return fmt.Errorf("[DEBUG] Error saving SendString to state for Monitor (%s): %s", d.Id(), err)
			}
			if err := d.Set("receive", m.ReceiveString); err != nil {
				return fmt.Errorf("[DEBUG] Error saving ReceiveString to state for Monitor (%s): %s", d.Id(), err)
			}
			d.Set("receive_disable", m.ReceiveDisable)
			d.Set("reverse", m.Reverse)
			d.Set("transparent", m.Transparent)
			d.Set("ip_dscp", m.IPDSCP)
			d.Set("time_until_up", m.TimeUntilUp)
			d.Set("manual_resume", m.ManualResume)
			d.Set("destination", m.Destination)
			d.Set("compatibility", m.Compatibility)
			d.Set("filename", m.Filename)
			d.Set("mode", m.Mode)
			d.Set("adaptive", m.Adaptive)
			d.Set("adaptive_limit", m.AdaptiveLimit)
			d.Set("name", name)

			return nil
		}
	}
	return fmt.Errorf("Couldn't find monitor %s", name)

}

func resourceBigipLtmMonitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching monitor " + name)

	monitors, err := client.Monitors()
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Monitor (%s) (%v) ", name, err)
		return false, err
	}
	if monitors == nil {
		log.Printf("[WARN] Monitor (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}
	for _, m := range monitors {
		if m.FullPath == name {
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
		Reverse:        d.Get("reverse").(string),
		Transparent:    d.Get("transparent").(string),
		IPDSCP:         d.Get("ip_dscp").(int),
		TimeUntilUp:    d.Get("time_until_up").(int),
		ManualResume:   d.Get("manual_resume").(string),
		Destination:    d.Get("destination").(string),
		Compatibility:  d.Get("compatibility").(string),
		Filename:       d.Get("filename").(string),
		Mode:           d.Get("mode").(string),
		Adaptive:       d.Get("adaptive").(string),
		AdaptiveLimit:  d.Get("adaptive_limit").(int),
	}

	err := client.ModifyMonitor(name, monitorParent(d.Get("parent").(string)), m)
	if err != nil {
		log.Printf("[ERROR] Unable to Update Monitor (%s) (%v) ", name, err)
		return err
	}

	return resourceBigipLtmMonitorRead(d, meta)
}

func resourceBigipLtmMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	parent := monitorParent(d.Get("parent").(string))
	log.Println("[Info] Deleting monitor " + name + "::" + parent)
	err := client.DeleteMonitor(name, parent)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Monitor (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func validateParent(v interface{}, k string) ([]string, []error) {
	p := v.(string)
	if p == "/Common/http" || p == "/Common/https" || p == "/Common/icmp" || p == "/Common/gateway-icmp" || p == "/Common/tcp" || p == "/Common/tcp-half-open" || p == "/Common/ftp" {
		return nil, nil
	}

	return nil, []error{fmt.Errorf("parent must be one of /Common/http, /Common/https, /Common/icmp, /Common/gateway-icmp, /Common/tcp-half-open, /Common/tcp, /Common/ftp")}
}

func monitorParent(s string) string {
	return strings.TrimPrefix(s, "/Common/")
}

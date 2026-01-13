package bigip

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipGtmServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmServerCreate,
		ReadContext:   resourceBigipGtmServerRead,
		UpdateContext: resourceBigipGtmServerUpdate,
		DeleteContext: resourceBigipGtmServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the GTM server",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				ForceNew:    true,
				Description: "Partition or tenant the server belongs to",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Datacenter the server belongs to",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the GTM server",
			},
			"product": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "bigip",
				Description: "Server type (bigip, generic-host, etc.)",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable or disable the GTM server",
			},
			"addresses": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configures IP addresses for the server",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address",
						},
						"device_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Device name for the address",
						},
						"translation": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "none",
							Description: "IP translation address",
						},
					},
				},
			},
			"monitor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Monitor assigned to the server",
			},
			"virtual_server_discovery": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "Virtual server discovery mode (enabled, disabled, enabled-no-delete)",
			},
			"link_discovery": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Link discovery mode (enabled, disabled, enabled-no-delete)",
			},
			"prober_preference": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "inherit",
				Description: "Prober preference (inside-datacenter, outside-datacenter, inherit, pool)",
			},
			"prober_fallback": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "inherit",
				Description: "Prober fallback (any-available, inside-datacenter, outside-datacenter, inherit, pool)",
			},
			"prober_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prober pool to use when prober_preference or prober_fallback is set to pool",
			},
			"expose_route_domains": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow the GTM server to expose route domains",
			},
			"iq_allow_path": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable iQuery path probing",
			},
			"iq_allow_service_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable iQuery service checking",
			},
			"iq_allow_snmp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable iQuery SNMP",
			},
			"limit_cpu_usage": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "CPU usage limit",
			},
			"limit_cpu_usage_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CPU usage limit status",
			},
			"limit_max_bps": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum bits per second",
			},
			"limit_max_bps_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Maximum bps status",
			},
			"limit_max_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum concurrent connections",
			},
			"limit_max_connections_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Maximum connections status",
			},
			"limit_max_pps": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum packets per second",
			},
			"limit_max_pps_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Maximum pps status",
			},
			"limit_mem_avail": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Available memory limit",
			},
			"limit_mem_avail_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Available memory status",
			},
		},
	}
}

func resourceBigipGtmServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[INFO] Creating GTM Server: %s", fullPath)

	server := &bigip.Server{
		Name:       name,
		Datacenter: d.Get("datacenter").(string),
		Product:    d.Get("product").(string),
	}

	// Handle addresses
	if v, ok := d.GetOk("addresses"); ok {
		addresses := v.([]interface{})
		server.Addresses = make([]bigip.ServerAddresses, len(addresses))
		for i, addr := range addresses {
			addrMap := addr.(map[string]interface{})
			server.Addresses[i] = bigip.ServerAddresses{
				Name:        addrMap["name"].(string),
				Device_name: addrMap["device_name"].(string),
				Translation: addrMap["translation"].(string),
			}
		}
	}

	if v, ok := d.GetOk("monitor"); ok {
		server.Monitor = v.(string)
	}

	server.Virtual_server_discovery = d.Get("virtual_server_discovery").(string)

	err := client.CreateGtmserver(server)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM Server (%s): %s", fullPath, err))
	}

	d.SetId(fullPath)

	return resourceBigipGtmServerUpdate(ctx, d, meta)
}

func resourceBigipGtmServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	fullPath := d.Id()
	log.Printf("[INFO] Reading GTM Server: %s", fullPath)

	// Parse partition and name from fullPath
	parts := strings.Split(strings.TrimPrefix(fullPath, "/"), "/")
	var name string
	if len(parts) == 2 {
		d.Set("partition", parts[0])
		name = parts[1]
	} else {
		name = fullPath
	}

	server, err := client.GetGtmserver(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve GTM Server %s: %v", fullPath, err)
		return diag.FromErr(err)
	}

	if server == nil {
		log.Printf("[WARN] GTM Server (%s) not found, removing from state", fullPath)
		d.SetId("")
		return nil
	}

	d.Set("name", server.Name)
	d.Set("datacenter", server.Datacenter)
	d.Set("description", server.Description)
	d.Set("product", server.Product)
	d.Set("enabled", server.Enabled)
	d.Set("monitor", server.Monitor)
	d.Set("virtual_server_discovery", server.Virtual_server_discovery)
	d.Set("link_discovery", server.LinkDiscovery)
	d.Set("prober_preference", server.ProberPreference)
	d.Set("prober_fallback", server.ProberFallback)
	d.Set("prober_pool", server.ProberPool)

	// Set boolean fields based on yes/no string values from API
	d.Set("expose_route_domains", server.ExposeRouteDomains == "yes")
	d.Set("iq_allow_path", server.IqAllowPath == "yes")
	d.Set("iq_allow_service_check", server.IqAllowServiceCheck == "yes")
	d.Set("iq_allow_snmp", server.IqAllowSnmp == "yes")

	// Set limit fields
	d.Set("limit_cpu_usage", server.LimitCpuUsage)
	d.Set("limit_cpu_usage_status", server.LimitCpuUsageStatus)
	d.Set("limit_max_bps", server.LimitMaxBps)
	d.Set("limit_max_bps_status", server.LimitMaxBpsStatus)
	d.Set("limit_max_connections", server.LimitMaxConnections)
	d.Set("limit_max_connections_status", server.LimitMaxConnectionsStatus)
	d.Set("limit_max_pps", server.LimitMaxPps)
	d.Set("limit_max_pps_status", server.LimitMaxPpsStatus)
	d.Set("limit_mem_avail", server.LimitMemAvail)
	d.Set("limit_mem_avail_status", server.LimitMemAvailStatus)

	// Handle addresses
	if len(server.Addresses) > 0 {
		addresses := make([]interface{}, len(server.Addresses))
		for i, addr := range server.Addresses {
			addresses[i] = map[string]interface{}{
				"name":        addr.Name,
				"device_name": addr.Device_name,
				"translation": addr.Translation,
			}
		}
		d.Set("addresses", addresses)
	}

	return nil
}

func resourceBigipGtmServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	fullPath := d.Id()
	name := d.Get("name").(string)

	log.Printf("[INFO] Updating GTM Server: %s", fullPath)

	server := &bigip.Server{
		Name:        name,
		Datacenter:  d.Get("datacenter").(string),
		Description: d.Get("description").(string),
		Product:     d.Get("product").(string),
	}

	// Handle enabled/disabled state
	enabled := d.Get("enabled").(bool)
	server.Enabled = enabled
	server.Disabled = !enabled

	// Handle addresses
	if v, ok := d.GetOk("addresses"); ok {
		addresses := v.([]interface{})
		server.Addresses = make([]bigip.ServerAddresses, len(addresses))
		for i, addr := range addresses {
			addrMap := addr.(map[string]interface{})
			server.Addresses[i] = bigip.ServerAddresses{
				Name:        addrMap["name"].(string),
				Device_name: addrMap["device_name"].(string),
				Translation: addrMap["translation"].(string),
			}
		}
	}

	if v, ok := d.GetOk("monitor"); ok {
		server.Monitor = v.(string)
	}

	server.Virtual_server_discovery = d.Get("virtual_server_discovery").(string)
	server.LinkDiscovery = d.Get("link_discovery").(string)
	server.ProberPreference = d.Get("prober_preference").(string)
	server.ProberFallback = d.Get("prober_fallback").(string)

	if v, ok := d.GetOk("prober_pool"); ok {
		server.ProberPool = v.(string)
	}

	// Convert boolean fields to yes/no strings for API
	if d.Get("expose_route_domains").(bool) {
		server.ExposeRouteDomains = "yes"
	} else {
		server.ExposeRouteDomains = "no"
	}
	if d.Get("iq_allow_path").(bool) {
		server.IqAllowPath = "yes"
	} else {
		server.IqAllowPath = "no"
	}
	if d.Get("iq_allow_service_check").(bool) {
		server.IqAllowServiceCheck = "yes"
	} else {
		server.IqAllowServiceCheck = "no"
	}
	if d.Get("iq_allow_snmp").(bool) {
		server.IqAllowSnmp = "yes"
	} else {
		server.IqAllowSnmp = "no"
	}

	// Set limit fields
	server.LimitCpuUsage = d.Get("limit_cpu_usage").(int)
	server.LimitCpuUsageStatus = d.Get("limit_cpu_usage_status").(string)
	server.LimitMaxBps = d.Get("limit_max_bps").(int)
	server.LimitMaxBpsStatus = d.Get("limit_max_bps_status").(string)
	server.LimitMaxConnections = d.Get("limit_max_connections").(int)
	server.LimitMaxConnectionsStatus = d.Get("limit_max_connections_status").(string)
	server.LimitMaxPps = d.Get("limit_max_pps").(int)
	server.LimitMaxPpsStatus = d.Get("limit_max_pps_status").(string)
	server.LimitMemAvail = d.Get("limit_mem_avail").(int)
	server.LimitMemAvailStatus = d.Get("limit_mem_avail_status").(string)

	err := client.UpdateGtmserver(name, server)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM Server (%s): %s", fullPath, err))
	}

	return resourceBigipGtmServerRead(ctx, d, meta)
}

func resourceBigipGtmServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting GTM Server: %s", name)

	err := client.DeleteGtmserver(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting GTM Server (%s): %s", name, err))
	}

	d.SetId("")
	return nil
}

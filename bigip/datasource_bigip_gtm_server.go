package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipGtmServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipGtmServerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the GTM server",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Datacenter the server belongs to",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the GTM server",
			},
			"product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server type (bigip, generic-host, etc.)",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the server is enabled",
			},
			"monitor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Monitor assigned to the server",
			},
			"virtual_server_discovery": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtual server discovery mode",
			},
			"link_discovery": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Link discovery mode",
			},
			"prober_preference": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prober preference",
			},
			"prober_fallback": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prober fallback",
			},
			"prober_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prober pool",
			},
			"addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IP addresses for the server",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address",
						},
						"device_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Device name for the address",
						},
						"translation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP translation address",
						},
					},
				},
			},
			"virtual_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Virtual servers configured on the GTM server",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the virtual server",
						},
						"destination": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination IP address and port",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the virtual server is enabled",
						},
						"translation_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Translation IP address for NAT",
						},
						"translation_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Translation port for NAT",
						},
						"monitor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health monitor for this virtual server",
						},
					},
				},
			},
		},
	}
}

func dataSourceBigipGtmServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	log.Printf("[DEBUG] Reading GTM Server data source: %s", name)

	server, err := client.GetGtmserver(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving GTM Server %s: %v", name, err))
	}
	if server == nil {
		return diag.FromErr(fmt.Errorf("GTM Server %s not found", name))
	}

	d.SetId(name)
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

	// Handle virtual servers
	if len(server.GTMVirtual_Server) > 0 {
		virtualServers := make([]interface{}, len(server.GTMVirtual_Server))
		for i, vs := range server.GTMVirtual_Server {
			translationAddr := vs.TranslationAddress
			if translationAddr == "none" {
				translationAddr = ""
			}
			virtualServers[i] = map[string]interface{}{
				"name":                vs.Name,
				"destination":         vs.Destination,
				"enabled":             vs.Enabled,
				"translation_address": translationAddr,
				"translation_port":    vs.TranslationPort,
				"monitor":             vs.Monitor,
			}
		}
		d.Set("virtual_servers", virtualServers)
	}

	return nil
}

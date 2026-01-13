package bigip

import (
	"context"
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipGtmWideip() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmWideipCreate,
		ReadContext:   resourceBigipGtmWideipRead,
		UpdateContext: resourceBigipGtmWideipUpdate,
		DeleteContext: resourceBigipGtmWideipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBigipGtmWideipImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the WideIP. Example: testwideip.local",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the type of WideIP (a, aaaa, cname, mx, naptr, srv)",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				Description: "Partition in which the WideIP resides",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User-defined description of the WideIP",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable or disable the WideIP",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disabled state of the WideIP",
			},
			"minimal_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "Specifies whether to minimize the response to the DNS query (enabled or disabled)",
			},
			"failure_rcode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "noerror",
				Description: "Specifies the DNS RCODE used when failure_rcode_response is enabled (noerror, formerr, servfail, nxdomain, notimp, refused, yxdomain, yxrrset, nxrrset, notauth, notzone)",
			},
			"failure_rcode_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Specifies whether to return a RCODE response to DNS queries when the WideIP is unavailable (enabled or disabled)",
			},
			"failure_rcode_ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the negative caching TTL of the SOA for the RCODE response",
			},
			"last_resort_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the last resort pool for the WideIP. Format: <type> <partition>/<pool_name> (e.g., 'a /Common/firstpool')",
			},
			"load_balancing_decision_log_verbosity": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Specifies the amount of detail logged when making load balancing decisions. Example: ['pool-selection']",
			},
			"persist_cidr_ipv4": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     32,
				Description: "Specifies the CIDR for IPv4 persistence",
			},
			"persist_cidr_ipv6": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     128,
				Description: "Specifies the CIDR for IPv6 persistence",
			},
			"persistence": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Specifies persistence for the WideIP (disabled or enabled)",
			},
			"pool_lb_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "round-robin",
				Description: "Specifies the load balancing mode for pools in the WideIP (round-robin, ratio, topology, global-availability)",
			},
			"ttl_persistence": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Specifies the TTL for the persistence of the WideIP",
			},
			"topology_prefer_edns0_client_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Specifies whether to prefer EDNS0 client subnet data for topology-based load balancing (enabled or disabled)",
			},
			"aliases": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Specifies alternate domain names for the WideIP",
			},
		},
	}
}

func resourceBigipGtmWideipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	partition := d.Get("partition").(string)

	log.Printf("[INFO] Creating GTM WideIP: %s (type: %s)", name, recordType)

	// Construct the full path for the WideIP
	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	config := &bigip.GTMWideIP{
		Name:      name,
		Partition: partition,
		FullPath:  fullPath,
	}

	// First, create the WideIP with just the name
	err := client.CreateGTMWideIP(config, recordType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM WideIP (%s): %s", name, err))
	}

	// Set the ID using the full path and type
	d.SetId(fmt.Sprintf("%s:%s", recordType, fullPath))

	// Now update with all the additional properties
	return resourceBigipGtmWideipUpdate(ctx, d, meta)
}

func resourceBigipGtmWideipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	// Get values from state
	id := d.Id()
	var name, recordType, partition string

	// Try to parse ID format: "type:/partition/name"
	if idParts := parseWideIPID(id); idParts != nil {
		recordType = idParts["type"]
		partition = idParts["partition"]
		name = idParts["name"]
	} else {
		// Fallback to getting from state
		name = d.Get("name").(string)
		recordType = d.Get("type").(string)
		partition = d.Get("partition").(string)
		if partition == "" {
			partition = "Common"
		}
	}

	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[INFO] Reading GTM WideIP: %s (type: %s)", fullPath, recordType)

	wideip, err := client.GetGTMWideIP(fullPath, recordType)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve GTM WideIP %s: %v", fullPath, err)
		return diag.FromErr(err)
	}

	if wideip == nil {
		log.Printf("[WARN] GTM WideIP (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	// Set all the attributes
	if wideip.Name != "" {
		d.Set("name", wideip.Name)
	}
	if wideip.Partition != "" {
		d.Set("partition", wideip.Partition)
	}
	d.Set("type", recordType)
	d.Set("description", wideip.Description)
	d.Set("enabled", wideip.Enabled)
	d.Set("disabled", wideip.Disabled)
	d.Set("failure_rcode", wideip.FailureRcode)
	d.Set("failure_rcode_response", wideip.FailureRcodeResponse)
	d.Set("failure_rcode_ttl", wideip.FailureRcodeTTL)
	d.Set("last_resort_pool", wideip.LastResortPool)
	d.Set("minimal_response", wideip.MinimalResponse)
	d.Set("persist_cidr_ipv4", wideip.PersistCidrIpv4)
	d.Set("persist_cidr_ipv6", wideip.PersistCidrIpv6)
	d.Set("persistence", wideip.Persistence)
	d.Set("pool_lb_mode", wideip.PoolLbMode)
	d.Set("ttl_persistence", wideip.TTLPersistence)
	d.Set("topology_prefer_edns0_client_subnet", wideip.TopologyPreferEdns0ClientSubnet)

	// Handle LoadBalancingDecisionLogVerbosity as array
	if len(wideip.LoadBalancingDecisionLogVerbosity) > 0 {
		d.Set("load_balancing_decision_log_verbosity", wideip.LoadBalancingDecisionLogVerbosity)
	}

	// Handle Aliases as array
	if len(wideip.Aliases) > 0 {
		d.Set("aliases", wideip.Aliases)
	}

	return nil
}

func resourceBigipGtmWideipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	partition := d.Get("partition").(string)

	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[INFO] Updating GTM WideIP: %s (type: %s)", fullPath, recordType)

	wideip := &bigip.GTMWideIP{
		Name:                            name,
		Partition:                       partition,
		FullPath:                        fullPath,
		Description:                     d.Get("description").(string),
		Enabled:                         d.Get("enabled").(bool),
		Disabled:                        d.Get("disabled").(bool),
		FailureRcode:                    d.Get("failure_rcode").(string),
		FailureRcodeResponse:            d.Get("failure_rcode_response").(string),
		FailureRcodeTTL:                 d.Get("failure_rcode_ttl").(int),
		LastResortPool:                  d.Get("last_resort_pool").(string),
		MinimalResponse:                 d.Get("minimal_response").(string),
		PersistCidrIpv4:                 d.Get("persist_cidr_ipv4").(int),
		PersistCidrIpv6:                 d.Get("persist_cidr_ipv6").(int),
		Persistence:                     d.Get("persistence").(string),
		PoolLbMode:                      d.Get("pool_lb_mode").(string),
		TopologyPreferEdns0ClientSubnet: d.Get("topology_prefer_edns0_client_subnet").(string),
		TTLPersistence:                  d.Get("ttl_persistence").(int),
	}

	// Handle LoadBalancingDecisionLogVerbosity as array
	if v, ok := d.GetOk("load_balancing_decision_log_verbosity"); ok {
		verbositySet := v.(*schema.Set)
		verbosityList := make([]string, 0, verbositySet.Len())
		for _, item := range verbositySet.List() {
			verbosityList = append(verbosityList, item.(string))
		}
		wideip.LoadBalancingDecisionLogVerbosity = verbosityList
	}

	// Handle Aliases as array
	if v, ok := d.GetOk("aliases"); ok {
		aliasesSet := v.(*schema.Set)
		aliasesList := make([]string, 0, aliasesSet.Len())
		for _, item := range aliasesSet.List() {
			aliasesList = append(aliasesList, item.(string))
		}
		wideip.Aliases = aliasesList
	}

	err := client.ModifyGTMWideIP(fullPath, wideip, recordType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM WideIP (%s): %s", fullPath, err))
	}

	return resourceBigipGtmWideipRead(ctx, d, meta)
}
func resourceBigipGtmWideipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	partition := d.Get("partition").(string)

	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[INFO] Deleting GTM WideIP: %s (type: %s)", fullPath, recordType)

	err := client.DeleteGTMWideIP(fullPath, recordType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting GTM WideIP (%s): %s", fullPath, err))
	}

	d.SetId("")
	return nil
}

func resourceBigipGtmWideipImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Import ID format: "type:/partition/name" or "type:name"
	id := d.Id()

	idParts := parseWideIPID(id)
	if idParts == nil {
		return nil, fmt.Errorf("invalid import ID format. Use 'type:/partition/name' (e.g., 'a:/Common/testwideip.local')")
	}

	d.Set("type", idParts["type"])
	d.Set("partition", idParts["partition"])
	d.Set("name", idParts["name"])

	// Set the ID in the correct format
	fullPath := fmt.Sprintf("/%s/%s", idParts["partition"], idParts["name"])
	d.SetId(fmt.Sprintf("%s:%s", idParts["type"], fullPath))

	// Read the resource to populate all attributes
	diags := resourceBigipGtmWideipRead(ctx, d, meta)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading GTM WideIP during import: %v", diags)
	}

	return []*schema.ResourceData{d}, nil
}

func parseWideIPID(id string) map[string]string {
	// Expected format: "type:/partition/name"
	// Example: "a:/Common/testwideip.local"

	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil
	}

	recordType := parts[0]
	path := parts[1]

	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	// Split path into partition and name
	pathParts := strings.SplitN(path, "/", 2)
	if len(pathParts) != 2 {
		// Try with Common as default partition
		return map[string]string{
			"type":      recordType,
			"partition": "Common",
			"name":      path,
		}
	}

	return map[string]string{
		"type":      recordType,
		"partition": pathParts[0],
		"name":      pathParts[1],
	}
}

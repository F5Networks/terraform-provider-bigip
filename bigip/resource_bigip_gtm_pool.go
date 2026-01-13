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

func resourceBigipGtmPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmPoolCreate,
		ReadContext:   resourceBigipGtmPoolRead,
		UpdateContext: resourceBigipGtmPoolUpdate,
		DeleteContext: resourceBigipGtmPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBigipGtmPoolImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the GTM pool",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of GTM pool (a, aaaa, cname, mx, naptr, srv)",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				ForceNew:    true,
				Description: "Partition in which the pool resides",
			},
			"alternate_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "round-robin",
				Description: "Specifies the load balancing mode to use if the preferred and alternate modes are unsuccessful",
			},
			"dynamic_ratio": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Enables or disables the dynamic ratio load balancing algorithm",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable or disable the pool",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disabled state of the pool",
			},
			"fallback_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any",
				Description: "Specifies the IPv4 or IPv6 address of the server to which the system directs requests when it cannot use one of its pools",
			},
			"fallback_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "return-to-dns",
				Description: "Specifies the load balancing mode that the system uses if the pool's preferred and alternate modes are unsuccessful",
			},
			"load_balancing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "round-robin",
				Description: "Specifies the preferred load balancing mode for the pool",
			},
			"manual_resume": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Specifies whether manual resume is enabled",
			},
			"max_answers_returned": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Specifies the maximum number of available virtual servers that the system lists in a response",
			},
			"monitor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the health monitor for the pool",
			},
			"qos_hit_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Specifies the weight for QoS hit ratio",
			},
			"qos_hops": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the weight for QoS hops",
			},
			"qos_kilobytes_second": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "Specifies the weight for QoS kilobytes per second",
			},
			"qos_lcs": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Specifies the weight for QoS link capacity",
			},
			"qos_packet_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Specifies the weight for QoS packet rate",
			},
			"qos_rtt": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     50,
				Description: "Specifies the weight for QoS round trip time",
			},
			"qos_topology": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the weight for QoS topology",
			},
			"qos_vs_capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the weight for QoS virtual server capacity",
			},
			"qos_vs_score": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the weight for QoS virtual server score",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Specifies the time to live (TTL) for the pool",
			},
			"limit_max_bps": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the maximum allowable data throughput rate in bits per second",
			},
			"limit_max_bps_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Enables or disables the limit_max_bps option",
			},
			"limit_max_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the maximum number of concurrent connections",
			},
			"limit_max_connections_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Enables or disables the limit_max_connections option",
			},
			"limit_max_pps": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the maximum allowable data transfer rate in packets per second",
			},
			"limit_max_pps_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Enables or disables the limit_max_pps option",
			},
			"min_members_up_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "off",
				Description: "Specifies whether the minimum number of members must be up for the pool to be active",
			},
			"min_members_up_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the minimum number of pool members that must be up",
			},
			"verify_member_availability": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "Specifies whether the system verifies the availability of pool members before sending traffic",
			},
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the pool member (format: <server_name>:<virtual_server_name>)",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Enable or disable the pool member",
						},
						"disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Disabled state of the pool member",
						},
						"ratio": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "Specifies the weight of the pool member for load balancing",
						},
						"member_order": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the order in which the member will be used",
						},
						"monitor": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "default",
							Description: "Specifies the health monitor for this pool member",
						},
						"limit_max_bps": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the maximum allowable data throughput rate for this member",
						},
						"limit_max_bps_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "disabled",
							Description: "Enables or disables the limit_max_bps option for this member",
						},
						"limit_max_connections": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the maximum number of concurrent connections for this member",
						},
						"limit_max_connections_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "disabled",
							Description: "Enables or disables the limit_max_connections option for this member",
						},
						"limit_max_pps": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the maximum allowable data transfer rate in packets per second for this member",
						},
						"limit_max_pps_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "disabled",
							Description: "Enables or disables the limit_max_pps option for this member",
						},
					},
				},
			},
		},
	}
}

func resourceBigipGtmPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	poolType := d.Get("type").(string)
	partition := d.Get("partition").(string)

	log.Printf("[INFO] Creating GTM Pool: %s (type: %s)", name, poolType)

	// Construct the full path for the Pool
	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	config := &bigip.GtmPool{
		Name:      name,
		Partition: partition,
		FullPath:  fullPath,
	}

	// First, create the Pool with just the name
	err := client.AddGTMPool(config, poolType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM Pool (%s): %s", name, err))
	}

	// Set the ID using the full path and type
	d.SetId(fmt.Sprintf("%s:%s", poolType, fullPath))

	// Now update with all the additional properties
	return resourceBigipGtmPoolUpdate(ctx, d, meta)
}

func resourceBigipGtmPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	// Get values from state
	name := d.Get("name").(string)
	poolType := d.Get("type").(string)
	partition := d.Get("partition").(string)

	// If any are empty, try to parse from ID during import
	if name == "" || poolType == "" || partition == "" {
		id := d.Id()
		if idParts := parseGtmPoolID(id); idParts != nil {
			if poolType == "" {
				poolType = idParts["type"]
			}
			if partition == "" {
				partition = idParts["partition"]
			}
			if name == "" {
				name = idParts["name"]
			}
		}
	}

	if partition == "" {
		partition = "Common"
	}

	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[INFO] Reading GTM Pool: %s (type: %s)", fullPath, poolType)

	pool, err := client.GetGTMPool(fullPath, poolType)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve GTM Pool %s: %v", fullPath, err)
		return diag.FromErr(err)
	}

	if pool == nil {
		log.Printf("[WARN] GTM Pool (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] GTM Pool %s has %d members", fullPath, len(pool.Members))
	if len(pool.Members) > 0 {
		log.Printf("[DEBUG] First member: %+v", pool.Members[0])
	}

	// Set all the attributes
	if pool.Name != "" {
		d.Set("name", pool.Name)
	}
	if pool.Partition != "" {
		d.Set("partition", pool.Partition)
	}
	d.Set("type", poolType)

	d.Set("alternate_mode", pool.AlternateMode)
	d.Set("dynamic_ratio", pool.DynamicRatio)
	d.Set("enabled", pool.Enabled)
	d.Set("disabled", pool.Disabled)
	d.Set("fallback_ip", pool.FallbackIp)
	d.Set("fallback_mode", pool.FallbackMode)
	d.Set("load_balancing_mode", pool.LoadBalancingMode)
	d.Set("manual_resume", pool.ManualResume)
	d.Set("max_answers_returned", pool.MaxAnswersReturned)
	d.Set("monitor", pool.Monitor)
	d.Set("qos_hit_ratio", pool.QosHitRatio)
	d.Set("qos_hops", pool.QosHops)
	d.Set("qos_kilobytes_second", pool.QosKilobytesSecond)
	d.Set("qos_lcs", pool.QosLcs)
	d.Set("qos_packet_rate", pool.QosPacketRate)
	d.Set("qos_rtt", pool.QosRtt)
	d.Set("qos_topology", pool.QosTopology)
	d.Set("qos_vs_capacity", pool.QosVsCapacity)
	d.Set("qos_vs_score", pool.QosVsScore)
	d.Set("ttl", pool.Ttl)
	d.Set("limit_max_bps", pool.LimitMaxBps)
	d.Set("limit_max_bps_status", pool.LimitMaxBpsStatus)
	d.Set("limit_max_connections", pool.LimitMaxConnections)
	d.Set("limit_max_connections_status", pool.LimitMaxConnectionsStatus)
	d.Set("limit_max_pps", pool.LimitMaxPps)
	d.Set("limit_max_pps_status", pool.LimitMaxPpsStatus)
	d.Set("min_members_up_mode", pool.MinMembersUpMode)
	d.Set("min_members_up_value", pool.MinMembersUpValue)
	d.Set("verify_member_availability", pool.VerifyMemberAvailability)

	// Handle Members
	if len(pool.Members) > 0 {
		members := make([]interface{}, 0, len(pool.Members))
		for _, member := range pool.Members {
			// Construct the member name from partition, subPath, and name
			// Format: /{partition}/{subPath}/{name}
			var memberName string
			if member.SubPath != "" {
				memberName = fmt.Sprintf("/%s/%s/%s", member.Partition, member.SubPath, member.Name)
			} else if member.Partition != "" {
				memberName = fmt.Sprintf("/%s/%s", member.Partition, member.Name)
			} else {
				memberName = member.Name
			}

			memberMap := map[string]interface{}{
				"name":                         memberName,
				"enabled":                      member.Enabled,
				"disabled":                     member.Disabled,
				"ratio":                        member.Ratio,
				"member_order":                 member.MemberOrder,
				"monitor":                      member.Monitor,
				"limit_max_bps":                member.LimitMaxBps,
				"limit_max_bps_status":         member.LimitMaxBpsStatus,
				"limit_max_connections":        member.LimitMaxConnections,
				"limit_max_connections_status": member.LimitMaxConnectionsStatus,
				"limit_max_pps":                member.LimitMaxPps,
				"limit_max_pps_status":         member.LimitMaxPpsStatus,
			}
			members = append(members, memberMap)
		}
		d.Set("members", members)
	}

	return nil
}

func resourceBigipGtmPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	poolType := d.Get("type").(string)
	partition := d.Get("partition").(string)

	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[INFO] Updating GTM Pool: %s (type: %s)", fullPath, poolType)

	pool := &bigip.GtmPool{
		Name:                      name,
		Partition:                 partition,
		FullPath:                  fullPath,
		AlternateMode:             d.Get("alternate_mode").(string),
		DynamicRatio:              d.Get("dynamic_ratio").(string),
		Enabled:                   d.Get("enabled").(bool),
		Disabled:                  d.Get("disabled").(bool),
		FallbackIp:                d.Get("fallback_ip").(string),
		FallbackMode:              d.Get("fallback_mode").(string),
		LoadBalancingMode:         d.Get("load_balancing_mode").(string),
		ManualResume:              d.Get("manual_resume").(string),
		MaxAnswersReturned:        d.Get("max_answers_returned").(int),
		Monitor:                   d.Get("monitor").(string),
		QosHitRatio:               d.Get("qos_hit_ratio").(int),
		QosHops:                   d.Get("qos_hops").(int),
		QosKilobytesSecond:        d.Get("qos_kilobytes_second").(int),
		QosLcs:                    d.Get("qos_lcs").(int),
		QosPacketRate:             d.Get("qos_packet_rate").(int),
		QosRtt:                    d.Get("qos_rtt").(int),
		QosTopology:               d.Get("qos_topology").(int),
		QosVsCapacity:             d.Get("qos_vs_capacity").(int),
		QosVsScore:                d.Get("qos_vs_score").(int),
		Ttl:                       d.Get("ttl").(int),
		LimitMaxBps:               d.Get("limit_max_bps").(int),
		LimitMaxBpsStatus:         d.Get("limit_max_bps_status").(string),
		LimitMaxConnections:       d.Get("limit_max_connections").(int),
		LimitMaxConnectionsStatus: d.Get("limit_max_connections_status").(string),
		LimitMaxPps:               d.Get("limit_max_pps").(int),
		LimitMaxPpsStatus:         d.Get("limit_max_pps_status").(string),
		MinMembersUpMode:          d.Get("min_members_up_mode").(string),
		MinMembersUpValue:         d.Get("min_members_up_value").(int),
		VerifyMemberAvailability:  d.Get("verify_member_availability").(string),
	}

	// Handle Members
	if v, ok := d.GetOk("members"); ok {
		membersSet := v.(*schema.Set)
		membersList := make([]bigip.GtmPoolMembers, 0, membersSet.Len())
		for _, item := range membersSet.List() {
			memberMap := item.(map[string]interface{})
			member := bigip.GtmPoolMembers{
				Name:                      memberMap["name"].(string),
				Enabled:                   memberMap["enabled"].(bool),
				Disabled:                  memberMap["disabled"].(bool),
				Ratio:                     memberMap["ratio"].(int),
				MemberOrder:               memberMap["member_order"].(int),
				Monitor:                   memberMap["monitor"].(string),
				LimitMaxBps:               memberMap["limit_max_bps"].(int),
				LimitMaxBpsStatus:         memberMap["limit_max_bps_status"].(string),
				LimitMaxConnections:       memberMap["limit_max_connections"].(int),
				LimitMaxConnectionsStatus: memberMap["limit_max_connections_status"].(string),
				LimitMaxPps:               memberMap["limit_max_pps"].(int),
				LimitMaxPpsStatus:         memberMap["limit_max_pps_status"].(string),
			}
			membersList = append(membersList, member)
		}
		pool.Members = membersList
	}

	err := client.ModifyGTMPool(fullPath, pool, poolType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM Pool (%s): %s", fullPath, err))
	}

	return resourceBigipGtmPoolRead(ctx, d, meta)
}

func resourceBigipGtmPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	poolType := d.Get("type").(string)
	partition := d.Get("partition").(string)

	fullPath := fmt.Sprintf("/%s/%s", partition, name)

	log.Printf("[INFO] Deleting GTM Pool: %s (type: %s)", fullPath, poolType)

	err := client.DeleteGTMPool(fullPath, poolType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting GTM Pool (%s): %s", fullPath, err))
	}

	d.SetId("")
	return nil
}

func resourceBigipGtmPoolImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Import ID format: "type:/partition/name" or "type:name"
	id := d.Id()

	idParts := parseGtmPoolID(id)
	if idParts == nil {
		return nil, fmt.Errorf("invalid import ID format. Use 'type:/partition/name' (e.g., 'a:/Common/firstpool')")
	}

	d.Set("type", idParts["type"])
	d.Set("partition", idParts["partition"])
	d.Set("name", idParts["name"])

	// Set the ID in the correct format
	fullPath := fmt.Sprintf("/%s/%s", idParts["partition"], idParts["name"])
	d.SetId(fmt.Sprintf("%s:%s", fullPath, idParts["type"]))

	// Read the resource to populate all attributes
	diags := resourceBigipGtmPoolRead(ctx, d, meta)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading GTM Pool during import: %v", diags)
	}

	return []*schema.ResourceData{d}, nil
}

func parseGtmPoolID(id string) map[string]string {
	// Expected format: "type:/partition/name"
	// Example: "a:/Common/firstpool"

	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil
	}

	poolType := parts[0]
	path := parts[1]

	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	// Split path into partition and name
	pathParts := strings.SplitN(path, "/", 2)
	if len(pathParts) != 2 {
		// Try with Common as default partition
		return map[string]string{
			"type":      poolType,
			"partition": "Common",
			"name":      path,
		}
	}

	return map[string]string{
		"type":      poolType,
		"partition": pathParts[0],
		"name":      pathParts[1],
	}
}

/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipNetSelfIP() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceBigipNetSelfIPCreate,
		ReadContext:   resourceBigipNetSelfIPRead,
		UpdateContext: resourceBigipNetSelfIPUpdate,
		DeleteContext: resourceBigipNetSelfIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the SelfIP",
			},

			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SelfIP IP address",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					old = strings.Replace(old, "%0", "", 1)
					new = strings.Replace(new, "%0", "", 1)
					return old == new
				},
			},

			"vlan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the vlan",
			},

			"traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group, defaults to traffic-group-local-only if not specified",
				Default:     "traffic-group-local-only",
			},

			"port_lockdown": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				// Computed: Read always populates this (defaulting to
				// ["none"] when BIG-IP's own allowService is absent/nil --
				// see resourceBigipNetSelfIPRead), even when the user
				// never configures port_lockdown at all. Without
				// Computed, Terraform treated that as drift on every
				// subsequent plan (state having a value an omitted-from-
				// config attribute "shouldn't"), permanently failing
				// terraform plan's empty-plan check after every apply.
				Computed:    true,
				Description: "port lockdown",
			},
		},
	}
}

func resourceBigipNetSelfIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	pss := &bigip.SelfIP{
		Name: name,
	}
	config := getNetSelfIPConfig(d, pss)

	log.Printf("[INFO] Creating SelfIP %s", name)

	err := client.CreateSelfIP(config)

	if err != nil {
		return diag.FromErr(fmt.Errorf("Error creating SelfIP %s: %v ", name, err))
	}

	d.SetId(name)

	return resourceBigipNetSelfIPRead(ctx, d, meta)
}

func resourceBigipNetSelfIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Reading SelfIP %s", name)

	selfIP, err := client.SelfIP(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving SelfIP %s: %v ", name, err))
	}
	if selfIP == nil {
		log.Printf("[DEBUG] SelfIP %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", selfIP.FullPath)
	_ = d.Set("vlan", selfIP.Vlan)
	_ = d.Set("ip", selfIP.Address)

	// Extract Traffic Group name from the full path (ignoring /Common/ prefix)
	regex := regexp.MustCompile(`\/Common\/(.+)`)
	_ = d.Set("traffic_group", selfIP.TrafficGroup)
	trafficGroup := regex.FindStringSubmatch(selfIP.TrafficGroup)
	if len(trafficGroup) > 0 {
		_ = d.Set("traffic_group", trafficGroup[1])
	}
	if selfIP.AllowService == nil {
		_ = d.Set("port_lockdown", []string{"none"})
	} else {
		// AllowService can come back from the API as a plain string (e.g.
		// "all", or occasionally an empty string), a []interface{} of
		// "service:port" entries, or nil. d.Set on this TypeList/TypeSet
		// field panics if given a bare (non-slice) value, so normalize all
		// shapes to a []string before setting.
		var portLockdown []string
		switch v := selfIP.AllowService.(type) {
		case string:
			if v != "" {
				portLockdown = []string{v}
			} else {
				portLockdown = []string{"none"}
			}
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					portLockdown = append(portLockdown, s)
				}
			}
		}
		_ = d.Set("port_lockdown", reorderPortLockdownToMatchConfig(d, portLockdown))
	}
	return nil
}

// reorderPortLockdownToMatchConfig reorders apiValues (the port_lockdown
// values freshly built from a SelfIP response, in whatever order BIG-IP
// returned them) to match the order "port_lockdown" currently appears in
// d's config/state. port_lockdown is a TypeList (order-sensitive) so
// Terraform diffs it positionally, but BIG-IP does not preserve submission
// order for a multi-entry allowService (confirmed directly against a live
// device: ["default","tcp:4040"] and ["tcp:4040","default"] are both
// stored and returned as ["tcp:4040","default"]) -- storing BIG-IP's own
// order directly here would cause permanent plan drift for any config
// listing more than one entry. Entries present in the API response but
// not found in the config's current order (e.g. on import) are appended
// at the end in their original (BIG-IP-returned) order.
func reorderPortLockdownToMatchConfig(d *schema.ResourceData, apiValues []string) []string {
	configRaw, ok := d.GetOk("port_lockdown")
	if !ok {
		return apiValues
	}
	configList, ok := configRaw.([]interface{})
	if !ok || len(configList) == 0 {
		return apiValues
	}

	remaining := append([]string{}, apiValues...)
	ordered := make([]string, 0, len(apiValues))
	for _, cfgEntry := range configList {
		cfgVal, ok := cfgEntry.(string)
		if !ok {
			continue
		}
		for i, v := range remaining {
			if v == cfgVal {
				ordered = append(ordered, v)
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}
	// Append any API-returned entries not present in the config's order
	// (e.g. on import), preserving their original relative order.
	ordered = append(ordered, remaining...)

	return ordered
}

func resourceBigipNetSelfIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Updating SelfIP %s", name)

	pss := &bigip.SelfIP{
		Name: name,
	}
	config := getNetSelfIPConfig(d, pss)

	err := client.ModifySelfIP(name, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error modifying SelfIP %s: %v ", name, err))
	}

	return resourceBigipNetSelfIPRead(ctx, d, meta)

}

func resourceBigipNetSelfIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Deleting SelfIP %s", name)

	err := client.DeleteSelfIP(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error deleting SelfIP %s: %v ", name, err))
	}

	d.SetId("")
	return nil
}

func getNetSelfIPConfig(d *schema.ResourceData, config *bigip.SelfIP) *bigip.SelfIP {
	var portLockdown interface{}
	p := d.Get("port_lockdown").([]interface{})

	if len(p) == 1 {
		// "all"/"none"/"default" are BIG-IP keywords, not real
		// "protocol:port" service entries -- when port_lockdown is a
		// single-element list containing one of them, BIG-IP's
		// allowService wire format expects the bare keyword string, not
		// a one-element array (confirmed directly against a live
		// device: allowService: "default" round-trips correctly, while
		// allowService: ["default"] alongside other entries is also
		// accepted, but a single ["default"] previously fell through to
		// the raw-list default case below and was sent verbatim,
		// producing a mismatch against what Read normalizes a bare
		// "default" response back into).
		switch p[0] {
		case "all":
			portLockdown = "all"
		case "none":
			portLockdown = nil
		case "default":
			portLockdown = "default"
		default:
			portLockdown = p
		}
	} else if len(p) > 1 {
		portLockdown = p
	}

	config.Address = d.Get("ip").(string)
	config.Vlan = d.Get("vlan").(string)
	config.TrafficGroup = d.Get("traffic_group").(string)
	config.AllowService = portLockdown

	return config
}

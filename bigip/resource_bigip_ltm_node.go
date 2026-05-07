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

func resourceBigipLtmNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmNodeCreate,
		ReadContext:   resourceBigipLtmNodeRead,
		UpdateContext: resourceBigipLtmNodeUpdate,
		DeleteContext: resourceBigipLtmNodeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the node",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the node",
				ForceNew:    true,
			},
			"rate_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the maximum number of connections per second allowed for a node or node address. The default value is 'disabled'.",
				Computed:    true,
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the maximum number of connections allowed for the node or node address.",
				Computed:    true,
			},
			"dynamic_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the dynamic ratio number for the node. Used for dynamic ratio load balancing. ",
				Computed:    true,
			},
			"ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the ratio number for the node.",
				Computed:    true,
			},
			"monitor": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/icmp",
				Description: "Specifies the name of the monitor or monitor rule that you want to associate with the node.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description of the node.",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the state of the node. Preferred values are `enabled`, `disabled`, or `forced_offline`. Legacy values `user-up` and `user-down` are accepted for backwards compatibility; in legacy mode pair them with the `session` field (`user-enabled` or `user-disabled`) to fully describe the desired state.",
			},
			"session": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Legacy field controlling whether the node accepts new sessions (`user-enabled` or `user-disabled`). Ignored when `state` is set to `enabled`/`disabled`/`forced_offline`, since those values control session implicitly.",
				Computed:    true,
			},
			"fqdn": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_family": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the node's address family. The default is 'unspecified', or IP-agnostic",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the fully qualified domain name of the node.",
						},
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the amount of time before sending the next DNS query.",
						},
						"downinterval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the number of attempts to resolve a domain name. The default is 5.",
						},
						"autopopulate": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies whether the node should scale to the IP address set returned by DNS.",
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmNodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	address := d.Get("address").(string)
	rateLimit := d.Get("rate_limit").(string)
	connectionLimit := d.Get("connection_limit").(int)
	dynamicRatio := d.Get("dynamic_ratio").(int)
	monitor := d.Get("monitor").(string)
	state := d.Get("state").(string)
	session := d.Get("session").(string)
	description := d.Get("description").(string)
	ratio := d.Get("ratio").(int)

	apiState, apiSession := translateNodeState(state, session)
	if isCanonicalNodeState(state) && session != "" && session != apiSession {
		log.Printf("[WARN] state=%q overrides explicit session=%q; canonical state values control session implicitly", state, session)
	}

	r := regexp.MustCompile("^((?:[0-9]{1,3}.){3}[0-9]{1,3})|(.*:[^%]*)$")

	log.Println("[INFO] Creating node " + name + "::" + address)

	nodeConfig := &bigip.Node{
		Name:            name,
		RateLimit:       rateLimit,
		ConnectionLimit: connectionLimit,
		DynamicRatio:    dynamicRatio,
		Monitor:         monitor,
		State:           apiState,
		Session:         apiSession,
		Description:     description,
		Ratio:           ratio,
	}

	if r.MatchString(address) {
		nodeConfig.Address = address
	} else {
		interval := d.Get("fqdn.0.interval").(string)
		addressFamily := d.Get("fqdn.0.address_family").(string)
		autopopulate := d.Get("fqdn.0.autopopulate").(string)
		downinterval := d.Get("fqdn.0.downinterval").(int)

		nodeConfig.FQDN.Name = address
		nodeConfig.FQDN.Interval = interval
		nodeConfig.FQDN.AddressFamily = addressFamily
		nodeConfig.FQDN.AutoPopulate = autopopulate
		nodeConfig.FQDN.DownInterval = downinterval
	}

	log.Printf("[DEBUG] config of Node to be add :%+v", nodeConfig)
	d.SetId(name)

	exist, _ := resourceBigipLtmNodeExists(d, meta)
	if !exist {
		if err := client.AddNode(nodeConfig); err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("error modifying node %s: %v", name, err))
		}
	}
	return resourceBigipLtmNodeRead(ctx, d, meta)
}

func resourceBigipLtmNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching node " + name)

	node, err := client.GetNode(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve node %s  %v :", name, err)
		return diag.FromErr(err)
	}
	if node == nil {
		log.Printf("[WARN] Node (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if node.FQDN.Name != "" {
		_ = d.Set("address", node.FQDN.Name)
	} else {
		// xxx.xxx.xxx.xxx(%x)
		// x:x(%x)
		regex := regexp.MustCompile(`((?:(?:[0-9]{1,3}\.){3}[0-9]{1,3})|(?:.*:[^%]*))(?:\%\d+)?`)
		address := regex.FindStringSubmatch(node.Address)
		log.Println("[INFO] Address: " + address[1])
		_ = d.Set("address", node.Address)
	}
	_ = d.Set("name", name)

	if err := d.Set("rate_limit", node.RateLimit); err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Error saving Monitor to state for Node (%s): %s", d.Id(), err))
	}

	priorState := d.Get("state").(string)
	state, session := nodeStateForRead(priorState, node.State, node.Session)
	_ = d.Set("state", state)
	_ = d.Set("session", session)
	_ = d.Set("connection_limit", node.ConnectionLimit)
	_ = d.Set("description", node.Description)
	_ = d.Set("dynamic_ratio", node.DynamicRatio)
	_ = d.Set("monitor", strings.TrimSpace(node.Monitor))
	_ = d.Set("ratio", node.Ratio)
	if _, ok := d.GetOk("fqdn"); ok {
		var fqdn []map[string]interface{}
		fqdnelements := map[string]interface{}{
			"interval":       node.FQDN.Interval,
			"downinterval":   node.FQDN.DownInterval,
			"autopopulate":   node.FQDN.AutoPopulate,
			"address_family": node.FQDN.AddressFamily,
		}
		fqdn = append(fqdn, fqdnelements)
		_ = d.Set("fqdn", fqdn)
	}
	return nil
}

func resourceBigipLtmNodeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching node " + name)

	node, err := client.GetNode(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve node %s  %v :", name, err)
		return false, err
	}

	if node == nil {
		log.Printf("[WARN] node (%s) not found, removing from state", d.Id())
		return false, nil
	}
	return true, nil
}

func resourceBigipLtmNodeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	address := d.Get("address").(string)
	r := regexp.MustCompile("^((?:[0-9]{1,3}.){3}[0-9]{1,3})|(.*:[^%]*)$")

	state := d.Get("state").(string)
	session := d.Get("session").(string)
	apiState, apiSession := translateNodeState(state, session)
	if isCanonicalNodeState(state) && session != "" && session != apiSession {
		log.Printf("[WARN] state=%q overrides explicit session=%q; canonical state values control session implicitly", state, session)
	}

	nodeConfig := &bigip.Node{
		ConnectionLimit: d.Get("connection_limit").(int),
		DynamicRatio:    d.Get("dynamic_ratio").(int),
		Monitor:         d.Get("monitor").(string),
		RateLimit:       d.Get("rate_limit").(string),
		State:           apiState,
		Session:         apiSession,
		Description:     d.Get("description").(string),
		Ratio:           d.Get("ratio").(int),
	}

	if r.MatchString(address) {
		nodeConfig.Address = address
	}

	if err := client.ModifyNode(name, nodeConfig); err != nil {
		return diag.FromErr(fmt.Errorf("error modifying node %s: %v", name, err))
	}

	return resourceBigipLtmNodeRead(ctx, d, meta)
}

// isCanonicalNodeState reports whether s is one of the preferred enabled/
// disabled/forced_offline values that map cleanly onto a (state, session)
// API tuple.
func isCanonicalNodeState(s string) bool {
	return s == "enabled" || s == "disabled" || s == "forced_offline"
}

// translateNodeState converts the user-facing state value into the
// (apiState, apiSession) tuple sent to BIG-IP. Canonical values expand into
// both fields; legacy values (user-up/user-down) and empty pass through with
// the user-supplied session.
func translateNodeState(state, session string) (string, string) {
	switch state {
	case "enabled":
		return "user-up", "user-enabled"
	case "disabled":
		return "user-up", "user-disabled"
	case "forced_offline":
		return "user-down", "user-disabled"
	default:
		return state, session
	}
}

// nodeStateForRead picks how to present the API's state/session pair back to
// Terraform state, keyed off the prior value already stored. Canonical or
// empty prior -> canonical form; legacy prior -> preserved legacy form so
// existing configs stay diff-free.
func nodeStateForRead(prior, apiState, apiSession string) (state, session string) {
	if apiSession == "monitor-enabled" || apiSession == "user-enabled" {
		session = "user-enabled"
	} else {
		session = "user-disabled"
	}
	if isCanonicalNodeState(prior) || prior == "" {
		switch {
		case apiState == "user-down":
			state = "forced_offline"
		case session == "user-disabled":
			state = "disabled"
		default:
			state = "enabled"
		}
		return
	}
	if apiState == "user-down" {
		state = "user-down"
	} else {
		state = "user-up"
	}
	return
}

func resourceBigipLtmNodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting node " + name)
	err := client.DeleteNode(name)

	if err != nil {
		log.Printf("[ERROR] Unable to Delete Node %s  %v : ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

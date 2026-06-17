package bigip

import (
	"context"
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// topologyEndpointSchema returns the schema shared by the ldns and server blocks.
func topologyEndpointSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"region", "datacenter", "pool", "subnet",
					"country", "state", "continent", "isp",
				}, false),
				Description: "The match type. Valid values: region, datacenter, pool, subnet, country, state, continent, isp",
			},
			"match_value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The match value (e.g. /Common/east-coast, 10.0.0.0/8, US, US/California, NA, /Common/my-pool)",
			},
			"match_negate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If true, the match is negated (uses 'not' prefix)",
			},
		},
	}
}

func resourceBigipGtmTopologyRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmTopologyRecordCreate,
		ReadContext:   resourceBigipGtmTopologyRecordRead,
		UpdateContext: resourceBigipGtmTopologyRecordUpdate,
		DeleteContext: resourceBigipGtmTopologyRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBigipGtmTopologyRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"ldns": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The LDNS (source) match criteria",
				Elem:        topologyEndpointSchema(),
			},
			"server": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The server (destination) match criteria",
				Elem:        topologyEndpointSchema(),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The order in which the topology record is evaluated. Lower values are evaluated first.",
			},
			"score": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The weight or preference given to this topology record. Higher scores indicate stronger preference.",
			},
		},
	}
}

// buildTopologyName constructs the BIG-IP topology name string from ldns and server blocks.
// The BIG-IP REST API requires the topology definition in the "name" field with the format:
//
//	"ldns: <type> <value> server: <type> <value>"
//
// For example: "ldns: region /Common/east-coast server: datacenter /Common/dc1"
// Negation is expressed as: "ldns: not region /Common/east-coast server: datacenter /Common/dc1"
func buildTopologyName(d *schema.ResourceData) string {
	ldnsList := d.Get("ldns").([]interface{})
	serverList := d.Get("server").([]interface{})

	ldns := ldnsList[0].(map[string]interface{})
	server := serverList[0].(map[string]interface{})

	ldnsStr := buildEndpointString(ldns)
	serverStr := buildEndpointString(server)

	return fmt.Sprintf("ldns: %s server: %s", ldnsStr, serverStr)
}

func buildEndpointString(endpoint map[string]interface{}) string {
	matchType := endpoint["match_type"].(string)
	matchValue := endpoint["match_value"].(string)
	negate := endpoint["match_negate"].(bool)

	if negate {
		return fmt.Sprintf("not %s %s", matchType, matchValue)
	}
	return fmt.Sprintf("%s %s", matchType, matchValue)
}

// parseTopologyName parses a BIG-IP topology name string back into ldns and server components.
// Input format: "ldns: [not] <type> <value> server: [not] <type> <value>"
func parseTopologyName(name string) (ldnsType, ldnsValue string, ldnsNegate bool, serverType, serverValue string, serverNegate bool, err error) {
	// Split on " server: " to separate ldns and server parts
	parts := strings.SplitN(name, " server: ", 2)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid topology name format: %s", name)
		return
	}

	ldnsPart := strings.TrimPrefix(parts[0], "ldns: ")
	serverPart := parts[1]

	ldnsType, ldnsValue, ldnsNegate, err = parseEndpointString(ldnsPart)
	if err != nil {
		return
	}

	serverType, serverValue, serverNegate, err = parseEndpointString(serverPart)
	return
}

func parseEndpointString(s string) (matchType, matchValue string, negate bool, err error) {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "not ") {
		negate = true
		s = strings.TrimPrefix(s, "not ")
	}

	// The first token is the type, everything after is the value
	var found bool
	matchType, matchValue, found = strings.Cut(s, " ")
	if !found {
		err = fmt.Errorf("invalid endpoint format: %s", s)
		return
	}
	return
}

func resourceBigipGtmTopologyRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	topologyName := buildTopologyName(d)

	log.Printf("[INFO] Creating GTM Topology Record: %s", topologyName)

	config := &bigip.GTMTopologyRecord{
		Name:        topologyName,
		Description: d.Get("description").(string),
		Order:       d.Get("order").(int),
		Score:       d.Get("score").(int),
	}

	err := client.CreateGTMTopologyRecord(config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM Topology Record: %v", err))
	}

	d.SetId(topologyName)

	return resourceBigipGtmTopologyRecordRead(ctx, d, meta)
}

func resourceBigipGtmTopologyRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	topologyName := d.Id()

	log.Printf("[INFO] Reading GTM Topology Record: %s", topologyName)

	record, err := client.GetGTMTopologyRecord(topologyName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving GTM Topology Record: %v", err))
	}
	if record == nil {
		log.Printf("[WARN] GTM Topology Record not found, removing from state")
		d.SetId("")
		return nil
	}

	// Parse the API name back into structured fields
	apiName := record.Name
	if apiName == "" {
		apiName = topologyName
	}

	ldnsType, ldnsValue, ldnsNegate, serverType, serverValue, serverNegate, parseErr := parseTopologyName(apiName)
	if parseErr != nil {
		return diag.FromErr(fmt.Errorf("error parsing GTM Topology Record name '%s': %v", apiName, parseErr))
	}

	ldns := []interface{}{
		map[string]interface{}{
			"match_type":   ldnsType,
			"match_value":  ldnsValue,
			"match_negate": ldnsNegate,
		},
	}
	server := []interface{}{
		map[string]interface{}{
			"match_type":   serverType,
			"match_value":  serverValue,
			"match_negate": serverNegate,
		},
	}

	d.Set("ldns", ldns)
	d.Set("server", server)
	d.Set("description", record.Description)
	d.Set("order", record.Order)
	d.Set("score", record.Score)

	return nil
}

func resourceBigipGtmTopologyRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	topologyName := d.Id()

	log.Printf("[INFO] Updating GTM Topology Record: %s", topologyName)

	config := &bigip.GTMTopologyRecord{
		Name:        topologyName,
		Description: d.Get("description").(string),
		Order:       d.Get("order").(int),
		Score:       d.Get("score").(int),
	}

	err := client.ModifyGTMTopologyRecord(topologyName, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM Topology Record: %v", err))
	}

	return resourceBigipGtmTopologyRecordRead(ctx, d, meta)
}

func resourceBigipGtmTopologyRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	topologyName := d.Id()

	log.Printf("[INFO] Deleting GTM Topology Record: %s", topologyName)

	err := client.DeleteGTMTopologyRecord(topologyName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting GTM Topology Record: %v", err))
	}

	d.SetId("")
	return nil
}

// resourceBigipGtmTopologyRecordImport handles terraform import by parsing the
// topology name string into structured state.
func resourceBigipGtmTopologyRecordImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	topologyName := d.Id()

	_, _, _, _, _, _, err := parseTopologyName(topologyName)
	if err != nil {
		return nil, fmt.Errorf("error parsing import ID '%s': %v. Expected format: 'ldns: <type> <value> server: <type> <value>'", topologyName, err)
	}

	return []*schema.ResourceData{d}, nil
}

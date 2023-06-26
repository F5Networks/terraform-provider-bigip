package bigip

import (
	"context"
	"encoding/json"
	"fmt"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceBigipWafEntityParameter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipWafEntityParameterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Entity Parameter.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the entity parameter.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the parameter is an explicit or a wildcard attribute.",
			},
			"value_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the valid type for the value of the attribute.",
			},
			"allow_empty_type": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether an empty value is allowed for a parameter.",
			},
			"allow_repeated_parameter_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether multiple parameter instances with the same name are allowed in one request.",
			},
			"attack_signatures_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether attack signatures and threat campaigns must be detected in a parameter's value.",
			},
			"check_max_value_length": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether a parameter has a restricted maximum length for value.",
			},
			"check_min_value_length": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether a parameter has a restricted minimum length for value.",
			},
			"data_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies data type of parameter's value.",
			},
			"enable_regular_expression": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether the parameter value includes the pattern defined in regularExpression.",
			},
			"is_base64": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether a parameter’s value contains a Base64 encoded string.",
			},
			"is_cookie": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether a parameter is located in the value of Cookie header.",
			},
			"is_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether a parameter is located in headers as one of the headers.",
			},
			"level": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"global", "url", "flow"}, false),
				Description:  "Specifies whether the parameter is associated with a URL, a flow, or neither.",
			},
			"mandatory": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether a parameter must exist in the request.",
			},
			"metachars_on_parameter_value_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether disallowed metacharacters must be detected in a parameter’s value.",
			},
			"parameter_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies location of parameter in request.",
			},
			"perform_staging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines the staging state of a parameter.",
			},
			"sensitive_parameter": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether a parameter is sensitive and must be not visible in logs nor in the user interface.",
			},
			"signature_overrides_disable": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "List of Attack Signature Ids which are disabled for this particular parameter.",
			},
			"url": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"json": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The JSON for WAF Entity Parameter.",
			},
		},
	}
}

func dataSourceBigipWafEntityParameterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	parameterName := d.Get("name").(string)
	entityParameter := &bigip.Parameter{
		Name: parameterName,
	}
	if d.Get("level").(string) == "url" {
		if _, OK := d.GetOk("url"); !OK {
			return diag.FromErr(fmt.Errorf("if level set to 'url',url object must be specificed"))
		}
	}
	getEPConfig(entityParameter, d)

	parameterJson, err := json.Marshal(entityParameter)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("json", string(parameterJson))
	d.SetId(parameterName)
	return nil
}

func getEPConfig(ep *bigip.Parameter, d *schema.ResourceData) {
	if d.Get("description") != nil {
		ep.Description = d.Get("description").(string)
	}
	if d.Get("type") != nil {
		ep.Type = d.Get("type").(string)
	}
	if d.Get("value_type") != nil {
		ep.ValueType = d.Get("value_type").(string)
	}
	if d.Get("allow_empty_value") != nil {
		ep.AllowEmptyValue = d.Get("allow_empty_value").(bool)
	}
	if d.Get("allow_repeated_parameter_name") != nil {
		ep.AllowRepeatedParameterName = d.Get("allow_repeated_parameter_name").(bool)
	}
	if d.Get("attack_signatures_check") != nil {
		ep.AttackSignaturesCheck = d.Get("attack_signatures_check").(bool)
	}
	if d.Get("check_max_value_length") != nil {
		ep.CheckMaxValueLength = d.Get("check_max_value_length").(bool)
	}
	if d.Get("check_min_value_length") != nil {
		ep.CheckMinValueLength = d.Get("check_min_value_length").(bool)
	}
	if d.Get("data_type") != nil {
		ep.DataType = d.Get("data_type").(string)
	}
	if d.Get("enable_regular_expression") != nil {
		ep.EnableRegularExpression = d.Get("enable_regular_expression").(bool)
	}
	if d.Get("is_base64") != nil {
		ep.IsBase64 = d.Get("is_base64").(bool)
	}
	if d.Get("is_cookie") != nil {
		ep.IsCookie = d.Get("is_cookie").(bool)
	}
	if d.Get("is_header") != nil {
		ep.IsHeader = d.Get("is_header").(bool)
	}
	if d.Get("level") != nil {
		ep.Level = d.Get("level").(string)
	}
	if d.Get("mandatory") != nil {
		ep.Mandatory = d.Get("mandatory").(bool)
	}
	if d.Get("metachars_on_parameter_value_check") != nil {
		ep.MetacharsOnParameterValueCheck = d.Get("metachars_on_parameter_value_check").(bool)
	}
	if d.Get("parameter_location") != nil {
		ep.ParameterLocation = d.Get("parameter_location").(string)
	}
	if d.Get("perform_staging") != nil {
		ep.PerformStaging = d.Get("perform_staging").(bool)
	}
	if d.Get("sensitive_parameter") != nil {
		ep.SensitiveParameter = d.Get("sensitive_parameter").(bool)
	}
	if d.Get("signature_overrides_disable") != nil {
		sigids := d.Get("signature_overrides_disable")
		var sigs []map[string]interface{}
		for _, s := range sigids.([]interface{}) {
			s1 := map[string]interface{}{
				"enabled":     false,
				"signatureId": s,
			}
			sigs = append(sigs, s1)
		}
		ep.SignatureOverrides = sigs
	}
	if urlVal, OK := d.GetOk("url"); OK {
		for _, v := range urlVal.([]interface{}) {
			ep1 := bigip.ParameterUrl{}
			ep1.Name = v.(map[string]interface{})["name"].(string)
			ep1.Method = v.(map[string]interface{})["method"].(string)
			ep1.Type = v.(map[string]interface{})["type"].(string)
			ep1.Protocol = v.(map[string]interface{})["protocol"].(string)
			ep.URL = ep1
		}
	}
}

package bigip

import (
	"encoding/json"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBigipWafEntityParameter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipWafEntityParameterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Entity Parameter",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the policy",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the entity parameter",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "explicit",
				Description: "",
			},
			"value_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "user-input",
				Description: "",
			},
			"allow_empty_type": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "",
			},
			"allow_repeated_parameter_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "",
			},
			"attack_signatures_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "",
			},
			"check_max_value_length": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"check_min_value_length": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"data_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "alpha-numeric",
				Description: "",
			},
			"enable_regular_expression": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"is_base64": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"is_cookie": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"is_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "global",
				Description: "",
			},
			"mandatory": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"metachars_on_parameter_value_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"parameter_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any",
				Description: "",
			},
			"perform_staging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "",
			},
			"sensitive_parameter": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"signature_overrides_disable": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "",
			},
			"parameter_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The payload of the WAF Entity Parameter",
			},
		},
	}
}

func dataSourceBigipWafEntityParameterRead(d *schema.ResourceData, meta interface{}) error {
	parameterName := d.Get("name").(string)
	entityParameter := &bigip.Parameter{
		Name: parameterName,
	}

	getEPConfig(entityParameter, d)

	json_payload, err := json.Marshal(entityParameter)
	if err != nil {
		return err
	}
	_ = d.Set("parameter_json", string(json_payload))
	d.SetId(parameterName)
	return nil
}

func getEPConfig(entityParameter *bigip.Parameter, d *schema.ResourceData) {
	if d.Get("description") != nil {
		entityParameter.Description = d.Get("description").(string)
	}
	if d.Get("type") != nil {
		entityParameter.Type = d.Get("type").(string)
	}
	if d.Get("value_type") != nil {
		entityParameter.ValueType = d.Get("value_type").(string)
	}
	if d.Get("allow_empty_value") != nil {
		entityParameter.AllowEmptyValue = d.Get("allow_empty_value").(bool)
	}
	if d.Get("allow_repeated_parameter_name") != nil {
		entityParameter.AllowRepeatedParameterName = d.Get("allow_repeated_parameter_name").(bool)
	}
	if d.Get("attack_signatures_check") != nil {
		entityParameter.AttackSignaturesCheck = d.Get("attack_signatures_check").(bool)
	}
	if d.Get("check_max_value_length") != nil {
		entityParameter.CheckMaxValueLength = d.Get("check_max_value_length").(bool)
	}
	if d.Get("check_min_value_length") != nil {
		entityParameter.CheckMinValueLength = d.Get("check_min_value_length").(bool)
	}
	if d.Get("data_type") != nil {
		entityParameter.DataType = d.Get("data_type").(string)
	}
	if d.Get("enable_regular_expression") != nil {
		entityParameter.EnableRegularExpression = d.Get("enable_regular_expression").(bool)
	}
	if d.Get("is_base64") != nil {
		entityParameter.IsBase64 = d.Get("is_base64").(bool)
	}
	if d.Get("is_cookie") != nil {
		entityParameter.IsCookie = d.Get("is_cookie").(bool)
	}
	if d.Get("is_header") != nil {
		entityParameter.IsHeader = d.Get("is_header").(bool)
	}
	if d.Get("level") != nil {
		entityParameter.Level = d.Get("level").(string)
	}
	if d.Get("mandatory") != nil {
		entityParameter.Mandatory = d.Get("mandatory").(bool)
	}
	if d.Get("metachars_on_parameter_value_check") != nil {
		entityParameter.MetacharsOnParameterValueCheck = d.Get("metachars_on_parameter_value_check").(bool)
	}
	if d.Get("parameter_location") != nil {
		entityParameter.ParameterLocation = d.Get("parameter_location").(string)
	}
	if d.Get("perform_staging") != nil {
		entityParameter.PerformStaging = d.Get("perform_staging").(bool)
	}
	if d.Get("sensitive_parameter") != nil {
		entityParameter.SensitiveParameter = d.Get("sensitive_parameter").(bool)
	}
}

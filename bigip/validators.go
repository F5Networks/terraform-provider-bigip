package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
)

//Validate the incoming set only contains values from the specified set
func validateSetValues(valid *schema.Set) schema.SchemaValidateFunc {
	return func(value interface{}, field string) (ws []string, errors []error) {
		if valid.Intersection(value.(*schema.Set)).Len() != value.(*schema.Set).Len() {
			errors = append(errors, fmt.Errorf("%q can only contain %v", field, value.(*schema.Set).List()))
		}
		return
	}
}

func validateStringValue(values []string) schema.SchemaValidateFunc {
	return func(value interface{}, field string) (ws []string, errors []error) {
		for _, v := range values {
			if v == value.(string) {
				return
			}
		}
		errors = append(errors, fmt.Errorf("%q must be one of %v", field, values))
		return
	}
}

func validateF5Name(value interface{}, field string) (ws []string, errors []error) {
	match, _ := regexp.MatchString("/[\\w_\\-.]+/[\\w_\\-.]+", value.(string))
	if !match {
		errors = append(errors, fmt.Errorf("%q must match /Partition/Name and contain letters, numbers or [._-]. e.g. /Common/my-pool", field))
	}
	return
}

package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
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

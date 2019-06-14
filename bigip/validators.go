package bigip

import (
	"fmt"
	"reflect"
	"regexp"

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

func validateF5Name(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
		break
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateF5Name", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^/[\\w_\\-.]+/[\\w_\\-.:]+$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match /Partition/Name and contain letters, numbers or [._-:]. e.g. /Common/my-pool", field))
		}
	}
	return
}

func validatePoolMemberName(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
		break
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validatePoolMemberName", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^\\/[\\w_\\-.]+\\/[\\w_\\-.]+:\\d+$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match /Partition/Node_Name:Port and contain letters, numbers or [._-]. e.g. /Common/node1:80", field))
		}
	}
	return
}

func validateEnabledDisabled(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
	case []string:
		values = value.([]string)
		break
	case *[]string:
		values = *(value.(*[]string))
		break
	case string:
		values = []string{value.(string)}
		break
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateEnabledDisabled", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^enabled$|^disabled$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as enabled or disabled", field))
		}
	}
	return
}

func validateReqPrefDisabled(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
	case []string:
		values = value.([]string)
	case *[]string:
		values = *(value.(*[]string))
	case string:
		values = []string{value.(string)}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateReqPrefDisabled", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^required$|^preferred$|^disabled$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as required, preferred, or disabled", field))
		}
	}
	return
}

func validateDataGroupType(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch value.(type) {
	case *schema.Set:
		values = setToStringSlice(value.(*schema.Set))
	case []string:
		values = value.([]string)
	case *[]string:
		values = *(value.(*[]string))
	case string:
		values = []string{value.(string)}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateDataGroupType", reflect.TypeOf(value)))
	}

	for _, v := range values {
		match, _ := regexp.MatchString("^string$|^ip$|^integer$", v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as string, ip, or integer", field))
		}
	}
	return
}

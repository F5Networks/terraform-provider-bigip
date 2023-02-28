/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"net"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateF5Name(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch val := value.(type) {
	case *schema.Set:
		values = setToStringSlice(val)
	case []string:
		values = val
	case *[]string:
		values = *(val)
	case string:
		values = []string{val}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateF5Name ", reflect.TypeOf(value)))
	}
	re := regexp.MustCompile(`^/[\w_\-.]+/[\w_\-.:]+$`)
	for _, v := range values {
		match := re.MatchString(v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match /Partition/Name and contain letters, numbers or [._-:]. e.g. /Common/my-pool", field))
		}
	}
	return
}

func validateF5NameWithDirectory(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch val := value.(type) {
	case *schema.Set:
		values = setToStringSlice(val)
	case []string:
		values = val
	case *[]string:
		values = *(val)
	case string:
		values = []string{val}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateF5Name ", reflect.TypeOf(value)))
	}
	re := regexp.MustCompile(`(^/[\w_\-.]+/[\w_\-.:]+/[\w_\-.:]+$)|(^/[\w_\-.]+/[\w_\-.:]+$)`)
	for _, v := range values {
		match := re.MatchString(v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match /Partition/Name or /Partition/Directory/Name  e.g. /Common/my-node or /Common/test/my-node", field))
		}
	}
	return
}

func validateVirtualAddressName(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch val := value.(type) {
	case *schema.Set:
		values = setToStringSlice(val)
	case []string:
		values = val
	case *[]string:
		values = *(val)
	case string:
		values = []string{val}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateVirtualAddressName", reflect.TypeOf(value)))
	}
	re := regexp.MustCompile(`^/[\w_\-.]+/[\w_\-.:]+[\%\d_]*$`)
	for _, v := range values {
		match := re.MatchString(v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match /Partition/Name and contain letters, numbers or [._-:%%]. e.g. /Common/172.16.124.156%%61", field))
		}
	}
	return
}

func validatePartitionName(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch val := value.(type) {
	case *schema.Set:
		values = setToStringSlice(val)
	case []string:
		values = val
	case *[]string:
		values = *(val)
	case string:
		values = []string{val}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validatePartitionName", reflect.TypeOf(value)))
	}
	re := regexp.MustCompile(`^[^/][^\s]+$`)
	for _, v := range values {
		match := re.MatchString(v)
		if !match {
			errors = append(errors, fmt.Errorf("%q name should not start with `/`, e.g Common [or] test-partition are valid ", field))
		}
	}
	return
}

// IsValidIP tests that the argument is a valid IP address.
func IsValidIP(value string) bool {
	return net.ParseIP(value) != nil
}

func validateEnabledDisabled(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch val := value.(type) {
	case *schema.Set:
		values = setToStringSlice(val)
	case []string:
		values = val
	case *[]string:
		values = *(val)
	case string:
		values = []string{val}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateEnabledDisabled", reflect.TypeOf(value)))
	}

	re := regexp.MustCompile("^enabled$|^disabled$")
	for _, v := range values {
		match := re.MatchString(v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as enabled or disabled", field))
		}
	}
	return
}

func validateDataGroupType(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch val := value.(type) {
	case *schema.Set:
		values = setToStringSlice(val)
	case []string:
		values = val
	case *[]string:
		values = *(val)
	case string:
		values = []string{val}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validateDataGroupType", reflect.TypeOf(value)))
	}

	re := regexp.MustCompile("^string$|^ip$|^integer$")
	for _, v := range values {
		match := re.MatchString(v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as string, ip, or integer", field))
		}
	}
	return
}

func validateAssignmentType(value interface{}, field string) (ws []string, errors []error) {
	var values []string
	switch val := value.(type) {
	case *schema.Set:
		values = setToStringSlice(val)
	case []string:
		values = val
	case *[]string:
		values = *(val)
	case string:
		values = []string{val}
	default:
		errors = append(errors, fmt.Errorf("Unknown type %v in validatePoolLicenseType", reflect.TypeOf(value)))
	}
	re := regexp.MustCompile("(?mi)^MANAGED$|^UNMANAGED$|^UNREACHABLE$")
	for _, v := range values {
		match := re.MatchString(v)
		if !match {
			errors = append(errors, fmt.Errorf("%q must match as MANAGED/UNMANAGED/UNREACHABLE", field))
		}
	}
	return
}

func getDeviceUri(str string) []string {
	re := regexp.MustCompile(`^(?:(?:(https?|s?ftp):)\/\/)([^:\/\s]+)(?::(\d*))?`)
	if len(re.FindStringSubmatch(str)) > 0 {
		return re.FindStringSubmatch(str)
	}
	return []string{}
}

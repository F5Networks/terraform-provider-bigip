/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestF5NameString(t *testing.T) {
	// test string => expected error count
	data := map[string]int{
		"/Common/foo":                           0,
		"/My-Partition_name/object-name_string": 0,
		"Common/foo":                            1,
		"/Common/foo/":                          1,
		"foo":                                   1,
		"//":                                    1,
		"/":                                     1,
	}
	for d, ec := range data {
		_, errs := validateF5Name(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

func TestF5NameSet(t *testing.T) {
	// test string => expected error count
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"/Common/foo", "/Common/bar"}): 0,
		makeStringSet(&[]string{"/Common/foo", "bar"}):         1,
		makeStringSet(&[]string{"foo", "bar"}):                 2,
	}

	for d, ec := range data {
		_, errs := validateF5Name(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

func TestF5NameList(t *testing.T) {
	// test string => expected error count
	data := map[*[]string]int{
		{"/Common/foo", "/Common/bar"}: 0,
		{"/Common/foo", "bar"}:         1,
		{"foo", "bar"}:                 2,
	}

	for d, ec := range data {
		_, errs := validateF5Name(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

func TestValidateEnabledDisabledString(t *testing.T) {
	data := map[*[]string]int{
		{"enabled"}:        0,
		{"disabled"}:       0,
		{"potato"}:         1,
		{"enabledpotato"}:  1,
		{"disabledpotato"}: 1,
	}

	for d, ec := range data {
		_, errs := validateEnabledDisabled(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

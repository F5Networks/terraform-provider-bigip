/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

//
//func TestValidateStringValue(t *testing.T) {
//	_, errors := validateStringValue([]string{"a", "b", "c"})("b", "field")
//	assert.Equal(t, 0, len(errors))
//}
//
//func TestValidateInvalidStringValue(t *testing.T) {
//	_, errors := validateStringValue([]string{"a", "b", "c"})("d", "field")
//	assert.Equal(t, 1, len(errors))
//	assert.Equal(t, "\"field\" must be one of [a b c]", errors[0].Error())
//}

func TestValidateF5NameString(t *testing.T) {
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

func TestValidateF5NameSet(t *testing.T) {
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

func TestValidateF5NameList(t *testing.T) {
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

func TestValidateF5NameDirString(t *testing.T) {
	// test string => expected error count
	data := map[string]int{
		"/myapp/Common/foo":                            0,
		"/my-dir/My-Partition_name/object-name_string": 0,
		"/myapp1/Common/foo":                           0,
		"/Common/foo/":                                 1,
		"foo":                                          1,
		"//":                                           1,
		"/":                                            1,
	}
	for d, ec := range data {
		_, errs := validateF5NameWithDirectory(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

func TestValidateF5NameDirSet(t *testing.T) {
	// test string => expected error count
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"/myapp/Common/foo", "/myapp/Common/bar"}): 0,
		makeStringSet(&[]string{"/myapp/Common/foo", "bar"}):               1,
		makeStringSet(&[]string{"foo", "bar"}):                             2,
	}

	for d, ec := range data {
		_, errs := validateF5NameWithDirectory(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

func TestValidateF5NameDirList(t *testing.T) {
	// test string => expected error count
	data := map[*[]string]int{
		{"/myapp/Common/foo", "/myapp/Common/bar"}: 0,
		{"/myapp/Common/foo", "bar"}:               1,
		{"foo", "bar"}:                             2,
	}

	for d, ec := range data {
		_, errs := validateF5NameWithDirectory(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateF5PartitionNameString(t *testing.T) {
	// test string => expected error count
	data := map[string]int{
		"common":            0,
		"My-Partition_name": 0,
		"Common":            0,
		"/Common/foo/":      1,
		"/foo":              1,
		"//":                1,
		"/":                 1,
	}
	for d, ec := range data {
		_, errs := validatePartitionName(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateF5PartitionNameSet(t *testing.T) {
	// test string => expected error count
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"Common", "test-dir"}):         0,
		makeStringSet(&[]string{"Mydir", "My-Partition_name"}): 0,
		makeStringSet(&[]string{"/foo/", "//bar"}):             2,
	}

	for d, ec := range data {
		_, errs := validatePartitionName(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateF5PartitionNameList(t *testing.T) {
	// test string => expected error count
	data := map[*[]string]int{
		{"Common", "test-dir"}: 0,
		{"mydir", "mypart1"}:   0,
		{"//", "/"}:            2,
	}

	for d, ec := range data {
		_, errs := validatePartitionName(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

func TestValidateEnabledDisabledString(t *testing.T) {
	data := map[string]int{
		"enabled":        0,
		"disabled":       0,
		"potato":         1,
		"enabledpotato":  1,
		"disabledpotato": 1,
	}

	for d, ec := range data {
		_, errs := validateEnabledDisabled(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateEnabledDisabledSet(t *testing.T) {
	// test string => expected error count
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"enabled", "disabled"}):           0,
		makeStringSet(&[]string{"disabled", "My-Partition_name"}): 1,
		makeStringSet(&[]string{"disabled23", "Enabled"}):         2,
	}

	for d, ec := range data {
		_, errs := validateEnabledDisabled(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateEnabledDisabledList(t *testing.T) {
	// test string => expected error count
	data := map[*[]string]int{
		{"enabled", "disabled"}:   0,
		{"disabled", "mypart1"}:   1,
		{"disabled23", "Enabled"}: 2,
	}

	for d, ec := range data {
		_, errs := validateEnabledDisabled(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateF5VirtualAddrsNameString(t *testing.T) {
	// test string => expected error count
	data := map[string]int{
		"/Common/virtual-address":    0,
		"/Common/virtual-address%22": 0,
		"/Testdir/virtual-address":   0,
		"/Common/foo/":               1,
		"common/virtual-address":     1,
		"/virtual-address%20":        1,
		"vaddr%20":                   1,
	}
	for d, ec := range data {
		_, errs := validateVirtualAddressName(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateF5VirtualAddrsNameSet(t *testing.T) {
	// test string => expected error count
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"/Common/virtual-address", "/Common/virtual-address%22"}):  0,
		makeStringSet(&[]string{"/Testdir/virtual-address", "/Testdir/virtual-address%0"}): 0,
		makeStringSet(&[]string{"vaddr%20", "//bar"}):                                      2,
	}

	for d, ec := range data {
		_, errs := validateVirtualAddressName(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateF5VirtualAddrsNameList(t *testing.T) {
	// test string => expected error count
	data := map[*[]string]int{
		{"/Common/virtual-address", "/Common/virtual-address%22"}:  0,
		{"/Testdir/virtual-address", "/Testdir/virtual-address%0"}: 0,
		{"vaddr%20", "//bar"}: 2,
	}

	for d, ec := range data {
		_, errs := validateVirtualAddressName(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

func TestValidateAssignTypeString(t *testing.T) {
	data := map[string]int{
		"MANAGED":        0,
		"UNMANAGED":      0,
		"UNREACHABLE":    0,
		"potato":         1,
		"enabledpotato":  1,
		"disabledpotato": 1,
	}

	for d, ec := range data {
		_, errs := validateAssignmentType(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateAssignTypeSet(t *testing.T) {
	// test string => expected error count
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"MANAGED", "UNMANAGED"}):     0,
		makeStringSet(&[]string{"UNREACHABLE", "UNMANAGED"}): 0,
		makeStringSet(&[]string{"disabled", "UNREACHABLE"}):  1,
		makeStringSet(&[]string{"disabled23", "Enabled"}):    2,
	}

	for d, ec := range data {
		_, errs := validateAssignmentType(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateAssignTypeList(t *testing.T) {
	// test string => expected error count
	data := map[*[]string]int{
		{"MANAGED", "UNMANAGED"}:   0,
		{"UNREACHABLE", "mypart1"}: 1,
		{"disabled23", "Enabled"}:  2,
	}

	for d, ec := range data {
		_, errs := validateAssignmentType(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

func TestValidateDataGroupTypeString(t *testing.T) {
	data := map[string]int{
		"string":         0,
		"ip":             0,
		"integer":        0,
		"potato":         1,
		"enabledpotato":  1,
		"disabledpotato": 1,
	}

	for d, ec := range data {
		_, errs := validateDataGroupType(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateDataGroupTypeSet(t *testing.T) {
	// test string => expected error count
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"string", "ip"}):           0,
		makeStringSet(&[]string{"ip", "integer"}):          0,
		makeStringSet(&[]string{"integer", "UNREACHABLE"}): 1,
		makeStringSet(&[]string{"disabled23", "Enabled"}):  2,
	}

	for d, ec := range data {
		_, errs := validateDataGroupType(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}
func TestValidateDataGroupTypeList(t *testing.T) {
	// test string => expected error count
	data := map[*[]string]int{
		{"string", "integer"}:     0,
		{"string", "mypart1"}:     1,
		{"disabled23", "Enabled"}: 2,
	}
	for d, ec := range data {
		_, errs := validateDataGroupType(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

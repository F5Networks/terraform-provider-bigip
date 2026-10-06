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

func TestF5NameSlice(t *testing.T) {
	// []string (not pointer) branch
	data := map[*[]string]int{
		{"/Common/foo", "/Common/bar"}: 0,
		{"/Common/foo", "bar"}:         1,
		{"foo", "bar"}:                 2,
		{}:                             0,
	}

	for d, ec := range data {
		_, errs := validateF5Name(*d, "testField")
		assert.Equal(t, ec, len(errs), "%v did not throw %d errors", *d, ec)
	}
}

func TestF5NameUnknownType(t *testing.T) {
	// default branch: unsupported type must yield exactly one error
	_, errs := validateF5Name(42, "testField")
	assert.Equal(t, 1, len(errs), "int input should produce one unknown-type error")

	_, errs = validateF5Name(nil, "testField")
	assert.Equal(t, 1, len(errs), "nil input should produce one unknown-type error")

	_, errs = validateF5Name(true, "testField")
	assert.Equal(t, 1, len(errs), "bool input should produce one unknown-type error")
}

func TestValidateF5NameWithDirectory(t *testing.T) {
	data := map[string]int{
		"/Common/my-node":      0,
		"/Common/test/my-node": 0,
		"/Partition/Dir/name":  0,
		"/Common/foo":          0,
		"Common/foo":           1,
		"foo":                  1,
		"/":                    1,
		"//":                   1,
		"/Common/":             1,
		"/a/b/c/d":             1,
		"":                     1,
	}
	for d, ec := range data {
		_, errs := validateF5NameWithDirectory(d, "testField")
		assert.Equal(t, ec, len(errs), "%q did not throw %d errors", d, ec)
	}
}

func TestValidateF5NameWithDirectorySet(t *testing.T) {
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"/Common/my-node", "/Common/test/my-node"}): 0,
		makeStringSet(&[]string{"/Common/my-node", "bad"}):                  1,
		makeStringSet(&[]string{"bad1", "bad2"}):                            2,
	}
	for d, ec := range data {
		_, errs := validateF5NameWithDirectory(d, "testField")
		assert.Equal(t, ec, len(errs), "set did not throw %d errors", ec)
	}
}

func TestValidateF5NameWithDirectorySlice(t *testing.T) {
	// []string and *[]string branches
	slice := []string{"/Common/my-node", "bad"}
	_, errs := validateF5NameWithDirectory(slice, "testField")
	assert.Equal(t, 1, len(errs))

	ptr := &[]string{"/Common/a/b", "/Common/c"}
	_, errs = validateF5NameWithDirectory(ptr, "testField")
	assert.Equal(t, 0, len(errs))
}

func TestValidateF5NameWithDirectoryUnknownType(t *testing.T) {
	_, errs := validateF5NameWithDirectory(3.14, "testField")
	assert.Equal(t, 1, len(errs), "float input should produce one unknown-type error")
}

func TestValidateVirtualAddressName(t *testing.T) {
	data := map[string]int{
		"/Common/172.16.124.156":    0,
		"/Common/172.16.124.156%61": 0,
		"/Common/1.2.3.4%1_2":       0,
		"/My_Partition/name":        0,
		"Common/1.2.3.4":            1,
		"/Common/1.2.3.4/":          1,
		"foo":                       1,
		"/":                         1,
		"":                          1,
	}
	for d, ec := range data {
		_, errs := validateVirtualAddressName(d, "testField")
		assert.Equal(t, ec, len(errs), "%q did not throw %d errors", d, ec)
	}
}

func TestValidateVirtualAddressNameSet(t *testing.T) {
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"/Common/1.2.3.4", "/Common/5.6.7.8%2"}): 0,
		makeStringSet(&[]string{"/Common/1.2.3.4", "bad"}):               1,
		makeStringSet(&[]string{"bad1", "bad2"}):                         2,
	}
	for d, ec := range data {
		_, errs := validateVirtualAddressName(d, "testField")
		assert.Equal(t, ec, len(errs), "set did not throw %d errors", ec)
	}
}

func TestValidateVirtualAddressNameSlice(t *testing.T) {
	slice := []string{"/Common/1.2.3.4", "bad"}
	_, errs := validateVirtualAddressName(slice, "testField")
	assert.Equal(t, 1, len(errs))

	ptr := &[]string{"/Common/1.2.3.4"}
	_, errs = validateVirtualAddressName(ptr, "testField")
	assert.Equal(t, 0, len(errs))
}

func TestValidateVirtualAddressNameUnknownType(t *testing.T) {
	_, errs := validateVirtualAddressName(42, "testField")
	assert.Equal(t, 1, len(errs), "int input should produce one unknown-type error")
}

func TestValidatePartitionName(t *testing.T) {
	data := map[string]int{
		"Common":         0,
		"test-partition": 0,
		"My_Partition.1": 0,
		"/Common":        1,
		"/":              1,
		"has space":      1,
		"a b":            1,
		"":               1,
		"x":              1,
	}
	for d, ec := range data {
		_, errs := validatePartitionName(d, "testField")
		assert.Equal(t, ec, len(errs), "%q did not throw %d errors", d, ec)
	}
}

func TestValidatePartitionNameSet(t *testing.T) {
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"Common", "test-partition"}): 0,
		makeStringSet(&[]string{"Common", "/bad"}):           1,
		makeStringSet(&[]string{"/bad1", "/bad2"}):           2,
	}
	for d, ec := range data {
		_, errs := validatePartitionName(d, "testField")
		assert.Equal(t, ec, len(errs), "set did not throw %d errors", ec)
	}
}

func TestValidatePartitionNameSlice(t *testing.T) {
	slice := []string{"Common", "/bad"}
	_, errs := validatePartitionName(slice, "testField")
	assert.Equal(t, 1, len(errs))

	ptr := &[]string{"Common", "other"}
	_, errs = validatePartitionName(ptr, "testField")
	assert.Equal(t, 0, len(errs))
}

func TestValidatePartitionNameUnknownType(t *testing.T) {
	_, errs := validatePartitionName(42, "testField")
	assert.Equal(t, 1, len(errs), "int input should produce one unknown-type error")
}

func TestIsValidIP(t *testing.T) {
	valid := []string{
		"192.168.1.1",
		"0.0.0.0",
		"255.255.255.255",
		"::1",
		"2001:db8::1",
		"fe80::1",
	}
	for _, v := range valid {
		assert.True(t, IsValidIP(v), "%q should be a valid IP", v)
	}

	invalid := []string{
		"",
		"not-an-ip",
		"256.256.256.256",
		"192.168.1",
		"192.168.1.1.1",
		"1.2.3.4/24",
		"gggg::1",
	}
	for _, v := range invalid {
		assert.False(t, IsValidIP(v), "%q should not be a valid IP", v)
	}
}

func TestValidateEnabledDisabledSet(t *testing.T) {
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"enabled", "disabled"}): 0,
		makeStringSet(&[]string{"enabled", "potato"}):   1,
		makeStringSet(&[]string{"foo", "bar"}):          2,
	}
	for d, ec := range data {
		_, errs := validateEnabledDisabled(d, "testField")
		assert.Equal(t, ec, len(errs), "set did not throw %d errors", ec)
	}
}

func TestValidateEnabledDisabledStringInput(t *testing.T) {
	_, errs := validateEnabledDisabled("enabled", "testField")
	assert.Equal(t, 0, len(errs))

	_, errs = validateEnabledDisabled("bogus", "testField")
	assert.Equal(t, 1, len(errs))
}

func TestValidateEnabledDisabledSlice(t *testing.T) {
	slice := []string{"enabled", "bogus"}
	_, errs := validateEnabledDisabled(slice, "testField")
	assert.Equal(t, 1, len(errs))
}

func TestValidateEnabledDisabledUnknownType(t *testing.T) {
	_, errs := validateEnabledDisabled(42, "testField")
	assert.Equal(t, 1, len(errs), "int input should produce one unknown-type error")
}

func TestValidateDataGroupType(t *testing.T) {
	data := map[string]int{
		"string":  0,
		"ip":      0,
		"integer": 0,
		"STRING":  1,
		"float":   1,
		"":        1,
		"strings": 1,
	}
	for d, ec := range data {
		_, errs := validateDataGroupType(d, "testField")
		assert.Equal(t, ec, len(errs), "%q did not throw %d errors", d, ec)
	}
}

func TestValidateDataGroupTypeSet(t *testing.T) {
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"string", "ip", "integer"}): 0,
		makeStringSet(&[]string{"string", "bogus"}):         1,
		makeStringSet(&[]string{"foo", "bar"}):              2,
	}
	for d, ec := range data {
		_, errs := validateDataGroupType(d, "testField")
		assert.Equal(t, ec, len(errs), "set did not throw %d errors", ec)
	}
}

func TestValidateDataGroupTypeSlice(t *testing.T) {
	slice := []string{"string", "bogus"}
	_, errs := validateDataGroupType(slice, "testField")
	assert.Equal(t, 1, len(errs))

	ptr := &[]string{"ip", "integer"}
	_, errs = validateDataGroupType(ptr, "testField")
	assert.Equal(t, 0, len(errs))
}

func TestValidateDataGroupTypeUnknownType(t *testing.T) {
	_, errs := validateDataGroupType(42, "testField")
	assert.Equal(t, 1, len(errs), "int input should produce one unknown-type error")
}

func TestValidateAssignmentType(t *testing.T) {
	data := map[string]int{
		"MANAGED":     0,
		"UNMANAGED":   0,
		"UNREACHABLE": 0,
		"managed":     0, // case-insensitive
		"unmanaged":   0,
		"unreachable": 0,
		"Managed":     0,
		"bogus":       1,
		"":            1,
		"MANAGE":      1,
	}
	for d, ec := range data {
		_, errs := validateAssignmentType(d, "testField")
		assert.Equal(t, ec, len(errs), "%q did not throw %d errors", d, ec)
	}
}

func TestValidateAssignmentTypeSet(t *testing.T) {
	data := map[*schema.Set]int{
		makeStringSet(&[]string{"MANAGED", "UNMANAGED"}): 0,
		makeStringSet(&[]string{"MANAGED", "bogus"}):     1,
		makeStringSet(&[]string{"foo", "bar"}):           2,
	}
	for d, ec := range data {
		_, errs := validateAssignmentType(d, "testField")
		assert.Equal(t, ec, len(errs), "set did not throw %d errors", ec)
	}
}

func TestValidateAssignmentTypeSlice(t *testing.T) {
	slice := []string{"MANAGED", "bogus"}
	_, errs := validateAssignmentType(slice, "testField")
	assert.Equal(t, 1, len(errs))

	ptr := &[]string{"UNMANAGED", "UNREACHABLE"}
	_, errs = validateAssignmentType(ptr, "testField")
	assert.Equal(t, 0, len(errs))
}

func TestValidateAssignmentTypeUnknownType(t *testing.T) {
	_, errs := validateAssignmentType(42, "testField")
	assert.Equal(t, 1, len(errs), "int input should produce one unknown-type error")
}

func TestGetDeviceUri(t *testing.T) {
	// Matching URIs return the full submatch slice (full match + 3 groups)
	matches := getDeviceUri("https://192.168.1.1:443/some/path")
	assert.Equal(t, 4, len(matches))
	assert.Equal(t, "https", matches[1])
	assert.Equal(t, "192.168.1.1", matches[2])
	assert.Equal(t, "443", matches[3])

	// http without port
	matches = getDeviceUri("http://example.com/path")
	assert.Equal(t, 4, len(matches))
	assert.Equal(t, "http", matches[1])
	assert.Equal(t, "example.com", matches[2])
	assert.Equal(t, "", matches[3])

	// ftp / sftp schemes are supported
	matches = getDeviceUri("ftp://host")
	assert.Equal(t, 4, len(matches))
	assert.Equal(t, "ftp", matches[1])

	matches = getDeviceUri("sftp://host:21")
	assert.Equal(t, 4, len(matches))
	assert.Equal(t, "sftp", matches[1])
	assert.Equal(t, "21", matches[3])

	// Non-matching inputs return an empty slice
	for _, s := range []string{"", "not-a-uri", "192.168.1.1", "tcp://host", "://host"} {
		assert.Equal(t, []string{}, getDeviceUri(s), "%q should return an empty slice", s)
	}
}

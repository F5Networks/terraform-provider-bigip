package bigip

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidStringValue(t *testing.T) {
	_, errors := validateStringValue([]string{"a", "b", "c"})("b", "field")
	assert.Equal(t, 0, len(errors))
}

func TestInvalidStringValue(t *testing.T) {
	_, errors := validateStringValue([]string{"a", "b", "c"})("d", "field")
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "\"field\" must be one of [a b c]", errors[0].Error())
}

func TestF5NameString(t *testing.T) {
	//test string => expected error count
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
	//test string => expected error count
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
	//test string => expected error count
	data := map[*[]string]int{
		&[]string{"/Common/foo", "/Common/bar"}: 0,
		&[]string{"/Common/foo", "bar"}:         1,
		&[]string{"foo", "bar"}:                 2,
	}

	for d, ec := range data {
		_, errs := validateF5Name(d, "testField")
		assert.Equal(t, ec, len(errs), "%s did not throw %d errors", d, ec)
	}
}

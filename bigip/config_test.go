package bigip

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigClient(t *testing.T) {
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

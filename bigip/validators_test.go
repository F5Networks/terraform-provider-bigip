package bigip

import (
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

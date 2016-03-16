package bigip

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
	"github.com/stretchr/testify/assert"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"bigip": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMapEntity(t *testing.T) {
	var a bigip.PolicyRuleAction
	m := map[string]interface{}{
		"name":      "foo",
		"asm":       true,
		"timeout":   1,
		"clonePool": "pool",
	}

	mapEntity(m, &a)

	assert.Equal(t, "foo", a.Name)
	assert.Equal(t, true, a.Asm)
	assert.Equal(t, 1, a.Timeout)
	assert.Equal(t, "pool", a.ClonePool)
}

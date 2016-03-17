package bigip

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapEntity(t *testing.T) {
	var c bigip.PolicyRuleCondition
	m := map[string]interface{}{
		"name":       "foo",
		"address":    true,
		"generation": 1,
		"values":     []interface{}{"biz", "baz"},
	}

	mapEntity(m, &c)

	assert.Equal(t, "foo", c.Name)
	assert.Equal(t, true, c.Address)
	assert.Equal(t, 1, c.Generation)
	assert.Equal(t, []string{"biz", "baz"}, c.Values)
}

func TestMapFromEntity(t *testing.T) {
	p := bigip.Policy{
		Name:     "policy",
		Strategy: "/Common/first-match",
		Controls: []string{"forwarding"},
		Requires: []string{"http"},
		Rules: []bigip.PolicyRule{
			bigip.PolicyRule{
				Name: "rule",
				Actions: []bigip.PolicyRuleAction{
					bigip.PolicyRuleAction{
						HttpUri: true,
						Value:   "/something",
					},
				},
				Conditions: []bigip.PolicyRuleCondition{
					bigip.PolicyRuleCondition{
						Name:       "foo",
						Generation: 1,
						Address:    true,
					},
					bigip.PolicyRuleCondition{
						Values: []string{"biz", "baz"},
					},
				},
			},
		},
	}

	d := resourceBigipLtmPolicy().TestResourceData()
	err := policyToData(&p, d)

	assert.Nil(t, err, err)
	assert.Equal(t, "/Common/first-match", d.Get("strategy").(string))
	assert.Equal(t, []string{"forwarding"}, setToStringSlice(d.Get("controls").(*schema.Set)))
	assert.Equal(t, []string{"http"}, setToStringSlice(d.Get("requires").(*schema.Set)))
	assert.Equal(t, "rule", d.Get("rule.0.name").(string))
	assert.Equal(t, true, d.Get("rule.0.condition.0.address").(bool))
	assert.Equal(t, []interface{}{"biz", "baz"}, d.Get("rule.0.condition.1.values").([]interface{}))
	assert.Equal(t, "/something", d.Get("rule.0.action.0.value").(string))
}

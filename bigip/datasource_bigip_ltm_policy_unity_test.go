package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func dataPolicyCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_ltm_policy" "test_policy" {
	  name = "/Common/test_policy"
	}	
	`, address)
}

func TestAccBigipPolicyUnit(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/policy/Common~test_policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected GET, got %s", r.Method)
		fmt.Fprintf(w, `
		{
			"name": "test_policy",
			"partition": "partition",
			"fullPath": "/Common/test_policy",
			"strategy": "/Common/best-match",
			"controls": ["caching"]
		}
		`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/Common~test_policy/rules", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected GET, got %s", r.Method)
		fmt.Fprintf(w, `
		{
			"items": [
				{
					"name": "test_rule",
					"fullPath": "test_rule",
					"actionsReference": {},
					"conditionsReference": {}
				}
			]
		}
		`)
	})

	mux.HandleFunc("/mgmt/tm/ltm/policy/Common~test_policy/rules/test_rule/actions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected GET, got %s", r.Method)
		fmt.Fprintf(w, `
		{
			"items": [
				{
					"name": "0",
					"fullPath": "0"
				}
			]
		}
		`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/Common~test_policy/rules/test_rule/conditions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected GET, got %s", r.Method)
		fmt.Fprintf(w, `
		{
			"items": [
				{
					"name": "0",
					"fullPath": "0"
				}
			]
		}
		`)
	})
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataPolicyCfg(server.URL),
			},
		},
	})
}

func TestAccBigipPolicyUnit_nil(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/policy/Common~test_policy", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code": 404}`, 404)
	})
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataPolicyCfg(server.URL),
			},
		},
	})
}

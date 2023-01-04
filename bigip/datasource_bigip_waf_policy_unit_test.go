package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func dataWAFPolicyCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_waf_policy" "test_policy" {
	  policy_id = "test_policy_id"
	}	
	`, address)
}

func TestAccBigipWAFPolicyUnit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/asm/policies/test_policy_id", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method GET, got %s", r.Method)

		fmt.Fprintf(w, `
		{
			"items": [
				{
					"name": "test_policy",
					"partition": "Common",
					"id": "test_policy_id"
				}
			]
		}
		`)
	})

	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method POST, got %s", r.Method)

		fmt.Fprintf(w, `
		{
			"status": "COMPLETED",
			"id": "1",
			"result": {}
		}
		`)
	})

	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method GET, got %s", r.Method)

		fmt.Fprintf(w, `
		{
			"status": "COMPLETED",
			"id": "1",
			"result": {
				"file": "{\"policy\":{\"name\": \"test_policy\",\"partition\": \"Common\",\"id\": \"test_policy_id\"}}"
			}
		}
		`)
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataWAFPolicyCfg(server.URL),
			},
		},
	})
}

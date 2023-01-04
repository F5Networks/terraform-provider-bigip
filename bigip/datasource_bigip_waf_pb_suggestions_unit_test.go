package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func dataWAFPbSuggestionsCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_waf_pb_suggestions" "pb_suggestions" {
	  policy_name            = "test_policy"
	  partition              = "Common"
	  minimum_learning_score = 7
	}	
	`, address)
}

func TestAccBigipWAFPbSuggestionsUnit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/asm/policies/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method GET, got %s", r.Method)

		fmt.Fprintf(w, `
		{
			"items": [
				{
					"name": "test_policy",
					"partition": "Common",
					"id": "xyz"
				}
			]
		}
		`)
	})

	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method POST, got %s", r.Method)

		fmt.Fprintf(w, `
		{
			"status": "COMPLETED",
			"id": "1",
			"result": {}
		}
		`)
	})

	counter := 0
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method GET, got %s", r.Method)

		if counter < 2 {
			fmt.Fprintf(w, `
		{
			"status": "IN PROGRESS",
			"id": "1",
			"result": {}
		}
		`)
			counter += 1
		} else {
			fmt.Fprintf(w, `
		{
			"status": "COMPLETED",
			"id": "1",
			"result": {}
		}
		`)
		}
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataWAFPbSuggestionsCfg(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bigip_waf_pb_suggestions.pb_suggestions", "policy_name", "test_policy"),
					resource.TestCheckResourceAttr("data.bigip_waf_pb_suggestions.pb_suggestions", "partition", "Common"),
					resource.TestCheckResourceAttr("data.bigip_waf_pb_suggestions.pb_suggestions", "minimum_learning_score", "7"),
				),
			},
		},
	})
}

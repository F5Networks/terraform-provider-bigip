package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func dataWAFSignaturesCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_waf_signatures" "test_signature" {
	  signature_id = 200104004
	}	
	`, address)
}

func TestAccBigipWAFSignaturesUnit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected 'GET', got '%s'", r.Method)

		fmt.Fprintf(w, `
		{
			"name": "asm",
			"fullPath": "asm",
			"level": "nominal"
		}
		`)
	})

	mux.HandleFunc("/mgmt/tm/asm/signatures/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected 'GET', got '%s'", r.Method)

		fmt.Fprintf(w, `
		{
			"items": [
				{
					"name": "test_sign",
					"id": "test_sign_id",
					"signatureId": 123456,
					"signatureType": "request",
					"accuracy": "low",
					"risk": "low"
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
				Config: dataWAFSignaturesCfg(server.URL),
				// Check: resource.ComposeTestCheckFunc(
				// 	resource.TestCheckResourceAttr("data.bigip_waf_policy.test_policy", "policy_id", "test_policy_id"),
				// ),
			},
		},
	})
}

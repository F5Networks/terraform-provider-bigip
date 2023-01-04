package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataPoolCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_ltm_pool" "test_pool" {
	  name      = "test_pool"
	  partition = "Common"
	}	
	`, address)
}

func TestAccBigipPoolUnit(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test_pool", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "test_pool",
			"partition": "Common",
			"fullPath": "/Common/test_pool"
		}
		`)
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataPoolCfg(server.URL),
			},
		},
	})
}

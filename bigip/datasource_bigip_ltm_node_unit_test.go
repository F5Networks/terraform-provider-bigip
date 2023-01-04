package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataNodeCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_ltm_node" "terraform_node" {
	  name      = "terraform_node"
	  partition = "Common"
	}	
	`, address)
}

func TestAccBigipNodeUnit(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~terraform_node", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "terraform_node",
			"partition": "Common",
			"fullPath": "/Common/terraform_node",
			"address": "1.2.3.4"
		}
		`)
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataNodeCfg(server.URL),
			},
		},
	})
}

func TestAccBigipNodeUnit_fqdn(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~terraform_node", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "terraform_node",
			"partition": "Common",
			"fullPath": "/Common/terraform_node",
			"address": "any6",
			"fqdn": {
				"tmName": "www.fqdn.com"
			}
		}
		`)
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataNodeCfg(server.URL),
			},
		},
	})
}

func TestAccBigipNodeUnit_nil(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~terraform_node", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code": 404}`, 404)
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataNodeCfg(server.URL),
			},
		},
	})
}

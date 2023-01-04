package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataGroupCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_ltm_datagroup" "images" {
		name      = "images"
		partition = "Common"
	  }	
	`, address)
}

func TestAccBigipDataGroupUnit(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/data-group/internal/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "images",
			"partition": "Common",
			"type": "string",
			"records": [
				{
					"name": ".bmp",
					"data": ""
				},
				{
					"name": ".gif",
					"data": ""
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
				Config: dataGroupCfg(server.URL),
			},
		},
	})
}

func TestAccBigipDataGroupUnit_nil(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/data-group/internal/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code": 404}`, 404)
		// w.WriteHeader(http.StatusNotFound)
	})
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataGroupCfg(server.URL),
			},
		},
	})
}

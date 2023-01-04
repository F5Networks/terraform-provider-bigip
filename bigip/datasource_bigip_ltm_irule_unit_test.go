package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataIruleCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_ltm_irule" "terraform_irule" {
	  name      = "terraform_irule"
	  partition = "Common"
	}	
	`, address)
}

func TestAccBigipIruleUnit(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~terraform_irule", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "terraform_irule",
			"partition": "Common",
			"apiAnonymous": "test irule",
			"fullPath": "/Common/terraform_irule"
		}
		`)
	})
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataIruleCfg(server.URL),
			},
		},
	})
}

func TestAccBigipIruleUnit_nil(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~terraform_irule", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code": 404}`, 404)
	})
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataIruleCfg(server.URL),
			},
		},
	})
}

package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataSSLCertificateCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_ssl_certificate" "test" {
	  name      = "terraform_ssl_certificate"
	  partition = "Common"
	}	
	`, address)
}

func TestAccBigipSSLCertificateUnit(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~terraform_ssl_certificate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"name": "terraform_ssl_certificate",
			"partition": "Common",
			"fullPath": "/Common/terraform_ssl_certificate"
		}
		`)
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSSLCertificateCfg(server.URL),
			},
		},
	})
}

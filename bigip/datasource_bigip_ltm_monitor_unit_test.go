package bigip

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataMonitorCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_ltm_monitor" "test-monitor" {
	  name      = "test-monitor"
	  partition = "Common"
	}	
	`, address)
}

func TestAccBigipMonitorUnit(t *testing.T) {
	monitors := []string{"http", "https", "icmp", "gateway-icmp", "tcp", "tcp-half-open", "ftp", "udp", "postgresql", "mysql", "mssql", "ldap"}
	setup()
	defer teardown()

	for _, name := range monitors {
		mux.HandleFunc(fmt.Sprintf("/mgmt/tm/ltm/monitor/%s", name), func(w http.ResponseWriter, r *http.Request) {
			req_url := r.URL.String()
			if strings.HasSuffix(req_url, "http") {
				fmt.Fprintf(w, `{
					"items": [
						{
							"name": "/Common/test-monitor",
							"partition": "Common",
							"fullPath": "/Common/test-monitor",
							"defaultsFrom": "/Common/http"
						}
					]
				}`)
			} else {
				fmt.Fprintf(w, `{"items":[]}`)
			}
		})
	}
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataMonitorCfg(server.URL),
			},
		},
	})
}

func TestAccBigipMonitorUnit_nil(t *testing.T) {
	monitors := []string{"http", "https", "icmp", "gateway-icmp", "tcp", "tcp-half-open", "ftp", "udp", "postgresql", "mysql", "mssql", "ldap"}
	setup()
	defer teardown()
	for _, name := range monitors {
		mux.HandleFunc(fmt.Sprintf("/mgmt/tm/ltm/monitor/%s", name), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"items": []}`)
		})
	}
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataMonitorCfg(server.URL),
			},
		},
	})
}

/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigiqRegkeyLicenseManageUnitInvalid(t *testing.T) {
	resourceName := "regkeypool_name"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBigiqRegkeyLicenseManageInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigiqRegkeyLicenseManageUnitCreate(t *testing.T) {
	resourceName := "regkeypool_name"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items": [{"id": "83880ee5-7a0d-46ff-a8f3-d59eae2e377e","name": "test-pool","sortName": "Purchased Pool"}]}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/utility/licenses", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items": [{"id": "83880ee5-7a0d-46ff-a8f3-d59eae2e377e","name": "test-pool","sortName": "Purchased Pool"}]}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/purchased-pool/licenses", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items": [{"id": "83880ee5-7a0d-46ff-a8f3-d59eae2e377e","name": "%s","sortName": "Purchased Pool"}]}`, resourceName)
	})

	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "licensePoolName": "%s",
    "command": "assign",
    "address": "127.0.0.1",
    "id": "d717c6a1-f3bd-46cb-8410-c6fda58940b9",
    "status": "STARTED",
    "userReference": {
        "link": "https://localhost/mgmt/shared/authz/users/lm"
    },
    "identityReferences": [
        {
            "link": "https://localhost/mgmt/shared/authz/users/lm"
        }
    ],
    "ownerMachineId": "b814fef0-2e8b-460d-af43-0d100dc50352",
    "taskWorkerGeneration": 1,
    "generation": 1,
    "lastUpdateMicros": 1497509153046646,
    "kind": "cm:device:tasks:licensing:pool:member-management:devicelicensingassignmenttaskstate",
    "selfLink": "https://localhost/mgmt/cm/device/tasks/licensing/pool/member-management/d717c6a1-f3bd-46cb-8410-c6fda58940b9"
}`, resourceName)
	})
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/d717c6a1-f3bd-46cb-8410-c6fda58940b9", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"address": "127.0.0.1",
    "command": "assign",
    "currentStep": "POLL_ASSIGNMENT_STATUS",
    "endDateTime": "2017-06-14T23:46:04.532-0700",
    "generation": 7,
    "id": "d717c6a1-f3bd-46cb-8410-c6fda58940b9",
    "identityReferences": [
        {
            "link": "https://localhost/mgmt/shared/authz/users/lm"
        }
    ],
    "kind": "cm:device:tasks:licensing:pool:member-management:devicelicensingassignmenttaskstate",
    "lastUpdateMicros": 1497509164582382,
    "licenseAssignmentReference": {
        "link": "https://localhost/mgmt/cm/device/licensing/pool/purchased-pool/licenses/9a79bcf5-906e-4418-83c2-190ea22b9ec8/member-management/9847bfb5-f0a4-474b-b978-04586fc6d17d"
    },
    "licensePoolName": "%s",
    "ownerMachineId": "b814fef0-2e8b-460d-af43-0d100dc50352",
    "selfLink": "https://localhost/mgmt/cm/device/tasks/licensing/pool/member-management/d717c6a1-f3bd-46cb-8410-c6fda58940b9",
    "skuKeyword1": "",
    "skuKeyword2": "",
    "startDateTime": "2017-06-14T23:45:53.063-0700",
    "status": "FINISHED",
    "userReference": {
        "link": "https://localhost/mgmt/shared/authz/users/lm"
    },
    "username": "lm"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:sys:license:licensestats",
    "selfLink": "https://localhost/mgmt/tm/sys/license?ver=16.1.0"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigiqRegkeyLicenseManageCreate(resourceName, server.URL),
			},
			{
				Config: testAccBigiqRegkeyLicenseManageModify(resourceName, server.URL),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

//
//func TestAccBigiqRegkeyLicenseManageUnitReadError(t *testing.T) {
//	resourceName := "regkeypool_name"
//	setup()
//	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//	})
//	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//		_, _ = fmt.Fprintf(w, `{}`)
//	})
//	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		_, _ = fmt.Fprintf(w, `{"name":"%s","destinationAddress":"3.10.11.2/32","ipsecPolicyReference":{},"sourceAddress":"2.10.11.12/32"}`, resourceName)
//	})
//	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector/~Common~test-traffic-selector", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		http.Error(w, "The requested IPsec Trafficselector (/Common/test-traffic-selector) was not found", http.StatusNotFound)
//	})
//
//	defer teardown()
//	resource.Test(t, resource.TestCase{
//		IsUnitTest: true,
//		Providers:  testProviders,
//		Steps: []resource.TestStep{
//			{
//				Config:      testAccBigiqRegkeyLicenseManageCreate(resourceName, server.URL),
//				ExpectError: regexp.MustCompile("HTTP 404 :: The requested IPsec Trafficselector \\(/Common/test-traffic-selector\\) was not found"),
//			},
//		},
//	})
//}
//
//func TestAccBigiqRegkeyLicenseManageUnitCreateError(t *testing.T) {
//	resourceName := "regkeypool_name"
//	httpDefault := "/Common/http"
//	setup()
//	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//	})
//	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//		_, _ = fmt.Fprintf(w, `{}`)
//	})
//	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		_, _ = fmt.Fprintf(w, `{"name":"/Common/testhttp##","defaultsFrom":"%s", "basicAuthRealm": "none"}`, httpDefault)
//		http.Error(w, "The requested object name (/Common/testravi##) is invalid", http.StatusNotFound)
//	})
//	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector/~Common~test-traffic-selector", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		http.Error(w, "The requested HTTP Profile (/Common/test-traffic-selector) was not found", http.StatusNotFound)
//	})
//
//	defer teardown()
//	resource.Test(t, resource.TestCase{
//		IsUnitTest: true,
//		Providers:  testProviders,
//		Steps: []resource.TestStep{
//			{
//				Config:      testAccBigiqRegkeyLicenseManageCreate(resourceName, server.URL),
//				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-traffic-selector\\) was not found"),
//			},
//		},
//	})
//}

func testAccBigiqRegkeyLicenseManageInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_common_license_manage_bigiq" "test_example" {
  license_poolname = "%s"
  key = "83880ee5-7a0d-46ff-a8f3-d59eae2e377e"
  assignment_type  = "MANAGED"
  bigiq_address    = ""
  bigiq_user       = ""
  bigiq_password   = ""
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testAccBigiqRegkeyLicenseManageCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_common_license_manage_bigiq" "test_example" {
  license_poolname = "%s"
  assignment_type  = "MANAGED"
  //key = "83880ee5-7a0d-46ff-a8f3-d59eae2e377e"
  bigiq_address    = "%s"
  bigiq_user       = ""
  bigiq_password   = ""
  bigiq_token_auth = false
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url, url)
}

func testAccBigiqRegkeyLicenseManageModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_common_license_manage_bigiq" "test_example" {
  license_poolname = "%s"
  assignment_type  = "MANAGED"
  //key = "83880ee5-7a0d-46ff-a8f3-d59eae2e377e"
  bigiq_address    = "%s"
  bigiq_user       = ""
  bigiq_password   = ""
  bigiq_token_auth = false
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url, url)
}

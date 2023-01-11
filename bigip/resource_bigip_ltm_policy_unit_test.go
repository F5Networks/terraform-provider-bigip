/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipLtmPolicyUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-policy"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmPolicyInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmPolicyUnitCreate(t *testing.T) {
	resourceName := "/Common/test-policy"
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
	mux.HandleFunc("/mgmt/tm/ltm/policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","strategy":"first-match"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method:%+v", r.Method)
		if r.Method == "DELETE" {
			_, _ = fmt.Fprintf(w, `{}`)
		}
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:policy:policystate",
    "name": "test-policy",
    "fullPath": "/Common/test-policy",
    "generation": 1,
    "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy?ver=16.1.0",
    "draftCopy": "/Common/Drafts/testpolicy1",
    "draftCopyReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy/~Common~Drafts~test-policy?ver=16.1.0"
    },
    "status": "published",
    "strategy": "/Common/first-match",
    "strategyReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy-strategy/~Common~first-match?ver=16.1.0"
    },
    "references": {},
    "rulesReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy/~Common~test-policy/rules?ver=16.1.0",
        "isSubcollection": true
    }
}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy/rules", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:policy:rules:rulescollectionstate",
    "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules?ver=16.1.0",
    "items": []}`)
	})
	//	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy/rules", func(w http.ResponseWriter, r *http.Request) {
	//		_, _ = fmt.Fprintf(w, `{
	//    "kind": "tm:ltm:policy:rules:rulescollectionstate",
	//    "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules?ver=16.1.0",
	//    "items": [
	//        {
	//            "kind": "tm:ltm:policy:rules:rulesstate",
	//            "name": "testrule",
	//            "fullPath": "testrule",
	//            "generation": 1,
	//            "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule?ver=16.1.0",
	//            "ordinal": 0,
	//            "actionsReference": {
	//                "link": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/actions?ver=16.1.0",
	//                "isSubcollection": true
	//            },
	//            "conditionsReference": {
	//                "link": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/conditions?ver=16.1.0",
	//                "isSubcollection": true
	//            }
	//        }
	//    ]
	//}`)
	//	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy/rules/testrule/actions", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:policy:rules:actions:actionscollectionstate",
    "selfLink": "https://localhost/mgmt/tm/ltm/policy/testpolicy1/rules/testrule/actions?ver=16.1.0",
    "items": []
}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy/rules/testrule/conditions", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:policy:rules:conditions:conditionscollectionstate",
    "selfLink": "https://localhost/mgmt/tm/ltm/policy/testpolicy1/rules/testrule/conditions?ver=16.1.0",
    "items": []
}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~~Common~test-policy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			_, _ = fmt.Fprintf(w, `{}`)
		}
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmPolicyCreate(resourceName, server.URL),
			},
			//{
			//	Config:             testBigipLtmPolicyModify(resourceName, server.URL),
			//	ExpectNonEmptyPlan: true,
			//},
		},
	})
}

func TestAccBigipLtmPolicyUnitCreateUpdate(t *testing.T) {
	var count = 0
	resourceName := "/Common/test-policy"
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
	mux.HandleFunc("/mgmt/tm/ltm/policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","strategy":"first-match"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:policy:policystate",
    "name": "test-policy",
    "fullPath": "/Common/test-policy",
    "generation": 12367,
    "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy?ver=16.1.0",
    "controls": [
        "forwarding"
    ],
    "lastModified": "2023-01-05T17:18:10Z",
    "requires": [
        "http"
    ],
    "status": "published",
    "strategy": "/Common/first-match",
    "strategyReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy-strategy/~Common~first-match?ver=16.1.0"
    },
    "references": {},
    "rulesReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy/~Common~test-policy/rules?ver=16.1.0",
        "isSubcollection": true
    }
}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~Drafts~test-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"Drafts/test-policy","publishedCopy":"","controls":["forwarding"],"requires":["http"],"strategy":"first-match","rulesReference":{"items":[{"name":"testrule","ordinal":0,"conditionsReference":{"items":[{"name":"0","equals":true,"httpHost":true,"request":true,"values":["acme.example.com"]}]},"actionsReference":{"items":[{"name":"0","httpReply":true,"location":"tcl:https://[HTTP::host][HTTP::uri]","redirect":true,"request":true}]}}]}}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy/rules", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
	   "kind": "tm:ltm:policy:rules:rulescollectionstate",
	   "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules?ver=16.1.0",
	   "items": [
	       {
	           "kind": "tm:ltm:policy:rules:rulesstate",
	           "name": "testrule",
	           "fullPath": "testrule",
	           "generation": 1,
	           "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule?ver=16.1.0",
	           "ordinal": 0,
	           "actionsReference": {
	               "link": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/actions?ver=16.1.0",
	               "isSubcollection": true
	           },
	           "conditionsReference": {
	               "link": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/conditions?ver=16.1.0",
	               "isSubcollection": true
	           }
	       }
	   ]
	}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy/rules/testrule/actions", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:policy:rules:actions:actionscollectionstate",
    "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/actions?ver=16.1.0",
    "items": [
        {
            "kind": "tm:ltm:policy:rules:actions:actionsstate",
            "name": "0",
            "fullPath": "0",
            "generation": 12366,
            "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/actions/0?ver=16.1.0",
            "code": 0,
            "expirySecs": 0,
            "httpReply": true,
            "length": 0,
            "location": "tcl:https://[HTTP::host][HTTP::uri]",
            "offset": 0,
            "port": 0,
            "redirect": true,
            "request": true,
            "status": 0,
            "timeout": 0,
            "vlanId": 0
        }
    ]
}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~Common~test-policy/rules/testrule/conditions", func(w http.ResponseWriter, r *http.Request) {
		if count == 0 || count == 1 || count == 2 || count == 3 || count == 4 {
			_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:policy:rules:conditions:conditionscollectionstate",
    "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/conditions?ver=16.1.0",
    "items": []}`)
		} else if count == 5 || count == 6 || count == 7 || count == 8 || count == 9 {
			_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:policy:rules:conditions:conditionscollectionstate",
    "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/conditions?ver=16.1.0",
    "items": [
        {
            "kind": "tm:ltm:policy:rules:conditions:conditionsstate",
            "name": "0",
            "fullPath": "0",
            "generation": 12366,
            "selfLink": "https://localhost/mgmt/tm/ltm/policy/test-policy/rules/testrule/conditions/0?ver=16.1.0",
            "all": true,
            "caseInsensitive": true,
            "equals": true,
            "external": true,
            "httpHost": true,
            "index": 0,
            "present": true,
            "remote": true,
            "request": true,
            "values": [
                "acme.example.com"
            ]
        }
    ]
}`)
		}
		count++
	})
	mux.HandleFunc("/mgmt/tm/ltm/policy/~~Common~test-policy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			_, _ = fmt.Fprintf(w, `{}`)
		}
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmPolicyCreateRule(resourceName, server.URL),
			},
			{
				Config: testBigipLtmPolicyModifyRule(resourceName, server.URL),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testBigipLtmPolicyInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_policy" "test-policy" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmPolicyCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_policy" "test-policy" {
  name    = "%s"
  strategy = "first-match"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmPolicyModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_policy" "test-policy" {
  name    = "%s"
  strategy = "first-match"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmPolicyCreateRule(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_policy" "test-policy" {
  name     = "%s"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule  {
    name = "testrule"
    action {
      redirect   = true
      connection = false
      location   = "tcl:https://[HTTP::host][HTTP::uri]"
      http_reply = true
    }
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmPolicyModifyRule(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_policy" "test-policy" {
  name     = "%s"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule  {
    name = "testrule"
    condition {
      http_host = true
      equals = true
      values  = ["acme.example.com"]
      request = true
    }
    action {
      redirect   = true
      connection = false
      location   = "tcl:https://[HTTP::host][HTTP::uri]"
      http_reply = true
    }
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

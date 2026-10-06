/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestResourceBigipEventServiceDiscoverySchema exercises the schema
// definition of the bigip_event_service_discovery resource without needing
// any BIG-IP connection.
func TestResourceBigipEventServiceDiscoverySchema(t *testing.T) {
	r := resourceServiceDiscovery()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	// taskid is required, ForceNew.
	taskidSchema, ok := r.Schema["taskid"]
	if !ok {
		t.Fatal("Expected 'taskid' field to exist in schema")
	}
	if !taskidSchema.Required {
		t.Error("Expected 'taskid' field to be required")
	}
	if !taskidSchema.ForceNew {
		t.Error("Expected 'taskid' field to be ForceNew")
	}

	// node is an optional TypeSet.
	nodeSchema, ok := r.Schema["node"]
	if !ok {
		t.Fatal("Expected 'node' field to exist in schema")
	}
	if nodeSchema.Required {
		t.Error("Expected 'node' field to be optional")
	}
	if nodeSchema.Type.String() != "TypeSet" {
		t.Errorf("Expected 'node' field to be TypeSet, got %v", nodeSchema.Type)
	}

	// Verify the nested "node" resource schema defines id/ip/port with the
	// expected types.
	nodeElem, ok := nodeSchema.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("Expected 'node' field Elem to be a *schema.Resource")
	}

	idSchema, ok := nodeElem.Schema["id"]
	if !ok {
		t.Fatal("Expected 'node.id' field to exist in schema")
	}
	if idSchema.Type != schema.TypeString {
		t.Errorf("Expected 'node.id' field to be TypeString, got %v", idSchema.Type)
	}

	ipSchema, ok := nodeElem.Schema["ip"]
	if !ok {
		t.Fatal("Expected 'node.ip' field to exist in schema")
	}
	if ipSchema.Type != schema.TypeString {
		t.Errorf("Expected 'node.ip' field to be TypeString, got %v", ipSchema.Type)
	}

	portSchema, ok := nodeElem.Schema["port"]
	if !ok {
		t.Fatal("Expected 'node.port' field to exist in schema")
	}
	if portSchema.Type != schema.TypeInt {
		t.Errorf("Expected 'node.port' field to be TypeInt, got %v", portSchema.Type)
	}

	// Verify CRUD functions are defined.
	if r.CreateContext == nil {
		t.Error("Expected CreateContext to be defined")
	}
	if r.ReadContext == nil {
		t.Error("Expected ReadContext to be defined")
	}
	if r.UpdateContext == nil {
		t.Error("Expected UpdateContext to be defined")
	}
	if r.DeleteContext == nil {
		t.Error("Expected DeleteContext to be defined")
	}

	// Verify Importer is defined.
	if r.Importer == nil {
		t.Error("Expected Importer to be defined")
	}
}

func testBigipEventServiceDiscoveryInvalid(taskid string) string {
	return fmt.Sprintf(`
resource "bigip_event_service_discovery" "test-sd" {
  taskid     = "%s"
  invalidkey = "foo"
}
`, taskid)
}

func TestUnitBigipEventServiceDiscoveryInvalid(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipEventServiceDiscoveryInvalid(taskid),
				ExpectError: regexp.MustCompile(`An argument named "invalidkey" is not expected here`),
			},
		},
	})
}

func testBigipEventServiceDiscoveryCreate(taskid string) string {
	return fmt.Sprintf(`
resource "bigip_event_service_discovery" "test-sd" {
  taskid = "%s"
  node {
    id   = "newNode1"
    ip   = "192.168.2.3"
    port = 8080
  }
}
`, taskid)
}

func TestUnitBigipEventServiceDiscoveryCreate(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskid)
			return
		}
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{
  "result": {
    "providerOptions": {
      "nodeList": [
        {
          "id": "newNode1",
          "ip": "192.168.2.3",
          "port": 8080
        }
      ]
    }
  }
}`)
	})
	// Destroy path: DeleteContext POSTs an empty node array to the
	// trailing-slash variant of the nodes endpoint using a raw *http.Client.
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipEventServiceDiscoveryCreate(taskid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_event_service_discovery.test-sd", "taskid", taskid),
				),
			},
		},
	})
}

func TestUnitBigipEventServiceDiscoveryCreateError(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, `{"code":400,"message":"invalid node configuration"}`, http.StatusBadRequest)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipEventServiceDiscoveryCreate(taskid),
				ExpectError: regexp.MustCompile(`error modifying node`),
			},
		},
	})
}

func TestUnitBigipEventServiceDiscoveryReadError(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskid)
			return
		}
		http.Error(w, `{"code":404,"message":"task not found"}`, http.StatusNotFound)
	})
	// Create sets the resource ID before Read runs, so even though Read
	// fails here, Terraform still considers the resource tainted/created and
	// will attempt to destroy it during test cleanup.
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipEventServiceDiscoveryCreate(taskid),
				ExpectError: regexp.MustCompile(`error Reading node`),
			},
		},
	})
}

func testBigipEventServiceDiscoveryCreateNoNodes(taskid string) string {
	return fmt.Sprintf(`
resource "bigip_event_service_discovery" "test-sd" {
  taskid = "%s"
}
`, taskid)
}

// TestUnitBigipEventServiceDiscoveryCreateEmptyNodeList exercises Create
// when no "node" blocks are configured. The "node" attribute is Optional, so
// resourceServiceDiscoveryCreate must handle a nil/empty nodeList and still
// call AddServiceDiscoveryNodes successfully.
func TestUnitBigipEventServiceDiscoveryCreateEmptyNodeList(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskid)
			return
		}
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{
  "result": {
    "providerOptions": {
      "nodeList": []
    }
  }
}`)
	})
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipEventServiceDiscoveryCreateNoNodes(taskid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_event_service_discovery.test-sd", "taskid", taskid),
					resource.TestCheckResourceAttr("bigip_event_service_discovery.test-sd", "node.#", "0"),
				),
			},
		},
	})
}

func testBigipEventServiceDiscoveryUpdate(taskid string) string {
	return fmt.Sprintf(`
resource "bigip_event_service_discovery" "test-sd" {
  taskid = "%s"
  node {
    id   = "newNode1"
    ip   = "192.168.2.3"
    port = 8080
  }
  node {
    id   = "newNode2"
    ip   = "192.168.2.4"
    port = 8080
  }
}
`, taskid)
}

func TestUnitBigipEventServiceDiscoveryUpdate(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	callCount := 0
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			callCount++
			_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskid)
			return
		}
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		if callCount <= 1 {
			_, _ = fmt.Fprintf(w, `{
  "result": {
    "providerOptions": {
      "nodeList": [
        {"id": "newNode1", "ip": "192.168.2.3", "port": 8080}
      ]
    }
  }
}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{
  "result": {
    "providerOptions": {
      "nodeList": [
        {"id": "newNode1", "ip": "192.168.2.3", "port": 8080},
        {"id": "newNode2", "ip": "192.168.2.4", "port": 8080}
      ]
    }
  }
}`)
	})
	// Destroy path after the final step.
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipEventServiceDiscoveryCreate(taskid),
			},
			{
				Config: testBigipEventServiceDiscoveryUpdate(taskid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_event_service_discovery.test-sd", "node.#", "2"),
				),
			},
		},
	})
}

// TestUnitBigipEventServiceDiscoveryUpdateError exercises the error path
// in resourceServiceDiscoveryUpdate (lines 122-125), where a failing
// AddServiceDiscoveryNodes call during Update must surface as an "error
// modifying node" diagnostic, mirroring the equivalent Create error test.
func TestUnitBigipEventServiceDiscoveryUpdateError(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	postCount := 0
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			postCount++
			// The first POST (Create) succeeds; the second POST (Update)
			// fails, simulating AddServiceDiscoveryNodes returning an error
			// when applying the updated node set.
			if postCount > 1 {
				http.Error(w, `{"code":400,"message":"invalid node configuration"}`, http.StatusBadRequest)
				return
			}
			_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskid)
			return
		}
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{
  "result": {
    "providerOptions": {
      "nodeList": [
        {"id": "newNode1", "ip": "192.168.2.3", "port": 8080}
      ]
    }
  }
}`)
	})
	// Create succeeds, so test cleanup will destroy the resource; Update
	// itself never reaches the DELETE endpoint, so this only needs to serve
	// the post-test destroy.
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipEventServiceDiscoveryCreate(taskid),
			},
			{
				Config:      testBigipEventServiceDiscoveryUpdate(taskid),
				ExpectError: regexp.MustCompile(`error modifying node`),
			},
		},
	})
}

func TestUnitBigipEventServiceDiscoveryDelete(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskid)
			return
		}
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{
  "result": {
    "providerOptions": {
      "nodeList": [
        {"id": "newNode1", "ip": "192.168.2.3", "port": 8080}
      ]
    }
  }
}`)
	})
	// Delete uses a raw *http.Client (not the go-bigip client) and POSTs an
	// empty array to the same nodes endpoint, addressed without the leading
	// slash duplication since clientBigip.Host already has no trailing slash.
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipEventServiceDiscoveryCreate(taskid),
			},
		},
	})
}

// TestUnitBigipEventServiceDiscoveryDeleteError exercises the status
// check in resourceServiceDiscoveryDelete (resp.Status != "200 OK"), which
// is otherwise untested. The mock returns a 500 on the first DELETE attempt
// (driven by the explicit Destroy step below) so the "error while
// Sending/Posting http request for Delete operation" diagnostic is
// exercised, and a 200 on any subsequent attempt so the SDK's own
// post-test cleanup destroy can still succeed.
func TestUnitBigipEventServiceDiscoveryDeleteError(t *testing.T) {
	taskid := "~Sample_event_sd~My_app~My_pool"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"id":"%s"}`, taskid)
			return
		}
		_, _ = fmt.Fprintf(w, `{
  "result": {
    "providerOptions": {
      "nodeList": [
        {"id": "newNode1", "ip": "192.168.2.3", "port": 8080}
      ]
    }
  }
}`)
	})
	deleteAttempts := 0
	mux.HandleFunc("/mgmt/shared/service-discovery/task/"+taskid+"/nodes/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		deleteAttempts++
		if deleteAttempts == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"code":500,"message":"internal server error"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipEventServiceDiscoveryCreate(taskid),
			},
			{
				Config:      testBigipEventServiceDiscoveryCreate(taskid),
				Destroy:     true,
				ExpectError: regexp.MustCompile(`error while Sending/Posting http request for Delete operation`),
			},
		},
	})
}

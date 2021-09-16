/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var resHttp2Name = "bigip_ltm_profile_http2"

func TestAccBigipLtmProfileHttp2CreateDefault(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-basic"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2DefaultConfig(TEST_PARTITION, TestHttp2Name, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
				Destroy: false,
			},
		},
	})
}

/*
TestAccBigipLtmProfileHttp2ModifyName Testcase added to check Name modification forces replacement
*/
func TestAccBigipLtmProfileHttp2ModifyName(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-old"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2DefaultConfig(TEST_PARTITION, TestHttp2Name, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config:             testaccbigipltmprofileHttp2DefaultConfig(TEST_PARTITION, fmt.Sprintf("/%s/%s", TEST_PARTITION, "test-http2-new"), instName),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(fmt.Sprintf("/%s/%s", TEST_PARTITION, "test-http2-new")),
					resource.TestCheckResourceAttr(resFullName, "name", fmt.Sprintf("/%s/%s", TEST_PARTITION, "test-http2-new")),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2UpdateConcurrentStreamsPerConnection(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-concurrentStreamsPerConnection"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "concurrent_streams_per_connection"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "concurrent_streams_per_connection", "40"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2UpdateConnectionIdleTimeout(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-connection-idle-timeout"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "connection_idle_timeout"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "connection_idle_timeout", "400"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2UpdateFrameSize(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-frame-size"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "frame_size"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "frame_size", "2058"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2UpdateHeaderTableSize(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-header-table-size"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "header_table_size"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "header_table_size", "5096"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2UpdateReceiveWindow(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-receive-window"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "receive_window"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "receive_window", "42"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2UpdateWriteSize(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-write-size"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "write_size"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "write_size", "16084"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2UpdateInsertHeader(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-insert-header"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "insert_header"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "insert_header", "enabled"),
				),
			},
		},
	})
}
func TestAccBigipLtmProfileHttp2UpdateEnforceTlsRequirements(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-enforce-tls-requirements"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "enforce_tls_requirements"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "enforce_tls_requirements", "enabled"),
				),
			},
		},
	})
}
func TestAccBigipLtmProfileHttp2UpdateIncludeContentLength(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-include-content-length"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "include_content_length"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "include_content_length", "enabled"),
				),
			},
		},
	})
}
func TestAccBigipLtmProfileHttp2UpdateInsertHeaderName(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-insert-header-name"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "insert_header_name"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, "insert_header_name", "X-HTTP2-Test"),
				),
			},
		},
	})
}
func TestAccBigipLtmProfileHttp2UpdateActivationModes(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-activation-modes"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
			},
			{
				Config: testaccbigipltmprofileHttp2UpdateParam(instName, "activation_modes"),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("activation_modes.%d", schema.HashString("always")), "always"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2Import(t *testing.T) {
	t.Parallel()
	var instName = "test-http2-import"
	var TestHttp2Name = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttp2Name, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileHttp2DefaultConfig(TEST_PARTITION, TestHttp2Name, instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttp2Name),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http2"),
				),
				ResourceName:      TestHttp2Name,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckHttp2Exists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetHttp2(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("http2 %s was not created ", name)
		}

		return nil
	}
}

func testCheckHttp2sDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_http2" {
			continue
		}

		name := rs.Primary.ID
		http2, err := client.GetHttp2(name)
		if err != nil {
			return err
		}
		if http2 != nil {
			return fmt.Errorf("http2 %s not destroyed. ", name)
		}
	}
	return nil
}

func testaccbigipltmprofileHttp2DefaultConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http2" "%[3]s" {
  name          = "%[2]s"
  defaults_from = "/%[1]s/http2"
}
`, partition, profileName, resourceName)
}

func testaccbigipltmprofileHttp2UpdateParam(instName, updateParam string) string {
	resPrefix := fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"
              defaults_from         = "/Common/http2"
		`, resHttp2Name, instName)
	switch updateParam {
	case "concurrent_streams_per_connection":
		resPrefix = fmt.Sprintf(`%s
			  concurrent_streams_per_connection = 40`, resPrefix)
	case "connection_idle_timeout":
		resPrefix = fmt.Sprintf(`%s
			  connection_idle_timeout = 400`, resPrefix)
	case "frame_size":
		resPrefix = fmt.Sprintf(`%s
			  frame_size = 2058`, resPrefix)
	case "header_table_size":
		resPrefix = fmt.Sprintf(`%s
			  header_table_size = 5096`, resPrefix)
	case "receive_window":
		resPrefix = fmt.Sprintf(`%s
			  receive_window = 42`, resPrefix)
	case "insert_header_name":
		resPrefix = fmt.Sprintf(`%s
			  insert_header_name = "X-HTTP2-Test"`, resPrefix)
	case "write_size":
		resPrefix = fmt.Sprintf(`%s
			  write_size = 16084`, resPrefix)
	case "insert_header":
		resPrefix = fmt.Sprintf(`%s
			  insert_header = "enabled"`, resPrefix)
	case "enforce_tls_requirements":
		resPrefix = fmt.Sprintf(`%s
			  enforce_tls_requirements = "enabled"`, resPrefix)
	case "include_content_length":
		resPrefix = fmt.Sprintf(`%s
			  include_content_length = "enabled"`, resPrefix)
	case "activation_modes":
		resPrefix = fmt.Sprintf(`%s
			  activation_modes = ["always"]`, resPrefix)
	default:
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}

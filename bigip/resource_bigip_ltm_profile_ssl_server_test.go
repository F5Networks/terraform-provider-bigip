/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resNameserver = "bigip_ltm_profile_server_ssl"

func TestAccBigipLtmProfileServerSsl_Default_create(t *testing.T) {
	t.Parallel()
	var instName = "test-ServerSsl"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resNameserver, instName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
					resource.TestCheckResourceAttr(resFullName, "alert_timeout", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, "authenticate_depth", "9"),
					resource.TestCheckResourceAttr(resFullName, "cache_size", "262144"),
					resource.TestCheckResourceAttr(resFullName, "ca_file", "none"),
					resource.TestCheckResourceAttr(resFullName, "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr(resFullName, "key", "/Common/default.key"),
					resource.TestCheckResourceAttr(resFullName, "chain", "none"),
					testCheckServerSSLDefaultsMatchStock(resFullName),
					resource.TestCheckResourceAttr(resFullName, "expire_cert_response_control", "drop"),
					resource.TestCheckResourceAttr(resFullName, "handshake_timeout", "10"),
					resource.TestCheckResourceAttr(resFullName, "mod_ssl_methods", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "mode", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "peer_cert_mode", "ignore"),
					resource.TestCheckResourceAttr(resFullName, "proxy_ssl", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "renegotiate_period", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "renegotiate_size", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "renegotiation", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "retain_certificate", "true"),
					resource.TestCheckResourceAttr(resFullName, "secure_renegotiation", "require-strict"),
					resource.TestCheckResourceAttr(resFullName, "server_name", "none"),
					resource.TestCheckResourceAttr(resFullName, "session_mirroring", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "session_ticket", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "sni_default", "false"),
					resource.TestCheckResourceAttr(resFullName, "sni_require", "false"),
					resource.TestCheckResourceAttr(resFullName, "ssl_forward_proxy", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "ssl_forward_proxy_bypass", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "ssl_sign_hash", "any"),
					resource.TestCheckResourceAttr(resFullName, "strict_resume", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "unclean_shutdown", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "untrusted_cert_response_control", "drop"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
func TestAccBigipLtmProfileServerSsl_UpdateAuthenticate(t *testing.T) {
	t.Parallel()
	var instName = "test-ServerSsl-UpdateAuthenticate"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resNameserver, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
				),
			},
			{
				Config: testAccBigipLtmProfileServerSsl_UpdateParam(instName, "authenticate"),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "always"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
func TestAccBigipLtmProfileServerSsl_UpdateTmoptions(t *testing.T) {
	t.Parallel()
	var instName = "test-ServerSsl-UpdateTmoptions"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resNameserver, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					// Newer BIG-IP defaults no longer consistently include inherited no-tlsv1.
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "dont-insert-empty-fragments"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
				),
			},
			{
				Config: testAccBigipLtmProfileServerSsl_UpdateParam(instName, "tm_options"),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "no-tlsv1.3"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileServerSsl_UpdateCipherGroup(t *testing.T) {
	t.Parallel()
	var instName = "test-ServerSsl-UpdateCipherGroup"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resNameserver, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
					testCheckServerSSLDefaultsMatchStock(resFullName),
				),
			},
			{
				Config: testAccBigipLtmProfileServerSsl_UpdateParam(instName, "cipher_group"),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
					resource.TestCheckResourceAttr(resFullName, "cipher_group", "/Common/f5-aes"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileServerSsl_import(t *testing.T) {
	var instName = "test-ServerSsl"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
				),
				ResourceName:      instFullName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckServerSslExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetServerSSLProfile(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("ServerSsl Profile %s was not created.", name)
		}

		return nil
	}
}

func testCheckServerSslDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resNameserver {
			continue
		}

		name := rs.Primary.ID
		ServerSsl, err := client.GetServerSSLProfile(name)
		if err != nil {
			return err
		}
		if ServerSsl != nil {
			return fmt.Errorf("ServerSsl Profile %s not destroyed.", name)
		}
	}
	return nil
}
func testaccbigipltmprofileserversslDefaultcreate(instName string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name = "/Common/%[2]s"
  //defaults_from = "/Common/serverssl"
}
		`, resNameserver, instName)
}

func testAccBigipLtmProfileServerSsl_UpdateParam(instName, updateParam string) string {
	resPrefix := fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"
			  defaults_from = "/Common/serverssl"`, resNameserver, instName)
	switch updateParam {
	case "authenticate":
		resPrefix = fmt.Sprintf(`%s
			  authenticate = "always"`, resPrefix)
	case "tm_options":
		resPrefix = fmt.Sprintf(`%s
			  tm_options = ["no-tlsv1.3"]`, resPrefix)
	case "cipher_group":
		resPrefix = fmt.Sprintf(`%s
			  cipher_group = "/Common/f5-aes"`, resPrefix)
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}

// TestAccBigipLtmProfileServerSsl_ChildInheritanceSurvivesUnrelatedUpdate is
// the bigip_ltm_profile_server_ssl analog of
// TestAccBigipLtmProfileClientSsl_ChildInheritanceSurvivesUnrelatedUpdate
// (SFDC #01262589, ENI SpA): a child profile that inherits tm_options from
// its parent via defaults_from -- and never sets tm_options itself -- must
// keep inheriting even after an unrelated attribute on the child (here,
// authenticate) is updated. Before the fix in getServerSslConfig, Read
// would backfill the child's state with the parent's inherited tm_options,
// and the very next Update (triggered by the unrelated change) would
// re-send that backfilled value to BIG-IP as an explicit "options" line,
// permanently breaking inheritance. The customer's original report was
// against bigip_ltm_profile_client_ssl only; this test is the server-ssl
// verification called for since defaults_from serverssl parent/child pairs
// (e.g. for SNI-based server-side cert selection) are a distinct, if less
// common, usage pattern that needed its own confirmation rather than
// assuming parity.
func TestAccBigipLtmProfileServerSsl_ChildInheritanceSurvivesUnrelatedUpdate(t *testing.T) {
	t.Parallel()
	parentName := "test-ServerSsl-InheritParent"
	childName := "test-ServerSsl-InheritChild"
	parentFullName := fmt.Sprintf("/%s/%s", TestPartition, parentName)
	childFullName := fmt.Sprintf("/%s/%s", TestPartition, childName)
	childResFullName := fmt.Sprintf("%s.%s", resNameserver, childName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				// Step 1: parent has tm_options set; child has defaults_from
				// pointing at the parent and does NOT set tm_options itself
				// in its own .tf. tm_options is Optional+Computed, so it is
				// expected (and correct) for the child's Terraform state to
				// reflect the inherited effective value here -- that is not
				// the bug. testCheckServerSslNoExplicitOptions instead
				// checks BIG-IP's own view of the child object directly:
				// the device-side tmOptions must still be "none", meaning
				// the child has not had an explicit options line written to
				// it and is purely inheriting from the parent.
				Config: testaccbigipltmprofileserversslInheritanceConfig(parentName, childName, "once"),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(parentFullName),
					testCheckServerSslExists(childFullName),
					resource.TestCheckResourceAttr(childResFullName, "defaults_from", parentFullName),
					testCheckServerSslNoExplicitOptions(childFullName),
				),
			},
			{
				// Step 2: update ONLY an unrelated attribute on the child
				// (authenticate). Before the fix, this Update would have
				// read the tm_options value Step 1's Read backfilled into
				// state (the parent's inherited value) and re-sent it to
				// BIG-IP as an explicit override, permanently breaking
				// inheritance. The device-side check below is the actual
				// regression assertion. Repeating the same Config a third
				// time with PlanOnly also confirms no further plan diff
				// appears for tm_options once the child has settled --
				// the specific acceptance-level guarantee called for
				// alongside the client-ssl fix.
				Config: testaccbigipltmprofileserversslInheritanceConfig(parentName, childName, "always"),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(childFullName),
					resource.TestCheckResourceAttr(childResFullName, "authenticate", "always"),
					testCheckServerSslNoExplicitOptions(childFullName),
				),
			},
			{
				// Step 3: re-apply the identical config from Step 2 with no
				// changes. PlanOnly asserts Terraform computes an empty
				// diff -- i.e. tm_options (and everything else) has
				// stabilized and does not flap/drift across repeated
				// applies on an inherited child profile.
				Config:   testaccbigipltmprofileserversslInheritanceConfig(parentName, childName, "always"),
				PlanOnly: true,
			},
		},
	})
}

func testaccbigipltmprofileserversslInheritanceConfig(parentName, childName, childAuthenticate string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name          = "/Common/%[2]s"
  defaults_from = "/Common/serverssl"
  tm_options    = ["no-tlsv1.3"]
}

resource "%[1]s" "%[3]s" {
  name          = "/Common/%[3]s"
  defaults_from = %[1]s.%[2]s.name
  authenticate  = "%[4]s"
  depends_on    = [%[1]s.%[2]s]
}
`, resNameserver, parentName, childName, childAuthenticate)
}

// explicitServerSslOptionsLineRegexp mirrors explicitOptionsLineRegexp
// (resource_bigip_ltm_profile_ssl_client_test.go): matches a tmsh-list
// "options { ... }" stanza line, tolerant of leading indentation and
// brace-spacing variance across TMOS versions/patches, anchored so it only
// matches "options" as a standalone tmsh keyword at the start of a
// (trimmed) line rather than wherever the substring "options" appears.
var explicitServerSslOptionsLineRegexp = regexp.MustCompile(`(?m)^\s*options\s*\{`)

// testCheckServerSslNoExplicitOptions is the bigip_ltm_profile_server_ssl
// analog of testCheckClientSslNoExplicitOptions: asserts, via "tmsh list",
// that the named server-ssl profile has no explicit "options { ... }" line.
//
// The plain iControl REST GET response's tmOptions field is NOT a reliable
// signal for this: BIG-IP's GET always echoes the *effective* (possibly
// inherited) tmOptions value for a profile with defaults_from set, even
// when nothing was ever explicitly configured on that specific object.
// "tmsh list" is the config-authoritative view: it only prints an
// "options { ... }" line when tmOptions has actually been written
// explicitly onto that object, which is exactly what must not happen here.
func testCheckServerSslNoExplicitOptions(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		resp, err := client.RunCommand(&bigip.BigipCommand{
			Command:     "run",
			UtilCmdArgs: fmt.Sprintf("-c \"tmsh list ltm profile server-ssl %s\"", name),
		})
		if err != nil {
			return err
		}
		if strings.Contains(resp.CommandResult, "not found") {
			return fmt.Errorf("ServerSsl Profile %s was not found: %s", name, resp.CommandResult)
		}
		if explicitServerSslOptionsLineRegexp.MatchString(resp.CommandResult) {
			return fmt.Errorf("expected profile %s to have no explicit \"options\" line (pure inheritance from parent), got tmsh output:\n%s", name, resp.CommandResult)
		}
		return nil
	}
}

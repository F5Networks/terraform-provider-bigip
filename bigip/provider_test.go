/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestPartition = "Common"

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.

// var testAccProviders map[string]*schema.Provider{}
// var testAccProvider *schema.Provider

var testAccProvider = Provider()
var testAccProviders = map[string]*schema.Provider{
	"bigip": testAccProvider,
}

func TestAccProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAcctPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("BIGIP_TEST_PARTITION"); v != "" {
		TestPartition = v
	}
	if os.Getenv("BIGIP_HOST") == "" {
		t.Fatal("Either BIGIP_TOKEN_AUTH + BIGIP_LOGIN_REF or BIGIP_USER, BIGIP_PASSWORD and BIGIP_HOST are required for tests.")
		return
	}
	if os.Getenv("BIGIP_TOKEN_VALUE") == "" && (os.Getenv("BIGIP_TOKEN_AUTH") == "" || os.Getenv("BIGIP_LOGIN_REF") == "") {
		for _, s := range [...]string{"BIGIP_USER", "BIGIP_PASSWORD"} {
			if os.Getenv(s) == "" {
				t.Fatal("Either BIGIP_TOKEN_AUTH + BIGIP_LOGIN_REF or BIGIP_USER, BIGIP_PASSWORD and BIGIP_HOST are required for tests.")
				return
			}
		}
	}
	waitForAcceptanceDeviceReady(t)
}

func testAcctUnitPreCheck(_ *testing.T, url string) {
	_ = os.Setenv("BIGIP_HOST", url)
	_ = os.Setenv("BIGIP_USER", "xxxx")
	_ = os.Setenv("BIGIP_PASSWORD", "xxx")
	_ = os.Setenv("BIGIP_TOKEN_AUTH", "false")
}

// testAcctBuildRawClient builds a *bigip.BigIP directly from the same
// environment variables the provider itself reads (see the table in
// AGENTS.md's "Provider configuration" section), independent of Terraform's
// Configure lifecycle. This is needed by PreCheck functions like
// testAcctPreCheckGtm: PreCheck runs before Terraform calls the provider's
// ConfigureContextFunc, so testAccProvider.Meta() is still nil at that
// point and can't be used to reach the device to check/fix prerequisite
// state (e.g. module provisioning) ahead of the test's own Config step.
func testAcctBuildRawClient(t *testing.T) *bigip.BigIP {
	t.Helper()
	waitForAcceptanceDeviceReady(t)

	verifyCertDisable := true
	if v := os.Getenv("BIGIP_VERIFY_CERT_DISABLE"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			t.Fatalf("invalid BIGIP_VERIFY_CERT_DISABLE value %q: %s", v, err)
		}
		verifyCertDisable = parsed
	}

	config := &bigip.Config{
		Address:            os.Getenv("BIGIP_HOST"),
		Port:               os.Getenv("BIGIP_PORT"),
		Username:           os.Getenv("BIGIP_USER"),
		Password:           os.Getenv("BIGIP_PASSWORD"),
		Token:              os.Getenv("BIGIP_TOKEN_VALUE"),
		CertVerifyDisable:  verifyCertDisable,
		TrustedCertificate: os.Getenv("BIGIP_TRUSTED_CERT_PATH"),
	}
	if tokenAuth, _ := strconv.ParseBool(os.Getenv("BIGIP_TOKEN_AUTH")); tokenAuth {
		config.LoginReference = os.Getenv("BIGIP_LOGIN_REF")
	}

	client, err := Client(config)
	if err != nil {
		t.Fatalf("unable to build BIG-IP client from environment: %s", err)
	}
	return client
}

func waitForAcceptanceDeviceReady(t *testing.T) {
	t.Helper()

	const (
		pollInterval = 5 * time.Second
		pollTimeout  = 10 * time.Minute
	)

	verifyCertDisable := true
	if v := os.Getenv("BIGIP_VERIFY_CERT_DISABLE"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			t.Fatalf("invalid BIGIP_VERIFY_CERT_DISABLE value %q: %s", v, err)
		}
		verifyCertDisable = parsed
	}

	config := &bigip.Config{
		Address:            os.Getenv("BIGIP_HOST"),
		Port:               os.Getenv("BIGIP_PORT"),
		Username:           os.Getenv("BIGIP_USER"),
		Password:           os.Getenv("BIGIP_PASSWORD"),
		Token:              os.Getenv("BIGIP_TOKEN_VALUE"),
		CertVerifyDisable:  verifyCertDisable,
		TrustedCertificate: os.Getenv("BIGIP_TRUSTED_CERT_PATH"),
	}
	if tokenAuth, _ := strconv.ParseBool(os.Getenv("BIGIP_TOKEN_AUTH")); tokenAuth {
		config.LoginReference = os.Getenv("BIGIP_LOGIN_REF")
	}
	requireFreshLogin := config.Username != "" && config.Password != ""

	var lastErr error
	deadline := time.Now().Add(pollTimeout)
	for time.Now().Before(deadline) {
		client, err := Client(config)
		if err == nil {
			if !requireFreshLogin || freshLoginSucceeds(client) {
				return
			}
			lastErr = fmt.Errorf("validate connection succeeded but fresh login is still unavailable")
		} else {
			lastErr = err
		}
		log.Printf("[DEBUG] acceptance precheck waiting for device management-plane readiness: %v", lastErr)
		time.Sleep(pollInterval)
	}

	t.Fatalf("timed out after %s waiting for BIG-IP management plane at %s to accept validation and fresh login: %v", pollTimeout, config.Address, lastErr)
}

// ensureModuleProvisioned checks whether the named BIG-IP module (e.g.
// "gtm", "asm") is currently provisioned above the "none" level, and if
// not, provisions it at the "nominal" level via the same
// client.ProvisionModule call the bigip_sys_provision resource uses.
// Provisioning a module restarts the daemon(s) that back it (e.g. gtmd),
// which does not happen synchronously with the provisioning PUT, so this
// polls client.Provisions(name) until the device reports the module is
// no longer at level "none", or fails the test after a timeout.
func ensureModuleProvisioned(t *testing.T, client *bigip.BigIP, name string) {
	t.Helper()

	prov, err := client.Provisions(name)
	if err != nil {
		if retryErr := waitForModuleProvisioningAccess(client); retryErr != nil {
			t.Fatalf("error checking %s provisioning status: device management plane did not recover: %s (initial error: %s)", name, retryErr, err)
		}
		prov, err = client.Provisions(name)
		if err != nil {
			t.Fatalf("error checking %s provisioning status after management plane recovered: %s", name, err)
		}
	}
	if prov != nil && prov.Level != "" && prov.Level != "none" {
		log.Printf("[INFO] %s module already provisioned at level %q", name, prov.Level)
		return
	}

	log.Printf("[INFO] %s module is not provisioned; provisioning at level \"nominal\"", name)
	if err := client.ProvisionModule(&bigip.Provision{Name: name, Level: "nominal"}); err != nil {
		t.Fatalf("error provisioning %s module: %s", name, err)
	}

	const (
		pollInterval = 5 * time.Second
		pollTimeout  = 10 * time.Minute
	)
	deadline := time.Now().Add(pollTimeout)
	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)
		prov, err = client.Provisions(name)
		if err != nil {
			// The device may be restarting daemons for the newly
			// provisioned module; treat errors here as transient and
			// keep polling until the timeout.
			log.Printf("[DEBUG] error polling %s provisioning status (will retry): %s", name, err)
			continue
		}
		if prov != nil && prov.Level != "" && prov.Level != "none" {
			if freshLoginSucceeds(client) {
				log.Printf("[INFO] %s module provisioned at level %q and fresh login succeeded", name, prov.Level)
				return
			}
			log.Printf("[DEBUG] %s module reports level %q but fresh login is still unavailable; waiting for device to settle", name, prov.Level)
		}
	}
	t.Fatalf("timed out after %s waiting for %s module provisioning and fresh-login readiness", pollTimeout, name)
}

func waitForModuleProvisioningAccess(client *bigip.BigIP) error {
	const (
		pollInterval = 5 * time.Second
		pollTimeout  = 10 * time.Minute
	)

	var lastErr error
	deadline := time.Now().Add(pollTimeout)
	for time.Now().Before(deadline) {
		if err := client.ValidateConnection(); err == nil {
			if freshLoginSucceeds(client) {
				return nil
			}
			lastErr = fmt.Errorf("validate connection succeeded but fresh login is still unavailable")
		} else {
			lastErr = err
		}
		log.Printf("[DEBUG] waiting for device management plane to recover before provisioning checks: %v", lastErr)
		time.Sleep(pollInterval)
	}

	return fmt.Errorf("timed out after %s waiting for device management-plane recovery: %w", pollTimeout, lastErr)
}

// testAcctPreCheckGtm wraps testAcctPreCheck and additionally ensures the
// target BIG-IP has the GTM module provisioned before a GTM acceptance
// test runs its Terraform config. GTM resources (datacenter, pool,
// wideip, server, monitor, ...) all require the module to be provisioned;
// on a BIG-IP where it isn't, every GTM API call fails with "Public URI
// path not registered" rather than a clear provisioning error, so this
// pre-check provisions GTM proactively instead of letting every GTM test
// fail with a confusing error on an unprovisioned device.
func testAcctPreCheckGtm(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	ensureModuleProvisioned(t, client, "gtm")
	ensureGtmMonitorTemplatesReady(t, client)
}

func testAcctPreCheckFast(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	if err := ensureFastHealthy(client); err != nil {
		t.Skipf("Skipping FAST acceptance test: %s", err)
	}
}

type fastInfoResponse struct {
	InstalledTemplates []struct {
		Name      string `json:"name,omitempty"`
		Supported bool   `json:"supported,omitempty"`
		Templates []struct {
			Name string `json:"name,omitempty"`
		} `json:"templates,omitempty"`
	} `json:"installedTemplates,omitempty"`
}

func ensureFastHealthy(client *bigip.BigIP) error {
	info, err := getFastInfo(client)
	if err != nil {
		return fmt.Errorf("unable to query FAST info: %w", err)
	}
	if !fastTemplateSupported(info, fastTmpl) {
		return fmt.Errorf("required FAST template %q is not installed and supported", fastTmpl)
	}

	tasks, err := getFastTasks(client)
	if err != nil {
		return fmt.Errorf("unable to query FAST tasks: %w", err)
	}

	const staleTaskAge = 5 * time.Minute
	now := time.Now().UTC()
	var stale []string
	for _, task := range tasks {
		if task.Code == 200 || task.Code >= 400 {
			continue
		}
		createdAt, parseErr := parseFastTaskTimestamp(task.Timestamp)
		if parseErr != nil {
			stale = append(stale, fmt.Sprintf("%s(code=%d,message=%q,timestamp=%q)", task.Id, task.Code, task.Message, task.Timestamp))
			continue
		}
		if now.Sub(createdAt) >= staleTaskAge {
			stale = append(stale, fmt.Sprintf("%s(code=%d,message=%q,age=%s)", task.Id, task.Code, task.Message, now.Sub(createdAt).Round(time.Second)))
		}
	}
	if len(stale) > 0 {
		return fmt.Errorf("FAST has stale non-terminal tasks: %s", strings.Join(stale, "; "))
	}

	return nil
}

func getFastInfo(client *bigip.BigIP) (*fastInfoResponse, error) {
	var info fastInfoResponse
	resp, err := client.APICall(&bigip.APIRequest{Method: http.MethodGet, URL: "mgmt/shared/fast/info"})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resp, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func getFastTasks(client *bigip.BigIP) ([]bigip.FastTask, error) {
	resp, err := client.APICall(&bigip.APIRequest{Method: http.MethodGet, URL: "mgmt/shared/fast/tasks"})
	if err != nil {
		return nil, err
	}
	var tasks []bigip.FastTask
	if err := json.Unmarshal(resp, &tasks); err == nil {
		return tasks, nil
	}
	var wrapped struct {
		Items []bigip.FastTask `json:"items,omitempty"`
	}
	if err := json.Unmarshal(resp, &wrapped); err != nil {
		return nil, err
	}
	return wrapped.Items, nil
}

func fastTemplateSupported(info *fastInfoResponse, templateName string) bool {
	for _, set := range info.InstalledTemplates {
		if !set.Supported {
			continue
		}
		for _, tmpl := range set.Templates {
			if tmpl.Name == templateName {
				return true
			}
		}
	}
	return false
}

func parseFastTaskTimestamp(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("empty FAST task timestamp")
	}
	return time.Parse(time.RFC3339, value)
}

func ensureGtmMonitorTemplatesReady(t *testing.T, client *bigip.BigIP) {
	t.Helper()

	requiredTemplates := []struct {
		fullPath string
		kind     string
	}{
		{fullPath: "/Common/http", kind: "http"},
		{fullPath: "/Common/https", kind: "https"},
		{fullPath: "/Common/tcp", kind: "tcp"},
		{fullPath: "/Common/postgresql", kind: "postgresql"},
		{fullPath: "/Common/bigip", kind: "bigip"},
	}

	const (
		pollInterval = 5 * time.Second
		pollTimeout  = 10 * time.Minute
	)

	deadline := time.Now().Add(pollTimeout)
	for _, template := range requiredTemplates {
		var lastErr error
		for time.Now().Before(deadline) {
			monitor, err := client.GetGtmMonitor(template.fullPath, template.kind)
			if err == nil && monitor != nil {
				log.Printf("[INFO] GTM monitor template %s (%s) is available", template.fullPath, template.kind)
				lastErr = nil
				break
			}
			if err != nil {
				lastErr = err
			} else {
				lastErr = fmt.Errorf("monitor template %s returned nil", template.fullPath)
			}
			if validateErr := client.ValidateConnection(); validateErr == nil {
				if freshLoginSucceeds(client) {
					log.Printf("[DEBUG] GTM monitor template %s (%s) still unavailable after successful fresh login: %v", template.fullPath, template.kind, lastErr)
				} else {
					lastErr = fmt.Errorf("validate connection succeeded but fresh login is still unavailable while waiting for GTM template %s", template.fullPath)
				}
			} else {
				lastErr = validateErr
			}
			log.Printf("[DEBUG] waiting for GTM monitor template %s (%s): %v", template.fullPath, template.kind, lastErr)
			time.Sleep(pollInterval)
		}
		if lastErr != nil {
			t.Fatalf("timed out after %s waiting for GTM monitor template %s (%s): %v", pollTimeout, template.fullPath, template.kind, lastErr)
		}
	}
}

// ensureRouteDomainExists checks whether a route domain with the given
// numeric ID already exists on the target BIG-IP, and if not, creates an
// empty one (no VLAN members) via the same client.CreateRouteDomain call
// the bigip_partition resource's route_domain_id attribute depends on. A
// freshly-provisioned BIG-IP only has the default route domain 0; any
// other route domain a test references (e.g. bigip_partition's
// route_domain_id = 2) must already exist on the device, or apply fails
// with "The requested route-domain (N) was not found."
func ensureRouteDomainExists(t *testing.T, client *bigip.BigIP, id int) {
	t.Helper()

	routeDomains, err := client.RouteDomains()
	if err != nil {
		t.Fatalf("error listing route domains: %s", err)
	}
	for _, rd := range routeDomains.RouteDomains {
		if rd.ID == id {
			log.Printf("[INFO] route domain %d already exists", id)
			return
		}
	}

	log.Printf("[INFO] route domain %d does not exist; creating it", id)
	name := strconv.Itoa(id)
	if err := client.CreateRouteDomain(name, id, true, ""); err != nil {
		t.Fatalf("error creating route domain %d: %s", id, err)
	}
}

// ensureGTMProberPoolExists checks whether a GTM prober pool with the
// given full-path name already exists, and if not, creates an empty one
// (no members) via client.CreateGTMProberPool. bigip_gtm_datacenter's
// prober_preference/prober_fallback = "pool" requires prober_pool to
// reference a real, existing GTM prober pool -- BIG-IP rejects it
// otherwise with "Select a valid Prober Pool value".
func ensureGTMProberPoolExists(t *testing.T, client *bigip.BigIP, fullPath string) {
	t.Helper()

	pool, err := client.GetGTMProberPool(fullPath)
	if err != nil {
		t.Fatalf("error checking GTM prober pool %s: %s", fullPath, err)
	}
	if pool != nil {
		log.Printf("[INFO] GTM prober pool %s already exists", fullPath)
		return
	}

	partition, name := partitionAndName(fullPath)
	log.Printf("[INFO] GTM prober pool %s does not exist; creating it", fullPath)
	if err := client.CreateGTMProberPool(&bigip.GTMProberPool{
		Name:      name,
		Partition: partition,
	}); err != nil {
		t.Fatalf("error creating GTM prober pool %s: %s", fullPath, err)
	}
}

// testAcctPreCheckGtmProberPool wraps testAcctPreCheckGtm and additionally
// ensures the "/Common/test-prober-pool" GTM prober pool
// TestAccBigipGtmDatacenter_withProberSettings's Terraform config
// references exists on the target BIG-IP. See ensureGTMProberPoolExists.
func testAcctPreCheckGtmProberPool(t *testing.T) {
	t.Helper()
	testAcctPreCheckGtm(t)

	client := testAcctBuildRawClient(t)
	ensureGTMProberPoolExists(t, client, "/Common/test-prober-pool")
}

// testAcctPreCheckRouteDomain2 wraps testAcctPreCheck and additionally
// ensures route domain 2 exists on the target BIG-IP before a test whose
// Terraform config references it (e.g. bigip_partition's
// route_domain_id = 2) runs. See ensureRouteDomainExists.
func testAcctPreCheckRouteDomain2(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	ensureRouteDomainExists(t, client, 2)
}

// testAcctPreCheckRouteDomain50 wraps testAcctPreCheck and additionally
// ensures route domain 50 exists on the target BIG-IP before a test whose
// Terraform config references it (e.g. bigip_ltm_virtual_address's
// name = "/Common/1.1.1.1%50"). See ensureRouteDomainExists.
func testAcctPreCheckRouteDomain50(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	ensureRouteDomainExists(t, client, 50)
}

// testAcctPreCheckRouteDomains wraps testAcctPreCheck and additionally
// ensures each of the given route domain IDs exists on the target
// BIG-IP. See ensureRouteDomainExists.
func testAcctPreCheckRouteDomains(ids ...int) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		testAcctPreCheck(t)

		client := testAcctBuildRawClient(t)
		for _, id := range ids {
			ensureRouteDomainExists(t, client, id)
		}
	}
}

// ensureTrafficMatchingCriteriaExists checks whether an LTM traffic
// matching criteria object with the given full-path name already exists,
// and if not, creates a minimal one. bigip_ltm_virtual_server's
// trafficmatching_criteria attribute requires a pre-existing object of
// this type; this provider has no resource to manage the type itself.
func ensureTrafficMatchingCriteriaExists(t *testing.T, client *bigip.BigIP, fullPath string) {
	t.Helper()

	existing, err := client.GetTrafficMatchingCriterion(fullPath)
	if err != nil {
		t.Fatalf("error checking traffic matching criteria %s: %s", fullPath, err)
	}
	if existing != nil {
		log.Printf("[INFO] traffic matching criteria %s already exists", fullPath)
		return
	}

	partition, name := partitionAndName(fullPath)
	log.Printf("[INFO] traffic matching criteria %s does not exist; creating it", fullPath)
	if err := client.CreateTrafficMatchingCriterion(&bigip.TrafficMatchingCriteria{
		Name:      name,
		Partition: partition,
	}); err != nil {
		t.Fatalf("error creating traffic matching criteria %s: %s", fullPath, err)
	}
}

// testAcctPreCheckTrafficMatchingCriteria wraps testAcctPreCheck and
// additionally ensures the "/Common/test-virtualserver_VS_TMC_OBJ" traffic
// matching criteria object TestAccBigipLtmVirtualServerTCIssue729's
// Terraform config references exists on the target BIG-IP. See
// ensureTrafficMatchingCriteriaExists.
func testAcctPreCheckTrafficMatchingCriteria(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	ensureTrafficMatchingCriteriaExists(t, client, "/Common/test-virtualserver_VS_TMC_OBJ")
}

// ensureWafPolicyExists checks whether an ASM/WAF security policy with the
// given name and partition already exists, and if not, creates a bare
// (TMOS "Fundamental" template default) one via
// client.CreateMinimalWafPolicy. Some LTM policy rule actions reference an
// ASM policy via "asm = true" + "policy = <full path>"; BIG-IP rejects
// that with "the requested policy action ... was not found" if the
// referenced ASM policy doesn't exist.
//
// Existence is checked via client.GetWafPolicyQuery, which (unlike this
// file's other ensureXExists helpers) returns a non-nil error rather than
// (nil, nil) when the policy isn't found -- any error here is treated as
// "not found, create it" rather than a hard failure, since
// GetWafPolicyQuery doesn't distinguish a real API error from a genuine
// not-found any other way.
func ensureWafPolicyExists(t *testing.T, client *bigip.BigIP, name, partition string) {
	t.Helper()

	if _, err := client.GetWafPolicyQuery(name, partition); err == nil {
		log.Printf("[INFO] WAF policy %s/%s already exists", partition, name)
		return
	}

	log.Printf("[INFO] WAF policy %s/%s does not exist; creating it", partition, name)
	if err := client.CreateMinimalWafPolicy(name, partition); err != nil {
		t.Fatalf("error creating WAF policy %s/%s: %s", partition, name, err)
	}
}

// testAcctPreCheckLtmPoolsAndWafPolicy wraps testAcctPreCheck and ensures
// both the given LTM pools and the "f5-waf-profile" ASM policy exist on
// the target BIG-IP, used by TestAccBigipLtmPolicyIssue737 (whose config
// references both pools and an ASM policy via policy rule actions).
func testAcctPreCheckLtmPoolsAndWafPolicy(fullPaths ...string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		testAcctPreCheck(t)

		client := testAcctBuildRawClient(t)
		for _, fullPath := range fullPaths {
			ensurePoolExists(t, client, fullPath)
		}
		ensureWafPolicyExists(t, client, "f5-waf-profile", "Common")
	}
}

// ensurePartitionExists checks whether a top-level partition with the
// given name already exists, and if not, creates it via
// client.CreatePartition (POST /mgmt/tm/auth/partition). A handful of
// bigip_ltm_policy acceptance tests reference policy names under a
// "/TEST" partition that must already exist on the device.
func ensurePartitionExists(t *testing.T, client *bigip.BigIP, name string) {
	t.Helper()

	existing, err := client.GetPartition(name)
	if err != nil {
		t.Fatalf("error checking partition %s: %s", name, err)
	}
	if existing != nil {
		log.Printf("[INFO] partition %s already exists", name)
		return
	}

	log.Printf("[INFO] partition %s does not exist; creating it", name)
	if err := client.CreatePartition(&bigip.Partition{Name: name}); err != nil {
		t.Fatalf("error creating partition %s: %s", name, err)
	}
}

// ensureSysFolderExists checks whether a sys folder with the given
// full-path name already exists, and if not, creates it via
// client.CreateSysFolder (POST /mgmt/tm/sys/folder). Used for subfolders
// within a partition (e.g. "/TEST/A1"); the partition itself
// ("/TEST") must already exist -- see ensurePartitionExists.
func ensureSysFolderExists(t *testing.T, client *bigip.BigIP, fullPath string) {
	t.Helper()

	existing, err := client.GetSysFolder(fullPath)
	if err != nil {
		t.Fatalf("error checking sys folder %s: %s", fullPath, err)
	}
	if existing != nil {
		log.Printf("[INFO] sys folder %s already exists", fullPath)
		return
	}

	partition, name := partitionAndName(fullPath)
	log.Printf("[INFO] sys folder %s does not exist; creating it", fullPath)
	if err := client.CreateSysFolder(&bigip.SysFolder{
		Name:      name,
		Partition: "/" + partition,
	}); err != nil {
		t.Fatalf("error creating sys folder %s: %s", fullPath, err)
	}
}

// testAcctPreCheckTestPartition wraps testAcctPreCheck and additionally
// ensures the "/TEST" partition exists on the target BIG-IP, used by
// several bigip_ltm_policy acceptance tests (Issue634_a, Issue132_c).
func testAcctPreCheckTestPartition(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	ensurePartitionExists(t, client, "TEST")
}

// testAcctPreCheckTestPartitionWithA1Folder wraps
// testAcctPreCheckTestPartition and additionally ensures the "/TEST/A1"
// subfolder exists, used by TestAccBigipLtmPolicy_Issue634.
func testAcctPreCheckTestPartitionWithA1Folder(t *testing.T) {
	t.Helper()
	testAcctPreCheckTestPartition(t)

	client := testAcctBuildRawClient(t)
	ensureSysFolderExists(t, client, "/TEST/A1")
}

// testAcctPreCheckPartition wraps testAcctPreCheck and additionally
// ensures the given top-level partition exists on the target BIG-IP. See
// ensurePartitionExists.
func testAcctPreCheckPartition(name string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		testAcctPreCheck(t)

		client := testAcctBuildRawClient(t)
		ensurePartitionExists(t, client, name)
	}
}

// testAcctPreCheckPartitionWithFolder wraps testAcctPreCheckPartition and
// additionally ensures the given subfolder (full path, e.g.
// "/TEST_iFile_300/A1TEST") exists under that partition. See
// ensureSysFolderExists.
func testAcctPreCheckPartitionWithFolder(partition, subfolderFullPath string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		testAcctPreCheckPartition(partition)(t)

		client := testAcctBuildRawClient(t)
		ensureSysFolderExists(t, client, subfolderFullPath)
	}
}

// ensureLtmPolicyExists checks whether a published LTM policy with the
// given full path (e.g. "/Common/test_policy" or
// "/Common/folder1/ecopolicy") already exists, and if not, creates a
// minimal one -- a single rule with one forward action to poolFullPath --
// via client.CreatePolicy + client.PublishPolicy. The draft/publish name
// construction here exactly mirrors dataToPolicy and
// resourceBigipLtmPolicyCreate in resource_bigip_ltm_policy.go (including
// their inconsistency between the partition-only and nested-subfolder
// cases -- a top-level "/Common/x" policy's draft is named "Drafts/x" with
// no partition prefix, while a nested "/Common/sub/x" policy's draft is
// named "/Common/sub/Drafts/x" with the full prefix); this is intentional
// replication of that existing behavior, not a design choice made here,
// so that policies created this way behave identically to ones created
// through the real resource. poolFullPath must already exist (see
// ensurePoolExists).
//
// Existence is checked via client.FetchPolicy rather than client.GetPolicy:
// GetPolicy has its own bug (it builds the request path as "partition~name"
// without the leading "~" iControlPath's mangled-path convention requires,
// e.g. "Common~test_policy" instead of "~Common~test_policy"), which
// always 404s even for a policy that does exist. FetchPolicy (used by the
// real datasource/resource Read paths) builds the path correctly.
func ensureLtmPolicyExists(t *testing.T, client *bigip.BigIP, fullPath, poolFullPath string) {
	t.Helper()

	polStr := strings.Split(fullPath, "/")
	var partition string
	if len(polStr[:len(polStr)-1]) > 1 {
		partition = strings.Join(polStr[:len(polStr)-1], "/")
	} else {
		partition = polStr[:len(polStr)-1][0]
	}
	policyName := polStr[len(polStr)-1]

	fetchFields := strings.Split(strings.Trim(fullPath, "/"), "/")
	existing, err := client.FetchPolicy(fetchFields...)
	if err != nil {
		t.Fatalf("error checking LTM policy %s: %s", fullPath, err)
	}
	if existing != nil {
		log.Printf("[INFO] LTM policy %s already exists", fullPath)
		return
	}

	log.Printf("[INFO] LTM policy %s does not exist; creating it", fullPath)
	var draftName string
	if partition == "/Common" {
		draftName = "Drafts/" + policyName
	} else {
		draftName = partition + "/Drafts/" + policyName
	}
	p := bigip.Policy{
		Name:     draftName,
		Strategy: "/Common/first-match",
		Requires: []string{"http"},
		Controls: []string{"forwarding"},
		Rules: []bigip.PolicyRule{
			{
				Name: "rule1",
				Actions: []bigip.PolicyRuleAction{
					{
						Name:    "0",
						Forward: true,
						Pool:    poolFullPath,
					},
				},
			},
		},
	}
	if err := client.CreatePolicy(&p); err != nil {
		t.Fatalf("error creating LTM policy draft %s: %s", fullPath, err)
	}
	publishedCopy := partition + "/Drafts/" + policyName
	if err := client.PublishPolicy(policyName, publishedCopy); err != nil {
		t.Fatalf("error publishing LTM policy %s: %s", fullPath, err)
	}
}

// testAcctPreCheckLtmPolicyFixtures wraps testAcctPreCheck and ensures the
// three fixture LTM policies the TestAccDataSourceBigipLtmPolicy_* tests
// read (but never create themselves, since they're pure data-source read
// tests) exist on the target BIG-IP: /Common/test_policy,
// /Common/folder1/ecopolicy, and /Common/testpolicy. See
// ensureLtmPolicyExists.
func testAcctPreCheckLtmPolicyFixtures(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	const poolFullPath = "/Common/test-pool-for-policy-datasource-fixtures"
	ensurePoolExists(t, client, poolFullPath)
	ensureSysFolderExists(t, client, "/Common/folder1")
	ensureLtmPolicyExists(t, client, "/Common/test_policy", poolFullPath)
	ensureLtmPolicyExists(t, client, "/Common/folder1/ecopolicy", poolFullPath)
	ensureLtmPolicyExists(t, client, "/Common/testpolicy", poolFullPath)
}

// ensurePoolExists checks whether an LTM pool with the given full-path
// name already exists, and if not, creates an empty one (no members) via
// client.CreatePool.
func ensurePoolExists(t *testing.T, client *bigip.BigIP, name string) {
	t.Helper()

	pool, err := client.GetPool(name)
	if err != nil {
		t.Fatalf("error checking pool %s: %s", name, err)
	}
	if pool != nil {
		log.Printf("[INFO] pool %s already exists", name)
		return
	}

	log.Printf("[INFO] pool %s does not exist; creating it", name)
	if err := client.CreatePool(name); err != nil {
		t.Fatalf("error creating pool %s: %s", name, err)
	}
}

// ensureServerSSLProfileExists checks whether a server-ssl profile with
// the given full-path name already exists, and if not, creates one
// inheriting from the stock "/Common/serverssl" profile.
func ensureServerSSLProfileExists(t *testing.T, client *bigip.BigIP, fullPath string) {
	t.Helper()

	profile, err := client.GetServerSSLProfile(fullPath)
	if err != nil {
		t.Fatalf("error checking server-ssl profile %s: %s", fullPath, err)
	}
	if profile != nil {
		log.Printf("[INFO] server-ssl profile %s already exists", fullPath)
		return
	}

	partition, name := partitionAndName(fullPath)
	log.Printf("[INFO] server-ssl profile %s does not exist; creating it", fullPath)
	if err := client.AddServerSSLProfile(&bigip.ServerSSLProfile{
		Name:         name,
		Partition:    partition,
		DefaultsFrom: "/Common/serverssl",
	}); err != nil {
		t.Fatalf("error creating server-ssl profile %s: %s", fullPath, err)
	}
}

func testCheckClientSSLDefaultsMatchStock(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		stock, err := client.GetClientSSLProfile("/Common/clientssl")
		if err != nil {
			return fmt.Errorf("error reading stock clientssl profile: %w", err)
		}
		if stock == nil {
			return fmt.Errorf("stock clientssl profile /Common/clientssl not found")
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		if got := rs.Primary.Attributes["cipher_group"]; got != stock.CipherGroup {
			return fmt.Errorf("cipher_group = %q, want stock /Common/clientssl cipher_group %q", got, stock.CipherGroup)
		}

		return nil
	}
}

func testCheckServerSSLDefaultsMatchStock(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		stock, err := client.GetServerSSLProfile("/Common/serverssl")
		if err != nil {
			return fmt.Errorf("error reading stock serverssl profile: %w", err)
		}
		if stock == nil {
			return fmt.Errorf("stock serverssl profile /Common/serverssl not found")
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		if got := rs.Primary.Attributes["cipher_group"]; got != stock.CipherGroup {
			return fmt.Errorf("cipher_group = %q, want stock /Common/serverssl cipher_group %q", got, stock.CipherGroup)
		}
		if got := rs.Primary.Attributes["ciphers"]; got != stock.Ciphers {
			return fmt.Errorf("ciphers = %q, want stock /Common/serverssl ciphers %q", got, stock.Ciphers)
		}

		return nil
	}
}

// partitionAndName splits a "/partition/name" full-path string into its
// two components. Falls back to ("Common", fullPath) if fullPath doesn't
// have the expected leading-slash, two-segment form.
func partitionAndName(fullPath string) (partition, name string) {
	trimmed := strings.TrimPrefix(fullPath, "/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return "Common", fullPath
	}
	return parts[0], parts[1]
}

// testAcctPreCheckSaasBotDefense wraps testAcctPreCheck and additionally
// ensures the prerequisites TestAccBigipSaasBotDefenseProfileTC1's
// Terraform config depends on exist on the target BIG-IP: the
// shape_protection_pool it references and the ssl_profile it references.
// (defaults_from's default, "/Common/bd", is itself a stock saas bd
// profile TMOS ships out of the box -- see AddSaasBotDefenseProfile's
// uriSaas/uriSaasBotDefense endpoint -- so no prerequisite is needed for
// it.)
func testAcctPreCheckSaasBotDefense(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	ensurePoolExists(t, client, "/Common/cs1.pool")
	ensureServerSSLProfileExists(t, client, "/Common/cloud-service-default-ssl")
}

// generateCAAndLeafCert generates a fresh CA key/certificate pair plus a
// leaf certificate signed by that same CA key. Returns (caCertPEM,
// leafCertPEM, leafKeyPEM).
func generateCAAndLeafCert(t *testing.T, caCommonName, leafCommonName string) (caCertPEM, leafCertPEM, leafKeyPEM []byte) {
	t.Helper()

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("error generating CA private key: %s", err)
	}
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: caCommonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("error creating self-signed CA certificate: %s", err)
	}
	caCertPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("error parsing freshly-generated CA certificate: %s", err)
	}

	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("error generating leaf private key: %s", err)
	}
	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: leafCommonName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, caCert, &leafKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("error creating leaf certificate: %s", err)
	}
	leafCertPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: leafDER})
	leafKeyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(leafKey)})

	return caCertPEM, leafCertPEM, leafKeyPEM
}

// installMyCA (re-)installs "/Common/MyCA" on the target BIG-IP using the
// given already-generated CA certificate PEM (see generateCAAndLeafCert;
// callers are expected to generate the CA and a matching leaf certificate
// together, then pass the CA half here and use the leaf half directly in
// their Terraform config -- see TestAccBigipSslCertificateOCSP for the
// pattern).
//
// Unlike the other ensureXExists helpers in this file, this always
// overwrites /Common/MyCA rather than skipping if it already exists:
// BIG-IP never exposes a certificate's private key for re-extraction, so
// if /Common/MyCA already exists from a prior run there is no way to sign
// a new leaf certificate against its original CA key. Regenerating both
// together as a matched pair on every run is the only way to guarantee
// BIG-IP's issuer_cert chain validation (which cryptographically checks
// that the issuer_cert object's key actually signed the target
// certificate, not just a name/DN match) succeeds.
func installMyCA(t *testing.T, client *bigip.BigIP, caCertPEM []byte) {
	t.Helper()

	const fullPath = "/Common/MyCA"
	partition, name := partitionAndName(fullPath)

	existing, err := client.GetCertificate(fullPath)
	if err != nil {
		t.Fatalf("error checking certificate %s: %s", fullPath, err)
	}
	if existing != nil {
		log.Printf("[INFO] certificate %s already exists; reinstalling with a freshly-generated CA so a matching leaf cert can be signed", fullPath)
		if err := client.DeleteCertificate(fullPath); err != nil {
			t.Fatalf("error deleting existing certificate %s before reinstalling: %s", fullPath, err)
		}
	}

	log.Printf("[INFO] installing certificate %s", fullPath)
	if err := client.UploadCertificate(string(caCertPEM), &bigip.Certificate{
		Name:      name,
		Partition: partition,
	}); err != nil {
		t.Fatalf("error installing certificate %s: %s", fullPath, err)
	}
}

// ensureOCSPResponderExists checks whether a sys crypto cert-validator
// ocsp responder with the given full-path name already exists, and if
// not, creates a minimal one using the stock "/Common/f5-aws-dns" DNS
// resolver (the same resolver TestAccBigipSysOCSP_create's own
// testSysOcspDNS fixture uses) via client.CreateOCSP.
func ensureOCSPResponderExists(t *testing.T, client *bigip.BigIP, fullPath string) {
	t.Helper()

	ocsp, err := client.GetOCSP(fullPath)
	if err != nil {
		t.Fatalf("error checking OCSP responder %s: %s", fullPath, err)
	}
	if ocsp != nil {
		log.Printf("[INFO] OCSP responder %s already exists", fullPath)
		return
	}

	partition, name := partitionAndName(fullPath)
	log.Printf("[INFO] OCSP responder %s does not exist; creating it", fullPath)
	if err := client.CreateOCSP(&bigip.OCSP{
		Name:        name,
		Partition:   partition,
		DnsResolver: "/Common/f5-aws-dns",
	}); err != nil {
		t.Fatalf("error creating OCSP responder %s: %s", fullPath, err)
	}
}

// testAcctPreCheckOCSP returns a PreCheck function that wraps
// testAcctPreCheck and additionally (re-)installs "/Common/MyCA" on the
// target BIG-IP using caCertPEM, and ensures the "/Common/testocsp1" OCSP
// responder exists. caCertPEM must be the CA half of a
// generateCAAndLeafCert pair generated by the caller; the leaf half should
// be embedded directly in the Terraform config passed as that same test's
// Steps[].Config (see TestAccBigipSslCertificateOCSP for the pattern) --
// see installMyCA's doc comment for why a static fixture file can't work
// here instead.
func testAcctPreCheckOCSP(caCertPEM []byte) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		testAcctPreCheck(t)

		client := testAcctBuildRawClient(t)
		installMyCA(t, client, caCertPEM)
		ensureOCSPResponderExists(t, client, "/Common/testocsp1")
	}
}

// ensureKeyCertPairExists checks whether a "<base>.crt" certificate
// already exists for the given base full-path name (e.g.
// "/Common/le-ssl" -> "/Common/le-ssl.crt"), and if not, generates a
// fresh self-signed key/certificate pair and installs both under that
// base name (key as "<base>.key", cert as "<base>.crt" -- BIG-IP's sys
// file ssl-key/ssl-cert objects always carry these suffixes on the
// actual stored object, even though bigip_sys_ocsp's signer_key/
// signer_cert attributes must reference them with the suffix included,
// e.g. "/Common/le-ssl.key"/"/Common/le-ssl.crt", not the schema's
// originally-documented (and incorrect) suffix-less "/Common/le-ssl"
// form). Used as a prerequisite for tests that need a signing key/cert
// pair to already exist but don't care about its actual contents or
// chain of trust (unlike installMyCA's /Common/MyCA, which must be
// regenerated on every run alongside a matching leaf cert -- see its own
// doc comment -- this one has no such constraint, since nothing here
// validates that the referenced key/cert actually signed anything).
func ensureKeyCertPairExists(t *testing.T, client *bigip.BigIP, baseFullPath string) {
	t.Helper()

	certFullPath := baseFullPath + ".crt"
	existingCert, err := client.GetCertificate(certFullPath)
	if err != nil {
		t.Fatalf("error checking certificate %s: %s", certFullPath, err)
	}
	if existingCert != nil {
		log.Printf("[INFO] certificate %s already exists", certFullPath)
		return
	}

	partition, name := partitionAndName(baseFullPath)
	_, certPEM, keyPEM := generateCAAndLeafCert(t, name, name)

	log.Printf("[INFO] certificate/key %s do not exist; creating them", baseFullPath)
	sourcePath, err := client.UploadKey(name+".key", string(keyPEM))
	if err != nil {
		t.Fatalf("error uploading key %s.key: %s", baseFullPath, err)
	}
	if err := client.AddKey(&bigip.Key{
		Name:       name + ".key",
		SourcePath: sourcePath,
		Partition:  partition,
	}); err != nil {
		t.Fatalf("error creating key %s.key: %s", baseFullPath, err)
	}
	if err := client.UploadCertificate(string(certPEM), &bigip.Certificate{
		Name:      name + ".crt",
		Partition: partition,
	}); err != nil {
		t.Fatalf("error creating certificate %s.crt: %s", baseFullPath, err)
	}
}

// testAcctPreCheckSysOcsp wraps testAcctPreCheck and additionally ensures
// the "/Common/le-ssl" key/cert pair TestAccBigipSysOCSP_create's config
// references for signer_key/signer_cert exists on the target BIG-IP. See
// ensureKeyCertPairExists.
func testAcctPreCheckSysOcsp(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	ensureKeyCertPairExists(t, client, "/Common/le-ssl")
	ensurePoolExists(t, client, "/Common/test-poolxyz")
}

// vcmpSupported reports whether the target BIG-IP has the /tm/vcmp module
// registered at all. VCMP (Virtual Clustered Multiprocessing) is a
// hardware/hypervisor-tier feature -- it is never available on a
// cloud-provisioned VE instance (e.g. this provider's OpenStack
// acceptance-test DUT), only on physical F5 chassis or a VCMP host. There
// is no prerequisite object that can be created to make it available, so
// unlike ensureModuleProvisioned/ensureRouteDomainExists this cannot be
// remediated -- callers should skip rather than fail when this returns
// false.
//
// A plain GET always returns an empty-body 404 whether or not vcmp is
// registered, so go-bigip's GetVcmpGuest (a GET) can't distinguish
// "unsupported platform" from "no guests configured yet". This bypasses
// go-bigip and checks the raw HTTP status of a GET against the vcmp
// module root directly, which is still a reliable signal even though the
// body itself is uninformative on this endpoint.
func vcmpSupported(t *testing.T, client *bigip.BigIP) bool {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, client.Host+"/mgmt/tm/vcmp", nil)
	if err != nil {
		t.Fatalf("error building vcmp capability probe request: %s", err)
	}
	if client.Token != "" {
		req.Header.Set("X-F5-Auth-Token", client.Token)
	} else {
		req.SetBasicAuth(client.User, client.Password)
	}

	httpClient := &http.Client{Transport: client.Transport}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("error probing vcmp module availability: %s", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode != http.StatusNotFound
}

// testAcctPreCheckVcmp wraps testAcctPreCheck and skips the test outright
// if the target BIG-IP doesn't have the vcmp module registered (see
// vcmpSupported), rather than letting it fail with a confusing "Public
// URI path not registered: /tm/vcmp/guest" error on a platform that can
// never support it.
func testAcctPreCheckVcmp(t *testing.T) {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	if !vcmpSupported(t, client) {
		t.Skip("skipping: vcmp module is not registered on this BIG-IP (VCMP requires physical F5 hardware or a VCMP host; unsupported on this platform)")
	}
}

// findFreeInterface returns the name of a non-"mgmt" interface that isn't
// currently a member of any VLAN, or "" if none is found. A net trunk
// cannot be created from an interface that's already assigned to a VLAN
// (BIG-IP rejects it: "It is already assigned to a vlan"), and on a
// resource-constrained DUT (e.g. a 2-physical-interface VE instance where
// both are already claimed by /Common/external and /Common/internal)
// there may genuinely be no interface available for a trunk at all.
func findFreeInterface(t *testing.T, client *bigip.BigIP) string {
	t.Helper()

	interfaces, err := client.Interfaces()
	if err != nil {
		t.Fatalf("error listing interfaces: %s", err)
	}

	vlans, err := client.Vlans()
	if err != nil {
		t.Fatalf("error listing vlans: %s", err)
	}

	claimed := make(map[string]bool)
	for _, vlan := range vlans.Vlans {
		vlanInterfaces, err := client.GetVlanInterfaces(vlan.FullPath)
		if err != nil {
			t.Fatalf("error listing interfaces for vlan %s: %s", vlan.FullPath, err)
		}
		for _, vi := range vlanInterfaces.VlanInterfaces {
			claimed[vi.Name] = true
		}
	}

	for _, iface := range interfaces.Interfaces {
		if iface.Name == "mgmt" {
			continue
		}
		if !claimed[iface.Name] {
			return iface.Name
		}
	}
	return ""
}

// testAcctPreCheckFreeInterface wraps testAcctPreCheck and skips the test
// outright if no interface on the target BIG-IP is available to bind a
// net trunk to (every non-mgmt interface already belongs to a VLAN). See
// findFreeInterface.
func testAcctPreCheckFreeInterface(t *testing.T) string {
	t.Helper()
	testAcctPreCheck(t)

	client := testAcctBuildRawClient(t)
	iface := findFreeInterface(t, client)
	if iface == "" {
		t.Skip("skipping: no interface is available to bind a trunk to on this BIG-IP (every non-mgmt interface already belongs to a VLAN)")
	}
	return iface
}

// ensureNodeDoesNotExist checks whether an LTM node with the given
// full-path name already exists, and if so, deletes it. Used defensively
// by tests with a hardcoded node address (e.g.
// TestAccBigipLtmNode_basic's 192.168.30.1) that was observed to
// intermittently collide with a leftover node from elsewhere when run
// as part of a larger parallel suite against a long-lived DUT, despite
// passing reliably (including repeated runs) in isolation -- the exact
// mechanism was not pinned down, so this clears any conflicting leftover
// defensively rather than leaving the test flaky.
func ensureNodeDoesNotExist(t *testing.T, client *bigip.BigIP, fullPath string) {
	t.Helper()

	node, err := client.GetNode(fullPath)
	if err != nil {
		t.Fatalf("error checking node %s: %s", fullPath, err)
	}
	if node == nil {
		return
	}

	log.Printf("[WARN] node %s already exists (unexpectedly); deleting it before test setup", fullPath)
	if err := client.DeleteNode(fullPath); err != nil {
		t.Fatalf("error deleting conflicting leftover node %s: %s", fullPath, err)
	}
}

// loadFixtureBytes returns the entire contents of the given file as a byte slice
func loadFixtureBytes(path string) []byte {
	contents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// loadFixtureString returns the entire contents of the given file as a string
func loadFixtureString(path string) string {
	return string(loadFixtureBytes(path))
}

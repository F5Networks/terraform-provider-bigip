package bigip

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

var TEST_PARTITION = "Common"

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

var testProviders = map[string]terraform.ResourceProvider{
	"bigip": Provider(),
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"bigip": testAccProvider,
	}
	if v := os.Getenv("BIGIP_TEST_PARTITION"); v != "" {
		TEST_PARTITION = v
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAcctPreCheck(t *testing.T) {
	if os.Getenv("BIGIP_TOKEN_AUTH") != "" && os.Getenv("BIGIP_LOGIN_REF") != "" {
		return
	}
	for _, s := range [...]string{"BIGIP_HOST", "BIGIP_USER", "BIGIP_PASSWORD"} {
		if os.Getenv(s) == "" {
			t.Fatal("Either BIGIP_TOKEN_AUTH + BIGIP_LOGIN_REF or BIGIP_USER, BIGIP_PASSWORD and BIGIP_HOST are required for tests.")
			return
		}
	}
}

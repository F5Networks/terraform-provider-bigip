package bigip

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigipSysIfileCreate(t *testing.T) {
	t.Parallel()
	ifileName := fmt.Sprintf("testitem-ifile-%d", time.Now().Unix())
	testAccIfileResource := fmt.Sprintf(`
resource "bigip_sys_ifile" "testifile" {
  name      = "%s"
  partition = "Common"
  sub_path  = "ravi"
  content   = "dummy content"
}
`, ifileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIfileResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "name", ifileName),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "sub_path", "ravi"),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "content", "dummy content"),
				),
			},
		},
	})
}

func TestAccBigipSysIfileUpdate(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	ifileName := fmt.Sprintf("testitem-ifile-update-%d", timestamp)

	// Initial configuration
	initialConfig := fmt.Sprintf(`
resource "bigip_sys_ifile" "testifile" {
  name      = "%s"
  partition = "Common"
  sub_path  = "ravi"
  content   = "initial content"
}
`, ifileName)

	// Updated configuration
	updatedConfig := fmt.Sprintf(`
resource "bigip_sys_ifile" "testifile" {
  name      = "%s"
  partition = "Common"
  sub_path  = "ravi"
  content   = "updated content with more data"
}
`, ifileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with initial values
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "name", ifileName),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "content", "initial content"),
				),
			},
			// Step 2: Update the content
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "name", ifileName),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "content", "updated content with more data"),
					// Verify other attributes remain unchanged
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "sub_path", "ravi"),
				),
			},
		},
	})
}

func TestAccBigipSysIfileNoChange(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	ifileName := fmt.Sprintf("testitem-ifile-nochange-%d", timestamp)

	// Initial configuration
	initialConfig := fmt.Sprintf(`
resource "bigip_sys_ifile" "testifile" {
  name      = "%s"
  partition = "Common"
  sub_path  = "ravi"
  content   = "initial content"
}
`, ifileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with initial values
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "name", ifileName),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "content", "initial content"),
				),
			},
			// Step 2: No change to the content
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "name", ifileName),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "content", "initial content"),
					// Verify other attributes remain unchanged
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "sub_path", "ravi"),
				),
			},
		},
	})
}

func TestAccBigipSysIfileCreateMinimal(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	ifileName := fmt.Sprintf("testitem-ifile-minimal-%d", timestamp)

	testAccIfileResourceMinimal := fmt.Sprintf(`
resource "bigip_sys_ifile" "testifile" {
  name    = "%s"
  content = "minimal test content"
}
`, ifileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIfileResourceMinimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "name", ifileName),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "partition", "Common"), // Default value
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "sub_path", ""),        // Empty when not specified
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "content", "minimal test content"),
					resource.TestCheckResourceAttrSet("bigip_sys_ifile.testifile", "checksum"), // Computed field
					resource.TestCheckResourceAttrSet("bigip_sys_ifile.testifile", "size"),     // Computed field
				),
			},
		},
	})
}

func TestAccBigipSysIfileImport(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	ifileName := fmt.Sprintf("testitem-ifile-import-%d", timestamp)

	testAccIfileResource := fmt.Sprintf(`
resource "bigip_sys_ifile" "testifile" {
  name      = "%s"
  partition = "Common"
  sub_path  = "ravi"
  content   = "import test content"
}
`, ifileName)

	fullPath := fmt.Sprintf("/Common/ravi/%s", ifileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create the resource
			{
				Config: testAccIfileResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "name", ifileName),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "sub_path", "ravi"),
					resource.TestCheckResourceAttr("bigip_sys_ifile.testifile", "content", "import test content"),
				),
			},
			// Step 2: Import the resource
			{
				ResourceName:            "bigip_sys_ifile.testifile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           fullPath,
				ImportStateVerifyIgnore: []string{"content"}, // Content is sensitive and may not be returned
			},
		},
	})
}

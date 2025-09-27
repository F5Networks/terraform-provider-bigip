package bigip

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigipLtmIfileCreateBasic(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	sysIfileName := fmt.Sprintf("testsys-ifile-%d", timestamp)
	ltmIfileName := fmt.Sprintf("testltm-ifile-%d", timestamp)

	testAccLtmIfileResource := fmt.Sprintf(`
resource "bigip_sys_ifile" "testsysifile" {
  name      = "%s"
  partition = "Common"
  content   = "system ifile content for ltm test"
}

resource "bigip_ltm_ifile" "testltmifile" {
  name      = "%s"
  partition = "Common"
  file_name = bigip_sys_ifile.testsysifile.id
}
`, sysIfileName, ltmIfileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLtmIfileResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "partition", "Common"),
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile", "id"),
					resource.TestCheckResourceAttrSet("bigip_ltm_ifile.testltmifile", "full_path"),
				),
			},
		},
	})
}

func TestAccBigipLtmIfileUpdate(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	sysIfileName1 := fmt.Sprintf("testsys-ifile1-%d", timestamp)
	sysIfileName2 := fmt.Sprintf("testsys-ifile2-%d", timestamp)
	ltmIfileName := fmt.Sprintf("testltm-ifile-update-%d", timestamp)

	// Initial configuration
	initialConfig := fmt.Sprintf(`
resource "bigip_sys_ifile" "testsysifile1" {
  name      = "%s"
  partition = "Common"
  content   = "first system ifile content"
}

resource "bigip_sys_ifile" "testsysifile2" {
  name      = "%s"
  partition = "Common"
  content   = "second system ifile content"
}

resource "bigip_ltm_ifile" "testltmifile" {
  name      = "%s"
  partition = "Common"
  file_name = bigip_sys_ifile.testsysifile1.id
}
`, sysIfileName1, sysIfileName2, ltmIfileName)

	// Updated configuration
	updatedConfig := fmt.Sprintf(`
resource "bigip_sys_ifile" "testsysifile1" {
  name      = "%s"
  partition = "Common"
  content   = "first system ifile content"
}

resource "bigip_sys_ifile" "testsysifile2" {
  name      = "%s"
  partition = "Common"
  content   = "second system ifile content"
}

resource "bigip_ltm_ifile" "testltmifile" {
  name      = "%s"
  partition = "Common"
  file_name = bigip_sys_ifile.testsysifile2.id
}
`, sysIfileName1, sysIfileName2, ltmIfileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with initial file reference
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile1", "id"),
				),
			},
			// Step 2: Update to reference different file
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile2", "id"),
					// Verify other attributes remain unchanged
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "partition", "Common"),
				),
			},
		},
	})
}

func TestAccBigipLtmIfileCreateMinimal(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	sysIfileName := fmt.Sprintf("testsys-ifile-minimal-%d", timestamp)
	ltmIfileName := fmt.Sprintf("testltm-ifile-minimal-%d", timestamp)

	testAccLtmIfileResourceMinimal := fmt.Sprintf(`
resource "bigip_sys_ifile" "testsysifile" {
  name    = "%s"
  content = "minimal system ifile content"
}

resource "bigip_ltm_ifile" "testltmifile" {
  name      = "%s"
  file_name = bigip_sys_ifile.testsysifile.id
}
`, sysIfileName, ltmIfileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLtmIfileResourceMinimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "partition", "Common"), // Default value
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile", "id"),
					resource.TestCheckResourceAttrSet("bigip_ltm_ifile.testltmifile", "full_path"), // Computed field
				),
			},
		},
	})
}

func TestAccBigipLtmIfileCreateWithSubPath(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	sysIfileName := fmt.Sprintf("testsys-ifile-subpath-%d", timestamp)
	ltmIfileName := fmt.Sprintf("testltm-ifile-subpath-%d", timestamp)

	testAccLtmIfileResourceSubPath := fmt.Sprintf(`
resource "bigip_sys_ifile" "testsysifile" {
  name      = "%s"
  partition = "TEST_iFile_300"
  sub_path  = "A1TEST"
  content   = "system ifile with subpath content"
}

resource "bigip_ltm_ifile" "testltmifile" {
  name      = "%s"
  partition = "TEST_iFile_300"
  sub_path  = "A1TEST"
  file_name = bigip_sys_ifile.testsysifile.id
}
`, sysIfileName, ltmIfileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLtmIfileResourceSubPath,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "partition", "TEST_iFile_300"),
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "sub_path", "A1TEST"),
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile", "id"),
					resource.TestCheckResourceAttrSet("bigip_ltm_ifile.testltmifile", "full_path"),
				),
			},
		},
	})
}

func TestAccBigipLtmIfileImport(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	sysIfileName := fmt.Sprintf("testsys-ifile-import-%d", timestamp)
	ltmIfileName := fmt.Sprintf("testltm-ifile-import-%d", timestamp)

	testAccLtmIfileResource := fmt.Sprintf(`
resource "bigip_sys_ifile" "testsysifile" {
  name      = "%s"
  partition = "Common"
  content   = "import test system ifile content"
}

resource "bigip_ltm_ifile" "testltmifile" {
  name      = "%s"
  partition = "Common"
  file_name = bigip_sys_ifile.testsysifile.id
}
`, sysIfileName, ltmIfileName)

	fullPath := fmt.Sprintf("/Common/%s", ltmIfileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create the resource
			{
				Config: testAccLtmIfileResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "partition", "Common"),
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile", "id"),
				),
			},
			// Step 2: Import the resource
			{
				ResourceName:      "bigip_ltm_ifile.testltmifile",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fullPath,
			},
		},
	})
}

func TestAccBigipLtmIfileNoChange(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	sysIfileName := fmt.Sprintf("testsys-ifile-nochange-%d", timestamp)
	ltmIfileName := fmt.Sprintf("testltm-ifile-nochange-%d", timestamp)

	// Configuration that will be applied twice
	config := fmt.Sprintf(`
resource "bigip_sys_ifile" "testsysifile" {
  name      = "%s"
  partition = "Common"
  content   = "no change test content"
}

resource "bigip_ltm_ifile" "testltmifile" {
  name      = "%s"
  partition = "Common"
  file_name = bigip_sys_ifile.testsysifile.id
}
`, sysIfileName, ltmIfileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with initial values
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile", "id"),
				),
			},
			// Step 2: Apply same configuration (no change)
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile", "id"),
					// Verify other attributes remain unchanged
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "partition", "Common"),
				),
			},
		},
	})
}

func TestAccBigipLtmIfileCreateDifferentPartition(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	sysIfileName := fmt.Sprintf("testsys-ifile-partition-%d", timestamp)
	ltmIfileName := fmt.Sprintf("testltm-ifile-partition-%d", timestamp)

	testAccLtmIfileResourcePartition := fmt.Sprintf(`
resource "bigip_sys_ifile" "testsysifile" {
  name      = "%s"
  partition = "TEST_iFile_300"
  content   = "partition test content"
}

resource "bigip_ltm_ifile" "testltmifile" {
  name      = "%s"
  partition = "TEST_iFile_300"
  file_name = bigip_sys_ifile.testsysifile.id
}
`, sysIfileName, ltmIfileName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLtmIfileResourcePartition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "name", ltmIfileName),
					resource.TestCheckResourceAttr("bigip_ltm_ifile.testltmifile", "partition", "TEST_iFile_300"),
					resource.TestCheckResourceAttrPair("bigip_ltm_ifile.testltmifile", "file_name", "bigip_sys_ifile.testsysifile", "id"),
					resource.TestCheckResourceAttrSet("bigip_ltm_ifile.testltmifile", "full_path"),
				),
			},
		},
	})
}

package bigip

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipFastTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipFastCreate,
		Read:   resourceBigipFastRead,
		Update: resourceBigipFastUpdate,
		Delete: resourceBigipFastDelete,
		Exists: resourceBigipFastExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of Fast template set",
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Location of the fast template set package on disk",
			},
			"md5_hash": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "MD5 hash of the fast template zip file",
			},
		},
	}
}

func resourceBigipFastCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	tmplPath := d.Get("source").(string)
	tmplName := filepath.Base(tmplPath)
	checksum := d.Get("md5_hash").(string)
	var name string
	if _, ok := d.GetOk("name"); ok {
		name = d.Get("name").(string)
	} else {
		name = strings.TrimSuffix(tmplName, ".zip")
	}

	log.Println("[INFO] Creating Fast Template Name " + name)
	log.Println("[INFO] Reading provided archive " + tmplName)

	file, fail := os.OpenFile(tmplPath, os.O_RDWR, 0644)
	if fail != nil {
		return fmt.Errorf("error in reading file: %s", fail)
	}

	err := client.UploadFastTemplate(file, name)

	defer file.Close()

	if err != nil {
		return fmt.Errorf("error in creating FAST template set (%s): %s", name, err)
	}
	_ = d.Set("md5_hash", checksum)
	d.SetId(name)
	if !client.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			Name:    "Terraform-provider-bigip",
			Version: client.UserAgent,
			Id:      uniqueID,
		}
		apiKey := os.Getenv("TEEM_API_KEY")
		teemDevice := f5teem.AnonymousClient(assetInfo, apiKey)
		f := map[string]interface{}{
			"Terraform Version": client.UserAgent,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_fast_template", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipFastRead(d, meta)
}

func resourceBigipFastRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	checksum := d.Get("md5_hash").(string)
	log.Println("[INFO] Reading Fast Template Set : " + name)

	template, err := client.GetTemplateSet(name)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Fast Template Set content: %+v", template)
	_ = d.Set("name", template.Name)
	_ = d.Set("md5_hash", checksum)

	return nil
}

func resourceBigipFastExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Checking Template Set " + name + " exists.")

	template, err := client.GetTemplateSet(name)

	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Fast template set (%s) (%v) ", name, err)
		return false, err
	}

	if template == nil {
		log.Printf("[WARN] Fast template set (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return template != nil, nil
}

func resourceBigipFastUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceBigipFastCreate(d, meta)
}

func resourceBigipFastDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting Fast Template Set " + name)
	err := client.DeleteTemplateSet(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Fast Template Set   (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

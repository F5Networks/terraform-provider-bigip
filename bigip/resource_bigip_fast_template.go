package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"os"
	"path/filepath"
	"strings"
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
				Description: "Name of Fast template set",
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Location of the fast template set package on disk",
			},
		},
	}
}

func resourceBigipFastCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	tmplPath := d.Get("source").(string)
	tmplName := filepath.Base(tmplPath)
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
	d.SetId(name)
	return resourceBigipFastRead(d, meta)
}

func resourceBigipFastRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading Fast Template Set : " + name)

	template, err := client.GetTemplateSet(name)

	log.Printf("[INFO] Fast Template Set content: %+v", template)

	d.Set("name", template.Name)

	if err != nil {
		return err
	}
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

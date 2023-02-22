package bigip

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipFastTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipFastCreate,
		ReadContext:   resourceBigipFastRead,
		UpdateContext: resourceBigipFastUpdate,
		DeleteContext: resourceBigipFastDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceBigipFastCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(fmt.Errorf("error in reading file: %s", fail))
	}

	err := client.UploadFastTemplate(file, name)

	defer file.Close()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error in creating FAST template set (%s): %s", name, err))
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
	return resourceBigipFastRead(ctx, d, meta)
}

func resourceBigipFastRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	checksum := d.Get("md5_hash").(string)
	log.Println("[INFO] Reading Fast Template Set : " + name)

	template, err := client.GetTemplateSet(name)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Fast Template Set content: %+v", template)
	_ = d.Set("name", template.Name)
	_ = d.Set("md5_hash", checksum)

	return nil
}

func resourceBigipFastUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceBigipFastCreate(ctx, d, meta)
}

func resourceBigipFastDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting Fast Template Set " + name)
	err := client.DeleteTemplateSet(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Fast Template Set   (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

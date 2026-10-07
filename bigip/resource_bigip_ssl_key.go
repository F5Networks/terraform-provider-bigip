package bigip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSslKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSslKeyCreate,
		ReadContext:   resourceBigipSslKeyRead,
		UpdateContext: resourceBigipSslKeyUpdate,
		DeleteContext: resourceBigipSslKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of SSL Certificate key with .key extension",
				ForceNew:    true,
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"content_wo"},
				Description:   "Content of SSL certificate key present on local Disk. Cannot be used with content_wo.",
			},
			"content_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				WriteOnly:     true,
				ConflictsWith: []string{"content"},
				Description:   "Write-only content of SSL certificate key. Not persisted to state. Cannot be used with content.",
			},
			"content_wo_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				WriteOnly:   true,
				Description: "Version number to trigger content_wo re-upload. Increment this value to force key re-upload.",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase on key.",
			},
			"partition": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Common",
				Description:  "Partition of ssl certificate key",
				ValidateFunc: validatePartitionName,
			},
			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Full Path Name of ssl key",
			},
		},
	}
}

func getWriteOnlyString(d *schema.ResourceData, attribute string) (string, bool, diag.Diagnostics) {
	value, diags := d.GetRawConfigAt(cty.GetAttrPath(attribute))
	if diags.HasError() {
		return "", false, diags
	}
	if value.IsNull() {
		return "", false, nil
	}
	if !value.IsKnown() {
		return "", false, diag.Errorf("%q must be known during apply", attribute)
	}
	return value.AsString(), true, nil
}

func resourceBigipSslKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Certificate Key Name " + name)

	// Handle both content and content_wo, ensuring at least one is provided
	var certpath string
	if contentVal, ok := d.GetOk("content"); ok {
		certpath = contentVal.(string)
	} else {
		contentWoVal, ok, diags := getWriteOnlyString(d, "content_wo")
		if diags.HasError() {
			return diags
		}
		if !ok {
			return diag.Errorf("either 'content' or 'content_wo' must be specified")
		}
		certpath = contentWoVal
	}

	partition := d.Get("partition").(string)
	passPhrase := d.Get("passphrase").(string)

	sourcePath, err := client.UploadKey(name, certpath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error in Uploading certificate key (%s): %s", name, err))
	}
	certkey := bigip.Key{
		Name:       name,
		SourcePath: sourcePath,
		Partition:  partition,
		Passphrase: passPhrase,
	}
	log.Printf("[DEBUG] certkey: %+v\n", certkey)
	err = client.AddKey(&certkey)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipSslKeyRead(ctx, d, meta)
}

func resourceBigipSslKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading Certificate key: " + name)
	/*if !strings.HasSuffix(name, ".key") {
		name = name + ".key"
	}*/
	partition := d.Get("partition").(string)
	if partition == "" {
		if !strings.HasPrefix(name, "/") {
			err := errors.New("the name must be in full_path format when partition is not specified")
			fmt.Print(err)
		}
	} else {
		if !strings.HasPrefix(name, "/") {
			name = "/" + partition + "/" + name
		}
	}
	certkey, err := client.GetKey(name)
	if err != nil && strings.Contains(err.Error(), "not found") {
		log.Printf("[WARN] SSL Key (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if certkey == nil {
		return diag.Errorf("Reading Certificate key failed with key:%v", certkey)
	}
	log.Printf("[INFO] SSL key content:%+v", certkey)
	_ = d.Set("name", certkey.Name)
	_ = d.Set("partition", certkey.Partition)
	_ = d.Set("full_path", certkey.FullPath)
	// Note: content and content_wo are write-only and not set here - they are never read from BIG-IP
	return nil
}

func resourceBigipSslKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Certificate key Name " + name)

	partition := d.Get("partition").(string)
	passPhrase := d.Get("passphrase").(string)

	// Check if content or content_wo has changed, or if content_wo_version has changed
	hasContentChange := d.HasChange("content")
	hasContentWoChange := d.HasChange("content_wo")
	hasVersionChange := d.HasChange("content_wo_version")

	// Only update if content, content_wo, or content_wo_version changed
	if hasContentChange || hasContentWoChange || hasVersionChange {
		var certpath string
		if contentVal, ok := d.GetOk("content"); ok {
			certpath = contentVal.(string)
		} else {
			contentWoVal, ok, diags := getWriteOnlyString(d, "content_wo")
			if diags.HasError() {
				return diags
			}
			if !ok {
				return diag.Errorf("either 'content' or 'content_wo' must be specified")
			}
			certpath = contentWoVal
		}

		sourcePath, err := client.UploadKey(name, certpath)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error in Uploading certificate key (%s): %s", name, err))
		}
		certkey := bigip.Key{
			Name:       name,
			SourcePath: sourcePath,
			Partition:  partition,
			Passphrase: passPhrase,
		}
		keyName := fmt.Sprintf("/%s/%s", partition, name)
		err = client.ModifyKey(keyName, &certkey)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if d.HasChange("passphrase") {
		// If only passphrase changed, update the key without re-uploading content
		keyName := fmt.Sprintf("/%s/%s", partition, name)
		certkey := bigip.Key{
			Name:       name,
			Partition:  partition,
			Passphrase: passPhrase,
		}
		err := client.ModifyKey(keyName, &certkey)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceBigipSslKeyRead(ctx, d, meta)
}

func resourceBigipSslKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting Certificate key" + name)
	/*if !strings.HasSuffix(name, ".key") {
		name = name + ".key"
	}*/
	partition := d.Get("partition").(string)
	name = "/" + partition + "/" + name
	err := client.DeleteKey(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Pool   (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

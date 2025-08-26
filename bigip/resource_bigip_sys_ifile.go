package bigip

import (
	"context"
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSysIfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysIfileCreate,
		ReadContext:   resourceBigipSysIfileRead,
		UpdateContext: resourceBigipSysIfileUpdate,
		DeleteContext: resourceBigipSysIfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the iFile.",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				Description: "Partition for the iFile.",
			},
			"sub_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subdirectory within the partition.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The content of the iFile.",
			},
			"checksum": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Checksum of the iFile content.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the iFile in bytes.",
			},
		},
	}
}

func resourceBigipSysIfileCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	subPath := d.Get("sub_path").(string)
	content := d.Get("content").(string)
	fileData := &bigip.IFile{
		Name:      name,
		Partition: partition,
		SubPath:   subPath,
	}
	fullPath := buildIFileFullPath(partition, subPath, name)
	err := client.ImportIfile(fileData, content, "POST")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error importing iFile: %v", err))
	}
	d.SetId(fullPath)
	return resourceBigipSysIfileRead(ctx, d, meta)
}

func resourceBigipSysIfileRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	fullPath := d.Id()
	ifile, err := client.GetIFile(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading iFile: %v", err))
	}
	if ifile == nil {
		return diag.Errorf("iFile not found: %s", fullPath)
	}
	log.Printf("[resourceBigipSysIfileRead][INFO] Fullpath iFile: %s", ifile.FullPath)
	_ = d.Set("name", ifile.Name)
	_ = d.Set("partition", ifile.Partition)
	_ = d.Set("sub_path", ifile.SubPath)
	// _ = d.Set("full_path", ifile.FullPath)
	_ = d.Set("checksum", ifile.Checksum)
	_ = d.Set("size", ifile.Size)
	return nil
}

func resourceBigipSysIfileUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	fullPath := d.Id()
	parts := strings.Split(fullPath, "/")
	name := parts[len(parts)-1]
	// name := d.Get("name").(string)
	content := d.Get("content").(string)
	// partition := d.Get("partition").(string)
	// subPath := d.Get("sub_path").(string)
	// fullPath := buildIFileFullPath(partition, subPath, name)
	fileData := &bigip.IFile{
		Name:     name,
		FullPath: fullPath,
		// Partition: partition,
		// SubPath:   subPath,
	}
	err := client.ImportIfile(fileData, content, "PUT")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error importing iFile: %v", err))
	}
	return resourceBigipSysIfileRead(ctx, d, meta)
}

func resourceBigipSysIfileDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	fullPath := d.Id()
	log.Printf("[INFO] Deleting iFile: %s", fullPath)
	err := client.DeleteIFile(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting iFile: %v", err))
	}
	d.SetId("")
	return nil
}

func buildIFileFullPath(partition, subPath, name string) string {
	path := "/" + partition
	if subPath != "" {
		path += "/" + subPath
	}
	path += "/" + name
	return path
}

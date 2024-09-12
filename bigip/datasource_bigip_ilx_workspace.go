package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigIPILXWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigIPILXWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the ILX Workspace",
			},
			"full_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full path of the ILX Workspace",
			},
			"generation": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The generation of the ILX Workspace",
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The self link of the ILX Workspace",
			},
			"node_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The node version of the ILX Workspace",
			},
			"staged_directory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The staged directory of the ILX Workspace",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the ILX Workspace",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The directory of the ILX Workspace",
				Elem:        ruleDataSchema(),
			},
			"extensions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The directory of the ILX Workspace",
				Elem:        extensionDataSchema(),
			},
		},
	}
}

func dataSourceBigIPILXWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	log.Printf("[INFO] Retrieving ILX Workspace %s", d.Get("name").(string))
	spc, err := client.GetWorkspace(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	diag := setILXWorkspaceResourceData(d, spc)
	if diag != nil {
		d.SetId("")
		return diag
	}
	log.Println("[INFO] Retrieved ILX Workspace")
	d.SetId(spc.FullPath)
	return nil
}

func setILXWorkspaceResourceData(d *schema.ResourceData, spc *bigip.ILXWorkspace) diag.Diagnostics {
	err := d.Set("name", spc.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("full_path", spc.FullPath)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("generation", spc.Generation)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("self_link", spc.SelfLink)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("node_version", spc.NodeVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("staged_directory", spc.StagedDirectory)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("version", spc.Version)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("extensions", flattenExtensions(spc.Extensions))
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("rules", flattenFiles(spc.Rules))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func extensionDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"files": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     fileSchema(),
			},
		},
	}
}

func ruleDataSchema() *schema.Resource {
	return fileSchema()
}

func fileSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func flattenExtensions(extensions []bigip.Extension) []map[string]any {
	if extensions == nil {
		return []map[string]any{}
	}

	var result []map[string]any
	for _, ext := range extensions {
		extMap := map[string]any{
			"name":  ext.Name,
			"files": flattenFiles(ext.Files),
		}
		result = append(result, extMap)
	}
	return result
}

func flattenFiles(files []bigip.File) []map[string]string {
	if files == nil {
		return []map[string]string{}
	}

	var result []map[string]string
	for _, file := range files {
		fileMap := map[string]string{
			"name": file.Name,
		}
		result = append(result, fileMap)
	}
	return result
}

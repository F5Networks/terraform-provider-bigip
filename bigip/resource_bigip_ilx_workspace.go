package bigip

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"os"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigIPILXWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceBigIPILXWorkspaceRead,
		CreateContext: resourceBigIPILXWorkspaceCreate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extension": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     extensionSchema(),
			},
		},
	}
}

func extensionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"extension_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"partition": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Common",
			},
			"index_source": {
				Type:          schema.TypeString,
				Optional:      true,
				ExactlyOneOf:  []string{"index_source", "index_source_dir"},
				ConflictsWith: []string{"index_source_dir"},
			},
			"package_source": {
				Type:          schema.TypeString,
				Optional:      true,
				ExactlyOneOf:  []string{"package_source", "package_source_dir"},
				ConflictsWith: []string{"package_source_dir"},
			},
			"index_source_dir": {
				Type:          schema.TypeString,
				Optional:      true,
				ExactlyOneOf:  []string{"index_source", "index_source_dir"},
				ConflictsWith: []string{"index_source"},
			},
			"package_source_dir": {
				Type:          schema.TypeString,
				Optional:      true,
				ExactlyOneOf:  []string{"package_source", "package_source_dir"},
				ConflictsWith: []string{"package_source"},
			},
			"force_upload": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"file_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The computed hash of the uploaded files in the workspace. Used to determine if the files have changed and need to be uploaded",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					iContent, _ := getFileContent(d.Get("index_source_path").(string), d.Get("index_source").(string), "", "")
					pContent, _ := getFileContent(d.Get("package_source_path").(string), d.Get("package_source").(string), "", "")
					currentHash := computeHash(iContent, pContent)
					return old == currentHash
				},
			},
		},
	}
}

func resourceBigIPILXWorkspaceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	name := d.Get("name").(string)
	workspace, err := client.GetWorkspace(ctx, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading workspace %w", err))
	}

	if workspace == nil {
		log.Printf("[DEBUG] workspace (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}

	d.SetId(workspace.Name)
	return nil
}

func resourceBigIPILXWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	err := client.CreateWorkspace(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating workspace %w", err))
	}
	if ext_block := d.Get("extension").(*schema.ResourceData); ext_block != nil {
		handleExtensionCreate(ctx, ext_block, client)
	}
	d.SetId(d.Get("name").(string))
	return nil
}

func handleExtensionCreate(ctx context.Context, ext *schema.ResourceData, client *bigip.BigIP) {
	opts := bigip.ExtensionConfig{
		WorkspaceName: ext.Get("name").(string),
		Name:          ext.Get("extension_name").(string),
		Partition:     ext.Get("partition").(string),
	}
	client.CreateExtension(ctx, opts)
	var iContent, pContent string
	var err error
	if iContent, err = handleSourcePath(ctx, ext, client, opts, "index_source_path", "index_source", client.WriteExtensionFile, bigip.IndexJS); err != nil {
		log.Fatalf("Error handling index source: %v", err)
	}
	if pContent, err = handleSourcePath(ctx, ext, client, opts, "package_source_path", "package_source", client.WriteExtensionFile, bigip.PackageJSON); err != nil {
		log.Fatalf("Error handling package source: %v", err)
	}
	hash := computeHash(iContent, pContent)
	ext.Set("file_hash", hash)
	compositeId := fmt.Sprintf("%s:%s:%s", opts.WorkspaceName, opts.Name, opts.Partition)
	ext.SetId(compositeId)
}

func handleSourcePath(ctx context.Context, ext *schema.ResourceData, client *bigip.BigIP, opts bigip.ExtensionConfig, pathKey, sourceKey string, writeFunc func(context.Context, bigip.ExtensionConfig, string, bigip.ExtensionFile) error, filename bigip.ExtensionFile) (string, error) {
	var content string
	var err error

	if path := ext.Get(pathKey).(string); path != "" {
		b, err := os.ReadFile(path)
		content = string(b)
		if err != nil {
			return content, fmt.Errorf("Error reading %s: %v", pathKey, err)
		}
	} else if source := ext.Get(sourceKey).(string); source != "" {
		content = source
	}

	if content != "" {
		err = writeFunc(ctx, opts, content, filename)
		if err != nil {
			return content, fmt.Errorf("Error writing %s: %v", pathKey, err)
		}
	}

	return content, nil
}

func getFileContent(ipath, icontent, ppath, pcontent string) (string, error) {
	var iContent, pContent string
	if ipath != "" {
		b, err := os.ReadFile(ipath)
		iContent = string(b)
		if err != nil {
			return "", fmt.Errorf("Error reading index source: %v", err)
		}
	} else if icontent != "" {
		iContent = icontent
	}
	if ppath != "" {
		b, err := os.ReadFile(ppath)
		pContent = string(b)
		if err != nil {
			return "", fmt.Errorf("Error reading package source: %v", err)
		} else if pcontent != "" {
			pContent = pcontent
		}
	}
	return computeHash(iContent, pContent), nil
}

func computeHash(iContent, pContent string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(iContent+pContent)))
}

func resourceBigIPILXWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	filehash := d.Get("file_hash").(string)

	return nil
}

func resourceBigIPILXWorkspaceDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	err := client.DeleteWorkspace(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting workspace %w", err))
	}
	d.SetId("")
	return nil
}

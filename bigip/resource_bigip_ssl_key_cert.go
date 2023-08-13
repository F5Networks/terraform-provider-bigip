package bigip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSSLKeyCert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSSLKeyCertCreate,
		ReadContext:   resourceBigipSSLKeyCertRead,
		UpdateContext: resourceBigipSSLKeyCertUpdate,
		DeleteContext: resourceBigipSSLKeyCertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the key.",
			},
			"key_content": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The content of the key.",
			},
			"key_full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Full Path Name of ssl key",
			},
			"cert_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cert.",
				ForceNew:    true,
			},
			"cert_content": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The content of the cert.",
			},
			"cert_full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Full Path Name of ssl certificate",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase on the key.",
			},
			"partition": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Common",
				Description:  "Partition on the ssl certificate and key.",
				ValidateFunc: validatePartitionName,
			},
		},
	}
}

func resourceBigipSSLKeyCertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	keyName := d.Get("key_name").(string)
	keyPath := d.Get("key_content").(string)
	partition := d.Get("partition").(string)
	passphrase := d.Get("passphrase").(string)
	certName := d.Get("cert_name").(string)
	certPath := d.Get("cert_content").(string)

	sourcePath, err := client.UploadKey(keyName, keyPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while uploading the ssl key: %v", err))
	}

	keyCfg := bigip.Key{
		Name:       keyName,
		SourcePath: sourcePath,
		Partition:  partition,
		Passphrase: passphrase,
	}

	err = client.AddKey(&keyCfg)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while adding the ssl key: %v", err))
	}
	err = client.UploadCertificate(certName, certPath, partition)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while uploading the ssl cert: %v", err))
	}

	id := keyName + "_" + certName
	d.SetId(id)
	return resourceBigipSSLKeyCertRead(ctx, d, meta)
}

func resourceBigipSSLKeyCertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	partition := d.Get("partition").(string)

	keyName := fqdn(partition, d.Get("key_name").(string))
	certName := fqdn(partition, d.Get("cert_name").(string))

	key, err := client.GetKey(keyName)
	if err != nil {
		diag.FromErr(err)
	}
	if key == nil {
		return diag.Errorf("reading ssl key failed with key: %v", key)
	}

	certificate, err := client.GetCertificate(certName)
	if err != nil {
		return diag.FromErr(err)
	}
	if certificate == nil {
		return diag.Errorf("reading certificate failed  :%+v", certificate)
	}

	d.Set("key_name", key.Name)
	d.Set("key_full_path", key.FullPath)
	d.Set("cert_name", certificate.Name)
	d.Set("cert_full_path", certificate.FullPath)
	d.Set("partition", key.Partition)

	return nil
}

func resourceBigipSSLKeyCertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	keyName := d.Get("key_name").(string)
	keyPath := d.Get("key_content").(string)
	partition := d.Get("partition").(string)
	passphrase := d.Get("passphrase").(string)
	certName := d.Get("cert_name").(string)
	certPath := d.Get("cert_content").(string)

	sourcePath, err := client.UploadKey(keyName, keyPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while trying to upload ssl key (%s): %s", keyName, err))
	}

	keyCfg := bigip.Key{
		Name:       keyName,
		SourcePath: sourcePath,
		Partition:  partition,
		Passphrase: passphrase,
	}

	keyFullPath := fmt.Sprintf("/%s/%s", partition, keyName)
	err = client.ModifyKey(keyFullPath, &keyCfg)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while trying to modify the ssl key (%s): %s", keyFullPath, err))
	}

	err = client.UpdateCertificate(certName, certPath, partition)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while updating the ssl certificate (%s): %s", certName, err))
	}

	return resourceBigipSSLKeyCertRead(ctx, d, meta)
}

func resourceBigipSSLKeyCertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Deleteing SSL Key and Certificate")
	keyName := d.Get("key_name").(string)
	partition := d.Get("partition").(string)
	certName := d.Get("cert_name").(string)

	log.Printf("[INFO] Deleteing SSL Key %s and Certificate %s", keyName, certName)

	keyFullPath := "/" + partition + "/" + keyName
	certFullPath := "/" + partition + "/" + certName

	err := client.DeleteKey(keyFullPath)
	if err != nil {
		log.Printf("[ERROR] unable to delete the ssl key (%s) (%v) ", keyFullPath, err)
	}

	err = client.DeleteCertificate(certFullPath)
	if err != nil {
		log.Printf("[ERROR] unable to delete the ssl certificate (%s) (%v) ", certFullPath, err)
	}

	d.SetId("")
	return nil
}

func fqdn(partition, name string) string {
	if partition == "" {
		if !strings.HasPrefix(name, "/") {
			err := errors.New("the key name must be in full_path format when partition is not specified")
			fmt.Print(err)
		}
	} else {
		if !strings.HasPrefix(name, "/") {
			name = "/" + partition + "/" + name
		}
	}

	return name
}

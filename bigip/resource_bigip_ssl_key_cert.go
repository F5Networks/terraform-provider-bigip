package bigip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
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
			"cert_monitoring_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type of monitoring used.",
			},
			"issuer_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the issuer certificate",
			},
			"cert_ocsp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the OCSP responder",
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

	t, err := client.StartTransaction()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while starting transaction: %v", err))
	}
	err = client.AddKey(&keyCfg)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while adding the ssl key: %v", err))
	}

	cert := &bigip.Certificate{
		Name:      certName,
		Partition: partition,
	}
	if val, ok := d.GetOk("cert_monitoring_type"); ok {
		cert.CertValidationOptions = []string{val.(string)}
	}
	if val, ok := d.GetOk("issuer_cert"); ok {
		cert.IssuerCert = val.(string)
	}

	err = client.UploadCertificate(certPath, cert)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while uploading the ssl cert: %v", err))
	}
	err = client.CommitTransaction(t.TransID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while ending transaction: %d", err))
	}

	if val, ok := d.GetOk("cert_ocsp"); ok {
		certValidState := &bigip.CertValidatorState{Name: val.(string)}
		certValidRef := &bigip.CertValidatorReference{}
		certValidRef.Items = append(certValidRef.Items, *certValidState)
		cert.CertValidatorRef = certValidRef
		err = client.UpdateCertificate(certPath, cert)
		if err != nil {
			log.Printf("[ERROR]Unable to add ocsp to the certificate:%v", err)
		}
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
	d.Set("issuer_cert", certificate.IssuerCert)
	if certificate.CertValidationOptions != nil && len(certificate.CertValidationOptions) > 0 {
		monitor_type := certificate.CertValidationOptions[0]
		_ = d.Set("cert_monitoring_type", monitor_type)
	}

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

	t, err := client.StartTransaction()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while trying to start transaction: %s", err))
	}
	err = client.ModifyKey(keyFullPath, &keyCfg)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while trying to modify the ssl key (%s): %s", keyFullPath, err))
	}

	cert := &bigip.Certificate{
		Name:      certName,
		Partition: partition,
	}
	if val, ok := d.GetOk("cert_monitoring_type"); ok {
		cert.CertValidationOptions = []string{val.(string)}
	}
	if val, ok := d.GetOk("issuer_cert"); ok {
		cert.IssuerCert = val.(string)
	}
	if val, ok := d.GetOk("cert_ocsp"); ok {
		certValidState := &bigip.CertValidatorState{Name: val.(string)}
		certValidRef := &bigip.CertValidatorReference{}
		certValidRef.Items = append(certValidRef.Items, *certValidState)
		cert.CertValidatorRef = certValidRef
	}

	err = client.UpdateCertificate(certPath, cert)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while updating the ssl certificate (%s): %s", certName, err))
	}
	err = client.CommitTransaction(t.TransID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while trying to end transaction: %s", err))
	}

	return resourceBigipSSLKeyCertRead(ctx, d, meta)
}

func resourceBigipSSLKeyCertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Deleting SSL Key and Certificate")
	keyName := d.Get("key_name").(string)
	partition := d.Get("partition").(string)
	certName := d.Get("cert_name").(string)

	log.Printf("[INFO] Deleting SSL Key %s and Certificate %s", keyName, certName)

	keyFullPath := "/" + partition + "/" + keyName
	certFullPath := "/" + partition + "/" + certName

	t, err := client.StartTransaction()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while starting transaction: %v", err))
	}

	err = client.DeleteKey(keyFullPath)
	if err != nil {
		log.Printf("[ERROR] unable to delete the ssl key (%s) (%v) ", keyFullPath, err)
	}

	err = client.DeleteCertificate(certFullPath)
	if err != nil {
		log.Printf("[ERROR] unable to delete the ssl certificate (%s) (%v) ", certFullPath, err)
	}

	err = client.CommitTransaction(t.TransID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while ending transaction: %v", err))
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

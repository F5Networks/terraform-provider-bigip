package bigip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
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
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_content_wo"},
				Description:   "The content of the key. Cannot be used with key_content_wo.",
			},
			"key_content_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				WriteOnly:     true,
				ConflictsWith: []string{"key_content"},
				Description:   "Write-only content of the key. Not persisted to state. Cannot be used with key_content.",
			},
			"key_content_wo_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				WriteOnly:   true,
				Description: "Version number to trigger key_content_wo re-upload. Increment to force key re-upload.",
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
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"cert_content_wo"},
				Description:   "The content of the cert. Cannot be used with cert_content_wo.",
			},
			"cert_content_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				WriteOnly:     true,
				ConflictsWith: []string{"cert_content"},
				Description:   "Write-only content of the cert. Not persisted to state. Cannot be used with cert_content.",
			},
			"cert_content_wo_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				WriteOnly:   true,
				Description: "Version number to trigger cert_content_wo re-upload. Increment to force cert re-upload.",
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
	partition := d.Get("partition").(string)
	passphrase := d.Get("passphrase").(string)
	certName := d.Get("cert_name").(string)

	// Handle both key_content and key_content_wo, ensuring at least one is provided
	var keyPath string
	if keyContentVal, ok := d.GetOk("key_content"); ok {
		keyPath = keyContentVal.(string)
	} else {
		keyContentWoVal, ok, diags := getWriteOnlyString(d, "key_content_wo")
		if diags.HasError() {
			return diags
		}
		if !ok {
			return diag.Errorf("either 'key_content' or 'key_content_wo' must be specified")
		}
		keyPath = keyContentWoVal
	}

	// Handle both cert_content and cert_content_wo, ensuring at least one is provided
	var certPath string
	if certContentVal, ok := d.GetOk("cert_content"); ok {
		certPath = certContentVal.(string)
	} else {
		certContentWoVal, ok, diags := getWriteOnlyString(d, "cert_content_wo")
		if diags.HasError() {
			return diags
		}
		if !ok {
			return diag.Errorf("either 'cert_content' or 'cert_content_wo' must be specified")
		}
		certPath = certContentWoVal
	}

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

	cert := &bigip.Certificate{
		Name:      certName,
		Partition: partition,
	}
	if val, ok := d.GetOk("cert_monitoring_type"); ok {
		cert.CertValidationOptions = []string{val.(string)}
	}

	if diags := func() diag.Diagnostics {
		mutex.Lock()
		defer mutex.Unlock()

		t, err := client.StartTransaction()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error while starting transaction: %v", err))
		}
		err = client.AddKey(&keyCfg)
		if err != nil {
			rollbackTransaction(client, t.TransID)
			return diag.FromErr(fmt.Errorf("error while adding the ssl key: %v", err))
		}

		err = client.UploadCertificate(certPath, cert)
		if err != nil {
			rollbackTransaction(client, t.TransID)
			return diag.FromErr(fmt.Errorf("error while uploading the ssl cert: %v", err))
		}
		err = client.CommitTransaction(t.TransID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error while ending transaction: %d", err))
		}
		return nil
	}(); diags != nil {
		return diags
	}

	// issuer_cert is silently ignored by BIG-IP when set on the initial
	// upload/create above (inside the transaction) -- it only takes
	// effect via a follow-up PATCH/PUT against the already-created
	// certificate, so it's applied here via ModifyCertificate instead.
	// ModifyCertificate requires the full "/partition/name" path (it does
	// not itself disambiguate a bare name against the partition field the
	// way AddCertificate's request body does).
	if val, ok := d.GetOk("issuer_cert"); ok {
		cert.IssuerCert = val.(string)
		fullPathCertName := fqdn(partition, certName)
		if err := client.ModifyCertificate(fullPathCertName, &bigip.Certificate{IssuerCert: cert.IssuerCert}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting issuer_cert on certificate (%s): %s", fullPathCertName, err))
		}
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

	_ = d.Set("key_name", key.Name)
	_ = d.Set("key_full_path", key.FullPath)
	_ = d.Set("cert_name", certificate.Name)
	_ = d.Set("cert_full_path", certificate.FullPath)
	_ = d.Set("partition", key.Partition)
	_ = d.Set("issuer_cert", certificate.IssuerCert)
	if len(certificate.CertValidationOptions) > 0 {
		monitor_type := certificate.CertValidationOptions[0]
		_ = d.Set("cert_monitoring_type", monitor_type)
	}
	// Note: key_content, key_content_wo, cert_content, cert_content_wo are write-only and not set here - they are never read from BIG-IP

	return nil
}

func resourceBigipSSLKeyCertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	keyName := d.Get("key_name").(string)
	partition := d.Get("partition").(string)
	passphrase := d.Get("passphrase").(string)
	certName := d.Get("cert_name").(string)

	keyFullPath := fmt.Sprintf("/%s/%s", partition, keyName)

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

	// Check if key content or version has changed
	hasKeyContentChange := d.HasChange("key_content")
	hasKeyContentWoChange := d.HasChange("key_content_wo")
	hasKeyVersionChange := d.HasChange("key_content_wo_version")

	// Check if cert content or version has changed
	hasCertContentChange := d.HasChange("cert_content")
	hasCertContentWoChange := d.HasChange("cert_content_wo")
	hasCertVersionChange := d.HasChange("cert_content_wo_version")

	if diags := func() diag.Diagnostics {
		mutex.Lock()
		defer mutex.Unlock()

		t, err := client.StartTransaction()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error while trying to start transaction: %s", err))
		}

		// Update key if needed.
		if hasKeyContentChange || hasKeyContentWoChange || hasKeyVersionChange {
			var keyPath string
			if keyContentVal, ok := d.GetOk("key_content"); ok {
				keyPath = keyContentVal.(string)
			} else {
				keyContentWoVal, ok, diags := getWriteOnlyString(d, "key_content_wo")
				if diags.HasError() {
					rollbackTransaction(client, t.TransID)
					return diags
				}
				if !ok {
					rollbackTransaction(client, t.TransID)
					return diag.Errorf("either 'key_content' or 'key_content_wo' must be specified")
				}
				keyPath = keyContentWoVal
			}

			sourcePath, err := client.UploadKey(keyName, keyPath)
			if err != nil {
				rollbackTransaction(client, t.TransID)
				return diag.FromErr(fmt.Errorf("error while trying to upload ssl key (%s): %s", keyName, err))
			}

			keyCfg := bigip.Key{
				Name:       keyName,
				SourcePath: sourcePath,
				Partition:  partition,
				Passphrase: passphrase,
			}
			if err := client.ModifyKey(keyFullPath, &keyCfg); err != nil {
				rollbackTransaction(client, t.TransID)
				return diag.FromErr(fmt.Errorf("error while trying to modify the ssl key (%s): %s", keyFullPath, err))
			}
		} else if d.HasChange("passphrase") {
			keyCfg := bigip.Key{Name: keyName, Partition: partition, Passphrase: passphrase}
			if err := client.ModifyKey(keyFullPath, &keyCfg); err != nil {
				rollbackTransaction(client, t.TransID)
				return diag.FromErr(fmt.Errorf("error while trying to modify the ssl key (%s): %s", keyFullPath, err))
			}
		}

		// Update certificate if needed.
		if hasCertContentChange || hasCertContentWoChange || hasCertVersionChange {
			var certPath string
			if certContentVal, ok := d.GetOk("cert_content"); ok {
				certPath = certContentVal.(string)
			} else {
				certContentWoVal, ok, diags := getWriteOnlyString(d, "cert_content_wo")
				if diags.HasError() {
					rollbackTransaction(client, t.TransID)
					return diags
				}
				if !ok {
					rollbackTransaction(client, t.TransID)
					return diag.Errorf("either 'cert_content' or 'cert_content_wo' must be specified")
				}
				certPath = certContentWoVal
			}
			if err := client.UpdateCertificate(certPath, cert); err != nil {
				rollbackTransaction(client, t.TransID)
				return diag.FromErr(fmt.Errorf("error while updating the ssl certificate (%s): %s", certName, err))
			}
		}

		if val, ok := d.GetOk("cert_ocsp"); ok {
			cert.CertValidatorRef = &bigip.CertValidatorReference{Items: []bigip.CertValidatorState{{Name: val.(string)}}}
		}

		if err := client.CommitTransaction(t.TransID); err != nil {
			return diag.FromErr(fmt.Errorf("error while trying to end transaction: %s", err))
		}
		return nil
	}(); diags != nil {
		return diags
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

// rollbackTransaction attempts to roll back an iControl REST transaction
// that was started with client.StartTransaction but should not be committed
// (e.g. because AddKey/ModifyKey/UploadCertificate failed partway through).
// go-bigip has no exported RollbackTransaction, only CommitTransaction, and
// its PATCH helper (patch) and URL builder (iControlPath) are unexported, so
// this issues the equivalent request directly via the exported APICall,
// mirroring exactly what CommitTransaction does but with
// {"state":"ROLLED_BACK"} instead of {"state":"VALIDATING"}. Errors are
// logged rather than returned: the caller is already in an error path and
// returning a second, unrelated error here would obscure the original
// failure. If the rollback request itself fails, the transaction is left to
// expire on its own via BIG-IP's transaction timeout instead of remaining
// open indefinitely.
func rollbackTransaction(client *bigip.BigIP, transID int64) {
	body, err := json.Marshal(map[string]interface{}{"state": "ROLLED_BACK"})
	if err != nil {
		log.Printf("[ERROR] unable to encode transaction rollback request for transaction %d: %v", transID, err)
		return
	}
	req := &bigip.APIRequest{
		Method:      "patch",
		URL:         "mgmt/tm/transaction/" + strconv.FormatInt(transID, 10),
		Body:        string(body),
		ContentType: "application/json",
	}
	if _, err := client.APICall(req); err != nil {
		log.Printf("[ERROR] unable to roll back transaction %d: %v", transID, err)
	}
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

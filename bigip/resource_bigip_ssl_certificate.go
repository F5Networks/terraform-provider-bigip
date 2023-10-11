package bigip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSslCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSslCertificateCreate,
		ReadContext:   resourceBigipSslCertificateRead,
		UpdateContext: resourceBigipSslCertificateUpdate,
		DeleteContext: resourceBigipSslCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of SSL Certificate with .crt extension",
				ForceNew:    true,
			},
			"content": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				//ForceNew:    true,
				Description: "Content of certificate on Disk",
			},
			"partition": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Common",
				Description:  "Partition of ssl certificate",
				ValidateFunc: validatePartitionName,
			},
			"monitoring_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type of monitoring used",
			},
			"issuer_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the issuer certificate",
			},
			"ocsp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the OCSP responder",
			},
			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Full Path Name of ssl certificate",
			},
		},
	}
}

func resourceBigipSslCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Certificate Name " + name)

	certPath := d.Get("content").(string)
	partition := d.Get("partition").(string)
	cert := &bigip.Certificate{
		Name:      name,
		Partition: partition,
	}

	if val, ok := d.GetOk("monitoring_type"); ok {
		cert.CertValidationOptions = []string{val.(string)}
	}
	if val, ok := d.GetOk("issuer_cert"); ok {
		cert.IssuerCert = val.(string)
	}

	err := client.UploadCertificate(certPath, cert)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error in Importing certificate (%s): %s", name, err))
	}
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
		err = teemDevice.Report(f, "bigip_ssl_certificate", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}

	if val, ok := d.GetOk("ocsp"); ok {
		certValidState := &bigip.CertValidatorState{Name: val.(string)}
		certValidRef := &bigip.CertValidatorReference{}
		certValidRef.Items = append(certValidRef.Items, *certValidState)
		cert.CertValidatorRef = certValidRef
		err = client.UpdateCertificate(certPath, cert)
		if err != nil {
			log.Printf("[ERROR]Unable to add ocsp to the certificate:%v", err)
		}
	}
	return resourceBigipSslCertificateRead(ctx, d, meta)
}

func resourceBigipSslCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading Certificate : " + name)
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

	certificate, err := client.GetCertificate(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if certificate == nil {
		return diag.Errorf("Reading Certificate Failed with certificate :%+v", certificate)
	}
	log.Printf("[INFO] Certificate content:%+v", certificate)
	_ = d.Set("name", certificate.Name)
	_ = d.Set("partition", certificate.Partition)
	_ = d.Set("full_path", certificate.FullPath)
	_ = d.Set("issuer_cert", certificate.IssuerCert)
	if certificate.CertValidationOptions != nil && len(certificate.CertValidationOptions) > 0 {
		monitorType := certificate.CertValidationOptions[0]
		_ = d.Set("monitoring_type", monitorType)
	}

	return nil
}

func resourceBigipSslCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Certificate Name " + name)
	certpath := d.Get("content").(string)
	partition := d.Get("partition").(string)

	cert := &bigip.Certificate{
		Name:      name,
		Partition: partition,
	}

	if val, ok := d.GetOk("monitoring_type"); ok {
		cert.CertValidationOptions = []string{val.(string)}
	}
	if val, ok := d.GetOk("issuer_cert"); ok {
		cert.IssuerCert = val.(string)
	}
	if val, ok := d.GetOk("ocsp"); ok {
		certValidState := &bigip.CertValidatorState{Name: val.(string)}
		certValidRef := &bigip.CertValidatorReference{}
		certValidRef.Items = append(certValidRef.Items, *certValidState)
		cert.CertValidatorRef = certValidRef
	}

	err := client.UpdateCertificate(certpath, cert)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error in Importing certificate (%s): %s", name, err))
	}

	return resourceBigipSslCertificateRead(ctx, d, meta)
}

func resourceBigipSslCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting Certificate " + name)
	partition := d.Get("partition").(string)
	/*if !strings.HasSuffix(name, ".crt") {
		name = name + ".crt"
	}*/
	name = "/" + partition + "/" + name
	err := client.DeleteCertificate(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Pool   (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

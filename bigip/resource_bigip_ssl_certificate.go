package bigip

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipSslCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSslCertificateCreate,
		Read:   resourceBigipSslCertificateRead,
		Update: resourceBigipSslCertificateUpdate,
		Delete: resourceBigipSslCertificateDelete,
		Exists: resourceBigipSslCertificateExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Full Path Name of ssl certificate",
			},
		},
	}
}

func resourceBigipSslCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Certificate Name " + name)

	certPath := d.Get("content").(string)
	partition := d.Get("partition").(string)
	err := client.UploadCertificate(name, certPath, partition)
	if err != nil {
		return fmt.Errorf("Error in Importing certificate (%s): %s ", name, err)
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
	return resourceBigipSslCertificateRead(d, meta)
}

func resourceBigipSslCertificateRead(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	log.Printf("[INFO] Certificate content:%+v", certificate)
	_ = d.Set("name", certificate.Name)
	_ = d.Set("partition", certificate.Partition)
	_ = d.Set("full_path", certificate.FullPath)

	return nil
}

func resourceBigipSslCertificateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Checking certificate " + name + " exists.")
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
		log.Printf("[ERROR] Unable to Retrieve certificate   (%s) (%v) ", name, err)
		return false, err
	}

	if certificate == nil {
		log.Printf("[WARN] certificate (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return certificate != nil, nil
}

func resourceBigipSslCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Certificate Name " + name)
	certpath := d.Get("content").(string)
	partition := d.Get("partition").(string)
	/*if !strings.HasSuffix(name, ".crt") {
		name = name + ".crt"
	}*/
	err := client.UpdateCertificate(name, certpath, partition)
	if err != nil {
		return fmt.Errorf("Error in Importing certificate (%s): %s ", name, err)
	}

	return resourceBigipSslCertificateRead(d, meta)
}

func resourceBigipSslCertificateDelete(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId("")
	return nil
}

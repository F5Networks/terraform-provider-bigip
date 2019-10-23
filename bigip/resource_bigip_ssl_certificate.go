package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
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
				Description: "Name of SSL Certificate",
				//ForceNew:    true,
			},
			"cert_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location of certificate on Disk",
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				Description: "Partition of ssl certificate",
			},
		},
	}
}

func resourceBigipSslCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Certificate Name " + name)
	certpath := d.Get("cert_path").(string)
	partition := d.Get("partition").(string)
	err := client.UploadCertificate(name, certpath, partition)
	if err != nil {
		return fmt.Errorf("Error in Importing certificate (%s): %s", name, err)
	}

	//err = resourceBigipSslCertificateUpdate(d, meta)
	//if err != nil {
	//	client.DeletePool(name)
	//	return err
	//}

	d.SetId(name)
	return resourceBigipSslCertificateRead(d, meta)
}

func resourceBigipSslCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading Certificate : " + name)
	partition := d.Get("partition").(string)
	name = "~" + partition + "~" + name
	certificate, err := client.GetCertificate(name)
	log.Printf("[INFO] Certificate content:%+v", certificate)
	//d.Set("name", certificate.Name)
	//d.Set("name", name)
	if err != nil {
		return err
	}
	return nil
}

func resourceBigipSslCertificateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Checking certificate " + name + " exists.")
	partition := d.Get("partition").(string)
	name = "~" + partition + "~" + name
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
	certpath := d.Get("cert_path").(string)
	partition := d.Get("partition").(string)
	err := client.UpdateCertificate(name, certpath, partition)
	if err != nil {
		return fmt.Errorf("Error in Importing certificate (%s): %s", name, err)
	}

	return resourceBigipSslCertificateRead(d, meta)
}

func resourceBigipSslCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting Certificate " + name)
	partition := d.Get("partition").(string)
	name = "~" + partition + "~" + name
	err := client.DeleteCertificate(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Pool   (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

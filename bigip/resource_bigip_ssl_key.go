package bigip

import (
	"errors"
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipSslKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSslKeyCreate,
		Read:   resourceBigipSslKeyRead,
		Update: resourceBigipSslKeyUpdate,
		Delete: resourceBigipSslKeyDelete,
		Exists: resourceBigipSslKeyExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of SSL Certificate key with .key extension",
				ForceNew:    true,
			},
			"content": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				//ForceNew:    true,
				Description: "Content of SSL certificate key present on local Disk",
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

func resourceBigipSslKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Certificate Key Name " + name)
	certpath := d.Get("content").(string)
	partition := d.Get("partition").(string)
	/*if !strings.HasSuffix(name, ".key") {
		name = name + ".key"
	}*/
	err := client.UploadKey(name, certpath, partition)
	if err != nil {
		return fmt.Errorf("Error in Importing certificate key (%s): %s ", name, err)
	}

	d.SetId(name)
	return resourceBigipSslKeyRead(d, meta)
}

func resourceBigipSslKeyRead(d *schema.ResourceData, meta interface{}) error {
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
	log.Printf("[INFO] SSL key content:%+v", certkey)
	_ = d.Set("name", certkey.Name)
	_ = d.Set("partition", certkey.Partition)
	_ = d.Set("full_path", certkey.FullPath)
	if err != nil {
		return err
	}
	return nil
}

func resourceBigipSslKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Checking certificate key" + name + " exists.")
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
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve certificate key (%s) (%v) ", name, err)
		return false, err
	}

	if certkey == nil {
		log.Printf("[WARN] certificate key(%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return certkey != nil, nil
}

func resourceBigipSslKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Certificate key Name " + name)
	certpath := d.Get("content").(string)
	/*if !strings.HasSuffix(name, ".key") {
		name = name + ".key"
	}*/
	partition := d.Get("partition").(string)
	err := client.UpdateKey(name, certpath, partition)
	if err != nil {
		return fmt.Errorf("Error in Importing certificate (%s): %s ", name, err)
	}

	return resourceBigipSslKeyRead(d, meta)
}

func resourceBigipSslKeyDelete(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId("")
	return nil
}

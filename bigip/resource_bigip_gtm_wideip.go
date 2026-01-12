package bigip

import (
	"fmt"
//	"log"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipGtmWideip() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmWideipCreate,
		Read:   resourceBigipGtmWideipRead,
		Update: resourceBigipGtmWideipUpdate,
		Delete: resourceBigipGtmWideipDelete,
//		Exists: resourceBigipGtmWideipExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the Wideip",
			},
                        "type": {
                                Type:        schema.TypeString,
                                Required:    true,
                        },

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
			},

			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
			},
			"app_service": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
			},

			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
			},

		},
	}
}
func resourceBigipGtmWideipCreate(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	d.SetId(name)
        gtmtype := d.Get("type").(string)
        err := client.AddGTMWideIP(name,gtmtype)
        if err != nil {
		return fmt.Errorf("Error creating wideip (%s): %s", name, err)
	}
        return nil
}

func resourceBigipGtmWideipRead(d *schema.ResourceData, meta interface{}) error {
 return nil
}

func resourceBigipGtmWideipUpdate(d *schema.ResourceData, meta interface{}) error {
  return nil
}
func resourceBigipGtmWideipDelete(d *schema.ResourceData, meta interface{}) error {
  return nil
}

package bigip

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipSysIapp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysIappCreate,
		Update: resourceBigipSysIappUpdate,
		Read:   resourceBigipSysIappRead,
		Delete: resourceBigipSysIappDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipSysIappImporter,
		},

		Schema: map[string]*schema.Schema{

			"jsonfile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the Iapp which needs to be Iappensed",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the Iapp which needs to be Iappensed",
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				Description: "Address of the Iapp which needs to be Iappensed",
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "This is iApp template for application objects",
				Description: "Address of the Iapp which needs to be Iappensed",
			},

			"devicegroup": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
				Description: "BIG-IP password",
			},
			"execute_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"inherited_devicegroup": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "true",
				Description: "BIG-IP password",
			},

			"inherited_traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "true",
				Description: "BIG-IP password",
			},
			"strict_updates": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "BIG-IP password",
			},

			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},

			"template_modified": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "BIG-IP password",
			},
			"template_prerequisite_errors": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},

			"traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/traffic-group-1",
				Description: "BIG-IP password",
			},
			"lists": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encrypted": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "no",
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
					},
				},
			},

			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"persists": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "true",
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
					},
				},
			},

			"tables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
						"column_names": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"encrypted_columns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},

						"rows": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"row": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},

			"variables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},

						"encrypted": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "no",
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
					},
				},
			},
		},
	}
}
func resourceBigipSysIappCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	log.Println("[INFO] Creating Iapp       " + name)
	p := dataToIapp(name, d)
	d.SetId(name)
	d.SetId(description)
	err := client.CreateIapp(&p)

	if err != nil {
		return err
	}
	d.SetId(name)
	//resourceBigipSysIappUpdate(d, meta)
	return resourceBigipSysIappRead(d, meta)
}

func resourceBigipSysIappUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating Iapp " + name)
	p := dataToIapp(name, d)
	err := client.UpdateIapp(name, &p)
	if err != nil {
		return err
	}
	return resourceBigipSysIappRead(d, meta)

}

func resourceBigipSysIappRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Iapp " + name)
	// Create a slice and append three strings to it.

	p, err := client.Iapp(name)
	log.Println(" Value of result in Read for iApp   *************** ", name)
	if err != nil {
		return err
	}
	if p == nil {
		log.Printf("[WARN] IApp (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	d.Set("partition", p.Partition)
	//d.Set("description", p.Description)
	d.Set("devicegroup", p.DeviceGroup)
	d.Set("execute_action", p.ExecuteAction)
	d.Set("inherited_devicegroup", p.InheritedDevicegroup)
	d.Set("inherited_traffic_group", p.InheritedTrafficGroup)
	d.Set("strict_updates", p.StrictUpdates)
	//d.Set("template", p.Template)
	d.Set("template_modified", p.TemplateModified)
	d.Set("template_prerequisite_errors", p.TemplatePrerequisiteErrors)
	d.Set("traffic_group", p.TrafficGroup)
	d.Set("tables", p.Tables)
	d.Set("lists", p.Lists)
	d.Set("variables", p.Variables)
	d.Set("metadata", p.Metadata)
	return nil
}

func resourceBigipSysIappDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteIapp(name)
	if err == nil {
		log.Printf("[WARN] IApp (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigipSysIappImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

// This function "IappToData...." helps to unmarshal json to Go struct
func IappToData(p *bigip.Iapp, d *schema.ResourceData) error {

	return nil
}

func dataToIapp(name string, d *schema.ResourceData) bigip.Iapp {
	var p bigip.Iapp

	jsonblob := []byte(d.Get("jsonfile").(string))
	err := json.Unmarshal(jsonblob, &p)
	if err != nil {
		fmt.Println("error", err)
	}
	return p
}

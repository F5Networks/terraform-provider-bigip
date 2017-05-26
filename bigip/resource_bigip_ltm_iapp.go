package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmiApp() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmiAppCreate,
		Update: resourceBigipLtmiAppUpdate,
		Read:   resourceBigipLtmiAppRead,
		Delete: resourceBigipLtmiAppDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmiAppImporter,
		},
		Schema: map[string]*schema.Schema{

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User defined description",
				//	ValidateFunc: validateF5Name,
			},
			"deviceGroup": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the device group that the application service is assigned to",
			},

			"executeAction": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Run the specified template action associated with the application.",
			},
			"inheritedDevicegroup": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Read-only.Use 'device-group default' or 'device-group non-default' to set this.",
			},
			"inheritedTrafficGroup": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "as its parent folder. Use 'traffic-group default' or 'traffic-group non-default' to set this..",
			},
			"tmPartition": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Displays the administrative partition within which the application resides.",
			},
			"strictUpdates": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "application objects can be modified or not",
			},
			"template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: " Defines configuration for the application",
			},
			"templateModified": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "template has been modified it shows that",
			},
			"templatePrerequisiteErrors": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: " Provides any missing prerequiste asscoiated with template",
			},
			"trafficGroup": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group that application is assign to.",
			},
			"variables": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Name of the traffic group that application is assign to.",
			},
			"encrypted": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group that application is assign to.",
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group that application is assign to.",
			},
			"tables": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Name of the traffic group that application is assign to.",
			},
			"columnNames": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group that application is assign to.",
			},
			"encryptedColumns": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group that application is assign to.",
			},
			"rows": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group that application is assign to.",
			},
		},
	}

}

func resourceBigipLtmiAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Get("description").(string)
	deviceGroup := d.Get("deviceGroup").(string)
	executeAction := d.Get("executeAction").(string)
	inheritedDevicegroup := d.Get("inheritedDevicegroup").(string)
	inheritedTrafficGroup := d.Get("inheritedTrafficGroup").(string)
	tmPartition := d.Get("tmPartition").(string)
	strictUpdates := d.Get("strictUpdates").(string)
	template := d.Get("template").(string)
	templateModified := d.Get("templateModified").(string)
	templatePrerequisiteErrors := d.Get("templatePrerequisiteErrors").(string)
	trafficGroup := d.Get("trafficGroup").(string)
	tables := setToStringSlice(d.Get("tables").(*schema.Set))
	columnNames := d.Get("columnNames").(string)
	encryptedColumns := d.Get("encryptedColumns").(string)
	rows := d.Get("rows").(string)
	variables := setToStringSlice(d.Get("variables").(*schema.Set))
	encrypted := d.Get("encrypted").(string)
	value := d.Get("value").(string)

	log.Println("[INFO] Creating Service ")

	err := client.CreateService(
		description,
		deviceGroup,
		executeAction,
		inheritedDevicegroup,
		inheritedTrafficGroup,
		tmPartition,
		strictUpdates,
		template,
		templateModified,
		templatePrerequisiteErrors,
		trafficGroup,
		tables,
		columnNames,
		encryptedColumns,
		rows,
		variables,
		encrypted,
		value,
	)

	if err != nil {
		return err
	}
	d.SetId(description)
	return nil
}

func resourceBigipLtmiAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Updating Service " + description)

	r := &bigip.Service{
		Description:                description,
		DeviceGroup:                d.Get("deviceGroup").(string),
		ExecuteAction:              d.Get("executeAction").(string),
		InheritedDevicegroup:       d.Get("inheritedDevicegroup").(string),
		InheritedTrafficGroup:      d.Get("inheritedTrafficGroup").(string),
		TmPartition:                d.Get("tmParition").(string),
		StrictUpdates:              d.Get("strictUpdates").(string),
		Template:                   d.Get("template").(string),
		TemplateModified:           d.Get("templateModified").(string),
		TemplatePrerequisiteErrors: d.Get("templatePrerequisiteErrors").(string),
		Tables:           setToStringSlice(d.Get("tables").(*schema.Set)),
		ColumnNames:      d.Get("columnNames").(string),
		EncryptedColumns: d.Get("encryptedColumns").(string),
		Rows:             d.Get("rows").(string),
		Variables:        setToStringSlice(d.Get("variables").(*schema.Set)),
		Encrypted:        d.Get("encrypted").(string),
		Value:            d.Get("value").(string),
	}

	return client.ModifyService(r)
}

func resourceBigipLtmiAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading Service " + description)

	members, err := client.Services()
	if err != nil {
		return err
	}

	d.Set("description", members.Description)
	d.Set("deviceGroup", members.DeviceGroup)
	d.Set("executeAction", members.ExecuteAction)
	d.Set("inheritedDevicegroup", members.InheritedDevicegroup)
	d.Set("inheritedTrafficGroup", members.InheritedTrafficGroup)
	d.Set("tmParition", members.TmPartition)
	d.Set("strictUpdates", members.StrictUpdates)
	d.Set("template", members.Template)
	d.Set("templateModified", members.TemplateModified)
	d.Set("templatePrerequisiteErrors", members.TemplatePrerequisiteErrors)
	d.Set("trafficGroup", members.TrafficGroup)
	d.Set("variables", members.Variables)
	d.Set("encrypted", members.Encrypted)
	d.Set("tables", members.Tables)
	d.Set("columnNames", members.ColumnNames)
	d.Set("encryptedColumns", members.EncryptedColumns)
	d.Set("rows", members.Rows)
	return nil
}

func resourceBigipLtmiAppDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmiAppImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

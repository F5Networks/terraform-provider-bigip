package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipRoleInfo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipRoleInfoCreate,
		ReadContext:   resourceBigipRoleInfoRead,
		UpdateContext: resourceBigipRoleInfoUpdate,
		DeleteContext: resourceBigipRoleInfoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the role info",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the role info",
			},
			"attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The attribute of the role info",
			},
			"console": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The console of the role info",
			},
			"deny": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The deny of the role info",
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The role of the role info",
			},
			"user_partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user partition of the role info",
			},
			"line_order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The line order of the role info",
			},
		},
	}
}

func resourceBigipRoleInfoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	name := d.Get("name").(string)
	attribute := d.Get("attribute").(string)
	console := d.Get("console").(string)
	deny := d.Get("deny").(string)
	description := d.Get("description").(string)
	lineOrder := d.Get("line_order").(int)
	role := d.Get("role").(string)
	userPartition := d.Get("user_partition").(string)

	roleInfo := &bigip.RoleInfo{
		Name:          name,
		Attribute:     attribute,
		Console:       console,
		Deny:          deny,
		Description:   description,
		LineOrder:     lineOrder,
		Role:          role,
		UserPartition: userPartition,
	}

	err := client.CreateRoleInfo(roleInfo)

	if err != nil {
		log.Printf("[ERROR] error while creating the role info: %s", name)
		return diag.FromErr(err)
	}

	d.SetId(name)
	return resourceBigipRoleInfoRead(ctx, d, m)
}

func resourceBigipRoleInfoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	name := d.Id()
	roleInfo, err := client.GetRoleInfo(name)

	if err != nil {
		log.Printf("[ERROR] error while reading the role info: %s", name)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Role Info: %+v", roleInfo)
	d.Set("name", roleInfo.Name)
	d.Set("attribute", roleInfo.Attribute)
	d.Set("console", roleInfo.Console)
	d.Set("deny", roleInfo.Deny)
	d.Set("description", roleInfo.Description)
	d.Set("line-order", roleInfo.LineOrder)
	d.Set("role", roleInfo.Role)
	d.Set("user-partition", roleInfo.UserPartition)

	return nil
}

func resourceBigipRoleInfoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*bigip.BigIP)
	name := d.Id()
	attribute := d.Get("attribute").(string)
	console := d.Get("console").(string)
	deny := d.Get("deny").(string)
	description := d.Get("description").(string)
	lineOrder := d.Get("line_order").(int)
	role := d.Get("role").(string)
	userPartition := d.Get("user_partition").(string)

	roleInfo := &bigip.RoleInfo{
		Name:          name,
		Attribute:     attribute,
		Console:       console,
		Deny:          deny,
		Description:   description,
		LineOrder:     lineOrder,
		Role:          role,
		UserPartition: userPartition,
	}
	err := client.ModifyRoleInfo(name, roleInfo)

	if err != nil {
		log.Printf("[ERROR] error while updating the role info: %s", name)
		return diag.FromErr(err)
	}

	return resourceBigipRoleInfoRead(ctx, d, m)
}

func resourceBigipRoleInfoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Id()
	client := m.(*bigip.BigIP)
	err := client.DeleteRoleInfo(name)
	if err != nil {
		log.Printf("[ERROR] error while deleting the role info: %s", name)
		return diag.FromErr(err)
	}
	return nil
}

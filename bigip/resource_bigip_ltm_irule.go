package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceBigipLtmIRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmIRuleCreate,
		Read:   resourceBigipLtmIRuleRead,
		Update: resourceBigipLtmIRuleUpdate,
		Delete: resourceBigipLtmIRuleDelete,
		Exists: resourceBigipLtmIRuleExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmIRuleImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the iRule",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},

			"irule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The iRule body",
				StateFunc: func(s interface{}) string {
					return strings.TrimSpace(s.(string))
				},
			},
		},
	}
}

func resourceBigipLtmIRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Creating iRule " + name)

	client.CreateIRule(name, d.Get("irule").(string))

	d.SetId(name)

	return resourceBigipLtmIRuleRead(d, meta)
}

func resourceBigipLtmIRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	irule, err := client.IRule(name)
	if err != nil {
		return err
	}
	if irule == nil {
		log.Printf("[WARN] irule (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err := d.Set("irule", irule.Rule); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IRule  to state for IRule (%s): %s", d.Id(), err)
	}
	d.Set("name", name)
	return nil
}

func resourceBigipLtmIRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching iRule " + name)

	irule, err := client.IRule(name)
	if err != nil {
		return false, err
	}

	return irule != nil, nil
}

func resourceBigipLtmIRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	r := &bigip.IRule{
		FullPath: name,
		Rule:     d.Get("irule").(string),
	}

	err := client.ModifyIRule(name, r)
	if err != nil {
		return err
	}
	return resourceBigipLtmIRuleRead(d, meta)
}

func resourceBigipLtmIRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteIRule(name)
	if err != nil {
		return err
	}
	if err == nil {
		log.Printf("[WARN] iRule (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigipLtmIRuleImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

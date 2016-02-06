package bigip

import (
	"log"

	"github.com/scottdware/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmIRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmIRuleCreate,
		Read:   resourceBigipLtmIRuleRead,
		Update: resourceBigipLtmIRuleUpdate,
		Delete: resourceBigipLtmIRuleDelete,
		Exists: resourceBigipLtmIRuleExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Name of the iRule",
				ForceNew: true,
			},

			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "LTM Partition",
				ForceNew: true,
			},

			"irule": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "The iRule body",
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
	if err != nil{
		return err
	}
	d.Set("partition", irule.Partition)
	d.Set("irule", irule.Rule)

	return nil
}

func resourceBigipLtmIRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching iRule " + name)

	_, err := client.IRule(name)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func resourceBigipLtmIRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	r := &bigip.IRule{
		Name: name,
		Partition: d.Get("partition").(string),
		Rule: d.Get("irule").(string),
	}

	return client.ModifyIRule(name, r)
}

func resourceBigipLtmIRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	return client.DeleteIRule(name)
}
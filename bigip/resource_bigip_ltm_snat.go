package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceBigipLtmSnat() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmSnatCreate,
		Update: resourceBigipLtmSnatUpdate,
		Read:   resourceBigipLtmSnatRead,
		Delete: resourceBigipLtmSnatDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmSnatImporter,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Snat list Name",
				//	ValidateFunc: validateF5Name,
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Which partition on BIG-IP",
			},

			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fullpath ",
			},

			"autolasthop": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP autolasthop",
			},
			"mirror": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"sourceport": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"translation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"snatpool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"vlansdisabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "BIG-IP password",
			},

			"origins": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
						"app_service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "app service",
							//ValidateFunc: validateF5Name,
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmSnatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating Snat" + name)

	p := dataToSnat(name, d)
	d.SetId(name)
	err := client.CreateSnat(&p)
	if err != nil {
		return err
	}
	return resourceBigipLtmSnatRead(d, meta)
}

func resourceBigipLtmSnatRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Fetching Ltm Snat " + name)
	p, err := client.GetSnat(name)
	if err != nil {
		d.SetId("")
		return err
	}
	if p == nil {
		log.Printf("[WARN] Snat  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("partition", p.Partition)
	d.Set("full_path", p.FullPath)
	if err := d.Set("full_path", p.FullPath); err != nil {
		return fmt.Errorf("[DEBUG] Error saving FullPath to state for Snat  (%s): %s", d.Id(), err)
	}
	d.Set("autolasthop", p.AutoLasthop)
	if err := d.Set("autolasthop", p.AutoLasthop); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AutoLasthop to state for Snat  (%s): %s", d.Id(), err)
	}
	d.Set("mirror", p.Mirror)
	d.Set("sourceport", p.SourcePort)
	if err := d.Set("sourceport", p.SourcePort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SourcePort to state for Snat  (%s): %s", d.Id(), err)
	}
	d.Set("translation", p.Translation)
	if err := d.Set("translation", p.Translation); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Translation to state for Snat  (%s): %s", d.Id(), err)
	}
	d.Set("snatpool", p.Snatpool)

	if err := d.Set("snatpool", p.Snatpool); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Snatpool to state for Snat  (%s): %s", d.Id(), err)
	}
	d.Set("vlansdisabled", p.VlansDisabled)

	if err != nil {
		return err
	}

	return SnatToData(p, d)
}

func resourceBigipLtmSnatExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching LtmSnat " + name)

	p, err := client.GetSnat(name)

	d.Set("partition", p.Partition)
	d.Set("full_path", p.FullPath)
	d.Set("autolasthop", p.AutoLasthop)
	d.Set("mirror", p.Mirror)
	d.Set("sourceport", p.SourcePort)
	d.Set("translation", p.Translation)
	d.Set("snatpool", p.Snatpool)
	d.Set("vlansdisabled", p.VlansDisabled)

	if err != nil {
		return false, err
	}

	return p != nil, nil
}

func resourceBigipLtmSnatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating LtmSnat " + name)
	p := dataToSnat(name, d)
	err := client.UpdateSnat(name, &p)
	if err != nil {
		return err
	}
	return resourceBigipLtmSnatRead(d, meta)
}

func resourceBigipLtmSnatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteSnat(name)
	if err != nil {
		return err
	}
	if err == nil {
		log.Printf("[WARN] Snat  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigipLtmSnatImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func dataToSnat(name string, d *schema.ResourceData) bigip.Snat {
	var p bigip.Snat

	p.Name = name
	p.Partition = d.Get("partition").(string)
	p.FullPath = d.Get("full_path").(string)
	p.AutoLasthop = d.Get("autolasthop").(string)
	p.Mirror = d.Get("mirror").(string)
	p.SourcePort = d.Get("sourceport").(string)
	p.Translation = d.Get("translation").(string)
	p.Snatpool = d.Get("snatpool").(string)
	p.VlansDisabled = d.Get("vlansdisabled").(bool)

	originsCount := d.Get("origins.#").(int)
	p.Origins = make([]bigip.Originsrecord, 0, originsCount)
	for i := 0; i < originsCount; i++ {
		var r bigip.Originsrecord
		log.Println("I am in dattosnat policy ", p, originsCount, i)
		prefix := fmt.Sprintf("origins.%d", i)
		r.Name = d.Get(prefix + ".name").(string)
		p.Origins = append(p.Origins, r)
	}

	log.Println("I am in DatatoSnat value of p                                                   ", p)

	return p
}

func SnatToData(p *bigip.Snat, d *schema.ResourceData) error {
	d.Set("partition", p.Partition)
	d.Set("full_path", p.FullPath)
	d.Set("autolasthop", p.AutoLasthop)
	d.Set("mirror", p.Mirror)
	d.Set("sourceport", p.SourcePort)
	d.Set("translation", p.Translation)
	d.Set("snatpool", p.Snatpool)
	d.Set("vlansdisabled", p.VlansDisabled)

	for i, r := range p.Origins {
		origins := fmt.Sprintf("origins.%d", i)
		d.Set(fmt.Sprintf("%s.name", origins), r.Name)
	}
	return nil
}

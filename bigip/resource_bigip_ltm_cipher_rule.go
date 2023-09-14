package bigip

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipLtmCipherRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmCipherRuleCreate,
		ReadContext:   resourceBigipLtmCipherRuleRead,
		UpdateContext: resourceBigipLtmCipherRuleUpdate,
		DeleteContext: resourceBigipLtmCipherRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The cipher rule name.",
				Required:    true,
			},
			"partition": {
				Type:        schema.TypeString,
				Description: "The partition name.",
				Optional:    true,
				Default:     "Common",
			},
			"cipher_suites": {
				Type:        schema.TypeString,
				Description: "The cipher suites.",
				Default:     "DEFAULT",
				Optional:    true,
			},
			"dh_groups": {
				Type:        schema.TypeString,
				Description: "The DH groups.",
				Optional:    true,
			},
			"signature_algorithms": {
				Type:        schema.TypeString,
				Description: "The signature algorithms.",
				Optional:    true,
			},
			"full_path": {
				Type:        schema.TypeString,
				Description: "The full path of the cipher rule.",
				Computed:    true,
			},
		},
	}
}

func resourceBigipLtmCipherRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	log.Println("[INFO] Creating Cipher Rule: ", name, " in partition: ", partition)
	cipherRule := &bigip.CipherRule{
		Name:                name,
		Partition:           partition,
		Cipher:              d.Get("cipher_suites").(string),
		DHGroups:            d.Get("dh_groups").(string),
		SignatureAlgorithms: d.Get("signature_algorithms").(string),
	}
	err := client.CreateCipherRule(cipherRule)
	if err != nil {
		return diag.FromErr(err)
	}
	fullPath := fmt.Sprintf("/%s/%s", partition, name)
	d.SetId(fullPath)
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
		err = teemDevice.Report(f, "bigip_ltm_cipher_rule", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipLtmCipherRuleRead(ctx, d, meta)
}

func resourceBigipLtmCipherRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	id := d.Id()
	id = strings.Replace(id, "/", "", 1)
	name_partition := strings.Split(id, "/")
	name := name_partition[1]
	partition := name_partition[0]

	log.Printf("----------------name_partition: %v------------------", name_partition)

	log.Println("[INFO] Reading Cipher Rule: ", name)
	cipherRule, err := client.GetCipherRule(name, partition)
	if err != nil {
		return diag.FromErr(err)
	}
	if cipherRule == nil {
		return diag.FromErr(fmt.Errorf("cipher Rule not found"))
	}
	fullPath := fmt.Sprintf("/%s/%s", partition, name)
	_ = d.Set("name", cipherRule.Name)
	_ = d.Set("partition", cipherRule.Partition)
	_ = d.Set("cipher_suites", cipherRule.Cipher)
	_ = d.Set("dh_groups", cipherRule.DHGroups)
	_ = d.Set("signature_algorithms", cipherRule.SignatureAlgorithms)
	_ = d.Set("full_path", fullPath)
	return nil
}

func resourceBigipLtmCipherRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	log.Println("[INFO] Updating Cipher Rule: ", name, " in partition: ", partition)
	cipherRule := &bigip.CipherRule{
		Name:                name,
		Partition:           partition,
		Cipher:              d.Get("cipher_suites").(string),
		DHGroups:            d.Get("dh_groups").(string),
		SignatureAlgorithms: d.Get("signature_algorithms").(string),
	}
	err := client.ModifyCipherRule(cipherRule)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceBigipLtmCipherRuleRead(ctx, d, meta)
}

func resourceBigipLtmCipherRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	log.Println("[INFO] Deleting Cipher Rule: ", name, " in partition: ", partition)
	err := client.DeleteCipherRule(name, partition)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

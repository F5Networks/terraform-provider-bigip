package bigip

import (
	"context"
	"encoding/json"
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

func resourceBigipSysOcsp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysOcspCreate,
		ReadContext:   resourceBigipSysOcspRead,
		UpdateContext: resourceBigipSysOcspUpdate,
		DeleteContext: resourceBigipSysOcspDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Specifies the name of the OCSP responder. It should be of the pattern '/partition/name'",
				Required:    true,
			},
			"proxy_server_pool": {
				Type:          schema.TypeString,
				Description:   "Specifies the proxy server pool the BIG-IP system uses to fetch the OCSP response. It should be of the pattern '/partition/pool-name'",
				ConflictsWith: []string{"dns_resolver"},
				Optional:      true,
			},
			"dns_resolver": {
				Type:          schema.TypeString,
				Description:   "Specifies the internal DNS resolver the BIG-IP system uses to fetch the OCSP response. It should be of the pattern '/partition/resolver-name'",
				ConflictsWith: []string{"proxy_server_pool"},
				Optional:      true,
			},
			"route_domain": {
				Type:        schema.TypeString,
				Description: "Specifies the route domain for the OCSP responder",
				Optional:    true,
			},
			"concurrent_connections_limit": {
				Type:        schema.TypeInt,
				Description: "Specifies the maximum number of connections per second allowed for the OCSP certificate validator",
				Optional:    true,
				Default:     50,
			},
			"responder_url": {
				Type:        schema.TypeString,
				Description: "Specifies the URL of the OCSP responder",
				Optional:    true,
			},
			"connection_timeout": {
				Type:        schema.TypeInt,
				Description: "Specifies the time interval that the BIG-IP system waits for before ending the connection to the OCSP responder, in seconds",
				Optional:    true,
				Default:     8,
			},
			"trusted_responders": {
				Type:        schema.TypeString,
				Description: "Specifies the certificates used for validating the OCSP response",
				Optional:    true,
			},
			"clock_skew": {
				Type:        schema.TypeInt,
				Description: "Specifies the tolerable absolute difference in the clocks of the responder and the BIG-IP system, in seconds",
				Optional:    true,
				Default:     300,
			},
			"status_age": {
				Type:        schema.TypeInt,
				Description: "Specifies the maximum allowed lag time that the BIG-IP system accepts for the 'thisUpdate' time in the OCSP response, in seconds",
				Optional:    true,
				Default:     0,
			},
			"strict_resp_cert_check": {
				Type:        schema.TypeString,
				Description: "Specifies whether the responder's certificate is checked for an OCSP signing extension",
				Optional:    true,
				Default:     "enabled",
			},
			"cache_timeout": {
				Type:        schema.TypeString,
				Description: "Specifies the lifetime of the OCSP response in the cache, in seconds",
				Optional:    true,
				Default:     "indefinite",
			},
			"cache_error_timeout": {
				Type:        schema.TypeInt,
				Description: "Specifies the lifetime of an error response in the cache, in seconds. This value must be greater than connection_timeout",
				Optional:    true,
				Default:     3600,
			},
			"signer_cert": {
				Type:        schema.TypeString,
				Description: "Specifies a certificate used to sign an OCSP request. It should be of the pattern '/partition/cert-name'",
				Optional:    true,
			},
			"signer_key": {
				Type:        schema.TypeString,
				Description: "Specifies a key used to sign an OCSP request. It should be of the pattern '/partition/key-name'",
				Optional:    true,
			},
			"passphrase": {
				Type:        schema.TypeString,
				Description: "Specifies a passphrase used to sign an OCSP request",
				Sensitive:   true,
				Optional:    true,
			},
			"sign_hash": {
				Type:        schema.TypeString,
				Description: "Specifies the hash algorithm used to sign an OCSP request",
				Default:     "sha256",
				Optional:    true,
			},
		},
	}
}

func resourceBigipSysOcspCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	ocsp := &bigip.OCSP{
		Name: name,
	}

	populateOcspConfig(ocsp, d)

	err := client.CreateOCSP(ocsp)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

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
		err = teemDevice.Report(f, "bigip_sys_ocsp", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}

	return resourceBigipSysOcspRead(ctx, d, meta)
}

func resourceBigipSysOcspRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	id := d.Id()
	id = strings.Trim(id, "/")
	splitArr := strings.Split(id, "/")
	if len(splitArr) != 2 {
		return diag.Errorf("Invalid ID %s", id)
	}

	name := splitArr[1]
	partition := splitArr[0]
	ocspFqdn := fmt.Sprintf("~%s~%s", partition, name)

	ocsp, err := client.GetOCSP(ocspFqdn)
	if err != nil {
		log.Printf("[ERROR] unable to retrieve ocsp %s: %s ", name, err)
		return diag.FromErr(err)
	}

	ocspJson, err := json.Marshal(ocsp)
	if err != nil {
		log.Printf("[ERROR] unable to marshal ocsp %s: %s ", name, err)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] ocsp response: %+v", string(ocspJson))

	d.Set("name", ocsp.FullPath)

	setOcspStateData(d, ocsp)

	return nil
}

func resourceBigipSysOcspUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	id := d.Id()
	id = strings.Trim(id, "/")
	splitArr := strings.Split(id, "/")
	if len(splitArr) != 2 {
		return diag.Errorf("Invalid ID %s", id)
	}

	name := splitArr[1]
	partition := splitArr[0]
	ocspFqdn := fmt.Sprintf("~%s~%s", partition, name)

	ocsp := &bigip.OCSP{
		Name: name,
	}
	populateOcspConfig(ocsp, d)

	err := client.ModifyOCSP(ocspFqdn, ocsp)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceBigipSysOcspRead(ctx, d, meta)
}

func resourceBigipSysOcspDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	id := d.Id()
	id = strings.Trim(id, "/")
	splitArr := strings.Split(id, "/")
	if len(splitArr) != 2 {
		return diag.Errorf("Invalid ID %s", id)
	}

	name := splitArr[1]
	partition := splitArr[0]
	ocspFqdn := fmt.Sprintf("~%s~%s", partition, name)

	err := client.DeleteOCSP(ocspFqdn)

	if err != nil {
		log.Printf("[ERROR] unable to delete ocsp %s: %s ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func populateOcspConfig(ocsp *bigip.OCSP, d *schema.ResourceData) {
	if v, ok := d.GetOk("proxy_server_pool"); ok {
		ocsp.ProxyServerPool = v.(string)
	}
	if v, ok := d.GetOk("dns_resolver"); ok {
		ocsp.DnsResolver = v.(string)
	}
	if v, ok := d.GetOk("concurrent_connections_limit"); ok {
		ocsp.ConcurrentConnectionsLimit = int64(v.(int))
	}
	if v, ok := d.GetOk("responder_url"); ok {
		ocsp.ResponderUrl = v.(string)
	}
	if v, ok := d.GetOk("route_domain"); ok {
		ocsp.RouteDomain = v.(string)
	}
	if v, ok := d.GetOk("connection_timeout"); ok {
		ocsp.ConnectionTimeout = int64(v.(int))
	}
	if v, ok := d.GetOk("trusted_responder"); ok {
		ocsp.TrustedResponders = v.(string)
	}
	if v, ok := d.GetOk("clock_skew"); ok {
		ocsp.ClockSkew = int64(v.(int))
	}
	if v, ok := d.GetOk("status_age"); ok {
		ocsp.StatusAge = int64(v.(int))
	}
	if v, ok := d.GetOk("strict_resp_cert_check"); ok {
		ocsp.StrictRespCertCheck = v.(string)
	}
	if v, ok := d.GetOk("cache_timeout"); ok {
		ocsp.CacheTimeout = v.(string)
	}
	if v, ok := d.GetOk("cache_error_timeout"); ok {
		ocsp.CacheErrorTimeout = int64(v.(int))
	}
	if v, ok := d.GetOk("signer_cert"); ok {
		ocsp.SignerCert = v.(string)
	}
	if v, ok := d.GetOk("signer_key"); ok {
		ocsp.SignerKey = v.(string)
	}
	if v, ok := d.GetOk("passphrase"); ok {
		ocsp.Passphrase = v.(string)
	}
	if v, ok := d.GetOk("sign_hash"); ok {
		ocsp.SignHash = v.(string)
	}
}

func setOcspStateData(d *schema.ResourceData, ocsp *bigip.OCSP) {
	if ocsp.ProxyServerPool != "" {
		d.Set("proxy_server_pool", ocsp.ProxyServerPool)
	} else {
		d.Set("dns_resolver", ocsp.DnsResolver)
	}
	if ocsp.RouteDomain != "" {
		d.Set("route_domain", ocsp.RouteDomain)
	}
	if ocsp.ResponderUrl != "" {
		d.Set("responder_url", ocsp.ResponderUrl)
	}
	if ocsp.TrustedResponders != "" {
		d.Set("trusted_responders", ocsp.TrustedResponders)
	}
	if ocsp.SignerCert != "" {
		d.Set("signer_cert", ocsp.SignerCert)
	}
	if ocsp.SignerKey != "" {
		d.Set("signer_key", ocsp.SignerKey)
	}

	d.Set("concurrent_connections_limit", ocsp.ConcurrentConnectionsLimit)
	d.Set("clock_skew", ocsp.ClockSkew)
	d.Set("status_age", ocsp.StatusAge)
	d.Set("cache_timeout", ocsp.CacheTimeout)
	d.Set("cache_error_timeout", ocsp.CacheErrorTimeout)
	d.Set("connection_timeout", ocsp.ConnectionTimeout)
	d.Set("strict_resp_cert_check", ocsp.StrictRespCertCheck)
	d.Set("sign_hash", ocsp.SignHash)
}

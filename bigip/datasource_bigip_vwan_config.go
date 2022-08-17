/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"

	"log"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	azureEnvErr = `Azure Environment is not set,please set Below Environment variables
AZURE_SUBSCRIPTION_ID
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
AZURE_TENANT_ID
STORAGE_ACCOUNT_NAME
STORAGE_ACCOUNT_KEY`
)

func dataSourceBigipVwanconfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipVwanconfigRead,
		Schema: map[string]*schema.Schema{
			"azure_vwan_resourcegroup": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name azure_vwan_resourcegroup",
			},
			"azure_vwan_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name azure_vwan_name",
			},
			"azure_vwan_vpnsite": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name azure_vwan_vpnsite",
			},
			"bigip_gw_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of BIGIP GW to vwan will establish tunnel",
			},
			"preshared_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "preshared_key used for establish tunnel",
			},
			"hub_address_space": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vWAN Hub address space",
			},
			"hub_connected_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "vWAN Hub connected subnets ",
			},
			"vwan_gw_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IP address of vWAN GW  that will establish tunnel",
			},
		},
	}
}
func dataSourceBigipVwanconfigRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	log.Println("[INFO] Reading VWAN Config for site:" + d.Get("azure_vwan_vpnsite").(string))

	if os.Getenv("AZURE_SUBSCRIPTION_ID") == "" || os.Getenv("AZURE_CLIENT_ID") == "" || os.Getenv("AZURE_CLIENT_SECRET") == "" || os.Getenv("AZURE_TENANT_ID") == "" || os.Getenv("STORAGE_ACCOUNT_NAME") == "" || os.Getenv("STORAGE_ACCOUNT_KEY") == "" {
		return fmt.Errorf("%s", azureEnvErr)
	}
	config := azureConfig{
		subscriptionID:    os.Getenv("AZURE_SUBSCRIPTION_ID"),
		clientID:          os.Getenv("AZURE_CLIENT_ID"),
		clientPassword:    os.Getenv("AZURE_CLIENT_SECRET"),
		tenantID:          os.Getenv("AZURE_TENANT_ID"),
		resourceGroupName: d.Get("azure_vwan_resourcegroup").(string),
		virtualWANName:    d.Get("azure_vwan_name").(string),
		siteName:          d.Get("azure_vwan_vpnsite").(string),
		accountName:       os.Getenv("STORAGE_ACCOUNT_NAME"),
		accountKey:        os.Getenv("STORAGE_ACCOUNT_KEY"),
	}
	res, err := DownloadVwanConfig(config)
	if err != nil {
		log.Printf("failed to download vpnClient Config: %+v", err.Error())
		return err
	}
	log.Printf("[DEBUG] Unmarshed Data : %+v", res)
	for _, v := range res {
		if v["vpnSiteConfiguration"].(map[string]interface{})["Name"] == config.siteName {
			log.Printf("[DEBUG] IPAddress : %+v", v["vpnSiteConfiguration"].(map[string]interface{})["IPAddress"])
			_ = d.Set("bigip_gw_ip", v["vpnSiteConfiguration"].(map[string]interface{})["IPAddress"])
			for _, vv := range v["vpnSiteConnections"].([]interface{}) {

				_ = d.Set("hub_address_space", vv.(map[string]interface{})["hubConfiguration"].(map[string]interface{})["AddressSpace"])
				var hubSubnet []string
				for _, val := range vv.(map[string]interface{})["hubConfiguration"].(map[string]interface{})["ConnectedSubnets"].([]interface{}) {
					hubSubnet = append(hubSubnet, val.(string))
				}
				log.Printf("[DEBUG] hub_connected_subnets : %+v", hubSubnet)
				_ = d.Set("hub_connected_subnets", hubSubnet)
				var vwanGWIPs []string
				vwanGWIPs = append(vwanGWIPs, vv.(map[string]interface{})["gatewayConfiguration"].(map[string]interface{})["IpAddresses"].(map[string]interface{})["Instance0"].(string))
				vwanGWIPs = append(vwanGWIPs, vv.(map[string]interface{})["gatewayConfiguration"].(map[string]interface{})["IpAddresses"].(map[string]interface{})["Instance1"].(string))
				log.Printf("[DEBUG] connectionConfiguration : %+v", vv.(map[string]interface{})["connectionConfiguration"])
				_ = d.Set("preshared_key", vv.(map[string]interface{})["connectionConfiguration"].(map[string]interface{})["PSK"])
				log.Printf("[DEBUG] vwan_gw_address : %+v", vwanGWIPs)
				_ = d.Set("vwan_gw_address", vwanGWIPs)
			}
		}
	}

	d.SetId(d.Get("azure_vwan_vpnsite").(string))

	return nil
}

const CountToEnd = 0

type azureConfig struct {
	subscriptionID    string
	clientID          string
	clientPassword    string
	tenantID          string
	resourceGroupName string
	virtualWANName    string
	siteName          string
	accountName       string
	accountKey        string
}

func DownloadVwanConfig(config azureConfig) ([]map[string]interface{}, error) {
	subscriptionID := config.subscriptionID
	clientID := config.clientID
	clientPassword := config.clientPassword
	tenantID := config.tenantID
	vpnconfigClient := network.NewVpnSitesConfigurationClient(subscriptionID)
	_, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	resourceGroupName := config.resourceGroupName
	virtualWANName := config.virtualWANName
	siteName := config.siteName
	siteId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnSites/%s", subscriptionID, resourceGroupName, siteName)
	siteList := []string{siteId}
	accountName := config.accountName
	accountKey := config.accountKey
	containerName := "myvpnsiteconfig"
	destFileName := "vpnconfigdownload.json"

	cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Printf("Azure Blob New SAS Failed %+v", err.Error())
		return nil, err
	}
	containerURL1 := azblob.NewContainerURL(*cURL, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	_, err = containerURL1.Create(context.Background(), nil, "")

	defer func() {
		if _, err := containerURL1.Delete(context.Background(), azblob.ContainerAccessConditions{}); err != nil {
			log.Printf("[DEBUG] Could not delete contrainer: %v", err)
		}
	}()

	if err != nil {
		log.Printf("NewContainerURL failed %+v", err.Error())
		return nil, err
	}
	containerURL, _ := url.Parse(containerURL1.String())

	cURL1, _ := url.Parse(fmt.Sprintf("%s/%s", containerURL, destFileName))

	// Set the desired SAS signature values and sign them with the shared key credentials to get the SAS query parameters.
	sasQueryParams, err := azblob.AccountSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,              // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:   azblob.AccountSASPermissions{Read: true, Write: true, List: true}.String(),
		Services:      azblob.AccountSASServices{Blob: true}.String(),
		ResourceTypes: azblob.AccountSASResourceTypes{Container: true, Object: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		log.Printf("SAS Query failed %+v", err.Error())
		return nil, err
	}

	qp := sasQueryParams.Encode()
	urlToSendToSomeone := fmt.Sprintf("%s?%s", cURL1, qp)
	// At this point, you can send the urlToSendToSomeone to someone via email or any other mechanism you choose.

	// ************************************************************************************************

	// When someone receives the URL, they access the SAS-protected resource with code like this:
	u, _ := url.Parse(urlToSendToSomeone)

	// Create an ServiceURL object that wraps the service URL (and its SAS) and a pipeline.
	// When using a SAS URLs, anonymous credentials are required.
	serviceURL := azblob.NewServiceURL(*u, azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	// Now, you can use this serviceURL just like any other to make requests of the resource.

	// You can parse a URL into its constituent parts:

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Here's how to create a blob with HTTP headers and metadata (I'm using the same metadata that was put on the container):
	blobURL := azblob.NewBlockBlobURL(*cURL1, p)

	sasUrl := serviceURL.String()
	log.Printf("[DEBUG] sasUrl : %+v", sasUrl)

	token, _ := CreateToken(tenantID, clientID, clientPassword)
	vpnconfigClient.Authorizer = autorest.NewBearerAuthorizer(token)

	_, err = vpnconfigClient.Download(
		context.Background(),
		resourceGroupName,
		virtualWANName,
		network.GetVpnSitesConfigurationRequest{
			OutputBlobSasURL: &sasUrl,
			VpnSites:         &siteList,
		})

	if err != nil {
		log.Printf("Push vWAN config to sasUrl Failed : %+v", err.Error())
		return nil, err
	}
	time.Sleep(10 * time.Second)

	destFile, err := os.Create(destFileName)
	if err != nil {
		log.Printf("Unable to create destination file %+v", err.Error())
		return nil, err
	}

	defer destFile.Close()

	// Perform download
	err = azblob.DownloadBlobToFile(context.Background(), blobURL.BlobURL, 0, CountToEnd, destFile,
		azblob.DownloadFromBlobOptions{})

	if err != nil {
		log.Printf("Download Blob to file failed %+v", err.Error())
		return nil, err
	}
	data, err := os.ReadFile("vpnconfigdownload.json")
	if err != nil {
		log.Printf("failed to Read vpnconfigdownload.json: %+v", err.Error())
		return nil, err
	}
	var result []map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateToken creates a service principal token
func CreateToken(tenantID, clientID, clientSecret string) (adal.OAuthTokenProvider, error) {
	const activeDirectoryEndpoint = "https://login.microsoftonline.com/"
	var token adal.OAuthTokenProvider

	oauthConfig, err := adal.NewOAuthConfig(activeDirectoryEndpoint, tenantID)
	if err != nil {
		return nil, err
	}

	// The resource for which the token is acquired
	activeDirectoryResourceID := "https://management.azure.com/"
	token, err = adal.NewServicePrincipalToken(
		*oauthConfig,
		clientID,
		clientSecret,
		activeDirectoryResourceID)

	return token, err
}

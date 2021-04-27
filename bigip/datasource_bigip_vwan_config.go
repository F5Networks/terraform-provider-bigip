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
	//"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"
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
			"azure_subsciption_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZURE_SUBSCRIPTION_ID", nil),
				Description: "Specifies the Azure subscription ID to use",
			},
			"azure_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZURE_CLIENT_ID", nil),
				Description: "Specifies the app client ID to use",
			},
			"azure_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZURE_CLIENT_SECRET", nil),
				Description: "Specifies the app secret to use",
			},
			"azure_tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZURE_TENANT_ID", nil),
				Description: "Specifies the Tenant to which to authenticate",
			},
			"storage_accounnt_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("STORAGE_ACCOUNT_NAME", nil),
				Description: "Specifies the Azure subscription ID to use",
			},
			"storage_accounnt_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("STORAGE_ACCOUNT_KEY", nil),
				Description: "Specifies the Azure subscription ID to use",
			},
		},
	}
}
func dataSourceBigipVwanconfigRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	d.SetId("")
	log.Println("[INFO] Reading VWAN Config for site:" + d.Get("azure_vwan_vpnsite").(string))
	config := azureConfig{
		subscriptionID:    d.Get("azure_subsciption_id").(string),
		clientID:          d.Get("azure_client_id").(string),
		clientPassword:    d.Get("azure_client_secret").(string),
		tenantID:          d.Get("azure_tenant_id").(string),
		resourceGroupName: d.Get("azure_vwan_resourcegroup").(string),
		virtualWANName:    d.Get("azure_vwan_name").(string),
		siteName:          d.Get("azure_vwan_vpnsite").(string),
		accountName:       d.Get("storage_accounnt_name").(string),
		accountKey:        d.Get("storage_accounnt_key").(string),
	}
	res, err := DownloadVwanConfig(config)
	if err != nil {
		log.Printf("failed to download vpnClient Config: %+v", err.Error())
	}
	log.Printf("[DEBUG] Unmarshed Data : %+v", res)
	//log.Println("[INFO] Reading VWAN Config : " + d.Get("azure_vwan_resourcegroup").(string))
	//DownloadVwanConfig()

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

type azureVwanGwconfig struct {
	edgeGwaddress    string
	vwanGwaddress    string
	vwanAddressSpace []string
	gwPsk            string
}

func DownloadVwanConfig(config azureConfig) ([]map[string]interface{}, error) {
	subscriptionID := config.subscriptionID
	clientID := config.clientID
	clientPassword := config.clientPassword
	tenantID := config.tenantID
	vpnconfigClient := network.NewVpnSitesConfigurationClient(subscriptionID)
	_, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	//defer resources.Cleanup(ctx)
	resourceGroupName := config.resourceGroupName
	virtualWANName := config.virtualWANName
	siteName := config.siteName
	siteId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnSites/%s", subscriptionID, resourceGroupName, siteName)
	siteList := []string{siteId}
	accountName := config.accountName
	accountKey := config.accountKey
	containerName := "myvpnsiteconfig"
	destFileName := "vpnconfigdownload.json"
	//strorageClient, err := storage.NewBasicClient(accountName, accountKey)
	//log.Printf("[DEBUG] strorageClient : %+v", strorageClient)

	// From the Azure portal, get your Storage account blob service URL endpoint.
	//cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/vpnsiteconfig/%s", accountName, destFileName))
	cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Printf("Azure Blob New SAS Failed %+v", err.Error())
		return nil, err
	}
	containerURL1 := azblob.NewContainerURL(*cURL, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	_, err = containerURL1.Create(context.Background(), nil, "")
	defer containerURL1.Delete(context.Background(), azblob.ContainerAccessConditions{})
	if err != nil {
		log.Printf("NewContainerURL failed %+v", err.Error())
		return nil, err
		//c.Fatal(err)
	}
	containerURL, _ := url.Parse(fmt.Sprintf("%s", containerURL1))

	cURL1, _ := url.Parse(fmt.Sprintf("%s/%s", containerURL, destFileName))
	//log.Printf("[DEBUG] cURL1  : %+v", cURL1)

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
	//blobURLParts := azblob.NewBlobURLParts(serviceURL.URL())

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// // From the Azure portal, get your Storage account blob service URL endpoint.
	// cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer", accountName))

	// Create an ServiceURL object that wraps the service URL and a request pipeline to making requests.
	//containerURL := azblob.NewContainerURL(*cURL, p)
	//log.Printf("containerURL : %+v", containerURL)
	// Here's how to create a blob with HTTP headers and metadata (I'm using the same metadata that was put on the container):
	blobURL := azblob.NewBlockBlobURL(*cURL1, p)
	//sasUrl := fmt.Sprintf("%s", blobURL.BlobURL)
	sasUrl := fmt.Sprintf("%s", serviceURL)
	log.Printf("[DEBUG] sasUrl : %+v", sasUrl)
	//log.Printf("blobURL : %+v", blobURL.BlobURL)

	token, _ := CreateToken(tenantID, clientID, clientPassword)
	vpnconfigClient.Authorizer = autorest.NewBearerAuthorizer(token)
	//log.Printf("[INFO]vpnconfigClient:%+v", vpnconfigClient.Authorizer)
	_, err = vpnconfigClient.Download(
		context.Background(),
		resourceGroupName,
		virtualWANName,
		network.GetVpnSitesConfigurationRequest{
			OutputBlobSasURL: &sasUrl,
			VpnSites:         &siteList,
		})
	//err = future1.WaitForCompletionRef(ctx, vpnconfigClient.Client)
	if err != nil {
		log.Printf("Push vWAN config to sasUrl Failed : %+v", err.Error())
		return nil, err
	}
	time.Sleep(10 * time.Second)

	destFile, err := os.Create(destFileName)
	defer destFile.Close()

	// Perform download
	err = azblob.DownloadBlobToFile(context.Background(), blobURL.BlobURL, 0, CountToEnd, destFile,
		azblob.DownloadFromBlobOptions{})

	if err != nil {
		log.Printf("Download Blob to file failed %+v", err.Error())
		return nil, err
	}
	data, err := ioutil.ReadFile("vpnconfigdownload.json")
	if err != nil {
		log.Printf("failed to Read vpnconfigdownload.json: %+v", err.Error())
		return nil, err
	}
	result := []map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//CreateToken creates a service principal token
func CreateToken(tenantID, clientID, clientSecret string) (adal.OAuthTokenProvider, error) {
	const activeDirectoryEndpoint = "https://login.microsoftonline.com/"
	var token adal.OAuthTokenProvider
	oauthConfig, err := adal.NewOAuthConfig(activeDirectoryEndpoint, tenantID)
	// The resource for which the token is acquired
	activeDirectoryResourceID := "https://management.azure.com/"
	token, err = adal.NewServicePrincipalToken(
		*oauthConfig,
		clientID,
		clientSecret,
		activeDirectoryResourceID)
	return token, err
}

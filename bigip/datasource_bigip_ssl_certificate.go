/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceBigipSslCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipSslCertificateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the certificate",
			},
			"partition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of resource group",
			},

			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate body",
				},
			},
     }
}
func dataSourceBigipSslCertificateRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*bigip.BigIP)
	d.SetId("")
	name := fmt.Sprintf("/%s/%s", d.Get("partition").(string), d.Get("name").(string))

	log.Println("[INFO] Reading Certificate : " + name)
	certificate, err := client.GetCertificate(name)
	log.Printf("[DEBUG] cert is :%v",certificate)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Certificate content:%+v", certificate)
	d.Set("name", certificate.Name)
	d.Set("partition", certificate.Partition)

	d.SetId(certificate.FullPath)

	return nil

}

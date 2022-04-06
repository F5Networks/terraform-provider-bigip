/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmProfileHttpcompress() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileHttpcompressCreate,
		Update: resourceBigipLtmProfileHttpcompressUpdate,
		Read:   resourceBigipLtmProfileHttpcompressRead,
		Delete: resourceBigipLtmProfileHttpcompressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
				Description:  "Name of the Httpcompress Profile",
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Httpcompress profile",
			},

			"uri_exclude": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Servers Address",
			},
			"uri_include": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Servers Address",
			},
			"content_type_include": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Specifies a list of content types for compression of HTTP Content-Type responses. Use a string list to specify a list of content types you want to compress.",
			},
			"content_type_exclude": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Specifies a list of content types for compression of HTTP Content-Type responses. Use a string list to specify a list of content types you want to exclude.",
			},
		},
	}
}

func resourceBigipLtmProfileHttpcompressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	uriExclude := setToStringSlice(d.Get("uri_exclude").(*schema.Set))
	uriInclude := setToStringSlice(d.Get("uri_include").(*schema.Set))
	contentTypeInclude := setToStringSlice(d.Get("content_type_include").(*schema.Set))
	contentTypeExclude := setToStringSlice(d.Get("content_type_exclude").(*schema.Set))

	log.Println("[INFO] Creating Httpcompress profile")

	/*	err := client.CreateHttpcompress(
			name,
			defaultsFrom,
			uriExclude,
			uriInclude,
			contentTypeInclude,
			contentTypeExclude,
		)
	*/
	r := &bigip.Httpcompress{
		Name:               name,
		DefaultsFrom:       defaultsFrom,
		UriExclude:         uriExclude,
		UriInclude:         uriInclude,
		ContentTypeInclude: contentTypeInclude,
		ContentTypeExclude: contentTypeExclude,
	}
	err := client.CreateHttpcompress(r)

	if err != nil {
		return fmt.Errorf("Error retrieving profile Http compress (%s): %s", name, err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileHttpcompressRead(d, meta)
}

func resourceBigipLtmProfileHttpcompressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	r := &bigip.Httpcompress{
		Name:               name,
		DefaultsFrom:       d.Get("defaults_from").(string),
		UriExclude:         setToStringSlice(d.Get("uri_exclude").(*schema.Set)),
		UriInclude:         setToStringSlice(d.Get("uri_include").(*schema.Set)),
		ContentTypeInclude: setToStringSlice(d.Get("content_type_include").(*schema.Set)),
		ContentTypeExclude: setToStringSlice(d.Get("content_type_exclude").(*schema.Set)),
	}

	err := client.ModifyHttpcompress(name, r)
	if err != nil {
		return fmt.Errorf("Error modifying  profile Http compress (%s): %s", name, err)
	}
	return resourceBigipLtmProfileHttpcompressRead(d, meta)
}

func resourceBigipLtmProfileHttpcompressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetHttpcompress(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Http Compress Profile (%s) (%v)", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Httpcompress Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	if _, ok := d.GetOk("defaults_from"); ok {
		if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
			return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("uri_include"); ok {
		if err := d.Set("uri_include", obj.UriInclude); err != nil {
			return fmt.Errorf("[DEBUG] Error saving UriInclude to state for  Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("uri_exclude"); ok {
		if err := d.Set("uri_exclude", obj.UriExclude); err != nil {
			return fmt.Errorf("[DEBUG] Error saving UriExclude to state for  Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("content_type_include"); ok {
		if err := d.Set("content_type_include", obj.ContentTypeInclude); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ContentTypeInclude to state for  Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("content_type_exclude"); ok {
		if err := d.Set("content_type_exclude", obj.ContentTypeExclude); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ContentTypeExclude to state for  Http Compress profile  (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func resourceBigipLtmProfileHttpcompressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Httpcompress Profile " + name)

	err := client.DeleteHttpcompress(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Httpcompress  (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

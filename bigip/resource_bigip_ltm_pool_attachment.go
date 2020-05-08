/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

func resourceBigipLtmPoolAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPoolAttachmentCreate,
		Read:   resourceBigipLtmPoolAttachmentRead,
		Delete: resourceBigipLtmPoolAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmPoolAttachmentImport,
		},

		Schema: map[string]*schema.Schema{
			"pool": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the pool",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},

			"node": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePoolMemberName,
				Description:  "Node to add/remove to/from the pool. Format /partition/node_name:port. e.g. /Common/node01:443",
			},
		},
	}
}

func resourceBigipLtmPoolAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	poolName := d.Get("pool").(string)
	nodeName := d.Get("node").(string)
	parts := strings.Split(nodeName, ":")
	node1, err := client.GetNode(parts[0])
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve node %s  %v :", nodeName, err)
		return err
	}
	if node1 == nil {
		log.Printf("[WARN] Node (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if node1.FQDN.Name != "" {
		config := &bigip.PoolMemberFqdn{
			Name: nodeName,
		}
		config.FQDN.Name = node1.FQDN.Name
		config.FQDN.Interval = node1.FQDN.Interval
		config.FQDN.AddressFamily = node1.FQDN.AddressFamily
		config.FQDN.AutoPopulate = node1.FQDN.AutoPopulate
		config.FQDN.DownInterval = node1.FQDN.DownInterval
		err = client.AddPoolMemberFQDN(poolName, config)
		if err != nil {
			return fmt.Errorf("Failure adding node %s to pool %s: %s", nodeName, poolName, err)
		}
		d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))
		return nil
	}
	log.Printf("[INFO] Adding node %s to pool: %s", nodeName, poolName)
	err = client.AddPoolMember(poolName, nodeName)
	if err != nil {
		return fmt.Errorf("Failure adding node %s to pool %s: %s", nodeName, poolName, err)
	}

	d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))

	return nil
}

func resourceBigipLtmPoolAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	poolName := d.Get("pool").(string)

	// only add the instance that was previously defined for this resource
	expected := d.Get("node").(string)

	pool, err := client.GetPool(poolName)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrive Pool (%s)  (%v) ", poolName, err)
		return err
	}
	if pool == nil {
		log.Printf("[WARN] Pool (%s) not found, removing from state", poolName)
		d.SetId("")
		return nil
	}

	nodes, err := client.PoolMembers(poolName)
	if err != nil {
		return fmt.Errorf("Error retrieving pool (%s) members: %s", poolName, err)
	}
	if nodes == nil {
		log.Printf("[WARN] Pool Members (%s) not found, removing from state", poolName)
		d.SetId("")
		return nil
	}

	// only set the instance Id that this resource manages
	found := false
	for _, node := range nodes.PoolMembers {
		if expected == node.FullPath {
			d.Set("node", expected)
			found = true
			break
		}
	}

	if !found {
		log.Printf("[WARN] Node %s is not a member of pool %s", expected, poolName)
		d.SetId("")
	}

	return nil
}

func resourceBigipLtmPoolAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	poolName := d.Get("pool").(string)
	nodeName := d.Get("node").(string)

	log.Printf("[INFO] Removing node %s from pool: %s", nodeName, poolName)

	err := client.DeletePoolMember(poolName, nodeName)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete PoolMember (%s)  (%s) ", nodeName, err)
		return fmt.Errorf("Failure removing node %s from pool %s: %s", nodeName, poolName, err)
	}
	d.SetId("")
	return nil
}

func resourceBigipLtmPoolAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*bigip.BigIP)

	var data map[string]string
	if err := json.Unmarshal([]byte(d.Id()), &data); err != nil {
		return nil, err
	}
	poolName, ok := data["pool"]
	if !ok {
		return nil, errors.New("missing pool name in input data")
	}
	expectedNode, ok := data["node"]
	if !ok {
		return nil, errors.New("missing node name in input data")
	}

	id := poolName + "-" + expectedNode

	pool, err := client.GetPool(poolName)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve pool %s from bigip: %v", poolName, err)
	}
	if pool == nil {
		return nil, fmt.Errorf("unable to find the pool %s in bigip", poolName)
	}

	nodes, err := client.PoolMembers(poolName)
	if err != nil {
		return nil, errors.New("error retrieving pool members")
	}

	// only set the instance Id that this resource manages
	found := false
	for _, node := range nodes.PoolMembers {
		if expectedNode == node.FullPath {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("cannot locate node %s in pool %s", expectedNode, poolName)
	}

	d.Set("pool", poolName)
	d.Set("node", expectedNode)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

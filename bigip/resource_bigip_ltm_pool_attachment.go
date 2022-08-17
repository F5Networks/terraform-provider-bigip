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
	"log"
	"regexp"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmPoolAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPoolAttachmentCreate,
		Read:   resourceBigipLtmPoolAttachmentRead,
		Update: resourceBigipLtmPoolAttachmentUpdate,
		Delete: resourceBigipLtmPoolAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmPoolAttachmentImport,
		},
		Schema: map[string]*schema.Schema{
			"pool": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the pool to be attached with pool members",
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
			},
			"node": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Poolmember to add/remove to/from the pool. Format node_address:port. e.g 1.1.1.1:80",
			},
			"ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the ratio weight to assign to the pool member. Valid values range from 1 through 65535. The default is 1, which means that each pool member has an equal ratio proportion.",
				Computed:    true,
			},
			"priority_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Specifies a number representing the priority group for the pool member. The default is 0, meaning that the member has no priority. To specify a priority, you must activate priority group usage when you create a new pool or when adding or removing pool members. " +
					"When activated, the system load balances traffic according to the priority group number assigned to the pool member. The higher the number, the higher the priority, so a member with a priority of 3 has higher priority than a member with a priority of 1.",
				Computed: true,
			},
			"connection_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Specifies a maximum established connection limit for a pool member or node. When the current connections count reaches this number, the system does not send additional connections to that pool member or node. The default is 0, meaning that there is no limit to the number of connections." +
					" When used with the weighted least connections load balancing methods, the system uses connection limits to determine the proportional load of each pool member or node. " +
					"This must be a value other than 0 when specified for the weighted least connections load balancing methods",
				Computed: true,
			},
			"connection_rate_limit": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Specifies the maximum number of connections-per-second allowed for a pool member. When the number of connections-per-second reaches the limit for a given pool member, the system drops (UDP) or resets (TCP) additional connection requests. " +
					"This helps detect Denial of Service attacks, where connection requests flood a pool member. Setting this to 0 turns off connection limits. The default is 0.",
				Computed: true,
			},
			"dynamic_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the dynamic ratio number for the node. Used for dynamic ratio load balancing. ",
				Computed:    true,
			},
			"fqdn_autopopulate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies whether the node should scale to the IP address set returned by DNS.",
			},
		},
	}
}

func resourceBigipLtmPoolAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	poolName := d.Get("pool").(string)
	nodeName := d.Get("node").(string)
	poolPartition := strings.Split(poolName, "/")[1]
	parts := strings.Split(nodeName, ":")
	log.Printf("[INFO][CREATE] Attaching Node :%+v to pool : %+v", nodeName, poolName)
	re := regexp.MustCompile(`/([a-zA-z0-9?_-]+)/([a-zA-z0-9.?_-]+):(\d+)`)
	match := re.FindStringSubmatch(nodeName)
	if match != nil {
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
				return fmt.Errorf("Failure adding node %s to pool %s: %s ", nodeName, poolName, err)
			}
			err = resourceBigipLtmPoolAttachmentUpdate(d, meta)
			if err != nil {
				return err
			}
			d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))
			return nil
		}
		log.Printf("[INFO][CREATE] Adding node : %+v to pool: %+v", nodeName, poolName)
		err = client.AddPoolMemberNode(poolName, nodeName)
		if err != nil {
			return fmt.Errorf("Failure adding node %s to pool %s: %s ", nodeName, poolName, err)
		}
		d.SetId(fmt.Sprintf("%s-%s", poolName, nodeName))
		err = resourceBigipLtmPoolAttachmentUpdate(d, meta)
		if err != nil {
			return err
		}
		return nil
	} else {
		log.Println("[DEBUG] creating node from pool attachment resource")
		ipNode := strings.Split(parts[0], "%")[0]
		config := &bigip.PoolMember{
			Name:      nodeName,
			Partition: poolPartition,
		}
		log.Printf("[DEUBG]: Node info:%+v with part:%+v", nodeName, parts[0])
		if !IsValidIP(ipNode) {
			var autoPopulate string
			if d.Get("fqdn_autopopulate").(string) == "" {
				autoPopulate = "enabled"
			} else {
				autoPopulate = d.Get("fqdn_autopopulate").(string)
			}
			config.FQDN.Name = ipNode
			config.FQDN.AutoPopulate = autoPopulate
		}
		log.Printf("[INFO] Adding Pool member (%s) to pool (%s)", nodeName, poolName)
		err := client.AddPoolMember(poolName, config)
		if err != nil {
			return fmt.Errorf("Failure adding node %s to pool %s: %s ", nodeName, poolName, err)
		}
		d.SetId(poolName)
		err = resourceBigipLtmPoolAttachmentUpdate(d, meta)
		if err != nil {
			return err
		}
	}
	return resourceBigipLtmPoolAttachmentRead(d, meta)
}

func resourceBigipLtmPoolAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	nodeName := d.Get("node").(string)
	log.Printf("[DEBUG][UPDATE] node name is :%s", nodeName)
	re := regexp.MustCompile(`/([a-zA-z0-9?_-]+)/([a-zA-z0-9.?_-]+):(\d+)`)
	match := re.FindStringSubmatch(nodeName)
	if match != nil {
		parts := strings.Split(nodeName, ":")
		node1, err := client.GetNode(parts[0])
		if err != nil {
			return err
		}
		if node1 == nil {
			log.Printf("[WARN] Node (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		poolMem := strings.Split(nodeName, ":")[0]
		nodeName1 := strings.Split(poolMem, "/")[2]
		poolName := d.Get("pool").(string)
		config := &bigip.PoolMember{
			Name:            nodeName1,
			FullPath:        nodeName,
			ConnectionLimit: d.Get("connection_limit").(int),
			DynamicRatio:    d.Get("dynamic_ratio").(int),
			PriorityGroup:   d.Get("priority_group").(int),
			RateLimit:       d.Get("connection_rate_limit").(string),
			Ratio:           d.Get("ratio").(int),
		}
		if node1.FQDN.Name != "" {
			log.Printf("[DEBUG] adding autopopulate for fqdn ")
			var autoPopulate string
			if d.Get("fqdn_autopopulate").(string) == "" {
				autoPopulate = "enabled"
			} else {
				autoPopulate = d.Get("fqdn_autopopulate").(string)
			}
			config.FQDN.Name = node1.FQDN.Name
			config.FQDN.AutoPopulate = autoPopulate
		}
		log.Printf("[DEBUG] [UPDATE] pool config :%+v", config)
		err = client.ModifyPoolMember(poolName, config)
		if err != nil {
			return fmt.Errorf("Failure adding node %s to pool %s: %s ", nodeName, poolName, err)
		}
	} else {
		poolName := d.Id()
		poolPartition := strings.Split(poolName, "/")[1]
		nodeName := d.Get("node").(string)
		parts := strings.Split(nodeName, ":")
		ipNode := strings.Split(parts[0], "%")[0]
		poolMem := fmt.Sprintf("/%s/%s", poolPartition, nodeName)
		config := &bigip.PoolMember{
			Name:            nodeName,
			FullPath:        poolMem,
			Address:         parts[0],
			ConnectionLimit: d.Get("connection_limit").(int),
			DynamicRatio:    d.Get("dynamic_ratio").(int),
			PriorityGroup:   d.Get("priority_group").(int),
			RateLimit:       d.Get("connection_rate_limit").(string),
			Ratio:           d.Get("ratio").(int),
		}
		log.Printf("[INFO] Modifying pool member (%+v) from pool (%+v)", poolMem, poolName)
		if !IsValidIP(ipNode) {
			var autoPopulate string
			if d.Get("fqdn_autopopulate").(string) == "" {
				autoPopulate = "enabled"
			} else {
				autoPopulate = d.Get("fqdn_autopopulate").(string)
			}
			config.FQDN.Name = ipNode
			config.FQDN.AutoPopulate = autoPopulate
		}
		log.Printf("[DEBUG] [UPDATE] pool config :%+v", config)
		err := client.ModifyPoolMember2(poolName, config)
		if err != nil {
			return fmt.Errorf("Failure adding node %s to pool %s: %s ", nodeName, poolName, err)
		}
	}
	return resourceBigipLtmPoolAttachmentRead(d, meta)
}

func resourceBigipLtmPoolAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var poolName string
	nodeName := d.Get("node").(string)
	log.Printf("[DEBUG] Reading node name is :%s", nodeName)
	re := regexp.MustCompile(`/([a-zA-z0-9?_-]+)/([a-zA-z0-9.?_-]+):(\d+)`)
	match := re.FindStringSubmatch(nodeName)
	if match != nil {
		poolName = d.Get("pool").(string)
	} else {
		poolName = d.Id()
	}

	// only add the instance that was previously defined for this resource
	expected := d.Get("node").(string)

	pool, err := client.GetPool(poolName)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Pool (%s)  (%v) ", poolName, err)
		return err
	}
	if pool == nil {
		log.Printf("[WARN] Pool (%s) not found, removing from state", poolName)
		d.SetId("")
		return nil
	}
	nodes, err := client.PoolMembers(poolName)
	if err != nil {
		return fmt.Errorf("Error retrieving pool (%s) members: %s ", poolName, err)
	}
	if nodes == nil {
		log.Printf("[WARN] Pool Members (%s) not found, removing from state", poolName)
		d.SetId("")
		return nil
	}
	// only set the instance Id that this resource manages
	found := false

	if match != nil {
		for _, node := range nodes.PoolMembers {
			if expected == node.FullPath {
				_ = d.Set("node", expected)
				found = true
				break
			}
		}
	} else {

		for _, node := range nodes.PoolMembers {
			if expected == node.Name {
				_ = d.Set("node", expected)
				found = true
				break
			}
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
		return fmt.Errorf("Failure removing node %s from pool %s: %s ", nodeName, poolName, err)
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

	_ = d.Set("pool", poolName)
	_ = d.Set("node", expectedNode)

	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

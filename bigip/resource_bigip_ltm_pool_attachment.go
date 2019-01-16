package bigip

import (
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmPoolAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPoolAttachmentCreate,
		Read:   resourceBigipLtmPoolAttachmentRead,
		Delete: resourceBigipLtmPoolAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

	log.Printf("[INFO] Adding node %s to pool: %s", nodeName, poolName)
	err := client.AddPoolMember(poolName, nodeName)
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

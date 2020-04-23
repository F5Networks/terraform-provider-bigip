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
	"strconv"
	"strings"
)

func resourceBigiqLicenseManage() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigiqLicenseManageCreate,
		Read:   resourceBigiqLicenseManageRead,
		Update: resourceBigiqLicenseManageUpdate,
		Delete: resourceBigiqLicenseManageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bigiq_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"bigiq_user": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_token_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Sensitive:   true,
				Default:     false,
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_TOKEN_AUTH", nil),
			},
			"bigiq_login_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "tmos",
				Description: "Login reference for token authentication (see BIG-IQ REST docs for details)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_LOGIN_REF", nil),
			},
			"pool_license_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePoolLicenseType,
				Description:  "This will specify Utility/regKey Licence pool type",
			},
			"assignment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAssignmentType,
				Description:  "Whether the specified device is a managed/un-managed/un-reachable device ",
			},
			"license_poolname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The registration key that you want to assign from the pool",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the rate at which this license usage is billed",
			},
			"unit_of_measure": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the rate at which this license usage is billed",
			},
			"skukeyword1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the rate at which this license usage is billed",
			},
			"skukeyword2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the rate at which this license usage is billed",
			},
			"hypervisor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Aws/Azure",
			},
			"tenant": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "optional description for the assignment in this field",
			},
			"device_license_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status of Licence Assignment",
			},
		},
	}
}

func resourceBigiqLicenseManageCreate(d *schema.ResourceData, meta interface{}) error {
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	bigipRef := meta.(*bigip.BigIP)
	var deviceIP []string
	var respID string
	deviceIP, _ = getDeviceUri(bigipRef.Host)
	devicePort, _ := strconv.Atoi(deviceIP[3])
	licensePoolName := d.Get("license_poolname").(string)
	poolInfo, err := bigiqRef.GetPoolType(licensePoolName)
	if err != nil {
		return err
	}
	if poolInfo == nil {
		return fmt.Errorf("there is no pool with specified name:%v", licensePoolName)
	}
	log.Printf("poolInfo:%+v", poolInfo)
	var licenseType string
	if poolInfo.SortName == "Registration Key Pool" {
		licenseType = poolInfo.SortName
	} else if poolInfo.SortName == "Utility" {
		licenseType = poolInfo.SortName
		if d.Get("unit_of_measure").(string) == "" {
			return fmt.Errorf("unit_of_measure is required parameter for %s licese type pool :%v", licenseType, licensePoolName)
		}
	}
	poolId, err := bigiqRef.GetRegkeyPoolId(licensePoolName)
	if err != nil {
		return fmt.Errorf("getting Poolid failed with :%v", err)
	}
	regKey := d.Get("key").(string)
	if regKey == "" {
		address := deviceIP[2]
		assignmentType := d.Get("assignment_type").(string)
		command := "assign"
		hyperVisor := d.Get("hypervisor").(string)
		macAddress := d.Get("mac_address").(string)
		skuKeyword1 := d.Get("skukeyword1").(string)
		skuKeyword2 := d.Get("skukeyword2").(string)
		tenant := d.Get("tenant").(string)
		unitOfMeasure := d.Get("unit_of_measure").(string)
		config := &bigip.LicenseParam{
			Address:         address,
			AssignmentType:  assignmentType,
			Command:         command,
			Hypervisor:      hyperVisor,
			LicensePoolName: licensePoolName,
			MacAddress:      macAddress,
			Password:        bigipRef.Password,
			SkuKeyword1:     skuKeyword1,
			SkuKeyword2:     skuKeyword2,
			Tenant:          tenant,
			UnitOfMeasure:   unitOfMeasure,
			User:            bigipRef.User,
		}
		taskID, err := bigiqRef.PostLicense(config)
		if err != nil {
			return fmt.Errorf("Error is : %v", err)
		}
		respID = taskID
	} else {
		assignmentType := d.Get("assignment_type").(string)
		if strings.ToUpper(assignmentType) == "MANAGED" {
			deviceID, err := bigiqRef.GetDeviceId(deviceIP[2])
			if (err != nil) && (deviceID == "") {
				return fmt.Errorf("getting deviceid failed with :%v", err)
			}
			deRef := bigip.DeviceRef{
				Link: deviceID,
			}
			config := &bigip.ManagedDevice{
				DeviceReference: deRef,
			}
			resp, err := bigiqRef.RegkeylicenseAssign(config, poolId, regKey)
			if err != nil {
				log.Printf("Assigning License failed from regKey Pool:%v", err)
				return err
			}
			respID = resp.ID
		} else if strings.ToUpper(assignmentType) == "UNMANAGED" {
			config := &bigip.UnmanagedDevice{
				DeviceAddress: deviceIP[2],
				Username:      bigipRef.User,
				Password:      bigipRef.Password,
				HTTPSPort:     devicePort,
			}
			//log.Printf("config2 = %+v", config)
			resp, err := bigiqRef.RegkeylicenseAssign(config, poolId, regKey)
			if err != nil {
				log.Printf("Assigning License failed from regKey Pool:%v", err)
				return err
			}
			//log.Printf("Resp from Post = %+v", resp)
			respID = resp.ID
		}
	}
	d.SetId(respID)
	return resourceBigiqLicenseManageRead(d, meta)
}

func resourceBigiqLicenseManageRead(d *schema.ResourceData, meta interface{}) error {
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	bigipRef := meta.(*bigip.BigIP)
	log.Printf("[DEBUG] bigipRef = %+v", bigipRef)
	memID := d.Id()
	//poolLicenseType := d.Get("pool_license_type").(string)
	poolName := d.Get("license_poolname").(string)
	regKey := d.Get("key").(string)
	poolId, err := bigiqRef.GetRegkeyPoolId(poolName)
	if err != nil && poolId == "" {
		log.Printf("Getting PoolID failed with :%v", err)
		return err
	}
	if regKey == "" {
		taskId := memID
		licenseStatus, err := bigiqRef.GetLicenseStatus(taskId)
		if err != nil {
			return fmt.Errorf("getting license status failed with : %v", err)
		}
		if licenseStatus["status"] == "FAILED" {
			d.SetId("")
			return fmt.Errorf("%s", licenseStatus["errorMessage"])
		}
		licenseAssignmentReference := licenseStatus["licenseAssignmentReference"].(map[string]interface{})["link"].(string)
		assignmentRef := strings.Split(licenseAssignmentReference, "/")
		deviceStatus, err := bigiqRef.GetDeviceLicenseStatus(assignmentRef[3:]...)
		d.Set("device_license_status", deviceStatus)
	} else {
		bigiqRef.GetMemberStatus(poolId, regKey, memID)
	}
	return nil
}

func resourceBigiqLicenseManageUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	return nil

}

func resourceBigiqLicenseManageDelete(d *schema.ResourceData, meta interface{}) error {
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	bigipRef := meta.(*bigip.BigIP)
	memID := d.Id()
	var poolId, regKey string
	if v, ok := d.GetOk("license_poolname"); ok {
		poolId, err = bigiqRef.GetRegkeyPoolId(v.(string))
		if (err != nil) && (poolId == "") {
			log.Printf("Getting PoolID failed with :%v", err)
			return err
		}
	}
	if v, ok := d.GetOk("key"); ok {
		regKey = v.(string)
	}
	var deviceIP []string
	deviceIP, _ = getDeviceUri(bigipRef.Host)
	devicePort, err := strconv.Atoi(deviceIP[3])
	if err == nil {
		return err
	}
	assignmentType := d.Get("assignment_type").(string)
	if regKey == "" {
		address := deviceIP[2]
		command := "revoke"
		hyperVisor := d.Get("hypervisor").(string)
		licensePoolName := d.Get("license_poolname").(string)
		macAddress := d.Get("mac_address").(string)
		skuKeyword1 := d.Get("skukeyword1").(string)
		skuKeyword2 := d.Get("skukeyword2").(string)
		tenant := d.Get("tenant").(string)
		unitOfMeasure := d.Get("unit_of_measure").(string)
		config := &bigip.LicenseParam{
			Address:         address,
			AssignmentType:  assignmentType,
			Command:         command,
			Hypervisor:      hyperVisor,
			LicensePoolName: licensePoolName,
			MacAddress:      macAddress,
			Password:        bigipRef.Password,
			SkuKeyword1:     skuKeyword1,
			SkuKeyword2:     skuKeyword2,
			Tenant:          tenant,
			UnitOfMeasure:   unitOfMeasure,
			User:            bigipRef.User,
		}
		_, err := bigiqRef.PostLicense(config)
		if err != nil {
			return fmt.Errorf("revoking license failed with : %v", err)
		}
	} else {
		if strings.ToUpper(assignmentType) == "MANAGED" {
			bigiqRef.RegkeylicenseRevoke(poolId, regKey, memID)
		} else if strings.ToUpper(assignmentType) == "UNMANAGED" {
			config := &struct {
				ID        string `json:"id"`
				Username  string `json:"username"`
				Password  string `json:"password"`
				HTTPSPort int    `json:"httpsPort,omitempty"`
			}{
				memID,
				bigipRef.User,
				bigipRef.Password,
				devicePort,
			}
			log.Printf("config = %+v", config)
			bigiqRef.LicenseRevoke(config, poolId, regKey, memID)
		}
	}
	d.SetId("")
	return nil
}

func connectBigIq(d *schema.ResourceData) (*bigip.BigIP, error) {
	bigiqConfig := Config{
		Address:  d.Get("bigiq_address").(string),
		Port:     d.Get("bigiq_port").(string),
		Username: d.Get("bigiq_user").(string),
		Password: d.Get("bigiq_password").(string),
	}
	if d.Get("bigiq_token_auth").(bool) {
		bigiqConfig.LoginReference = d.Get("bigiq_login_ref").(string)
	}
	return bigiqConfig.Client()
}

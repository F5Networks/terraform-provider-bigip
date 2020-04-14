package bigip

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	uriRegkey      = "regkey"
	uriLicenses    = "licenses"
	uriResolver    = "resolver"
	uriDevicegroup = "device-groups"
	uriCmBigip     = "cm-bigip-allBigIpDevices"
	uriDevice      = "device"
	uriMembers     = "members"
)

type DeviceRef struct {
	Link string `json:"link"`
}
type ManagedDevice struct {
	DeviceReference DeviceRef `json:"deviceReference"`
}

type unmanagedDevice struct {
	DeviceAddress string `json:"deviceAddress"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	HTTPSPort     int    `json:"httpsPort"`
}

type regKeyPools struct {
	//Items      []struct {
	//	ID       string `json:"id"`
	//	Name     string `json:"name"`
	//	SortName string `json:"sortName"`
	//} `json:"items"`
	RegKeyPoollist []regKeyPool `json:"items"`
}

type regKeyPool struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SortName string `json:"sortName"`
}

type devicesList struct {
	DevicesInfo []deviceInfo `json:"items"`
}
type deviceInfo struct {
	Address           string `json:"address"`
	DeviceURI         string `json:"deviceUri"`
	Hostname          string `json:"hostname"`
	HTTPSPort         int    `json:"httpsPort"`
	IsClustered       bool   `json:"isClustered"`
	MachineID         string `json:"machineId"`
	ManagementAddress string `json:"managementAddress"`
	McpDeviceName     string `json:"mcpDeviceName"`
	Product           string `json:"product"`
	SelfLink          string `json:"selfLink"`
	State             string `json:"state"`
	UUID              string `json:"uuid"`
	Version           string `json:"version"`
}

type MembersList struct {
	Members []memberDetail `json:"items"`
}

type memberDetail struct {
	AssignmentType  string `json:"assignmentType"`
	DeviceAddress   string `json:"deviceAddress"`
	DeviceMachineID string `json:"deviceMachineId"`
	DeviceName      string `json:"deviceName"`
	ID              string `json:"id"`
	Message         string `json:"message"`
	Status          string `json:"status"`
}

type regKeyAssignStatus struct {
	ID             string `json:"id"`
	DeviceAddress  string `json:"deviceAddress"`
	AssignmentType string `json:"assignmentType"`
	DeviceName     string `json:"deviceName"`
	Status         string `json:"status"`
}

func (b *BigIP) GetRegPools() (*regKeyPools, error) {
	var self regKeyPools
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIP) GetRegPool(poolid string) (*regKeyPool, error) {
	var self regKeyPool
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolid)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIP) GetManagedDevices() (*devicesList, error) {
	var self devicesList
	err, _ := b.getForEntity(&self, uriMgmt, uriShared, uriResolver, uriDevicegroup, uriCmBigip, uriDevices)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIP) GetRegkeyPoolId(poolName string) (string, error) {
	var self regKeyPools
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return "", err
	}
	for _, pool := range self.RegKeyPoollist {
		if pool.Name == poolName {
			return pool.ID, nil
		}
	}
	return "", nil
}

func (b *BigIP) RegkeylicenseAssign(config interface{}, poolId string, regKey string) (*memberDetail, error) {
	resp, err := b.postReq(config, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers)
	if err != nil {
		return nil, err
	}
	var resp1 regKeyAssignStatus
	err = json.Unmarshal(resp, &resp1)
	if err != nil {
		return nil, err
	}
	return b.GetMemberStatus(poolId, regKey, resp1.ID)
}

func (b *BigIP) GetMemberStatus(poolId, regKey, memId string) (*memberDetail, error) {
	var self memberDetail
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return nil, err
	}
	for self.Status != "LICENSED" {
		log.Printf("Member status:%+v", self.Status)
		if self.Status == "INSTALLATION_FAILED" {
			return &self, fmt.Errorf("INSTALLATION_FAILED with %s", self.Message)
		}
		return b.GetMemberStatus(poolId, regKey, memId)
	}
	return &self, nil
}
func (b *BigIP) RegkeylicenseRevoke(poolId, regKey, memId string) error {
	log.Printf("Deleting License for Member:%+v", memId)
	_, err := b.deleteReq(uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	r1 := make(map[string]interface{})
	err, _ = b.getForEntity(&r1, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	log.Printf("Response after delete:%+v", r1)
	return nil
}

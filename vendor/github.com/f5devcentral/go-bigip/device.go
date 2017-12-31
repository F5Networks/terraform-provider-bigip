package bigip

import (
	"encoding/json"
	"log"
)

//  LIC contains device license for BIG-IP system.
type LICs struct {
	LIC []LIC `json:"items"`
}

// VirtualAddress contains information about each individual virtual address.
type LIC struct {
	DeviceAddress string
	Username      string
	Password      string
}

type LicensePools struct {
	LicensePool []LicensePool `json:"items"`
}

type LicensePool struct {
	Items []struct {
		Uuid string `json:"Uuid,omitempty"`
	}
}

type LICDTO struct {
	DeviceAddress string `json:"deviceAddress,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
}

type Devicenames struct {
	Devicenames []Devicename `json:"items"`
}

type Devicename struct {
	Command string `json:"command,omitempty"`
	Name    string `json:"name,omitempty"`
	Target  string `json:"target,omitempty"`
}

type Devices struct {
	Devices []Device `json:"items"`
}

type Device struct {
	ConfigsyncIp      string `json:"configsyncIp,omitempty"`
	Name              string `json:"name,omitempty"`
	MirrorIp          string `json:"mirrorIp,omitempty"`
	MirrorSecondaryIp string `json:"mirrorSecondaryIp,omitempty"`
}

type Devicegroups struct {
	Devicegroups []Devicegroup `json:"items"`
}

type Devicegroup struct {
	AutoSync                     string
	Name                         string
	Partition                    string
	Description                  string
	Type                         string
	FullLoadOnSync               string
	SaveOnAutoSync               string
	NetworkFailover              string
	IncrementalConfigSyncSizeMax int
	Deviceb                      []Devicerecord
}
type devicegroupDTO struct {
	AutoSync                     string `json:"autoSync,omitempty"`
	Name                         string `json:"name,omitempty"`
	Partition                    string `json:"partition,omitempty"`
	Description                  string `json:"description,omitempty"`
	Type                         string `json:"type,omitempty"`
	FullLoadOnSync               string `json:"fullLoadOnSync,omitempty"`
	SaveOnAutoSync               string `json:"saveOnAutoSync,omitempty"`
	NetworkFailover              string `json:"networkFailover,omitempty"`
	IncrementalConfigSyncSizeMax int    `json:"incrementalConfigSyncSizeMax,omitempty"`
	Deviceb                      struct {
		Items []Devicerecord `json:"items,omitempty"`
	} `json:"deviceReference,omitempty"`
}

type Devicerecords struct {
	Items []Originsrecord `json:"items,omitempty"`
}

type Devicerecord struct {
	SetSyncLeader bool   `json:"setSyncLeader"`
	Name          string `json:"name"`
}

func (p *Devicegroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(devicegroupDTO{
		Name:                         p.Name,
		Partition:                    p.Partition,
		AutoSync:                     p.AutoSync,
		Description:                  p.Description,
		Type:                         p.Type,
		FullLoadOnSync:               p.FullLoadOnSync,
		SaveOnAutoSync:               p.SaveOnAutoSync,
		NetworkFailover:              p.NetworkFailover,
		IncrementalConfigSyncSizeMax: p.IncrementalConfigSyncSizeMax,
		Deviceb: struct {
			Items []Devicerecord `json:"items,omitempty"`
		}{Items: p.Deviceb},
	})
}

func (p *Devicegroup) UnmarshalJSON(b []byte) error {
	var dto devicegroupDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	p.Name = dto.Name
	p.Partition = dto.Partition
	p.AutoSync = dto.AutoSync
	p.Description = dto.Description
	p.Type = dto.Type
	p.FullLoadOnSync = dto.FullLoadOnSync
	p.SaveOnAutoSync = dto.SaveOnAutoSync
	p.NetworkFailover = dto.NetworkFailover
	p.IncrementalConfigSyncSizeMax = dto.IncrementalConfigSyncSizeMax
	p.Deviceb = dto.Deviceb.Items

	return nil
}

// https://10.192.74.80/mgmt/cm/device/licensing/pool/purchased-pool/licenses
// The above command will spit out license uuid and which should be mapped uriUuid
const (
	uriMgmt          = "mgmt"
	uriCm            = "cm"
	uriDiv           = "device"
	uriDG            = "device-group"
	uriLins          = "licensing"
	uriPoo           = "pool"
	uriPur           = "purchased-pool"
	uriLicn          = "licenses"
	uriMemb          = "members"
	uriUtility       = "utility"
	uriOfferings     = "offerings"
	uriF5BIGMSPBT10G = "f37c66e0-a80d-43e8-924b-3bbe9fe96bbe"
)

func (p *LIC) MarshalJSON() ([]byte, error) {
	var dto LICDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *LIC) UnmarshalJSON(b []byte) error {
	var dto LICDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (b *BigIP) getLicensePool() (*LicensePool, error) {
	var licensePool LicensePool
	err, _ := b.getForEntity(&licensePool, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn)
	if err != nil {
		return nil, err
	}
	// for loop over all returned license pools to check which one has available licenses
	// getAvailablePool(member[index_of_array].Uuid)
	// At the end change return statement to return only the UUID string of the one where license
	// is available
	return &licensePool, nil
}

// VirtualAddresses returns a list of virtual addresses.
func (b *BigIP) LIC() (*LIC, error) {
	var va LIC
	licensePool, licensePoolErr := b.getLicensePool()
	if licensePoolErr != nil {
		return nil, licensePoolErr
	}
	err, _ := b.getForEntity(&va, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn, licensePool.Items[0].Uuid, uriMemb)
	if err != nil {
		return nil, err
	}
	return &va, nil
}

func (b *BigIP) CreateLIC(deviceAddress string, username string, password string) error {
	config := &LIC{
		DeviceAddress: deviceAddress,
		Username:      username,
		Password:      password,
	}

	licensePool, licensePoolErr := b.getLicensePool()
	if licensePoolErr != nil {
		return licensePoolErr
	}

	return b.post(config, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn, licensePool.Items[0].Uuid, uriMemb)
}

func (b *BigIP) ModifyLIC(config *LIC) error {
	licensePool, licensePoolErr := b.getLicensePool()
	if licensePoolErr != nil {
		return licensePoolErr
	}
	return b.post(config, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn, licensePool.Items[0].Uuid, uriMemb)
}

func (b *BigIP) LICs() (*LIC, error) {
	var members LIC
	licensePool, licensePoolErr := b.getLicensePool()
	if licensePoolErr != nil {
		return nil, licensePoolErr
	}
	err, _ := b.getForEntity(&members, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn, licensePool.Items[0].Uuid, uriMemb)

	if err != nil {
		return nil, err
	}

	return &members, nil
}

func (b *BigIP) CreateDevice(name, configsyncIp, mirrorIp, mirrorSecondaryIp string) error {
	config := &Device{
		Name:              name,
		ConfigsyncIp:      configsyncIp,
		MirrorIp:          mirrorIp,
		MirrorSecondaryIp: mirrorSecondaryIp,
	}

	return b.post(config, uriCm, uriDiv)
}

func (b *BigIP) ModifyDevice(config *Device) error {
	return b.put(config, uriCm, uriDiv)
}

func (b *BigIP) DeleteDevice(name string) error {
	return b.delete(uriCm, uriDiv, name)
}

func (b *BigIP) Devices(name string) (*Device, error) {
	var device Device
	err, _ := b.getForEntity(&device, uriCm, uriDiv, name)

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (b *BigIP) CreateDevicegroup(p *Devicegroup) error {
	log.Println(" what is the complete payload    ", p)
	return b.post(p, uriCm, uriDG)
}

func (b *BigIP) UpdateDevicegroup(name string, p *Devicegroup) error {
	return b.put(p, uriCm, uriDG, name)
}
func (b *BigIP) ModifyDevicegroup(config *Devicegroup) error {
	return b.put(config, uriCm, uriDG)
}

func (b *BigIP) Devicegroups(name string) (*Devicegroup, error) {
	var devicegroup Devicegroup
	err, _ := b.getForEntity(&devicegroup, uriCm, uriDiv, name)

	if err != nil {
		return nil, err
	}

	return &devicegroup, nil
}

/* func (b *BigIP) Devicegroups() (*Devicegroup, error) {
	var devicegroup Devicegroup
	err, _ := b.getForEntity(&devicegroup, uriCm, uriDG)

	if err != nil {
		return nil, err
	}
	return &devicegroup, nil
}
*/
func (b *BigIP) DeleteDevicegroup(name string) error {
	return b.delete(uriCm, uriDG, name)
}

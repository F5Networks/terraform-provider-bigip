package bigip

import "encoding/json"

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

// https://10.192.74.80/mgmt/cm/device/licensing/pool/purchased-pool/licenses
// The above command will spit out license uuid and which should be mapped uriUuid
const (
	uriMgmt = "mgmt"
	uriCm   = "cm"
	uriDiv  = "device"
	uriLins = "licensing"
	uriPoo  = "pool"
	uriPur  = "purchased-pool"
	uriLicn = "licenses"
	uriMemb = "members"
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
	// is availble
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

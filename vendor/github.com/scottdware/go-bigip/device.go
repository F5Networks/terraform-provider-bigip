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

type LICDTO struct {
	DeviceAddress string `json:"deviceAddress,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
}

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

const (
	uriMgmt = "mgmt"
	uriCm   = "cm"
	uriDiv  = "device"
	uriLins = "licensing"
	uriPoo  = "pool"
	uriPur  = "purchased-pool"
	uriLicn = "licenses"
	uriUuid = "e0a94ea6-e859-4bec-961d-261c91ef85ad"
	uriMemb = "members"
)

// VirtualAddresses returns a list of virtual addresses.
func (b *BigIP) LIC() (*LIC, error) {
	var va LIC
	err, _ := b.getForEntity(&va, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn, uriUuid, uriMemb)
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

	return b.post(config, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn, uriUuid, uriMemb)
}

func (b *BigIP) ModifyLIC(config *LIC) error {
	return b.post(config, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn, uriUuid, uriMemb)
}

func (b *BigIP) LICs() (*LIC, error) {
	var members LIC
	err, _ := b.getForEntity(&members, uriMgmt, uriCm, uriDiv, uriLins, uriPoo, uriPur, uriLicn, uriUuid, uriMemb)

	if err != nil {
		return nil, err
	}

	return &members, nil
}

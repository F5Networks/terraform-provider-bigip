package bigip

import "encoding/json"
import "log"



type Datacenters struct {
	Datacenters []Datacenter `json:"items"`
}

type Datacenter struct {
	Name  string
	App_service  string
	Description  string
	Disabled bool
	Enabled bool
	Prober_pool  string
}

type datacenterDTO struct {
	Name  string `json:"name,omitempty"`
	App_service  string `json:"appService,omitempty"`
	Description string `json:"description,omitempty"`
	Disabled  bool `json:"disabled,omitempty"`
	Enabled  bool `json:"enabled,omitempty"`
	Prober_pool  string `json:"proberPool,omitempty"`
}

func (p *Datacenter) MarshalJSON() ([]byte, error) {
	var dto datacenterDTO
	return json.Marshal(dto)
}

func (p *Datacenter) UnmarshalJSON(b []byte) error {
	var dto datacenterDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return nil
}

const (
	uriGtm       = "gtm"
	uriDatacenter = "datacenter"
)

func (b *BigIP) Datacenters() (*Datacenter, error) {
	var datacenter Datacenter
	err, _ := b.getForEntity(&datacenter, uriGtm, uriDatacenter)

	if err != nil {
		return nil, err
	}

	return &datacenter, nil
}

func (b *BigIP) CreateDatacenter(name, description, app_service string, enabled, disabled bool, prober_pool string) error {
	config := &Datacenter{
		Name:    name,
		Description:  description,
		App_service: app_service,
		Enabled: enabled,
		Disabled: disabled,
		Prober_pool: prober_pool,
	}
	log.Printf("I am %#v\n  here", config)
	return b.post(config, uriGtm, uriDatacenter)
}

func (b *BigIP) ModifyDatacenter(*Datacenter) error {
	return b.patch(uriGtm, uriDatacenter)
}


func (b *BigIP) DeleteDatacenter(name string) error {
	return b.delete(uriGtm, uriDatacenter, name)
}

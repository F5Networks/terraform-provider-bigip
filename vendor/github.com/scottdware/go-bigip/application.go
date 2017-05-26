package bigip

import "encoding/json"

//  LIC contains device license for BIG-IP system.
type Services struct {
	Service []Service `json:"items"`
}

// VirtualAddress contains information about each individual virtual address.
type Service struct {
	Description                string
	DeviceGroup                string
	ExecuteAction              string
	InheritedDevicegroup       string
	InheritedTrafficGroup      string
	TmPartition                string
	StrictUpdates              string
	Template                   string
	TemplateModified           string
	TemplatePrerequisiteErrors string
	TrafficGroup               string
	Variables                  []string
	Encrypted                  string
	Value                      string
	Tables                     []string
	ColumnNames                string
	EncryptedColumns           string
	Rows                       string
}

type ServiceDTO struct {
	Description                string   `json:"description,omitempty"`
	DeviceGroup                string   `json:"deviceGroup,omitempty"`
	ExecuteAction              string   `json:"executeAction,omitempty"`
	InheritedDevicegroup       string   `json:"inheritedDevicegroup,omitempty"`
	InheritedTrafficGroup      string   `json:"inheritedTrafficGroup,omitempty"`
	TmPartition                string   `json:"tmPartition,omitempty"`
	StrictUpdates              string   `json:"strictUpdates,omitempty"`
	Template                   string   `json:"template,omitempty"`
	TemplateModified           string   `json:"templateModified,omitempty"`
	TemplatePrerequisiteErrors string   `json:"templatePrerequisiteErrors,omitempty"`
	TrafficGroup               string   `json:"trafficGroup,omitempty"`
	Variables                  []string `json:"variables,omitempty"`
	Encrypted                  string   `json:"encrypted,omitempty"`
	Value                      string   `json:"value,omitempty"`
	Tables                     []string `json:"tables,omitempty"`
	ColumnNames                string   `json:"columnNames,omitempty"`
	EncryptedColumns           string   `json:"encryptedColumns,omitempty"`
	Rows                       string   `json:"rows,omitempty"`
}

func (p *Service) MarshalJSON() ([]byte, error) {
	var dto ServiceDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Service) UnmarshalJSON(b []byte) error {
	var dto ServiceDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

const (
	uriSyss = "sys"
	uriApp  = "application"
	uriServ = "service"
)

//  returns a list of iApps
func (b *BigIP) Service() (*Service, error) {
	var va Service
	err, _ := b.getForEntity(&va, uriSyss, uriApp, uriServ)
	if err != nil {
		return nil, err
	}
	return &va, nil
}

func (b *BigIP) CreateService(description string, deviceGroup string, executeAction string, inheritedDevicegroup string, inheritedTrafficGroup string, tmPartition string, strictUpdates string, template string, templateModified string, templatePrerequisiteErrors string, trafficGroup string, tables []string, columnNames string, encryptedColumns string, rows string, variables []string, encrypted string, value string) error {
	config := &Service{
		Description:                description,
		DeviceGroup:                deviceGroup,
		ExecuteAction:              executeAction,
		InheritedDevicegroup:       inheritedDevicegroup,
		InheritedTrafficGroup:      inheritedTrafficGroup,
		TmPartition:                tmPartition,
		StrictUpdates:              strictUpdates,
		Template:                   template,
		TemplateModified:           templateModified,
		TemplatePrerequisiteErrors: templatePrerequisiteErrors,
		TrafficGroup:               trafficGroup,
		Variables:                  variables,
		Encrypted:                  encrypted,
		Value:                      value,
		Tables:                     tables,
		ColumnNames:                columnNames,
		EncryptedColumns:           encryptedColumns,
		Rows:                       rows,
	}

	return b.post(config, uriSyss, uriApp, uriServ)
}

func (b *BigIP) ModifyService(config *Service) error {
	return b.post(config, uriSyss, uriApp, uriServ)
}

func (b *BigIP) Services() (*Service, error) {
	var members Service
	err, _ := b.getForEntity(&members, uriSyss, uriApp, uriServ)

	if err != nil {
		return nil, err
	}

	return &members, nil
}

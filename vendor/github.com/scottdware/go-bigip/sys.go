package bigip

import "encoding/json"

type NTPs struct {
	NTPs []NTP `json:"items"`
}

type NTP struct {
	Description string   `json:"description,omitempty"`
	Servers     []string `json:"servers,omitempty"`
	Timezone    string   `json:"timezone,omitempty"`
}

type DNSs struct {
	DNSs []DNS `json:"items"`
}

type DNS struct {
	Description  string   `json:"description,omitempty"`
	NameServers  []string `json:"nameServers,omitempty"`
	NumberOfDots int      `json:"numberOfDots,omitempty"`
	Search       []string `json:"search,omitempty"`
}

type Provisions struct {
	Provisions []Provision `json:"items"`
}

type Provision struct {
	Name        string `json:"name,omitempty"`
	FullPath    string `json:"fullPath,omitempty"`
	CpuRatio    int    `json:"cpuRatio,omitempty"`
	DiskRatio   int    `json:"diskRatio,omitempty"`
	Level       string `json:"level,omitempty"`
	MemoryRatio int    `json:"memoryRatio,omitempty"`
}

type RemoteServer struct {
	Name       string `json:"name,omitempty"`
	Host       string `json:"host,omitempty"`
	RemotePort int    `json:"remotePort,omitempty"`
}

type Syslogs struct {
	Syslogs []Syslog `json:"items"`
}

type Syslog struct {
	AuthPrivFrom  string
	RemoteServers []RemoteServer
}

type SyslogDTO struct {
	AuthPrivFrom  string `json:"authPrivFrom,omitempty"`
	RemoteServers struct {
		Items []RemoteServer `json:"items,omitempty"`
	} `json:"remoteServers,omitempty"`
}

type SNMPs struct {
	SNMPs []SNMP `json:"items"`
}

type SNMP struct {
	SysContact       string   `json:"sysContact,omitempty"`
	SysLocation      string   `json:"sysLocation,omitempty"`
	AllowedAddresses []string `json:"allowedAddresses,omitempty"`
}

type TRAPs struct {
	SNMPs []SNMP `json:"items"`
}

type TRAP struct {
	Name                     string `json:"name,omitempty"`
	AuthPasswordEncrypted    string `json:"authPasswordEncrypted,omitempty"`
	AuthProtocol             string `json:"authProtocol,omitempty"`
	Community                string `json:"community,omitempty"`
	Description              string `json:"description,omitempty"`
	EngineId                 string `json:"engineId,omitempty"`
	Host                     string `json:"host,omitempty"`
	Port                     int    `json:"port,omitempty"`
	PrivacyPassword          string `json:"privacyPassword,omitempty"`
	PrivacyPasswordEncrypted string `json:"privacyPasswordEncrypted,omitempty"`
	PrivacyProtocol          string `json:"privacyProtocol,omitempty"`
	SecurityLevel            string `json:"securityLevel,omitempty"`
	SecurityName             string `json:"SecurityName,omitempty"`
	Version                  string `json:"version,omitempty"`
}

const (
	uriSys       = "sys"
	uriNtp       = "ntp"
	uriDNS       = "dns"
	uriProvision = "provision"
	uriAfm       = "afm"
	uriAsm       = "asm"
	uriApm       = "apm"
	uriGtm       = "gtm"
	uriAvr       = "avr"
	uriIlx       = "ilx"
	uriSyslog    = "syslog"
	uriSnmp      = "snmp"
	uriTraps     = "traps"
)

func (p *Syslog) MarshalJSON() ([]byte, error) {
	var dto SyslogDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Syslog) UnmarshalJSON(b []byte) error {
	var dto SyslogDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (b *BigIP) CreateNTP(description string, servers []string, timezone string) error {
	config := &NTP{
		Description: description,
		Servers:     servers,
		Timezone:    timezone,
	}

	return b.patch(config, uriSys, uriNtp)
}

func (b *BigIP) ModifyNTP(config *NTP) error {
	return b.put(config, uriSys, uriNtp)
}

func (b *BigIP) NTPs() (*NTP, error) {
	var ntp NTP
	err, _ := b.getForEntity(&ntp, uriSys, uriNtp)

	if err != nil {
		return nil, err
	}

	return &ntp, nil
}

func (b *BigIP) CreateDNS(description string, nameservers []string, numberofdots int, search []string) error {
	config := &DNS{
		Description:  description,
		NameServers:  nameservers,
		NumberOfDots: numberofdots,
		Search:       search,
	}

	return b.patch(config, uriSys, uriDNS)
}

func (b *BigIP) ModifyDNS(config *DNS) error {
	return b.put(config, uriSys, uriDNS)
}

func (b *BigIP) DNSs() (*DNS, error) {
	var dns DNS
	err, _ := b.getForEntity(&dns, uriSys, uriDNS)

	if err != nil {
		return nil, err
	}

	return &dns, nil
}

func (b *BigIP) CreateProvision(name string, fullPath string, cpuRatio int, diskRatio int, level string, memoryRatio int) error {
	config := &Provision{
		Name:        name,
		FullPath:    fullPath,
		CpuRatio:    cpuRatio,
		DiskRatio:   diskRatio,
		Level:       level,
		MemoryRatio: memoryRatio,
	}
	if name == "/Common/asm" {
		return b.put(config, uriSys, uriProvision, uriAsm)
	}
	if name == "/Common/afm" {
		return b.put(config, uriSys, uriProvision, uriAfm)
	}
	if name == "/Common/gtm" {
		return b.put(config, uriSys, uriProvision, uriGtm)
	}

	if name == "/Common/apm" {
		return b.put(config, uriSys, uriProvision, uriApm)
	}

	if name == "/Common/avr" {
		return b.put(config, uriSys, uriProvision, uriAvr)
	}
	if name == "/Common/ilx" {
		return b.put(config, uriSys, uriProvision, uriIlx)
	}
	return nil
}

func (b *BigIP) ModifyProvision(config *Provision) error {

	return b.put(config, uriSys, uriProvision, uriAfm)
}

func (b *BigIP) DeleteProvision(name string) error {
	return b.delete(uriSys, uriProvision, uriIlx, name)
}

func (b *BigIP) Provisions() (*Provision, error) {
	var provision Provision
	err, _ := b.getForEntity(&provision, uriProvision, uriAfm)

	if err != nil {
		return nil, err
	}

	return &provision, nil
}

func (b *BigIP) Syslogs() (*Syslog, error) {
	var syslog Syslog
	err, _ := b.getForEntity(&syslog, uriSys, uriSyslog)

	if err != nil {
		return nil, err
	}

	return &syslog, nil
}

func (b *BigIP) CreateSyslog(authPrivFrom string, remoteServers []RemoteServer) error {
	config := &Syslog{
		AuthPrivFrom:  authPrivFrom,
		RemoteServers: remoteServers,
	}

	return b.patch(config, uriSys, uriSyslog)
}

func (b *BigIP) CreateSNMP(sysContact string, sysLocation string, allowedAddresses []string) error {
	config := &SNMP{
		SysContact:       sysContact,
		SysLocation:      sysLocation,
		AllowedAddresses: allowedAddresses,
	}

	return b.patch(config, uriSys, uriSnmp)
}

func (b *BigIP) ModifySNMP(config *SNMP) error {
	return b.put(config, uriSys, uriSnmp)
}

func (b *BigIP) SNMPs() (*SNMP, error) {
	var snmp SNMP
	err, _ := b.getForEntity(&snmp, uriSys, uriSnmp)

	if err != nil {
		return nil, err
	}

	return &snmp, nil
}

func (b *BigIP) CreateTRAP(name string, authPasswordEncrypted string, authProtocol string, community string, description string, engineId string, host string, port int, privacyPassword string, privacyPasswordEncrypted string, privacyProtocol string, securityLevel string, securityName string, version string) error {
	config := &TRAP{
		Name: name,
		AuthPasswordEncrypted:    authPasswordEncrypted,
		AuthProtocol:             authProtocol,
		Community:                community,
		Description:              description,
		EngineId:                 engineId,
		Host:                     host,
		Port:                     port,
		PrivacyPassword:          privacyPassword,
		PrivacyPasswordEncrypted: privacyPasswordEncrypted,
		PrivacyProtocol:          privacyProtocol,
		SecurityLevel:            securityLevel,
		SecurityName:             securityName,
		Version:                  version,
	}

	return b.post(config, uriSys, uriSnmp, uriTraps)
}

func (b *BigIP) ModifyTRAP(config *TRAP) error {
	return b.put(config, uriSys, uriSnmp, uriTraps)
}

func (b *BigIP) TRAPs() (*TRAP, error) {
	var traps TRAP
	err, _ := b.getForEntity(&traps, uriSys, uriSnmp, uriTraps)

	if err != nil {
		return nil, err
	}

	return &traps, nil
}

func (b *BigIP) DeleteTRAP(name string) error {
	return b.delete(uriSys, uriSnmp, uriTraps, name)
}

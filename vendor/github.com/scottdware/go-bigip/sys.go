package bigip

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
)

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

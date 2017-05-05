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
	uriSys = "sys"
	uriNtp = "ntp"
	uriDNS = "dns"
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

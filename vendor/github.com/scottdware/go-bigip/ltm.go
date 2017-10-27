package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Nodes contains a list of every node on the BIG-IP system.
type Nodes struct {
	Nodes []Node `json:"items"`
}

// Node contains information about each individual node. You can use all
// of these fields when modifying a node.
type Node struct {
	Name            string `json:"name,omitempty"`
	Partition       string `json:"partition,omitempty"`
	FullPath        string `json:"fullPath,omitempty"`
	Generation      int    `json:"generation,omitempty"`
	Address         string `json:"address,omitempty"`
	ConnectionLimit int    `json:"connectionLimit,omitempty"`
	DynamicRatio    int    `json:"dynamicRatio,omitempty"`
	Logging         string `json:"logging,omitempty"`
	Monitor         string `json:"monitor,omitempty"`
	RateLimit       string `json:"rateLimit,omitempty"`
	Ratio           int    `json:"ratio,omitempty"`
	Session         string `json:"session,omitempty"`
	State           string `json:"state,omitempty"`
}

// Pools contains a list of pools on the BIG-IP system.
type Pools struct {
	Pools []Pool `json:"items"`
}

// Pool contains information about each pool. You can use all of these
// fields when modifying a pool.
type Pool struct {
	Name                   string
	Partition              string
	FullPath               string
	Generation             int
	AllowNAT               bool
	AllowSNAT              bool
	IgnorePersistedWeight  bool
	IPTOSToClient          string
	IPTOSToServer          string
	LinkQoSToClient        string
	LinkQoSToServer        string
	LoadBalancingMode      string
	MinActiveMembers       int
	MinUpMembers           int
	MinUpMembersAction     string
	MinUpMembersChecking   string
	Monitor                string
	QueueDepthLimit        int
	QueueOnConnectionLimit string
	QueueTimeLimit         int
	ReselectTries          int
	SlowRampTime           int
}

// Pool transfer object so we can mask the bool data munging
type poolDTO struct {
	Name                   string `json:"name,omitempty"`
	Partition              string `json:"partition,omitempty"`
	FullPath               string `json:"fullPath,omitempty"`
	Generation             int    `json:"generation,omitempty"`
	AllowNAT               string `json:"allowNat,omitempty" bool:"yes"`
	AllowSNAT              string `json:"allowSnat,omitempty" bool:"yes"`
	IgnorePersistedWeight  string `json:"ignorePersistedWeight,omitempty" bool:"enabled"`
	IPTOSToClient          string `json:"ipTosToClient,omitempty"`
	IPTOSToServer          string `json:"ipTosToServer,omitempty"`
	LinkQoSToClient        string `json:"linkQosToClient,omitempty"`
	LinkQoSToServer        string `json:"linkQosToServer,omitempty"`
	LoadBalancingMode      string `json:"loadBalancingMode,omitempty"`
	MinActiveMembers       int    `json:"minActiveMembers,omitempty"`
	MinUpMembers           int    `json:"minUpMembers,omitempty"`
	MinUpMembersAction     string `json:"minUpMembersAction,omitempty"`
	MinUpMembersChecking   string `json:"minUpMembersChecking,omitempty"`
	Monitor                string `json:"monitor"`
	QueueDepthLimit        int    `json:"queueDepthLimit,omitempty"`
	QueueOnConnectionLimit string `json:"queueOnConnectionLimit,omitempty"`
	QueueTimeLimit         int    `json:"queueTimeLimit,omitempty"`
	ReselectTries          int    `json:"reselectTries,omitempty"`
	SlowRampTime           int    `json:"slowRampTime,omitempty"`
}

func (p *Pool) MarshalJSON() ([]byte, error) {
	var dto poolDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Pool) UnmarshalJSON(b []byte) error {
	var dto poolDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

// poolMember is used only when adding members to a pool.
type poolMember struct {
	Name string `json:"name"`
}

// VirtualServers contains a list of all virtual servers on the BIG-IP system.
type VirtualServers struct {
	VirtualServers []VirtualServer `json:"items"`
}

// VirtualServer contains information about each individual virtual server.
type VirtualServer struct {
	Name                     string `json:"name,omitempty"`
	Partition                string `json:"partition,omitempty"`
	FullPath                 string `json:"fullPath,omitempty"`
	Generation               int    `json:"generation,omitempty"`
	AddressStatus            string `json:"addressStatus,omitempty"`
	AutoLastHop              string `json:"autoLastHop,omitempty"`
	CMPEnabled               string `json:"cmpEnabled,omitempty"`
	ConnectionLimit          int    `json:"connectionLimit,omitempty"`
	Destination              string `json:"destination,omitempty"`
	Enabled                  bool   `json:"enabled,omitempty"`
	GTMScore                 int    `json:"gtmScore,omitempty"`
	IPProtocol               string `json:"ipProtocol,omitempty"`
	Mask                     string `json:"mask,omitempty"`
	Mirror                   string `json:"mirror,omitempty"`
	MobileAppTunnel          string `json:"mobileAppTunnel,omitempty"`
	NAT64                    string `json:"nat64,omitempty"`
	Pool                     string `json:"pool,omitempty"`
	RateLimit                string `json:"rateLimit,omitempty"`
	RateLimitDestinationMask int    `json:"rateLimitDstMask,omitempty"`
	RateLimitMode            string `json:"rateLimitMode,omitempty"`
	RateLimitSourceMask      int    `json:"rateLimitSrcMask,omitempty"`
	Source                   string `json:"source,omitempty"`
	SourceAddressTranslation struct {
		Type string `json:"type,omitempty"`
		Pool string `json:"pool,omitempty"`
	} `json:"sourceAddressTranslation,omitempty"`
	SourcePort       string    `json:"sourcePort,omitempty"`
	SYNCookieStatus  string    `json:"synCookieStatus,omitempty"`
	TranslateAddress string    `json:"translateAddress,omitempty"`
	TranslatePort    string    `json:"translatePort,omitempty"`
	VlansDisabled    bool      `json:"vlansDisabled,omitempty"`
	VSIndex          int       `json:"vsIndex,omitempty"`
	Rules            []string  `json:"rules,omitempty"`
	Profiles         []Profile `json:"profiles,omitempty"`
	Policies         []string  `json:"policies,omitempty"`
}

// VirtualAddresses contains a list of all virtual addresses on the BIG-IP system.
type VirtualAddresses struct {
	VirtualAddresses []VirtualAddress `json:"items"`
}

// VirtualAddress contains information about each individual virtual address.
type VirtualAddress struct {
	Name                  string
	Partition             string
	FullPath              string
	Generation            int
	Address               string
	ARP                   bool
	AutoDelete            bool
	ConnectionLimit       int
	Enabled               bool
	Floating              bool
	ICMPEcho              bool
	InheritedTrafficGroup bool
	Mask                  string
	RouteAdvertisement    bool
	ServerScope           string
	TrafficGroup          string
	Unit                  int
}

type virtualAddressDTO struct {
	Name                  string `json:"name"`
	Partition             string `json:"partition,omitempty"`
	FullPath              string `json:"fullPath,omitempty"`
	Generation            int    `json:"generation,omitempty"`
	Address               string `json:"address,omitempty"`
	ARP                   string `json:"arp,omitempty" bool:"enabled"`
	AutoDelete            string `json:"autoDelete,omitempty" bool:"true"`
	ConnectionLimit       int    `json:"connectionLimit,omitempty"`
	Enabled               string `json:"enabled,omitempty" bool:"yes"`
	Floating              string `json:"floating,omitempty" bool:"enabled"`
	ICMPEcho              string `json:"icmpEcho,omitempty" bool:"enabled"`
	InheritedTrafficGroup string `json:"inheritedTrafficGroup,omitempty" bool:"yes"`
	Mask                  string `json:"mask,omitempty"`
	RouteAdvertisement    string `json:"routeAdvertisement,omitempty" bool:"enabled"`
	ServerScope           string `json:"serverScope,omitempty"`
	TrafficGroup          string `json:"trafficGroup,omitempty"`
	Unit                  int    `json:"unit,omitempty"`
}

type Policies struct {
	Policies []Policy `json:"items"`
}

type VirtualServerPolicies struct {
	PolicyRef Policies `json:"policiesReference"`
}

type Policy struct {
	Name      string
	Partition string
	FullPath  string
	Controls  []string
	Requires  []string
	Strategy  string
	Rules     []PolicyRule
}

type policyDTO struct {
	Name      string   `json:"name"`
	Partition string   `json:"partition,omitempty"`
	Controls  []string `json:"controls,omitempty"`
	Requires  []string `json:"requires,omitempty"`
	Strategy  string   `json:"strategy,omitempty"`
	FullPath  string   `json:"fullPath,omitempty"`
	Rules     struct {
		Items []PolicyRule `json:"items,omitempty"`
	} `json:"rulesReference,omitempty"`
}

func (p *Policy) MarshalJSON() ([]byte, error) {
	return json.Marshal(policyDTO{
		Name:      p.Name,
		Partition: p.Partition,
		Controls:  p.Controls,
		Requires:  p.Requires,
		Strategy:  p.Strategy,
		FullPath:  p.FullPath,
		Rules: struct {
			Items []PolicyRule `json:"items,omitempty"`
		}{Items: p.Rules},
	})
}

func (p *Policy) UnmarshalJSON(b []byte) error {
	var dto policyDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	p.Name = dto.Name
	p.Partition = dto.Partition
	p.Controls = dto.Controls
	p.Requires = dto.Requires
	p.Strategy = dto.Strategy
	p.Rules = dto.Rules.Items
	p.FullPath = dto.FullPath

	return nil
}

type PolicyRules struct {
	Items []PolicyRule `json:"items,omitempty"`
}

type PolicyRule struct {
	Name       string
	FullPath   string
	Ordinal    int
	Conditions []PolicyRuleCondition
	Actions    []PolicyRuleAction
}

type policyRuleDTO struct {
	Name       string `json:"name"`
	Ordinal    int    `json:"ordinal"`
	FullPath   string `json:"fullPath,omitempty"`
	Conditions struct {
		Items []PolicyRuleCondition `json:"items,omitempty"`
	} `json:"conditionsReference,omitempty"`
	Actions struct {
		Items []PolicyRuleAction `json:"items,omitempty"`
	} `json:"actionsReference,omitempty"`
}

func (p *PolicyRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(policyRuleDTO{
		Name:     p.Name,
		Ordinal:  p.Ordinal,
		FullPath: p.FullPath,
		Conditions: struct {
			Items []PolicyRuleCondition `json:"items,omitempty"`
		}{Items: p.Conditions},
		Actions: struct {
			Items []PolicyRuleAction `json:"items,omitempty"`
		}{Items: p.Actions},
	})
}

func (p *PolicyRule) UnmarshalJSON(b []byte) error {
	var dto policyRuleDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	p.Name = dto.Name
	p.Ordinal = dto.Ordinal
	p.Actions = dto.Actions.Items
	p.Conditions = dto.Conditions.Items
	p.FullPath = dto.FullPath

	return nil
}

type PolicyRuleActions struct {
	Items []PolicyRuleAction `json:"items"`
}

type PolicyRuleAction struct {
	Name               string `json:"name,omitempty"`
	AppService         string `json:"appService,omitempty"`
	Application        string `json:"application,omitempty"`
	Asm                bool   `json:"asm,omitempty"`
	Avr                bool   `json:"avr,omitempty"`
	Cache              bool   `json:"cache,omitempty"`
	Carp               bool   `json:"carp,omitempty"`
	Category           string `json:"category,omitempty"`
	Classify           bool   `json:"classify,omitempty"`
	ClonePool          string `json:"clonePool,omitempty"`
	Code               int    `json:"code,omitempty"`
	Compress           bool   `json:"compress,omitempty"`
	Content            string `json:"content,omitempty"`
	CookieHash         bool   `json:"cookieHash,omitempty"`
	CookieInsert       bool   `json:"cookieInsert,omitempty"`
	CookiePassive      bool   `json:"cookiePassive,omitempty"`
	CookieRewrite      bool   `json:"cookieRewrite,omitempty"`
	Decompress         bool   `json:"decompress,omitempty"`
	Defer              bool   `json:"defer,omitempty"`
	DestinationAddress bool   `json:"destinationAddress,omitempty"`
	Disable            bool   `json:"disable,omitempty"`
	Domain             string `json:"domain,omitempty"`
	Enable             bool   `json:"enable,omitempty"`
	Expiry             string `json:"expiry,omitempty"`
	ExpirySecs         int    `json:"expirySecs,omitempty"`
	Expression         string `json:"expression,omitempty"`
	Extension          string `json:"extension,omitempty"`
	Facility           string `json:"facility,omitempty"`
	Forward            bool   `json:"forward,omitempty"`
	FromProfile        string `json:"fromProfile,omitempty"`
	Hash               bool   `json:"hash,omitempty"`
	Host               string `json:"host,omitempty"`
	Http               bool   `json:"http,omitempty"`
	HttpBasicAuth      bool   `json:"httpBasicAuth,omitempty"`
	HttpCookie         bool   `json:"httpCookie,omitempty"`
	HttpHeader         bool   `json:"httpHeader,omitempty"`
	HttpHost           bool   `json:"httpHost,omitempty"`
	HttpReferer        bool   `json:"httpReferer,omitempty"`
	HttpReply          bool   `json:"httpReply,omitempty"`
	HttpSetCookie      bool   `json:"httpSetCookie,omitempty"`
	HttpUri            bool   `json:"httpUri,omitempty"`
	Ifile              string `json:"ifile,omitempty"`
	Insert             bool   `json:"insert,omitempty"`
	InternalVirtual    string `json:"internalVirtual,omitempty"`
	IpAddress          string `json:"ipAddress,omitempty"`
	Key                string `json:"key,omitempty"`
	L7dos              bool   `json:"l7dos,omitempty"`
	Length             int    `json:"length,omitempty"`
	Location           string `json:"location,omitempty"`
	Log                bool   `json:"log,omitempty"`
	LtmPolicy          bool   `json:"ltmPolicy,omitempty"`
	Member             string `json:"member,omitempty"`
	Message            string `json:"message,omitempty"`
	TmName             string `json:"tmName,omitempty"`
	Netmask            string `json:"netmask,omitempty"`
	Nexthop            string `json:"nexthop,omitempty"`
	Node               string `json:"node,omitempty"`
	Offset             int    `json:"offset,omitempty"`
	Path               string `json:"path,omitempty"`
	Pem                bool   `json:"pem,omitempty"`
	Persist            bool   `json:"persist,omitempty"`
	Pin                bool   `json:"pin,omitempty"`
	Policy             string `json:"policy,omitempty"`
	Pool               string `json:"pool,omitempty"`
	Port               int    `json:"port,omitempty"`
	Priority           string `json:"priority,omitempty"`
	Profile            string `json:"profile,omitempty"`
	Protocol           string `json:"protocol,omitempty"`
	QueryString        string `json:"queryString,omitempty"`
	Rateclass          string `json:"rateclass,omitempty"`
	Redirect           bool   `json:"redirect,omitempty"`
	Remove             bool   `json:"remove,omitempty"`
	Replace            bool   `json:"replace,omitempty"`
	Request            bool   `json:"request,omitempty"`
	RequestAdapt       bool   `json:"requestAdapt,omitempty"`
	Reset              bool   `json:"reset,omitempty"`
	Response           bool   `json:"response,omitempty"`
	ResponseAdapt      bool   `json:"responseAdapt,omitempty"`
	Scheme             string `json:"scheme,omitempty"`
	Script             string `json:"script,omitempty"`
	Select             bool   `json:"select,omitempty"`
	ServerSsl          bool   `json:"serverSsl,omitempty"`
	SetVariable        bool   `json:"setVariable,omitempty"`
	Snat               string `json:"snat,omitempty"`
	Snatpool           string `json:"snatpool,omitempty"`
	SourceAddress      bool   `json:"sourceAddress,omitempty"`
	SslClientHello     bool   `json:"sslClientHello,omitempty"`
	SslServerHandshake bool   `json:"sslServerHandshake,omitempty"`
	SslServerHello     bool   `json:"sslServerHello,omitempty"`
	SslSessionId       bool   `json:"sslSessionId,omitempty"`
	Status             int    `json:"status,omitempty"`
	Tcl                bool   `json:"tcl,omitempty"`
	TcpNagle           bool   `json:"tcpNagle,omitempty"`
	Text               string `json:"text,omitempty"`
	Timeout            int    `json:"timeout,omitempty"`
	Uie                bool   `json:"uie,omitempty"`
	Universal          bool   `json:"universal,omitempty"`
	Value              string `json:"value,omitempty"`
	Virtual            string `json:"virtual,omitempty"`
	Vlan               string `json:"vlan,omitempty"`
	VlanId             int    `json:"vlanId,omitempty"`
	Wam                bool   `json:"wam,omitempty"`
	Write              bool   `json:"write,omitempty"`
}

type PolicyRuleConditions struct {
	Items []PolicyRuleCondition `json:"items"`
}

type PolicyRuleCondition struct {
	Name                  string   `json:"name,omitempty"`
	Generation            int      `json:"generation,omitempty"`
	Address               bool     `json:"address,omitempty"`
	All                   bool     `json:"all,omitempty"`
	AppService            string   `json:"appService,omitempty"`
	BrowserType           bool     `json:"browserType,omitempty"`
	BrowserVersion        bool     `json:"browserVersion,omitempty"`
	CaseInsensitive       bool     `json:"caseInsensitive,omitempty"`
	CaseSensitive         bool     `json:"caseSensitive,omitempty"`
	Cipher                bool     `json:"cipher,omitempty"`
	CipherBits            bool     `json:"cipherBits,omitempty"`
	ClientSsl             bool     `json:"clientSsl,omitempty"`
	Code                  bool     `json:"code,omitempty"`
	CommonName            bool     `json:"commonName,omitempty"`
	Contains              bool     `json:"contains,omitempty"`
	Continent             bool     `json:"continent,omitempty"`
	CountryCode           bool     `json:"countryCode,omitempty"`
	CountryName           bool     `json:"countryName,omitempty"`
	CpuUsage              bool     `json:"cpuUsage,omitempty"`
	DeviceMake            bool     `json:"deviceMake,omitempty"`
	DeviceModel           bool     `json:"deviceModel,omitempty"`
	Domain                bool     `json:"domain,omitempty"`
	EndsWith              bool     `json:"endsWith,omitempty"`
	Equals                bool     `json:"equals,omitempty"`
	Expiry                bool     `json:"expiry,omitempty"`
	Extension             bool     `json:"extension,omitempty"`
	External              bool     `json:"external,omitempty"`
	Geoip                 bool     `json:"geoip,omitempty"`
	Greater               bool     `json:"greater,omitempty"`
	GreaterOrEqual        bool     `json:"greaterOrEqual,omitempty"`
	Host                  bool     `json:"host,omitempty"`
	HttpBasicAuth         bool     `json:"httpBasicAuth,omitempty"`
	HttpCookie            bool     `json:"httpCookie,omitempty"`
	HttpHeader            bool     `json:"httpHeader,omitempty"`
	HttpHost              bool     `json:"httpHost,omitempty"`
	HttpMethod            bool     `json:"httpMethod,omitempty"`
	HttpReferer           bool     `json:"httpReferer,omitempty"`
	HttpSetCookie         bool     `json:"httpSetCookie,omitempty"`
	HttpStatus            bool     `json:"httpStatus,omitempty"`
	HttpUri               bool     `json:"httpUri,omitempty"`
	HttpUserAgent         bool     `json:"httpUserAgent,omitempty"`
	HttpVersion           bool     `json:"httpVersion,omitempty"`
	Index                 int      `json:"index,omitempty"`
	Internal              bool     `json:"internal,omitempty"`
	Isp                   bool     `json:"isp,omitempty"`
	Last_15secs           bool     `json:"last_15secs,omitempty"`
	Last_1min             bool     `json:"last_1min,omitempty"`
	Last_5mins            bool     `json:"last_5mins,omitempty"`
	Less                  bool     `json:"less,omitempty"`
	LessOrEqual           bool     `json:"lessOrEqual,omitempty"`
	Local                 bool     `json:"local,omitempty"`
	Major                 bool     `json:"major,omitempty"`
	Matches               bool     `json:"matches,omitempty"`
	Minor                 bool     `json:"minor,omitempty"`
	Missing               bool     `json:"missing,omitempty"`
	Mss                   bool     `json:"mss,omitempty"`
	TmName                string   `json:"tmName,omitempty"`
	Not                   bool     `json:"not,omitempty"`
	Org                   bool     `json:"org,omitempty"`
	Password              bool     `json:"password,omitempty"`
	Path                  bool     `json:"path,omitempty"`
	PathSegment           bool     `json:"pathSegment,omitempty"`
	Port                  bool     `json:"port,omitempty"`
	Present               bool     `json:"present,omitempty"`
	Protocol              bool     `json:"protocol,omitempty"`
	QueryParameter        bool     `json:"queryParameter,omitempty"`
	QueryString           bool     `json:"queryString,omitempty"`
	RegionCode            bool     `json:"regionCode,omitempty"`
	RegionName            bool     `json:"regionName,omitempty"`
	Remote                bool     `json:"remote,omitempty"`
	Request               bool     `json:"request,omitempty"`
	Response              bool     `json:"response,omitempty"`
	RouteDomain           bool     `json:"routeDomain,omitempty"`
	Rtt                   bool     `json:"rtt,omitempty"`
	Scheme                bool     `json:"scheme,omitempty"`
	ServerName            bool     `json:"serverName,omitempty"`
	SslCert               bool     `json:"sslCert,omitempty"`
	SslClientHello        bool     `json:"sslClientHello,omitempty"`
	SslExtension          bool     `json:"sslExtension,omitempty"`
	SslServerHandshake    bool     `json:"sslServerHandshake,omitempty"`
	SslServerHello        bool     `json:"sslServerHello,omitempty"`
	StartsWith            bool     `json:"startsWith,omitempty"`
	Tcp                   bool     `json:"tcp,omitempty"`
	Text                  bool     `json:"text,omitempty"`
	UnnamedQueryParameter bool     `json:"unnamedQueryParameter,omitempty"`
	UserAgentToken        bool     `json:"userAgentToken,omitempty"`
	Username              bool     `json:"username,omitempty"`
	Value                 bool     `json:"value,omitempty"`
	Values                []string `json:"values,omitempty"`
	Version               bool     `json:"version,omitempty"`
	Vlan                  bool     `json:"vlan,omitempty"`
	VlanId                bool     `json:"vlanId,omitempty"`
}

func (p *VirtualAddress) MarshalJSON() ([]byte, error) {
	var dto virtualAddressDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *VirtualAddress) UnmarshalJSON(b []byte) error {
	var dto virtualAddressDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

// Monitors contains a list of all monitors on the BIG-IP system.
type Monitors struct {
	Monitors []Monitor `json:"items"`
}

// Monitor contains information about each individual monitor.
type Monitor struct {
	Name           string
	Partition      string
	FullPath       string
	Generation     int
	ParentMonitor  string
	Description    string
	Destination    string
	Interval       int
	IPDSCP         int
	ManualResume   bool
	Password       string
	ReceiveString  string
	ReceiveDisable string
	Reverse        bool
	SendString     string
	TimeUntilUp    int
	Timeout        int
	Transparent    bool
	UpInterval     int
	Username       string
}

type monitorDTO struct {
	Name           string `json:"name,omitempty"`
	Partition      string `json:"partition,omitempty"`
	FullPath       string `json:"fullPath,omitempty"`
	Generation     int    `json:"generation,omitempty"`
	ParentMonitor  string `json:"defaultsFrom,omitempty"`
	Description    string `json:"description,omitempty"`
	Destination    string `json:"destination,omitempty"`
	Interval       int    `json:"interval,omitempty"`
	IPDSCP         int    `json:"ipDscp,omitempty"`
	ManualResume   string `json:"manualResume,omitempty" bool:"enabled"`
	Password       string `json:"password,omitempty"`
	ReceiveString  string `json:"recv,omitempty"`
	ReceiveDisable string `json:"recvDisable,omitempty"`
	Reverse        string `json:"reverse,omitempty" bool:"enabled"`
	SendString     string `json:"send,omitempty"`
	TimeUntilUp    int    `json:"timeUntilUp,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
	Transparent    string `json:"transparent,omitempty" bool:"enabled"`
	UpInterval     int    `json:"upInterval,omitempty"`
	Username       string `json:"username,omitempty"`
}

type Profiles struct {
	Profiles []Profile `json:"items"`
}

type Profile struct {
	Name      string `json:"name,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
	Partition string `json:"partition,omitempty"`
	Context   string `json:"context,omitempty"`
}

type IRules struct {
	IRules []IRule `json:"items"`
}

type IRule struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
	Rule      string `json:"apiAnonymous,omitempty"`
}

type oneconnectDTO struct {
	Name                string `json:"name,omitempty"`
	Partition           string `json:"partition,omitempty"`
	DefaultsFrom        string `json:"defaultsFrom,omitempty"`
	IdleTimeoutOverride string `json:"idleTimeoutOverride,omitempty"`
	MaxAge              int    `json:"maxAge,omitempty"`
	MaxReuse            int    `json:"maxReuse,omitempty"`
	MaxSize             int    `json:"maxSize,omitempty"`
	SourceMask          string `json:"sourceMask,omitempty"`
	SharePools          string `json:"sharePools,omitempty"`
}
type Oneconnects struct {
	Oneconnects []Oneconnect `json:"items"`
}

type Oneconnect struct {
	Name                string
	Partition           string
	DefaultsFrom        string
	IdleTimeoutOverride string
	MaxAge              int
	MaxReuse            int
	MaxSize             int
	SourceMask          string
	SharePools          string
}

type tcpDTO struct {
	Name              string `json:"name,omitempty"`
	Partition         string `json:"partition,omitempty"`
	DefaultsFrom      string `json:"defaultsFrom,omitempty"`
	IdleTimeout       int    `json:"idleTimeout,omitempty"`
	CloseWaitTimeout  int    `json:"closeWaitTimeout,omitempty"`
	FinWait_2Timeout  int    `json:"finWait_2Timeout,omitempty"`
	FinWaitTimeout    int    `json:"finWaitTimeout,omitempty"`
	KeepAliveInterval int    `json:"keepAliveInterval,omitempty"`
	DeferredAccept    string `json:"deferredAccept,omitempty"`
	FastOpen          string `json:"fastOpen,omitempty"`
}

type Tcps struct {
	Tcps []Tcp `json:"items"`
}

type Tcp struct {
	Name              string
	Partition         string
	DefaultsFrom      string
	IdleTimeout       int
	CloseWaitTimeout  int
	FinWait_2Timeout  int
	FinWaitTimeout    int
	KeepAliveInterval int
	DeferredAccept    string
	FastOpen          string
}

type fasthttpDTO struct {
	Name                        string `json:"name,omitempty"`
	DefaultsFrom                string `json:"defaultsFrom,omitempty"`
	IdleTimeout                 int    `json:"idleTimeout,omitempty"`
	ConnpoolIdleTimeoutOverride int    `json:"connpoolIdleTimeoutOverride,omitempty"`
	ConnpoolMaxReuse            int    `json:"connpoolMaxReuse,omitempty"`
	ConnpoolMaxSize             int    `json:"connpoolMaxSize,omitempty"`
	ConnpoolMinSize             int    `json:"connpoolMinSize,omitempty"`
	ConnpoolReplenish           string `json:"connpoolReplenish,omitempty"`
	ConnpoolStep                int    `json:"deferredAccept,omitempty"`
	ForceHttp_10Response        string `json:"forceHttp_10Response,omitempty"`
	MaxHeaderSize               int    `json:"maxHeaderSize,omitempty"`
}

type Fasthttps struct {
	Fasthttps []Fasthttp `json:"items"`
}

type Fasthttp struct {
	Name                        string
	DefaultsFrom                string
	IdleTimeout                 int
	ConnpoolIdleTimeoutOverride int
	ConnpoolMaxReuse            int
	ConnpoolMaxSize             int
	ConnpoolMinSize             int
	ConnpoolReplenish           string
	ConnpoolStep                int
	ForceHttp_10Response        string
	MaxHeaderSize               int
}

type fastl4DTO struct {
	Name                  string `json:"name,omitempty"`
	DefaultsFrom          string `json:"defaultsFrom,omitempty"`
	Partition             string `json:"partition,omitempty"`
	ExplicitFlowMigration string `json:"explicitFlowMigration,omitempty"`
	HardwareSynCookie     string `json:"hardwareSynCookie,omitem"`
	IdleTimeout           string `json:"idleTimeout,omitempty"`
	ClientTimeout         int    `json:"clientTimeout,omitempty"`
	IpTosToClient         string `json:"ipTosToClient,omitempty"`
	IpTosToServer         string `json:"ipTosToServer,omitempty"`
	KeepAliveInterval     string `json:"keepAliveInterval,omitempty"`
}

type Fastl4s struct {
	Fastl4s []Fastl4 `json:"items"`
}

type Fastl4 struct {
	Name                  string
	Partition             string
	DefaultsFrom          string
	ExplicitFlowMigration string
	HardwareSynCookie     string
	IdleTimeout           string
	ClientTimeout         int
	IpTosToClient         string
	IpTosToServer         string
	KeepAliveInterval     string
}

type httpcompressDTO struct {
	Name         string   `json:"name,omitempty"`
	DefaultsFrom string   `json:"defaultsFrom,omitempty"`
	UriExclude   []string `json:"uriExclude,omitempty"`
	UriInclude   []string `json:"uriInclude,omitempty"`
}

type Httpcompresss struct {
	Httpcompresss []Httpcompress `json:"items"`
}

type Httpcompress struct {
	Name         string
	DefaultsFrom string
	UriExclude   []string
	UriInclude   []string
}

type http2DTO struct {
	Name                           string   `json:"name,omitempty"`
	DefaultsFrom                   string   `json:"defaultsFrom,omitempty"`
	ConcurrentStreamsPerConnection int      `json:"concurrentStreamsPerConnection,omitempty"`
	ConnectionIdleTimeout          int      `json:"connectionIdleTimeout,omitempty"`
	HeaderTableSize                int      `json:"headerTableSize,omitempty"`
	ActivationModes                []string `json:"activationModes,omitempty"`
}

type Http2s struct {
	Http2s []Http2 `json:"items"`
}

type Http2 struct {
	Name                           string
	DefaultsFrom                   string
	ConcurrentStreamsPerConnection int
	ConnectionIdleTimeout          int
	HeaderTableSize                int
	ActivationModes                []string
}

type Recordss struct {
	Recordss []Records `json:"items"`
}

type Records struct {
	Name string
	Data string
}
type RecordsDTO struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}

type Datagroup struct {
	Name    string
	Type    string
	Records []Records
}

type Datagroups struct {
	Datagroups []Datagroup `json:"items"`
}

type DatagroupDTO struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Records struct {
		Items []Records `json:"items,omitempty"`
	} `json:"records,omitempty"`
}

func (p *Datagroup) MarshalJSON() ([]byte, error) {
	var dto DatagroupDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Datagroup) UnmarshalJSON(b []byte) error {
	var dto DatagroupDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (p *Records) MarshalJSON() ([]byte, error) {
	var dto RecordsDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Records) UnmarshalJSON(b []byte) error {
	var dto RecordsDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (p *Monitor) MarshalJSON() ([]byte, error) {
	var dto monitorDTO
	marshal(&dto, p)
	if strings.Contains(dto.SendString, "\r\n") {
		dto.SendString = strings.Replace(dto.SendString, "\r\n", "\\r\\n", -1)
	}
	return json.Marshal(dto)
}

func (p *Monitor) UnmarshalJSON(b []byte) error {
	var dto monitorDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (p *Oneconnect) MarshalJSON() ([]byte, error) {
	var dto oneconnectDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Oneconnect) UnmarshalJSON(b []byte) error {
	var dto oneconnectDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (p *Tcp) MarshalJSON() ([]byte, error) {
	var dto tcpDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Tcp) UnmarshalJSON(b []byte) error {
	var dto tcpDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (p *Fasthttp) MarshalJSON() ([]byte, error) {
	var dto fasthttpDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Fasthttp) UnmarshalJSON(b []byte) error {
	var dto fasthttpDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (p *Fastl4) MarshalJSON() ([]byte, error) {
	var dto fastl4DTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Fastl4) UnmarshalJSON(b []byte) error {
	var dto fastl4DTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (p *Httpcompress) MarshalJSON() ([]byte, error) {
	var dto httpcompressDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Httpcompress) UnmarshalJSON(b []byte) error {
	var dto httpcompressDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

func (p *Http2) MarshalJSON() ([]byte, error) {
	var dto http2DTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Http2) UnmarshalJSON(b []byte) error {
	var dto http2DTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

type Snats struct {
	Snats []Snat `json:"items"`
}

// VirtualAddress contains information about each individual virtual address.
type Snat struct {
	Name          string
	Partition     string
	FullPath      string
	AutoLasthop   string
	Mirror        bool
	SourcePort    string
	Translation   string
	Snatpool      string
	VlansDisabled bool
	Origins       []string
}

type snatDTO struct {
	Name          string   `json:"name"`
	Partition     string   `json:"partition,omitempty"`
	FullPath      string   `json:"fullPath,omitempty"`
	AutoLasthop   string   `json:"autoLastHop,omitempty"`
	Mirror        bool     `json:"mirror,omitempty" bool:"disabled"`
	SourcePort    string   `json:"sourePort,omitempty"`
	Translation   string   `json:"translation,omitempty"`
	Snatpool      string   `json:"snatpool,omitempty"`
	VlansDisabled bool     `json:"vlansDisabled,omitempty" bool:"disabled"`
	Origins       []string `json:"origins,omitempty"`
}

func (p *Snat) MarshalJSON() ([]byte, error) {
	var dto snatDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Snat) UnmarshalJSON(b []byte) error {
	var dto snatDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	p.Origins = dto.Origins
	//return marshal(p, &dto)
	return nil
}

type Snatpools struct {
	Snatpools []Snatpool `json:"items"`
}

// Snatpool structure
type Snatpool struct {
	Name      string
	Partition string
	Members   []string
}

type snatpoolDTO struct {
	Name      string   `json:"name"`
	Partition string   `json:"partition,omitempty"`
	Members   []string `json:"members,omitempty"`
}

func (p *Snatpool) MarshalJSON() ([]byte, error) {
	var dto snatpoolDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Snatpool) UnmarshalJSON(b []byte) error {
	var dto snatpoolDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	p.Members = dto.Members
	//return marshal(p, &dto)
	return nil
}

const (
	uriLtm            = "ltm"
	uriNode           = "node"
	uriPool           = "pool"
	uriVirtual        = "virtual"
	uriVirtualAddress = "virtual-address"
	uriMonitor        = "monitor"
	uriIRule          = "rule"
	uriDatagroup      = "data-group"
	uriInternal       = "internal"
	uriPolicy         = "policy"
	uriProfile        = "profile"
	uriOneconnect     = "one-connect"
	ENABLED           = "enable"
	DISABLED          = "disable"
	CONTEXT_SERVER    = "serverside"
	CONTEXT_CLIENT    = "clientside"
	CONTEXT_ALL       = "all"
	uriTcp            = "tcp"
	uriFasthttp       = "fasthttp"
	uriFastl4         = "fastl4"
	uriHttpcompress   = "http-compression"
	uriHttp2          = "http2"
	uriSnat           = "snat"
	uriSnatpool       = "snatpool"
)

var cidr = map[string]string{
	"0":  "0.0.0.0",
	"1":  "128.0.0.0",
	"2":  "192.0.0.0",
	"3":  "224.0.0.0",
	"4":  "240.0.0.0",
	"5":  "248.0.0.0",
	"6":  "252.0.0.0",
	"7":  "254.0.0.0",
	"8":  "255.0.0.0",
	"9":  "255.128.0.0",
	"10": "255.192.0.0",
	"11": "255.224.0.0",
	"12": "255.240.0.0",
	"13": "255.248.0.0",
	"14": "255.252.0.0",
	"15": "255.254.0.0",
	"16": "255.255.0.0",
	"17": "255.255.128.0",
	"18": "255.255.192.0",
	"19": "255.255.224.0",
	"20": "255.255.240.0",
	"21": "255.255.248.0",
	"22": "255.255.252.0",
	"23": "255.255.254.0",
	"24": "255.255.255.0",
	"25": "255.255.255.128",
	"26": "255.255.255.192",
	"27": "255.255.255.224",
	"28": "255.255.255.240",
	"29": "255.255.255.248",
	"30": "255.255.255.252",
	"31": "255.255.255.254",
	"32": "255.255.255.255",
}

// Nodes returns a list of nodes.
func (b *BigIP) Nodes() (*Nodes, error) {
	var nodes Nodes
	err, _ := b.getForEntity(&nodes, uriLtm, uriNode)
	if err != nil {
		return nil, err
	}

	return &nodes, nil
}

// CreateNode adds a new node to the BIG-IP system.
func (b *BigIP) CreateNode(name, address string) error {
	config := &Node{
		Name:    name,
		Address: address,
	}

	return b.post(config, uriLtm, uriNode)
}

// Get a Node by name. Returns nil if the node does not exist
func (b *BigIP) GetNode(name string) (*Node, error) {
	var node Node
	err, ok := b.getForEntity(&node, uriLtm, uriNode, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &node, nil
}

// DeleteNode removes a node.
func (b *BigIP) DeleteNode(name string) error {
	return b.delete(uriLtm, uriNode, name)
}

// ModifyNode allows you to change any attribute of a node. Fields that
// can be modified are referenced in the Node struct.
func (b *BigIP) ModifyNode(name string, config *Node) error {
	return b.put(config, uriLtm, uriNode, name)
}

// NodeStatus changes the status of a node. <state> can be either
// "enable" or "disable".
func (b *BigIP) NodeStatus(name, state string) error {
	config := &Node{}

	switch state {
	case "enable":
		// config.State = "unchecked"
		config.Session = "user-enabled"
	case "disable":
		// config.State = "unchecked"
		config.Session = "user-disabled"
		// case "offline":
		// 	config.State = "user-down"
		// 	config.Session = "user-disabled"
	}

	return b.put(config, uriLtm, uriNode, name)
}

// Pools returns a list of pools.
func (b *BigIP) Pools() (*Pools, error) {
	var pools Pools
	err, _ := b.getForEntity(&pools, uriLtm, uriPool)
	if err != nil {
		return nil, err
	}

	return &pools, nil
}

// PoolMembers returns a list of pool members for the given pool.
func (b *BigIP) PoolMembers(name string) ([]string, error) {
	var nodes Nodes
	members := []string{}
	errString := []string{}
	err, _ := b.getForEntity(&nodes, uriLtm, uriPool, name, "members")
	if err != nil {
		return errString, err
	}

	for _, m := range nodes.Nodes {
		members = append(members, m.Name)
	}

	return members, nil
}

// AddPoolMember adds a node/member to the given pool. <member> must be in the form
// of <node>:<port>, i.e.: "web-server1:443".
func (b *BigIP) AddPoolMember(pool, member string) error {
	config := &poolMember{
		Name: member,
	}

	return b.post(config, uriLtm, uriPool, pool, "members")
}

// DeletePoolMember removes a member from the given pool. <member> must be in the form
// of <node>:<port>, i.e.: "web-server1:443".
func (b *BigIP) DeletePoolMember(pool, member string) error {
	return b.delete(uriLtm, uriPool, pool, "members", member)
}

// PoolMemberStatus changes the status of a pool member. <state> can be either
// "enable" or "disable". <member> must be in the form of <node>:<port>,
// i.e.: "web-server1:443".
func (b *BigIP) PoolMemberStatus(pool, member, state string) error {
	config := &Node{}

	switch state {
	case "enable":
		// config.State = "unchecked"
		config.Session = "user-enabled"
	case "disable":
		// config.State = "unchecked"
		config.Session = "user-disabled"
		// case "offline":
		// 	config.State = "user-down"
		// 	config.Session = "user-disabled"
	}

	return b.put(config, uriLtm, uriPool, pool, "members", member)
}

// CreatePool adds a new pool to the BIG-IP system.
func (b *BigIP) CreatePool(name string) error {
	config := &Pool{
		Name: name,
	}

	return b.post(config, uriLtm, uriPool)
}

// Get a Pool by name. Returns nil if the Pool does not exist
func (b *BigIP) GetPool(name string) (*Pool, error) {
	var pool Pool
	err, ok := b.getForEntity(&pool, uriLtm, uriPool, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &pool, nil
}

// DeletePool removes a pool.
func (b *BigIP) DeletePool(name string) error {
	return b.delete(uriLtm, uriPool, name)
}

// ModifyPool allows you to change any attribute of a pool. Fields that
// can be modified are referenced in the Pool struct.
func (b *BigIP) ModifyPool(name string, config *Pool) error {
	return b.put(config, uriLtm, uriPool, name)
}

// VirtualServers returns a list of virtual servers.
func (b *BigIP) VirtualServers() (*VirtualServers, error) {
	var vs VirtualServers
	err, _ := b.getForEntity(&vs, uriLtm, uriVirtual)
	if err != nil {
		return nil, err
	}

	return &vs, nil
}

// CreateVirtualServer adds a new virtual server to the BIG-IP system. <mask> can either be
// in CIDR notation or decimal, i.e.: "24" or "255.255.255.0". A CIDR mask of "0" is the same
// as "0.0.0.0".
func (b *BigIP) CreateVirtualServer(name, destination, mask, pool string, port int) error {
	subnetMask := cidr[mask]

	if strings.Contains(mask, ".") {
		subnetMask = mask
	}

	config := &VirtualServer{
		Name:        name,
		Destination: fmt.Sprintf("%s:%d", destination, port),
		Mask:        subnetMask,
		Pool:        pool,
	}

	return b.post(config, uriLtm, uriVirtual)
}

// Get a VirtualServer by name. Returns nil if the VirtualServer does not exist
func (b *BigIP) GetVirtualServer(name string) (*VirtualServer, error) {
	var vs VirtualServer
	err, ok := b.getForEntity(&vs, uriLtm, uriVirtual, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	profiles, err := b.VirtualServerProfiles(name)
	if err != nil {
		return nil, err
	}
	vs.Profiles = profiles.Profiles

	policy_names, err := b.VirtualServerPolicyNames(name)
	if err != nil {
		return nil, err
	}
	vs.Policies = policy_names

	return &vs, nil
}

// DeleteVirtualServer removes a virtual server.
func (b *BigIP) DeleteVirtualServer(name string) error {
	return b.delete(uriLtm, uriVirtual, name)
}

// ModifyVirtualServer allows you to change any attribute of a virtual server. Fields that
// can be modified are referenced in the VirtualServer struct.
func (b *BigIP) ModifyVirtualServer(name string, config *VirtualServer) error {
	return b.put(config, uriLtm, uriVirtual, name)
}

// VirtualServerProfiles gets the profiles currently associated with a virtual server.
func (b *BigIP) VirtualServerProfiles(vs string) (*Profiles, error) {
	var p Profiles
	err, ok := b.getForEntity(&p, uriLtm, uriVirtual, vs, "profiles")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &p, nil
}

//Get the names of policies associated with a particular virtual server
func (b *BigIP) VirtualServerPolicyNames(vs string) ([]string, error) {
	var policies VirtualServerPolicies
	err, _ := b.getForEntity(&policies, uriLtm, uriVirtual, vs, "policies")
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", policies)
	retval := make([]string, 0, len(policies.PolicyRef.Policies))
	for _, p := range policies.PolicyRef.Policies {
		retval = append(retval, p.FullPath)
	}
	return retval, nil
}

// VirtualAddresses returns a list of virtual addresses.
func (b *BigIP) VirtualAddresses() (*VirtualAddresses, error) {
	var va VirtualAddresses
	err, _ := b.getForEntity(&va, uriLtm, uriVirtualAddress)
	if err != nil {
		return nil, err
	}
	return &va, nil
}

func (b *BigIP) CreateVirtualAddress(vaddr string, config *VirtualAddress) error {
	config.Name = vaddr
	return b.post(config, uriLtm, uriVirtualAddress)
}

// VirtualAddressStatus changes the status of a virtual address. <state> can be either
// "enable" or "disable".
func (b *BigIP) VirtualAddressStatus(vaddr, state string) error {
	config := &VirtualAddress{}
	config.Enabled = (state == ENABLED)
	return b.put(config, uriLtm, uriVirtualAddress, vaddr)
}

// ModifyVirtualAddress allows you to change any attribute of a virtual address. Fields that
// can be modified are referenced in the VirtualAddress struct.
func (b *BigIP) ModifyVirtualAddress(vaddr string, config *VirtualAddress) error {
	return b.put(config, uriLtm, uriVirtualAddress, vaddr)
}

func (b *BigIP) DeleteVirtualAddress(vaddr string) error {
	return b.delete(uriLtm, uriVirtualAddress, vaddr)
}

// Monitors returns a list of all HTTP, HTTPS, Gateway ICMP, and ICMP monitors.
func (b *BigIP) Monitors() ([]Monitor, error) {
	var monitors []Monitor
	monitorUris := []string{"http", "https", "icmp", "gateway-icmp"}

	for _, name := range monitorUris {
		var m Monitors
		err, _ := b.getForEntity(&m, uriLtm, uriMonitor, name)
		if err != nil {
			return nil, err
		}
		for _, monitor := range m.Monitors {
			monitors = append(monitors, monitor)
		}
	}

	return monitors, nil
}

// CreateMonitor adds a new monitor to the BIG-IP system. <parent> must be one of "http", "https",
// "icmp", or "gateway icmp".
func (b *BigIP) CreateMonitor(name, parent string, interval, timeout int, send, receive string) error {
	config := &Monitor{
		Name:          name,
		ParentMonitor: parent,
		Interval:      interval,
		Timeout:       timeout,
		SendString:    send,
		ReceiveString: receive,
	}

	return b.AddMonitor(config)
}

// Create a monitor by supplying a config
func (b *BigIP) AddMonitor(config *Monitor) error {
	if strings.Contains(config.ParentMonitor, "gateway") {
		config.ParentMonitor = "gateway_icmp"
	}

	return b.post(config, uriLtm, uriMonitor, config.ParentMonitor)
}

// DeleteMonitor removes a monitor.
func (b *BigIP) DeleteMonitor(name, parent string) error {
	return b.delete(uriLtm, uriMonitor, parent, name)
}

// ModifyMonitor allows you to change any attribute of a monitor. <parent> must be
// one of "http", "https", "icmp", or "gateway icmp". Fields that
// can be modified are referenced in the Monitor struct.
func (b *BigIP) ModifyMonitor(name, parent string, config *Monitor) error {
	if strings.Contains(config.ParentMonitor, "gateway") {
		config.ParentMonitor = "gateway_icmp"
	}

	return b.put(config, uriLtm, uriMonitor, parent, name)
}

// AddMonitorToPool assigns the monitor, <monitor> to the given <pool>.
func (b *BigIP) AddMonitorToPool(monitor, pool string) error {
	config := &Pool{
		Monitor: monitor,
	}

	return b.put(config, uriLtm, uriPool, pool)
}

// IRules returns a list of irules
func (b *BigIP) IRules() (*IRules, error) {
	var rules IRules
	err, _ := b.getForEntity(&rules, uriLtm, uriIRule)
	if err != nil {
		return nil, err
	}

	return &rules, nil
}

// IRule returns information about the given iRule.
func (b *BigIP) IRule(name string) (*IRule, error) {
	var rule IRule
	err, ok := b.getForEntity(&rule, uriLtm, uriIRule, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &rule, nil
}

// CreateIRule creates a new iRule on the system.
func (b *BigIP) CreateIRule(name, rule string) error {
	irule := &IRule{
		Name: name,
		Rule: rule,
	}
	return b.post(irule, uriLtm, uriIRule)
}

// DeleteIRule removes an iRule from the system.
func (b *BigIP) DeleteIRule(name string) error {
	return b.delete(uriLtm, uriIRule, name)
}

// ModifyIRule updates the given iRule with any changed values.
func (b *BigIP) ModifyIRule(name string, irule *IRule) error {
	irule.Name = name
	return b.put(irule, uriLtm, uriIRule, name)
}

func (b *BigIP) Policies() (*Policies, error) {
	var p Policies
	err, _ := b.getForEntity(&p, uriLtm, uriPolicy)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

//Load a fully policy definition. Policies seem to be best dealt with as one big entity.
func (b *BigIP) GetPolicy(name string) (*Policy, error) {
	var p Policy
	err, ok := b.getForEntity(&p, uriLtm, uriPolicy, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	var rules PolicyRules
	err, _ = b.getForEntity(&rules, uriLtm, uriPolicy, name, "rules")
	if err != nil {
		return nil, err
	}
	p.Rules = rules.Items

	for i, _ := range p.Rules {
		var a PolicyRuleActions
		var c PolicyRuleConditions

		err, _ = b.getForEntity(&a, uriLtm, uriPolicy, name, "rules", p.Rules[i].Name, "actions")
		if err != nil {
			return nil, err
		}
		err, _ = b.getForEntity(&c, uriLtm, uriPolicy, name, "rules", p.Rules[i].Name, "conditions")
		if err != nil {
			return nil, err
		}
		p.Rules[i].Actions = a.Items
		p.Rules[i].Conditions = c.Items
	}

	return &p, nil
}

func normalizePolicy(p *Policy) {
	//f5 doesn't seem to automatically handle setting the ordinal
	for ri, _ := range p.Rules {
		p.Rules[ri].Ordinal = ri
		for ai, _ := range p.Rules[ri].Actions {
			p.Rules[ri].Actions[ai].Name = fmt.Sprintf("%d", ai)
		}
		for ci, _ := range p.Rules[ri].Conditions {
			p.Rules[ri].Conditions[ci].Name = fmt.Sprintf("%d", ci)
		}
	}
}

//Create a new policy. It is not necessary to set the Ordinal fields on subcollections.
func (b *BigIP) CreatePolicy(p *Policy) error {
	normalizePolicy(p)
	return b.post(p, uriLtm, uriPolicy)
}

//Update an existing policy.
func (b *BigIP) UpdatePolicy(name string, p *Policy) error {
	normalizePolicy(p)
	return b.put(p, uriLtm, uriPolicy, name)
}

//Delete a policy by name.
func (b *BigIP) DeletePolicy(name string) error {
	return b.delete(uriLtm, uriPolicy, name)
}

// Oneconnect profile creation
func (b *BigIP) CreateOneconnect(name, idleTimeoutOverride, partition, defaultsFrom, sharePools, sourceMask string, maxAge, maxReuse, maxSize int) error {
	oneconnect := &Oneconnect{
		Name:                name,
		IdleTimeoutOverride: idleTimeoutOverride,
		Partition:           partition,
		DefaultsFrom:        defaultsFrom,
		SharePools:          sharePools,
		SourceMask:          sourceMask,
		MaxAge:              maxAge,
		MaxReuse:            maxReuse,
		MaxSize:             maxSize,
	}
	return b.post(oneconnect, uriLtm, uriProfile, uriOneconnect)
}

func (b *BigIP) Oneconnect(name string) (*Oneconnects, error) {
	var oneconnects Oneconnects
	err, _ := b.getForEntity(&oneconnects, uriLtm, uriProfile, uriOneconnect)

	if err != nil {
		return nil, err
	}

	return &oneconnects, nil
}

// DeleteOneconnect removes an OneConnect profile from the system.
func (b *BigIP) DeleteOneconnect(name string) error {
	return b.delete(uriLtm, uriProfile, uriOneconnect, name)
}

// ModifyOneconnect updates the given Oneconnect profile with any changed values.
func (b *BigIP) ModifyOneconnect(name string, oneconnect *Oneconnect) error {
	oneconnect.Name = name
	return b.put(oneconnect, uriLtm, uriProfile, uriOneconnect, name)
}

// Create TCP profile for WAN or LAN

func (b *BigIP) CreateTcp(name, partition, defaultsFrom string, idleTimeout, closeWaitTimeout, finWait_2Timeout, finWaitTimeout, keepAliveInterval int, deferredAccept, fastOpen string) error {
	tcp := &Tcp{
		Name:              name,
		Partition:         partition,
		DefaultsFrom:      defaultsFrom,
		IdleTimeout:       idleTimeout,
		CloseWaitTimeout:  closeWaitTimeout,
		FinWait_2Timeout:  finWait_2Timeout,
		FinWaitTimeout:    finWaitTimeout,
		KeepAliveInterval: keepAliveInterval,
		DeferredAccept:    deferredAccept,
		FastOpen:          fastOpen,
	}
	return b.post(tcp, uriLtm, uriProfile, uriTcp)
}

// DeleteOneconnect removes an OneConnect profile from the system.
func (b *BigIP) DeleteTcp(name string) error {
	return b.delete(uriLtm, uriProfile, uriTcp, name)
}

// ModifyTcp updates the given Oneconnect profile with any changed values.
func (b *BigIP) ModifyTcp(name string, tcp *Tcp) error {
	tcp.Name = name
	return b.put(tcp, uriLtm, uriProfile, uriTcp, name)
}

func (b *BigIP) Tcp(name string) (*Tcps, error) {
	var tcps Tcps
	err, _ := b.getForEntity(&tcps, uriLtm, uriProfile, uriTcp)

	if err != nil {
		return nil, err
	}

	return &tcps, nil
}

func (b *BigIP) CreateFasthttp(name, defaultsFrom string, idleTimeout, connpoolIdleTimeoutOverride, connpoolMaxReuse, connpoolMaxSize, connpoolMinSize int, connpoolReplenish string, connpoolStep int, forceHttp_10Response string, maxHeaderSize int) error {
	fasthttp := &Fasthttp{
		Name:                        name,
		DefaultsFrom:                defaultsFrom,
		IdleTimeout:                 idleTimeout,
		ConnpoolIdleTimeoutOverride: connpoolIdleTimeoutOverride,
		ConnpoolMaxReuse:            connpoolMaxReuse,
		ConnpoolMaxSize:             connpoolMaxSize,
		ConnpoolMinSize:             connpoolMinSize,
		ConnpoolReplenish:           connpoolReplenish,
		ConnpoolStep:                connpoolStep,
		ForceHttp_10Response:        forceHttp_10Response,
		MaxHeaderSize:               maxHeaderSize,
	}
	return b.post(fasthttp, uriLtm, uriProfile, uriFasthttp)
}

// Delete Fast http removes an Fasthttp profile from the system.
func (b *BigIP) DeleteFasthttp(name string) error {
	return b.delete(uriLtm, uriProfile, uriFasthttp, name)
}

// ModifyFasthttp updates the given Fasthttp profile with any changed values.
func (b *BigIP) ModifyFasthttp(name string, fasthttp *Fasthttp) error {
	fasthttp.Name = name
	return b.put(fasthttp, uriLtm, uriProfile, uriFasthttp, name)
}

func (b *BigIP) Fasthttp() (*Fasthttps, error) {
	var fasthttps Fasthttps
	err, _ := b.getForEntity(&fasthttps, uriLtm, uriProfile, uriFasthttp)

	if err != nil {
		return nil, err
	}

	return &fasthttps, nil
}

func (b *BigIP) CreateFastl4(name, partition, defaultsFrom string, clientTimeout int, explicitFlowMigration, hardwareSynCookie string, idleTimeout string, ipTosToClient, ipTosToServer, keepAliveInterval string) error {
	fastl4 := &Fastl4{
		Name:                  name,
		Partition:             partition,
		DefaultsFrom:          defaultsFrom,
		ClientTimeout:         clientTimeout,
		ExplicitFlowMigration: explicitFlowMigration,
		HardwareSynCookie:     hardwareSynCookie,
		IdleTimeout:           idleTimeout,
		IpTosToClient:         ipTosToClient,
		IpTosToServer:         ipTosToServer,
		KeepAliveInterval:     keepAliveInterval,
	}
	return b.post(fastl4, uriLtm, uriProfile, uriFastl4)
}

// Delete Fast http removes an Fasthttp profile from the system.
func (b *BigIP) DeleteFastl4(name string) error {
	return b.delete(uriLtm, uriProfile, uriFastl4, name)
}

// ModifyFastl4 updates the given Fastl4 profile with any changed values.
func (b *BigIP) ModifyFastl4(name string, fastl4 *Fastl4) error {
	fastl4.Name = name
	return b.put(fastl4, uriLtm, uriProfile, uriFastl4, name)
}

func (b *BigIP) Fastl4(name string) (*Fastl4s, error) {
	var fastl4s Fastl4s
	err, _ := b.getForEntity(&fastl4s, uriLtm, uriProfile, uriFastl4)

	if err != nil {
		return nil, err
	}

	return &fastl4s, nil
}

// ===============

func (b *BigIP) CreateHttpcompress(name, defaultsFrom string, uriExclude, uriInclude []string) error {
	httpcompress := &Httpcompress{
		Name:         name,
		DefaultsFrom: defaultsFrom,
		UriExclude:   uriExclude,
		UriInclude:   uriInclude,
	}
	return b.post(httpcompress, uriLtm, uriProfile, uriHttpcompress)
}

// Delete Fast http removes an Fasthttp profile from the system.
func (b *BigIP) DeleteHttpcompress(name string) error {
	return b.delete(uriLtm, uriProfile, uriHttpcompress, name)
}

// ModifyFastl4 updates the given Fastl4 profile with any changed values.
func (b *BigIP) ModifyHttpcompress(name string, httpcompress *Httpcompress) error {
	httpcompress.Name = name
	return b.put(httpcompress, uriLtm, uriProfile, uriHttpcompress, name)
}

func (b *BigIP) Httpcompress() (*Httpcompresss, error) {
	var httpcompresss Httpcompresss
	err, _ := b.getForEntity(&httpcompresss, uriLtm, uriProfile, uriHttpcompress)

	if err != nil {
		return nil, err
	}

	return &httpcompresss, nil
}

func (b *BigIP) CreateHttp2(name, defaultsFrom string, concurrentStreamsPerConnection, connectionIdleTimeout, headerTableSize int, activationModes []string) error {
	http2 := &Http2{
		Name:                           name,
		DefaultsFrom:                   defaultsFrom,
		ConcurrentStreamsPerConnection: concurrentStreamsPerConnection,
		ConnectionIdleTimeout:          connectionIdleTimeout,
		HeaderTableSize:                headerTableSize,
		ActivationModes:                activationModes,
	}
	return b.post(http2, uriLtm, uriProfile, uriHttp2)
}

// Delete  http2 removes an http2 profile from the system.
func (b *BigIP) DeleteHttp2(name string) error {
	return b.delete(uriLtm, uriProfile, uriHttp2, name)
}

// Modify http2 updates the given http2 profile with any changed values.
func (b *BigIP) ModifyHttp2(name string, http2 *Http2) error {
	http2.Name = name
	return b.put(http2, uriLtm, uriProfile, uriHttp2, name)
}

func (b *BigIP) Http2() (*Http2s, error) {
	var http2s Http2s
	err, _ := b.getForEntity(&http2s, uriLtm, uriProfile, uriHttp2)

	if err != nil {
		return nil, err
	}

	return &http2s, nil
}

// Datagroups returns a list of datagroups.
func (b *BigIP) Datagroups(name string) (*Datagroups, error) {
	var datagroups Datagroups
	err, _ := b.getForEntity(&datagroups, uriLtm, uriDatagroup, uriInternal, name)

	if err != nil {
		return nil, err
	}

	return &datagroups, nil
}

// CreateDatagroup adds a new Datagroup to the BIG-IP system.
func (b *BigIP) CreateDatagroup(typo, name string, records []Records) error {
	//func (b *BigIP) CreateDatagroup(typo, name string) error {
	config := &Datagroup{
		Type:    typo,
		Name:    name,
		Records: records,
	}

	return b.post(config, uriLtm, uriDatagroup, uriInternal)
}
func (b *BigIP) Records() (*Records, error) {
	var records Records
	err, _ := b.getForEntity(&records, uriLtm, uriDatagroup, uriInternal)

	if err != nil {
		return nil, err
	}

	return &records, nil
}
func (b *BigIP) DeleteDatagroup(name string) error {

	return b.delete(uriLtm, uriDatagroup, uriInternal, name)
}

func (b *BigIP) AddRecords(name, rname, data string) error {
	config := &Records{}

	config.Name = rname
	config.Data = data

	return b.post(config, uriLtm, uriDatagroup, uriInternal, "records")
}

// Snats returns a list of snat
func (b *BigIP) Snats(name string) (*Snats, error) {
	var snats Snats
	err, _ := b.getForEntity(&snats, uriLtm, uriSnat, name)

	if err != nil {
		return nil, err
	}

	return &snats, nil
}

func (b *BigIP) CreateSnat(name, partition, autoLastHop, sourcePort, translation, snatpool string, vlansDisabled, mirror bool, origins []string) error {
	snat := &Snat{
		Name:          name,
		Partition:     partition,
		AutoLasthop:   autoLastHop,
		SourcePort:    sourcePort,
		Translation:   translation,
		Snatpool:      snatpool,
		Mirror:        mirror,
		VlansDisabled: vlansDisabled,
		Origins:       origins,
	}
	log.Println("[INFO] Creating snat  ", snat)
	return b.post(snat, uriLtm, uriSnat)
}

func (b *BigIP) ModifySnat(config *Snat) error {
	return b.patch(config, uriLtm, uriSnat)
}

// Get a Snat list  by name. Returns nil if the node does not exist
func (b *BigIP) GetSnat(name string) (*Snat, error) {
	var snat Snat
	err, ok := b.getForEntity(&snat, uriLtm, uriSnat, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &snat, nil
}

func (b *BigIP) DeleteSnat(name string) error {
	return b.delete(uriLtm, uriSnat, name)
}

// Snats returns a list of snat
func (b *BigIP) Snatpools(name string) (*Snatpools, error) {
	var snatpools Snatpools
	err, _ := b.getForEntity(&snatpools, uriLtm, uriSnatpool, name)

	if err != nil {
		return nil, err
	}

	return &snatpools, nil
}

func (b *BigIP) CreateSnatpool(name, partition string, members []string) error {
	snatpool := &Snatpool{
		Name:      name,
		Partition: partition,
		Members:   members,
	}
	log.Println("[INFO] Creating snatpool ", snatpool)
	return b.post(snatpool, uriLtm, uriSnatpool)
}

func (b *BigIP) ModifySnatpool(config *Snatpool) error {
	return b.put(config, uriLtm, uriSnatpool)
}

// Get a Snat list  by name. Returns nil if the node does not exist
func (b *BigIP) GetSnatpool(name string) (*Snatpool, error) {
	var snatpool Snatpool
	err, ok := b.getForEntity(&snatpool, uriLtm, uriSnatpool, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &snatpool, nil
}

func (b *BigIP) DeleteSnatpool(name string) error {
	return b.delete(uriLtm, uriSnatpool, name)
}

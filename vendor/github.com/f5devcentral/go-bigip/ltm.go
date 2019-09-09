package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// ServerSSLProfiles
// Documentation: https://devcentral.f5.com/wiki/iControlREST.APIRef_tm_ltm_profile_server-ssl.ashx

// ServerSSLProfiles contains a list of every server-ssl profile on the BIG-IP system.
type ServerSSLProfiles struct {
	ServerSSLProfiles []ServerSSLProfile `json:"items"`
}

// ServerSSLProfile contains information about each server-ssl profile. You can use all
// of these fields when modifying a server-ssl profile.
type ServerSSLProfile struct {
	Name                         string   `json:"name,omitempty"`
	Partition                    string   `json:"partition,omitempty"`
	FullPath                     string   `json:"fullPath,omitempty"`
	Generation                   int      `json:"generation,omitempty"`
	AlertTimeout                 string   `json:"alertTimeout,omitempty"`
	Authenticate                 string   `json:"authenticate,omitempty"`
	AuthenticateDepth            int      `json:"authenticateDepth,omitempty"`
	CaFile                       string   `json:"caFile,omitempty"`
	CacheSize                    int      `json:"cacheSize,omitempty"`
	CacheTimeout                 int      `json:"cacheTimeout,omitempty"`
	Cert                         string   `json:"cert,omitempty"`
	Chain                        string   `json:"chain,omitempty"`
	Ciphers                      string   `json:"ciphers,omitempty"`
	DefaultsFrom                 string   `json:"defaultsFrom,omitempty"`
	ExpireCertResponseControl    string   `json:"expireCertResponseControl,omitempty"`
	GenericAlert                 string   `json:"genericAlert,omitempty"`
	HandshakeTimeout             string   `json:"handshakeTimeout,omitempty"`
	Key                          string   `json:"key,omitempty"`
	ModSslMethods                string   `json:"modSslMethods,omitempty"`
	Mode                         string   `json:"mode,omitempty"`
	TmOptions                    []string `json:"tmOptions,omitempty"`
	Passphrase                   string   `json:"passphrase,omitempty"`
	PeerCertMode                 string   `json:"peerCertMode,omitempty"`
	ProxySsl                     string   `json:"proxySsl,omitempty"`
	RenegotiatePeriod            string   `json:"renegotiatePeriod,omitempty"`
	RenegotiateSize              string   `json:"renegotiateSize,omitempty"`
	Renegotiation                string   `json:"renegotiation,omitempty"`
	RetainCertificate            string   `json:"retainCertificate,omitempty"`
	SecureRenegotiation          string   `json:"secureRenegotiation,omitempty"`
	ServerName                   string   `json:"serverName,omitempty"`
	SessionMirroring             string   `json:"sessionMirroring,omitempty"`
	SessionTicket                string   `json:"sessionTicket,omitempty"`
	SniDefault                   string   `json:"sniDefault,omitempty"`
	SniRequire                   string   `json:"sniRequire,omitempty"`
	SslForwardProxy              string   `json:"sslForwardProxy,omitempty"`
	SslForwardProxyBypass        string   `json:"sslForwardProxyBypass,omitempty"`
	SslSignHash                  string   `json:"sslSignHash,omitempty"`
	StrictResume                 string   `json:"strictResume,omitempty"`
	UncleanShutdown              string   `json:"uncleanShutdown,omitempty"`
	UntrustedCertResponseControl string   `json:"untrustedCertResponseControl,omitempty"`
}

// ClientSSLProfiles
// Documentation: https://devcentral.f5.com/wiki/iControlREST.APIRef_tm_ltm_profile_client-ssl.ashx

// ClientSSLProfiles contains a list of every client-ssl profile on the BIG-IP system.
type ClientSSLProfiles struct {
	ClientSSLProfiles []ClientSSLProfile `json:"items"`
}

// ClientSSLProfile contains information about each client-ssl profile. You can use all
// of these fields when modifying a client-ssl profile.
type ClientSSLProfile struct {
	Name              string `json:"name,omitempty"`
	Partition         string `json:"partition,omitempty"`
	FullPath          string `json:"fullPath,omitempty"`
	Generation        int    `json:"generation,omitempty"`
	AlertTimeout      string `json:"alertTimeout,omitempty"`
	AllowNonSsl       string `json:"allowNonSsl,omitempty"`
	Authenticate      string `json:"authenticate,omitempty"`
	AuthenticateDepth int    `json:"authenticateDepth,omitempty"`
	CaFile            string `json:"caFile,omitempty"`
	CacheSize         int    `json:"cacheSize,omitempty"`
	CacheTimeout      int    `json:"cacheTimeout,omitempty"`
	Cert              string `json:"cert,omitempty"`
	CertKeyChain      []struct {
		Name       string `json:"name,omitempty"`
		Cert       string `json:"cert,omitempty"`
		Chain      string `json:"chain,omitempty"`
		Key        string `json:"key,omitempty"`
		Passphrase string `json:"passphrase,omitempty"`
	} `json:"certKeyChain,omitempty"`
	CertExtensionIncludes           []string `json:"certExtensionIncludes,omitempty"`
	CertLifespan                    int      `json:"certLifespan,omitempty"`
	CertLookupByIpaddrPort          string   `json:"certLookupByIpaddrPort,omitempty"`
	Chain                           string   `json:"chain,omitempty"`
	Ciphers                         string   `json:"ciphers,omitempty"`
	ClientCertCa                    string   `json:"clientCertCa,omitempty"`
	CrlFile                         string   `json:"crlFile,omitempty"`
	DefaultsFrom                    string   `json:"defaultsFrom,omitempty"`
	ForwardProxyBypassDefaultAction string   `json:"forwardProxyBypassDefaultAction,omitempty"`
	GenericAlert                    string   `json:"genericAlert,omitempty"`
	HandshakeTimeout                string   `json:"handshakeTimeout,omitempty"`
	InheritCertkeychain             string   `json:"inheritCertkeychain,omitempty"`
	Key                             string   `json:"key,omitempty"`
	ModSslMethods                   string   `json:"modSslMethods,omitempty"`
	Mode                            string   `json:"mode,omitempty"`
	TmOptions                       []string `json:"tmOptions,omitempty"`
	Passphrase                      string   `json:"passphrase,omitempty"`
	PeerCertMode                    string   `json:"peerCertMode,omitempty"`
	ProxyCaCert                     string   `json:"proxyCaCert,omitempty"`
	ProxyCaKey                      string   `json:"proxyCaKey,omitempty"`
	ProxyCaPassphrase               string   `json:"proxyCaPassphrase,omitempty"`
	ProxySsl                        string   `json:"proxySsl,omitempty"`
	ProxySslPassthrough             string   `json:"proxySslPassthrough,omitempty"`
	RenegotiatePeriod               string   `json:"renegotiatePeriod,omitempty"`
	RenegotiateSize                 string   `json:"renegotiateSize,omitempty"`
	Renegotiation                   string   `json:"renegotiation,omitempty"`
	RetainCertificate               string   `json:"retainCertificate,omitempty"`
	SecureRenegotiation             string   `json:"secureRenegotiation,omitempty"`
	ServerName                      string   `json:"serverName,omitempty"`
	SessionMirroring                string   `json:"sessionMirroring,omitempty"`
	SessionTicket                   string   `json:"sessionTicket,omitempty"`
	SniDefault                      string   `json:"sniDefault,omitempty"`
	SniRequire                      string   `json:"sniRequire,omitempty"`
	SslForwardProxy                 string   `json:"sslForwardProxy,omitempty"`
	SslForwardProxyBypass           string   `json:"sslForwardProxyBypass,omitempty"`
	SslSignHash                     string   `json:"sslSignHash,omitempty"`
	StrictResume                    string   `json:"strictResume,omitempty"`
	UncleanShutdown                 string   `json:"uncleanShutdown,omitempty"`
}

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
	Description     string `json:"description,omitempty"`
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
	FQDN            struct {
		AddressFamily string `json:"addressFamily,omitempty"`
		AutoPopulate  string `json:"autopopulate,omitempty"`
		DownInterval  int    `json:"downInterval,omitempty"`
		Interval      string `json:"interval,omitempty"`
		Name          string `json:"tmName,omitempty"`
	} `json:"fqdn,omitempty"`
}

// DataGroups contains a list of data groups on the BIG-IP system.
type DataGroups struct {
	DataGroups []DataGroup `json:"items"`
}

// DataGroups contains information about each data group.
type DataGroup struct {
	Name       string
	Partition  string
	FullPath   string
	Generation int
	Type       string
	Records    []DataGroupRecord
}

type DataGroupRecord struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}

type dataGroupDTO struct {
	Name       string            `json:"name,omitempty"`
	Partition  string            `json:"partition,omitempty"`
	FullPath   string            `json:"fullPath,omitempty"`
	Generation int               `json:"generation,omitempty"`
	Type       string            `json:"type,omitempty"`
	Records    []DataGroupRecord `json:"records,omitempty"`
}

func (p *DataGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(dataGroupDTO{
		Name:       p.Name,
		Partition:  p.Partition,
		FullPath:   p.FullPath,
		Generation: p.Generation,
		Type:       p.Type,
		Records:    p.Records,
	})
}

func (p *DataGroup) UnmarshalJSON(b []byte) error {
	var dto dataGroupDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	p.Name = dto.Name
	p.Partition = dto.Partition
	p.Type = dto.Type
	p.FullPath = dto.FullPath
	p.Generation = dto.Generation
	p.Records = dto.Records
	return nil
}

// SnatPools contains a list of every snatpool on the BIG-IP system.
type SnatPools struct {
	SnatPools []SnatPool `json:"items"`
}

// SnatPool contains information about each individual snatpool. You can use all
// of these fields when modifying a snatpool.
type SnatPool struct {
	Name       string   `json:"name,omitempty"`
	Partition  string   `json:"partition,omitempty"`
	FullPath   string   `json:"fullPath,omitempty"`
	Generation int      `json:"generation,omitempty"`
	Members    []string `json:"members,omitempty"`
}

// Pools contains a list of pools on the BIG-IP system.
type Pools struct {
	Pools []Pool `json:"items"`
}

// Pool contains information about each pool. You can use all of these
// fields when modifying a pool.
type Pool struct {
	Name                   string `json:"name,omitempty"`
	Partition              string `json:"partition,omitempty"`
	FullPath               string `json:"fullPath,omitempty"`
	Description            string `json:"description,omitempty"`
	Generation             int    `json:"generation,omitempty"`
	AllowNAT               string `json:"allowNat,omitempty"`
	AllowSNAT              string `json:"allowSnat,omitempty"`
	IgnorePersistedWeight  string `json:"ignorePersistedWeight,omitempty"`
	IPTOSToClient          string `json:"ipTosToClient,omitempty"`
	IPTOSToServer          string `json:"ipTosToServer,omitempty"`
	LinkQoSToClient        string `json:"linkQosToClient,omitempty"`
	LinkQoSToServer        string `json:"linkQosToServer,omitempty"`
	LoadBalancingMode      string `json:"loadBalancingMode,omitempty"`
	MinActiveMembers       int    `json:"minActiveMembers,omitempty"`
	MinUpMembers           int    `json:"minUpMembers,omitempty"`
	MinUpMembersAction     string `json:"minUpMembersAction,omitempty"`
	MinUpMembersChecking   string `json:"minUpMembersChecking,omitempty"`
	Monitor                string `json:"monitor,omitempty"`
	QueueDepthLimit        int    `json:"queueDepthLimit,omitempty"`
	QueueOnConnectionLimit string `json:"queueOnConnectionLimit,omitempty"`
	QueueTimeLimit         int    `json:"queueTimeLimit,omitempty"`
	ReselectTries          int    `json:"reselectTries"`
	ServiceDownAction      string `json:"serviceDownAction,omitempty"`
	SlowRampTime           int    `json:"slowRampTime"`
}

// Pool Members contains a list of pool members within a pool on the BIG-IP system.
type PoolMembers struct {
	PoolMembers []PoolMember `json:"items"`
}

// poolMember is used only when adding members to a pool.
type poolMember struct {
	Name string `json:"name"`
}

// poolMembers is used only when modifying members on a pool.
type poolMembers struct {
	Members []PoolMember `json:"members"`
}

// Pool Member contains information about each individual member in a pool. You can use all
// of these fields when modifying a pool member.
type PoolMember struct {
	Name            string `json:"name,omitempty"`
	Partition       string `json:"partition,omitempty"`
	FullPath        string `json:"fullPath,omitempty"`
	Generation      int    `json:"generation,omitempty"`
	Address         string `json:"address,omitempty"`
	ConnectionLimit int    `json:"connectionLimit,omitempty"`
	DynamicRatio    int    `json:"dynamicRatio,omitempty"`
	InheritProfile  string `json:"inheritProfile,omitempty"`
	Logging         string `json:"logging,omitempty"`
	Monitor         string `json:"monitor,omitempty"`
	PriorityGroup   int    `json:"priorityGroup,omitempty"`
	RateLimit       string `json:"rateLimit,omitempty"`
	Ratio           int    `json:"ratio,omitempty"`
	Session         string `json:"session,omitempty"`
	State           string `json:"state,omitempty"`
}

// Pool transfer object so we can mask the bool data munging
type poolDTO struct {
	Name                   string `json:"name,omitempty"`
	Partition              string `json:"partition,omitempty"`
	FullPath               string `json:"fullPath,omitempty"`
	Description            string `json:"description,omitempty"`
	Generation             int    `json:"generation,omitempty"`
	AllowNAT               string `json:"allowNat,omitempty"`
	AllowSNAT              string `json:"allowSnat,omitempty"`
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
	Monitor                string `json:"monitor,omitempty"`
	QueueDepthLimit        int    `json:"queueDepthLimit,omitempty"`
	QueueOnConnectionLimit string `json:"queueOnConnectionLimit,omitempty"`
	QueueTimeLimit         int    `json:"queueTimeLimit,omitempty"`
	ReselectTries          int    `json:"reselectTries"`
	ServiceDownAction      string `json:"serviceDownAction,omitempty"`
	SlowRampTime           int    `json:"slowRampTime"`
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

// PersistenceProfiles contains of list of persistence profiles
type PersistenceProfiles struct {
	PersistenceProfiles []PersistenceProfile `json:"items"`
}

// PersistenceProfile is a base for all persistence profiles
type PersistenceProfile struct {
	AppService              string `json:"appService,omitempty"`
	DefaultsFrom            string `json:"defaultsFrom,omitempty"`
	Description             string `json:"description,omitempty"`
	FullPath                string `json:"fullPath,omitempty"`
	MatchAcrossPools        string `json:"matchAcrossPools,omitempty"`
	MatchAcrossServices     string `json:"matchAcrossServices,omitempty"`
	MatchAcrossVirtuals     string `json:"matchAcrossVirtuals,omitempty"`
	Method                  string `json:"method,omitempty"`
	Mirror                  string `json:"mirror,omitempty"`
	Mode                    string `json:"mode,omitempty"`
	Name                    string `json:"name,omitempty"`
	OverrideConnectionLimit string `json:"overrideConnectionLimit,omitempty"`
	Partition               string `json:"partition,omitempty"`
	TmPartition             string `json:"tmPartition,omitempty"`
	Timeout                 string `json:"timeout,omitempty"`
}

// CookiePersistenceProfiles contains a list of all cookies profiles
type CookiePersistenceProfiles struct {
	CookiePersistenceProfiles []CookiePersistenceProfile `json:"items"`
}

// CookiePersistenceProfile Defines a single cookie profile
type CookiePersistenceProfile struct {
	PersistenceProfile
	AlwaysSend                 string `json:"alwaysSend,omitempty"`
	CookieEncryption           string `json:"cookieEncryption,omitempty"`
	CookieEncryptionPassphrase string `json:"cookieEncryptionPassphrase,omitempty"`
	CookieName                 string `json:"cookieName,omitempty"`
	Expiration                 string `json:"expiration,omitempty"`
	HashLength                 int    `json:"hashLength,omitempty"`
	HashOffset                 int    `json:"hashOffset,omitempty"`
	HTTPOnly                   string `json:"httponly,omitempty"`
	Secure                     string `json:"secure,omitempty"`
}

// DestAddrPersistenceProfiles contains a list of all dest-addr profiles
type DestAddrPersistenceProfiles struct {
	DestAddrPersistenceProfiles []DestAddrPersistenceProfile `json:"items"`
}

// DestAddrPersistenceProfile Defines a single dest-addr profile
type DestAddrPersistenceProfile struct {
	PersistenceProfile
	HashAlgorithm string `json:"hashAlgorithm,omitempty"`
	Mask          string `json:"mask,omitempty"`
}

// HashPersistenceProfiles contains a list of all hash profiles
type HashPersistenceProfiles struct {
	HashPersistenceProfiles []HashPersistenceProfile `json:"items"`
}

// HashPersistenceProfile Defines a single hash profile
type HashPersistenceProfile struct {
	PersistenceProfile
	HashAlgorithm    string `json:"hashAlgorithm,omitempty"`
	HashBufferLimit  int    `json:"hashBufferLimit,omitempty"`
	HashEndPattern   int    `json:"hashEndPattern,omitempty"`
	HashLength       int    `json:"hashLength,omitempty"`
	HashOffset       int    `json:"hashOffset,omitempty"`
	HashStartPattern int    `json:"hashStartPattern,omitempty"`
}

// HostPersistenceProfiles contains a list of all host profiles
type HostPersistenceProfiles struct {
	HostPersistenceProfiles []HostPersistenceProfile `json:"items"`
}

// HostPersistenceProfile Defines a single host profile
type HostPersistenceProfile struct {
	PersistenceProfile
}

// MSRDPPersistenceProfiles contains a list of all msrdp profiles
type MSRDPPersistenceProfiles struct {
	MSRDPPersistenceProfiles []MSRDPPersistenceProfile `json:"items"`
}

// MSRDPPersistenceProfile Defines a single msrdp profile
type MSRDPPersistenceProfile struct {
	PersistenceProfile
	HasSessionDir string `json:"hasSessionDir,omitempty"`
}

// SIPPersistenceProfiles contains a list of all sip profiles
type SIPPersistenceProfiles struct {
	SIPPersistenceProfiles []SIPPersistenceProfile `json:"items"`
}

// SIPPersistenceProfile Defines a single sip profile
type SIPPersistenceProfile struct {
	PersistenceProfile
	SIPInfo string `json:"sipInfo,omitempty"`
}

// SourceAddrPersistenceProfiles contains a list of all source-addr profiles
type SourceAddrPersistenceProfiles struct {
	SourceAddrPersistenceProfiles []SourceAddrPersistenceProfile `json:"items"`
}

// SourceAddrPersistenceProfile Defines a single source-addr profile
type SourceAddrPersistenceProfile struct {
	PersistenceProfile
	HashAlgorithm   string `json:"hashAlgorithm,omitempty"`
	MapProxies      string `json:"mapProxies,omitempty"`
	MapProxyAddress string `json:"mapProxyAddress,omitempty"`
	MapProxyClass   string `json:"mapProxyClass,omitempty"`
	Mask            string `json:"mask,omitempty"`
}

// SSLPersistenceProfiles contains a list of all ssl profiles
type SSLPersistenceProfiles struct {
	SSLPersistenceProfiles []SSLPersistenceProfile `json:"items"`
}

// SSLPersistenceProfile Defines a single ssl profile
type SSLPersistenceProfile struct {
	PersistenceProfile
}

// UniversalPersistenceProfiles contains a list of all universal profiles
type UniversalPersistenceProfiles struct {
	SSLPersistenceProfiles []UniversalPersistenceProfile `json:"items"`
}

// UniversalPersistenceProfile Defines a single universal profile
type UniversalPersistenceProfile struct {
	PersistenceProfile
	Rule string `json:"rule,omitempty"`
}

// VirtualServers contains a list of all virtual servers on the BIG-IP system.
type VirtualServers struct {
	VirtualServers []VirtualServer `json:"items"`
}

// VirtualServer contains information about each individual virtual server.
type VirtualServer struct {
	Name                       string `json:"name,omitempty"`
	Partition                  string `json:"partition,omitempty"`
	FullPath                   string `json:"fullPath,omitempty"`
	Generation                 int    `json:"generation,omitempty"`
	Description                string `json:"description,omitempty"`
	AddressStatus              string `json:"addressStatus,omitempty"`
	AutoLastHop                string `json:"autoLastHop,omitempty"`
	CMPEnabled                 string `json:"cmpEnabled,omitempty"`
	ConnectionLimit            int    `json:"connectionLimit,omitempty"`
	Destination                string `json:"destination,omitempty"`
	Enabled                    bool   `json:"enabled,omitempty"`
	Disabled                   bool   `json:"disabled,omitempty"`
	GTMScore                   int    `json:"gtmScore,omitempty"`
	FallbackPersistenceProfile string `json:"fallbackPersistence,omitempty"`
	IPProtocol                 string `json:"ipProtocol,omitempty"`
	Mask                       string `json:"mask,omitempty"`
	Mirror                     string `json:"mirror,omitempty"`
	MobileAppTunnel            string `json:"mobileAppTunnel,omitempty"`
	NAT64                      string `json:"nat64,omitempty"`
	Pool                       string `json:"pool,omitempty"`
	RateLimit                  string `json:"rateLimit,omitempty"`
	RateLimitDestinationMask   int    `json:"rateLimitDstMask,omitempty"`
	RateLimitMode              string `json:"rateLimitMode,omitempty"`
	RateLimitSourceMask        int    `json:"rateLimitSrcMask,omitempty"`
	Source                     string `json:"source,omitempty"`
	SourceAddressTranslation   struct {
		Type string `json:"type,omitempty"`
		Pool string `json:"pool,omitempty"`
	} `json:"sourceAddressTranslation,omitempty"`
	SourcePort          string    `json:"sourcePort,omitempty"`
	SYNCookieStatus     string    `json:"synCookieStatus,omitempty"`
	TranslateAddress    string    `json:"translateAddress,omitempty"`
	TranslatePort       string    `json:"translatePort,omitempty"`
	VlansEnabled        bool      `json:"vlansEnabled,omitempty"`
	VSIndex             int       `json:"vsIndex,omitempty"`
	Vlans               []string  `json:"vlans,omitempty"`
	Rules               []string  `json:"rules,omitempty"`
	PersistenceProfiles []Profile `json:"persist,omitempty"`
	Profiles            []Profile `json:"profiles,omitempty"`
	Policies            []string  `json:"policies,omitempty"`
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
	PolicyRef []VirtualServerPolicy `json:"items"`
}

type PolicyPublish struct {
	Name        string
	Command     string
}
type PolicyPublishDTO struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
}
func (p *PolicyPublish) MarshalJSON() ([]byte, error) {
	return json.Marshal(PolicyPublishDTO{
		Name:        p.Name,
		Command:     p.Command,
	})
}
func (p *PolicyPublish) UnmarshalJSON(b []byte) error {
	var dto PolicyPublishDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	p.Name = dto.Name
	p.Command = dto.Command
	return nil
}
type Policy struct {
	Name        string
	PublishCopy string
	Partition   string
	FullPath    string
	Controls    []string
	Requires    []string
	Strategy    string
	Rules       []PolicyRule
}
type policyDTO struct {
	Name        string   `json:"name"`
	PublishCopy string   `json:"publishedCopy"`
	Partition   string   `json:"partition,omitempty"`
	Controls    []string `json:"controls,omitempty"`
	Requires    []string `json:"requires,omitempty"`
	Strategy    string   `json:"strategy,omitempty"`
	FullPath    string   `json:"fullPath,omitempty"`
	Rules       struct {
		Items []PolicyRule `json:"items,omitempty"`
	} `json:"rulesReference,omitempty"`
}

func (p *Policy) MarshalJSON() ([]byte, error) {
	return json.Marshal(policyDTO{
		Name:        p.Name,
		PublishCopy: p.PublishCopy,
		Partition:   p.Partition,
		Controls:    p.Controls,
		Requires:    p.Requires,
		Strategy:    p.Strategy,
		FullPath:    p.FullPath,
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
	p.PublishCopy = dto.PublishCopy
	p.Partition = dto.Partition
	p.Controls = dto.Controls
	p.Requires = dto.Requires
	p.Strategy = dto.Strategy
	p.Rules = dto.Rules.Items
	p.FullPath = dto.FullPath

	return nil
}

type VirtualServerPolicy struct {
	Name        string
	Partition   string
	FullPath    string
}
type VirtualServerPolicyDTO struct {
	Name        string   `json:"name"`
	Partition   string   `json:"partition,omitempty"`
	FullPath    string   `json:"fullPath,omitempty"`
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
	DefaultsFrom   string
	FullPath       string
	Generation     int
	ParentMonitor  string
	Description    string
	Destination    string
	Interval       int
	IPDSCP         int
	ManualResume   string
	Password       string
	ReceiveString  string
	ReceiveDisable string
	Reverse        string
	SendString     string
	TimeUntilUp    int
	Timeout        int
	Transparent    string
	UpInterval     int
	Username       string
	Compatibility  string
	Filename       string
	Mode           string
	Adaptive       string
	AdaptiveLimit  int
}

type monitorDTO struct {
	Name           string `json:"name,omitempty"`
	Partition      string `json:"partition,omitempty"`
	DefaultsFrom   string `json:"defaultsFrom,omitempty"`
	FullPath       string `json:"fullPath,omitempty"`
	Generation     int    `json:"generation,omitempty"`
	ParentMonitor  string `json:"defaultsFrom,omitempty"`
	Description    string `json:"description,omitempty"`
	Destination    string `json:"destination,omitempty"`
	Interval       int    `json:"interval,omitempty"`
	IPDSCP         int    `json:"ipDscp,omitempty"`
	ManualResume   string `json:"manualResume,omitempty"`
	Password       string `json:"password,omitempty"`
	ReceiveString  string `json:"recv,omitempty"`
	ReceiveDisable string `json:"recvDisable,omitempty"`
	Reverse        string `json:"reverse,omitempty"`
	SendString     string `json:"send,omitempty"`
	TimeUntilUp    int    `json:"timeUntilUp,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
	Transparent    string `json:"transparent,omitempty"`
	UpInterval     int    `json:"upInterval,omitempty"`
	Username       string `json:"username,omitempty"`
	Compatibility  string `json:"compatibility,omitempty"`
	Filename       string `json:"filename,omitempty"`
	Mode           string `json:"mode,omitempty"`
	Adaptive       string `json:"adaptive,omitempty"`
	AdaptiveLimit  int    `json:"adaptiveLimit,omitempty"`
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
	Name               string   `json:"name,omitempty"`
	DefaultsFrom       string   `json:"defaultsFrom,omitempty"`
	UriExclude         []string `json:"uriExclude,omitempty"`
	UriInclude         []string `json:"uriInclude,omitempty"`
	ContentTypeInclude []string `json:"contentTypeInclude,omitempty"`
	ContentTypeExclude []string `json:"contentTypeExclude,omitempty"`
}

type Httpcompresss struct {
	Httpcompresss []Httpcompress `json:"items"`
}

type Httpcompress struct {
	Name               string
	DefaultsFrom       string
	UriExclude         []string
	UriInclude         []string
	ContentTypeInclude []string
	ContentTypeExclude []string
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
	return jsonMarshal(dto)
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
	Mirror        string
	SourcePort    string
	Translation   string
	Snatpool      string
	VlansDisabled bool
	Vlans         []string
	Origins       []Originsrecord
}

type snatDTO struct {
	Name          string   `json:"name"`
	Partition     string   `json:"partition,omitempty"`
	FullPath      string   `json:"fullPath,omitempty"`
	AutoLasthop   string   `json:"autoLastHop,omitempty"`
	Mirror        string   `json:"mirror,omitempty"`
	SourcePort    string   `json:"sourePort,omitempty"`
	Translation   string   `json:"translation,omitempty"`
	Snatpool      string   `json:"snatpool,omitempty"`
	Vlans         []string `json:"vlans,omitempty"`
	VlansDisabled bool     `json:"vlansDisabled,omitempty" bool:"disabled"`
	Origins       struct {
		Items []Originsrecord `json:"items,omitempty"`
	} `json:"originsReference,omitempty"`
}

type Originsrecords struct {
	Items []Originsrecord `json:"items,omitempty"`
}

type Originsrecord struct {
	Name        string `json:"name"`
	app_service string `json:"appService,omitempty"`
}

func (p *Snat) MarshalJSON() ([]byte, error) {
	return json.Marshal(snatDTO{
		Name:          p.Name,
		Partition:     p.Partition,
		FullPath:      p.FullPath,
		Mirror:        p.Mirror,
		SourcePort:    p.SourcePort,
		Translation:   p.Translation,
		Snatpool:      p.Snatpool,
		VlansDisabled: p.VlansDisabled,
		Vlans:         p.Vlans,
		Origins: struct {
			Items []Originsrecord `json:"items,omitempty"`
		}{Items: p.Origins},
	})
}

func (p *Snat) UnmarshalJSON(b []byte) error {
	var dto snatDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	p.Name = dto.Name
	p.Partition = dto.Partition
	p.FullPath = dto.FullPath
	p.AutoLasthop = dto.AutoLasthop
	p.Mirror = dto.Mirror
	p.SourcePort = dto.SourcePort
	p.Translation = dto.Translation
	p.Snatpool = dto.Snatpool
	p.VlansDisabled = dto.VlansDisabled
	p.Vlans = dto.Vlans
	p.Origins = dto.Origins.Items

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

// TcpProfiles contains a list of every tcp profile on the BIG-IP system.
type TcpProfiles struct {
	TcpProfiles []TcpProfile `json:"items"`
}

type TcpProfile struct {
	Abc                      string `json:"abc,omitempty"`
	AckOnPush                string `json:"ackOnPush,omitempty"`
	AppService               string `json:"appService,omitempty"`
	AutoProxyBufferSize      string `json:"autoProxyBufferSize,omitempty"`
	AutoReceiveWindowSize    string `json:"autoReceiveWindowSize,omitempty"`
	AutoSendBufferSize       string `json:"autoSendBufferSize,omitempty"`
	CloseWaitTimeout         int    `json:"closeWaitTimeout,omitempty"`
	CmetricsCache            string `json:"cmetricsCache,omitempty"`
	CmetricsCacheTimeout     int    `json:"cmetricsCacheTimeout,omitempty"`
	CongestionControl        string `json:"congestionControl,omitempty"`
	DefaultsFrom             string `json:"defaultsFrom,omitempty"`
	DeferredAccept           string `json:"deferredAccept,omitempty"`
	DelayWindowControl       string `json:"delayWindowControl,omitempty"`
	DelayedAcks              string `json:"delayedAcks,omitempty"`
	Description              string `json:"description,omitempty"`
	Dsack                    string `json:"dsack,omitempty"`
	EarlyRetransmit          string `json:"earlyRetransmit,omitempty"`
	Ecn                      string `json:"ecn,omitempty"`
	EnhancedLossRecovery     string `json:"enhancedLossRecovery,omitempty"`
	FastOpen                 string `json:"fastOpen,omitempty"`
	FastOpenCookieExpiration int    `json:"fastOpenCookieExpiration,omitempty"`
	FinWait_2Timeout         int    `json:"finWait_2Timeout,omitempty"`
	FinWaitTimeout           int    `json:"finWaitTimeout,omitempty"`
	HardwareSynCookie        string `json:"hardwareSynCookie,omitempty"`
	IdleTimeout              int    `json:"idleTimeout,omitempty"`
	InitCwnd                 int    `json:"initCwnd,omitempty"`
	InitRwnd                 int    `json:"initRwnd,omitempty"`
	IpDfMode                 string `json:"ipDfMode,omitempty"`
	IpTosToClient            string `json:"ipTosToClient,omitempty"`
	IpTtlMode                string `json:"ipTtlMode,omitempty"`
	IpTtlV4                  int    `json:"ipTtlV4,omitempty"`
	IpTtlV6                  int    `json:"ipTtlV6,omitempty"`
	KeepAliveInterval        int    `json:"keepAliveInterval,omitempty"`
	LimitedTransmit          string `json:"limitedTransmit,omitempty"`
	LinkQosToClient          string `json:"linkQosToClient,omitempty"`
	MaxRetrans               int    `json:"maxRetrans,omitempty"`
	MaxSegmentSize           int    `json:"maxSegmentSize,omitempty"`
	Md5Signature             string `json:"md5Signature,omitempty"`
	Md5SignaturePassphrase   string `json:"md5SignaturePassphrase,omitempty"`
	MinimumRto               int    `json:"minimumRto,omitempty"`
	Mptcp                    string `json:"mptcp,omitempty"`
	MptcpCsum                string `json:"mptcpCsum,omitempty"`
	MptcpCsumVerify          string `json:"mptcpCsumVerify,omitempty"`
	MptcpDebug               string `json:"mptcpDebug,omitempty"`
	MptcpFallback            string `json:"mptcpFallback,omitempty"`
	MptcpFastjoin            string `json:"mptcpFastjoin,omitempty"`
	MptcpIdleTimeout         int    `json:"mptcpIdleTimeout,omitempty"`
	MptcpJoinMax             int    `json:"mptcpJoinMax,omitempty"`
	MptcpMakeafterbreak      string `json:"mptcpMakeafterbreak,omitempty"`
	MptcpNojoindssack        string `json:"mptcpNojoindssack,omitempty"`
	MptcpRtomax              int    `json:"mptcpRtomax,omitempty"`
	MptcpRxmitmin            int    `json:"mptcpRxmitmin,omitempty"`
	MptcpSubflowmax          int    `json:"mptcpSubflowmax,omitempty"`
	MptcpTimeout             int    `json:"mptcpTimeout,omitempty"`
	Nagle                    string `json:"nagle,omitempty"`
	Name                     string `json:"name,omitempty"`
	TmPartition              string `json:"tmPartition,omitempty"`
	PktLossIgnoreBurst       int    `json:"pktLossIgnoreBurst,omitempty"`
	PktLossIgnoreRate        int    `json:"pktLossIgnoreRate,omitempty"`
	ProxyBufferHigh          int    `json:"proxyBufferHigh,omitempty"`
	ProxyBufferLow           int    `json:"proxyBufferLow,omitempty"`
	ProxyMss                 string `json:"proxyMss,omitempty"`
	ProxyOptions             string `json:"proxyOptions,omitempty"`
	RatePace                 string `json:"ratePace,omitempty"`
	RatePaceMaxRate          int    `json:"ratePaceMaxRate,omitempty"`
	ReceiveWindowSize        int    `json:"receiveWindowSize,omitempty"`
	ResetOnTimeout           string `json:"resetOnTimeout,omitempty"`
	RexmtThresh              int    `json:"rexmtThresh,omitempty"`
	SelectiveAcks            string `json:"selectiveAcks,omitempty"`
	SelectiveNack            string `json:"selectiveNack,omitempty"`
	SendBufferSize           int    `json:"sendBufferSize,omitempty"`
	SlowStart                string `json:"slowStart,omitempty"`
	SynCookieEnable          string `json:"synCookieEnable,omitempty"`
	SynCookieWhitelist       string `json:"synCookieWhitelist,omitempty"`
	SynMaxRetrans            int    `json:"synMaxRetrans,omitempty"`
	SynRtoBase               int    `json:"synRtoBase,omitempty"`
	TailLossProbe            string `json:"tailLossProbe,omitempty"`
	TcpOptions               string `json:"tcpOptions,omitempty"`
	TimeWaitRecycle          string `json:"timeWaitRecycle,omitempty"`
	TimeWaitTimeout          string `json:"timeWaitTimeout,omitempty"`
	Timestamps               string `json:"timestamps,omitempty"`
	VerifiedAccept           string `json:"verifiedAccept,omitempty"`
}

// UdpProfiles contains a list of every tcp profile on the BIG-IP system.
type UdpProfiles struct {
	UdpProfiles []UdpProfile `json:"items"`
}

type UdpProfile struct {
	AllowNoPayload        string `json:"allowNoPayload,omitempty"`
	AppService            string `json:"appService,omitempty"`
	BufferMaxBytes        int    `json:"bufferMaxBytes,omitempty"`
	BufferMaxPackets      int    `json:"bufferMaxPackets,omitempty"`
	DatagramLoadBalancing string `json:"datagramLoadBalancing,omitempty"`
	DefaultsFrom          string `json:"defaultsFrom,omitempty"`
	Description           string `json:"description,omitempty"`
	IdleTimeout           string `json:"idleTimeout,omitempty"`
	IpDfMode              string `json:"ipDfMode,omitempty"`
	IpTosToClient         string `json:"ipTosToClient,omitempty"`
	IpTtlMode             string `json:"ipTtlMode,omitempty"`
	IpTtlV4               int    `json:"ipTtlV4,omitempty"`
	IpTtlV6               int    `json:"ipTtlV6,omitempty"`
	LinkQosToClient       string `json:"linkQosToClient,omitempty"`
	Name                  string `json:"name,omitempty"`
	NoChecksum            string `json:"noChecksum,omitempty"`
	TmPartition           string `json:"tmPartition,omitempty"`
	ProxyMss              string `json:"proxyMss,omitempty"`
}

type HttpProfiles struct {
	HttpProfiles []HttpProfile `json:"items"`
}

type HttpProfile struct {
	AcceptXff                 string   `json:"acceptXff,omitempty"`
	AppService                string   `json:"appService,omitempty"`
	BasicAuthRealm            string   `json:"basicAuthRealm,omitempty"`
	DefaultsFrom              string   `json:"defaultsFrom,omitempty"`
	Description               string   `json:"description,omitempty"`
	EncryptCookieSecret       string   `json:"encryptCookieSecret,omitempty"`
	EncryptCookies            []string `json:"encryptCookies,omitempty"`
	FallbackHost              string   `json:"fallbackHost,omitempty"`
	FallbackStatusCodes       []string `json:"fallbackStatusCodes,omitempty"`
	HeaderErase               string   `json:"headerErase,omitempty"`
	HeaderInsert              string   `json:"headerInsert,omitempty"`
	InsertXforwardedFor       string   `json:"insertXforwardedFor,omitempty"`
	LwsSeparator              string   `json:"lwsSeparator,omitempty"`
	LwsWidth                  int      `json:"lwsWidth,omitempty"`
	Name                      string   `json:"name,omitempty"`
	OneconnectTransformations string   `json:"oneconnectTransformations,omitempty"`
	TmPartition               string   `json:"tmPartition,omitempty"`
	ProxyType                 string   `json:"proxyType,omitempty"`
	RedirectRewrite           string   `json:"redirectRewrite,omitempty"`
	RequestChunking           string   `json:"requestChunking,omitempty"`
	ResponseChunking          string   `json:"responseChunking,omitempty"`
	ResponseHeadersPermitted  []string `json:"responseHeadersPermitted,omitempty"`
	ServerAgentName           string   `json:"serverAgentName,omitempty"`
	ViaHostName               string   `json:"viaHostName,omitempty"`
	ViaRequest                string   `json:"viaRequest,omitempty"`
	ViaResponse               string   `json:"viaResponse,omitempty"`
	XffAlternativeNames       []string `json:"xffAlternativeNames,omitempty"`
}

type OneconnectProfiles struct {
	OneconnectProfiles []OneconnectProfile `json:"items"`
}

type OneconnectProfile struct {
	AppService          string `json:"appService,omitempty"`
	DefaultsFrom        string `json:"defaultsFrom,omitempty"`
	Description         string `json:"description,omitempty"`
	IdleTimeoutOverride string `json:"idleTimeoutOverride,omitempty"`
	LimitType           string `json:"limitType,omitempty"`
	MaxAge              int    `json:"maxAge,omitempty"`
	MaxReuse            int    `json:"maxReuse,omitempty"`
	MaxSize             int    `json:"maxSize,omitempty"`
	Name                string `json:"name,omitempty"`
	TmPartition         string `json:"tmPartition,omitempty"`
	SharePools          string `json:"sharePools,omitempty"`
	SourceMask          string `json:"sourceMask,omitempty"`
}

type HttpCompressionProfiles struct {
	HttpCompressionProfiles []HttpCompressionProfile `json:"items"`
}

type HttpCompressionProfile struct {
	AllowHttp_10       string   `json:"allowHttp_10,omitempty"`
	AppService         string   `json:"appService,omitempty"`
	BrowserWorkarounds string   `json:"browserWorkarounds,omitempty"`
	BufferSize         int      `json:"bufferSize,omitempty"`
	ContentTypeExclude []string `json:"contentTypeExclude,omitempty"`
	ContentTypeInclude []string `json:"contentTypeInclude,omitempty"`
	CpuSaver           string   `json:"cpuSaver,omitempty"`
	CpuSaverHigh       int      `json:"cpuSaverHigh,omitempty"`
	CpuSaverLow        int      `json:"cpuSaverLow,omitempty"`
	DefaultsFrom       string   `json:"defaultsFrom,omitempty"`
	Description        string   `json:"description,omitempty"`
	GzipLevel          int      `json:"gzipLevel,omitempty"`
	GzipMemoryLevel    int      `json:"gzipMemoryLevel,omitempty"`
	GzipWindowSize     int      `json:"gzipWindowSize,omitempty"`
	KeepAcceptEncoding string   `json:"keepAcceptEncoding,omitempty"`
	MethodPrefer       string   `json:"methodPrefer,omitempty"`
	MinSize            int      `json:"minSize,omitempty"`
	Name               string   `json:"name,omitempty"`
	TmPartition        string   `json:"tmPartition,omitempty"`
	Selective          string   `json:"selective,omitempty"`
	UriExclude         []string `json:"uriExclude,omitempty"`
	UriInclude         []string `json:"uriInclude,omitempty"`
	VaryHeader         string   `json:"varyHeader,omitempty"`
}

const (
	uriLtm            = "ltm"
	uriNode           = "node"
	uriPool           = "pool"
	uriPoolMember     = "members"
	uriProfile        = "profile"
	uriServerSSL      = "server-ssl"
	uriClientSSL      = "client-ssl"
	uriVirtual        = "virtual"
	uriVirtualAddress = "virtual-address"
	uriSnatPool       = "snatpool"
	uriMonitor        = "monitor"
	uriIRule          = "rule"
	uriDatagroup      = "data-group"
	uriInternal       = "internal"
	uriPolicy         = "policy"
	uriOneconnect     = "one-connect"
	uriPersistence    = "persistence"
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
	uriCookie         = "cookie"
	uriDestAddr       = "dest-addr"
	uriHash           = "hash"
	uriHost           = "host"
	uriMSRDP          = "msrdp"
	uriSIP            = "sip"
	uriSourceAddr     = "source-addr"
	uriSSL            = "ssl"
	uriUniversal      = "universal"
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

// SnatPools returns a list of snatpools.
func (b *BigIP) SnatPools() (*SnatPools, error) {
	var snatPools SnatPools
	err, _ := b.getForEntity(&snatPools, uriLtm, uriSnatPool)
	if err != nil {
		return nil, err
	}

	return &snatPools, nil
}

// CreateSnatPool adds a new snatpool to the BIG-IP system.
func (b *BigIP) CreateSnatPool(name string, members []string) error {
	config := &SnatPool{
		Name:    name,
		Members: members,
	}

	return b.post(config, uriLtm, uriSnatPool)
}

// AddSnatPool adds a new snatpool by config to the BIG-IP system.
func (b *BigIP) AddSnatPool(config *SnatPool) error {

	return b.post(config, uriLtm, uriSnatPool)
}

// GetSnatPool retrieves a SnatPool by name. Returns nil if the snatpool does not exist
func (b *BigIP) GetSnatPool(name string) (*SnatPool, error) {
	var snatPool SnatPool
	err, _ := b.getForEntity(&snatPool, uriLtm, uriSnatPool, name)
	if err != nil {
		return nil, err
	}
	//if !ok {
	//	return nil, nil
	//}

	return &snatPool, nil
}

// DeleteSnatPool removes a snatpool.
func (b *BigIP) DeleteSnatPool(name string) error {
	return b.delete(uriLtm, uriSnatPool, name)
}

// ModifySnatPool allows you to change any attribute of a snatpool. Fields that
// can be modified are referenced in the Snatpool struct.
func (b *BigIP) ModifySnatPool(name string, config *SnatPool) error {
	return b.put(config, uriLtm, uriSnatPool, name)
}

// ServerSSLProfiles returns a list of server-ssl profiles.
func (b *BigIP) ServerSSLProfiles() (*ServerSSLProfiles, error) {
	var serverSSLProfiles ServerSSLProfiles
	err, _ := b.getForEntity(&serverSSLProfiles, uriLtm, uriProfile, uriServerSSL)
	if err != nil {
		return nil, err
	}

	return &serverSSLProfiles, nil
}

// GetServerSSLProfile gets a server-ssl profile by name. Returns nil if the server-ssl profile does not exist
func (b *BigIP) GetServerSSLProfile(name string) (*ServerSSLProfile, error) {
	var serverSSLProfile ServerSSLProfile
	err, ok := b.getForEntity(&serverSSLProfile, uriLtm, uriProfile, uriServerSSL, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &serverSSLProfile, nil
}

// CreateServerSSLProfile creates a new server-ssl profile on the BIG-IP system.
func (b *BigIP) CreateServerSSLProfile(name string, parent string) error {
	config := &ServerSSLProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriServerSSL)
}

// AddServerSSLProfile adds a new server-ssl profile on the BIG-IP system.
func (b *BigIP) AddServerSSLProfile(config *ServerSSLProfile) error {
	return b.post(config, uriLtm, uriProfile, uriServerSSL)
}

// DeleteServerSSLProfile removes a server-ssl profile.
func (b *BigIP) DeleteServerSSLProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriServerSSL, name)
}

// ModifyServerSSLProfile allows you to change any attribute of a sever-ssl profile.
// Fields that can be modified are referenced in the VirtualServer struct.
func (b *BigIP) ModifyServerSSLProfile(name string, config *ServerSSLProfile) error {
	return b.patch(config, uriLtm, uriProfile, uriServerSSL, name)
}

// ClientSSLProfiles returns a list of client-ssl profiles.
func (b *BigIP) ClientSSLProfiles() (*ClientSSLProfiles, error) {
	var clientSSLProfiles ClientSSLProfiles
	err, _ := b.getForEntity(&clientSSLProfiles, uriLtm, uriProfile, uriClientSSL)
	if err != nil {
		return nil, err
	}

	return &clientSSLProfiles, nil
}

// GetClientSSLProfile gets a client-ssl profile by name. Returns nil if the client-ssl profile does not exist
func (b *BigIP) GetClientSSLProfile(name string) (*ClientSSLProfile, error) {
	var clientSSLProfile ClientSSLProfile
	err, ok := b.getForEntity(&clientSSLProfile, uriLtm, uriProfile, uriClientSSL, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &clientSSLProfile, nil
}

// CreateClientSSLProfile creates a new client-ssl profile on the BIG-IP system.
func (b *BigIP) CreateClientSSLProfile(name string, parent string) error {
	config := &ClientSSLProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriClientSSL)
}

// AddClientSSLProfile adds a new client-ssl profile on the BIG-IP system.
func (b *BigIP) AddClientSSLProfile(config *ClientSSLProfile) error {
	return b.post(config, uriLtm, uriProfile, uriClientSSL)
}

// DeleteClientSSLProfile removes a client-ssl profile.
func (b *BigIP) DeleteClientSSLProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriClientSSL, name)
}

// ModifyClientSSLProfile allows you to change any attribute of a client-ssl profile.
// Fields that can be modified are referenced in the ClientSSLProfile struct.
func (b *BigIP) ModifyClientSSLProfile(name string, config *ClientSSLProfile) error {
	return b.patch(config, uriLtm, uriProfile, uriClientSSL, name)
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

// AddNode adds a new node to the BIG-IP system using the Node Spec
func (b *BigIP) AddNode(config *Node) error {
	return b.post(config, uriLtm, uriNode)
}

// CreateNode adds a new IP based node to the BIG-IP system.
func (b *BigIP) CreateNode(name, address, rate_limit string, connection_limit, dynamic_ratio int, monitor, state ,description string) error {
	config := &Node{
		Name:            name,
		Address:         address,
		RateLimit:       rate_limit,
		ConnectionLimit: connection_limit,
		DynamicRatio:    dynamic_ratio,
		Monitor:         monitor,
		State:           state,
		Description:     description,
	}

	return b.post(config, uriLtm, uriNode)
}

// CreateFQDNNode adds a new FQDN based node to the BIG-IP system.
func (b *BigIP) CreateFQDNNode(name, address, rate_limit string, connection_limit, dynamic_ratio int, monitor, state, description, interval, address_family, autopopulate string, downinterval int) error {
	config := &Node{
		Name:            name,
		RateLimit:       rate_limit,
		ConnectionLimit: connection_limit,
		DynamicRatio:    dynamic_ratio,
		Monitor:         monitor,
		State:           state,
		Description:     description,
	}
	config.FQDN.Name = address
	config.FQDN.Interval = interval
	config.FQDN.AddressFamily = address_family
	config.FQDN.AutoPopulate = autopopulate
	config.FQDN.DownInterval = downinterval

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

// InternalDataGroups returns a list of internal data groups.
func (b *BigIP) InternalDataGroups() (*DataGroups, error) {
	var dataGroups DataGroups
	err, _ := b.getForEntity(&dataGroups, uriLtm, uriDatagroup, uriInternal)
	if err != nil {
		return nil, err
	}

	return &dataGroups, nil
}

func (b *BigIP) AddInternalDataGroup(config *DataGroup) error {
	return b.post(config, uriLtm, uriDatagroup, uriInternal)
}

func (b *BigIP) DeleteInternalDataGroup(name string) error {
	return b.delete(uriLtm, uriDatagroup, uriInternal, name)
}

// Modify a named internal data group, REPLACING all the records
func (b *BigIP) ModifyInternalDataGroupRecords(name string, records []DataGroupRecord) error {
	config := &DataGroup{
		Records: records,
	}
	return b.put(config, uriLtm, uriDatagroup, uriInternal, name)
}

// Get an internal data group by name, returns nil if the data group does not exist
func (b *BigIP) GetInternalDataGroup(name string) (*DataGroup, error) {
	var datagroup DataGroup
	err, ok := b.getForEntity(&datagroup, uriLtm, uriDatagroup, uriInternal, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &datagroup, nil
}

// Get the internal data group records for a named internal data group
func (b *BigIP) GetInternalDataGroupRecords(name string) (*[]DataGroupRecord, error) {
	var dataGroup DataGroup
	err, _ := b.getForEntity(&dataGroup, uriLtm, uriDatagroup, uriInternal, name)
	if err != nil {
		return nil, err
	}

	return &dataGroup.Records, nil
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
func (b *BigIP) PoolMembers(name string) (*PoolMembers, error) {
	var poolMembers PoolMembers
	err, _ := b.getForEntity(&poolMembers, uriLtm, uriPool, name, uriPoolMember)
	if err != nil {
		return nil, err
	}

	return &poolMembers, nil
}

// AddPoolMember adds a node/member to the given pool. <member> must be in the form
// of <node>:<port>, i.e.: "web-server1:443".
func (b *BigIP) AddPoolMember(pool, member string) error {
	config := &poolMember{
		Name: member,
	}

	return b.post(config, uriLtm, uriPool, pool, uriPoolMember)
}

// GetPoolMember returns the details of a member in the specified pool.
func (b *BigIP) GetPoolMember(pool string, member string) (*PoolMember, error) {
	var poolMember PoolMember
	err, ok := b.getForEntity(&poolMember, uriLtm, uriPool, pool, uriPoolMember, member)

	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &poolMember, nil
}

// CreatePoolMember creates a pool member for the specified pool.
func (b *BigIP) CreatePoolMember(pool string, config *PoolMember) error {
	return b.post(config, uriLtm, uriPool, pool, uriPoolMember)
}

// ModifyPoolMember will update the configuration of a particular pool member.
func (b *BigIP) ModifyPoolMember(pool string, config *PoolMember) error {
	member := config.FullPath
	// These fields are not used when modifying a pool member; so omit them.
	config.Name = ""
	config.Partition = ""
	config.FullPath = ""

	// This cannot be modified for an existing pool member.
	config.Address = ""

	return b.put(config, uriLtm, uriPool, pool, uriPoolMember, member)
}

// UpdatePoolMembers does a replace-all-with for the members of a pool.
func (b *BigIP) UpdatePoolMembers(pool string, pm *[]PoolMember) error {
	config := &poolMembers{
		Members: *pm,
	}
	return b.put(config, uriLtm, uriPool, pool)
}

// RemovePoolMember removes a pool member from the specified pool.
func (b *BigIP) RemovePoolMember(pool string, config *PoolMember) error {
	member := config.FullPath
	return b.delete(uriLtm, uriPool, pool, uriPoolMember, member)
}

// DeletePoolMember removes a member from the given pool. <member> must be in the form
// of <node>:<port>, i.e.: "web-server1:443".
func (b *BigIP) DeletePoolMember(pool string, member string) error {
	return b.delete(uriLtm, uriPool, pool, uriPoolMember, member)
}

// PoolMemberStatus changes the status of a pool member. <state> can be either
// "enable" or "disable". <member> must be in the form of <node>:<port>,
// i.e.: "web-server1:443".
func (b *BigIP) PoolMemberStatus(pool string, member string, state string) error {
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

	return b.put(config, uriLtm, uriPool, pool, uriPoolMember, member)
}

// CreatePool adds a new pool to the BIG-IP system by name.
func (b *BigIP) CreatePool(name string) error {
	config := &Pool{
		Name: name,
	}

	return b.post(config, uriLtm, uriPool)
}

// AddPool creates a new pool on the BIG-IP system.
func (b *BigIP) AddPool(config *Pool) error {
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
func (b *BigIP) CreateVirtualServer(name, destination, mask, pool string, vlans_enabled bool, port int, translate_address, translate_port string) error {
	subnetMask := cidr[mask]

	if strings.Contains(mask, ".") {
		subnetMask = mask
	}

	config := &VirtualServer{
		Name:             name,
		Destination:      fmt.Sprintf("%s:%d", destination, port),
		Mask:             subnetMask,
		Pool:             pool,
		TranslateAddress: translate_address,
		TranslatePort:    translate_port,
	}

	return b.post(config, uriLtm, uriVirtual)
}

// AddVirtualServer adds a new virtual server by config to the BIG-IP system.
func (b *BigIP) AddVirtualServer(config *VirtualServer) error {
	return b.post(config, uriLtm, uriVirtual)
}

// GetVirtualServer retrieves a virtual server by name. Returns nil if the virtual server does not exist
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
	retval := make([]string, 0, len(policies.PolicyRef))
	for _, p := range policies.PolicyRef {
		retval = append(retval, p.Name)
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

// GetVirtualAddress retrieves a VirtualAddress by name. Returns nil if the VirtualAddress does not exist
func (b *BigIP) GetVirtualAddress(vaddr string) (*VirtualAddress, error) {
	var virtualAddress VirtualAddress
	err, _ := b.getForEntity(&virtualAddress, uriLtm, uriVirtualAddress, vaddr)
	if err != nil {
		return nil, err
	}
	return &virtualAddress, nil
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

// Monitors returns a list of all HTTP, HTTPS, Gateway ICMP, ICMP, and TCP monitors.
func (b *BigIP) Monitors() ([]Monitor, error) {
	var monitors []Monitor
	monitorUris := []string{"http", "https", "icmp", "gateway-icmp", "tcp", "tcp-half-open", "ftp"}

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
// "icmp", "gateway icmp", or "tcp".
func (b *BigIP) CreateMonitor(name, parent, defaults_from string, interval, timeout int, send, receive, receive_disable, compatibility string) error {
	config := &Monitor{
		Name:           name,
		ParentMonitor:  parent,
		DefaultsFrom:   defaults_from,
		Interval:       interval,
		Timeout:        timeout,
		SendString:     send,
		ReceiveString:  receive,
		ReceiveDisable: receive_disable,
		Compatibility:  compatibility,
	}

	return b.AddMonitor(config)
}

// Create a monitor by supplying a config
func (b *BigIP) AddMonitor(config *Monitor) error {
	if strings.Contains(config.ParentMonitor, "gateway") {
		config.ParentMonitor = "gatewayIcmp"
	}
	if strings.Contains(config.ParentMonitor, "tcp-half-open") {
		config.ParentMonitor = "tcp-half-open"
	}
	if strings.Contains(config.ParentMonitor, "ftp") {
		config.ParentMonitor = "ftp"
	}
	return b.post(config, uriLtm, uriMonitor, config.ParentMonitor)
}

// GetVirtualServer retrieves a monitor by name. Returns nil if the monitor does not exist
func (b *BigIP) GetMonitor(name string, parent string) (*Monitor, error) {
	// Add a verification that type is an accepted monitor type
	var monitor Monitor
	err, ok := b.getForEntity(&monitor, uriLtm, uriMonitor, parent, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &monitor, nil
}

// DeleteMonitor removes a monitor.
func (b *BigIP) DeleteMonitor(name, parent string) error {
	return b.delete(uriLtm, uriMonitor, parent, name)
}

// ModifyMonitor allows you to change any attribute of a monitor. <parent> must be
// one of "http", "https", "icmp", "gateway icmp", or "tcp". Fields that
// can be modified are referenced in the Monitor struct.
func (b *BigIP) ModifyMonitor(name, parent string, config *Monitor) error {
	if strings.Contains(config.ParentMonitor, "gateway") {
		config.ParentMonitor = "gatewayIcmp"
	}
	if strings.Contains(config.ParentMonitor, "tcp-half-open") {
		config.ParentMonitor = "tcp-half-open"
	}
	if strings.Contains(config.ParentMonitor, "ftp") {
		config.ParentMonitor = "ftp"
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
	values := []string{}
	values = append(values, "Drafts/")
	values = append(values, name)
	// Join three strings into one.
	//result := strings.Join(values, "")
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

func (b *BigIP) PublishPolicy(name, publish string) error {
	config := &PolicyPublish{
		Name: publish,
		Command: "publish",
	}
	values := []string{}
	values = append(values, "~Common~Drafts~")
	values = append(values, name)
	// Join three strings into one.
	result := strings.Join(values, "")

	log.Println("  ================== here in publish ", result, publish)

	return b.post(config, uriLtm, uriPolicy)
}

//Update an existing policy.
func (b *BigIP) UpdatePolicy(name string, p *Policy) error {
	normalizePolicy(p)
	values := []string{}
	values = append(values, "Drafts/")
	values = append(values, name)
	// Join three strings into one.
	result := strings.Join(values, "")
	return b.put(p, uriLtm, uriPolicy, result)
}

//Delete a policy by name.
func (b *BigIP) DeletePolicy(name string) error {
	values := []string{}
	values = append(values, "Drafts/")
	values = append(values, name)
	// Join three strings into one.
	//result := strings.Join(values, "")
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

func (b *BigIP) GetOneconnect(name string) (*Oneconnect, error) {
	var oneconnect Oneconnect
	err, ok := b.getForEntity(&oneconnect, uriLtm, uriProfile, uriOneconnect, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &oneconnect, nil
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

func (b *BigIP) GetTcp(name string) (*Tcp, error) {
	var tcp Tcp
	err, ok := b.getForEntity(&tcp, uriLtm, uriProfile, uriTcp, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &tcp, nil
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

func (b *BigIP) GetFasthttp(name string) (*Fasthttp, error) {
	var fasthttp Fasthttp
	err, _ := b.getForEntity(&fasthttp, uriLtm, uriProfile, uriFasthttp)

	if err != nil {
		return nil, err
	}

	return &fasthttp, nil
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

func (b *BigIP) GetFastl4(name string) (*Fastl4, error) {
	var fastl4 Fastl4
	err, ok := b.getForEntity(&fastl4, uriLtm, uriProfile, uriFastl4, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &fastl4, nil
}

// ===============

func (b *BigIP) CreateHttpcompress(name, defaultsFrom string, uriExclude, uriInclude, contentTypeInclude, contentTypeExclude []string) error {
	httpcompress := &Httpcompress{
		Name:               name,
		DefaultsFrom:       defaultsFrom,
		UriExclude:         uriExclude,
		UriInclude:         uriInclude,
		ContentTypeInclude: contentTypeInclude,
		ContentTypeExclude: contentTypeExclude,
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

func (b *BigIP) GetHttpcompress(name string) (*Httpcompress, error) {
	var httpcompress Httpcompress
	err, ok := b.getForEntity(&httpcompress, uriLtm, uriProfile, uriHttpcompress, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &httpcompress, nil
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

func (b *BigIP) GetHttp2(name string) (*Http2, error) {
	var http2 Http2
	err, ok := b.getForEntity(&http2, uriLtm, uriProfile, uriHttp2, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &http2, nil
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
/*func (b *BigIP) Snats(name string) (*Snats, error) {
	var snats Snats
	err, _ := b.getForEntity(&snats, uriLtm, uriSnat, name)

	if err != nil {
		return nil, err
	}

	return &snats, nil
}*/

/*func (b *BigIP) CreateSnat(name, partition, autoLastHop, sourcePort, translation, snatpool, mirror string, vlansDisabled bool, origins []string) error {
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
} */
func (b *BigIP) CreateSnat(p *Snat) error {
	log.Println(" what is the complete payload    ", p)
	return b.post(p, uriLtm, uriSnat)
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

func (b *BigIP) UpdateSnat(name string, p *Snat) error {
	return b.put(p, uriLtm, uriSnat, name)
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

// CookiePersistenceProfiles returns a list of cookie persist profiles
func (b *BigIP) CookiePersistenceProfiles() (*CookiePersistenceProfiles, error) {
	var cookiePersistenceProfiles CookiePersistenceProfiles
	err, _ := b.getForEntity(&cookiePersistenceProfiles, uriLtm, uriPersistence, uriCookie)
	if err != nil {
		return nil, err
	}

	return &cookiePersistenceProfiles, nil
}

// GetCookiePersistenceProfile gets a single cookie persist profile by name
func (b *BigIP) GetCookiePersistenceProfile(name string) (*CookiePersistenceProfile, error) {
	var cookiePersistenceProfile CookiePersistenceProfile
	err, ok := b.getForEntity(&cookiePersistenceProfile, uriLtm, uriPersistence, uriCookie, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &cookiePersistenceProfile, nil
}

// CreateCookiePersistenceProfile creates a new cookie persist profile on the BIG-IP system.
func (b *BigIP) CreateCookiePersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriCookie)
}

// AddCookiePersistenceProfile adds a cookie persist profile to the BIG-IP system
func (b *BigIP) AddCookiePersistenceProfile(config *CookiePersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriCookie)
}

// DeleteCookiePersistenceProfile removes a cookie persist profile.
func (b *BigIP) DeleteCookiePersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriCookie, name)
}

// ModifyCookiePersistenceProfile allows you to change any attribute of a cookie persist profile.
// Fields that can be modified are referenced in the CookiePersistenceProfile struct.
func (b *BigIP) ModifyCookiePersistenceProfile(name string, config *CookiePersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriCookie, name)
}

// DestAddrPersistenceProfiles returns a list of dest-addr persist profiles
func (b *BigIP) DestAddrPersistenceProfiles() (*DestAddrPersistenceProfiles, error) {
	var destAddrPersistenceProfiles DestAddrPersistenceProfiles
	err, _ := b.getForEntity(&destAddrPersistenceProfiles, uriLtm, uriPersistence, uriDestAddr)
	if err != nil {
		return nil, err
	}

	return &destAddrPersistenceProfiles, nil
}

// GetDestAddrPersistenceProfile gets a single dest-addr persist profile by name
func (b *BigIP) GetDestAddrPersistenceProfile(name string) (*DestAddrPersistenceProfile, error) {
	var destAddrPersistenceProfile DestAddrPersistenceProfile
	err, ok := b.getForEntity(&destAddrPersistenceProfile, uriLtm, uriPersistence, uriDestAddr, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &destAddrPersistenceProfile, nil
}

// CreateDestAddrPersistenceProfile creates a new dest-addr persist profile on the BIG-IP system.
func (b *BigIP) CreateDestAddrPersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriDestAddr)
}

// AddDestAddrPersistenceProfile adds a dest-addr persist profile to the BIG-IP system
func (b *BigIP) AddDestAddrPersistenceProfile(config *DestAddrPersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriDestAddr)
}

// DeleteDestAddrPersistenceProfile removes a dest-addr persist profile.
func (b *BigIP) DeleteDestAddrPersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriDestAddr, name)
}

// ModifyDestAddrPersistenceProfile allows you to change any attribute of a dest-addr persist profile.
// Fields that can be modified are referenced in the DestAddrPersistenceProfile struct.
func (b *BigIP) ModifyDestAddrPersistenceProfile(name string, config *DestAddrPersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriDestAddr, name)
}

// HashPersistenceProfiles returns a list of hash persist profiles
func (b *BigIP) HashPersistenceProfiles() (*HashPersistenceProfiles, error) {
	var hashPersistenceProfiles HashPersistenceProfiles
	err, _ := b.getForEntity(&hashPersistenceProfiles, uriLtm, uriPersistence, uriHash)
	if err != nil {
		return nil, err
	}

	return &hashPersistenceProfiles, nil
}

// GetHashPersistenceProfile gets a single hash persist profile by name
func (b *BigIP) GetHashPersistenceProfile(name string) (*HashPersistenceProfile, error) {
	var hashPersistenceProfile HashPersistenceProfile
	err, ok := b.getForEntity(&hashPersistenceProfile, uriLtm, uriPersistence, uriHash, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &hashPersistenceProfile, nil
}

// CreateHashPersistenceProfile creates a new hash persist profile on the BIG-IP system.
func (b *BigIP) CreateHashPersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriHash)
}

// AddHashPersistenceProfile adds a hash persist profile to the BIG-IP system
func (b *BigIP) AddHashPersistenceProfile(config *HashPersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriHash)
}

// DeleteHashPersistenceProfile removes a dest-addr persist profile.
func (b *BigIP) DeleteHashPersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriHash, name)
}

// ModifyHashPersistenceProfile allows you to change any attribute of a hash persist profile.
// Fields that can be modified are referenced in the HashPersistenceProfile struct.
func (b *BigIP) ModifyHashPersistenceProfile(name string, config *HashPersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriHash, name)
}

// HostPersistenceProfiles returns a list of host persist profiles
func (b *BigIP) HostPersistenceProfiles() (*HostPersistenceProfiles, error) {
	var hostPersistenceProfiles HostPersistenceProfiles
	err, _ := b.getForEntity(&hostPersistenceProfiles, uriLtm, uriPersistence, uriHost)
	if err != nil {
		return nil, err
	}

	return &hostPersistenceProfiles, nil
}

// GetHostPersistenceProfile gets a single host persist profile by name
func (b *BigIP) GetHostPersistenceProfile(name string) (*HostPersistenceProfile, error) {
	var hostPersistenceProfile HostPersistenceProfile
	err, ok := b.getForEntity(&hostPersistenceProfile, uriLtm, uriPersistence, uriHost, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &hostPersistenceProfile, nil
}

// CreateHostPersistenceProfile creates a new host persist profile on the BIG-IP system.
func (b *BigIP) CreateHostPersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriHost)
}

// AddHostPersistenceProfile adds a host persist profile to the BIG-IP system
func (b *BigIP) AddHostPersistenceProfile(config *HostPersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriHost)
}

// DeleteHashHostPersistenceProfile removes a host persist profile.
func (b *BigIP) DeleteHashHostPersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriHost, name)
}

// ModifyHostPersistenceProfile allows you to change any attribute of a host persist profile.
// Fields that can be modified are referenced in the HostPersistenceProfile struct.
func (b *BigIP) ModifyHostPersistenceProfile(name string, config *HostPersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriHost, name)
}

// MSRDPPersistenceProfiles returns a list of msrdp persist profiles
func (b *BigIP) MSRDPPersistenceProfiles() (*MSRDPPersistenceProfiles, error) {
	var msrdpPersistenceProfiles MSRDPPersistenceProfiles
	err, _ := b.getForEntity(&msrdpPersistenceProfiles, uriLtm, uriPersistence, uriMSRDP)
	if err != nil {
		return nil, err
	}

	return &msrdpPersistenceProfiles, nil
}

// GetMSRDPPersistenceProfile gets a single msrdp persist profile by name
func (b *BigIP) GetMSRDPPersistenceProfile(name string) (*MSRDPPersistenceProfile, error) {
	var msrdpPersistenceProfile MSRDPPersistenceProfile
	err, ok := b.getForEntity(&msrdpPersistenceProfile, uriLtm, uriPersistence, uriMSRDP, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &msrdpPersistenceProfile, nil
}

// CreateMSRDPPersistenceProfile creates a new msrdp persist profile on the BIG-IP system.
func (b *BigIP) CreateMSRDPPersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriMSRDP)
}

// AddMSRDPPersistenceProfile adds a msrdp persist profile to the BIG-IP system
func (b *BigIP) AddMSRDPPersistenceProfile(config *MSRDPPersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriMSRDP)
}

// DeleteMSRDPPersistenceProfile removes a msrdp persist profile.
func (b *BigIP) DeleteMSRDPPersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriMSRDP, name)
}

// ModifyMSRDPPersistenceProfile allows you to change any attribute of a msrdp persist profile.
// Fields that can be modified are referenced in the MSRDPPersistenceProfile struct.
func (b *BigIP) ModifyMSRDPPersistenceProfile(name string, config *MSRDPPersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriMSRDP, name)
}

// SIPPersistenceProfiles returns a list of sip persist profiles
func (b *BigIP) SIPPersistenceProfiles() (*SIPPersistenceProfiles, error) {
	var sipPersistenceProfiles SIPPersistenceProfiles
	err, _ := b.getForEntity(&sipPersistenceProfiles, uriLtm, uriPersistence, uriSIP)
	if err != nil {
		return nil, err
	}

	return &sipPersistenceProfiles, nil
}

// GetSIPPersistenceProfile gets a single sip persist profile by name
func (b *BigIP) GetSIPPersistenceProfile(name string) (*SIPPersistenceProfile, error) {
	var sipPersistenceProfile SIPPersistenceProfile
	err, ok := b.getForEntity(&sipPersistenceProfile, uriLtm, uriPersistence, uriSIP, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &sipPersistenceProfile, nil
}

// CreateSIPPersistenceProfile creates a new sip persist profile on the BIG-IP system.
func (b *BigIP) CreateSIPPersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriSIP)
}

// AddSIPPersistenceProfile adds a sip persist profile to the BIG-IP system
func (b *BigIP) AddSIPPersistenceProfile(config *SIPPersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriSIP)
}

// DeleteSIPPersistenceProfile removes a sip persist profile.
func (b *BigIP) DeleteSIPPersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriSIP, name)
}

// ModifySIPPersistenceProfile allows you to change any attribute of a sip persist profile.
// Fields that can be modified are referenced in the SIPPersistenceProfile struct.
func (b *BigIP) ModifySIPPersistenceProfile(name string, config *SIPPersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriSIP, name)
}

// SourceAddrPersistenceProfiles returns a list of source-addr persist profiles
func (b *BigIP) SourceAddrPersistenceProfiles() (*SourceAddrPersistenceProfiles, error) {
	var sourceAddrPersistenceProfiles SourceAddrPersistenceProfiles
	err, _ := b.getForEntity(&sourceAddrPersistenceProfiles, uriLtm, uriPersistence, uriSourceAddr)
	if err != nil {
		return nil, err
	}

	return &sourceAddrPersistenceProfiles, nil
}

// GetSourceAddrPersistenceProfile gets a single source-addr persist profile by name
func (b *BigIP) GetSourceAddrPersistenceProfile(name string) (*SourceAddrPersistenceProfile, error) {
	var sourceAddrPersistenceProfile SourceAddrPersistenceProfile
	err, ok := b.getForEntity(&sourceAddrPersistenceProfile, uriLtm, uriPersistence, uriSourceAddr, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &sourceAddrPersistenceProfile, nil
}

// CreateSourceAddrPersistenceProfile creates a new source-addr persist profile on the BIG-IP system.
func (b *BigIP) CreateSourceAddrPersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriSourceAddr)
}

// AddSourceAddrPersistenceProfile adds a source-addr persist profile to the BIG-IP system
func (b *BigIP) AddSourceAddrPersistenceProfile(config *SourceAddrPersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriSourceAddr)
}

// DeleteSourceAddrPersistenceProfile removes a source-addr persist profile.
func (b *BigIP) DeleteSourceAddrPersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriSourceAddr, name)
}

// ModifySourceAddrPersistenceProfile allows you to change any attribute of a source-addr persist profile.
// Fields that can be modified are referenced in the SourceAddrPersistenceProfile struct.
func (b *BigIP) ModifySourceAddrPersistenceProfile(name string, config *SourceAddrPersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriSourceAddr, name)
}

// SSLPersistenceProfiles returns a list of ssl persist profiles
func (b *BigIP) SSLPersistenceProfiles() (*SSLPersistenceProfiles, error) {
	var sslPersistenceProfiles SSLPersistenceProfiles
	err, _ := b.getForEntity(&sslPersistenceProfiles, uriLtm, uriPersistence, uriSSL)
	if err != nil {
		return nil, err
	}

	return &sslPersistenceProfiles, nil
}

// GetSSLPersistenceProfile gets a single ssl persist profile by name
func (b *BigIP) GetSSLPersistenceProfile(name string) (*SSLPersistenceProfile, error) {
	var sslPersistenceProfile SSLPersistenceProfile
	err, ok := b.getForEntity(&sslPersistenceProfile, uriLtm, uriPersistence, uriSSL, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &sslPersistenceProfile, nil
}

// CreateSSLPersistenceProfile creates a new ssl persist profile on the BIG-IP system.
func (b *BigIP) CreateSSLPersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriSSL)
}

// AddSSLPersistenceProfile adds a ssl persist profile to the BIG-IP system
func (b *BigIP) AddSSLPersistenceProfile(config *SSLPersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriSSL)
}

// DeleteSSLPersistenceProfile removes a ssl persist profile.
func (b *BigIP) DeleteSSLPersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriSSL, name)
}

// ModifySSLPersistenceProfile allows you to change any attribute of a ssl persist profile.
// Fields that can be modified are referenced in the SSLPersistenceProfile struct.
func (b *BigIP) ModifySSLPersistenceProfile(name string, config *SSLPersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriSSL, name)
}

// UniversalPersistenceProfiles returns a list of universal persist profiles
func (b *BigIP) UniversalPersistenceProfiles() (*UniversalPersistenceProfiles, error) {
	var universalPersistenceProfiles UniversalPersistenceProfiles
	err, _ := b.getForEntity(&universalPersistenceProfiles, uriLtm, uriPersistence, uriUniversal)
	if err != nil {
		return nil, err
	}

	return &universalPersistenceProfiles, nil
}

// GetUniversalPersistenceProfile gets a single universal persist profile by name
func (b *BigIP) GetUniversalPersistenceProfile(name string) (*UniversalPersistenceProfile, error) {
	var universalPersistenceProfile UniversalPersistenceProfile
	err, ok := b.getForEntity(&universalPersistenceProfile, uriLtm, uriPersistence, uriUniversal, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &universalPersistenceProfile, nil
}

// CreateUniversalPersistenceProfile creates a new universal persist profile on the BIG-IP system.
func (b *BigIP) CreateUniversalPersistenceProfile(name string, parent string) error {
	config := &PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriPersistence, uriUniversal)
}

// AddUniversalPersistenceProfile adds a universal persist profile to the BIG-IP system
func (b *BigIP) AddUniversalPersistenceProfile(config *UniversalPersistenceProfile) error {
	return b.post(config, uriLtm, uriPersistence, uriUniversal)
}

// DeleteUniversalPersistenceProfile removes a universal persist profile.
func (b *BigIP) DeleteUniversalPersistenceProfile(name string) error {
	return b.delete(uriLtm, uriPersistence, uriUniversal, name)
}

// ModifyUniversalPersistenceProfile allows you to change any attribute of a universal persist profile.
// Fields that can be modified are referenced in the UniversalPersistenceProfile struct.
func (b *BigIP) ModifyUniversalPersistenceProfile(name string, config *UniversalPersistenceProfile) error {
	return b.put(config, uriLtm, uriPersistence, uriUniversal, name)
}

// HttpProfiles returns a list of HTTP profiles
func (b *BigIP) HttpProfiles() (*HttpProfiles, error) {
	var httpProfiles HttpProfiles
	err, _ := b.getForEntity(&httpProfiles, uriLtm, uriProfile, uriHttp)
	if err != nil {
		return nil, err
	}

	return &httpProfiles, nil
}

func (b *BigIP) GetHttpProfile(name string) (*HttpProfile, error) {
	var httpProfile HttpProfile
	err, ok := b.getForEntity(&httpProfile, uriLtm, uriProfile, uriHttp, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &httpProfile, nil
}

// CreateHttpProfile creates a new http profile on the BIG-IP system.
func (b *BigIP) CreateHttpProfile(name string, parent string) error {
	config := &HttpProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriHttp)
}

func (b *BigIP) AddHttpProfile(config *HttpProfile) error {
	return b.post(config, uriLtm, uriProfile, uriHttp)
}

// DeleteHttpProfile removes a http profile.
func (b *BigIP) DeleteHttpProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriHttp, name)
}

// ModifyHttpProfile allows you to change any attribute of a http profile.
// Fields that can be modified are referenced in the HttpProfile struct.
func (b *BigIP) ModifyHttpProfile(name string, config *HttpProfile) error {
	return b.put(config, uriLtm, uriProfile, uriHttp, name)
}

// OneconnectProfiles returns a list of HTTP profiles
func (b *BigIP) OneconnectProfiles() (*OneconnectProfiles, error) {
	var oneconnectProfiles OneconnectProfiles
	err, _ := b.getForEntity(&oneconnectProfiles, uriLtm, uriProfile, uriOneconnect)
	if err != nil {
		return nil, err
	}

	return &oneconnectProfiles, nil
}

func (b *BigIP) GetOneconnectProfile(name string) (*OneconnectProfile, error) {
	var oneconnectProfile OneconnectProfile
	err, ok := b.getForEntity(&oneconnectProfile, uriLtm, uriProfile, uriOneconnect, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &oneconnectProfile, nil
}

// CreateOneconnectProfile creates a new http profile on the BIG-IP system.
func (b *BigIP) CreateOneconnectProfile(name string, parent string) error {
	config := &OneconnectProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriOneconnect)
}

func (b *BigIP) AddOneconnectProfile(config *OneconnectProfile) error {
	return b.post(config, uriLtm, uriProfile, uriOneconnect)
}

// DeleteOneconnectProfile removes a http profile.
func (b *BigIP) DeleteOneconnectProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriOneconnect, name)
}

// ModifyOneconnectProfile allows you to change any attribute of a http profile.
// Fields that can be modified are referenced in the OneconnectProfile struct.
func (b *BigIP) ModifyOneconnectProfile(name string, config *OneconnectProfile) error {
	return b.put(config, uriLtm, uriProfile, uriOneconnect, name)
}

// HttpCompressionProfiles returns a list of HTTP profiles
func (b *BigIP) HttpCompressionProfiles() (*HttpCompressionProfiles, error) {
	var httpCompressionProfiles HttpCompressionProfiles
	err, _ := b.getForEntity(&httpCompressionProfiles, uriLtm, uriProfile, uriHttpcompress)
	if err != nil {
		return nil, err
	}

	return &httpCompressionProfiles, nil
}

func (b *BigIP) GetHttpCompressionProfile(name string) (*HttpCompressionProfile, error) {
	var httpCompressionProfile HttpCompressionProfile
	err, ok := b.getForEntity(&httpCompressionProfile, uriLtm, uriProfile, uriHttpcompress, name)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return &httpCompressionProfile, nil
}

// CreateHttpCompressionProfile creates a new http profile on the BIG-IP system.
func (b *BigIP) CreateHttpCompressionProfile(name string, parent string) error {
	config := &HttpCompressionProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriHttpcompress)
}

func (b *BigIP) AddHttpCompressionProfile(config *HttpCompressionProfile) error {
	return b.post(config, uriLtm, uriProfile, uriHttpcompress)
}

// DeleteHttpCompressionProfile removes a http profile.
func (b *BigIP) DeleteHttpCompressionProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriHttpcompress, name)
}

// ModifyHttpCompressionProfile allows you to change any attribute of a http profile.
// Fields that can be modified are referenced in the HttpCompressionProfile struct.
func (b *BigIP) ModifyHttpCompressionProfile(name string, config *HttpCompressionProfile) error {
	return b.put(config, uriLtm, uriProfile, uriHttpcompress, name)
}

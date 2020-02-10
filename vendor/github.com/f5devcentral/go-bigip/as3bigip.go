package bigip

import (
	//"encoding/json"
	"fmt"
	//	"log"
	//"time"
)

const (
	as3SharedApplication = "Shared"
	DEFAULT_PARTITION    = "Sample_01"
	//uriMgmt              = "mgmt"
	//uriShared            = "shared"
	//uriAppsvcs           = "appsvcs"
	uriDeclare = "declare"
)

type as3JSONWithArbKeys map[string]interface{}

// TODO: Need to remove omitempty tag for the mandatory fields
// as3JSONDeclaration maps to ADC in AS3 Resources
type as3ADC as3JSONWithArbKeys

// as3Tenant maps to Tenant in AS3 Resources
type as3Tenant as3JSONWithArbKeys

// as3Application maps to Application in AS3 Resources
type as3Application as3JSONWithArbKeys

type as3Main struct {
	Class       string      `json:"class"`
	Action      string      `json:"action,omitempty"`
	Persist     bool        `json:"persist,omitempty"`
	Declaration interface{} `json:"declaration"`
}
type As3Main struct {
	Class       string      `json:"class"`
	Action      string      `json:"action,omitempty"`
	Persist     bool        `json:"persist,omitempty"`
	Declaration interface{} `json:"declaration"`
}

// as3Pool maps to Pool in AS3 Resources
type as3Pool struct {
	Class                string          `json:"class,omitempty"`
	Monitors             []string        `json:"monitors,omitempty"`
	Members              []as3PoolMember `json:"members,omitempty"`
	LoadBalancingMode    string          `json:"loadBalancingMode,omitempty"`
	MinimumMembersActive int             `json:"minimumMembersActive,omitempty"`
	ReselectTries        int             `json:"reselectTries,omitempty"`
	ServiceDownAction    string          `json:"serviceDownAction,omitempty"`
	SlowRampTime         int             `json:"slowRampTime,omitempty"`
	MinimumMonitors      int             `json:"minimumMonitors,omitempty"`
}

// as3PoolMember maps to Pool_Member in AS3 Resources
type as3PoolMember struct {
	ServicePort      int      `json:"servicePort,omitempty"`
	ServerAddresses  []string `json:"serverAddresses,omitempty"`
	Enable           bool     `json:"enable,omitempty"`
	ConnectionLimit  int      `json:"connectionLimit,omitempty"`
	RateLimit        int      `json:"rateLimit,omitempty"`
	DynamicRatio     int      `json:"dynamicRatio,omitempty"`
	Ratio            int      `json:"ratio,omitempty"`
	PriorityGroup    int      `json:"priorityGroup,omitempty"`
	AdminState       string   `json:"adminState,omitempty"`
	AddressDiscovery string   `json:"addressDiscovery,omitempty"`
	ShareNodes       bool     `json:"shareNodes,omitempty"`
}

// as3Pool maps to Pool in AS3 Resources
type As3Pool struct {
	Class                string          `json:"class,omitempty"`
	Monitors             []string        `json:"monitors,omitempty"`
	Members              []As3PoolMember `json:"members,omitempty"`
	LoadBalancingMode    string          `json:"loadBalancingMode,omitempty"`
	MinimumMembersActive int             `json:"minimumMembersActive,omitempty"`
	ReselectTries        int             `json:"reselectTries,omitempty"`
	ServiceDownAction    string          `json:"serviceDownAction,omitempty"`
	SlowRampTime         int             `json:"slowRampTime,omitempty"`
	MinimumMonitors      int             `json:"minimumMonitors,omitempty"`
}

// as3PoolMember maps to Pool_Member in AS3 Resources
type As3PoolMember struct {
	ServicePort      int      `json:"servicePort,omitempty"`
	ServerAddresses  []string `json:"serverAddresses,omitempty"`
	Enable           bool     `json:"enable,omitempty"`
	ConnectionLimit  int      `json:"connectionLimit,omitempty"`
	RateLimit        int      `json:"rateLimit,omitempty"`
	DynamicRatio     int      `json:"dynamicRatio,omitempty"`
	Ratio            int      `json:"ratio,omitempty"`
	PriorityGroup    int      `json:"priorityGroup,omitempty"`
	AdminState       string   `json:"adminState,omitempty"`
	AddressDiscovery string   `json:"addressDiscovery,omitempty"`
	ShareNodes       bool     `json:"shareNodes,omitempty"`
}

// as3Service maps to the following in AS3 Resources
// - Service_HTTP
// - Service_HTTPS
// - Service_TCP
// - Service_UDP
type as3Service struct {
	Class                  string   `json:"class,omitempty"`
	VirtualAddresses       []string `json:"virtualAddresses,omitempty"`
	Pool                   string   `json:"pool,omitempty"`
	VirtualPort            int      `json:"virtualPort,omitempty"`
	PersistenceMethods     []string `json:"persistenceMethods,omitempty"`
	ProfileHTTP            string   `json:"profileHTTP,omitempty"`
	Layer4                 string   `json:"layer4,omitempty"`
	ProfileTCP             string   `json:"profileTCP,omitempty"`
	Enable                 bool     `json:"enable,omitempty"`
	MaxConnections         int      `json:"maxConnections,omitempty"`
	Snat                   string   `json:"snat,omitempty"`
	AddressStatus          bool     `json:"addressStatus,omitempty"`
	Mirroring              string   `json:"mirroring,omitempty"`
	LastHop                string   `json:"lastHop,omitempty"`
	TranslateClientPort    bool     `json:"translateClientPort,omitempty"`
	TranslateServerAddress bool     `json:"translateServerAddress,omitempty"`
	TranslateServerPort    bool     `json:"translateServerPort,omitempty"`
	Nat64Enabled           bool     `json:"nat64Enabled,omitempty"`
}

// As3Service maps to the following in AS3 Resources
// - Service_HTTP
// - Service_HTTPS
// - Service_TCP
// - Service_UDP
type As3Service struct {
	Class                  string   `json:"class,omitempty"`
	VirtualAddresses       []string `json:"virtualAddresses,omitempty"`
	Pool                   string   `json:"pool,omitempty"`
	VirtualPort            int      `json:"virtualPort,omitempty"`
	PersistenceMethods     []string `json:"persistenceMethods,omitempty"`
	ProfileHTTP            string   `json:"profileHTTP,omitempty"`
	Layer4                 string   `json:"layer4,omitempty"`
	ProfileTCP             string   `json:"profileTCP,omitempty"`
	Enable                 bool     `json:"enable,omitempty"`
	MaxConnections         int      `json:"maxConnections,omitempty"`
	Snat                   string   `json:"snat,omitempty"`
	AddressStatus          bool     `json:"addressStatus,omitempty"`
	Mirroring              string   `json:"mirroring,omitempty"`
	LastHop                string   `json:"lastHop,omitempty"`
	TranslateClientPort    bool     `json:"translateClientPort,omitempty"`
	TranslateServerAddress bool     `json:"translateServerAddress,omitempty"`
	TranslateServerPort    bool     `json:"translateServerPort,omitempty"`
	Nat64Enabled           bool     `json:"nat64Enabled,omitempty"`
	ServerTLS              string   `json:"serverTLS,omitempty"`
}

type TlsClient struct {
	Class                   string `json:"class,omitempty"`
	AllowExpiredCRL         bool   `json:"allowExpiredCRL,omitempty"`
	AuthenticationFrequency string `json:"authenticationFrequency,omitempty"`
	C3dCertificateAuthority string `json:"c3dCertificateAuthority,omitempty"`
	C3dCertificateLifespan  int    `json:"c3dCertificateLifespan,omitempty"`
	ClientCertificate       string `json:"clientCertificate,omitempty"`
	// C3dCertificateExtensions  string   `json:"c3dCertificateExtensions,omitempty"`
	Ciphers             string `json:"ciphers,omitempty"`
	C3dEnabled          bool   `json:"c3dEnabled,omitempty"`
	IgnoreExpired       bool   `json:"ignoreExpired,omitempty"`
	IgnoreUntrusted     bool   `json:"ignoreUntrusted,omitempty"`
	Label               string `json:"label,omitempty"`
	LdapStartTLS        string `json:"ldapStartTLS,omitempty"`
	Remark              string `json:"remark,omitempty"`
	SendSNI             string `json:"sendSNI,omitempty"`
	ServerName          string `json:"serverName,omitempty"`
	SessionTickets      bool   `json:"sessionTickets,omitempty"`
	TrustCA             string `json:"trustCA,omitempty"`
	ValidateCertificate bool   `json:"validateCertificate,omitempty"`
}

type TlsServer struct {
	Class        string            `json:"class,omitempty"`
	Certificates []As3Certificates `json:"certificates,omitempty"`
}

type As3Certificates struct {
	Certificate string `json:"certificate,omitempty"`
	MatchToSNI  string `json:"matchToSNI,omitempty"`
}

type As3Certificate struct {
	Class       string `json:"class,omitempty"`
	Remark      string `json:"remark,omitempty"`
	Label       string `json:"label,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	ChainCA     string `json:"chainCA,omitempty"`
	Pkcs12      string `json:"pkcs12,omitempty"`
	PrivateKey  string `json:"privateKey,omitempty"`
}

type As3Passphrase struct {
	AllowReuse    bool   `json:"allowReuse,omitempty"`
	Ciphertext    string `json:"ciphertext,omitempty"`
	IgnoreChanges bool   `json:"ignoreChanges,omitempty"`
	MiniJWE       bool   `json:"miniJWE,omitempty"`
	Protected     string `json:"protected,omitempty"`
	ReuseFrom     string `json:"reuseFrom,omitempty"`
}

type typeADC struct {
	Class         string `json:"class"`
	SchemaVersion string `json:"schemaVersion"`
	ID            string `json:"id"`
	Label         string `json:"label"`
	Remark        string `json:"remark"`
	UpdateMode    string `json:"updateMode"`
}

type applicationType struct {
	Class           string   `json:"class"`
	Template        string   `json:"template"`
	MyVirtualServer struct{} `json:"MyVirtualServer"`
	Enable          bool     `json:"enable"`
}
type ApplicationType struct {
	Class           string   `json:"class"`
	Template        string   `json:"template"`
	MyVirtualServer struct{} `json:"MyVirtualServer"`
	Enable          bool     `json:"enable"`
}

type tenantType struct {
	TenantName       string
	AppName          string
	ServiceName      string
	PoolName         string
	VirtualAddresses []string
	PoolMembers      []string
}
type TenantType struct {
	Class         string      `json:"class"`
	MyApplication interface{} `json:"myApplication"`
}

type As3Object struct {
	TenantList []string
}

var As3Tenant As3Object

/*func main() {
	fmt.Printf("Return Struct:%+v\n",generateAS3Declaration())
	j, err:= json.MarshalIndent(generateAS3Declaration(),"","\t")
	fmt.Printf("My result:%+v\nerror:%v\n",string(j), err)
	//fmt.Printf("Error:%v\n",err)
}
*/

func (appClass as3Application) createAs3Pool() {
	// Create Pool
	sharedPool := &as3Pool{}
	sharedPool.Class = "Pool"
	SA := []string{"192.0.2.10", "192.0.2.11"}
	var poolMember as3PoolMember
	poolMember.ServicePort = 80
	poolMember.ServerAddresses = SA
	var memberList []as3PoolMember
	sharedPool.Members = append(memberList, poolMember)
	sharedPool.Monitors = []string{"http"}
	//return sharedPool
	appClass["web_pool2"] = sharedPool
	//fmt.Printf("SharedApp :%+v\n",appClass)
}

func (appClass as3Application) createAs3Service() {
	// Create ServiceClass
	svc := &as3Service{}
	svc.Layer4 = "tcp"
	//svc.Source = "0.0.0.0/0"
	svc.TranslateServerAddress = true
	svc.TranslateServerPort = true
	svc.Class = "Service_HTTP"
	svc.VirtualAddresses = []string{"10.0.2.10"}
	svc.VirtualPort = 80
	svc.Snat = "auto"
	svc.Pool = "web_pool2"
	appClass["Virtual_server"] = svc
	//fmt.Printf("SharedApp :%+v\n",appClass)
	//return appClass

}
func generateAS3Declaration() *as3Main {
	// Create Shared as3Application object
	sharedApp := as3Application{}
	sharedApp["class"] = "Application"
	sharedApp["template"] = "shared"
	sharedApp.createAs3Pool()
	sharedApp.createAs3Service()
	fmt.Printf("SharedApp :%+v\n", sharedApp)
	// Create AS3 Tenant
	tenant := as3Tenant{
		"class":              "Tenant",
		as3SharedApplication: sharedApp,
	}
	as3ADCJson := as3ADC{
		"class":           "ADC",
		"schemaVersion":   "3.15.0",
		DEFAULT_PARTITION: tenant,
	}
	as3JSONDecl := &as3Main{}
	as3JSONDecl.Class = "AS3"
	as3JSONDecl.Declaration = as3ADCJson
	return as3JSONDecl
}

func (b *BigIP) PostAs3Bigip(as3NewJson *As3Main) error {
	//as3NewJson := generateAS3Declaration()
	//	log.Printf("AS3 Struct inside PostAs3Bigip:%+v",as3NewJson)
	return b.post(as3NewJson, uriMgmt, uriShared, uriAppsvcs, uriDeclare)
}

func (b *BigIP) DeleteAs3Bigip(tenantName string) error {
	//tenantName := "Sample_01"
	return b.delete(uriMgmt, uriShared, uriAppsvcs, uriDeclare, tenantName)
}

/*
func (b *BigIP) DeleteAs3Bigip(tenantName string) error {
	//tenantName := "Sample_01"
	return b.delete(uriMgmt, uriShared, uriAppsvcs, uriDeclare, tenantName)
}
*/

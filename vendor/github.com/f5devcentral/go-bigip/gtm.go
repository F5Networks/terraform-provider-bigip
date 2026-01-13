/*
Original work Copyright © 2015 Scott Ware
Modifications Copyright 2019 F5 Networks Inc
Licensed under the Apache License, Version 2.0 (the "License");
You may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.
*/
package bigip

import (
	"encoding/json"
	"log"
)

const (
	uriGtm        = "gtm"
	uriServer     = "server"
	uriDatacenter = "datacenter"
	uriGtmmonitor = "monitor"
	uriPoolA      = "pool/a"
	uriWideIp     = "wideip"
)

type Datacenters struct {
	Datacenters []GTMDatacenter `json:"items"`
}

type GTMWideIP struct {
	Name                              string   `json:"name,omitempty"`
	Partition                         string   `json:"partition,omitempty"`
	FullPath                          string   `json:"fullPath,omitempty"`
	Generation                        int      `json:"generation,omitempty"`
	AppService                        string   `json:"appService,omitempty"`
	Description                       string   `json:"description,omitempty"`
	Disabled                          bool     `json:"disabled,omitempty"`
	Enabled                           bool     `json:"enabled,omitempty"`
	FailureRcode                      string   `json:"failureRcode,omitempty"`
	FailureRcodeResponse              string   `json:"failureRcodeResponse,omitempty"`
	FailureRcodeTTL                   int      `json:"failureRcodeTtl,omitempty"`
	LastResortPool                    string   `json:"lastResortPool,omitempty"`
	LoadBalancingDecisionLogVerbosity []string `json:"loadBalancingDecisionLogVerbosity,omitempty"`
	MinimalResponse                   string   `json:"minimalResponse,omitempty"`
	PersistCidrIpv4                   int      `json:"persistCidrIpv4,omitempty"`
	PersistCidrIpv6                   int      `json:"persistCidrIpv6,omitempty"`
	Persistence                       string   `json:"persistence,omitempty"`
	PoolLbMode                        string   `json:"poolLbMode,omitempty"`
	TopologyPreferEdns0ClientSubnet   string   `json:"topologyPreferEdns0ClientSubnet,omitempty"`
	TTLPersistence                    int      `json:"ttlPersistence,omitempty"`
	Aliases                           []string `json:"aliases,omitempty"`

	// Not in the spec, but returned by the API
	// Setting this field atomically updates all members.
	//Pools *[]GTMWideIPPool `json:"pools,omitempty"`
}

// type Datacenter struct {
// 	Name        string `json:"name,omitempty"`
// 	Description string `json:"description,omitempty"`
// 	Contact     string `json:"contact,omitempty"`
// 	App_service string `json:"appService,omitempty"`
// 	Disabled    bool   `json:"disabled,omitempty"`
// 	Enabled     bool   `json:"enabled,omitempty"`
// 	Prober_pool string `json:"proberPool,omitempty"`
// }

type GTMDatacenter struct {
	Name             string `json:"name,omitempty"`
	Partition        string `json:"partition,omitempty"`
	FullPath         string `json:"fullPath,omitempty"`
	Contact          string `json:"contact,omitempty"`
	Description      string `json:"description,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	Disabled         bool   `json:"disabled,omitempty"`
	Location         string `json:"location,omitempty"`
	ProberFallback   string `json:"proberFallback,omitempty"`
	ProberPreference string `json:"proberPreference,omitempty"`
}

type Gtmmonitors struct {
	Gtmmonitors []Gtmmonitor `json:"items"`
}

type Gtmmonitor struct {
	Name          string `json:"name,omitempty"`
	Defaults_from string `json:"defaultsFrom,omitempty"`
	Interval      int    `json:"interval,omitempty"`
	Probe_timeout int    `json:"probeTimeout,omitempty"`
	Recv          string `json:"recv,omitempty"`
	Send          string `json:"send,omitempty"`
}

type Servers struct {
	Servers []Server `json:"items"`
}

type Server struct {
	Name                      string
	Datacenter                string
	Description               string
	Monitor                   string
	Virtual_server_discovery  string
	Product                   string
	Enabled                   bool
	Disabled                  bool
	Addresses                 []ServerAddresses
	GTMVirtual_Server         []VSrecord
	ExposeRouteDomains        string
	IqAllowPath               string
	IqAllowServiceCheck       string
	IqAllowSnmp               string
	LinkDiscovery             string
	ProberFallback            string
	ProberPreference          string
	ProberPool                string
	LimitCpuUsage             int
	LimitCpuUsageStatus       string
	LimitMaxBps               int
	LimitMaxBpsStatus         string
	LimitMaxConnections       int
	LimitMaxConnectionsStatus string
	LimitMaxPps               int
	LimitMaxPpsStatus         string
	LimitMemAvail             int
	LimitMemAvailStatus       string
}

type serverDTO struct {
	Name                      string `json:"name"`
	Datacenter                string `json:"datacenter,omitempty"`
	Description               string `json:"description,omitempty"`
	Monitor                   string `json:"monitor,omitempty"`
	Virtual_server_discovery  string `json:"virtualServerDiscovery,omitempty"`
	Product                   string `json:"product,omitempty"`
	Enabled                   bool   `json:"enabled,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	ExposeRouteDomains        string `json:"exposeRouteDomains,omitempty"`
	IqAllowPath               string `json:"iqAllowPath,omitempty"`
	IqAllowServiceCheck       string `json:"iqAllowServiceCheck,omitempty"`
	IqAllowSnmp               string `json:"iqAllowSnmp,omitempty"`
	LinkDiscovery             string `json:"linkDiscovery,omitempty"`
	ProberFallback            string `json:"proberFallback,omitempty"`
	ProberPreference          string `json:"proberPreference,omitempty"`
	ProberPool                string `json:"proberPool,omitempty"`
	LimitCpuUsage             int    `json:"limitCpuUsage,omitempty"`
	LimitCpuUsageStatus       string `json:"limitCpuUsageStatus,omitempty"`
	LimitMaxBps               int    `json:"limitMaxBps,omitempty"`
	LimitMaxBpsStatus         string `json:"limitMaxBpsStatus,omitempty"`
	LimitMaxConnections       int    `json:"limitMaxConnections,omitempty"`
	LimitMaxConnectionsStatus string `json:"limitMaxConnectionsStatus,omitempty"`
	LimitMaxPps               int    `json:"limitMaxPps,omitempty"`
	LimitMaxPpsStatus         string `json:"limitMaxPpsStatus,omitempty"`
	LimitMemAvail             int    `json:"limitMemAvail,omitempty"`
	LimitMemAvailStatus       string `json:"limitMemAvailStatus,omitempty"`
	Addresses                 struct {
		Items []ServerAddresses `json:"items,omitempty"`
	} `json:"addressesReference,omitempty"`
	GTMVirtual_Server struct {
		Items []VSrecord `json:"items,omitempty"`
	} `json:"virtualServersReference,omitempty"`
}

func (p *Server) MarshalJSON() ([]byte, error) {
	return json.Marshal(serverDTO{
		Name:                      p.Name,
		Datacenter:                p.Datacenter,
		Description:               p.Description,
		Monitor:                   p.Monitor,
		Virtual_server_discovery:  p.Virtual_server_discovery,
		Product:                   p.Product,
		Enabled:                   p.Enabled,
		Disabled:                  p.Disabled,
		ExposeRouteDomains:        p.ExposeRouteDomains,
		IqAllowPath:               p.IqAllowPath,
		IqAllowServiceCheck:       p.IqAllowServiceCheck,
		IqAllowSnmp:               p.IqAllowSnmp,
		LinkDiscovery:             p.LinkDiscovery,
		ProberFallback:            p.ProberFallback,
		ProberPreference:          p.ProberPreference,
		ProberPool:                p.ProberPool,
		LimitCpuUsage:             p.LimitCpuUsage,
		LimitCpuUsageStatus:       p.LimitCpuUsageStatus,
		LimitMaxBps:               p.LimitMaxBps,
		LimitMaxBpsStatus:         p.LimitMaxBpsStatus,
		LimitMaxConnections:       p.LimitMaxConnections,
		LimitMaxConnectionsStatus: p.LimitMaxConnectionsStatus,
		LimitMaxPps:               p.LimitMaxPps,
		LimitMaxPpsStatus:         p.LimitMaxPpsStatus,
		LimitMemAvail:             p.LimitMemAvail,
		LimitMemAvailStatus:       p.LimitMemAvailStatus,
		Addresses: struct {
			Items []ServerAddresses `json:"items,omitempty"`
		}{Items: p.Addresses},
		GTMVirtual_Server: struct {
			Items []VSrecord `json:"items,omitempty"`
		}{Items: p.GTMVirtual_Server},
	})
}

func (p *Server) UnmarshalJSON(b []byte) error {
	var dto serverDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	p.Name = dto.Name
	p.Datacenter = dto.Datacenter
	p.Description = dto.Description
	p.Monitor = dto.Monitor
	p.Virtual_server_discovery = dto.Virtual_server_discovery
	p.Product = dto.Product
	p.Enabled = dto.Enabled
	p.Disabled = dto.Disabled
	p.ExposeRouteDomains = dto.ExposeRouteDomains
	p.IqAllowPath = dto.IqAllowPath
	p.IqAllowServiceCheck = dto.IqAllowServiceCheck
	p.IqAllowSnmp = dto.IqAllowSnmp
	p.LinkDiscovery = dto.LinkDiscovery
	p.ProberFallback = dto.ProberFallback
	p.ProberPreference = dto.ProberPreference
	p.ProberPool = dto.ProberPool
	p.LimitCpuUsage = dto.LimitCpuUsage
	p.LimitCpuUsageStatus = dto.LimitCpuUsageStatus
	p.LimitMaxBps = dto.LimitMaxBps
	p.LimitMaxBpsStatus = dto.LimitMaxBpsStatus
	p.LimitMaxConnections = dto.LimitMaxConnections
	p.LimitMaxConnectionsStatus = dto.LimitMaxConnectionsStatus
	p.LimitMaxPps = dto.LimitMaxPps
	p.LimitMaxPpsStatus = dto.LimitMaxPpsStatus
	p.LimitMemAvail = dto.LimitMemAvail
	p.LimitMemAvailStatus = dto.LimitMemAvailStatus
	p.Addresses = dto.Addresses.Items
	p.GTMVirtual_Server = dto.GTMVirtual_Server.Items
	return nil
}

type ServerAddressess struct {
	Items []ServerAddresses `json:"items,omitempty"`
}

type ServerAddresses struct {
	Name        string `json:"name"`
	Device_name string `json:"deviceName,omitempty"`
	Translation string `json:"translation,omitempty"`
}

type VSrecords struct {
	Items []VSrecord `json:"items,omitempty"`
}

type VSrecord struct {
	Name        string `json:"name"`
	Destination string `json:"destination,omitempty"`
}

type GtmPools struct {
	Pools []GtmPool `json:"items"`
}

type GtmPool struct {
	Name                      string `json:"name,omitempty"`
	Partition                 string `json:"partition,omitempty"`
	FullPath                  string `json:"fullPath,omitempty"`
	Generation                int    `json:"generation,omitempty"`
	AlternateMode             string `json:"alternateMode,omitempty"`
	DynamicRatio              string `json:"dynamicRatio,omitempty"`
	Enabled                   bool   `json:"enabled,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	FallbackIp                string `json:"fallbackIp,omitempty"`
	FallbackMode              string `json:"fallbackMode,omitempty"`
	LoadBalancingMode         string `json:"loadBalancingMode,omitempty"`
	ManualResume              string `json:"manualResume,omitempty"`
	MaxAnswersReturned        int    `json:"maxAnswersReturned,omitempty"`
	Monitor                   string `json:"monitor,omitempty"`
	QosHitRatio               int    `json:"qosHitRatio,omitempty"`
	QosHops                   int    `json:"qosHops,omitempty"`
	QosKilobytesSecond        int    `json:"qosKilobytesSecond,omitempty"`
	QosLcs                    int    `json:"qosLcs,omitempty"`
	QosPacketRate             int    `json:"qosPacketRate,omitempty"`
	QosRtt                    int    `json:"qosRtt,omitempty"`
	QosTopology               int    `json:"qosTopology,omitempty"`
	QosVsCapacity             int    `json:"qosVsCapacity,omitempty"`
	QosVsScore                int    `json:"qosVsScore,omitempty"`
	Ttl                       int    `json:"ttl,omitempty"`
	LimitMaxBps               int    `json:"limitMaxBps,omitempty"`
	LimitMaxBpsStatus         string `json:"limitMaxBpsStatus,omitempty"`
	LimitMaxConnections       int    `json:"limitMaxConnections,omitempty"`
	LimitMaxConnectionsStatus string `json:"limitMaxConnectionsStatus,omitempty"`
	LimitMaxPps               int    `json:"limitMaxPps,omitempty"`
	LimitMaxPpsStatus         string `json:"limitMaxPpsStatus,omitempty"`
	MinMembersUpMode          string `json:"minMembersUpMode,omitempty"`
	MinMembersUpValue         int    `json:"minMembersUpValue,omitempty"`
	VerifyMemberAvailability  string `json:"verifyMemberAvailability,omitempty"`
	Members                   []GtmPoolMembers

	// Legacy fields for backward compatibility
	Load_balancing_mode  string `json:"load_balancing_mode,omitempty"`
	Max_answers_returned int    `json:"max_answers_returned,omitempty"`
	Alternate_mode       string `json:"alternate_mode,omitempty"`
	Fallback_ip          string `json:"fallback_ip,omitempty"`
	Fallback_mode        string `json:"fallback_mode,omitempty"`
}

type gtmPoolDTO struct {
	Name                      string `json:"name,omitempty"`
	Partition                 string `json:"partition,omitempty"`
	FullPath                  string `json:"fullPath,omitempty"`
	Generation                int    `json:"generation,omitempty"`
	AlternateMode             string `json:"alternateMode,omitempty"`
	DynamicRatio              string `json:"dynamicRatio,omitempty"`
	Enabled                   bool   `json:"enabled,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	FallbackIp                string `json:"fallbackIp,omitempty"`
	FallbackMode              string `json:"fallbackMode,omitempty"`
	LoadBalancingMode         string `json:"loadBalancingMode,omitempty"`
	ManualResume              string `json:"manualResume,omitempty"`
	MaxAnswersReturned        int    `json:"maxAnswersReturned,omitempty"`
	Monitor                   string `json:"monitor,omitempty"`
	QosHitRatio               int    `json:"qosHitRatio,omitempty"`
	QosHops                   int    `json:"qosHops,omitempty"`
	QosKilobytesSecond        int    `json:"qosKilobytesSecond,omitempty"`
	QosLcs                    int    `json:"qosLcs,omitempty"`
	QosPacketRate             int    `json:"qosPacketRate,omitempty"`
	QosRtt                    int    `json:"qosRtt,omitempty"`
	QosTopology               int    `json:"qosTopology,omitempty"`
	QosVsCapacity             int    `json:"qosVsCapacity,omitempty"`
	QosVsScore                int    `json:"qosVsScore,omitempty"`
	Ttl                       int    `json:"ttl,omitempty"`
	LimitMaxBps               int    `json:"limitMaxBps,omitempty"`
	LimitMaxBpsStatus         string `json:"limitMaxBpsStatus,omitempty"`
	LimitMaxConnections       int    `json:"limitMaxConnections,omitempty"`
	LimitMaxConnectionsStatus string `json:"limitMaxConnectionsStatus,omitempty"`
	LimitMaxPps               int    `json:"limitMaxPps,omitempty"`
	LimitMaxPpsStatus         string `json:"limitMaxPpsStatus,omitempty"`
	MinMembersUpMode          string `json:"minMembersUpMode,omitempty"`
	MinMembersUpValue         int    `json:"minMembersUpValue,omitempty"`
	VerifyMemberAvailability  string `json:"verifyMemberAvailability,omitempty"`
	MembersReference          struct {
		Items []GtmPoolMembers `json:"items,omitempty"`
	} `json:"membersReference,omitempty"`

	// Legacy fields for backward compatibility
	Load_balancing_mode  string `json:"load_balancing_mode,omitempty"`
	Max_answers_returned int    `json:"max_answers_returned,omitempty"`
	Alternate_mode       string `json:"alternate_mode,omitempty"`
	Fallback_ip          string `json:"fallback_ip,omitempty"`
	Fallback_mode        string `json:"fallback_mode,omitempty"`
}

func (p *GtmPool) MarshalJSON() ([]byte, error) {
	return json.Marshal(gtmPoolDTO{
		Name:                      p.Name,
		Partition:                 p.Partition,
		FullPath:                  p.FullPath,
		Generation:                p.Generation,
		AlternateMode:             p.AlternateMode,
		DynamicRatio:              p.DynamicRatio,
		Enabled:                   p.Enabled,
		Disabled:                  p.Disabled,
		FallbackIp:                p.FallbackIp,
		FallbackMode:              p.FallbackMode,
		LoadBalancingMode:         p.LoadBalancingMode,
		ManualResume:              p.ManualResume,
		MaxAnswersReturned:        p.MaxAnswersReturned,
		Monitor:                   p.Monitor,
		QosHitRatio:               p.QosHitRatio,
		QosHops:                   p.QosHops,
		QosKilobytesSecond:        p.QosKilobytesSecond,
		QosLcs:                    p.QosLcs,
		QosPacketRate:             p.QosPacketRate,
		QosRtt:                    p.QosRtt,
		QosTopology:               p.QosTopology,
		QosVsCapacity:             p.QosVsCapacity,
		QosVsScore:                p.QosVsScore,
		Ttl:                       p.Ttl,
		LimitMaxBps:               p.LimitMaxBps,
		LimitMaxBpsStatus:         p.LimitMaxBpsStatus,
		LimitMaxConnections:       p.LimitMaxConnections,
		LimitMaxConnectionsStatus: p.LimitMaxConnectionsStatus,
		LimitMaxPps:               p.LimitMaxPps,
		LimitMaxPpsStatus:         p.LimitMaxPpsStatus,
		MinMembersUpMode:          p.MinMembersUpMode,
		MinMembersUpValue:         p.MinMembersUpValue,
		VerifyMemberAvailability:  p.VerifyMemberAvailability,
		MembersReference: struct {
			Items []GtmPoolMembers `json:"items,omitempty"`
		}{Items: p.Members},
		Load_balancing_mode:  p.Load_balancing_mode,
		Max_answers_returned: p.Max_answers_returned,
		Alternate_mode:       p.Alternate_mode,
		Fallback_ip:          p.Fallback_ip,
		Fallback_mode:        p.Fallback_mode,
	})
}

func (p *GtmPool) UnmarshalJSON(b []byte) error {
	var dto gtmPoolDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	p.Name = dto.Name
	p.Partition = dto.Partition
	p.FullPath = dto.FullPath
	p.Generation = dto.Generation
	p.AlternateMode = dto.AlternateMode
	p.DynamicRatio = dto.DynamicRatio
	p.Enabled = dto.Enabled
	p.Disabled = dto.Disabled
	p.FallbackIp = dto.FallbackIp
	p.FallbackMode = dto.FallbackMode
	p.LoadBalancingMode = dto.LoadBalancingMode
	p.ManualResume = dto.ManualResume
	p.MaxAnswersReturned = dto.MaxAnswersReturned
	p.Monitor = dto.Monitor
	p.QosHitRatio = dto.QosHitRatio
	p.QosHops = dto.QosHops
	p.QosKilobytesSecond = dto.QosKilobytesSecond
	p.QosLcs = dto.QosLcs
	p.QosPacketRate = dto.QosPacketRate
	p.QosRtt = dto.QosRtt
	p.QosTopology = dto.QosTopology
	p.QosVsCapacity = dto.QosVsCapacity
	p.QosVsScore = dto.QosVsScore
	p.Ttl = dto.Ttl
	p.LimitMaxBps = dto.LimitMaxBps
	p.LimitMaxBpsStatus = dto.LimitMaxBpsStatus
	p.LimitMaxConnections = dto.LimitMaxConnections
	p.LimitMaxConnectionsStatus = dto.LimitMaxConnectionsStatus
	p.LimitMaxPps = dto.LimitMaxPps
	p.LimitMaxPpsStatus = dto.LimitMaxPpsStatus
	p.MinMembersUpMode = dto.MinMembersUpMode
	p.MinMembersUpValue = dto.MinMembersUpValue
	p.VerifyMemberAvailability = dto.VerifyMemberAvailability
	p.Members = dto.MembersReference.Items
	p.Load_balancing_mode = dto.Load_balancing_mode
	p.Max_answers_returned = dto.Max_answers_returned
	p.Alternate_mode = dto.Alternate_mode
	p.Fallback_ip = dto.Fallback_ip
	p.Fallback_mode = dto.Fallback_mode
	return nil
}

type GtmPoolMembers struct {
	Name                      string `json:"name,omitempty"`
	Partition                 string `json:"partition,omitempty"`
	SubPath                   string `json:"subPath,omitempty"`
	FullPath                  string `json:"fullPath,omitempty"`
	Generation                int    `json:"generation,omitempty"`
	Enabled                   bool   `json:"enabled,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	LimitMaxBps               int    `json:"limitMaxBps,omitempty"`
	LimitMaxBpsStatus         string `json:"limitMaxBpsStatus,omitempty"`
	LimitMaxConnections       int    `json:"limitMaxConnections,omitempty"`
	LimitMaxConnectionsStatus string `json:"limitMaxConnectionsStatus,omitempty"`
	LimitMaxPps               int    `json:"limitMaxPps,omitempty"`
	LimitMaxPpsStatus         string `json:"limitMaxPpsStatus,omitempty"`
	MemberOrder               int    `json:"memberOrder,omitempty"`
	Monitor                   string `json:"monitor,omitempty"`
	Ratio                     int    `json:"ratio,omitempty"`
}

func (b *BigIP) Gtmmonitors() (*Gtmmonitor, error) {
	var gtmmonitor Gtmmonitor
	err, _ := b.getForEntity(&gtmmonitor, uriGtm, uriGtmmonitor, uriHttp)

	if err != nil {
		return nil, err
	}

	return &gtmmonitor, nil
}
func (b *BigIP) CreateGtmmonitor(name, defaults_from string, interval, probeTimeout int, recv, send string) error {
	config := &Gtmmonitor{
		Name:          name,
		Defaults_from: defaults_from,
		Interval:      interval,
		Probe_timeout: probeTimeout,
		Recv:          recv,
		Send:          send,
	}
	return b.post(config, uriGtm, uriGtmmonitor, uriHttp)
}

func (b *BigIP) ModifyGtmmonitor(*Gtmmonitor) error {
	return b.patch(uriGtm, uriGtmmonitor, uriHttp)
}

func (b *BigIP) DeleteGtmmonitor(name string) error {
	return b.delete(uriGtm, uriGtmmonitor, uriHttp, name)
}

func (b *BigIP) CreateGtmserver(p *Server) error {
	log.Println(" what is the complete payload    ", p)
	return b.post(p, uriGtm, uriServer)
}

// Update an existing policy.
func (b *BigIP) UpdateGtmserver(name string, p *Server) error {
	return b.put(p, uriGtm, uriServer, name)
}

// Delete a policy by name.
func (b *BigIP) DeleteGtmserver(name string) error {
	return b.delete(uriGtm, uriServer, name)
}

func (b *BigIP) GetGtmserver(name string) (*Server, error) {
	var p Server
	err, ok := b.getForEntity(&p, uriGtm, uriServer, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &p, nil
}

// func (b *BigIP) CreatePool_a(name, monitor, load_balancing_mode string, max_answers_returned int, alternate_mode, fallback_ip, fallback_mode string, members []string) error {
// 	config := &Pool{
// 		Name:                 name,
// 		Monitor:              monitor,
// 		Load_balancing_mode:  load_balancing_mode,
// 		Max_answers_returned: max_answers_returned,
// 		Alternate_mode:       alternate_mode,
// 		Fallback_ip:          fallback_ip,
// 		Fallback_mode:        fallback_mode,
// 		Members:              members,
// 	}
// 	log.Println("in poola now", config)
// 	return b.patch(config, uriGtm, uriPoolA)
// }

// func (b *BigIP) ModifyPool_a(config *GtmPool) error {
// 	return b.put(config, uriGtm, uriPoolA)
// }

// func (b *BigIP) Pool_as() (*GtmPool, error) {
// 	var pool GtmPool
// 	err, _ := b.getForEntity(&pool, uriGtm, uriPoolA)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &pool, nil
// }

// AddGTMPool creates a new GTM pool
func (b *BigIP) AddGTMPool(config *GtmPool, poolType string) error {
	return b.post(config, uriGtm, "pool", poolType)
}

// GetGTMPool retrieves a GTM pool by name
func (b *BigIP) GetGTMPool(fullPath string, poolType string) (*GtmPool, error) {
	var pool GtmPool
	err, ok := b.getForEntity(&pool, uriGtm, "pool", poolType, fullPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	// Fetch members separately using the same path pattern as LTM pools
	var membersResponse struct {
		Items []GtmPoolMembers `json:"items"`
	}
	log.Printf("[DEBUG] Attempting to fetch members for GTM pool: %s, type: %s, path: gtm/pool/%s/%s/members", fullPath, poolType, poolType, fullPath)
	err, ok = b.getForEntity(&membersResponse, uriGtm, "pool", poolType, fullPath, "members")
	if err != nil {
		log.Printf("[DEBUG] Error fetching GTM pool members for %s: %v", fullPath, err)
		// Don't fail - pool might not have members or endpoint might not be accessible
	} else if ok {
		log.Printf("[DEBUG] Members response OK, found %d items", len(membersResponse.Items))
		if len(membersResponse.Items) > 0 {
			pool.Members = membersResponse.Items
			log.Printf("[DEBUG] Fetched %d members for GTM pool %s", len(membersResponse.Items), fullPath)
		} else {
			log.Printf("[DEBUG] No members found for GTM pool %s", fullPath)
		}
	} else {
		log.Printf("[DEBUG] Members fetch returned ok=false for pool %s", fullPath)
	}

	return &pool, nil
}

// ModifyGTMPool updates a GTM pool
func (b *BigIP) ModifyGTMPool(fullPath string, config *GtmPool, poolType string) error {
	return b.put(config, uriGtm, "pool", poolType, fullPath)
}

// DeleteGTMPool removes a GTM pool
func (b *BigIP) DeleteGTMPool(fullPath string, poolType string) error {
	return b.delete(uriGtm, "pool", poolType, fullPath)
}

func (b *BigIP) CreateGTMWideIP(config *GTMWideIP, recordType string) error {
	return b.post(config, uriGtm, uriWideIp, recordType)
}

func (b *BigIP) DeleteGTMWideIP(fullPath string, recordType string) error {
	return b.delete(uriGtm, uriWideIp, recordType, fullPath)
}
func (b *BigIP) ModifyGTMWideIP(fullPath string, config *GTMWideIP, recordType string) error {
	return b.put(config, uriGtm, uriWideIp, recordType, fullPath)
}

// GetGTMWideIP get's a WideIP by name
func (b *BigIP) GetGTMWideIP(name string, recordType string) (*GTMWideIP, error) {
	var w GTMWideIP

	err, ok := b.getForEntity(&w, uriGtm, uriWideIp, string(recordType), name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &w, nil
}

// CreateGTMDatacenter creates a new GTM datacenter
func (b *BigIP) CreateGTMDatacenter(config *GTMDatacenter) error {
	return b.post(config, uriGtm, uriDatacenter)
}

// GetGTMDatacenter retrieves a GTM datacenter by full path
func (b *BigIP) GetGTMDatacenter(fullPath string) (*GTMDatacenter, error) {
	var dc GTMDatacenter
	err, ok := b.getForEntity(&dc, uriGtm, uriDatacenter, fullPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &dc, nil
}

// ModifyGTMDatacenter updates a GTM datacenter
func (b *BigIP) ModifyGTMDatacenter(fullPath string, config *GTMDatacenter) error {
	return b.put(config, uriGtm, uriDatacenter, fullPath)
}

// DeleteGTMDatacenter removes a GTM datacenter
func (b *BigIP) DeleteGTMDatacenter(fullPath string) error {
	return b.delete(uriGtm, uriDatacenter, fullPath)
}

// func (b *BigIP) GetDatacenters() (*GTMDatacenter, error) {
// 	var datacenter GTMDatacenter
// 	err, _ := b.getForEntity(&datacenter, uriGtm, uriDatacenter)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &datacenter, nil
// }

// func (b *BigIP) CreateDatacenter(name, description, contact, app_service string, enabled, disabled bool, prober_pool string) error {
// 	config := &GTMDatacenter{
// 		Name:        name,
// 		Description: description,
// 		Contact:     contact,
// 		App_service: app_service,
// 		Enabled:     enabled,
// 		Disabled:    disabled,
// 		Prober_pool: prober_pool,
// 	}
// 	return b.post(config, uriGtm, uriDatacenter)
// }

// func (b *BigIP) ModifyDatacenter(*GTMDatacenter) error {
// 	return b.patch(uriGtm, uriDatacenter)
// }

// func (b *BigIP) DeleteDatacenter(name string) error {
// 	return b.delete(uriGtm, uriDatacenter, name)
// }

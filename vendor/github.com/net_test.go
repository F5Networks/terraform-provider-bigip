package bigip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
)

type NetTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter, *http.Request)
}

func (s *NetTestSuite) SetupSuite() {
	s.Server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		s.LastRequestBody = string(body)
		s.LastRequest = r
		if s.ResponseFunc != nil {
			s.ResponseFunc(w, r)
		}
	}))

	s.Client = NewSession(s.Server.URL, "", "", nil)
}

func (s *NetTestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *NetTestSuite) SetupTest() {
	s.ResponseFunc = nil
	s.LastRequest = nil
}

func TestNetSuite(t *testing.T) {
	suite.Run(t, new(NetTestSuite))
}

func (s *NetTestSuite) TestGetInterfaces() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:interface:interfacecollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/interface?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:interface:interfacestate",
      "name": "1.1",
      "fullPath": "1.1",
      "generation": 53,
      "selfLink": "https://localhost/mgmt/tm/net/interface/1.1?ver=11.5.1",
      "bundle": "not-supported",
      "enabled": true,
      "forceGigabitFiber": "disabled",
      "ifIndex": 48,
      "lldpAdmin": "txonly",
      "lldpTlvmap": 130943,
      "macAddress": "00:00:00:00:00:00",
      "mediaFixed": "10000T-FD",
      "mediaMax": "10000T-FD",
      "mediaSfp": "auto",
      "mtu": 9198,
      "preferPort": "sfp",
      "sflow": {
        "pollInterval": 0,
        "pollIntervalGlobal": "yes"
      },
      "stp": "enabled",
      "stpAutoEdgePort": "enabled",
      "stpEdgePort": "true",
      "stpLinkType": "auto"
    }]}`))
	}

	i, err := s.Client.Interfaces()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/interface", "")
	assert.Equal(s.T(), 1, len(i.Interfaces))
	assert.Equal(s.T(), "1.1", i.Interfaces[0].Name)
}

func (s *NetTestSuite) TestAddInterfaceToVLan() {
	err := s.Client.AddInterfaceToVlan("vlan-name", "iface-name", false)

	assert.Nil(s.T(), err)
	assertRestCall(s,
		"PUT",
		"/mgmt/tm/net/vlan/vlan-name/interfaces",
		`{"name":"iface-name", "untagged":true}`)
}

func (s *NetTestSuite) TestSelfIPs() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:self:selfcollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/self?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:self:selfstate",
      "name": "0.0.0.0",
      "partition": "Common",
      "fullPath": "/Common/0.0.0.0",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/self/~Common~0.0.0.0?ver=11.5.1",
      "address": "0.0.0.0/20",
      "floating": "disabled",
      "inheritedTrafficGroup": "false",
      "trafficGroup": "/Common/traffic-group-local-only",
      "unit": 0,
      "vlan": "/Common/vlan-name",
      "allowService": "all"
    }]}`))
	}

	ips, err := s.Client.SelfIPs()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/self", "")
	assert.Equal(s.T(), 1, len(ips.SelfIPs))
	assert.Equal(s.T(), "0.0.0.0", ips.SelfIPs[0].Name)

}

func (s *NetTestSuite) TestCreateSelfIP() {
	err := s.Client.CreateSelfIP("0.0.0.0", "0.0.0.0/20", "vlan")

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/self", `{"name":"0.0.0.0","address":"0.0.0.0/20", "vlan":"vlan"}`)
}

func (s *NetTestSuite) TestDeleteSelfIP() {
	err := s.Client.DeleteSelfIP("0.0.0.0")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/self/0.0.0.0", "")
}

func (s *NetTestSuite) TestModifySelfIP() {
	ip := &SelfIP{Address: "0.0.0.0/24"}

	err := s.Client.ModifySelfIP("0.0.0.0", ip)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/self/0.0.0.0", "")
}

func (s *NetTestSuite) TestTrunks() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:self:selfcollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/trunk?ver=11.5.1",
  "items": [
  {
  	"name" : "trunk-name",
  	"id" : 1
  }
  ]}`))
	}

	trunks, err := s.Client.Trunks()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/trunk", "")
	assert.Equal(s.T(), 1, len(trunks.Trunks))
}

func (s *NetTestSuite) TestCreateTrunk() {
	err := s.Client.CreateTrunk("name", "ifaces", true)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/trunk", `{"name":"name", "interfaces":["ifaces"], "lacp":"enabled"}`)
}

func (s *NetTestSuite) TestDeleteTrunk() {
	err := s.Client.DeleteTrunk("name")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/trunk/name", "")
}

func (s *NetTestSuite) TestModifyTrunk() {
	trunk := &Trunk{Name: "name", LACP: "enabled"}

	err := s.Client.ModifyTrunk("name", trunk)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/trunk/name", `{"name":"name", "lacp":"enabled"}`)
}

func (s *NetTestSuite) TestVlans() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:vlan:vlancollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/vlan?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:vlan:vlanstate",
      "name": "vlan",
      "partition": "Common",
      "fullPath": "/Common/vlan",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/vlan/~Common~vlan?ver=11.5.1",
      "autoLasthop": "default",
      "cmpHash": "default",
      "dagRoundRobin": "disabled",
      "failsafe": "disabled",
      "failsafeAction": "failover-restart-tm",
      "failsafeTimeout": 90,
      "ifIndex": 80,
      "learning": "enable-forward",
      "mtu": 1500,
      "sflow": {
        "pollInterval": 0,
        "pollIntervalGlobal": "yes",
        "samplingRate": 0,
        "samplingRateGlobal": "yes"
      },
      "sourceChecking": "disabled",
      "tag": 10,
      "interfacesReference": {
        "link": "https://localhost/mgmt/tm/net/vlan/~Common~vlan/interfaces?ver=11.5.1",
        "isSubcollection": true
      }
    }]}`))
	}

	vlans, err := s.Client.Vlans()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/vlan", "")
	assert.Equal(s.T(), 1, len(vlans.Vlans))
	assert.Equal(s.T(), "vlan", vlans.Vlans[0].Name)
	assert.Equal(s.T(), 0, vlans.Vlans[0].SFlow.PollInterval)
}

func (s *NetTestSuite) TestCreateVLan() {
	err := s.Client.CreateVlan("name", 1)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/vlan", `{"name":"name", "tag":1, "sflow":{}}`)
}

func (s *NetTestSuite) TestDeleteVLan() {
	err := s.Client.DeleteVlan("name")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/vlan/name", "")
}

func (s *NetTestSuite) TestModifyVLan() {
	vlan := &Vlan{MTU: 1500}

	err := s.Client.ModifyVlan("name", vlan)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/vlan/name", `{"mtu":1500, "sflow":{}}`)
}

func (s *NetTestSuite) TestRoutes() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:route:routecollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/route?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:route:routestate",
      "name": "default_route",
      "partition": "Common",
      "fullPath": "/Common/default_route",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/route/~Common~default_route?ver=11.5.1",
      "gw": "0.0.0.0",
      "mtu": 0,
      "network": "default"
    }
  ]
}`))
	}

	routes, err := s.Client.Routes()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/route", "")
	assert.Equal(s.T(), 1, len(routes.Routes))
	assert.Equal(s.T(), "default_route", routes.Routes[0].Name)
}

func (s *NetTestSuite) TestCreateRoute() {
	err := s.Client.CreateRoute("default_route", "default", "0.0.0.0")

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route", `{"name":"default_route", "network":"default", "gw":"0.0.0.0"}`)
}

func (s *NetTestSuite) TestDeleteRoute() {
	err := s.Client.DeleteRoute("default_route")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/route/default_route", "")
}

func (s *NetTestSuite) TestModifyRoute() {
	route := &Route{Gateway: "1.1.1.1"}

	err := s.Client.ModifyRoute("default_route", route)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/route/default_route", `{"gw":"1.1.1.1"}`)
}

func (s *NetTestSuite) TestRouteDomains() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:route-domain:route-domaincollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/route-domain?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:route-domain:route-domainstate",
      "name": "0",
      "partition": "Common",
      "fullPath": "/Common/0",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/route-domain/~Common~0?ver=11.5.1",
      "id": 0,
      "strict": "enabled",
      "vlans": [
        "/Common/http-tunnel",
        "/Common/socks-tunnel"
      ]
    }
  ]
}`))
	}

	routes, err := s.Client.RouteDomains()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/route-domain", "")
	assert.Equal(s.T(), 1, len(routes.RouteDomains))
	assert.Equal(s.T(), 2, len(routes.RouteDomains[0].Vlans))
	assert.Equal(s.T(), "0", routes.RouteDomains[0].Name)
}

func (s *NetTestSuite) TestCreateRouteDomain() {
	err := s.Client.CreateRouteDomain("name", 1, false, "vlan1,vlan2")

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route-domain", `{"name":"name", "id":1, "vlans":["vlan1","vlan2"], "strict":"disabled"}`)
}

func (s *NetTestSuite) TestDeleteRouteDomain() {
	err := s.Client.DeleteRouteDomain("name")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/route-domain/name", "")
}

func (s *NetTestSuite) TestModifyRouteDomain() {
	route := &RouteDomain{Name: "name", ID: 1, Strict: "enabled"}

	err := s.Client.ModifyRouteDomain("name", route)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/route-domain/name", `{"name":"name", "id": 1, "strict" : "enabled"}`)
}

func assertRestCall(s *NetTestSuite, method, path, body string) {
	assert.Equal(s.T(), method, s.LastRequest.Method)
	assert.Equal(s.T(), path, s.LastRequest.URL.Path)
	if body != "" {
		assert.JSONEq(s.T(), body, s.LastRequestBody)
	}
}

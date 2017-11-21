package bigip

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LTMTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter, *http.Request)
}

func (s *LTMTestSuite) SetupSuite() {
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

func (s *LTMTestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *LTMTestSuite) SetupTest() {
	s.ResponseFunc = nil
	s.LastRequest = nil
}

func TestLtmSuite(t *testing.T) {
	suite.Run(t, new(LTMTestSuite))
}

func (s *LTMTestSuite) TestIRules() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"items" :
			[
				{"name":"rule1","apiAnonymous":"rule1"},
				{"name":"rule2","apiAnonymous":"this\nis\nrule2"}
			]}`))
	}

	rules, err := s.Client.IRules()

	assert.Nil(s.T(), err, "Error loading rules")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriIRule), s.LastRequest.URL.Path, "Wrong uri to fetch rules")
	assert.Equal(s.T(), 2, len(rules.IRules), "Wrong number of rules")
	assert.Equal(s.T(), "rule1", rules.IRules[0].Name)
	assert.Equal(s.T(), "rule2", rules.IRules[1].Name)
	assert.Equal(s.T(), `this
is
rule2`, rules.IRules[1].Rule, "Multiline rule not unmarshalled")
}

func (s *LTMTestSuite) TestCreateIRule() {
	s.Client.CreateIRule("rule1", `when CLIENT_ACCEPTED { log local0. "test"}`)

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriIRule), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"rule1","apiAnonymous":"when CLIENT_ACCEPTED { log local0. \"test\"}"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestModifyIRule() {
	s.Client.ModifyIRule("rule1", &IRule{Rule: "modified"})

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriIRule, "rule1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"rule1","apiAnonymous":"modified"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDeleteIRule() {
	s.Client.DeleteIRule("rule1")

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriIRule, "rule1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}

func (s *LTMTestSuite) TestModifyVirtualAddress() {
	d := &VirtualAddress{}
	s.Client.ModifyVirtualAddress("address1", d)

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriVirtualAddress, "address1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
}

func (s *LTMTestSuite) TestGetPolicies() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:ltm:policy:policycollectionstate",
  "selfLink": "https://localhost/mgmt/tm/ltm/policy?ver=11.5.1",
  "items": [
    {
      "kind": "tm:ltm:policy:policystate",
      "name": "policy1",
      "partition": "Common",
      "fullPath": "/Common/policy1",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/ltm/policy/~Common~policy1?ver=11.5.1",
      "controls": [
        "forwarding"
      ],
      "requires": [
        "http",
        "client-ssl"
      ],
      "strategy": "/Common/first-match",
      "rulesReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy/~Common~policy1/rules?ver=11.5.1",
        "isSubcollection": true
      }
    },
    {
      "kind": "tm:ltm:policy:policystate",
      "name": "policy2",
      "partition": "Common",
      "fullPath": "/Common/policy2",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/ltm/policy/~Common~policy2?ver=11.5.1",
      "controls": [
        "classification"
      ],
      "requires": [
        "ssl-persistence"
      ],
      "strategy": "/Common/first-match",
      "rulesReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy/~Common~policy2/rules?ver=11.5.1",
        "isSubcollection": true
      }
    }
    ]}`))
	}

	p, e := s.Client.Policies()

	assert.Nil(s.T(), e, "Fetching policy list should not return an error")
	assert.Equal(s.T(), 2, len(p.Policies), "Wrong number of policies returned")
	assert.Equal(s.T(), "policy1", p.Policies[0].Name)
	assert.Equal(s.T(), "Common", p.Policies[0].Partition)
	assert.Equal(s.T(), "/Common/first-match", p.Policies[0].Strategy)
	assert.EqualValues(s.T(), []string{"http", "client-ssl"}, p.Policies[0].Requires)
	assert.EqualValues(s.T(), []string{"forwarding"}, p.Policies[0].Controls)
}

func (s *LTMTestSuite) TestGetPolicy() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "rules") {
			w.Write([]byte(`{
			  "kind": "tm:ltm:policy:rules:rulescollectionstate",
			  "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules?ver=11.5.1",
			  "items": [
			    {
			      "kind": "tm:ltm:policy:rules:rulesstate",
			      "name": "rule1",
			      "fullPath": "rule1",
			      "generation": 144344,
			      "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1?ver=11.5.1",
			      "ordinal": 0,
			      "actionsReference": {
				"link": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/actions?ver=11.5.1",
				"isSubcollection": true
			      },
			      "conditionsReference": {
				"link": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/conditions?ver=11.5.1",
				"isSubcollection": true
			      }
			    }
			  ]
			}`))
		} else if strings.HasSuffix(r.URL.Path, "actions") {
			w.Write([]byte(`{
			  "kind": "tm:ltm:policy:rules:actions:actionscollectionstate",
			  "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/actions?ver=11.5.1",
			  "items": [
			    {
			      "kind": "tm:ltm:policy:rules:actions:actionsstate",
			      "name": "0",
			      "fullPath": "0",
			      "generation": 144344,
			      "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/actions/0?ver=11.5.1",
			      "code": 0,
			      "forward": true,
			      "pool": "/Common/sorry_server",
			      "port": 0,
			      "request": true,
			      "select": true,
			      "status": 0,
			      "vlanId": 0
			    }
			  ]
			}`))
		} else if strings.HasSuffix(r.URL.Path, "conditions") {
			w.Write([]byte(`{
			  "kind": "tm:ltm:policy:rules:conditions:conditionscollectionstate",
			  "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/conditions?ver=11.5.1",
			  "items": [
			    {
			      "kind": "tm:ltm:policy:rules:conditions:conditionsstate",
			      "name": "0",
			      "fullPath": "0",
			      "generation": 144344,
			      "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/conditions/0?ver=11.5.1",
			      "all": true,
			      "caseInsensitive": true,
			      "external": true,
			      "httpUri": true,
			      "index": 0,
			      "present": true,
			      "remote": true,
			      "request": true,
			      "startsWith": true,
			      "values": [
				"/foo"
			      ]
			    }]
			    }`))
		} else {
			w.Write([]byte(`{
			  "kind": "tm:ltm:policy:policystate",
			  "name": "my_policy",
			  "fullPath": "my_policy",
			  "generation": 144344,
			  "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy?ver=11.5.1",
			  "controls": [
			    "forwarding"
			  ],
			  "requires": [
			    "http"
			  ],
			  "strategy": "/Common/first-match",
			  "rulesReference": {
			    "link": "https://localhost/mgmt/tm/ltm/policy/~Common~my_policy/rules?ver=11.5.1",
			    "isSubcollection": true
			  }
			}`))
		}
	}

	p, err := s.Client.GetPolicy("my_policy")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "my_policy", p.Name)
	assert.Equal(s.T(), 1, len(p.Rules), "Not enough rules")
	assert.Equal(s.T(), 1, len(p.Rules[0].Actions), "Not enough actions")
	assert.Equal(s.T(), "/Common/sorry_server", p.Rules[0].Actions[0].Pool)
	assert.Equal(s.T(), 1, len(p.Rules[0].Conditions), "Not enough conditions")
	assert.Equal(s.T(), []string{"/foo"}, p.Rules[0].Conditions[0].Values)
}

func (s *LTMTestSuite) TestGetNonExistentPolicy() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}

	p, err := s.Client.GetPolicy("asdf")

	assert.NotNil(s.T(), err, "nil error returned")
	assert.Nil(s.T(), p)
	assert.True(s.T(), strings.HasPrefix(err.Error(), "HTTP 404"), err.Error())
}

func (s *LTMTestSuite) TestCreatePolicy() {
	p := Policy{
		Name:     "test",
		Controls: []string{"forwarding"},
		Requires: []string{"http"},
		Strategy: "/Common/first-match",
		Rules: []PolicyRule{
			PolicyRule{
				Name: "rule1",
				Actions: []PolicyRuleAction{
					PolicyRuleAction{
						Forward: true,
						Pool:    "somepool",
					},
					PolicyRuleAction{},
				},
				Conditions: []PolicyRuleCondition{
					PolicyRuleCondition{
						CaseInsensitive: true,
						Values:          []string{"/foo", "/bar"},
					},
					PolicyRuleCondition{},
				},
			},
			PolicyRule{
				Name: "rule2",
			},
		},
	}

	s.Client.CreatePolicy(&p)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriPolicy), s.LastRequest.URL.Path)
	assert.JSONEq(s.T(), `{"name":"test",
		"controls":["forwarding"],
		"requires":["http"],
		"strategy":"/Common/first-match",
		"rulesReference" : {"items":[
			{ "name":"rule1",
			  "ordinal":0,
			  "actionsReference":{ "items":[
			  	{
			  		"name":"0",
			  		"pool":"somepool",
			  		"forward":true
			  	},
			  	{
			  		"name":"1"
			  	}
			  ]},
			  "conditionsReference" : {"items":[
			  	{
			  		"name":"0",
			  		"caseInsensitive":true,
			  		"values":["/foo","/bar"]
			  	},
			  	{
			  		"name":"1"
			  	}
			  ]}
			},
			{ "name":"rule2", "ordinal":1, "actionsReference":{}, "conditionsReference":{}}]}}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestUpdatePolicy() {
	//TODO: test more stuff
	s.Client.UpdatePolicy("foo", &Policy{})

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/foo", uriLtm, uriPolicy), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestDeletePolicy() {
	s.Client.DeletePolicy("foo")

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/foo", uriLtm, uriPolicy), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestCreateVitualAddress() {

	s.Client.CreateVirtualAddress("test-va", &VirtualAddress{Address: "10.10.10.10", ARP: true, AutoDelete: false})

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriVirtualAddress), s.LastRequest.URL.Path)
	assert.JSONEq(s.T(), `
	{"name":"test-va",
	"arp":"enabled",
	"autoDelete":"false",
	"address" : "10.10.10.10",
	"enabled":"no",
	"floating":"disabled",
	"icmpEcho":"disabled",
	"inheritedTrafficGroup":"no",
	"routeAdvertisement":"disabled"}`, s.LastRequestBody)

}

func (s *LTMTestSuite) TestDeleteVitualAddress() {

	s.Client.DeleteVirtualAddress("test-va")

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/test-va", uriLtm, uriVirtualAddress), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestCreateVirtualServer() {
	s.Client.CreateVirtualServer("/Common/test-vs", "10.10.10.10", "255.255.255.255", "/Common/test-pool", 80)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriVirtual), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"/Common/test-vs","destination":"10.10.10.10:80","mask":"255.255.255.255","pool":"/Common/test-pool","sourceAddressTranslation":{}}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestAddVirtualServer() {
	config := &VirtualServer{
		Name:        "/Common/test-vs",
		Destination: "10.10.10.10:80",
		Pool:        "/Common/test-pool",
	}

	s.Client.AddVirtualServer(config)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriVirtual), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"/Common/test-vs","destination":"10.10.10.10:80","pool":"/Common/test-pool","sourceAddressTranslation":{}}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestModifyVirtualServer() {
	vs := &VirtualServer{
		Name: "test",
		Profiles: []Profile{
			Profile{Name: "/Common/tcp", Context: CONTEXT_CLIENT},
			Profile{Name: "/Common/tcp", Context: CONTEXT_SERVER}},
		//TODO: test more
	}

	s.Client.ModifyVirtualServer("test", vs)

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/test", uriLtm, uriVirtual), s.LastRequest.URL.Path)
	assert.JSONEq(s.T(), `
	{"name":"test",
	"sourceAddressTranslation":{},
	"profiles":[
		{"name":"/Common/tcp","context":"clientside"},
		{"name":"/Common/tcp","context":"serverside"}
	]
	}`, s.LastRequestBody)

}

func (s *LTMTestSuite) TestDeleteVirtualServer() {
	s.Client.DeleteVirtualServer("/Common/test-vs")

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriVirtual, "~Common~test-vs"), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestCreatePool() {
	name := "/Common/test-pool"

	s.Client.CreatePool(name)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriPool), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"/Common/test-pool"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestAddPool() {
	config := &Pool{
		Name:      "test-pool",
		Partition: "Common",
		Monitor:   "/Common/http",
	}

	s.Client.AddPool(config)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriPool), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-pool","partition":"Common","monitor":"/Common/http"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestGetPool() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "kind": "tm:ltm:pool:poolstate",
    "name": "test-pool",
    "partition": "Common",
    "fullPath": "/Common/test-pool",
    "generation": 3882,
    "selfLink": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool?ver=11.5.3",
    "allowNat": "yes",
    "allowSnat": "yes",
    "ignorePersistedWeight": "disabled",
    "ipTosToClient": "pass-through",
    "ipTosToServer": "pass-through",
    "linkQosToClient": "pass-through",
    "linkQosToServer": "pass-through",
    "loadBalancingMode": "round-robin",
    "minActiveMembers": 0,
    "minUpMembers": 0,
    "minUpMembersAction": "failover",
    "minUpMembersChecking": "disabled",
    "monitor": "/Common/http ",
    "queueDepthLimit": 0,
    "queueOnConnectionLimit": "disabled",
    "queueTimeLimit": 0,
    "reselectTries": 0,
    "slowRampTime": 10,
    "membersReference": {
        "link": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool/members?ver=11.5.3",
        "isSubcollection": true
				}
		}`))
	}

	p, err := s.Client.GetPool("/Common/test-pool")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "test-pool", p.Name)
	assert.Equal(s.T(), "Common", p.Partition)
	assert.Equal(s.T(), "/Common/test-pool", p.FullPath)
	assert.Equal(s.T(), "/Common/http ", p.Monitor)
}

func (s *LTMTestSuite) TestPools() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
				"kind": "tm:ltm:pool:poolcollectionstate",
				"selfLink": "https://localhost/mgmt/tm/ltm/pool?ver=11.5.3",
				"items": [
						{
								"kind": "tm:ltm:pool:poolstate",
								"name": "test-pool",
								"partition": "Common",
								"fullPath": "/Common/test-pool",
								"generation": 3882,
								"selfLink": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool?ver=11.5.3",
								"allowNat": "yes",
								"allowSnat": "yes",
								"ignorePersistedWeight": "disabled",
								"ipTosToClient": "pass-through",
								"ipTosToServer": "pass-through",
								"linkQosToClient": "pass-through",
								"linkQosToServer": "pass-through",
								"loadBalancingMode": "round-robin",
								"minActiveMembers": 0,
								"minUpMembers": 0,
								"minUpMembersAction": "failover",
								"minUpMembersChecking": "disabled",
								"monitor": "/Common/http ",
								"queueDepthLimit": 0,
								"queueOnConnectionLimit": "disabled",
								"queueTimeLimit": 0,
								"reselectTries": 0,
								"slowRampTime": 10,
								"membersReference": {
										"link": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool/members?ver=11.5.3",
										"isSubcollection": true
								}
						},
						{
								"kind": "tm:ltm:pool:poolstate",
								"name": "test-pool2",
								"partition": "Common",
								"fullPath": "/Common/test-pool2",
								"generation": 3886,
								"selfLink": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool2?ver=11.5.3",
								"allowNat": "no",
								"allowSnat": "no",
								"ignorePersistedWeight": "disabled",
								"ipTosToClient": "pass-through",
								"ipTosToServer": "pass-through",
								"linkQosToClient": "pass-through",
								"linkQosToServer": "pass-through",
								"loadBalancingMode": "round-robin",
								"minActiveMembers": 0,
								"minUpMembers": 0,
								"minUpMembersAction": "failover",
								"minUpMembersChecking": "disabled",
								"queueDepthLimit": 0,
								"queueOnConnectionLimit": "disabled",
								"queueTimeLimit": 0,
								"reselectTries": 0,
								"slowRampTime": 10,
								"membersReference": {
										"link": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool2/members?ver=11.5.3",
										"isSubcollection": true
								}
						}
				]
		}`))
	}

	p, err := s.Client.Pools()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriPool), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "/Common/test-pool", p.Pools[0].FullPath)
	assert.Equal(s.T(), "/Common/test-pool2", p.Pools[1].FullPath)
}

func (s *LTMTestSuite) TestDeletePool() {
	name := "/Common/test-pool"
	s.Client.DeletePool(name)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool"), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestModifyPool() {
	config := &Pool{
		Name:              "test-pool",
		Partition:         "Common",
		Monitor:           "/Common/http",
		LoadBalancingMode: "round-robin",
		AllowSNAT:         "yes",
		AllowNAT:          "yes",
	}

	s.Client.ModifyPool("/Common/test-pool", config)

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-pool","partition":"Common","allowNat":"yes","allowSnat":"yes","loadBalancingMode":"round-robin","monitor":"/Common/http"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestAddPoolMember() {
	pool := "/Common/test-pool"
	poolmember := "/Common/test-pool-member"

	s.Client.AddPoolMember(pool, poolmember)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool", uriPoolMember), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestCreatePoolMember() {
	pool := "/Common/test-pool"
	config := &PoolMember{
		Name:      "test-pool-member",
		Partition: "Common",
		Monitor:   "/Common/icmp",
	}

	s.Client.CreatePoolMember(pool, config)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool", uriPoolMember), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-pool-member","partition":"Common","monitor":"/Common/icmp"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestGetPoolMember() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "kind": "tm:ltm:pool:members:membersstate",
    "name": "test-pool-member:80",
    "partition": "Common",
    "fullPath": "/Common/test-pool-member:80",
    "generation": 5124,
    "selfLink": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~test-pool-member:80?ver=11.5.3",
    "address": "10.10.20.30",
    "connectionLimit": 0,
    "dynamicRatio": 1,
    "inheritProfile": "disabled",
    "logging": "disabled",
    "monitor": "default",
    "priorityGroup": 0,
    "rateLimit": "disabled",
    "ratio": 1,
    "session": "monitor-enabled",
    "state": "down"
		}`))
	}

	pool := "/Common/test-pool"
	poolmember := "/Common/test-pool-member:80"
	p, err := s.Client.GetPoolMember(pool, poolmember)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool", uriPoolMember, "~Common~test-pool-member:80"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "test-pool-member:80", p.Name)
	assert.Equal(s.T(), "Common", p.Partition)
	assert.Equal(s.T(), "/Common/test-pool-member:80", p.FullPath)
	assert.Equal(s.T(), "default", p.Monitor)
}

func (s *LTMTestSuite) TestPoolMembers() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "kind": "tm:ltm:pool:members:memberscollectionstate",
    "selfLink": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool/members?ver=11.5.3",
    "items": [
        {
            "kind": "tm:ltm:pool:members:membersstate",
            "name": "test-pool-member-1:80",
            "partition": "Common",
            "fullPath": "/Common/test-pool-member-1:80",
            "generation": 5124,
            "selfLink": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~test-pool-member-1:80?ver=11.5.3",
            "address": "10.156.153.10",
            "connectionLimit": 0,
            "dynamicRatio": 1,
            "inheritProfile": "enabled",
            "logging": "disabled",
            "monitor": "default",
            "priorityGroup": 0,
            "rateLimit": "disabled",
            "ratio": 1,
            "session": "monitor-enabled",
            "state": "down"
        },
        {
            "kind": "tm:ltm:pool:members:membersstate",
            "name": "test-pool-member-2:80",
            "partition": "Common",
            "fullPath": "/Common/test-pool-member-2:80",
            "generation": 3882,
            "selfLink": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~test-pool-member-2:80?ver=11.5.3",
            "address": "10.10.20.30",
            "connectionLimit": 0,
            "dynamicRatio": 1,
            "inheritProfile": "disabled",
            "logging": "disabled",
            "monitor": "default",
            "priorityGroup": 0,
            "rateLimit": "disabled",
            "ratio": 1,
            "session": "monitor-enabled",
            "state": "down"
        },
        {
            "kind": "tm:ltm:pool:members:membersstate",
            "name": "test-pool-member-3:80",
            "partition": "Common",
            "fullPath": "/Common/test-pool-member-3:80",
            "generation": 3862,
            "selfLink": "https://localhost/mgmt/tm/ltm/pool/~Common~test-pool/members/~Common~test-pool-member-3:80?ver=11.5.3",
            "address": "10.10.20.40",
            "connectionLimit": 0,
            "dynamicRatio": 1,
            "inheritProfile": "enabled",
            "logging": "disabled",
            "monitor": "/Common/http ",
            "priorityGroup": 0,
            "rateLimit": "disabled",
            "ratio": 1,
            "session": "monitor-enabled",
            "state": "down"
        }
    ]
		}`))
	}

	pool := "/Common/test-pool"
	p, err := s.Client.PoolMembers(pool)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool", uriPoolMember), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "/Common/test-pool-member-1:80", p.PoolMembers[0].FullPath)
	assert.Equal(s.T(), "/Common/test-pool-member-2:80", p.PoolMembers[1].FullPath)
	assert.Equal(s.T(), "/Common/test-pool-member-3:80", p.PoolMembers[2].FullPath)
}

func (s *LTMTestSuite) TestDeletePoolMember() {

	pool := "/Common/test-pool"
	poolmember := "/Common/test-pool-member:80"
	s.Client.DeletePoolMember(pool, poolmember)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool", uriPoolMember, "~Common~test-pool-member:80"), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestRemovePoolMember() {
	pool := "/Common/test-pool"
	config := &PoolMember{
		Name:      "test-pool-member",
		Partition: "Common",
		Monitor:   "/Common/icmp",
	}
	s.Client.RemovePoolMember(pool, config)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s/", uriLtm, uriPool, "~Common~test-pool", uriPoolMember), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestModifyPoolMember() {
	pool := "/Common/test-pool"
	config := &PoolMember{
		Name:      "test-pool-member:80",
		Partition: "Common",
		FullPath:  "/Common/test-pool-member:80",
		Monitor:   "/Common/icmp",
	}

	s.Client.ModifyPoolMember(pool, config)

	fmt.Println(s.LastRequest.URL)
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s/%s", uriLtm, uriPool, "~Common~test-pool", uriPoolMember, "~Common~test-pool-member:80"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"monitor":"/Common/icmp"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestCreateMonitor() {
	config := &Monitor{
		Name:          "test-web-monitor",
		ParentMonitor: "http",
		Interval:      15,
		Timeout:       5,
		SendString:    "GET /\r\n",
		ReceiveString: "200 OK",
	}

	s.Client.CreateMonitor(config.Name, config.ParentMonitor, config.Interval, config.Timeout, config.SendString, config.ReceiveString)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriMonitor, config.ParentMonitor), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-web-monitor","defaultsFrom":"http","interval":15,"manualResume":"disabled","recv":"200 OK","reverse":"disabled","send":"GET /\\r\\n","timeout":5,"transparent":"disabled"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestCreateMonitorSpecialCharacters() {
	config := &Monitor{
		Name:          "test-web-monitor",
		ParentMonitor: "http",
		Interval:      15,
		Timeout:       5,
		SendString:    "GET /test&parms=1<2>3\r\n",
		ReceiveString: "Response &<>",
	}

	s.Client.CreateMonitor(config.Name, config.ParentMonitor, config.Interval, config.Timeout, config.SendString, config.ReceiveString)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriMonitor, config.ParentMonitor), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-web-monitor","defaultsFrom":"http","interval":15,"manualResume":"disabled","recv":"Response &<>","reverse":"disabled","send":"GET /test&parms=1<2>3\\r\\n","timeout":5,"transparent":"disabled"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestAddMonitor() {
	config := &Monitor{
		Name:          "test-web-monitor",
		ParentMonitor: "http",
		Interval:      15,
		Timeout:       5,
		SendString:    "GET /\r\n",
		ReceiveString: "200 OK",
		Username:      "monitoring",
		Password:      "monitoring",
	}

	s.Client.AddMonitor(config)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriMonitor, config.ParentMonitor), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-web-monitor","defaultsFrom":"http","interval":15,"manualResume":"disabled","password":"monitoring","recv":"200 OK","reverse":"disabled","send":"GET /\\r\\n","timeout":5,"transparent":"disabled","username":"monitoring"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestGetMonitor() {
	config := &Monitor{
		Name:          "test-web-monitor",
		ParentMonitor: "http",
	}

	s.Client.GetMonitor(config.Name, config.ParentMonitor)

	assert.Equal(s.T(), "GET", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriMonitor, config.ParentMonitor, config.Name), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestDeleteMonitor() {
	config := &Monitor{
		Name:          "test-web-monitor",
		ParentMonitor: "http",
	}

	s.Client.DeleteMonitor(config.Name, config.ParentMonitor)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriMonitor, config.ParentMonitor, config.Name), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestVirtualServerPolicies() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
		  "kind": "tm:ltm:virtual:policies:policiescollectionstate",
		  "selfLink": "https://localhost/mgmt/tm/ltm/virtual/foo/policies?ver=11.5.1",
		  "policiesReference": {
		    "link": "https://localhost/mgmt/tm/ltm/virtual/foo/policies?ver=11.5.1",
		    "isSubcollection": true,
		    "items": [
		      {
			"kind": "tm:ltm:virtual:policies:policiesstate",
			"name": "policy1",
			"partition": "Common",
			"fullPath": "/Common/policy1",
			"generation": 1,
			"selfLink": "https://localhost/mgmt/tm/ltm/virtual/foo/policies/~Common~policy1?ver=11.5.1"
		      },
		      {
			"kind": "tm:ltm:virtual:policies:policiesstate",
			"name": "policy2",
			"partition": "Common",
			"fullPath": "/Common/policy2",
			"generation": 1,
			"selfLink": "https://localhost/mgmt/tm/ltm/virtual/foo/policies/~Common~policy2?ver=11.5.1"
		      }
		    ]
		  }
		}`))
	}

	p, err := s.Client.VirtualServerPolicyNames("foo")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/foo/policies", uriLtm, uriVirtual), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "/Common/policy1", p[0])
	assert.Equal(s.T(), "/Common/policy2", p[1])
}

func (s *LTMTestSuite) TestInternalDataGroups() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
		  "kind": "tm:ltm:data-group:internal:internalcollectionstate",
		  "selfLink": "https://localhost/mgmt/tm/ltm/data-group/internal?ver=12.1.2",
		  "items": [
		    {
		      "kind": "tm:ltm:data-group:internal:internalstate",
		      "name": "some_data_group",
		      "partition": "Common",
		      "fullPath": "/Common/api.tv2.dk_host_pool_map",
		      "generation": 2552,
		      "selfLink": "https://localhost/mgmt/tm/ltm/data-group/internal/~Common~api.tv2.dk_host_pool_map?ver=12.1.2",
		      "type": "string",
		      "records": [
		        {
		          "name": "jens.medister.api.tv2.dk-hest",
		          "data": "pool-medister"
		        }
		      ]
		    },
		    {
		      "kind": "tm:ltm:data-group:internal:internalstate",
		      "name": "jenkins_whitelisted_paths",
		      "partition": "Common",
		      "fullPath": "/Common/jenkins_whitelisted_paths",
		      "generation": 41,
		      "selfLink": "https://localhost/mgmt/tm/ltm/data-group/internal/~Common~jenkins_whitelisted_paths?ver=12.1.2",
		      "type": "string",
		      "records": [
		        {
		          "name": "/medister",
		          "data": "1"
		        }
		      ]
		    }
		  ]
		}`))
	}

	g, err := s.Client.InternalDataGroups()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "GET", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriDatagroup, uriInternal), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "some_data_group", g.DataGroups[0].Name)
	assert.Equal(s.T(), "jenkins_whitelisted_paths", g.DataGroups[1].Name)
	assert.Equal(s.T(), "/medister", g.DataGroups[1].Records[0].Name)
}

func (s *LTMTestSuite) TestAddInternalDataGroup() {
	config := &DataGroup{
		Name: "test-datagroup",
		Type: "string",
		Records: []DataGroupRecord{
			DataGroupRecord{
				Name: "name1",
				Data: "data1",
			},
			DataGroupRecord{
				Name: "name2",
				Data: "data2",
			},
		},
	}

	s.Client.AddInternalDataGroup(config)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriDatagroup, uriInternal), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-datagroup","type":"string","records":[{"name":"name1","data":"data1"},{"name":"name2","data":"data2"}]}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestModifyInternalDataGroupRecords() {
	dataGroup := "test"

	records := &[]DataGroupRecord{
		DataGroupRecord{
			Name: "name1",
			Data: "data1",
		},
		DataGroupRecord{
			Name: "name42",
			Data: "data42",
		},
	}

	s.Client.ModifyInternalDataGroupRecords(dataGroup, records)

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriDatagroup, uriInternal, dataGroup), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"records":[{"name":"name1","data":"data1"},{"name":"name42","data":"data42"}]}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDeleteInternalDataGroup() {
	dataGroup := "test"
	s.Client.DeleteInternalDataGroup(dataGroup)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriDatagroup, uriInternal, dataGroup), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestGetInternalDataGroupRecords() {
	dataGroup := "test"
	s.Client.GetInternalDataGroupRecords(dataGroup)

	assert.Equal(s.T(), "GET", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriDatagroup, uriInternal, dataGroup), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestSnatPools() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"kind": "tm:ltm:snatpool:snatpoolcollectionstate",
			"selfLink": "https://localhost/mgmt/tm/ltm/snatpool?ver=11.5.3",
			"items": [{
				"kind": "tm:ltm:snatpool:snatpoolstate",
				"name": "mySnatPool",
				"partition": "Common",
				"fullPath": "/Common/mySnatPool",
				"generation": 419,
				"selfLink": "https://localhost/mgmt/tm/ltm/snatpool/~Common~mySnatPool?ver=11.5.3",
				"members": [
					"/Common/10.0.0.1"
				]
			}, {
				"kind": "tm:ltm:snatpool:snatpoolstate",
				"name": "mySnatPool2",
				"partition": "Common",
				"fullPath": "/Common/mySnatPool2",
				"generation": 477,
				"selfLink": "https://localhost/mgmt/tm/ltm/snatpool/~Common~mySnatPool2?ver=11.5.3",
				"members": [
					"/Common/10.0.0.2"
				]
			}]
		}`))
	}

	g, err := s.Client.SnatPools()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "GET", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriSnatPool), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "mySnatPool", g.SnatPools[0].Name)
	assert.Equal(s.T(), "mySnatPool2", g.SnatPools[1].Name)
}

func (s *LTMTestSuite) TestCreateSnatPool() {

	s.Client.CreateSnatPool("mySnatPool", []string{"10.1.1.1", "10.2.2.2"})

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriSnatPool), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"mySnatPool","members":["10.1.1.1","10.2.2.2"]}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestAddSnatPool() {

	mySnatPool := &SnatPool{Name: "mySnatPool", Members: []string{"10.1.1.1", "10.2.2.2"}}

	s.Client.AddSnatPool(mySnatPool)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriSnatPool), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"mySnatPool","members":["10.1.1.1","10.2.2.2"]}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestModifySnatPool() {

	snatPool := "mySnatPool"

	myModifedSnatPool := &SnatPool{Members: []string{"10.0.0.1", "10.0.0.2"}}

	s.Client.ModifySnatPool(snatPool, myModifedSnatPool)

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriSnatPool, snatPool), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"members":["10.0.0.1","10.0.0.2"]}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDeleteSnatPool() {
	snatPool := "mySnatPool"
	s.Client.DeleteSnatPool(snatPool)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriSnatPool, snatPool), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestServerSSLProfiles() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"kind": "tm:ltm:profile:server-ssl:server-sslcollectionstate",
			"selfLink": "https://localhost/mgmt/tm/ltm/profile/server-ssl?ver=11.5.3",
			"items": [{
				"kind": "tm:ltm:profile:server-ssl:server-sslstate",
				"name": "myServerSSL",
				"partition": "Common",
				"fullPath": "/Common/myServerSSL",
				"generation": 784,
				"selfLink": "https://localhost/mgmt/tm/ltm/profile/server-ssl/~Common~myServerSSL?ver=11.5.3",
				"alertTimeout": "10",
				"authenticate": "once",
				"authenticateDepth": 9,
				"cacheSize": 262144,
				"cacheTimeout": 3600,
				"cert": "/Common/default.crt",
				"ciphers": "DEFAULT",
				"defaultsFrom": "/Common/serverssl",
				"expireCertResponseControl": "drop",
				"genericAlert": "enabled",
				"handshakeTimeout": "10",
				"key": "/Common/default.key",
				"modSslMethods": "disabled",
				"mode": "enabled",
				"tmOptions": [
					"dont-insert-empty-fragments"
				],
				"peerCertMode": "ignore",
				"proxySsl": "disabled",
				"renegotiatePeriod": "indefinite",
				"renegotiateSize": "indefinite",
				"renegotiation": "enabled",
				"retainCertificate": "true",
				"secureRenegotiation": "require-strict",
				"serverName": "myserver.contoso.com",
				"sessionTicket": "disabled",
				"sniDefault": "false",
				"sniRequire": "false",
				"sslForwardProxy": "disabled",
				"sslForwardProxyBypass": "disabled",
				"sslSignHash": "any",
				"strictResume": "disabled",
				"uncleanShutdown": "enabled",
				"untrustedCertResponseControl": "drop"
			}, {
				"kind": "tm:ltm:profile:server-ssl:server-sslstate",
				"name": "serverssl",
				"partition": "Common",
				"fullPath": "/Common/serverssl",
				"generation": 1,
				"selfLink": "https://localhost/mgmt/tm/ltm/profile/server-ssl/~Common~serverssl?ver=11.5.3",
				"alertTimeout": "10",
				"authenticate": "once",
				"authenticateDepth": 9,
				"cacheSize": 262144,
				"cacheTimeout": 3600,
				"ciphers": "DEFAULT",
				"expireCertResponseControl": "drop",
				"genericAlert": "enabled",
				"handshakeTimeout": "10",
				"modSslMethods": "disabled",
				"mode": "enabled",
				"tmOptions": [
					"dont-insert-empty-fragments"
				],
				"peerCertMode": "ignore",
				"proxySsl": "disabled",
				"renegotiatePeriod": "indefinite",
				"renegotiateSize": "indefinite",
				"renegotiation": "enabled",
				"retainCertificate": "true",
				"secureRenegotiation": "require-strict",
				"sessionTicket": "disabled",
				"sniDefault": "false",
				"sniRequire": "false",
				"sslForwardProxy": "disabled",
				"sslForwardProxyBypass": "disabled",
				"sslSignHash": "any",
				"strictResume": "disabled",
				"uncleanShutdown": "enabled",
				"untrustedCertResponseControl": "drop"
			}, {
				"kind": "tm:ltm:profile:server-ssl:server-sslstate",
				"name": "serverssl-insecure-compatible",
				"partition": "Common",
				"fullPath": "/Common/serverssl-insecure-compatible",
				"generation": 1,
				"selfLink": "https://localhost/mgmt/tm/ltm/profile/server-ssl/~Common~serverssl-insecure-compatible?ver=11.5.3",
				"alertTimeout": "10",
				"authenticate": "once",
				"authenticateDepth": 9,
				"cacheSize": 262144,
				"cacheTimeout": 3600,
				"ciphers": "!SSLv2:!EXPORT:!DH:RSA+RC4:RSA+AES:RSA+DES:RSA+3DES:ECDHE+AES:ECDHE+3DES:@SPEED",
				"defaultsFrom": "/Common/serverssl",
				"expireCertResponseControl": "drop",
				"genericAlert": "enabled",
				"handshakeTimeout": "10",
				"modSslMethods": "disabled",
				"mode": "enabled",
				"tmOptions": [
					"dont-insert-empty-fragments"
				],
				"peerCertMode": "ignore",
				"proxySsl": "disabled",
				"renegotiatePeriod": "indefinite",
				"renegotiateSize": "indefinite",
				"renegotiation": "enabled",
				"retainCertificate": "true",
				"secureRenegotiation": "request",
				"sessionTicket": "disabled",
				"sniDefault": "false",
				"sniRequire": "false",
				"sslForwardProxy": "disabled",
				"sslForwardProxyBypass": "disabled",
				"sslSignHash": "any",
				"strictResume": "disabled",
				"uncleanShutdown": "enabled",
				"untrustedCertResponseControl": "drop"
			}]
		}`))
	}

	g, err := s.Client.ServerSSLProfiles()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "GET", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriProfile, uriServerSSL), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "myServerSSL", g.ServerSSLProfiles[0].Name)
	assert.Equal(s.T(), "serverssl", g.ServerSSLProfiles[1].Name)
	assert.Equal(s.T(), "serverssl-insecure-compatible", g.ServerSSLProfiles[2].Name)
}

func (s *LTMTestSuite) TestCreateServerSSLProfile() {
	s.Client.CreateServerSSLProfile("myServerSSL", "/Common/serverssl")

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriProfile, uriServerSSL), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"myServerSSL","defaultsFrom":"/Common/serverssl"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestModifyServerSSLProfile() {
	serverSSLProfile := "myServerSSL"
	myModifedServerSSLProfile := &ServerSSLProfile{Mode: "enabled", Cert: "/Common/default.crt", Key: "/Common/default.key", ServerName: "myserver.contoso.com"}

	s.Client.ModifyServerSSLProfile(serverSSLProfile, myModifedServerSSLProfile)

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriProfile, uriServerSSL, serverSSLProfile), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"cert":"/Common/default.crt","key":"/Common/default.key","mode":"enabled","serverName":"myserver.contoso.com"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDeleteServerSSLProfile() {
	serverSSLProfile := "myServerSSL"
	s.Client.DeleteServerSSLProfile(serverSSLProfile)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriProfile, uriServerSSL, serverSSLProfile), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestClientSSLProfiles() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"kind": "tm:ltm:profile:client-ssl:client-sslcollectionstate",
			"selfLink": "https://localhost/mgmt/tm/ltm/profile/client-ssl?ver=11.5.3",
			"items": [{
				"kind": "tm:ltm:profile:client-ssl:client-sslstate",
				"name": "clientssl",
				"partition": "Common",
				"fullPath": "/Common/clientssl",
				"generation": 1,
				"selfLink": "https://localhost/mgmt/tm/ltm/profile/client-ssl/~Common~clientssl?ver=11.5.3",
				"alertTimeout": "10",
				"allowNonSsl": "disabled",
				"authenticate": "once",
				"authenticateDepth": 9,
				"cacheSize": 262144,
				"cacheTimeout": 3600,
				"cert": "/Common/default.crt",
				"certExtensionIncludes": [
					"basic-constraints",
					"subject-alternative-name"
				],
				"certLifespan": 30,
				"certLookupByIpaddrPort": "disabled",
				"ciphers": "DEFAULT",
				"forwardProxyBypassDefaultAction": "intercept",
				"genericAlert": "enabled",
				"handshakeTimeout": "10",
				"inheritCertkeychain": "false",
				"key": "/Common/default.key",
				"modSslMethods": "disabled",
				"mode": "enabled",
				"tmOptions": [
					"dont-insert-empty-fragments"
				],
				"peerCertMode": "ignore",
				"proxySsl": "disabled",
				"renegotiateMaxRecordDelay": "indefinite",
				"renegotiatePeriod": "indefinite",
				"renegotiateSize": "indefinite",
				"renegotiation": "enabled",
				"retainCertificate": "true",
				"secureRenegotiation": "require",
				"sessionTicket": "disabled",
				"sniDefault": "false",
				"sniRequire": "false",
				"sslForwardProxy": "disabled",
				"sslForwardProxyBypass": "disabled",
				"sslSignHash": "any",
				"strictResume": "disabled",
				"uncleanShutdown": "enabled",
				"certKeyChain": [{
					"name": "default",
					"cert": "/Common/default.crt",
					"key": "/Common/default.key"
				}]
			}, {
				"kind": "tm:ltm:profile:client-ssl:client-sslstate",
				"name": "clientssl-insecure-compatible",
				"partition": "Common",
				"fullPath": "/Common/clientssl-insecure-compatible",
				"generation": 1,
				"selfLink": "https://localhost/mgmt/tm/ltm/profile/client-ssl/~Common~clientssl-insecure-compatible?ver=11.5.3",
				"alertTimeout": "10",
				"allowNonSsl": "disabled",
				"authenticate": "once",
				"authenticateDepth": 9,
				"cacheSize": 262144,
				"cacheTimeout": 3600,
				"cert": "/Common/default.crt",
				"certExtensionIncludes": [
					"basic-constraints",
					"subject-alternative-name"
				],
				"certLifespan": 30,
				"certLookupByIpaddrPort": "disabled",
				"ciphers": "!SSLv2:ALL:!DH:!ADH:!EDH:@SPEED",
				"defaultsFrom": "/Common/clientssl",
				"forwardProxyBypassDefaultAction": "intercept",
				"genericAlert": "enabled",
				"handshakeTimeout": "10",
				"inheritCertkeychain": "true",
				"key": "/Common/default.key",
				"modSslMethods": "disabled",
				"mode": "enabled",
				"tmOptions": [
					"dont-insert-empty-fragments"
				],
				"peerCertMode": "ignore",
				"proxySsl": "disabled",
				"renegotiateMaxRecordDelay": "indefinite",
				"renegotiatePeriod": "indefinite",
				"renegotiateSize": "indefinite",
				"renegotiation": "enabled",
				"retainCertificate": "true",
				"secureRenegotiation": "request",
				"sessionTicket": "disabled",
				"sniDefault": "false",
				"sniRequire": "false",
				"sslForwardProxy": "disabled",
				"sslForwardProxyBypass": "disabled",
				"sslSignHash": "any",
				"strictResume": "disabled",
				"uncleanShutdown": "enabled",
				"certKeyChain": [{
					"name": "default",
					"cert": "/Common/default.crt",
					"key": "/Common/default.key"
				}]
			}, {
				"kind": "tm:ltm:profile:client-ssl:client-sslstate",
				"name": "myClientSSL",
				"partition": "Common",
				"fullPath": "/Common/myClientSSL",
				"generation": 727,
				"selfLink": "https://localhost/mgmt/tm/ltm/profile/client-ssl/~Common~myClientSSL?ver=11.5.3",
				"alertTimeout": "10",
				"allowNonSsl": "disabled",
				"authenticate": "once",
				"authenticateDepth": 9,
				"cacheSize": 262144,
				"cacheTimeout": 3600,
				"certExtensionIncludes": [
					"basic-constraints",
					"subject-alternative-name"
				],
				"certLifespan": 30,
				"certLookupByIpaddrPort": "disabled",
				"ciphers": "DEFAULT",
				"defaultsFrom": "/Common/clientssl",
				"forwardProxyBypassDefaultAction": "intercept",
				"genericAlert": "enabled",
				"handshakeTimeout": "10",
				"inheritCertkeychain": "false",
				"modSslMethods": "disabled",
				"mode": "enabled",
				"tmOptions": [
					"dont-insert-empty-fragments"
				],
				"peerCertMode": "ignore",
				"proxySsl": "disabled",
				"renegotiateMaxRecordDelay": "indefinite",
				"renegotiatePeriod": "indefinite",
				"renegotiateSize": "indefinite",
				"renegotiation": "enabled",
				"retainCertificate": "true",
				"secureRenegotiation": "require",
				"sessionTicket": "disabled",
				"sniDefault": "false",
				"sniRequire": "false",
				"sslForwardProxy": "disabled",
				"sslForwardProxyBypass": "disabled",
				"sslSignHash": "any",
				"strictResume": "disabled",
				"uncleanShutdown": "enabled",
				"certKeyChain": [{
					"name": "\"\""
				}]
			}, {
				"kind": "tm:ltm:profile:client-ssl:client-sslstate",
				"name": "wom-default-clientssl",
				"partition": "Common",
				"fullPath": "/Common/wom-default-clientssl",
				"generation": 1,
				"selfLink": "https://localhost/mgmt/tm/ltm/profile/client-ssl/~Common~wom-default-clientssl?ver=11.5.3",
				"alertTimeout": "10",
				"allowNonSsl": "enabled",
				"authenticate": "once",
				"authenticateDepth": 9,
				"cacheSize": 262144,
				"cacheTimeout": 3600,
				"cert": "/Common/default.crt",
				"certExtensionIncludes": [
					"basic-constraints",
					"subject-alternative-name"
				],
				"certLifespan": 30,
				"certLookupByIpaddrPort": "disabled",
				"ciphers": "DEFAULT",
				"defaultsFrom": "/Common/clientssl",
				"forwardProxyBypassDefaultAction": "intercept",
				"genericAlert": "enabled",
				"handshakeTimeout": "10",
				"inheritCertkeychain": "true",
				"key": "/Common/default.key",
				"modSslMethods": "disabled",
				"mode": "enabled",
				"tmOptions": [
					"dont-insert-empty-fragments"
				],
				"peerCertMode": "ignore",
				"proxySsl": "disabled",
				"renegotiateMaxRecordDelay": "indefinite",
				"renegotiatePeriod": "indefinite",
				"renegotiateSize": "indefinite",
				"renegotiation": "enabled",
				"retainCertificate": "true",
				"secureRenegotiation": "require",
				"sessionTicket": "disabled",
				"sniDefault": "false",
				"sniRequire": "false",
				"sslForwardProxy": "disabled",
				"sslForwardProxyBypass": "disabled",
				"sslSignHash": "any",
				"strictResume": "disabled",
				"uncleanShutdown": "enabled",
				"certKeyChain": [{
					"name": "default",
					"cert": "/Common/default.crt",
					"key": "/Common/default.key"
				}]
			}]
		}`))
	}

	g, err := s.Client.ClientSSLProfiles()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "GET", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriProfile, uriClientSSL), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "clientssl", g.ClientSSLProfiles[0].Name)
	assert.Equal(s.T(), "clientssl-insecure-compatible", g.ClientSSLProfiles[1].Name)
	assert.Equal(s.T(), "myClientSSL", g.ClientSSLProfiles[2].Name)
	assert.Equal(s.T(), "wom-default-clientssl", g.ClientSSLProfiles[3].Name)
}

func (s *LTMTestSuite) TestCreateClientSSLProfile() {
	s.Client.CreateClientSSLProfile("myClientSSL", "/Common/clientssl")

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriProfile, uriClientSSL), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"myClientSSL","defaultsFrom":"/Common/clientssl"}`, s.LastRequestBody)
}

// Add additional test for more complex certKeyChain

func (s *LTMTestSuite) TestModifyClientSSLProfile() {
	clientSSLProfile := "myClientSSL"
	myModifedClientSSLProfile := &ClientSSLProfile{Mode: "enabled", Cert: "/Common/default.crt", Key: "/Common/default.key", ServerName: "myserver.contoso.com"}

	s.Client.ModifyClientSSLProfile(clientSSLProfile, myModifedClientSSLProfile)

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriProfile, uriClientSSL, clientSSLProfile), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"cert":"/Common/default.crt","key":"/Common/default.key","mode":"enabled","serverName":"myserver.contoso.com"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDeleteClientSSLProfile() {
	clientSSLProfile := "myClientSSL"
	s.Client.DeleteClientSSLProfile(clientSSLProfile)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriProfile, uriClientSSL, clientSSLProfile), s.LastRequest.URL.Path)
}

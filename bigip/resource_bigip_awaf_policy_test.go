/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipAWFPolicyUnitInvalid(t *testing.T) {
	resourceName := "mytestpolicy"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBigipAwafPolicyInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipAWFPolicyUnitCreate(t *testing.T) {
	resourceName := "mytestpolicy"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"kind": "tm:sys:provision:provisionstate",
    "name": "asm",
    "fullPath": "asm",
    "selfLink": "https://localhost/mgmt/tm/sys/provision/asm?ver=16.1.0",
    "cpuRatio": 0,
    "diskRatio": 0,
    "level": "nominal",
    "memoryRatio": 0}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/file-transfer/uploads/mytestpolicy.json", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Range", "[0-416/417]")
		r.Header.Add("Content-Type", "[application/octet-stream]")
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/import-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
            "isBase64": false,
            "executionStartTime": "2022-12-29T05:42:25Z",
            "status": "COMPLETED",
            "lastUpdateMicros": 1.672292549e+15,
            "getPolicyAttributesOnly": false,
            "fullPath": "/Common/mytestpolicy",
            "kind": "tm:asm:tasks:import-policy:import-policy-taskstate",
            "selfLink": "https://localhost/mgmt/tm/asm/tasks/import-policy/PW6zgmf9L7C4I8UdUjaWgw?ver=16.1.0",
            "filename": "mytestpolicy.json",
            "endTime": "2022-12-29T05:42:30Z",
            "id": "PW6zgmf9L7C4I8UdUjaWgw",
            "startTime": "2022-12-29T05:42:25Z",
            "retainInheritanceSettings": false,
            "result": {
                "policyReference": {
                    "link": "https://localhost/mgmt/tm/asm/policies/?ver=16.1.0"
                },
                "message": "The operation was completed successfully. The security policy name is '/Common/mytestpolicy'. Policy Template set to API Security."
            }
        }`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/import-policy/PW6zgmf9L7C4I8UdUjaWgw", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
            "isBase64": false,
            "executionStartTime": "2022-12-29T05:42:25Z",
            "status": "COMPLETED",
            "lastUpdateMicros": 1.672292549e+15,
            "getPolicyAttributesOnly": false,
            "fullPath": "/Common/mytestpolicy",
            "kind": "tm:asm:tasks:import-policy:import-policy-taskstate",
            "selfLink": "https://localhost/mgmt/tm/asm/tasks/import-policy/PW6zgmf9L7C4I8UdUjaWgw?ver=16.1.0",
            "filename": "mytestpolicy.json",
            "endTime": "2022-12-29T05:42:30Z",
            "id": "PW6zgmf9L7C4I8UdUjaWgw",
            "startTime": "2022-12-29T05:42:25Z",
            "retainInheritanceSettings": false,
            "result": {
                "policyReference": {
                    "link": "https://localhost/mgmt/tm/asm/policies/?ver=16.1.0"
                },
                "message": "The operation was completed successfully. The security policy name is '/Common/mytestpolicy'. Policy Template set to API Security."
            }
        }`)
	})
	http.HandleFunc("/mgmt/tm/asm/policies/?$filter=contains(name,'mytestpolicy')+and+contains(partition,'Common')", func(w http.ResponseWriter, r *http.Request) {
	})

	mux.HandleFunc("/mgmt/tm/asm/policies/?$filter=contains(name,'mytestpolicy')+and+contains(partition,'Common')", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "*" {
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, `{}`)
		}
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/apply-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})

	//	req := httptest.NewRequest(http.MethodGet, "/mgmt/tm/asm/policies/?$filter=contains(name,'mytestpolicy')+and+contains(partition,'Common')", nil)
	//	t.Logf("URL:%+v", req.URL.RawPath)
	//	w := httptest.NewRecorder()
	//	w.WriteHeader(200)
	//	w.WriteString(`{
	//    "kind": "tm:asm:policies:policycollectionstate",
	//    "selfLink": "https://localhost/mgmt/tm/asm/policies?ver=16.1.0&$filter=contains%28name%2C%27testraviawaf%27%29%20and%20contains%28partition%2C%27Common%27%29",
	//    "totalItems": 1,
	//    "items": [
	//        {
	//            "signatureRequirementReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/signature-requirements?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "plainTextProfileReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/plain-text-profiles?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "enablePassiveMode": false,
	//            "behavioralEnforcementReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/behavioral-enforcement?ver=16.1.0"
	//            },
	//            "dataGuardReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/data-guard?ver=16.1.0"
	//            },
	//            "createdDatetime": "2022-12-05T08:12:29Z",
	//            "databaseProtectionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/database-protection?ver=16.1.0"
	//            },
	//            "cookieSettingsReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/cookie-settings?ver=16.1.0"
	//            },
	//            "csrfUrlReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/csrf-urls?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "versionLastChange": " Security Policy /Common/testraviawaf [add]: Parent Policy was set to empty value.\nType was set to Security.\nEncoding Selected was set to true.\nApplication Language was set to utf-8.\nActive was set to false.\nPolicy Name was set to /Common/testraviawaf. { audit: policy = /Common/testraviawaf, component = tsconfd }",
	//            "name": "testraviawaf",
	//            "caseInsensitive": false,
	//            "headerSettingsReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/header-settings?ver=16.1.0"
	//            },
	//            "sectionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/sections?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "flowReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/flows?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "loginPageReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/login-pages?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "description": "Rapid Deployment-1",
	//            "fullPath": "/Common/testraviawaf",
	//            "policyBuilderParameterReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-parameter?ver=16.1.0"
	//            },
	//            "hasParent": false,
	//            "threatCampaignReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/threat-campaigns?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "partition": "Common",
	//            "managedByBewaf": false,
	//            "csrfProtectionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/csrf-protection?ver=16.1.0"
	//            },
	//            "graphqlProfileReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/graphql-profiles?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "policyAntivirusReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/antivirus?ver=16.1.0"
	//            },
	//            "kind": "tm:asm:policies:policystate",
	//            "virtualServers": [],
	//            "policyBuilderCookieReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-cookie?ver=16.1.0"
	//            },
	//            "ipIntelligenceReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/ip-intelligence?ver=16.1.0"
	//            },
	//            "protocolIndependent": false,
	//            "sessionAwarenessSettingsReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/session-tracking?ver=16.1.0"
	//            },
	//            "policyBuilderUrlReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-url?ver=16.1.0"
	//            },
	//            "policyBuilderServerTechnologiesReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-server-technologies?ver=16.1.0"
	//            },
	//            "policyBuilderFiletypeReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-filetype?ver=16.1.0"
	//            },
	//            "signatureSetReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/signature-sets?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "parameterReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/parameters?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "applicationLanguage": "utf-8",
	//            "enforcementMode": "transparent",
	//            "loginEnforcementReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/login-enforcement?ver=16.1.0"
	//            },
	//            "openApiFileReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/open-api-files?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "navigationParameterReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/navigation-parameters?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "gwtProfileReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/gwt-profiles?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "webhookReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/webhooks?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "whitelistIpReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/whitelist-ips?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "historyRevisionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/history-revisions?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "policyBuilderReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder?ver=16.1.0"
	//            },
	//            "responsePageReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/response-pages?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "vulnerabilityAssessmentReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/vulnerability-assessment?ver=16.1.0"
	//            },
	//            "serverTechnologyReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/server-technologies?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "cookieReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/cookies?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "blockingSettingReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/blocking-settings?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "hostNameReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/host-names?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "versionDeviceName": "ecosyshydbigip16.com",
	//            "selfLink": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ?ver=16.1.0",
	//            "threatCampaignSettingReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/threat-campaign-settings?ver=16.1.0"
	//            },
	//            "signatureReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/signatures?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "policyBuilderRedirectionProtectionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-redirection-protection?ver=16.1.0"
	//            },
	//            "filetypeReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/filetypes?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "id": "Xnm5-sgcAnOxSsK2uDZTkQ",
	//            "modifierName": "",
	//            "manualVirtualServers": [],
	//            "versionDatetime": "2022-12-05T08:12:30Z",
	//            "ssrfHostReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/ssrf-hosts?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "subPath": "/Common",
	//            "sessionTrackingStatusReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/session-tracking-statuses?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "active": false,
	//            "auditLogReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/audit-logs?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "disallowedGeolocationReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/disallowed-geolocations?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "redirectionProtectionDomainReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/redirection-protection-domains?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "type": "security",
	//            "signatureSettingReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/signature-settings?ver=16.1.0"
	//            },
	//            "deceptionResponsePageReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/deception-response-pages?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "websocketUrlReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/websocket-urls?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "xmlProfileReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/xml-profiles?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "methodReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/methods?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "vulnerabilityReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/vulnerabilities?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "redirectionProtectionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/redirection-protection?ver=16.1.0"
	//            },
	//            "policyBuilderSessionsAndLoginsReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-sessions-and-logins?ver=16.1.0"
	//            },
	//            "templateReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policy-templates/EzpBNMs9gbVsF5uuiBjYDw?ver=16.1.0",
	//                "title": "Rapid Deployment Policy"
	//            },
	//            "policyBuilderHeaderReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-header?ver=16.1.0"
	//            },
	//            "creatorName": "tsconfd",
	//            "urlReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/urls?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "headerReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/headers?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "actionItemReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/action-items?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "microserviceReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/microservices?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "xmlValidationFileReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/xml-validation-files?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "lastUpdateMicros": 0,
	//            "jsonProfileReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/json-profiles?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "bruteForceAttackPreventionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/brute-force-attack-preventions?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "disabledActionItemReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/disabled-action-items?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "jsonValidationFileReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/json-validation-files?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "extractionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/extractions?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "characterSetReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/character-sets?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "suggestionReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/suggestions?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "deceptionSettingsReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/deception-settings?ver=16.1.0"
	//            },
	//            "isModified": false,
	//            "sensitiveParameterReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/sensitive-parameters?ver=16.1.0",
	//                "isSubCollection": true
	//            },
	//            "generalReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/general?ver=16.1.0"
	//            },
	//            "versionPolicyName": "/Common/testraviawaf",
	//            "policyBuilderCentralConfigurationReference": {
	//                "link": "https://localhost/mgmt/tm/asm/policies/Xnm5-sgcAnOxSsK2uDZTkQ/policy-builder-central-configuration?ver=16.1.0"
	//            }
	//        }
	//    ]
	//}`)
	//	res := w.Result()
	//	defer res.Body.Close()
	//	data, _ := io.ReadAll(res.Body)
	//	if string(data) != "ABC" {
	//		t.Errorf("expected ABC got %v", string(data))
	//	}

	//mux.HandleFunc("/mgmt/tm/asm/policies/?$filter=contains(name,'mytestpolicy')+and+contains(partition,'Common')", func(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.String() == "/mgmt/tm/asm/policies/?$filter=contains(name,'mytestpolicy')+and+contains(partition,'Common')" {
	//		w.WriteHeader(200)
	//	}
	//	_, _ = fmt.Fprintf(w, `{"items": [{
	//	"signatureRequirementReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-requirements?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"plainTextProfileReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/plain-text-profiles?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"enablePassiveMode": false,
	//	"behavioralEnforcementReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/behavioral-enforcement?ver=16.1.0"
	//	},
	//	"dataGuardReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/data-guard?ver=16.1.0"
	//	},
	//	"createdDatetime": "2022-12-28T10:27:41Z",
	//	"databaseProtectionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/database-protection?ver=16.1.0"
	//	},
	//	"cookieSettingsReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/cookie-settings?ver=16.1.0"
	//	},
	//	"csrfUrlReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/csrf-urls?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"versionLastChange": "Valid Host Name testhostname [add]: Include Subdomains was set to enabled.\nHost Name was set to testhostname. { audit: policy = /Common/mytestpolicy.app/mytestpolicy, username = admin, client IP = 172.18.236.59 }",
	//	"name": "mytestpolicy",
	//	"caseInsensitive": true,
	//	"headerSettingsReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/header-settings?ver=16.1.0"
	//	},
	//	"sectionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/sections?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"flowReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/flows?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"loginPageReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/login-pages?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"description": "Generic template for OWA Exchange 2016",
	//	"fullPath": "/Common/mytestpolicy.app/mytestpolicy",
	//	"policyBuilderParameterReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-parameter?ver=16.1.0"
	//	},
	//	"hasParent": false,
	//	"threatCampaignReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/threat-campaigns?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"partition": "Common",
	//	"managedByBewaf": false,
	//	"csrfProtectionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/csrf-protection?ver=16.1.0"
	//	},
	//	"graphqlProfileReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/graphql-profiles?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"policyAntivirusReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/antivirus?ver=16.1.0"
	//	},
	//	"kind": "tm:asm:policies:policystate",
	//	"virtualServers": [
	//	"/Common/mytestpolicy.app/mytestpolicy_vs"
	//	],
	//	"policyBuilderCookieReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-cookie?ver=16.1.0"
	//	},
	//	"ipIntelligenceReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/ip-intelligence?ver=16.1.0"
	//	},
	//	"protocolIndependent": true,
	//	"sessionAwarenessSettingsReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/session-tracking?ver=16.1.0"
	//	},
	//	"policyBuilderUrlReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-url?ver=16.1.0"
	//	},
	//	"policyBuilderServerTechnologiesReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-server-technologies?ver=16.1.0"
	//	},
	//	"policyBuilderFiletypeReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-filetype?ver=16.1.0"
	//	},
	//	"signatureSetReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-sets?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"parameterReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/parameters?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"applicationLanguage": "utf-8",
	//	"enforcementMode": "blocking",
	//	"loginEnforcementReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/login-enforcement?ver=16.1.0"
	//	},
	//	"openApiFileReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/open-api-files?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"navigationParameterReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/navigation-parameters?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"applicationService": "/Common/mytestpolicy.app/mytestpolicy",
	//	"gwtProfileReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/gwt-profiles?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"webhookReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/webhooks?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"whitelistIpReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/whitelist-ips?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"historyRevisionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/history-revisions?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"policyBuilderReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder?ver=16.1.0"
	//	},
	//	"responsePageReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/response-pages?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"vulnerabilityAssessmentReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/vulnerability-assessment?ver=16.1.0"
	//	},
	//	"serverTechnologyReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/server-technologies?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"cookieReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/cookies?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"blockingSettingReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/blocking-settings?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"hostNameReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/host-names?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"versionDeviceName": "ecosyshydbigip16.com",
	//	"selfLink": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg?ver=16.1.0",
	//	"threatCampaignSettingReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/threat-campaign-settings?ver=16.1.0"
	//	},
	//	"signatureReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signatures?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"policyBuilderRedirectionProtectionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-redirection-protection?ver=16.1.0"
	//	},
	//	"filetypeReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/filetypes?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"id": "BLDetCAoCevvoh-INfzTFg",
	//	"modifierName": "admin",
	//	"manualVirtualServers": [],
	//	"versionDatetime": "2022-12-28T12:40:28Z",
	//	"ssrfHostReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/ssrf-hosts?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"subPath": "/Common/mytestpolicy.app",
	//	"sessionTrackingStatusReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/session-tracking-statuses?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"active": true,
	//	"auditLogReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/audit-logs?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"disallowedGeolocationReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/disallowed-geolocations?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"redirectionProtectionDomainReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/redirection-protection-domains?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"applicationServiceManagedUpdatesOnly": false,
	//	"type": "security",
	//	"signatureSettingReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-settings?ver=16.1.0"
	//	},
	//	"deceptionResponsePageReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/deception-response-pages?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"websocketUrlReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/websocket-urls?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"xmlProfileReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/xml-profiles?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"methodReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/methods?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"vulnerabilityReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/vulnerabilities?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"redirectionProtectionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/redirection-protection?ver=16.1.0"
	//	},
	//	"policyBuilderSessionsAndLoginsReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-sessions-and-logins?ver=16.1.0"
	//	},
	//	"templateReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policy-templates/clM9mdnuWVeuPOydwWNWnA?ver=16.1.0",
	//	"title": "OWA Exchange 2016"
	//	},
	//	"policyBuilderHeaderReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-header?ver=16.1.0"
	//	},
	//	"creatorName": "admin",
	//	"urlReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/urls?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"headerReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/headers?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"actionItemReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/action-items?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"microserviceReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/microservices?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"xmlValidationFileReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/xml-validation-files?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"lastUpdateMicros": 1.672231238e+15,
	//	"jsonProfileReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/json-profiles?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"bruteForceAttackPreventionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/brute-force-attack-preventions?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"disabledActionItemReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/disabled-action-items?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"jsonValidationFileReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/json-validation-files?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"extractionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/extractions?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"characterSetReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/character-sets?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"suggestionReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/suggestions?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"deceptionSettingsReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/deception-settings?ver=16.1.0"
	//	},
	//	"isModified": false,
	//	"sensitiveParameterReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/sensitive-parameters?ver=16.1.0",
	//	"isSubCollection": true
	//	},
	//	"generalReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/general?ver=16.1.0"
	//	},
	//	"versionPolicyName": "/Common/mytestpolicy.app/mytestpolicy",
	//	"policyBuilderCentralConfigurationReference": {
	//	"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-central-configuration?ver=16.1.0"}
	//	}]}`)
	//})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipAwafPolicyCreate(resourceName, server.URL),
			},
			//{
			//	Config: testAccBigipAwafPolicyModify(resourceName, server.URL),
			//	//ExpectNonEmptyPlan: true,
			//},
		},
	})
}

//
//func TestAccBigipAWFPolicyUnitReadError(t *testing.T) {
//	resourceName := "mytestpolicy"
//	setup()
//	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//	})
//	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//		_, _ = fmt.Fprintf(w, `{}`)
//	})
//	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(202)
//		_, _ = fmt.Fprintf(w, `{
//    "id": "50ce5959-a256-463d-92e9-eee11b20d229",
//    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b",
//    "result": {
//        "class": "Result",
//        "code": 202,
//        "status": "RUNNING",
//        "message": "processing"
//    },
//    "declaration": {
//        "schemaVersion": "1.20.0",
//        "class": "Device",
//        "async": true,
//        "label": "my BIG-IP declaration for declarative onboarding",
//        "Common": {
//            "class": "Tenant",
//            "hostname": "bigip1.example.com",
//            "ravinder": {
//                "class": "User",
//                "userType": "regular",
//                "partitionAccess": {
//                    "Common": {
//                        "role": "guest"
//                    }
//                },
//                "shell": "tmsh"
//            }
//        }
//    }
//}`)
//	})
//	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229", func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(202)
//		_, _ = fmt.Fprintf(w, `{
//    "id": "50ce5959-a256-463d-92e9-eee11b20d229",
//    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229",
//    "result": {
//        "class": "Result",
//        "code": 202,
//        "status": "ROLLING_BACK",
//        "message": "invalid config - rolling back",
//        "errors": [
//            "01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\""
//        ]
//    },
//    "declaration": {
//        "schemaVersion": "1.20.0",
//        "class": "Device",
//        "async": true,
//        "label": "my BIG-IP declaration for declarative onboarding",
//        "Common": {
//            "class": "Tenant",
//            "hostname": "bigip1.example.com",
//            "ravinder": {
//                "class": "User",
//                "userType": "regular",
//                "partitionAccess": {
//                    "Common": {
//                        "role": "guest"
//                    }
//                },
//                "shell": "tmsh"
//            }
//        }
//    }
//}`)
//	})
//	defer teardown()
//	resource.Test(t, resource.TestCase{
//		IsUnitTest: true,
//		Providers:  testProviders,
//		Steps: []resource.TestStep{
//			{
//				Config:      testAccBigipAwafPolicyCreate(resourceName, server.URL),
//				ExpectError: regexp.MustCompile("Error while reading the response body :map\\[class:Result code:202 errors:\\[01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\"] message:invalid config - rolling back status:ROLLING_BACK]"),
//			},
//		},
//	})
//}

//
//func TestAccBigipAWFPolicyUnitReadError(t *testing.T) {
//	resourceName := "regkeypool_name"
//	setup()
//	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//	})
//	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//		_, _ = fmt.Fprintf(w, `{}`)
//	})
//	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		_, _ = fmt.Fprintf(w, `{"name":"%s","destinationAddress":"3.10.11.2/32","ipsecPolicyReference":{},"sourceAddress":"2.10.11.12/32"}`, resourceName)
//	})
//	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector/~Common~test-traffic-selector", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		http.Error(w, "The requested IPsec Trafficselector (/Common/test-traffic-selector) was not found", http.StatusNotFound)
//	})
//
//	defer teardown()
//	resource.Test(t, resource.TestCase{
//		IsUnitTest: true,
//		Providers:  testProviders,
//		Steps: []resource.TestStep{
//			{
//				Config:      testAccBigipAwafPolicyCreate(resourceName, server.URL),
//				ExpectError: regexp.MustCompile("HTTP 404 :: The requested IPsec Trafficselector \\(/Common/test-traffic-selector\\) was not found"),
//			},
//		},
//	})
//}
//
//func TestAccBigipAWFPolicyUnitCreateError(t *testing.T) {
//	resourceName := "regkeypool_name"
//	httpDefault := "/Common/http"
//	setup()
//	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//	})
//	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//		_, _ = fmt.Fprintf(w, `{}`)
//	})
//	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		_, _ = fmt.Fprintf(w, `{"name":"/Common/testhttp##","defaultsFrom":"%s", "basicAuthRealm": "none"}`, httpDefault)
//		http.Error(w, "The requested object name (/Common/testravi##) is invalid", http.StatusNotFound)
//	})
//	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector/~Common~test-traffic-selector", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		http.Error(w, "The requested HTTP Profile (/Common/test-traffic-selector) was not found", http.StatusNotFound)
//	})
//
//	defer teardown()
//	resource.Test(t, resource.TestCase{
//		IsUnitTest: true,
//		Providers:  testProviders,
//		Steps: []resource.TestStep{
//			{
//				Config:      testAccBigipAwafPolicyCreate(resourceName, server.URL),
//				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-traffic-selector\\) was not found"),
//			},
//		},
//	})
//}

func testAccBigipAwafPolicyInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_waf_policy" "test-awaf" {
  name                 = "%s"
  partition            = "Common"
  template_name        = "POLICY_TEMPLATE_API_SECURITY"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  description          = "Rapid Deployment-2"
  invalidkey = "foo"
  policy_builder {
    learning_mode = "disabled"
  }
  server_technologies = ["MySQL", "Unix/Linux", "MongoDB"]
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testAccBigipAwafPolicyCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_waf_policy" "test-awaf" {
  name                 = "%s"
  partition            = "Common"
  template_name        = "POLICY_TEMPLATE_API_SECURITY"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  description          = "Rapid Deployment-2"
  policy_builder {
    learning_mode = "disabled"
  }
  server_technologies = ["MySQL", "Unix/Linux", "MongoDB"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testAccBigipAwafPolicyModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_waf_policy" "test-awaf" {
  name                 = "%s"
  partition            = "Common"
  template_name        = "POLICY_TEMPLATE_API_SECURITY"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  description          = "Rapid Deployment-2"
  policy_builder {
    learning_mode = "disabled"
  }
  server_technologies = ["MySQL", "Unix/Linux", "MongoDB"]
}

provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

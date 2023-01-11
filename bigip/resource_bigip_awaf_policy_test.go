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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/mgmt/tm/asm/policies/?$filter=contains(name,'mytestpolicy')+and+contains(partition,'Common')" {
			_, _ = fmt.Fprintf(w, `{"items": [{
			"signatureRequirementReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-requirements?ver=16.1.0",
			"isSubCollection": true
			},
			"plainTextProfileReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/plain-text-profiles?ver=16.1.0",
			"isSubCollection": true
			},
			"enablePassiveMode": false,
			"behavioralEnforcementReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/behavioral-enforcement?ver=16.1.0"
			},
			"dataGuardReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/data-guard?ver=16.1.0"
			},
			"createdDatetime": "2022-12-28T10:27:41Z",
			"databaseProtectionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/database-protection?ver=16.1.0"
			},
			"cookieSettingsReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/cookie-settings?ver=16.1.0"
			},
			"csrfUrlReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/csrf-urls?ver=16.1.0",
			"isSubCollection": true
			},
			"versionLastChange": "Valid Host Name testhostname [add]: Include Subdomains was set to enabled.\nHost Name was set to testhostname. { audit: policy = /Common/mytestpolicy.app/mytestpolicy, username = admin, client IP = 172.18.236.59 }",
			"name": "mytestpolicy",
			"caseInsensitive": true,
			"headerSettingsReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/header-settings?ver=16.1.0"
			},
			"sectionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/sections?ver=16.1.0",
			"isSubCollection": true
			},
			"flowReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/flows?ver=16.1.0",
			"isSubCollection": true
			},
			"loginPageReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/login-pages?ver=16.1.0",
			"isSubCollection": true
			},
			"description": "Generic template for OWA Exchange 2016",
			"fullPath": "/Common/mytestpolicy.app/mytestpolicy",
			"policyBuilderParameterReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-parameter?ver=16.1.0"
			},
			"hasParent": false,
			"threatCampaignReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/threat-campaigns?ver=16.1.0",
			"isSubCollection": true
			},
			"partition": "Common",
			"managedByBewaf": false,
			"csrfProtectionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/csrf-protection?ver=16.1.0"
			},
			"graphqlProfileReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/graphql-profiles?ver=16.1.0",
			"isSubCollection": true
			},
			"policyAntivirusReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/antivirus?ver=16.1.0"
			},
			"kind": "tm:asm:policies:policystate",
			"virtualServers": [
			"/Common/mytestpolicy.app/mytestpolicy_vs"
			],
			"policyBuilderCookieReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-cookie?ver=16.1.0"
			},
			"ipIntelligenceReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/ip-intelligence?ver=16.1.0"
			},
			"protocolIndependent": true,
			"sessionAwarenessSettingsReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/session-tracking?ver=16.1.0"
			},
			"policyBuilderUrlReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-url?ver=16.1.0"
			},
			"policyBuilderServerTechnologiesReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-server-technologies?ver=16.1.0"
			},
			"policyBuilderFiletypeReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-filetype?ver=16.1.0"
			},
			"signatureSetReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-sets?ver=16.1.0",
			"isSubCollection": true
			},
			"parameterReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/parameters?ver=16.1.0",
			"isSubCollection": true
			},
			"applicationLanguage": "utf-8",
			"enforcementMode": "blocking",
			"loginEnforcementReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/login-enforcement?ver=16.1.0"
			},
			"openApiFileReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/open-api-files?ver=16.1.0",
			"isSubCollection": true
			},
			"navigationParameterReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/navigation-parameters?ver=16.1.0",
			"isSubCollection": true
			},
			"applicationService": "/Common/mytestpolicy.app/mytestpolicy",
			"gwtProfileReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/gwt-profiles?ver=16.1.0",
			"isSubCollection": true
			},
			"webhookReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/webhooks?ver=16.1.0",
			"isSubCollection": true
			},
			"whitelistIpReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/whitelist-ips?ver=16.1.0",
			"isSubCollection": true
			},
			"historyRevisionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/history-revisions?ver=16.1.0",
			"isSubCollection": true
			},
			"policyBuilderReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder?ver=16.1.0"
			},
			"responsePageReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/response-pages?ver=16.1.0",
			"isSubCollection": true
			},
			"vulnerabilityAssessmentReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/vulnerability-assessment?ver=16.1.0"
			},
			"serverTechnologyReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/server-technologies?ver=16.1.0",
			"isSubCollection": true
			},
			"cookieReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/cookies?ver=16.1.0",
			"isSubCollection": true
			},
			"blockingSettingReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/blocking-settings?ver=16.1.0",
			"isSubCollection": true
			},
			"hostNameReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/host-names?ver=16.1.0",
			"isSubCollection": true
			},
			"versionDeviceName": "ecosyshydbigip16.com",
			"selfLink": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg?ver=16.1.0",
			"threatCampaignSettingReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/threat-campaign-settings?ver=16.1.0"
			},
			"signatureReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signatures?ver=16.1.0",
			"isSubCollection": true
			},
			"policyBuilderRedirectionProtectionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-redirection-protection?ver=16.1.0"
			},
			"filetypeReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/filetypes?ver=16.1.0",
			"isSubCollection": true
			},
			"id": "BLDetCAoCevvoh-INfzTFg",
			"modifierName": "admin",
			"manualVirtualServers": [],
			"versionDatetime": "2022-12-28T12:40:28Z",
			"ssrfHostReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/ssrf-hosts?ver=16.1.0",
			"isSubCollection": true
			},
			"subPath": "/Common/mytestpolicy.app",
			"sessionTrackingStatusReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/session-tracking-statuses?ver=16.1.0",
			"isSubCollection": true
			},
			"active": true,
			"auditLogReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/audit-logs?ver=16.1.0",
			"isSubCollection": true
			},
			"disallowedGeolocationReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/disallowed-geolocations?ver=16.1.0",
			"isSubCollection": true
			},
			"redirectionProtectionDomainReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/redirection-protection-domains?ver=16.1.0",
			"isSubCollection": true
			},
			"applicationServiceManagedUpdatesOnly": false,
			"type": "security",
			"signatureSettingReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-settings?ver=16.1.0"
			},
			"deceptionResponsePageReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/deception-response-pages?ver=16.1.0",
			"isSubCollection": true
			},
			"websocketUrlReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/websocket-urls?ver=16.1.0",
			"isSubCollection": true
			},
			"xmlProfileReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/xml-profiles?ver=16.1.0",
			"isSubCollection": true
			},
			"methodReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/methods?ver=16.1.0",
			"isSubCollection": true
			},
			"vulnerabilityReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/vulnerabilities?ver=16.1.0",
			"isSubCollection": true
			},
			"redirectionProtectionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/redirection-protection?ver=16.1.0"
			},
			"policyBuilderSessionsAndLoginsReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-sessions-and-logins?ver=16.1.0"
			},
			"templateReference": {
			"link": "https://localhost/mgmt/tm/asm/policy-templates/clM9mdnuWVeuPOydwWNWnA?ver=16.1.0",
			"title": "OWA Exchange 2016"
			},
			"policyBuilderHeaderReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-header?ver=16.1.0"
			},
			"creatorName": "admin",
			"urlReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/urls?ver=16.1.0",
			"isSubCollection": true
			},
			"headerReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/headers?ver=16.1.0",
			"isSubCollection": true
			},
			"actionItemReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/action-items?ver=16.1.0",
			"isSubCollection": true
			},
			"microserviceReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/microservices?ver=16.1.0",
			"isSubCollection": true
			},
			"xmlValidationFileReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/xml-validation-files?ver=16.1.0",
			"isSubCollection": true
			},
			"lastUpdateMicros": 1.672231238e+15,
			"jsonProfileReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/json-profiles?ver=16.1.0",
			"isSubCollection": true
			},
			"bruteForceAttackPreventionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/brute-force-attack-preventions?ver=16.1.0",
			"isSubCollection": true
			},
			"disabledActionItemReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/disabled-action-items?ver=16.1.0",
			"isSubCollection": true
			},
			"jsonValidationFileReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/json-validation-files?ver=16.1.0",
			"isSubCollection": true
			},
			"extractionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/extractions?ver=16.1.0",
			"isSubCollection": true
			},
			"characterSetReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/character-sets?ver=16.1.0",
			"isSubCollection": true
			},
			"suggestionReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/suggestions?ver=16.1.0",
			"isSubCollection": true
			},
			"deceptionSettingsReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/deception-settings?ver=16.1.0"
			},
			"isModified": false,
			"sensitiveParameterReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/sensitive-parameters?ver=16.1.0",
			"isSubCollection": true
			},
			"generalReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/general?ver=16.1.0"
			},
			"versionPolicyName": "/Common/mytestpolicy.app/mytestpolicy",
			"policyBuilderCentralConfigurationReference": {
			"link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-central-configuration?ver=16.1.0"}
			}]}`)
		}
	})
	mux.HandleFunc("/mgmt/tm/asm/policies/?$filter=contains(name,'mytestpolicy')+and+contains(partition,'Common')", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "*" {
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, `{}`)
		}
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/apply-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:asm:tasks:apply-policy:apply-policy-taskstate",
    "selfLink": "https://localhost/mgmt/tm/asm/tasks/apply-policy/MgSRYOSqS5ohDSjGTns52A?ver=16.1.0",
    "status": "NEW",
    "id": "MgSRYOSqS5ohDSjGTns52A",
    "startTime": "2023-01-04T07:45:15Z",
    "lastUpdateMicros": 1.672818315e+15
}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/apply-policy/MgSRYOSqS5ohDSjGTns52A", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "executionStartTime": "2022-07-05T12:04:52Z",
    "status": "COMPLETED",
    "lastUpdateMicros": 1.657022692e+15,
    "kind": "tm:asm:tasks:apply-policy:apply-policy-taskstate",
    "selfLink": "https://localhost/mgmt/tm/asm/tasks/apply-policy/MgSRYOSqS5ohDSjGTns52A?ver=16.1.0",
    "policyReference": {
        "link": "https://localhost/mgmt/tm/asm/policies/LieFcG9wmFVfllNRyw144Q?ver=16.1.0",
        "fullPath": "/Common/mytestpolicy"
    },
    "endTime": "2022-07-05T12:04:53Z",
    "startTime": "2022-07-05T12:04:52Z",
    "id": "MgSRYOSqS5ohDSjGTns52A"
}`)
	})

	mux.HandleFunc("/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
  "signatureRequirementReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-requirements?ver=16.1.0",
    "isSubCollection": true
  },
  "plainTextProfileReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/plain-text-profiles?ver=16.1.0",
    "isSubCollection": true
  },
  "enablePassiveMode": false,
  "behavioralEnforcementReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/behavioral-enforcement?ver=16.1.0"
  },
  "dataGuardReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/data-guard?ver=16.1.0"
  },
  "createdDatetime": "2022-12-05T08:12:24Z",
  "databaseProtectionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/database-protection?ver=16.1.0"
  },
  "cookieSettingsReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/cookie-settings?ver=16.1.0"
  },
  "csrfUrlReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/csrf-urls?ver=16.1.0",
    "isSubCollection": true
  },
  "versionLastChange": " Security Policy /Common/mytestpolicy [add]: Parent Policy was set to empty value.\nType was set to Security.\nEncoding Selected was set to true.\nApplication Language was set to utf-8.\nActive was set to false.\nPolicy Name was set to /Common/mytestpolicy. { audit: policy = /Common/mytestpolicy, component = tsconfd }",
  "name": "mytestpolicy",
  "caseInsensitive": false,
  "headerSettingsReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/header-settings?ver=16.1.0"
  },
  "sectionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/sections?ver=16.1.0",
    "isSubCollection": true
  },
  "flowReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/flows?ver=16.1.0",
    "isSubCollection": true
  },
  "loginPageReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/login-pages?ver=16.1.0",
    "isSubCollection": true
  },
  "description": "Rapid Deployment-1",
  "fullPath": "/Common/mytestpolicy",
  "policyBuilderParameterReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-parameter?ver=16.1.0"
  },
  "hasParent": false,
  "threatCampaignReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/threat-campaigns?ver=16.1.0",
    "isSubCollection": true
  },
  "partition": "Common",
  "managedByBewaf": false,
  "csrfProtectionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/csrf-protection?ver=16.1.0"
  },
  "graphqlProfileReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/graphql-profiles?ver=16.1.0",
    "isSubCollection": true
  },
  "policyAntivirusReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/antivirus?ver=16.1.0"
  },
  "kind": "tm:asm:policies:policystate",
  "virtualServers": [],
  "policyBuilderCookieReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-cookie?ver=16.1.0"
  },
  "ipIntelligenceReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/ip-intelligence?ver=16.1.0"
  },
  "protocolIndependent": false,
  "sessionAwarenessSettingsReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/session-tracking?ver=16.1.0"
  },
  "policyBuilderUrlReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-url?ver=16.1.0"
  },
  "policyBuilderServerTechnologiesReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-server-technologies?ver=16.1.0"
  },
  "policyBuilderFiletypeReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-filetype?ver=16.1.0"
  },
  "signatureSetReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-sets?ver=16.1.0",
    "isSubCollection": true
  },
  "parameterReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/parameters?ver=16.1.0",
    "isSubCollection": true
  },
  "applicationLanguage": "utf-8",
  "enforcementMode": "transparent",
  "loginEnforcementReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/login-enforcement?ver=16.1.0"
  },
  "openApiFileReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/open-api-files?ver=16.1.0",
    "isSubCollection": true
  },
  "navigationParameterReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/navigation-parameters?ver=16.1.0",
    "isSubCollection": true
  },
  "gwtProfileReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/gwt-profiles?ver=16.1.0",
    "isSubCollection": true
  },
  "webhookReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/webhooks?ver=16.1.0",
    "isSubCollection": true
  },
  "whitelistIpReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/whitelist-ips?ver=16.1.0",
    "isSubCollection": true
  },
  "historyRevisionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/history-revisions?ver=16.1.0",
    "isSubCollection": true
  },
  "policyBuilderReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder?ver=16.1.0"
  },
  "responsePageReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/response-pages?ver=16.1.0",
    "isSubCollection": true
  },
  "vulnerabilityAssessmentReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/vulnerability-assessment?ver=16.1.0"
  },
  "serverTechnologyReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/server-technologies?ver=16.1.0",
    "isSubCollection": true
  },
  "cookieReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/cookies?ver=16.1.0",
    "isSubCollection": true
  },
  "blockingSettingReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/blocking-settings?ver=16.1.0",
    "isSubCollection": true
  },
  "hostNameReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/host-names?ver=16.1.0",
    "isSubCollection": true
  },
  "versionDeviceName": "ecosyshydbigip16.com",
  "selfLink": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg?ver=16.1.0",
  "threatCampaignSettingReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/threat-campaign-settings?ver=16.1.0"
  },
  "signatureReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signatures?ver=16.1.0",
    "isSubCollection": true
  },
  "policyBuilderRedirectionProtectionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-redirection-protection?ver=16.1.0"
  },
  "filetypeReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/filetypes?ver=16.1.0",
    "isSubCollection": true
  },
  "id": "BLDetCAoCevvoh-INfzTFg",
  "modifierName": "",
  "manualVirtualServers": [],
  "versionDatetime": "2022-12-05T08:12:29Z",
  "ssrfHostReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/ssrf-hosts?ver=16.1.0",
    "isSubCollection": true
  },
  "subPath": "/Common",
  "sessionTrackingStatusReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/session-tracking-statuses?ver=16.1.0",
    "isSubCollection": true
  },
  "active": false,
  "auditLogReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/audit-logs?ver=16.1.0",
    "isSubCollection": true
  },
  "disallowedGeolocationReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/disallowed-geolocations?ver=16.1.0",
    "isSubCollection": true
  },
  "redirectionProtectionDomainReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/redirection-protection-domains?ver=16.1.0",
    "isSubCollection": true
  },
  "type": "security",
  "signatureSettingReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/signature-settings?ver=16.1.0"
  },
  "deceptionResponsePageReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/deception-response-pages?ver=16.1.0",
    "isSubCollection": true
  },
  "websocketUrlReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/websocket-urls?ver=16.1.0",
    "isSubCollection": true
  },
  "xmlProfileReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/xml-profiles?ver=16.1.0",
    "isSubCollection": true
  },
  "methodReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/methods?ver=16.1.0",
    "isSubCollection": true
  },
  "vulnerabilityReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/vulnerabilities?ver=16.1.0",
    "isSubCollection": true
  },
  "redirectionProtectionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/redirection-protection?ver=16.1.0"
  },
  "policyBuilderSessionsAndLoginsReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-sessions-and-logins?ver=16.1.0"
  },
  "templateReference": {
    "link": "https://localhost/mgmt/tm/asm/policy-templates/EzpBNMs9gbVsF5uuiBjYDw?ver=16.1.0",
    "title": "Rapid Deployment Policy"
  },
  "policyBuilderHeaderReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-header?ver=16.1.0"
  },
  "creatorName": "tsconfd",
  "urlReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/urls?ver=16.1.0",
    "isSubCollection": true
  },
  "headerReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/headers?ver=16.1.0",
    "isSubCollection": true
  },
  "actionItemReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/action-items?ver=16.1.0",
    "isSubCollection": true
  },
  "microserviceReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/microservices?ver=16.1.0",
    "isSubCollection": true
  },
  "xmlValidationFileReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/xml-validation-files?ver=16.1.0",
    "isSubCollection": true
  },
  "lastUpdateMicros": 0,
  "jsonProfileReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/json-profiles?ver=16.1.0",
    "isSubCollection": true
  },
  "bruteForceAttackPreventionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/brute-force-attack-preventions?ver=16.1.0",
    "isSubCollection": true
  },
  "disabledActionItemReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/disabled-action-items?ver=16.1.0",
    "isSubCollection": true
  },
  "jsonValidationFileReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/json-validation-files?ver=16.1.0",
    "isSubCollection": true
  },
  "extractionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/extractions?ver=16.1.0",
    "isSubCollection": true
  },
  "characterSetReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/character-sets?ver=16.1.0",
    "isSubCollection": true
  },
  "suggestionReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/suggestions?ver=16.1.0",
    "isSubCollection": true
  },
  "deceptionSettingsReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/deception-settings?ver=16.1.0"
  },
  "isModified": false,
  "sensitiveParameterReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/sensitive-parameters?ver=16.1.0",
    "isSubCollection": true
  },
  "generalReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/general?ver=16.1.0"
  },
  "versionPolicyName": "/Common/mytestpolicy",
  "policyBuilderCentralConfigurationReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/BLDetCAoCevvoh-INfzTFg/policy-builder-central-configuration?ver=16.1.0"
  }
}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
  "isBase64": false,
  "inline": true,
  "minimal": false,
  "status": "NEW",
  "lastUpdateMicros": 1.672819199e+15,
  "exportSuggestions": false,
  "includeVulnerabilityAssessmentConfigurationAndData": true,
  "kind": "tm:asm:tasks:export-policy:export-policy-taskstate",
  "selfLink": "https://localhost/mgmt/tm/asm/tasks/export-policy/uWHzxaTZGWL0X3ft7lYPEQ?ver=16.1.0",
  "format": "json",
  "policyReference": {
    "link": "https://localhost/mgmt/tm/asm/policies/LieFcG9wmFVfllNRyw144Q?ver=16.1.0",
    "fullPath": "/Common/mytestpolicy"
  },
  "id": "uWHzxaTZGWL0X3ft7lYPEQ",
  "startTime": "2023-01-04T07:59:59Z"
}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy/uWHzxaTZGWL0X3ft7lYPEQ", func(w http.ResponseWriter, r *http.Request) {
		filedata := string([]byte(`{\n   \"policy\" : {\n      \"antivirus\" : {\n         \"inspectHttpUploads\" : false\n      },\n      \"applicationLanguage\" : \"utf-8\",\n      \"behavioral-enforcement\" : {\n         \"behavioralEnforcementViolations\" : [\n            {\n               \"name\" : \"VIOL_CONVICTION\"\n            },\n            {\n               \"name\" : \"VIOL_THREAT_ANALYSIS\"\n            },\n            {\n               \"name\" : \"VIOL_BLOCKING_CONDITION\"\n            },\n            {\n               \"name\" : \"VIOL_THREAT_CAMPAIGN\"\n            },\n            {\n               \"name\" : \"VIOL_BLACKLISTED_IP\"\n            },\n            {\n               \"name\" : \"VIOL_GEOLOCATION\"\n            }\n         ],\n         \"enableBehavioralEnforcement\" : false,\n         \"enableBlockingCveSignatures\" : true,\n         \"enableBlockingHighAccuracySignatures\" : true,\n         \"enableBlockingLikelyMaliciousTransactions\" : true,\n         \"enableBlockingSuspiciousTransactions\" : false,\n         \"enableBlockingViolations\" : true\n      },\n      \"blocking-settings\" : {\n         \"evasions\" : [\n            {\n               \"description\" : \"Bad unescape\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Apache whitespace\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Bare byte decoding\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"IIS Unicode codepoints\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"IIS backslashes\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"%u decoding\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Multiple decoding\",\n               \"enabled\" : false,\n               \"learn\" : true,\n               \"maxDecodingPasses\" : 3\n            },\n            {\n               \"description\" : \"Directory traversals\",\n               \"enabled\" : false,\n               \"learn\" : true\n            }\n         ],\n         \"http-protocols\" : [\n            {\n               \"description\" : \"Multiple host headers\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Check maximum number of parameters\",\n               \"enabled\" : false,\n               \"learn\" : true,\n               \"maxParams\" : 500\n            },\n            {\n               \"description\" : \"Bad host header value\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Check maximum number of headers\",\n               \"enabled\" : false,\n               \"learn\" : true,\n               \"maxHeaders\" : 20\n            },\n            {\n               \"description\" : \"Unparsable request content\",\n               \"enabled\" : true\n            },\n            {\n               \"description\" : \"High ASCII characters in headers\",\n               \"enabled\" : false,\n               \"learn\" : false\n            },\n            {\n               \"description\" : \"Null in request\",\n               \"enabled\" : true\n            },\n            {\n               \"description\" : \"Bad HTTP version\",\n               \"enabled\" : true\n            },\n            {\n               \"description\" : \"Content length should be a positive number\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Host header contains IP address\",\n               \"enabled\" : false,\n               \"learn\" : false\n            },\n            {\n               \"description\" : \"CRLF characters before request start\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"No Host header in HTTP/1.1 request\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Bad multipart parameters parsing\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Bad multipart/form-data request parsing\",\n               \"enabled\" : false,\n               \"learn\" : false\n            },\n            {\n               \"description\" : \"Body in GET or HEAD requests\",\n               \"enabled\" : false,\n               \"learn\" : false\n            },\n            {\n               \"description\" : \"Chunked request with Content-Length header\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Several Content-Length headers\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Header name with no header value\",\n               \"enabled\" : false,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"POST request with Content-Length: 0\",\n               \"enabled\" : false,\n               \"learn\" : false\n            }\n         ],\n         \"violations\" : [\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Request length exceeds defined buffer size\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_REQUEST_MAX_LENGTH\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal parameter value length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_VALUE_LENGTH\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal URL\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_URL\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal session ID in URL\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_DYNAMIC_SESSION\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Login URL bypassed\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_LOGIN_URL_BYPASSED\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"GraphQL data does not comply with format settings\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_GRAPHQL_FORMAT\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal file type\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_FILETYPE\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Web Services Security failure\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_XML_WEB_SERVICES_SECURITY\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"ASM Cookie Hijacking\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_ASM_COOKIE_HIJACKING\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Cookie not RFC-compliant\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_COOKIE_MALFORMED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"XML data does not comply with format settings\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_XML_FORMAT\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"DataSafe Data Integrity\",\n               \"name\" : \"VIOL_DATA_INTEGRITY\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Login URL expired\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_LOGIN_URL_EXPIRED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal flow to URL\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_FLOW\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Blocking Condition Detected\",\n               \"name\" : \"VIOL_BLOCKING_CONDITION\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Virus detected\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_VIRUS\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Server-side access to disallowed host\",\n               \"name\" : \"VIOL_SERVER_SIDE_HOST\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Failed to convert character\",\n               \"name\" : \"VIOL_ENCODING\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal meta character in value\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_VALUE_METACHAR\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Expired timestamp\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_COOKIE_EXPIRED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"SOAP method not allowed\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_XML_SOAP_METHOD\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal parameter data type\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_DATA_TYPE\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Threat Campaign detected\",\n               \"name\" : \"VIOL_THREAT_CAMPAIGN\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Modified ASM cookie\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_ASM_COOKIE_MODIFIED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal meta character in parameter name\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_NAME_METACHAR\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal entry point\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_FLOW_ENTRY_POINT\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Access from malicious IP address\",\n               \"name\" : \"VIOL_MALICIOUS_IP\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal cross-origin request\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_CROSS_ORIGIN_REQUEST\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"HTTP protocol compliance failed\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_HTTP_PROTOCOL\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal number of mandatory parameters\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_FLOW_MANDATORY_PARAMS\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Plain text data does not comply with format settings\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PLAINTEXT_FORMAT\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Illegal parameter location\",\n               \"name\" : \"VIOL_PARAMETER_LOCATION\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal header length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_HEADER_LENGTH\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal meta character in URL\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_URL_METACHAR\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Mandatory parameter is missing\",\n               \"name\" : \"VIOL_MANDATORY_PARAMETER\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal POST data length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_POST_DATA_LENGTH\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal host name\",\n               \"name\" : \"VIOL_HOSTNAME\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal attachment in SOAP message\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_XML_SOAP_ATTACHMENT\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal request length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_REQUEST_LENGTH\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Data Guard: Information leakage detected\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_DATA_GUARD\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Null character found in WebSocket text message\",\n               \"name\" : \"VIOL_WEBSOCKET_TEXT_NULL_VALUE\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Text content found in binary only WebSocket\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_WEBSOCKET_TEXT_MESSAGE_NOT_ALLOWED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal WebSocket frame length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_WEBSOCKET_FRAME_LENGTH\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Illegal HTTP status in response\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_HTTP_RESPONSE_STATUS\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Mandatory HTTP header is missing\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_MANDATORY_HEADER\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal URL length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_URL_LENGTH\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"IP is blacklisted\",\n               \"name\" : \"VIOL_BLACKLISTED_IP\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Mitigation action determined by Threat Analysis Platform\",\n               \"name\" : \"VIOL_THREAT_ANALYSIS\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal static parameter value\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_STATIC_VALUE\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Mask not found in client frame\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_WEBSOCKET_FRAME_MASKING\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal number of frames per message\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_WEBSOCKET_FRAMES_PER_MESSAGE_COUNT\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal meta character in header\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_HEADER_METACHAR\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Illegal Base64 value\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_PARAMETER_VALUE_BASE64\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Modified domain cookie(s)\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_COOKIE_MODIFIED\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"JSON data does not comply with JSON schema\",\n               \"name\" : \"VIOL_JSON_SCHEMA\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Access from disallowed Geolocation\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_GEOLOCATION\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"CSRF attack detected\",\n               \"name\" : \"VIOL_CSRF\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal parameter numeric value\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_NUMERIC_VALUE\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Disallowed file upload content detected\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_FILE_UPLOAD\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Malformed GraphQL data\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_GRAPHQL_MALFORMED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Bad WebSocket handshake request\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_WEBSOCKET_BAD_REQUEST\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Parameter value does not comply with regular expression\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_VALUE_REGEXP\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Malformed JSON data\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_JSON_MALFORMED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Null in multi-part parameter value\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_MULTIPART_NULL_VALUE\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Failure in WebSocket framing protocol\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_WEBSOCKET_FRAMING_PROTOCOL\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Host name mismatch\",\n               \"name\" : \"VIOL_HOSTNAME_MISMATCH\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Violation Rating Need Examination detected\",\n               \"name\" : \"VIOL_RATING_NEED_EXAMINATION\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Leaked Credentials Detection\",\n               \"name\" : \"VIOL_LEAKED_CREDENTIALS\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Access from disallowed User/Session/IP/Device ID\",\n               \"name\" : \"VIOL_SESSION_AWARENESS\"\n            },\n            {\n               \"description\" : \"Attack signature detected\",\n               \"name\" : \"VIOL_ATTACK_SIGNATURE\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Illegal repeated header\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_HEADER_REPEATED\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Malformed XML data\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_XML_MALFORMED\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Illegal method\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_METHOD\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal empty parameter value\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_EMPTY_VALUE\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"XML data does not comply with schema or WSDL document\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_XML_SCHEMA\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Bad Actor Convicted\",\n               \"name\" : \"VIOL_CONVICTION\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Violation Rating Threat detected\",\n               \"name\" : \"VIOL_RATING_THREAT\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"CSRF authentication expired\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_CSRF_EXPIRED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"GWT data does not comply with format settings\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_GWT_FORMAT\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Malformed GWT data\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_GWT_MALFORMED\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Disallowed file upload content detected in body\",\n               \"name\" : \"VIOL_FILE_UPLOAD_IN_BODY\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"JSON data does not comply with format settings\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_JSON_FORMAT\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal WebSocket binary message length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_WEBSOCKET_BINARY_MESSAGE_LENGTH\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Bad Actor Detected\",\n               \"name\" : \"VIOL_MALICIOUS_DEVICE\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal request content type\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_URL_CONTENT_TYPE\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal dynamic parameter value\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_DYNAMIC_VALUE\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal query string length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_QUERY_STRING_LENGTH\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal WebSocket extension\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_WEBSOCKET_EXTENSION\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal cookie length\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_COOKIE_LENGTH\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"GraphQL introspection query\",\n               \"name\" : \"VIOL_GRAPHQL_INTROSPECTION_QUERY\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal repeated parameter name\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER_REPEATED\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal parameter array value\",\n               \"name\" : \"VIOL_PARAMETER_ARRAY_VALUE\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Evasion technique detected\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_EVASION\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Binary content found in text only WebSocket\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_WEBSOCKET_BINARY_MESSAGE_NOT_ALLOWED\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Brute Force: Maximum login attempts are exceeded\",\n               \"name\" : \"VIOL_BRUTE_FORCE\"\n            },\n            {\n               \"alarm\" : true,\n               \"block\" : true,\n               \"description\" : \"Illegal redirection attempt\",\n               \"learn\" : true,\n               \"name\" : \"VIOL_REDIRECT\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Mandatory request body is missing\",\n               \"name\" : \"VIOL_MANDATORY_REQUEST_BODY\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal parameter\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_PARAMETER\"\n            },\n            {\n               \"alarm\" : false,\n               \"block\" : false,\n               \"description\" : \"Illegal query string or POST data\",\n               \"learn\" : false,\n               \"name\" : \"VIOL_FLOW_DISALLOWED_INPUT\"\n            }\n         ],\n         \"web-services-securities\" : [\n            {\n               \"description\" : \"UnSigned Timestamp\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Timestamp expiration is too far in the future\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Expired Timestamp\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Invalid Timestamp\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Missing Timestamp\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Verification Error\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Signing Error\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Encryption Error\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Decryption Error\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Certificate Error\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Certificate Expired\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Malformed Error\",\n               \"enabled\" : true,\n               \"learn\" : true\n            },\n            {\n               \"description\" : \"Internal Error\",\n               \"enabled\" : true,\n               \"learn\" : true\n            }\n         ]\n      },\n      \"brute-force-attack-preventions\" : [\n         {\n            \"bruteForceProtectionForAllLoginPages\" : false,\n            \"captchaBypassCriteria\" : {\n               \"action\" : \"alarm-and-drop\",\n               \"enabled\" : true,\n               \"threshold\" : 5\n            },\n            \"clientSideIntegrityBypassCriteria\" : {\n               \"action\" : \"alarm-and-captcha\",\n               \"enabled\" : true,\n               \"threshold\" : 3\n            },\n            \"detectionCriteria\" : {\n               \"action\" : \"alarm-and-captcha\",\n               \"credentialsStuffingMatchesReached\" : 100,\n               \"detectCredentialsStuffingAttack\" : true,\n               \"detectDistributedBruteForceAttack\" : true,\n               \"failedLoginAttemptsRateReached\" : 100\n            },\n            \"leakedCredentialsCriteria\" : {\n               \"action\" : \"alarm-and-blocking-page\",\n               \"enabled\" : false\n            },\n            \"loginAttemptsFromTheSameDeviceId\" : {\n               \"action\" : \"alarm-and-captcha\",\n               \"enabled\" : false,\n               \"threshold\" : 3\n            },\n            \"loginAttemptsFromTheSameIp\" : {\n               \"action\" : \"alarm-and-captcha\",\n               \"enabled\" : true,\n               \"threshold\" : 20\n            },\n            \"loginAttemptsFromTheSameUser\" : {\n               \"action\" : \"alarm-and-captcha\",\n               \"enabled\" : true,\n               \"threshold\" : 3\n            },\n            \"measurementPeriod\" : 900,\n            \"preventionDuration\" : \"3600\",\n            \"reEnableLoginAfter\" : 3600,\n            \"sourceBasedProtectionDetectionPeriod\" : 3600\n         }\n      ],\n      \"caseInsensitive\" : false,\n      \"character-sets\" : [\n         {\n            \"characterSet\" : [\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x0\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x1\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x2\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x3\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x4\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x5\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x6\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x7\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x8\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x9\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0xa\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0xb\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0xc\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0xd\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0xe\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0xf\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x10\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x11\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x12\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x13\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x14\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x15\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x16\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x17\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x18\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x19\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x1a\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x1b\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x1c\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x1d\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x1e\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x1f\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x20\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x21\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x22\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x23\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x24\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x25\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x26\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x27\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x28\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x29\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x2a\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x2b\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x2c\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x2d\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x2e\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x2f\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x30\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x31\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x32\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x33\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x34\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x35\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x36\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x37\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x38\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x39\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x3a\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x3b\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x3c\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x3d\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x3e\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x3f\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x40\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x41\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x42\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x43\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x44\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x45\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x46\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x47\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x48\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x49\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x4a\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x4b\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x4c\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x4d\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x4e\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x4f\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x50\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x51\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x52\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x53\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x54\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x55\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x56\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x57\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x58\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x59\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x5a\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x5b\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x5c\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x5d\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x5e\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x5f\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x60\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x61\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x62\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x63\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x64\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x65\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x66\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x67\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x68\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x69\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x6a\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x6b\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x6c\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x6d\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x6e\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x6f\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x70\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x71\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x72\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x73\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x74\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x75\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x76\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x77\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x78\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x79\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x7a\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x7b\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x7c\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x7d\"\n               },\n               {\n                  \"isAllowed\" : true,\n                  \"metachar\" : \"0x7e\"\n               },\n               {\n                  \"isAllowed\" : false,\n                  \"metachar\" : \"0x7f\"\n               }\n            ],\n            \"characterSetType\" : \"plain-text-content\"\n         }\n      ],\n      \"cookie-settings\" : {\n         \"maximumCookieHeaderLength\" : \"any\"\n      },\n      \"cookies\" : [\n         {\n            \"accessibleOnlyThroughTheHttpProtocol\" : false,\n            \"attackSignaturesCheck\" : true,\n            \"enforcementType\" : \"allow\",\n            \"insertSameSiteAttribute\" : \"none\",\n            \"isBase64\" : false,\n            \"maskValueInLogs\" : false,\n            \"name\" : \"*\",\n            \"performStaging\" : true,\n            \"securedOverHttpsConnection\" : false,\n            \"type\" : \"wildcard\",\n            \"wildcardOrder\" : 1\n         }\n      ],\n      \"csrf-protection\" : {\n         \"enabled\" : false\n      },\n      \"csrf-urls\" : [\n         {\n            \"enforcementAction\" : \"verify-csrf-token\",\n            \"method\" : \"POST\",\n            \"requiredParameters\" : \"ignore\",\n            \"url\" : \"*\",\n            \"wildcardOrder\" : 1\n         }\n      ],\n      \"data-guard\" : {\n         \"enabled\" : false,\n         \"enforcementMode\" : \"ignore-urls-in-list\"\n      },\n      \"database-protection\" : {\n         \"databaseProtectionEnabled\" : false,\n         \"userSource\" : \"apm\"\n      },\n      \"deception-settings\" : {\n         \"enableCustomResponses\" : false,\n         \"enableResponsePageByAttackType\" : true,\n         \"serverTechnologyName\" : \"Nginx\"\n      },\n      \"description\" : \"Rapid Deployment-1\",\n      \"enablePassiveMode\" : false,\n      \"enforcementMode\" : \"transparent\",\n      \"filetypes\" : [\n         {\n            \"allowed\" : false,\n            \"name\" : \"bin\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"exe\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"shtml\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"pol\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"sys\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"p7c\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"bck\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"printer\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"cfg\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"dat\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"tmp\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"conf\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"save\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"cgi\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"bat\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"pfx\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"reg\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"bkp\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"msi\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"config\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"ida\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"temp\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"htw\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"shtm\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"wmz\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"com\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"pem\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"old\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"cmd\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"idc\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"crt\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"hta\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"sav\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"log\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"idq\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"key\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"stm\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"der\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"nws\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"eml\"\n         },\n         {\n            \"allowed\" : true,\n            \"checkPostDataLength\" : false,\n            \"checkQueryStringLength\" : false,\n            \"checkRequestLength\" : false,\n            \"checkUrlLength\" : false,\n            \"name\" : \"*\",\n            \"performStaging\" : false,\n            \"responseCheck\" : false,\n            \"type\" : \"wildcard\",\n            \"wildcardOrder\" : 1\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"dll\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"bak\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"ini\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"cer\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"htr\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"p12\"\n         },\n         {\n            \"allowed\" : false,\n            \"name\" : \"p7b\"\n         }\n      ],\n      \"fullPath\" : \"/Common/mytestpolicy\",\n      \"general\" : {\n         \"allowedResponseCodes\" : [\n            400,\n            401,\n            404,\n            407,\n            417,\n            503\n         ],\n         \"enableEventCorrelation\" : true,\n         \"enforcementReadinessPeriod\" : 7,\n         \"maskCreditCardNumbersInRequest\" : true,\n         \"pathParameterHandling\" : \"as-parameters\",\n         \"triggerAsmIruleEvent\" : \"disabled\",\n         \"trustXff\" : false,\n         \"useDynamicSessionIdInUrl\" : false\n      },\n      \"graphql-profiles\" : [\n         {\n            \"attackSignaturesCheck\" : true,\n            \"defenseAttributes\" : {\n               \"allowIntrospectionQueries\" : false,\n               \"maximumBatchedQueries\" : \"any\",\n               \"maximumStructureDepth\" : \"any\",\n               \"maximumTotalLength\" : \"any\",\n               \"maximumValueLength\" : \"any\",\n               \"tolerateParsingWarnings\" : true\n            },\n            \"description\" : \"Default GraphQL Profile\",\n            \"metacharElementCheck\" : true,\n            \"name\" : \"Default\"\n         }\n      ],\n      \"gwt-profiles\" : [\n         {\n            \"attackSignaturesCheck\" : true,\n            \"defenseAttributes\" : {\n               \"maximumTotalLengthOfGWTData\" : 10000,\n               \"maximumValueLength\" : 100,\n               \"tolerateGWTParsingWarnings\" : true\n            },\n            \"description\" : \"Default GWT Profile\",\n            \"metacharElementCheck\" : true,\n            \"name\" : \"Default\"\n         }\n      ],\n      \"header-settings\" : {\n         \"maximumHttpHeaderLength\" : \"any\"\n      },\n      \"headers\" : [\n         {\n            \"allowRepeatedOccurrences\" : false,\n            \"base64Decoding\" : false,\n            \"checkSignatures\" : true,\n            \"htmlNormalization\" : false,\n            \"mandatory\" : false,\n            \"maskValueInLogs\" : false,\n            \"name\" : \"transfer-encoding\",\n            \"normalizationViolations\" : false,\n            \"percentDecoding\" : false,\n            \"type\" : \"explicit\",\n            \"urlNormalization\" : false\n         },\n         {\n            \"allowRepeatedOccurrences\" : true,\n            \"base64Decoding\" : false,\n            \"checkSignatures\" : true,\n            \"htmlNormalization\" : false,\n            \"mandatory\" : false,\n            \"maskValueInLogs\" : true,\n            \"name\" : \"authorization\",\n            \"normalizationViolations\" : false,\n            \"percentDecoding\" : true,\n            \"type\" : \"explicit\",\n            \"urlNormalization\" : false\n         },\n         {\n            \"allowRepeatedOccurrences\" : true,\n            \"checkSignatures\" : false,\n            \"mandatory\" : false,\n            \"maskValueInLogs\" : false,\n            \"name\" : \"cookie\",\n            \"type\" : \"explicit\"\n         },\n         {\n            \"allowRepeatedOccurrences\" : true,\n            \"base64Decoding\" : false,\n            \"checkSignatures\" : true,\n            \"htmlNormalization\" : false,\n            \"mandatory\" : false,\n            \"maskValueInLogs\" : false,\n            \"name\" : \"*\",\n            \"normalizationViolations\" : false,\n            \"percentDecoding\" : true,\n            \"type\" : \"wildcard\",\n            \"urlNormalization\" : false,\n            \"wildcardOrder\" : 1\n         },\n         {\n            \"allowRepeatedOccurrences\" : true,\n            \"base64Decoding\" : false,\n            \"checkSignatures\" : true,\n            \"htmlNormalization\" : false,\n            \"mandatory\" : false,\n            \"maskValueInLogs\" : false,\n            \"name\" : \"referer\",\n            \"normalizationViolations\" : true,\n            \"percentDecoding\" : false,\n            \"type\" : \"explicit\",\n            \"urlNormalization\" : true\n         }\n      ],\n      \"ip-intelligence\" : {\n         \"enabled\" : false\n      },\n      \"json-profiles\" : [\n         {\n            \"defenseAttributes\" : {\n               \"maximumArrayLength\" : \"any\",\n               \"maximumStructureDepth\" : \"any\",\n               \"maximumTotalLengthOfJSONData\" : \"any\",\n               \"maximumValueLength\" : \"any\",\n               \"tolerateJSONParsingWarnings\" : true\n            },\n            \"description\" : \"Default JSON Profile\",\n            \"handleJsonValuesAsParameters\" : true,\n            \"hasValidationFiles\" : false,\n            \"name\" : \"Default\",\n            \"validationFiles\" : []\n         }\n      ],\n      \"login-enforcement\" : {\n         \"expirationTimePeriod\" : \"disabled\"\n      },\n      \"methods\" : [\n         {\n            \"actAsMethod\" : \"POST\",\n            \"name\" : \"POST\"\n         },\n         {\n            \"actAsMethod\" : \"GET\",\n            \"name\" : \"HEAD\"\n         },\n         {\n            \"actAsMethod\" : \"GET\",\n            \"name\" : \"GET\"\n         }\n      ],\n      \"name\" : \"mytestpolicy\",\n      \"parameters\" : [\n         {\n            \"allowEmptyValue\" : true,\n            \"allowRepeatedParameterName\" : true,\n            \"attackSignaturesCheck\" : true,\n            \"checkMaxValueLength\" : false,\n            \"checkMetachars\" : false,\n            \"isBase64\" : false,\n            \"isCookie\" : false,\n            \"isHeader\" : false,\n            \"level\" : \"global\",\n            \"metacharsOnParameterValueCheck\" : false,\n            \"name\" : \"*\",\n            \"parameterLocation\" : \"any\",\n            \"performStaging\" : false,\n            \"sensitiveParameter\" : false,\n            \"type\" : \"wildcard\",\n            \"valueType\" : \"auto-detect\",\n            \"wildcardOrder\" : 1\n         },\n         {\n            \"allowEmptyValue\" : true,\n            \"allowRepeatedParameterName\" : false,\n            \"isCookie\" : false,\n            \"isHeader\" : false,\n            \"level\" : \"global\",\n            \"mandatory\" : false,\n            \"name\" : \"__VIEWSTATE\",\n            \"parameterLocation\" : \"any\",\n            \"performStaging\" : false,\n            \"sensitiveParameter\" : false,\n            \"type\" : \"explicit\",\n            \"valueType\" : \"ignore\"\n         }\n      ],\n      \"plain-text-profiles\" : [\n         {\n            \"attackSignaturesCheck\" : true,\n            \"defenseAttributes\" : {\n               \"maximumLineLength\" : \"any\",\n               \"maximumTotalLength\" : \"any\",\n               \"performPercentDecoding\" : false\n            },\n            \"description\" : \"Default Plain Text Profile\",\n            \"metacharElementCheck\" : false,\n            \"name\" : \"Default\"\n         }\n      ],\n      \"policy-builder\" : {\n         \"enableFullPolicyInspection\" : true,\n         \"enableTrustedTrafficSiteChangeTracking\" : true,\n         \"enableUntrustedTrafficSiteChangeTracking\" : true,\n         \"inactiveEntityInactivityDurationInDays\" : 90,\n         \"learnFromResponses\" : false,\n         \"learnInactiveEntities\" : true,\n         \"learnOnlyFromNonBotTraffic\" : true,\n         \"learningMode\" : \"manual\",\n         \"responseStatusCodes\" : [\n            \"1xx\",\n            \"3xx\",\n            \"2xx\"\n         ],\n         \"trafficTighten\" : {\n            \"maxModificationSuggestionScore\" : 50,\n            \"minDaysBetweenSamples\" : 1,\n            \"totalRequests\" : 15000\n         },\n         \"trustAllIps\" : false,\n         \"trustedTrafficLoosen\" : {\n            \"differentSources\" : 1,\n            \"maxDaysBetweenSamples\" : 7,\n            \"minHoursBetweenSamples\" : 0\n         },\n         \"trustedTrafficSiteChangeTracking\" : {\n            \"differentSources\" : 1,\n            \"maxDaysBetweenSamples\" : 7,\n            \"minMinutesBetweenSamples\" : 0\n         },\n         \"untrustedTrafficLoosen\" : {\n            \"differentSources\" : 20,\n            \"maxDaysBetweenSamples\" : 7,\n            \"minHoursBetweenSamples\" : 1\n         },\n         \"untrustedTrafficSiteChangeTracking\" : {\n            \"differentSources\" : 10,\n            \"maxDaysBetweenSamples\" : 7,\n            \"minMinutesBetweenSamples\" : 20\n         }\n      },\n      \"policy-builder-central-configuration\" : {\n         \"buildingMode\" : \"local\",\n         \"eventCorrelationMode\" : \"local\"\n      },\n      \"policy-builder-cookie\" : {\n         \"collapseCookieOccurrences\" : 10,\n         \"collapseCookiesIntoOneEntity\" : true,\n         \"enforceUnmodifiedCookies\" : true,\n         \"learnExplicitCookies\" : \"selective\",\n         \"maximumCookies\" : 100\n      },\n      \"policy-builder-filetype\" : {\n         \"learnExplicitFiletypes\" : \"never\",\n         \"maximumFileTypes\" : 100\n      },\n      \"policy-builder-header\" : {\n         \"maximumHosts\" : 10000,\n         \"validHostNames\" : false\n      },\n      \"policy-builder-parameter\" : {\n         \"classifyParameters\" : false,\n         \"collapseParametersIntoOneEntity\" : false,\n         \"dynamicParameters\" : {\n            \"allHiddenFields\" : false,\n            \"formParameters\" : false,\n            \"linkParameters\" : false,\n            \"uniqueValueSets\" : 10\n         },\n         \"learnExplicitParameters\" : \"never\",\n         \"maximumParameters\" : 10000,\n         \"parameterLearningLevel\" : \"global\",\n         \"parametersIntegerValue\" : false\n      },\n      \"policy-builder-redirection-protection\" : {\n         \"learnExplicitRedirectionDomains\" : \"never\",\n         \"maximumRedirectionDomains\" : 100\n      },\n      \"policy-builder-server-technologies\" : {\n         \"enableServerTechnologiesDetection\" : false\n      },\n      \"policy-builder-sessions-and-logins\" : {\n         \"learnLoginPage\" : false\n      },\n      \"policy-builder-url\" : {\n         \"classifyUrls\" : false,\n         \"classifyWebsocketUrls\" : false,\n         \"collapseUrlDepth\" : 2,\n         \"collapseUrlOccurrences\" : 500,\n         \"collapseUrlsIntoOneEntity\" : true,\n         \"learnExplicitUrls\" : \"never\",\n         \"learnExplicitWebsocketUrls\" : \"never\",\n         \"learnMethodsOnUrls\" : false,\n         \"maximumUrls\" : 10000,\n         \"maximumWebsocketUrls\" : 100,\n         \"wildcardUrlFiletypes\" : [\n            \"pcx\",\n            \"ico\",\n            \"pdf\",\n            \"swf\",\n            \"gif\",\n            \"png\",\n            \"jpeg\",\n            \"bmp\",\n            \"wav\",\n            \"jpg\"\n         ]\n      },\n      \"protocolIndependent\" : false,\n      \"redirection-protection\" : {\n         \"redirectionProtectionEnabled\" : false\n      },\n      \"redirection-protection-domains\" : [\n         {\n            \"domainName\" : \"*\",\n            \"type\" : \"wildcard\",\n            \"wildcardOrder\" : 1\n         }\n      ],\n      \"response-pages\" : [\n         {\n            \"responseActionType\" : \"erase-cookies\",\n            \"responsePageType\" : \"hijack\"\n         },\n         {\n            \"responseActionType\" : \"default\",\n            \"responsePageType\" : \"failed-login-honeypot\"\n         },\n         {\n            \"ajaxActionType\" : \"alert-popup\",\n            \"ajaxPopupMessage\" : \"These username and password were found in the Leaked Credentials Data Base. They maybe used by attackers to compromise your account. Please change password.\",\n            \"responsePageType\" : \"leaked-credentials-ajax\"\n         },\n         {\n            \"responseActionType\" : \"default\",\n            \"responsePageType\" : \"default\"\n         },\n         {\n            \"ajaxActionType\" : \"alert-popup\",\n            \"ajaxPopupMessage\" : \"The requested URL was rejected. Please consult with your administrator. Your support ID is: <%TS.request.ID()%>\",\n            \"responsePageType\" : \"ajax-login\"\n         },\n         {\n            \"responseActionType\" : \"default\",\n            \"responsePageType\" : \"captcha\"\n         },\n         {\n            \"responseActionType\" : \"default\",\n            \"responsePageType\" : \"graphql\"\n         },\n         {\n            \"responseActionType\" : \"default\",\n            \"responsePageType\" : \"leaked-credentials\"\n         },\n         {\n            \"responseActionType\" : \"soap-fault\",\n            \"responsePageType\" : \"xml\"\n         },\n         {\n            \"responseActionType\" : \"default\",\n            \"responsePageType\" : \"captcha-fail\"\n         },\n         {\n            \"responseActionType\" : \"default\",\n            \"responsePageType\" : \"mobile\"\n         },\n         {\n            \"ajaxActionType\" : \"alert-popup\",\n            \"ajaxEnabled\" : false,\n            \"ajaxPopupMessage\" : \"The requested URL was rejected. Please consult with your administrator. Your support ID is: <%TS.request.ID()%>\",\n            \"responsePageType\" : \"ajax\"\n         },\n         {\n            \"ajaxActionType\" : \"alert-popup\",\n            \"ajaxPopupMessage\" : \"Login Failed. Username or password is incorrect. Please try to log in again.\",\n            \"responsePageType\" : \"failed-login-honeypot-ajax\"\n         },\n         {\n            \"responseActionType\" : \"default\",\n            \"responsePageType\" : \"persistent-flow\"\n         }\n      ],\n      \"sensitive-parameters\" : [\n         {\n            \"name\" : \"password\"\n         }\n      ],\n      \"session-tracking\" : {\n         \"delayBlocking\" : {},\n         \"sessionTrackingConfiguration\" : {\n            \"enableSessionAwareness\" : false,\n            \"enableTrackingSessionHijackingByDeviceId\" : false\n         }\n      },\n      \"signature-sets\" : [\n         {\n            \"alarm\" : true,\n            \"block\" : true,\n            \"learn\" : true,\n            \"name\" : \"Generic Detection Signatures (High/Medium Accuracy)\"\n         }\n      ],\n      \"signature-settings\" : {\n         \"attackSignatureFalsePositiveMode\" : \"disabled\",\n         \"minimumAccuracyForAutoAddedSignatures\" : \"medium\",\n         \"placeSignaturesInStaging\" : true,\n         \"signatureStaging\" : true\n      },\n      \"softwareVersion\" : \"16.1.0\",\n      \"template\" : {\n         \"name\" : \"POLICY_TEMPLATE_RAPID_DEPLOYMENT\"\n      },\n      \"threat-campaign-settings\" : {\n         \"threatCampaignEnforcementReadinessPeriod\" : 1,\n         \"threatCampaignStaging\" : false\n      },\n      \"type\" : \"security\",\n      \"urls\" : [\n         {\n            \"attackSignaturesCheck\" : true,\n            \"clickjackingProtection\" : false,\n            \"description\" : \"\",\n            \"disallowFileUploadOfExecutables\" : false,\n            \"html5CrossOriginRequestsEnforcement\" : {\n               \"enforcementMode\" : \"disabled\"\n            },\n            \"isAllowed\" : true,\n            \"mandatoryBody\" : false,\n            \"metacharsOnUrlCheck\" : false,\n            \"method\" : \"*\",\n            \"methodsOverrideOnUrlCheck\" : false,\n            \"name\" : \"*\",\n            \"performStaging\" : false,\n            \"protocol\" : \"http\",\n            \"type\" : \"wildcard\",\n            \"urlContentProfiles\" : [\n               {\n                  \"headerName\" : \"*\",\n                  \"headerOrder\" : \"default\",\n                  \"headerValue\" : \"*\",\n                  \"type\" : \"apply-value-and-content-signatures\"\n               },\n               {\n                  \"headerName\" : \"Content-Type\",\n                  \"headerOrder\" : \"1\",\n                  \"headerValue\" : \"*form*\",\n                  \"type\" : \"form-data\"\n               },\n               {\n                  \"contentProfile\" : {\n                     \"name\" : \"Default\"\n                  },\n                  \"headerName\" : \"Content-Type\",\n                  \"headerOrder\" : \"2\",\n                  \"headerValue\" : \"*json*\",\n                  \"type\" : \"json\"\n               },\n               {\n                  \"contentProfile\" : {\n                     \"name\" : \"Default\"\n                  },\n                  \"headerName\" : \"Content-Type\",\n                  \"headerOrder\" : \"3\",\n                  \"headerValue\" : \"*xml*\",\n                  \"type\" : \"xml\"\n               }\n            ],\n            \"wildcardIncludesSlash\" : true,\n            \"wildcardOrder\" : 2\n         },\n         {\n            \"attackSignaturesCheck\" : true,\n            \"clickjackingProtection\" : false,\n            \"description\" : \"\",\n            \"disallowFileUploadOfExecutables\" : false,\n            \"html5CrossOriginRequestsEnforcement\" : {\n               \"enforcementMode\" : \"disabled\"\n            },\n            \"isAllowed\" : true,\n            \"mandatoryBody\" : false,\n            \"metacharsOnUrlCheck\" : false,\n            \"method\" : \"*\",\n            \"methodsOverrideOnUrlCheck\" : false,\n            \"name\" : \"*\",\n            \"performStaging\" : false,\n            \"protocol\" : \"https\",\n            \"type\" : \"wildcard\",\n            \"urlContentProfiles\" : [\n               {\n                  \"headerName\" : \"*\",\n                  \"headerOrder\" : \"default\",\n                  \"headerValue\" : \"*\",\n                  \"type\" : \"apply-value-and-content-signatures\"\n               },\n               {\n                  \"headerName\" : \"Content-Type\",\n                  \"headerOrder\" : \"1\",\n                  \"headerValue\" : \"*form*\",\n                  \"type\" : \"form-data\"\n               },\n               {\n                  \"contentProfile\" : {\n                     \"name\" : \"Default\"\n                  },\n                  \"headerName\" : \"Content-Type\",\n                  \"headerOrder\" : \"2\",\n                  \"headerValue\" : \"*json*\",\n                  \"type\" : \"json\"\n               },\n               {\n                  \"contentProfile\" : {\n                     \"name\" : \"Default\"\n                  },\n                  \"headerName\" : \"Content-Type\",\n                  \"headerOrder\" : \"3\",\n                  \"headerValue\" : \"*xml*\",\n                  \"type\" : \"xml\"\n               }\n            ],\n            \"wildcardIncludesSlash\" : true,\n            \"wildcardOrder\" : 1\n         }\n      ],\n      \"websocket-urls\" : [\n         {\n            \"allowBinaryMessage\" : true,\n            \"allowJsonMessage\" : true,\n            \"allowTextMessage\" : true,\n            \"checkBinaryMessageMaxSize\" : false,\n            \"checkMessageFrameMaxCount\" : false,\n            \"checkMessageFrameMaxSize\" : false,\n            \"checkPayload\" : true,\n            \"description\" : \"\",\n            \"html5CrossOriginRequestsEnforcement\" : {\n               \"enforcementMode\" : \"disabled\"\n            },\n            \"isAllowed\" : true,\n            \"jsonProfile\" : {\n               \"name\" : \"Default\"\n            },\n            \"metacharsOnWebsocketUrlCheck\" : false,\n            \"name\" : \"*\",\n            \"performStaging\" : false,\n            \"plainTextProfile\" : {\n               \"name\" : \"Default\"\n            },\n            \"protocol\" : \"wss\",\n            \"type\" : \"wildcard\",\n            \"unsupportedExtensions\" : \"remove\",\n            \"wildcardIncludesSlash\" : true,\n            \"wildcardOrder\" : 2\n         },\n         {\n            \"allowBinaryMessage\" : true,\n            \"allowJsonMessage\" : true,\n            \"allowTextMessage\" : true,\n            \"checkBinaryMessageMaxSize\" : false,\n            \"checkMessageFrameMaxCount\" : false,\n            \"checkMessageFrameMaxSize\" : false,\n            \"checkPayload\" : true,\n            \"description\" : \"\",\n            \"html5CrossOriginRequestsEnforcement\" : {\n               \"enforcementMode\" : \"disabled\"\n            },\n            \"isAllowed\" : true,\n            \"jsonProfile\" : {\n               \"name\" : \"Default\"\n            },\n            \"metacharsOnWebsocketUrlCheck\" : false,\n            \"name\" : \"*\",\n            \"performStaging\" : false,\n            \"plainTextProfile\" : {\n               \"name\" : \"Default\"\n            },\n            \"protocol\" : \"ws\",\n            \"type\" : \"wildcard\",\n            \"unsupportedExtensions\" : \"remove\",\n            \"wildcardIncludesSlash\" : true,\n            \"wildcardOrder\" : 1\n         }\n      ],\n      \"xml-profiles\" : [\n         {\n            \"attachmentsInSoapMessages\" : false,\n            \"attackSignaturesCheck\" : true,\n            \"defenseAttributes\" : {\n               \"allowCDATA\" : true,\n               \"allowDTDs\" : true,\n               \"allowExternalReferences\" : true,\n               \"allowProcessingInstructions\" : true,\n               \"maximumAttributeValueLength\" : \"any\",\n               \"maximumAttributesPerElement\" : \"any\",\n               \"maximumChildrenPerElement\" : \"any\",\n               \"maximumDocumentDepth\" : \"any\",\n               \"maximumDocumentSize\" : \"any\",\n               \"maximumElements\" : \"any\",\n               \"maximumNSDeclarations\" : \"any\",\n               \"maximumNameLength\" : \"any\",\n               \"maximumNamespaceLength\" : \"any\",\n               \"tolerateCloseTagShorthand\" : true,\n               \"tolerateLeadingWhiteSpace\" : true,\n               \"tolerateNumericNames\" : true\n            },\n            \"description\" : \"Default XML Profile\",\n            \"enableWss\" : false,\n            \"followSchemaLinks\" : false,\n            \"inspectSoapAttachments\" : false,\n            \"metacharAttributeCheck\" : false,\n            \"metacharElementCheck\" : false,\n            \"name\" : \"Default\",\n            \"useXmlResponsePage\" : false,\n            \"validationFiles\" : [],\n            \"validationSoapActionHeader\" : false\n         }\n      ]\n   }\n}\n`))
		_, _ = fmt.Fprintf(w, `{
    "isBase64": false,
    "minimal": false,
    "status": "COMPLETED",
    "includeVulnerabilityAssessmentConfigurationAndData": true,
    "kind": "tm:asm:tasks:export-policy:export-policy-taskstate",
    "selfLink": "https://localhost/mgmt/tm/asm/tasks/export-policy/uWHzxaTZGWL0X3ft7lYPEQ?ver=16.1.0",
    "policyReference": {
        "link": "https://localhost/mgmt/tm/asm/policies/LieFcG9wmFVfllNRyw144Q?ver=16.1.0",
        "fullPath": "/Common/mytestpolicy"
    },
    "endTime": "2023-01-04T08:00:01Z",
    "startTime": "2023-01-04T07:59:59Z",
    "id": "uWHzxaTZGWL0X3ft7lYPEQ",
    "inline": true,
    "executionStartTime": "2023-01-04T07:59:59Z",
    "lastUpdateMicros": 1.672819201e+15,
    "exportSuggestions": false,
    "format": "json",
    "result": {
        "fileSize": 72092,
        "file": "%v",
        "message": "Policy '/Common/mytestpolicy' was successfully exported."
    }
}`, filedata)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipAwafPolicyCreate(resourceName, server.URL),
			},
			{
				Config:             testAccBigipAwafPolicyModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
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
  template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
  application_language = "utf-8"
  enforcement_mode     = "transparent"
  description          = "Rapid Deployment-1"
  policy_builder {
    learning_mode = "disabled"
  }
  //server_technologies = ["MySQL", "Unix/Linux", "MongoDB"]
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
  template_name        = "POLICY_TEMPLATE_RAPID_DEPLOYMENT"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  description          = "Rapid Deployment-1"
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

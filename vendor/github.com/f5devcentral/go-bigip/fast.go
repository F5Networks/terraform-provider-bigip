/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"log"
)

type FastParameters struct {
	TenantName      string   `json:"tenant_name,omitempty"`
	ApplicationName string   `json:"application_name,omitempty"`
	VirtualPort     int      `json:"virtual_port,omitempty"`
	VirtualAddress  string   `json:"virtual_address,omitempty"`
	ServerPort      int      `json:"server_port,omitempty"`
	ServerAddresses []string `json:"server_addresses,omitempty"`
}

type FastTemplateType map[string]interface{}

type Fasttemplate struct {
	Name       string          `json:"name,omitempty"`
	Parameters *FastParameters `json:"parameters,omitempty"`
}

const (
	uriFast         = "fast"
	uriApplications = "applications"
	uriTasks        = "tasks"
)

func (b *BigIP) CreateFastTemplate(template *Fasttemplate) error {
	//b.getfastTask()
	return b.post(template, uriMgmt, uriShared, uriFast, uriApplications)
}

func (b *BigIP) GetFastTemplate(tenantName string, applicationName string) (*Fasttemplate, error) {
	var fastAppResult FastTemplateType
	fastAppResult = make(map[string]interface{})
	err, _ := b.getForEntity(&fastAppResult, uriMgmt, uriShared, uriFast, uriApplications, tenantName, applicationName)
	if err != nil {
		return nil, err
	}
	fastParams, err := maptoStruct(fastAppResult["constants"].(map[string]interface{})["fast"].(map[string]interface{})["view"])
	if err != nil {
		return nil, err
	}
	fastTemdata := &Fasttemplate{
		Name:       fastAppResult["constants"].(map[string]interface{})["fast"].(map[string]interface{})["template"].(string),
		Parameters: fastParams,
	}
	log.Printf("[DEBUG] Structure data for getFastApp:%+v", fastTemdata)
	return fastTemdata, nil
}

func (b *BigIP) DeleteFastTemplate(tenantName string, applicationName string) error {
	return b.delete(uriMgmt, uriShared, uriFast, uriApplications, tenantName, applicationName)
}

func maptoStruct(body interface{}) (*FastParameters, error) {
	jsonbody, err := jsonMarshal(body)
	if err != nil {
		return nil, err
	}
	fastData := FastParameters{}
	if err := json.Unmarshal(jsonbody, &fastData); err != nil {
		return nil, err
	}
	return &fastData, nil

}

func (b *BigIP) getfastTask() error {

	//b.getForEntity(&fastAppResult, uriMgmt, uriShared, uriFast, uriTasks)
	//jsonbody, err := jsonMarshal(body)
	//if err != nil {
	//      return nil, err
	//}
	return nil

}

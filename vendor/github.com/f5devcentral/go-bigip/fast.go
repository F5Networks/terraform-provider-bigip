/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"log"
	"time"
)

type FastParameters struct {
	TenantName      string   `json:"tenant_name,omitempty"`
	ApplicationName string   `json:"application_name,omitempty"`
	VirtualPort     int      `json:"virtual_port,omitempty"`
	VirtualAddress  string   `json:"virtual_address,omitempty"`
	ServerPort      int      `json:"server_port,omitempty"`
	ServerAddresses []string `json:"server_addresses,omitempty"`
}
type FastTaskType struct {
	Code       int64    `json:"code,omitempty"`
	ID         string   `json:"id,omitempty"`
	Message    string   `json:"message,omitempty"`
	Name       string   `json:"name,omitempty"`
	Parameters struct{} `json:"parameters,omitempty"`
}

type fastAppResponse struct {
	Code    int64 `json:"code,omitempty"`
	Message struct {
		ID         string   `json:"id,omitempty"`
		Name       string   `json:"name,omitempty"`
		Parameters struct{} `json:"parameters,omitempty"`
	} `json:"message,omitempty"`
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
	err := b.postFastTemplate(template)
	if err != nil {
		return err
	}
	return nil
}

func (b *BigIP) postFastTemplate(template *Fasttemplate) error {
	//b.getfastTaskid()
	resp, err := b.fastPost(template, uriMgmt, uriShared, uriFast, uriApplications)
	if err != nil {
		return err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respCode, respID := respRef["code"], respRef["message"].(map[string]interface{})["id"].(string)
	//log.Printf("[DEBUG]Code = %v,ID = %v", respCode, respID)
	log.Printf("[DEBUG]Creating Application with ID  = %v", respID)
	for respCode != 200 {
		fastTask, err := b.getfastTaskstatus(respID)
		if err != nil {
			return err
		}
		respCode = fastTask.Code
		if fastTask.Code == 200 {
			log.Printf("[DEBUG]Sucessfully Created Application with ID  = %v", respID)
			break // break here
		}
		if fastTask.Code == 503 {
			log.Printf("[DEBUG] Failed  Creating Application with ID  = %v", respID)
			log.Printf("[DEBUG] Sleeping for 1 Sec")
			time.Sleep(1 * time.Second)
			//break
			return b.postFastTemplate(template)
		}
	}
	return nil
	//return b.post(template, uriMgmt, uriShared, uriFast, uriApplications)
}

func (b *BigIP) GetFastTemplate(tenantName string, applicationName string) (*Fasttemplate, error) {
	var fastAppResult FastTemplateType
	fastAppResult = make(map[string]interface{})
	err, _ := b.getForEntity(&fastAppResult, uriMgmt, uriShared, uriFast, uriApplications, tenantName, applicationName)
	if err != nil {
		return nil, err
	}
	/*
		if len(fastAppResult) == 0 {
			log.Printf("[DEBUG] Sleeping for 1 Second to reflect application in BIGIP")
			time.Sleep(1 * time.Second)
			err, _ := b.getForEntity(&fastAppResult, uriMgmt, uriShared, uriFast, uriApplications, tenantName, applicationName)
			if err != nil {
				return nil, err
			}
		}*/

	fastParams, err := maptoStruct(fastAppResult["constants"].(map[string]interface{})["fast"].(map[string]interface{})["view"])
	if err != nil {
		return nil, err
	}
	fastTemdata := &Fasttemplate{
		Name:       fastAppResult["constants"].(map[string]interface{})["fast"].(map[string]interface{})["template"].(string),
		Parameters: fastParams,
	}
	//log.Printf("[DEBUG] Structure data for getFastApp:%+v", fastTemdata)
	return fastTemdata, nil
}

//func (b *BigIP) DeleteFastTemplate(tenantName string, applicationName string) error {
//	return b.delete(uriMgmt, uriShared, uriFast, uriApplications, tenantName, applicationName)
//}

func (b *BigIP) DeleteFastTemplate(tenantName string, applicationName string) error {
	resp, err := b.fastDelete(uriMgmt, uriShared, uriFast, uriApplications, tenantName, applicationName)
	if err != nil {
		return err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respCode, respID := respRef["status"], respRef["body"].(map[string]interface{})["id"].(string)
	log.Printf("[DEBUG]Deleting tenantName = %v,applicationName = %v with ID=%v", tenantName, applicationName, respID)

	for respCode != 200 {
		fastTask, err := b.getfastTaskstatus(respID)
		if err != nil {
			return err
		}
		respCode = fastTask.Code
		//log.Printf("[DEBUG]Code = %v,respCode = %v,ID = %v", fastTask.Code, respCode, fastTask.ID)
		if fastTask.Code == 200 {
			log.Printf("[DEBUG]Delete Success for tenantName = %v,applicationName = %v with ID=%v", tenantName, applicationName, fastTask.ID)
			break // break here
		}
		if fastTask.Code == 503 {
			log.Printf("[DEBUG]Delete Failed for tenantName = %v,applicationName = %v with ID=%v", tenantName, applicationName, fastTask.ID)
			log.Printf("[DEBUG] Waiting for 2 Sec")
			time.Sleep(2000 * time.Millisecond)
			//break
			return b.DeleteFastTemplate(tenantName, applicationName)
		}
	}
	return nil
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

func (b *BigIP) getfastTaskid() error {
	var taskList []FastTaskType
	err, _ := b.getForEntity(&taskList, uriMgmt, uriShared, uriFast, uriTasks)
	if err != nil {
		return err
	}
	for l := range taskList {
		if taskList[l].Message == "in progress" {
			//log.Printf("Id = %v, Name = %v,Code = %v", taskList[l].ID, taskList[l].Message, taskList[l].Code)
			//time.Sleep(1 * time.Second)
		}
	}
	return nil
}
func (b *BigIP) getfastTaskstatus(id string) (*FastTaskType, error) {
	var taskList FastTaskType
	err, _ := b.getForEntity(&taskList, uriMgmt, uriShared, uriFast, uriTasks, id)
	if err != nil {
		return nil, err
	}
	return &taskList, nil
}

func (b *BigIP) pollingStatus(id string) bool {
	var taskList FastTaskType
	err, _ := b.getForEntity(&taskList, uriMgmt, uriShared, uriFast, uriTasks, id)
	if err != nil {
		return false
	}
	if taskList.Code != 200 && taskList.Code != 503 {
		time.Sleep(1 * time.Second)
		return b.pollingStatus(id)
	}
	if taskList.Code == 503 {
		return false
	}
	return true
	//return &taskList, nil
}

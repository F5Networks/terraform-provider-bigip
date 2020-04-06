package bigip

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const as3SchemaLatestURL = "https://raw.githubusercontent.com/F5Networks/f5-appsvcs-extension/master/schema/latest/as3-schema.json"

type as3Validate struct {
	as3SchemaURL    string
	as3SchemaLatest string
}

func ValidateAS3Template(as3ExampleJson string) bool {
	myAs3 := &as3Validate{
		as3SchemaLatestURL,
		"",
	}
	err := myAs3.fetchAS3Schema()
	if err != nil {
		fmt.Errorf("As3 Schema Fetch failed: %s", err)
		return false
	}

	schemaLoader := gojsonschema.NewStringLoader(myAs3.as3SchemaLatest)
	//schemaLoader := gojsonschema.NewReferenceLoader("file:///Users/chinthalapalli/go/src/github.com/Practice/as3-schema-3.13.2-1-cis.json")
	documentLoader := gojsonschema.NewStringLoader(as3ExampleJson)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Errorf("%s", err)
		return false
	}
	if !result.Valid() {
		log.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			log.Printf("- %s\n", desc)
		}
		return false
	}
	return true
}

func (as3 *as3Validate) fetchAS3Schema() error {
	res, resErr := http.Get(as3.as3SchemaURL)
	if resErr != nil {
		log.Printf("Error while fetching latest as3 schema : %v", resErr)
		return resErr
	}
	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Unable to read the as3 template from json response body : %v", err)
			return err
		}
		defer res.Body.Close()
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			log.Printf("Unable to unmarshal json response body : %v", err)
			return err
		}
		jsonMap["$id"] = as3SchemaLatestURL
		byteJSON, err := json.Marshal(jsonMap)
		if err != nil {
			log.Printf("Unable to marshal : %v", err)
			return err
		}
		as3.as3SchemaLatest = string(byteJSON)
		return err
	}
	return nil
}

type As3AllTaskType struct {
	Items []As3TaskType `json:"items,omitempty"`
}
type As3TaskType struct {
	ID string `json:"id,omitempty"`
	//Declaration struct{} `json:"declaration,omitempty"`
	Results []Results1 `json:"results,omitempty"`
}
type Results1 struct {
	Code      int64  `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	LineCount int64  `json:"lineCount,omitempty"`
	Host      string `json:"host,omitempty"`
	Tenant    string `json:"tenant,omitempty"`
	RunTime   int64  `json:"runTime,omitempty"`
}

func (b *BigIP) PostAs3Bigip(as3NewJson string) error {
	resp, err := b.postReq(as3NewJson, uriMgmt, uriShared, uriAppsvcs, uriAsyncDeclare)
	if err != nil {
		return err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respID := respRef["id"].(string)
	taskStatus, err := b.getas3Taskstatus(respID)
	respCode := taskStatus.Results[0].Code
	log.Printf("[DEBUG]Code = %v,ID = %v", respCode, respID)
	for respCode != 200 {
		fastTask, err := b.getas3Taskstatus(respID)
		if err != nil {
			return err
		}
		respCode = fastTask.Results[0].Code
		if respCode == 200 {
			log.Printf("[DEBUG]Sucessfully Created Application with ID  = %v", respID)
			break // break here
		}
		if respCode == 503 {
			taskIds, err := b.getas3Taskid()
			if err != nil {
				return err
			}
			for _, id := range taskIds {
				if b.pollingStatus(id) {
					return b.PostAs3Bigip(as3NewJson)
				}
			}
		}
	}

	return nil
}
func (b *BigIP) DeleteAs3Bigip(tenantName string) error {
	tenant := tenantName + "?async=true"
	resp, err := b.deleteReq(uriMgmt, uriShared, uriAppsvcs, uriDeclare, tenant)
	if err != nil {
		return err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respID := respRef["id"].(string)
	taskStatus, err := b.getas3Taskstatus(respID)
	respCode := taskStatus.Results[0].Code
	log.Printf("[DEBUG]Delete Code = %v,ID = %v", respCode, respID)
	for respCode != 200 {
		fastTask, err := b.getas3Taskstatus(respID)
		if err != nil {
			return err
		}
		respCode = fastTask.Results[0].Code
		if respCode == 200 {
			log.Printf("[DEBUG]Sucessfully Deleted Application with ID  = %v", respID)
			break // break here
		}
		if respCode == 503 {
			taskIds, err := b.getas3Taskid()
			if err != nil {
				return err
			}
			for _, id := range taskIds {
				if b.pollingStatus(id) {
					return b.DeleteAs3Bigip(tenantName)
				}
			}
		}
	}

	return nil

}
func (b *BigIP) ModifyAs3(name string, as3_json string) error {
	tenant := name + "?async=true"
	resp, err := b.fastPatch(as3_json, uriMgmt, uriShared, uriAppsvcs, uriDeclare, tenant)
	if err != nil {
		return err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respID := respRef["id"].(string)
	taskStatus, err := b.getas3Taskstatus(respID)
	respCode := taskStatus.Results[0].Code
	for respCode != 200 {
		fastTask, err := b.getas3Taskstatus(respID)
		if err != nil {
			return err
		}
		respCode = fastTask.Results[0].Code
		if respCode == 200 {
			log.Printf("[DEBUG]Sucessfully Modified Application with ID  = %v", respID)
			break // break here
		}
		if respCode == 503 {
			taskIds, err := b.getas3Taskid()
			if err != nil {
				return err
			}
			for _, id := range taskIds {
				if b.pollingStatus(id) {
					return b.ModifyAs3(name, as3_json)
				}
			}
		}
	}

	return nil

}
func (b *BigIP) GetAs3(name string) (string, error) {
	//name = name + "?show=base"
	as3Json := make(map[string]interface{})
	as3Json["class"] = "AS3"
	as3Json["action"] = "deploy"
	as3Json["persist"] = true
	adcJson := make(map[string]interface{})
	//as3json, err, ok := b.getForEntityas3(uriMgmt, uriShared, uriAppsvcs, uriDeclare, name)
	err, ok := b.getForEntity(&adcJson, uriMgmt, uriShared, uriAppsvcs, uriDeclare, name)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	delete(adcJson, "updateMode")
	delete(adcJson, "controls")
	as3Json["declaration"] = adcJson
	out, _ := json.Marshal(as3Json)
	as3String := string(out)
	log.Printf("[DEBUG] As3 response string :%+v", as3String)
	return as3String, nil
}
func (b *BigIP) getas3Taskstatus(id string) (*As3TaskType, error) {
	var taskList As3TaskType
	err, _ := b.getForEntity(&taskList, uriMgmt, uriShared, uriAppsvcs, uriTask, id)
	if err != nil {
		return nil, err
	}
	return &taskList, nil

}
func (b *BigIP) getas3Taskid() ([]string, error) {
	var taskList As3AllTaskType
	var taskIDs []string
	err, _ := b.getForEntity(&taskList, uriMgmt, uriShared, uriAppsvcs, uriTask)
	if err != nil {
		return taskIDs, err
	}
	for l := range taskList.Items {
		if taskList.Items[l].Results[0].Message == "in progress" {
			taskIDs = append(taskIDs, taskList.Items[l].ID)
		}
	}
	return taskIDs, nil
}
func (b *BigIP) pollingStatus(id string) bool {
	var taskList As3TaskType
	err, _ := b.getForEntity(&taskList, uriMgmt, uriShared, uriAppsvcs, uriTask, id)
	if err != nil {
		return false
	}
	if taskList.Results[0].Code != 200 && taskList.Results[0].Code != 503 {
		time.Sleep(1 * time.Second)
		return b.pollingStatus(id)
	}
	if taskList.Results[0].Code == 503 {
		return false
	}
	return true
}
func (b *BigIP) GetTenantList(body interface{}) []string {
	s := make([]string, 1)
	as3json := body.(string)
	resp := []byte(as3json)
	jsonRef := make(map[string]interface{})
	json.Unmarshal(resp, &jsonRef)
	for _, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok {
			for k, v := range rec {
				if _, ok := v.(map[string]interface{}); ok {
					log.Println(k)
					s = append(s, k)
				}
			}
		}
	}
	return s
}
func (b *BigIP) AddTeemAgent(body interface{}) string {
	var s string
	as3json := body.(string)
	resp := []byte(as3json)
	jsonRef := make(map[string]interface{})
	json.Unmarshal(resp, &jsonRef)
	//jsonRef["controls"] = map[string]interface{}{"class": "Controls", "userAgent": "Terraform Configured AS3"}
	for _, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok {
			rec["controls"] = map[string]interface{}{"class": "Controls", "userAgent": "Terraform Configured AS3"}
		}
	}
	jsonData, err := json.Marshal(jsonRef)
	if err != nil {
		log.Println(err)
	}
	s = string(jsonData)
	return s
}

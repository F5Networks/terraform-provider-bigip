package bigip

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const as3SchemaLatestURL = "https://raw.githubusercontent.com/F5Networks/f5-appsvcs-extension/master/schema/latest/as3-schema.json"
const doSchemaLatestURL = "https://raw.githubusercontent.com/F5Networks/terraform-provider-bigip/master/schemas/doschema.json"

const (
	uriSha          = "shared"
	uriAppsvcs      = "appsvcs"
	uriDecl         = "declare"
	uriInfo         = "info"
	uriTask         = "task"
	uriDeclare      = "declare"
	uriAsyncDeclare = "declare?async=true"
)

type as3Validate struct {
	as3SchemaURL    string
	as3SchemaLatest string
}

type doValidate struct {
	doSchemaURL    string
	doSchemaLatest string
}

type as3Version struct {
	Version       string `json:"version"`
	Release       string `json:"release"`
	SchemaCurrent string `json:"schemaCurrent"`
	SchemaMinimum string `json:"schemaMinimum"`
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

func ValidateDOTemplate(doExampleJson string) bool {
	myDO := &doValidate{
		doSchemaLatestURL,
		"",
	}
	log.Printf("[DEBUG] validating DO json against DO schema")
	err := myDO.fetchDOSchema()
	if err != nil {
		fmt.Errorf("DO Schema Fetch failed: %s", err)
		return false
	}

	schemaLoader := gojsonschema.NewStringLoader(myDO.doSchemaLatest)
	documentLoader := gojsonschema.NewStringLoader(doExampleJson)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Errorf("%s", err)
		return false
	}
	if !result.Valid() {
		log.Printf("DO JSON is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			log.Printf("- %s\n", desc)
		}
		return false
	} else {
		log.Printf("[DEBUG] DO Json  is valid\n")
	}
	return true
}

func (do *doValidate) fetchDOSchema() error {
	res, resErr := http.Get(do.doSchemaURL)
	if resErr != nil {
		log.Printf("Error while fetching latest DO schema : %v", resErr)
		return resErr
	}
	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Unable to read the DO template from json response body : %v", err)
			return err
		}
		defer res.Body.Close()
		do.doSchemaLatest = string(body)
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

func (b *BigIP) PostAs3Bigip(as3NewJson string, tenantFilter string) (error, string) {
	tenant := tenantFilter + "?async=true"
	successfulTenants := make([]string, 0)
	resp, err := b.postReq(as3NewJson, uriMgmt, uriShared, uriAppsvcs, uriDeclare, tenant)
	if err != nil {
		return err, ""
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
			return err, ""
		}
		respCode = fastTask.Results[0].Code
		if respCode != 0 && respCode != 503 {
			tenant_list, tenant_count := b.GetTenantList(as3NewJson)
			if tenantCompare(tenant_list, tenantFilter) == 1 {
				i := tenant_count - 1
				success_count := 0
				for i >= 0 {
					if fastTask.Results[i].Code == 200 {
						successfulTenants = append(successfulTenants, fastTask.Results[i].Tenant)
						success_count++
					}
					if fastTask.Results[i].Code >= 400 {
						log.Printf("[ERROR] : HTTP %d :: %s for tenant %v", fastTask.Results[i].Code, fastTask.Results[i].Message, fastTask.Results[i].Tenant)
					}
					i = i - 1
				}
				if success_count == tenant_count {
					log.Printf("[DEBUG]Sucessfully Created Application with ID  = %v", respID)
					break // break here
				} else if success_count == 0 {
					return errors.New(fmt.Sprintf("Tenant Creation failed")), ""
				} else {
					finallist := strings.Join(successfulTenants[:], ",")
					return errors.New(fmt.Sprintf("Partial Success")), finallist
				}
			}
			if respCode == 200 {
				log.Printf("[DEBUG]Sucessfully Created Application with ID  = %v", respID)
				break // break here
			}
			if respCode >= 400 {
				return errors.New(fmt.Sprintf("Tenant Creation failed")), ""
			}
		}
		if respCode == 503 {
			taskIds, err := b.getas3Taskid()
			if err != nil {
				return err, ""
			}
			for _, id := range taskIds {
				if b.pollingStatus(id) {
					return b.PostAs3Bigip(as3NewJson, tenantFilter)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	return nil, strings.Join(successfulTenants[:], ",")
}
/*func (b *BigIP) DeleteAs3Bigip(tenantName string) error {
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
		time.Sleep(3 * time.Second)
	}

	return nil

}*/




func (b *BigIP) DeleteAs3Bigip(tenantName string) (error, string) {
        tenant := tenantName + "?async=true"
        failedTenants := make([]string, 0)
        resp, err := b.deleteReq(uriMgmt, uriShared, uriAppsvcs, uriDeclare, tenant)
        if err != nil {
                return err, ""
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
                        return err, ""
                }
                respCode = fastTask.Results[0].Code
                if respCode != 0 && respCode != 503 {
                        tenant_count := len(strings.Split(tenantName, ","))
                        if tenant_count != 1 {
                                i := tenant_count - 1
                                success_count := 0
                                for i >= 0 {
                                        if fastTask.Results[i].Code == 200 {
                                                success_count++
                                        }
                                        if fastTask.Results[i].Code >= 400 {
                                                failedTenants = append(failedTenants, fastTask.Results[i].Tenant)
                                                log.Printf("[ERROR] : HTTP %d :: %s for tenant %v", fastTask.Results[i].Code, fastTask.Results[i].Message, fastTask.Results[i].Tenant)
                                        }
                                        i = i - 1
                                }
                                if success_count == tenant_count {
                                        log.Printf("[DEBUG]Sucessfully Deleted Application with ID  = %v", respID)
                                        break // break here
                                } else if success_count == 0 {
                                        return errors.New(fmt.Sprintf("Tenant Deletion failed")), ""
                                } else {
                                        finallist := strings.Join(failedTenants[:], ",")
                                        return errors.New(fmt.Sprintf("Partial Success")), finallist
                                }
                        }
                        if respCode == 200 {
                                log.Printf("[DEBUG]Sucessfully Deleted Application with ID  = %v", respID)
                                break // break here
                        }
                        if respCode >= 400 {
                                return errors.New(fmt.Sprintf("Tenant Deletion failed")), ""
                        }
                }
                if respCode == 503 {
                        taskIds, err := b.getas3Taskid()
                        if err != nil {
                                return err, ""
                        }
                        for _, id := range taskIds {
                                if b.pollingStatus(id) {
                                        return b.DeleteAs3Bigip(tenantName)
                                }
                        }
                }
                time.Sleep(3 * time.Second)
        }

        return nil, ""

}
func (b *BigIP) ModifyAs3(tenantFilter string, as3_json string) error {
	tenant := tenantFilter + "?async=true"
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
					return b.ModifyAs3(tenantFilter, as3_json)
				}
			}
		}
	}

	return nil

}
func (b *BigIP) GetAs3(name string) (string, error) {
	as3Json := make(map[string]interface{})
	as3Json["class"] = "AS3"
	as3Json["action"] = "deploy"
	as3Json["persist"] = true
	adcJson := make(map[string]interface{})
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
	return as3String, nil
}
func (b *BigIP) getAs3version() (*as3Version, error) {
	var as3Ver as3Version
	err, _ := b.getForEntity(&as3Ver, uriMgmt, uriShared, uriAppsvcs, uriInfo)
	if err != nil {
		return nil, err
	}
	return &as3Ver, nil
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
func (b *BigIP) GetTenantList(body interface{}) (string, int) {
	s := make([]string, 0)
	as3json := body.(string)
	resp := []byte(as3json)
	jsonRef := make(map[string]interface{})
	json.Unmarshal(resp, &jsonRef)
	for key, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
			for k, v := range rec {
				if rec2, ok := v.(map[string]interface{}); ok {
					found := 0
					for k1, v1 := range rec2 {
						if k1 == "class" && v1 == "Tenant" {
							found = 1
						}
					}
					if found == 1 {
						s = append(s, k)
					}
				}
			}
		}
	}
	tenant_list := strings.Join(s[:], ",")
	return tenant_list, len(s)
}
func (b *BigIP) AddTeemAgent(body interface{}) (string, error) {
	var s string
	as3json := body.(string)
	resp := []byte(as3json)
	jsonRef := make(map[string]interface{})
	json.Unmarshal(resp, &jsonRef)
	//jsonRef["controls"] = map[string]interface{}{"class": "Controls", "userAgent": "Terraform Configured AS3"}
	as3ver, err := b.getAs3version()
	if err != nil {
		return "", fmt.Errorf("Getting AS3 Version failed with %v", err)
	}
	if as3ver.Version == "" {
		return "", fmt.Errorf("Getting AS3 Version failed,please check AS3 installed?")
	}
	log.Printf("[DEBUG] AS3 Version:%+v", as3ver.Version)
	log.Printf("[DEBUG] Terraform Version:%+v", b.UserAgent)
	//userAgent, err := getVersion("/usr/local/bin/terraform")
	//log.Printf("[DEBUG] Terraform version:%+v", userAgent)
	res1 := strings.Split(as3ver.Version, ".")
	for _, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok {
			if intConvert(res1[0]) > 3 || intConvert(res1[1]) >= 18 {
				rec["controls"] = map[string]interface{}{"class": "Controls", "userAgent": b.UserAgent}
			}
		}
	}
	jsonData, err := json.Marshal(jsonRef)
	if err != nil {
		//log.Println(err)
		return "", fmt.Errorf("Getting AS3 Version failed with %v", err)
	}
	s = string(jsonData)
	return s, nil
}
func intConvert(v interface{}) int {
	if s, err := strconv.Atoi(v.(string)); err == nil {
		return s
	}
	return 0
}
func getVersion(tfBinary string) (string, error) {
	var versionRegex = regexp.MustCompile("Terraform v(.*?)(\\s.*)?\n")
	out, err := exec.Command(tfBinary, "version").Output()
	if err != nil {
		return "", err
	}
	versionOutput := string(out)
	match := versionRegex.FindStringSubmatch(versionOutput)
	ua := fmt.Sprintf("Terraform/%s", match[1])
	return ua, nil
}
func (b *BigIP) TenantDifference(slice1 []string, slice2 []string) string {
	var diff []string
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}

		}
	}
	diff_tenant_list := strings.Join(diff[:], ",")
	return diff_tenant_list
}
func tenantCompare(t1 string, t2 string) int {
	tenantList1 := strings.Split(t1, ",")
	tenantList2 := strings.Split(t2, ",")
	if len(tenantList1) == len(tenantList2) {
		return 1
	}
	return 0
}

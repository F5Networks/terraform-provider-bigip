package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	uriFast     = "fast"
	uriFasttask = "tasks"
	uriTempl    = "templatesets"
	uriFastApp  = "applications"
)

type FastPayload struct {
	Name       string                 `json:"name,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type FastTask struct {
	Id          string                 `json:"id,omitempty"`
	Code        int64                  `json:"code,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Tenant      string                 `json:"tenant,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Application string                 `json:"application,omitempty"`
	Operation   string                 `json:"operation,omitempty"`
}

type FastTemplateSet struct {
	Name            string        `json:"name,omitempty"`
	Hash            string        `json:"hash,omitempty"`
	Supported       bool          `json:"supported,omitempty"`
	Templates       []TmplArrType `json:"templates,omitempty"`
	Schemas         []TmplArrType `json:"schemas,omitempty"`
	Enabled         bool          `json:"enabled,omitempty"`
	UpdateAvailable bool          `json:"updateAvailable,omitempty"`
}

type TmplArrType struct {
	Name string `json:"name,omitempty"`
	Hash string `json:"hash,omitempty"`
}

type FastTCPJson struct {
	Tenant            string         `json:"tenant_name,omitempty"`
	Application       string         `json:"app_name,omitempty"`
	VirtualAddress    string         `json:"virtual_address,omitempty"`
	VirtualPort       interface{}    `json:"virtual_port,omitempty"`
	SnatEnable        bool           `json:"enable_snat,omitempty"`
	SnatAutomap       bool           `json:"snat_automap"`
	MakeSnatPool      bool           `json:"make_snatpool"`
	SnatPoolName      string         `json:"snatpool_name,omitempty"`
	SnatAddresses     []string       `json:"snat_addresses,omitempty"`
	PoolEnable        bool           `json:"enable_pool"`
	MakePool          bool           `json:"make_pool"`
	PoolName          string         `json:"pool_name,omitempty"`
	PoolMembers       []FastHttpPool `json:"pool_members,omitempty"`
	LoadBalancingMode string         `json:"load_balancing_mode,omitempty"`
	SlowRampTime      int            `json:"slow_ramp_time,omitempty"`
	MonitorEnable     bool           `json:"enable_monitor,omitempty"`
	MakeMonitor       bool           `json:"make_monitor"`
	TCPMonitor        string         `json:"monitor_name,omitempty"`
	MonitorInterval   int            `json:"monitor_interval,omitempty"`
}

type FastHttpJson struct {
	Tenant                 string         `json:"tenant_name,omitempty"`
	Application            string         `json:"app_name,omitempty"`
	VirtualAddress         string         `json:"virtual_address,omitempty"`
	VirtualPort            interface{}    `json:"virtual_port,omitempty"`
	SnatEnable             bool           `json:"enable_snat,omitempty"`
	SnatAutomap            bool           `json:"snat_automap"`
	MakeSnatPool           bool           `json:"make_snatpool"`
	SnatPoolName           string         `json:"snatpool_name,omitempty"`
	SnatAddresses          []string       `json:"snat_addresses,omitempty"`
	PoolEnable             bool           `json:"enable_pool"`
	MakePool               bool           `json:"make_pool"`
	TlsServerEnable        bool           `json:"enable_tls_server"`
	TlsClientEnable        bool           `json:"enable_tls_client"`
	TlsServerProfileCreate bool           `json:"make_tls_server_profile"`
	TlsServerProfileName   string         `json:"tls_server_profile_name,omitempty"`
	TlsCertName            string         `json:"tls_cert_name,omitempty"`
	TlsKeyName             string         `json:"tls_key_name,omitempty"`
	PoolName               string         `json:"pool_name,omitempty"`
	PoolMembers            []FastHttpPool `json:"pool_members,omitempty"`
	LoadBalancingMode      string         `json:"load_balancing_mode,omitempty"`
	SlowRampTime           int            `json:"slow_ramp_time,omitempty"`
	MonitorEnable          bool           `json:"enable_monitor,omitempty"`
	MakeMonitor            bool           `json:"make_monitor"`
	HTTPMonitor            string         `json:"monitor_name_http,omitempty"`
	HTTPSMonitor           string         `json:"monitor_name,omitempty"`
	MonitorAuth            bool           `json:"monitor_credentials"`
	MonitorUsername        string         `json:"monitor_username,omitempty"`
	MonitorPassword        string         `json:"monitor_passphrase,omitempty"`
	MonitorInterval        int            `json:"monitor_interval,omitempty"`
	MonitorSendString      string         `json:"monitor_send_string,omitempty"`
	MonitorResponse        string         `json:"monitor_expected_response,omitempty"`
}

type FastHttpPool struct {
	ServerAddresses []string `json:"serverAddresses,omitempty"`
	ServicePort     int      `json:"servicePort,omitempty"`
	ConnectionLimit int      `json:"connectionLimit,omitempty"`
	PriorityGroup   int      `json:"priorityGroup,omitempty"`
	ShareNodes      bool     `json:"shareNodes,omitempty"`
}

// UploadFastTemplate copies a template set from local disk to BIGIP
func (b *BigIP) UploadFastTemplate(tmplpath *os.File, tmplname string) error {
	_, err := b.UploadFastTemp(tmplpath, tmplname)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]Template Path:%+v", tmplpath.Name())
	payload := FastTemplateSet{
		Name: tmplname,
	}
	err = b.AddTemplateSet(&payload)
	if err != nil {
		return err
	}
	return nil
}

// AddTemplateSet installs a template set.
func (b *BigIP) AddTemplateSet(tmpl *FastTemplateSet) error {
	return b.post(tmpl, uriMgmt, uriSha, uriFast, uriTempl)
}

// GetTemplateSet retrieves a Template set by name. Returns nil if the Template set does not exist
func (b *BigIP) GetTemplateSet(name string) (*FastTemplateSet, error) {
	var tmpl FastTemplateSet
	err, ok := b.getForEntity(&tmpl, uriMgmt, uriSha, uriFast, uriTempl, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &tmpl, nil
}

// DeleteTemplateSet removes a template set.
func (b *BigIP) DeleteTemplateSet(name string) error {
	return b.delete(uriMgmt, uriSha, uriFast, uriTempl, name)
}

// GetFastApp retrieves a Application set by tenant and app name. Returns nil if the application does not exist
func (b *BigIP) GetFastApp(tenant, app string) (string, error) {
	var out []byte
	fastJson := make(map[string]interface{})
	err, ok := b.getForEntity(&fastJson, uriMgmt, uriShared, uriFast, uriFastApp, tenant, app)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	for key, value := range fastJson {
		if rec, ok := value.(map[string]interface{}); ok && key == "constants" {
			for k, v := range rec {
				if rec2, ok := v.(map[string]interface{}); ok && k == "fast" {
					for k1, v1 := range rec2 {
						if rec3, ok := v1.(map[string]interface{}); ok {
							if k1 == "view" {
								out, _ = json.Marshal(rec3)
							}
						}
					}
				}

			}
		}
	}
	fastString := string(out)

	return fastString, nil
}

// PostFastAppBigip used for posting FAST json file to BIGIP
func (b *BigIP) PostFastAppBigip(body, fastTemplate, userAgent string) (tenant, app string, err error) {
	param := []byte(body)
	jsonRef := make(map[string]interface{})
	json.Unmarshal(param, &jsonRef)
	payload := &FastPayload{
		Name:       fastTemplate,
		Parameters: jsonRef,
	}
	log.Printf("[DEBUG]payload = %+v", payload)
	resp, err := b.postReq(payload, uriMgmt, uriShared, uriFast, uriFastApp, userAgent)
	if err != nil {
		return "", "", err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respID := respRef["message"].([]interface{})[0].(map[string]interface{})["id"].(string)
	taskStatus, err := b.getFastTaskStatus(respID)
	if err != nil {
		return "", "", err
	}
	respCode := taskStatus.Code
	log.Printf("[DEBUG]Initial response code = %+v,ID = %+v", respCode, respID)
	for respCode != 200 {
		fastTask, err := b.getFastTaskStatus(respID)
		if err != nil {
			return "", "", err
		}
		respCode = fastTask.Code
		log.Printf("[DEBUG]Response code = %+v,ID = %+v", respCode, respID)
		if respCode == 200 {
			log.Printf("[DEBUG]Sucessfully Created Application with ID  = %v", respID)
			break // break here
		}
		if respCode >= 400 {
			return "", "", fmt.Errorf("FAST Application creation failed with :%+v", fastTask.Message)
		}
		time.Sleep(3 * time.Second)
	}
	return taskStatus.Tenant, taskStatus.Application, err
}

// ModifyFastAppBigip used for updating FAST application on BIGIP
func (b *BigIP) ModifyFastAppBigip(body, fastTenant, fastApp string) error {
	param := []byte(body)
	jsonRef := make(map[string]interface{})
	json.Unmarshal(param, &jsonRef)
	payload := &FastPayload{
		Parameters: jsonRef,
	}
	resp, err := b.fastPatch(payload, uriMgmt, uriShared, uriFast, uriFastApp, fastTenant, fastApp)
	if err != nil {
		return err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respID := respRef["message"].([]interface{})[0].(map[string]interface{})["id"].(string)
	taskStatus, err := b.getFastTaskStatus(respID)
	if err != nil {
		return err
	}
	respCode := taskStatus.Code
	log.Printf("[DEBUG]Code = %+v,ID = %+v", respCode, respID)
	for respCode != 200 {
		fastTask, err := b.getFastTaskStatus(respID)
		if err != nil {
			return err
		}
		respCode = fastTask.Code
		if respCode == 200 {
			log.Printf("[DEBUG]Sucessfully Modified Application with ID  = %v", respID)
			break // break here
		}
		if respCode >= 400 {
			return fmt.Errorf("FAST Application update failed with :%+v", fastTask.Message)
			//return fmt.Errorf("FAST Application update failed")
		}
		time.Sleep(3 * time.Second)
	}
	return err
}

// DeleteFastAppBigip used for deleting FAST application on BIGIP
func (b *BigIP) DeleteFastAppBigip(fastTenant, fastApp string) error {
	resp, err := b.deleteReq(uriMgmt, uriShared, uriFast, uriFastApp, fastTenant, fastApp)
	if err != nil {
		return err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respID := respRef["id"].(string)
	taskStatus, err := b.getFastTaskStatus(respID)
	if err != nil {
		return err
	}
	respCode := taskStatus.Code
	log.Printf("[DEBUG]Code = %+v,ID = %+v", respCode, respID)
	for respCode != 200 {
		fastTask, err := b.getFastTaskStatus(respID)
		if err != nil {
			return err
		}
		respCode = fastTask.Code
		if respCode == 200 {
			log.Printf("[DEBUG]Sucessfully Deleted Application with ID  = %v", respID)
			break // break here
		}
		if respCode >= 400 {
			return fmt.Errorf("FAST Application deletion failed")
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}

// getFastTaskStatus used to obtain status of async task from BIGIP
func (b *BigIP) getFastTaskStatus(id string) (*FastTask, error) {
	var taskList FastTask
	err, _ := b.getForEntity(&taskList, uriMgmt, uriShared, uriFast, uriFasttask, id)
	if err != nil {
		return nil, err
	}
	return &taskList, nil
}

// Upload a file
func (b *BigIP) UploadFastTemp(f *os.File, tmpName string) (*Upload, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return b.Upload(f, info.Size(), uriShared, uriFileTransfer, uriUploads, fmt.Sprintf("%s.zip", tmpName))
}

package bigip

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	uriWafPol       = "policies"
	uriUrls         = "urls"
	uriParams       = "parameters"
	uriWafSign      = "signatures"
	uriImportpolicy = "import-policy"
	uriExpPb        = "export-suggestions"
)

type ApplywafPolicy struct {
	Filename string `json:"filename,omitempty"`
	Policy   struct {
		FullPath string `json:"fullPath,omitempty"`
	} `json:"policy,omitempty"`
}

type PbExport struct {
	Status  string                 `json:"status,omitempty"`
	Task_id string                 `json:"id,omitempty"`
	Result  map[string]interface{} `json:"result,omitempty"`
}

type WafQueriedPolicies struct {
	WafPolicyList []WafQueriedPolicy `json:"items"`
}

type WafQueriedPolicy struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	Policy_id string `json:"id,omitempty"`
}

type Signatures struct {
	Signatures []Signature `json:"items"`
}

type Signature struct {
	Name        string `json:"name,omitempty"`
	ResourceId  string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	SignatureId int    `json:"signatureId,omitempty"`
	Type        string `json:"signatureType,omitempty"`
	Accuracy    string `json:"accuracy,omitempty"`
	Risk        string `json:"risk,omitempty"`
}

type WafUrlJson struct {
	Name                      string            `json:"name,omitempty"`
	Description               string            `json:"description,omitempty"`
	Type                      string            `json:"type,omitempty"`
	Protocol                  string            `json:"protocol,omitempty"`
	Method                    string            `json:"method,omitempty"`
	PerformStaging            bool              `json:"performStaging,omitempty"`
	SignatureOverrides        []WafUrlSig       `json:"signatureOverrides,omitempty"`
	MethodOverrides           []MethodOverrides `json:"methodOverrides,omitempty"`
	AttackSignaturesCheck     bool              `json:"attackSignaturesCheck,omitempty"`
	IsAllowed                 bool              `json:"isAllowed,omitempty"`
	MethodsOverrideOnUrlCheck bool              `json:"methodsOverrideOnUrlCheck,omitempty"`
}

type MethodOverrides struct {
	Allowed bool   `json:"allowed"` // as we can supply true and false, omitempty would automatically remove allowed = false which we do not want
	Method  string `json:"method,omitempty"`
}

type WafUrlSig struct {
	Enabled bool `json:"enabled"` // as we can supply true and false, omitempty would automatically remove allowed = false which we do not want
	Id      int  `json:"signatureId,omitempty"`
}

type WafPolicies struct {
	WafPolicies []WafPolicy `json:"items,omitempty"`
}

type WafPolicy struct {
	Name        string `json:"name,omitempty"`
	Partition   string `json:"partition,omitempty"`
	Description string `json:"description,omitempty"`
	FullPath    string `json:"fullPath,omitempty"`
	ID          string `json:"id,omitempty"`
	Template    struct {
		Name string `json:"name,omitempty"`
	} `json:"template,omitempty"`
	HasParent           bool   `json:"hasParent,omitempty"`
	ApplicationLanguage string `json,"applicationLanguage,omitempty"`
	EnablePassiveMode   bool   `json:"enablePassiveMode,omitempty"`
	ProtocolIndependent bool   `json:"protocolIndependent,omitempty"`
	CaseInsensitive     bool   `json:"caseInsensitive,omitempty"`
	EnforcementMode     string `json:"enforcementMode,omitempty"`
	Type                string `json:"type,omitempty"`
	//Parameters          []WafEntityParameter `json:"parameters,omitempty"`
	ServerTechnologies []struct {
		ServerTechnologyName string `json:"serverTechnologyName,omitempty"`
	} `json:"server-technologies,omitempty"`
	Urls           []WafUrlJson  `json:"urls,omitempty"`
	VirtualServers []interface{} `json:"virtualServers,omitempty"`
}

type ImportStatus struct {
	IsBase64                  bool   `json:"isBase64,omitempty"`
	Status                    string `json:"status"`
	GetPolicyAttributesOnly   bool   `json:"getPolicyAttributesOnly,omitempty"`
	Filename                  string `json:"filename"`
	ID                        string `json:"id"`
	RetainInheritanceSettings bool   `json:"retainInheritanceSettings"`
	Result                    struct {
		Message string `json:"message"`
	} `json:"result,omitempty"`
}

func (b *BigIP) GetWafSignature(signatureid int) (*Signatures, error) {
	var signature Signatures
	var query = fmt.Sprintf("?$filter=signatureId+eq+%d", signatureid)
	err, _ := b.getForEntity(&signature, uriMgmt, uriTm, uriAsm, uriWafSign, query)
	if err != nil {
		return nil, err
	}
	return &signature, nil
}

func (b *BigIP) GetWafPolicyId(policyName, partition string) (string, error) {
	var self WafQueriedPolicies
	query := fmt.Sprintf("?$filter=contains(name,'%s')+and+contains(partition,'%s')&$select=name,partition,id", policyName, partition)
	err, _ := b.getForEntity(&self, uriMgmt, uriTm, uriAsm, uriWafPol, query)

	if err != nil {
		return "", err
	}

	for _, policy := range self.WafPolicyList {
		if policy.Name == policyName && policy.Partition == partition {
			return policy.Policy_id, nil
		}
	}

	return "", fmt.Errorf("could not get the policy ID")
}

func (b *BigIP) PostPbExport(payload interface{}) (*PbExport, error) {
	var export PbExport
	resp, err := b.postReq(payload, uriMgmt, uriTm, uriAsm, uriTasks, uriExpPb)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp, &export)
	return &export, nil
}
func (b *BigIP) GetWafPbExportResult(id string) (*PbExport, error) {
	var pbexport PbExport
	err, _ := b.getForEntity(&pbexport, uriMgmt, uriShared, uriFast, uriFasttask, id)
	if err != nil {
		return nil, err
	}
	return &pbexport, nil
}

func (b *BigIP) GetWafPolicyQuery(wafPolicyName string) (*WafPolicy, error) {
	var wafPolicies WafPolicies
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("fullPath eq '%s'", wafPolicyName))
	var query = fmt.Sprintf("?$%v", params.Encode())
	err, _ := b.getForEntity(&wafPolicies, uriMgmt, uriTm, uriAsm, uriWafPol, query)
	if err != nil {
		return nil, err
	}
	if len(wafPolicies.WafPolicies) == 0 {
		return nil, fmt.Errorf("[ERROR] WafPolicy: %+v not found", wafPolicyName)
	}
	// if successful filter query will return a list with a single item
	wafPolicy := wafPolicies.WafPolicies[0]

	return &wafPolicy, nil
}

func (b *BigIP) GetWafPolicy(policyID string) (*WafPolicy, error) {
	var wafPolicy WafPolicy
	err, _ := b.getForEntity(&wafPolicy, uriMgmt, uriTm, uriAsm, uriWafPol, policyID)
	if err != nil {
		return nil, err
	}
	return &wafPolicy, nil
}

func (b *BigIP) GetImportStatus(taskId string) error {
	var importStatus ImportStatus
	err, _ := b.getForEntity(&importStatus, uriMgmt, uriTm, uriAsm, uriTasks, uriImportpolicy, taskId)
	if err != nil {
		return err
	}
	if importStatus.Status == "COMPLETED" {
		return nil
	}
	if importStatus.Status == "FAILURE" {
		return fmt.Errorf("[ERROR] WafPolicy import failed with :%+v", importStatus.Result)
	}
	if importStatus.Status == "STARTED" {
		time.Sleep(5 * time.Second)
		return b.GetImportStatus(taskId)
	}
	return nil
}

// DeleteWafPolicy removes waf Policy
func (b *BigIP) DeleteWafPolicy(policyId string) error {
	return b.delete(uriMgmt, uriTm, uriAsm, uriWafPol, policyId)
}

// ImportAwafJson import Awaf Json from local machine to BIGIP
func (b *BigIP) ImportAwafJson(awafPolicyName, awafJsonContent string) (string, error) {
	certbyte := []byte(awafJsonContent)
	policyName := awafPolicyName[strings.LastIndex(awafPolicyName, "/")+1:]
	_, err := b.UploadAsmBytes(certbyte, fmt.Sprintf("%s.json", policyName))
	if err != nil {
		return "", err
	}
	policyPath := struct {
		FullPath string `json:"fullPath,omitempty"`
	}{
		FullPath: awafPolicyName,
	}
	applywaf := ApplywafPolicy{
		Filename: fmt.Sprintf("%s.json", policyName),
		Policy:   policyPath,
	}
	resp, err := b.postReq(applywaf, uriMgmt, uriTm, uriAsm, uriTasks, uriImportpolicy)
	if err != nil {
		return "", err
	}
	var taskStatus ImportStatus
	err = json.Unmarshal(resp, &taskStatus)
	if err != nil {
		return "", err
	}
	return taskStatus.ID, nil
}

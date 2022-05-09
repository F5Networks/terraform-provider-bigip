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

type PbExport struct {
	Status  string                 `json:"status,omitempty"`
	Task_id string                 `json:"id,omitempty"`
	Result  map[string]interface{} `json:"result,omitempty"`
}

type WafEntityUrl struct {
	Name                            string `json:"name,omitempty"`
	WildcardOrder                   int    `json:"wildcardOrder,omitempty"`
	Protocol                        string `json:"protocol,omitempty"`
	Method                          string `json:"method,omitempty"`
	Type                            string `json:"type,omitempty"`
	AttackSignaturesCheck           bool   `json:"attackSignaturesCheck,omitempty"`
	MetacharsOnURLCheck             bool   `json:"metacharsOnUrlCheck,omitempty"`
	CanChangeDomainCookie           bool   `json:"canChangeDomainCookie,omitempty"`
	ClickjackingProtection          bool   `json:"clickjackingProtection,omitempty"`
	IsAllowed                       bool   `json:"isAllowed,omitempty"`
	MandatoryBody                   bool   `json:"mandatoryBody"`
	DisallowFileUploadOfExecutables bool   `json:"disallowFileUploadOfExecutables,omitempty"`
	AllowRenderingInFrames          string `json:"allowRenderingInFrames,omitempty"`
}

type WafEntityParameters struct {
	WafEntityParametersList []WafEntityParameter `json:"items"`
}

type WafEntityParameter struct {
	Name                           string `json:"name,omitempty"`
	Description                    string `json:"description,omitempty"`
	Type                           string `json:"type,omitempty"`
	ValueType                      string `json:"valueType,omitempty"`
	AllowEmptyValue                bool   `json:"allowEmptyValue,omitempty"`
	AllowRepeatedParameterName     bool   `json:"allowRepeatedParameterName,omitempty"`
	AttackSignaturesCheck          bool   `json:"attackSignaturesCheck,omitempty"`
	CheckMaxValueLength            bool   `json:"checkMaxValueLength,omitempty"`
	CheckMinValueLength            bool   `json:"checkMinValueLength,omitempty"`
	DataType                       string `json:"dataType,omitempty"`
	EnableRegularExpression        bool   `json:"enableRegularExpression,omitempty"`
	IsBase64                       bool   `json:"isBase64,omitempty"`
	IsCookie                       bool   `json:"isCookie,omitempty"`
	IsHeader                       bool   `json:"isHeader,omitempty"`
	Level                          string `json:"level,omitempty"`
	Mandatory                      bool   `json:"mandatory,omitempty"`
	MetacharsOnParameterValueCheck bool   `json:"metacharsOnParameterValueCheck,omitempty"`
	ParameterLocation              string `json:"parameterLocation,omitempty"`
	PerformStaging                 bool   `json:"performStaging,omitempty"`
	SensitiveParameter             bool   `json:"sensitiveParameter,omitempty"`
	SignatureOverrides_Disable     []int  `json:"signatureOverrides_disable,omitempty"`
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
	HasParent           bool                 `json:"hasParent,omitempty"`
	ApplicationLanguage string               `json,"applicationLanguage,omitempty"`
	EnablePassiveMode   bool                 `json:"enablePassiveMode,omitempty"`
	ProtocolIndependent bool                 `json:"protocolIndependent,omitempty"`
	CaseInsensitive     bool                 `json:"caseInsensitive,omitempty"`
	EnforcementMode     string               `json:"enforcementMode,omitempty"`
	Type                string               `json:"type,omitempty"`
	Parameters          []WafEntityParameter `json:"parameters,omitempty"`
	ServerTechnologies  []struct {
		ServerTechnologyName string `json:"serverTechnologyName,omitempty"`
	} `json:"server-technologies,omitempty"`
	Urls           []WafEntityUrl `json:"urls,omitempty"`
	VirtualServers []interface{}  `json:"virtualServers,omitempty"`
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

type ApplywafPolicy struct {
	Filename string `json:"filename,omitempty"`
	Policy   struct {
		FullPath string `json:"fullPath,omitempty"`
	} `json:"policy,omitempty"`
}

type WafEntityURLs struct {
	WafEntityURLList []WafEntityURL `json:"items"`
}

type WafEntityURL struct {
	Name               string         `json:"name,omitempty"`
	Description        string         `json:"description,omitempty"`
	Type               string         `json:"type,omitempty"`
	Protocol           string         `json:"protocol,omitempty"`
	Method             string         `json:"method,omitempty"`
	MethodOverrides    string         `json:"methodOverrides,omitempty"`
	PerformStaging     bool           `json:"performStaging,omitempty"`
	SignatureOverrides []SignatureIDs `json:"signatureOverrides,omitempty"`
}

type SignatureIDs struct {
	SignatureReference []SigIDs
	Enabled            bool `json:"enabled,omitempty"`
}

type SigIDs struct {
	Link          string `json:"link,omitempty"`
	IsUserDefined bool   `json:"isUserDefined,omitempty"`
	Name          string `json:"name,omitempty"`
	SignatureId   int    `json:"signatureId,omitempty"`
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

func (b *BigIP) WafEntityParameters(policyId string) (*WafEntityParameters, error) {
	var self WafEntityParameters
	err, _ := b.getForEntity(&self, uriMgmt, uriTm, uriAsm, uriWafPol, uriParams)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIP) WafEntityUrls(policyId string) (*WafEntityURLs, error) {
	var self WafEntityURLs
	err, _ := b.getForEntity(&self, uriMgmt, uriTm, uriAsm, uriWafPol, uriUrls)
	if err != nil {
		return nil, err
	}
	return &self, nil
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
	err, _ := b.getForEntity(id, uriMgmt, uriShared, uriFast, uriFasttask, id)
	if err != nil {
		return nil, err
	}
	return &pbexport, nil
}

func (b *BigIP) CreateWafEntityUrl(config *WafEntityURL, policyId string) error {
	return b.post(config, uriMgmt, uriTm, uriAsm, uriWafPol, policyId, uriUrls)
}

func (b *BigIP) ModifyWafEntityUrl(config *WafEntityURL, urlId, policyId string) error {
	return b.patch(config, uriMgmt, uriTm, uriAsm, uriWafPol, policyId, uriUrls, urlId)
}

func (b *BigIP) DeleteWafEntityUrl(urlId, policyId string) error {
	return b.delete(uriMgmt, uriTm, uriAsm, uriWafPol, policyId, uriUrls, urlId)
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

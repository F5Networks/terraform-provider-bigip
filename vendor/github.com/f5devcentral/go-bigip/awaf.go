package bigip

import (
	"encoding/json"
	"fmt"
)

const (
	uriWafPol  = "policies"
	uriUrls    = "urls"
	uriParams  = "parameters"
	uriWafSign = "signatures"
	uriExpPb   = "export-suggestions"
)

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

type Parameter struct {
	Name                           string                   `json:"name,omitempty"`
	Description                    string                   `json:"description,omitempty"`
	Type                           string                   `json:"type,omitempty"`
	ValueType                      string                   `json:"valueType,omitempty"`
	AllowEmptyValue                bool                     `json:"allowEmptyValue,omitempty"`
	AllowRepeatedParameterName     bool                     `json:"allowRepeatedParameterName,omitempty"`
	AttackSignaturesCheck          bool                     `json:"attackSignaturesCheck,omitempty"`
	CheckMaxValueLength            bool                     `json:"checkMaxValueLength,omitempty"`
	CheckMinValueLength            bool                     `json:"checkMinValueLength,omitempty"`
	DataType                       string                   `json:"dataType,omitempty"`
	EnableRegularExpression        bool                     `json:"enableRegularExpression,omitempty"`
	IsBase64                       bool                     `json:"isBase64,omitempty"`
	IsCookie                       bool                     `json:"isCookie,omitempty"`
	IsHeader                       bool                     `json:"isHeader,omitempty"`
	Level                          string                   `json:"level,omitempty"`
	Mandatory                      bool                     `json:"mandatory,omitempty"`
	MetacharsOnParameterValueCheck bool                     `json:"metacharsOnParameterValueCheck,omitempty"`
	ParameterLocation              string                   `json:"parameterLocation,omitempty"`
	PerformStaging                 bool                     `json:"performStaging,omitempty"`
	SensitiveParameter             bool                     `json:"sensitiveParameter,omitempty"`
	SignatureOverrides             []map[string]interface{} `json:"signatureOverrides,omitempty"`
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

func (b *BigIP) GetWafEntityParameters(name, policyId string) (*Parameter, error) {
	var entityparam Parameter
	var query = fmt.Sprintf("?$filter=name+eq+%s", name)
	err, _ := b.getForEntity(&entityparam, uriMgmt, uriTm, uriAsm, uriWafPol, policyId, uriParams, query)
	if err != nil {
		return nil, err
	}
	return &entityparam, nil
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

package bigip

import "fmt"

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
	WafPolicyList []WafPolicy `json:"items"`
}

type WafPolicy struct {
	Name                  string        `json:"name,omitempty"`
	Partition             string        `json:"partition,omitempty"`
	Description           string        `json:"description,omitempty"`
	Template              string        `json:"template,omitempty"`
	Tags                  string        `json:"tags"`
	ApplicationLanguage   string        `json,"applicationLanguage,omitempty"`
	PolicyEnforcement     string        `json:"policyEnforcement,omitempty"`
	ServerTechnoDetection bool          `json:"serverTechnodetection,omitempty"`
	ServerTechnologies    []string      `json:"serverTechnologies,omitempty"`
	PolicyAdjustments     []interface{} `json:"polcyAdjustments,omitempty"`
	PolicyModifications   []interface{} `json:"policyModifications,omitempty"`
}

type WafEntityURLs struct {
	WafEntityURLList []WafEntityURL `json:"items"`
}

type WafEntityURL struct {
	Name                       string `json:"name,omitempty"`
	Description                string `json:"description,omitempty"`
	Type                       string `json:"type,omitempty"`
	Protocol                   string `json:"protocol,omitempty"`
	Method                     string `json:"method,omitempty"`
	MethodOverrides            string `json:"methodOverrides,omitempty"`
	PerformStaging             bool   `json:"performStaging,omitempty"`
	SignatureOverrides_disable bool   `json:"signatureOverrides_disable,omitempty"`
}

// This method is not correct as of now, it tries to access keys that are not there in WafPolicy struct yet
func (b *BigIP) GetPolicyId(policyName string) (string, error) {
	var self WafPolicies
	err, _ := b.getForEntity(&self, uriMgmt, uriTm, uriAsm, "policies")

	if err != nil {
		return "", err
	}

	for _, policy := range self.WafPolicyList {
		if policy.FullPath == "policyName" {
			return policy.Id, nil
		}
	}

	return "", fmt.Errorf("could not get the policy ID")
}

func (b *BigIP) WafEntityParameters(policyId string) (*WafEntityParameters, error) {
	var self WafEntityParameters
	err, _ := b.getForEntity(&self, uriMgmt, uriTm, uriAsm, "policies", "parameters")
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIP) GetEntityParameters(policyId, parameterId string) (*WafEntityParameter, error) {
	var wafEntityParameter WafEntityParameter
	err, _ := b.getForEntity(wafEntityParameter, uriMgmt, uriTm, uriAsm, "policies", policyId, "parameters", parameterId)
	if err != nil {
		return nil, err
	}
	return &wafEntityParameter, nil
}

func (b *BigIP) CreateWafEntityParameter(config *WafEntityParameter, policyId string) error {
	return b.post(config, uriMgmt, uriTm, uriAsm, "policies", policyId, "parameters")
}

func (b *BigIP) ModifyWafEntityParameter(config *WafEntityParameter, policyId string) error {
	return b.patch(config, uriMgmt, uriTm, uriAsm, "policies", policyId, "parameters")
}

func (b *BigIP) DeleteWafEntityParameter(parameterId, policyId string) error {
	return b.delete(uriMgmt, uriTm, uriAsm, "policies", policyId, "parameters", parameterId)
}

func (b *BigIP) GetWafPolicies(policyId string) (*WafPolicy, error) {
	var wafPolicy WafPolicy
	err, _ := b.getForEntity(&wafPolicy, uriMgmt, uriTm, uriAsm, "policies", policyId)

	if err != nil {
		return nil, err
	}

	return &wafPolicy, nil
}

func (b *BigIP) CreateWafPolicy(config *WafPolicy) error {
	return b.post(config, uriMgmt, uriTm, uriAsm, "policies")
}

func (b *BigIP) ModifyWafPolicy(config *WafPolicy, policyId string) error {
	return b.patch(config, uriMgmt, uriTm, uriAsm, "policies", policyId)
}

func (b *BigIP) DeleteWafPolicy(config *WafPolicy, policyId string) error {
	return b.delete(uriMgmt, uriTm, uriAsm, "policies", policyId)
}

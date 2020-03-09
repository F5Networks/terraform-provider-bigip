/*
Original work Copyright © 2015 Scott Ware
Modifications Copyright 2019 F5 Networks Inc
Licensed under the Apache License, Version 2.0 (the "License");
You may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.
*/
package bigip

type FastParameters struct {
                TenantName               string `json:"tenant_name,omitempty"`
                ApplicationName          string `json:"application_name,omitempty"`
                VirtualPort              int `json:"virtual_port,omitempty"`
                VirtualAddress            string `json:"virtual_address,omitempty"`
                ServerPort                int   `json:"server_port,omitempty"`
                ServerAddresses          []string `json:"server_addresses,omitempty"`
}

type Fasttemplate struct {
        Name                      string `json:"name,omitempty"`
        Parameters                FastParameters `json:"parameters,omitempty"`
}

const (
        uriFast           = "fast"
        uriApplications   = "applications"
)

func (b *BigIP) CreateFastTemplate(template *Fasttemplate) error {
	 return b.post(template, uriMgmt, uriShared, uriFast, uriApplications) }

func (b *BigIP) DeleteFastTemplate(tenantName string, applicationName string) error {
         return b.post(uriMgmt, uriShared, uriFast, uriApplications, tenantName, applicationName ) }

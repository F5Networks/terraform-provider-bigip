package bigip

import (
	"log"
	"os"
)

const (
	uriFast = "fast"
	uriTempl = "templatesets"
)


type FastTemplateSet struct {
	Name               		string  `json:"name,omitempty"`
	Hash    		   		string	`json:"hash,omitempty"`
	Supported 				bool 	`json:"supported,omitempty"`
	Templates 				[]TmplArrType `json:"templates,omitempty"`
	Schemas 				[]TmplArrType `json:"schemas,omitempty"`
	Enabled 				bool `json:"enabled,omitempty"`
	UpdateAvailable 		bool `json:"updateAvailable,omitempty"`
}

type TmplArrType struct {
	Name	string	`json:"name,omitempty"`
	Hash    string	`json:"hash,omitempty"`
}

// UploadFastTemplate copies a template set from local disk to BIGIP
func (b *BigIP) UploadFastTemplate(tmplpath *os.File, tmplname string) error {
	_, err := b.UploadFile(tmplpath)
	if err != nil {
		return err
	}
	log.Println("string:", tmplpath)
	payload := FastTemplateSet{
		Name:       tmplname,
	}
	log.Printf("%+v\n", payload)
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
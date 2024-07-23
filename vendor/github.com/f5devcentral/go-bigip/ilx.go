package bigip

import "fmt"

type ILXWorkspace struct {
	Name            string      `json:"name,omitempty"`
	FullPath        string      `json:"fullPath,omitempty"`
	Generation      int         `json:"generation,omitempty"`
	SelfLink        string      `json:"selfLink,omitempty"`
	NodeVersion     string      `json:"nodeVersion,omitempty"`
	StagedDirectory string      `json:"stagedDirectory,omitempty"`
	Version         string      `json:"version,omitempty"`
	Extensions      []Extension `json:"extensions,omitempty"`
	Rules           []File      `json:"rules,omitempty"`
}

type File struct {
	Name string `json:"name,omitempty"`
}

type Extension struct {
	Name  string `json:"name,omitempty"`
	Files []File `json:"files,omitempty"`
}

func (b *BigIP) GetWorkspace(name string) (*ILXWorkspace, error) {
	spc := &ILXWorkspace{}
	err, exists := b.getForEntity(spc, uriMgmt, uriTm, uriIlx, uriWorkspace, name)
	if !exists {
		return nil, fmt.Errorf("workspace does not exist: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting ILX Workspace: %w", err)
	}
	return spc, nil
}

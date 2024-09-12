package bigip

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

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
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

type Extension struct {
	Name  string `json:"name,omitempty"`
	Files []File `json:"files,omitempty"`
}

func (b *BigIP) GetWorkspace(ctx context.Context, path string) (*ILXWorkspace, error) {
	spc := &ILXWorkspace{}
	err, exists := b.getForEntity(spc, uriMgmt, uriTm, uriIlx, uriWorkspace, path)
	if !exists {
		return nil, fmt.Errorf("workspace does not exist: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting ILX Workspace: %w", err)
	}

	return spc, nil
}

func (b *BigIP) CreateWorkspace(ctx context.Context, path string) error {
	err := b.post(ILXWorkspace{Name: path}, uriMgmt, uriTm, uriIlx, uriWorkspace, "")
	if err != nil {
		return fmt.Errorf("error creating ILX Workspace: %w", err)
	}

	return nil
}

func (b *BigIP) DeleteWorkspace(ctx context.Context, name string) error {
	err := b.delete(uriMgmt, uriTm, uriIlx, uriWorkspace, name)
	if err != nil {
		return fmt.Errorf("error deleting ILX Workspace: %w", err)
	}
	return nil
}

type ExtensionConfig struct {
	Name          string `json:"name,omitempty"`
	Partition     string `json:"partition,omitempty"`
	WorkspaceName string `json:"workspaceName,omitempty"`
}

func (b *BigIP) CreateExtension(ctx context.Context, opts ExtensionConfig) error {
	err := b.post(ILXWorkspace{Name: opts.WorkspaceName}, uriMgmt, uriTm, uriIlx, uriWorkspace+"?options=extension,"+opts.Name)
	if err != nil {
		return fmt.Errorf("error creating ILX Extension: %w", err)
	}
	return nil
}

// UploadExtensionFiles uploads the files in the given directory to the BIG-IP system
// Only index.js and package.json files are uploaded as they are the only mutable files.
func (b *BigIP) UploadExtensionFiles(ctx context.Context, opts ExtensionConfig, path string) error {
	destination := fmt.Sprintf("%s/%s/%s/extensions/%s/", WORKSPACE_UPLOAD_PATH, opts.Partition, opts.WorkspaceName, opts.Name)
	files, err := readFilesFromDirectory(path)
	if err != nil {
		return err
	}
	err = b.uploadFilesToDestination(files, destination)
	if err != nil {
		return err
	}
	return nil
}

type ExtensionFile int64

const (
	PackageJSON ExtensionFile = iota
	IndexJS
)

func (e ExtensionFile) String() string {
	switch e {
	case PackageJSON:
		return "package.json"
	case IndexJS:
		return "index.js"
	}
	return "unknown"
}

func (b *BigIP) WriteExtensionFile(ctx context.Context, opts ExtensionConfig, content string, filename ExtensionFile) error {
	destination := fmt.Sprintf("%s/%s/%s/extensions/%s/%s", WORKSPACE_UPLOAD_PATH, opts.Partition, opts.WorkspaceName, opts.Name, filename)
	err := b.WriteFile(content, destination)
	if err != nil {
		return fmt.Errorf("error uploading packagejson: %w", err)
	}
	return nil
}

// Reads extension files for a specified ExtensionConfig. Ignores node_modules folder
func (b *BigIP) ReadExtensionFiles(ctx context.Context, opts ExtensionConfig) ([]File, error) {
	destination := fmt.Sprintf("%s/%s/%s/extensions/%s/", WORKSPACE_UPLOAD_PATH, opts.Partition, opts.WorkspaceName, opts.Name)
	ignored := []string{"node_modules"}
	files, err := b.getFilesFromDestination(destination, ignored)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (b *BigIP) UploadRuleFiles(ctx context.Context, opts ExtensionConfig, path string) error {
	destination := fmt.Sprintf("%s/%s/%s/rules/", WORKSPACE_UPLOAD_PATH, opts.Partition, opts.WorkspaceName)
	files, err := readFilesFromDirectory(path)
	if err != nil {
		return err
	}
	if err = b.uploadFilesToDestination(files, destination); err != nil {
		return err
	}
	return nil
}

func (b *BigIP) uploadFilesToDestination(files []*os.File, destination string) error {
	uploadedFilePaths, err := b.uploadFiles(files)
	if err != nil {
		return err
	}
	for _, uploadedFilePath := range uploadedFilePaths {
		err := b.runCatCommand(uploadedFilePath, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func (b *BigIP) getFilesFromDestination(destination string, ignoreList []string) ([]File, error) {
	files := []File{}
	command := BigipCommand{
		Command:     "run",
		UtilCmdArgs: fmt.Sprintf("-c 'ls %s'", destination),
	}
	output, err := b.RunCommand(&command)
	if err != nil {
		return nil, fmt.Errorf("error running command: %w", err)
	}
	split := strings.Split(output.CommandResult, "\n")
	split = removeEmpty(split)
	for _, line := range strings.Split(output.CommandResult, "\n") {
		// Check if the current file is in the ignore list
		if contains(ignoreList, strings.TrimSpace(line)) || strings.TrimSpace(line) == "" {
			continue
		}
		fileContentCommand := BigipCommand{
			Command:     "run",
			UtilCmdArgs: fmt.Sprintf("-c 'cat %s/%s'", destination, strings.TrimSpace(line)),
		}
		fileContent, err := b.RunCommand(&fileContentCommand)
		if err != nil {
			return nil, fmt.Errorf("error running command: %w", err)
		}
		files = append(files, File{Name: line, Content: fileContent.CommandResult})
	}
	return files, nil
}

func readFilesFromDirectory(path string) ([]*os.File, error) {
	fileDirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}
	files := []*os.File{}
	for _, fileDir := range fileDirs {
		if fileDir.IsDir() {
			continue
		}
		f, err := fileFromDirEntry(fileDir, path)
		if err != nil {
			return nil, fmt.Errorf("error getting file from directory entry: %w", err)
		}
		files = append(files, f)
	}
	return files, nil
}

func (b *BigIP) uploadFiles(files []*os.File) ([]string, error) {
	uploadedFilePaths := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "index.js") || strings.HasSuffix(file.Name(), "package.json") {
			res, err := b.UploadFile(file)
			if err != nil {
				return nil, fmt.Errorf("error uploading file: %w", err)
			}
			uploadedFilePaths = append(uploadedFilePaths, res.LocalFilePath)
		}
	}

	return uploadedFilePaths, nil
}

func (b *BigIP) runCatCommand(uploadedFilePath, destination string) error {
	fileName := filepath.Base(uploadedFilePath)
	command := BigipCommand{
		Command:     "run",
		UtilCmdArgs: fmt.Sprintf("-c 'cat %s > %s'", uploadedFilePath, destination+fileName),
	}
	_, err := b.RunCommand(&command)
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	return nil
}

func fileFromDirEntry(entry fs.DirEntry, dir string) (*os.File, error) {
	path := filepath.Join(dir, entry.Name())

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

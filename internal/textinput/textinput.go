package textinput

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"text/template"
)

//go:embed template.go.tmpl
var tempModule string

const tempNewInstanceFunc = "NewInstance"

type DynamicStructLoader struct {
	tempDir    string
	tempFile   string
	tempPlugin string
	StructText string
	StructName string
	plugin     *plugin.Plugin
	loaded     bool
}

func (d *DynamicStructLoader) Done() error {
	if d.tempPlugin != "" {
		os.RemoveAll(d.tempDir)
	}

	if d.tempFile != "" {
		os.Remove(d.tempFile)
	}

	if d.tempPlugin != "" {
		os.Remove(d.tempPlugin)
	}

	return nil
}

func (d *DynamicStructLoader) Load() error {
	tempDirPath, err := os.MkdirTemp("", "structwrap-*")
	if err != nil {
		return err
	}

	tempFilePath, err := writeTemporaryModule(tempDirPath, d.StructText, d.StructName)
	if err != nil {
		return err
	}

	tempPlugin, err := buildPlugin(tempFilePath)
	if err != nil {
		return err
	}

	plugin, err := loadPlugin(tempPlugin)
	if err != nil {
		return err
	}

	d.plugin = plugin
	d.tempPlugin = tempPlugin
	d.tempFile = tempFilePath
	d.tempDir = tempDirPath
	d.loaded = true

	return nil
}

func (d *DynamicStructLoader) NewInstance() (any, error) {
	if !d.loaded {
		return nil, fmt.Errorf("loader not loaded")
	}

	outFn, err := lookupSymbol(d.plugin, tempNewInstanceFunc)
	return outFn(), err
}

func NewDynamicStructLoader(structText, structName string) *DynamicStructLoader {
	return &DynamicStructLoader{
		StructText: structText,
		StructName: structName,
	}
}

func writeTemporaryModule(tempDir, structText, structName string) (string, error) {
	tmpFile, err := os.CreateTemp(tempDir, "structwrap-*.go")
	if err != nil {
		return "", err
	}

	defer tmpFile.Close()

	tmpl, err := template.New("wrap").Parse(tempModule)
	if err != nil {
		return "", err
	}

	input := map[string]string{
		"Structs":    structText,
		"StructName": structName,
	}

	if err = tmpl.Execute(tmpFile, input); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func buildPlugin(tempFileName string) (string, error) {
	dir := filepath.Dir(tempFileName)
	pluginPath := filepath.Join(dir, "dyn.so")

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginPath, tempFileName)
	cmd.Stdout = os.Stdout // or os.Stdout for debugging
	cmd.Stderr = os.Stderr // or os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build plugin: %w", err)
	}

	return pluginPath, nil
}

func loadPlugin(pluginPath string) (*plugin.Plugin, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	return p, nil
}

func lookupSymbol(p *plugin.Plugin, name string) (func() any, error) {
	sym, err := p.Lookup(name)
	if err != nil {
		return nil, fmt.Errorf("%s symbol not found: %w", name, err)
	}

	convertSymbol, ok := sym.(func() any)
	if !ok {
		return nil, fmt.Errorf("%s has wrong signature", name)
	}

	return convertSymbol, nil
}

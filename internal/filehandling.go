package 

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Resource struct {
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	Recommended   bool     `yaml:"recommended"`
	Level         int      `yaml:"level"`
	Categories    []string `yaml:"categories"`
	Documentation string   `yaml:"docs"`
	IsSoftware    bool     `yaml:"software"`
}

type ResourcesFile struct {
	Resources []Resource `yaml:"resources"`
}

func loadResources(filename string) (*ResourcesFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var resourcesFile ResourcesFile
	err = yaml.Unmarshal(data, &resourcesFile)
	if err != nil {
		return nil, err
	}

	return &resourcesFile, nil
}

package file

import (
	"os"

	"github.com/go-yaml/yaml"
)

// Resource struct for the contents of the yaml file
type Resource struct {
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	Recommended   bool     `yaml:"recommended"`
	Level         int      `yaml:"level"`
	Categories    []string `yaml:"categories"`
	Documentation string   `yaml:"docs"`
	IsSoftware    bool     `yaml:"software"`
}

// ResourcesFile definition of the yaml file itself
type ResourcesFile struct {
	Resources []Resource `yaml:"resources"`
}

// LoadResources reading the file and unmarshaling the yaml
func LoadResources(filename string) (*ResourcesFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var resource ResourcesFile
	err = yaml.Unmarshal(data, &resource)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

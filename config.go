package main

import (
	"bytes"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Config represents the load test configuration
type Config struct {
	LoadTest struct {
		Count int `yaml:"count"`
		Delay int `yaml:"delay"`
	} `yaml:"loadTest"`
	Randomization struct {
		NamePrefix   string `yaml:"namePrefix"`
		SuffixLength int    `yaml:"suffixLength"`
	} `yaml:"randomization"`
	Resource string                 `yaml:"resource"`
	Template map[string]interface{} `yaml:"template"`
}

// TemplateData holds variables for template rendering
type TemplateData struct {
	RandomName string
}

// LoadConfig loads configuration from YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GenerateObject creates a kubernetes object from template with randomized fields
func (c *Config) GenerateObject() (*unstructured.Unstructured, error) {
	randomName := GenerateRandomName(c.Randomization.NamePrefix, c.Randomization.SuffixLength)
	templateData := TemplateData{
		RandomName: randomName,
	}

	templateBytes, err := yaml.Marshal(c.Template)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("k8s-object").Parse(string(templateBytes))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	err = yaml.Unmarshal(buf.Bytes(), &obj)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: obj}, nil
}

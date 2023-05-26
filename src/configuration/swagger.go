package configuration

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type CSwagger struct {
	Version     string `yaml:"version"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	BasePath    string `yaml:"basepath"`
}

func (c *Configuration) LoadSwaggerConf() {
	yamlFile, err := os.ReadFile("swagger.yaml")
	if err != nil {
		log.Fatalf("error opening swagger configuration file swagger.yaml: %s", err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &c.Swagger)
	if err != nil {
		log.Fatalf("swagger configuration unmarshal failed: %s", err.Error())
	}
}

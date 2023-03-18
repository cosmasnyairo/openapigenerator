package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type OpenApiSpec struct {
	Openapiversion string `yaml:"openapi"`
	Servers        []Servers
	Paths          []OpenApiPath `yaml:"paths"`
}

type OpenApiPath struct {
	Method      string `yaml:"openapi"`
	Integration []OpenApiPathIntegration
}

type OpenApiPathIntegration struct {
	Url        string `yaml:"openapi"`
	Httpmethod string `yaml:"openapi"`
}

type Apispec struct {
	Name                  string    `yaml:"name"`
	Description           string    `yaml:"description"`
	EndpointConfiguration string    `yaml:"endpointConfiguration"`
	Servers               []Servers `yaml:"servers"`
	Corsheaders           []string  `yaml:"corsheaders"`
}

type Servers struct {
	Url string `yaml:"url"`
}

type Pathspec struct {
	Name            string   `yaml:"name,omitempty"`
	Uri             string   `yaml:"uri,omitempty"`
	Methods         []string `yaml:"methods,omitempty"`
	Cors            bool     `yaml:"cors,omitempty"`
	Queryparameters []string `yaml:"queryparameters,omitempty"`
}

type Model[T any] struct {
	Paths T `yaml:"paths"`
}

func main() {
	apispec := Apispec{}
	apispec.getApispec()

	mergePaths([]Pathspec{})
	// for _, paths := range mergedpaths {
	// 	fmt.Println(paths)
	// }

	// openapi := OpenApiSpec{
	// 	Openapiversion: "3.0.1",
	// 	Servers:        apispec.Servers,
	// }
	// yaml, err := yaml.Marshal(&openapi)
	// onError("Unable to unmarshall", err)
	// fmt.Println(custom)
}

func mergePaths(pathspec []Pathspec) []Pathspec {
	pathbyte := []byte{}
	files, err := os.ReadDir("services")
	onError("Unable to read directory", err)

	for _, v := range files {
		yamlfile, err := os.ReadFile("services/" + v.Name())
		onError("Unable to read files in directory", err)

		err = yaml.Unmarshal(yamlfile, &pathspec)
		onError("Unable to unmarshall", err)

		yaml, err := yaml.Marshal(&pathspec)
		onError("Unable to marshall", err)

		pathbyte = append(pathbyte[:], yaml[:]...)
	}

	err = yaml.Unmarshal(pathbyte, &pathspec)
	onError("Unable to unmarshall", err)

	for _, paths := range pathspec {
		paths.convertoOpenAPi()
	}

	return pathspec
}

func (path *Pathspec) convertoOpenAPi() {

}

func onError(message string, err error) {
	if err != nil {
		log.Fatalf("%v %v", message, err)
	}
}

func (a *Apispec) getApispec() {
	yamlfile, err := os.ReadFile("apispec.yaml")
	onError("Unable to openfile", err)

	err = yaml.Unmarshal(yamlfile, &a)
	onError("Unable to unmarshall", err)
}

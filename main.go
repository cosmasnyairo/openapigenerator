package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

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

type OpenAPI struct {
	OpenAPI string                 `yaml:"openapi"`
	Info    OpenApiSpecInfo        `yaml:"info"`
	Servers []Servers              `yaml:"servers,omitempty"`
	Paths   map[string]OpenApiPath `yaml:"paths,omitempty"`
	// Components OpenApiSpecComponents  `yaml:"components,omitempty"`
	// Security   []map[string][]string  `yaml:"security,omitempty"`
	// Tags       []*Tag                      `yaml:"tags,omitempty"`

}

type OpenApiSpecInfo struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`
	Contact     struct {
		Name  string `yaml:"name,omitempty"`
		Email string `yaml:"email,omitempty"`
		URL   string `yaml:"url,omitempty"`
	} `yaml:"contact,omitempty"`
}

// type OpenApiSpecComponents struct {
// 	// securitySchemes
// }

type OpenApiPath struct {
	Ref                  string               `yaml:"$ref,omitempty"`
	Summary              string               `yaml:"summary,omitempty"`
	Description          string               `yaml:"description,omitempty"`
	OpenApiPathOperation OpenApiPathOperation `yaml:"get,omitempty,flow"`

	// Parameters  []*OpenApiPathParameter `yaml:"parameters,omitempty"`
}

type OpenApiPathOperation struct {
	Tags        []string `yaml:"tags,omitempty"`
	Summary     string   `yaml:"summary,omitempty"`
	Description string   `yaml:"description,omitempty"`
	// RequestBody OpenApiPathRequestBody `yaml:"requestBody,omitempty"`
	// Responses   OpenApiPathResponses   `yaml:"responses,omitempty"`
	Deprecated         bool              `yaml:"deprecated,omitempty"`
	GatewayIntegration map[string]string `yaml:"x-amazon-apigateway-integration,omitempty"`
}

type OpenApiPathSpec struct {
	Apis struct {
		Get struct {
			MethodResponses             interface{} `yaml:"methodresponses"`
			AmazonApiGatewayIntegration struct {
				URI           string      `yaml:"uri"`
				ConnectionID  string      `yaml:"connectionId"`
				HTTPMethod    string      `yaml:"httpMethod"`
				MethodDetails interface{} `yaml:"methoddetails"`
			} `yaml:"x-amazon-apigateway-integration"`
		} `yaml:"get"`
		Options interface{} `yaml:"options"`
	}
}

// type OpenApiResource struct {
// 	Resource    OpenApiResource
// 	Method      string `yaml:"openapi"`
// 	Integration []OpenApiPathIntegration
// }

// type OpenApiPathIntegration struct {
// 	Url        string `yaml:"openapi"`
// 	Httpmethod string `yaml:"openapi"`
// }

func main() {
	apispec := Apispec{}
	apispec.getApispec()

	pathspec := mergePaths([]Pathspec{})

	generateOpenApiPaths(pathspec)

	var pathitems []map[string]any

	// for _, paths := range pathspec {
	// 	pathitems = append(pathitems, paths.generateOpenApiPaths())
	// }

	generated, _ := yaml.Marshal(&pathitems)

	err := os.WriteFile("generated-paths.yaml", generated, 0644)
	onError("Unable to write", err)

	// fmt.Println(string(generated))

	// pathspec, _ :=yaml.Marshal(mergePaths([]Pathspec{}))

	// err := os.WriteFile("generated.yaml", merged, 0644)
	// onError("Unable to write", err)

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
	return pathspec
}

func generateOpenApiPaths(pathspec []Pathspec) map[string]any {

	var methods []byte
	for _, paths := range pathspec {
		// generate operations
		for _, method := range paths.Methods {
			methodOperation := map[string]any{
				strings.ToLower(method): map[string]any{
					"responses": map[string]string{
						"test": "test",
					},
					"x-amazon-apigateway-integration": map[string]any{
						"uri":          "https://$${stageVariables.url}/apis",
						"connectionId": "$${stageVariables.vpcLink}",
						"httpMethod":   strings.ToUpper(method),
						"methoddetails": map[string]string{
							"test": "test",
						},
					},
				},
			}
			generated, _ := yaml.Marshal(&methodOperation)
			methods = append(methods[:], generated[:]...)
		}

	}

	err := os.WriteFile("generated-paths1.yaml", methods, 0644)
	onError("Unable to write", err)
	// TODO prevent conversion to a list here

	pathitem := map[string]any{
		"path.Name": string(methods),
	}

	return pathitem

}

func (a *Apispec) getApispec() {
	yamlfile, err := os.ReadFile("apispec.yaml")
	onError("Unable to openfile", err)

	err = yaml.Unmarshal(yamlfile, &a)
	onError("Unable to unmarshall", err)

}

func onError(message string, err error) {
	if err != nil {
		log.Fatalf("%v %v", message, err)
	}
}

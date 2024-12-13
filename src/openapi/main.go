package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type OpenApiTagEntity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OpenApiServer struct {
	Url string `json:"url"`
}

type OpenApiPathRequestBody struct {
	Content struct {
		Json struct {
			Schema interface{} `json:"schema"`
		} `json:"application/json"`
	} `json:"content"`
}

type OpenApiParam struct {
	Name        string  `json:"name"`
	Location    string  `json:"in"`
	Description string  `json:"description"`
	Ref         *string `json:"$ref"`
	IsRequired  bool    `json:"required"`
}

// TODO: make type assertion
// type OpenApiSchemaInlineProp struct {
// 	Description string    `json:"description"`
// 	Enum        *[]string `json:"enum"`
// 	Example     any    `json:"example"`
// 	Type        string    `json:"type"`
// }

// type OpenApiSchemaReferencedProp struct {
// 	Description string `json:"description"`
// 	Ref         string `json:"$ref"`
// }

// type OpenApiInlineSchema struct {
// 	Properties map[string]interface{} `json:"properties"`
// 	Required   *[]string              `json:"required"`
// }

// TODO: make type assertion
// type OpenApiPathDiscriminatedSchema struct {
// 	Discriminator struct {
// 		Mapping map[string]string `json:"mapping"`
// 	} `json:"discriminator"`
// 	AnyOf []struct {
// 		Ref string `json:"$ref"`
// 	} `json:"anyOf"`
// }

//	type OpenApiPathReferencedSchema struct {
//		Ref string `json:"$ref"`
//	}

type OpenApiGenericProp struct {
	Default     *any      `json:"default"`
	Description *string   `json:"description"`
	Enum        *[]string `json:"enum"`
	Example     *any      `json:"example"`
	Ref         *string   `json:"$ref"`
	Type        *string   `json:"type"`
}

type OpenApiGenericSchema struct {
	Description *string                        `json:"description"`
	Properties  *map[string]OpenApiGenericProp `json:"properties"`
	Ref         *string                        `json:"$ref"`
	Required    *[]string                      `json:"required"`
	Type        *string                        `json:"type"`
}

type OpenApiResponse struct {
	Description string `json:"description"`
	Content     *struct {
		Json struct {
			Schema OpenApiGenericSchema `json:"schema"`
		} `json:"application/json"`
	} `json:"content"`
}

type OpenApiPath struct {
	Description  string                     `json:"description"`
	IsDeprecated bool                       `json:"deprecated"`
	Operation    string                     `json:"operationId"`
	Params       []OpenApiParam             `json:"parameters"`
	RequestBody  *OpenApiPathRequestBody    `json:"requestBody"`
	Responses    map[string]OpenApiResponse `json:"responses"`
	Summary      string                     `json:"summary"`
	Tags         []string                   `json:"tags"`
}

type OpenApiComponents struct {
	Schemas *map[string]OpenApiGenericSchema `json:"schemas"`
}

type OpenApiSpecification struct {
	Components OpenApiComponents                 `json:"components"`
	Entities   []OpenApiTagEntity                `json:"tags"`
	Paths      map[string]map[string]OpenApiPath `json:"paths"`
	Servers    []OpenApiServer                   `json:"servers"`
	Version    string                            `json:"openapi"`
}

func ReadFile(path string) (*OpenApiSpecification, error) {
	data, err := os.ReadFile(path)

	if errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("Specification file not found")
	}

	var specification OpenApiSpecification

	if err := json.Unmarshal(data, &specification); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &specification, nil
}

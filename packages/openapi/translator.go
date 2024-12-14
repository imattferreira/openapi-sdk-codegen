package openapi

import (
	"strconv"
	"strings"
)

type TranslatedProp struct {
	Name        string
	IsRequired  bool
	Default     *any
	Enum        []any
	Example     any
	Type        string
	Description string
}

type TranslatedParam struct {
	Location string
	TranslatedProp
}

type TranslatedRequest struct {
	Body      []TranslatedProp
	Method    string
	Operation string
	Params    []TranslatedParam
}

type TranslatedResponse struct {
	Description string
	Body        []TranslatedProp
	Status      int64
}

type TranslatedResponses struct {
	Errors  []TranslatedResponse
	Success *TranslatedResponse
}

type TranslatedEndpoint struct {
	Description  string
	Entity       string
	IsDeprecated bool
	Path         string
	Request      TranslatedRequest
	Response     TranslatedResponses
	Summary      string
}

type TranslatedSpecification struct {
	DomainUrls []string
	Endpoints  []TranslatedEndpoint
	Entities   []string
}

type Translator struct {
	components map[string][]TranslatedProp
}

func (t *Translator) DomainUrls(servers []OpenApiServer) []string {
	var urls []string

	for _, server := range servers {
		urls = append(urls, server.Url)
	}

	return urls
}

func (t *Translator) Property(name string, schema OpenApiGenericSchema) TranslatedProp {
	var example any = nil
	schemaType := ""

	if schema.Example != nil {
		example = *schema.Example
	}

	if schema.Type != nil {
		schemaType = *schema.Type
	}

	return TranslatedProp{
		Name:        name,
		IsRequired:  false,
		Default:     nil,
		Enum:        schema.Enum,
		Example:     example,
		Type:        schemaType,
		Description: *schema.Description,
	}
}

func (t *Translator) Schema(name string, schema OpenApiGenericSchema) []TranslatedProp {
	var translated []TranslatedProp

	if schema.Type != nil && *schema.Type != "object" {
		prop := t.Property(name, schema)
		translated = append(translated, prop)

		return translated
	}

	if schema.Type != nil && *schema.Type == "object" {
		for name, propSchema := range *schema.Properties {
			isRequired := false

			for _, prop := range *schema.Required {
				if prop == name {
					isRequired = true
					break
				}
			}

			prop := t.Property(name, propSchema)
			prop.IsRequired = isRequired
			translated = append(translated, prop)
		}

		return translated
	}

	if schema.Ref != nil {
		translated = t.GetComponentFromRef(*schema.Ref)
	}

	return translated
}

func (t *Translator) GetComponentFromRef(ref string) []TranslatedProp {
	var component []TranslatedProp

	if len(ref) == 0 {
		return component
	}

	parts := strings.Split(ref, "/")
	name := parts[len(parts)-1]
	component = t.components[name]

	return component
}

func (t *Translator) Components(c OpenApiComponents) map[string][]TranslatedProp {
	components := make(map[string][]TranslatedProp)

	for name, props := range c.Schemas {
		components[name] = t.Schema(name, props)
	}

	return components
}

func (t *Translator) Param(param OpenApiParam) TranslatedParam {
	paramType := ""
	isRequired := param.IsRequired || false
	var defaultValue any = nil
	var example any = nil

	if param.Default != nil {
		defaultValue = param.Default
	}

	if param.Example != nil {
		example = param.Example
	}

	if param.Type != nil {
		paramType = *param.Type

	}

	if param.Ref != nil {
		component := t.GetComponentFromRef(*param.Ref)

		if len(component) > 0 {
			prop := component[0]
			example = prop.Example
			paramType = prop.Type

			if prop.Default != nil {
				defaultValue = prop.Default
			}
		}
	}

	return TranslatedParam{
		Location: param.Location,
		TranslatedProp: TranslatedProp{
			Name:        param.Name,
			IsRequired:  isRequired,
			Description: param.Description,
			Type:        paramType,
			Default:     &defaultValue,
			Example:     example,
		},
	}
}

func (t *Translator) Body(requestBody *OpenApiPathRequestBody) []TranslatedProp {
	var body []TranslatedProp

	if requestBody == nil {
		return body
	}

	if requestBody.Content.Json.Schema.Type != nil && *requestBody.Content.Json.Schema.Type == "object" {
		for name, propSchema := range *requestBody.Content.Json.Schema.Properties {
			isRequired := false

			for _, prop := range *requestBody.Content.Json.Schema.Required {
				if prop == name {
					isRequired = true
					break
				}
			}

			prop := t.Property(name, propSchema)
			prop.IsRequired = isRequired
			body = append(body, prop)
		}

		return body
	}

	if requestBody.Content.Json.Schema.Ref != nil {
		body = t.GetComponentFromRef(*requestBody.Content.Json.Schema.Ref)
	}

	return body
}

func (t *Translator) Response(status int64, res OpenApiResponse) TranslatedResponse {
	var body []TranslatedProp

	if res.Content == nil {
		return TranslatedResponse{
			Description: res.Description,
			Status:      status,
			Body:        body,
		}
	}

	if res.Content.Json.Schema.Type != nil && *res.Content.Json.Schema.Type == "object" {
		for name, propSchema := range *res.Content.Json.Schema.Properties {
			isRequired := false

			for _, prop := range *res.Content.Json.Schema.Required {
				if prop == name {
					isRequired = true
					break
				}
			}

			prop := t.Property(name, propSchema)
			prop.IsRequired = isRequired
			body = append(body, prop)
		}
	}

	if res.Content.Json.Schema.Ref != nil {
		body = t.GetComponentFromRef(*res.Content.Json.Schema.Ref)
	}

	return TranslatedResponse{
		Description: res.Description,
		Status:      status,
		Body:        body,
	}
}

func (t *Translator) Endpoints(paths map[string]map[string]OpenApiPath) []TranslatedEndpoint {
	var endpoints []TranslatedEndpoint

	for path, methods := range paths {
		for method, props := range methods {
			isDeprecated := props.IsDeprecated || false
			entity := camelfy(props.Tags[0])
			body := t.Body(props.RequestBody)

			var params []TranslatedParam

			for _, param := range props.Params {
				params = append(params, t.Param(param))
			}

			request := TranslatedRequest{
				Method:    method,
				Body:      body,
				Operation: props.Operation,
				Params:    params,
			}

			var errors []TranslatedResponse
			var success *TranslatedResponse = nil

			for s, res := range props.Responses {
				status, err := strconv.ParseInt(s, 10, 64)

				if err != nil {
					continue
				}

				if success == nil && status <= 200 && status < 300 {
					r := t.Response(status, res)
					success = &r
					continue
				}

				r := t.Response(status, res)
				errors = append(errors, r)
			}

			response := TranslatedResponses{
				Errors:  errors,
				Success: success,
			}

			endpoints = append(endpoints, TranslatedEndpoint{
				IsDeprecated: isDeprecated,
				Description:  props.Description,
				Entity:       entity,
				Path:         path,
				Request:      request,
				Response:     response,
				Summary:      props.Summary,
			})
		}
	}

	return endpoints
}

func (t *Translator) Entities(tags []OpenApiTagEntity) []string {
	var entities []string

	for _, tag := range tags {
		entities = append(entities, camelfy(tag.Name))
	}

	return entities
}

func (t *Translator) Translate(specification *OpenApiSpecification) TranslatedSpecification {
	t.components = t.Components(specification.Components)
	domainUrls := t.DomainUrls(specification.Servers)
	entities := t.Entities(specification.Entities)
	endpoints := t.Endpoints(specification.Paths)

	return TranslatedSpecification{
		DomainUrls: domainUrls,
		Entities:   entities,
		Endpoints:  endpoints,
	}
}

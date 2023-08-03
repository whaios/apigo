package apifox

import (
	"github.com/go-openapi/spec"
	"github.com/whaios/apigo/parser"
)

func NewOpenApi2() *spec.Swagger {
	api := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger: "2.0",
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "Apigo",
					Version:     "1.0.0",
					Description: "解析 Go 代码文件中的注释生成 Api 文档。",
				},
			},
			Host:     "",
			BasePath: "",
			Paths: &spec.Paths{
				Paths: make(map[string]spec.PathItem),
			},
		},
	}
	return api
}

func convtParamType(in string) string {
	switch in {
	case parser.ParamTypeForm:
		return "formData" // 在 Apifox 中使用
	default:
		return in
	}
}

func convtParameters(in string, items []parser.Parameter) []spec.Parameter {
	parameters := make([]spec.Parameter, 0)
	for _, item := range items {
		param := spec.Parameter{
			SimpleSchema: spec.SimpleSchema{
				Type:    item.Type,
				Example: item.Example,
			},
			ParamProps: spec.ParamProps{
				Name:        item.Name,
				In:          convtParamType(in),
				Required:    item.Required,
				Description: item.Description,
			},
		}
		parameters = append(parameters, param)
	}
	return parameters
}

func OpenApi2AddPaths(api *spec.Swagger, apiItems []parser.ApiItem) {
	for _, apiItem := range apiItems {
		parameters := make([]spec.Parameter, 0)
		{
			if len(apiItem.Parameters.Path) > 0 {
				params := convtParameters(parser.ParamTypePath, apiItem.Parameters.Path)
				parameters = append(parameters, params...)
			}
			if len(apiItem.Parameters.Query) > 0 {
				params := convtParameters(parser.ParamTypeQuery, apiItem.Parameters.Query)
				parameters = append(parameters, params...)
			}
			if len(apiItem.Parameters.Header) > 0 {
				params := convtParameters(parser.ParamTypeHeader, apiItem.Parameters.Header)
				parameters = append(parameters, params...)
			}
			if len(apiItem.Parameters.Cookie) > 0 {
				params := convtParameters(parser.ParamTypeCookie, apiItem.Parameters.Cookie)
				parameters = append(parameters, params...)
			}
			if len(apiItem.Parameters.FormData) > 0 {
				params := convtParameters(parser.ParamTypeForm, apiItem.Parameters.FormData)
				parameters = append(parameters, params...)
			}
			if apiItem.Parameters.JsonSchema != nil {
				param := spec.Parameter{
					ParamProps: spec.ParamProps{
						Name:   parser.SchemaGetTypeFullName(apiItem.Parameters.JsonSchema),
						In:     parser.ParamTypeBody,
						Schema: apiItem.Parameters.JsonSchema,
					},
				}
				parameters = append(parameters, param)
			}
		}

		statusCodeResponses := make(map[int]spec.Response)
		for _, item := range apiItem.Responses {
			resp := spec.Response{
				ResponseProps: spec.ResponseProps{
					Description: item.Name,
					Schema:      item.JsonSchema,
				},
			}
			statusCodeResponses[item.Code] = resp
		}

		consumes, produces := make([]string, 0), make([]string, 0)
		{
			if apiItem.Parameters.BodyType != "" {
				consumes = append(consumes, apiItem.Parameters.BodyType)
			}
			if len(consumes) == 0 {
				consumes = append(consumes, parser.BodyTypeNone)
			}
			if apiItem.ContentType != "" {
				produces = append(produces, apiItem.ContentType)
			}
			if len(produces) == 0 {
				produces = append(produces, parser.BodyTypeJSON)
			}
		}

		operation := &spec.Operation{
			VendorExtensible: spec.VendorExtensible{
				Extensions: map[string]interface{}{
					XFolder: apiItem.Folder,
					XStatus: apiItem.Status,
				},
			},
			OperationProps: spec.OperationProps{
				Summary:     apiItem.Title,
				Description: apiItem.Description,
				Consumes:    consumes,
				Parameters:  parameters,
				Produces:    produces,
				Responses: &spec.Responses{
					ResponsesProps: spec.ResponsesProps{
						StatusCodeResponses: statusCodeResponses,
					},
				},
			},
		}

		pathItem, ok := api.Paths.Paths[apiItem.Path]
		if !ok {
			pathItem = spec.PathItem{
				PathItemProps: spec.PathItemProps{},
			}
		}
		switch apiItem.Method {
		case parser.MethodGet:
			pathItem.PathItemProps.Get = operation
		case parser.MethodPut:
			pathItem.PathItemProps.Put = operation
		case parser.MethodPost:
			pathItem.PathItemProps.Post = operation
		case parser.MethodDelete:
			pathItem.PathItemProps.Delete = operation
		case parser.MethodOptions:
			pathItem.PathItemProps.Options = operation
		case parser.MethodHead:
			pathItem.PathItemProps.Head = operation
		case parser.MethodPatch:
			pathItem.PathItemProps.Patch = operation
		}
		api.Paths.Paths[apiItem.Path] = pathItem
	}
}

package middleware

import (
	"context"
	"io"
	"net/http"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
)

func yamlBodyDecoder(r io.Reader, h http.Header, s *openapi3.SchemaRef, fn openapi3filter.EncodingFn) (interface{}, error) {
	return "", nil
}

func OpenapiInputValidator(openApiFile string) gin.HandlerFunc {
	
	return func(c *gin.Context) {
		httpReq := c.Request
		ctx := context.Background()
		openapi3filter.RegisterBodyDecoder("text/yaml", yamlBodyDecoder)

		loader := openapi3.Loader{Context: ctx}
		doc, err := loader.LoadFromFile(openApiFile)
		if err != nil {
			fmt.Println(err.Error())
		}

		err = doc.Validate(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}

		router, err := legacyrouter.NewRouter(doc)//.WithSwaggerFromFile(openapiFile)
		if err != nil {
			fmt.Println(err.Error())
		}

		route, pathParams, _ := router.FindRoute(httpReq)

		var requestValidationInput *openapi3filter.RequestValidationInput

		// Validate Request
		requestValidationInput = &openapi3filter.RequestValidationInput{
			Request: httpReq,
			PathParams: pathParams,
			Route: route,
		}

		if erro := openapi3filter.ValidateRequest(ctx, requestValidationInput); erro != nil {
			// panic(err)
			fmt.Errorf(erro.Error())
			fmt.Println("Fail")
		}

		c.Next()
	}
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
)

func main() {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx}
	doc, _ := loader.LoadFromFile("./sriov-dp.yml")
	err := doc.Validate(ctx)
	if err != nil {
		panic(err)
	}

	router, err := legacyrouter.NewRouter(doc)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(http.MethodGet, strings.TrimSpace("http://localhost:8000/sriov-dp"), nil)
	if err != nil {
		panic(err)
	}

	log.Println(httpReq)
	// Find route
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		panic(err)
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		panic(err)
	}

	var (
		respStatus      = 200
		respContentType = "application/json"
		respBody        = bytes.NewBufferString(`{}`)
	)

	log.Println("Response:", respStatus)
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{respContentType}},
	}
	if respBody != nil {
		data, err := json.Marshal(respBody)
		if err != nil {
			panic(err)
		}
		log.Println(string(data))
		responseValidationInput.SetBodyBytes(data)
	}

	// Validate response.
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}
}

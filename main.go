package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
)

func main() {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx}
	doc, _ := loader.LoadFromFile("./sriov-cni.yml")
	err := doc.Validate(ctx)
	if err != nil {
		panic(err)
	}

	router, err := legacyrouter.NewRouter(doc)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(http.MethodGet, strings.TrimSpace("http://127.0.0.1:8000/sriov-cni"), nil)
	if err != nil {
		panic(err)
	}
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
		// respBody        = bytes.NewBufferString(`{}`)
	)

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{respContentType}},
	}

	resp, err := http.Get("http://127.0.0.1:8000/sriov-cni")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
	responseValidationInput.SetBodyBytes(output)

	// Validate response.
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}
}

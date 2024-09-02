//go:build go1.22

// Package httphandler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package httphandler

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
)

// AnalysisRequest defines model for AnalysisRequest.
type AnalysisRequest struct {
	Answer          string `json:"answer"`
	IsUserSatisfied bool   `json:"is_user_satisfied"`
	Query           string `json:"query"`
}

// AnalysisResponse defines model for AnalysisResponse.
type AnalysisResponse struct {
	Answer          string `json:"answer"`
	Id              int    `json:"id"`
	IsUserSatisfied bool   `json:"is_user_satisfied"`
	Query           string `json:"query"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetAnalyzesParams defines parameters for GetAnalyzes.
type GetAnalyzesParams struct {
	Query           *string `form:"query,omitempty" json:"query,omitempty"`
	Answer          *string `form:"answer,omitempty" json:"answer,omitempty"`
	IsUserSatisfied *bool   `form:"is_user_satisfied,omitempty" json:"is_user_satisfied,omitempty"`
	Limit           *int    `form:"limit,omitempty" json:"limit,omitempty"`
	Offset          *int    `form:"offset,omitempty" json:"offset,omitempty"`
}

// AddAnalysisJSONRequestBody defines body for AddAnalysis for application/json ContentType.
type AddAnalysisJSONRequestBody = AnalysisRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get analyzes
	// (GET /analyzes)
	GetAnalyzes(w http.ResponseWriter, r *http.Request, params GetAnalyzesParams)
	// Add analysis
	// (POST /analyzes)
	AddAnalysis(w http.ResponseWriter, r *http.Request)
	// Get analysis by id
	// (GET /analyzes/{id})
	GetAnalysisById(w http.ResponseWriter, r *http.Request, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetAnalyzes operation middleware
func (siw *ServerInterfaceWrapper) GetAnalyzes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAnalyzesParams

	// ------------- Optional query parameter "query" -------------

	err = runtime.BindQueryParameter("form", true, false, "query", r.URL.Query(), &params.Query)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "query", Err: err})
		return
	}

	// ------------- Optional query parameter "answer" -------------

	err = runtime.BindQueryParameter("form", true, false, "answer", r.URL.Query(), &params.Answer)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "answer", Err: err})
		return
	}

	// ------------- Optional query parameter "is_user_satisfied" -------------

	err = runtime.BindQueryParameter("form", true, false, "is_user_satisfied", r.URL.Query(), &params.IsUserSatisfied)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "is_user_satisfied", Err: err})
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", r.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAnalyzes(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// AddAnalysis operation middleware
func (siw *ServerInterfaceWrapper) AddAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddAnalysis(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetAnalysisById operation middleware
func (siw *ServerInterfaceWrapper) GetAnalysisById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", r.PathValue("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAnalysisById(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/analyzes", wrapper.GetAnalyzes)
	m.HandleFunc("POST "+options.BaseURL+"/analyzes", wrapper.AddAnalysis)
	m.HandleFunc("GET "+options.BaseURL+"/analyzes/{id}", wrapper.GetAnalysisById)

	return m
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xW32/bNhD+Vwhuj17ktH0o9OYMw2asHYLusQgKWjzb10okc3dyphn634eTbCWOlBpp",
	"BywPezJF3i9+991H720RqxQDBGGb7y0XW6hct1wEVzaM/AFua2DRrUQxAQlCZ+AC3wHpCv5yVSrB5vYD",
	"OG9ki2wcCRYlmK1I4jzLvtQroAACfIEx87HgbAMiGDY/sTgS8JmdWWmShmEhDBvbzizyp5qBPrET5DWC",
	"P8knVMPgtIqxBBfU67YGak4r+y3eGYlmA2IO+cwdytb8PtQ1Tt/OLMFtjaRpPx6izo4XnyruZogRV5+h",
	"EC3mHkhOMTD8V0ieQnf56vVghEFgA/TS8EZvZ98K+i9EkZ5GHPT4tGAFgYIrDQPtgExvcq7G3mpcgdph",
	"WEfN4YELwiQYg83t4nppZOvEJIo79MAmRWZcYYnSKGTOe+OC76Bzyp2/gU1cm7tIXzBsehjfvzNV9FB2",
	"KKJ0N3j/zhy5Zq774GQW10s7szsg7tNfXswv5opQTBBcQpvb193WzCYn2w6d7JhWPzbQzb6i5/QKS29z",
	"+yvI4mijjuQqECC2+ce9Rc1z7FtwFTz47AVGA45wnXYc+v5szzFRJoIMJH4qSokVypTnMDRPecb1muGM",
	"642yqSdpB/ar+Vx/ihgEQoe7S6nEokM++8zawv2DgI915L5tKFB1ix8J1ja3P2T3Sp8dZD4bSVM7ENkR",
	"uaaj+wSzTxn9Z10UwKzOb76r/smp3LkSvTliO57Gs+VdOW/o8Ih1h2tXl/KsMr+G4anSTORfTgqL2nFd",
	"VU6FU+dpmHUtMkWemLqF98eO2V6GgOUq+uZfu8zjR7891TvV/3ZE2cuxyP1M4AT8/5T4dkosumfg0G09",
	"GlQ526Nvz0ozI181S/+EPKvWP5BKbx+3+WuC+72q9TxJOiM4b8bs+yOKWcc6+Jc/7/pWrxqDvo/Uu/R9",
	"qqm0udV/fXmWlbFw5Tay5G/nb+eZS5jtLm170/4TAAD//+qyPFTCCwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

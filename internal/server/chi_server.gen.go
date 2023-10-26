// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

const (
	Ionos_tokenScopes = "ionos_token.Scopes"
)

// DBaaSQuota defines model for DBaaSQuota.
type DBaaSQuota struct {
	Limits DBaaSQuotaFields `json:"Limits"`
	Usage  DBaaSQuotaFields `json:"Usage"`
}

// DBaaSQuotaFields defines model for DBaaSQuotaFields.
type DBaaSQuotaFields struct {
	CPU              int64 `json:"CPU"`
	Memory           int64 `json:"Memory"`
	MongoClusters    int64 `json:"MongoClusters"`
	PostgresClusters int64 `json:"PostgresClusters"`
	Storage          int64 `json:"Storage"`
}

// DNSQuota defines model for DNSQuota.
type DNSQuota struct {
	Limits DNSQuotaFields `json:"Limits"`
	Usage  DNSQuotaFields `json:"Usage"`
}

// DNSQuotaFields defines model for DNSQuotaFields.
type DNSQuotaFields struct {
	Records        int64 `json:"Records"`
	SecondaryZones int64 `json:"SecondaryZones"`
	Zones          int64 `json:"Zones"`
}

// Quotas defines model for Quotas.
type Quotas struct {
	DBaaS *DBaaSQuota `json:"DBaaS,omitempty"`
	DNS   *DNSQuota   `json:"DNS,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /health)
	GetHealth(w http.ResponseWriter, r *http.Request)

	// (GET /quotas)
	GetQuotas(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /health)
func (_ Unimplemented) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /quotas)
func (_ Unimplemented) GetQuotas(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetHealth operation middleware
func (siw *ServerInterfaceWrapper) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHealth(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetQuotas operation middleware
func (siw *ServerInterfaceWrapper) GetQuotas(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, Ionos_tokenScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetQuotas(w, r)
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
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
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

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/health", wrapper.GetHealth)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/quotas", wrapper.GetQuotas)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xVTXPTMBD9K56Fo4lTUjj4VgqFUEgCocMMmRxUeWOrtbWutM6Q6fi/M5KdhHxQXMop",
	"drT79r2n9e49SCpK0qjZQnwPBm1J2qJ/eWcMGfcgSTNqdo+MPzkqc6G0e7Myw0L4/1clQgyWjdIp1HUd",
	"QoJWGlWyIg0xnE2GgRR5HqBHrUP4UhELu4cvyjJXUrik6MbSXpXnBhcQw7NoyzpqTm3Uwh0rLWVVVLlg",
	"TII7HxUsyASOjEWzVBJdWtjW8YzevhFi6hHdW2moRMOqseWTKlTj1kN0tggXCvPEOsVXVqT4+MQ6BIN3",
	"lTKYQDxrUcI1j3m4Np+ub1CyK3SAcSDifHLlfhZkCsEQg9L8+hQ2UEozpuiv6TMWZFZdg0mndJ5XltHY",
	"jjkTspwatI9MmzKZ1s6/Ru85uMvyCIPQ+7PRvi121OzREztl9G99Mnpyl4we7pGvKMkkne8DJelEmNUP",
	"0tg1qXvsnrwmMdxwPKh/TPB24uwK9Z9L9++y9a7rBXnue1zctEFZGcWrqQtveCjSZF8w3aKfe9coDJqL",
	"tTMfv3+Ddkg5qOZ061TGXDbTT+kFufzdKXiBLDO0gUhTg+nBMAxStUQdCCmp0hwsDBUBZxgMx6PxNDjP",
	"qUqCs8nQ1VOcu4KNm8G0GaEQwhKNbWr1eye9vrOJStSiVBDDoNfvDSCEUnDmxUYZipwz95giH/L94I8D",
	"maG8BY9k/FoYJhDDe+TmHMLdnfWy3++4sVBXheuk8eVvvfLn/TW+dHpeNfDHrn1DI2r2Zu1BortNzx1V",
	"+d9u5cCfttmP+/OwgPUqDeG0f9JZ7yPd2X4BEM/2en82r+d1E2KWfiXM7qEyedvmcRTlJEWekeV4MBgM",
	"oJ7XvwIAAP//A7pPIssIAAA=",
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

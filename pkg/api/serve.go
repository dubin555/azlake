package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/dubin555/azlake/pkg/api/apigen"
	"github.com/dubin555/azlake/pkg/api/apiutil"
	"github.com/dubin555/azlake/pkg/httputil"
	"github.com/dubin555/azlake/pkg/logging"
)

const (
	LoggerServiceName = "rest_api"
	extensionValidationExcludeBody = "x-validation-exclude-body"
)

func Serve(
	logger logging.Logger,
) *chi.Mux {
	logger.Info("initialize OpenAPI server")
	swagger, err := apigen.GetSwagger()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	apiRouter := r.With(
		OapiRequestValidatorWithOptions(swagger, &openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		}),
		httputil.LoggingMiddleware(
			httputil.RequestIDHeaderName,
			logging.Fields{logging.ServiceNameFieldKey: LoggerServiceName},
			"",
			false,
			false,
		),
	)
	controller := NewController()
	apigen.HandlerFromMuxWithBaseURL(controller, apiRouter, apiutil.BaseURL)

	r.Mount("/_health", httputil.ServeHealth())
	r.Mount("/metrics", promhttp.Handler())
	r.Mount("/_pprof/", httputil.ServePPROF("/_pprof/"))
	r.Mount("/openapi.json", http.HandlerFunc(swaggerSpecHandler))
	r.Mount(apiutil.BaseURL, http.HandlerFunc(InvalidAPIEndpointHandler))

	return r
}

func swaggerSpecHandler(w http.ResponseWriter, _ *http.Request) {
	reader, err := apigen.GetSwaggerSpecReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.Copy(w, reader)
}

func OapiRequestValidatorWithOptions(swagger *openapi3.T, options *openapi3filter.Options) func(http.Handler) http.Handler {
	router, err := legacy.NewRouter(swagger)
	if err != nil {
		panic(err)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, m, err := router.FindRoute(r)
			if err != nil {
				writeError(w, r, http.StatusBadRequest, err.Error())
				return
			}
			r = r.WithContext(logging.AddFields(r.Context(), logging.Fields{"operation_id": route.Operation.OperationID}))
			statusCode, err := validateRequest(r, route, m, options)
			if err != nil {
				writeError(w, r, statusCode, err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func validateRequest(r *http.Request, route *routers.Route, pathParams map[string]string, options *openapi3filter.Options) (int, error) {
	if _, ok := route.Operation.Extensions[extensionValidationExcludeBody]; ok {
		o := *options
		o.ExcludeRequestBody = true
		options = &o
	}
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
		Options:    options,
	}
	if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
		var reqErr *openapi3filter.RequestError
		if errors.As(err, &reqErr) {
			return http.StatusBadRequest, err
		}
		var seqErr *openapi3filter.SecurityRequirementsError
		if errors.As(err, &seqErr) {
			return http.StatusUnauthorized, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func InvalidAPIEndpointHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusNotFound, ErrInvalidAPIEndpoint)
}

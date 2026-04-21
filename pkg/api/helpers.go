package api

import (
	"net/http"

	"github.com/dubin555/azlake/pkg/httputil"
)

func writeError(w http.ResponseWriter, r *http.Request, code int, v any) {
	httputil.WriteAPIError(w, r, code, v)
}

func writeResponse(w http.ResponseWriter, r *http.Request, code int, response any) {
	httputil.WriteAPIResponse(w, r, code, response)
}

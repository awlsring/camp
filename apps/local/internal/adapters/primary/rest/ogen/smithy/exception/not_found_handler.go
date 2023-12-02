package exception

import (
	"net/http"
)

var _ http.HandlerFunc = UnknownOperationHandler

func UnknownOperationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(SmithyErrorTypeHeader, ExceptionTypeUnknownOperationException.String())
	w.WriteHeader(404)
	w.Write([]byte("{\"message\":\"Unknown operation\"}"))
}

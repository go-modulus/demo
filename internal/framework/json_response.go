package framework

import (
	"encoding/json"
	"net/http"
)

type JsonResponseWriter interface {
	Success(w http.ResponseWriter, r *http.Request, response ActionResponse)
	Error(w http.ResponseWriter, r *http.Request, response ActionResponse)
}

type DefaultJsonResponseWriter struct {
	logger Logger
}

func NewJsonResponseWriter(logger Logger) JsonResponseWriter {
	return &DefaultJsonResponseWriter{logger: logger}
}

func (j *DefaultJsonResponseWriter) Success(w http.ResponseWriter, r *http.Request, response ActionResponse) {
	w.WriteHeader(response.StatusCode)
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(response.Response)
	if err != nil {
		ctx := r.Context()
		j.logger.Error(ctx, "Error happened in JSON marshal. Err: %s", err)
	}
	_, _ = w.Write(jsonResp)
}

func (j *DefaultJsonResponseWriter) Error(w http.ResponseWriter, r *http.Request, response ActionResponse) {
	w.WriteHeader(response.StatusCode)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]interface{})
	resp["error"] = "Unknown error"
	if response.Error != nil {
		resp["error"] = response.Error.Error()
	}

	if len(response.Error.ValidationErrors) > 0 {
		vErrors := make([]map[string]string, len(response.Error.ValidationErrors))
		for i, validationError := range response.Error.ValidationErrors {
			vErrors[i] = map[string]string{
				"id":      string(validationError.Identifier),
				"field":   string(validationError.Field),
				"message": validationError.Err,
			}
		}
		resp["invalidInputs"] = vErrors
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		ctx := r.Context()
		j.logger.Error(ctx, "Error happened in JSON marshal. Err: %s", err)
		_, _ = w.Write([]byte(`{"error": "Error happened in JSON marshal."}`))
		return
	}
	_, _ = w.Write(jsonResp)
}

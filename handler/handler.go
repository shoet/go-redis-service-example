package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Index struct{}

func (*Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}
	RespondJSON(w, http.StatusOK, resp)
}

type NotFound struct{}

func (*NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message string `json:"message"`
	}{
		Message: ErrorMessageNotFound,
	}
	RespondJSON(w, http.StatusNotFound, resp)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

const ErrorMessageNotFound = "Not Found"
const ErrorMessageInternalServerError = "Internal Server Error"

func RespondJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	jsonBody, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResp := ErrorResponse{Message: ErrorMessageInternalServerError}
		err := json.NewEncoder(w).Encode(errResp)
		if err != nil {
			log.Printf("failed to encode error response: %v", err)
		}
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonBody)
	return
}

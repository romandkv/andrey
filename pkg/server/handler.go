package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const contextKey = Key("modified_request")

type Key string

func errorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, "{\"error\":\"%s\"}", message)
}

type FinalHandler struct{}

func (h FinalHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("Executing query")
	encoder := json.NewEncoder(w)
	encoder.Encode(req.Context().Value(contextKey))
}

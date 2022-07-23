package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type IMarshaler interface {
	Execute(http.Handler) http.Handler
}

type Marshaler struct{}

func (mutator *Marshaler) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing Marshaler")
		request := map[string]interface{}{}
		decoder := json.NewDecoder(r.Body)
		decoder.UseNumber()
		err := decoder.Decode(&request)
		if err != nil {
			log.Println(err)
			errorResponse(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), contextKey, request))
		next.ServeHTTP(w, r)
	})
}

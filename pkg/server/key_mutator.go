package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type JsonKeyMutatorMiddleware struct {
	rules *RewriteRules
}

func (mutator *JsonKeyMutatorMiddleware) SetRules(rules *RewriteRules) {
	mutator.rules = rules
}

func (mutator *JsonKeyMutatorMiddleware) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing JsonKeyMutatorMiddleware")
		if mutator.rules == nil {
			next.ServeHTTP(w, r)
			return
		}
		request := r.Context().Value(contextKey).(map[string]interface{})
		err := mutator.ReplaceByRules(request, 0)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), contextKey, request))
		next.ServeHTTP(w, r)
	})
}

func (mutator *JsonKeyMutatorMiddleware) ReplaceByRules(
	request map[string]interface{},
	level int,
) error {
	key, last := mutator.rules.GetKey(level)
	value, ok := request[key]
	if !ok {
		return fmt.Errorf("missing required key: %s", key)
	}
	if !last {
		untypedSections, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("%s should be an array", key)
		}
		for _, untypedSection := range untypedSections {
			section, ok := untypedSection.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%s should be an array", key)
			}
			/// ??? should we throw err in case of no key (just a part of slice)
			if err := mutator.ReplaceByRules(section, level+1); err != nil {
				return err
			}
		}
		return nil
	}
	request[mutator.rules.NewKey] = value
	delete(request, key)
	return nil
}

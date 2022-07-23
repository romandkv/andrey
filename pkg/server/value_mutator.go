package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RewriteRules struct {
	SourceKey string `json:"key"`
	key       []string
	OldValue  interface{} `json:"oldValue"`
	NewValue  interface{} `json:"newValue"`
	NewKey    string      `json:"newKey"`
}

func NewRewriteRules(key string, oldValue, newValue interface{}) *RewriteRules {
	rewriteRules := RewriteRules{
		OldValue: oldValue,
		NewValue: newValue,
	}
	rewriteRules.SetKey(key)
	return &rewriteRules
}

func NewRewriteRulesKey(key, newKey string) *RewriteRules {
	rewriteRules := RewriteRules{
		NewKey: newKey,
	}
	rewriteRules.SetKey(key)
	return &rewriteRules
}

func (r *RewriteRules) SetKey(key string) {
	r.key = strings.Split(key, ".")
}

func (r *RewriteRules) GetKey(level int) (key string, last bool) {
	return r.key[level], level == len(r.key)-1
}

type IJsonMutatorMiddleware interface {
	Execute(http.Handler) http.Handler
	SetRules(*RewriteRules)
}

type JsonValueMutatorMiddleware struct {
	rules *RewriteRules
}

func (mutator *JsonValueMutatorMiddleware) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing JsonValueMutatorMiddleware")
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

func (mutator *JsonValueMutatorMiddleware) ReplaceByRules(
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
	}
	if mutator.rules.OldValue != value {
		return nil
	}
	request[key] = mutator.rules.NewValue
	return nil
}

func (mutator *JsonValueMutatorMiddleware) SetRules(rules *RewriteRules) {
	mutator.rules = rules
}

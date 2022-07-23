package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const rewriteRulesPath = "config/rewrite_rules.json"

type server struct {
	port                       string
	handler                    http.Handler
	jsonValueMutatorMiddleware IJsonMutatorMiddleware
	jsonKeyMutatorMiddleware   IJsonMutatorMiddleware
	marshalMiddleware          IMarshaler
}

func NewServer(
	port string,
	jsonValueMutatorMiddleware IJsonMutatorMiddleware,
	jsonKeyMutatorMiddleware IJsonMutatorMiddleware,
	marshalMiddleware IMarshaler,
	handler http.Handler,
) *server {
	return &server{
		port:                       port,
		jsonValueMutatorMiddleware: jsonValueMutatorMiddleware,
		jsonKeyMutatorMiddleware:   jsonKeyMutatorMiddleware,
		marshalMiddleware:          marshalMiddleware,
		handler:                    handler,
	}
}

func GetRewriteRules() (*RewriteRules, error) {
	path := os.Getenv("REWRITE_RULES_PATH")
	if path == "" {
		path = rewriteRulesPath
	}
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	rules := RewriteRules{}
	err = decoder.Decode(&rules)
	if err != nil {
		return nil, err
	}
	rules.SetKey(rules.SourceKey)
	return &rules, nil
}

func GetRewriteRulesKey() (*RewriteRules, error) {
	path := os.Getenv("REWRITE_RULES_KEY_PATH")
	if path == "" {
		path = rewriteRulesPath
	}
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	rules := RewriteRules{}
	err = decoder.Decode(&rules)
	if err != nil {
		return nil, err
	}
	rules.SetKey(rules.SourceKey)
	return &rules, nil
}

func (s server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/", s.handler)
	rules, err := GetRewriteRules()
	if err != nil {
		return err
	}
	rulesKey, err := GetRewriteRules()
	if err != nil {
		return err
	}
	s.jsonValueMutatorMiddleware.SetRules(rules)
	s.jsonKeyMutatorMiddleware.SetRules(rulesKey)
	err = http.ListenAndServe(
		fmt.Sprintf(":%s", s.port),
		s.marshalMiddleware.Execute(
			s.jsonValueMutatorMiddleware.Execute(
				s.jsonKeyMutatorMiddleware.Execute(s.handler),
			),
		),
	)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	return nil
}

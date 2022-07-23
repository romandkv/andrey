package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndToEnd(t *testing.T) {
	svr := NewServer(
		"8080",
		&JsonValueMutatorMiddleware{},
		&JsonKeyMutatorMiddleware{},
		&Marshaler{},
		FinalHandler{},
	)
	os.Setenv("REWRITE_RULES_PATH", "../../config/rewrite_rules.json")
	os.Setenv("REWRITE_RULES_KEY_PATH", "../../config/rewrite_rules_key.json")
	go func() {
		svr.Run()
	}()
	time.Sleep(time.Second)
	Composite = append(Composite, append(Tests, TestsKey...)...)
	for _, tst := range Composite {
		svr.jsonValueMutatorMiddleware.SetRules(tst.rules)
		svr.jsonKeyMutatorMiddleware.SetRules(tst.rulesKey)
		t.Run(tst.description, func(t *testing.T) {
			resp, err := http.Post("http://127.0.0.1:8080", "application/json", strings.NewReader(tst.jsonString))
			require.NoError(t, err)
			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tst.err != nil {
				assert.Equal(t, fmt.Sprintf("{\"error\":\"%s\"}", tst.err.Error()), string(bodyBytes))
				return
			}
			assert.JSONEq(t, tst.expectJsonString, string(bodyBytes))
		})
	}
}

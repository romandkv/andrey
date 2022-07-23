package server

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Composite = []struct {
	description      string
	rules            *RewriteRules
	rulesKey         *RewriteRules
	jsonString       string
	expectJsonString string
	err              error
}{
	{
		rulesKey: nil,
		rules:    nil,
		jsonString: `{
			"foo":123,
			"bar":"bar"
		}`,
		expectJsonString: `{
			"foo":123,
			"bar":"bar"
		}`,
	},
	{
		rulesKey: NewRewriteRulesKey(
			"foo",
			"gg",
		),
		rules: NewRewriteRules(
			"foo",
			json.Number("123"),
			"yyy",
		),
		jsonString: `{
			"foo":123,
			"bar":"bar"
		}`,
		expectJsonString: `{
			"gg":"yyy",
			"bar":"bar"
		}`,
	},
}

var TestsKey = []struct {
	description      string
	rules            *RewriteRules
	rulesKey         *RewriteRules
	jsonString       string
	expectJsonString string
	err              error
}{
	{
		rulesKey: NewRewriteRulesKey(
			"foo",
			"bar",
		),
		jsonString:       `{}`,
		expectJsonString: `{}`,
		err:              errors.New("missing required key: foo"),
	},
	{
		rulesKey: NewRewriteRulesKey(
			"foo",
			"gg",
		),
		jsonString: `{
			"foo":123,
			"bar":"bar"
		}`,
		expectJsonString: `{
			"gg":123,
			"bar":"bar"
		}`,
	},
	{
		rulesKey: NewRewriteRulesKey(
			"foo",
			"gg",
		),
		jsonString: `{
			"rr":123,
			"bar":"bar"
		}`,
		expectJsonString: `{
			"rr":123,
			"bar":"bar"
		}`,
		err: errors.New("missing required key: foo"),
	},
	{
		rulesKey: NewRewriteRulesKey(
			"foo",
			"gg",
		),
		jsonString: `{
			"rr":123,
			"bar":"bar"
		}`,
		expectJsonString: `{
			"rr":123,
			"bar":"bar"
		}`,
		err: errors.New("missing required key: foo"),
	},
	{
		rulesKey: NewRewriteRulesKey(
			"foo.bar.q",
			"bar",
		),
		jsonString: `{
			"foo":[
				{
					"bar":[]
				},
				{
					"bar":[
						{
							"q":1
						}
					]
				}
			]
		}`,
		expectJsonString: `{
			"foo":[
				{
					"bar":[]
				},
				{
					"bar":[
						{
							"bar":1
						}
					]
				}
			]
		}`,
	},
	{
		rulesKey: NewRewriteRulesKey(
			"foo.bar",
			"gg",
		),
		jsonString: `{
			"foo":[
				{
					"bar":123
				},
				{
					"bar":"not 123"
				}
			]
		}`,
		expectJsonString: `{
			"foo":[
				{
					"gg":123
				},
				{
					"gg":"not 123"
				}
			]
		}`,
	},
}

var Tests = []struct {
	description      string
	rules            *RewriteRules
	rulesKey         *RewriteRules
	jsonString       string
	expectJsonString string
	err              error
}{
	{
		rules: NewRewriteRules(
			"foo",
			json.Number("123"),
			"bar",
		),
		jsonString:       `{}`,
		expectJsonString: `{}`,
		err:              errors.New("missing required key: foo"),
	},
	{
		rules: NewRewriteRules(
			"foo",
			json.Number("123"),
			"bar",
		),
		jsonString: `{
			"foo":123,
			"bar":"bar"
		}`,
		expectJsonString: `{
			"foo":"bar",
			"bar":"bar"
		}`,
	},

	{
		rules: NewRewriteRules(
			"foo",
			json.Number("1"),
			"bar",
		),
		jsonString: `{
			"foo":123,
			"bar":"bar"
		}`,
		expectJsonString: `{
			"foo":123,
			"bar":"bar"
		}`,
	},
	{
		/// ??? should we throw err in case of no key (just a part of slice)
		rules: NewRewriteRules(
			"foo.bar",
			json.Number("123"),
			"bar",
		),
		jsonString: `{
			"foo":[
				{
					"qwe":123
				},
				{
					"bar":"not 123"
				}
			]
		}`,
		expectJsonString: `{
			"foo":[
				{
					"qwe":123
				},
				{
					"bar":"not 123"
				}
			]
		}`,
		err: errors.New("missing required key: bar"),
	},
	{
		rules: NewRewriteRules(
			"foo.bar",
			json.Number("123"),
			"bar",
		),
		jsonString: `{
			"foo":[
				{
					"bar":123
				},
				{
					"bar":"not 123"
				}
			]
		}`,
		expectJsonString: `{
			"foo":[
				{
					"bar":"bar"
				},
				{
					"bar":"not 123"
				}
			]
		}`,
	},
	{
		rules: NewRewriteRules(
			"foo.bar",
			"not 123",
			"bar",
		),
		jsonString: `{
			"foo":[
				{
					"bar":123
				},
				{
					"bar":"not 123"
				}
			]
		}`,
		expectJsonString: `{
			"foo":[
				{
					"bar":123
				},
				{
					"bar":"bar"
				}
			]
		}`,
	},
	{
		rules: NewRewriteRules(
			"foo.bar",
			json.Number("not"),
			"bar",
		),
		jsonString: `{
			"foo":[
				{
					"bar":123
				},
				{
					"bar":"not 123"
				}
			]
		}`,
		expectJsonString: `{
			"foo":[
				{
					"bar":123
				},
				{
					"bar":"not 123"
				}
			]
		}`,
	},
	{
		rules: NewRewriteRules(
			"foo.bar.q",
			json.Number("1"),
			"bar",
		),
		jsonString: `{
			"foo":[
				{
					"bar":[]
				},
				{
					"bar":[
						{
							"q":1
						}
					]
				}
			]
		}`,
		expectJsonString: `{
			"foo":[
				{
					"bar":[]
				},
				{
					"bar":[
						{
							"q":"bar"
						}
					]
				}
			]
		}`,
	},
}

func Test_ReplaceByRules(t *testing.T) {
	mutator := &JsonValueMutatorMiddleware{}
	for _, test := range Tests {
		mutator.SetRules(test.rules)
		request := map[string]interface{}{}
		decoder := json.NewDecoder(strings.NewReader(test.jsonString))
		decoder.UseNumber()
		assert.NoError(t, decoder.Decode(&request))
		t.Run(test.description, func(t *testing.T) {
			assert.Equal(t, test.err, mutator.ReplaceByRules(request, 0))
			if test.err != nil {
				return
			}
			data, err := json.Marshal(request)
			assert.NoError(t, err)
			test.expectJsonString = strings.Replace(test.expectJsonString, "\t", "", -1)
			test.expectJsonString = strings.Replace(test.expectJsonString, "\n", "", -1)
			assert.JSONEq(t, test.expectJsonString, string(data))
		})
	}
}

func Test_ReplaceByRulesKey(t *testing.T) {
	mutator := &JsonKeyMutatorMiddleware{}
	for _, test := range TestsKey {
		mutator.SetRules(test.rulesKey)
		request := map[string]interface{}{}
		decoder := json.NewDecoder(strings.NewReader(test.jsonString))
		decoder.UseNumber()
		assert.NoError(t, decoder.Decode(&request))
		t.Run(test.description, func(t *testing.T) {
			assert.Equal(t, test.err, mutator.ReplaceByRules(request, 0))
			if test.err != nil {
				return
			}
			data, err := json.Marshal(request)
			assert.NoError(t, err)
			test.expectJsonString = strings.Replace(test.expectJsonString, "\t", "", -1)
			test.expectJsonString = strings.Replace(test.expectJsonString, "\n", "", -1)
			assert.JSONEq(t, test.expectJsonString, string(data))
		})
	}
}

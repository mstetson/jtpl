package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tcs := []struct {
		name    string
		useText bool
		tpl     string
		data    string
		out     string
		err     bool
	}{
		{
			name: "empty",
		},
		{
			name: "t1",
			tpl:  `<a>{{range .}}<b c="{{.}}"/>{{end}}</a>`,
			data: `[1234, "<>&\""]`,
			out:  `<a><b c="1234"/><b c="&lt;&gt;&amp;&#34;"/></a>`,
		},
		{
			name:    "csv",
			useText: true,
			tpl: `
				{{- range . -}}
					{{- csv .a}},{{csv .b}}{{"\n"}}
				{{- end -}}
			`,
			data: `[{"a":1,"b":{"foo":"bar"}},{"a":",\"<>"}]`,
			out:  "1,map[foo:bar]\n\",\"\"<>\",<nil>\n",
		},
	}
	for _, tc := range tcs {
		var b strings.Builder
		err := run(&b, tc.tpl, strings.NewReader(tc.data), tc.useText)
		if err != nil {
			if !tc.err {
				t.Errorf("%s: unexpected error: %v", tc.name, err)
			}
			continue
		}
		if tc.err {
			t.Errorf("%s: unexpected success; output: %q", tc.name, b.String())
			continue
		}
		if b.String() != tc.out {
			t.Errorf("%s: mismatched output %q != %q", tc.name, b.String(), tc.out)
		}
	}
}

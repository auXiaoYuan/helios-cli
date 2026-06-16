package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRootCmd_HasDomains(t *testing.T) {
	root := NewRootCmd()
	want := map[string]bool{"backend": false, "ai": false, "finance": false}
	for _, sub := range root.Commands() {
		if _, ok := want[sub.Name()]; ok {
			want[sub.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected domain %q to be registered", name)
		}
	}
}

func TestRootCmd_DispatchesSubcommand(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"backend mocktest", []string{"backend", "mocktest"}, "backend.mocktest"},
		{"backend devcode", []string{"backend", "devcode"}, "backend.devcode"},
		{"ai sealtoken", []string{"ai", "sealtoken"}, "ai.sealtoken"},
		{"ai deepresearch", []string{"ai", "deepresearch"}, "ai.deepresearch"},
		{"finance buy", []string{"finance", "buy"}, "finance.buy"},
		{"finance sell", []string{"finance", "sell"}, "finance.sell"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := NewRootCmd()
			buf := &bytes.Buffer{}
			root.SetOut(buf)
			root.SetErr(buf)
			root.SetArgs(tc.args)
			if err := root.Execute(); err != nil {
				t.Fatalf("execute %v: %v", tc.args, err)
			}
			if !strings.Contains(buf.String(), tc.want) {
				t.Errorf("expected output to contain %q, got %q", tc.want, buf.String())
			}
		})
	}
}

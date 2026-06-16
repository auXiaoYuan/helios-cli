package finance

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewFinanceCmd_HasSubcommands(t *testing.T) {
	cmd := NewFinanceCmd()
	if cmd.Use != "finance" {
		t.Fatalf("expected Use=finance, got %q", cmd.Use)
	}
	want := map[string]bool{"buy": false, "sell": false}
	for _, sub := range cmd.Commands() {
		if _, ok := want[sub.Name()]; ok {
			want[sub.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestBuyCmd_Run(t *testing.T) {
	cmd := NewBuyCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "finance.buy") {
		t.Errorf("output missing buy marker, got %q", buf.String())
	}
}

func TestSellCmd_Run(t *testing.T) {
	cmd := NewSellCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "finance.sell") {
		t.Errorf("output missing sell marker, got %q", buf.String())
	}
}

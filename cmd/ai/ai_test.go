package ai

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewAICmd_HasSubcommands(t *testing.T) {
	cmd := NewAICmd()
	if cmd.Use != "ai" {
		t.Fatalf("expected Use=ai, got %q", cmd.Use)
	}
	want := map[string]bool{"sealtoken": false, "deepresearch": false}
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

func TestSealTokenCmd_Run(t *testing.T) {
	cmd := NewSealTokenCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ai.sealtoken") {
		t.Errorf("output missing sealtoken marker, got %q", buf.String())
	}
}

func TestDeepResearchCmd_Run(t *testing.T) {
	cmd := NewDeepResearchCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ai.deepresearch") {
		t.Errorf("output missing deepresearch marker, got %q", buf.String())
	}
}

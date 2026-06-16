package backend

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewBackendCmd_HasSubcommands(t *testing.T) {
	cmd := NewBackendCmd()
	if cmd.Use != "backend" {
		t.Fatalf("expected Use=backend, got %q", cmd.Use)
	}
	want := map[string]bool{"mocktest": false, "devcode": false}
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

func TestMockTestCmd_Run(t *testing.T) {
	cmd := NewMockTestCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "backend.mocktest") {
		t.Errorf("output missing mocktest marker, got %q", buf.String())
	}
}

func TestDevCodeCmd_Run(t *testing.T) {
	cmd := NewDevCodeCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "backend.devcode") {
		t.Errorf("output missing devcode marker, got %q", buf.String())
	}
}

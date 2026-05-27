package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAuthCloneURL(t *testing.T) {
	t.Parallel()

	m := NewManager("/tmp/repos", "https://plydot.dev", "glpat-test-token")

	got := m.authCloneURL("https://plydot.dev/group/project.git")
	want := "https://oauth2:glpat-test-token@plydot.dev/group/project.git"
	if got != want {
		t.Fatalf("authCloneURL() = %q, want %q", got, want)
	}

	empty := NewManager("/tmp/repos", "https://plydot.dev", "")
	if got := empty.authCloneURL("https://plydot.dev/group/project.git"); got != "https://plydot.dev/group/project.git" {
		t.Fatalf("authCloneURL() without token = %q", got)
	}
}

func TestGitEnvUsesBasicAuth(t *testing.T) {
	t.Parallel()

	m := NewManager("/tmp/repos", "https://plydot.dev", "secret-token")
	env := m.GitEnv()

	var headerValue string
	for _, entry := range env {
		if entry == "GIT_CONFIG_VALUE_0=Authorization: Basic b2F1dGgyOnNlY3JldC10b2tlbg==" {
			headerValue = entry
		}
	}
	if headerValue == "" {
		t.Fatalf("GitEnv() missing basic auth header, got env: %v", env)
	}
}

func TestRepoRecloneReason(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, repoPath string)
		reason string
	}{
		{
			name:   "missing git dir",
			setup:  func(_ *testing.T, _ string) {},
			reason: "missing_git_dir",
		},
		{
			name: "normal clone",
			setup: func(t *testing.T, repoPath string) {
				t.Helper()
				writeGitConfig(t, repoPath, `[remote "origin"]
	url = https://example.com/group/project.git
	fetch = +refs/heads/*:refs/remotes/origin/*
`)
			},
			reason: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoPath := t.TempDir()
			tt.setup(t, repoPath)

			if got := repoRecloneReason(repoPath); got != tt.reason {
				t.Fatalf("repoRecloneReason() = %q, want %q", got, tt.reason)
			}
		})
	}
}

func writeGitConfig(t *testing.T, repoPath, config string) {
	t.Helper()

	gitDir := filepath.Join(repoPath, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", gitDir, err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte(config), 0o644); err != nil {
		t.Fatalf("WriteFile(config): %v", err)
	}
}

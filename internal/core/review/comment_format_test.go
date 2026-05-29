package review

import (
	"testing"

	"github.com/antlss/gitlab-review-agent/internal/domain"
)

func TestFormatComment(t *testing.T) {
	t.Parallel()

	c := &domain.ParsedComment{
		FilePath:      "auth/src/session-store.js",
		LineNumber:    67,
		ReviewComment: "Confirmed: The usedHandoffs Set in markHandoffUsed can grow without bound.",
		Category:      domain.CategoryPerformance,
		Severity:      domain.SeverityHigh,
		Suggestion:    "this should not appear",
	}

	got := FormatComment(c, "setTimeout(() => this.usedHandoffs.delete(code), 5 * 60 * 1000);")
	want := "File: `auth/src/session-store.js:67`\n" +
		"snippet:\n" +
		"```javascript\n" +
		"setTimeout(() => this.usedHandoffs.delete(code), 5 * 60 * 1000);\n" +
		"```\n\n" +
		"The usedHandoffs Set in markHandoffUsed can grow without bound."

	if got != want {
		t.Fatalf("FormatComment() =\n%q\nwant:\n%q", got, want)
	}
}

func TestExtractSnippetAtLine(t *testing.T) {
	t.Parallel()

	content := []byte("line one\nline two\nline three\n")
	if got := extractSnippetAtLine(content, 2); got != "line two" {
		t.Fatalf("extractSnippetAtLine() = %q", got)
	}
	if got := extractSnippetAtLine(content, 99); got != "" {
		t.Fatalf("extractSnippetAtLine(out of range) = %q", got)
	}
}

func TestNormalizeReviewComment(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"Confirmed: issue text": "issue text",
		"confirmed: issue text": "issue text",
		"Plain issue text":      "Plain issue text",
	}
	for input, want := range cases {
		if got := normalizeReviewComment(input); got != want {
			t.Fatalf("normalizeReviewComment(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestSnippetLanguage(t *testing.T) {
	t.Parallel()

	if got := snippetLanguage("auth/src/session-store.js"); got != "javascript" {
		t.Fatalf("snippetLanguage(.js) = %q", got)
	}
	if got := snippetLanguage("internal/core/review/pipeline.go"); got != "go" {
		t.Fatalf("snippetLanguage(.go) = %q", got)
	}
}

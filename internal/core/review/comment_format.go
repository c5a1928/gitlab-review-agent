package review

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/antlss/gitlab-review-agent/internal/domain"
)

type fileReaderAtSHA interface {
	ReadFileAtSHA(ctx context.Context, projectID int64, sha, filePath string) ([]byte, error)
}

var confirmedPrefixRE = regexp.MustCompile(`(?i)^confirmed:\s*`)

func SnippetAtLine(ctx context.Context, reader fileReaderAtSHA, projectID int64, sha, filePath string, lineNum int) string {
	if reader == nil || sha == "" || filePath == "" || lineNum <= 0 {
		return ""
	}
	content, err := reader.ReadFileAtSHA(ctx, projectID, sha, filePath)
	if err != nil {
		return ""
	}
	return extractSnippetAtLine(content, lineNum)
}

func extractSnippetAtLine(content []byte, lineNum int) string {
	lines := strings.Split(string(content), "\n")
	if lineNum > len(lines) {
		return ""
	}
	return strings.TrimRight(lines[lineNum-1], "\r")
}

func snippetLanguage(filePath string) string {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".js", ".mjs", ".cjs":
		return "javascript"
	case ".jsx":
		return "jsx"
	case ".ts":
		return "typescript"
	case ".tsx":
		return "tsx"
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".rb":
		return "ruby"
	case ".java":
		return "java"
	case ".kt", ".kts":
		return "kotlin"
	case ".rs":
		return "rust"
	case ".php":
		return "php"
	case ".cs":
		return "csharp"
	case ".cpp", ".cc", ".cxx", ".hpp":
		return "cpp"
	case ".c", ".h":
		return "c"
	case ".swift":
		return "swift"
	case ".sql":
		return "sql"
	case ".sh", ".bash":
		return "bash"
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	case ".vue":
		return "vue"
	case ".scss":
		return "scss"
	case ".css":
		return "css"
	case ".html", ".htm":
		return "html"
	default:
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")
		if ext == "" {
			return "text"
		}
		return ext
	}
}

func normalizeReviewComment(comment string) string {
	return strings.TrimSpace(confirmedPrefixRE.ReplaceAllString(strings.TrimSpace(comment), ""))
}

func FormatComment(c *domain.ParsedComment, snippet string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "File: `%s:%d`\n", c.FilePath, c.LineNumber)

	snippet = strings.TrimSpace(snippet)
	if snippet != "" {
		sb.WriteString("snippet:\n")
		fmt.Fprintf(&sb, "```%s\n%s\n```\n\n", snippetLanguage(c.FilePath), snippet)
	} else {
		sb.WriteByte('\n')
	}

	sb.WriteString(normalizeReviewComment(c.ReviewComment))
	return sb.String()
}

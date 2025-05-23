//go:generate go run ./cmd/gen_errors/gen.go
package errorz

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode"
)

var (
	titleCache sync.Map // map[string]string
)

var ErrInvalidDomainName = errors.New("domain name must be non-empty and alphanumeric")

// GenerateMarkdownContent builds Markdown content for a given domain and its errors.
func GenerateMarkdownContent(domain string, errors map[string]ErrorDefinition) (string, error) {
	if strings.TrimSpace(domain) == "" || strings.ContainsAny(domain, " ./\\") {
		return "", ErrInvalidDomainName
	}

	if len(errors) == 0 {
		return "", ErrLenErrors
	}

	// Sort error codes alphabetically for consistent ordering
	var codes []string
	for code := range errors {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	var builder strings.Builder
	// Estimate rough capacity: header + rows + details (~300 bytes per error)
	builder.Grow(500 + len(codes)*300)

	builder.WriteString(NormalizeMarkdownTitle(domain))

	// Write Markdown header
	builder.WriteString("| Code | Msg | HTTP | Category | Severity | Retryable |\n")
	builder.WriteString("|:------:|:-----:|:------:|:----------:|:----------:|:-----------:|\n")

	// Write each error row
	for _, code := range codes {
		errDef := errors[code]
		builder.WriteString(fmt.Sprintf(
			"| %s | %s | %d | %s | %s | %t |\n",
			errDef.Code,
			EscapeMarkdownInline(errDef.Msg),
			errDef.HTTPStatus,
			errDef.Category,
			errDef.Severity,
			errDef.IsRetryable,
		))
	}

	builder.WriteString("\n---\n\n")

	// Full error details
	for _, code := range codes {
		errDef := errors[code]
		builder.WriteString(fmt.Sprintf("## %s\n\n", code))
		builder.WriteString(fmt.Sprintf("- **Message**: %s\n", EscapeMarkdownBlock(errDef.Msg)))
		builder.WriteString(fmt.Sprintf("- **Cause**: %s\n", EscapeMarkdownBlock(errDef.Cause)))
		builder.WriteString(fmt.Sprintf("- **Solution**: %s\n", EscapeMarkdownBlock(errDef.Solution)))
		builder.WriteString(fmt.Sprintf("- **HTTP Status**: %d\n", errDef.HTTPStatus))
		builder.WriteString(fmt.Sprintf("- **Category**: %s\n", errDef.Category))
		builder.WriteString(fmt.Sprintf("- **Severity**: %s\n", errDef.Severity))
		builder.WriteString(fmt.Sprintf("- **Retryable**: %t\n", errDef.IsRetryable))

		if len(errDef.Tags) > 0 {
			builder.WriteString("- **Tags**: ")
			builder.WriteString("`" + strings.Join(errDef.Tags, "`, `") + "`\n\n")
		}

	}

	output := builder.String()
	output = strings.TrimRight(output, "\n") + "\n"
	return output, nil

}

// EscapeMarkdownInline escapes Markdown inline content (e.g., table cells)
func EscapeMarkdownInline(text string) string {
	return strings.ReplaceAll(text, "|", "\\|")
}

// EscapeMarkdownBlock escapes Markdown block content (e.g., details)
func EscapeMarkdownBlock(text string) string {
	return strings.ReplaceAll(text, "`", "\\`")
}

// NormalizeMarkdownTitle formats the domain into a Markdown header with each
// hyphen-separated part capitalized (first letter uppercase, rest lowercase).
// Results are cached for improved performance and concurrency safety.
func NormalizeMarkdownTitle(domain string) string {
	if cached, ok := titleCache.Load(domain); ok {
		return cached.(string)
	}

	parts := strings.Split(domain, "-")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}

		p = strings.ToLower(p)
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}

	result := "# " + strings.Join(parts, "-") + " Errors\n\n"
	titleCache.Store(domain, result)

	return result
}

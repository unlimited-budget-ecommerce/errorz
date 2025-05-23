//go:generate go run ./cmd/gen_errors/gen.go
package errorz

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var ErrLenErrors = errors.New("no error definitions provided")

// GenerateGoContent generates the Go code content from error definitions.
func GenerateGoContent(errors map[string]ErrorDefinition) (string, error) {
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

	// Header
	builder.WriteString("// Code generated by errorz/gen_go.go; DO NOT EDIT.\n")
	builder.WriteString("package errorz\n\n")

	// Error struct definition
	builder.WriteString("// Error represents a centralized error definition.\n")
	builder.WriteString("type Error struct {\n")
	builder.WriteString("\tCode        string\n")
	builder.WriteString("\tMsg         string\n")
	builder.WriteString("\tCause       string\n")
	builder.WriteString("\tHTTPStatus  int\n")
	builder.WriteString("\tCategory    string\n")
	builder.WriteString("\tSeverity    string\n")
	builder.WriteString("\tIsRetryable bool\n")
	builder.WriteString("\tSolution    string\n")
	builder.WriteString("\tTags        []string\n")
	builder.WriteString("}\n\n")

	// Individual error variables
	builder.WriteString("var (\n")
	for _, code := range codes {
		errDef := errors[code]
		builder.WriteString(fmt.Sprintf("\t%s = &Error{\n", code))
		builder.WriteString(fmt.Sprintf("\t\tCode: \"%s\",\n", Escape(errDef.Code)))
		builder.WriteString(fmt.Sprintf("\t\tMsg: \"%s\",\n", Escape(errDef.Msg)))
		builder.WriteString(fmt.Sprintf("\t\tCause: \"%s\",\n", Escape(errDef.Cause)))
		builder.WriteString(fmt.Sprintf("\t\tHTTPStatus: %d,\n", errDef.HTTPStatus))
		builder.WriteString(fmt.Sprintf("\t\tCategory: \"%s\",\n", Escape(errDef.Category)))
		builder.WriteString(fmt.Sprintf("\t\tSeverity: \"%s\",\n", Escape(errDef.Severity)))
		builder.WriteString(fmt.Sprintf("\t\tIsRetryable: %t,\n", errDef.IsRetryable))
		builder.WriteString(fmt.Sprintf("\t\tSolution: \"%s\",\n", Escape(errDef.Solution)))
		builder.WriteString("\t\tTags: []string{\n")
		for _, tag := range errDef.Tags {
			builder.WriteString(fmt.Sprintf("\t\t\t\"%s\",\n", Escape(tag)))
		}
		builder.WriteString("\t\t},\n")
		builder.WriteString("\t}\n")
	}
	builder.WriteString(")\n\n")

	// ErrorMap maps error code to Error object for fast lookup
	builder.WriteString("var ErrorMap = map[string]*Error{\n")
	for _, code := range codes {
		builder.WriteString(fmt.Sprintf("\t\"%s\": %s,\n", code, code))
	}
	builder.WriteString("}\n")

	return builder.String(), nil
}

func Escape(s string) string {
	replacer := strings.NewReplacer(
		`"`, `\"`,
		`\`, `\\`,
		"\n", `\n`,
	)

	return replacer.Replace(s)
}

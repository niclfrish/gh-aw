package cli

import (
	"strings"

	"github.com/github/gh-aw/pkg/logger"
)

var rateLimitFieldsCodemodLog = logger.New("cli:codemod_rate_limit_fields")

// getRateLimitFieldsCodemod creates a codemod that renames deprecated
// rate-limit.max and rate-limit.window fields to max-runs and max-runs-window.
func getRateLimitFieldsCodemod() Codemod {
	return Codemod{
		ID:           "rate-limit-fields-migration",
		Name:         "Rename rate-limit.max and rate-limit.window",
		Description:  "Renames deprecated rate-limit fields: max -> max-runs, window -> max-runs-window.",
		IntroducedIn: "1.0.0",
		Apply: func(content string, frontmatter map[string]any) (string, bool, error) {
			if !hasDeprecatedRateLimitFields(frontmatter) {
				return content, false, nil
			}
			newContent, applied, err := applyFrontmatterLineTransform(content, renameRateLimitFields)
			if applied {
				rateLimitFieldsCodemodLog.Print("Renamed deprecated rate-limit fields")
			}
			return newContent, applied, err
		},
	}
}

func hasDeprecatedRateLimitFields(frontmatter map[string]any) bool {
	rateLimitValue, hasRateLimit := frontmatter["rate-limit"]
	if !hasRateLimit {
		return false
	}

	rateLimitMap, ok := rateLimitValue.(map[string]any)
	if !ok {
		return false
	}

	_, hasMax := rateLimitMap["max"]
	_, hasMaxRuns := rateLimitMap["max-runs"]
	_, hasWindow := rateLimitMap["window"]
	_, hasMaxRunsWindow := rateLimitMap["max-runs-window"]

	return (hasMax && !hasMaxRuns) || (hasWindow && !hasMaxRunsWindow)
}

func renameRateLimitFields(lines []string) ([]string, bool) {
	result := make([]string, 0, len(lines))
	modified := false

	inRateLimit := false
	var rateLimitIndent string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			result = append(result, line)
			continue
		}

		if !strings.HasPrefix(trimmed, "#") && inRateLimit && hasExitedBlock(line, rateLimitIndent) {
			inRateLimit = false
		}

		if strings.HasPrefix(trimmed, "rate-limit:") {
			inRateLimit = true
			rateLimitIndent = getIndentation(line)
			result = append(result, line)
			continue
		}

		if inRateLimit {
			lineIndent := getIndentation(line)
			if isDescendant(lineIndent, rateLimitIndent) && strings.HasPrefix(trimmed, "max:") {
				newLine, replaced := findAndReplaceInLine(line, "max", "max-runs")
				if replaced {
					result = append(result, newLine)
					modified = true
					rateLimitFieldsCodemodLog.Printf("Renamed rate-limit.max on line %d", i+1)
					continue
				}
			}

			if isDescendant(lineIndent, rateLimitIndent) && strings.HasPrefix(trimmed, "window:") {
				newLine, replaced := findAndReplaceInLine(line, "window", "max-runs-window")
				if replaced {
					result = append(result, newLine)
					modified = true
					rateLimitFieldsCodemodLog.Printf("Renamed rate-limit.window on line %d", i+1)
					continue
				}
			}
		}

		result = append(result, line)
	}

	return result, modified
}

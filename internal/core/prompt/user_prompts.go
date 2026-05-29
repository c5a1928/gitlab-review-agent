package prompt

import "fmt"

// ─── Reviewer User Message Components ───────────────────────────────────────────

const ReviewerAllDiffsHeader = "## All Changed Files — Diffs Pre-loaded\n"
const ReviewerAllDiffsFooter = "\nAll diffs are above. Do NOT call `get_multi_diff` for these files — use tools only to read context/dependency files.\n\n"
const ReviewerAllDiffsInstruction = "Begin your investigation. Analyze the pre-loaded diffs immediately."

const ReviewerHighRiskDiffsHeader = "## High-Risk Files — Diffs Pre-loaded\n"
const ReviewerHighRiskDiffsFooter = "\nHigh-risk diffs are above. Use `get_multi_diff` for any remaining files not shown.\n\n"
const ReviewerHighRiskDiffsInstruction = "Begin your investigation. Analyze the pre-loaded diffs first, then use tools for remaining files."

const ReviewerNoDiffsInstruction = "\nBegin your investigation. Start with HIGH RISK files using `get_multi_diff`."

// ─── Agent Budget Messages ──────────────────────────────────────────────────────

func BudgetWarning(iteration, max int) string {
	return fmt.Sprintf(
		"BUDGET WARNING: You have used %d/%d tool calls. "+
			"STOP exploring and emit '=== FINAL REVIEW ===' NOW with all findings collected so far. "+
			"Only make 1 more tool call if you have a HIGH-confidence lead on a critical/high severity bug.",
		iteration, max,
	)
}

func BudgetExhausted(max int) string {
	return fmt.Sprintf(
		"BUDGET EXHAUSTED. You have used all %d allowed tool calls. "+
			"Emit '=== FINAL REVIEW ===' immediately followed by your JSON output.",
		max,
	)
}

const AgentNudge = "Please either make a tool call to gather more information, or emit '=== FINAL REVIEW ===' followed by your JSON review output if you are ready."

// ─── External-Facing Localized Text ─────────────────────────────────────────────
// These strings appear in GitLab comments/replies and are language-dependent.

func SuggestionLabel(lang ResponseLanguage) string {
	switch lang {
	case LangVI:
		return "Gợi ý sửa:"
	case LangJA:
		return "修正案:"
	default:
		return "Suggested fix:"
	}
}

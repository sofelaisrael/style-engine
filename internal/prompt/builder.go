package prompt

import (
	"fmt"
	"strings"

	"github.com/sofelaisrael/style-engine/internal/styles"
)

func BuildSystemPrompt(profile *styles.Profile, intensity float64) string {
	var b strings.Builder

	b.WriteString("You are a text transformation engine.\n\n")
	b.WriteString("Your task is to preserve the meaning of the original text while transforming its expression into the specified style.\n\n")
	b.WriteString(fmt.Sprintf("STYLE: %s\n", profile.Name))
	if profile.Description != "" {
		b.WriteString(fmt.Sprintf("DESCRIPTION: %s\n", profile.Description))
	}
	if profile.Period != "" {
		b.WriteString(fmt.Sprintf("PERIOD: %s\n", profile.Period))
	}

	intensityDesc := "moderate"
	switch {
	case intensity <= 0.3:
		intensityDesc = "subtle \u2014 use light touches of the style, keep most of the original structure"
	case intensity <= 0.6:
		intensityDesc = "moderate \u2014 blend the style with the original, clear stylistic changes"
	case intensity <= 0.8:
		intensityDesc = "strong \u2014 the style should dominate, only preserve core meaning"
	default:
		intensityDesc = "extreme \u2014 fully commit to the style, maximum stylistic transformation"
	}
	b.WriteString(fmt.Sprintf("\nINTENSITY: %s (%.0f%%)\n", intensityDesc, intensity*100))

	b.WriteString("\nSTYLE CHARACTERISTICS:\n")
	for _, c := range profile.Characteristics {
		b.WriteString(fmt.Sprintf("- %s\n", c))
	}

	if len(profile.Tone) > 0 {
		b.WriteString(fmt.Sprintf("\nTONE: %s\n", strings.Join(profile.Tone, ", ")))
	}
	if len(profile.Rhetoric) > 0 {
		b.WriteString(fmt.Sprintf("RHETORICAL DEVICES: %s\n", strings.Join(profile.Rhetoric, ", ")))
	}

	if len(profile.Examples) > 0 {
		b.WriteString("\nSTYLE EXAMPLES (follow these patterns):\n")
		for _, ex := range profile.Examples {
			b.WriteString(fmt.Sprintf("Input: %s\nOutput: %s\n\n", ex.Input, ex.Output))
		}
	}

	b.WriteString("RULES:\n")
	b.WriteString("- Preserve the original meaning exactly\n")
	b.WriteString("- Do not invent events, emotions, or details not in the original\n")
	b.WriteString("- Do not turn a simple sentence into a completely different story\n")
	b.WriteString("- Return ONLY the transformed text, nothing else\n")
	b.WriteString("- No explanations, no quotes around the output, just the transformed text\n")

	return b.String()
}

func BuildUserPrompt(text string) string {
	return fmt.Sprintf("Transform this text:\n\n\"%s\"", text)
}

func BuildRetryPrompt(original, previousOutput string, score float64, profile *styles.Profile) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("The previous output scored %.0f%% compliance. Improve it.\n\n", score*100))
	b.WriteString(fmt.Sprintf("Previous output: %s\n\n", previousOutput))
	b.WriteString("Issues to fix:\n")

	outputLower := strings.ToLower(previousOutput)

	if len(profile.Vocabulary.Pronouns) > 0 {
		hasArchaic := false
		for _, p := range profile.Vocabulary.Pronouns {
			if strings.Contains(outputLower, p) {
				hasArchaic = true
				break
			}
		}
		if !hasArchaic {
			b.WriteString("- Use the style's characteristic vocabulary (pronouns, verb forms)\n")
		}
	}

	if len(profile.Tone) > 0 {
		b.WriteString(fmt.Sprintf("- Make the tone more %s\n", strings.Join(profile.Tone, ", ")))
	}

	if profile.Syntax.PreferredSentenceComplexity == "high" && len(strings.Fields(previousOutput)) < 8 {
		b.WriteString("- Use longer, more complex sentence structures\n")
	}

	b.WriteString(fmt.Sprintf("\nOriginal text: \"%s\"\n", original))
	b.WriteString("\nReturn ONLY the improved transformed text.")

	return b.String()
}

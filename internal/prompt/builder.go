package prompt

import (
	"fmt"
	"strings"

	"github.com/sofelaisrael/style-engine/internal/styles"
)

func BuildSystemPrompt(profile *styles.Profile) string {
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

	return b.String()
}

func BuildUserPrompt(text string) string {
	return fmt.Sprintf("Transform this text:\n\n\"%s\"", text)
}

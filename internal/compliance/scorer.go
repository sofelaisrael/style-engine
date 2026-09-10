package compliance

import (
	"math"
	"strings"

	"github.com/sofelaisrael/style-engine/internal/styles"
)

func Score(output string, profile *styles.Profile) float64 {
	outputLower := strings.ToLower(output)
	score := 0.0
	weights := 0.0

	vocabScore := scoreVocabulary(outputLower, profile)
	score += vocabScore * 3.0
	weights += 3.0

	toneScore := scoreTone(outputLower, profile)
	score += toneScore * 2.0
	weights += 2.0

	lengthScore := scoreLength(output, profile)
	score += lengthScore * 1.0
	weights += 1.0

	if weights > 0 {
		return math.Round((score/weights)*100) / 100
	}
	return 0
}

func scoreVocabulary(output string, profile *styles.Profile) float64 {
	if len(profile.Vocabulary.Pronouns) == 0 && len(profile.Vocabulary.Verbs) == 0 && len(profile.Vocabulary.ArchaicWords) == 0 && len(profile.Vocabulary.Exclamations) == 0 {
		return 0.7
	}

	totalPreferred := 0
	found := 0

	for _, w := range profile.Vocabulary.Pronouns {
		totalPreferred++
		if strings.Contains(output, w) {
			found++
		}
	}
	for _, w := range profile.Vocabulary.Verbs {
		totalPreferred++
		if strings.Contains(output, w) {
			found++
		}
	}
	for _, w := range profile.Vocabulary.ArchaicWords {
		totalPreferred++
		if strings.Contains(output, w) {
			found++
		}
	}
	for _, w := range profile.Vocabulary.Exclamations {
		totalPreferred++
		if strings.Contains(output, w) {
			found++
		}
	}
	for _, w := range profile.Vocabulary.NauticalWords {
		totalPreferred++
		if strings.Contains(output, w) {
			found++
		}
	}
	for _, w := range profile.Vocabulary.Greetings {
		totalPreferred++
		if strings.Contains(output, w) {
			found++
		}
	}

	if totalPreferred == 0 {
		return 0.7
	}

	ratio := float64(found) / float64(totalPreferred)
	return math.Min(ratio*2, 1.0)
}

func scoreTone(output string, profile *styles.Profile) float64 {
	if len(profile.Tone) == 0 {
		return 0.7
	}

	toneKeywords := map[string][]string{
		"dramatic":    {"alas", "verily", "forsooth", "hath", "doth", "nay", "ay", "prithee"},
		"theatrical":  {"o ", "how ", "what ", "surely", "indeed"},
		"aggressive":  {"arr", "blimey", "avast", "shiver", "plunder"},
		"informal":    {"yer", "ye", "matey", "fer", "be "},
		"expressive":  {"!", "how ", "what ", "surely"},
		"boastful":    {"finest", "greatest", "never", "always", "mightiest"},
		"reckless":    {"damn", "hell", "blast", "curse"},
	}

	matched := 0
	for _, tone := range profile.Tone {
		if keywords, ok := toneKeywords[tone]; ok {
			for _, kw := range keywords {
				if strings.Contains(output, kw) {
					matched++
					break
				}
			}
		}
	}

	if len(profile.Tone) == 0 {
		return 0.7
	}

	return math.Min(float64(matched)/float64(len(profile.Tone))*2, 1.0)
}

func scoreLength(output string, profile *styles.Profile) float64 {
	words := strings.Fields(output)
	wordCount := len(words)

	switch profile.Syntax.PreferredSentenceComplexity {
	case "high":
		if wordCount >= 8 {
			return 1.0
		}
		return float64(wordCount) / 8.0
	case "medium":
		if wordCount >= 5 {
			return 1.0
		}
		return float64(wordCount) / 5.0
	default:
		return 0.7
	}
}
package compliance

import (
	"math"
	"strings"
	"unicode"

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

	structureScore := scoreStructure(outputLower, profile)
	score += structureScore * 1.5
	weights += 1.5

	lengthScore := scoreLength(output, profile)
	score += lengthScore * 1.0
	weights += 1.0

	if weights > 0 {
		return math.Round((score/weights)*100) / 100
	}
	return 0
}

func scoreVocabulary(output string, profile *styles.Profile) float64 {
	totalPreferred := 0
	found := 0

	allWords := append(profile.Vocabulary.Pronouns, profile.Vocabulary.Verbs...)
	allWords = append(allWords, profile.Vocabulary.ArchaicWords...)
	allWords = append(allWords, profile.Vocabulary.Exclamations...)
	allWords = append(allWords, profile.Vocabulary.NauticalWords...)
	allWords = append(allWords, profile.Vocabulary.Greetings...)

	for _, w := range allWords {
		if w == "" {
			continue
		}
		totalPreferred++
		if strings.Contains(output, w) {
			found++
		}
	}

	if totalPreferred == 0 {
		return 0.5
	}

	baseRatio := float64(found) / float64(totalPreferred)

	patternBonus := 0.0
	if strings.Contains(output, "mine ") || strings.Contains(output, "mine,") {
		patternBonus += 0.15
	}
	if hasArchaicEndings(output) {
		patternBonus += 0.15
	}
	if strings.Contains(output, "hath") || strings.Contains(output, "doth") || strings.Contains(output, "dost") || strings.Contains(output, "art ") {
		patternBonus += 0.1
	}
	if strings.Contains(output, "comest") || strings.Contains(output, "dost") || strings.Contains(output, "wilt") || strings.Contains(output, "shalt") {
		patternBonus += 0.1
	}

	result := baseRatio*0.6 + patternBonus + 0.25
	return math.Min(result, 1.0)
}

func hasArchaicEndings(output string) bool {
	archaicEndings := []string{"eth", "est", "eth.", "est."}
	for _, ending := range archaicEndings {
		if strings.HasSuffix(output, ending) || strings.Contains(output, ending+" ") {
			return true
		}
	}
	words := strings.Fields(output)
	for _, w := range words {
		w = strings.Trim(w, ".,!?;:\"'")
		if len(w) > 3 {
			suffix := w[len(w)-3:]
			if suffix == "eth" || suffix == "est" {
				return true
			}
		}
	}
	return false
}

func scoreTone(output string, profile *styles.Profile) float64 {
	if len(profile.Tone) == 0 {
		return 0.5
	}

	tonePatterns := map[string][]string{
		"formal":      {"indeed", "certainly", "assuredly", "truly", "most ", "quite ", "very ", "indubitably", "verily"},
		"elaborate":   {"considerable", "inconsiderable", "indubitably", "assiduous", "perchance", "wherefore", "henceforth", "notwithstanding", "aforementioned"},
		"restrained":  {", however,", "nevertheless", "notwithstanding", "whilst", "although", "perhaps", "it seems", "one might"},
		"moral":       {"virtue", "conscience", "duty", "honour", "righteous", "prudence", "propriety", "conscientious"},
		"vague":       {"approximately", "various", "several", "some", "multiple", "potentially", "arguably", "generally"},
		"bureaucratic": {"pursuant", "herein", "thereof", "therein", "hereby", "henceforth", "whereas", "accordingly"},
		"dramatic":    {"!", "alas", "how ", "what ", "surely", "nay ", "cannot", "must not", "shall not"},
		"unfiltered":  {"literally", "honestly", "no cap", "fr fr", "deadass", "lowkey", "highkey", "bestie"},
		"chaotic":     {"bruh", "wait what", "no way", "i'm dead", "i can't", "why is", "what even"},
		"enthusiastic": {"great", "amazing", "love", "excited", "awesome", "fantastic", "perfect", "incredible"},
		"professional": {"furthermore", "therefore", "however", "consequently", "regarding", "pursuant"},
		"casual":      {"hey", "yeah", "gonna", "wanna", "kinda", "pretty much", "honestly", "lol"},
		"dry":         {"presumably", "allegedly", "apparently", "somewhat", "rather", "apparently", "so to speak"},
		"theatrical":  {"o ", "how ", "what ", "surely", "indeed", "most ", "thus "},
		"aggressive":  {"arr", "blimey", "avast", "shiver", "plunder", "damn", "curse"},
		"informal":    {"yer", "ye", "matey", "fer ", "be ", "ain't", "gonna", "wanna"},
		"expressive":  {"!", "how ", "what ", "surely", "most ", "very "},
		"boastful":    {"finest", "greatest", "never ", "always ", "mightiest", "supreme"},
		"reckless":    {"damn", "hell", "blast", "curse", "devil"},
	}

	matched := 0
	for _, tone := range profile.Tone {
		if patterns, ok := tonePatterns[tone]; ok {
			for _, p := range patterns {
				if strings.Contains(output, p) {
					matched++
					break
				}
			}
		}
	}

	if len(profile.Tone) == 0 {
		return 0.5
	}

	ratio := float64(matched) / float64(len(profile.Tone))
	return math.Min(ratio*1.5+0.3, 1.0)
}

func scoreStructure(output string, profile *styles.Profile) float64 {
	score := 0.5

	words := strings.Fields(output)
	wordCount := len(words)

	sentences := splitSentences(output)
	avgSentenceLen := 0.0
	if len(sentences) > 0 {
		for _, s := range sentences {
			avgSentenceLen += float64(len(strings.Fields(s)))
		}
		avgSentenceLen /= float64(len(sentences))
	}

	switch profile.Syntax.PreferredSentenceComplexity {
	case "high":
		if avgSentenceLen >= 8 {
			score += 0.3
		} else if avgSentenceLen >= 5 {
			score += 0.15
		}
		if wordCount >= 8 {
			score += 0.1
		}
	case "medium":
		if avgSentenceLen >= 5 {
			score += 0.2
		}
		if wordCount >= 5 {
			score += 0.1
		}
	}

	if profile.Syntax.AllowInversion {
		inversionMarkers := []string{"thus ", "hence ", "wherefore", "thither", "hither", "so too", "neither ", "nor "}
		for _, m := range inversionMarkers {
			if strings.Contains(output, m) {
				score += 0.1
				break
			}
		}
	}

	if profile.Syntax.UseSemicolons && strings.Contains(output, ";") {
		score += 0.05
	}
	if profile.Syntax.UseEmDashes && (strings.Contains(output, "—") || strings.Contains(output, " - ")) {
		score += 0.05
	}

	return math.Min(score, 1.0)
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
		if wordCount >= 3 {
			return 1.0
		}
		return float64(wordCount) / 3.0
	}
}

func splitSentences(text string) []string {
	var sentences []string
	current := strings.Builder{}

	for _, r := range text {
		current.WriteRune(r)
		if r == '.' || r == '!' || r == '?' {
			s := strings.TrimSpace(current.String())
			if s != "" {
				sentences = append(sentences, s)
			}
			current.Reset()
		}
	}

	if current.Len() > 0 {
		s := strings.TrimSpace(current.String())
		if s != "" {
			sentences = append(sentences, s)
		}
	}

	return sentences
}

func countWords(text string) int {
	return len(strings.Fields(text))
}

func isPunct(r rune) bool {
	return unicode.IsPunct(r)
}

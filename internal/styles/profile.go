package styles

type Vocabulary struct {
	Pronouns      []string `json:"pronouns,omitempty"`
	Verbs         []string `json:"verbs,omitempty"`
	ArchaicWords  []string `json:"archaic_words,omitempty"`
	Greetings     []string `json:"greetings,omitempty"`
	Exclamations  []string `json:"exclamations,omitempty"`
	NauticalWords []string `json:"nautical_words,omitempty"`
}

type Syntax struct {
	AllowInversion              bool   `json:"allow_inversion"`
	PreferredSentenceComplexity string `json:"preferred_sentence_complexity"`
	UseSemicolons               bool   `json:"use_semicolons"`
	UseEmDashes                 bool   `json:"use_em_dashes"`
}

type Example struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Profile struct {
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Period          string     `json:"period,omitempty"`
	Vocabulary      Vocabulary `json:"vocabulary"`
	Syntax          Syntax     `json:"syntax"`
	Rhetoric        []string   `json:"rhetoric"`
	Tone            []string   `json:"tone"`
	Characteristics []string   `json:"characteristics"`
	Examples        []Example  `json:"examples"`
}

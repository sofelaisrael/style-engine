package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sofelaisrael/style-engine/internal/compliance"
	"github.com/sofelaisrael/style-engine/internal/llm"
	"github.com/sofelaisrael/style-engine/internal/prompt"
	"github.com/sofelaisrael/style-engine/internal/styles"
)

type TransformRequest struct {
	Text      string  `json:"text"`
	Style     string  `json:"style"`
	Intensity float64 `json:"intensity"`
	MaxWords  int     `json:"max_words,omitempty"`
}

type TransformResponse struct {
	Original    string  `json:"original"`
	Transformed string  `json:"transformed"`
	Style       string  `json:"style"`
	Intensity   float64 `json:"intensity"`
	Score       float64 `json:"compliance_score"`
	Retries     int     `json:"retries"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.MarshalIndent(v, "", "  ")
	w.Write(data)
	w.Write([]byte("\n"))
}

func HandleTransform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		writeError(w, "text is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Style) == "" {
		writeError(w, "style is required", http.StatusBadRequest)
		return
	}
	if req.Intensity <= 0 || req.Intensity > 1 {
		req.Intensity = 0.7
	}

	profile, err := styles.LoadProfile(req.Style)
	if err != nil {
		writeError(w, fmt.Sprintf("style not found: %s", req.Style), http.StatusBadRequest)
		return
	}

	client := llm.NewClient()
	systemPrompt := prompt.BuildSystemPrompt(profile, req.Intensity, req.MaxWords)
	userPrompt := prompt.BuildUserPrompt(req.Text)

	var result string
	var score float64
	var retries int

	inputWords := len(strings.Fields(req.Text))
	threshold := 0.6
	if inputWords <= 30 {
		threshold = 0.35
	}

	for attempt := 0; attempt < 3; attempt++ {
		llmMaxTokens := 1024
		if req.MaxWords > 0 {
			llmMaxTokens = req.MaxWords * 2
		}
		result, err = client.Transform(systemPrompt, userPrompt, llmMaxTokens)
		if err != nil {
			writeError(w, fmt.Sprintf("LLM error: %v", err), http.StatusInternalServerError)
			return
		}

		score = compliance.Score(result, profile)
		retries = attempt

		if score >= threshold {
			break
		}

		userPrompt = prompt.BuildRetryPrompt(req.Text, result, score, profile, req.MaxWords)
	}

	resp := TransformResponse{
		Original:    req.Text,
		Transformed: result,
		Style:       req.Style,
		Intensity:   req.Intensity,
		Score:       score,
		Retries:     retries,
	}

	writeJSON(w, resp)
}

func HandleListStyles(w http.ResponseWriter, r *http.Request) {
	names, err := styles.ListStyles()
	if err != nil {
		writeError(w, "could not list styles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string][]string{"styles": names})
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "internal/api/index.html")
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data, _ := json.MarshalIndent(ErrorResponse{Error: msg}, "", "  ")
	w.Write(data)
	w.Write([]byte("\n"))
}

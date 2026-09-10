package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sofelaisrael/style-engine/internal/llm"
	"github.com/sofelaisrael/style-engine/internal/prompt"
	"github.com/sofelaisrael/style-engine/internal/styles"
)

func main() {
	styleFlag := flag.String("style", "", "Style to apply (e.g. shakespeare, pirate)")
	listFlag := flag.Bool("list", false, "List available styles")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: style [text] --style <name>\n\n")
		fmt.Fprintf(os.Stderr, "Transform text using a style profile.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  --style string    Style to apply\n")
		fmt.Fprintf(os.Stderr, "  --list            List available styles\n")
		fmt.Fprintf(os.Stderr, "  --help            Show this help\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, '  style "I finished the project." --style shakespeare\n')
		fmt.Fprintf(os.Stderr, '  style "Are you coming?" --style pirate\n')
		fmt.Fprintf(os.Stderr, "  style --list\n")
	}
	flag.Parse()

	if *listFlag {
		names, err := styles.ListStyles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Available styles:")
		for _, n := range names {
			fmt.Printf("  - %s\n", n)
		}
		return
	}

	if *styleFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --style is required\n")
		flag.Usage()
		os.Exit(1)
	}

	text := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if text == "" {
		fmt.Fprintf(os.Stderr, "Error: provide text to transform\n")
		flag.Usage()
		os.Exit(1)
	}

	profile, err := styles.LoadProfile(*styleFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	systemPrompt := prompt.BuildSystemPrompt(profile)
	userPrompt := prompt.BuildUserPrompt(text)

	client := llm.NewClient()
	result, err := client.Transform(systemPrompt, userPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
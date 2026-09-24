package pipeline

import (
	"cicd/model"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func ValidatePipeline(p model.Pipeline) error {
	if p.Version <= 0 {
		return errors.New("version must be positive")
	}
	for name, job := range p.Pipeline {
		fmt.Printf("%q", name)
		// if !isValidJobName(name) {
		// 	return fmt.Errorf("invalid job name: %s", name)
		// }
		if job.Name != "docker" && job.Name != "go" && job.Name != "kubernetes" {
			return fmt.Errorf("invalid job Name: %s", job.Name)
		}
		if job.SupplyChain != nil && job.SupplyChain.Enabled {
			if job.SupplyChain.TokenRef == "" {
				return errors.New("token_ref required when supply_chain is enabled")
			}
		}
		// Validate 'needs' references exist in pipeline
		for _, need := range job.Needs {
			if _, exists := p.Pipeline[need]; !exists {
				return fmt.Errorf("job %s needs non-existent job: %s", name, need)
			}
		}
	}
	return nil
}

func isValidJobName(name string) bool {
	// pattern := `^[a-zA-Z]+$` // Regular expression for matching only alphabets
	if name == "build" || name == "lint" || name == "test" || name == "deploy" {

		return true
	}
	return false
	// Trim whitespace from the input string
	// cleanInput := TrimQuotesAndHiddenSpaces(name)
	// fmt.Printf("Job name => %q", cleanInput)

	// Validate if the cleaned input matches the pattern
	// re, _ := regexp.Compile(pattern)
	// matched := re.MatchString(cleanInput)

	// return matched
}

func TrimQuotesAndHiddenSpaces(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		// 1. Remove double quotes (and optionally single quotes/backticks)
		if r == '"' || r == '\'' || r == '`' {
			return true
		}
		// 2. Remove standard & Unicode whitespace (spaces, tabs, \n, non-breaking space \u00A0)
		if unicode.IsSpace(r) {
			return true
		}
		// 3. Remove zero-width spaces (\u200B), BOM (\uFEFF), or control characters
		if r == '\u200B' || r == '\uFEFF' || unicode.IsControl(r) {
			return true
		}

		return false
	})
}

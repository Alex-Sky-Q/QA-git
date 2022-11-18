package greeter

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func Greet(name string, hour int) string {
	greeting := "Good night"

	// Bug 2. Unclear requirements, but possible unhandled inputs
	if hour >= 24 || hour < 0 {
		return "Hour is out of range, should be between 0 and 23"
	}

	if hour >= 6 && hour < 12 {
		greeting = "Good morning"
		// Bug 1. 18 should not be included (< instead of <=)
	} else if hour >= 12 && hour < 18 {
		greeting = "Hello"
	} else if hour >= 18 && hour < 22 {
		greeting = "Good evening"
	}
	trimmedName := strings.Trim(name, " ")

	// Bug 3. Unclear requirements, but possible unhandled several language specific rules
	if strings.Contains(trimmedName, " da ") {
		nameParts := strings.Split(trimmedName, " ")
		var tempParts []string
		for _, p := range nameParts {
			if p != "da" {
				tempParts = append(tempParts, cases.Title(language.Italian).String(p))
			} else {
				tempParts = append(tempParts, p)
			}
		}
		return fmt.Sprintf("%s %s!", greeting, strings.Join(tempParts, " "))
	}

	return fmt.Sprintf("%s %s!", greeting, strings.Title(trimmedName))
}

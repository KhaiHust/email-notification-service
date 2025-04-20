package utils

import "regexp"

func ExtractVariablesBySection(subject, body string) map[string][]string {
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)

	getUniqueMatches := func(text string) []string {
		matches := re.FindAllStringSubmatch(text, -1)
		set := make(map[string]bool)
		var result []string
		for _, match := range matches {
			v := match[1]
			if !set[v] {
				set[v] = true
				result = append(result, v)
			}
		}
		return result
	}

	return map[string][]string{
		"subject": getUniqueMatches(subject),
		"body":    getUniqueMatches(body),
	}
}

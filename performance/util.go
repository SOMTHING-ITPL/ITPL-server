package performance

import "strings"

func ParsingCast(cast string) ([]string, error) {
	cleaned := strings.ReplaceAll(cast, " ", "")

	if strings.HasSuffix(cleaned, "등") {
		cleaned = strings.TrimSuffix(cleaned, "등")
	}

	if cleaned == "" {
		return []string{}, nil
	}

	parts := strings.Split(cleaned, ",")
	return parts, nil
}

func ParsingKeyword(keyword string) ([]string, error) {
	if keyword == "" {
		return []string{}, nil
	}
	parts := strings.Split(keyword, "|")
	return parts, nil
}

package greeter

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var driveFileIDPatterns = []*regexp.Regexp{
	regexp.MustCompile(`[?&]id=([a-zA-Z0-9_-]+)`),
	regexp.MustCompile(`/file/d/([a-zA-Z0-9_-]+)`),
	regexp.MustCompile(`/d/([a-zA-Z0-9_-]+)`),
}

func ExtractDriveFileID(url string) string {
	for _, p := range driveFileIDPatterns {
		if m := p.FindStringSubmatch(url); len(m) >= 2 {
			return m[1]
		}
	}
	return ""
}

func NormalizeEmailBase(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func NormalizeSurnameBase(lastName string) string {
	s := strings.ToLower(strings.TrimSpace(lastName))

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	out, _, err := transform.String(t, s)
	if err != nil {
		out = s
	}

	var b strings.Builder
	for _, r := range out {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-':
			b.WriteRune('_')
		}
	}

	result := b.String()
	re := regexp.MustCompile(`_+`)
	result = re.ReplaceAllString(result, "_")
	result = strings.Trim(result, "_")
	return result
}

func NormalizeManifestPhotoURL(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}

	if id := ExtractDriveFileID(trimmed); id != "" {
		return "https://drive.google.com/uc?export=view&id=" + id
	}

	return trimmed
}

func NormalizeManifestSourceURL(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}

	if id := ExtractDriveFileID(trimmed); id != "" {
		return "https://drive.google.com/uc?export=download&id=" + id
	}

	return trimmed
}

package bearkeek

import (
	"regexp"
	"strings"
)

var (
	compTagRegex = regexp.MustCompile(`-?#[^\s][^#]+[^\s-]#`)
	tagRegex     = regexp.MustCompile(`-?#[\w/]+`)
	splitRegex   = regexp.MustCompile(`\s+`)
)

type ParseResult struct {
	Tags      []MatchingTag
	Terms     []string
	IsTagLast bool
	LastTag   string
	RawButTag string
}

func Parse(s string) ParseResult {
	spaceIndex := strings.LastIndex(s, " ")
	hashIndex := strings.LastIndex(s, "#")

	res := ParseResult{IsTagLast: hashIndex > spaceIndex}

	if res.IsTagLast {
		res.LastTag = s[hashIndex+1:]
		res.RawButTag = s[:hashIndex]
	}

	compTags := compTagRegex.FindAllString(s, -1)
	s = compTagRegex.ReplaceAllString(s, "")

	tags := tagRegex.FindAllString(s, -1)
	s = tagRegex.ReplaceAllString(s, "")

	savetags := make([]MatchingTag, 0, len(compTags)+len(tags))

	var (
		exclude bool
		shift   int
		cut     int
	)

	for _, tag := range append(tags, compTags...) {
		shift = 1
		cut = len(tag)
		exclude = false

		if tag[0] == '-' {
			exclude = true
			shift = 2
		}

		if tag[len(tag)-1] == '#' {
			cut = len(tag) - 1
		}

		savetags = append(savetags, MatchingTag{Name: tag[shift:cut], Exclude: exclude})
	}

	res.Tags = savetags
	if len(s) > 0 {
		res.Terms = splitRegex.Split(strings.TrimSpace(s), -1)
	}

	return res
}

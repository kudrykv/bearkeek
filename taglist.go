package bearkeek

import "strings"

type taglist []string

func (t taglist) exact(in []string, exclude bool) bool {
	for _, lookupTag := range in {
		hit := false

		for _, ourTag := range t {
			if lookupTag == ourTag {
				hit = true

				continue
			}

			if strings.Contains(ourTag, lookupTag) && len(ourTag) > len(lookupTag) {
				hit = false

				break
			}
		}

		if !exclude && !hit {
			return false
		}

		if exclude && hit {
			return false
		}
	}

	return true
}

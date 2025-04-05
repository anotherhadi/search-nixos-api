package indexer

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/anotherhadi/search-nixos-api/indexer/darwin"
	"github.com/anotherhadi/search-nixos-api/indexer/homemanager"
	"github.com/anotherhadi/search-nixos-api/indexer/nixos"
	"github.com/anotherhadi/search-nixos-api/indexer/nixpkgs"
	"github.com/anotherhadi/search-nixos-api/indexer/nur"
)

func DeleteNonMatchingItems(keys Keys, pattern string) Keys {
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Err: ", err)
		return keys
	}

	var filtered Keys
	for _, k := range keys {
		if re.MatchString(k.Name) {
			filtered = append(filtered, k)
		}
	}
	return filtered
}

func (keys Keys) Search(
	query string,
	exclude []string,
) (results Keys) {
	query = strings.TrimSpace(query)
	if query == "" {
		return keys
	}

	results = slices.Clone(keys)

	var patterns []string
	if !slices.Contains(exclude, "nixpkgs") {
		patterns = append(patterns, `^`+nixpkgs.Prefix+`.*`)
	}
	if !slices.Contains(exclude, "nixos") {
		patterns = append(patterns, `^`+nixos.Prefix+`.*`)
	}
	if !slices.Contains(exclude, "homemanager") {
		patterns = append(patterns, `^`+homemanager.Prefix+`.*`)
	}
	if !slices.Contains(exclude, "darwin") {
		patterns = append(patterns, `^`+darwin.Prefix+`.*`)
	}
	if !slices.Contains(exclude, "nur") {
		patterns = append(patterns, `^`+nur.Prefix+`.*`)
	}
	if len(patterns) > 0 {
		pattern := strings.Join(patterns, "|")
		results = DeleteNonMatchingItems(results, pattern)
	} else {
		return Keys{}
	}

	for _, term := range strings.Fields(query) {
		regex := `(?i)`
		endWith := ``
		if strings.HasPrefix(term, "^") {
			term = strings.TrimPrefix(term, "^")
			regex += "^"
		}
		if strings.HasSuffix(term, "$") {
			term = strings.TrimSuffix(term, "$")
			endWith = "$"
		}
		regex += regexp.QuoteMeta(term)
		regex += endWith
		results = DeleteNonMatchingItems(results, regex)
	}

	slices.SortFunc(results, func(a, b Key) int {
		a.Key = strings.TrimPrefix(a.Key, "services.")
		b.Key = strings.TrimPrefix(b.Key, "services.")
		a.Key = strings.TrimPrefix(a.Key, "programs.")
		b.Key = strings.TrimPrefix(b.Key, "programs.")
		if len(a.Key) < len(b.Key) {
			return -1
		} else if len(a.Key) > len(b.Key) {
			return 1
		}
		return 0
	})

	return
}

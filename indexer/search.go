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

func DeleteNonMatchingItems(keys Keys, pattern string, onName ...bool) Keys {
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Err: ", err)
		return keys
	}

	if len(onName) > 0 && onName[0] {
		var filtered Keys
		for _, k := range keys {
			if re.MatchString(k.Name) {
				filtered = append(filtered, k)
			}
		}
		return filtered
	} else {
		var filtered Keys
		for _, k := range keys {
			if re.MatchString(k.Key) {
				filtered = append(filtered, k)
			}
		}
		return filtered
	}
}

// TODO: Remove this function
func Contains(list []string, item string) bool {
	return slices.Contains(list, item)
}

func (keys Keys) Search(
	query string,
	exclude []string,
) (results Keys) {
	query = strings.TrimSpace(query)
	if query == "" {
		return keys
	}

	if strings.HasPrefix(query, "package ") {
		query = strings.TrimPrefix(query, "package ")
		exclude = []string{"nixos", "homemanager", "darwin"}
	} else if strings.HasPrefix(query, "option ") {
		query = strings.TrimPrefix(query, "option ")
		exclude = []string{"nixpkgs", "nur"}
	}

	results = slices.Clone(keys)

	var patterns []string
	if !Contains(exclude, "nixpkgs") {
		patterns = append(patterns, `^`+nixpkgs.Prefix+`.*`)
	}
	if !Contains(exclude, "nixos") {
		patterns = append(patterns, `^`+nixos.Prefix+`.*`)
	}
	if !Contains(exclude, "homemanager") {
		patterns = append(patterns, `^`+homemanager.Prefix+`.*`)
	}
	if !Contains(exclude, "darwin") {
		patterns = append(patterns, `^`+darwin.Prefix+`.*`)
	}
	if !Contains(exclude, "nur") {
		patterns = append(patterns, `^`+nur.Prefix+`.*`)
	}
	if len(patterns) > 0 {
		pattern := strings.Join(patterns, "|")
		results = DeleteNonMatchingItems(results, pattern, true)
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

package indexer

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

// PackageOrOption represents a package or option result.
type PackageOrOption struct {
	Type        string // "package" or "option"
	Source      string
	Key         string
	Description string

	Broken   bool
	Insecure bool
}

// packageRemoveNotMatching filters packages based on a regex pattern.
// If onlyOnKey is true, the regex is matched only against the key.
func packageRemoveNotMatching(i Packages, pattern string, onlyOnKey bool) Packages {
	res := Packages{}

	// Special search for maintainer
	if strings.HasPrefix(pattern, "?maintainer=") {
		pattern = strings.TrimPrefix(pattern, "?maintainer=")
		for key, pkg := range i {
			for _, maintainer := range pkg.Maintainers {
				if strings.EqualFold(maintainer.GitHub, pattern) {
					res[key] = pkg
					break
				}
			}
		}
		return res
	} else if strings.HasPrefix(pattern, "?broken") {
		for key, pkg := range i {
			if pkg.Broken {
				res[key] = pkg
			}
		}
		return res
	}

	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return i
	}
	if onlyOnKey {
		for key, pkg := range i {
			if re.MatchString(key) {
				res[key] = pkg
			}
		}
	} else {
		for key, pkg := range i {
			// Combining source, literal " package " and key for matching
			if re.MatchString(pkg.Source + " package " + key) {
				res[key] = pkg
			}
		}
	}
	return res
}

// optionRemoveNotMatching filters options based on a regex pattern.
// If onlyOnKey is true, the regex is matched only against the key.
func optionRemoveNotMatching(i Options, pattern string, onlyOnKey bool) Options {
	res := Options{}
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return i
	}
	if onlyOnKey {
		for key, opt := range i {
			if re.MatchString(key) {
				res[key] = opt
			}
		}
	} else {
		for key, opt := range i {
			// Combining source, literal " option " and key for matching
			if re.MatchString(opt.Source + " option " + key) {
				res[key] = opt
			}
		}
	}
	return res
}

// removeNotMaching filters the entire index for both packages and options.
func removeNotMaching(i Index, regex string, onlyOnKey bool) Index {
	i.Nixpkgs = packageRemoveNotMatching(i.Nixpkgs, regex, onlyOnKey)
	i.Nur = packageRemoveNotMatching(i.Nur, regex, onlyOnKey)
	i.Nixos = optionRemoveNotMatching(i.Nixos, regex, onlyOnKey)
	i.Homemanager = optionRemoveNotMatching(i.Homemanager, regex, onlyOnKey)
	i.Darwin = optionRemoveNotMatching(i.Darwin, regex, onlyOnKey)
	return i
}

// stripPrefix removes a possible prefix (e.g. "services." or "programs.") from the key.
func stripPrefix(key string) string {
	key = strings.TrimPrefix(key, "services.")
	key = strings.TrimPrefix(key, "programs.")
	key = strings.TrimSuffix(key, ".enable")
	key = strings.TrimSuffix(key, ".settings")
	return key
}

// Search performs a search on the index based on the provided query and exclude list.
// It returns a sorted slice of PackageOrOption results.
func (index Index) Search(query string) []PackageOrOption {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}

	results := Index{}

	fields := strings.Fields(query)

	if len(fields) == 0 {
		return []PackageOrOption{}
	}

	if fields[0] == "package" {
		results.Nixpkgs = index.Nixpkgs
		results.Nur = index.Nur
		results.Nixos = Options{}
		results.Homemanager = Options{}
		results.Darwin = Options{}
	} else if fields[0] == "option" {
		results.Nixpkgs = Packages{}
		results.Nur = Packages{}
		results.Nixos = index.Nixos
		results.Homemanager = index.Homemanager
		results.Darwin = index.Darwin
	} else {
		results.Nixpkgs = index.Nixpkgs
		results.Nur = index.Nur
		results.Nixos = index.Nixos
		results.Homemanager = index.Homemanager
		results.Darwin = index.Darwin
	}

	// Exclude sections if specified.
	if slices.Contains(fields, "!nixos") {
		results.Nixos = Options{}
		fields = slices.Delete(
			fields,
			slices.Index(fields, "!nixos"),
			slices.Index(fields, "!nixos")+1,
		)
	}
	if slices.Contains(fields, "!nixpkgs") {
		results.Nixpkgs = Packages{}
		fields = slices.Delete(
			fields,
			slices.Index(fields, "!nixpkgs"),
			slices.Index(fields, "!nixpkgs")+1,
		)
	}
	if slices.Contains(fields, "!nur") {
		results.Nur = Packages{}
		fields = slices.Delete(fields, slices.Index(fields, "!nur"), slices.Index(fields, "!nur")+1)
	}
	if slices.Contains(fields, "!home-manager") {
		results.Homemanager = Options{}
		fields = slices.Delete(
			fields,
			slices.Index(fields, "!home-manager"),
			slices.Index(fields, "!home-manager")+1,
		)
	}
	if slices.Contains(fields, "!darwin") {
		results.Darwin = Options{}
		fields = slices.Delete(
			fields,
			slices.Index(fields, "!darwin"),
			slices.Index(fields, "!darwin")+1,
		)
	}

	// Process each search term (using Fields to handle spaces efficiently).
	for _, term := range fields {
		regex := ""
		onlyOnKey := false
		endWith := ""
		if strings.HasPrefix(term, "^") {
			term = strings.TrimPrefix(term, "^")
			regex += "^"
			onlyOnKey = true
		}
		if strings.HasSuffix(term, "$") {
			term = strings.TrimSuffix(term, "$")
			endWith = "$"
			onlyOnKey = true
		}
		regex += regexp.QuoteMeta(term) + endWith
		results = removeNotMaching(results, regex, onlyOnKey)
	}

	// Combine packages and options into a single slice.
	var items []PackageOrOption
	for key, opt := range results.Nixos {
		items = append(
			items,
			PackageOrOption{
				Type:        "option",
				Source:      opt.Source,
				Key:         key,
				Description: opt.Description,
			},
		)
	}
	for key, opt := range results.Homemanager {
		items = append(
			items,
			PackageOrOption{
				Type:        "option",
				Source:      opt.Source,
				Key:         key,
				Description: opt.Description,
			},
		)
	}
	for key, opt := range results.Darwin {
		items = append(
			items,
			PackageOrOption{
				Type:        "option",
				Source:      opt.Source,
				Key:         key,
				Description: opt.Description,
			},
		)
	}
	for key, pkg := range results.Nixpkgs {
		items = append(
			items,
			PackageOrOption{
				Type:        "package",
				Source:      pkg.Source,
				Key:         key,
				Description: pkg.Description,
				Broken:      pkg.Broken,
				Insecure:    pkg.Insecure,
			},
		)
	}
	for key, pkg := range results.Nur {
		items = append(
			items,
			PackageOrOption{
				Type:        "package",
				Source:      pkg.Source,
				Key:         key,
				Description: pkg.Description,
				Broken:      pkg.Broken,
				Insecure:    pkg.Insecure,
			},
		)
	}

	// Naive sorting by the length of the key after stripping the prefix.
	sort.Slice(items, func(i, j int) bool {
		return len(stripPrefix(items[i].Key)) < len(stripPrefix(items[j].Key))
	})

	return items
}

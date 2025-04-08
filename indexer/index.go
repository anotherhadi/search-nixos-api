package indexer

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/anotherhadi/search-nixos-api/indexer/darwin"
	"github.com/anotherhadi/search-nixos-api/indexer/homemanager"
	"github.com/anotherhadi/search-nixos-api/indexer/nixos"
	"github.com/anotherhadi/search-nixos-api/indexer/nixpkgs"
	"github.com/anotherhadi/search-nixos-api/indexer/nur"
)

func simplifyPlatform(pkgs Package) Package {
	pkgs.PlatformsSimplify = []string{}
	platforms := []string{"darwin", "linux", "windows", "freebsd", "cygwin"}
	for _, p := range pkgs.Platforms {
		for _, platform := range platforms {
			if strings.Contains(p, platform) {
				if slices.Contains(pkgs.PlatformsSimplify, platform) {
					continue
				}
				pkgs.PlatformsSimplify = append(pkgs.PlatformsSimplify, platform)
				break
			}
		}
	}
	return pkgs
}

func dlNixos() Options {
	file := "nixos.json"
	log.Println("Downloading " + file + "...")
	jsonObject := map[string]nixos.Package{}
	err := downloadAndReadRelease(file, &jsonObject)
	if err != nil {
		log.Println(err)
	}
	log.Println("Downloaded " + file + " successfully")
	log.Println("Parsing " + file + "...")
	options := Options{}
	for k, v := range jsonObject {
		opt := Option{
			Source:       "nixpkgs",
			Type:         v.Type,
			Description:  v.Description,
			Declarations: v.Declarations,
			Default:      v.Default.Text,
			Example:      v.Example.Text,
		}
		options[k] = opt
	}
	log.Println("Parsed " + file + " successfully")
	return options
}

func dlHomemanager() Options {
	file := "home-manager.json"
	log.Println("Downloading " + file + "...")
	jsonObject := map[string]homemanager.Package{}
	err := downloadAndReadRelease(file, &jsonObject)
	if err != nil {
		log.Println(err)
	}
	log.Println("Downloaded " + file + " successfully")
	log.Println("Parsing " + file + "...")
	options := Options{}
	for k, v := range jsonObject {
		opt := Option{
			Source:       "home-manager",
			Type:         v.Type,
			Description:  v.Description,
			Declarations: []string{},
			Default:      v.Default.Text,
			Example:      v.Example.Text,
		}
		for _, d := range v.Declarations {
			opt.Declarations = append(opt.Declarations, d.URL)
		}
		options[k] = opt
	}
	log.Println("Parsed " + file + " successfully")
	return options
}

func dlDarwin() Options {
	file := "darwin.json"
	log.Println("Downloading " + file + "...")
	jsonObject := darwin.Darwin{}
	err := downloadAndReadRelease(file, &jsonObject)
	if err != nil {
		log.Println(err)
	}
	log.Println("Downloaded " + file + " successfully")
	log.Println("Parsing " + file + "...")
	options := Options{}
	for k, v := range jsonObject.Packages {
		opt := Option{
			Source:       "darwin",
			Type:         v.Type,
			Description:  v.Description,
			Declarations: v.DeclaredBy,
			Default:      v.Default,
			Example:      v.Example,
		}
		options[k] = opt
	}
	log.Println("Parsed " + file + " successfully")
	return options
}

func dlNixpkgs() Packages {
	file := "nixpkgs.json"
	log.Println("Downloading " + file + "...")
	jsonObject := nixpkgs.Nixpkgs{}
	err := downloadAndReadRelease(file, &jsonObject)
	if err != nil {
		log.Println(err)
	}
	log.Println("Downloaded " + file + " successfully")
	log.Println("Parsing " + file + "...")
	packages := Packages{}
	for k, v := range jsonObject.Packages {
		pkg := Package{
			Source:          "nixpkgs",
			Name:            v.Meta.Name,
			Version:         v.Version,
			Description:     v.Meta.Description,
			LongDescription: v.Meta.LongDescription,
			MainProgram:     v.Meta.MainProgram,
			Licenses:        []License{},
			Maintainers:     []Maintainer{},
			Broken:          v.Meta.Broken,
			Unfree:          v.Meta.Unfree,
			Position:        v.Meta.Position,
			PositionUrl:     v.Meta.Position,
		}

		pkg.PositionUrl = "https://github.com/NixOS/nixpkgs/blob/nixos-unstable/" + strings.Replace(
			v.Meta.Position,
			":",
			"#L",
			1,
		)
		if v.Meta.KnownVulnerabilities != nil {
			pkg.KnownVulnerabilities = v.Meta.KnownVulnerabilities
			pkg.Vulnerable = true
		} else {
			pkg.KnownVulnerabilities = []string{}
		}
		if v.Meta.Homepages != nil {
			pkg.Homepages = v.Meta.Homepages
		} else {
			pkg.Homepages = []string{}
		}
		if v.Meta.Platforms != nil {
			pkg.Platforms = v.Meta.Platforms
		} else {
			pkg.Platforms = []string{}
		}
		for _, l := range v.Meta.Licenses {
			pkg.Licenses = append(pkg.Licenses, License{
				FullName: l.FullName,
				Free:     l.Free,
				SpdxID:   l.SpdxID,
			})
		}
		for _, m := range v.Meta.Maintainers {
			pkg.Maintainers = append(pkg.Maintainers, Maintainer{
				Name:     m.Name,
				Email:    m.Email,
				GitHub:   m.GitHub,
				GithubId: m.GithubId,
			})
		}
		packages[k] = simplifyPlatform(pkg)
	}
	log.Println("Parsed " + file + " successfully")
	return packages
}

func dlNur() Packages {
	file := "nur.json"
	log.Println("Downloading " + file + "...")
	jsonObject := nur.Nur{}
	err := downloadAndReadRelease(file, &jsonObject)
	if err != nil {
		log.Println(err)
	}
	log.Println("Downloaded " + file + " successfully")
	log.Println("Parsing " + file + "...")
	packages := Packages{}
	for k, v := range jsonObject.Packages {
		pkg := Package{
			Source:          "nur",
			Name:            v.Meta.Name,
			Version:         v.Version,
			Description:     v.Meta.Description,
			LongDescription: v.Meta.LongDescription,
			MainProgram:     v.Meta.MainProgram,
			Licenses:        []License{},
			Maintainers:     []Maintainer{},
			Broken:          v.Meta.Broken,
			Unfree:          v.Meta.Unfree,
			Position:        v.Meta.Position,
			PositionUrl:     v.Meta.Position,
		}
		if v.Meta.KnownVulnerabilities != nil {
			pkg.KnownVulnerabilities = v.Meta.KnownVulnerabilities
			pkg.Vulnerable = true
		} else {
			pkg.KnownVulnerabilities = []string{}
		}
		if v.Meta.Homepages != nil {
			pkg.Homepages = v.Meta.Homepages
		} else {
			pkg.Homepages = []string{}
		}
		if v.Meta.Platforms != nil {
			pkg.Platforms = v.Meta.Platforms
		} else {
			pkg.Platforms = []string{}
		}
		for _, l := range v.Meta.Licenses {
			pkg.Licenses = append(pkg.Licenses, License{
				FullName: l.FullName,
				Free:     l.Free,
				SpdxID:   l.SpdxID,
			})
		}
		for _, m := range v.Meta.Maintainers {
			pkg.Maintainers = append(pkg.Maintainers, Maintainer{
				Name:     m.Name,
				Email:    m.Email,
				GitHub:   m.GitHub,
				GithubId: m.GithubId,
			})
		}
		packages[k] = simplifyPlatform(pkg)
	}
	log.Println("Parsed " + file + " successfully")
	return packages
}

func DownloadReleases(path string) {
	log.Println("Downloading releases...")
	index := Index{}

	index.Darwin = dlDarwin()
	index.Nixpkgs = dlNixpkgs()
	index.Nur = dlNur()
	index.Nixos = dlNixos()
	index.Homemanager = dlHomemanager()

	log.Println("Downloading version")
	rc, err := downloadRelease("version")
	if err != nil {
		log.Println(err)
		return
	}
	defer rc.Close()
	content, err := io.ReadAll(rc)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Writing index.json...")
	index.Info = map[string]string{
		"version":            string(content),
		"last-updated":       time.Now().Format(time.RFC3339),
		"nixos-length":       strconv.Itoa(len(index.Nixos)),
		"nixpkgs-length":     strconv.Itoa(len(index.Nixpkgs)),
		"nur-length":         strconv.Itoa(len(index.Nur)),
		"darwin-length":      strconv.Itoa(len(index.Darwin)),
		"homemanager-length": strconv.Itoa(len(index.Homemanager)),
	}

	indexFile, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer indexFile.Close()
	content, err = json.MarshalIndent(index, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = indexFile.Write(content)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Index downloaded and saved to", path)
	log.Println("Info:", index.Info)
	log.Println("Index file size:", len(content), "bytes")
}

func GetIndex(path string) (index Index) {
	if !DoesFileExist(path) {
		DownloadReleases(path)
	}

	log.Println("Opening index.json...")
	indexFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer indexFile.Close()
	content, err := io.ReadAll(indexFile)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Parsing index.json...")
	err = json.Unmarshal(content, &index)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Index opened successfully")
	return index
}

package indexer

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/anotherhadi/search-nixos-api/indexer/darwin"
	"github.com/anotherhadi/search-nixos-api/indexer/homemanager"
	"github.com/anotherhadi/search-nixos-api/indexer/nixos"
	"github.com/anotherhadi/search-nixos-api/indexer/nixpkgs"
	"github.com/anotherhadi/search-nixos-api/indexer/nur"
)

type Key struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Keys []Key

func (k Keys) Len() int {
	return len(k)
}

func (k Keys) String(i int) string {
	return k[i].Name
}

type Index struct {
	Info map[string]string `json:"info"`
	Keys Keys              `json:"keys"`

	Nixos       map[string]nixos.Package       `json:"nixos"`
	Nixpkgs     map[string]nixpkgs.Package     `json:"nixpkgs"`
	Homemanager map[string]homemanager.Package `json:"homemanager"`
	Darwin      map[string]darwin.Package      `json:"darwin"`
	Nur         map[string]nur.Package         `json:"nur"`
}

func DownloadReleases(path string) {
	log.Println("Downloading releases...")
	index := Index{}

	log.Println("Downloading nixos.json...")
	err := downloadAndReadRelease("nixos.json", &index.Nixos)
	if err != nil {
		log.Println(err)
	}

	log.Println("Downloading nixpkgs.json...")
	nixpkgsJson := nixpkgs.Nixpkgs{}
	err = downloadAndReadRelease("nixpkgs.json", &nixpkgsJson)
	if err != nil {
		log.Println(err)
	}
	index.Nixpkgs = nixpkgsJson.Packages

	log.Println("Downloading homemanager.json...")
	err = downloadAndReadRelease("home-manager.json", &index.Homemanager)
	if err != nil {
		log.Println(err)
	}

	log.Println("Downloading darwin.json...")
	darwinJson := darwin.Darwin{}
	err = downloadAndReadRelease("darwin.json", &darwinJson)
	if err != nil {
		log.Println(err)
	}
	index.Darwin = darwinJson.Packages

	log.Println("Downloading nur.json...")
	nurJson := nur.Nur{}
	err = downloadAndReadRelease("nur.json", &nurJson)
	if err != nil {
		log.Println(err)
	}
	index.Nur = nurJson.Packages

	log.Println("Downloading info")
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
		"version":        string(content),
		"last-updated":   time.Now().Format(time.RFC3339),
		"nixos-length":   strconv.Itoa(len(index.Nixos)),
		"nixpkgs-length": strconv.Itoa(len(index.Nixpkgs)),
		"hm-length":      strconv.Itoa(len(index.Homemanager)),
		"darwin-length":  strconv.Itoa(len(index.Darwin)),
		"nur-length":     strconv.Itoa(len(index.Nur)),
	}

	index.Keys = []Key{}
	for key := range index.Nixos {
		index.Keys = append(index.Keys, Key{
			Key:         key,
			Name:        nixos.Prefix + key,
			Description: index.Nixos[key].Description,
		})
	}

	for key := range index.Nixpkgs {
		index.Keys = append(index.Keys, Key{
			Key:         key,
			Name:        nixpkgs.Prefix + key,
			Description: index.Nixpkgs[key].Meta.Description,
		})
	}

	for key := range index.Homemanager {
		index.Keys = append(index.Keys, Key{
			Key:         key,
			Name:        homemanager.Prefix + key,
			Description: index.Homemanager[key].Description,
		})
	}

	for key := range index.Darwin {
		index.Keys = append(index.Keys, Key{
			Key:         key,
			Name:        darwin.Prefix + key,
			Description: index.Darwin[key].Description,
		})
	}

	for key := range index.Nur {
		index.Keys = append(index.Keys, Key{
			Key:         key,
			Name:        nur.Prefix + key,
			Description: index.Nur[key].Meta.Description,
		})
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
	log.Println("Version:", index.Info["version"])
	log.Println("Index file size:", len(content), "bytes")
	log.Println("Nixpkgs version:", index.Info["version"])
	log.Println("Nixpkgs length:", len(nixpkgsJson.Packages))
	log.Println("Nixos length:", index.Info["nixos-length"])
	log.Println("Homemanager length:", index.Info["hm-length"])
	log.Println("Darwin length:", index.Info["darwin-length"])
	log.Println("Nur length:", index.Info["nur-length"])
}

func GetIndex(path string) (index Index) {
	if !DoesFileExist(path) {
		DownloadReleases(path)
	}

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
	err = json.Unmarshal(content, &index)
	if err != nil {
		log.Println(err)
		return
	}

	return index
}

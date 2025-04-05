package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DoesFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getFileFromUrl(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get file from url: %s", url)
	}
	return resp.Body, nil
}

const url = "https://github.com/anotherhadi/nix-json/releases/latest/download/"

func downloadRelease(filename string) (io.ReadCloser, error) {
	return getFileFromUrl(url + filename)
}

func downloadAndReadRelease(filename string, unmarshalTo any) error {
	log.Println("Downloading", filename, "...")
	jsonfile, err := downloadRelease(filename)
	if err != nil {
		return err
	}
	defer jsonfile.Close()

	log.Println("Reading", filename, "...")
	content, err := io.ReadAll(jsonfile)
	if err != nil {
		return err
	}
	log.Println("Read", filename, "successfully :", len(content), "bytes")
	log.Println("Unmarshalling", filename, "...")

	err = json.Unmarshal(content, &unmarshalTo)
	if err != nil {
		return err
	}
	log.Println("Unmarshalled", filename, "successfully")

	return nil
}

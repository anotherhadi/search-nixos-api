package main

import (
	"os"
	"strings"
	"time"

	"github.com/anotherhadi/search-nixos-api/indexer"
	"github.com/anotherhadi/search-nixos-api/indexer/homemanager"
	"github.com/anotherhadi/search-nixos-api/indexer/nixos"
	"github.com/anotherhadi/search-nixos-api/indexer/nixpkgs"
	"github.com/gin-gonic/gin"
)

func main() {
	// Running in production mode by default
	production := os.Getenv("PRODUCTION")
	if production == "" {
		production = "true"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	indexPath := os.Getenv("INDEX_PATH")
	if indexPath == "" {
		if production == "true" {
			indexPath = "/var/lib/search-nixos-api/index.json"
		} else {
			indexPath = "./index.json"
		}
	}

	interval := os.Getenv("INDEX_INTERVAL")
	if interval == "" {
		interval = "12h"
	}
	intervalTime, err := time.ParseDuration(interval)
	if err != nil {
		panic(err)
	}

	index := indexer.GetIndex(indexPath)

	// Update the index every {interval} hours
	go func() {
		for {
			time.Sleep(intervalTime)
			indexer.DownloadReleases(indexPath)
			index = indexer.GetIndex(indexPath)
		}
	}()

	if strings.ToLower(production) == "true" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the Search NixOS API"})
	})

	r.GET("/index.json", func(c *gin.Context) {
		c.JSON(200, index)
	})

	r.GET("/stats", func(c *gin.Context) {
		c.JSON(200, index.Info)
	})

	// Search endpoint
	// curl -X GET "http://localhost:8080/search?q=kitty&exclude=nixos,homemanager"
	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		excludeStr := c.Query("exclude")
		exclude := []string{}
		if excludeStr != "" {
			exclude = strings.Split(excludeStr, ",")
		}
		if query == "" {
			c.JSON(400, gin.H{"error": "Query parameter 'q' is required"})
		} else {
			results := index.Keys.Search(query, exclude)
			c.JSON(200, results)
		}
	})

	r.GET(nixpkgs.Prefix+":q", func(c *gin.Context) {
		query := c.Param("q")
		if result, found := index.Nixpkgs[query]; found {
			c.JSON(200, result)
		} else {
			c.JSON(404, gin.H{"error": "Not found"})
		}
	})

	r.GET(nixos.Prefix+":q", func(c *gin.Context) {
		query := c.Param("q")
		if result, found := index.Nixos[query]; found {
			c.JSON(200, result)
		} else {
			c.JSON(404, gin.H{"error": "Not found"})
		}
	})

	r.GET(homemanager.Prefix+":q", func(c *gin.Context) {
		query := c.Param("q")
		if result, found := index.Homemanager[query]; found {
			c.JSON(200, result)
		} else {
			c.JSON(404, gin.H{"error": "Not found"})
		}
	})

	err = r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

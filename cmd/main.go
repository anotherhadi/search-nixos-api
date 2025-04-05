package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anotherhadi/search-nixos-api/indexer"
	"github.com/anotherhadi/search-nixos-api/indexer/darwin"
	"github.com/anotherhadi/search-nixos-api/indexer/homemanager"
	"github.com/anotherhadi/search-nixos-api/indexer/nixos"
	"github.com/anotherhadi/search-nixos-api/indexer/nixpkgs"
	"github.com/anotherhadi/search-nixos-api/indexer/nur"
	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
	}))

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
		page := c.Query("page")
		if page == "" {
			page = "1"
		}
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid page number"})
			return
		} else if pageInt < 1 {
			c.JSON(400, gin.H{"error": "Page number must be greater than 0"})
			return
		}
		perPage := c.Query("per_page")
		if perPage == "" {
			perPage = "20"
		}
		perPageInt, err := strconv.Atoi(perPage)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid per_page number"})
			return
		} else if perPageInt < 1 {
			c.JSON(400, gin.H{"error": "per_page number must be greater than 0"})
			return
		}

		excludeStr := c.Query("exclude")
		exclude := []string{}
		if excludeStr != "" {
			exclude = strings.Split(excludeStr, ",")
		}
		if query == "" {
			c.JSON(400, gin.H{"error": "Query parameter 'q' is required"})
		} else {
			results := index.Keys.Search(query, exclude)
			if len(results) == 0 {
				c.JSON(200, gin.H{"results": indexer.Keys{}, "total": 0, "page": pageInt, "per_page": perPageInt, "totalPages": 1})
				return
			}
			total := len(results)
			totalPages := (total + perPageInt - 1) / perPageInt
			if total > perPageInt {
				start := (pageInt - 1) * perPageInt
				end := start + perPageInt
				end = min(end, total)

				results = results[start:end]
			}
			c.JSON(200, gin.H{"results": results, "total": total, "totalPages": totalPages, "page": pageInt, "per_page": perPageInt})
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

	r.GET(darwin.Prefix+":q", func(c *gin.Context) {
		query := c.Param("q")
		if result, found := index.Darwin[query]; found {
			c.JSON(200, result)
		} else {
			c.JSON(404, gin.H{"error": "Not found"})
		}
	})

	r.GET(nur.Prefix+":q", func(c *gin.Context) {
		query := c.Param("q")
		if result, found := index.Nur[query]; found {
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

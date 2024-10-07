package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type URL struct {
	Route string `json:"route"`
	URL   string `json:"url"`
}

var urls = []URL{
	{Route: "/foo", URL: "https://github.com/donschenck"},
	{Route: "/bar", URL: "https://developers.redhat.com"},
}

func main() {
	router := gin.Default()

	// Define a GET endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"urlshortener-redirector": "v1.0.0",
		})
	})

	// Define a GET endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": "OK",
		})
	})

	// Get a URL for a route
	router.GET("/:route", func(c *gin.Context) {
		route := "/" + c.Param("route")
		url := getURL(route)
		if url != "" {
			c.JSON(http.StatusOK, url)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "url not found"})
	})

	// Start the server
	router.Run(":8080")
}

func getURL(routeId string) string {
	connStr := "user=shorties password=shorties dbname=urls host=shorties sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Query for a single row
	var url string

	err = db.QueryRow("SELECT url FROM routes WHERE route = '$1'", 1).Scan(&routeId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows were returned!")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("ROUTE: %s, URL: %s\n", routeId, url)
	}
	return url
}

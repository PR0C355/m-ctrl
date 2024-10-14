package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Mission struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text"`
}

var missions = []Mission{}

func main() {
	// Load missions from JSON into application memory
	loadMissionsIntoMemory("/mnt/media")

	// Create & Configure gin router
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	router.GET("/missions/get/all", getMissions)
	router.GET("/missions/get/random", getRandomMission)
	router.GET("/missions/get/unique", getUniqueMission)
	router.GET("/missions/get/:id", getMissionByID)

	// Start the server
	router.Run("0.0.0.0:8080")
}

func loadMissionsIntoMemory(path string) {

	// Open the directory
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the files
	for _, e := range dirs {

		targetPath := filepath.Join(path, e.Name(), "data.json")

		// Check if the file exists
		if _, err := os.Stat(targetPath); err == nil {

			// Open the file
			file, _ := os.Open(targetPath)
			defer file.Close()

			// Read the file
			byteValue, err := os.ReadFile(targetPath)
			if err != nil {
				fmt.Println("Error reading JSON file:", err)
				return
			}

			// Create a variable of the struct type
			var mission Mission

			// Unmarshal the JSON data into the struct
			err = json.Unmarshal(byteValue, &mission)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return
			}

			// Use the struct
			fmt.Printf("ID: %s, Title: %s, URL: %s, Text: %s\n\n", mission.ID, mission.Title, mission.URL, mission.Text)
			missions = append(missions, mission)

		} else {
			fmt.Printf("File %s does not exist\n", targetPath)
		}

	}

}

func generateDailyHash(sessionToken string) string {
	// Generate a hash based on the session token and the current date
	currentDate := time.Now().Format("2006-01-02")

	// Combine the session token and the current date
	combinedString := sessionToken + currentDate

	// Generate a SHA256 hash of the combined string
	hash := sha256.Sum256([]byte(combinedString))

	// Return the hash as a hex string
	return hex.EncodeToString(hash[:])
}

func getMissions(c *gin.Context) {
	// Return the missions as JSON
	c.IndentedJSON(http.StatusOK, missions)
}

func getMissionByID(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Loop through the missions to find the one with the matching ID
	for _, a := range missions {

		// If the ID matches, return the mission
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	// If no mission with the ID is found, return an error
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
}

func getRandomMission(c *gin.Context) {
	// Use the current time as a seed for the random number generator
	c.IndentedJSON(http.StatusOK, missions[rand.Intn(len(missions))])
}

func getUniqueMission(c *gin.Context) {

	// Get the session token from the query string
	sessionToken := c.Query("token")

	// If the session token is empty, return an error
	if sessionToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session Token Required."})
		return
	}

	// Generate a daily hash based on the session token
	dailyHash := generateDailyHash(sessionToken)

	// Use the first 8 characters of the hash as a seed for the random number generator
	seed, _ := strconv.ParseInt(dailyHash[:8], 16, 64)

	// Create a new random number generator with the seed
	r := rand.New(rand.NewSource(seed))

	// Generate a random index based on the length of the missions
	randomIndex := r.Intn(len(missions))

	// Return the random mission
	c.IndentedJSON(http.StatusOK, missions[randomIndex])
}

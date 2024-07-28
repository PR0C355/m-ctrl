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
	loadMissionsIntoMemory("/mnt/media")
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8000", "http://mission.tumi.dev", "https://mission.tumi.dev"}
	router.Use(cors.New(config))

	router.GET("/missions/get/all", getMissions)
	router.GET("/missions/get/random", getRandomMission)
	router.GET("/missions/get/unique", getUniqueMission)
	router.GET("/missions/get/:id", getMissionByID)
	//router.POST("/missions/add", addMission)
	router.Run("0.0.0.0:8080")
}

func loadMissionsIntoMemory(path string) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range dirs {
		targetPath := filepath.Join(path, e.Name(), "data.json")

		if _, err := os.Stat(targetPath); err == nil {
			file, _ := os.Open(targetPath)
			defer file.Close()

			// byteValue, err := io.ReadAll(file)
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

			// Now you can use the struct
			fmt.Printf("ID: %s, Title: %s, URL: %s, Text: %s\n\n", mission.ID, mission.Title, mission.URL, mission.Text)
			missions = append(missions, mission)
		} else {
			fmt.Printf("File %s does not exist\n", targetPath)
		}

	}

}

func generateDailyHash(sessionToken string) string {
	currentDate := time.Now().Format("2006-01-02")
	combinedString := sessionToken + currentDate
	hash := sha256.Sum256([]byte(combinedString))
	return hex.EncodeToString(hash[:])
}

func getMissions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, missions)
}

func getMissionByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range missions {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
}

func getRandomMission(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, missions[rand.Intn(len(missions))])
}

func getUniqueMission(c *gin.Context) {
	sessionToken := c.Query("token")
	if sessionToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session Token Required."})
		return
	}

	dailyHash := generateDailyHash(sessionToken)

	// Use the first 8 characters of the hash as a seed for the random number generator
	seed, _ := strconv.ParseInt(dailyHash[:8], 16, 64)
	r := rand.New(rand.NewSource(seed))

	randomIndex := r.Intn(len(missions))
	c.IndentedJSON(http.StatusOK, missions[randomIndex])
}

// func addMission(c *gin.Context) {
// 	var newMission Mission

// 	if err := c.BindJSON(&newMission); err != nil {
// 		return
// 	}

// 	missions = append(missions, newMission)
// 	c.IndentedJSON(http.StatusCreated, newMission)
// }

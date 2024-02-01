package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Wikipedia URLs for the topic of interest
	urls := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
	}

	// Create a new collector
	c := colly.NewCollector()

	// Create a slice that stores content from Wikipedia URLs
	var sliceTexts []string

	// On every html element that has an arritbute call callback (see lines 17-18 on colly basic.go)
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// Extract text and append to the slice
		sliceTexts = append(sliceTexts, e.Text)
	})

	// Set up error handling (see lines 20-22 on colly error_handing.go)
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Use Colly's 'Visit' method to fetch and process each URL
	for _, url := range urls {
		err := c.Visit(url)
		if err != nil {
			// Log error with details about URL if error occurs
			log.Println("Error visiting URL:", url, "\nError:", err)
		}
	}

	// Create a JSON lines file and write the extracted text
	file, err := os.Create("extractedoutput.jl")
	if err != nil {
		// Log fatal error and exit program if error occurs
		log.Fatal("Error creating output file:", err)
	}
	defer file.Close()

	// Create JSON encoder and write to specific file
	enc := json.NewEncoder(file)
	// Encode text as JSON object with single key 'text' and its corresponding value
	for _, text := range sliceTexts {
		err := enc.Encode(map[string]string{"text": text})
		// Log error if one occurs
		if err != nil {
			log.Println("Error encoding JSON:", err)
		}
	}

	log.Println("Extraction complete. Check 'extractedoutput.jl' for the results.")
}

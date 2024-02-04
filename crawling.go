package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jaytaylor/html2text"
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
	var stringsSlice []string

	// Initialize info card as a map
	infoCard := make(map[string]string)

	// Define a regular expression to clean up values
	re := regexp.MustCompile(`\[\d+\]`)

	// On every HTML element that has an attribute, call the callback
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// Extract and clean title
		title := strings.TrimSpace(e.ChildText("#firstHeading"))
		organizationName := strings.TrimSpace(e.ChildText("#mw-content-text > div > table.infobox.vcard > caption"))

		// Extract and clean strings
		e.ForEach("#mw-content-text > div > p", func(_ int, node *colly.HTMLElement) {
			text, err := html2text.FromString(node.Text)
			if err == nil {
				stringsSlice = append(stringsSlice, strings.TrimSpace(text))
			}
		})

		// Extract info card values
		e.ForEach("#mw-content-text > div > table.infobox.vcard tr", func(_ int, row *colly.HTMLElement) {
			item := strings.TrimSpace(row.ChildText("th"))
			value := strings.TrimSpace(row.ChildText("td"))

			// Clean up values
			value = re.ReplaceAllString(value, "")
			value = strings.TrimSpace(value)

			// Handle special cases, e.g., Website
			if item == "Website" {
				infoCard[item] = value
			} else {
				infoCard[item] = strings.Join(strings.Fields(value), " ")
			}
		})

		// Write to JSON file
		file, err := os.Create("extractedoutput.jl")
		if err != nil {
			log.Fatal("Error creating output file:", err)
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		err = enc.Encode(map[string]interface{}{
			"Title":             title,
			"Organization_name": organizationName,
			"Info_card":         infoCard,
			"Strings":           stringsSlice,
		})
		if err != nil {
			log.Println("Error encoding JSON:", err)
		}
	})

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Use Colly's 'Visit' method to fetch and process each URL
	for _, url := range urls {
		err := c.Visit(url)
		if err != nil {
			log.Println("Error visiting URL:", url, "\nError:", err)
		}
	}

	log.Println("Extraction complete. Check 'extractedoutput.jl' for the results.")
}

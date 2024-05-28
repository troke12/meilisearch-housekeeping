package main

import (
	"fmt"
	"log"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

func main() {
	// Replace these with your MeiliSearch instance URL and index UID
	meilisearchURL := "http://localhost:7700"
	indexUID := "Alerts"
	timestampField := "@timestamp"

	// Calculate the timestamp for March 31, 2023
	startDate := time.Date(2023, time.March, 31, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, time.October, 31, 23, 59, 59, 999999999, time.UTC)

	// Create a MeiliSearch client
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: meilisearchURL,
	})

	// Initialize batch variables
	batchSize := int64(1000)
	offset := int64(0)

	// Loop through documents and delete in batches
	for {
		// Fetch documents in batches
		searchResult, err := client.Index(indexUID).Search("", &meilisearch.SearchRequest{
			Limit:  batchSize,
			Offset: offset,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Break the loop if no more documents
		if len(searchResult.Hits) == 0 {
			break
		}

		// Loop through documents in the batch
		for _, hit := range searchResult.Hits {
			doc, ok := hit.(map[string]interface{})
			if !ok {
				log.Fatal("Error converting hit to map[string]interface{}")
			}

			// Check if the timestamp field exists
			timestamp, ok := doc[timestampField]
			if !ok {
				fmt.Println("Timestamp field not found in the document")
				continue
			}

			// Attempt to parse the timestamp
			createdTime, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", timestamp))
			if err != nil {
				fmt.Printf("Error parsing timestamp: %v\n", err)
				continue
			}

			// Check if the document is in the date range
			if createdTime.After(startDate) && createdTime.Before(endDate) {
				// Delete the document
				_, err := client.Index(indexUID).DeleteDocument(doc["id"].(string))
				if err != nil {
					fmt.Printf("Error deleting document with ID %s: %v\n", doc["id"], err)
				} else {
					fmt.Printf("Deleted document with ID: %s\n", doc["id"])
				}
			}
		}

		// Update offset for the next batch
		offset += batchSize
	}
}

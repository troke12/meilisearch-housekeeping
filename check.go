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

	// Create a MeiliSearch client
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: meilisearchURL,
	})

	// Fetch all documents in the "Alerts" index
	searchResult, err := client.Index(indexUID).Search("", &meilisearch.SearchRequest{
		Limit: 100000,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Convert hits to []map[string]interface{}
	var allDocuments []map[string]interface{}
	for _, hit := range searchResult.Hits {
		doc, ok := hit.(map[string]interface{})
		if !ok {
			log.Fatal("Error converting hit to map[string]interface{}")
		}
		allDocuments = append(allDocuments, doc)
	}

	// Print information about each document
	for _, doc := range allDocuments {
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

		// Print information about the document
		fmt.Printf("Document ID: %s\n", doc["id"])
		fmt.Printf("Timestamp: %v\n", timestamp)
		fmt.Printf("Parsed Timestamp: %v\n", createdTime)
		fmt.Println("----------")
	}

	// Print total number of documents
	fmt.Printf("Total Documents: %d\n", len(allDocuments))
}

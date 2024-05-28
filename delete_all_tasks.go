package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// Replace this with your MeiliSearch instance URL
	meilisearchURL := "http://localhost:7700"

	// Create an HTTP client
	client := &http.Client{}

	// Fetch all tasks
	tasks, err := getAllTasks(meilisearchURL, client)
	if err != nil {
		log.Fatal(err)
	}

	// Delete each task
	for _, task := range tasks {
		err := deleteTask(meilisearchURL, task.UID, client)
		if err != nil {
			fmt.Printf("Error deleting task with UID %s: %v\n", task.UID, err)
		} else {
			fmt.Printf("Deleted task with UID: %s\n", task.UID)
		}
	}
}

// Task represents a MeiliSearch task
type Task struct {
	UID string `json:"uid"`
}

// getAllTasks fetches all tasks from MeiliSearch
func getAllTasks(meilisearchURL string, client *http.Client) ([]Task, error) {
	url := fmt.Sprintf("%s/api/indexes/_all/tasks", meilisearchURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tasks []Task
	if err := decodeJSON(resp.Body, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// deleteTask deletes a task from MeiliSearch
func deleteTask(meilisearchURL, taskUID string, client *http.Client) error {
	url := fmt.Sprintf("%s/api/indexes/_all/tasks/%s", meilisearchURL, taskUID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// decodeJSON decodes JSON response into a target interface
func decodeJSON(reader io.Reader, target interface{}) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(target)
}

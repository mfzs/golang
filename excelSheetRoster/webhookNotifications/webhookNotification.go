package webhookNotification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"net/http"
)

// SendWebhookNotification sends the emails.txt content as a message to a Microsoft Teams channel webhook URL
func SendWebhookNotification(webhookURL string, filePath string) error {
	// Open the emails.txt file
	emailFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer emailFile.Close()

	// Read the content of the file using io.ReadAll
	emailContent, err := io.ReadAll(emailFile)
	if err != nil {
		return fmt.Errorf("failed to read file content: %v", err)
	}

	// Prepare the JSON payload for Teams webhook with markdown formatting
	payload := map[string]interface{}{
		"type":    "MessageCard",
		"context": "http://schema.org/extensions",
		"summary": "Roster Notification", // This can be a short summary
		"text":    fmt.Sprintf("```\n%s\n```", string(emailContent)),  // Wrap the content in markdown code block
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	// Send the POST request with the JSON data
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read and log the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	log.Printf("Webhook Response: %s", string(respBody))
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"io"

	"github.com/joho/godotenv"
	"github.com/go-resty/resty/v2"
)

func loadEnv() {
	if err := godotenv.Load("env_file"); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func getAccessToken() (string, error) {
    clientID := os.Getenv("CLIENT_ID")
    clientSecret := os.Getenv("CLIENT_SECRET")
    tenantID := os.Getenv("TENANT_ID")

    // Use only .default scope for client credentials flow
    tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

    formData := url.Values{
        "grant_type":    {"client_credentials"},
        "client_id":     {clientID},
        "client_secret": {clientSecret},
        "scope":         {"https://graph.microsoft.com/.default"},
    }

    resp, err := http.PostForm(tokenURL, formData)
    if err != nil {
        return "", fmt.Errorf("failed to get access token: %v", err)
    }
    defer resp.Body.Close()

    // Log the status and body for debugging
    fmt.Println("Response Status:", resp.Status)
    body, _ := io.ReadAll(resp.Body)
    // fmt.Println("Response Body:", string(body))

    var tokenResponse struct {
        AccessToken string `json:"access_token"`
    }
    if err := json.Unmarshal(body, &tokenResponse); err != nil {
        return "", fmt.Errorf("failed to decode token response: %v", err)
    }

    if tokenResponse.AccessToken == "" {
        return "", fmt.Errorf("access token is empty")
    }

    // fmt.Println("Access Token:", tokenResponse.AccessToken)

    return tokenResponse.AccessToken, nil
}




func downloadFile(accessToken, fileID, destFilePath string) error {
	client := resty.New()

	// Correct URL for downloading files from OneDrive or SharePoint
	userID := "" // Replace with actual user ID
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/drive/items/%s/content", userID, fileID)

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(url)

	if err != nil {
		return err
	}

	outFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(resp.Body())
	if err != nil {
		return err
	}

	fmt.Println("File downloaded successfully.")
	return nil
}

func main() {
	loadEnv()

	fileID := "" // Replace with actual file ID
	destFilePath := ""

	// fmt.Println("CLIENT_ID:", os.Getenv("CLIENT_ID"))
	// fmt.Println("CLIENT_SECRET:", os.Getenv("CLIENT_SECRET"))
	// fmt.Println("TENANT_ID:", os.Getenv("TENANT_ID"))
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	err = downloadFile(accessToken, fileID, destFilePath)
	if err != nil {
		log.Fatalf("Failed to download the file: %v", err)
	}
}

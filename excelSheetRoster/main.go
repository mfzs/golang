package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"excelSheetRoster/webhookNotifications"
	"excelSheetRoster/emailNotifications" 
	"github.com/xuri/excelize/v2"
)

func main() {

	err := loadEnv("env_file")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Get values from the environment
	filePath := os.Getenv("FILE_PATH")
	webhookURL := os.Getenv("WEBHOOK_URL")
	
	// List of sheet names to process
	sheetNames := []string{"Infra oncall", "Diagnostics oncall", "Search oncall"}

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open the Excel file: %v", err)
	}

	// Get the testing week's starting date (Monday) + 1 week
	currentTime := time.Now()
	currentWeekStart := currentTime.AddDate(0, 0, -int(currentTime.Weekday())+1)	
	currentWeekFormatted := currentWeekStart.Format("1/2/06") 

	currentWeekEnd := currentWeekStart.AddDate(0, 0, 6)
	currentWeekFormattedEnd := currentWeekEnd.Format("1/2/06")

	// Create or open the output file
	outputFile, err := os.Create("emails.txt")
	if err != nil {
		log.Fatalf("Failed to create the output file: %v", err)
	}
	defer outputFile.Close()

	// Write the testing week's date range to the file
	outputFile.WriteString(fmt.Sprintf("Roster week starting from: %s to: %s", currentWeekFormatted,currentWeekFormattedEnd))

	// Iterate over each sheet
	for _, sheetName := range sheetNames {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("Failed to read rows from sheet %s: %v\n", sheetName, err)
			continue
		}

		// Find column indexes dynamically
		headerRow := rows[0]
		oncallIndex := -1
		rosterIndex := -1
		phoneIndex := -1

		for i, col := range headerRow {
			colName := strings.ToLower(strings.TrimSpace(col))
			if colName == "oncall" {
				oncallIndex = i
			} else if colName == "roster" {
				rosterIndex = i
			} else if strings.Contains(colName, "phone") {
				phoneIndex = i
			}
		}

		// Validate column indexes
		if oncallIndex == -1 || rosterIndex == -1 || phoneIndex == -1 {
			outputFile.WriteString(fmt.Sprintf("%s: Missing required columns (Oncall, Roster, Phone numbers)\n", sheetName))
			continue
		}

		// Look for the testing week's on-call email
		var oncallEmail, phoneNumber string

		// Find the On-call email for the testing week
		for i := 1; i < len(rows); i++ {
			row := rows[i]

			// Skip rows with insufficient data
			if len(row) <= oncallIndex {
				continue
			}

			// Parse and normalize the WeekStarting column
			weekStarting := strings.ReplaceAll(row[0], "-", "/") // Replace "-" with "/"
			parsedWeekStarting, err := time.Parse("1/2/06", weekStarting)
			if err != nil {
				continue
			}

			// Check if the WeekStarting matches the testing week's start
			if parsedWeekStarting.Format("1/2/06") == currentWeekFormatted {
				oncallEmail = strings.TrimSpace(row[oncallIndex])
				break
			}
		}

		// Skip if no email is found for this week
		if oncallEmail == "" {
			outputFile.WriteString(fmt.Sprintf("\n\n%s: No on-call email found for this week\n", sheetName))
			continue
		}

		// Find the phone number for the on-call email in the "Roster" section
		for i := 1; i < len(rows); i++ {
			row := rows[i]

			// Skip rows with insufficient data
			if len(row) <= rosterIndex || len(row) <= phoneIndex {
				continue
			}

			// If a match is found, save the phone number to the output file
			if strings.Contains(strings.ToLower(row[rosterIndex]), strings.ToLower(oncallEmail)) {
				phoneNumber = strings.TrimSpace(row[phoneIndex])
				break
			}
		}

		// Write the results to the output file
		if phoneNumber != "" {
			outputFile.WriteString(fmt.Sprintf("\n\n%s: %s \nPhone: %s", sheetName, oncallEmail, phoneNumber))
		} else {
			outputFile.WriteString(fmt.Sprintf("\n\n%s %s (Phone: N/A)\n", sheetName, oncallEmail))
		}
	}

	// Send email notifications
	err = emailNotification.SendAlertToEmails("emails.txt")
	if err != nil {
		log.Fatalf("Failed to send alert: %v", err)
	}

	// Send notification in teams's channel

	err = webhookNotification.SendWebhookNotification(webhookURL, "emails.txt")
	if err != nil {
		log.Fatalf("Failed to send webhook notification: %v", err)
	}
}

func loadEnv(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse each line of the env file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return scanner.Err()
}
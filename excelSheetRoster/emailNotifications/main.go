package emailNotification

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

// Configuration variables
var (
	smtpHost       string
	smtpPort       string
	senderName     string
	senderEmail    string
	senderPassword string
	disclaimer     = "\n\nNOTE: If you are unavailable during this time, please contact your manager immediately.\n"
	ccEmails       = []string{
		"firojsiddique100@gmail.com",
	}
)

// Load environment variables from .env file
func loadEnvVariables() error {
	err := godotenv.Load("env_file")
	if err != nil {
		return fmt.Errorf("error loading env_file")
	}

	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	senderName = os.Getenv("SENDER_NAME")
	senderEmail = os.Getenv("SENDER_EMAIL")
	senderPassword = os.Getenv("SENDER_PASSWORD")

	if smtpHost == "" || smtpPort == "" || senderName == "" || senderEmail == "" || senderPassword == "" {
		return fmt.Errorf("missing required environment variables")
	}
	return nil
}

// Extract valid email addresses from a file
func extractEmailsFromFile(filename string) ([]string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	emails := []string{}
	scanner := bufio.NewScanner(file)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)

	var fileContent strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		fileContent.WriteString(line + "\n")

		matches := emailRegex.FindAllString(line, -1)
		emails = append(emails, matches...)
	}

	if err := scanner.Err(); err != nil {
		return nil, "", err
	}

	return emails, fileContent.String(), nil
}

// Send a single alert to all recipients
func sendAlertToAll(recipients []string, message string, ccEmails []string) error {
	// Connect to SMTP server
	serverAddress := smtpHost + ":" + smtpPort
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %v", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", senderName, senderPassword, smtpHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Set sender
	if err := client.Mail(senderEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// Add all recipients (To + Cc)
	allRecipients := append(recipients, ccEmails...)
	for _, recipient := range allRecipients {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add recipient %s: %v", recipient, err)
		}
	}

	// Start email data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data command: %v", err)
	}
	defer writer.Close()

	// Headers
	toHeader := strings.Join(recipients, ", ")
	ccHeader := strings.Join(ccEmails, ", ")
	headers := fmt.Sprintf("From: %s\nTo: %s\nCc: %s\nSubject: Roster Notification\n\n",
		senderEmail, toHeader, ccHeader)

	// Body
	body := fmt.Sprintf("Hello N7 Folks,\n\n%s\n%s\nBest regards,\nAlert System", message, disclaimer)

	// Write the email
	_, err = writer.Write([]byte(headers + body))
	if err != nil {
		return fmt.Errorf("failed to write email body: %v", err)
	}

	return nil
}

// Public function to send alerts from file
func SendAlertToEmails(filePath string) error {
	// Load env vars
	if err := loadEnvVariables(); err != nil {
		return err
	}

	// Extract email list and message
	emails, fileContent, err := extractEmailsFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to extract emails: %v", err)
	}

	if len(emails) == 0 {
		return fmt.Errorf("no email addresses found in file")
	}

	// Send one alert to all recipients
	err = sendAlertToAll(emails, fileContent, ccEmails)
	if err != nil {
		return fmt.Errorf("failed to send alert: %v", err)
	}

	log.Printf("Alert sent successfully to all recipients")
	return nil
}

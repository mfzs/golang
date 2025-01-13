package emailNotification

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"log"
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
		"pratiksha.somani@neuron7.ai", 
		"abinaya.govindan@neuron7.ai",
		"chetan.kalyan@neuron7.ai",
		"gyan.ranjan@neuron7.ai",
		"infra@neuron7.ai",
	}
)

// Load environment variables from .env file
func loadEnvVariables() error {
	err := godotenv.Load("env_file") // Replace with your .env file name if it's different
	if err != nil {
		return fmt.Errorf("error loading env_file")
	}

	// Read environment variables
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	senderName = os.Getenv("SENDER_NAME")
	senderEmail = os.Getenv("SENDER_EMAIL")
	senderPassword = os.Getenv("SENDER_PASSWORD")

	// Ensure the required environment variables are set
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
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	var fileContent strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		fileContent.WriteString(line + "\n") // Append the line to file content

		matches := emailRegex.FindAllString(line, -1) // Extract all email addresses
		emails = append(emails, matches...)
	}

	if err := scanner.Err(); err != nil {
		return nil, "", err
	}

	return emails, fileContent.String(), nil
}

// Send alert to a specific email
func sendAlert(email string, message string, ccEmails []string) error {
	// Connect to the SMTP server
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

	// Set sender and recipients
	if err := client.Mail(senderEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err := client.Rcpt(email); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Add CC recipients
	for _, cc := range ccEmails {
		if err := client.Rcpt(cc); err != nil {
			return fmt.Errorf("failed to set CC recipient: %v", err)
		}
	}

	// Send email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data command: %v", err)
	}
	defer writer.Close()

	// Construct email headers and body
	ccHeader := strings.Join(ccEmails, ", ")
	headers := fmt.Sprintf("From: %s\nTo: %s\nCc: %s\nSubject: Roster Notification\n\n",
		senderEmail, email, ccHeader)
	body := fmt.Sprintf("Hello N7 Folks,\n\n%s\n%s\nBest regards,\nAlert System", message, disclaimer)

	// Write headers and body to the email
	_, err = writer.Write([]byte(headers + body))
	if err != nil {
		return fmt.Errorf("failed to write email body: %v", err)
	}

	return nil
}

// Send alerts to all valid emails from the provided file
func SendAlertToEmails(filePath string) error {
	// Load environment variables
	if err := loadEnvVariables(); err != nil {
		return err
	}

	// Extract emails and file content from the provided file
	emails, fileContent, err := extractEmailsFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read emails from file: %v", err)
	}

	// Add file content to alert message
	alertMessageWithContent := fileContent

	// Send alerts to each valid email and include CC
	for _, email := range emails {
		if err := sendAlert(email, alertMessageWithContent, ccEmails); err != nil {
			log.Printf("Failed to send alert to %s: %v\n", email, err)
		} else {
			log.Printf("Alert sent successfully to %s\n", email)
		}
	}

	return nil
}

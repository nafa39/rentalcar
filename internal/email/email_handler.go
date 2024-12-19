package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailHandler struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderName   string
	SMTPPassword string
	SMTPUser     string
}

// // NewEmailHandler initializes a new instance of EmailHandler with environment variables
// func NewEmailHandler() *EmailHandler {
// 	// Log the environment variables to confirm their values
// 	log.Println("SMTP Host: ", os.Getenv("SMTP_HOST"))
// 	log.Println("SMTP Port: ", os.Getenv("SMTP_PORT"))
// 	log.Println("Sender Email: ", os.Getenv("SMTP_SENDER_EMAIL"))
// 	log.Println("Sender Name: ", os.Getenv("SMTP_SENDER_NAME"))
// 	log.Println("SMTP Password: ", os.Getenv("SMTP_PASSWORD"))
// 	log.Println("SMTP Username: ", os.Getenv("SMTP_USERNAME"))
// 	return &EmailHandler{}
// }

// SendEmail sends an email using SMTP
func SendEmail(toEmail, subject, body string) error {

	SMTPHost := os.Getenv("SMTP_HOST")
	SMTPPort := os.Getenv("SMTP_PORT")
	SenderEmail := os.Getenv("SMTP_SENDER_EMAIL")
	SenderName := os.Getenv("SMTP_SENDER_NAME")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")
	SMTPUser := os.Getenv("SMTP_USERNAME")

	auth := smtp.PlainAuth("", SMTPUser, SMTPPassword, SMTPHost)

	message := []byte(fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		SenderName, SenderEmail, toEmail, subject, body,
	))

	address := fmt.Sprintf("%s:%s", SMTPHost, SMTPPort)

	if err := smtp.SendMail(address, auth, SenderEmail, []string{toEmail}, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

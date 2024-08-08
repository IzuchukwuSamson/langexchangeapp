package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"regexp"
	"text/template"
	"time"
)

type smtpDetails struct {
	User, Password, Host, Port, From string
}

// Function to validate email format
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	// Check if the email matches the regular expression
	return re.MatchString(email)
}

// getSmtpDetails returns the smtp details from the app environment
func getSmtpDetails() smtpDetails {
	return smtpDetails{
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		From:     os.Getenv("MAIL_FROM"),
	}
}

type EmailType int

const (
	ResetPassword EmailType = iota
	VerifyEmail
	SendAdminLink
)

type EmailInfo struct {
	FirstName, Email string
}

type EmailData struct {
	Code, Name, Link, Email string
}

type AdminEmailData struct {
	Email, Link string
}

// SendTokenMail sends an email containing the random string to the specified receiver
func sendTokenMail(to, subject, templateFile string, data interface{}) <-chan error {
	det := getSmtpDetails()
	auth := smtp.PlainAuth("", det.User, det.Password, det.Host)

	// Parse the specified template file
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		errs := make(chan error, 1)
		errs <- err
		return errs
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\n" + "Content-Type: text/html; charset=\"UTF-8\";\n" +
		fmt.Sprintf("From: %s\n", det.From) + fmt.Sprintf("To: %s\n", to) + "\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	err = t.Execute(&body, data)
	if err != nil {
		errs := make(chan error, 1)
		errs <- err
		return errs
	}

	sendErr := smtp.SendMail(fmt.Sprintf("%s:%s", det.Host, det.Port), auth, det.User, []string{to}, body.Bytes())
	errs := make(chan error, 1)
	errs <- sendErr
	return errs
}

func generateExpiryDate(minutes int) string {
	expiryTime := time.Now().Add(time.Duration(minutes) * time.Minute)
	return expiryTime.Format(time.RFC3339) // or use another format if needed
}

func SendEmail(emailType EmailType, user EmailInfo, token string) <-chan error {
	errs := make(chan error, 1) // Create a channel to return errors

	go func() {
		defer close(errs) // Close the channel when done

		var subject, templateFile string
		var data EmailData

		switch emailType {
		case ResetPassword:
			subject = "Reset Your Password"
			templateFile = "templates/reset_password.html"
			data = EmailData{
				Code: token,
				Name: user.FirstName,
			}
		case VerifyEmail:
			subject = "Verify Your Account"
			templateFile = "templates/verify_email.html"
			data = EmailData{
				Code: token,
				Name: user.FirstName,
			}
		case SendAdminLink:
			subject = "Welcome to the Datenest Admin Team"
			templateFile = "templates/verify_admin_email.html"
			// Generate expiry date
			expiryDate := generateExpiryDate(10)
			data = EmailData{
				Link:  fmt.Sprintf("https://datenest.net/admin?token=%s&expires=%s", token, expiryDate),
				Email: user.Email,
			}
		default:
			errs <- fmt.Errorf("invalid email type")
			return
		}

		// Send the email and handle error channel
		errChan := sendTokenMail(user.Email, subject, templateFile, data)
		if err := <-errChan; err != nil {
			errs <- err // Send error to channel
		} else {
			errs <- nil // Send nil if no error
		}
	}()

	return errs
}

package alert

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Alert is an interface declaration for all types of alerts.
type Alert interface {
	Send() (string, []error)
}

// EmailAlert represents the data model for any email alert. With this model, we should be able to
//   to send an HTML structured message to an end user.
type EmailAlert struct {
	From      *mail.Email
	Message   string
	Subject   string
	Timestamp time.Time
	To        []*mail.Email
}

// NewEmailAlert is a factory function for producing email alert objects.
func NewEmailAlert(fromName, fromEmail, toName, toEmail, subject, html string) *EmailAlert {
	t := time.Now()
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	return &EmailAlert{
		From:      from,
		To:        []*mail.Email{to},
		Timestamp: t,
		Subject:   subject,
		Message:   html,
	}
}

// Send uses sendgrid to send out an email, and returns any output and email
func (e *EmailAlert) Send() (string, []error) {
	var errs []error
	if os.Getenv("SENDGRID_API_KEY") == "" {
		return "", append(errs, errors.New("SENDGRID_API_KEY not defined, can't send email alert. Please set the environment variable and try again"))
	}
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	var output string

	for _, to := range e.To {
		message := mail.NewSingleEmail(e.From, e.Subject, to, e.Message, e.Message)
		response, err := client.Send(message)
		if err != nil {
			errs = append(errs, err)
		} else {
			output += fmt.Sprintf("%s\tAlert sent: %s to %s\tStatus: %d\tBody: %s\n", e.Timestamp, e.Subject, to, response.StatusCode, response.Body)
		}
	}
	return output, errs
}

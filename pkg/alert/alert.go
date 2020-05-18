package alert

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Alert is an interface declaration for all types of alerts.
type Alert interface {
	Evaluate(attributeData string) bool // Evalute
	GetName() string
	GetURIs() []string
	Send() (string, []error)
}

// EmailAlert represents the data model for any email alert. With this model, we should be able to
//   to send an HTML structured message to an end user.
type EmailAlert struct {
	Attribute string // Determines what piece of scan data is used to evaluate against
	From      *mail.Email
	Message   string
	Name      string
	Regex     string // Used in Evaluate to determine if we should alert
	Subject   string
	Timestamp time.Time
	To        []*mail.Email
	URIs      []string
}

// NewEmailAlert is a factory function for producing email alert objects.
// TODO - refactor this mess
func NewEmailAlert(name, fromName, fromEmail, toName, toEmail, subject, html string) *EmailAlert {
	t := time.Now()
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	return &EmailAlert{
		From:      from,
		Message:   html,
		Name:      name,
		Subject:   subject,
		Timestamp: t,
		To:        []*mail.Email{to},
	}
}

// Evaluate takes the expected input, and returns true if our regex matches a pattern against it
// TODO - proper error handling
func (e *EmailAlert) Evaluate(input string) bool {
	match, _ := regexp.MatchString(e.Regex, input)
	return match
}

// GetName returns the alert name to identify the alert
func (e *EmailAlert) GetName() string {
	return e.Name
}

// GetURI returns the alert name to identify the alert
func (e *EmailAlert) GetURIs() []string {
	return e.URIs
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
			output += fmt.Sprintf("%s\nAlert sent: %s to %s - %s\tStatus: %d\n∂Body: %s\n", e.Timestamp, e.Subject, to.Name, to.Address, response.StatusCode, response.Body)
		}
	}
	return output, errs
}

package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailAlert(t *testing.T) {
	alert := NewEmailAlert("Ryan Hartje", "ryan@ryanhartje.com", "Ryan Hartje", "hartjepc@gmail.com", "Testing", "<H1>HOLY BANANAS</H1>")
	output, errs := alert.Send()
	for _, err := range errs {
		assert.NoErrorf(t, err, "Received an error when sending an email alert: %s", err.Error())
	}
	assert.NotEqual(t, output, "", "Expected output to be received when testing email sending")

	t.Logf("output: %s", output)
}

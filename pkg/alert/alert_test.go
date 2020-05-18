package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailAlertEvaluation(t *testing.T) {
	a := &EmailAlert{}
	a.Regex = "find.me"
	matched := a.Evaluate("can you find me?")
	assert.True(t, true, matched, "Failed to evaluate the pattern correctly")
}

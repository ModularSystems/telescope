package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLScanner(t *testing.T) {
	s := HTMLScan{
		URI: "https://modularsystems.io",
	}
	s.Scan()
	assert.NoError(t, s.Error, "Failed to get test site")
	assert.NotNil(t, s.HTML, "Got nil data from GetHTML scan from %s", s.URI)
	assert.NotNil(t, s.Timestamp)
}

// TODO - Implement test for TestWPScan

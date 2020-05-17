package scanner

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLScanner(t *testing.T) {
	s := HTMLScan{
		URI: "https://google.com",
	}
	s.Scan()
	assert.NoError(t, s.Error, "Failed to get google")
	assert.NotNil(t, s.HTML, "Got nil data from GetHTML scanner when scanning google")
	assert.NotNil(t, s.Timestamp)
}

// TestWPScanner checks for the existance of WPVULNDB_API_KEY, and if its there, clears it so that we don't waste requests
// as the free version only allows 50 requests/day
func TestWPScanner(t *testing.T) {
	keyCache := os.Getenv("WPVULNDB_API_KEY")
	if keyCache != "" {
		os.Setenv("WPVULNDB_API_KEY", "")
	}

	w := WPScan{
		URI: "https://mamalovesyacookies.com",
	}
	w.Scan()

	assert.Nil(t, w.Error)
	assert.Equal(t, w.Stderr, "")
	assert.NotEqual(t, w.Stdout, "")
	assert.NotNil(t, w.Timestamp)
	if keyCache != "" {
		os.Setenv("WPVULNDB_API_KEY", keyCache)
	}
}

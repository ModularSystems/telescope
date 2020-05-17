package scan

import "time"

// Scanner sets up a common interface for a variety of scans
type Scanner interface {
	GetTimestamp() time.Time
	GetName() string
	IsEligible(time.Time) bool // Used to determine if a scan is eligible to run
	Scan() bool                // Scan returns true when a scan is performed
}

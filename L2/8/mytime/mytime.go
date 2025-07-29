package mytime

import (
	"fmt"
	"time"

	"github.com/beevik/ntp"
)

// Now prints the current time from an NTP server.
func Now() (time.Time, error) {
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get time from ntp server: %w", err)
	}
	return ntpTime, nil
}

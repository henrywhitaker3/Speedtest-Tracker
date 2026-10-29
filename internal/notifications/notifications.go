// Package notifications
package notifications

import (
	"fmt"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
)

type Phase int

const (
	PreRun Phase = iota
	PostRun

	StartingTitle = "Starting Speedtest"
	SuccessTitle  = "Finished Speedtest"
	ErrorTitle    = "Failed Speedtest"
)

func (p Phase) String() string {
	switch p {
	case PreRun:
		return string(config.PreRun)
	case PostRun:
		return string(config.PostRun)
	default:
		return fmt.Sprintf("%d", p)
	}
}

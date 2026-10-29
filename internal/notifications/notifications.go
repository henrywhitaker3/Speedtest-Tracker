// Package notifications
package notifications

type Phase int

const (
	PreRun Phase = iota
	PostRun

	StartingTitle = "Starting Speedtest"
	SuccessTitle  = "Finished Speedtest"
	ErrorTitle    = "Failed Speedtest"
)

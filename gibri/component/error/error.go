package error

import gibrierror "github.com/amirphl/gitsi/gibri/error"

var (
	FailedToJoinCall = gibrierror.New(gibrierror.Session, "Failed to join the call")
	ChromeHung       = gibrierror.New(gibrierror.Session, "Chrome Hung")
	NoMediaReceived  = gibrierror.New(gibrierror.Session, "No media received")
	ICEFailed        = gibrierror.New(gibrierror.Session, "ICE failed")
)

package state

type State string

const (
	StartingUp       State = "starting up"
	Running          State = "running"
	Finished         State = "finished"
	FailedToJoinCall State = "failed to join call"
	ChromeHung       State = "chrome hung"
	NoMediaReceived  State = "no media received"
	ICEFailed        State = "ICE failed"
	Undefined        State = "undefined"
)

func (c State) String() string {
	return string(c)
}

func FromString(s string) State {
	switch s {
	case StartingUp.String():
		return StartingUp
	case Running.String():
		return Running
	case Finished.String():
		return Finished
	case FailedToJoinCall.String():
		return FailedToJoinCall
	case ChromeHung.String():
		return ChromeHung
	case NoMediaReceived.String():
		return NoMediaReceived
	case ICEFailed.String():
		return ICEFailed
	default:
		return Undefined
	}
}

// enum class ErrorScope

package state

type State int

const (
	StartingUp State = iota
	Running
	Finished
	FailedToJoinCall
	ChromeHung
	NoMediaReceived
	ICEFailed
)

func (c State) String() string {
	switch c {
	case StartingUp:
		return "starting up"
	case Running:
		return "running"
	case Finished:
		return "finished"
	case FailedToJoinCall:
		return "failed to join call"
	case ChromeHung:
		return "chrome hung"
	case NoMediaReceived:
		return "no media received"
	case ICEFailed:
		return "ICE failed"
	default:
		return "undefined"
	}
}

// enum class ErrorScope

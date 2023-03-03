package event

type Event int

const (
	CallJoined Event = iota
	FailedToJoinCall
	CallEmpty
	ChromeHung
	NoMediaReceived
	ICEFailed
	LocalParticipantKicked
	OK
)

func (s Event) String() string {
	switch s {
	case CallJoined:
		return "call joined"
	case FailedToJoinCall:
		return "failed to join call"
	case CallEmpty:
		return "call empty"
	case ChromeHung:
		return "chrome hung"
	case NoMediaReceived:
		return "no media received"
	case ICEFailed:
		return "ice failed"
	case LocalParticipantKicked:
		return "local participant kicked"
	default:
		return "ok"
	}
}

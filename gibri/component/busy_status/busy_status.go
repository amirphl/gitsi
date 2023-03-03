package busystatus

type BusyStatus int

const (
	Busy BusyStatus = iota
	Idle
	Expired
)

func (b BusyStatus) String() string {
	switch b {
	case Busy:
		return "busy"
	case Idle:
		return "idle"
	case Expired:
		return "expired"
	default:
		return "Undefined"
	}
}

package failurereason

type FailureReason int

const (
	Busy FailureReason = iota
	Error
	Undefined
)

func (f FailureReason) String() string {
	switch f {
	case Busy:
		return "busy"
	case Error:
		return "error"
	default:
		return "undefined"
	}
}

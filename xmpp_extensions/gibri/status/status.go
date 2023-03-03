package status

type Status int

const (
	On Status = iota
	Off
	Pending
	Undefined
)

func (s Status) String() string {
	switch s {
	case On:
		return "on"
	case Off:
		return "off"
	case Pending:
		return "pending"
	default:
		return "undefined"
	}
}

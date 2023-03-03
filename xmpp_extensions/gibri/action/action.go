package action

type Action int

const (
	Start Action = iota
	Stop
	Undefined
)

func (a Action) String() string {
	switch a {
	case Start:
		return "start"
	case Stop:
		return "stop"
	default:
		return "undefined"
	}
}

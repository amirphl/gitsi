package recordingmode

type RecordingMode int

const (
	File RecordingMode = iota
	Stream
	Undefined
)

func (r RecordingMode) String() string {
	switch r {
	case File:
		return "file"
	case Stream:
		return "stream"
	default:
		return "undefined"
	}
}

package statuscheck

import (
	"log"
	"time"

	"github.com/amirphl/gitsi/gibri/component/event"
	"github.com/amirphl/gitsi/gibri/selenium/page"
	"github.com/amirphl/gitsi/gibri/selenium/state_machine"
)

const (
	DefaultEmptyCallTimeout = 30 * time.Second // TODO "jibri.call-status-checks.default-call-empty-timeout".from(Config.configSource)
)

type CallSc interface {
	run(page.CallPage) event.Event
}

type EmptyCallSc interface {
	CallSc
}

type ICEConnectionSc interface {
	CallSc
}

type LocalParticipantKickedSc interface {
	CallSc
}

type MediaReceivedSc interface {
	CallSc
}

type emptyCallSc struct {
	tracker statemachine.TransitionTimeTracker
	timeout time.Duration
}

type iceConnectionSc struct {
	lastSuccessTime time.Time
	timeout         time.Duration
}

type localParticipantKickedSc struct {
}

type mediaReceivedSc struct {
	// The timestamp at which we last saw that all clients transitioned to muted
	tracker         statemachine.TransitionTimeTracker
	lastMediaTime   time.Time
	noMediaTimeout  time.Duration
	allMutedTimeout time.Duration
}

func (e *emptyCallSc) run(callPage page.CallPage) event.Event {
	e.tracker.MaybeUpdate(callPage.IsEmpty())

	if e.tracker.ExceededTimeout(e.timeout) {
		t := e.tracker.Timestamp()

		log.Printf(
			"Info: Call has been empty since %v (%v ago). Returning CallEmpty event\n",
			t,
			time.Since(t),
		)

		return event.CallEmpty
	}

	return event.OK
}

func (i *iceConnectionSc) run(callPage page.CallPage) event.Event {
	now := time.Now()

	if callPage.IsEmpty() || callPage.IsICEConnected() {
		i.lastSuccessTime = now
		return event.OK
	}

	if now.Sub(i.lastSuccessTime) > i.timeout {
		log.Printf("Warn: ICE has failed and not recovered in %v\n", i.timeout)
		return event.ICEFailed
	}

	return event.OK
}

func (l *localParticipantKickedSc) run(callPage page.CallPage) event.Event {
	if callPage.IsLocalParticipantKicked() {
		log.Println("Info: Local participant was kicked, returning LocalParticipantKicked event")
		return event.LocalParticipantKicked
	}

	return event.OK
}

func (m *mediaReceivedSc) run(callPage page.CallPage) event.Event {
	now := time.Now()
	bitrates := callPage.Bitrates()
	numParticipants := callPage.NumParticipants() - 1 // mines Jibri
	numMutedParticipants := callPage.NumRemoteParticipantsMuted()
	numJigasiParticipants := callPage.NumRemoteParticipantsJigasi()

	// We don't get any mute state for Jigasi participants, so to prevent timing out when only Jigasi participants
	// may be speaking, always count them as "muted"
	allClientsMuted := (numMutedParticipants + numJigasiParticipants) == numParticipants

	log.Printf(
		"receive bitrates: %v, num participants: %d, numMutedParticipants: %d, numJigasis: %d, all clients muted? %t\n",
		bitrates,
		numParticipants,
		numMutedParticipants,
		numJigasiParticipants,
		allClientsMuted,
	)

	m.tracker.MaybeUpdate(allClientsMuted)
	downloadBitrate := int64(0)

	if val, ok := bitrates["download"]; ok {
		downloadBitrate = val.(int64)
	}

	// If all clients are muted, register it as 'receiving media': that way when clients unmute
	// we'll get the full noMediaTimeout duration before timing out due to lack of media.
	if downloadBitrate != 0 || allClientsMuted {
		m.lastMediaTime = now
	}

	timeSinceLastMedia := now.Sub(m.lastMediaTime)

	// There are a couple possible outcomes here:
	// 1) All clients are muted, but have been muted for longer than allMutedTimeout so
	//     we'll exit the call gracefully (CallEmpty)
	// 2) No media has flowed for longer than noMediaTimeout and all clients are not
	//     muted so we'll exit with an error (NoMediaReceived)
	// 3) If neither of the above are true, we're fine and no event has occurred
	if m.tracker.ExceededTimeout(m.allMutedTimeout) {
		return event.CallEmpty
	}
	if timeSinceLastMedia > m.noMediaTimeout && !allClientsMuted {
		return event.NoMediaReceived
	}

	return event.OK
}

func NewEmptyCallSc(timeout time.Duration) EmptyCallSc {
	tracker := statemachine.NewTransitionTimeTracker()
	log.Printf("Info: Starting empty call check with a timeout of %v\n", timeout)

	return &emptyCallSc{
		tracker: tracker,
		timeout: timeout,
	}
}

func NewICEConnectionSc() ICEConnectionSc {
	return &iceConnectionSc{
		lastSuccessTime: time.Now(),
		timeout:         30 * time.Second, // TODO "jibri.call-status-checks.ice-connection-timeout".from(Config.configSource)
	}
}

func NewLocalParticipantKickedSc() LocalParticipantKickedSc {
	log.Println("Info: Starting local participant kicked out call check")

	return &localParticipantKickedSc{}
}

func NewMediaReceivedSc() MediaReceivedSc {
	return &mediaReceivedSc{
		tracker:         statemachine.NewTransitionTimeTracker(),
		lastMediaTime:   time.Now(),
		noMediaTimeout:  3 * time.Minute,  // "jibri.call-status-checks.no-media-timeout".from(Config.configSource)
		allMutedTimeout: 10 * time.Minute, // "jibri.call-status-checks.all-muted-timeout".from(Config.configSource)
	}
}

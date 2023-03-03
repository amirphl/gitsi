package gibri

import (
	"github.com/amirphl/gitsi/xmpp_extensions/gibri/action"
	"github.com/amirphl/gitsi/xmpp_extensions/gibri/failure_reason"
	"github.com/amirphl/gitsi/xmpp_extensions/gibri/recording_mode"
	"github.com/amirphl/gitsi/xmpp_extensions/gibri/status"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
)

// TODO whole this file

const (
	ActionAttrName             = "action"
	DisplayNameAttrName        = "displayname"
	Element                    = "jibri"
	Namespace                  = "http://jitsi.org/protocol/jibri"
	SipAddressAttrName         = "sipaddress"
	StatusAttrName             = "status"
	FailureReasonAttrName      = "failure_reason"
	ShouldRetryAttrName        = "should_retry"
	StreamIdAttrName           = "streamid"
	YoutubeBroadcastIdAttrName = "you_tube_broadcast_id"
	SessionIdAttrName          = "session_id"
	AppDataAttrName            = "app_data"
	RecordingModeAttrName      = "recording_mode"
	RoomAfterName              = "room"
)

type IQ interface {
	ChildElementBuilder()
	StanzaId() string
	From() jid.JID
}

type iq struct {
	stanzaIQ      stanza.IQ
	action        action.Action
	recordingMode recordingmode.RecordingMode
	failureReason failurereason.FailureReason
	status        status.Status
	// displayName        string
	// sipAddress         string
	// streamId           string
	// youTubeBroadCastId string
	sessionId string
	// appData            string
	// shouldRetry        bool
	// room               jid.JID
}

// TODO
func (g *iq) ChildElementBuilder() {

}

func (g *iq) StanzaId() string {
	return g.stanzaIQ.ID
}

func (g *iq) From() jid.JID {
	return g.stanzaIQ.From
}

func NewIQ() IQ {
	// TODO default val for other fields
	return &iq{
		action:        action.Undefined,
		recordingMode: recordingmode.Undefined,
		failureReason: failurereason.Undefined,
		status:        status.Undefined,
	}
}

func CreateResult(request IQ, sessionId string) IQ {
	// TODO default val for other fields
	return &iq{
		stanzaIQ: stanza.IQ{
			Type: stanza.ResultIQ,
			ID:   request.StanzaId(),
			To:   request.From(),
		},
		sessionId: sessionId,
	}
}

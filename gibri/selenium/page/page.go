package page

import (
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

// TODO All functions require heavy testing.
// TODO Do scripts work without trimMargin() called?

type AbstractPage interface {
	Visit(url string) bool
}

type CallPage interface {
	AbstractPage
	NumParticipants() int32
	IsEmpty() bool
	Stats() map[string]interface{}
	Bitrates() map[string]interface{}
	InjectParticipantTrackerScript() bool
	InjectLocalParticipantTrackerScript() bool
	Participants() []map[string]interface{}
	NumRemoteParticipantsJigasi() int32
	IsICEConnected() bool
	IsLocalParticipantKicked() bool
	NumRemoteParticipantsMuted() int32
	AddToPresence(key, val string) bool
	SendPresence() bool
	Leave() bool
}

type HomePage interface {
	AbstractPage
}

type callPage struct {
	driver selenium.WebDriver
}

type homePage struct {
	driver selenium.WebDriver
}

func visit(driver selenium.WebDriver, url string) bool {
	log.Printf("Info: Visiting url %s\n", url)

	start := time.Now()
	err := driver.Get(url)

	if err != nil {
		log.Printf("Error: Failed to visited url %s\n", url)
		return false
	}

	elapsed := time.Since(start)

	log.Printf("Info: Waited %v for driver to load page\n", elapsed)

	return true
}

func checkRoomIsJoined(wd selenium.WebDriver) (bool, error) {
	script := `
		try {
			return APP.conference._room.isJoined();
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := wd.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Debug: Not joined yet: %v\n", res)
		return false, err // TODO Return error?
	}

	val, ok := res.(bool)

	if !ok {
		log.Printf("Debug: Not joined yet: %v\n", val)
		return false, err // TODO Return error?
	}

	return val, nil
}

func (c *callPage) Visit(url string) bool {
	start := time.Now()

	if !visit(c.driver, url) {
		return false
	}

	err := c.driver.WaitWithTimeoutAndInterval(
		checkRoomIsJoined,
		30*time.Second,
		5*time.Second, // TODO Adjust the interval. What is default interval in java version?
	)

	if err != nil {
		log.Println("Error: Timed out waiting for call page to load")
		return false
	}

	elapsed := time.Since(start)

	log.Printf("Info: Waited %v to join the conference\n", elapsed)

	return true
}

func (c *callPage) NumParticipants() int32 {
	script :=
		`
		try {
			return APP.conference.membersCount;
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running GetNumParticipants script: %v\n", err)
		return 1
	}

	val, ok := res.(int32)

	if !ok {
		log.Printf("Warn: parsing GetNumParticipants script result: %v\n", res)
		return 1
	}

	return val
}

func (c *callPage) IsEmpty() bool {
	return c.NumParticipants() == 1
}

func (c *callPage) Stats() map[string]interface{} {
	script :=
		`
		try {
			return APP.conference.getStats();
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg
	stats := make(map[string]interface{})

	if err != nil {
		log.Printf("Error: running GetStats script: %v\n", err)
		return stats
	}

	// TODO performance: cost of copy
	s := reflect.ValueOf(res)

	if s.Kind() != reflect.Map {
		log.Printf("Warn: parsing GetStats script result: %v\n", res)
		return stats
	}

	iter := s.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		stats[k.String()] = v.Interface()
	}

	return stats
}

func (c *callPage) Bitrates() map[string]interface{} {
	stats := c.Stats()
	res, ok := stats["bitrate"]
	bitrates := make(map[string]interface{})

	if !ok {
		return bitrates
	}

	// TODO performance: cost of copy
	s := reflect.ValueOf(res)

	if s.Kind() != reflect.Map {
		return bitrates
	}

	iter := s.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		bitrates[k.String()] = v.Interface()
	}

	return bitrates
}

func (c *callPage) InjectParticipantTrackerScript() bool {
	script :=
		`
		try {
			window._jibriParticipants = [];
			const existingMembers = APP.conference._room.room.members || {};
			const existingMemberJids = Object.keys(existingMembers);
			console.log("There were " + existingMemberJids.length + " existing members");
			existingMemberJids.forEach(jid => {
				const existingMember = existingMembers[jid];
				if (existingMember.identity) {
					console.log("Member ", existingMember, " has identity, adding");
					window._jibriParticipants.push(existingMember.identity);
				} else {
					console.log("Member ", existingMember.jid, " has no identity, skipping");
				}
			});
			APP.conference._room.room.addListener(
				"xmpp.muc_member_joined",
				(from, nick, role, hidden, statsid, status, identity) => {
					console.log("Jibri got MUC_MEMBER_JOINED: ", from, identity);
					if (!hidden && identity) {
						window._jibriParticipants.push(identity);
					}
				}
			);
			return true;
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running InjectParticipantTrackerScript script: %v\n", err)
		return false
	}

	val, ok := res.(bool)

	if !ok {
		log.Printf("Warn: parsing InjectParticipantTrackerScript script result: %v\n", res)
		return false
	}

	return val
}

func (c *callPage) InjectLocalParticipantTrackerScript() bool {
	script :=
		`
		try {
			window._isLocalParticipantKicked=false
			
			APP.conference._room.room.addListener(
				"xmpp.kicked",
				(isSelfPresence, actorId, kickedParticipantId, reason) => {
					console.log("Jibri got a KICKED event: ", isSelfPresence, actorId, kickedParticipantId, reason);
					if (isSelfPresence) {
						window._isLocalParticipantKicked=true
					}
				}
			);
			
			return true;
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running InjectLocalParticipantTrackerScript script: %v\n", err)
		return false
	}

	val, ok := res.(bool)

	if !ok {
		log.Printf("Warn: parsing InjectLocalParticipantTrackerScript script result: %v\n", res)
		return false
	}

	return val
}

func (c *callPage) Participants() []map[string]interface{} {
	script :=
		`
		try {
			return window._jibriParticipants;
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg
	participants := []map[string]interface{}{}

	if err != nil {
		log.Printf("Error: running GetParticipants script: %v\n", err)
		return participants
	}

	// TODO performance: cost of copy
	s := reflect.ValueOf(res)

	if s.Kind() != reflect.Slice {
		log.Printf("Warn: parsing GetParticipants script result: %v\n", res)
		return participants
	}

	for i := 0; i < s.Len(); i++ {
		m := s.Index(i)

		if m.Kind() == reflect.Map {
			p := make(map[string]interface{})
			iter := m.MapRange()

			for iter.Next() {
				k := iter.Key()
				v := iter.Value()
				p[k.String()] = v.Interface()
			}

			participants = append(participants, p)
		}
	}

	return participants
}

func (c *callPage) NumRemoteParticipantsJigasi() int32 {
	script :=
		`
		try {
			return APP.conference._room.getParticipants()
				.filter(participant => participant.getProperty("features_jigasi") == true)
				.length;
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running numRemoteParticipantsJigasi script: %v\n", res)
		return 0
	}

	val, ok := res.(int32)

	if !ok {
		log.Printf("Warn: running numRemoteParticipantsJigasi script result: %v\n", res)
		return 0
	}

	return val
}

func (c *callPage) IsICEConnected() bool {
	script :=
		`
		try {
			return APP.conference.getConnectionState();
		} catch(e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running IsICEConnected script: %v\n", err)
		return false
	}

	val, ok := res.(string)

	if !ok {
		log.Printf("Warn: parsing IsICEConnected script result: %v\n", res)
		return false
	}

	isConnected := strings.ToLower(val) == "connected"

	if !isConnected {
		log.Printf("ICE not connected: %s\n", val)
	}

	return isConnected
}

func (c *callPage) IsLocalParticipantKicked() bool {
	script :=
		`
		try {
			return window._isLocalParticipantKicked;
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running IsLocalParticipantKicked script: %v\n", err)
		return false
	}

	val, ok := res.(bool)

	if !ok {
		log.Printf("Warn: parsing IsLocalParticipantKicked script result: %v\n", res)
		return false
	}

	return val
}

func (c *callPage) NumRemoteParticipantsMuted() int32 {
	script :=
		`
		try {
			return APP.conference._room.getParticipants()
				.filter(participant => participant.isAudioMuted() && participant.isVideoMuted())
				.length;
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running NumRemoteParticipantsMuted script: %v\n", err)
		return 0
	}

	val, ok := res.(int32)

	if !ok {
		log.Printf("Warn: parsing NumRemoteParticipantsMuted script result: %v\n", res)
		return 0
	}

	return val
}

func (c *callPage) AddToPresence(key, val string) bool {
	script :=
		`
		try {
			APP.conference._room.room.addToPresence(
				'$key',
				{
					value: '$value'
				}
			);
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running AddToPresence script: %v\n", err)
		return false
	}

	_, ok := res.(string)

	return !ok
}

func (c *callPage) SendPresence() bool {
	script :=
		`
		try {
			APP.conference._room.room.sendPresence();
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running SendPresence script: %v\n", err)
		return false
	}

	_, ok := res.(string)

	return !ok
}

func (c *callPage) Leave() bool {
	script :=
		`
		try {
			return APP.conference._room.leave();
		} catch (e) {
			return e.message;
		}
		`
	script = strings.Trim(script, " ")
	res, err := c.driver.ExecuteScript(script, nil) // TODO nil arg

	if err != nil {
		log.Printf("Error: running Leave script: %v\n", err)
		return false
	}

	err = c.driver.WaitWithTimeoutAndInterval(
		func(_ selenium.WebDriver) (bool, error) {
			return c.NumParticipants() == 1, nil
		},
		2*time.Second,
		2*time.Second, // TODO adjust value
	)

	if err != nil {
		log.Printf("Error: checking num participants == 1: %v\n", err)
		return false
	}

	_, ok := res.(string)

	return !ok
}

func (h *homePage) Visit(url string) bool {
	return visit(h.driver, url)
}

func NewCallPage() CallPage {
	// TODO PageFactory.initElements(driver, this)
	return &callPage{}
}

func NewHomePage() HomePage {
	// TODO PageFactory.initElements(driver, this)
	return &homePage{}
}

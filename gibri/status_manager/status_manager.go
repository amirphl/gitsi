package statusmanager

import (
	"log"
	"sync"

	"github.com/amirphl/gitsi/gibri/collections/list"
	"github.com/amirphl/gitsi/gibri/component/busy_status"
	"github.com/amirphl/gitsi/gibri/component/health_status"
	"github.com/amirphl/gitsi/gibri/component/status"
	"github.com/amirphl/gitsi/xmpp_extensions/gibri/failure_reason"
	gibristatus "github.com/amirphl/gitsi/xmpp_extensions/gibri/status"
)

type OverallStatus struct {
	busystatus.BusyStatus
	healthstatus.OverallHealth
}

type Failure struct {
	reason failurereason.FailureReason // TODO fix // TODO May be null
	err    error                       // TODO May be null
}

type SessionStatus struct {
	status      gibristatus.Status // TODO
	failure     Failure            // TODO
	sessionId   string
	sipAddress  string
	shouldRetry bool
}

type StatusManager interface {
	status.Publisher
	BusyStatus() busystatus.BusyStatus
	SetBusyStatus(busystatus.BusyStatus)
	HealthStatus() healthstatus.HealthStatus
	OverallHealth() healthstatus.OverallHealth
	OverallStatus() OverallStatus
}

type statusManager struct {
	statusHandlers     list.List // TODO copyonwritearraylist
	subComponentHealth sync.Map  // TODO concurrency, ref of val
	busyStatus         busystatus.BusyStatus
	lock               sync.Mutex
}

func (g *statusManager) AddHandler(h status.Handler) {
	g.statusHandlers.Add(h)
}

func (g *statusManager) Publish(s interface{}) {
	f := func(item interface{}) bool {
		return item.(status.Handler).Handle(s)
	}

	g.statusHandlers.RetainAll(f)
}

func (g *statusManager) HealthStatus() healthstatus.HealthStatus {
	hs := healthstatus.Healthy

	// TODO concurrency
	// TODO duplicate code
	g.subComponentHealth.Range(func(_, v interface{}) bool {
		c := v.(healthstatus.HealthDetail)
		hs = hs.And(c.HealthStatus)

		return true
	})

	return hs
}

func (g *statusManager) OverallHealth() healthstatus.OverallHealth {
	hd := make(healthstatus.HealthDetails)
	hs := healthstatus.Healthy

	// TODO concurrency
	g.subComponentHealth.Range(func(k, v interface{}) bool {
		c := v.(healthstatus.HealthDetail)
		hs = hs.And(c.HealthStatus)
		hd[k.(string)] = c

		return true
	})

	return healthstatus.OverallHealth{
		HealthStatus:  hs,
		HealthDetails: hd,
	}
}

func (g *statusManager) BusyStatus() busystatus.BusyStatus {
	return g.busyStatus
}

func (g *statusManager) OverallStatus() OverallStatus {
	return OverallStatus{
		g.BusyStatus(),
		g.OverallHealth(),
	}
}

func (g *statusManager) SetBusyStatus(newStatus busystatus.BusyStatus) {
	oldStatus := g.busyStatus
	g.busyStatus = newStatus

	if oldStatus != newStatus {
		log.Printf("Info: Busy status has changed: %s -> %s\n", oldStatus.String(), newStatus.String())
		g.Publish(g.OverallStatus())
	}
}

func (g *statusManager) UpdateHealth(
	componentName string,
	healthStatus healthstatus.HealthStatus,
	detail string,
) {
	g.lock.Lock()

	log.Printf("Info: Received component health update: %s has status %s (detail: %s)\n",
		componentName,
		healthStatus.String(),
		detail,
	)

	oldHealthStatus := g.HealthStatus()
	g.subComponentHealth.Store(componentName, healthstatus.HealthDetail{
		HealthStatus: healthStatus,
		Detail:       detail,
	})
	newHealthStatus := g.HealthStatus()

	if oldHealthStatus != newHealthStatus {
		log.Printf("Info: Health status has changed: %s -> %s\n", oldHealthStatus.String(), newHealthStatus.String())
		g.Publish(g.OverallStatus())
	}

	g.lock.Unlock()
}

func New() StatusManager {
	return &statusManager{
		statusHandlers: list.New(),
		busyStatus:     busystatus.Idle,
	}
}

package healthstatus

type HealthStatus int

type HealthDetail struct {
	HealthStatus
	Detail string
}

type HealthDetails map[string]HealthDetail

type OverallHealth struct {
	HealthStatus
	HealthDetails
}

const (
	Healthy HealthStatus = iota
	Unhealthy
)

func (h HealthStatus) And(other HealthStatus) HealthStatus {
	if h == Healthy && other == Healthy {
		return Healthy
	}

	return Unhealthy
}

func (h HealthStatus) String() string {
	switch h {
	case Healthy:
		return "healthy"
	case Unhealthy:
		return "unhealthy"
	default:
		return "undefined"
	}
}

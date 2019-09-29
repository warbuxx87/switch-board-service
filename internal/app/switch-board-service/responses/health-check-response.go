package responses

type HealthCheckResponse struct {
	CurrentlyHealthy         bool
	AuthenticationErrors     int
	DatabaseConnectionErrors int
}

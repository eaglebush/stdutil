package stdutil

// CommonServiceInfo provides contextual information on new data struct declaration
type CommonServiceInfo struct {
	Login         string // Login of the API user
	ApplicationID string // Application ID of the API used
	ServiceID     string // Service ID of the service
	HelperID      string // Helper ID of the datahelperlite implementation
}

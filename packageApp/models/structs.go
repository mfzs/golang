package models

// User struct for login credentials
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Server struct for endpoint data
type Server struct {
	Name       string
	Env        string
	URL        string
	StatusCode int
}

// InfraResource struct for infra resources
type InfraResource struct {
	ServiceName string
	URL         string
	Owner       string
}

// ClientResources struct for client data
type ClientResources struct {
	ClientName string `json:"client_name"`
	Resources  string `json:"resources"`
}

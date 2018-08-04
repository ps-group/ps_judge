package restapi

// Route - represents single route on server
type Route struct {
	Method  string
	Pattern string
	Handler MethodHandler
}

// RouterConfig - keeps routes, API prefix and other router parameters
type RouterConfig struct {
	Routes    []Route
	APIPrefix string
}

// ServiceConfig - keeps REST API server configuration
type ServiceConfig struct {
	RouterConfig RouterConfig
	ServerURL    string
	LogFileName  string
	// Context - user-defined context object passed to all requests as 1st argument
	Context interface{}
}

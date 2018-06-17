package main

const (
	// BuilderAPIPrefix - API prefix added to all URLS
	BuilderAPIPrefix = "/api/v1"
)

type routeJSON struct {
	Method      string
	Pattern     string
	HandlerFunc APIHandler
}

type routesJSON []routeJSON

// See BuilderAPIPrefix
var jsonRoutes = routesJSON{
	routeJSON{
		"GET",
		"/build/{uuid}",
		getBuildInfo,
	},
	routeJSON{
		"POST",
		"/build/new",
		createBuild,
	},
	routeJSON{
		"POST",
		"/testcase/new",
		createTestCase,
	},
}

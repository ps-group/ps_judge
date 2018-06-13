package main

type routeJSON struct {
	Method      string
	Pattern     string
	HandlerFunc APIHandler
}

type routesJSON []routeJSON

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
}

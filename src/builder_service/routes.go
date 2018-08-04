package main

import (
	"ps-group/restapi"
)

const (
	// BuilderAPIPrefix - API prefix added to all URLS
	BuilderAPIPrefix = "/api/v1"
)

var g_routes = restapi.RouterConfig{
	[]restapi.Route{
		restapi.Route{
			"GET",
			"/build/report/{uuid}",
			getBuildReport,
		},
		restapi.Route{
			"GET",
			"/build/status/{uuid}",
			getBuildStatus,
		},
		restapi.Route{
			"POST",
			"/build/new",
			createBuild,
		},
		restapi.Route{
			"POST",
			"/testcase/new",
			createTestCase,
		},
	},
	BuilderAPIPrefix,
}

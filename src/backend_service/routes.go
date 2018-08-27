package main

import (
	"ps-group/restapi"
)

const (
	// BackendAPIPrefix - API prefix added to all URLS
	BackendAPIPrefix = "/api/v1"
)

var routes = restapi.RouterConfig{
	[]restapi.Route{
		restapi.Route{
			"POST",
			"/user/login",
			toRESTHandler(loginUser),
		},
		restapi.Route{
			"GET",
			"/user/{id}/info",
			toRESTHandler(getUserInfo),
		},
		restapi.Route{
			"GET",
			"/user/{id}/solutions",
			toRESTHandler(getUserSolutions),
		},
		restapi.Route{
			"POST",
			"/user/{id}/commit",
			toRESTHandler(commitSolution),
		},
		restapi.Route{
			"GET",
			"/contest/{id}/assignments",
			toRESTHandler(getContestAssignments),
		},
		restapi.Route{
			"GET",
			"/assignment/{id}",
			toRESTHandler(getAssignmentInfo),
		},
		restapi.Route{
			"POST",
			"/contest/create",
			toRESTHandler(createContest),
		},
		restapi.Route{
			"POST",
			"/user/create",
			toRESTHandler(createUser),
		},
		restapi.Route{
			"POST",
			"/assignment/create",
			toRESTHandler(createAssignment),
		},
		restapi.Route{
			"POST",
			"/testcase/create",
			toRESTHandler(createTestCase),
		},
	},
	BackendAPIPrefix,
}

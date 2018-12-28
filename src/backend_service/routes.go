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
			loginUser,
		},
		restapi.Route{
			"GET",
			"/user/{id}/info",
			getUserInfo,
		},
		restapi.Route{
			"GET",
			"/admin/contests",
			getAdminContestList,
		},
		restapi.Route{
			"GET",
			"/user/{user_id}/contest/list",
			getUserContestList,
		},
		restapi.Route{
			"GET",
			"/user/{user_id}/contest/{contest_id}/solutions",
			getUserContestSolutions,
		},
		restapi.Route{
			"GET",
			"/contest/{id}/results",
			getContestResults,
		},
		restapi.Route{
			"POST",
			"/user/{id}/commit",
			commitSolution,
		},
		restapi.Route{
			"GET",
			"/commit/{id}/report",
			getCommitReport,
		},
		restapi.Route{
			"GET",
			"/contest/{id}/assignments",
			getContestAssignments,
		},
		restapi.Route{
			"GET",
			"/assignment/{id}",
			getAssignmentInfo,
		},
		restapi.Route{
			"POST",
			"/contest/create",
			createContest,
		},
		restapi.Route{
			"POST",
			"/user/create",
			createUser,
		},
		restapi.Route{
			"POST",
			"/assignment/create",
			createAssignment,
		},
		restapi.Route{
			"POST",
			"/testcase/create",
			createTestCase,
		},
	},
	BackendAPIPrefix,
}

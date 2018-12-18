import BackendApi from "../data/backendapi.mjs";

export function getNavbarData(req, res, next) 
{
    const contests = [];

    for (const contest of await backendApi.getUserContestList(userId))
    {
        contests.push({
            'id': verifyInt(contest['id']),
            'title': verifyString(contest['title']),
        });
    }

    req.contests = contests;
}
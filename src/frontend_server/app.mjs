import BackendApi from "./src/data/backendapi.mjs";
import { hashPassword } from "./src/data/password.mjs";
import { verifyInt, verifyString, verifyArray } from "./src/validate.mjs";
import { readConfig } from './src/config.mjs';

import express from 'express';
import bodyParser   from 'body-parser';

import cookieSession from 'cookie-session';

import { checkAuth } from './src/middlewares/auth.mjs'
import { checkRoles } from './src/middlewares/roles.mjs'
import { log } from './src/middlewares/log.mjs'

import { redirectAuthorized } from './src/handlers/main.mjs'
import { convertMarkdown , create_uuid} from './src/handlers/student.mjs'

const app = express();
const SESSION_SECRET = '7pv0OvUy';

const config = readConfig('frontend_server.json');
const backendApi = new BackendApi(config);

// should be adding to middleware's
async function getUserContests(userId)
{
    const contests = [];

    for (const contest of await backendApi.getUserContestList(userId))
    {
        contests.push({
            'id': verifyInt(contest['id']),
            'title': verifyString(contest['title']),
        });
    }

    return contests;
}

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: false,
}));

app.use(cookieSession({
    secret: SESSION_SECRET
}));

app.use(express.static('public'));
app.set('views', 'public');

app.use(checkAuth);

app.get('/', (req, res) => {
    redirectAuthorized(req, res);
});

// handle some URLs
app.post('/login', async(req, res) => {
    const username = req.body['username'];
    const rawPassword = req.body['password'];
    const passwordHash = hashPassword(rawPassword);
    const loginInfo = await backendApi.loginUser(username, passwordHash);

    if (loginInfo.succeed)
    {
        req.session.auth = true;
        req.session.userId = verifyInt(loginInfo['user_id']);
        req.session.username = verifyString(loginInfo['user']['username']);
        req.session.roles = verifyArray(loginInfo['user']['roles']);

        redirectAuthorized(req, res);
    }
    else
    {
        return res.render('tpl/login.ejs', { loginFailed: true });
    }
});

app.get('/login', (req, res) => {
    if (req.session.auth)
    {
        redirectAuthorized(req, res);
    }
    else
    {
        return res.render('tpl/login.ejs', { loginFailed: false });
    }
});

app.get('/student', async(req, res) => {
    const contests = await getUserContests(req.session.userId);

    return res.redirect(`/contest/${contests[0].id}/assignments`);
});

app.get('/contest/:id/assignments', async(req, res) => {
    const contests = await getUserContests(req.session.userId);

    const contestId = parseInt(verifyString(req.params.id));
    const assignments = [];

    const contestTitle = contests.find(contest => contest.id === contestId).title;

    for (const assignment of await backendApi.getContestAssignments(contestId))
    {
        assignments.push({
            'id': verifyInt(assignment['id']),
            'uuid': verifyString(assignment['uuid']),
            'title': verifyString(assignment['title']),
        });
    }
    
    const options = {
        'page': {
            'navbar': {
                'contests': contests,
            },
            'content': {
                'contest': {
                    'id': contestId,
                    'title': contestTitle,
                    'assignments': assignments,
                }
            }
        }
    };

    return res.render('tpl/student/contest_assignments.ejs', options);
});

app.get('/student/contest/:id/solutions', async(req, res) => {
    const contests = await getUserContests(req.session.userId);

    const contestId = parseInt(verifyString(req.params.id));
    const contestTitle = contests.find(contest => contest.id === contestId).title;

    const solutions = [];

    let totalScore = 0;

    for (const solution of await backendApi.getUserContestSolutions(req.session.userId, contestId))
    {
        const score = verifyInt(solution['score']);
        totalScore = totalScore + score;

        solutions.push({
            'assignment_id': verifyInt(solution['assignment_id']),
            'assignment_title': verifyString(solution['assignment_title']),
            'commit_id': verifyInt(solution['commit_id']),
            'score': score,
            'build_status': verifyString(solution['build_status']),
        });
    }
    
    const options = {
        'page': {
            'navbar': {
                'contests': contests,
            },
            'content': {
                'contest': {
                    'id': contestId,
                    'title': contestTitle,
                    'solutions': solutions,
                    'totalScore': totalScore,
                }
            }
        }
    };
    
    return res.render('tpl/student/contest_solutions.ejs', options);
});

app.get('/contest/:contestId/assignment/:assignmentId', async(req, res) => {
    const contests = await getUserContests(req.session.userId);

    const contestId = parseInt(verifyString(req.params.contestId));
    const assignmentId = parseInt(verifyString(req.params.assignmentId));
    const info = await backendApi.getAssignmentInfo(assignmentId);
    
    const articleHtml = convertMarkdown(verifyString(info['description']));

    const options = {
        'page': {
            'navbar': {
                'contests': contests,
            },
            'content': {
                'contestId': contestId,
                'assignment': {
                    'id': verifyInt(info['id']),
                    'title': verifyString(info['title']),
                    'info': articleHtml,
                }
            }
        }
    };

    return res.render('tpl/assignment.ejs', options);
});

app.post('/student/commit', async(req, res) => {
    const userId = verifyInt(req.session.userId);
    const assignmentId = parseInt(verifyString(req.body['assignmentId']));
    const contestId = parseInt(verifyString(req.body['contestId']));
    const language = verifyString(req.body['language']);
    const source = verifyString(req.body['source']);
    const uuid = create_uuid();

    await backendApi.commitSolution(userId, uuid, assignmentId, language, source);
    await res.redirect(`/student/contest/${contestId}/solutions`);
});

app.get('/student/commit/:commitId', async(req, res) => {
    const contests = await getUserContests(req.session.userId);

    const commitId = parseInt(verifyString(req.params.commitId));

    const report = await backendApi.getCommitReport(commitId);

    const options = {
        'page': {
            'navbar': {
                'contests': contests,
            },
            'content': {
                'commit' : {
                    'status': verifyString(report['status']),
                    'buildLog': verifyString(report['build_log']),
                    'testsLog': verifyString(report['tests_log']),
                    'testsPassed': verifyInt(report['tests_passed']),
                    'testsTotal': verifyInt(report['tests_total']),   
                }
            }
        }
    };

    return res.render('tpl/student/commit.ejs', options);
});

const server = app.listen(config.port, (error) => {
    if (error) return console.log(`Error: ${error}`);

    console.log(`Server listening on port ${server.address().port}`);
});

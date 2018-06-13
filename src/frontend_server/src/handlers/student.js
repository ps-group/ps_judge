const util = require('util');
const ejs = require('ejs');
const url = require('url');
const basehandler = require('./basehandler');
const showdown  = require('showdown');
const uuidv1 = require('uuid/v1');
const assert = require('assert');
const routes = require('../routes');

class Student extends basehandler.BaseHandler
{
    /**
     * @param {context.Context} context
     */
    constructor(context)
    {
        super(context);
    }

    async home(request, response) 
    {
        const userInfo = await this._fetchUser(request, response);
        if (userInfo != null)
        {
            const contestId = userInfo['active_contest_id'];
            const assignments = await this._repository.getAssignmentsBriefInfo(contestId);
            const options = {
                'page': {
                    'assignments': assignments
                }
            };
            return this._render('./tpl/student_home.ejs', options, response);
        }
    }

    async assignment(request, response) 
    {
        const userInfo = await this._fetchUser(request, response);
        if (userInfo != null)
        {
            const query = url.parse(request.url, true).query;
            const assignmentId = parseInt(query.id);
            const assignment = await this._repository.getAssignmentInfo(assignmentId);

            // TODO: we can speedup markdown convertion by moving it to client side.
            const articleHtml = this._convertMarkdown(assignment.article);
            const options = {
                'page': {
                    'assignment_id': assignmentId,
                    'title': assignment.title,
                    'assignment_info': articleHtml
                }
            };
            return this._render('./tpl/student_assignment.ejs', options, response);
        }
    }

    async commit(request, response) 
    {
        const userInfo = await this._fetchUser(request, response);
        if (userInfo)
        {
            const assignmentId = request.body.assignmentId;
            const source = request.body.source;
            const uuid = uuidv1().split('-').join('');
            assert(source);
            assert(assignmentId);
            assert(uuid);

            this._initRepository();
            let solutionInfo = await this._repository.getSolutionInfo(userInfo.id, assignmentId);
            if (solutionInfo == null)
            {
                await this._repository.createSolution(userInfo.id, assignmentId);
                solutionInfo = await this._repository.getSolutionInfo(userInfo.id, assignmentId);
            }

            this._repository.createCommit(solutionInfo.id, uuid, source);
            this._redirect(routes.STUDENT_SOLUTION_URL + '?id=' + solutionInfo.id, response);
        }
    }

    _convertMarkdown(markdown)
    {
        const converter = new showdown.Converter();
        return converter.makeHtml(markdown);
    }

    _fetchUser(request, response)
    {
        if (this._checkAuth(request, response))
        {
            this._initSession();
            this._initRepository();
            return this._repository.getUserAuthInfo(this._session.username);
        }
        return null;
    }
}

module.exports = Student;

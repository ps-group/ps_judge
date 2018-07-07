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
     * @param {Request} request
     * @param {Response} response
     */
    constructor(context, request, response)
    {
        super(context, request, response);
    }

    async home() 
    {
        const userInfo = await this._fetchUser();
        if (userInfo != null)
        {
            const contestId = userInfo['active_contest_id'];
            const assignments = await this.repository.getAssignmentsBriefInfo(contestId);
            const options = {
                'page': {
                    'assignments': assignments
                }
            };
            return this._render('./tpl/student_home.ejs', options);
        }
    }

    async assignment() 
    {
        const userInfo = await this._fetchUser();
        if (userInfo != null)
        {
            const query = url.parse(this.request.url, true).query;
            const assignmentId = parseInt(query.id);
            const assignment = await this.repository.getAssignmentInfo(assignmentId);

            // TODO: we can speedup markdown convertion by moving it to client side.
            const articleHtml = this._convertMarkdown(assignment.article);
            const options = {
                'page': {
                    'assignment_id': assignmentId,
                    'title': assignment.title,
                    'assignment_info': articleHtml
                }
            };
            return this._render('./tpl/student_assignment.ejs', options);
        }
    }

    async commit() 
    {
        const userInfo = await this._fetchUser();
        if (userInfo)
        {
            const assignmentId = this.request.body.assignmentId;
            const source = this.request.body.source;
            const uuid = uuidv1().split('-').join('');
            assert(source);
            assert(assignmentId);
            assert(uuid);

            let solutionInfo = await this.repository.getSolutionInfo(userInfo.id, assignmentId);
            if (solutionInfo == null)
            {
                await this.repository.createSolution(userInfo.id, assignmentId);
                solutionInfo = await this.repository.getSolutionInfo(userInfo.id, assignmentId);
            }

            await this.repository.createCommit(solutionInfo.id, uuid, source);
            await this._redirect(routes.STUDENT_SOLUTION_URL + '?id=' + solutionInfo.id);
        }
    }

    _convertMarkdown(markdown)
    {
        const converter = new showdown.Converter();
        return converter.makeHtml(markdown);
    }

    async _fetchUser()
    {
        const checked = await this._checkAuth()
        if (checked)
        {
            return this.repository.getUserAuthInfo(this.session.username);
        }
        return null;
    }
}

module.exports = Student;

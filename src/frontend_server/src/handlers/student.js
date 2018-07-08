const url = require('url');
const basehandler = require('./basehandler');
const showdown  = require('showdown');
const uuidv1 = require('uuid/v1');
const assert = require('assert');
const routes = require('../routes');
const repository = require('../data/repository');

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
        const userInfo = await this._fetchStudentUser();
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
        const userInfo = await this._fetchStudentUser();
        if (userInfo != null)
        {
            const query = url.parse(this.request.url, true).query;
            const assignmentId = parseInt(query["id"]);
            const assignment = await this.repository.getAssignmentFullInfo(assignmentId);

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

    async solutions()
    {
        const userInfo = await this._fetchStudentUser();
        if (userInfo != null)
        {
            const contestId = userInfo['active_contest_id'];
            const solutions = await this.repository.getUserSolutions(userInfo.id);
            const assignments = await this.repository.getAssignmentsBriefInfo(contestId);
            const assignmentTitles = {};
            for (let assignment of assignments)
            {
                assignmentTitles[assignment['id']] = assignment['title'];
            }

            const infos = [];
            for (let solution of solutions)
            {
                const solutionId = solution['id'];
                const assignmentId = solution['assignment_id'];
                const commitInfo = await this.repository.getLastCommitInfo(solutionId);
                const info = {
                    'assignment_id': assignmentId,
                    'assignment_title': assignmentTitles[assignmentId],
                    'score': solution['score'],
                    'commit_id': commitInfo['id'],
                    'build_status': commitInfo['build_status'],
                };
                infos.push(info);
            }

            const options = {
                'page': {
                    'solutions': infos
                }
            };
            return this._render('./tpl/student_solutions.ejs', options);
        }
    }

    async commit() 
    {
        const userInfo = await this._fetchStudentUser();
        if (userInfo)
        {
            const assignmentId = this.request.body.assignmentId;
            const language = this._validateLanguage(this.request.body.language);
            const source = this.request.body.source;
            const uuid = this._create_uuid();
            assert(source);
            assert(assignmentId);
            assert(uuid);

            let solutionInfo = await this.repository.getSolutionInfo(userInfo.id, assignmentId);
            if (solutionInfo == null)
            {
                await this.repository.createSolution(userInfo.id, assignmentId);
                solutionInfo = await this.repository.getSolutionInfo(userInfo.id, assignmentId);
            }
            const assignmentInfo = await this.repository.getAssignmentFullInfo(assignmentId);

            await this.repository.createCommit(solutionInfo.id, uuid);
            await this.backendApi.registerNewBuild(uuid, assignmentInfo.uuid, language, source);
            await this._redirect(routes.STUDENT_SOLUTIONS_URL);
        }
    }

    /**
     * Returns if language valid or throws
     * @param {string} language - language enumeration string
     * @returns {string} returned string is the same as the language param
     */
    _validateLanguage(language)
    {
        const languages = ['c++', 'pascal'];
        if (languages.includes(language))
        {
            return language;
        }
        throw new Error(`unknown language: ${language}`);
    }

    _convertMarkdown(markdown)
    {
        const converter = new showdown.Converter();
        return converter.makeHtml(markdown);
    }

    /**
     * @returns {string} new uuid
     */
    _create_uuid()
    {
        return uuidv1().split('-').join('');
    }

    /**
     * Fetches current user info and ensures it has student role.
     * Returns null if user is not authorized.
     */
    async _fetchStudentUser()
    {
        const userInfo = await this._fetchUser();
        if (userInfo)
        {
            const roles = '' + userInfo['roles'];
            if (roles.indexOf(repository.ROLE_STUDENT) < 0)
            {
                throw new Error('user has no privilegies to view this page');
            }
        }
        return userInfo;
    }
}

module.exports = Student;

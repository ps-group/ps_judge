import { parse } from 'url';
import { BaseHandler } from './basehandler.mjs';
import showdown from 'showdown';
import uuidv1 from 'uuid/v1';
import assert from 'assert';
import { STUDENT_SOLUTIONS_URL } from '../routes';
import { ROLE_STUDENT } from '../data/roles';
import { verifyInt, verifyString } from '../validate.mjs';

export default class StudentHandler extends BaseHandler
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
        if (!await this._fetchStudentUser())
        {
            return;
        }

        const userInfo = await this._backend.getUserInfo(this.session.userId);
        const infos = [];
        for (const assignment of await this._backend.getContestAssignments(verifyInt(userInfo['contest_id'])))
        {
            infos.push({
                'id': verifyInt(assignment['id']),
                'contest_id': verifyInt(assignment['contest_id']),
                'uuid': verifyString(assignment['uuid']),
                'title': verifyString(assignment['title']),
            });
        }
        const options = {
            'page': {
                'assignments': infos
            }
        };
        return this._render('./tpl/student_home.ejs', options);
    }

    async assignment() 
    {
        if (!await this._fetchStudentUser())
        {
            return;
        }

        const query = parse(this.request.url, true).query;
        const assignmentId = verifyInt(parseInt(query["id"]));
        const info = await this._backend.getAssignmentInfo(assignmentId);
        const articleHtml = this._convertMarkdown(verifyString(info['description']));
        const options = {
            'page': {
                'assignment_id': verifyInt(info['id']),
                'title': verifyString(info['title']),
                'assignment_info': articleHtml
            }
        };
        return this._render('./tpl/student_assignment.ejs', options);
    }

    async solutions()
    {
        if (!await this._fetchStudentUser())
        {
            return;
        }

        const infos = [];
        for (let solution of await this._backend.getUserSolutions(this.session.userId))
        {
            infos.push({
                'assignment_id': verifyInt(solution['assignment_id']),
                'assignment_title': verifyString(solution['assignment_title']),
                'commit_id': verifyInt(solution['commit_id']),
                'score': verifyInt(solution['score']),
                'build_status': verifyString(solution['build_status']),
            });
        }
        const options = {
            'page': {
                'solutions': infos
            }
        };
        return this._render('./tpl/student_solutions.ejs', options);
    }

    async commit() 
    {
        if (!await this._fetchStudentUser())
        {
            return;
        }

        const userId = verifyInt(this.session.userId);
        const assignmentId = parseInt(verifyString(this.request.body['assignmentId']));
        const language = verifyString(this.request.body['language']);
        const source = verifyString(this.request.body['source']);
        const uuid = this._create_uuid();

        await this._backend.commitSolution(userId, uuid, assignmentId, language, source);
        await this._redirect(STUDENT_SOLUTIONS_URL);
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
     * Ensures that user authorized and has student role.
     * Redirects to login and returns false if user is not authorized.
     */
    async _fetchStudentUser()
    {
        const authorized = await this._checkAuth();
        if (authorized)
        {
            const roles = this.session.roles;
            if (roles.indexOf(ROLE_STUDENT) < 0)
            {
                throw new Error('user has no privilegies to view this page');
            }
            return true;
        }
        return false;
    }
}

const connector = require('./connector.js')

const ROLE_ADMIN = 'admin';
const ROLE_STUDENT = 'student';
const ROLE_JUDGE = 'judge';
const assert = require('assert');

/**
 * This repository class can access frontend server database.
 */
class FrontendRepository
{
    /**
     * Creates repository which can access frontend server database.
     * @param {connector.Connector} connector
     */
    constructor(connector)
    {
        assert(connector !== undefined);
        this.connector = connector;
    }

    /**
     * Creates new user.
     * @param {string} name
     * @param {string} password
     * @param {string} role
     */
    async createUser(name, password, role)
    {
        const knownRoles = { ROLE_ADMIN, ROLE_STUDENT, ROLE_JUDGE };
        if (knownRoles.indexOf(role) == -1)
        {
            throw new Error('unknown user role argument: ' + role);
        }
        const sql = 'INSERT INTO user (username, password, roles) VALUES (?, ?, ?)';
        await this.connector.query(sql, [name, password, role]);
    }

    /**
     * Creates new contest
     * @param {string} title - short contest title
     * @param {Date} startTime - contest start time
     * @param {Date} endTime - must be greater than startTime
     */
    async createContest(title, startTime, endTime)
    {
        if (endTime.getTime() <= startTime.getTime())
        {
            throw new Error('endTime cannot be greater than startTime');
        }
        const sql = 'INSERT INTO contest (title, start_time, end_time) VALUES (?, ?, ?)';
        await this.connector.query(sql, [title, startTime, endTime]);
    }

    /**
     * Creates new assignment in given contest.
     * @param {number} contestId - database id of assignment contest
     * @param {string} title - short title
     * @param {string} article - article text in Markdown format
     */
    async createAssignment(contestId, title, article)
    {
        const sql = 'INSERT INTO assignment (contest_id, title, article) VALUES (?, ?, ?)';
        await this.connector.query(sql, [contestId, title, article]);
    }

    /**
     * Creates new solution, which connects user and assignment
     * @param {number} userId - database id of user which owns solution
     * @param {number} assignmentId - database id of assignment which needs this solution
     */
    async createSolution(userId, assignmentId)
    {
        const sql = 'INSERT INTO solution (user_id, assignment_id, score) VALUES (?, ?, ?)';
        const score = 0;
        await this.connector.query(sql, [userId, assignmentId, score]);
    }

    /**
     * Creates new solution source code commit
     * @param {number} solutionId - database id of commit solution
     * @param {string} uuid - global unique identifier of given commit
     * @param {string} source - commit source code
     */
    async createCommit(solutionId, uuid, source)
    {
        const sql = 'INSERT INTO commit (solution_id, uuid, source) VALUES (?, ?, ?)';
        await this.connector.query(sql, [solutionId, uuid, source]);
    }

    /**
     * Creates new commit review
     * @param {number} commitId - database id of commit
     * @param {number} reviewerId - database id of user which does review
     */
    async createReview(commitId, reviewerId)
    {
        const sql = 'INSERT INTO review (commit_id, reviewer_id, score) VALUES (?, ?, ?)';
        const score = 0;
        await this.connector.query(sql, [commitId, reviewerId, score]);
    }

    /**
     * Returns info from contest with given id.
     * @param {number} contestId - database id of the programming contest
     */
    async getContestInfo(contestId)
    {
        const sql = 'SELECT title, start_time, end_time FROM contest WHERE id = ?';
        return await this.connector.query(sql, [contestId]);
    }

    /**
     * Retuns assignments attached to given contest.
     * @param {number} contestId - database id of the programming context
     */
    async getAssignmentsBriefInfo(contestId)
    {
        const sql = 'SELECT id, title FROM assignment WHERE contest_id = ?';
        return await this.connector.query(sql, [contestId]);
    }

    /**
     * Return full information about assignment.
     * @param {*} assignmentId - database id of the assignment
     */
    async getAssignmentInfo(assignmentId)
    {
        const sql = 'SELECT id, title, article FROM assignment WHERE id = ?';
        const infos = await this.connector.query(sql, [assignmentId]);
        if (infos.length == 0)
        {
            return null;
        }
        return infos[0];
    }

    /**
     * Returns id, password, contest_id and role for user with given usernam
     * @param {string} username - username used to search
     */
    async getUserAuthInfo(username)
    {
        const sql = 'SELECT id, password, active_contest_id, roles FROM user WHERE username=?';
        const infos = await this.connector.query(sql, [username]);
        if (infos.length == 0)
        {
            return null;
        }
        assert(infos.length == 1);
        return infos[0];
    }

    /**
     * Returns solution information
     * @param {number} userId - user which made solution
     * @param {number} assignmentId - solution target assignment
     */
    async getSolutionInfo(userId, assignmentId)
    {
        const sql = 'SELECT id, score FROM solution WHERE user_id=? AND assignment_id=?';
        const infos = await this.connector.query(sql, [userId, assignmentId]);
        if (infos.length == 0)
        {
            return null;
        }
        assert(infos.length == 1);
        return infos[0];
    }
}

module.exports.ROLE_ADMIN = ROLE_ADMIN;
module.exports.ROLE_STUDENT = ROLE_STUDENT;
module.exports.ROLE_JUDGE = ROLE_JUDGE;
module.exports.FrontendRepository = FrontendRepository;

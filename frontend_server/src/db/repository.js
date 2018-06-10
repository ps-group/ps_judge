const connector = require('./connector.js')

const ROLE_ADMIN = 'admin';
const ROLE_STUDENT = 'student';
const ROLE_JUDGE = 'judge';

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
}

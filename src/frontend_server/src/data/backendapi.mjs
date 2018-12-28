import assert from 'assert';
import { ApiClient } from './apiclient.mjs';
import { verifyInt } from '../validate.mjs';
import { verifyString } from '../validate.mjs';

/**
 * @typedef {Object} UserInfo
 * @property {string} username
 * @property {Array<string>} roles
 * @property {number} contest_id
 */

/**
 * @typedef {Object} LoginInfo
 * @property {boolean} succeed
 * @property {number} user_id
 * @property {UserInfo} user
 */

/**
 * @typedef {Object} ContestInfo
 * @property {number} id
 * @property {string} title
 */

/**
 * @typedef {Object} BriefSolutionInfo
 * @property {number} assignment_id
 * @property {string} assignment_title
 * @property {number} score
 * @property {number} commit_id
 * @property {string} build_status
 */

 /**
  * @typedef {Object} AssignmentInfo
  * @property {number} id
  * @property {number} contest_id
  * @property {string} uuid
  * @property {string} title
  */

 /**
  * @typedef {Object} FullAssignmentInfo
  * @property {number} id
  * @property {number} contest_id
  * @property {string} uuid
  * @property {string} title
  * @property {string} description
  */

  /**
   * @typedef {Object} CommitReport
   * @property {string} uuid
   * @property {string} status
   * @property {string} exception
   * @property {string} build_log
   * @property {string} tests_log
   * @property {number} tests_passed
   * @property {number} tests_total
   */

   /**
    * @typedef {Object} ContestResult
    * @property {string} username
    * @property {number} score
    * @property {string} assignment_title
    */

export default class BackendApi
{
    /**
     * @param {config.Config} config
     */
    constructor(config)
    {
        this._client = new ApiClient(config.backendURL);
    }

    /**
     * Asks backend to login user
     * @param {string} usename
     * @param {string} passwordHash
     * @returns {!Promise<LoginInfo>}
     */
    loginUser(username, passwordHash)
    {
        const params = {
            'username': username,
            'password_hash': passwordHash,
        };
        return this._client.sendPost('user/login', params);
    }

    /**
     * Queries user information by user ID.
     * @param {number} userId - unique user ID
     * @returns {!Promise<UserInfo>}
     */
    getUserInfo(userId)
    {
        userId = verifyInt(userId);
        return this._client.sendGet(`user/${userId}/info`);
    }

    /**
     * Queries user solutions list.
     * @param {number} userId - unique user ID
     * @returns {!Promise<Array<BriefSolutionInfo>>}
     */
    getUserSolutions(userId)
    {
        userId = verifyInt(userId);
        return this._client.sendGet(`user/${userId}/solutions`);
    }

    /**
     * Commits solution source code to the backend.
     * @param {number} userId
     * @param {string} uuid - uuid of build
     * @param {number} assignmentId
     * @param {string} language
     * @param {string} source
     * @returns {Promise<undefined>}
     */
    async commitSolution(userId, uuid, assignmentId, language, source)
    {
        userId = verifyInt(userId);
        const response = await this._client.sendPost(`user/${userId}/commit`, {
            'uuid': uuid,
            'assignment_id': assignmentId,
            'language': language,
            'source': source,
        });
        assert(verifyString(response['uuid']) == uuid);
    }

    /**
     * @returns {!Promise<Array<ContestInfo>>}
     */
    getAdminContestList()
    {
        return this._client.sendGet(`admin/contests`);
    }

    /**
     * @param userId 
     * @returns {!Promise<Array<ContestInfo>>}
     */
    getUserContestList(userId)
    {
        userId = verifyInt(userId);
        return this._client.sendGet(`user/${userId}/contest/list`);
    }

    /**
     * @param {number} userId
     * @param {number} contestId 
     * @returns {!Promise<Array<BriefSolutionInfo>>}
     */
    getUserContestSolutions(userId, contestId)
    {
        userId = verifyInt(userId);
        contestId = verifyInt(contestId);
        return  this._client.sendGet(`user/${userId}/contest/${contestId}/solutions`)
    }

    /**
     * Queries list of assignments for the given contest.
     * @param {number} contestId
     * @returns {Array<AssignmentInfo>} list of assignments
     */
    async getContestAssignments(contestId)
    {
        contestId = verifyInt(contestId);
        return await this._client.sendGet(`contest/${contestId}/assignments`);
    }

    /**
     * Queries detailed information about given assignment
     * @param {number} assignmentId
     * @returns {FullAssignmentInfo} assignment info
     */
    async getAssignmentInfo(assignmentId)
    {
        assignmentId = verifyInt(assignmentId);
        return await this._client.sendGet(`assignment/${assignmentId}`);
    }

    /**
     * @param {number} commitId
     * @returns {CommitReport}
     */
    async getCommitReport(commitId)
    {
        commitId = verifyInt(commitId);
        return await this._client.sendGet(`commit/${commitId}/report`)
    }

    /**
     * @param {number} contestId
     * @returns {<Array<ContestResult>>}
     */
    async getContestResults(contestId)
    {
        contestId = verifyInt(contestId);
        return await this._client.sendGet(`contest/${contestId}/results`)
    }

    /**
     * 
     * @param {string} title
     * @param {number} maxReviews 
     */
    async createContest(title, maxReviews)
    {
        const params = {
            'title': title,
            'max_reviews': parseInt(maxReviews),
        };

        const contestId = await this._client.sendPost('contest/create', params);
    }

    /**
     * 
     * @param {string} uuid 
     * @param {string} title 
     * @param {number} contestId 
     * @param {string} description 
     */
    async createAssignment(uuid, title, contestId, description)
    {
        const params = {
            'uuid': uuid,
            'contest_id': contestId,
            'title': title,
            'description': description,
        }

        const assignmentId = await this._client.sendPost('assignment/create', params);
    }
}

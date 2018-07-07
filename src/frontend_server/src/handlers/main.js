const util = require('util');
const ejs = require('ejs');
const assert = require('assert');
const basehandler = require('./basehandler');
const repository = require('../data/repository');
const appsession = require('../data/appsession');
const routes = require('../routes');
const password = require('../data/password');

class Main extends basehandler.BaseHandler
{
    /**
     * @param {context.Context} context
     */
    constructor(context)
    {
        super(context);
    }

    async index(request, response) 
    {
        const checked = await this._checkAuth(request, response)
        if (checked)
        {
            return this._redirectAuthorized(request, response);
        }
    }

    async login(request, response)
    {
        if (request.method != 'POST')
        {
            if (this._hasAuth(request, response))
            {
                this._redirectAuthorized(request, response);
            }
            else
            {
                return this._render('./tpl/login.ejs', { loginFailed: false }, response);
            }
        }
        else
        {
            const username = request.body['username'];
            const rawPassword = request.body['password'];
            const passwordHash = password.hashPassword(rawPassword);

            // TODO: do not compare with rawPassword - it's unsafe.
            this._initRepository();
            const info = await this._repository.getUserAuthInfo(username);
            if (info != null && (info['password'] == rawPassword || info['password'] == passwordHash))
            {
                this._initSession(request);
                this._session.authorized = true;
                this._session.username = username;
                return this._redirectAuthorized(request, response);
            }
            else
            {
                return this._render('./tpl/login.ejs', { loginFailed: true }, response);
            }
        }
    }

    async _redirectAuthorized(request, response)
    {
        this._initSession(request);
        this._initRepository();

        const info = await this._repository.getUserAuthInfo(this._session.username);
        const roles = '' + info['roles'];
        if (roles.indexOf(repository.ROLE_ADMIN) >= 0)
        {
            return this._redirect(routes.ADMIN_HOME_URL, response);
        }
        else if (roles.indexOf(repository.ROLE_JUDGE) >= 0)
        {
            return this._redirect(routes.JUDGE_HOME_URL, response);
        }
        else if (roles.indexOf(repository.ROLE_STUDENT) >= 0)
        {
            return this._redirect(routes.STUDENT_HOME_URL, response);
        }
        else
        {
            throw new Error('user has incorrect roles: ' + roles);
        }
    }
}

module.exports = Main;

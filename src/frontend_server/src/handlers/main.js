const basehandler = require('./basehandler');
const repository = require('../data/repository');
const routes = require('../routes');
const password = require('../data/password');

class Main extends basehandler.BaseHandler
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

    async index() 
    {
        const checked = await this._checkAuth()
        if (checked)
        {
            return this._redirectAuthorized();
        }
    }

    async login()
    {
        if (this.request.method != 'POST')
        {
            if (this._hasAuth(this.request))
            {
                this._redirectAuthorized(this.request);
            }
            else
            {
                return this._render('./tpl/login.ejs', { loginFailed: false });
            }
        }
        else
        {
            const username = this.request.body['username'];
            const rawPassword = this.request.body['password'];
            const passwordHash = password.hashPassword(rawPassword);

            const info = await this.repository.getUserAuthInfo(username);
            if (info != null && (info['password'] == passwordHash))
            {
                this.session.authorized = true;
                this.session.username = username;
                return this._redirectAuthorized();
            }
            else
            {
                return this._render('./tpl/login.ejs', { loginFailed: true });
            }
        }
    }

    async _redirectAuthorized()
    {
        const info = await this.repository.getUserAuthInfo(this.session.username);
        const roles = '' + info['roles'];
        if (roles.indexOf(repository.ROLE_ADMIN) >= 0)
        {
            return this._redirect(routes.ADMIN_HOME_URL);
        }
        else if (roles.indexOf(repository.ROLE_JUDGE) >= 0)
        {
            return this._redirect(routes.JUDGE_HOME_URL);
        }
        else if (roles.indexOf(repository.ROLE_STUDENT) >= 0)
        {
            return this._redirect(routes.STUDENT_HOME_URL);
        }
        else
        {
            throw new Error('user has incorrect roles: ' + roles);
        }
    }
}

module.exports = Main;

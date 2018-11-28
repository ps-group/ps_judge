import { BaseHandler } from './basehandler.mjs';
import { ROLE_ADMIN, ROLE_JUDGE, ROLE_STUDENT } from '../data/roles.mjs';
import { ADMIN_HOME_URL, JUDGE_HOME_URL, STUDENT_HOME_URL } from '../routes';
import { hashPassword } from '../data/password.mjs';
import { verifyInt, verifyString, verifyArray } from '../validate.mjs';

export default class MainHandler extends BaseHandler
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
        console.log('On index!');
        const checked = await this._checkAuth();
        if (checked)
        {
            return this._redirectAuthorized();
        }
    }

    async login()
    {
        //console.log('On login!');
        if (this.request.method != 'POST')
        {
            //console.log('Post request!');
            if (this.session.authorized)
            {
                console.log('Authtorized!');
                this._redirectAuthorized(this.request);
            }
            else
            {
                console.log('Not authtorized!');
                return this._render('./tpl/login.ejs', { loginFailed: false });
            }
        }
        else
        {
            const username = this.request.body['username'];
            const rawPassword = this.request.body['password'];
            const passwordHash = hashPassword(rawPassword);
            const loginInfo = await this._backend.loginUser(username, passwordHash);
            if (loginInfo.succeed)
            {
                this.session.authorized = true;
                this.session.userId = verifyInt(loginInfo['user_id']);
                this.session.username = verifyString(loginInfo['user']['username']);
                this.session.roles = verifyArray(loginInfo['user']['roles']);
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
        const roles = this.session.roles;
        if (roles.indexOf(ROLE_ADMIN) >= 0)
        {
            return this._redirect(ADMIN_HOME_URL);
        }
        else if (roles.indexOf(ROLE_JUDGE) >= 0)
        {
            return this._redirect(JUDGE_HOME_URL);
        }
        else if (roles.indexOf(ROLE_STUDENT) >= 0)
        {
            return this._redirect(STUDENT_HOME_URL);
        }
        else
        {
            throw new Error('user has incorrect roles: ' + roles);
        }
    }
}

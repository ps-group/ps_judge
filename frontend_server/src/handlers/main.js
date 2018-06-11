const util = require('util');
const ejs = require('ejs');
const basehandler = require('./basehandler');
const repository = require('../data/repository');
const appsession = require('../data/appsession');
const routes = require('../routes');

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
        if (this._checkAuth(request, response))
        {
            const repo = this._context.connectDB();
            const session = new appsession.AppSession(request);
            const info = repo.getUserAuthInfo(session.username);
            if (info.roles.contains(repository.ROLE_ADMIN))
            {
                this._redirect(routes.ADMIN_HOME_URL)
            }
            else if (info.roles.contains(repository.ROLE_JUDGE))
            {
                this._redirect(routes.JUDGE_HOME_URL)
            }
            else if (info.roles.contains(repository.ROLE_STUDENT))
            {
                this._redirect(routes.STUDENT_HOME_URL)
            }
            throw new Error('user has incorrect roles: ' + info.roles);
        }
    }
}

module.exports = Main;

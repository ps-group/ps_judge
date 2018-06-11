const util = require('util');
const ejs = require('ejs');
const routes = require('../routes');
const appsession = require('../data/appsession');
const assert = require('assert');

const renderFileAsync = util.promisify(ejs.renderFile);

class BaseHandler
{
    /**
     * @param {context.Context} context
     */
    constructor(context)
    {
        this._context = context;
    }

    /**
     * Renders HTML page
     * @param {*} tplPath - path to template which should be rendered
     * @param {Object} data - key/value mapping for page data
     * @param {*} response - server response object
     */
    async _render(tplPath, data, response)
    {
        assert(data !== undefined);
        assert(response !== undefined);

        const html = await renderFileAsync(tplPath, data);
        response.writeHead(200, {'Content-Type': 'text/html'});
        response.write(html);
        response.end();
    }

    /**
     * Redirects user to given URL
     * @param {string} url - URL which should be opened instead of current URL.
     * @param {*} response - server response object
     */
    async _redirect(url, response)
    {
        response.writeHead(301, { 'Location': url });
        response.end();
    }

    /**
     * Checks if user authorized and redirects to login page if
     * @param {*} request - client request object
     * @param {*} response - server response object
     */
    _checkAuth(request, response, next)
    {
        const session = new appsession.AppSession(request);
        if (!session.authorized)
        {
            this._redirect(routes.LOGIN_URL, response);
            return false;
        }
        return true;
    }
}

module.exports.BaseHandler = BaseHandler;

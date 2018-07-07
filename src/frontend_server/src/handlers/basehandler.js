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
        /**
         * @property {repository.FrontendRepository}
         */
        this._repository = null;
        /**
         * @property {appsession.AppSession}
         */
        this._session = null;
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
        response.writeHead(301, { 'Location': url, 'Cache-Control': 'no-store' });
        response.end();
    }

    /**
     * Returns true if user authorized
     * @param {*} request - client request object
     * @param {*} response - server response object
     */
    _hasAuth(request)
    {
        this._initSession(request);
        return this._session.authorized;
    }

    /**
     * Checks if user authorized and redirects to login page if
     * @param {*} request - client request object
     * @param {*} response - server response object
     */
    async _checkAuth(request, response)
    {
        if (!this._hasAuth(request))
        {
            await this._redirect(routes.LOGIN_URL, response);
            return false;
        }
        return true;
    }

    /**
     * Initializes user session lazily
     */
    _initSession(request)
    {
        if (!this._session)
        {
            this._session = new appsession.AppSession(request);
        }
        assert(this._session);
    }

    /**
     * Initializes repository lazily.
     */
    _initRepository()
    {
        if (!this._repository)
        {
            this._repository = this._context.connectDB();
        }
        assert(this._repository);
    }
}

module.exports.BaseHandler = BaseHandler;

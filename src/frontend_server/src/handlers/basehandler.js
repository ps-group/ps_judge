const util = require('util');
const ejs = require('ejs');
const routes = require('../routes');
const assert = require('assert');
const backendapi = require('../data/backendapi');

const renderFileAsync = util.promisify(ejs.renderFile);

class BaseHandler
{
    /**
     * @param {context.Context} context
     * @param {Request} request
     * @param {Response} response
     */
    constructor(context, request, response)
    {
        /**
         * @property {context.Context}
         */
        this._context = context;
        /**
         * @property {http.Request}
         */
        this._request = request;
        /**
         * @property {http.Response}
         */
        this._response = response;
        /**
         * @property {repository.FrontendRepository}
         */
        this._repository = null;
        /**
         * @property {backendapi.BackendApi}
         */
        this._backendApi = null;
    }

    /**
     * Renders HTML page
     * @param {string} tplPath - path to template which should be rendered
     * @param {Object} data - key/value mapping for page data
     */
    async _render(tplPath, data)
    {
        assert(data !== undefined);

        const html = await renderFileAsync(tplPath, data);
        this._response.writeHead(200, {'Content-Type': 'text/html'});
        this._response.write(html);
        this._response.end();
    }

    /**
     * Redirects user to given URL
     * @param {string} url - URL which should be opened instead of current URL.
     * @param {*} response - server response object
     */
    async _redirect(url)
    {
        this._response.writeHead(301, { 'Location': url, 'Cache-Control': 'no-store' });
        this._response.end();
    }

    /**
     * Returns true if user authorized
     */
    _hasAuth()
    {
        return this.session.authorized;
    }

    /**
     * Checks if user authorized and redirects to login page if
     */
    async _checkAuth()
    {
        if (!this._hasAuth())
        {
            await this._redirect(routes.LOGIN_URL);
            return false;
        }
        return true;
    }

    /**
     * Fetches current user info or returns null if user is not authorized.
     */
    async _fetchUser()
    {
        const checked = await this._checkAuth()
        if (checked)
        {
            return this.repository.getUserAuthInfo(this.session.username);
        }
        return null;
    }

    /**
     * @returns {Request}
     */
    get request()
    {
        return this._request;
    }

    /**
     * @returns {Response}
     */
    get response()
    {
        return this._response;
    }

    /**
     * Returns user session object.
     */
    get session()
    {
        return this._request.session;
    }

    /**
     * @returns {backendapi.BackendApi}
     */
    get backendApi()
    {
        if (!this._backendApi)
        {
            this._backendApi = new backendapi.BackendApi(this._context.config);
        }
        return this._backendApi;
    }

    /**
     * Initializes repository lazily.
     * @returns {repository.FrontendRepository}
     */
    get repository()
    {
        if (!this._repository)
        {
            this._repository = this._context.connectDB();
        }
        assert(this._repository);
        return this._repository;
    }
}

module.exports.BaseHandler = BaseHandler;

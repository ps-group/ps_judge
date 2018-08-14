import { promisify } from 'util';
import ejs from 'ejs';
import { LOGIN_URL } from '../routes';
import assert from 'assert';
import BackendApi from '../data/backendapi.mjs';

const renderFileAsync = promisify(ejs.renderFile);

export class BaseHandler
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
         * @property {BackendApi}
         */
        this._backend = this._context.backend;
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
     * Checks if user authorized and redirects to login page if
     */
    async _checkAuth()
    {
        if (!this.session.authorized)
        {
            await this._redirect(LOGIN_URL);
            return false;
        }
        return true;
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
    get backend()
    {
        return this._context.backend;
    }
}

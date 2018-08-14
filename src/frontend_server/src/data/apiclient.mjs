import { request as _request } from 'http';
import assert from 'assert';
import url from 'url';

export class ApiClient
{
    /**
     * @param {string} apiURL 
     */
    constructor(apiURL)
    {
        const components = url.parse(apiURL);
        this._options = {
            hostname: components.hostname,
            port: components.port,
            pathname: components.pathname,
        };
    }

    /**
     * @param {string} method - REST API method name which will be mapped to URL
     * @returns {Promise<Object>} - response parsed as JSON into Object
     */
    sendGet(method)
    {
        return new Promise((resolve, reject) => {
            const options = {
                host: this._options.hostname,
                port: this._options.port,
                path: this._options.pathname + method,
                method: 'GET'
            };
            const request = _request(options, (response) => {
                try
                {
                    if (response.statusCode != 200) {
                        throw new Error(`backend api GET '${method}' returns ${response.statusCode}`);
                    }
                    response.setEncoding('utf8');
                    response.on('data', function (chunk) {
                        try
                        {
                            const value = JSON.parse(chunk);
                            resolve(value);
                        }
                        catch (error)
                        {
                            reject(error);
                        }
                    });
                }
                catch (error)
                {
                    reject(error);
                }
            });
            request.end();
        });
    }

    /**
     * @param {string} method - REST API method name which will be mapped to URL
     * @param {string} payload - payload which will be serialized as JSON
     * @returns {Promise<Object>} - response parsed as JSON into Object
     */
    sendPost(method, payload)
    {
        return new Promise((resolve, reject) => {
            assert(typeof(payload) == 'object');
            const data = JSON.stringify(payload);
            const options = {
                host: this._options.hostname,
                port: this._options.port,
                path: this._options.pathname + method,
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            };
            const request = _request(options, (response) => {
                try
                {
                    if (response.statusCode != 200) {
                        throw new Error(`backend api POST '${method}' returns ${response.statusCode}`);
                    }
                    response.setEncoding('utf8');
                    response.on('data', function (chunk) {
                        try
                        {
                            const value = JSON.parse(chunk);
                            resolve(value);
                        }
                        catch (error)
                        {
                            reject(error);
                        }
                    });
                }
                catch (error)
                {
                    reject(error);
                }
            });
            request.write(data);
            request.end();
        });
    }
}

const http = require('http');
const util = require('util');
const url = require('url');
const assert = require('assert');

class BackendApi
{
    /**
     * @param {config.Config} config 
     */
    constructor(config)
    {
        this._backendHost = 'localhost';
        this._apiPrefix = '/api/v1/';
        this._backendPort = config.backendPort;
    }

    /**
     * Registers new solution build.
     * @param {string} uuid - build UUID
     * @param {string} assignmentUuid - UUID of an assignment which contains test-cases
     * @param {string} language - source code language enumeration string
     * @param {string} source - source code
     * @returns {Promise<string>} - build UUID
     */
    async registerNewBuild(uuid, assignmentUuid, language, source)
    {
        const params = {
            'uuid': uuid,
            'assignment_uuid': assignmentUuid,
            'language': language,
            'source': source
        }
        const response = await this._sendPost('build/new', params);
        assert(uuid == String(response['uuid']));
        console.log('BackendApi: registered build with uuid ' + uuid);

        return uuid;
    }

    /**
     * Register new test case for assignment solutions.
     * @param {sting} uuid - test case UUID
     * @param {sting} assignmentUuid - UUID of an assignment which contains test-case
     * @param {sting} input - given input
     * @param {sting} expected - expected output
     * @returns {Promise<string>} - test case UUID
     */
    async registerTestCase(uuid, assignmentUuid, input, expected)
    {
        const params = {
            'uuid': uuid,
            'assignment_uuid': assignmentUuid,
            'input': input,
            'expected': expected
        }
        const response = await this._sendPost('testcase/new', params);
        assert(uuid == String(response['uuid']));
        console.log('BackendApi: registered test case with uuid ' + uuid);

        return uuid;
    }

    /**
     * Queries current build status.
     * @param {string} uuid - build UUID
     * @returns string - build status enumeration string
     */
    async getBuildStatus(uuid)
    {
        const response = await this._sendGet('build/status/' + uuid);
        return response.get('status');
    }

    /**
     * Queries finished build report, throws if build is not finished.
     * @param {string} uuid - build UUID
     * @returns Object - build report object
     */
    async getBuildReport(uuid)
    {
        const response = await this._sendGet('build/report/' + uuid);
        return response;
    }

    /**
     * @param {string} method - REST API method name which will be mapped to URL
     * @returns {Promise<Object>} - response parsed as JSON into Object
     */
    _sendGet(method)
    {
        return new Promise((resolve, reject) => {
            const options = {
                host: this._backendHost,
                port: this._backendPort,
                path: this._apiPrefix + method,
                method: 'GET'
            };
            const request = http.request(options, (response) => {
                try
                {
                    if (response.statusCode != 200) {
                        throw new Error(`backend api call '${method}' returns ${res.statusCode}`);
                    }
                    const json = response.read();
                    const value = JSON.parse(json);
                    resolve(value);
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
    _sendPost(method, payload)
    {
        return new Promise((resolve, reject) => {
            assert(typeof(payload) == 'object');
            const data = JSON.stringify(payload);
            const options = {
                host: this._backendHost,
                port: this._backendPort,
                path: this._apiPrefix + method,
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            };
            const request = http.request(options, (response) => {
                try
                {
                    if (response.statusCode != 200) {
                        throw new Error(`backend api call '${method}' returns ${res.statusCode}`);
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

module.exports.BackendApi = BackendApi;

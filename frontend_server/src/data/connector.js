const mysql = require('mysql');
const util = require('util');
const assert = require('assert');

class Connector
{
    /**
     * Creates new MySQL database connector using given configuration.
     * @param {config.Config} conf 
     */
    constructor(conf)
    {
        this.connector =  mysql.createConnection({
            host: conf.dbHost,
            user: conf.dbUser,
            password: conf.dbPassword,
            database: conf.dbName
        });
    }

    async connect()
    {
        return util.promisify((cb) => {
            this.connector.connect(cb);
        })();
    }

    /**
     * Executes query and returns Promise to result.
     * @param {string} sql - SQL code to execute
     * @param {Array} values - values to fill SQL code "?" placeholders
     */
    async query(sql, values)
    {
        assert(sql !== undefined);
        assert(values !== undefined);
        return util.promisify((cb) => {
            return this.connector.query(sql, values, cb);
        })();
    }
}

module.exports.Connector = Connector;

const mysql = require('mysql');
const util = require('util');

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
        this.connectImpl = util.promisify(this.connector.connect);
        this.queryImpl = util.promisify(this.connector.query);
    }

    async connect()
    {
        return this.connectImpl();
    }

    /**
     * Executes query and returns Promise to result.
     * @param {string} sql - SQL code to execute
     * @param {Array} values - values to fill SQL code "?" placeholders
     */
    async query(sql, values)
    {
        return this.queryImpl();
    }
}

module.exports.Connector = Connector;

const fs = require('fs');
const util = require('util');

class Config
{
    constructor(jsonObject)
    {
        const dbAuth = jsonObject["db_auth"];

        this.dbName = jsonObject["db_name"];
        this.dbHost = dbAuth["host"];
        this.dbUser = dbAuth["user"];
        this.dbPassword = dbAuth["password"];
    }
}

/**
 * @param {string} path 
 */
async function readConfig(path)
{
    const readFileAsync = util.promisify(fs.readFile);
    const jsonText = await readFileAsync(path);
    const jsonObject = JSON.parse(jsonText);

    return new Config(jsonObject);
}

module.exports.Config = Config;
module.exports.readConfig = readConfig;

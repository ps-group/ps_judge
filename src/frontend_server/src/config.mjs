import { readFile } from 'fs';
import { promisify } from 'util';

export class Config
{
    constructor(jsonObject)
    {
        const dbAuth = jsonObject["db_auth"];
        this.dbName = jsonObject["db_name"];
        this.dbHost = dbAuth["host"];
        this.dbUser = dbAuth["user"];
        this.dbPassword = dbAuth["password"];

        this.backendURL = jsonObject['backend_url'];

        this.port = jsonObject["port"];
        this.backendPort = jsonObject["backend_port"];
    }
}

/**
 * @param {string} path 
 */
export async function readConfig(path)
{
    const readFileAsync = promisify(readFile);
    const jsonText = await readFileAsync(path);
    const jsonObject = JSON.parse(jsonText);

    return new Config(jsonObject);
}

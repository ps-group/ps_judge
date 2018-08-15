import { readFile } from 'fs';
import { promisify } from 'util';

export class Config
{
    constructor(jsonObject)
    {
        this.backendURL = jsonObject['backend_url'];
        this.port = jsonObject["port"];
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

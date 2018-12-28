import { readFileSync } from 'fs';

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
export function readConfig(path)
{
    const jsonText = readFileSync(path);
    const jsonObject = JSON.parse(jsonText);

    return new Config(jsonObject);
}
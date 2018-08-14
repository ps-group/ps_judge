import BackendApi from './data/backendapi.mjs';

export class AppContext
{
    /**
     * Creates application context with given configuration.
     * @param {config.Config} config - application configuration
     */
    constructor(config)
    {
        this.backend = new BackendApi(config);
        this.config = config;
    }
}

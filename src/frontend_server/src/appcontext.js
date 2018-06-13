const connector = require('./data/connector');
const repository = require('./data/repository');

class AppContext
{
    /**
     * Creates application context with given configuration.
     * @param {config.Config} config - application configuration
     */
    constructor(config)
    {
        this.config = config;
        this._connector = new connector.Connector(config);
    }

    /**
     * @returns repository.FrontendRepository
     */
    connectDB()
    {
        return new repository.FrontendRepository(this._connector);
    }
}

module.exports.AppContext = AppContext;

const connect = require('connect');
const compression = require('compression');
const cookieSession = require('cookie-session');
const serveStatic = require('serve-static');
const bodyParser = require('body-parser');
const http = require('http');
const url = require('url');
const fs = require('fs');
const util = require('util');
const routes = require('./routes');
const Router = require('./router');
const appcontext = require('./appcontext');
const config = require('./config');
const messagerouter = require('./data/messagerouter.js');
const buildlistener = require('./listeners/buildlistener.js');

const SESSION_SECRET = '7pv0OvUy';

class Server
{
    constructor()
    {
        this.readFileAsync = util.promisify(fs.readFile);
    }

    async start()
    {
        this.config = await config.readConfig('frontend_server.json');
        this.context = new appcontext.AppContext(this.config);
        this.createServer({
            "port": this.config.port,
            "routes": routes.ROUTES
        });
        this.messageRouter = new messagerouter.MessageRouter();
        this.buildListener = new buildlistener.BuildListener(this.context);
        this.messageRouter.consumeBuildFinished((info) => {
            const uuid = String(info['key']);
            const succeed = Boolean(info['succeed']);
            this.buildListener.onBuildFinished(uuid, succeed);
        });
    }

    createServer(options)
    {
        this.router = new Router(options.routes);

        // gzip/deflate outgoing responses
        this.app = connect();
        this.app.use(compression());

        // store session state in browser cookie
        // TODO: store session data in Redis
        this.app.use(cookieSession({
            secret: 'SESSION_SECRET'
        }));

        // parse urlencoded request bodies into req.body
        this.app.use(bodyParser.urlencoded({extended: false}));

        // respond to all requests with application-specific handlers
        this.app.use((request, response, next) => {
            this.handleCustom(request, response, next)
        });

        // serve static files
        this.app.use(serveStatic('./www/', {
            'dotfiles': 'ignore',
        }));

        console.log('starting on http://localhost:' + options.port + '/');
        http.createServer(this.app).listen(options.port);
    }

    handleCustom(request, response, next)
    {
        const path = url.parse(request.url).pathname;
        const route = this.router.find(path);
        if (route !== null)
        {
            this.handleRoute(route, request, response);
        }
        else
        {
            next();
        }
    }

    async handleRoute(route, request, response)
    {
        try
        {
            const handlerClass = require('./handlers/' + route.handler);
            const hanler = new handlerClass(this.context, request, response);
            await hanler[route.action]();
        }
        catch (error)
        {
            const reason = this._prettyPrintError(error);
            console.error(`internal error when handling '${request.url}': ${reason}`);
            response.writeHead(500);
            response.end();
        }
    }

    _prettyPrintError(error)
    {
        try
        {
            if (typeof(error) == 'object')
            {
                const stack = error.stack || '';
                if (stack !== '')
                {
                    return stack;
                }
                return error.message || '';
            }
            return '' + error
        }
        catch (nextError)
        {
            return `failed to print error (${nextError})`
        }
    }
}

const server = new Server();
server.start();

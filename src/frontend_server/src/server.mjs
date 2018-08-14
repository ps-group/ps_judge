import connect from 'connect';
import compression from 'compression';
import cookieSession from 'cookie-session';
import serveStatic from 'serve-static';
import bodyParser from 'body-parser';
import { createServer as _createServer } from 'http';
import { parse } from 'url';
import { readFile } from 'fs';
import { promisify } from 'util';
import { ROUTES } from './routes.mjs';
import Router from './router.mjs';
import { AppContext } from './appcontext.mjs';
import { readConfig } from './config.mjs';

const SESSION_SECRET = '7pv0OvUy';

export class Server
{
    constructor()
    {
        this.readFileAsync = promisify(readFile);
    }

    async start()
    {
        this.config = await readConfig('frontend_server.json');
        this.context = new AppContext(this.config);
        this.createServer({
            "port": this.config.port,
            "routes": ROUTES
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
            secret: SESSION_SECRET
        }));

        // parse urlencoded request bodies into req.body
        this.app.use(bodyParser.urlencoded({extended: false}));

        // respond to all requests with application-specific handlers
        this.app.use((request, response, next) => {
            this.handleCustom(request, response, next);
        });

        // serve static files
        this.app.use(serveStatic('./www/', {
            'dotfiles': 'ignore',
        }));

        console.log('starting on http://localhost:' + options.port + '/');
        _createServer(this.app).listen(options.port);
    }

    handleCustom(request, response, next)
    {
        const path = parse(request.url).pathname;
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
            const module = await import('./handlers/' + route.handler);
            const hanler = new module.default(this.context, request, response);
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
            return '' + error;
        }
        catch (nextError)
        {
            return `failed to print error (${nextError})`;
        }
    }
}

export function runServer()
{
    const server = new Server();
    server.start();
}


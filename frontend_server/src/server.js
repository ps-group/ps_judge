const http = require('http');
const url = require('url');
const fs = require('fs');
const ejs = require('ejs');
const util = require('util');
const routes = require('./routes');
const Router = require('./router');
const mime = require('mime-types');
const connector = require('./db/connector');
const config = require('./config')

class Server
{
    constructor()
    {
        this.readFileAsync = util.promisify(fs.readFile);
    }

    async start()
    {
        this.config = await config.readConfig('config.json');
        this.connector = new connector.Connector(this.config);
        this.createServer({
            "port": config.port,
            "routes": routes
        });
    }

    createServer(options)
    {
        this.router = new Router(options.routes);
        http.createServer((request, response) => {
            this.handle(request, response)
        }).listen(options.port);
    }

    handle(request, response)
    {
        const path = url.parse(request.url).pathname;
        const route = this.router.find(path);
        if (route !== null)
        {
            this.handleRoute(route, request, response);
        }
        else
        {
            this.handleStatic(path, response);
        }
    }

    async handleRoute(route, request, response)
    {
        try
        {
            const handler = require('./handlers/' + route.handler);
            await handler[route.action](request, response);
        }
        catch (err)
        {
            console.error('internal error: ', err);
            response.writeHead(500);
            response.end();
        }
    }

    async handleStatic(path, response)
    {
        try
        {
            const staticPath = './www/' + path
            const contentType = mime.lookup(staticPath) || 'application/octet-stream';
            response.writeHead(200, {'Content-Type': contentType});
            const data = await this.readFileAsync(staticPath);
            response.write(data);
            response.end();
        }
        catch (err)
        {
            response.writeHead(404, {'Content-Type': 'text/html'});
            response.end("404 Not Found");
        }
    }
}

const server = new Server();
server.start();


class Router
{
    constructor(routes)
    {
        this.routes = routes;
    }

    find(path)
    {
        if (path in this.routes)
        {
            return this.routes[path];
        }
        return null;
    }
}

module.exports = Router;

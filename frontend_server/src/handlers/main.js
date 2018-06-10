const util = require('util');
const ejs = require('ejs');

class Main {
    static async index(request, response) 
    {
        const renderFileAsync = util.promisify(ejs.renderFile);
        const html = await renderFileAsync('./tpl/index.ejs');
        response.writeHead(200, {'Content-Type': 'text/html'});
        response.write(html);
        response.end();
    }
}

module.exports = Main;

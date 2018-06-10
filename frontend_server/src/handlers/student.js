const util = require('util');
const ejs = require('ejs');

class Student {
    constructor()
    {
        this.renderFileAsync = util.promisify(ejs.renderFile);
    }

    static async home(request, response) 
    {
        const html = await this.renderFileAsync('./tpl/student_home.ejs');
        response.writeHead(200, {'Content-Type': 'text/html'});
        response.write(html);
        response.end();
    }

    static async assignment(request, response) 
    {
        const html = await this.renderFileAsync('./tpl/student_assignment.ejs');
        response.writeHead(200, {'Content-Type': 'text/html'});
        response.write(html);
        response.end();
    }

    static async commit(request, response) 
    {
        // TODO: commit changes
    }
}

module.exports = Student;

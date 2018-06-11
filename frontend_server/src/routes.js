
module.exports.LOGIN_URL = '/login';
module.exports.ADMIN_HOME_URL = '/admin';
module.exports.STUDENT_HOME_URL = '/student';
module.exports.JUDGE_HOME_URL = '/judge';

module.exports.ROUTES = {
    "/": {
        "handler": "main",
        "action": "index"
    },
    STUDENT_HOME_URL: {
        "handler": "student",
        "action": "home"
    },
    "/student/assignment": {
        "handler": "student",
        "action": "assignment"
    },
    "/student/commit": {
        "handler": "student",
        "action": "commit"
    }
}

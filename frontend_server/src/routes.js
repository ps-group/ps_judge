const LOGIN_URL = '/login';
const ADMIN_HOME_URL = '/admin';
const STUDENT_HOME_URL = '/student';
const JUDGE_HOME_URL = '/judge';

const ROUTES = {};
ROUTES['/'] = {
    "handler": "main",
    "action": "index"
};
ROUTES[LOGIN_URL] = {
    "handler": "main",
    "action": "login"
};
ROUTES['/login/send'] = {
    "handler": "main",
    "action": "loginSend"
};
ROUTES[STUDENT_HOME_URL] = {
    "handler": "student",
    "action": "home"
};
ROUTES['/student/assignment'] = {
    "handler": "student",
    "action": "assignment"
};
ROUTES['/student/commit'] = {
    "handler": "student",
    "action": "commit"
};

module.exports.ROUTES = ROUTES;
module.exports.LOGIN_URL =  LOGIN_URL;
module.exports.ADMIN_HOME_URL = ADMIN_HOME_URL;
module.exports.STUDENT_HOME_URL = STUDENT_HOME_URL;
module.exports.JUDGE_HOME_URL = JUDGE_HOME_URL;

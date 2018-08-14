export const LOGIN_URL = '/login';
export const ADMIN_HOME_URL = '/admin';
export const STUDENT_HOME_URL = '/student';
export const STUDENT_SOLUTIONS_URL = '/student/solutions';
export const JUDGE_HOME_URL = '/judge';

export const ROUTES = {};
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
ROUTES[STUDENT_SOLUTIONS_URL] = {
    "handler": "student",
    "action": "solutions"
};
ROUTES['/student/assignment'] = {
    "handler": "student",
    "action": "assignment"
};
ROUTES['/student/solution'] = {
    "handler": "student",
    "action": "solution"
};
ROUTES['/student/commit'] = {
    "handler": "student",
    "action": "commit"
};

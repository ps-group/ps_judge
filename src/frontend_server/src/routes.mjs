export const HOME_URL = '/';
export const LOGIN_URL = '/login';
export const LOGIN_SEND_URL = '/login/send';
export const STUDENT_HOME_URL = '/student';
export const STUDENT_SOLUTIONS_URL = '/student/solutions';
export const STUDENT_ASSIGNMENT_URL = '/student/assignment';
export const STUDENT_SOLUTION_URL = '/student/solution';
export const STUDENT_COMMIT_URL = '/student/commit';
export const JUDGE_HOME_URL = '/judge';
export const ADMIN_HOME_URL = '/admin';

export const ROUTES = {};
ROUTES[HOME_URL] = {
    "handler": "main",
    "action": "index"
};
ROUTES[LOGIN_URL] = {
    "handler": "main",
    "action": "login"
};
ROUTES[LOGIN_SEND_URL] = {
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
ROUTES[STUDENT_ASSIGNMENT_URL] = {
    "handler": "student",
    "action": "assignment"
};
ROUTES[STUDENT_SOLUTION_URL] = {
    "handler": "student",
    "action": "solution"
};
ROUTES[STUDENT_COMMIT_URL] = {
    "handler": "student",
    "action": "commit"
};
ROUTES[ADMIN_HOME_URL] = {
    "handler": "admin",
    "action": "home"
};

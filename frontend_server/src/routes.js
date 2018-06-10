
module.exports = {
    "/": {
        "handler": "main",
        "action": "index"
    },
    "/student": {
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

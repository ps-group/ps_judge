const password = require('./data/password')

const pass = process.argv[2]
console.log("pass=", pass)
const hash = password.hashPassword(pass)
console.log("hash=", hash)

import { hashPassword } from './data/password.mjs';

const pass = process.argv[2];
const hash = hashPassword(pass);

console.log("pass=", pass);
console.log("hash=", hash);

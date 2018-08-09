#!/usr/bin/env python3

import json
import uuid
import os
import time

from test_runner import TestScenario, run_test_scenarios

SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
BACKEND_API_URL = 'http://localhost:8081/api/v1/'

PASCAL_SOURCE = """PROGRAM APLUSB;
VAR
  a, b: INTEGER;
BEGIN
  READLN(a);
  READLN(b);
  WRITELN(a+b);
END.
"""

TEST_USERNAME = 'Martin'
TEST_PASSWORD_HASH = '4476d6a3edee189e699ca3c2cfd80905abc8d999954a08d1c504e6ae437cc28dd4194a1051d84bf5cb2cfc19e09e339ce7f2ff83fec56b07ee39ec2205c2adba'

class BackendTestScenario(TestScenario):
    def __init__(self):
        super().__init__(BACKEND_API_URL)

class LoginScenario(BackendTestScenario):
    def run(self):
        self.login_bad_password()
        self.login_bad_username()
        login_info = self.login_ok()
        user_id = login_info['user_id']
        user_info = self.get_user_info(user_id)
        assert user_info == login_info['user']

    def login_ok(self):
        response = self.post_json('user/login', {
            'username': TEST_USERNAME,
            'password_hash': TEST_PASSWORD_HASH
        })
        assert response['succeed'] == True
        assert response['user']['username'] == TEST_USERNAME
        assert int(response['user_id']) != 0
        assert int(response['user']['contest_id']) != 0
        return response

    def get_user_info(self, user_id):
        response = self.get_json('user/{0}/info'.format(str(user_id)))
        assert int(response['contest_id']) != 0
        return response

    def login_bad_password(self):
        response = self.post_json('user/login', {
            'username': TEST_USERNAME,
            "password_hash": TEST_PASSWORD_HASH + "abba",
        })
        assert response['succeed'] == False

    def login_bad_username(self):
        response = self.post_json('user/login', {
            'username': TEST_USERNAME + "abba",
            "password_hash": TEST_PASSWORD_HASH,
        })
        assert response['succeed'] == False

class ViewCommitAndWaitScenario(BackendTestScenario):
    def run(self):
        login_info = self.login_ok()
        contest_id = login_info['user']['contest_id']
        pass

    def login_ok(self):
        response = self.post_json('user/login', {
            'username': TEST_USERNAME,
            'password_hash': TEST_PASSWORD_HASH
        })
        return response

    def view_contest(self, contest_id):
        response = self.get_json('contest/{0}/assignments'.format(str(contest_id)))

if __name__ == "__main__":
    run_test_scenarios([
        LoginScenario,
        ViewCommitAndWaitScenario,
    ])

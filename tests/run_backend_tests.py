#!/usr/bin/env python3

import json
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

class ViewAndCommitScenario(BackendTestScenario):
    def run(self):
        login_info = self.login_ok()
        user_id = login_info['user_id']
        contest_id = login_info['user']['contest_id']
        assignments = self.view_contest(contest_id)
        assignment = self.get_aplusb_assignment(assignments)
        assignment_id = assignment['id']
        self.check_assignment_info(assignment)

        build_uuid = self.commit_aplusb_solution(user_id, assignment_id)
        
        solutions = self.get_user_solutions(user_id)
        solution = self.get_solution_of_assignment(solutions, assignment_id)
        assert solution['assignment_title'] == 'A+B Problem'

    def login_ok(self):
        response = self.post_json('user/login', {
            'username': TEST_USERNAME,
            'password_hash': TEST_PASSWORD_HASH
        })
        return response

    def get_solution_of_assignment(self, solutions, assignment_id):
        for solution in solutions:
            if solution['assignment_id'] == assignment_id:
                return solution
        raise RuntimeError('solution for assignment #{0} not found in list'.format(str(assignment_id)))

    def view_contest(self, contest_id):
        response = self.get_json('contest/{0}/assignments'.format(str(contest_id)))
        assert isinstance(response, list)
        assert len(response) > 0
        return response

    def get_aplusb_assignment(self, assignments):
        for assignment in assignments:
            if assignment['title'] == 'A+B Problem':
                return assignment
        raise RuntimeError('A+B Problem not fund in assignments list')

    def commit_aplusb_solution(self, user_id, assignment_id):
        params = {
            'uuid': self.create_uuid(),
            'assignment_id': assignment_id,
            'language': 'pascal',
            'source': PASCAL_SOURCE
        }
        response = self.post_json('user/{0}/commit'.format(str(user_id)), params)
        assert isinstance(response['uuid'], str)
        return response['uuid']

    def check_assignment_info(self, assignment):
        response = self.get_json('assignment/{0}'.format(str(assignment['id'])))
        assert assignment['id'] == response['id']
        assert assignment['title'] == response['title']
        assert assignment['contest_id'] == response['contest_id']
        assert assignment['uuid'] == response['uuid']
        assert isinstance(response['description'], str)

    def get_user_solutions(self, user_id):
        response = self.get_json('user/{0}/solutions'.format(str(user_id)))
        for solution in response:
            assert isinstance(solution['assignment_id'], int)
            assert isinstance(solution['assignment_title'], str)
            assert isinstance(solution['commit_id'], int)
            assert isinstance(solution['score'], int)
            assert isinstance(solution['build_status'], str)
        return response

if __name__ == "__main__":
    run_test_scenarios([
        LoginScenario,
        ViewAndCommitScenario,
    ])
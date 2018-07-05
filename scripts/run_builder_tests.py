#!/usr/bin/env python3

import json
import uuid
import os
import subprocess
import requests
import time

SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
API_URL_PREFIX = 'http://localhost:8081/api/v1/'
PASCAL_SOURCE = """PROGRAM APLUSB;
VAR
  a, b: INTEGER;
BEGIN
  READLN(a);
  READLN(b);
  WRITELN(a+b);
END.
"""

class RegisterBuildScenario:
    def __init__(self):
        self.assignment_uuid = self.create_uuid()

    def run(self):
        self.register_test_case()
        build_uuid = self.register_new_build()

        while True:
            response = self.get_build_status(build_uuid)
            if response["status"] != "pending":
                break
            time.sleep(1.0)
        print('build {0} finished'.format(build_uuid))
        report = self.get_build_report(build_uuid)
        print('build {0} report:\n{1}'.format(build_uuid, json.dumps(report, indent=2)))

    def register_new_build(self):
        uuid = self.create_uuid()
        request = {
            'uuid': uuid,
            'assignment_uuid': self.assignment_uuid,
            'language': "pascal",
            'source': PASCAL_SOURCE,
        }
        response = self.json_api_post('build/new', request)
        print('registered build ' + uuid)
        assert response.get('uuid') == uuid
        return uuid

    def register_test_case(self):
        uuid = self.create_uuid()
        request = {
            'uuid': uuid,
            'assignment_uuid': self.assignment_uuid,
            'input': '1\n2\n',
            'expected': '3\n',
        }
        response = self.json_api_post('testcase/new', request)
        print('registered test case ' + uuid)
        assert response.get('uuid') == uuid
        return uuid

    def get_build_status(self, uuid):
        response = self.json_api_get('build/status/' + uuid)
        print('got info for build ' + uuid)
        assert response.get('uuid') == uuid
        assert response.get('status') != ''
        return response

    def get_build_report(self, uuid):
        response = self.json_api_get('build/report/' + uuid)
        print('got report for build ' + uuid)
        assert response.get('uuid') == uuid
        return response

    def json_api_post(self, method, request_dict):
        url = API_URL_PREFIX + method
        data = json.dumps(request_dict, indent=2)
        headers = {
            'Content-Type': 'application/json'
        }
        response = requests.post(url, data, headers=headers)
        return response.json()

    def json_api_get(self, query):
        url = API_URL_PREFIX + query
        response = requests.get(url)
        if response.status_code != 200:
            reason = "status code {0} on URL '{1}'\n{2}".format(response.status_code, url, response.text)
            raise RuntimeError(reason)
        return response.json()

    def exec_process(self, cmd):
        process = subprocess.Popen(cmd, cwd=SCRIPT_DIR, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout, stderr = process.communicate()
        if process.returncode != 0:
            raise subprocess.CalledProcessError(process.returncode, cmd, stderr)
        return stdout

    def create_uuid(self):
        """
        Returns string like '9fe2c4e93f654fdbb24c02b15259716c'
        """
        return uuid.uuid4().hex

def main():
    scenario_classes = [
        RegisterBuildScenario
    ]
    for scenario_class in scenario_classes:
        scenario = scenario_class()
        scenario.run()

if __name__ == "__main__":
    main()

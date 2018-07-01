#!/usr/bin/env python3

import json
import uuid
import os
import subprocess
import requests

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
        build_uuid = self.register_new_build()
        self.call_build_info(build_uuid)

    def register_new_build(self):
        build_uuid = self.create_uuid()
        request = {
            'uuid': build_uuid,
            'assignment_uuid': self.assignment_uuid,
            'language': "pascal",
            'source': PASCAL_SOURCE
        }
        response = self.json_api_post('build/new', request)
        assert response.get('uuid') == build_uuid
        return build_uuid

    def call_build_info(self, build_uuid):
        response = self.json_api_get('build/' + build_uuid)
        assert response.get('uuid') == build_uuid
        assert response.get('status') == 'pending'
        assert response.get('score') == 0
        assert response.get('details') == ""
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

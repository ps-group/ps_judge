#!/usr/bin/env python3

from __future__ import print_function
import json
import uuid
import os
import time

from test_runner import TestScenario, run_test_scenarios

SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
BUILDER_API_URL = 'http://localhost:9092/api/v1/'
PASCAL_SOURCE = """PROGRAM APLUSB;
VAR
  a, b: INTEGER;
BEGIN
  READLN(a);
  READLN(b);
  WRITELN(a+b);
END.
"""

class BuilderTestScenario(TestScenario):
    def __init__(self):
        super().__init__(BUILDER_API_URL)

class RegisterBuildScenario(BuilderTestScenario):
    def __init__(self):
        super().__init__()
        self.assignment_uuid = self.create_uuid()

    def run(self):
        self.register_test_case()
        build_uuid = self.register_new_build()
        for _ in range(0, 20):
            response = self.get_build_status(build_uuid)
            if response["status"] != "pending":
                break
            time.sleep(1.0)
        else:
            raise RuntimeError('build timeout exceed')
        print('build {0} finished'.format(build_uuid))
        report = self.get_build_report(build_uuid)
        print('build {0} report:\n{1}'.format(build_uuid, json.dumps(report, indent=2)))

    def register_new_build(self):
        uuid = self.create_uuid()
        response = self.post_json('build/new', {
            'uuid': uuid,
            'assignment_uuid': self.assignment_uuid,
            'language': "pascal",
            'source': PASCAL_SOURCE,
        })
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
        response = self.post_json('testcase/new', request)
        print('registered test case ' + uuid)
        assert response.get('uuid') == uuid
        return uuid

    def get_build_status(self, uuid):
        response = self.get_json('build/status/' + uuid)
        print('got info for build ' + uuid)
        assert response.get('uuid') == uuid
        assert response.get('status') != ''
        return response

    def get_build_report(self, uuid):
        response = self.get_json('build/report/' + uuid)
        print('got report for build ' + uuid)
        assert response.get('uuid') == uuid
        return response

def main():
    run_test_scenarios([
        RegisterBuildScenario
    ])

if __name__ == "__main__":
    main()

from __future__ import print_function
import json
import uuid
import os
import subprocess

SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
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
        info = self.call_build_info(build_uuid)
        print('build status=', info('status'))

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
        return response

    def json_api_post(self, method, request_dict):
        url = 'localhost:8081/api/v1/' + method
        data_path = os.path.join(SCRIPT_DIR, '__data__.json')
        json_str = json.dumps(request_dict, indent=2)
        with open(data_path, 'w') as data_file:
            data_file.write(json_str)
        cmd = ['curl', '-X', 'POST', '-H', 'Content-Type: application/json', '--data-binary', '@' + data_path, url]
        stdout = self.exec_process(cmd)
        return json.loads(stdout)

    def json_api_get(self, query):
        url = 'localhost:8081/api/v1/' + query
        cmd = ['curl', url]
        stdout = self.exec_process(cmd)
        return json.loads(stdout)

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

import json
import requests
import sys
import traceback
import uuid

class TestScenario:
    def __init__(self, api_url):
        self.api_url = str(api_url)

    def post_json(self, method, request_dict):
        print("  call {0}".format(method))
        url = self.api_url+ method
        data = json.dumps(request_dict, indent=2)
        headers = {
            'Content-Type': 'application/json'
        }
        response = requests.post(url, data, headers=headers)
        if len(response.text) == 0:
            return None
        response_json = response.json()
        self.throw_if_error_response(response_json)
        return response_json

    def get_json(self, query):
        url = self.api_url + query
        response = requests.get(url)
        if response.status_code != 200:
            reason = "status code {0} on URL '{1}'\n{2}".format(response.status_code, url, response.text)
            raise RuntimeError(reason)
        response_json = response.json()
        self.throw_if_error_response(response_json)
        return response_json

    def throw_if_error_response(self, response):
        if isinstance(response, dict) and response.get('error') is not None:
            reason = response.get('error', dict()).get('text', '')
            raise RuntimeError('method failed: ' + reason)

    def dump_response(self, response):
        print('response:\n{0}'.format(json.dumps(response, indent=2)))

    def create_uuid(self):
        """
        Returns string like '9fe2c4e93f654fdbb24c02b15259716c'
        """
        return uuid.uuid4().hex

def run_test_scenarios(scenario_classes):
    total_count = 0
    succeed_count = 0
    for scenario_class in scenario_classes:
        print("Running {0}...".format(scenario_class.__name__))
        scenario = scenario_class()
        try:
            total_count += 1
            scenario.run()
            succeed_count += 1
        except:
            print("*** ERROR IN SCENARIO '{0}' ***".format(scenario_class.__name__), file=sys.stderr)
            traceback.print_exc(limit=4, file=sys.stderr, chain=False)
    report = "Finished {0} tests, {1} succeed, {2} failed".format(total_count, succeed_count, total_count - succeed_count)
    print("".ljust(len(report), "."))
    print(report)

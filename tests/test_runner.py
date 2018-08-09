import json
import requests
import sys
import traceback

class TestScenario:
    def __init__(self, api_url):
        self.api_url = str(api_url)

    def post_json(self, method, request_dict):
        url = self.api_url+ method
        data = json.dumps(request_dict, indent=2)
        headers = {
            'Content-Type': 'application/json'
        }
        response = requests.post(url, data, headers=headers)
        return response.json()

    def get_json(self, query):
        url = self.api_url + query
        response = requests.get(url)
        if response.status_code != 200:
            reason = "status code {0} on URL '{1}'\n{2}".format(response.status_code, url, response.text)
            raise RuntimeError(reason)
        return response.json()

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
            type, value, _ = sys.exc_info()
            print("*** ERROR IN SCENARIO '{0}' ***".format(scenario_class.__name__), file=sys.stderr)
            print("{0}: {1}".format(type.__name__, value), file=sys.stderr)
    report = "Finished {0} tests, {1} succeed, {2} failed".format(total_count, succeed_count, total_count - succeed_count)
    print("".ljust(len(report), "."))
    print(report)

#!/usr/bin/env python3

import os
from subprocess import Popen, check_call
import test_config_master

SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
REPO_DIR = os.path.dirname(SCRIPT_DIR)

check_call("bash " + SCRIPT_DIR + "/update_test_db", shell=True)

test_config_master.create_test_config()

backend_service_process = Popen("exec " + REPO_DIR + '/bin/backend_service', shell=True)
builder_service_process = Popen("exec " + REPO_DIR + '/bin/builder_service', shell=True)

check_call(SCRIPT_DIR + "/run_backend_tests.py", shell=True)
check_call(SCRIPT_DIR + "/run_builder_tests.py", shell=True)

backend_service_process.kill()
builder_service_process.kill()

test_config_master.remove_dev_config()
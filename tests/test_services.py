#!/usr/bin/env python3

import os
from subprocess import Popen
import test_config_master


SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
REPO_DIR = os.path.dirname(SCRIPT_DIR)

os.system("bash " + SCRIPT_DIR + "/update_test_db")

test_config_master.create_test_config()

backend_service_process = Popen("exec " + REPO_DIR + '/bin/backend_service', shell=True)
builder_service_process = Popen("exec " + REPO_DIR + '/bin/builder_service', shell=True)

os.system(SCRIPT_DIR + "/run_backend_tests.py")
os.system(SCRIPT_DIR + "/run_builder_tests.py")

backend_service_process.kill()
builder_service_process.kill()

test_config_master.remove_dev_config()
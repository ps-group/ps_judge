#!/usr/bin/env python3

import os
import json

backend_port = 8081
builder_port = 8082
rabbitmq_port = 5672

test_backend_port = 9091
test_builder_port = 9092
test_rabbitmq_port = 6672

db_username = "psjudge"
db_password = "1234"

test_frontend_db = "psjudge_frontend_test"
test_builder_db = "psjudge_builder_test"

SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
REPO_DIR = os.path.dirname(SCRIPT_DIR)

BIN_DIR = os.path.join(REPO_DIR, 'bin')
if not os.path.exists(BIN_DIR):
    os.makedirs(BIN_DIR)

BACKEND_CONFIG_PATH = os.path.join(BIN_DIR, 'backend_service.json')
BUILDER_CONFIG_PATH = os.path.join(BIN_DIR, 'builder_service.json')

with open(REPO_DIR + '/bin/backend_service.json') as f:
    dev_config = json.load(f)

class DevConfig:
    backend_port = dev_config["backend_url"].split(':')[1]
    builder_port = dev_config["builder_url"].split(':')[1]
    rabbitmq_port = dev_config["amqp_socket"].split(':')[3][0:-1]

def write_config(filepath, values):
    content = json.dumps(values, indent=4)
    with open(filepath, 'w') as f:
        f.write(content)
    print('written {0}'.format(filepath))

def create_test_config():
    backend_config = {
        "mysql_db": test_frontend_db,
        "mysql_host": "127.0.0.1",
        "mysql_user": db_username,
        "mysql_password": db_password,
        "backend_url": "127.0.0.1:{0}".format(str(test_backend_port)),
        "builder_url": "127.0.0.1:{0}".format(str(test_builder_port)),
        "log_file_name": "backend_service.log",
        "amqp_socket": "amqp://guest:guest@localhost:{0}/".format(str(test_rabbitmq_port))
    }
    write_config(BACKEND_CONFIG_PATH, backend_config)

    builder_config = {
        "mysql_db": test_builder_db ,
        "mysql_host": "127.0.0.1",
        "mysql_user": db_username,
        "mysql_password": db_password,
        "builder_url": "127.0.0.1:{0}".format(str(test_builder_port)),
        "log_file_name": "builder_service.log",
        "amqp_socket": "amqp://guest:guest@localhost:{0}/".format(str(test_rabbitmq_port))
    }
    write_config(BUILDER_CONFIG_PATH, builder_config)

def remove_dev_config():
    backend_config = {
        "mysql_db": "psjudge_frontend",
        "mysql_host": "127.0.0.1",
        "mysql_user": db_username,
        "mysql_password": db_password,
        "backend_url": "127.0.0.1:{0}".format(str(DevConfig.backend_port)),
        "builder_url": "127.0.0.1:{0}".format(str(DevConfig.builder_port)),
        "log_file_name": "backend_service.log",
        "amqp_socket": "amqp://guest:guest@localhost:{0}/".format(str(DevConfig.rabbitmq_port))
    }
    write_config(BACKEND_CONFIG_PATH, backend_config)

    builder_config = {
        "mysql_db": "psjudge_builder",
        "mysql_host": "127.0.0.1",
        "mysql_user": db_username,
        "mysql_password": db_password,
        "builder_url": "127.0.0.1:{0}".format(str(DevConfig.builder_port)),
        "log_file_name": "builder_service.log",
        "amqp_socket": "amqp://guest:guest@localhost:{0}/".format(str(DevConfig.rabbitmq_port))
    }
    write_config(BUILDER_CONFIG_PATH, builder_config)

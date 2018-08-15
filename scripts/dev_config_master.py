#!/usr/bin/env python3

import os
import json

SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
REPO_DIR = os.path.dirname(SCRIPT_DIR)

BIN_DIR = os.path.join(REPO_DIR, 'bin')
if not os.path.exists(BIN_DIR):
    os.makedirs(BIN_DIR)

FRONTEND_CONFIG_PATH = os.path.join(REPO_DIR, 'src', 'frontend_server', 'frontend_server.json')
BACKEND_CONFIG_PATH = os.path.join(BIN_DIR, 'backend_service.json')
BUILDER_CONFIG_PATH = os.path.join(BIN_DIR, 'builder_service.json')


def input_port(hint, default_value):
    value = input(hint)
    if len(value) == 0:
        return default_value
    value = int(value)
    if value <= 0 or value >= 2**16:
        raise RuntimeError('invalid URL port value: ' + value)
    return value


def write_config(filepath, values):
    content = json.dumps(values, indent=4)
    with open(filepath, 'w') as f:
        f.write(content)
    print('written {0}'.format(filepath))


db_username = input('database user name: ')
db_password = input('database password: ')
frontend_port = input_port('frontend service port (default: 8080): ', 8080)
backend_port = input_port('backend service port (default: 8081): ', 8081)
builder_port = input_port('builder service port (default: 8082): ', 8082)
rabbitmq_port = input_port('RabbitMQ port (default: 5672): ', 5672)

frontend_config = {
    "port": frontend_port,
    "backend_url": "http://localhost:{0}/api/v1/".format(str(backend_port))
}
write_config(FRONTEND_CONFIG_PATH, frontend_config)

backend_config = {
    "mysql_db": "psjudge_frontend",
    "mysql_host": "127.0.0.1",
    "mysql_user": db_username,
    "mysql_password": db_password,
    "backend_url": "127.0.0.1:{0}".format(str(backend_port)),
    "builder_url": "127.0.0.1:{0}".format(str(builder_port)),
    "log_file_name": "backend_service.log",
    "amqp_socket": "amqp://guest:guest@localhost:{0}/".format(str(rabbitmq_port))
}
write_config(BACKEND_CONFIG_PATH, backend_config)

builder_config = {
    "mysql_db": "psjudge_builder",
    "mysql_host": "127.0.0.1",
    "mysql_user": db_username,
    "mysql_password": db_password,
    "builder_url": "127.0.0.1:{0}".format(str(builder_port)),
    "log_file_name": "builder_service.log",
    "amqp_socket": "amqp://guest:guest@localhost:{0}/".format(str(rabbitmq_port))
}
write_config(BUILDER_CONFIG_PATH, builder_config)

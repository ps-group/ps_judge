# PSJudge - simple online judge for PS-Group

Work in progress implementation.

## Build and Run

See [Project Setup](docs/setup.md).

## Architecture

There are following components:

* `frontend_server` serves HTML/JS/CSS for users
* `builder_service` builds solutions and runs assignment input/output tests

There are following interaction routes:

* `frontend_server` monitors "build finished" events from RabbitMQ
* `builder_service` provides JSON REST API for `frontend_server` and posts "build finished" events to RabbitMQ

Following diagram shows how student solutions served by judge:

![Diagram](docs/psjudge-events.png)

# Project Setup

>You need Linux (Ubuntu recommended for development), MySQL, Node.js, Go and RabbitMQ to build and run this project.

## Install Dependencies and Build

* Run Bash script `scripts\install_deps` to install third-party dependencies
* Run Bash script `scripts\build` to build all services

## Setup MySQL

MySQL setup can be done with the following MySQL commands (don't forget to replace `<password>` with password):

```sql
CREATE USER 'psjudge'@'localhost' IDENTIFIED BY '<password>';
CREATE DATABASE `psjudge_builder`;
GRANT ALL PRIVILEGES ON psjudge_builder.* TO 'psjudge'@'localhost' WITH GRANT OPTION;
CREATE DATABASE `psjudge_frontend`;
GRANT ALL PRIVILEGES ON psjudge_frontend.* TO 'psjudge'@'localhost' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

After that, you can type `exit` to exit MySQL shell.

## Setup Database

Run the following Bash scripts to setup database:

```
scripts/update_builder_model
scripts/update_frontend_model
```

* The first script creates/updates database scheme
* The second script creates/updates database scheme and installs sandbox data listed in `frontend_sandbox_data.sql`
  * The script creates user `Martin` with password `2018`
  
## Create config file

You can run all services on localhost, but still should assign different ports. We recommend to use following URLs:

* `localhost:8080` for frontend service
* `localhost:8081` for backend service
* `localhost:8082` for builder service

1. Go to the `src/frontend_server/` and create `frontend_service.json` configuration file. Look at `src/frontend_server/src/config.mjs` and create config with matching fields.
2. Go to the `bin` directory and create `backend_service.json`
  * See `src/backend_service/config.go` and create config with matching fields.
3. Go to the `bin` directory and create `builder_service.json`
  * See `src/builder_service/config.go` and create config with matching fields.

## Run Tests

Now, you can run integration tests to check that everything is OK.

Backend and Builder tested with the following Python scripts:

```bash
tests/run_backend_tests.py
tests/run_builder_tests.py
```

Frontend has no automatic tests and can be tested manually in browser.

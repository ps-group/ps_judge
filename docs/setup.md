# Project setup

>This document is in work in progress state

## MySQL setup

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

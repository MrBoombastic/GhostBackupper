# GhostBackupper

![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/mrboombasticultra/ghostbackupper)
![Docker Pulls](https://img.shields.io/docker/pulls/mrboombasticultra/ghostbackupper)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/mrboombasticultra/ghostbackupper)
[![Go Report Card](https://goreportcard.com/badge/github.com/MrBoombastic/GhostBackupper)](https://goreportcard.com/report/github.com/MrBoombastic/GhostBackupper)
[![CodeFactor](https://www.codefactor.io/repository/github/mrboombastic/ghostbackupper/badge)](https://www.codefactor.io/repository/github/mrboombastic/ghostbackupper)

Back up your Ghost CMS instance easily!

![Preview 1](/docs/gb1.png)

This tool can work in two modes:

- CLI
- Docker

## CLI

```bash
"./ghostbackupper backup --db_host $DB_HOST --db_password $DB_PASSWORD --db_user $DB_USER --db_database $DB_DATABASE --db_port $DB_PORT --content $CONTENT --output $OUTPUT --mega_login $MEGA_LOGIN --mega_password $MEGA_PASSWORD"
```
## Docker `run`
Use similarly as above. [Official Docker reference here.](https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e---env---env-file)
## Docker Compose

Recommended Docker Compose configuration - everything in Docker

```yaml
version: '3.1'

services:
  ghost:
    image: ghost:5-alpine
    restart: always
    ports:
      - 2368:2368
    environment:
      database__client: mysql
      database__connection__host: db
      database__connection__user: ghost-123
      database__connection__password: superpassword
      database__connection__database: ghost_prod
      url: https://example.com
    volumes:
      - /var/www/ghost/content:/var/lib/ghost/content

  db:
    image: mysql:8.0
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: evenmoresuperpassword
      MYSQL_DATABASE: ghost_prod
      MYSQL_PASSWORD: superpassword
      MYSQL_USER: ghost-123
    volumes:
      - /home/ubuntu/docker/db:/var/lib/mysql

  backupper:
    image: mrboombasticultra/ghostbackupper:latest
    restart: no
    environment:
      DB_HOST: db
      DB_PASSWORD: superpassword
      DB_DATABASE: ghost_prod
      DB_USER: ghost-123
      DB_PORT: 3306
      CONTENT: /var/www/ghost/content
      OUTPUT: backup.tar.gz
      MEGA_LOGIN: user@mail.com
      MEGA_PASSWORD: megapassword
    volumes:
      - /var/www/ghost/content:/var/www/ghost/content
```
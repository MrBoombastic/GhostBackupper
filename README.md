# GhostBackupper

This tool can work in two modes:

- CLI
- Docker

## CLI

Recommended Docker Compose configuration

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
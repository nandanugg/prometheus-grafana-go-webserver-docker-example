
### Prometheus Grafana integration with Golang Webserver as a Backend, compiled with docker compose example

#### Prerequisites
- Docker
- Golang
- Postgres database

#### Env variables
```
export DB_NAME=
export DB_PORT=
export DB_HOST=
export DB_USERNAME=
export DB_PASSWORD=
export DB_PARAMS="sslmode=disable"
export ENV=development
```
> ⚠️ Do not use `.env` file, set your own environment variable in your machine to mimic the production env, or just use `dotenv` or `autoenv`
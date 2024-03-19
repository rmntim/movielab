# movielab

RESTful API in Golang with JWT authentication and PostgreSQL as storage.

## Usage

First of all you have to put up postgres db.
You can either do it yourself or use [docker-compose.yml](./docker-compose.yml) file provided
by typing in your terminal:

```sh
docker compose up db -d
```

This will create a container with port 5432, username and password `postgres` and database name `movielab`.
Export `DATABASE_URL` environment variable somewhere, app will need it to connect to postgres instance.

Then you can run the app by running
```sh
task # or task run
```
or build it with
```sh
task build
./build/server
```

Note that you have to specify config path either in `--config` flag or with `CONFIG_PATH` environment variable.

## Docker

You can run the app with single command by typing:

```sh
docker compose up -d
```

You still need to specify `DATABASE_URL` environment variable.

## Testing

All generated mocks are already available, but for consistency you can run

```sh
task generate
```
before running the tests with
```sh
task test
```

## Docs

There is OpenAPI specification available in `/docs` route and [openapi.yaml](./api/openapi.yaml) file.
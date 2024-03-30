# `make run` only run only the backend server
.PHONY: run

# `make runDockerCompose` run the backend server, prometheus and grafana with docker-compose
.PHONY: runDockerCompose

# `make runDockerComposeBuildBackendService` rebuild and run the backend server, prometheus and grafana with docker-compose
.PHONY: runDockerComposeBuildBackendService

# `make resetGrafana` reset the grafana dashboard
.PHONY: resetGrafana

# `make runDb` run the postgres database via docker
.PHONY: runDb

# `make migrateNew` create a new migration
.PHONY: migrateNew

# `make migrateUp` apply the migrations
.PHONY: migrateUp

# `make migrateDown` rollback the migrations
.PHONY: migrateDown

# BEFORE STARTING
# Make sure you already run `make runDb` and `make migrateUp` before running the application

run:
	@echo "Running the application..."
	go run cmd/main.go

runDockerCompose:
	@echo "Running the application with docker-compose..."
	DB_HOST=host.docker.internal ENV=production docker-compose up

runDockerComposeBuildBackendService:
	@echo "Rebuilding and running the application with docker-compose..."
	DB_HOST=host.docker.internal ENV=production docker-compose up --build backend 

resetGrafana:
	@echo "Resetting Grafana..."
	docker-compose rm -f grafana
	docker volume rm -f prometheusgrafanaexample_grafana_data

runDb:
	@echo "Creating database..."
	docker run --rm -e POSTGRES_DB=postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16-alpine

migrateNew:
	@echo "Creating new migration..."
	@read -p "Enter the name of the migration: " name; \
	migrate create -ext sql -dir db/migrations $$name

migrateUp:
	@echo "Applying migrations..."
	migrate -path db/migrations -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" up

migrateDown:
	@echo "Rolling back migrations..."
	migrate -path db/migrations -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" down 1

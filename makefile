.PHONY: run
.PHONY: runDockerCompose
.PHONY: runDockerComposeBuildBackendService
.PHONY: resetGrafana
.PHONY: exportGrafana
.PHONY: importGrafana 
.PHONY: runDb
.PHONY: migrateNew
.PHONY: migrateUp
.PHONY: migrateDown

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

exportGrafana:
	@echo "Exporting Grafana dashboard..."
	mkdir -p data_sources && curl -s "http://localhost:3000/api/datasources"  -u admin:admin|jq -c -M '.[]'|split -l 1 - data_sources/

importGrafana:
	@echo "Importing Grafana dashboard..."
	for i in data_sources/*; do \
		curl -X "POST" "http://localhost:3000/api/datasources" \
		-H "Content-Type: application/json" \
		--user admin:admin \
		--data-binary @$i
	done
	
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

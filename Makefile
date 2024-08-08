APP_NAME=go_worker

# DOCKER TASKS
# Build the container
build: ## Build the container
	docker build -t $(APP_NAME) .

dcb:
	docker compose build

dcu:
	docker compose up

rnode:
	docker compose run --rm phoenix cd assets && yarn install 

dci:
	docker exec -ti $(APP_NAME) /bin/bash

dbprd:
	docker compose -f 'docker-compose-prod.yml' up --build

dbdev:
	docker compose -f 'docker-compose.yml' up --build
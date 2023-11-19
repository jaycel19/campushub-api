DSN="host=localhost port=5432 user=root password=root dbname=campushubdb sslmode=disable timezone=UTC connect_timeout=5"
PORT=8080
SECRET=adsfncnvkq08ew7r098djfqew89r098d7f987qerlakdj
DB_DOCKER_CONTAINER=campushub_db
BINARY_NAME=campushubapi
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=AKIAZHO2L2OKOV2MW2MM
AWS_SECRET_ACCESS_KEY=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


postgres:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine
createdb:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=root --owner=root campushubdb

stop_containers:
	@echo "Stopping other docker containers"
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers..."; \
		 docker stop $$(docker ps -q); \
	else \
		echo "no active containers found..."; \
	fi

start-docker:
	docker start ${DB_DOCKER_CONTAINER}

create_migrations:
	sqlx migrate add -r $(name)

migrate_up:
	sqlx migrate run --database-url "postgres://root:root@localhost:5432/campushubdb?sslmode=disable"

migrate_down:
	sqlx migrate revert --database-url "postgres://root:root@localhost:5432/campushubdb?sslmode=disable"

build:
	@echo "Building backend api binary"
	go build -o ${BINARY_NAME} cmd/server/*.go
	@echo "Binary built!"

run: build
	@echo "Start api"
	@env PORT=${PORT} DSN=${DSN} ./${BINARY_NAME} &
	@echo "api started!"

stop:
	@echo "Stopping backend"
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped backend"

start: run

restart: stop start

postgres_restart: stop_containers start-docker

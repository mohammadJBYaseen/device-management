export PGUSER ?= postgres
export PGPASSWORD ?= postgres
export PGHOST ?= localhost
export DB_NAME ?= device_management
export SSL_MODE ?= disable

run-ci:
	make dev-db-docker
	make db-ci
	make dockerized-app
dev-db-docker:
	docker run -d --name pg --env=POSTGRES_PASSWORD=$(PGPASSWORD) --health-cmd "pg_isready -U $(PGUSER)" --health-interval 10s --health-timeout 5s --health-retries 5 \
		-p 5432:5432 -v ./db/setup/ci_setup.sql:/docker-entrypoint-initdb.d/init.sql -d postgres:17
	@printf "Waiting for Postgres container.";
	@while [ $$(docker inspect --format='{{.State.Health.Status}}' pg) != "healthy" ]; do \
		sleep 2; \
		printf '.'; \
	done
dev-db-docker-delete:
	docker container stop pg
	docker container rm pg
db-create:
	migrate -database postgres://$(PGUSER):$(PGPASSWORD)@$(PGHOST):5432/$(DB_NAME)
db-ci:
	cat ./db/setup/ci_setup.sql| docker exec -i pg psql -U dbowner -d postgres
	cat ./db/setup/db_setup.sql| docker exec -i pg psql -U dbowner -d postgres
	make db-migrate
	cat ./db/setup/dev_setup.sql| docker exec -i pg psql -U dbowner -d postgres

db-clean:
	make db-
	cat ./db/setup/db_remove.sql| docker exec -i pg psql -U dbowner -d postgres
	make db-ci
db-migrate:
	#Docker CLI variant
	docker run --rm -v ./db/migrations:/migrations --network host migrate/migrate -path=/migrations -database postgres://$(PGUSER):$(PGPASSWORD)@$(PGHOST):5432/$(DB_NAME)?sslmode=$(SSL_MODE) up
db+:
	migrate -source file://db/migrations -database postgres://$(PGUSER):$(PGPASSWORD)@$(PGHOST):5432/$(DB_NAME)?sslmode=$(SSL_MODE) up
db-:
	migrate -source file://db/migrations -database postgres://$(PGUSER):$(PGPASSWORD)@$(PGHOST):5432/$(DB_NAME)?sslmode=$(SSL_MODE) down -all
dockerized-app:
	docker build -t device-management . && docker run --name device-management -p 8080:8080 --env=DATABASE_HOST=host.docker.internal --env=SERVER_HOST=0.0.0.0 --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --restart=no --label='org.opencontainers.image.ref.name=ubuntu' --label='org.opencontainers.image.version=24.04' --runtime=runc -d  device-management

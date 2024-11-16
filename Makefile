export PGUSER ?= postgres
export PGPASSWORD ?= postgres
export PGHOST ?= localhost
export DB_NAME ?= device_management
export SSL_MODE ?= disable

run-ci:
	make dev-db-docker
	make db-ci

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
	make db+
	cat ./db/setup/dev_setup.sql| docker exec -i pg psql -U dbowner -d postgres

db-clean:
	make db-
	cat ./db/setup/db_remove.sql| docker exec -i pg psql -U dbowner -d postgres
	make db-ci
db+:
	migrate -source file://db/migrations -database postgres://$(PGUSER):$(PGPASSWORD)@$(PGHOST):5432/$(DB_NAME)?sslmode=$(SSL_MODE) up
db-:
	migrate -source file://db/migrations -database postgres://$(PGUSER):$(PGPASSWORD)@$(PGHOST):5432/$(DB_NAME)?sslmode=$(SSL_MODE) down -all

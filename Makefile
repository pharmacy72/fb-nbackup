DOCKER_COMPOSE_SRC_LOCAL=build/docker.local/docker-compose.yml
DOCKER_COMPOSE_CMD_LOCAL=docker-compose -f ${DOCKER_COMPOSE_SRC_LOCAL}
DOCKER_COMPOSE_SRC_CI=build/docker.ci/docker-compose.yml
DOCKER_COMPOSE_CMD_CI=docker-compose -f ${DOCKER_COMPOSE_SRC_CI}

ISC_PASSWORD=023RsdTf4UI123
ISC_USER=fbuser
ISC_DATABASE=NBEXAMPLE

run-local: stop-local
	${DOCKER_COMPOSE_CMD_LOCAL} up -d
	${DOCKER_COMPOSE_CMD_LOCAL} exec -T fb /usr/local/firebird/bin/isql -e -user ${ISC_USER} -password ${ISC_PASSWORD} -database ${ISC_DATABASE} < ./scripts/create.db.sql

stop-local:
	${DOCKER_COMPOSE_CMD_LOCAL} down -v --remove-orphans

run-for-test: stop-for-test
	${DOCKER_COMPOSE_CMD_CI} up -d
	${DOCKER_COMPOSE_CMD_CI} exec -T fb /usr/local/firebird/bin/isql -e -user ${ISC_USER} -password ${ISC_PASSWORD} -database ${ISC_DATABASE} < ./scripts/create.db.sql

stop-for-test:
	${DOCKER_COMPOSE_CMD_CI} down -v --remove-orphans

test:
	go test -v ./...

lint:
	golangci-lint run ./...

FIREBIRD_DATABASE=
func-test: stop-for-test run-for-test
	go test -v --tags=functest -run ^TestFunctional\$
	$(MAKE) stop-for-test
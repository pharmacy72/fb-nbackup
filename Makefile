DOCKER_COMPOSE_SRC_LOCAL=build/docker.local/docker-compose.yml
DOCKER_COMPOSE_CMD_LOCAL=docker-compose -f ${DOCKER_COMPOSE_SRC_LOCAL}
DOCKER_COMPOSE_SRC_CI=build/docker.ci/docker-compose.yml
DOCKER_COMPOSE_CMD_CI=docker-compose -f ${DOCKER_COMPOSE_SRC_CI}
run-local: stop-local
	${DOCKER_COMPOSE_CMD_LOCAL} up -d
stop-local:
	${DOCKER_COMPOSE_CMD_LOCAL} down -v --remove-orphans

run-for-test: stop-for-test
	${DOCKER_COMPOSE_CMD_CI} up -d
stop-for-test:
	${DOCKER_COMPOSE_CMD_CI} down -v --remove-orphans

test:
	go test -v ./...

lint:
	golangci-lint run ./...

func-test: stop-for-test run-for-test
	go test -v --tags=functest -run ^TestFunctional\$
	$(MAKE) stop-for-test
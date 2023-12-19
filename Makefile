DOCKER_IMAGE_TEST_NAME := metrics-agent-tests
DOCKER_CONTAINER_NAME := metrics-agent-tests-container



up.debug:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d

stop:
	docker-compose stop

logs.agent:
	docker-compose logs -f agent

logs.server:
	docker-compose logs -f server

up.prod:
	docker-compose up --build -d

up.prod.server:
	docker-compose up server --build

restart:
	docker-compose stop
	make up.prod

tests.remote:
	docker build -t $(DOCKER_IMAGE_TEST_NAME) -f Docker/Dockerfile-tests --target=tests .
	docker run --rm -it --name ${DOCKER_CONTAINER_NAME} $(DOCKER_IMAGE_TEST_NAME) bash


run.server:
	docker-compose run --rm --build server sh

run.agent:
	docker-compose run --rm --build server sh

tests.server:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml run --rm server go test ./...

tests.agent:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml run --rm agent go test ./...



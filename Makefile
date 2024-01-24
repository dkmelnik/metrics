DOCKER_IMAGE_TEST_NAME := metrics-agent-tests
DOCKER_CONTAINER_NAME := metrics-agent-tests-container


up.prod:
	docker-compose up --build -d

up.debug:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

logs.agent:
	docker-compose logs -f agent

logs.server:
	docker-compose logs -f server

tests.local:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml run --rm server go test ./...

migrate.create:
	docker-compose exec server migrate create -tz Europe/Moscow -ext sql -dir ./migrations/pg ${name}

migrate.lite.create:
	docker-compose exec server migrate create -tz Europe/Moscow -ext sql -dir ./migrations/sqlite ${name}





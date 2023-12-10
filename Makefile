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

restart:
	docker-compose stop
	make up.prod

tests.agent:
	docker build -t metrics-agent-tests -f Docker/Dockerfile-agent --target=tests .
	docker run --rm --name metrics-agent-tests-container metrics-agent-tests
	#docker run --rm -it --name metrics-agent-tests-container metrics-agent-tests bash

tests.server:
	docker build -t metrics-server-tests -f Docker/Dockerfile-server --target=tests .
	docker run --rm --name metrics-server-tests-container metrics-server-tests
	#docker run --rm -it --name metrics-server-tests-container metrics-server-tests bash

run.server:
	docker-compose run --rm --build server sh




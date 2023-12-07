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






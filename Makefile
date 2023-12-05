up.debug:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

stop:
	docker-compose stop

logs.agent:
	docker-compose logs -f agent

up.prod:
	docker-compose up --build






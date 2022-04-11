.PHONY: up_dev	
up_dev:
	docker-compose -f docker-compose.dev.yml up -d --build

.PHONY: down_dev
down_dev:
	docker-compose -f docker-compose.dev.yml down -v
	
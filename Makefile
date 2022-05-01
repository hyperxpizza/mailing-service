.PHONY: up_dev	
up_dev:
	docker-compose -f docker-compose.dev.yml up -d --build

.PHONY: down_dev
down_dev:
	docker-compose -f docker-compose.dev.yml down -v

.PHONY: psql
psql:
	docker exec -it mailing-service_postgres_1 psql -d mailingServiceDB -U mailingServiceDBuser
	
.PHONY: test
test:
	go test -v ./tests/ --run TestMailingServer --config=/home/hyperxpizza/dev/golang/reusable-microservices/mailing-service/config.dev.json
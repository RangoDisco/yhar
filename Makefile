build:
	docker compose build --no-cache
up:
	docker compose up -d
start: build up

build-prod:
	docker compose -f compose.yaml -f compose.prod.yaml build --no-cache
up-prod:
	docker compose -f compose.yaml -f compose.prod.yaml up -d

start-prod: build-prod up-prod

start-test:
	docker compose -f compose.test.yaml build --no-cache && \
    docker compose -f compose.test.yaml up -d && \
    DSN="user=yhar password=yhar dbname=yhar host=localhost port=54020 sslmode=disable" JWT_SECRET="SECRETJWT" REFRESH_SECRET="SECRETREFRESH"  go test -v ./...
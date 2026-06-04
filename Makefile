start-cold-dev:
	docker compose build --no-cache && \
    docker compose up -d
start-dev:
	docker compose -f compose.yaml -f compose.override.yaml up -d
start-cold-prod:
	docker compose -f compose.yaml -f compose.prod.yaml build --no-cache && \
    docker compose -f compose.yaml -f compose.prod.yaml up -d
start-cold-test:
	docker compose -f compose.test.yaml build --no-cache && \
    docker compose -f compose.test.yaml up -d && \
    DSN="user=yhar password=yhar dbname=yhar host=localhost port=54020 sslmode=disable" JWT_SECRET="SECRETJWT" REFRESH_SECRET="SECRETREFRESH"  go test -v ./...
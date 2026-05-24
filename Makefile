start-cold-dev:
	docker compose build --no-cache && \
    docker compose up -d
start-dev:
	docker compose -f compose.yaml -f compose.override.yaml up -d
start-cold-prod:
	docker compose -f compose.yaml -f compose.prod.yaml build --no-cache && \
    docker compose -f compose.yaml -f compose.prod.yaml up -d
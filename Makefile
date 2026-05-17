docker-dev-start:
	docker compose -f compose.yaml -f compose.override.yaml build --no-cache && \
    docker compose -f compose.yaml -f compose.override.yaml up -d
docker-prod-start:
	docker compose -f compose.yaml -f compose.prod.yaml build --no-cache && \
    docker compose -f compose.yaml -f compose.prod.yaml up -d
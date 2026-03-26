# Yhar

![GitHub commit activity](https://img.shields.io/github/commit-activity/w/rangodisco/yhar)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/rangodisco/yhar/latest)
![GitHub last commit](https://img.shields.io/github/last-commit/rangodisco/yhar)

# What it is

Self-hosted scrobbler + scrobbler database.

### Testing the app (only the API for now)

```shell
# Prepare database
docker compose -f compose.test.yml up -d

# Run tests
go test -v ./...
```
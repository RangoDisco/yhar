# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.26-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GOOS=linux go build -o /yhar-api ./cmd/

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:3.21 AS build-release-stage

WORKDIR /

COPY --from=build-stage /yhar-api /yhar-api

EXPOSE 8080

ENTRYPOINT ["/yhar-api"]
# Build the application from source
FROM golang:1.26 AS gin-build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /yhar-api ./cmd/

FROM gin-build-stage as gin-dev
ENV GIN_MODE=debug
RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN go install github.com/go-delve/delve/cmd/dlv@latest
# Run the air command in the directory where our code will live
WORKDIR /app
CMD ["air"]

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:3.21 AS gin-release
COPY --from=gin-build-stage /yhar-api /yhar-api
ENV GIN_MODE=release
WORKDIR /
EXPOSE 8080
ENTRYPOINT ["/yhar-api"]
FROM golang:1.24 AS base-stage

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

# Run the tests in the container
FROM build-stage AS run-test-stage

RUN go test -v ./...

# Build the application from source
FROM base-stage AS build-stage

WORKDIR /app

COPY . .

RUN go build -o ./AuthService

EXPOSE 9090

CMD [ "/app/AuthService" ]

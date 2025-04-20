ENV=DEVELOPMENT
# ENV=PRODUCTION
POSTGRES_DSN='host=localhost user=postgres password=password dbname=authservice port=5432 sslmode=disable TimeZone=America/Los_Angeles'

build:
	go build -o bin/authservice ./main.go

test:
	go test -v ./... -coverprofile=coverage.out

run:
	./bin/authservice

build-image:
	docker build -t authservice .

up:
	docker compose -p authservice --env-file ./.env up -d

down:
	docker compose -p authservice down
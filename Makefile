format:
	go fmt ./...

lint:
	golangci-lint run

test-unit:
	cd api && go test ./... -v

start-all-services:
	cd e2e && docker-compose build && docker-compose up -d

stop-all-services:
	cd e2e && docker-compose down

test-e2e: start-all-services
	sleep 5
	cd e2e && go test ./... -v
	make stop-all-services
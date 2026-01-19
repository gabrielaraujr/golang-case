.PHONY: configure logs build test clean

configure:
	docker compose up -d --build

logs:
	docker compose logs -f

build: build-account build-risk

build-account:
	cd account && go build -o bin/account cmd/main.go

build-risk:
	cd risk-analysis && go build -o bin/risk-analysis cmd/main.go

run-account:
	cd account && go run cmd/main.go

run-risk:
	cd risk-analysis && go run cmd/main.go

test: test-account test-risk

test-account:
	cd account && go test ./...

test-risk:
	cd risk-analysis && go test ./...

clean:
	docker compose down -v
	rm -rf account/bin risk-analysis/bin
	go clean -cache

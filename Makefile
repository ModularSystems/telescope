build:
	mkdir -p bin
	go build -o bin/telescope ./cmd/telescope/main.go

docker:
	docker build -t modularsystems/telescope:latest .

test:
	go test -v pkg/scan/*

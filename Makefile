build:
	mkdir -p bin
	go build -o bin/telescope ./cmd/telescope/

test:
	go test -v pkg/scanner/*

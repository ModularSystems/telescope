build:
	mkdir -p bin
	go build -o bin/telescope ./cmd

test:
	go test -v pkg/scanner/*

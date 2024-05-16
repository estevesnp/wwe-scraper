build:
	@go build -o bin/scrapewwe main.go

run:build
	@bin/scrapewwe

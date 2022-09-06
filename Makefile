build:
	go build -o bin/ ./cmd/youtube-fetch-search/main.go

run:
	go run ./cmd/youtube-fetch-search/main.go

build-docker:
	docker build -t youtube-fetch-and-search .

run-docker:
	docker run -p 8080:8080 --network="host" youtube-fetch-and-search  
build:
	go build -o bin/dcache

run: build
	./bin/dcache

follower_run: build
	./bin/dcache --listenaddr :4000 --leaderaddr :3000

tests:
	go test -v ./...
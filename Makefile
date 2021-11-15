
build-image:
	docker build -t zerohash/vwap:latest .

go-run:
	go run cmd/server.go

docker-run:
	docker run --rm zerohash/vwap:latest

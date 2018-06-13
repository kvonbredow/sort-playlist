all: docker

docker: sort-playlist
	docker build -t sort-playlist .

sort-playlist: main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sort-playlist .

clean:
	rm -f sort-playlist
	docker image rm -f sort-playlist:latest

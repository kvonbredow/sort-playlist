all: docker

docker: get-songinfo
	docker build -t get-songinfo .

get-songinfo: main.go ../get-songinfo.pb.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o get-songinfo .

clean:
	rm -f get-songinfo
	docker image rm -f get-songinfo:latest

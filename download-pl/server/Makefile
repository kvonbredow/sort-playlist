all: docker

docker: download-pl ca-certificates.crt
	docker build -t download-pl .

download-pl: main.go ../download-pl.pb.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o download-pl .

clean:
	rm -f download-pl
	docker image rm -f download-pl:latest

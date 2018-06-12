all: docker

docker: download-pl
	docker build -t download-pl .

download-pl: main.go
	go build -a -installsuffix cgo -o download-pl .

clean:
	rm -f download-pl
	docker image rm -f download-pl:latest

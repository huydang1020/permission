start:
	go build 
	./permission start
cdb:
	go build 
	./permission createDb
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o permission .	
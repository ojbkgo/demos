all: local remote

local:
	go build -o vtun ./cmd/vtun
	chmod 0777 vtun

remote:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tund ./cmd/tund
	chmod 0777 tund
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vtun-remote ./cmd/vtun
	#go build -o tund-local ./cmd/tund
	#chmod 0777 tund tund-local
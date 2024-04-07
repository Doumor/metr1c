build :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .
clean :
	rm metr1c

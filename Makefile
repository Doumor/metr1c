build :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .
	tar -czf metr1c.tar.gz metr1c metr1c.service
	rm metr1c
clean :
	rm metr1c.tar.gz

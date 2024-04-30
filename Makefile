build_dev :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-w" .

build :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" .

tar : build
	tar -v -czf metr1c.tar.gz metr1c metr1c.service
	rm -v metr1c

clean :
	rm -v -f metr1c.tar.gz metr1c

install : build
	mkdir -v -p /opt/metr1c
	install -v -m 755 ./metr1c /opt/metr1c/metr1c
	/usr/bin/install -v -b -S .bak -m 750 -o root -g root ./metr1c.service /etc/systemd/user/
	@echo "Now open the /etc/systemd/user/metr1c.service file and edit variables"

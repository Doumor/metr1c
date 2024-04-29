build_dev :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .

build :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w" .

tar :
	tar -czf metr1c.tar.gz metr1c metr1c.service
	rm metr1c

clean :
	rm -f metr1c.tar.gz metr1c

install : build
	mkdir -p /opt/metr1c
	install -m 755 ./metr1c /opt/metr1c/metr1c
	/usr/bin/install -b -S .bak -m 750 -o root -g root ./metr1c.service /etc/systemd/user/
	@echo "Now open the /etc/systemd/system/metr1c.service file and edit varriables"

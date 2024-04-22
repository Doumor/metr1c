build :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .
	tar -czf metr1c.tar.gz metr1c metr1c.service
	rm metr1c
clean :
	rm metr1c.tar.gz
install : build
	tar -zxvf ./metr1c.tar.gz
	rm ./metr1c.tar.gz
	mkdir -p /opt/metr1c
	install -m 755 ./metr1c /opt/metr1c/metr1c
	/usr/bin/install -b -s -S .bak -m 750 -o root -g root ./metr1c.service /etc/systemd/system/
	echo "Now open the /etc/systemd/system/metr1c.service file and edit varriables"

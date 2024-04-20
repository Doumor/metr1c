build :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .
	tar -czf metr1c.tar.gz metr1c metr1c.service
	rm metr1c
clean :
	rm metr1c.tar.gz
install :
	tar -zxvf ./metr1c.tar.gz
	rm ./metr1c.tar.gz
	mkdir /opt/metr1c
	mv ./metr1c /opt/metr1c/metr1c
	mv -n metr1c.service /etc/systemd/system/
	chown root:root /etc/systemd/system/metr1c.service
	chmod 750 /etc/systemd/system/metr1c.service

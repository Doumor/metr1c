TAG=$(shell git describe --abbrev=0 2> /dev/null || echo "0.0.1")
HASH=$(shell git rev-parse --verify --short HEAD)
VERSION="${TAG}-${HASH}"

build_dev :
	@printf "building version %s, stripped\n" "${VERSION}"
	@printf "building version %s, stripped\n" "${VERSION}"
        @CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v \
        -ldflags "-X main.Version=${VERSION} -w" .

metr1c :
	@printf "building version %s, stripped\n" "${VERSION}"
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v \
	-ldflags "-X main.Version=${VERSION} -s -w" .

tar : clean metr1c
	tar -v -czf metr1c.tar.gz metr1c metr1c.service
	rm -v metr1c

clean :
	rm -v -f metr1c.tar.gz metr1c

install : build
	mkdir -v -p /opt/metr1c
	install -v -m 755 ./metr1c /opt/metr1c/metr1c
	/usr/bin/install -v -b -S .bak -m 750 -o root -g root ./metr1c.service /etc/systemd/user/
	@echo "Now open the /etc/systemd/user/metr1c.service file and edit variables"

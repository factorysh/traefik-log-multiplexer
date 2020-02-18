GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

binary: bin
	go build -ldflags "-X github.com/factorysh/traefik-log-multiplxer/version.version=$(GIT_VERSION)" \
		-o bin/stfm

bin:
	mkdir -p bin

test:
	go test -v -cover github.com/factorysh/traefik-log-multiplexer/route
	go test -v -cover github.com/factorysh/traefik-log-multiplexer/fusebox
ifeq ($(shell uname),Darwin)
OS=darwin
else
OS=linux
endif
TRAEFIK1=1.7.30
TRAEFIK2=2.4.8

test:
	@echo "Nope"

bin:
	mkdir -p bin

log:
	mkdir -p log

up: log
	docker-compose up

multiplexer:
	CONFIG=./traefik-demultiplexer.yml ../bin/stfm

clean:
	rm -rf bin

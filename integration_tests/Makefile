ifeq ($(shell uname),Darwin)
OS=darwin
else
OS=linux
endif
TRAEFIK1=1.7.30
TRAEFIK2=2.4.8
FLUENTBIT=1.7.4

test:
	@echo "Nope"

bin:
	mkdir -p bin

log:
	mkdir -p log

bin/traefik: bin
	@echo ${OS}
	curl -o bin/traefik -L https://github.com/traefik/traefik/releases/download/v${TRAEFIK1}/traefik_${OS}-amd64
	file bin/traefik
	chmod +x bin/traefik

traefik: bin/traefik log
	./bin/traefik -c traefik.cfg

up:
	docker-compose up

clean:
	rm -rf bin

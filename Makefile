prefix?=$(HOME)/.bin

fetch:
	go get

install: fetch
	go install

test: test/shunt.sh fetch
	@./test/shunt.sh --verbose ./test/dockme_test.sh

test/shunt.sh:
	mkdir -p test
	curl -sSL https://raw.githubusercontent.com/odb/shunt/master/shunt.sh > /tmp/shunt.sh
	install -m 0755 /tmp/shunt.sh test

build:
	go run dockme.go --config Buildme.yml

docker/build: fetch clean
	bash ./scripts/build.all.bash
	find builds | xargs chmod 777

clean:
	rm -rf builds
	rm -rf test/shunt.sh

.PHONY: test

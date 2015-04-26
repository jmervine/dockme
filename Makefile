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
	find builds | grep -v "*.md5" | xargs chmod 755
	find builds -type f -name "*.md5" | xargs chmod 644

clean:
	rm -rf builds
	rm -rf test/shunt.sh

examples:
	go run dockme.go -D --save --sudo -T ruby -C Dockme.yml.example
	go run dockme.go -D --save --sudo -T node -C ./examples/SudoNode.yml
	go run dockme.go -D --save -T node -C ./examples/Node.yml
	go run dockme.go -D --save -T nodebox -C ./examples/Nodebox.yml
	go run dockme.go -D --save -T ruby -C ./examples/Ruby.yml
	go run dockme.go -D --save -T rails -C ./examples/Rails.yml
	go run dockme.go -D --save -T rails -C ./examples/Rails.yml
	go run dockme.go -D --save -T python2 -C ./examples/Python2.yml
	go run dockme.go -D --save -T python3 -C ./examples/Python3.yml

finalize: clean test build examples

.PHONY: test examples

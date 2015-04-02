prefix?=$(HOME)/.bin

install: $(prefix)/dockme

update: uninstall install

uninstall:
	rm $(prefix)/dockme

$(prefix)/dockme:
	install -m 0700 dockme $@

test: test/shunt.sh
	@./test/shunt.sh --verbose ./test/dockme_test.sh

test/shunt.sh:
	mkdir -p test
	curl -sSL https://raw.githubusercontent.com/odb/shunt/master/shunt.sh > /tmp/shunt.sh
	install -m 0755 /tmp/shunt.sh test

.PHONY: test

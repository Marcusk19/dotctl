clean:
	rm -r test/bender_test
	rm -r tmp

sandbox:
	mkdir -p ./tmp/
	cp -r ~/.config/ ./tmp/config

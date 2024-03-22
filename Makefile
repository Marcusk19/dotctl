clean:
	rm -rf test/bender_test 2> /dev/null
	rm -rf tmp 2> /dev/null

sandbox:
	mkdir -p ./tmp/ 2> /dev/null
	cp -r ~/.config/ ./tmp/config 2> /dev/null

unit-test:
	TESTING=true go test -v ./test

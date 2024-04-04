clean:
	rm -rf test/dotctl_test 2> /dev/null
	rm -rf tmp 2> /dev/null

sandbox:
	mkdir -p ./tmp/ 2> /dev/null
	cp -r ~/.config/ ./tmp/config 2> /dev/null

unit-test:
	TESTING=true go test -v ./test
	rm -rf test/dotctl_test 2> /dev/null

pre-commit-hooks:
	pre-commit autoupdate
	pre-commit install

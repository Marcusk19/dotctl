unit-test:
	TESTING=true go test -v ./test
	rm -rf test/dotctl_test 2> /dev/null

install:
	go build
	cp dotctl /usr/local/bin

pre-commit-hooks:
	pre-commit autoupdate
	pre-commit install

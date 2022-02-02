
main:
	go run main.go ${ARG}

test:
	./test.sh

clean:
	rm -f main tmp*

.PHONY: test clean
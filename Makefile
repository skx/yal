#
# Simple Makefile to run our lisp-based tests, and update our examples/
# index
#

ALL: test-lisp test-go update-examples


#
# Build our binary
#
yal:
	go build .

#
# Run our lisp-based tests
#
test-lisp: yal
	./yal examples/lisp-tests.lisp | _misc/tapview


#
# Run our go-based tests
#
test-go:
	go test ./...


#
# Update our list of examples.
#
update-examples: yal
	cd examples && make

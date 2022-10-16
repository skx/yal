#!/bin/bash

# Run our golang tests
go test ./... -race

# Run the perl script to look for function orders
.github/test-ordering.pl || exit 1

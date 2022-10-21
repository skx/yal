[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/yal)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/yal)](https://goreportcard.com/report/github.com/skx/yal)
[![license](https://img.shields.io/github/license/skx/yal.svg)](https://github.com/skx/yal/blob/master/LICENSE)

* [yet another lisp](#yet-another-lisp)
* [Building / Installing](#building--installing)
* [Standard Library](#standard-library)
* [Usage](#usage)
* [Examples](#examples)
* [Fuzz Testing](#fuzz-testing)
* [Benchmark](#benchmark)
* [See Also](#see-also)


# yet another lisp


* [A brief introduction to using this lisp](INTRODUCTION.md).
  * Getting started setting variables, defining functions, etc.
* [A list of primitives we have implemented](PRIMITIVES.md).
  * This describes the functions we support, whether implemented in lisp or golang.
  * For example `(car)`, `(cdr)`, `(file:lines)`, `(shell)`, etc.



## Building / Installing

If you have [the yal repository](https://github.com/skx/yal) cloned locally then
you should be able to build and install in the standard way:

```sh
$ go build .
$ go install .
```

If you don't have the repository installed, but you have a working golang toolset then installation should be as simple as:

```sh
$ go install github.com/skx/yal@latest
```

If neither of those options suit, you may download the most recent binary from our [release page](https://github.com/skx/yal/releases).  Remember that if you're running a Mac you'll need to remove the quarantine flag which _protects you_ from unsigned binaries, for example:

```sh
% xattr  -d com.apple.quarantine yal-darwin-amd64
% chmod 755 com.apple.quarantine yal-darwin-amd64
```



## Usage

Once installed there are two ways to execute code:

* By specifying an expressions on the command-line:
  * `yal -e '(print (os))'`
* By passing the name of a file containing lisp code to read and execute:
  * `yal test.lisp`

The yal interpreter allows (optional) documentation to be attached to functions, both those implemented in golang and those written in lisp, there is another command-line flag to dump that information from the standard library and built-in functions:

* `yal -h [regexp]`
  * By default this will show the help for all available functions, in sorted order.
  * If you specify any regular expressions then any entry which matches the given patterns will be displayed.

Finally if you've downloaded a binary release from [our release page](https://github.com/skx/yal/releases) the `-v`flag will show you what version you're running:

```sh
% yal-darwin-amd64 -v
v0.11.0 f21d032e812ee6eadad5eac23f079a11f5e1041a
```



## Standard Library

When user-code is executed, whether a simple statement supplied via the command-line, or read from a file, a standard-library is loaded from beneath the directory:

* [stdlib/stdlib/](stdlib/stdlib/)


Our standard-library consists of primitive functions such as `(map..)`, `(min..)` and similar, is written in 100% yal-lisp.

The standard library may be entirely excluded via the use of the environmental varilable `YAL_STDLIB_EXCLUDE_ALL`:

```
$ yal  -e "(print (hms))"
22:30:57

$ YAL_STDLIB_EXCLUDE_ALL=true yal  -e "(print (hms))"
Error running: error expanding argument [hms] for call to (print ..):
  ERROR{argument 'hms' not a function}
```

If you prefer you may exclude specific _parts_ of the standard library, by specify the filenames you wish to exclude separated by commas:

```
$ YAL_STDLIB_EXCLUDE=date,type-checks yal  -e "(print (hms))"
22:30:57
```



## Examples

A reasonable amount of sample code can be found in the various included examples:

* [test.lisp](test.lisp) shows many things.
* [fibonacci.list](fibonacci.lisp) calculate the first 25 numbers of the Fibonacci sequence.
* [fizzbuzz.lisp](fizzbuzz.lisp) is a standalone sample of solving the fizzbuzz problem.
* [mtest.lisp](mtest.lisp) shows some macro examples.

As noted there is a standard-library of functions which are loaded along with any user-supplied script - that library of functions may also provide a useful reference and example of yal-code:

* [stdlib/stdlib/](stdlib/stdlib/)

Running any of our supplied examples should produce useful output for reference.  For example here's the result of running the `fibonacci.lisp` file:

```lisp
$ yal fibonacci.lisp
1st fibonacci number is 1
2nd fibonacci number is 1
3rd fibonacci number is 2
4th fibonacci number is 3
5th fibonacci number is 5
6th fibonacci number is 8
7th fibonacci number is 13
8th fibonacci number is 21
9th fibonacci number is 34
10th fibonacci number is 55
11th fibonacci number is 89
12th fibonacci number is 144
13th fibonacci number is 233
14th fibonacci number is 377
15th fibonacci number is 610
16th fibonacci number is 987
17th fibonacci number is 1597
18th fibonacci number is 2584
19th fibonacci number is 4181
20th fibonacci number is 6765
21st fibonacci number is 10946
22nd fibonacci number is 17711

```



## Fuzz Testing

The project has 100% test-coverage of all the internal packages, using the standard go facilities you can run those test-cases:

```sh
go test ./...
```

In addition to the static-tests there is also support for the integrated fuzz-testing facility which became available with go 1.18+.  Fuzz-testing essentially feeds the interpreter random input and hopes to discover crashes.

You can launch a series of fuzz-tests like so:

```sh
go test -fuzztime=300s -parallel=1 -fuzz=FuzzYAL -v
```

Sample output will look like this:

```
=== FUZZ  FuzzYAL
...
fuzz: elapsed: 56m54s, execs: 163176 (0/sec), new interesting: 108 (total: 658)
fuzz: elapsed: 56m57s, execs: 163176 (0/sec), new interesting: 108 (total: 658)
fuzz: elapsed: 57m0s, execs: 163183 (2/sec), new interesting: 109 (total: 659)
fuzz: elapsed: 57m3s, execs: 163433 (83/sec), new interesting: 110 (total: 660)
..
```

If you find a crash then it is either a bug which needs to be fixed, or a false-positive (i.e. a function reports an error which is expected) in which case the fuzz-test should be updated to add it to the list of known-OK results.  (For example "division by zero" is a fatal error, so that's a known-OK result).




## Benchmark

There is a simple benchmark included within this repository, computing the factorial of 100, to run this execute execute:

```sh
$ go test -run=Bench -bench=.
```

To run the benchmark for longer add `-benchtime=30s`, or similar, to the command-line.

I also put together an external comparison of my toy scripting languages here:

* [Toy Language Benchmarks](https://github.com/skx/toy-language-benchmarks)

This shows that the Lisp implementation isn't so slow, although it is not the fasted of the scripting languages I've implemented.



## See Also

This repository was put together after [experimenting with a scripting language](https://github.com/skx/monkey/), an [evaluation engine](https://github.com/skx/evalfilter/), putting together a [TCL-like scripting language](https://github.com/skx/critical), writing a [BASIC interpreter](https://github.com/skx/gobasic) and creating [tutorial-style FORTH interpreter](https://github.com/skx/foth).

I've also played around with a couple of compilers which might be interesting to refer to:

* Brainfuck compiler:
   * [https://github.com/skx/bfcc/](https://github.com/skx/bfcc/)
* A math-compiler:
  * [https://github.com/skx/math-compiler](https://github.com/skx/math-compiler)

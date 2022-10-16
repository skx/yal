[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/yal)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/yal)](https://goreportcard.com/report/github.com/skx/yal)
[![license](https://img.shields.io/github/license/skx/yal.svg)](https://github.com/skx/yal/blob/master/LICENSE)

* [yet another lisp](#yet-another-lisp)
* [Special Features](#special-features)
* [Building / Installing](#building--installing)
* [Usage](#usage)
* [Examples](#examples)
* [Omissions](#omissions)
* [Fuzz Testing](#fuzz-testing)
* [Benchmark](#benchmark)
* [References](#references)
* [See Also](#see-also)


# yet another lisp


* [A brief introduction to using this lisp](INTRODUCTION.md)
  * Getting started setting variables, defining functions, etc.
* [A list of primitives we have implemented](PRIMITIVES.md)
  * This describes the functions we support, whether implemented in lisp or golang.
  * For example `(car)`, `(cdr)`, `(file:lines)`, `(shell)`, etc.



## Building / Installing

If you have [this repository](https://github.com/skx/yal) cloned locally then
you should be able to build and install in the standard way:

```sh
$ go build .
$ go install .
```

If you don't have the repository installed, but you have a working golang environment then installation should be as simple as:

```sh
$ go install github.com/skx/yal@latest
```

Finally if neither of those options suit, you may download the most recent binary from our [release page](https://github.com/skx/yal/releases).



## Usage

Once installed there are two ways to execute code:

* By specifying an expressions on the command-line:
  * `yal -e "(print (os))"`
* By passing the name of a file containing lisp code to read and execute:
  * `yal test.lisp`

As our interpreter allows documentation to be attached to functions, both those implemented in golang and those written in lisp, we also have a flag to dump that information:

* `yal -h`
  * Shows all functions which contain help-text, in sorted order.
  * Examples are included where available.




## Examples

A reasonable amount of sample code can be found in the various included examples:

* [test.lisp](test.lisp) shows many things.
* [fibonacci.list](fibonacci.lisp) calculate the first 25 numbers of the Fibonacci sequence.
* [fizzbuzz.lisp](fizzbuzz.lisp) is a standalone sample of solving the fizzbuzz problem.
* [mtest.lisp](mtest.lisp) shows some macro examples.

As noted there is a standard-library of functions which are loaded along with any user-supplied script.  These functions are implemented in lisp and also serve as a demonstration of syntax and features:

* [stdlib/stdlib.lisp](stdlib/stdlib.lisp)
* [stdlib/mal.lisp](stdlib/mal.lisp)

Running these example will produce output, for example:

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

I also put together a small comparison of toy scripting languages available here:

* [Toy Language Benchmarks](https://github.com/skx/toy-language-benchmarks)

This shows that the Lisp implementation isn't so slow, although it is not the fasted of the scripting languages I've implemented.




## References

* https://github.com/thesephist/klisp/blob/main/lib/klisp.klisp
  * Very helpful _inspiration_ for writing primitives in Lisp.
* https://github.com/kanaka/mal/
  * Make A Lisp, very helpful for the quoting, unquoting, and macro magic.
* https://lispcookbook.github.io/cl-cookbook/macros.html
  * The Common Lisp Cookbook – Macros
* http://soft.vub.ac.be/~pcostanz/documents/08/macros.pdf
  * The source of the cute "while" macro, and a good read beyond that.




## See Also

This repository was put together after [experimenting with a scripting language](https://github.com/skx/monkey/), an [evaluation engine](https://github.com/skx/evalfilter/), putting together a [TCL-like scripting language](https://github.com/skx/critical), writing a [BASIC interpreter](https://github.com/skx/gobasic) and creating [tutorial-style FORTH interpreter](https://github.com/skx/foth).

I've also played around with a couple of compilers which might be interesting to refer to:

* Brainfuck compiler:
   * [https://github.com/skx/bfcc/](https://github.com/skx/bfcc/)
* A math-compiler:
  * [https://github.com/skx/math-compiler](https://github.com/skx/math-compiler)

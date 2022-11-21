[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/yal)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/yal)](https://goreportcard.com/report/github.com/skx/yal)
[![license](https://img.shields.io/github/license/skx/yal.svg)](https://github.com/skx/yal/blob/master/LICENSE)

* [yet another lisp](#yet-another-lisp)
* [Building / Installing](#building--installing)
* [Standard Library](#standard-library)
* [Usage](#usage)
  * [Integrated Help](#integrated-help)
  * [REPL Helper](#repl-helper)
* [Examples](#examples)
* [Fuzz Testing](#fuzz-testing)
* [Benchmark](#benchmark)
* [See Also](#see-also)


# yet another lisp


* [A brief introduction to using this lisp](INTRODUCTION.md).
  * Getting started setting variables, defining functions, etc.
  * This includes documentation on enhanced features such as
    * Hashes.
    * Structures.
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

If neither of those options suit, you may download the most recent binary from our [release page](https://github.com/skx/yal/releases).

Remember that if you're running a Mac you'll need to remove the quarantine flag which _protects you_ from unsigned binaries, for example:

```sh
% xattr  -d com.apple.quarantine yal-darwin-amd64
% chmod 755 com.apple.quarantine yal-darwin-amd64
```



## Usage

Once installed there are three ways to execute code:

* By specifying an expression to execute upon the command-line:
  * `yal -e '(print (os))'`
* By passing the name of a file containing lisp code to read and execute:
  * `yal examples/test.lisp`
* By launching the interpreter with zero arguments, which will start the interactive REPL mode.
  * If present the file `~/.yalrc` is loaded before the REPL starts.
  * Here is a sample [.yalrc](.yalrc) file which shows the kind of thing you might wish to do.

Finally if you've downloaded a binary release from [our release page](https://github.com/skx/yal/releases) the `-v` flag will show you what version you're running:

```sh
% yal-darwin-amd64 -v
v0.11.0 f21d032e812ee6eadad5eac23f079a11f5e1041a
```


### Integrated Help

The yal interpreter allows (optional) documentation to be attached to functions, both those implemented in the core, and those which are added in lisp.

You can view the help output by launching with the `-h` flag:

    $ yal -h

By default all the help-text contained within the standard-library, and our built-in primitives, will be shown.  You may limit the display to specific function(s) by supplying an arbitrary number of regular expression, for example:

    $ yal -h count execute
    count (arg)
    ===========
    count is an alias for length.

    load-file (filename)
    ====================
    Load and execute the contents of the supplied filename.

When you specify a regular expression, or more than one, the matches will be applied to the complete documentation for each function.  So the term "foo" will match the term "foo" inside the explanation of the function, the argument list, and the function name itself.

A good example of the broad matching would include the term "length":

    $ yal -h length | grep -B1 ==
    apply-pairs (lst:list fun:function)
    ===================================
    --
    count (arg)
    ===========
    --
    length (arg)
    ============
    --
    pad:left (str add len)
    ======================
    --
    pad:right (str add len)
    =======================
    --
    repeated (n:number x)
    =====================
    --
    strlen (str:string)
    ===================



### REPL Helper

If you wish to get command-line completion, history, etc, within the REPL-environment you might consider using the `rlwrap` tool.

First of all output a list of the names of each of the built-in function:

     $ yal -e "(apply (env) (lambda (x) (print (get x :name))))" > functions.txt

Now launch the REPL with completion on those names:

     $ rlwrap --file functions.txt ./yal




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

If you prefer you may exclude specific _parts_ of the standard library, by specifying a comma-separated list of regular expressions:

```
$ YAL_STDLIB_EXCLUDE=date,type-checks yal  -e "(print (hms))"
22:30:57
```

Here the regular expressions will be matched against the name of the file(s) in the [standard library directory](stdlib/stdlib/).



## Examples

A reasonable amount of sample code can be found beneath the [examples/](examples/) directory, including:

* [examples/fibonacci.list](examples/fibonacci.lisp)
  * Calculate the first 25 numbers of the Fibonacci sequence.
* [examples/fizzbuzz.lisp](examples/fizzbuzz.lisp)
  * A standalone sample of solving the fizzbuzz problem.
* [examples/mtest.lisp](examples/mtest.lisp)
  * Shows simple some macro examples, but see [examples/lisp-tests.lisp](examples&lisp-tests.lisp) for a more useful example.
    * This uses macros in an interesting way.
    * It is also used to actually test the various Lisp-methods we've implemented.
* [examples/sorting.lisp](examples/sorting.lisp)
  * Demonstrates writing & benchmarking sorting-routines.
* [examples/test.lisp](examples/test.lisp)
  * A misc. collection of sample code, functions, and notes.

As noted there is a standard-library of functions which are loaded along with any user-supplied script - that library of functions may also provide a useful reference and example of yal-code:

* [stdlib/stdlib/](stdlib/stdlib/)

The standard-library contains its own series of test-cases written in Lisp:

* [examples/lisp-tests.lisp](examples/lisp-tests.lisp)

The lisp-tests.lisp file contains a simple macro for defining test-cases, and uses that to good effect to test a range of our lisp-implemented primitives.



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

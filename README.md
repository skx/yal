[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/yal)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/yal)](https://goreportcard.com/report/github.com/skx/yal)
[![license](https://img.shields.io/github/license/skx/yal.svg)](https://github.com/skx/yal/blob/master/LICENSE)

* [yet another lisp](#yet-another-lisp)
* [Brief Overview](#brief-overview)
* [Features](#features)
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

Another trivial/toy Lisp implementation in Go.




## Brief Overview

To set a global variable use `set!`:

    (set! foo "bar")

To start a new scope, with local variables, use `let*`:

    (let (foo "bar"
         baz  "bart")
      (print "foo is %s" foo)
      (print "baz is %s" baz)
      ;...
    )

To define a function use `set!` with `fn*`:

    (set! fact (fn* (n)
      (if (<= n 1)
        1
          (* n (fact (- n 1))))))

To define a macro use `defmacro!`:

    (defmacro! debug (fn* (x) `(print "Variable '%s' has value %s" '~x ~x)))

    ; Using this macro
    (set! foo "steve")
    (debug foo)

That concludes the brief overview, note that `lambda` can be used as a synonym for `fn*`, and other synonyms exist.  In the interests of simplicity they're not covered here.





## Features

We have a reasonable number of functions implemented, either in our golang core or in our standard-library (which is implemented in yal itself):

* Support for strings, numbers, errors, lists, hashes, etc.
  * `#t` is the true symbol, `#f` is false.
  * `true` and `false` are available as synonyms.
* Hash operations:
  * Hashes are literals like this `{ :name "Steve" :location "Helsinki" }`
  * Hash functions are `contains?`, `get`, `keys`, `set`, & `vals`.
    * Note that keys are returned in sorted order.  Values are returned in order of their sorted keys too.
* List operations:
  * `car`, `cdr`, `cons`, `list`, & `sort`.
* Logical operations:
  * `and`, & `or`.
* Mathematical operations:
  * `+`, `-`, `*`, `/`, `#`, & `%`.
* String operations:
  * `chr`, `join`, `match` (regular-expression matching),`ord`, & `split`.
* Comparison functions:
  * `<`, `<=`, `>`, `>=`, `=`, & `eq`.
* Misc features:
  * `arch`, `getenv`, `os`, `print`, `slurp`, `sprintf`, `str` & `type`
* Special forms:
  * `begin`, `define`, `do`, `env`, `eval`, `gensym`, `if`, `lambda`, `let`, `macroexpand`, `read`, `set!`, `quote`, & `quasiquote`.
* Error handling:
  * `error`, `try`, and `catch` - as demonstrated in [try.lisp](try.lisp).
* Tail recursion optimization.
* MAL compatability:
  * `do` can be used as a synonym for `begin`.
  * `def!` can be used as a synonym for `define`.
  * `defmacro!` is used to define macros.
  * `fn*` can be used as a synonym for `lambda`.
  * `let*` can be used to define a local scope.

Building upon those primitives we have a larger standard-library of functions written in Lisp such as:

* `abs`, `apply`, `append`, `cond`, `count`, `false?`, `filter`, `lower`, `map`, `min`, `max`, `ms`, `nat`, `neg`, `now`, `nth`, `reduce`, `repeat`, `reverse`, `seq`, `strlen`, `true?`, `upper`, `when`, `while`, etc.

Although the lists above should be up to date you can check the definitions to see what is currently available:

* Primitives implemented in go:
  * [builtins/builtins.go](builtins/builtins.go)
* Primitives implemented in 100% pure lisp:
  * [stdlib/stdlib.lisp](stdlib/stdlib.lisp)
  * The code in this file is essentially **prepended** to any script that is supplied upon the command-line.



## Special Features

There are a couple of areas where we've implemented special/unusual things:

* Access to command-line arguments, when used via a shebang.
  * See [args.lisp](args.lisp) for an example.
* Introspection via the `(env)` function, which will return details of all variables/functions in the environment.
  * Allowing dynamic invocation shown in [dynamic.lisp](dynamic.lisp) and other neat things.
* Support for hashes as well as lists/strings/numbers/etc.
  * A hash looks like this `{ :name "Steve" :location "Helsinki" }`
  * Sample code is visible in [hash.lisp](hash.lisp).
* Type checking for function parameters.
  * Via a `:type` suffix.  For example `(lambda (a:string b:number) ..`.
* Support for macros.
  * See [mtest.lisp](mtest.lisp) for some simple tests/usage examples.
  * The standard library uses macros to implement the `(cond)`, and `(while)` functions, for example.


Here's an example of type-checking on a parameter value, in this case a list is required, via the `:list` suffix:

```lisp
(define blah (lambda (a:list) (print "I received the list %s" a)))

(blah '(1 2 3))    ; => "I received the list (1 2 3)"
(blah #f)          ; => Error running: argument a to blah was supposed to be list, but got false
(blah 3)           ; => Error running: argument a to blah was supposed to be list, but got 3
```

The following type suffixes are permitted and match what you'd expect:

* `:any`
* `:boolean`
* `:error`
* `:function`
* `:hash`
* `:list`
* `:nil`
* `:number`
* `:string`
* `:symbol`

If multiple types are permitted then just keep appending things, for example:

* `(define blah (lambda (a:list:number)  (print "I was given a list OR a number: %s" a)))`
  * Allows either a list, or a number.




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



## Usage

Once you've built, and optinall installed, the CLI driver there are two ways to execute code:

* By specifying sexpressions on the command-line.
  * `yal -e "(print (os))"`
* By passing the name of a file to read and execute.
  * `yal test.lisp`



## Examples

A reasonable amount of sample code can be found in the various included examples:

* [test.lisp](test.lisp) shows many things.
* [fizzbuzz.lisp](fizzbuzz.lisp) is a standalone sample of solving the fizzbuzz problem.
* [mtest.lisp](mtest.lisp) shows some macro examples.




## Fuzz Testing

The project has 100% test-coverage of all the internal packages, using the standard facilities you can run those test-cases:

```sh
go test ./...
```

In addition to that there is support for the integrated fuzz-testing which is available with go 1.18+, which essentially feeds the interpreter random input and hopes to discover crashes.

You can launch the fuzz-testing like so:

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

There is a simple benchmark included within this repository, computing the factorial of 100, to run this execute:

```sh
$ go test -run=Bench -bench=.
goos: linux
goarch: amd64
pkg: github.com/skx/yal
cpu: AMD A10-6800K APU with Radeon(tm) HD Graphics
BenchmarkGoFactorial-4    	 4752786	       248.1 ns/op
BenchmarkYALFactorial-4   	    1250	    908525 ns/op
PASS
ok  	github.com/skx/yal	2.679s
```

For longer runs add `-benchtime=30s`, or similar, to the command-line.

Here you see that the lisp version is approximately 3000% slower than the pure golang implementation.  There is a small comparison of my toy scripting languages available here:

* [Toy Language Benchmarks](https://github.com/skx/toy-language-benchmarks)

This shows that the Lisp implementation here isn't so bad!



## References

* https://github.com/thesephist/klisp/blob/main/lib/klisp.klisp
  * Very helpful "inspiration" for writing primitives in Lisp.
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

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/yal)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/yal)](https://goreportcard.com/report/github.com/skx/yal)
[![license](https://img.shields.io/github/license/skx/yal.svg)](https://github.com/skx/yal/blob/master/LICENSE)

* [yet another lisp](#yet-another-lisp)
* [Special Features](#special-features)
* [Building / Installing](#building--installing)
* [Examples](#examples)
* [Features](#features)
* [Omissions](#omissions)
* [Fuzz Testing](#fuzz-testing)
* [Benchmark](#benchmark)
* [References](#references)


# yet another lisp

Another trivial/toy Lisp implementation in Go.


## Special Features

Although this implementation is clearly derived from the [make a lisp](https://github.com/kanaka/mal/) series there are a couple of areas where we've implemented special/unusual things:

* Access to command-line arguments, when used via a shebang.
  * See [args.lisp](args.lisp) for an example.
* Introspection via the `(env)` function, which will return details of all variables/functions in the environment.
  * Allowing dynamic invocation shown in [dynamic.lisp](dynamic.lisp) and other neat things.
* Support for hashes as well as lists/strings/numbers/etc.
  * A hash looks like this `{ :name "Steve" :location "Helsinki" }`
  * Sample code is visible in [hash.lisp](hash.lisp).
* Optional parameters for functions.
  * Any parameter which is prefixed by `&` is optional, and if not specified then `nil` is assumed.
* Type checking for function parameters.
  * Via a `:type` suffix.  For example `(lambda (a:string b:number) ..`.
* Support for macros.
  * See [mtest.lisp](mtest.lisp) for some simple tests/usage examples.

Here's what optional parameters, inspired by Emacs, look like in practice:

```lisp
(define foo (lambda (a &b &c)  (print "A:%s B:%s C:%s\n", a b c)))

(foo 1 2 3)  ; => "A:1 B:2 C:3"
(foo 1 2)    ; => "A:1 B:2 C:nil"
(foo 1)      ; => "A:1 B:nil C:nil"
```

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

Once installed you can execute a file containing your lisp like so:

```sh
$ yal test.lisp
```



## Examples

A reasonable amount of sample code can be found in [test.lisp](test.lisp), but as a quick example we have a [fizzbuzz.lisp](fizzbuzz.lisp) sample:

```lisp
;;
;; This is a simple FizzBuzz example, which we can execute.
;;
;; You'll see here that we can define functions, that we have
;; primitives such as "zero?" and that we have a built-in "cond"
;; function too.
;;
;; cond here will take a list, which is processed in pairs:
;;
;;  (cond
;;    (quote
;;      TEST1  ACTION1
;;      TEST2  ACTION2
;;    )
;;  )
;;
;; For each pair (e.g. `TEST1 ACTION1`) we run the first statement, and if
;; the result is `true` we evaluate the action, and stop.
;;
;; When the test returns nil/false/similar then we continue running until
;; we do get success.  That means it is important to end with something that
;; will always succeed.
;;
;; `(quote) is used to ensure we don't evaluate the list in advance of the
;; statement.
;;

;; Is the given number divisible by 3?
;;
;; Note that we add ":number" to the end of the argument, which means
;; a fatal error will be raised if we invoke this function with a non-number,
;; for example:
;;
;;   (divByThree "Steve")
;;   (divByThree true)
;;
(define divByThree (lambda (n:number) (zero? (% n 3))))

;; Is the given number divisible by 5?
(define divByFive  (lambda (n:number) (zero? (% n 5))))

;; Run the fizz-buzz test for the given number, N
;;
;; NOTE: `and` takes a list here.
;;
(define fizz (lambda (n:number)
  (cond
    (quote
      (and (list (divByThree n) (divByFive n)))  (print "fizzbuzz")
      (divByThree n)                             (print "fizz")
      (divByFive  n)                             (print "buzz")
      #t                                         (print n)))))


;; Apply the function fizz, for each number 1-50
(apply (nat 51) fizz)
```



## Features

We have a reasonable number of functions implemented, either in our golang core or in our standard-library (which is implemented in yal itself):

* Support for strings, numbers, errors, lists, hashes, etc.
  * `#t` is the true symbol, `#f` is false.
  * `true` and `false` are available as synonyms.
* Hash operations:
  * Hashes are literals like this `{ :name "Steve" :location "Helsinki" }`
  * Hash functions are `get`, `keys`, & `set`.
* List operations:
  * `car`, `cdr`, `list`, & `sort`.
* Logical operations
  * `and`, & `or`.
* Mathematical operations:
  * `+`, `-`, `*`, `/`, `#`, & `%`.
* String operations:
  * `join`, `match` (regular-expression matching), & `split`.
* Comparison functions:
  * `<`, `<=`, `>`, `>=`, `=`, & `eq`.
* Misc features
  * `error`, `getenv`, `str`, `print`, & `type`
* Special forms
  * `begin`, `cond`, `define`, `env`, `eval`, `if`, `lambda`, `let`, `read`, `set!`, `quote`,
* Tail recursion optimization.

Building upon those primitives we have a larger standard-library of functions written in Lisp such as:

* `abs`, `apply`, `append`, `filter`, `lower`, `map`, `min`, `max`, `nat`, `neg`, `now`, `nth`, `reduce`, `reverse`, `seq`, `upper`, etc.

Although the lists above should be up to date you can check the definitions to see what is currently available:

* Primitives implemented in go:
  * [builtins/builtins.go](builtins/builtins.go)
* Primitives implemented in 100% pure lisp:
  * [stdlib/stdlib.lisp](stdlib/stdlib.lisp)
  * The code in this file is essentially **prepended** to any script that is supplied upon the command-line.




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

There is a simple benchmark included, comparing the time taken to run 100! in pure golang, and in our interpreted lisp.

To run this:

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

Here you see that the lisp version is approximately 3000% slower than the pure golang implementation.

For longer runs add `-benchtime=30s`, or similar, to the command-line.



## References

* https://github.com/thesephist/klisp/blob/main/lib/klisp.klisp
  * Very helpful "inspiration" for writing primitives in Lisp.
* https://github.com/kanaka/mal/
  * Make A Lisp, very helpful for the quoting, unquoting, and macro magic.

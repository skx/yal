[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/yal)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/yal)](https://goreportcard.com/report/github.com/skx/yal)
[![license](https://img.shields.io/github/license/skx/yal.svg)](https://github.com/skx/yal/blob/master/LICENSE)

# yet another lisp

Another trivial/toy Lisp implementation in Go.


## Special Features

Although this implementation is clearly derived from the [make a lisp](https://github.com/kanaka/mal/) series there are a couple of areas where we've implemented special/unusual things:

* Type checking
* Optional parameters

Both of these are visible only in the definition of functions.

Any parameter which is prefixed by `&` becomes optional, as you can see in this example:

```lisp
(define foo (lambda (a &b &c)  (print "A:%s B:%s C:%s\n", a b c)))

(foo 1 2 3)  ; => "A:1 B:2 C:3"
(foo 1 2)    ; => "A:1 B:2 C:nil"
(foo 1)      ; => "A:1 B:nil C:nil"
```

Similarly any parameter may have an optional ":type" annotation, and at run-time a wrong type will cause an error:

```lisp
(define blah (lambda (a:list) (print "I received the list %s" a)))

(blah '(1 2 3))    ; => "I received the list (1 2 3)"
(blah #f)          ; => Error running: argument a to blah was supposed to be list, but got false
(blah 3)           ; => Error running: argument a to blah was supposed to be list, but got 3
```

You may restrict the types to numbers (`:number`), lists (`:list`), strings (`:string`), or declare that any type is OK with `:any` (which is the same as not having an annotation at all!



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
(define fizz (lambda (n:number)
  (cond
    (quote
      (and (divByThree n) (divByFive n))  (print "fizzbuzz")
      (divByThree n)                      (print "fizz")
      (divByFive  n)                      (print "buzz")
      #t                                  (print n)))))


;; Apply the function fizz, for each number 1-50
(apply (nat 51) fizz)
```



# Features

We have a reasonable number of functions implemented in our golang core:

* Support for strings, numbers, errors, lists, etc.
  * `#t` is the true symbol, `#f` is false, though `true` and `false` are synonyms.
* List operations:
  * `car`, `cdr`, `list`, & `sort`.
* Logical operations
  * `and`, & `or`.
* Mathematical operations:
  * `+`, `-`, `*`, `/`, `#`, & `%`.
* String operations:
  * `split` and `join`.
* Comparison functions:
  * `<`, `<=`, `>`, `>=`, `=`, & `eq`.
* Misc features
  * `error`, `str`, `print`, & `type`
* Special forms
  * `begin`, `cond`, `define`, `eval`, `if`, `lambda`, `let`,  `set!`, `quote`,
* Tail recursion optimization.

Building upon those primitives we have a larger standard-library of functions written in Lisp such as:

* `abs`, `apply`, `append`, `filter`, `map`, `min`, `max`, `nat`, `neg`, `nth`, `reduce`, `reverse`, `seq`, etc.

Although the lists above should be up to date you can check the definitions to see what is currently available:

* Primitives implemented in go:
  * [builtins/builtins.go](builtins/builtins.go)
* Primitives implemented in 100% pure lisp:
  * [stdlib/stdlib.lisp](stdlib/stdlib.lisp)
  * The code in this file is essentially **prepended** to any script that is supplied upon the command-line.



## Omissions

Notable omissions here:

* No macros.
* No vectors/hashes/records.


## Fuzz Testing

If you're working with go 1.18+ you'll be able to run a fuzz-test of the interpreter, without the need for any external tools.

Run the included driver like so:

```sh
go test -fuzztime=300s -parallel=1 -fuzz=FuzzYAL -v
```

Sample output will look like this:

```
=== FUZZ  FuzzYAL
fuzz: elapsed: 0s, gathering baseline coverage: 0/38 completed
fuzz: elapsed: 0s, gathering baseline coverage: 38/38 completed, now fuzzing with 1 workers
fuzz: elapsed: 3s, execs: 39 (13/sec), new interesting: 0 (total: 38)
fuzz: elapsed: 6s, execs: 39 (0/sec), new interesting: 0 (total: 38)
fuzz: elapsed: 9s, execs: 39 (0/sec), new interesting: 0 (total: 38)
fuzz: elapsed: 12s, execs: 39 (0/sec), new interesting: 0 (total: 38)
fuzz: elapsed: 15s, execs: 39 (0/sec), new interesting: 0 (total: 38)
fuzz: elapsed: 18s, execs: 39 (0/sec), new interesting: 0 (total: 38)
fuzz: elapsed: 21s, execs: 39 (0/sec), new interesting: 0 (total: 38)
fuzz: elapsed: 24s, execs: 39 (0/sec), new interesting: 0 (total: 38)
fuzz: elapsed: 27s, execs: 39 (0/sec), new interesting: 0 (total: 38)
```

If you find a crash then it is either a bug which needs to be fixed, or a false-positive (i.e. a function reports an error which is expected) in which case the fuzz-test should be updated to add it to the list of known-OK results.  (For example "division by zero" is a fatal error, so that's a known-OK result).




## References

* https://github.com/thesephist/klisp/blob/main/lib/klisp.klisp
  * Very helpful "inspiration" for writing primitives in Lisp.

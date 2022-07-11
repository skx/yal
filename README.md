[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/yal)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/yal)](https://goreportcard.com/report/github.com/skx/yal)
[![license](https://img.shields.io/github/license/skx/yal.svg)](https://github.com/skx/yal/blob/master/LICENSE)

# yet another lisp

Another trivial/toy Lisp implementation in Go.



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
;; This is a simple example of the YAL interpreter
;;
;; You'll see here that we can define functions, that we have
;; primitives such as "zero?" and that we have a built-in "cond"
;; function too.
;;
;; cond here will take a list, which is processed in pairs:
;;
;;  (cond
;;    (quote
;;      EVAL1  ACTION1
;;      EVAL2  ACTION2
;;    )
;;  )
;;
;; We take each pair "EVAL1 ACTION1", or "EVAL2 ACTION2", and if the
;; result of evaluating the first part is true we run the action.
;;
;; If not we continue down the list.  Quote is used to ensure we don't
;; evaluate the list in advance.
;;

;; Is the given number divisible by 3?
(define divByThree (lambda (n) (zero? (% n 3))))

;; Is the given number divisible by 5?
(define divByFive  (lambda (n) (zero? (% n 5))))

;; Run the fizz-buzz test for the given number, N
(define fizz (lambda (n)
  (cond
    (quote
      (and (list (divByThree n) (divByFive n)))  (print "fizzbuzz")
      (divByThree n)                             (print "fizz")
      (divByFive  n)                             (print "buzz")
      #t                                         (print n)))))


;; Apply the function fizz, for each number 1-50
(apply (nat 50) fizz)
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
  * `+`, `-`, `*`, `/`, & `%`
* String operations:
  * `split` and `join`.
* Comparison functions:
  * `<`, `<=`, `>`, `>=`, `=`, & `eq`.
* Misc features
  * `str`, `print`, & `type`
* Special forms
  * `begin`, `define`, `if`, `lambda`, `let`,  `set!`, `quote`,
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



## References

* https://github.com/thesephist/klisp/blob/main/lib/klisp.klisp
  * Very helpful "inspiration" for writing primitives in Lisp.
* [mal - Make a Lisp](https://github.com/kanaka/mal/)

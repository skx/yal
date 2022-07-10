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

A reasonable amount of sample code can be found in [test.lisp](test.lisp), but as a quick example:

```lisp
;; A useful helper to apply a given function to each element of a list.
(define each (lambda (lst fun)
  (if (nil? lst)
    ()
      (begin
         (fun (car lst))
         (each (cdr lst) fun)))))

;; Now create a utility function to square a number
(define sq (lambda (x) (* x x)))

;; For each item in the list, print it, and the associated square.
;; Awesome!  Much Wow!
(each (list 1 2 3 4 5 6 7 8 9 10)
      (lambda (x)
        (print "%s squared is %s" x (sq x))))

```



# Features

We have a reasonable number of functions implemented in our golang core:

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

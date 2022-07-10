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

We have a decent core of functions:

* List operations `list`, `car`, `cdr`, etc.
* Mathematical operations which work with a variable number of arguments `+`, `-`, `*`, `/`
* Comparison functions `<`, `>`, etc.
* Conditionals via `if`, functions via `define`/`lambda`.
* Tail recursion optimization.
* Output via `print`, with support for format-strings.
* Decent range of standard functions `apply`, `filter`, `map`, `min`, `max`, `reduce`, etc.

Our primitives are implemented in either golang, or 100% pure lisp, and
you can inspect both sets of code:

* Primitives implemented in go:
  * [builtins/builtins.go](builtins/builtins.go)
* Primitives implemented in 100% pure lisp:
  * [stdlib/stdlib.lisp](stdlib/stdlib.lisp)
  * This code is essentially **prepended** to any script that is supplied upon the command-line.



## Omissions

Notable omissions here:

* No vectors.
* No eval.
* No macros.




## References

* https://github.com/thesephist/klisp/blob/main/lib/klisp.klisp
* [mal - Make a Lisp](https://github.com/kanaka/mal/)
* [(How to Write a (Lisp) Interpreter (in Python))](http://norvig.com/lispy.html)

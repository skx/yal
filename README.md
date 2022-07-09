# yet another lisp

Another trivial/toy Lisp implementation in Go.

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

We have a small core of built-in functions, including those you'd expect:

* List operations `list`, `car`, `cdr`, etc.
* Mathematical operations which work with a variable number of arguments `+`, `-`, `*`, `/`
* Comparison functions `<`, `>`, etc.
* Conditionals via `if`, functions via `define`/`lambda`.
* Tail recursion optimization.
* Output via `print`, with support for format-strings.


## References

* https://github.com/thesephist/klisp/blob/main/lib/klisp.klisp
* [mal - Make a Lisp](https://github.com/kanaka/mal/)
* [(How to Write a (Lisp) Interpreter (in Python))](http://norvig.com/lispy.html)

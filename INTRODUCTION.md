

# Brief Yal Introduction

To set the contents of a variable use `set!`:

    (set! foo "bar")

To start a new scope, with local variables, use `let*`:

    (let* (foo "bar"
           baz  "bart")
      (print "foo is %s" foo)
      (print "baz is %s" baz)
      ;...
    )

Inside a `let*` block, or a function, you'll only set the values of _local_
variables.  If you wish to affect something globally use the three-argument
form of `set!`:

    (let* (foo "bar"
           baz  "bart")
      (set! global "updated" true)
      ;..
    )

To define a function use `set!` with `fn*`:

    (set! fact (fn* (n)
      (if (<= n 1)
        1
          (* n (fact (- n 1))))))

Optionally you may write some help/usage information in your definition:

    (def! gcd (fn* (m n)
      "Return the greatest common divisor between the two arguments."
      (if (= (% m n) 0) n (gcd n (% m n)))))

Help information can be retrieved at runtime, for usage:

    (print (help print))

To define a macro use `defmacro!`:

    (defmacro! debug (fn* (x) `(print "Variable '%s' has value %s" '~x ~x)))

You might use this macro like so:

    (set! foo "steve")
    (debug foo)

That concludes the brief overview, note that `lambda` can be used as a synonym for `fn*`, and other synonyms exist.  In the interests of simplicity they're not covered here.


## See Also

* [README.md](README.md)
  * More details of the project.
* [PRIMITIVES.md](PRIMITIVES.md)
  * The list of built-in functions, whether implemented in Golang or YAL.

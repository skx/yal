* [Brief Yal Introduction](#brief-yal-introduction)
* [See Also](#see-also)




# Brief Yal Introduction

Yal is a typical toy lisp with support for numbers, strings, characters, hashes and structures.



## Primitive Types

Primitive types work as you would expect:

* Numbers can be written as integers in decimal, binary, or hex:
  * `(print 3)`
  * `(print 0xff)`
  * `(print 0b1010)`
* Floating point numbers are also supported:
  * `(print 3.4)`
* Strings are just encoded literally, and escaped characters are honored:
  * `(print "Hello, world\n")`
* Characters are written with a `#\` prefix:
  * `(print #\*)`
* Lists are written using parenthesis to group them:
  * `(print (list 1 2 3))`



## Other Types

We support hashes, which are key/value pairs, written between `{` and `}` pairs:

```lisp
(print { name "Steve" age (- 2022 1976) } )
```

Functions exist for getting/setting fields by name, and for iterating over the keys, values, or key/value pairs, contained in a given hash.

We also support structures, which are syntactical sugar for hashes, along with the autogeneration of some simple helper methods.

To define a "person" with three fields you'd write:

```lisp
(struct person name age address)
```

Once this `struct` has been defined it can be populated via the constructor:

```lisp
(person "Steve" "18" "123 Fake Street")
```

Values which are not supplied will default to `nil`:

```lisp
(person "Steve" "18")
```

The structure's fields can be accessed, and updated via generated methods, named after the type of the object, and the field involved:

```
; Set the variable "me" to be a new instance of the person-structure.
(set! me (person "Steve" "18" "123 Fake Street"))

; Change the adddress
(person.address me "999 Fake Lane")

; Confirm it worked
(print (person.address me))
```



## IF

`if` is a standard of lisp, and we support it.  We also implement `cond` - as a macro - and other functions for looping, ranging, and calling functions against lists.

The handling of `if` has been enhanced to allow **multiple** expressions to be executed in the "else" branch.  A standard `if` would look like this:

    (if true
      (print "This is executed")
     (print "This is not"))

We also allow:

    (if false
      (print "This is not executed")
     (print "This is executed")
     (print "This is executed too")
     (print "This is also executed")
     (print "This is executed as well ..")
     )

Here you see _multiple_ else-expressions are executed.  Only one expression is executed if the test succeeds, but the alternates have no limit - the return value of the expression is the return value of the last executed expression.



## Variables

To set the contents of a variable use `set!` which we saw above:

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



## Functions

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

As shown in the examples above parameters are named, rather than specifying
them as distinct symbols it is also possible to specify a default value by
expressing the parameters as a list (of two items only):

    (set! greet (fn* ( (name "World") )
      "Greet the supplied name, use the default if a name is not supplied."
      (print "Hello, %s" name)))

That would operate like so:

    (greet "Steve") ; "Hello, Steve"
    (greet)         ; "Hello, World"



## Macros
To define a macro use `defmacro!`:

    (defmacro! debug (fn* (x) `(print "Variable '%s' has value %s" '~x ~x)))

You might use this macro like so:

    (set! foo "steve")
    (debug foo)

That concludes the brief overview, note that `lambda` can be used as a synonym for `fn*`, and other synonyms exist.  In the interests of simplicity they're not covered here.



## See Also

* [LSP.md](LSP.md)
  * LSP support.
* [README.md](README.md)
  * More details of the project.
* [PRIMITIVES.md](PRIMITIVES.md)
  * The list of built-in functions, whether implemented in Golang or YAL.

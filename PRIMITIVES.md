# Primitives

We have a reasonable number of functions implemented, either in our golang core or in the standard-library which is loaded ahead of all user-scripts (the standard-library is implemented in lisp).

* Basic types include strings, numbers, errors, lists, hashes, etc.
  * `#t` is the true symbol, `#f` is false.
  * `true` and `false` are available as synonyms.
* Comparison functions:
  * `<`, `<=`, `>`, `>=`, `=`, & `eq`.
* Date and time functions:
  * `(date)` and `(time)` are implemented in our core application, but individual fields are made available via our standard-library:
    * `(year)`, `(month)`, `(day)`, `(weekday)`, `(hour)`, `(minute)`, `(second)` ‚ `(hms)`.
  * Each of these are demonstrated in [time.lisp](time.lisp).
* Error handling:
  * `error`, `try`, and `catch` - as demonstrated in [try.lisp](try.lisp).
* Hash operations:
  * Hashes are literals like this `{ :name "Steve" :location "Helsinki" }`
  * Hash functions are `contains?`, `get`, `keys`, `set`, & `vals`.
    * Note that keys are returned in sorted order.  Values are returned in order of their sorted keys too.
* Help functions:
  * `help` will return the supplied help text for any functions which provide it - which includes all of our built-in functions and large parts of our standard-library.
  * **NOTE**: This does not (yet?) include help for special forms such as `(let* ..)`, `(if ..)`, etc.
* List operations:
  * `butlast`, `car`, `cdr`, `cons`, `drop`, `every`, `first`, `last`, `list`, `repeated`, `sort`, & `take`.
* Logical operations:
  * `and`, & `or`.
* Mathematical operations:
  * `+`, `-`, `*`, `/`, `#`, & `%`.
* Macro support.
  * See [mtest.lisp](mtest.lisp) for some simple tests/usage examples.
  * The standard library uses macros to implement `(cond)`, `(while)`, and other functions.
* Platform features:
  * `arch`, `getenv`, `os`, `print`, `slurp`, `sprintf`, `str` & `type`
* String operations:
  * `chr`, `join`, `match` (regular-expression matching),`ord`, & `split`.
* Special forms:
  * `define`, `do`, `env`, `eval`, `gensym`, `if`, `lambda`, `let*`, `macroexpand`, `read`, `set!`, `quote`, & `quasiquote`.
* System functions:
  * Run commands via `(shell)`, find files via `(glob)`
  * Test files via `(directory?)`, `(exists?)` & `(file?)`.
* Tail call optimization.
* MAL compatibility:
  * `def!` can be used as a synonym for `define`.
  * `defmacro!` is used to define macros.
  * `fn*` can be used as a synonym for `lambda`.
  * `let*` can be used to define a local scope.
  * Several functions/macros that are expected to be present can be found in [stdlib/mal.lisp](stdlib/mal.lisp).

Building upon those primitives we have a larger standard-library of functions written in Lisp such as:

* `abs`, `apply`, `append`, `cond`, `count`, `false?`, `filter`, `lower`, `map`, `min`, `max`, `ms`, `nat`, `neg`, `now`, `nth`, `reduce`, `repeat`, `reverse`, `seq`, `strlen`, `true?`, `upper`, `when`, `while`, etc.

Although the lists above should be up to date you can check the definitions to see what is currently available:

* Primitives implemented in go:
  * [builtins/builtins.go](builtins/builtins.go)
* Primitives implemented in 100% pure lisp:
  * [stdlib/stdlib.lisp](stdlib/stdlib.lisp)
  * [stdlib/mal.lisp](stdlib/mal.lisp)
  * The code in these files is essentially **prepended** to any script that is supplied upon the command-line.



## Special Features

There are a couple of areas where we've implemented special/unusual things:

* Access to command-line arguments, when used via a shebang.
  * See [args.lisp](args.lisp) for an example.
* Introspection via the `(env)` function, which will return details of all variables/functions in the environment.
  * Allowing dynamic invocation shown in [dynamic.lisp](dynamic.lisp) and other neat things.
  * This includes help-information for both built-in and user-written functions.
* Support for hashes as well as lists/strings/numbers/etc.
  * A hash looks like this `{ :name "Steve" :location "Helsinki" }`
  * Sample code is visible in [hash.lisp](hash.lisp).
* Type checking for function parameters.
  * Via a `:type` suffix.  For example `(lambda (a:string b:number) ..`.

Here's an example of type-checking on a parameter value, in this case a list is required, via the `:list` suffix:

```lisp
(set! blah (fn* (a:list) (print "I received the list %s" a)))

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

* `(set! blah (fn* (a:list:number)  (print "I was given a list OR a number: %s" a)))`
  * Allows either a list, or a number.



## See Also

* [README.md](README.md)
  * More details of the project.
* [INTRODUCTION.md](INTRODUCTION.md)
  * Getting started setting variables, defining functions, etc.

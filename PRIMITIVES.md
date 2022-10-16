# Primitives

Here is a list of all the primitives which are available to yal users.

Note that you might need to consult the source of the standard-library, or
the help file, to see further details.  This is just intended as a summary.


## Symbols

The only notable special symbols are the strings which represent boolean values:

* `#t`
  * `true` is an alias.
* `#f`
  * `false` is an alias.


## Special Forms

Special forms are things that are built into the core interpreter, and include:

* `alias`
* `def!`
  * `define` is an alias.
* `defmacro!`
* `do`
* `env`
* `eval`
* `fn*`
  * `lambda` is an alias.
* `if`
* `let*`
* `macroexpand`
* `quote`
* `read`
* `set!`
* `try`


## Core Primitives

Core primitives are those that can be overridden, and are implemented in golang, in the [builtins/builtins.go](builtins/builtins.go) file.

Things you'll find here include:

* `#`
  * Exponent function.
* `%`
  * Modulus function.
* `*`
  * Multiplication function.
* `+`
  * Addition function.
* `-`
  * Subtraction function.
* `/`
  * Division function.
* `<`
  * Less-than function.
* `=`
  * Numerical comparison function.
* `arch`
  * Return the operating system architecture.
* `car`
  * Return the first item of a list.
* `cdr`
  * Return all items of the list, except the first.
* `chr`
  * Return the ASCII character of the given number.
* `cons`
  * Join two specified lists.
* `contains?`
  * Does the specified hash contain the given key?
* `date`
  * Return details of today's date, as a list.
* `directory?`
  * Does the given path represent something that exists, and is a directory?
* `directory:entries`
  * Return all entries beneath a given directory, recursively.
* `eq`
  * Equality test, handling arbitrary types.
* `error`
  * Return an error.
* `exists?`
  * Does the given path exist?
* `file?`
  * Does the given path exist, and is it not a directory?
* `file:lines`
  * Return the contents of the given file, as a list of strings.
* `file:read`
  * Return the contents of the given file, as a string.
* `file:stat`
  * Return details of the given path.
* `gensym`
  * Generate, and return, a unique symbol.  Useful for macro definitions.
* `get`
  * Get the given key from the specified hash.
* `getenv`
  * Read and return the given value from the environment.
* `glob`
  * Return the list of filenames matching the specified pattern.
* `help`
  * Return help for the specified function, either built-in or lisp.
* `join`
  * Convert every element of the supplied list into a string, and return the joined result.
* `keys`
  * Return the keys present in the specified hash.
  * Note that these are returned in sorted order.
* `list`
  * Create a new list.
* `match`
  * Perform a regular expression test.
* `ms`
  * Return the time, in milliseconds.
* `nil?`
  * Is the given value nil, or an empty list?
* now`
  * Return the number of seconds past the Unix Epoch.
* `ord`
  * Return the ASCII code of the specified character.
* `os`
  * Return a string describing the current operating-system.
* `print`
  * Output the specified string, or format string + values.
* `set`
  * Update the value of the specified hash-key.
* `shell`
  * Run a command via the shell, and return STDOUT and STDERR contents it generated.
* `sort`
  * Sort the given list.
* `split`
  * Split the given string, by the specified character.
* `sprintf`
  * Generate a string, using a format-string.
* `str`
  * Convert the specified parameter to a string.
* `time`
  * Return values relating to the current time, as a list.
* `type`
  * Return the type of the given object.
* `vals`
  * Return the values contained within the given hash.
  * Note that this returns things in the order of the sorted-keys.



## Standard Library

The standard library consists of routines, and helpers, which are written in 100% yal itself.

The implementation of these primitives can be found in the following two files:

* [stdlib/stdlib.lisp](stdlib/stdlib.lisp)
* [stdlib/mal.lisp](stdlib/mal.lisp)

The code in those files is essentially **prepended** to any script that is supplied upon the command-line.

Functions here include:

* !
* <=
* >
* >=
* abs
* and
* append
* apply
* apply-hash
* boolean?
* butlast
* concat
* date:day
* date:month
* date:weekday
* date:year
* dec
* directory:walk
* drop
* error?
* even?
* every
* file:stat:gid
* file:stat:mode
* file:stat:size
* file:stat:uid
* filter
* first
* function?
* hash?
* hms
* inc
* last
* length
* list?
* lower
* lower-table
* macro?
* map
* max
* min
* nat
* neg
* neg?
* nth
* number?
* odd?
* one?
* or
* pos?
* range
* reduce
* repeat
* repeated
* rest
* reverse
* seq
* (set!
* sign
* sqrt
* string?
* strlen
* symbol?
* take
* time:hms
* time:hour
* time:minute
* time:second
* translate
* upper
* upper-table
* zero?



# Old Content

* **TODO** Revisit this, and probably remove.

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

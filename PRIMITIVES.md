
* [Primitives](#primitives)
  * [Symbols](#symbols)
  * [Special Forms](#special-forms)
  * [Core Primitives](#core-primitives)
  * [Standard Library](#standard-library)
* [Type Checking](#type-checking)
* [See Also](#see-also)


# Primitives

Here is a list of all the primitives which are available to yal users.

Note that you might need to consult the source of the standard-library, or
the help file, to see further details.  This is just intended as a summary.


## Symbols

The only notable special symbols are the following strings which
represent the nil value. and ourboolean values.

* `nil`
  * The nil value.
* `#t`
  * `true` is also available as an alias.
* `#f`
  * `false` is also available as an alias.

In the future we _might_ support characters, via \#A, etc.



## Special Forms

Special forms are things that are built into the core interpreter, and include:

* `alias`
  * Define function aliases, this is used whenever we rename/change things in the standard-library to avoid breaking user scripts.
* `catch`.
  * Demonstrated in [try.lisp](try.lisp).
* `def!`
  * `define` is an alias.
* `defmacro!`
  * Demonstrated in [mtest.lisp](mtest.lisp).
* `do`
  * Execute each statement in the list.
* `env`
  * Env allows introspection of the current environment.
  * Demonstrated in [dynamic.lisp](dynamic.lisp)
* `eval`
  * Execute the given expression.
* `fn*`
  * `lambda` is an alias.
* `if`
  * Our conditional operation.
* `let*`
  * Create a new scope, with locally bound variables.
* `macroexpand`
  * Expand the given macro.
* `quote`
  * Return the argument without evaluating it.
* `read`
  * Read a form from the specified string.
* `set!`
  * Set the value of a variable.
* `try`
  * Error-catching warpper, demonstrated in [try.lisp](try.lisp).


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
  * Note that if only a single value is specified the reciprocal is returned - i.e. "(/ 3)" is equal to "1/3".
* `/=`
  * Numerical inequality test, if any argument is the same as another return false, otherwise if all arguments are unique return true.
* `<`
  * Less-than function.
* `=`
  * Numerical comparison function.
  * Note that multiple arguments are supported, not just two.
* `arch`
  * Return the operating system architecture.
* `car`
  * Return the first item of a list.
* `cdr`
  * Return all items of the list, except the first.
* `chr`
  * Return the ASCII character of the given number.
* `cons`
  * Add the element to the start of the given (potentialy empty) list.
* `contains?`
  * Does the specified hash contain the given key?
* `date`
  * Return details of today's date, as a list.
  * Demonstrated in [time.lisp](time.lisp).
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
* `file:write`
  * Write the specified content to the provided path.
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
* `now`
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
  * Run a command via the shell, and return STDOUT and STDERR it generated.
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
  * Demonstrated in [time.lisp](time.lisp).
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

* `!`
  * Logical "not".
* `<=`
  * Is the first number less than, or equal to, the second?
* `>`
  * Is the first number greater than the second?
* `>=`
  * Is the first number greater than, or equal to, the second?
* `abs`
  * Return the absolute value of the specified number.
* `and`
  * Logical operator, are all elements true?
* `append`
  * Append the given entry to the specified list.
* `apply`
  * Call the specified function on every element in the supplied list.
* `apply-hash`
  * Call the specified function against every key present in the specified hash.
* `apply-pairs`
  * Call the specified function on every two elements of the given list, as pairs.
* `boolean?`
  * Is the given thing a boolean?
* `butlast`
  * Return all elements of the supplied list, except for the last.
* `concat`
  * Join the specified lists.
* `date:day`
  * Return the current day of the month, via the output of `date`.
* `date:month`
  * Return the current month, via the output of `date`.
* `date:weekday`
  * Return the current day of the week, via the output of `date`.
* `date:year`
  * Return the current year, via the output of `date`.
* `dec`
  * Decrease the given thing by one.
* `directory:walk`
  * Invoke the specified callback, with every path-name contained beneath the specified directory - recursively.
* `drop`
  * Remove the specified number of elements from the provided list.
* `error?`
  * Is the given thing an error?
* `even?`
  * Is the given number even?
* `every`
  * Return true if applying the specified function to every element of the list returns a true result.
* `file:stat:gid`
  * Return the GID of the path, from the information provided by `(file:stat)`.
* `file:stat:mode`
  * Return the mode of the path, from the information provided by `(file:stat)`.
* `file:stat:size`
  * Return the size of the path, from the information provided by `(file:stat)`.
* `file:stat:uid`
  * Return the UID of the path, from the information provided by `(file:stat)`.
* `file:which`
  * Locate the specified binary's location, upon the users' PATH.
  * NOTE: This is almost certainly Unix/Linux/Darwin only, and will fail upon Windows systems.
* `file:write`
  * Write the specified content to the given path.
* `filter`
  * Remove every element from the given list, unless the function returns true.
* `first`
  * Return the first element of the given list.
  * This is the same as `car`.
* `function?`
  * Is the given thing a function?
* `hash?`
  * Is the given thing a hash?
* `inc`
  * Increment the given variable.
* `last`
  * Return the last element of the specified list.
* `length`
  * Return the length of the specified list.
* `list?`
  * Is the given thing a list?
* `lower`
  * Return an lower-case version of the specified string.
* `lower-table`
  * A translation table for converting an upper-case character to lower-case.
* `macro?`
  * Is the given thing a macro?
* `map`
  * Return the results of applying the specified function to every element of the given list.
* `max`
  * Return the maximum value in the specified list.
* `min`
  * Return the maximum value in the specified list.
* `nat`
  * Return the list of natural numbers 1 to N.
* `neg`
  * Negate the given number, and return it.
* `neg?`
  * Is the given number negative?
* `nth`
  * Return the Nth element of the list.
* `number?`
  * Is the given thing a number?
* `odd?`
  * Is the given number odd?
* `one?`
  * Is the given number equal to one?
* `or`
  * Logical operator, are any elements true?
* `pos?`
  * Is the given number positive?
* `range`
  * Return a list of numbers between the given start/end, using the specified step-size.
* `reduce`
  * Our reduce function, with the list, function and accumulator.
* `repeat`
  * Run the given body N times.
* `repeated`
  * Return a list of length N whose elements are all X.
* `rest`
  * Return the rest of the list, except the first element.
  * This is the same as `cdr`.
* `reverse`
  * Reverse the contents of the specified list.
* `seq`
  * Return a list of numbers from 0 to N.
* `sign`
  * Return the sign of the given number.  (1 for positive, -1 for negative).
* `sqrt`
  * Return the square-root of the supplied number.
* `string?`
  * Is the given thing a string?
* `strlen`
  * Return the length of the specified string.
* `symbol?`
  * Is the given thing a symbol?
* `take`
  * Take only the first N items from the specified list.
* `time:hms`
  * Return the time in HH:MM:SS format, as a string.
* `time:hour`
  * Return the current hour, as found from `(time)`.
* `time:minute`
  * Return the current minute, as found from `(time)`.
* `time:second`
  * Return the current second, as found from `(time)`.
* `translate`
  * Translate a string of characters, via a lookup table.
  * Used by `lower`, and `upper`.
* `upper`
  * Return an upper-case version of the specified string.
* `upper-table`
  * A translation table for converting a lower-case character to upper-case.
* `zero?`
  * Is the given number zero?


# Type Checking

Type checking is optional, but supported for function parameters via a `:type` suffix.  Here's an example of type-checking on a parameter value, in this case a list is required, via the `:list` suffix:

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


# See Also

* [README.md](README.md)
  * More details of the project.
* [INTRODUCTION.md](INTRODUCTION.md)
  * Getting started setting variables, defining functions, etc.

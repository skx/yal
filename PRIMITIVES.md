* [Primitives](#primitives)
  * [Symbols](#symbols)
  * [Special Forms](#special-forms)
  * [Core Primitives](#core-primitives)
  * [Structure Methods](#structure-methods)
  * [Standard Library](#standard-library)
* [Type Checking](#type-checking)
* [Testing](#testing)
* [See Also](#see-also)




# Primitives

Here is a list of all the primitives which are available to yal users.

Note that you might need to consult the source of the standard-library, or
the help file, to see further details.  This document is primarily intended
as a quick summary, and might lapse behind reality at times.



## Symbols

The only notable special symbols are the following strings which
represent the nil value. and our boolean values.

* `nil`
  * The nil value.
* `#t`
  * `true` is also available as an alias.
* `#f`
  * `false` is also available as an alias.

Characters are specified via the `#\X` syntax, for escaped characters you just need to add the escape:

* `#\a` -> "a"
* `#\b` -> "b"
* ..
* `#\X` -> "X"
* `#\\n` -> newline
* `#\\t` -> tab



## Special Forms

Special forms are things that are built into the core interpreter, and handled specially.

You can receive a full-list of special forms via `(specials)`, this list will include:

* `$`
  * Allow running a shell-command, including pipelines, and return the output as either a string or a list of strings.
* `alias`
  * Define function aliases, this is used whenever we rename/change things in the standard-library to avoid breaking user scripts.
* `catch`.
  * Demonstrated in [examples/try.lisp](examples/try.lisp).
* `def!`
  * `define` is an alias.
* `defmacro!`
  * Demonstrated in [examples/mtest.lisp](examples/mtest.lisp).
* `do`
  * Execute each statement in the list.
* `env`
  * Env allows introspection of the current environment.
  * Demonstrated in [examples/dynamic.lisp](examples/dynamic.lisp)
* `eval`
  * Execute the given expression.
* `exit`
  * Terminate the interpreter, optionally with a given numeric status-code.
* `fn*`
  * `lambda` is an alias.
* `forever`
  * Run the supplied list of statements forever, never terminating, without recursion.
* `if`
  * Our conditional operation.
  * Note that we support multiple "else" statements, if the condition is not true.
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
* `stdlib`
  * Return the names of functions/macros defined in the standard-library.
  * This is any function defined between a call to `stdlib-start` and `stdlib-end`.
* `stdlib-end` (internal)
  * Mark ourselves as no longer loading the standard library.
* `stdlib-start` (internal)
  * Mark ourselves as loading the standard library, such that new functions/macros will be appended to the `stdlib` list.
* `struct`
  * Define a structure.
* `symbol`
  * Create a new symbol from the given string.
* `try`
  * Error-catching warpper, demonstrated in [examples/try.lisp](examples/try.lisp).



## Core Primitives

Core primitives are those that can be overridden, and are implemented in golang, in the [builtins/builtins.go](builtins/builtins.go) file.

You can receive a full-list of these via `(builtins)`, this list will include:

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
* `acos`
  * Trig. function.
* `arch`
  * Return the operating system architecture.
* `asin`
  * Trig. function.
* `atan`
  * Trig. function.
* `base`
  * Convert the specified integer to a string, in the given base.
* `body`
  * Return the body of a lisp-function.
* `builtins`
  * Return the list of built-in functions, implemented in golang.
* `car`
  * Return the first item of a list.
* `cdr`
  * Return all items of the list, except the first.
* `char=`
  * Return true if the supplied values are characters, equal in value.
* `char<`
  * Return true if the first character is less than the second.
* `char<=`
  * Return true if the first character is less than, or equal to the second.
* `char>`
  * Return true if the first character is greater than the second.
* `char>=`
  * Return true if the first character is greater than, or equal to the second.
* `chr`
  * Return the ASCII character of the given number.
* `cons`
  * Add the element to the start of the given (potentialy empty) list.
* `contains?`
  * Does the specified hash contain the given key?
* `cos`
  * Trig. function.
* `cosh`
  * Trig. function.
* `date`
  * Return details of today's date, as a list.
  * Demonstrated in [examples/time.lisp](examples/time.lisp).
* `dec2bin`
  * Convert the specified integer to a binary string.
* `dec2hex`
  * Convert the specified integer to a hexadecimal string.
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
* `explode`
  * Convert the supplied string to a list of characters.
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
* `md5`
  * Return the MD5 digest of the given string.
* `ms`
  * Return the time, in milliseconds.
* `nil?`
  * Is the given value nil, or an empty list?
* `nth`
  * Return the nth element of the supplied list.
* `now`
  * Return the number of seconds past the Unix Epoch.
* `number`
  * Convert the specified string to a number.  We accept base 2, 10, and 16.
  * Use the appropriate prefix in your input, for example "0b10101", or "0xFF".
* `ord`
  * Return the ASCII code of the specified character, or the first character of the supplied string.
* `os`
  * Return a string describing the current operating-system.
* `pad:left`
  * Pad the specified string to the given length, by prepending to it.
* `pad:right`
  * Pad the specified string to the given length, by appending to it.
* `print`
  * Output the specified string, or format string + values.
* `set`
  * Update the value of the specified hash-key.
* `sha1`
  * Return the SHA1 digest of the given string.
* `sha256`
  * Return the SHA256 digest of the given string.
* `shell`
  * Run a command via the shell, and return STDOUT and STDERR it generated.
* `sin`
  * Trig. function.
* `sinh`
  * Trig. function.
* `sort`
  * Sort the given list.
* `source`
  * Return the source of a lisp-function.
* `specials`
  * Return the list of built-in special-forms, implemented in golang.
* `split`
  * Split the given string, by the specified character.
* `sprintf`
  * Generate a string, using a format-string.
* `stack:empty?`
  * Is the stack empty?  If so return true, else false.
* `stack:push`
  * Add an item to the given stack.
* `stack:pop`
  * Add an item to the given stack.
* `stack:size`
  * Return the size of the given stack.
* `str`
  * Convert the specified parameter to a string.
* `string<`
  * Return true if the first string is less than the second.
* `string<=`
  * Return true if the first string is less than, or equal to the second.
* `string>`
  * Return true if the first string is greater than the second.
* `string>=`
  * Return true if the first string is greater than, or equal to the second.
* `tan`
  * Trig. function.
* `tanh`
  * Trig. function.
* `time`
  * Return values relating to the current time, as a list.
  * Demonstrated in [examples/time.lisp](examples/time.lisp).
* `type`
  * Return the type of the given object.
* `vals`
  * Return the values contained within the given hash.
  * Note that this returns things in the order of the sorted-keys.



## Structure Methods

A structure is a minimal wrapper over a hash, but when a structure is
defined several methods are created.  Assuming a person-structure has
been defined like so:

```lisp
(struct person name age address)
```

There is now a new structure, named `person` with three fields `name`, `age`, and `address` which can be instantiated.

To help operate upon this structure several methods have also been created:

* `(person "name" "age" "address")`
  * Constructor method, which returns a new struct instance.
  * If the number of arguments is less than the number of object-fields they will be left unset (i.e. nil).
* `(person? obj)`
  * Returns true if the given object is an instance of the person struct.
* `(person.name obj [new-value])`
  * Accessor/Mutator for the name-field in the given struct instance.
* `(person.age obj [new-value])`
  * Accessor/Mutator for the age-field in the given struct instance.
* `(person.address obj [new-value])`
  * Accessor/Mutator for the address-field in the given struct instance.



## Standard Library

The standard library consists of routines, and helpers, which are written in 100% yal itself.

The implementation of these primitives can be found in the following directory:

* [stdlib/stdlib/](stdlib/stdlib/)

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
* `find`
  * Return the offset(s) at which the given item occurs in the list, if at all.
* `first`
  * Return the first element of the given list.
  * This is the same as `car`.
* `flatten`
  * Convert a list of nested lists to a single list, flattening it.
* `function?`
  * Is the given thing a function?
* `hash?`
  * Is the given thing a hash?
* `inc`
  * Increment the given variable.
* `intersection`
  * Return those elements in common in the specified pair of lists.
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
* `map-pairs`
    * Return the results of applying the specified function to every pair of elements in the given list.
* `max`
  * Return the maximum value in the specified list.
* `mean`
  * Return the average of the numbers in the specified list.
* `member`
  * Return true if the specified item is contained within the given list.
* `member?`
  * Return true if the specified item is contained within the given list.
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
* `occurrences`
  * Count how many times the given item appears in the specified list.
* `odd?`
  * Is the given number odd?
* `one?`
  * Is the given number equal to one?
* `or`
  * Logical operator, are any elements true?
* `pi`
  * Return the value of PI - calculated via `atan` as per [this reference](https://en.m.wikibooks.org/wiki/Trigonometry/Calculating_Pi).
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
* `replace`
  * Replace any item of a list matching BEFORE with AFTER.
* `require`
  * Allow loading a file, by name.
* `rest`
  * Return the rest of the list, except the first element.
  * This is the same as `cdr`.
* `reverse`
  * Reverse the contents of the specified list.
* `seq`
  * Return a list of numbers from 0 to N.
* `sign`
  * Return the sign of the given number.  (1 for positive, -1 for negative).
* `sort-by`
  * Sort the given list, using the supplied comparison method.
* `sqrt`
  * Return the square-root of the supplied number.
* `string?`
  * Is the given thing a string?
* `strlen`
  * Return the length of the specified string.
* `substr`
  * Return part of the specified string, identified by offset and length.
* `sum`
  * Return the sum of the numbers in the specified list.
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
* `union`
  * Return a list of all items in the specified two lists - without duplicates.
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




# Testing

There is a simple set of tests written in Lisp, using a macro to define them easily, which can be viewed:

* [examples/lisp-tests.lisp](examples/lisp-tests.lisp)

Adding new tests is easy enough that this file should be updated over time with new test-cases.




# See Also

* [LSP.md](LSP.md)
  * LSP support.
* [README.md](README.md)
  * More details of the project.
* [INTRODUCTION.md](INTRODUCTION.md)
  * Getting started setting variables, defining functions, etc.

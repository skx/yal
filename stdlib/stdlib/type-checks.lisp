;;; type-checks.lisp - Type-comparisions for given objects

;; There is a built in `type` function which returns the type of an object.
;;
;; Use this to define some simple methods to test argument-types
(set! boolean?  (fn* (x)
                     "Returns true if the argument specified is a boolean value."
                     (eq (type x) "boolean")))

(set! error?    (fn* (x)
                     "Returns true if the argument specified is an error-value."
                     (eq (type x) "error")))

(set! function? (fn* (x) "Returns true if the argument specified is a function, either a built-in function, or a user-written one."
                     (or
                      (list
                       (eq (type x) "procedure(lisp)")
                       (eq (type x) "procedure(golang)")))))

(set! hash?     (fn* (x)
                     "Returns true if the argument specified is a hash."
                     (eq (type x) "hash")))

(set! macro?    (fn* (x)
                     "Returns true if the argument specified is a macro."
                     (eq (type x) "macro")))

(set! list?     (fn* (x)
                     "Returns true if the argument specified is a list."
                     (eq (type x) "list")))

(set! number?   (fn* (x)
                     "Returns true if the argument specified is a number."
                     (eq (type x) "number")))

(set! string?   (fn* (x)
                     "Returns true if the argument specified is a string."
                     (eq (type x) "string")))

(set! symbol?   (fn* (x)
                     "Returns true if the argument specified is a symbol."
                     (eq (type x) "symbol")))

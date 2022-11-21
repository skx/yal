;;; try.lisp - Demonstrate our error-handling, with try/catch.

;;
;; This file demonstrates our try/catch behaviour, which allows
;; catching errors at runtime, and continuing execution.
;;


(try
 (print "OK")
 (catch e
   (print "We expected no error to be thrown, but we got one:%s" e)))


(try
 (print (/ 1 0))
 (catch e
   (print "Expected error caught, when attempting division by zero:%s" e)))


(try
 (try "foo")
 (catch e
   (print "Expected error caught, when calling '(try)' with bogus arguments:%s" e)))


(try
  (nth () 1)
  (catch e
   (print "Expected error caught, when accessing beyond the end of a list:%s" e)))

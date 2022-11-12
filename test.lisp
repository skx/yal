;;; test.lisp - Simple feature-tests/demonstrations of our system.

;;
;; This is a sample input file for our minimal lisp interpreter.
;;
;; We use it to demonstrate and test some basic features.
;;
;; NOTE: A lot of the things called here are defined in the standard
;; library, which is pre-pended to all loaded-scripts.


;; Instead of just (+ 1 2) we allow multiple args
(print "Our mathematical functions allow 2+ arguments, e.g: %s = %d"
  (quote (+ 1 2 3 4 5 6)) (+ 1 2 3 4 5 6))


;;
;; Use our "repeat" function, from the standard library, to run a block
;; N/10 times.  The number of the attempt is given as a parameter.
;;
(repeat 10 (lambda (n) (print "I'm in a loop %d" n)))

;;
;; Use our "while" function, from the standard library, to run a block
;; of code N/5 times.
;;
(let* (a 5)
  (while (> a 0)
    (do
     (print "(while) loop - iteration %d" a)
     (set! a (- a 1) true))))


;; Define a function, `fact`, to calculate factorials.
(set! fact (fn* (n)
  (if (<= n 1)
    1
      (* n (fact (- n 1))))))




;; Return the number of ms a function invokation took.
(set! benchmark (fn* (fn)
                     "Run the specified function, while recording the time
it took to execute.  Return that time, in ms."
                     (let* (start-ms (ms)
                                     _ (fn)
                                     end-ms (ms))
                       (- end-ms start-ms))))

;; Invoke the factorial function, using apply
;;
;; Calculate the factorial of "big numbers" mostly as a test of the
;; `now` function which times how long it took.
(apply (list 1 10 100 1000 10000 50000 100000)
       (lambda (x)
         (print "Calculating %d factorial took %dms"
           x
           (benchmark (lambda () (fact x))))))


; Split a string into a list, reverse it, and join it
(let* (input "Steve Kemp")
   (print "Starting string: %s" input)
   (print "Reversed string: %s" (join (reverse (split "Steve Kemp" "")))))


;; Define a variable "foo => 0"
;; but then change it, and show that result
(let* (foo 0)
  (print "foo is set to %d" foo)
  (set! foo 3)
  (print "foo is now set to %d" foo))

;;Now we're outside the scope of the `let` so `foo` is nil
(if foo
  (print "something weird happened!")
     (print "foo is unset now, outside the scope of the `let`"))


;; Define another function, and invoke it
(set! sum2 (fn* (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc)))))
(print "Sum of 1-100: %d" (sum2 100 0))

;; Now create a utility function to square a number
(set! sq (fn* (x) (* x x)))

;; For each item in the range 1-10, print it, and the associated square.
;; Awesome!  Much Wow!
(apply (nat 10)
      (lambda (x)
        (print "%d\tsquared is %d" x (sq x))))

;; Test our some of our earlier functions against a range of numbers
(apply (list -2 -1 0 1 2 3 4 5)
  (lambda (x)
    (do
      (if (neg? x)  (print "%d is negative" x))
      (if (zero? x) (print "%d is ZERO"     x))
      (if (even? x) (print "%d is EVEN"     x))
      (if (odd? x)  (print "%d is ODD"      x)))))

;; Test that we can get the correct type of each of our primitives
(apply (list 1 "steve" (list 1 2 3) true #t false #f nil boolean? print)
  (lambda (x)
    (print "'%s' has type '%s'" (str x) (type x))))


;;
;; Show even numbers via the filter-function.
;;
(print "Even numbers from 0-10: %s" (filter (nat 11) even?))

;;
;; And again with square numbers.
;;
(print "Squared numbers from 0-10: %s" (map (nat 11) sq))


;;
;; Setup a list of integers, and do a few things with it.
;;
(let* (vals '(32 92 109 903 31 3 -93 -31 -17 -3))
  (print "Working with the list: %s " vals)
  (print "\tBiggest item is %d"       (max vals))
  (print "\tSmallest item is %d"      (min vals))
  (print "\tReversed list is %s "     (reverse vals))
  (print "\tSorted list is %s "       (sort vals))
  (print "\tFirst item is %d "        (first vals))
  (print "\tRemaining items %s "      (rest vals)))


;;
;; A simple assertion function
;;
(set! assert (fn* (result msg)
  (if result ()
    (print "ASSERT failed - %s" msg))))

;;
;; Make some basic tests using our assert function.
;;
(assert (function? print)  "(function? print) failed")
(assert (function? assert) "(function? assert) failed")

(assert (eq 6  (+ 1 2 3))     "1+2+3 != 6")
(assert (eq 24 (* 2 3 4))     "2*3*4 != 24")
(assert (eq 70 (- 100 10 20)) "100-10-20 != 70")

(assert (eq (type type)   "procedure(golang)")  "(type type)")
(assert (eq (type assert) "procedure(lisp)")    "(type assert)")
(assert (eq (type 1)    "number")               "(type number)")
(assert (eq (type "me") "string")               "(type string)")
(assert (eq (type (list 1 2)) "list")           "(type list)")

(assert (neg? -3)            "negative number detected")
(assert (! (neg? 0) )        "zero is not negative")
(assert (! (neg? 30) )       "a positive number is not negative")
(assert (= (abs -3) (abs 3)) "abs(-3) == 3")

(assert (= (fact 1) 1) "1! = 1")
(assert (= (fact 2) 2) "2! = 2")
(assert (= (fact 3) 6) "3! = 6")

(assert (< 3 30)       "3 < 30")
(assert (! (< 30 30))  "30 < 30")
(assert (<= 30 30)     "30 < 30")
(assert (>  30 20)     "30 > 20")

;; nth starts counting at zero which is perhaps surprising.
(assert (= (nth (list 10 20 30) 0) 10) "Got the first item of the list.")
(assert (= (nth (list 10 20 30) 1) 20) "Got the second item of the list.")


;; We have a built-in eval function, which operates upon symbols, or strings.
(set! e "(+ 3 4)")
(print "Eval of '%s' resulted in %d" e (eval e))
(print "Eval of '%s' resulted in %d" "(+ 40 2)" (eval "(+ 40 2)"))

;; Simple test of `cond`
(set! a 6)
(cond
    (> a 20) (print "A > 20")
    (> a 15) (print "A > 15")
    true     (print "A is %d" a)
)

;;
;; Trivial Read/Eval pair
;;
(print "The answer to life, the universe, and everything is %d!\n"
  (eval (read "(* 6 7)")))

;; Upper-case and lower-casing of strings
(print "%s" (upper "hello, world"))
(print "%s" (lower "Hello, World; in LOWER-case."))

;; All done! -> In red :)
(set! red (fn* (msg) (sprintf "\e[0;31m%s\e[0m" msg)))
(print (red "All done!"))

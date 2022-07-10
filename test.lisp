;;
;; This is a sample input file for our minimal lisp interpreter.
;;
;; We use it to demonstrates and test our the basic features.
;;
;; NOTE: A lot of the things called here are defined in the standard
;; library, which is pre-pended to all loaded-scripts.


;; Instead of just (+ 1 2) we allow multiple args
(print "Our mathematical functions allow 2+ arguments, e.g: %s = %s"
  (quote (+ 1 2 3 4 5 6)) (+ 1 2 3 4 5 6))

;; Define a function, `fact`, to calculate factorials.
(define fact (lambda (n)
  (if (<= n 1)
    1
      (* n (fact (- n 1))))))

;; Invoke the factorial function, using apply
(apply (list 1 2 3 4 5 6 7 8 9 10)
  (lambda (x)
    (print "%s! => %s" x (fact x))))


                                        ; Split a string into a list, reverse it, and join it
(let ((input "Steve Kemp"))
  (begin
   (print "Starting string: %s" input)
   (print "Reversed string: %s" (join (reverse (split "Steve Kemp" ""))))))


;; Define a variable "foo => 0"
;; but then change it, and show that result
(let ((foo 0))
   (begin
      (print "foo is set to %s" foo)
      (set! foo 3)
      (print "foo is now set to %s" foo)))

;;Now we're outside the scope of the `let` so `foo` is nil
(if foo
  (print "something weird happened!")
     (print "foo is unset now, outside the scope of the `let`"))


;; Define another function, and invoke it
(define sum2 (lambda (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc)))))
(print "Sum of 1-100: %s" (sum2 100 0))

;; Now create a utility function to square a number
(define sq (lambda (x) (* x x)))

;; For each item in the range 1-10, print it, and the associated square.
;; Awesome!  Much Wow!
(apply (nat 11)
      (lambda (x)
        (print "%s\tsquared is %s" x (sq x))))

;; Test our some of our earlier functions against a range of numbers
(apply (list -2 -1 0 1 2 3 4 5)
  (lambda (x)
    (begin
      (if (neg? x)  (print "%s is negative" x))
      (if (zero? x) (print "%s is ZERO"     x))
      (if (even? x) (print "%s is EVEN"     x))
      (if (odd? x)  (print "%s is ODD"      x)))))

;; Test that we can get the correct type of each of our primitives
(apply (list 1 "steve" (list 1 2 3) #t #f nil boolean? print)
  (lambda (x)
    (print "'%s' has type '%s'" x (type x))))

;; Test the nth function
;;
;; nth starts counting at zero which is perhaps surprising.
(if (= (nth (list 10 20 30 40 50) 0) 10)
    (print "Got the first item of the list."))

(if (= (nth (list 10 20 30 40 50) 1) 20)
 (print "Got the second item of the list."))

;;
;; Show even numbers via the filter-function.
;;
(print "Even numbers from 0-10: %s"
       (filter (nat 11) (lambda (x) (even? x))))

;;
;; And again with square numbers.
;;
(print "Squared numbers from 0-10: %s"
       (map (nat 11) (lambda (x) (sq x))))


;;
;; Setup a list of integers, and do a few things with it.
;;
(let ((vals '(32 92 109 903 31 3 -93 -31 -17 -3)))
  (begin
     (print "Working with the list: %s " vals)
     (print "\tBiggest item is %s"       (max vals))
     (print "\tSmallest item is %s"      (min vals))
     (print "\tReversed list is %s "     (reverse vals))
     (print "\tSorted list is %s "       (sort vals))
     (print "\tFirst item is %s "        (first vals))
     (print "\tRemaining items %s "      (rest vals))
   ))

;
; all done
;
(print "All done")

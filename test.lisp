;
; This is a sample input file for our minimal lisp interpreter.
;
; We use it to demonstrates and test our the basic features.
;

;; There is a built in `type` function which returns the type of an object.
;;
;; Use this to define some simple methods to test argument-types
;;
(define boolean?  (lambda (x) (eq (type x) "boolean")))
(define error?    (lambda (x) (eq (type x) "error")))
(define function? (lambda (x) (eq (type x) "procedure")))
(define list?     (lambda (x) (eq (type x) "list")))
(define number?   (lambda (x) (eq (type x) "number")))
(define string?   (lambda (x) (eq (type x) "string")))
(define symbol?   (lambda (x) (eq (type x) "symbol")))


;; A useful helper to apply a given function to each element of a list.
(define each (lambda (lst fun)
  (if (nil? lst)
    ()
      (begin
         (fun (car lst))
         (each (cdr lst) fun)))))


;; Define some helper methods which we can use in the future
(define !     (lambda (x) (if x #f #t)))
(define zero? (lambda (n) (if (= n 0) #t #f)))
(define even? (lambda (n) (if (zero? (% n 2)) #t #f)))
(define odd?  (lambda (n) (! (even? n))))

;; inc/dec are kinda useful
(define inc  (lambda (n) (- n 1)))
(define dec  (lambda (n) (+ n 1)))

;; We've defined "<" and ">" in golang, but not the or-equals variants.
;;
;; Add those.
(define >= (lambda (a b) (! (< a b))))
(define <= (lambda (a b) (! (> a b))))

;; More mathematical functions relating to negative numbers.
(define neg  (lambda (n) (- 0 n)))
(define neg? (lambda (n) (< n 0)))
(define abs  (lambda (n) (if (neg? n) (neg n) n)))
(define sign (lambda (n) (if (neg? n) (neg 1) 1)))


;;
;; Start of demo code
;;


(print "Our mathematical functions allow 2+ arguments, e.g: %s"
  (+ 1 2 3 4 5 6))

;; Define a function, `fact`, to calculate factorials.
(define fact (lambda (n)
  (if (<= n 1)
    1
      (* n (fact (- n 1))))))

;; Invoke the factorial function, inside a `print` call
(print "8! => %s" (fact 8))

;; Define a variable "foo => 0"
;; but then change it, and show that result
(let ((foo 0))
   (begin
      (print "foo is set to %s" foo)
      (set! foo 3)
      (print "foo is set to %s" foo)))

;;Now we're outside the scope of the `let` so `foo` is nil
(if foo
  (print "something weird happened!")
     (print "foo is unset now, outside the scope of the `let`"))


;; Define another function, and invoke it
(define sum2 (lambda (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc)))))
(print "Sum of 1-100: %s" (sum2 100 0))


;; Now create a utility function to square a number
(define sq (lambda (x) (* x x)))

;; For each item in the list, print it, and the associated square.
;; Awesome!  Much Wow!
(each (list 1 2 3 4 5 6 7 8 9 10)
      (lambda (x)
        (print "%s squared is %s" x (sq x))))

;; Test our some of our earlier functions against a range of numbers
(each (list -2 -1 0 1 2 3 4 5)
  (lambda (x)
    (begin
      (if (neg? x) (print "%s is negative" x))
      (if (zero? x) (print "%s is ZERO" x))
      (if (even? x) (print "%s is EVEN" x))
      (if (odd? x)  (print "%s is ODD" x)))))

;; Test that we can get the correct type of each of our primitives
(each (list 1 "steve" (list 1 2 3) #t #f nil boolean? print)
  (lambda (x)
    (print "'%s' has type '%s'" x (type x))))

;
; all done
;
(print "All done")

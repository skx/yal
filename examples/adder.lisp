;;; adder.lisp - Demonstrate creating an adder with closures.


;; Generate a function that adds two numbers.
;;
;; The first number is a constant, which is set when
;; the function is called.
;;
(set! make-adder (fn* (n)
   (fn* (m) (+ n m))))

(set! doubler (fn* (f) (lambda (x) (f x x))))
(print ((doubler *) 4))

;;
;; Here's a function which uses a closure to keep
;; returning an incremented value each time it is called
;;
(set! counter (fn* (m)
                 (fn* ()
                      (do
                          (set! m (+ m 1) true)
                          m))))


;; Now create two adders.
;;
;; Here we see how the first value, N, is set.
;;
;; The second value in the definition above, M, is
;; set when the generated function is called.
(set! addFive (make-adder 5))
(set! addTen  (make-adder 10))


;;
;; Finally we can invoke our generated functions,
;; which work as you'd expect:
;;
(print "(+ 10 5) => %d" (addFive 10))
(print "(+ 10 3) => %d" (addTen   3))


;; We can create a counter using the counter function
;; we defined above too:
(set! count (counter 0))

;; Now call that ten times
(repeat 10 (lambda (n)
             (print "Counter shows: %d" (count))))

;; And a final run
(print "Last counter returned: %d" (count))
